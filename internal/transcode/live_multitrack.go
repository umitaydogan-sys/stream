package transcode

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/output/hls"
)

type liveTrackBootstrap struct {
	videoSeq map[uint8]*media.Packet
	audioSeq *media.Packet
}

type liveDirectSession struct {
	streamKey      string
	outputDir      string
	parentMuxer    *hls.Muxer
	primaryTrackID uint8
	videoSeq       map[uint8]*media.Packet
	audioSeq       *media.Packet
	variants       map[uint8]*directVariant
	lastMasterAt   time.Time
	mu             sync.Mutex
}

type directVariant struct {
	trackID       uint8
	streamName    string
	playlistPath  string
	width         int
	height        int
	bitrate       int64
	packetBytes   int64
	windowStarted time.Time
	muxer         *hls.StreamMuxer
}

func newLiveTrackBootstrap() *liveTrackBootstrap {
	return &liveTrackBootstrap{
		videoSeq: make(map[uint8]*media.Packet),
	}
}

func (m *Manager) cacheLiveBootstrap(streamKey string, pkt *media.Packet) *liveTrackBootstrap {
	m.mu.Lock()
	defer m.mu.Unlock()

	boot := m.liveBoot[streamKey]
	if boot == nil {
		boot = newLiveTrackBootstrap()
		m.liveBoot[streamKey] = boot
	}
	if pkt == nil {
		return boot
	}

	if pkt.IsSequenceHeader {
		switch pkt.Type {
		case media.PacketTypeVideo:
			boot.videoSeq[pkt.TrackID] = pkt.Clone()
		case media.PacketTypeAudio:
			boot.audioSeq = pkt.Clone()
		}
	}
	return boot
}

func (m *Manager) getLiveDirectSession(streamKey string) *liveDirectSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.liveDirect[streamKey]
}

func (m *Manager) tryActivateDirectMultitrackSession(streamKey string, pkt *media.Packet, boot *liveTrackBootstrap) *liveDirectSession {
	if pkt == nil || pkt.Type != media.PacketTypeVideo || pkt.TrackID == 0 || !pkt.IsEnhanced {
		return nil
	}

	opts := m.getLiveOptions(streamKey)
	if !(opts.ABREnabled && opts.MasterEnabled) {
		return nil
	}

	m.mu.RLock()
	if session := m.liveDirect[streamKey]; session != nil {
		m.mu.RUnlock()
		return session
	}
	m.mu.RUnlock()

	rootDir := filepath.Join(m.GetLiveOutputDir(), streamKey)
	session := newLiveDirectSession(streamKey, rootDir, boot, pkt.TrackID)
	hadDash := m.stopLivePipelineForDirectSwitch(streamKey)
	cleanupLiveOutputDir(rootDir)
	session.bootstrapKnownTracks()

	m.mu.Lock()
	if existing := m.liveDirect[streamKey]; existing != nil {
		m.mu.Unlock()
		return existing
	}
	m.liveDirect[streamKey] = session
	m.mu.Unlock()

	log.Printf("[TC] OBS multitrack dogrudan ABR HLS aktif edildi: %s", streamKey)
	if hadDash {
		go func() {
			time.Sleep(1500 * time.Millisecond)
			if _, err := m.StartLiveDASH(streamKey); err != nil {
				log.Printf("[TC] Direct multitrack sonrasi DASH yeniden baslatilamadi (%s): %v", streamKey, err)
			}
		}()
	}
	return session
}

func (m *Manager) stopLivePipelineForDirectSwitch(streamKey string) bool {
	m.mu.Lock()
	job := m.liveJobs[streamKey]
	if job != nil {
		delete(m.liveJobs, streamKey)
	}
	dashJob := m.liveDash[streamKey]
	if dashJob != nil {
		delete(m.liveDash, streamKey)
	}
	m.mu.Unlock()

	if job != nil {
		job.Status = "completed"
		job.closeInput()
		if job.cancel != nil {
			job.cancel()
		}
	}
	if dashJob != nil {
		dashJob.Status = "completed"
		if dashJob.cancel != nil {
			dashJob.cancel()
		}
	}
	return dashJob != nil
}

