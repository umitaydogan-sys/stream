package mp4

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/stream"
	"github.com/fluxstream/fluxstream/internal/transcode"
)

// Server serves fMP4 and WebM progressive streaming
type Server struct {
	manager   *stream.Manager
	transcode *transcode.Manager
	httpPort  int
	jobsMu    sync.Mutex
	jobs      map[string]*liveJob
}

// NewServer creates a new MP4/WebM streaming server
func NewServer(manager *stream.Manager, tcManager *transcode.Manager, httpPort int) *Server {
	return &Server{
		manager:   manager,
		transcode: tcManager,
		httpPort:  httpPort,
		jobs:      make(map[string]*liveJob),
	}
}

// HandleFMP4 serves fragmented MP4 streaming via chunked transfer
func (s *Server) HandleFMP4(w http.ResponseWriter, r *http.Request) {
	key := extractMediaKey(r.URL.Path, "/mp4/", ".mp4")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	if s.serveFFmpegOutput(w, r, key, "mp4") {
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Subscribe to stream
	subID := fmt.Sprintf("fmp4_%s_%d", r.RemoteAddr, time.Now().UnixNano())
	sub := s.manager.Subscribe(key, subID, 256)
	if sub == nil {
		http.Error(w, "Subscribe failed", http.StatusInternalServerError)
		return
	}
	defer s.manager.Unsubscribe(key, subID)

	log.Printf("[fMP4] İzleyici bağlandı: %s -> %s", r.RemoteAddr, key)

	// Write ftyp
	ftyp := buildBox("ftyp", []byte("iso6\x00\x00\x02\x00iso6mp41"))
	w.Write(ftyp)
	flusher.Flush()

	// Track state
	var videoSeqHeader []byte
	var audioSeqHeader []byte
	initSent := false
	seqNum := 0
	var samples []fmp4Sample

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}

			if pkt.IsSequenceHeader {
				if pkt.Type == media.PacketTypeVideo {
					videoSeqHeader = make([]byte, len(pkt.Data))
					copy(videoSeqHeader, pkt.Data)
				} else if pkt.Type == media.PacketTypeAudio {
					audioSeqHeader = make([]byte, len(pkt.Data))
					copy(audioSeqHeader, pkt.Data)
				}
				if !initSent && videoSeqHeader != nil {
					moov := buildInitMoov(videoSeqHeader, audioSeqHeader)
					w.Write(moov)
					flusher.Flush()
					initSent = true
				}
				continue
			}

			if !initSent {
				continue
			}

			// Accumulate samples
			data := pkt.Data
			if pkt.Type == media.PacketTypeVideo && len(data) > 5 {
				data = data[5:]
			} else if pkt.Type == media.PacketTypeAudio && len(data) > 2 {
				data = data[2:]
			}

			samples = append(samples, fmp4Sample{
				data:       data,
				timestamp:  pkt.Timestamp,
				isVideo:    pkt.Type == media.PacketTypeVideo,
				isKeyframe: pkt.IsKeyframe,
			})

			// Flush fragment on keyframe (GOP-based fragmentation)
			if pkt.Type == media.PacketTypeVideo && pkt.IsKeyframe && len(samples) > 1 {
				fragment := buildFragment(seqNum, samples[:len(samples)-1])
				w.Write(fragment)
				flusher.Flush()
				seqNum++
				samples = samples[len(samples)-1:]
			}

		case <-sub.Done:
			return
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) serveFFmpegOutput(w http.ResponseWriter, r *http.Request, streamKey, format string) bool {
	inputURL := s.waitForLiveManifestURL(streamKey, 4*time.Second)
	if inputURL == "" {
		return false
	}

	ffPath, err := s.transcode.DetectFFmpeg()
	if err != nil {
		log.Printf("[%s] FFmpeg bulunamadi, native fallback kullaniliyor: %v", strings.ToUpper(format), err)
		return false
	}

	job, err := s.getOrStartLiveJob(streamKey, format, ffPath, inputURL)
	if err != nil {
		log.Printf("[%s] shared live job baslatilamadi, native fallback kullaniliyor: %v", strings.ToUpper(format), err)
		return false
	}

	w.Header().Set("Content-Type", contentTypeForFormat(format))
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	var flusher func()
	if f, ok := w.(http.Flusher); ok {
		flusher = f.Flush
	}
	subID := fmt.Sprintf("%s_%s_%d", format, r.RemoteAddr, time.Now().UnixNano())
	job.serve(w, flusher, r.Context().Done(), subID)
	return true
}

