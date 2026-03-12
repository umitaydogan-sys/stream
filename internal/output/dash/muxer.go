package dash

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
)

// Muxer manages all DASH outputs
type Muxer struct {
	outputDir string
	streams   map[string]*StreamMuxer
	mu        sync.RWMutex
}

// StreamMuxer handles DASH/CMAF output for a single stream
type StreamMuxer struct {
	streamKey       string
	outputDir       string
	segmentDuration time.Duration
	maxSegments     int

	currentVideo   []byte
	currentAudio   []byte
	segIdx         int
	segStart       uint32
	segStarted     bool
	segStartedAt   time.Time
	segments       []segmentInfo
	videoSeqHeader []byte
	audioSeqHeader []byte
	videoTrackID   uint32
	audioTrackID   uint32
	hasInit        bool
	mu             sync.Mutex
}

type segmentInfo struct {
	Index    int
	Duration float64
	Video    string
	Audio    string
}

// NewMuxer creates a new DASH muxer
func NewMuxer(outputDir string) *Muxer {
	return &Muxer{
		outputDir: outputDir,
		streams:   make(map[string]*StreamMuxer),
	}
}

// AddStream creates DASH output for a stream key
func (m *Muxer) AddStream(streamKey string) *StreamMuxer {
	m.mu.Lock()
	defer m.mu.Unlock()

	dir := filepath.Join(m.outputDir, streamKey)
	os.MkdirAll(dir, 0755)

	sm := &StreamMuxer{
		streamKey:       streamKey,
		outputDir:       dir,
		segmentDuration: 2 * time.Second,
		maxSegments:     6,
		videoTrackID:    1,
		audioTrackID:    2,
	}
	m.streams[streamKey] = sm
	return sm
}

// RemoveStream removes DASH output
func (m *Muxer) RemoveStream(streamKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sm, ok := m.streams[streamKey]; ok {
		sm.Close()
		delete(m.streams, streamKey)
	}
}

// GetOutputDir returns the dash output directory
func (m *Muxer) GetOutputDir() string {
	return m.outputDir
}

// WritePacket writes a media packet to DASH/CMAF output
func (sm *StreamMuxer) WritePacket(pkt *media.Packet) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Store sequence headers
	if pkt.IsSequenceHeader {
		if pkt.Type == media.PacketTypeVideo {
			sm.videoSeqHeader = make([]byte, len(pkt.Data))
			copy(sm.videoSeqHeader, pkt.Data)
		} else if pkt.Type == media.PacketTypeAudio {
			sm.audioSeqHeader = make([]byte, len(pkt.Data))
			copy(sm.audioSeqHeader, pkt.Data)
		}
		if !sm.hasInit && sm.videoSeqHeader != nil {
			sm.writeInitSegment()
			sm.hasInit = true
		}
		return nil
	}

	// Check segment split on keyframe
	if sm.shouldSplit(pkt) {
		sm.flushSegment(pkt.Timestamp)
	}

	if !sm.segStarted {
		sm.segStart = pkt.Timestamp
		sm.segStarted = true
		sm.segStartedAt = time.Now()
	}

	// Buffer media data as fMP4 mdat content
	data := pkt.Data
	if pkt.Type == media.PacketTypeVideo && len(data) > 5 {
		data = data[5:] // strip FLV AVC header
	} else if pkt.Type == media.PacketTypeAudio && len(data) > 2 {
		data = data[2:] // strip FLV AAC header
	}

	// Build moof+mdat sample (simplified CMAF)
	delta := uint32(0)
	if pkt.Timestamp >= sm.segStart {
		delta = pkt.Timestamp - sm.segStart
	}
	sample := buildSample(data, delta, pkt.IsKeyframe)

	if pkt.Type == media.PacketTypeVideo {
		sm.currentVideo = append(sm.currentVideo, sample...)
	} else if pkt.Type == media.PacketTypeAudio {
		sm.currentAudio = append(sm.currentAudio, sample...)
	}

	return nil
}

func (sm *StreamMuxer) shouldSplit(pkt *media.Packet) bool {
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
	if sm.videoSeqHeader != nil {
		return pkt.Type == media.PacketTypeVideo && pkt.IsKeyframe
	}
	return true
}

