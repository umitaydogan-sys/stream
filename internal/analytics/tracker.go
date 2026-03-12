package analytics

import (
	"encoding/json"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
)

// Event represents an analytics event.
type Event struct {
	Type      string    `json:"type"`
	StreamKey string    `json:"stream_key"`
	ViewerID  string    `json:"viewer_id,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Country   string    `json:"country,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	Format    string    `json:"format,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Duration  float64   `json:"duration,omitempty"`
	BytesSent int64     `json:"bytes_sent,omitempty"`
}

// StreamStats represents analytics for a single stream.
type StreamStats struct {
	StreamKey      string         `json:"stream_key"`
	TotalViewers   int64          `json:"total_viewers"`
	PeakViewers    int64          `json:"peak_viewers"`
	CurrentViewers int            `json:"current_viewers"`
	TotalDuration  float64        `json:"total_duration"`
	TotalBytes     int64          `json:"total_bytes"`
	ByFormat       map[string]int `json:"by_format"`
	ByCountry      map[string]int `json:"by_country"`
	HourlyViewers  map[int]int    `json:"hourly_viewers"`
	LastUpdated    time.Time      `json:"last_updated"`
}

// Dashboard represents the analytics dashboard data.
type Dashboard struct {
	TotalStreams      int             `json:"total_streams"`
	TotalViewers      int64           `json:"total_viewers"`
	CurrentViewers    int             `json:"current_viewers"`
	PeakConcurrent    int             `json:"peak_concurrent"`
	TotalBandwidth    int64           `json:"total_bandwidth"`
	TopStreams        []StreamRank    `json:"top_streams"`
	ViewersByFormat   map[string]int  `json:"viewers_by_format"`
	ViewersByCountry  map[string]int  `json:"viewers_by_country"`
	ViewersTimeline   []TimelinePoint `json:"viewers_timeline"`
	BandwidthTimeline []TimelinePoint `json:"bandwidth_timeline"`
}

// StreamRank represents a stream in the top-streams list.
type StreamRank struct {
	StreamKey  string `json:"stream_key"`
	StreamName string `json:"stream_name"`
	Viewers    int    `json:"viewers"`
	Duration   string `json:"duration"`
}

// TimelinePoint represents a data point in a timeline.
type TimelinePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int64     `json:"value"`
}

// ViewerSession represents an active viewer session snapshot.
type ViewerSession struct {
	StreamKey       string    `json:"stream_key"`
	StreamName      string    `json:"stream_name"`
	ViewerID        string    `json:"viewer_id"`
	IP              string    `json:"ip"`
	Country         string    `json:"country"`
	UserAgent       string    `json:"user_agent"`
	Format          string    `json:"format"`
	StartedAt       time.Time `json:"started_at"`
	LastSeen        time.Time `json:"last_seen"`
	BytesSent       int64     `json:"bytes_sent"`
	DurationSeconds int64     `json:"duration_seconds"`
}

type viewerSessionInternal struct {
	StreamKey string
	ViewerID  string
	IP        string
	Country   string
	UserAgent string
	Format    string
	StartedAt time.Time
	LastSeen  time.Time
	BytesSent int64
}

// Tracker tracks analytics events and generates stats.
type Tracker struct {
	events         []Event
	streamStats    map[string]*StreamStats
	streamNames    map[string]string
	sessions       map[string]*viewerSessionInternal
	peakConcurrent int
	sessionTTL     time.Duration
	mu             sync.RWMutex
}

// NewTracker creates a new analytics tracker.
func NewTracker() *Tracker {
	return &Tracker{
		streamStats: make(map[string]*StreamStats),
		streamNames: make(map[string]string),
		sessions:    make(map[string]*viewerSessionInternal),
		sessionTTL:  20 * time.Second,
	}
}

// RegisterStreamName stores a friendly name for a stream.
func (t *Tracker) RegisterStreamName(streamKey, streamName string) {
	streamKey = strings.TrimSpace(streamKey)
	streamName = strings.TrimSpace(streamName)
	if streamKey == "" || streamName == "" {
		return
	}
	t.mu.Lock()
	t.streamNames[streamKey] = streamName
	t.mu.Unlock()
}

