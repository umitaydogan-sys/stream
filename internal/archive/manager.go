package archive

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/recording"
	"github.com/fluxstream/fluxstream/internal/storage"
)

type TargetSettings struct {
	Enabled                 bool
	Configured              bool
	Provider                string
	ProviderVariant         string
	DisplayName             string
	LocalDir                string
	Endpoint                string
	Region                  string
	Bucket                  string
	AccessKey               string
	SecretKey               string
	RcloneRemote            string
	RclonePath              string
	RcloneConfigPath        string
	SFTPHost                string
	SFTPPort                int
	SFTPUser                string
	SFTPRemoteDir           string
	SFTPKeyPath             string
	SFTPDisableHostKeyCheck bool
	Prefix                  string
	UsePathStyle            bool
	PublicBaseURL           string
	AutoUpload              bool
	DeleteLocal             bool
	ScanMinutes             int
	BatchSize               int
	Schedule                string
	TargetTier              string
	ColdAfterDays           int
}

type Settings struct {
	Enabled                 bool
	Configured              bool
	RecordingsEnabled       bool
	BackupsEnabled          bool
	Mode                    string
	BackupUseSameTarget     bool
	Recording               TargetSettings
	Backup                  TargetSettings

	// Compatibility mirrors for the primary recording target.
	Provider          string
	ProviderVariant   string
	LocalDir          string
	Endpoint          string
	Region            string
	Bucket            string
	AccessKey         string
	SecretKey         string
	SFTPHost          string
	SFTPPort          int
	SFTPUser          string
	SFTPRemoteDir     string
	SFTPKeyPath       string
	Prefix            string
	UsePathStyle      bool
	PublicBaseURL     string
	AutoUpload        bool
	DeleteLocal       bool
	ScanMinutes       int
	BatchSize         int
	BackupAutoUpload  bool
	BackupDeleteLocal bool
	BackupScanMinutes int
	BackupBatchSize   int
}

type Summary struct {
	Enabled                 bool      `json:"enabled"`
	Configured              bool      `json:"configured"`
	RecordingsEnabled       bool      `json:"recordings_enabled"`
	BackupsEnabled          bool      `json:"backups_enabled"`
	Mode                    string    `json:"mode,omitempty"`
	BackupUseSameTarget     bool      `json:"backup_use_same_target"`
	Provider                string    `json:"provider"`
	ProviderVariant         string    `json:"provider_variant,omitempty"`
	AutoUpload              bool      `json:"auto_upload"`
	DeleteLocal             bool      `json:"delete_local_after_upload"`
	BackupAutoUpload        bool      `json:"backup_auto_upload"`
	BackupDeleteLocal       bool      `json:"backup_delete_local_after_upload"`
	RecordingConfigured     bool      `json:"recording_configured"`
	BackupConfigured        bool      `json:"backup_configured"`
	RecordingProvider       string    `json:"recording_provider"`
	RecordingProviderVariant string   `json:"recording_provider_variant,omitempty"`
	BackupProvider          string    `json:"backup_provider"`
	BackupProviderVariant   string    `json:"backup_provider_variant,omitempty"`
	RecordingSchedule       string    `json:"recording_schedule,omitempty"`
	BackupSchedule          string    `json:"backup_schedule,omitempty"`
	RecordingTargetTier     string    `json:"recording_target_tier,omitempty"`
	BackupTargetTier        string    `json:"backup_target_tier,omitempty"`
	LocalDir                string    `json:"local_dir,omitempty"`
	Endpoint                string    `json:"endpoint,omitempty"`
	Bucket                  string    `json:"bucket,omitempty"`
	Prefix                  string    `json:"prefix,omitempty"`
	BackupLocalDir          string    `json:"backup_local_dir,omitempty"`
	BackupEndpoint          string    `json:"backup_endpoint,omitempty"`
	BackupBucket            string    `json:"backup_bucket,omitempty"`
	BackupPrefix            string    `json:"backup_prefix,omitempty"`
	Items                   int       `json:"items"`
	ErrorItems              int       `json:"error_items"`
	LocalDeletedItems       int       `json:"local_deleted_items"`
	BackupItems             int       `json:"backup_items"`
	BackupErrorItems        int       `json:"backup_error_items"`
	BackupLocalDeletedItems int       `json:"backup_local_deleted_items"`
	LastSyncAt              time.Time `json:"last_sync_at,omitempty"`
	LastError               string    `json:"last_error,omitempty"`
	LastBackupSyncAt        time.Time `json:"last_backup_sync_at,omitempty"`
	LastBackupError         string    `json:"last_backup_error,omitempty"`
}

type Manager struct {
	cfg        *config.Manager
	db         *storage.SQLiteDB
	recordings *recording.Manager
	dataDir    string

	mu               sync.RWMutex
	lastSyncAt       time.Time
	lastError        string
	lastBackupSyncAt time.Time
	lastBackupError  string
}