func (sm *StreamMuxer) flushSegment(ts uint32) {
	if !sm.segStarted || (len(sm.currentVideo) == 0 && len(sm.currentAudio) == 0) {
		return
	}

	deltaMS := uint32(0)
	if ts >= sm.segStart {
		deltaMS = ts - sm.segStart
	}
	duration := float64(deltaMS) / 1000.0
	if duration <= 0 {
		duration = time.Since(sm.segStartedAt).Seconds()
	}
	if duration <= 0 {
		duration = sm.segmentDuration.Seconds()
	}

	// Write video segment
	if len(sm.currentVideo) > 0 {
		vFile := fmt.Sprintf("video_%d.m4s", sm.segIdx)
		writeSegmentFile(filepath.Join(sm.outputDir, vFile), sm.currentVideo, sm.segIdx, sm.videoTrackID, sm.segStart)
		sm.currentVideo = nil
	}

	// Write audio segment
	if len(sm.currentAudio) > 0 {
		aFile := fmt.Sprintf("audio_%d.m4s", sm.segIdx)
		writeSegmentFile(filepath.Join(sm.outputDir, aFile), sm.currentAudio, sm.segIdx, sm.audioTrackID, sm.segStart)
		sm.currentAudio = nil
	}

	sm.segments = append(sm.segments, segmentInfo{
		Index:    sm.segIdx,
		Duration: duration,
		Video:    fmt.Sprintf("video_%d.m4s", sm.segIdx),
		Audio:    fmt.Sprintf("audio_%d.m4s", sm.segIdx),
	})

	// Cleanup old segments
	for len(sm.segments) > sm.maxSegments {
		old := sm.segments[0]
		os.Remove(filepath.Join(sm.outputDir, old.Video))
		os.Remove(filepath.Join(sm.outputDir, old.Audio))
		sm.segments = sm.segments[1:]
	}

	sm.segIdx++
	sm.segStart = ts
	sm.segStarted = true
	sm.segStartedAt = time.Now()
	sm.writeMPD()
}

// writeInitSegment writes the fMP4 initialization segment (ftyp + moov)
func (sm *StreamMuxer) writeInitSegment() {
	// ftyp box
	ftyp := buildBox("ftyp", []byte("iso6\x00\x00\x02\x00iso6mp41dash"))

	// Simplified moov with video track
	moov := buildFMP4Moov(sm.videoTrackID, sm.audioTrackID, sm.videoSeqHeader, sm.audioSeqHeader)

	initData := append(ftyp, moov...)
	path := filepath.Join(sm.outputDir, "init.mp4")
	os.WriteFile(path, initData, 0644)
}

// writeMPD writes the DASH manifest (MPD)
func (sm *StreamMuxer) writeMPD() {
	targetDur := sm.segmentDuration.Seconds()
	for _, s := range sm.segments {
		if s.Duration > targetDur {
			targetDur = s.Duration
		}
	}

	minBufferTime := fmt.Sprintf("PT%.1fS", targetDur)

	mpd := `<?xml version="1.0" encoding="UTF-8"?>` + "\n"
	mpd += `<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" type="dynamic" ` +
		`minimumUpdatePeriod="PT2S" minBufferTime="` + minBufferTime + `" ` +
		`availabilityStartTime="` + time.Now().UTC().Format(time.RFC3339) + `" ` +
		`profiles="urn:mpeg:dash:profile:isoff-live:2011">` + "\n"
	mpd += `  <Period>` + "\n"

	// Video AdaptationSet
	mpd += `    <AdaptationSet mimeType="video/mp4" codecs="avc1.64001f" startWithSAP="1">` + "\n"
	mpd += `      <Representation id="video" bandwidth="2000000">` + "\n"
	mpd += `        <SegmentTemplate media="video_$Number$.m4s" initialization="init.mp4" startNumber="` +
		fmt.Sprintf("%d", sm.firstSegNum()) + `" timescale="1000">` + "\n"
	mpd += `          <SegmentTimeline>` + "\n"
	for _, seg := range sm.segments {
		mpd += fmt.Sprintf(`            <S d="%d"/>`, int(seg.Duration*1000)) + "\n"
	}
	mpd += `          </SegmentTimeline>` + "\n"
	mpd += `        </SegmentTemplate>` + "\n"
	mpd += `      </Representation>` + "\n"
	mpd += `    </AdaptationSet>` + "\n"

	// Audio AdaptationSet
	mpd += `    <AdaptationSet mimeType="audio/mp4" codecs="mp4a.40.2" startWithSAP="1">` + "\n"
	mpd += `      <Representation id="audio" bandwidth="128000">` + "\n"
	mpd += `        <SegmentTemplate media="audio_$Number$.m4s" initialization="init.mp4" startNumber="` +
		fmt.Sprintf("%d", sm.firstSegNum()) + `" timescale="1000">` + "\n"
	mpd += `          <SegmentTimeline>` + "\n"
	for _, seg := range sm.segments {
		mpd += fmt.Sprintf(`            <S d="%d"/>`, int(seg.Duration*1000)) + "\n"
	}
	mpd += `          </SegmentTimeline>` + "\n"
	mpd += `        </SegmentTemplate>` + "\n"
	mpd += `      </Representation>` + "\n"
	mpd += `    </AdaptationSet>` + "\n"

	mpd += `  </Period>` + "\n"
	mpd += `</MPD>` + "\n"

	os.WriteFile(filepath.Join(sm.outputDir, "manifest.mpd"), []byte(mpd), 0644)
}

func (sm *StreamMuxer) firstSegNum() int {
	if len(sm.segments) > 0 {
		return sm.segments[0].Index
	}
	return 0
}

// Close finalizes the DASH output
func (sm *StreamMuxer) Close() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if !sm.segStarted {
		log.Printf("[DASH] Muxer kapatıldı: %s", sm.streamKey)
		return
	}
	sm.flushSegment(sm.segStart + uint32(sm.segmentDuration.Milliseconds()))
	log.Printf("[DASH] Muxer kapatıldı: %s", sm.streamKey)
}

