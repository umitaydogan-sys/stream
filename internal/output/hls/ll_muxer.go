package hls

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	ts "github.com/fluxstream/fluxstream/internal/media/container/ts"
)

// LLMuxer manages Low-Latency HLS outputs for all streams
type LLMuxer struct {
	outputDir string
	streams   map[string]*LLStreamMuxer
	mu        sync.RWMutex
}

// LLStreamMuxer handles LL-HLS output for a single stream
type LLStreamMuxer struct {
	streamKey       string
	outputDir       string
	tsMuxer         *ts.Muxer
	partDuration    time.Duration // partial segment duration (~200ms)
	segmentDuration time.Duration // full segment (~2s)
	maxSegments     int

	currentPart     *os.File
	currentPartIdx  int
	currentSegIdx   int
	partStart       uint32
	partStartedAt   time.Time
	segStart        uint32
	segStartedAt    time.Time
	partBytes       int
	partStarted     bool
	segStarted      bool
	hasVideo        bool
	partIndependent bool
	parts           []partInfo
	segments        []llSegmentInfo
	videoSeqHeader  []byte
	audioSeqHeader  []byte
	msn             int // media sequence number
	mu              sync.Mutex
}

type partInfo struct {
	Index       int
	Duration    float64
	Filename    string
	Independent bool // starts with keyframe
	Size        int
}

type llSegmentInfo struct {
	Index    int
	Duration float64
	Parts    []partInfo
	Filename string
}

// NewLLMuxer creates a new LL-HLS muxer
func NewLLMuxer(outputDir string) *LLMuxer {
	return &LLMuxer{
		outputDir: outputDir,
		streams:   make(map[string]*LLStreamMuxer),
	}
}

// AddStream creates an LL-HLS output for a stream key
func (m *LLMuxer) AddStream(streamKey string) *LLStreamMuxer {
	m.mu.Lock()
	defer m.mu.Unlock()

	dir := filepath.Join(m.outputDir, streamKey)
	os.MkdirAll(dir, 0755)

	sm := &LLStreamMuxer{
		streamKey:       streamKey,
		outputDir:       dir,
		tsMuxer:         ts.NewMuxer(),
		partDuration:    200 * time.Millisecond,
		segmentDuration: 2 * time.Second,
		maxSegments:     10,
	}
	m.streams[streamKey] = sm
	return sm
}

// RemoveStream removes LL-HLS output
func (m *LLMuxer) RemoveStream(streamKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sm, ok := m.streams[streamKey]; ok {
		sm.Close()
		delete(m.streams, streamKey)
	}
}

// WritePacket writes a media packet to LL-HLS output
func (sm *LLStreamMuxer) WritePacket(pkt *media.Packet) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if pkt.IsSequenceHeader {
		if pkt.Type == media.PacketTypeVideo {
			sm.videoSeqHeader = make([]byte, len(pkt.Data))
			copy(sm.videoSeqHeader, pkt.Data)
			sm.hasVideo = true
		} else if pkt.Type == media.PacketTypeAudio {
			sm.audioSeqHeader = make([]byte, len(pkt.Data))
			copy(sm.audioSeqHeader, pkt.Data)
		}
		return nil
	}
	if pkt.Type == media.PacketTypeVideo {
		sm.hasVideo = true
	}

	// Check if we need new partial segment
	if sm.currentPart == nil || sm.shouldSplitPart(pkt) {
		isSegBoundary := sm.shouldSplitSegment(pkt)
		if err := sm.finalizePart(pkt.Timestamp, isSegBoundary); err != nil {
			return err
		}
		if err := sm.startNewPart(pkt.Timestamp, pkt.IsKeyframe); err != nil {
			return err
		}
	}

	// Mux packet to TS
	mediaPkt := pkt.Clone()
	if pkt.Type == media.PacketTypeVideo && len(pkt.Data) > 5 {
		mediaPkt.Data = pkt.Data[5:]
	} else if pkt.Type == media.PacketTypeAudio && len(pkt.Data) > 2 {
		mediaPkt.Data = pkt.Data[2:]
	}

	tsData := sm.tsMuxer.MuxPacket(mediaPkt)
	if tsData == nil {
		return nil
	}

	n, err := sm.currentPart.Write(tsData)
	if err != nil {
		return fmt.Errorf("write LL part: %w", err)
	}
	sm.partBytes += n
	return nil
}

