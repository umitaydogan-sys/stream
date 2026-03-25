package main

import (
	"encoding/json"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/storage"
)

const playerTelemetryTTL = 45 * time.Second
const playerTelemetryPersistEvery = 15 * time.Second

type playerTelemetryPayload struct {
	StreamKey          string  `json:"stream_key"`
	SessionID          string  `json:"session_id"`
	Page               string  `json:"page"`
	PreferredFormat    string  `json:"preferred_format"`
	ActiveSourceKind   string  `json:"active_source_kind"`
	SourceOverride     string  `json:"source_override"`
	Quality            string  `json:"quality"`
	SelectedAudioTrack string  `json:"selected_audio_track"`
	SelectedAudioLabel string  `json:"selected_audio_label"`
	PlaybackSeconds    float64 `json:"playback_seconds"`
	BufferSeconds      float64 `json:"buffer_seconds"`
	StallCount         int     `json:"stall_count"`
	Recoveries         int     `json:"recoveries"`
	QualityTransitions int     `json:"quality_transitions"`
	AudioSwitches      int     `json:"audio_switches"`
	LastError          string  `json:"last_error"`
	Reconnect          string  `json:"reconnect"`
	Offline            bool    `json:"offline"`
	Waiting            bool    `json:"waiting"`
	DebugEnabled       bool    `json:"debug_enabled"`
}

type playerTelemetryCollector struct {
	mu      sync.Mutex
	db      *storage.SQLiteDB
	streams map[string]*playerTelemetryStream
}

type playerTelemetryStream struct {
	StreamKey               string
	Sessions                map[string]*playerTelemetrySession
	Reports                 int64
	TotalStalls             int64
	TotalRecoveries         int64
	TotalQualityTransitions int64
	TotalAudioSwitches      int64
	LastError               string
	LastUpdate              time.Time
	LastPersist             time.Time
}

type playerTelemetrySession struct {
	SessionID          string
	Page               string
	PreferredFormat    string
	ActiveSourceKind   string
	SourceOverride     string
	Quality            string
	SelectedAudioTrack string
	SelectedAudioLabel string
	PlaybackSeconds    float64
	BufferSeconds      float64
	StallCount         int
	Recoveries         int
	QualityTransitions int
	AudioSwitches      int
	LastError          string
	Reconnect          string
	Offline            bool
	Waiting            bool
	DebugEnabled       bool
	LastSeen           time.Time
	RemoteAddr         string
	UserAgent          string
}

type playerTelemetrySnapshot struct {
	StreamKey               string                           `json:"stream_key"`
	ActiveSessions          int                              `json:"active_sessions"`
	WaitingSessions         int                              `json:"waiting_sessions"`
	OfflineSessions         int                              `json:"offline_sessions"`
	DebugSessions           int                              `json:"debug_sessions"`
	Reports                 int64                            `json:"reports"`
	TotalStalls             int64                            `json:"total_stalls"`
	TotalRecoveries         int64                            `json:"total_recoveries"`
	TotalQualityTransitions int64                            `json:"total_quality_transitions"`
	TotalAudioSwitches      int64                            `json:"total_audio_switches"`
	AverageBufferSeconds    float64                          `json:"average_buffer_seconds"`
	AveragePlayback         float64                          `json:"average_playback_seconds"`
	LastError               string                           `json:"last_error"`
	LastUpdate              time.Time                        `json:"last_update"`
	Sources                 map[string]int                   `json:"sources"`
	Formats                 map[string]int                   `json:"formats"`
	Pages                   map[string]int                   `json:"pages"`
	Qualities               map[string]int                   `json:"qualities"`
	AudioTracks             map[string]int                   `json:"audio_tracks"`
	Sessions                []playerTelemetrySessionSnapshot `json:"sessions"`
}

type playerTelemetrySessionSnapshot struct {
	SessionID          string    `json:"session_id"`
	Page               string    `json:"page"`
	PreferredFormat    string    `json:"preferred_format"`
	ActiveSourceKind   string    `json:"active_source_kind"`
	SourceOverride     string    `json:"source_override"`
	Quality            string    `json:"quality"`
	SelectedAudioTrack string    `json:"selected_audio_track"`
	SelectedAudioLabel string    `json:"selected_audio_label"`
	PlaybackSeconds    float64   `json:"playback_seconds"`
	BufferSeconds      float64   `json:"buffer_seconds"`
	StallCount         int       `json:"stall_count"`
	Recoveries         int       `json:"recoveries"`
	QualityTransitions int       `json:"quality_transitions"`
	AudioSwitches      int       `json:"audio_switches"`
	LastError          string    `json:"last_error"`
	Reconnect          string    `json:"reconnect"`
	Offline            bool      `json:"offline"`
	Waiting            bool      `json:"waiting"`
	DebugEnabled       bool      `json:"debug_enabled"`
	LastSeen           time.Time `json:"last_seen"`
	LastSeenAgoSec     int       `json:"last_seen_ago_sec"`
	RemoteAddr         string    `json:"remote_addr"`
	UserAgent          string    `json:"user_agent"`
}

