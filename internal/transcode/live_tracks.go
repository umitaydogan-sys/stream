package transcode

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
)

type LiveTrackSnapshot struct {
	StreamKey           string          `json:"stream_key"`
	DirectMode          bool            `json:"direct_mode"`
	DefaultVideoTrackID int             `json:"default_video_track_id"`
	DefaultAudioTrackID int             `json:"default_audio_track_id"`
	ActiveVideoTrackID  int             `json:"active_video_track_id"`
	ActiveAudioTrackID  int             `json:"active_audio_track_id"`
	UpdatedAt           time.Time       `json:"updated_at"`
	VideoTracks         []LiveTrackInfo `json:"video_tracks"`
	AudioTracks         []LiveTrackInfo `json:"audio_tracks"`
}

type LiveTrackInfo struct {
	TrackID      int       `json:"track_id"`
	Kind         string    `json:"kind"`
	Codec        string    `json:"codec"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	SampleRate   int       `json:"sample_rate,omitempty"`
	Channels     int       `json:"channels,omitempty"`
	Bitrate      int64     `json:"bitrate,omitempty"`
	Packets      int64     `json:"packets"`
	Bytes        int64     `json:"bytes"`
	LastSeen     time.Time `json:"last_seen"`
	LastSeenAgo  int       `json:"last_seen_ago_sec"`
	Enhanced     bool      `json:"enhanced"`
	IsDefault    bool      `json:"is_default"`
	IsActive     bool      `json:"is_active"`
	PlaylistPath string    `json:"playlist_path,omitempty"`
	DisplayLabel string    `json:"display_label"`
}

type liveTrackRegistry struct {
	streamKey           string
	video               map[uint8]*liveTrackState
	audio               map[uint8]*liveTrackState
	defaultVideoTrackID uint8
	defaultAudioTrackID uint8
	updatedAt           time.Time
}

type liveTrackState struct {
	trackID       uint8
	kind          string
	codec         string
	width         int
	height        int
	sampleRate    int
	channels      int
	bitrate       int64
	packets       int64
	bytes         int64
	lastSeen      time.Time
	enhanced      bool
	playlistPath  string
	windowStarted time.Time
	windowBytes   int64
}

func newLiveTrackRegistry(streamKey string) *liveTrackRegistry {
	return &liveTrackRegistry{
		streamKey: streamKey,
		video:     make(map[uint8]*liveTrackState),
		audio:     make(map[uint8]*liveTrackState),
	}
}

func (m *Manager) ensureTrackRegistryLocked(streamKey string) *liveTrackRegistry {
	if m.liveTracks == nil {
		m.liveTracks = make(map[string]*liveTrackRegistry)
	}
	reg := m.liveTracks[streamKey]
	if reg == nil {
		reg = newLiveTrackRegistry(streamKey)
		m.liveTracks[streamKey] = reg
	}
	return reg
}

func (m *Manager) observeTrackPacket(streamKey string, pkt *media.Packet) {
	if pkt == nil || strings.TrimSpace(streamKey) == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	reg := m.ensureTrackRegistryLocked(streamKey)
	reg.observe(pkt, time.Now())
}

func (m *Manager) SetStreamTrackDefaults(streamKey string, videoTrackID, audioTrackID uint8) {
	streamKey = strings.TrimSpace(streamKey)
	if streamKey == "" {
		return
	}

	var direct *liveDirectSession

	m.mu.Lock()
	reg := m.ensureTrackRegistryLocked(streamKey)
	reg.defaultVideoTrackID = videoTrackID
	reg.defaultAudioTrackID = audioTrackID
	if opts, ok := m.streamOpts[streamKey]; ok {
		opts.DefaultVideoTrackID = videoTrackID
		opts.DefaultAudioTrackID = audioTrackID
		m.streamOpts[streamKey] = opts
	}
	direct = m.liveDirect[streamKey]
	m.mu.Unlock()

	if direct != nil {
		if videoTrackID != 0 {
			direct.setPreferredVideoTrack(videoTrackID)
		} else {
			direct.setPreferredVideoTrack(0)
		}
		if audioTrackID != 0 {
			direct.setSelectedAudioTrack(audioTrackID)
		} else {
			direct.setSelectedAudioTrack(0)
		}
	}
}

func (m *Manager) GetLiveTrackSnapshot(streamKey string) LiveTrackSnapshot {
	now := time.Now()

	m.mu.RLock()
	reg := m.liveTracks[strings.TrimSpace(streamKey)]
	direct := m.liveDirect[strings.TrimSpace(streamKey)]
	m.mu.RUnlock()

	snapshot := LiveTrackSnapshot{
		StreamKey:   strings.TrimSpace(streamKey),
		VideoTracks: []LiveTrackInfo{},
		AudioTracks: []LiveTrackInfo{},
	}
	if reg == nil {
		if direct != nil {
			snapshot.DirectMode = true
			snapshot.ActiveVideoTrackID = int(direct.getPrimaryTrackID())
			snapshot.ActiveAudioTrackID = int(direct.getSelectedAudioTrackID())
			snapshot.DefaultVideoTrackID = int(direct.getPreferredVideoTrackID())
			snapshot.DefaultAudioTrackID = int(direct.getPreferredAudioTrackID())
			snapshot.UpdatedAt = now
		}
		return snapshot
	}

	snapshot.DefaultVideoTrackID = int(reg.defaultVideoTrackID)
	snapshot.DefaultAudioTrackID = int(reg.defaultAudioTrackID)
	snapshot.UpdatedAt = reg.updatedAt
	if direct != nil {
		snapshot.DirectMode = true
		snapshot.ActiveVideoTrackID = int(direct.getPrimaryTrackID())
		snapshot.ActiveAudioTrackID = int(direct.getSelectedAudioTrackID())
		if snapshot.DefaultVideoTrackID == 0 {
			snapshot.DefaultVideoTrackID = int(direct.getPreferredVideoTrackID())
		}
		if snapshot.DefaultAudioTrackID == 0 {
			snapshot.DefaultAudioTrackID = int(direct.getPreferredAudioTrackID())
		}
	}

	videoIDs := make([]int, 0, len(reg.video))
	for trackID := range reg.video {
		videoIDs = append(videoIDs, int(trackID))
	}
	sort.Ints(videoIDs)
	for _, rawTrackID := range videoIDs {
		state := reg.video[uint8(rawTrackID)]
		if state == nil {
			continue
		}
		info := state.snapshot(now)
		info.IsDefault = snapshot.DefaultVideoTrackID > 0 && info.TrackID == snapshot.DefaultVideoTrackID
		info.IsActive = snapshot.ActiveVideoTrackID > 0 && info.TrackID == snapshot.ActiveVideoTrackID
		snapshot.VideoTracks = append(snapshot.VideoTracks, info)
	}

	audioIDs := make([]int, 0, len(reg.audio))
	for trackID := range reg.audio {
		audioIDs = append(audioIDs, int(trackID))
	}
	sort.Ints(audioIDs)
	for _, rawTrackID := range audioIDs {
		state := reg.audio[uint8(rawTrackID)]
		if state == nil {
			continue
		}
		info := state.snapshot(now)
		info.IsDefault = snapshot.DefaultAudioTrackID > 0 && info.TrackID == snapshot.DefaultAudioTrackID
		info.IsActive = snapshot.ActiveAudioTrackID > 0 && info.TrackID == snapshot.ActiveAudioTrackID
		snapshot.AudioTracks = append(snapshot.AudioTracks, info)
	}

	return snapshot
}

func (m *Manager) ResolveLiveAudioPlaylistPath(streamKey string, requestedTrackID uint8) string {
	streamKey = strings.TrimSpace(streamKey)
	if streamKey == "" {
		return ""
	}
	m.mu.RLock()
	direct := m.liveDirect[streamKey]
	m.mu.RUnlock()
	if direct == nil {
		return "audio.m3u8"
	}
	return direct.resolveAudioPlaylistPath(requestedTrackID)
}

func (r *liveTrackRegistry) observe(pkt *media.Packet, now time.Time) {
	if pkt == nil {
		return
	}
	r.updatedAt = now
	trackID := pkt.TrackID
	var bucket map[uint8]*liveTrackState
	var kind string
	switch pkt.Type {
	case media.PacketTypeVideo:
		bucket = r.video
		kind = "video"
	case media.PacketTypeAudio:
		bucket = r.audio
		kind = "audio"
	default:
		return
	}
	state := bucket[trackID]
	if state == nil {
		state = &liveTrackState{trackID: trackID, kind: kind}
		bucket[trackID] = state
	}
	state.observe(pkt, now)
}

func (s *liveTrackState) observe(pkt *media.Packet, now time.Time) {
	s.lastSeen = now
	s.packets++
	s.bytes += int64(len(pkt.Data))
	s.enhanced = s.enhanced || pkt.IsEnhanced

	if s.windowStarted.IsZero() {
		s.windowStarted = now
	}
	s.windowBytes += int64(len(pkt.Data))
	if elapsed := now.Sub(s.windowStarted); elapsed >= 2*time.Second {
		s.bitrate = int64(float64(s.windowBytes*8) / elapsed.Seconds())
		s.windowStarted = now
		s.windowBytes = 0
	}

	switch pkt.Type {
	case media.PacketTypeVideo:
		if pkt.FourCC == "avc1" || s.codec == "" {
			s.codec = "H.264"
		}
		if pkt.IsSequenceHeader {
			width, height := parseAVCSequenceHeaderDimensions(pkt.Data)
			if width > 0 {
				s.width = width
			}
			if height > 0 {
				s.height = height
			}
		}
	case media.PacketTypeAudio:
		if pkt.FourCC == "mp4a" || s.codec == "" {
			s.codec = "AAC"
		}
		if pkt.IsSequenceHeader {
			sampleRate, channels := parseAACSequenceHeader(pkt.Data)
			if sampleRate > 0 {
				s.sampleRate = sampleRate
			}
			if channels > 0 {
				s.channels = channels
			}
		}
	}
}

func (s *liveTrackState) snapshot(now time.Time) LiveTrackInfo {
	label := s.codec
	if s.kind == "video" {
		switch {
		case s.height > 0:
			label = label + " " + itoaSafe(s.height) + "p"
		case s.width > 0 && s.height > 0:
			label = label + " " + itoaSafe(s.width) + "x" + itoaSafe(s.height)
		}
	} else if s.sampleRate > 0 {
		label = label + " " + itoaSafe(s.sampleRate) + " Hz"
	}
	label = strings.TrimSpace(label)
	if label == "" {
		label = strings.Title(s.kind) + " Track " + itoaSafe(int(s.trackID))
	}
	return LiveTrackInfo{
		TrackID:      int(s.trackID),
		Kind:         s.kind,
		Codec:        s.codec,
		Width:        s.width,
		Height:       s.height,
		SampleRate:   s.sampleRate,
		Channels:     s.channels,
		Bitrate:      s.bitrate,
		Packets:      s.packets,
		Bytes:        s.bytes,
		LastSeen:     s.lastSeen,
		LastSeenAgo:  int(now.Sub(s.lastSeen).Seconds()),
		Enhanced:     s.enhanced,
		PlaylistPath: s.playlistPath,
		DisplayLabel: label,
	}
}

func parseAACSequenceHeader(data []byte) (int, int) {
	if len(data) < 4 {
		return 0, 0
	}
	asc := data[2:]
	if len(asc) < 2 {
		return 0, 0
	}
	freqIdx := int(((asc[0] & 0x07) << 1) | ((asc[1] >> 7) & 0x01))
	chCfg := int((asc[1] >> 3) & 0x0F)
	frequencies := []int{96000, 88200, 64000, 48000, 44100, 32000, 24000, 22050, 16000, 12000, 11025, 8000, 7350}
	sampleRate := 0
	if freqIdx >= 0 && freqIdx < len(frequencies) {
		sampleRate = frequencies[freqIdx]
	}
	return sampleRate, chCfg
}

func itoaSafe(value int) string {
	return strconv.Itoa(value)
}
