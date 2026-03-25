package audio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/stream"
	"github.com/fluxstream/fluxstream/internal/transcode"
)

// Server serves all audio output formats
type Server struct {
	manager   *stream.Manager
	transcode *transcode.Manager
	httpPort  int
}

// NewServer creates a new audio output server
func NewServer(manager *stream.Manager, tcManager *transcode.Manager, httpPort int) *Server {
	return &Server{
		manager:   manager,
		transcode: tcManager,
		httpPort:  httpPort,
	}
}

// HandleMP3 serves MP3 audio stream (Icecast compatible)
func (s *Server) HandleMP3(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/audio/mp3/")
	s.handleMP3Like(w, r, key)
}

// HandleIcecast serves an Icecast-compatible MP3 endpoint on the main HTTP server.
func (s *Server) HandleIcecast(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/icecast/")
	s.handleMP3Like(w, r, key)
}

func (s *Server) handleMP3Like(w http.ResponseWriter, r *http.Request, key string) {
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	if s.serveFFmpegAudioOutput(w, r, key, "mp3", true) {
		return
	}

	// Icecast-compatible headers
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("icy-name", "FluxStream - "+key)
	w.Header().Set("icy-genre", "Live")
	w.Header().Set("icy-br", "128")

	s.streamAudio(w, r, key, "mp3", func(pkt *media.Packet) []byte {
		if pkt.Type != media.PacketTypeAudio {
			return nil
		}
		// Strip FLV AAC header, return raw audio
		if len(pkt.Data) > 2 {
			return pkt.Data[2:]
		}
		return nil
	})
}

// HandleAAC serves raw AAC audio stream
func (s *Server) HandleAAC(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/audio/aac/")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	if s.serveFFmpegAudioOutput(w, r, key, "aac", false) {
		return
	}

	w.Header().Set("Content-Type", "audio/aac")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.streamAudio(w, r, key, "aac", func(pkt *media.Packet) []byte {
		if pkt.Type != media.PacketTypeAudio {
			return nil
		}
		if pkt.IsSequenceHeader {
			return nil
		}
		// Strip FLV header and add ADTS header
		if len(pkt.Data) > 2 {
			raw := pkt.Data[2:]
			return addADTSHeader(raw)
		}
		return nil
	})
}

// HandleOpus serves Opus audio via OGG container
func (s *Server) HandleOpus(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/audio/opus/")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "audio/ogg; codecs=opus")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Write OGG stream header
	oggHeader := buildOGGOpusHeader()
	w.Write(oggHeader)
	flusher.Flush()

	subID := fmt.Sprintf("opus_%s_%d", r.RemoteAddr, time.Now().UnixNano())
	sub := s.manager.Subscribe(key, subID, 256)
	if sub == nil {
		http.Error(w, "Subscribe failed", http.StatusInternalServerError)
		return
	}
	defer s.manager.Unsubscribe(key, subID)

	granulePos := uint64(0)
	pageSeqNum := uint32(2) // after header pages

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}
			if pkt.Type != media.PacketTypeAudio || pkt.IsSequenceHeader {
				continue
			}
			data := pkt.Data
			if len(data) > 2 {
				data = data[2:]
			}

			granulePos += 960 // 20ms at 48kHz
			oggPage := buildOGGPage(data, granulePos, 0, pageSeqNum)
			pageSeqNum++
			w.Write(oggPage)
			flusher.Flush()

		case <-sub.Done:
			return
		case <-r.Context().Done():
			return
		}
	}
}

// HandleOGG serves OGG Vorbis audio
func (s *Server) HandleOGG(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/audio/ogg/")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	if s.serveFFmpegAudioOutput(w, r, key, "ogg", false) {
		return
	}

	w.Header().Set("Content-Type", "audio/ogg")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.streamAudio(w, r, key, "ogg", func(pkt *media.Packet) []byte {
		if pkt.Type != media.PacketTypeAudio || pkt.IsSequenceHeader {
			return nil
		}
		if len(pkt.Data) > 2 {
			return pkt.Data[2:]
		}
		return nil
	})
}