type storeObject struct {
	Key  string
	URL  string
	ETag string
	Size int64
}

type storeClient interface {
	UploadFile(ctx context.Context, objectKey, localPath, contentType string) (storeObject, error)
	DownloadFile(ctx context.Context, objectKey, destPath string) (int64, error)
	DeleteObject(ctx context.Context, objectKey string) error
}

func NewManager(cfg *config.Manager, db *storage.SQLiteDB, recordings *recording.Manager, dataDir string) *Manager {
	return &Manager{
		cfg:        cfg,
		db:         db,
		recordings: recordings,
		dataDir:    dataDir,
	}
}

func (m *Manager) Settings() Settings {
	return m.settingsWithOverrides(nil)
}

func (m *Manager) settingsWithOverrides(overrides map[string]string) Settings {
	get := func(key, defaultVal string) string {
		if overrides != nil {
			if value, ok := overrides[key]; ok {
				return value
			}
		}
		return m.cfg.Get(key, defaultVal)
	}
	settings := Settings{
		Mode:                normalizeUIMode(get("archive_ui_mode", "simple")),
		RecordingsEnabled:   parseBoolValue(get("archive_enabled", "false"), false),
		BackupsEnabled:      parseBoolValue(get("backup_archive_enabled", "false"), false),
		BackupUseSameTarget: parseBoolValue(get("backup_archive_use_same_target", "true"), true),
	}
	recordingTarget := buildTargetSettings("archive", filepath.Join(m.dataDir, "archive"), get)
	recordingTarget.Enabled = settings.RecordingsEnabled
	backupTarget := buildTargetSettings("backup_archive", filepath.Join(m.dataDir, "archive-backups"), get)
	backupTarget.Enabled = settings.BackupsEnabled
	if settings.BackupUseSameTarget {
		backupTarget = cloneTargetForBackup(recordingTarget, get)
		backupTarget.Enabled = settings.BackupsEnabled
	}
	settings.Recording = recordingTarget
	settings.Backup = backupTarget
	settings.Provider = recordingTarget.Provider
	settings.ProviderVariant = recordingTarget.ProviderVariant
	settings.LocalDir = recordingTarget.LocalDir
	settings.Endpoint = recordingTarget.Endpoint
	settings.Region = recordingTarget.Region
	settings.Bucket = recordingTarget.Bucket
	settings.AccessKey = recordingTarget.AccessKey
	settings.SecretKey = recordingTarget.SecretKey
	settings.SFTPHost = recordingTarget.SFTPHost
	settings.SFTPPort = recordingTarget.SFTPPort
	settings.SFTPUser = recordingTarget.SFTPUser
	settings.SFTPRemoteDir = recordingTarget.SFTPRemoteDir
	settings.SFTPKeyPath = recordingTarget.SFTPKeyPath
	settings.Prefix = recordingTarget.Prefix
	settings.UsePathStyle = recordingTarget.UsePathStyle
	settings.PublicBaseURL = recordingTarget.PublicBaseURL
	settings.AutoUpload = recordingTarget.AutoUpload
	settings.DeleteLocal = recordingTarget.DeleteLocal
	settings.ScanMinutes = recordingTarget.ScanMinutes
	settings.BatchSize = recordingTarget.BatchSize
	settings.BackupAutoUpload = backupTarget.AutoUpload
	settings.BackupDeleteLocal = backupTarget.DeleteLocal
	settings.BackupScanMinutes = backupTarget.ScanMinutes
	settings.BackupBatchSize = backupTarget.BatchSize
	settings.Configured = (!settings.RecordingsEnabled || recordingTarget.Configured) && (!settings.BackupsEnabled || backupTarget.Configured)
	settings.Enabled = (settings.RecordingsEnabled && recordingTarget.Configured) || (settings.BackupsEnabled && backupTarget.Configured)
	return settings
}

func (m *Manager) Enabled() bool {
	settings := m.Settings()
	return settings.Recording.Configured && settings.RecordingsEnabled
}