func newLiveDirectSession(streamKey, outputDir string, boot *liveTrackBootstrap, fallbackPrimary uint8) *liveDirectSession {
	session := &liveDirectSession{
		streamKey:   streamKey,
		outputDir:   outputDir,
		parentMuxer: hls.NewMuxer(outputDir),
		videoSeq:    make(map[uint8]*media.Packet),
		variants:    make(map[uint8]*directVariant),
	}
	if boot != nil {
		for trackID, seq := range boot.videoSeq {
			session.videoSeq[trackID] = seq.Clone()
		}
		if boot.audioSeq != nil {
			session.audioSeq = boot.audioSeq.Clone()
		}
	}
	if trackID, ok := session.pickPrimaryTrack(); ok {
		session.primaryTrackID = trackID
	} else {
		session.primaryTrackID = fallbackPrimary
	}
	return session
}

func (s *liveDirectSession) pickPrimaryTrack() (uint8, bool) {
	if len(s.videoSeq) == 0 {
		return 0, false
	}
	type candidate struct {
		trackID   uint8
		width     int
		height    int
		bandwidth int64
	}
	candidates := make([]candidate, 0, len(s.videoSeq))
	for trackID, seq := range s.videoSeq {
		width, height := 0, 0
		if seq != nil {
			width, height = parseAVCSequenceHeaderDimensions(seq.Data)
		}
		candidates = append(candidates, candidate{
			trackID:   trackID,
			width:     width,
			height:    height,
			bandwidth: fallbackVariantBandwidth(width, height),
		})
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].bandwidth == candidates[j].bandwidth {
			if candidates[i].height == candidates[j].height {
				if candidates[i].width == candidates[j].width {
					return candidates[i].trackID < candidates[j].trackID
				}
				return candidates[i].width < candidates[j].width
			}
			return candidates[i].height < candidates[j].height
		}
		return candidates[i].bandwidth < candidates[j].bandwidth
	})
	return candidates[0].trackID, true
}

func (s *liveDirectSession) bootstrapKnownTracks() {
	s.mu.Lock()
	defer s.mu.Unlock()

	trackIDs := make([]int, 0, len(s.videoSeq))
	for trackID := range s.videoSeq {
		trackIDs = append(trackIDs, int(trackID))
	}
	sort.Ints(trackIDs)
	for _, rawTrackID := range trackIDs {
		trackID := uint8(rawTrackID)
		seq := s.videoSeq[trackID]
		if seq == nil {
			continue
		}
		s.ensureVariantLocked(trackID)
	}
	s.writeMasterPlaylistLocked(true)
}

func (s *liveDirectSession) writePacket(pkt *media.Packet) {
	if pkt == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if pkt.IsSequenceHeader {
		switch pkt.Type {
		case media.PacketTypeVideo:
			s.videoSeq[pkt.TrackID] = pkt.Clone()
			if variant := s.ensureVariantLocked(pkt.TrackID); variant != nil {
				_ = variant.muxer.WritePacket(pkt.Clone())
				s.writeMasterPlaylistLocked(true)
			}
			return
		case media.PacketTypeAudio:
			s.audioSeq = pkt.Clone()
		}
	}

	switch pkt.Type {
	case media.PacketTypeAudio:
		for _, variant := range s.variants {
			_ = variant.muxer.WritePacket(pkt.Clone())
		}
	case media.PacketTypeVideo:
		variant := s.ensureVariantLocked(pkt.TrackID)
		if variant == nil {
			return
		}
		variant.observePacket(pkt)
		_ = variant.muxer.WritePacket(pkt.Clone())
		s.writeMasterPlaylistLocked(false)
	}
}

func (s *liveDirectSession) ensureVariantLocked(trackID uint8) *directVariant {
	if variant, ok := s.variants[trackID]; ok {
		return variant
	}

	videoSeq := s.videoSeq[trackID]
	if videoSeq == nil {
		return nil
	}

	width, height := parseAVCSequenceHeaderDimensions(videoSeq.Data)
	streamName := ""
	playlistPath := "index.m3u8"
	if trackID != s.primaryTrackID {
		streamName = variantStreamName(trackID, width, height)
		playlistPath = filepath.ToSlash(filepath.Join(streamName, "index.m3u8"))
	}
	muxer := s.parentMuxer.AddStream(streamName)
	if s.audioSeq != nil {
		_ = muxer.WritePacket(s.audioSeq.Clone())
	}
	_ = muxer.WritePacket(videoSeq.Clone())

	variant := &directVariant{
		trackID:      trackID,
		streamName:   streamName,
		playlistPath: playlistPath,
		width:        width,
		height:       height,
		muxer:        muxer,
	}
	s.variants[trackID] = variant
	return variant
}