func (sm *LLStreamMuxer) shouldSplitPart(pkt *media.Packet) bool {
	if !sm.partStarted {
		return false
	}
	elapsedTS := time.Duration(0)
	if pkt.Timestamp >= sm.partStart {
		elapsedTS = time.Duration(pkt.Timestamp-sm.partStart) * time.Millisecond
	}
	elapsedWall := time.Since(sm.partStartedAt)
	elapsed := elapsedTS
	if elapsedWall > elapsed {
		elapsed = elapsedWall
	}
	return elapsed >= sm.partDuration
}

func (sm *LLStreamMuxer) shouldSplitSegment(pkt *media.Packet) bool {
	if !sm.segStarted {
		return false
	}
	elapsedTS := time.Duration(0)
	if pkt.Timestamp >= sm.segStart {
		elapsedTS = time.Duration(pkt.Timestamp-sm.segStart) * time.Millisecond
	}
	elapsedWall := time.Since(sm.segStartedAt)
	elapsed := elapsedTS
	if elapsedWall > elapsed {
		elapsed = elapsedWall
	}
	if elapsed < sm.segmentDuration {
		return false
	}
	if sm.hasVideo {
		return pkt.Type == media.PacketTypeVideo && pkt.IsKeyframe
	}
	return true
}

func (sm *LLStreamMuxer) finalizePart(ts uint32, isSegEnd bool) error {
	if sm.currentPart != nil {
		sm.currentPart.Close()

		deltaMS := uint32(0)
		if ts >= sm.partStart {
			deltaMS = ts - sm.partStart
		}
		duration := float64(deltaMS) / 1000.0
		if duration <= 0 {
			duration = time.Since(sm.partStartedAt).Seconds()
		}
		if duration <= 0 {
			duration = sm.partDuration.Seconds()
		}

		pi := partInfo{
			Index:       sm.currentPartIdx,
			Duration:    duration,
			Filename:    fmt.Sprintf("part_%d_%d.ts", sm.currentSegIdx, sm.currentPartIdx),
			Independent: sm.partIndependent,
			Size:        sm.partBytes,
		}
		sm.parts = append(sm.parts, pi)
		sm.currentPartIdx++
		sm.partStarted = false
		sm.partIndependent = false
	}

	if isSegEnd && len(sm.parts) > 0 {
		segDeltaMS := uint32(0)
		if ts >= sm.segStart {
			segDeltaMS = ts - sm.segStart
		}
		segDuration := float64(segDeltaMS) / 1000.0
		if segDuration <= 0 {
			segDuration = time.Since(sm.segStartedAt).Seconds()
		}
		if segDuration <= 0 {
			segDuration = sm.segmentDuration.Seconds()
		}

		sm.segments = append(sm.segments, llSegmentInfo{
			Index:    sm.currentSegIdx,
			Duration: segDuration,
			Parts:    sm.parts,
			Filename: fmt.Sprintf("ll_seg_%d.ts", sm.currentSegIdx),
		})

		// Concatenate parts into full segment
		sm.concatenateSegment()

		sm.parts = nil
		sm.currentPartIdx = 0

		// Cleanup old segments
		for len(sm.segments) > sm.maxSegments {
			old := sm.segments[0]
			os.Remove(filepath.Join(sm.outputDir, old.Filename))
			for _, p := range old.Parts {
				os.Remove(filepath.Join(sm.outputDir, p.Filename))
			}
			sm.segments = sm.segments[1:]
		}

		sm.currentSegIdx++
		sm.segStart = ts
		sm.segStarted = true
		sm.segStartedAt = time.Now()
		sm.msn++
	}

	sm.currentPart = nil
	sm.writePlaylist()
	return nil
}