// HandleWAV serves uncompressed WAV audio (PCM)
func (s *Server) HandleWAV(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/audio/wav/")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	if s.serveFFmpegAudioOutput(w, r, key, "wav", false) {
		return
	}

	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Write WAV header with unknown size for streaming
	wavHeader := buildWAVHeader(44100, 2, 16)
	w.Write(wavHeader)
	flusher.Flush()

	s.streamAudio(w, r, key, "wav", func(pkt *media.Packet) []byte {
		if pkt.Type != media.PacketTypeAudio || pkt.IsSequenceHeader {
			return nil
		}
		if len(pkt.Data) > 2 {
			return pkt.Data[2:]
		}
		return nil
	})
}

// HandleFLAC serves FLAC audio stream
func (s *Server) HandleFLAC(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/audio/flac/")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	if s.serveFFmpegAudioOutput(w, r, key, "flac", false) {
		return
	}

	w.Header().Set("Content-Type", "audio/flac")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.streamAudio(w, r, key, "flac", func(pkt *media.Packet) []byte {
		if pkt.Type != media.PacketTypeAudio || pkt.IsSequenceHeader {
			return nil
		}
		if len(pkt.Data) > 2 {
			return pkt.Data[2:]
		}
		return nil
	})
}

// HandleHLSAudio serves audio-only HLS playlist
func (s *Server) HandleHLSAudio(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/audio/hls/")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	// Redirect to HLS with audio-only variant
	targetPlaylist := "audio.m3u8"
	if s.transcode != nil {
		var requestedTrackID uint8
		if raw := strings.TrimSpace(r.URL.Query().Get("track")); raw != "" {
			var parsed int
			fmt.Sscanf(raw, "%d", &parsed)
			if parsed > 0 && parsed <= 255 {
				requestedTrackID = uint8(parsed)
			}
		}
		if resolved := s.transcode.ResolveLiveAudioPlaylistPath(key, requestedTrackID); strings.TrimSpace(resolved) != "" {
			targetPlaylist = strings.TrimLeft(resolved, "/")
		}
	}
	target := fmt.Sprintf("/hls/%s/%s", key, targetPlaylist)
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}

// HandleDASHAudio serves audio-only DASH manifest
func (s *Server) HandleDASHAudio(w http.ResponseWriter, r *http.Request) {
	key := extractKey(r.URL.Path, "/audio/dash/")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	// Redirect to DASH with audio-only adaptation set
	target := fmt.Sprintf("/dash/%s/audio.mpd", key)
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusTemporaryRedirect)
}

// ── common streaming helper ──

func (s *Server) streamAudio(w http.ResponseWriter, r *http.Request, key, format string, transform func(*media.Packet) []byte) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	subID := fmt.Sprintf("audio_%s_%s_%d", format, r.RemoteAddr, time.Now().UnixNano())
	sub := s.manager.Subscribe(key, subID, 256)
	if sub == nil {
		http.Error(w, "Subscribe failed", http.StatusInternalServerError)
		return
	}
	defer s.manager.Unsubscribe(key, subID)

	log.Printf("[Audio-%s] İzleyici bağlandı: %s -> %s", strings.ToUpper(format), r.RemoteAddr, key)

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}
			data := transform(pkt)
			if data == nil {
				continue
			}
			if _, err := w.Write(data); err != nil {
				return
			}
			flusher.Flush()

		case <-sub.Done:
			return
		case <-r.Context().Done():
			return
		}
	}
}

// ── ADTS Header ──

func addADTSHeader(raw []byte) []byte {
	// ADTS header for AAC-LC, 44100Hz, stereo
	frameLen := len(raw) + 7
	header := make([]byte, 7)
	header[0] = 0xFF
	header[1] = 0xF1 // MPEG4, Layer 0, no CRC
	header[2] = 0x50 // AAC-LC, 44100Hz, private=0
	header[3] = byte(0x80 | ((frameLen >> 11) & 0x03))
	header[4] = byte((frameLen >> 3) & 0xFF)
	header[5] = byte(((frameLen & 0x07) << 5) | 0x1F)
	header[6] = 0xFC

	return append(header, raw...)
}

// ── WAV Header ──