func (m *Manager) Summary() Summary {
	settings := m.Settings()
	items, _ := m.db.ListRecordingArchives("", 0)
	backupItems, _ := m.db.ListBackupArchives(0)
	summary := Summary{
		Enabled:                  settings.Enabled,
		Configured:               settings.Configured,
		RecordingsEnabled:        settings.RecordingsEnabled,
		BackupsEnabled:           settings.BackupsEnabled,
		Mode:                     settings.Mode,
		BackupUseSameTarget:      settings.BackupUseSameTarget,
		Provider:                 providerLabel(settings.Recording.Provider, settings.Recording.ProviderVariant),
		ProviderVariant:          settings.Recording.ProviderVariant,
		AutoUpload:               settings.Recording.AutoUpload,
		DeleteLocal:              settings.Recording.DeleteLocal,
		BackupAutoUpload:         settings.Backup.AutoUpload,
		BackupDeleteLocal:        settings.Backup.DeleteLocal,
		RecordingConfigured:      settings.Recording.Configured,
		BackupConfigured:         settings.Backup.Configured,
		RecordingProvider:        providerLabel(settings.Recording.Provider, settings.Recording.ProviderVariant),
		RecordingProviderVariant: settings.Recording.ProviderVariant,
		BackupProvider:           providerLabel(settings.Backup.Provider, settings.Backup.ProviderVariant),
		BackupProviderVariant:    settings.Backup.ProviderVariant,
		RecordingSchedule:        settings.Recording.Schedule,
		BackupSchedule:           settings.Backup.Schedule,
		RecordingTargetTier:      settings.Recording.TargetTier,
		BackupTargetTier:         settings.Backup.TargetTier,
		LocalDir:                 settings.Recording.LocalDir,
		Endpoint:                 settings.Recording.Endpoint,
		Bucket:                   settings.Recording.Bucket,
		Prefix:                   settings.Recording.Prefix,
		BackupLocalDir:           settings.Backup.LocalDir,
		BackupEndpoint:           settings.Backup.Endpoint,
		BackupBucket:             settings.Backup.Bucket,
		BackupPrefix:             settings.Backup.Prefix,
		Items:                    len(items),
		BackupItems:              len(backupItems),
	}
	for _, item := range items {
		if strings.EqualFold(item.Status, "error") {
			summary.ErrorItems++
		}
		if item.LocalDeleted {
			summary.LocalDeletedItems++
		}
	}
	for _, item := range backupItems {
		if strings.EqualFold(item.Status, "error") {
			summary.BackupErrorItems++
		}
		if item.LocalDeleted {
			summary.BackupLocalDeletedItems++
		}
	}
	m.mu.RLock()
	summary.LastSyncAt = m.lastSyncAt
	summary.LastError = m.lastError
	summary.LastBackupSyncAt = m.lastBackupSyncAt
	summary.LastBackupError = m.lastBackupError
	m.mu.RUnlock()
	return summary
}

func parseBoolValue(raw string, defaultVal bool) bool {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return defaultVal
	}
	switch raw {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultVal
	}
}

func parseIntValue(raw string, defaultVal int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return defaultVal
	}
	return n
}

func normalizeUIMode(raw string) string {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "advanced", "gelismis":
		return "advanced"
	default:
		return "simple"
	}
}

func normalizeSchedule(raw string, defaultVal string) string {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "manual", "immediate", "hourly", "daily", "weekly":
		return strings.TrimSpace(strings.ToLower(raw))
	default:
		return defaultVal
	}
}

func providerFromVariant(provider, variant string) string {
	provider = strings.TrimSpace(strings.ToLower(provider))
	variant = strings.TrimSpace(strings.ToLower(variant))
	if provider != "" && provider != "disabled" {
		return provider
	}
	switch variant {
	case "local":
		return "local"
	case "aws_s3", "cloudflare_r2", "backblaze_b2", "wasabi", "digitalocean_spaces", "linode_object_storage", "scaleway_object_storage", "idrive_e2", "ceph_rgw":
		return "s3"
	case "minio":
		return "minio"
	case "sftp":
		return "sftp"
	case "google_drive", "onedrive", "dropbox", "google_cloud_storage", "azure_blob", "box", "pcloud", "mega", "webdav", "nextcloud":
		return "rclone"
	default:
		return "disabled"
	}
}

func defaultVariantForProvider(provider string) string {
	switch strings.TrimSpace(strings.ToLower(provider)) {
	case "local":
		return "local"
	case "s3":
		return "aws_s3"
	case "minio":
		return "minio"
	case "sftp":
		return "sftp"
	case "rclone":
		return "google_drive"
	default:
		return "local"
	}
}