// TrackPlayback updates or creates a viewer session from a playback heartbeat.
func (t *Tracker) TrackPlayback(streamKey, viewerID, format, ip, country, userAgent string, bytesSent int64) {
	streamKey = strings.TrimSpace(streamKey)
	viewerID = strings.TrimSpace(viewerID)
	if streamKey == "" || viewerID == "" {
		return
	}

	now := time.Now()

	t.mu.Lock()
	defer t.mu.Unlock()

	t.pruneExpiredLocked(now)

	stats := t.ensureStreamStatsLocked(streamKey)
	sKey := sessionKey(streamKey, viewerID)
	sess, exists := t.sessions[sKey]
	if !exists {
		country = normalizeCountry(country, ip)
		sess = &viewerSessionInternal{
			StreamKey: streamKey,
			ViewerID:  viewerID,
			IP:        ip,
			Country:   country,
			UserAgent: userAgent,
			Format:    strings.TrimSpace(format),
			StartedAt: now,
			LastSeen:  now,
		}
		t.sessions[sKey] = sess

		stats.TotalViewers++
		stats.CurrentViewers++
		stats.HourlyViewers[now.Hour()]++
		if sess.Format != "" {
			stats.ByFormat[sess.Format]++
		}
		if sess.Country != "" {
			stats.ByCountry[sess.Country]++
		}
		if int64(stats.CurrentViewers) > stats.PeakViewers {
			stats.PeakViewers = int64(stats.CurrentViewers)
		}
		if current := len(t.sessions); current > t.peakConcurrent {
			t.peakConcurrent = current
		}
		t.appendEventLocked(Event{
			Type:      "viewer_join",
			StreamKey: streamKey,
			ViewerID:  viewerID,
			IP:        ip,
			Country:   sess.Country,
			UserAgent: userAgent,
			Format:    sess.Format,
			Timestamp: now,
		})
	} else {
		sess.LastSeen = now
		if sess.IP == "" {
			sess.IP = ip
		}
		if sess.UserAgent == "" {
			sess.UserAgent = userAgent
		}
		if sess.Country == "" {
			sess.Country = normalizeCountry(country, ip)
		}
		if sess.Format == "" && strings.TrimSpace(format) != "" {
			sess.Format = strings.TrimSpace(format)
		}
	}

	if bytesSent > 0 {
		sess.BytesSent += bytesSent
		stats.TotalBytes += bytesSent
	}
	stats.LastUpdated = now
}

// EndPlayback closes an active session explicitly.
func (t *Tracker) EndPlayback(streamKey, viewerID string, bytesSent int64) {
	streamKey = strings.TrimSpace(streamKey)
	viewerID = strings.TrimSpace(viewerID)
	if streamKey == "" || viewerID == "" {
		return
	}

	now := time.Now()

	t.mu.Lock()
	defer t.mu.Unlock()

	if sess, ok := t.sessions[sessionKey(streamKey, viewerID)]; ok {
		t.endSessionLocked(sess, now, bytesSent)
	}
}

// TrackEvent preserves compatibility with older callers.
func (t *Tracker) TrackEvent(evt Event) {
	switch evt.Type {
	case "viewer_join":
		t.TrackPlayback(evt.StreamKey, evt.ViewerID, evt.Format, evt.IP, evt.Country, evt.UserAgent, evt.BytesSent)
	case "viewer_leave":
		t.EndPlayback(evt.StreamKey, evt.ViewerID, evt.BytesSent)
	default:
		t.mu.Lock()
		t.appendEventLocked(evt)
		t.mu.Unlock()
	}
}

// GetDashboard returns the analytics dashboard.
func (t *Tracker) GetDashboard() Dashboard {
	now := time.Now()

	t.mu.Lock()
	t.pruneExpiredLocked(now)

	dash := Dashboard{
		TotalStreams:     len(t.streamStats),
		PeakConcurrent:   t.peakConcurrent,
		ViewersByFormat:  make(map[string]int),
		ViewersByCountry: make(map[string]int),
	}

	for _, stats := range t.streamStats {
		dash.TotalViewers += stats.TotalViewers
		dash.CurrentViewers += stats.CurrentViewers
		dash.TotalBandwidth += stats.TotalBytes
		for f, c := range stats.ByFormat {
			dash.ViewersByFormat[f] += c
		}
		for co, c := range stats.ByCountry {
			dash.ViewersByCountry[co] += c
		}
		if stats.CurrentViewers > 0 || stats.TotalViewers > 0 {
			dash.TopStreams = append(dash.TopStreams, StreamRank{
				StreamKey:  stats.StreamKey,
				StreamName: t.streamNames[stats.StreamKey],
				Viewers:    stats.CurrentViewers,
				Duration:   time.Since(stats.LastUpdated).Round(time.Second).String(),
			})
		}
	}

	sort.Slice(dash.TopStreams, func(i, j int) bool {
		if dash.TopStreams[i].Viewers == dash.TopStreams[j].Viewers {
			return dash.TopStreams[i].StreamKey < dash.TopStreams[j].StreamKey
		}
		return dash.TopStreams[i].Viewers > dash.TopStreams[j].Viewers
	})

	for i := 23; i >= 0; i-- {
		ts := now.Add(-time.Duration(i) * time.Hour).Truncate(time.Hour)
		viewers := int64(0)
		bandwidth := int64(0)
		for _, e := range t.events {
			if e.Timestamp.Before(ts) || !e.Timestamp.Before(ts.Add(time.Hour)) {
				continue
			}
			if e.Type == "viewer_join" {
				viewers++
			}
			bandwidth += e.BytesSent
		}
		dash.ViewersTimeline = append(dash.ViewersTimeline, TimelinePoint{Timestamp: ts, Value: viewers})
		dash.BandwidthTimeline = append(dash.BandwidthTimeline, TimelinePoint{Timestamp: ts, Value: bandwidth})
	}

	t.mu.Unlock()

	return dash
}

