package hls

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	ts "github.com/fluxstream/fluxstream/internal/media/container/ts"
)

// Muxer handles HLS output for all streams
type Muxer struct {
	outputDir string
	streams   map[string]*StreamMuxer
	mu        sync.RWMutex
}

// StreamMuxer handles HLS output for a single stream
type StreamMuxer struct {
	streamKey       string
	outputDir       string
	tsMuxer         *ts.Muxer
	segmentDuration time.Duration
	maxSegments     int

	currentSegment   *os.File
	currentSegIdx    int
	segmentStart     uint32
	segmentStartedAt time.Time
	segmentBytes     int
	segments         []SegmentInfo
	hasVideo         bool
	hasAudio         bool
	videoSeqHeader   []byte
	audioSeqHeader   []byte
	videoConfigNALU  []byte
	aacProfile       int
	aacFreqIndex     int
	aacChannelCfg    int
	mu               sync.Mutex
}

// SegmentInfo describes an HLS segment
type SegmentInfo struct {
	Index    int
	Duration float64
	Filename string
	Size     int
}

// NewMuxer creates a new HLS muxer
func NewMuxer(outputDir string) *Muxer {
	return &Muxer{
		outputDir: outputDir,
		streams:   make(map[string]*StreamMuxer),
	}
}

// AddStream creates an HLS output for a stream key
func (m *Muxer) AddStream(streamKey string) *StreamMuxer {
	m.mu.Lock()
	defer m.mu.Unlock()

	dir := filepath.Join(m.outputDir, streamKey)
	os.MkdirAll(dir, 0755)

	sm := &StreamMuxer{
		streamKey:       streamKey,
		outputDir:       dir,
		tsMuxer:         ts.NewMuxer(),
		segmentDuration: 2 * time.Second,
		maxSegments:     10,
		currentSegIdx:   0,
		aacProfile:      1, // AAC-LC
		aacFreqIndex:    4, // 44100Hz
		aacChannelCfg:   2, // stereo
	}

	m.streams[streamKey] = sm
	return sm
}

// RemoveStream removes and cleans up HLS output for a stream
func (m *Muxer) RemoveStream(streamKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sm, ok := m.streams[streamKey]; ok {
		sm.Close()
		delete(m.streams, streamKey)
	}
}

// GetStream returns the stream muxer for a key
func (m *Muxer) GetStream(streamKey string) *StreamMuxer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.streams[streamKey]
}

// GetOutputDir returns the output directory
func (m *Muxer) GetOutputDir() string {
	return m.outputDir
}

// WritePacket writes a media packet to the HLS output
func (sm *StreamMuxer) WritePacket(pkt *media.Packet) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Store sequence headers (codec config)
	if pkt.IsSequenceHeader {
		if pkt.Type == media.PacketTypeVideo {
			sm.videoSeqHeader = make([]byte, len(pkt.Data))
			copy(sm.videoSeqHeader, pkt.Data)
			sm.videoConfigNALU = parseAVCConfigToAnnexB(pkt.Data)
			sm.hasVideo = true
		} else if pkt.Type == media.PacketTypeAudio {
			sm.audioSeqHeader = make([]byte, len(pkt.Data))
			copy(sm.audioSeqHeader, pkt.Data)
			sm.parseAACAudioSpecificConfig(pkt.Data)
			sm.hasAudio = true
		}
		return nil
	}
	if pkt.Type == media.PacketTypeVideo {
		sm.hasVideo = true
	}

	// Check if we need to start a new segment
	if sm.currentSegment == nil || sm.shouldSplit(pkt) {
		if err := sm.startNewSegment(pkt.Timestamp); err != nil {
			return err
		}
	}

	// Convert FLV payload to MPEG-TS friendly payload.
	mediaPkt := pkt.Clone()
	if pkt.Type == media.PacketTypeVideo {
		// RTMP/FLV AVC payload: [1 byte frame+codec][1 byte avcPacketType][3 byte cts][NALUs(length-prefixed)]
		if len(pkt.Data) <= 5 {
			return nil
		}
		annexB := avccToAnnexB(pkt.Data[5:])
		if pkt.IsKeyframe && len(sm.videoConfigNALU) > 0 {
			mediaPkt.Data = append(append([]byte{}, sm.videoConfigNALU...), annexB...)
		} else {
			mediaPkt.Data = annexB
		}
	} else if pkt.Type == media.PacketTypeAudio {
		if len(pkt.Data) < 2 {
			return nil
		}
		codecID := (pkt.Data[0] >> 4) & 0x0F
		if codecID == byte(media.AudioCodecAAC) {
			if len(pkt.Data) <= 2 {
				return nil
			}
			mediaPkt.Data = addADTSHeader(pkt.Data[2:], sm.aacProfile, sm.aacFreqIndex, sm.aacChannelCfg)
		} else {
			mediaPkt.Data = pkt.Data[1:]
		}
	}

	tsData := sm.tsMuxer.MuxPacket(mediaPkt)
	if tsData == nil {
		return nil
	}

	n, err := sm.currentSegment.Write(tsData)
	if err != nil {
		return fmt.Errorf("write TS: %w", err)
	}
	sm.segmentBytes += n

	return nil
}