func (s *Server) getOrStartLiveJob(streamKey, format, ffmpegPath, inputURL string) (*liveJob, error) {
	jobKey := format + ":" + streamKey

	s.jobsMu.Lock()
	if job, ok := s.jobs[jobKey]; ok && !job.closed {
		s.jobsMu.Unlock()
		return job, nil
	}
	s.jobsMu.Unlock()

	var job *liveJob
	job = newLiveJob(streamKey, format, ffmpegPath, inputURL, contentTypeForFormat(format), func(key string) {
		s.jobsMu.Lock()
		if existing, ok := s.jobs[key]; ok && existing == job {
			delete(s.jobs, key)
		}
		s.jobsMu.Unlock()
	})

	if err := job.start(); err != nil {
		return nil, err
	}

	s.jobsMu.Lock()
	if existing, ok := s.jobs[jobKey]; ok && !existing.closed {
		s.jobsMu.Unlock()
		job.stop()
		return existing, nil
	}
	s.jobs[jobKey] = job
	s.jobsMu.Unlock()

	log.Printf("[%s] shared live job baslatildi: %s", strings.ToUpper(format), streamKey)
	return job, nil
}

func (s *Server) liveManifestPath(streamKey string) string {
	if s.transcode == nil {
		return ""
	}
	return s.transcode.GetLiveManifestPath(streamKey)
}

func (s *Server) waitForLiveManifestPath(streamKey string, timeout time.Duration) string {
	if s.transcode == nil {
		return ""
	}
	return s.transcode.WaitForLiveManifestPath(streamKey, timeout)
}

func (s *Server) waitForLiveManifestURL(streamKey string, timeout time.Duration) string {
	if s.transcode == nil {
		return ""
	}
	return s.transcode.WaitForLiveManifestURL(streamKey, timeout)
}

func buildFFmpegStreamArgs(inputURL, format string) []string {
	base := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-user_agent", "FluxStreamInternal/2.0",
		"-avioflags", "direct",
		"-fflags", "nobuffer",
		"-analyzeduration", "0",
		"-probesize", "32k",
		"-i", inputURL,
		"-map", "0:v:0",
		"-map", "0:a:0?",
	}
	switch format {
	case "mp4":
		return append(base,
			"-c:v", "copy",
			"-c:a", "aac",
			"-ar", "48000",
			"-ac", "2",
			"-b:a", "128k",
			"-movflags", "frag_keyframe+empty_moov+default_base_moof+separate_moof+dash+omit_tfhd_offset",
			"-frag_duration", "1000000",
			"-flush_packets", "1",
			"-muxdelay", "0",
			"-muxpreload", "0",
			"-f", "mp4",
			"pipe:1",
		)
	case "webm":
		return append(base,
			"-c:v", "libvpx",
			"-deadline", "realtime",
			"-quality", "realtime",
			"-cpu-used", "6",
			"-threads", "4",
			"-row-mt", "1",
			"-g", "30",
			"-b:v", "2500k",
			"-c:a", "libopus",
			"-b:a", "128k",
			"-cluster_time_limit", "1000",
			"-f", "webm",
			"pipe:1",
		)
	default:
		return nil
	}
}

func contentTypeForFormat(format string) string {
	switch format {
	case "mp4":
		return "video/mp4"
	case "webm":
		return "video/webm"
	default:
		return "application/octet-stream"
	}
}

