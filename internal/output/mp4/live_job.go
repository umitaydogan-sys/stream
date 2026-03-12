package mp4

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

const (
	liveJobIdleTimeout = 20 * time.Second
	maxLiveBoxSize     = 16 * 1024 * 1024
)

type liveChunk struct {
	seq    uint64
	data   []byte
	header bool
}

type liveSubscriber struct {
	id     string
	ch     chan liveChunk
	done   chan struct{}
	minSeq uint64
}

type liveChunkParser interface {
	Push(data []byte) []liveChunk
	Flush() []liveChunk
}

type liveJob struct {
	streamKey   string
	format      string
	ffmpegPath  string
	inputURL    string
	contentType string
	onStop      func(string)

	ctx    context.Context
	cancel context.CancelFunc
	cmd    *exec.Cmd
	logFile *os.File

	mu        sync.Mutex
	subs      map[string]*liveSubscriber
	header    [][]byte
	seq       uint64
	idleTimer *time.Timer
	closed    bool
}

func newLiveJob(streamKey, format, ffmpegPath, inputURL, contentType string, onStop func(string)) *liveJob {
	return &liveJob{
		streamKey:   streamKey,
		format:      format,
		ffmpegPath:  ffmpegPath,
		inputURL:    inputURL,
		contentType: contentType,
		onStop:      onStop,
		subs:        make(map[string]*liveSubscriber),
	}
}

func (j *liveJob) jobKey() string {
	return j.format + ":" + j.streamKey
}

func (j *liveJob) start() error {
	j.ctx, j.cancel = context.WithCancel(context.Background())
	args := buildFFmpegStreamArgs(j.inputURL, j.format)
	if len(args) == 0 {
		j.cancel()
		return io.EOF
	}

	cmd := exec.CommandContext(j.ctx, j.ffmpegPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		j.cancel()
		return err
	}
	logFile, err := os.CreateTemp("", "fluxstream-"+j.format+"-"+j.streamKey+"-*.log")
	if err == nil {
		j.logFile = logFile
		cmd.Stderr = logFile
	} else {
		cmd.Stderr = io.Discard
	}

	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		if j.logFile != nil {
			_ = j.logFile.Close()
			j.logFile = nil
		}
		j.cancel()
		return err
	}

	j.cmd = cmd
	go j.run(stdout)
	return nil
}

func (j *liveJob) serve(w io.Writer, flusher func(), done <-chan struct{}, subID string) {
	sub, header := j.addSubscriber(subID)
	defer j.removeSubscriber(subID)

	for _, chunk := range header {
		if _, err := w.Write(chunk); err != nil {
			return
		}
		if flusher != nil {
			flusher()
		}
	}

	for {
		select {
		case chunk, ok := <-sub.ch:
			if !ok {
				return
			}
			if !chunk.header && chunk.seq < sub.minSeq {
				continue
			}
			if _, err := w.Write(chunk.data); err != nil {
				return
			}
			if flusher != nil {
				flusher()
			}
		case <-sub.done:
			return
		case <-done:
			return
		}
	}
}

func (j *liveJob) addSubscriber(subID string) (*liveSubscriber, [][]byte) {
	j.mu.Lock()
	defer j.mu.Unlock()

	if j.idleTimer != nil {
		j.idleTimer.Stop()
		j.idleTimer = nil
	}

	minSeq := j.seq + 1
	if j.seq == 0 {
		minSeq = 1
	}

	sub := &liveSubscriber{
		id:     subID,
		ch:     make(chan liveChunk, 256),
		done:   make(chan struct{}),
		minSeq: minSeq,
	}
	j.subs[subID] = sub

	header := make([][]byte, 0, len(j.header))
	for _, chunk := range j.header {
		header = append(header, chunk)
	}
	return sub, header
}

func (j *liveJob) removeSubscriber(subID string) {
	j.mu.Lock()
	sub, ok := j.subs[subID]
	if ok {
		delete(j.subs, subID)
		close(sub.done)
		close(sub.ch)
	}
	if len(j.subs) == 0 && !j.closed && j.idleTimer == nil {
		j.idleTimer = time.AfterFunc(liveJobIdleTimeout, j.stop)
	}
	j.mu.Unlock()
}

func (j *liveJob) stop() {
	j.mu.Lock()
	if j.closed {
		j.mu.Unlock()
		return
	}
	j.closed = true
	cancel := j.cancel
	j.mu.Unlock()

	if cancel != nil {
		cancel()
	}
}

func (j *liveJob) run(stdout io.ReadCloser) {
	defer stdout.Close()

	parser := newLiveChunkParser(j.format)
	buf := make([]byte, 64*1024)

	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			for _, chunk := range parser.Push(buf[:n]) {
				j.broadcast(chunk)
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("[%s] live job okuma hatasi (%s): %v", j.format, j.streamKey, err)
			}
			break
		}
	}

	for _, chunk := range parser.Flush() {
		j.broadcast(chunk)
	}

	if j.cmd != nil {
		_ = j.cmd.Wait()
	}
	j.cleanup()
}