func buildTargetSettings(prefix string, defaultLocalDir string, get func(string, string) string) TargetSettings {
	provider := providerFromVariant(get(prefix+"_provider", ""), get(prefix+"_provider_variant", ""))
	variant := strings.TrimSpace(strings.ToLower(get(prefix+"_provider_variant", "")))
	if variant == "" {
		variant = defaultVariantForProvider(provider)
	}
	target := TargetSettings{
		Provider:                provider,
		ProviderVariant:         variant,
		LocalDir:                strings.TrimSpace(get(prefix+"_local_dir", defaultLocalDir)),
		Endpoint:                strings.TrimRight(strings.TrimSpace(get(prefix+"_endpoint", "")), "/"),
		Region:                  strings.TrimSpace(get(prefix+"_region", "us-east-1")),
		Bucket:                  strings.TrimSpace(get(prefix+"_bucket", "")),
		AccessKey:               strings.TrimSpace(get(prefix+"_access_key", "")),
		SecretKey:               strings.TrimSpace(get(prefix+"_secret_key", "")),
		RcloneRemote:            strings.TrimSpace(get(prefix+"_rclone_remote", "")),
		RclonePath:              strings.Trim(strings.TrimSpace(get(prefix+"_rclone_path", "")), "/"),
		RcloneConfigPath:        strings.TrimSpace(get(prefix+"_rclone_config_path", "")),
		SFTPHost:                strings.TrimSpace(get(prefix+"_sftp_host", "")),
		SFTPPort:                parseIntValue(get(prefix+"_sftp_port", "22"), 22),
		SFTPUser:                strings.TrimSpace(get(prefix+"_sftp_user", "")),
		SFTPRemoteDir:           strings.TrimSpace(get(prefix+"_sftp_remote_dir", "")),
		SFTPKeyPath:             strings.TrimSpace(get(prefix+"_sftp_key_path", "")),
		SFTPDisableHostKeyCheck: parseBoolValue(get(prefix+"_sftp_disable_host_key_check", "false"), false),
		Prefix:                  strings.Trim(strings.TrimSpace(get(prefix+"_prefix", "fluxstream")), "/"),
		UsePathStyle:            parseBoolValue(get(prefix+"_use_path_style", "true"), true),
		PublicBaseURL:           strings.TrimRight(strings.TrimSpace(get(prefix+"_public_base_url", "")), "/"),
		AutoUpload:              parseBoolValue(get(prefix+"_auto_upload", "false"), false),
		DeleteLocal:             parseBoolValue(get(prefix+"_delete_local_after_upload", "false"), false),
		ScanMinutes:             parseIntValue(get(prefix+"_scan_interval_minutes", "10"), 10),
		BatchSize:               parseIntValue(get(prefix+"_batch_size", "3"), 3),
		Schedule:                normalizeSchedule(get(prefix+"_schedule", "manual"), "manual"),
		TargetTier:              strings.TrimSpace(strings.ToLower(get(prefix+"_target_tier", "standard"))),
		ColdAfterDays:           parseIntValue(get(prefix+"_cold_after_days", "30"), 30),
	}
	if target.SFTPPort <= 0 {
		target.SFTPPort = 22
	}
	if target.ScanMinutes <= 0 {
		target.ScanMinutes = 10
	}
	if target.BatchSize <= 0 {
		target.BatchSize = 3
	}
	if target.Provider == "local" && target.LocalDir == "" {
		target.LocalDir = defaultLocalDir
	}
	target.DisplayName = providerLabel(target.Provider, target.ProviderVariant)
	switch target.Provider {
	case "local":
		target.Configured = strings.TrimSpace(target.LocalDir) != ""
	case "s3", "minio":
		target.Configured = target.Endpoint != "" && target.Bucket != "" && target.AccessKey != "" && target.SecretKey != ""
	case "sftp":
		target.Configured = target.SFTPHost != "" && target.SFTPUser != "" && target.SFTPRemoteDir != ""
	case "rclone":
		target.Configured = target.RcloneRemote != ""
	default:
		target.Configured = false
	}
	return target
}

func cloneTargetForBackup(source TargetSettings, get func(string, string) string) TargetSettings {
	clone := source
	clone.Provider = providerFromVariant(source.Provider, source.ProviderVariant)
	clone.ProviderVariant = strings.TrimSpace(strings.ToLower(get("backup_archive_provider_variant", source.ProviderVariant)))
	if clone.ProviderVariant == "" {
		clone.ProviderVariant = source.ProviderVariant
	}
	if override := strings.TrimSpace(get("backup_archive_prefix", clone.Prefix)); override != "" {
		clone.Prefix = strings.Trim(override, "/")
	}
	if override := strings.TrimSpace(get("backup_archive_public_base_url", clone.PublicBaseURL)); override != "" {
		clone.PublicBaseURL = strings.TrimRight(override, "/")
	}
	clone.AutoUpload = parseBoolValue(get("backup_archive_auto_upload", "false"), false)
	clone.DeleteLocal = parseBoolValue(get("backup_archive_delete_local_after_upload", "false"), false)
	clone.ScanMinutes = parseIntValue(get("backup_archive_scan_interval_minutes", "30"), 30)
	if clone.ScanMinutes <= 0 {
		clone.ScanMinutes = 30
	}
	clone.BatchSize = parseIntValue(get("backup_archive_batch_size", "2"), 2)
	if clone.BatchSize <= 0 {
		clone.BatchSize = 2
	}
	clone.Schedule = normalizeSchedule(get("backup_archive_schedule", "weekly"), "weekly")
	clone.TargetTier = strings.TrimSpace(strings.ToLower(get("backup_archive_target_tier", "cold")))
	clone.ColdAfterDays = parseIntValue(get("backup_archive_cold_after_days", "7"), 7)
	return clone
}