func extractMediaKey(path, prefix, ext string) string {
	key := strings.TrimPrefix(path, prefix)
	key = strings.Trim(key, "/")
	if key == "" {
		return ""
	}
	key = strings.Split(key, "/")[0]
	return strings.TrimSuffix(key, ext)
}

// HandleWebM serves WebM (VP8/VP9 + Opus) streaming
func (s *Server) HandleWebM(w http.ResponseWriter, r *http.Request) {
	key := extractMediaKey(r.URL.Path, "/webm/", ".webm")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	if s.serveFFmpegOutput(w, r, key, "webm") {
		return
	}

	w.Header().Set("Content-Type", "video/webm")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	subID := fmt.Sprintf("webm_%s_%d", r.RemoteAddr, time.Now().UnixNano())
	sub := s.manager.Subscribe(key, subID, 256)
	if sub == nil {
		http.Error(w, "Subscribe failed", http.StatusInternalServerError)
		return
	}
	defer s.manager.Unsubscribe(key, subID)

	log.Printf("[WebM] İzleyici bağlandı: %s -> %s", r.RemoteAddr, key)

	// Write WebM/EBML header
	webmHeader := buildWebMHeader()
	w.Write(webmHeader)
	flusher.Flush()

	clusterOpen := false
	var clusterStart uint32

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}

			if pkt.IsSequenceHeader {
				continue
			}

			data := pkt.Data
			if pkt.Type == media.PacketTypeVideo && len(data) > 5 {
				data = data[5:]
			} else if pkt.Type == media.PacketTypeAudio && len(data) > 2 {
				data = data[2:]
			}

			// Start new cluster on keyframe
			if pkt.Type == media.PacketTypeVideo && pkt.IsKeyframe {
				clusterOpen = true
				clusterStart = pkt.Timestamp

				cluster := buildWebMClusterHeader(pkt.Timestamp)
				w.Write(cluster)
			}

			if !clusterOpen {
				continue
			}

			// Write SimpleBlock
			relativeTS := int16(pkt.Timestamp - clusterStart)
			trackNum := byte(1) // video
			if pkt.Type == media.PacketTypeAudio {
				trackNum = 2
			}
			block := buildWebMSimpleBlock(trackNum, relativeTS, data, pkt.IsKeyframe)
			w.Write(block)
			flusher.Flush()

		case <-sub.Done:
			return
		case <-r.Context().Done():
			return
		}
	}
}

// ── fMP4 helpers ──

type fmp4Sample struct {
	data       []byte
	timestamp  uint32
	isVideo    bool
	isKeyframe bool
}

func buildBox(boxType string, data []byte) []byte {
	size := uint32(8 + len(data))
	buf := make([]byte, size)
	binary.BigEndian.PutUint32(buf[0:4], size)
	copy(buf[4:8], boxType)
	copy(buf[8:], data)
	return buf
}

func buildInitMoov(videoSeqHeader, audioSeqHeader []byte) []byte {
	mvhd := buildBox("mvhd", make([]byte, 100))

	trak := buildMinimalTrak(1, "video")
	moovContent := append(mvhd, trak...)

	if audioSeqHeader != nil {
		audioTrak := buildMinimalTrak(2, "audio")
		moovContent = append(moovContent, audioTrak...)
	}

	// mvex
	trexData := make([]byte, 24)
	binary.BigEndian.PutUint32(trexData[4:8], 1)
	binary.BigEndian.PutUint32(trexData[8:12], 1)
	trex := buildBox("trex", trexData)
	mvex := buildBox("mvex", trex)
	moovContent = append(moovContent, mvex...)

	return buildBox("moov", moovContent)
}