func buildWAVHeader(sampleRate, channels, bitsPerSample int) []byte {
	header := make([]byte, 44)
	copy(header[0:4], "RIFF")
	binary.LittleEndian.PutUint32(header[4:8], 0xFFFFFFFF) // unknown size for streaming
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16) // chunk size
	binary.LittleEndian.PutUint16(header[20:22], 1)  // PCM
	binary.LittleEndian.PutUint16(header[22:24], uint16(channels))
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))
	byteRate := sampleRate * channels * bitsPerSample / 8
	binary.LittleEndian.PutUint32(header[28:32], uint32(byteRate))
	blockAlign := channels * bitsPerSample / 8
	binary.LittleEndian.PutUint16(header[32:34], uint16(blockAlign))
	binary.LittleEndian.PutUint16(header[34:36], uint16(bitsPerSample))
	copy(header[36:40], "data")
	binary.LittleEndian.PutUint32(header[40:44], 0xFFFFFFFF) // unknown size
	return header
}

// ── OGG Page Builder ──

func buildOGGOpusHeader() []byte {
	// Page 1: OpusHead
	opusHead := make([]byte, 19)
	copy(opusHead[0:8], "OpusHead")
	opusHead[8] = 1                                       // Version
	opusHead[9] = 2                                       // Channels
	binary.LittleEndian.PutUint16(opusHead[10:12], 0)     // Pre-skip
	binary.LittleEndian.PutUint32(opusHead[12:16], 48000) // Sample rate
	binary.LittleEndian.PutUint16(opusHead[16:18], 0)     // Output gain
	opusHead[18] = 0                                      // Channel mapping

	page1 := buildOGGPage(opusHead, 0, 0x02, 0) // BOS

	// Page 2: OpusTags
	tags := make([]byte, 0, 60)
	vendor := "FluxStream"
	vendorLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(vendorLen, uint32(len(vendor)))
	tags = append(tags, []byte("OpusTags")...)
	tags = append(tags, vendorLen...)
	tags = append(tags, []byte(vendor)...)
	commentCount := make([]byte, 4)
	binary.LittleEndian.PutUint32(commentCount, 0)
	tags = append(tags, commentCount...)

	page2 := buildOGGPage(tags, 0, 0, 1)

	return append(page1, page2...)
}

func buildOGGPage(data []byte, granulePos uint64, headerType byte, pageSeqNum uint32) []byte {
	// OGG page header (27 bytes + segment table)
	segments := (len(data) + 254) / 255
	if segments == 0 {
		segments = 1
	}

	header := make([]byte, 27+segments)
	copy(header[0:4], "OggS")
	header[4] = 0 // Version
	header[5] = headerType

	binary.LittleEndian.PutUint64(header[6:14], granulePos)
	binary.LittleEndian.PutUint32(header[14:18], 0) // Serial number
	binary.LittleEndian.PutUint32(header[18:22], pageSeqNum)
	// CRC will be 0 (simplified)
	binary.LittleEndian.PutUint32(header[22:26], 0)
	header[26] = byte(segments)

	// Segment table
	remaining := len(data)
	for i := 0; i < segments; i++ {
		if remaining >= 255 {
			header[27+i] = 255
			remaining -= 255
		} else {
			header[27+i] = byte(remaining)
			remaining = 0
		}
	}

	return append(header, data...)
}

func extractKey(path, prefix string) string {
	key := strings.TrimPrefix(path, prefix)
	key = strings.Trim(key, "/")
	if key == "" {
		return ""
	}
	key = strings.Split(key, "/")[0]
	for _, ext := range []string{".mp3", ".aac", ".ogg", ".wav", ".flac"} {
		key = strings.TrimSuffix(key, ext)
	}
	return key
}

