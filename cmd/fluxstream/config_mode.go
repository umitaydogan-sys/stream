package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/storage"
)

func handleConfigMode(args []string) (bool, error) {
	if len(args) < 1 || !strings.EqualFold(args[0], "config") {
		return false, nil
	}
	if len(args) < 3 || !strings.EqualFold(args[1], "set") {
		return true, fmt.Errorf("kullanim: fluxstream.exe config set key=value ...")
	}

	execPath, err := os.Executable()
	if err != nil {
		return true, fmt.Errorf("executable path alinamadi: %w", err)
	}

	dataDir := filepath.Join(filepath.Dir(execPath), "data")
	if err := ensureDataDirs(dataDir); err != nil {
		return true, err
	}

	db, err := storage.NewSQLiteDB(filepath.Join(dataDir, "fluxstream.db"))
	if err != nil {
		return true, err
	}
	defer db.Close()

	cfg := config.NewManager(db)
	if err := cfg.LoadDefaults(); err != nil {
		return true, err
	}

	for _, arg := range args[2:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return true, fmt.Errorf("gecersiz arguman: %s", arg)
		}
		key := strings.TrimSpace(parts[0])
		value := parts[1]
		if key == "" {
			return true, fmt.Errorf("bos config key: %s", arg)
		}
		if err := cfg.Set(key, value, configCategoryForKey(key)); err != nil {
			return true, err
		}
	}

	fmt.Println("Config values applied.")
	return true, nil
}

func configCategoryForKey(key string) string {
	switch key {
	case "server_name", "http_port", "https_port", "language", "theme", "timezone", "setup_completed", "guided_mode_enabled":
		return "general"
	case "embed_domain", "embed_http_port", "embed_https_port", "embed_use_https":
		return "embed"
	case "rtmp_enabled", "rtmp_port", "rtmps_enabled", "rtmps_port", "srt_enabled", "srt_port", "rtsp_enabled", "rtsp_port", "webrtc_enabled", "webrtc_port", "mpegts_enabled", "mpegts_port", "http_push_enabled", "http_push_port", "http_push_token", "rtp_enabled", "rtp_port":
		return "protocols"
	case "ssl_enabled", "ssl_mode", "ssl_cert_path", "ssl_key_path", "ssl_le_domain", "ssl_le_email",
		"stream_ssl_mode", "stream_ssl_cert_path", "stream_ssl_key_path", "stream_ssl_le_domain", "stream_ssl_le_email":
		return "ssl"
	case "hls_enabled", "hls_segment_duration", "hls_playlist_length", "hls_ll_enabled", "dash_enabled", "dash_segment_duration", "httpflv_enabled", "httpflv_gop_cache", "whep_enabled", "relay_enabled", "rtsp_out_enabled", "rtsp_out_port", "rtp_out_enabled", "srt_out_enabled", "srt_out_port", "tsudp_out_enabled", "fmp4_enabled", "webm_enabled", "mp3_enabled", "mp3_bitrate", "aac_out_enabled", "icecast_enabled", "icecast_port", "abr_enabled", "abr_profile_set", "abr_master_enabled", "abr_audio_passthrough", "abr_profiles_json", "player_quality_selector":
		return "outputs"
	case "ffmpeg_path", "gpu_accel", "transcode_live_hls_enabled", "transcode_live_dash_enabled", "transcode_mode", "transcode_cpu_limit":
		return "transcode"
	case "stream_key_required", "token_enabled", "token_duration", "rate_limit", "token_secret", "twfa_enabled":
		return "security"
	case "storage_max_gb", "storage_auto_clean", "recordings_retention_days", "recordings_keep_latest",
		"archive_enabled", "archive_ui_mode", "archive_provider", "archive_provider_variant", "archive_local_dir",
		"archive_endpoint", "archive_region", "archive_bucket", "archive_access_key", "archive_secret_key",
		"archive_rclone_remote", "archive_rclone_path", "archive_rclone_config_path", "archive_prefix",
		"archive_use_path_style", "archive_public_base_url", "archive_auto_upload", "archive_schedule",
		"archive_target_tier", "archive_cold_after_days", "archive_delete_local_after_upload",
		"archive_scan_interval_minutes", "archive_batch_size", "archive_sftp_host", "archive_sftp_port",
		"archive_sftp_user", "archive_sftp_remote_dir", "archive_sftp_key_path",
		"archive_sftp_disable_host_key_check", "backup_archive_use_same_target", "backup_archive_enabled",
		"backup_archive_provider", "backup_archive_provider_variant", "backup_archive_local_dir",
		"backup_archive_endpoint", "backup_archive_region", "backup_archive_bucket", "backup_archive_access_key",
		"backup_archive_secret_key", "backup_archive_rclone_remote", "backup_archive_rclone_path",
		"backup_archive_rclone_config_path", "backup_archive_sftp_host", "backup_archive_sftp_port",
		"backup_archive_sftp_user", "backup_archive_sftp_remote_dir", "backup_archive_sftp_key_path",
		"backup_archive_sftp_disable_host_key_check", "backup_archive_prefix", "backup_archive_use_path_style",
		"backup_archive_public_base_url", "backup_archive_auto_upload", "backup_archive_schedule",
		"backup_archive_target_tier", "backup_archive_cold_after_days",
		"backup_archive_delete_local_after_upload", "backup_archive_scan_interval_minutes",
		"backup_archive_batch_size",
		"analytics_retention_days", "maintenance_auto_cleanup", "maintenance_cleanup_interval":
		return "storage"
	case "analytics_persist_enabled", "analytics_snapshot_interval":
		return "analytics"
	case "alerts_enabled", "alerts_disk_threshold_gb", "alerts_memory_threshold_mb", "alerts_cert_days",
		"alerts_qoe_stalls_threshold", "alerts_qoe_buffer_seconds", "alerts_qoe_buffer_warn_seconds",
		"alerts_qoe_buffer_critical_seconds", "alerts_qoe_waiting_sessions", "alerts_qoe_waiting_ratio_percent",
		"alerts_qoe_offline_sessions", "alerts_qoe_offline_ratio_percent", "alerts_qoe_transition_ratio_threshold",
		"alerts_qoe_audio_ratio_threshold", "diagnostics_enabled":
		return "health"
	case "recording_enabled", "recording_format", "recording_max_hours":
		return "recording"
	default:
		return "general"
	}
}

func ensureDataDirs(dataDir string) error {
	dirs := []string{
		dataDir,
		filepath.Join(dataDir, "hls"),
		filepath.Join(dataDir, "dash"),
		filepath.Join(dataDir, "recordings"),
		filepath.Join(dataDir, "backups"),
		filepath.Join(dataDir, "thumbnails"),
		filepath.Join(dataDir, "license"),
		filepath.Join(dataDir, "certs"),
		filepath.Join(dataDir, "certs", "web"),
		filepath.Join(dataDir, "certs", "stream"),
		filepath.Join(dataDir, "certs", "acme"),
		filepath.Join(dataDir, "logs"),
		filepath.Join(dataDir, "players"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("dizin olusturulamadi %s: %w", d, err)
		}
	}
	return nil
}
