package storage

import "time"

// Config represents a key-value configuration entry stored in SQLite
type Config struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Category  string    `json:"category"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents an admin/operator/viewer user
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // admin, operator, viewer
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLogin    time.Time `json:"last_login,omitempty"`
}

// Stream represents a configured stream
type Stream struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	StreamKey     string    `json:"stream_key"`
	Status        string    `json:"status"` // offline, live, waiting
	IngestProto   string    `json:"ingest_proto,omitempty"`
	OutputFormats string    `json:"output_formats"` // JSON array: ["hls","dash","flv"]
	PolicyJSON    string    `json:"policy_json,omitempty"`
	MaxViewers    int       `json:"max_viewers,omitempty"`
	MaxBitrate    int       `json:"max_bitrate,omitempty"`
	RecordEnabled bool      `json:"record_enabled"`
	RecordFormat  string    `json:"record_format,omitempty"` // mp4, ts, mkv, flv
	Password      string    `json:"password,omitempty"`
	DomainLock    string    `json:"domain_lock,omitempty"`
	IPWhitelist   string    `json:"ip_whitelist,omitempty"`
	ThumbnailPath string    `json:"thumbnail_path,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	StartedAt     time.Time `json:"started_at,omitempty"`
	ViewerCount   int       `json:"viewer_count"`
	InputBitrate  int64     `json:"input_bitrate"`
	InputCodec    string    `json:"input_codec,omitempty"`
	InputWidth    int       `json:"input_width,omitempty"`
	InputHeight   int       `json:"input_height,omitempty"`
	InputFPS      float64   `json:"input_fps,omitempty"`
}

// Viewer represents an active viewer connection
type Viewer struct {
	ID        int64     `json:"id"`
	StreamID  int64     `json:"stream_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Protocol  string    `json:"protocol"` // hls, dash, flv, webrtc, rtsp, etc.
	Country   string    `json:"country,omitempty"`
	City      string    `json:"city,omitempty"`
	StartedAt time.Time `json:"started_at"`
	Bandwidth int64     `json:"bandwidth"`
}

// BannedIP represents a banned IP address
type BannedIP struct {
	ID       int64     `json:"id"`
	IP       string    `json:"ip"`
	Reason   string    `json:"reason,omitempty"`
	BannedAt time.Time `json:"banned_at"`
}

// LogEntry represents a log entry
type LogEntry struct {
	ID        int64     `json:"id"`
	Level     string    `json:"level"` // INFO, WARN, ERROR, DEBUG
	Component string    `json:"component"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// PlayerTemplate represents a player appearance profile
type PlayerTemplate struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	BackgroundCSS string    `json:"background_css"`
	ControlBarCSS string    `json:"control_bar_css"`
	PlayButtonCSS string    `json:"play_button_css"`
	LogoURL       string    `json:"logo_url,omitempty"`
	LogoPosition  string    `json:"logo_position,omitempty"`
	LogoOpacity   float64   `json:"logo_opacity"`
	WatermarkText string    `json:"watermark_text,omitempty"`
	ShowTitle     bool      `json:"show_title"`
	ShowLiveBadge bool      `json:"show_live_badge"`
	Theme         string    `json:"theme"` // dark, light, custom
	CustomCSS     string    `json:"custom_css,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// EmbedDefaults represents default embed generation settings
type EmbedDefaults struct {
	Domain      string `json:"domain"`
	HTTPPort    int    `json:"http_port"`
	HTTPSPort   int    `json:"https_port"`
	RTMPPort    int    `json:"rtmp_port"`
	RTSPPort    int    `json:"rtsp_port"`
	SRTPort     int    `json:"srt_port"`
	WebRTCPort  int    `json:"webrtc_port"`
	IcecastPort int    `json:"icecast_port"`
	UseHTTPS    bool   `json:"use_https"`
}

// ServerStats represents current server statistics
type ServerStats struct {
	ActiveStreams int     `json:"active_streams"`
	TotalViewers  int     `json:"total_viewers"`
	BandwidthIn   int64   `json:"bandwidth_in"`
	BandwidthOut  int64   `json:"bandwidth_out"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsedMB  int64   `json:"memory_used_mb"`
	MemoryTotalMB int64   `json:"memory_total_mb"`
	UptimeSeconds int64   `json:"uptime_seconds"`
}

// AnalyticsSnapshot stores a persisted dashboard snapshot.
type AnalyticsSnapshot struct {
	ID               int64     `json:"id"`
	Timestamp        time.Time `json:"timestamp"`
	TotalStreams     int       `json:"total_streams"`
	TotalViewers     int64     `json:"total_viewers"`
	CurrentViewers   int       `json:"current_viewers"`
	PeakConcurrent   int       `json:"peak_concurrent"`
	TotalBandwidth   int64     `json:"total_bandwidth"`
	ViewersByFormat  string    `json:"viewers_by_format"`
	ViewersByCountry string    `json:"viewers_by_country"`
}