func (m *Manager) targetSettingsFor(role string, overrides map[string]string) (TargetSettings, error) {
	settings := m.settingsWithOverrides(overrides)
	switch strings.TrimSpace(strings.ToLower(role)) {
	case "recordings", "recording", "archive":
		return settings.Recording, nil
	case "backups", "backup":
		return settings.Backup, nil
	default:
		return TargetSettings{}, fmt.Errorf("gecersiz hedef rolu")
	}
}

func (m *Manager) TestConnection(ctx context.Context, role string, overrides map[string]string) (map[string]interface{}, error) {
	target, err := m.targetSettingsFor(role, overrides)
	if err != nil {
		return nil, err
	}
	if !target.Configured {
		return nil, fmt.Errorf("hedef ayarlari eksik")
	}
	client, err := m.newClient(target)
	if err != nil {
		return nil, err
	}
	tmpFile, err := os.CreateTemp("", "fluxstream-archive-probe-*.txt")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())
	payload := []byte("fluxstream archive connection probe\n")
	if _, err := tmpFile.Write(payload); err != nil {
		_ = tmpFile.Close()
		return nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}
	objectKey := path.Join(strings.Trim(target.Prefix, "/"), "_probe", fmt.Sprintf("%s-%d.txt", sanitizeObjectPart(role), time.Now().UnixNano()))
	obj, err := client.UploadFile(ctx, objectKey, tmpFile.Name(), "text/plain")
	if err != nil {
		return nil, err
	}
	if err := client.DeleteObject(ctx, objectKey); err != nil {
		return nil, fmt.Errorf("baglanti denemesi upload oldu ama test dosyasi silinemedi: %w", err)
	}
	return map[string]interface{}{
		"provider": providerLabel(target.Provider, target.ProviderVariant),
		"engine":   target.Provider,
		"variant":  target.ProviderVariant,
		"object":   obj.Key,
		"url":      obj.URL,
		"size":     obj.Size,
	}, nil
}

func (m *Manager) ShouldRunRecordingSchedule(now time.Time) bool {
	settings := m.Settings()
	return scheduleDue(settings.Recording.Schedule, settings.Recording.ScanMinutes, m.lastRecordingSyncAt(), now)
}

func (m *Manager) ShouldRunBackupSchedule(now time.Time) bool {
	settings := m.Settings()
	return scheduleDue(settings.Backup.Schedule, settings.Backup.ScanMinutes, m.lastBackupSyncTimestamp(), now)
}

func (m *Manager) lastRecordingSyncAt() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastSyncAt
}

func (m *Manager) lastBackupSyncTimestamp() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastBackupSyncAt
}

func scheduleDue(schedule string, scanMinutes int, lastRun, now time.Time) bool {
	switch normalizeSchedule(schedule, "manual") {
	case "manual":
		return false
	case "immediate":
		return true
	case "hourly":
		if lastRun.IsZero() {
			return true
		}
		return now.Sub(lastRun) >= time.Hour
	case "daily":
		if lastRun.IsZero() {
			return true
		}
		return now.Sub(lastRun) >= 24*time.Hour
	case "weekly":
		if lastRun.IsZero() {
			return true
		}
		return now.Sub(lastRun) >= 7*24*time.Hour
	default:
		if scanMinutes <= 0 {
			return true
		}
		if lastRun.IsZero() {
			return true
		}
		return now.Sub(lastRun) >= time.Duration(scanMinutes)*time.Minute
	}
}

func (m *Manager) ListArchives() ([]storage.RecordingArchive, error) {
	return m.db.ListRecordingArchives("", 0)
}

func (m *Manager) ArchiveRecording(ctx context.Context, streamKey, filename string) (*storage.RecordingArchive, error) {
	settings := m.Settings()
	target := settings.Recording
	if !target.Configured || !settings.RecordingsEnabled {
		return nil, fmt.Errorf("arsivleme etkin degil")
	}
	if m.recordings == nil {
		return nil, fmt.Errorf("kayit yoneticisi hazir degil")
	}
	streamKey = strings.TrimSpace(streamKey)
	filename = filepath.Base(strings.TrimSpace(filename))
	if streamKey == "" || filename == "" {
		return nil, fmt.Errorf("stream_key ve filename gerekli")
	}
	localPath := m.recordings.RecordingFilePath(streamKey, filename)
	info, err := os.Stat(localPath)
	if err != nil {
		return nil, err
	}
	client, err := m.newClient(target)
	if err != nil {
		return nil, err
	}
	objectKey := buildObjectKey(target.Prefix, streamKey, filename)
	item := &storage.RecordingArchive{
		StreamKey: streamKey,
		Filename:  filename,
		Format:    strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), "."),
		Provider:  providerLabel(target.Provider, target.ProviderVariant),
		Bucket:    target.Bucket,
		Endpoint:  target.Endpoint,
		ObjectKey: objectKey,
		Size:      info.Size(),
		Status:    "archived",
	}
	obj, err := client.UploadFile(ctx, objectKey, localPath, detectArchiveContentType(filename))
	if err != nil {
		item.Status = "error"
		item.LastError = err.Error()
		_ = m.db.UpsertRecordingArchive(item)
		m.setLastError(err)
		return nil, err
	}
	item.ObjectURL = obj.URL
	item.ETag = obj.ETag
	item.Size = obj.Size
	item.ArchivedAt = time.Now()
	if target.DeleteLocal {
		if err := os.Remove(localPath); err == nil {
			item.LocalDeleted = true
		}
	}
	if err := m.db.UpsertRecordingArchive(item); err != nil {
		return nil, err
	}
	m.markSyncSuccess()
	return m.db.GetRecordingArchive(streamKey, filename)
}