func (s *liveDirectSession) writeMasterPlaylistLocked(force bool) {
	if len(s.variants) == 0 {
		return
	}
	if !force && time.Since(s.lastMasterAt) < time.Second {
		return
	}

	type item struct {
		variant *directVariant
	}
	items := make([]item, 0, len(s.variants))
	for _, variant := range s.variants {
		target := filepath.Join(s.outputDir, filepath.FromSlash(variant.playlistPath))
		if info, err := os.Stat(target); err == nil && !info.IsDir() && info.Size() > 0 {
			items = append(items, item{variant: variant})
		}
	}
	if len(items) == 0 {
		return
	}

	sort.Slice(items, func(i, j int) bool {
		leftBandwidth := items[i].variant.bandwidthEstimate()
		rightBandwidth := items[j].variant.bandwidthEstimate()
		if leftBandwidth == rightBandwidth {
			if items[i].variant.height == items[j].variant.height {
				return items[i].variant.trackID < items[j].variant.trackID
			}
			return items[i].variant.height < items[j].variant.height
		}
		return leftBandwidth < rightBandwidth
	})

	var b strings.Builder
	b.WriteString("#EXTM3U\n")
	b.WriteString("#EXT-X-VERSION:3\n")
	for _, item := range items {
		v := item.variant
		bandwidth := v.bandwidthEstimate()
		if bandwidth <= 0 {
			bandwidth = int(fallbackVariantBandwidth(v.width, v.height))
		}
		codecs := `avc1.64001f,mp4a.40.2`
		if v.width > 0 && v.height > 0 {
			b.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,AVERAGE-BANDWIDTH=%d,RESOLUTION=%dx%d,CODECS=\"%s\"\n", bandwidth, bandwidth, v.width, v.height, codecs))
		} else {
			b.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,AVERAGE-BANDWIDTH=%d,CODECS=\"%s\"\n", bandwidth, bandwidth, codecs))
		}
		b.WriteString(v.playlistPath)
		b.WriteString("\n")
	}

	tmpPath := filepath.Join(s.outputDir, "master.m3u8.tmp")
	finalPath := filepath.Join(s.outputDir, "master.m3u8")
	if err := os.WriteFile(tmpPath, []byte(b.String()), 0644); err == nil {
		_ = os.Rename(tmpPath, finalPath)
	}
	s.lastMasterAt = time.Now()
}

func playlistLooksStable(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	totalDurations := 0
	shortDurations := 0
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "#EXTINF:") {
			continue
		}
		totalDurations++
		value := strings.TrimSuffix(strings.TrimPrefix(line, "#EXTINF:"), ",")
		duration, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			continue
		}
		if duration > 0 && duration < 0.25 {
			shortDurations++
		}
	}

	if totalDurations == 0 {
		return false
	}
	return shortDurations == 0
}

func (s *liveDirectSession) close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, variant := range s.variants {
		if variant != nil && variant.muxer != nil {
			variant.muxer.Close()
		}
	}
}

func (v *directVariant) observePacket(pkt *media.Packet) {
	if pkt == nil || pkt.Type != media.PacketTypeVideo || pkt.IsSequenceHeader {
		return
	}
	if v.windowStarted.IsZero() {
		v.windowStarted = time.Now()
	}
	v.packetBytes += int64(len(pkt.Data))
	if elapsed := time.Since(v.windowStarted); elapsed >= 2*time.Second {
		v.bitrate = int64(float64(v.packetBytes*8)/elapsed.Seconds()) + 128000
		v.packetBytes = 0
		v.windowStarted = time.Now()
	}
}

func (v *directVariant) bandwidthEstimate() int {
	if v.bitrate > 0 {
		return int(v.bitrate)
	}
	return int(fallbackVariantBandwidth(v.width, v.height))
}

func variantStreamName(trackID uint8, width, height int) string {
	switch {
	case height > 0:
		return fmt.Sprintf("obs_%dp_t%d", height, trackID)
	case width > 0:
		return fmt.Sprintf("obs_%dw_t%d", width, trackID)
	default:
		return fmt.Sprintf("obs_track_%d", trackID)
	}
}