// ── fMP4 Box Helpers ──

func buildBox(boxType string, data []byte) []byte {
	size := uint32(8 + len(data))
	buf := make([]byte, size)
	binary.BigEndian.PutUint32(buf[0:4], size)
	copy(buf[4:8], boxType)
	copy(buf[8:], data)
	return buf
}

func buildSample(data []byte, dts uint32, isKey bool) []byte {
	// NAL length prefix (4 bytes big-endian)
	nalLen := make([]byte, 4)
	binary.BigEndian.PutUint32(nalLen, uint32(len(data)))
	return append(nalLen, data...)
}

func writeSegmentFile(path string, samples []byte, seqNum int, trackID uint32, baseDecodeTime uint32) {
	// styp box
	styp := buildBox("styp", []byte("msdh\x00\x00\x00\x00msdhmsix"))

	// moof box
	moof := buildMoof(seqNum, trackID, baseDecodeTime, len(samples))

	// mdat box
	mdat := buildBox("mdat", samples)

	data := append(styp, moof...)
	data = append(data, mdat...)
	os.WriteFile(path, data, 0644)
}

func buildMoof(seqNum int, trackID uint32, baseDecodeTime uint32, dataSize int) []byte {
	// mfhd box (Movie Fragment Header)
	mfhdData := make([]byte, 8)
	binary.BigEndian.PutUint32(mfhdData[4:8], uint32(seqNum))
	mfhd := buildBox("mfhd", mfhdData)

	// tfhd box (Track Fragment Header)
	tfhdData := make([]byte, 12)
	binary.BigEndian.PutUint32(tfhdData[0:4], 0x020000) // flags: default-base-is-moof
	binary.BigEndian.PutUint32(tfhdData[4:8], trackID)
	tfhd := buildBox("tfhd", tfhdData)

	// tfdt box (Track Fragment Decode Time)
	tfdtData := make([]byte, 12)
	binary.BigEndian.PutUint32(tfdtData[0:4], 0x01000000) // version 1
	binary.BigEndian.PutUint64(tfdtData[4:12], uint64(baseDecodeTime))
	tfdt := buildBox("tfdt", tfdtData)

	// trun box (Track Run)
	trunData := make([]byte, 12)
	binary.BigEndian.PutUint32(trunData[0:4], 0x000201) // flags: data-offset, sample-size
	binary.BigEndian.PutUint32(trunData[4:8], 1)        // sample count
	binary.BigEndian.PutUint32(trunData[8:12], uint32(dataSize))
	trun := buildBox("trun", trunData)

	// traf box (Track Fragment)
	trafContent := append(tfhd, tfdt...)
	trafContent = append(trafContent, trun...)
	traf := buildBox("traf", trafContent)

	// moof box
	moofContent := append(mfhd, traf...)
	return buildBox("moof", moofContent)
}

func buildFMP4Moov(videoTrackID, audioTrackID uint32, videoSeqHeader, audioSeqHeader []byte) []byte {
	// Simplified moov box with mvhd + video trak + audio trak + mvex
	mvhd := buildBox("mvhd", make([]byte, 100)) // minimal mvhd

	// Video trak (minimal)
	videoTrak := buildMinimalTrak(videoTrackID, "video", videoSeqHeader)

	// Audio trak (minimal)
	audioTrak := buildMinimalTrak(audioTrackID, "audio", audioSeqHeader)

	// mvex with trex for each track
	trex1Data := make([]byte, 24)
	binary.BigEndian.PutUint32(trex1Data[4:8], videoTrackID)
	binary.BigEndian.PutUint32(trex1Data[8:12], 1)  // default sample description
	binary.BigEndian.PutUint32(trex1Data[16:20], 0) // default sample size
	trex1 := buildBox("trex", trex1Data)

	trex2Data := make([]byte, 24)
	binary.BigEndian.PutUint32(trex2Data[4:8], audioTrackID)
	binary.BigEndian.PutUint32(trex2Data[8:12], 1)
	trex2 := buildBox("trex", trex2Data)

	mvex := buildBox("mvex", append(trex1, trex2...))

	moovContent := append(mvhd, videoTrak...)
	moovContent = append(moovContent, audioTrak...)
	moovContent = append(moovContent, mvex...)
	return buildBox("moov", moovContent)
}

func buildMinimalTrak(trackID uint32, mediaType string, seqHeader []byte) []byte {
	// tkhd
	tkhdData := make([]byte, 84)
	binary.BigEndian.PutUint32(tkhdData[0:4], 0x00000003) // flags enabled+in_movie
	binary.BigEndian.PutUint32(tkhdData[12:16], trackID)
	if mediaType == "video" {
		binary.BigEndian.PutUint32(tkhdData[76:80], 1920<<16) // width
		binary.BigEndian.PutUint32(tkhdData[80:84], 1080<<16) // height
	}
	tkhd := buildBox("tkhd", tkhdData)

	// mdia with mdhd + hdlr + minf
	mdhdData := make([]byte, 24)
	binary.BigEndian.PutUint32(mdhdData[12:16], 90000) // timescale
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

	// minf with stbl (empty sample table for fMP4)
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