func (m *Manager) RestoreRecording(ctx context.Context, streamKey, filename string) (*storage.RecordingArchive, error) {
	settings := m.Settings()
	target := settings.Recording
	if !target.Configured || !settings.RecordingsEnabled {
		return nil, fmt.Errorf("arsivleme etkin degil")
	}
	streamKey = strings.TrimSpace(streamKey)
	filename = filepath.Base(strings.TrimSpace(filename))
	if streamKey == "" || filename == "" {
		return nil, fmt.Errorf("stream_key ve filename gerekli")
	}
	item, err := m.db.GetRecordingArchive(streamKey, filename)
	if err != nil {
		return nil, err
	}
	if item == nil || strings.TrimSpace(item.ObjectKey) == "" {
		return nil, fmt.Errorf("arsiv kaydi bulunamadi")
	}
	client, err := m.newClient(target)
	if err != nil {
		return nil, err
	}
	if m.recordings == nil {
		return nil, fmt.Errorf("kayit yoneticisi hazir degil")
	}
	targetPath := m.recordings.RecordingFilePath(streamKey, filename)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return nil, err
	}
	size, err := client.DownloadFile(ctx, item.ObjectKey, targetPath)
	if err != nil {
		item.Status = "error"
		item.LastError = err.Error()
		_ = m.db.UpsertRecordingArchive(item)
		m.setLastError(err)
		return nil, err
	}
	item.Size = size
	item.Status = "archived"
	item.LastError = ""
	item.LocalDeleted = false
	item.RestoredAt = time.Now()
	if err := m.db.UpsertRecordingArchive(item); err != nil {
		return nil, err
	}
	m.markSyncSuccess()
	return m.db.GetRecordingArchive(streamKey, filename)
}

func (m *Manager) SyncPending(ctx context.Context, limit int) (int, error) {
	settings := m.Settings()
	target := settings.Recording
	if !target.Configured || !settings.RecordingsEnabled || !target.AutoUpload || m.recordings == nil {
		return 0, nil
	}
	if limit <= 0 {
		limit = target.BatchSize
	}
	files, err := m.recordings.ListAllRecordingFiles()
	if err != nil {
		m.setLastError(err)
		return 0, err
	}
	archives, err := m.db.ListRecordingArchives("", 0)
	if err != nil {
		m.setLastError(err)
		return 0, err
	}
	known := make(map[string]storage.RecordingArchive, len(archives))
	for _, item := range archives {
		known[item.StreamKey+"::"+item.Filename] = item
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.Before(files[j].ModTime)
	})
	uploaded := 0
	for _, file := range files {
		key := file.StreamKey + "::" + file.Name
		if existing, ok := known[key]; ok && strings.EqualFold(existing.Status, "archived") && strings.TrimSpace(existing.ObjectKey) != "" {
			continue
		}
		if _, err := m.ArchiveRecording(ctx, file.StreamKey, file.Name); err != nil {
			return uploaded, err
		}
		uploaded++
		if uploaded >= limit {
			break
		}
	}
	m.markSyncSuccess()
	return uploaded, nil
}

func (m *Manager) newClient(settings TargetSettings) (storeClient, error) {
	switch settings.Provider {
	case "local":
		return &localStore{
			rootDir:       settings.LocalDir,
			publicBaseURL: settings.PublicBaseURL,
		}, nil
	case "s3", "minio":
		return &s3Store{
			endpoint:      settings.Endpoint,
			region:        settings.Region,
			bucket:        settings.Bucket,
			accessKey:     settings.AccessKey,
			secretKey:     settings.SecretKey,
			usePathStyle:  settings.UsePathStyle,
			publicBaseURL: settings.PublicBaseURL,
			httpClient:    &http.Client{Timeout: 2 * time.Minute},
		}, nil
	case "sftp":
		return &sftpStore{
			host:                settings.SFTPHost,
			port:                settings.SFTPPort,
			user:                settings.SFTPUser,
			remoteDir:           settings.SFTPRemoteDir,
			keyPath:             settings.SFTPKeyPath,
			disableHostKeyCheck: settings.SFTPDisableHostKeyCheck,
			publicBaseURL:       settings.PublicBaseURL,
		}, nil
	case "rclone":
		return &rcloneStore{
			remote:        settings.RcloneRemote,
			remotePath:    settings.RclonePath,
			configPath:    settings.RcloneConfigPath,
			publicBaseURL: settings.PublicBaseURL,
		}, nil
	default:
		return nil, fmt.Errorf("desteklenmeyen arsiv saglayicisi: %s", settings.Provider)
	}
}