func fallbackVariantBandwidth(width, height int) int64 {
	switch {
	case height >= 1080 || width >= 1920:
		return 5000000
	case height >= 720 || width >= 1280:
		return 3000000
	case height >= 480 || width >= 854:
		return 1400000
	case height >= 360 || width >= 640:
		return 800000
	default:
		return 600000
	}
}

func parseAVCSequenceHeaderDimensions(data []byte) (int, int) {
	sps := extractAVCConfigSPS(data)
	if len(sps) == 0 {
		return 0, 0
	}
	return parseH264SPSDimensions(sps)
}

func extractAVCConfigSPS(data []byte) []byte {
	if len(data) <= 10 {
		return nil
	}
	config := data[5:]
	if len(config) < 7 {
		return nil
	}
	pos := 6
	numSPS := int(config[5] & 0x1F)
	if numSPS <= 0 {
		return nil
	}
	if pos+2 > len(config) {
		return nil
	}
	l := int(config[pos])<<8 | int(config[pos+1])
	pos += 2
	if l <= 0 || pos+l > len(config) {
		return nil
	}
	return append([]byte{}, config[pos:pos+l]...)
}

func parseH264SPSDimensions(sps []byte) (int, int) {
	if len(sps) == 0 {
		return 0, 0
	}
	rbsp := removeEmulationBytes(sps)
	br := &bitReader{data: rbsp}

	if _, ok := br.readBits(8); !ok { // nal header
		return 0, 0
	}
	profileIDC, ok := br.readBits(8)
	if !ok {
		return 0, 0
	}
	if _, ok = br.readBits(8); !ok { // constraints
		return 0, 0
	}
	if _, ok = br.readBits(8); !ok { // level idc
		return 0, 0
	}
	if _, ok = br.readUE(); !ok { // sps id
		return 0, 0
	}

	chromaFormatIDC := uint(1)
	if isHighProfile(profileIDC) {
		if chromaFormatIDC, ok = br.readUE(); !ok {
			return 0, 0
		}
		if chromaFormatIDC == 3 {
			if _, ok = br.readBit(); !ok {
				return 0, 0
			}
		}
		if _, ok = br.readUE(); !ok {
			return 0, 0
		}
		if _, ok = br.readUE(); !ok {
			return 0, 0
		}
		if _, ok = br.readBit(); !ok {
			return 0, 0
		}
		seqScalingMatrixPresent, ok := br.readBit()
		if !ok {
			return 0, 0
		}
		if seqScalingMatrixPresent == 1 {
			count := 8
			if chromaFormatIDC == 3 {
				count = 12
			}
			for i := 0; i < count; i++ {
				present, ok := br.readBit()
				if !ok {
					return 0, 0
				}
				if present == 1 {
					size := 16
					if i >= 6 {
						size = 64
					}
					if !skipScalingList(br, size) {
						return 0, 0
					}
				}
			}
		}
	}

	if _, ok = br.readUE(); !ok {
		return 0, 0
	}
	picOrderCntType, ok := br.readUE()
	if !ok {
		return 0, 0
	}
	if picOrderCntType == 0 {
		if _, ok = br.readUE(); !ok {
			return 0, 0
		}
	} else if picOrderCntType == 1 {
		if _, ok = br.readBit(); !ok {
			return 0, 0
		}
		if _, ok = br.readSE(); !ok {
			return 0, 0
		}
		if _, ok = br.readSE(); !ok {
			return 0, 0
		}
		numRefFramesCycle, ok := br.readUE()
		if !ok {
			return 0, 0
		}
		for i := uint(0); i < numRefFramesCycle; i++ {
			if _, ok = br.readSE(); !ok {
				return 0, 0
			}
		}
	}

	if _, ok = br.readUE(); !ok {
		return 0, 0
	}
	if _, ok = br.readBit(); !ok {
		return 0, 0
	}
	picWidthInMbsMinus1, ok := br.readUE()
	if !ok {
		return 0, 0
	}
	picHeightInMapUnitsMinus1, ok := br.readUE()
	if !ok {
		return 0, 0
	}
	frameMbsOnlyFlag, ok := br.readBit()
	if !ok {
		return 0, 0
	}
	if frameMbsOnlyFlag == 0 {
		if _, ok = br.readBit(); !ok {
			return 0, 0
		}
	}
	if _, ok = br.readBit(); !ok {
		return 0, 0
	}
	frameCroppingFlag, ok := br.readBit()
	if !ok {
		return 0, 0
	}

	var cropLeft, cropRight, cropTop, cropBottom uint
	if frameCroppingFlag == 1 {
		if cropLeft, ok = br.readUE(); !ok {
			return 0, 0
		}
		if cropRight, ok = br.readUE(); !ok {
			return 0, 0
		}
		if cropTop, ok = br.readUE(); !ok {
			return 0, 0
		}
		if cropBottom, ok = br.readUE(); !ok {
			return 0, 0
		}
	}

	width := int((picWidthInMbsMinus1 + 1) * 16)
	height := int((2 - uint(frameMbsOnlyFlag)) * (picHeightInMapUnitsMinus1 + 1) * 16)

	cropUnitX, cropUnitY := h264CropUnits(chromaFormatIDC, frameMbsOnlyFlag == 1)
	width -= int((cropLeft + cropRight) * cropUnitX)
	height -= int((cropTop + cropBottom) * cropUnitY)
	if width < 0 || height < 0 {
		return 0, 0
	}
	return width, height
}

