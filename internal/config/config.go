package config

import (
	"encoding/json"
	"strconv"

	"github.com/fluxstream/fluxstream/internal/storage"
)

// Manager handles configuration loaded from SQLite
type Manager struct {
	db *storage.SQLiteDB
}

// NewManager creates a new configuration manager
func NewManager(db *storage.SQLiteDB) *Manager {
	return &Manager{db: db}
}

// LoadDefaults inserts default configuration values if they don't exist
func (m *Manager) LoadDefaults() error {
	defaultABRProfiles, _ := json.Marshal(map[string][]map[string]interface{}{
		"balanced": {
			{"name": "1080p", "width": 1920, "height": 1080, "bitrate": "4500k", "max_bitrate": "5000k", "buf_size": "9000k", "preset": "fast", "fps": 30, "audio_rate": "192k"},
			{"name": "720p", "width": 1280, "height": 720, "bitrate": "2500k", "max_bitrate": "3000k", "buf_size": "5000k", "preset": "fast", "fps": 30, "audio_rate": "128k"},
			{"name": "480p", "width": 854, "height": 480, "bitrate": "1000k", "max_bitrate": "1200k", "buf_size": "2000k", "preset": "fast", "fps": 30, "audio_rate": "96k"},
			{"name": "360p", "width": 640, "height": 360, "bitrate": "600k", "max_bitrate": "700k", "buf_size": "1200k", "preset": "fast", "fps": 25, "audio_rate": "64k"},
		},
		"mobile": {
			{"name": "720p", "width": 1280, "height": 720, "bitrate": "1800k", "max_bitrate": "2200k", "buf_size": "3600k", "preset": "fast", "fps": 30, "audio_rate": "128k"},
			{"name": "480p", "width": 854, "height": 480, "bitrate": "900k", "max_bitrate": "1100k", "buf_size": "1800k", "preset": "fast", "fps": 30, "audio_rate": "96k"},
			{"name": "360p", "width": 640, "height": 360, "bitrate": "500k", "max_bitrate": "650k", "buf_size": "1000k", "preset": "fast", "fps": 25, "audio_rate": "64k"},
		},
		"radio": {
			{"name": "audio", "width": 0, "height": 0, "bitrate": "0", "max_bitrate": "0", "buf_size": "0", "preset": "fast", "fps": 0, "audio_rate": "128k"},
		},
	})

	defaults := map[string]struct {
		Value    string
		Category string
	}{
		// General
		"server_name":     {"FluxStream", "general"},
		"http_port":       {"8844", "general"},
		"https_port":      {"443", "general"},
		"language":        {"tr", "general"},
		"theme":           {"dark", "general"},
		"timezone":        {"Europe/Istanbul", "general"},
		"setup_completed": {"false", "general"},

		// RTMP
		"rtmp_enabled":    {"true", "protocols"},
		"rtmp_port":       {"1935", "protocols"},
		"rtmp_chunk_size": {"4096", "protocols"},
		"rtmp_gop_cache":  {"true", "protocols"},
		"rtmp_max_conns":  {"100", "protocols"},

		// RTMPS
		"rtmps_enabled": {"false", "protocols"},
		"rtmps_port":    {"1936", "protocols"},

		// SRT
		"srt_enabled": {"false", "protocols"},
		"srt_port":    {"9000", "protocols"},
		"srt_latency": {"120", "protocols"},

		// RTP
		"rtp_enabled": {"false", "protocols"},
		"rtp_port":    {"5004", "protocols"},

		// RTSP
		"rtsp_enabled": {"false", "protocols"},
		"rtsp_port":    {"8554", "protocols"},

		// WebRTC
		"webrtc_enabled": {"false", "protocols"},
		"webrtc_port":    {"8855", "protocols"},

		// MPEG-TS
		"mpegts_enabled": {"false", "protocols"},
		"mpegts_port":    {"9001", "protocols"},

		// HTTP Push
		"http_push_enabled": {"false", "protocols"},
		"http_push_port":    {"8850", "protocols"},
		"http_push_token":   {"", "protocols"},

		// HLS Output
		"hls_enabled":          {"true", "outputs"},
		"hls_segment_duration": {"2", "outputs"},
		"hls_playlist_length":  {"6", "outputs"},
		"hls_ll_enabled":       {"false", "outputs"},

		// DASH Output
		"dash_enabled":            {"false", "outputs"},
		"dash_segment_duration":   {"2", "outputs"},
		"abr_enabled":             {"false", "outputs"},
		"abr_profile_set":         {"balanced", "outputs"},
		"abr_master_enabled":      {"true", "outputs"},
		"abr_audio_passthrough":   {"false", "outputs"},
		"abr_profiles_json":       {string(defaultABRProfiles), "outputs"},
		"player_quality_selector": {"true", "outputs"},

		// HTTP-FLV Output
		"httpflv_enabled":   {"false", "outputs"},
		"httpflv_gop_cache": {"true", "outputs"},

		// WebRTC Output
		"whep_enabled": {"false", "outputs"},

		// RTMP Relay
		"relay_enabled": {"true", "outputs"},

		// RTSP Output
		"rtsp_out_enabled": {"false", "outputs"},
		"rtsp_out_port":    {"8555", "outputs"},

		// RTP Output
		"rtp_out_enabled": {"false", "outputs"},

		// SRT Output
		"srt_out_enabled": {"false", "outputs"},
		"srt_out_port":    {"9010", "outputs"},

		// MPEG-TS UDP Output
		"tsudp_out_enabled": {"false", "outputs"},

		// fMP4 Output
		"fmp4_enabled": {"true", "outputs"},

		// WebM Output
		"webm_enabled": {"true", "outputs"},

		// Audio Outputs
		"mp3_enabled":     {"false", "outputs"},
		"mp3_bitrate":     {"128", "outputs"},
		"aac_out_enabled": {"false", "outputs"},
		"icecast_enabled": {"false", "outputs"},
		"icecast_port":    {"8000", "outputs"},

		// SSL
		"ssl_enabled":          {"false", "ssl"},
		"ssl_mode":             {"file", "ssl"},
		"ssl_cert_path":        {"", "ssl"},
		"ssl_key_path":         {"", "ssl"},
		"ssl_le_domain":        {"", "ssl"},
		"ssl_le_email":         {"", "ssl"},
		"stream_ssl_mode":      {"file", "ssl"},
		"stream_ssl_cert_path": {"", "ssl"},
		"stream_ssl_key_path":  {"", "ssl"},
		"stream_ssl_le_domain": {"", "ssl"},
		"stream_ssl_le_email":  {"", "ssl"},

		// Security
		"stream_key_required": {"true", "security"},
		"token_enabled":       {"false", "security"},
		"token_duration":      {"60", "security"},
		"rate_limit":          {"100", "security"},

		// Storage
		"storage_max_gb":               {"50", "storage"},
		"storage_auto_clean":           {"30", "storage"},
		"recordings_retention_days":    {"30", "storage"},
		"recordings_keep_latest":       {"10", "storage"},
		"analytics_retention_days":     {"30", "storage"},
		"maintenance_auto_cleanup":     {"true", "storage"},
		"maintenance_cleanup_interval": {"6", "storage"},

		// Transcode
		"ffmpeg_path":                 {"ffmpeg", "transcode"},
		"gpu_accel":                   {"none", "transcode"},
		"transcode_live_hls_enabled":  {"true", "transcode"},
		"transcode_live_dash_enabled": {"true", "transcode"},
		"transcode_mode":              {"balanced", "transcode"},
		"transcode_cpu_limit":         {"0", "transcode"},

		// Recording
		"recording_enabled":   {"true", "recording"},
		"recording_format":    {"ts", "recording"},
		"recording_max_hours": {"24", "recording"},

		// Security extras
		"token_secret": {"", "security"},
		"twfa_enabled": {"false", "security"},

		// Embed Defaults
		"embed_domain":     {"localhost", "embed"},
		"embed_http_port":  {"8844", "embed"},
		"embed_https_port": {"443", "embed"},
		"embed_use_https":  {"false", "embed"},

		// Analytics / Health / Diagnostics
		"analytics_persist_enabled":       {"true", "analytics"},
		"analytics_snapshot_interval":     {"60", "analytics"},
		"track_analytics_enabled":         {"true", "analytics"},
		"track_analytics_interval":        {"20", "analytics"},
		"player_telemetry_retention_days": {"30", "analytics"},
		"track_analytics_retention_days":  {"30", "analytics"},
		"alerts_enabled":                  {"true", "health"},
		"alerts_disk_threshold_gb":        {"5", "health"},
		"alerts_memory_threshold_mb":      {"2048", "health"},
		"alerts_cert_days":                {"21", "health"},
		"alerts_qoe_stalls_threshold":     {"6", "health"},
		"alerts_qoe_buffer_seconds":       {"1", "health"},
		"alerts_qoe_waiting_sessions":     {"2", "health"},
		"alerts_qoe_offline_sessions":     {"1", "health"},
		"diagnostics_enabled":             {"true", "health"},
		"guided_mode_enabled":             {"true", "general"},
	}

	for key, d := range defaults {
		existing, _ := m.db.GetConfig(key)
		if existing == "" {
			if err := m.db.SetConfig(key, d.Value, d.Category); err != nil {
				return err
			}
		}
	}
	return nil
}

// Get returns a config value
func (m *Manager) Get(key, defaultVal string) string {
	val, err := m.db.GetConfig(key)
	if err != nil || val == "" {
		return defaultVal
	}
	return val
}

// Set sets a config value
func (m *Manager) Set(key, value, category string) error {
	return m.db.SetConfig(key, value, category)
}

// GetInt returns a config value as int
func (m *Manager) GetInt(key string, defaultVal int) int {
	val := m.Get(key, "")
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

// GetBool returns a config value as bool
func (m *Manager) GetBool(key string, defaultVal bool) bool {
	val := m.Get(key, "")
	if val == "" {
		return defaultVal
	}
	return val == "true" || val == "1"
}

// GetAll returns all config as a map
func (m *Manager) GetAll() (map[string]string, error) {
	return m.db.GetAllConfig()
}

// GetByCategory returns configs for a specific category
func (m *Manager) GetByCategory(category string) ([]storage.Config, error) {
	return m.db.GetConfigByCategory(category)
}