func buildMinimalTrak(trackID uint32, mediaType string) []byte {
	tkhdData := make([]byte, 84)
	binary.BigEndian.PutUint32(tkhdData[0:4], 0x00000003)
	binary.BigEndian.PutUint32(tkhdData[12:16], trackID)
	if mediaType == "video" {
		binary.BigEndian.PutUint32(tkhdData[76:80], 1920<<16)
		binary.BigEndian.PutUint32(tkhdData[80:84], 1080<<16)
	}
	tkhd := buildBox("tkhd", tkhdData)

	mdhdData := make([]byte, 24)
	binary.BigEndian.PutUint32(mdhdData[12:16], 90000)
	mdhd := buildBox("mdhd", mdhdData)

	var hdlrType string
	if mediaType == "video" {
		hdlrType = "vide"
	} else {
		hdlrType = "soun"
	}
	hdlrData := make([]byte, 25)
	copy(hdlrData[4:8], hdlrType)
	hdlr := buildBox("hdlr", hdlrData)

	stbl := buildBox("stbl",
		append(append(append(
			buildBox("stsd", make([]byte, 8)),
			buildBox("stts", make([]byte, 8))...),
			buildBox("stsc", make([]byte, 8))...),
			buildBox("stsz", make([]byte, 12))...))

	minf := buildBox("minf", stbl)
	mdia := buildBox("mdia", append(append(mdhd, hdlr...), minf...))
	return buildBox("trak", append(tkhd, mdia...))
}

func buildFragment(seqNum int, samples []fmp4Sample) []byte {
	if len(samples) == 0 {
		return nil
	}

	// Concatenate all sample data
	var mdatContent []byte
	for _, s := range samples {
		mdatContent = append(mdatContent, s.data...)
	}

	// mfhd
	mfhdData := make([]byte, 8)
	binary.BigEndian.PutUint32(mfhdData[4:8], uint32(seqNum))
	mfhd := buildBox("mfhd", mfhdData)

	// tfhd
	tfhdData := make([]byte, 12)
	binary.BigEndian.PutUint32(tfhdData[0:4], 0x020000)
	binary.BigEndian.PutUint32(tfhdData[4:8], 1)
	tfhd := buildBox("tfhd", tfhdData)

	// tfdt
	tfdtData := make([]byte, 12)
	binary.BigEndian.PutUint32(tfdtData[0:4], 0x01000000)
	binary.BigEndian.PutUint64(tfdtData[4:12], uint64(samples[0].timestamp))
	tfdt := buildBox("tfdt", tfdtData)

	// trun with per-sample sizes
	trunFlags := uint32(0x000201) // data-offset-present, sample-size-present
	trunData := make([]byte, 8+4*len(samples))
	binary.BigEndian.PutUint32(trunData[0:4], trunFlags)
	binary.BigEndian.PutUint32(trunData[4:8], uint32(len(samples)))
	for i, s := range samples {
		binary.BigEndian.PutUint32(trunData[8+i*4:12+i*4], uint32(len(s.data)))
	}
	trun := buildBox("trun", trunData)

	traf := buildBox("traf", append(append(tfhd, tfdt...), trun...))
	moof := buildBox("moof", append(mfhd, traf...))
	mdat := buildBox("mdat", mdatContent)

	return append(moof, mdat...)
}

// ── WebM/EBML helpers ──

func buildWebMHeader() []byte {
	// EBML Header
	ebml := []byte{
		0x1A, 0x45, 0xDF, 0xA3, // EBML element
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1F, // size
		0x42, 0x86, 0x81, 0x01, // EBMLVersion: 1
		0x42, 0xF7, 0x81, 0x01, // EBMLReadVersion: 1
		0x42, 0xF2, 0x81, 0x04, // EBMLMaxIDLength: 4
		0x42, 0xF3, 0x81, 0x08, // EBMLMaxSizeLength: 8
		0x42, 0x82, 0x84, 0x77, 0x65, 0x62, 0x6D, // DocType: "webm"
		0x42, 0x87, 0x81, 0x04, // DocTypeVersion: 4
		0x42, 0x85, 0x81, 0x02, // DocTypeReadVersion: 2
	}

	// Segment (unknown size for live streaming)
	segment := []byte{
		0x18, 0x53, 0x80, 0x67, // Segment element
		0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Unknown size
	}

	// Info
	info := []byte{
		0x15, 0x49, 0xA9, 0x66, // Info element
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, // size
		0x2A, 0xD7, 0xB1, 0x88, // TimestampScale
		0x00, 0x0F, 0x42, 0x40, 0x00, 0x00, 0x00, 0x00, // 1000000 ns = 1ms
		0x4D, 0x80, 0x87, // MuxingApp
		0x46, 0x6C, 0x75, 0x78, 0x53, 0x74, 0x72, // "FluxStr"
	}

	// Tracks
	tracks := buildWebMTracks()

	header := append(ebml, segment...)
	header = append(header, info...)
	header = append(header, tracks...)
	return header
}