func newPlayerTelemetryCollector() *playerTelemetryCollector {
	return &playerTelemetryCollector{
		streams: make(map[string]*playerTelemetryStream),
	}
}

func (c *playerTelemetryCollector) SetDB(db *storage.SQLiteDB) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.db = db
}

func (c *playerTelemetryCollector) Record(payload playerTelemetryPayload, remoteAddr, userAgent string) {
	streamKey := strings.TrimSpace(payload.StreamKey)
	sessionID := strings.TrimSpace(payload.SessionID)
	if streamKey == "" || sessionID == "" {
		return
	}

	now := time.Now()
	payload = normalizeTelemetryPayload(payload)

	var persistSample *storage.PlayerTelemetrySample

	c.mu.Lock()

	c.cleanupLocked(now)

	state := c.streams[streamKey]
	if state == nil {
		state = &playerTelemetryStream{
			StreamKey: streamKey,
			Sessions:  make(map[string]*playerTelemetrySession),
		}
		c.streams[streamKey] = state
	}

	prev := state.Sessions[sessionID]
	if prev != nil {
		if delta := payload.StallCount - prev.StallCount; delta > 0 {
			state.TotalStalls += int64(delta)
		}
		if delta := payload.Recoveries - prev.Recoveries; delta > 0 {
			state.TotalRecoveries += int64(delta)
		}
		if delta := payload.QualityTransitions - prev.QualityTransitions; delta > 0 {
			state.TotalQualityTransitions += int64(delta)
		}
		if delta := payload.AudioSwitches - prev.AudioSwitches; delta > 0 {
			state.TotalAudioSwitches += int64(delta)
		}
	} else {
		state.TotalStalls += int64(maxInt(payload.StallCount, 0))
		state.TotalRecoveries += int64(maxInt(payload.Recoveries, 0))
		state.TotalQualityTransitions += int64(maxInt(payload.QualityTransitions, 0))
		state.TotalAudioSwitches += int64(maxInt(payload.AudioSwitches, 0))
	}

	state.Reports++
	state.LastUpdate = now
	if payload.LastError != "" && payload.LastError != "-" {
		state.LastError = payload.LastError
	}

	state.Sessions[sessionID] = &playerTelemetrySession{
		SessionID:          sessionID,
		Page:               payload.Page,
		PreferredFormat:    payload.PreferredFormat,
		ActiveSourceKind:   payload.ActiveSourceKind,
		SourceOverride:     payload.SourceOverride,
		Quality:            payload.Quality,
		SelectedAudioTrack: payload.SelectedAudioTrack,
		SelectedAudioLabel: payload.SelectedAudioLabel,
		PlaybackSeconds:    payload.PlaybackSeconds,
		BufferSeconds:      payload.BufferSeconds,
		StallCount:         payload.StallCount,
		Recoveries:         payload.Recoveries,
		QualityTransitions: payload.QualityTransitions,
		AudioSwitches:      payload.AudioSwitches,
		LastError:          payload.LastError,
		Reconnect:          payload.Reconnect,
		Offline:            payload.Offline,
		Waiting:            payload.Waiting,
		DebugEnabled:       payload.DebugEnabled,
		LastSeen:           now,
		RemoteAddr:         compactLabel(stripPort(remoteAddr), 80),
		UserAgent:          compactLabel(userAgent, 160),
	}

	if c.db != nil && (state.LastPersist.IsZero() || now.Sub(state.LastPersist) >= playerTelemetryPersistEvery) {
		persistSample = snapshotToTelemetrySample(c.snapshotForStateLocked(state, now))
		state.LastPersist = now
	}
	c.mu.Unlock()

	if c.db != nil && persistSample != nil {
		if err := c.db.SavePlayerTelemetrySample(persistSample); err != nil {
			// Persistence is best-effort; runtime telemetry should continue even if snapshot writes fail.
		}
	}
}