func (m *Manager) markSyncSuccess() {
	m.mu.Lock()
	m.lastSyncAt = time.Now()
	m.lastError = ""
	m.mu.Unlock()
}

func (m *Manager) setLastError(err error) {
	if err == nil {
		return
	}
	m.mu.Lock()
	m.lastError = err.Error()
	m.mu.Unlock()
}

func (m *Manager) markBackupSyncSuccess() {
	m.mu.Lock()
	m.lastBackupSyncAt = time.Now()
	m.lastBackupError = ""
	m.mu.Unlock()
}

func (m *Manager) setLastBackupError(err error) {
	if err == nil {
		return
	}
	m.mu.Lock()
	m.lastBackupError = err.Error()
	m.mu.Unlock()
}

type localStore struct {
	rootDir       string
	publicBaseURL string
}

func (s *localStore) UploadFile(_ context.Context, objectKey, localPath, _ string) (storeObject, error) {
	if err := os.MkdirAll(s.rootDir, 0755); err != nil {
		return storeObject{}, err
	}
	src, err := os.Open(localPath)
	if err != nil {
		return storeObject{}, err
	}
	defer src.Close()

	destPath := filepath.Join(s.rootDir, filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return storeObject{}, err
	}
	dst, err := os.Create(destPath)
	if err != nil {
		return storeObject{}, err
	}
	defer dst.Close()

	hasher := sha256.New()
	written, err := io.Copy(io.MultiWriter(dst, hasher), src)
	if err != nil {
		return storeObject{}, err
	}
	url := ""
	if strings.TrimSpace(s.publicBaseURL) != "" {
		url = strings.TrimRight(s.publicBaseURL, "/") + "/" + strings.TrimLeft(objectKey, "/")
	}
	return storeObject{
		Key:  objectKey,
		URL:  url,
		ETag: hex.EncodeToString(hasher.Sum(nil)),
		Size: written,
	}, nil
}

func (s *localStore) DownloadFile(_ context.Context, objectKey, destPath string) (int64, error) {
	src, err := os.Open(filepath.Join(s.rootDir, filepath.FromSlash(objectKey)))
	if err != nil {
		return 0, err
	}
	defer src.Close()
	tmpPath := destPath + ".restore"
	dst, err := os.Create(tmpPath)
	if err != nil {
		return 0, err
	}
	written, err := io.Copy(dst, src)
	closeErr := dst.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	return written, nil
}

func (s *localStore) DeleteObject(_ context.Context, objectKey string) error {
	targetPath := filepath.Join(s.rootDir, filepath.FromSlash(objectKey))
	if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

type s3Store struct {
	endpoint      string
	region        string
	bucket        string
	accessKey     string
	secretKey     string
	usePathStyle  bool
	publicBaseURL string
	httpClient    *http.Client
}

func (s *s3Store) UploadFile(ctx context.Context, objectKey, localPath, contentType string) (storeObject, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return storeObject{}, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return storeObject{}, err
	}
	payloadHash, err := hashSeeker(file)
	if err != nil {
		return storeObject{}, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return storeObject{}, err
	}
	req, canonicalURI, err := s.newSignedRequest(ctx, http.MethodPut, objectKey, nil, file, payloadHash)
	if err != nil {
		return storeObject{}, err
	}
	req.ContentLength = info.Size()
	req.Header.Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return storeObject{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		return storeObject{}, fmt.Errorf("object storage upload hatasi (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return storeObject{
		Key:  objectKey,
		URL:  s.objectURL(objectKey, canonicalURI),
		ETag: strings.Trim(resp.Header.Get("ETag"), `"`),
		Size: info.Size(),
	}, nil
}

func (s *s3Store) DownloadFile(ctx context.Context, objectKey, destPath string) (int64, error) {
	req, _, err := s.newSignedRequest(ctx, http.MethodGet, objectKey, nil, nil, emptySHA256)
	if err != nil {
		return 0, err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		return 0, fmt.Errorf("object storage indirme hatasi (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	tmpPath := destPath + ".restore"
	dst, err := os.Create(tmpPath)
	if err != nil {
		return 0, err
	}
	written, err := io.Copy(dst, resp.Body)
	closeErr := dst.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	return written, nil
}

func (s *s3Store) DeleteObject(ctx context.Context, objectKey string) error {
	req, _, err := s.newSignedRequest(ctx, http.MethodDelete, objectKey, nil, nil, emptySHA256)
	if err != nil {
		return err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted {
		return nil
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	return fmt.Errorf("object storage silme hatasi (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
}

func (s *s3Store) newSignedRequest(ctx context.Context, method, objectKey string, query url.Values, body io.Reader, payloadHash string) (*http.Request, string, error) {
	endpointURL, host, canonicalURI, canonicalQuery, err := s.buildObjectURL(objectKey, query)
	if err != nil {
		return nil, "", err
	}
	req, err := http.NewRequestWithContext(ctx, method, endpointURL, body)
	if err != nil {
		return nil, "", err
	}
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	req.Host = host
	req.Header.Set("Host", host)
	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)

	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", host, payloadHash, amzDate)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalRequest := strings.Join([]string{
		method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")
	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, s.region)
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		credentialScope,
		hashString(canonicalRequest),
	}, "\n")
	signingKey := deriveAWS4SigningKey(s.secretKey, dateStamp, s.region, "s3")
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))
	req.Header.Set("Authorization", fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKey, credentialScope, signedHeaders, signature,
	))
	return req, canonicalURI, nil
}