// parseAVCConfigToAnnexB extracts SPS/PPS from AVCDecoderConfigurationRecord and returns Annex-B bytes.
func parseAVCConfigToAnnexB(data []byte) []byte {
	if len(data) < 11 {
		return nil
	}
	cfg := data
	if len(data) > 5 {
		cfg = data[5:]
	}
	if len(cfg) < 7 {
		return nil
	}

	pos := 6
	numSPS := int(cfg[5] & 0x1F)
	out := make([]byte, 0, 256)

	for i := 0; i < numSPS; i++ {
		if pos+2 > len(cfg) {
			return out
		}
		l := int(binary.BigEndian.Uint16(cfg[pos : pos+2]))
		pos += 2
		if l <= 0 || pos+l > len(cfg) {
			return out
		}
		out = append(out, 0x00, 0x00, 0x00, 0x01)
		out = append(out, cfg[pos:pos+l]...)
		pos += l
	}

	if pos >= len(cfg) {
		return out
	}

	numPPS := int(cfg[pos])
	pos++
	for i := 0; i < numPPS; i++ {
		if pos+2 > len(cfg) {
			return out
		}
		l := int(binary.BigEndian.Uint16(cfg[pos : pos+2]))
		pos += 2
		if l <= 0 || pos+l > len(cfg) {
			return out
		}
		out = append(out, 0x00, 0x00, 0x00, 0x01)
		out = append(out, cfg[pos:pos+l]...)
		pos += l
	}

	return out
}

// avccToAnnexB converts AVC NAL units from AVCC length-prefix format to Annex-B start-code format.
func avccToAnnexB(data []byte) []byte {
	if len(data) < 4 {
		return data
	}
	out := make([]byte, 0, len(data)+32)
	pos := 0

	for pos+4 <= len(data) {
		naluLen := int(binary.BigEndian.Uint32(data[pos : pos+4]))
		pos += 4
		if naluLen <= 0 || pos+naluLen > len(data) {
			// Fallback: if parsing fails, keep original payload.
			return data
		}
		out = append(out, 0x00, 0x00, 0x00, 0x01)
		out = append(out, data[pos:pos+naluLen]...)
		pos += naluLen
	}

	if len(out) == 0 {
		return data
	}
	return out
}