// GetStreamStats returns stats for a specific stream.
func (t *Tracker) GetStreamStats(streamKey string) *StreamStats {
	now := time.Now()

	t.mu.Lock()
	defer t.mu.Unlock()

	t.pruneExpiredLocked(now)
	stats := t.streamStats[streamKey]
	if stats == nil {
		return nil
	}
	copyStats := *stats
	copyStats.ByFormat = cloneMap(stats.ByFormat)
	copyStats.ByCountry = cloneMap(stats.ByCountry)
	copyStats.HourlyViewers = cloneMap(stats.HourlyViewers)
	return &copyStats
}

// GetViewerSessions returns active viewer sessions.
func (t *Tracker) GetViewerSessions() []ViewerSession {
	now := time.Now()

	t.mu.Lock()
	defer t.mu.Unlock()

	t.pruneExpiredLocked(now)

	out := make([]ViewerSession, 0, len(t.sessions))
	for _, sess := range t.sessions {
		out = append(out, ViewerSession{
			StreamKey:       sess.StreamKey,
			StreamName:      t.streamNames[sess.StreamKey],
			ViewerID:        sess.ViewerID,
			IP:              sess.IP,
			Country:         sess.Country,
			UserAgent:       sess.UserAgent,
			Format:          sess.Format,
			StartedAt:       sess.StartedAt,
			LastSeen:        sess.LastSeen,
			BytesSent:       sess.BytesSent,
			DurationSeconds: int64(now.Sub(sess.StartedAt).Seconds()),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].StreamKey == out[j].StreamKey {
			return out[i].StartedAt.Before(out[j].StartedAt)
		}
		return out[i].StreamKey < out[j].StreamKey
	})

	return out
}

// CurrentViewersByStream returns current active viewer counts per stream.
func (t *Tracker) CurrentViewersByStream() map[string]int {
	now := time.Now()

	t.mu.Lock()
	defer t.mu.Unlock()

	t.pruneExpiredLocked(now)
	out := make(map[string]int, len(t.streamStats))
	for key, stats := range t.streamStats {
		out[key] = stats.CurrentViewers
	}
	return out
}

// MarshalDashboard returns JSON dashboard data.
func (t *Tracker) MarshalDashboard() ([]byte, error) {
	return json.Marshal(t.GetDashboard())
}

func (t *Tracker) ensureStreamStatsLocked(streamKey string) *StreamStats {
	stats, ok := t.streamStats[streamKey]
	if !ok {
		stats = &StreamStats{
			StreamKey:     streamKey,
			ByFormat:      make(map[string]int),
			ByCountry:     make(map[string]int),
			HourlyViewers: make(map[int]int),
		}
		t.streamStats[streamKey] = stats
	}
	return stats
}

func (t *Tracker) pruneExpiredLocked(now time.Time) {
	cutoff := now.Add(-t.sessionTTL)
	for _, sess := range t.sessions {
		if sess.LastSeen.Before(cutoff) {
			t.endSessionLocked(sess, now, 0)
		}
	}

	eventCutoff := now.Add(-24 * time.Hour)
	keep := t.events[:0]
	for _, evt := range t.events {
		if evt.Timestamp.After(eventCutoff) {
			keep = append(keep, evt)
		}
	}
	t.events = keep
}

func (t *Tracker) endSessionLocked(sess *viewerSessionInternal, now time.Time, extraBytes int64) {
	if sess == nil {
		return
	}

	stats := t.ensureStreamStatsLocked(sess.StreamKey)
	if extraBytes > 0 {
		sess.BytesSent += extraBytes
		stats.TotalBytes += extraBytes
	}
	if stats.CurrentViewers > 0 {
		stats.CurrentViewers--
	}
	stats.TotalDuration += now.Sub(sess.StartedAt).Seconds()
	stats.LastUpdated = now

	delete(t.sessions, sessionKey(sess.StreamKey, sess.ViewerID))
	t.appendEventLocked(Event{
		Type:      "viewer_leave",
		StreamKey: sess.StreamKey,
		ViewerID:  sess.ViewerID,
		IP:        sess.IP,
		Country:   sess.Country,
		UserAgent: sess.UserAgent,
		Format:    sess.Format,
		Timestamp: now,
		Duration:  now.Sub(sess.StartedAt).Seconds(),
		BytesSent: sess.BytesSent,
	})
}

func (t *Tracker) appendEventLocked(evt Event) {
	if evt.Timestamp.IsZero() {
		evt.Timestamp = time.Now()
	}
	t.events = append(t.events, evt)
}

func sessionKey(streamKey, viewerID string) string {
	return streamKey + "|" + viewerID
}

func normalizeCountry(country, ip string) string {
	country = strings.TrimSpace(country)
	if country != "" {
		return country
	}
	host := strings.TrimSpace(ip)
	if host == "" {
		return "Bilinmiyor"
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	addr := net.ParseIP(host)
	if addr == nil {
		return "Bilinmiyor"
	}
	if addr.IsLoopback() || addr.IsPrivate() {
		return "Yerel Ag"
	}
	return "Bilinmiyor"
}

func cloneMap[T comparable](in map[T]int) map[T]int {
	out := make(map[T]int, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