func (c *playerTelemetryCollector) Snapshot(streamKey string) playerTelemetrySnapshot {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cleanupLocked(now)

	state := c.streams[strings.TrimSpace(streamKey)]
	if state == nil {
		return playerTelemetrySnapshot{
			StreamKey:   strings.TrimSpace(streamKey),
			Sources:     map[string]int{},
			Formats:     map[string]int{},
			Pages:       map[string]int{},
			Qualities:   map[string]int{},
			AudioTracks: map[string]int{},
			Sessions:    []playerTelemetrySessionSnapshot{},
		}
	}

	return c.snapshotForStateLocked(state, now)
}

func (c *playerTelemetryCollector) snapshotForStateLocked(state *playerTelemetryStream, now time.Time) playerTelemetrySnapshot {
	snapshot := playerTelemetrySnapshot{
		StreamKey:               state.StreamKey,
		Reports:                 state.Reports,
		TotalStalls:             state.TotalStalls,
		TotalRecoveries:         state.TotalRecoveries,
		TotalQualityTransitions: state.TotalQualityTransitions,
		TotalAudioSwitches:      state.TotalAudioSwitches,
		LastError:               state.LastError,
		LastUpdate:              state.LastUpdate,
		Sources:                 make(map[string]int),
		Formats:                 make(map[string]int),
		Pages:                   make(map[string]int),
		Qualities:               make(map[string]int),
		AudioTracks:             make(map[string]int),
		Sessions:                make([]playerTelemetrySessionSnapshot, 0, len(state.Sessions)),
	}

	var totalBuffer float64
	var totalPlayback float64
	for _, session := range state.Sessions {
		if now.Sub(session.LastSeen) > playerTelemetryTTL {
			continue
		}
		snapshot.ActiveSessions++
		totalBuffer += session.BufferSeconds
		totalPlayback += session.PlaybackSeconds
		if session.Waiting {
			snapshot.WaitingSessions++
		}
		if session.Offline {
			snapshot.OfflineSessions++
		}
		if session.DebugEnabled {
			snapshot.DebugSessions++
		}
		snapshot.Sources[telemetryLabel(session.ActiveSourceKind, "-")]++
		snapshot.Formats[telemetryLabel(session.PreferredFormat, "auto")]++
		snapshot.Pages[telemetryLabel(session.Page, "player")]++
		snapshot.Qualities[telemetryLabel(session.Quality, "-")]++
		snapshot.AudioTracks[telemetryLabel(session.SelectedAudioLabel, telemetryLabel(session.SelectedAudioTrack, "-"))]++
		snapshot.Sessions = append(snapshot.Sessions, playerTelemetrySessionSnapshot{
			SessionID:          session.SessionID,
			Page:               telemetryLabel(session.Page, "player"),
			PreferredFormat:    telemetryLabel(session.PreferredFormat, "auto"),
			ActiveSourceKind:   telemetryLabel(session.ActiveSourceKind, "-"),
			SourceOverride:     telemetryLabel(session.SourceOverride, "auto"),
			Quality:            telemetryLabel(session.Quality, "-"),
			SelectedAudioTrack: telemetryLabel(session.SelectedAudioTrack, "-"),
			SelectedAudioLabel: telemetryLabel(session.SelectedAudioLabel, telemetryLabel(session.SelectedAudioTrack, "-")),
			PlaybackSeconds:    session.PlaybackSeconds,
			BufferSeconds:      session.BufferSeconds,
			StallCount:         session.StallCount,
			Recoveries:         session.Recoveries,
			QualityTransitions: session.QualityTransitions,
			AudioSwitches:      session.AudioSwitches,
			LastError:          telemetryLabel(session.LastError, "-"),
			Reconnect:          telemetryLabel(session.Reconnect, "-"),
			Offline:            session.Offline,
			Waiting:            session.Waiting,
			DebugEnabled:       session.DebugEnabled,
			LastSeen:           session.LastSeen,
			LastSeenAgoSec:     int(now.Sub(session.LastSeen).Seconds()),
			RemoteAddr:         telemetryLabel(session.RemoteAddr, "-"),
			UserAgent:          telemetryLabel(session.UserAgent, "-"),
		})
	}

	if snapshot.ActiveSessions > 0 {
		snapshot.AverageBufferSeconds = totalBuffer / float64(snapshot.ActiveSessions)
		snapshot.AveragePlayback = totalPlayback / float64(snapshot.ActiveSessions)
	}

	sort.Slice(snapshot.Sessions, func(i, j int) bool {
		return snapshot.Sessions[i].LastSeen.After(snapshot.Sessions[j].LastSeen)
	})
	if len(snapshot.Sessions) > 8 {
		snapshot.Sessions = snapshot.Sessions[:8]
	}

	return snapshot
}