func (sm *StreamMuxer) parseAACAudioSpecificConfig(data []byte) {
	// FLV AAC sequence header: [sound byte][aacPacketType=0][AudioSpecificConfig...]
	if len(data) < 4 {
		return
	}
	asc := data[2:]
	if len(asc) < 2 {
		return
	}
	audioObjectType := int((asc[0] >> 3) & 0x1F)
	freqIdx := int(((asc[0] & 0x07) << 1) | ((asc[1] >> 7) & 0x01))
	chCfg := int((asc[1] >> 3) & 0x0F)

	if audioObjectType >= 2 {
		sm.aacProfile = audioObjectType - 1
	}
	if freqIdx >= 0 && freqIdx <= 12 {
		sm.aacFreqIndex = freqIdx
	}
	if chCfg > 0 && chCfg <= 7 {
		sm.aacChannelCfg = chCfg
	}
}

func addADTSHeader(raw []byte, profile, freqIdx, chCfg int) []byte {
	frameLen := len(raw) + 7
	adts := make([]byte, 7)
	adts[0] = 0xFF
	adts[1] = 0xF1
	adts[2] = byte((profile<<6)&0xC0 | (freqIdx<<2)&0x3C | (chCfg>>2)&0x01)
	adts[3] = byte((chCfg&0x03)<<6 | (frameLen>>11)&0x03)
	adts[4] = byte((frameLen >> 3) & 0xFF)
	adts[5] = byte((frameLen&0x07)<<5 | 0x1F)
	adts[6] = 0xFC
	return append(adts, raw...)
}

func (sm *StreamMuxer) shouldSplit(pkt *media.Packet) bool {
	if sm.currentSegment == nil {
		return false
	}

	elapsedTS := time.Duration(0)
	if pkt.Timestamp >= sm.segmentStart {
		elapsedTS = time.Duration(pkt.Timestamp-sm.segmentStart) * time.Millisecond
	}

	// Use timestamp-based elapsed time as the primary signal.
	// Only fall back to wall-clock when timestamps appear stuck (advancing
	// less than 250 ms while real time exceeds the segment duration).
	// This prevents micro-segments caused by timestamp jitter while still
	// providing a safety net for truly stale timestamp streams.
	elapsed := elapsedTS
	if elapsedTS < 250*time.Millisecond {
		elapsedWall := time.Since(sm.segmentStartedAt)
		if elapsedWall > sm.segmentDuration {
			elapsed = elapsedWall
		}
	}

	if elapsed < sm.segmentDuration {
		return false
	}

	if sm.hasVideo {
		// For video streams split only on keyframes to keep segments decodable.
		return pkt.Type == media.PacketTypeVideo && pkt.IsKeyframe
	}
	return true
}

func (sm *StreamMuxer) startNewSegment(timestamp uint32) error {
	// Close previous segment
	if sm.currentSegment != nil {
		sm.currentSegment.Close()

		// Record segment info – prefer timestamp-based duration but fall
		// back to wall-clock when timestamps look implausible (< 250 ms
		// while wall-clock shows a realistic value).
		deltaMS := uint32(0)
		if timestamp >= sm.segmentStart {
			deltaMS = timestamp - sm.segmentStart
		}
		duration := float64(deltaMS) / 1000.0
		wallDuration := time.Since(sm.segmentStartedAt).Seconds()
		if duration < 0.25 && wallDuration >= 0.5 {
			duration = wallDuration
		}
		if duration <= 0 {
			duration = wallDuration
		}
		if duration <= 0 {
			duration = sm.segmentDuration.Seconds()
		}
		sm.segments = append(sm.segments, SegmentInfo{
			Index:    sm.currentSegIdx,
			Duration: duration,
			Filename: fmt.Sprintf("seg_%d.ts", sm.currentSegIdx),
			Size:     sm.segmentBytes,
		})

		// Clean up old segments
		for len(sm.segments) > sm.maxSegments {
			old := sm.segments[0]
			os.Remove(filepath.Join(sm.outputDir, old.Filename))
			sm.segments = sm.segments[1:]
		}

		sm.currentSegIdx++
	}

	// Create new segment file
	filename := fmt.Sprintf("seg_%d.ts", sm.currentSegIdx)
	f, err := os.Create(filepath.Join(sm.outputDir, filename))
	if err != nil {
		return fmt.Errorf("create segment: %w", err)
	}

	sm.currentSegment = f
	sm.segmentStart = timestamp
	sm.segmentStartedAt = time.Now()
	sm.segmentBytes = 0

	// Write PAT/PMT at the start of each segment
	patpmt := sm.tsMuxer.GeneratePatPmt()
	n, _ := f.Write(patpmt)
	sm.segmentBytes += n

	// Update playlist
	sm.writePlaylist()

	return nil
}