func (s *s3Store) buildObjectURL(objectKey string, query url.Values) (string, string, string, string, error) {
	base, err := url.Parse(s.endpoint)
	if err != nil {
		return "", "", "", "", err
	}
	cleanKey := strings.TrimLeft(objectKey, "/")
	basePath := strings.TrimRight(base.EscapedPath(), "/")
	objectPath := escapeObjectPath(cleanKey)
	canonicalURI := ""
	if s.usePathStyle {
		canonicalURI = joinURLPath(basePath, "/"+url.PathEscape(s.bucket)+"/"+objectPath)
	} else {
		base.Host = s.bucket + "." + base.Host
		canonicalURI = joinURLPath(basePath, "/"+objectPath)
	}
	base.Path = canonicalURI
	if query == nil {
		query = url.Values{}
	}
	base.RawQuery = query.Encode()
	return base.String(), base.Host, canonicalURI, query.Encode(), nil
}

func (s *s3Store) objectURL(objectKey, canonicalURI string) string {
	if strings.TrimSpace(s.publicBaseURL) != "" {
		return strings.TrimRight(s.publicBaseURL, "/") + "/" + strings.TrimLeft(objectKey, "/")
	}
	base, err := url.Parse(s.endpoint)
	if err != nil {
		return ""
	}
	if s.usePathStyle {
		base.Path = canonicalURI
	} else {
		base.Host = s.bucket + "." + base.Host
		base.Path = canonicalURI
	}
	base.RawQuery = ""
	return base.String()
}

func buildObjectKey(prefix, streamKey, filename string) string {
	parts := []string{}
	if prefix = strings.Trim(prefix, "/"); prefix != "" {
		parts = append(parts, prefix)
	}
	parts = append(parts, "recordings", sanitizeObjectPart(streamKey), sanitizeObjectPart(filepath.Base(filename)))
	return path.Join(parts...)
}

func detectArchiveContentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".flv":
		return "video/x-flv"
	case ".ts":
		return "video/mp2t"
	case ".mp3":
		return "audio/mpeg"
	case ".aac":
		return "audio/aac"
	case ".ogg":
		return "audio/ogg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	default:
		return "application/octet-stream"
	}
}

func providerLabel(raw, variant string) string {
	variant = strings.TrimSpace(strings.ToLower(variant))
	if variant != "" {
		return variant
	}
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "minio":
		return "minio"
	default:
		return strings.ToLower(strings.TrimSpace(raw))
	}
}

func sanitizeObjectPart(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}
	value = strings.ReplaceAll(value, "\\", "_")
	value = strings.ReplaceAll(value, "..", "_")
	value = strings.ReplaceAll(value, " ", "_")
	return strings.Trim(value, "/")
}

func escapeObjectPath(objectKey string) string {
	parts := strings.Split(strings.TrimLeft(objectKey, "/"), "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

func joinURLPath(basePath, suffix string) string {
	basePath = strings.TrimRight(basePath, "/")
	suffix = strings.TrimLeft(suffix, "/")
	if basePath == "" {
		return "/" + suffix
	}
	if suffix == "" {
		return basePath
	}
	return basePath + "/" + suffix
}

func hashSeeker(rs io.ReadSeeker) (string, error) {
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, rs); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func hashString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(data))
	return mac.Sum(nil)
}

func deriveAWS4SigningKey(secret, dateStamp, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), dateStamp)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	return hmacSHA256(kService, "aws4_request")
}

var emptySHA256 = hashString("")