func (c *playerTelemetryCollector) cleanupLocked(now time.Time) {
	for streamKey, state := range c.streams {
		for sessionID, session := range state.Sessions {
			if now.Sub(session.LastSeen) > playerTelemetryTTL {
				delete(state.Sessions, sessionID)
			}
		}
		if len(state.Sessions) == 0 && !state.LastUpdate.IsZero() && now.Sub(state.LastUpdate) > playerTelemetryTTL {
			delete(c.streams, streamKey)
		}
	}
}

func normalizeTelemetryPayload(payload playerTelemetryPayload) playerTelemetryPayload {
	payload.StreamKey = compactLabel(strings.TrimSpace(payload.StreamKey), 96)
	payload.SessionID = compactLabel(strings.TrimSpace(payload.SessionID), 96)
	payload.Page = compactLabel(strings.TrimSpace(payload.Page), 24)
	payload.PreferredFormat = compactLabel(strings.TrimSpace(payload.PreferredFormat), 24)
	payload.ActiveSourceKind = compactLabel(strings.TrimSpace(payload.ActiveSourceKind), 24)
	payload.SourceOverride = compactLabel(strings.TrimSpace(payload.SourceOverride), 24)
	payload.Quality = compactLabel(strings.TrimSpace(payload.Quality), 80)
	payload.SelectedAudioTrack = compactLabel(strings.TrimSpace(payload.SelectedAudioTrack), 32)
	payload.SelectedAudioLabel = compactLabel(strings.TrimSpace(payload.SelectedAudioLabel), 80)
	payload.LastError = compactLabel(strings.TrimSpace(payload.LastError), 120)
	payload.Reconnect = compactLabel(strings.TrimSpace(payload.Reconnect), 48)
	if payload.PlaybackSeconds < 0 {
		payload.PlaybackSeconds = 0
	}
	if payload.BufferSeconds < 0 {
		payload.BufferSeconds = 0
	}
	if payload.StallCount < 0 {
		payload.StallCount = 0
	}
	if payload.Recoveries < 0 {
		payload.Recoveries = 0
	}
	if payload.QualityTransitions < 0 {
		payload.QualityTransitions = 0
	}
	if payload.AudioSwitches < 0 {
		payload.AudioSwitches = 0
	}
	return payload
}

func telemetryLabel(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func compactLabel(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit]
}

func stripPort(addr string) string {
	if host, _, err := net.SplitHostPort(strings.TrimSpace(addr)); err == nil {
		return host
	}
	return strings.TrimSpace(addr)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func snapshotToTelemetrySample(snapshot playerTelemetrySnapshot) *storage.PlayerTelemetrySample {
	sourcesJSON, _ := json.Marshal(snapshot.Sources)
	formatsJSON, _ := json.Marshal(snapshot.Formats)
	pagesJSON, _ := json.Marshal(snapshot.Pages)
	qualitiesJSON, _ := json.Marshal(snapshot.Qualities)
	audioTracksJSON, _ := json.Marshal(snapshot.AudioTracks)
	return &storage.PlayerTelemetrySample{
		StreamKey:               snapshot.StreamKey,
		ActiveSessions:          snapshot.ActiveSessions,
		WaitingSessions:         snapshot.WaitingSessions,
		OfflineSessions:         snapshot.OfflineSessions,
		DebugSessions:           snapshot.DebugSessions,
		TotalStalls:             snapshot.TotalStalls,
		TotalRecoveries:         snapshot.TotalRecoveries,
		TotalQualityTransitions: snapshot.TotalQualityTransitions,
		TotalAudioSwitches:      snapshot.TotalAudioSwitches,
		AverageBufferSeconds:    snapshot.AverageBufferSeconds,
		AveragePlaybackSeconds:  snapshot.AveragePlayback,
		LastError:               snapshot.LastError,
		SourcesJSON:             string(sourcesJSON),
		FormatsJSON:             string(formatsJSON),
		PagesJSON:               string(pagesJSON),
		QualitiesJSON:           string(qualitiesJSON),
		AudioTracksJSON:         string(audioTracksJSON),
		CreatedAt:               snapshot.LastUpdate,
	}
}