// writePlaylist generates the M3U8 playlist file
func (sm *StreamMuxer) writePlaylist() error {
	playlistPath := filepath.Join(sm.outputDir, "index.m3u8")
	segments := sm.playlistSegments()

	// Calculate target duration
	targetDuration := sm.segmentDuration.Seconds()
	for _, seg := range segments {
		if seg.Duration > targetDuration {
			targetDuration = seg.Duration
		}
	}

	content := "#EXTM3U\n"
	content += "#EXT-X-VERSION:3\n"
	content += fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", int(targetDuration)+1)

	if len(segments) > 0 {
		content += fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", segments[0].Index)
	} else {
		content += "#EXT-X-MEDIA-SEQUENCE:0\n"
	}

	for _, seg := range segments {
		content += fmt.Sprintf("#EXTINF:%.3f,\n", seg.Duration)
		content += seg.Filename + "\n"
	}

	return os.WriteFile(playlistPath, []byte(content), 0644)
}

// Close closes the current segment
func (sm *StreamMuxer) Close() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.currentSegment != nil {
		// Finalize the last segment
		if sm.segmentBytes > 0 {
			duration := time.Since(sm.segmentStartedAt).Seconds()
			if duration <= 0 {
				duration = sm.segmentDuration.Seconds()
			}
			sm.segments = append(sm.segments, SegmentInfo{
				Index:    sm.currentSegIdx,
				Duration: duration,
				Filename: fmt.Sprintf("seg_%d.ts", sm.currentSegIdx),
				Size:     sm.segmentBytes,
			})
		}
		sm.currentSegment.Close()
		sm.currentSegment = nil

		// Write final playlist with ENDLIST
		sm.writeEndPlaylist()
	}
}

func (sm *StreamMuxer) writeEndPlaylist() {
	playlistPath := filepath.Join(sm.outputDir, "index.m3u8")
	segments := sm.playlistSegments()
	targetDuration := sm.segmentDuration.Seconds()
	for _, seg := range segments {
		if seg.Duration > targetDuration {
			targetDuration = seg.Duration
		}
	}

	content := "#EXTM3U\n"
	content += "#EXT-X-VERSION:3\n"
	content += fmt.Sprintf("#EXT-X-TARGETDURATION:%d\n", int(targetDuration)+1)
	if len(segments) > 0 {
		content += fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%d\n", segments[0].Index)
	}
	for _, seg := range segments {
		content += fmt.Sprintf("#EXTINF:%.3f,\n", seg.Duration)
		content += seg.Filename + "\n"
	}
	content += "#EXT-X-ENDLIST\n"

	os.WriteFile(playlistPath, []byte(content), 0644)
}

func (sm *StreamMuxer) playlistSegments() []SegmentInfo {
	if len(sm.segments) == 0 {
		return nil
	}
	out := make([]SegmentInfo, 0, len(sm.segments))
	for _, seg := range sm.segments {
		if info, err := os.Stat(filepath.Join(sm.outputDir, seg.Filename)); err == nil && !info.IsDir() {
			out = append(out, seg)
		}
	}
	return out
}

// GetPlaylistPath returns the path to the M3U8 playlist
func (sm *StreamMuxer) GetPlaylistPath() string {
	return filepath.Join(sm.outputDir, "index.m3u8")
}

// IsActive returns true if the muxer is actively writing
func (sm *StreamMuxer) IsActive() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.currentSegment != nil
}