func buildWebMTracks() []byte {
	// Video track entry
	videoTrack := []byte{
		0xAE, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1A, // TrackEntry
		0xD7, 0x81, 0x01, // TrackNumber: 1
		0x73, 0xC5, 0x81, 0x01, // TrackUID: 1
		0x83, 0x81, 0x01, // TrackType: video
		0x86, 0x86, // CodecID
		0x56, 0x5F, 0x56, 0x50, 0x39, 0x00, // "V_VP9\0"
		0xE0, 0x86, // Video element
		0xB0, 0x82, 0x07, 0x80, // PixelWidth: 1920
		0xBA, 0x82, 0x04, 0x38, // PixelHeight: 1080
	}

	// Audio track entry
	audioTrack := []byte{
		0xAE, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, // TrackEntry
		0xD7, 0x81, 0x02, // TrackNumber: 2
		0x73, 0xC5, 0x81, 0x02, // TrackUID: 2
		0x83, 0x81, 0x02, // TrackType: audio
		0x86, 0x87, // CodecID
		0x41, 0x5F, 0x4F, 0x50, 0x55, 0x53, 0x00, // "A_OPUS\0"
		0xE1, 0x84, // Audio element
		0xB5, 0x82, 0xBB, 0x80, // SamplingFrequency: 48000
	}

	content := append(videoTrack, audioTrack...)
	// Tracks header
	header := []byte{0x16, 0x54, 0xAE, 0x6B} // Tracks element ID
	sizeBytes := encodeEBMLSize(uint64(len(content)))
	header = append(header, sizeBytes...)
	return append(header, content...)
}

func buildWebMClusterHeader(timestamp uint32) []byte {
	cluster := []byte{
		0x1F, 0x43, 0xB6, 0x75, // Cluster element
		0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // Unknown size (live)
		0xE7, 0x84, // Timestamp element
	}
	ts := make([]byte, 4)
	binary.BigEndian.PutUint32(ts, timestamp)
	return append(cluster, ts...)
}

func buildWebMSimpleBlock(trackNum byte, relativeTS int16, data []byte, isKeyframe bool) []byte {
	// SimpleBlock = Element ID + size + tracknum + timestamp + flags + data
	blockData := make([]byte, 4+len(data))
	blockData[0] = 0x80 | trackNum // Track number (EBML coded)
	blockData[1] = byte(relativeTS >> 8)
	blockData[2] = byte(relativeTS)
	flags := byte(0)
	if isKeyframe {
		flags |= 0x80
	}
	blockData[3] = flags
	copy(blockData[4:], data)

	// SimpleBlock element (0xA3)
	header := []byte{0xA3}
	sizeBytes := encodeEBMLSize(uint64(len(blockData)))
	header = append(header, sizeBytes...)
	return append(header, blockData...)
}

func encodeEBMLSize(size uint64) []byte {
	if size < 0x7F {
		return []byte{byte(size) | 0x80}
	}
	if size < 0x3FFF {
		return []byte{byte(size>>8) | 0x40, byte(size)}
	}
	if size < 0x1FFFFF {
		return []byte{byte(size>>16) | 0x20, byte(size >> 8), byte(size)}
	}
	// 4-byte size
	return []byte{byte(size>>24) | 0x10, byte(size >> 16), byte(size >> 8), byte(size)}
}

// FMP4Stats tracks fMP4 viewer count
type FMP4Stats struct {
	mu      sync.RWMutex
	viewers int
}