func (j *liveJob) cleanup() {
	j.mu.Lock()
	if j.idleTimer != nil {
		j.idleTimer.Stop()
		j.idleTimer = nil
	}
	subs := make([]*liveSubscriber, 0, len(j.subs))
	for _, sub := range j.subs {
		subs = append(subs, sub)
	}
	j.subs = make(map[string]*liveSubscriber)
	j.closed = true
	j.mu.Unlock()

	for _, sub := range subs {
		close(sub.done)
		close(sub.ch)
	}

	if j.cancel != nil {
		j.cancel()
	}
	if j.logFile != nil {
		_ = j.logFile.Close()
		j.logFile = nil
	}
	if j.onStop != nil {
		j.onStop(j.jobKey())
	}
}

func (j *liveJob) broadcast(chunk liveChunk) {
	if len(chunk.data) == 0 {
		return
	}

	j.mu.Lock()
	if j.closed {
		j.mu.Unlock()
		return
	}
	if chunk.header {
		j.header = append(j.header, chunk.data)
	} else if chunk.seq > j.seq {
		j.seq = chunk.seq
	}

	var slowSubs []*liveSubscriber
	for id, sub := range j.subs {
		if !chunk.header && chunk.seq < sub.minSeq {
			continue
		}
		select {
		case sub.ch <- chunk:
		default:
			log.Printf("[%s] yavas izleyici dusuruldu: %s -> %s", j.format, id, j.streamKey)
			slowSubs = append(slowSubs, sub)
			delete(j.subs, id)
		}
	}
	shouldStop := len(j.subs) == 0 && j.idleTimer == nil
	if shouldStop {
		j.idleTimer = time.AfterFunc(liveJobIdleTimeout, j.stop)
	}
	j.mu.Unlock()

	for _, sub := range slowSubs {
		close(sub.done)
		close(sub.ch)
	}
}

func newLiveChunkParser(format string) liveChunkParser {
	switch format {
	case "webm":
		return &webmLiveParser{}
	default:
		return &mp4LiveParser{}
	}
}

type mp4LiveParser struct {
	buffer []byte
	seq    uint64
}

func (p *mp4LiveParser) Push(data []byte) []liveChunk {
	p.buffer = append(p.buffer, data...)
	var chunks []liveChunk

	for {
		if len(p.buffer) < 8 {
			break
		}

		size := int(binary.BigEndian.Uint32(p.buffer[0:4]))
		headerSize := 8
		if size == 1 {
			if len(p.buffer) < 16 {
				break
			}
			size = int(binary.BigEndian.Uint64(p.buffer[8:16]))
			headerSize = 16
		}
		if size < headerSize || size > maxLiveBoxSize {
			p.buffer = p.buffer[1:]
			continue
		}
		if len(p.buffer) < size {
			break
		}

		box := append([]byte(nil), p.buffer[:size]...)
		p.buffer = p.buffer[size:]
		boxType := string(box[4:8])

		switch boxType {
		case "ftyp", "moov":
			chunks = append(chunks, liveChunk{header: true, data: box})
		case "moof":
			p.seq++
			chunks = append(chunks, liveChunk{seq: p.seq, data: box})
		default:
			if p.seq == 0 {
				p.seq = 1
			}
			chunks = append(chunks, liveChunk{seq: p.seq, data: box})
		}
	}

	return chunks
}

func (p *mp4LiveParser) Flush() []liveChunk {
	return nil
}

var webmClusterID = []byte{0x1F, 0x43, 0xB6, 0x75}

type webmLiveParser struct {
	buffer     []byte
	headerSent bool
	seq        uint64
}

func (p *webmLiveParser) Push(data []byte) []liveChunk {
	p.buffer = append(p.buffer, data...)
	var chunks []liveChunk

	if !p.headerSent {
		idx := bytes.Index(p.buffer, webmClusterID)
		if idx == -1 {
			return nil
		}
		if idx > 0 {
			chunks = append(chunks, liveChunk{
				header: true,
				data:   append([]byte(nil), p.buffer[:idx]...),
			})
		}
		p.buffer = p.buffer[idx:]
		p.headerSent = true
	}

	for {
		if len(p.buffer) < len(webmClusterID)+1 {
			break
		}
		next := bytes.Index(p.buffer[len(webmClusterID):], webmClusterID)
		if next == -1 {
			break
		}
		next += len(webmClusterID)
		cluster := append([]byte(nil), p.buffer[:next]...)
		p.buffer = p.buffer[next:]
		p.seq++
		chunks = append(chunks, liveChunk{seq: p.seq, data: cluster})
	}

	return chunks
}

func (p *webmLiveParser) Flush() []liveChunk {
	if !p.headerSent {
		if len(p.buffer) == 0 {
			return nil
		}
		return []liveChunk{{
			header: true,
			data:   append([]byte(nil), p.buffer...),
		}}
	}
	if len(p.buffer) == 0 {
		return nil
	}
	p.seq++
	return []liveChunk{{
		seq:  p.seq,
		data: append([]byte(nil), p.buffer...),
	}}
}