func (sm *LLStreamMuxer) startNewPart(ts uint32, isKeyframe bool) error {
	if !sm.segStarted {
		sm.segStart = ts
		sm.segStarted = true
		sm.segStartedAt = time.Now()
	}

	filename := fmt.Sprintf("part_%d_%d.ts", sm.currentSegIdx, sm.currentPartIdx)
	f, err := os.Create(filepath.Join(sm.outputDir, filename))
	if err != nil {
		return fmt.Errorf("create LL part: %w", err)
	}

	sm.currentPart = f
	sm.partStart = ts
	sm.partStartedAt = time.Now()
	sm.partBytes = 0
	sm.partStarted = true
	sm.partIndependent = isKeyframe

	// Write PAT/PMT at start of each part
	patpmt := sm.tsMuxer.GeneratePatPmt()
	n, _ := f.Write(patpmt)
	sm.partBytes += n

	return nil
}

func (sm *LLStreamMuxer) concatenateSegment() {
	seg := sm.segments[len(sm.segments)-1]
	fullPath := filepath.Join(sm.outputDir, seg.Filename)

	f, err := os.Create(fullPath)
	if err != nil {
		return
	}
	defer f.Close()

	for _, p := range seg.Parts {
		data, err := os.ReadFile(filepath.Join(sm.outputDir, p.Filename))
		if err == nil {
			f.Write(data)
		}
	}
}

// writePlaylist generates the LL-HLS M3U8 playlist
func (sm *LLStreamMuxer) writePlaylist() {
	playlistPath := filepath.Join(sm.outputDir, "ll.m3u8")

	targetDur := sm.segmentDuration.Seconds()
	partTarget := sm.partDuration.Seconds()

	content := "#EXTM3U\n"
	content += "#EXT-X-VERSION:6\n"
	content += fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", int(targetDur)+1)
	content += "#EXT-X-SERVER-CONTROL:CAN-BLOCK-RELOAD=YES,CAN-SKIP-UNTIL=" +
		fmt.Sprintf("%.1f", targetDur*6) + ",PART-HOLD-BACK=" +
		fmt.Sprintf("%.1f", partTarget*3) + "\n"
	content += fmt.Sprintf("#EXT-X-PART-INF:PART-TARGET=%.3f\n", partTarget)

	if len(sm.segments) > 0 {
		content += fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", sm.segments[0].Index)
	}

	// Write completed segments
	for _, seg := range sm.segments {
		for _, p := range seg.Parts {
			independent := ""
			if p.Independent {
				independent = ",INDEPENDENT=YES"
			}
			content += fmt.Sprintf("#EXT-X-PART:DURATION=%.3f,URI=\"%s\"%s\n", p.Duration, p.Filename, independent)
		}
		content += fmt.Sprintf("#EXTINF:%.3f,\n", seg.Duration)
		content += seg.Filename + "\n"
	}

	// Write in-progress parts (current segment)
	for _, p := range sm.parts {
		independent := ""
		if p.Independent {
			independent = ",INDEPENDENT=YES"
		}
		content += fmt.Sprintf("#EXT-X-PART:DURATION=%.3f,URI=\"%s\"%s\n", p.Duration, p.Filename, independent)
	}

	// Preload hint for next part
	nextPart := fmt.Sprintf("part_%d_%d.ts", sm.currentSegIdx, sm.currentPartIdx)
	content += fmt.Sprintf("#EXT-X-PRELOAD-HINT:TYPE=PART,URI=\"%s\"\n", nextPart)

	os.WriteFile(playlistPath, []byte(content), 0644)
}

// Close finalizes the LL-HLS output
func (sm *LLStreamMuxer) Close() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.currentPart != nil {
		finalTS := sm.partStart + uint32(sm.partDuration.Milliseconds())
		_ = sm.finalizePart(finalTS, true)
	}

	// Write final playlist with ENDLIST
	playlistPath := filepath.Join(sm.outputDir, "ll.m3u8")
	if data, err := os.ReadFile(playlistPath); err == nil {
		data = append(data, []byte("#EXT-X-ENDLIST\n")...)
		os.WriteFile(playlistPath, data, 0644)
	}
}