// PlayerTelemetrySample stores a persisted QoE snapshot for a stream.
type PlayerTelemetrySample struct {
	ID                      int64     `json:"id"`
	StreamKey               string    `json:"stream_key"`
	ActiveSessions          int       `json:"active_sessions"`
	WaitingSessions         int       `json:"waiting_sessions"`
	OfflineSessions         int       `json:"offline_sessions"`
	DebugSessions           int       `json:"debug_sessions"`
	TotalStalls             int64     `json:"total_stalls"`
	TotalRecoveries         int64     `json:"total_recoveries"`
	TotalQualityTransitions int64     `json:"total_quality_transitions"`
	TotalAudioSwitches      int64     `json:"total_audio_switches"`
	AverageBufferSeconds    float64   `json:"average_buffer_seconds"`
	AveragePlaybackSeconds  float64   `json:"average_playback_seconds"`
	LastError               string    `json:"last_error"`
	SourcesJSON             string    `json:"sources_json"`
	FormatsJSON             string    `json:"formats_json"`
	PagesJSON               string    `json:"pages_json"`
	QualitiesJSON           string    `json:"qualities_json"`
	AudioTracksJSON         string    `json:"audio_tracks_json"`
	CreatedAt               time.Time `json:"created_at"`
}

// TrackTelemetrySample stores a persisted bitrate and runtime snapshot for a live track.
type TrackTelemetrySample struct {
	ID           int64     `json:"id"`
	StreamKey    string    `json:"stream_key"`
	TrackID      int       `json:"track_id"`
	Kind         string    `json:"kind"`
	Codec        string    `json:"codec"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	SampleRate   int       `json:"sample_rate,omitempty"`
	Channels     int       `json:"channels,omitempty"`
	Bitrate      int64     `json:"bitrate"`
	Packets      int64     `json:"packets"`
	Bytes        int64     `json:"bytes"`
	IsDefault    bool      `json:"is_default"`
	IsActive     bool      `json:"is_active"`
	DisplayLabel string    `json:"display_label"`
	CreatedAt    time.Time `json:"created_at"`
}

// RecordingArchive stores archive/object-storage metadata for a recording file.
type RecordingArchive struct {
	ID           int64     `json:"id"`
	StreamKey    string    `json:"stream_key"`
	Filename     string    `json:"filename"`
	Format       string    `json:"format"`
	Provider     string    `json:"provider"`
	Bucket       string    `json:"bucket,omitempty"`
	Endpoint     string    `json:"endpoint,omitempty"`
	ObjectKey    string    `json:"object_key"`
	ObjectURL    string    `json:"object_url,omitempty"`
	ETag         string    `json:"etag,omitempty"`
	Size         int64     `json:"size"`
	Status       string    `json:"status"`
	LastError    string    `json:"last_error,omitempty"`
	LocalDeleted bool      `json:"local_deleted"`
	ArchivedAt   time.Time `json:"archived_at,omitempty"`
	RestoredAt   time.Time `json:"restored_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BackupArchive stores archive/object-storage metadata for a system backup file.
type BackupArchive struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Provider          string    `json:"provider"`
	Bucket            string    `json:"bucket,omitempty"`
	Endpoint          string    `json:"endpoint,omitempty"`
	ObjectKey         string    `json:"object_key"`
	ObjectURL         string    `json:"object_url,omitempty"`
	ETag              string    `json:"etag,omitempty"`
	Size              int64     `json:"size"`
	IncludeRecordings bool      `json:"include_recordings"`
	Status            string    `json:"status"`
	LastError         string    `json:"last_error,omitempty"`
	LocalDeleted      bool      `json:"local_deleted"`
	ArchivedAt        time.Time `json:"archived_at,omitempty"`
	RestoredAt        time.Time `json:"restored_at,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// EmbedProfile stores reusable embed/link generation presets per stream.
type EmbedProfile struct {
	ID            int64     `json:"id"`
	StreamKey     string    `json:"stream_key"`
	Name          string    `json:"name"`
	UseCase       string    `json:"use_case"`
	Mode          string    `json:"mode"`
	PrimaryFormat string    `json:"primary_format"`
	Width         int       `json:"width"`
	Height        int       `json:"height"`
	Theme         string    `json:"theme"`
	OptionsJSON   string    `json:"options_json"`
	BrandingJSON  string    `json:"branding_json"`
	SecurityJSON  string    `json:"security_json"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ABRProfile stores saved adaptive bitrate profile libraries.
type ABRProfile struct {
	ID          int64     `json:"id"`
	ProfileSet  string    `json:"profile_set"`
	Name        string    `json:"name"`
	Scope       string    `json:"scope"` // global or stream
	StreamKey   string    `json:"stream_key,omitempty"`
	Description string    `json:"description,omitempty"`
	Preset      string    `json:"preset,omitempty"`
	ProfilesJSON string   `json:"profiles_json"`
	SummaryJSON string    `json:"summary_json,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