func isHighProfile(profileIDC uint) bool {
	switch profileIDC {
	case 100, 110, 122, 244, 44, 83, 86, 118, 128, 138, 139, 134, 135:
		return true
	default:
		return false
	}
}

func h264CropUnits(chromaFormatIDC uint, frameMbsOnly bool) (uint, uint) {
	switch chromaFormatIDC {
	case 0:
		if frameMbsOnly {
			return 1, 2
		}
		return 1, 4
	case 1:
		if frameMbsOnly {
			return 2, 2
		}
		return 2, 4
	case 2:
		if frameMbsOnly {
			return 2, 1
		}
		return 2, 2
	case 3:
		if frameMbsOnly {
			return 1, 1
		}
		return 1, 2
	default:
		return 1, 2
	}
}

func skipScalingList(br *bitReader, size int) bool {
	lastScale := 8
	nextScale := 8
	for i := 0; i < size; i++ {
		if nextScale != 0 {
			deltaScale, ok := br.readSE()
			if !ok {
				return false
			}
			nextScale = (lastScale + int(deltaScale) + 256) % 256
		}
		if nextScale != 0 {
			lastScale = nextScale
		}
	}
	return true
}

func removeEmulationBytes(data []byte) []byte {
	if len(data) < 3 {
		return append([]byte{}, data...)
	}
	out := make([]byte, 0, len(data))
	for i := 0; i < len(data); i++ {
		if i+2 < len(data) && data[i] == 0x00 && data[i+1] == 0x00 && data[i+2] == 0x03 {
			out = append(out, 0x00, 0x00)
			i += 2
			continue
		}
		out = append(out, data[i])
	}
	return out
}

type bitReader struct {
	data []byte
	pos  int
}

func (br *bitReader) readBit() (uint, bool) {
	if br.pos >= len(br.data)*8 {
		return 0, false
	}
	bytePos := br.pos / 8
	bitPos := 7 - (br.pos % 8)
	br.pos++
	return uint((br.data[bytePos] >> bitPos) & 0x01), true
}

func (br *bitReader) readBits(count int) (uint, bool) {
	var value uint
	for i := 0; i < count; i++ {
		bit, ok := br.readBit()
		if !ok {
			return 0, false
		}
		value = (value << 1) | bit
	}
	return value, true
}

func (br *bitReader) readUE() (uint, bool) {
	zeros := 0
	for {
		bit, ok := br.readBit()
		if !ok {
			return 0, false
		}
		if bit == 0 {
			zeros++
			continue
		}
		break
	}
	if zeros == 0 {
		return 0, true
	}
	value, ok := br.readBits(zeros)
	if !ok {
		return 0, false
	}
	return (1<<zeros - 1) + value, true
}

func (br *bitReader) readSE() (int, bool) {
	value, ok := br.readUE()
	if !ok {
		return 0, false
	}
	codeNum := int(value)
	if codeNum%2 == 0 {
		return -(codeNum / 2), true
	}
	return (codeNum + 1) / 2, true
}