func (s *Server) serveFFmpegAudioOutput(w http.ResponseWriter, r *http.Request, key, format string, withICY bool) bool {
	if s.transcode == nil || s.httpPort <= 0 {
		return false
	}

	inputURL := s.waitForLiveManifestURL(key, 4*time.Second)
	if inputURL == "" {
		return false
	}

	ffPath, err := s.transcode.DetectFFmpeg()
	if err != nil {
		log.Printf("[Audio-%s] FFmpeg bulunamadi, native fallback kullaniliyor: %v", strings.ToUpper(format), err)
		return false
	}

	args := buildFFmpegAudioArgs(inputURL, format)
	if len(args) == 0 {
		return false
	}

	cmd := exec.CommandContext(r.Context(), ffPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("[Audio-%s] stdout pipe acilamadi: %v", strings.ToUpper(format), err)
		return false
	}
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		log.Printf("[Audio-%s] FFmpeg baslatilamadi: %v", strings.ToUpper(format), err)
		return false
	}
	defer func() {
		_ = stdout.Close()
		_ = cmd.Wait()
	}()

	w.Header().Set("Content-Type", audioContentTypeForFormat(format))
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", key+"."+audioExtensionForFormat(format)))
	if withICY {
		w.Header().Set("icy-name", "FluxStream - "+key)
		w.Header().Set("icy-genre", "Live")
		w.Header().Set("icy-br", "128")
	}
	w.WriteHeader(http.StatusOK)

	flusher, _ := w.(http.Flusher)
	if flusher != nil {
		flusher.Flush()
	}

	_, copyErr := io.Copy(flushWriter{w: w, flusher: flusher}, stdout)
	if copyErr != nil && r.Context().Err() == nil {
		log.Printf("[Audio-%s] yayin kopyalama hatasi (%s): %v", strings.ToUpper(format), key, copyErr)
	}
	return true
}

type flushWriter struct {
	w       io.Writer
	flusher http.Flusher
}

func (fw flushWriter) Write(p []byte) (int, error) {
	n, err := fw.w.Write(p)
	if err == nil && n > 0 && fw.flusher != nil {
		fw.flusher.Flush()
	}
	return n, err
}

func buildFFmpegAudioArgs(inputURL, format string) []string {
	base := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-user_agent", "FluxStreamInternal/2.0",
		"-avioflags", "direct",
		"-fflags", "nobuffer",
		"-analyzeduration", "0",
		"-probesize", "32k",
		"-i", inputURL,
		"-map", "0:a:0?",
		"-vn",
		"-flush_packets", "1",
	}

	switch format {
	case "mp3":
		return append(base,
			"-c:a", "libmp3lame",
			"-b:a", "128k",
			"-ar", "44100",
			"-ac", "2",
			"-f", "mp3",
			"pipe:1",
		)
	case "aac":
		return append(base,
			"-c:a", "aac",
			"-b:a", "128k",
			"-ar", "44100",
			"-ac", "2",
			"-f", "adts",
			"pipe:1",
		)
	case "ogg":
		return append(base,
			"-c:a", "libvorbis",
			"-b:a", "160k",
			"-ar", "44100",
			"-ac", "2",
			"-f", "ogg",
			"pipe:1",
		)
	case "wav":
		return append(base,
			"-c:a", "pcm_s16le",
			"-ar", "44100",
			"-ac", "2",
			"-f", "wav",
			"pipe:1",
		)
	case "flac":
		return append(base,
			"-c:a", "flac",
			"-compression_level", "5",
			"-ar", "44100",
			"-ac", "2",
			"-f", "flac",
			"pipe:1",
		)
	default:
		return nil
	}
}

func audioContentTypeForFormat(format string) string {
	switch format {
	case "mp3":
		return "audio/mpeg"
	case "aac":
		return "audio/aac"
	case "ogg":
		return "audio/ogg"
	case "wav":
		return "audio/wav"
	case "flac":
		return "audio/flac"
	default:
		return "application/octet-stream"
	}
}

func audioExtensionForFormat(format string) string {
	switch format {
	case "mp3", "aac", "ogg", "wav", "flac":
		return format
	default:
		return "bin"
	}
}

func (s *Server) liveManifestPath(key string) string {
	if s.transcode == nil {
		return ""
	}
	return s.transcode.GetLiveManifestPath(key)
}

func (s *Server) waitForLiveManifestURL(key string, timeout time.Duration) string {
	if s.transcode == nil || s.httpPort <= 0 {
		return ""
	}
	return s.transcode.WaitForLiveManifestURL(key, timeout)
}

// AudioStats provides audio listener counts
type AudioStats struct {
	mu        sync.RWMutex
	listeners map[string]int // format -> count
}

// NewAudioStats creates a new audio stats tracker
func NewAudioStats() *AudioStats {
	return &AudioStats{
		listeners: make(map[string]int),
	}
}
