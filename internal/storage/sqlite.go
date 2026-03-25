package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/textutil"
	_ "modernc.org/sqlite"
)

// SQLiteDB wraps a SQLite database connection
type SQLiteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates/opens a SQLite database and runs migrations
func NewSQLiteDB(path string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite is single-writer
	db.SetMaxIdleConns(1)

	s := &SQLiteDB{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}
	return s, nil
}

// Close closes the database
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// DB returns the underlying sql.DB
func (s *SQLiteDB) DB() *sql.DB {
	return s.db
}

// ExportBackupSnapshot writes a consistent SQLite snapshot to destPath.
func (s *SQLiteDB) ExportBackupSnapshot(destPath string) error {
	escaped := strings.ReplaceAll(destPath, "'", "''")
	if _, err := s.db.Exec("PRAGMA wal_checkpoint(FULL)"); err != nil {
		return err
	}
	_, err := s.db.Exec(fmt.Sprintf("VACUUM INTO '%s'", escaped))
	return err
}

func (s *SQLiteDB) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS config (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL DEFAULT '',
			category TEXT NOT NULL DEFAULT 'general',
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'admin',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_login DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS streams (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			stream_key TEXT UNIQUE NOT NULL,
			status TEXT NOT NULL DEFAULT 'offline',
			ingest_proto TEXT DEFAULT '',
			output_formats TEXT DEFAULT '["hls"]',
			policy_json TEXT DEFAULT '',
			max_viewers INTEGER DEFAULT 0,
			max_bitrate INTEGER DEFAULT 0,
			record_enabled INTEGER DEFAULT 0,
			record_format TEXT DEFAULT 'mp4',
			password TEXT DEFAULT '',
			domain_lock TEXT DEFAULT '',
			ip_whitelist TEXT DEFAULT '',
			thumbnail_path TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			started_at DATETIME,
			viewer_count INTEGER DEFAULT 0,
			input_bitrate INTEGER DEFAULT 0,
			input_codec TEXT DEFAULT '',
			input_width INTEGER DEFAULT 0,
			input_height INTEGER DEFAULT 0,
			input_fps REAL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS viewers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stream_id INTEGER NOT NULL,
			ip TEXT NOT NULL,
			user_agent TEXT DEFAULT '',
			protocol TEXT DEFAULT 'hls',
			country TEXT DEFAULT '',
			city TEXT DEFAULT '',
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			bandwidth INTEGER DEFAULT 0,
			FOREIGN KEY (stream_id) REFERENCES streams(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS banned_ips (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT UNIQUE NOT NULL,
			reason TEXT DEFAULT '',
			banned_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level TEXT NOT NULL DEFAULT 'INFO',
			component TEXT NOT NULL DEFAULT 'system',
			message TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS player_templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			background_css TEXT DEFAULT '',
			control_bar_css TEXT DEFAULT '',
			play_button_css TEXT DEFAULT '',
			logo_url TEXT DEFAULT '',
			logo_position TEXT DEFAULT 'top-right',
			logo_opacity REAL DEFAULT 1.0,
			watermark_text TEXT DEFAULT '',
			show_title INTEGER DEFAULT 1,
			show_live_badge INTEGER DEFAULT 1,
			theme TEXT DEFAULT 'dark',
			custom_css TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS analytics_snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			total_streams INTEGER DEFAULT 0,
			total_viewers INTEGER DEFAULT 0,
			current_viewers INTEGER DEFAULT 0,
			peak_concurrent INTEGER DEFAULT 0,
			total_bandwidth INTEGER DEFAULT 0,
			viewers_by_format TEXT DEFAULT '{}',
			viewers_by_country TEXT DEFAULT '{}'
		)`,
		`CREATE TABLE IF NOT EXISTS player_telemetry_samples (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stream_key TEXT NOT NULL,
			active_sessions INTEGER DEFAULT 0,
			waiting_sessions INTEGER DEFAULT 0,
			offline_sessions INTEGER DEFAULT 0,
			debug_sessions INTEGER DEFAULT 0,
			total_stalls INTEGER DEFAULT 0,
			total_recoveries INTEGER DEFAULT 0,
			average_buffer_seconds REAL DEFAULT 0,
			average_playback_seconds REAL DEFAULT 0,
			last_error TEXT DEFAULT '',
			sources_json TEXT DEFAULT '{}',
			formats_json TEXT DEFAULT '{}',
			pages_json TEXT DEFAULT '{}',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS track_telemetry_samples (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			stream_key TEXT NOT NULL,
			track_id INTEGER NOT NULL,
			kind TEXT NOT NULL DEFAULT 'video',
			codec TEXT DEFAULT '',
			width INTEGER DEFAULT 0,
			height INTEGER DEFAULT 0,
			sample_rate INTEGER DEFAULT 0,
			channels INTEGER DEFAULT 0,
			bitrate INTEGER DEFAULT 0,
			packets INTEGER DEFAULT 0,
			bytes INTEGER DEFAULT 0,
			is_default INTEGER DEFAULT 0,
			is_active INTEGER DEFAULT 0,
			display_label TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_created ON logs(created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level)`,
		`CREATE INDEX IF NOT EXISTS idx_streams_key ON streams(stream_key)`,
		`CREATE INDEX IF NOT EXISTS idx_streams_status ON streams(status)`,
		`CREATE INDEX IF NOT EXISTS idx_viewers_stream ON viewers(stream_id)`,
		`CREATE INDEX IF NOT EXISTS idx_analytics_snapshots_ts ON analytics_snapshots(timestamp DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_player_telemetry_stream_created ON player_telemetry_samples(stream_key, created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_track_telemetry_stream_created ON track_telemetry_samples(stream_key, created_at DESC)`,
	}

	for _, m := range migrations {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("migration exec: %w", err)
		}
	}
	// Backward-compatible additive migrations for existing installs.
	_ = s.ensureColumn("streams", "policy_json", "TEXT NOT NULL DEFAULT ''")
	return nil
}

func (s *SQLiteDB) ensureColumn(table, column, definition string) error {
	rows, err := s.db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return err
		}
		if strings.EqualFold(name, column) {
			return nil
		}
	}
	_, err = s.db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, definition))
	return err
}

// ─── Config Operations ───────────────────────────────────────

func (s *SQLiteDB) GetConfig(key string) (string, error) {
	var value string
	err := s.db.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (s *SQLiteDB) SetConfig(key, value, category string) error {
	_, err := s.db.Exec(
		`INSERT INTO config (key, value, category, updated_at) VALUES (?, ?, ?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value, category = excluded.category, updated_at = excluded.updated_at`,
		key, value, category, time.Now(),
	)
	return err
}

func (s *SQLiteDB) GetConfigByCategory(category string) ([]Config, error) {
	rows, err := s.db.Query("SELECT key, value, category, updated_at FROM config WHERE category = ? ORDER BY key", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []Config
	for rows.Next() {
		var c Config
		if err := rows.Scan(&c.Key, &c.Value, &c.Category, &c.UpdatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
}

func (s *SQLiteDB) GetAllConfig() (map[string]string, error) {
	rows, err := s.db.Query("SELECT key, value FROM config")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

// ─── User Operations ─────────────────────────────────────────

func (s *SQLiteDB) CreateUser(username, passwordHash, role string) (int64, error) {
	result, err := s.db.Exec(
		"INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		username, passwordHash, role,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *SQLiteDB) GetUserByUsername(username string) (*User, error) {
	u := &User{}
	err := s.db.QueryRow(
		"SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE username = ?",
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (s *SQLiteDB) GetUsers() ([]User, error) {
	rows, err := s.db.Query("SELECT id, username, role, created_at, updated_at FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *SQLiteDB) UpdateUserLogin(id int64) error {
	_, err := s.db.Exec("UPDATE users SET last_login = ? WHERE id = ?", time.Now(), id)
	return err
}

// FixCorruptedUsers removes users whose password hash is corrupted (wrong length)
// and resets setup_completed so the wizard runs again.
func (s *SQLiteDB) FixCorruptedUsers() (int, error) {
	// A valid HMAC-SHA256 hex hash is exactly 64 characters
	result, err := s.db.Exec("DELETE FROM users WHERE length(password_hash) <> 64")
	if err != nil {
		return 0, err
	}
	affected, _ := result.RowsAffected()
	if affected > 0 {
		// Reset setup so wizard runs again
		s.db.Exec("DELETE FROM config WHERE key = 'setup_completed'")
	}
	return int(affected), nil
}

func (s *SQLiteDB) UpdateUserPassword(id int64, passwordHash string) error {
	_, err := s.db.Exec("UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?", passwordHash, time.Now(), id)
	return err
}

// ─── Stream Operations ───────────────────────────────────────

func (s *SQLiteDB) CreateStream(st *Stream) (int64, error) {
	result, err := s.db.Exec(
		`INSERT INTO streams (name, description, stream_key, status, output_formats, policy_json, max_viewers, max_bitrate, record_enabled, record_format, password, domain_lock, ip_whitelist)
		 VALUES (?, ?, ?, 'offline', ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		st.Name, st.Description, st.StreamKey, st.OutputFormats, st.PolicyJSON, st.MaxViewers, st.MaxBitrate,
		st.RecordEnabled, st.RecordFormat, st.Password, st.DomainLock, st.IPWhitelist,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *SQLiteDB) GetStreamByKey(key string) (*Stream, error) {
	st := &Stream{}
	err := s.db.QueryRow(
		`SELECT id, name, description, stream_key, status, ingest_proto, output_formats, policy_json,
		 max_viewers, max_bitrate, record_enabled, record_format, password, domain_lock,
		 ip_whitelist, thumbnail_path, created_at, updated_at, viewer_count, input_bitrate,
		 input_codec, input_width, input_height, input_fps
		 FROM streams WHERE stream_key = ?`, key,
	).Scan(
		&st.ID, &st.Name, &st.Description, &st.StreamKey, &st.Status, &st.IngestProto,
		&st.OutputFormats, &st.PolicyJSON, &st.MaxViewers, &st.MaxBitrate, &st.RecordEnabled, &st.RecordFormat,
		&st.Password, &st.DomainLock, &st.IPWhitelist, &st.ThumbnailPath, &st.CreatedAt,
		&st.UpdatedAt, &st.ViewerCount, &st.InputBitrate, &st.InputCodec, &st.InputWidth,
		&st.InputHeight, &st.InputFPS,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return st, err
}

func (s *SQLiteDB) GetStreamByID(id int64) (*Stream, error) {
	st := &Stream{}
	err := s.db.QueryRow(
		`SELECT id, name, description, stream_key, status, ingest_proto, output_formats, policy_json,
		 max_viewers, max_bitrate, record_enabled, record_format, password, domain_lock,
		 ip_whitelist, thumbnail_path, created_at, updated_at, viewer_count, input_bitrate,
		 input_codec, input_width, input_height, input_fps
		 FROM streams WHERE id = ?`, id,
	).Scan(
		&st.ID, &st.Name, &st.Description, &st.StreamKey, &st.Status, &st.IngestProto,
		&st.OutputFormats, &st.PolicyJSON, &st.MaxViewers, &st.MaxBitrate, &st.RecordEnabled, &st.RecordFormat,
		&st.Password, &st.DomainLock, &st.IPWhitelist, &st.ThumbnailPath, &st.CreatedAt,
		&st.UpdatedAt, &st.ViewerCount, &st.InputBitrate, &st.InputCodec, &st.InputWidth,
		&st.InputHeight, &st.InputFPS,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return st, err
}

func (s *SQLiteDB) GetAllStreams() ([]Stream, error) {
	rows, err := s.db.Query(
		`SELECT id, name, description, stream_key, status, ingest_proto, output_formats, policy_json,
		 max_viewers, max_bitrate, record_enabled, record_format, thumbnail_path,
		 created_at, updated_at, viewer_count, input_bitrate, input_codec, input_width, input_height, input_fps
		 FROM streams ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streams []Stream
	for rows.Next() {
		var st Stream
		if err := rows.Scan(
			&st.ID, &st.Name, &st.Description, &st.StreamKey, &st.Status, &st.IngestProto,
			&st.OutputFormats, &st.PolicyJSON, &st.MaxViewers, &st.MaxBitrate, &st.RecordEnabled, &st.RecordFormat,
			&st.ThumbnailPath, &st.CreatedAt, &st.UpdatedAt, &st.ViewerCount, &st.InputBitrate,
			&st.InputCodec, &st.InputWidth, &st.InputHeight, &st.InputFPS,
		); err != nil {
			return nil, err
		}
		streams = append(streams, st)
	}
	return streams, nil
}

func (s *SQLiteDB) UpdateStreamStatus(key, status, ingestProto string) error {
	now := time.Now()
	if status == "live" {
		_, err := s.db.Exec(
			"UPDATE streams SET status = ?, ingest_proto = ?, started_at = ?, updated_at = ? WHERE stream_key = ?",
			status, ingestProto, now, now, key,
		)
		return err
	}
	_, err := s.db.Exec(
		`UPDATE streams
		 SET status = ?, ingest_proto = '', started_at = NULL, updated_at = ?,
		     viewer_count = 0, input_bitrate = 0, input_codec = '',
		     input_width = 0, input_height = 0, input_fps = 0
		 WHERE stream_key = ?`,
		status, now, key,
	)
	return err
}

func (s *SQLiteDB) ResetRuntimeStreamState() (int64, error) {
	result, err := s.db.Exec(
		`UPDATE streams
		 SET status = 'offline', ingest_proto = '', started_at = NULL, updated_at = ?,
		     viewer_count = 0, input_bitrate = 0, input_codec = '',
		     input_width = 0, input_height = 0, input_fps = 0
		 WHERE status != 'offline'
		    OR ingest_proto != ''
		    OR viewer_count != 0
		    OR input_bitrate != 0
		    OR input_codec != ''
		    OR input_width != 0
		    OR input_height != 0
		    OR input_fps != 0`,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}
	affected, _ := result.RowsAffected()
	return affected, nil
}

func (s *SQLiteDB) UpdateStreamMeta(key string, codec string, width, height int, fps float64, bitrate int64) error {
	_, err := s.db.Exec(
		"UPDATE streams SET input_codec = ?, input_width = ?, input_height = ?, input_fps = ?, input_bitrate = ?, updated_at = ? WHERE stream_key = ?",
		codec, width, height, fps, bitrate, time.Now(), key,
	)
	return err
}

func (s *SQLiteDB) UpdateStream(st *Stream) error {
	_, err := s.db.Exec(
		`UPDATE streams SET name = ?, description = ?, output_formats = ?, policy_json = ?, max_viewers = ?,
		 max_bitrate = ?, record_enabled = ?, record_format = ?, password = ?,
		 domain_lock = ?, ip_whitelist = ?, updated_at = ? WHERE id = ?`,
		st.Name, st.Description, st.OutputFormats, st.PolicyJSON, st.MaxViewers, st.MaxBitrate,
		st.RecordEnabled, st.RecordFormat, st.Password, st.DomainLock, st.IPWhitelist,
		time.Now(), st.ID,
	)
	return err
}

func (s *SQLiteDB) DeleteStream(id int64) error {
	_, err := s.db.Exec("DELETE FROM streams WHERE id = ?", id)
	return err
}

func (s *SQLiteDB) UpdateViewerCount(key string, count int) error {
	_, err := s.db.Exec("UPDATE streams SET viewer_count = ? WHERE stream_key = ?", count, key)
	return err
}

// ─── Log Operations ──────────────────────────────────────────

func (s *SQLiteDB) AddLog(level, component, message string) error {
	level = textutil.FixLegacyUTF8String(level)
	component = textutil.FixLegacyUTF8String(component)
	message = textutil.FixLegacyUTF8String(message)
	_, err := s.db.Exec(
		"INSERT INTO logs (level, component, message) VALUES (?, ?, ?)",
		level, component, message,
	)
	return err
}

func (s *SQLiteDB) GetLogs(limit int, level, component string) ([]LogEntry, error) {
	query := "SELECT id, level, component, message, created_at FROM logs WHERE 1=1"
	args := []interface{}{}

	if level != "" {
		query += " AND level = ?"
		args = append(args, level)
	}
	if component != "" {
		query += " AND component = ?"
		args = append(args, component)
	}
	query += " ORDER BY id DESC LIMIT ?"
	args = append(args, limit)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.ID, &l.Level, &l.Component, &l.Message, &l.CreatedAt); err != nil {
			return nil, err
		}
		l.Level = textutil.FixLegacyUTF8String(l.Level)
		l.Component = textutil.FixLegacyUTF8String(l.Component)
		l.Message = textutil.FixLegacyUTF8String(l.Message)
		logs = append(logs, l)
	}
	return logs, nil
}

func (s *SQLiteDB) SavePlayerTelemetrySample(sample *PlayerTelemetrySample) error {
	if sample == nil || strings.TrimSpace(sample.StreamKey) == "" {
		return nil
	}
	lastError := textutil.FixLegacyUTF8String(sample.LastError)
	_, err := s.db.Exec(
		`INSERT INTO player_telemetry_samples
		(stream_key, active_sessions, waiting_sessions, offline_sessions, debug_sessions,
		 total_stalls, total_recoveries, average_buffer_seconds, average_playback_seconds,
		 last_error, sources_json, formats_json, pages_json, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sample.StreamKey,
		sample.ActiveSessions,
		sample.WaitingSessions,
		sample.OfflineSessions,
		sample.DebugSessions,
		sample.TotalStalls,
		sample.TotalRecoveries,
		sample.AverageBufferSeconds,
		sample.AveragePlaybackSeconds,
		lastError,
		sample.SourcesJSON,
		sample.FormatsJSON,
		sample.PagesJSON,
		time.Now(),
	)
	return err
}

func (s *SQLiteDB) GetPlayerTelemetrySamples(streamKey string, limit int) ([]PlayerTelemetrySample, error) {
	streamKey = strings.TrimSpace(streamKey)
	if streamKey == "" {
		return []PlayerTelemetrySample{}, nil
	}
	if limit <= 0 {
		limit = 48
	}
	rows, err := s.db.Query(
		`SELECT id, stream_key, active_sessions, waiting_sessions, offline_sessions, debug_sessions,
		        total_stalls, total_recoveries, average_buffer_seconds, average_playback_seconds,
		        last_error, sources_json, formats_json, pages_json, created_at
		   FROM player_telemetry_samples
		  WHERE stream_key = ?
		  ORDER BY created_at DESC
		  LIMIT ?`,
		streamKey, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	samples := make([]PlayerTelemetrySample, 0, limit)
	for rows.Next() {
		var sample PlayerTelemetrySample
		if err := rows.Scan(
			&sample.ID,
			&sample.StreamKey,
			&sample.ActiveSessions,
			&sample.WaitingSessions,
			&sample.OfflineSessions,
			&sample.DebugSessions,
			&sample.TotalStalls,
			&sample.TotalRecoveries,
			&sample.AverageBufferSeconds,
			&sample.AveragePlaybackSeconds,
			&sample.LastError,
			&sample.SourcesJSON,
			&sample.FormatsJSON,
			&sample.PagesJSON,
			&sample.CreatedAt,
		); err != nil {
			return nil, err
		}
		sample.LastError = textutil.FixLegacyUTF8String(sample.LastError)
		samples = append(samples, sample)
	}

	for i, j := 0, len(samples)-1; i < j; i, j = i+1, j-1 {
		samples[i], samples[j] = samples[j], samples[i]
	}
	return samples, nil
}

func (s *SQLiteDB) SaveTrackTelemetrySamples(samples []TrackTelemetrySample) error {
	if len(samples) == 0 {
		return nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`INSERT INTO track_telemetry_samples
		(stream_key, track_id, kind, codec, width, height, sample_rate, channels, bitrate,
		 packets, bytes, is_default, is_active, display_label, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, sample := range samples {
		if strings.TrimSpace(sample.StreamKey) == "" || sample.TrackID <= 0 {
			continue
		}
		if _, err := stmt.Exec(
			sample.StreamKey,
			sample.TrackID,
			textutil.FixLegacyUTF8String(strings.TrimSpace(sample.Kind)),
			textutil.FixLegacyUTF8String(strings.TrimSpace(sample.Codec)),
			sample.Width,
			sample.Height,
			sample.SampleRate,
			sample.Channels,
			sample.Bitrate,
			sample.Packets,
			sample.Bytes,
			boolToInt(sample.IsDefault),
			boolToInt(sample.IsActive),
			textutil.FixLegacyUTF8String(strings.TrimSpace(sample.DisplayLabel)),
			now,
		); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *SQLiteDB) GetTrackTelemetrySamples(streamKey string, limit int) ([]TrackTelemetrySample, error) {
	streamKey = strings.TrimSpace(streamKey)
	if streamKey == "" {
		return []TrackTelemetrySample{}, nil
	}
	if limit <= 0 {
		limit = 120
	}
	rows, err := s.db.Query(
		`SELECT id, stream_key, track_id, kind, codec, width, height, sample_rate, channels,
		        bitrate, packets, bytes, is_default, is_active, display_label, created_at
		   FROM track_telemetry_samples
		  WHERE stream_key = ?
		  ORDER BY created_at DESC, track_id ASC
		  LIMIT ?`,
		streamKey, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	samples := make([]TrackTelemetrySample, 0, limit)
	for rows.Next() {
		var sample TrackTelemetrySample
		var isDefault int
		var isActive int
		if err := rows.Scan(
			&sample.ID,
			&sample.StreamKey,
			&sample.TrackID,
			&sample.Kind,
			&sample.Codec,
			&sample.Width,
			&sample.Height,
			&sample.SampleRate,
			&sample.Channels,
			&sample.Bitrate,
			&sample.Packets,
			&sample.Bytes,
			&isDefault,
			&isActive,
			&sample.DisplayLabel,
			&sample.CreatedAt,
		); err != nil {
			return nil, err
		}
		sample.Kind = textutil.FixLegacyUTF8String(sample.Kind)
		sample.Codec = textutil.FixLegacyUTF8String(sample.Codec)
		sample.DisplayLabel = textutil.FixLegacyUTF8String(sample.DisplayLabel)
		sample.IsDefault = isDefault == 1
		sample.IsActive = isActive == 1
		samples = append(samples, sample)
	}
	for i, j := 0, len(samples)-1; i < j; i, j = i+1, j-1 {
		samples[i], samples[j] = samples[j], samples[i]
	}
	return samples, nil
}

func (s *SQLiteDB) ClearLogs() error {
	_, err := s.db.Exec("DELETE FROM logs")
	return err
}

func (s *SQLiteDB) SaveAnalyticsSnapshot(snapshot *AnalyticsSnapshot) error {
	_, err := s.db.Exec(
		`INSERT INTO analytics_snapshots
		(timestamp, total_streams, total_viewers, current_viewers, peak_concurrent, total_bandwidth, viewers_by_format, viewers_by_country)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		snapshot.Timestamp,
		snapshot.TotalStreams,
		snapshot.TotalViewers,
		snapshot.CurrentViewers,
		snapshot.PeakConcurrent,
		snapshot.TotalBandwidth,
		snapshot.ViewersByFormat,
		snapshot.ViewersByCountry,
	)
	return err
}

func (s *SQLiteDB) GetAnalyticsSnapshots(limit int) ([]AnalyticsSnapshot, error) {
	if limit <= 0 {
		limit = 24
	}
	rows, err := s.db.Query(
		`SELECT id, timestamp, total_streams, total_viewers, current_viewers, peak_concurrent, total_bandwidth, viewers_by_format, viewers_by_country
		 FROM analytics_snapshots ORDER BY timestamp DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]AnalyticsSnapshot, 0, limit)
	for rows.Next() {
		var item AnalyticsSnapshot
		if err := rows.Scan(
			&item.ID,
			&item.Timestamp,
			&item.TotalStreams,
			&item.TotalViewers,
			&item.CurrentViewers,
			&item.PeakConcurrent,
			&item.TotalBandwidth,
			&item.ViewersByFormat,
			&item.ViewersByCountry,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *SQLiteDB) GetAnalyticsSnapshotsSince(since time.Time, limit int) ([]AnalyticsSnapshot, error) {
	query := `SELECT id, timestamp, total_streams, total_viewers, current_viewers, peak_concurrent, total_bandwidth, viewers_by_format, viewers_by_country
		 FROM analytics_snapshots WHERE timestamp >= ? ORDER BY timestamp ASC`
	args := []interface{}{since}
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]AnalyticsSnapshot, 0, 256)
	for rows.Next() {
		var item AnalyticsSnapshot
		if err := rows.Scan(
			&item.ID,
			&item.Timestamp,
			&item.TotalStreams,
			&item.TotalViewers,
			&item.CurrentViewers,
			&item.PeakConcurrent,
			&item.TotalBandwidth,
			&item.ViewersByFormat,
			&item.ViewersByCountry,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *SQLiteDB) CleanupAnalyticsSnapshots(maxAge time.Duration) (int64, error) {
	if maxAge <= 0 {
		return 0, nil
	}
	res, err := s.db.Exec("DELETE FROM analytics_snapshots WHERE timestamp < ?", time.Now().Add(-maxAge))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteDB) CleanupPlayerTelemetrySamples(maxAge time.Duration) (int64, error) {
	if maxAge <= 0 {
		return 0, nil
	}
	res, err := s.db.Exec("DELETE FROM player_telemetry_samples WHERE created_at < ?", time.Now().Add(-maxAge))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteDB) CleanupTrackTelemetrySamples(maxAge time.Duration) (int64, error) {
	if maxAge <= 0 {
		return 0, nil
	}
	res, err := s.db.Exec("DELETE FROM track_telemetry_samples WHERE created_at < ?", time.Now().Add(-maxAge))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// ─── Player Template Operations ──────────────────────────────

func (s *SQLiteDB) CreatePlayerTemplate(pt *PlayerTemplate) (int64, error) {
	result, err := s.db.Exec(
		`INSERT INTO player_templates (name, background_css, control_bar_css, play_button_css,
		 logo_url, logo_position, logo_opacity, watermark_text, show_title, show_live_badge,
		 theme, custom_css) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		pt.Name, pt.BackgroundCSS, pt.ControlBarCSS, pt.PlayButtonCSS,
		pt.LogoURL, pt.LogoPosition, pt.LogoOpacity, pt.WatermarkText,
		pt.ShowTitle, pt.ShowLiveBadge, pt.Theme, pt.CustomCSS,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *SQLiteDB) GetPlayerTemplates() ([]PlayerTemplate, error) {
	rows, err := s.db.Query(
		`SELECT id, name, background_css, control_bar_css, play_button_css,
		 logo_url, logo_position, logo_opacity, watermark_text, show_title, show_live_badge,
		 theme, custom_css, created_at, updated_at FROM player_templates ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []PlayerTemplate
	for rows.Next() {
		var pt PlayerTemplate
		if err := rows.Scan(
			&pt.ID, &pt.Name, &pt.BackgroundCSS, &pt.ControlBarCSS, &pt.PlayButtonCSS,
			&pt.LogoURL, &pt.LogoPosition, &pt.LogoOpacity, &pt.WatermarkText,
			&pt.ShowTitle, &pt.ShowLiveBadge, &pt.Theme, &pt.CustomCSS,
			&pt.CreatedAt, &pt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		templates = append(templates, pt)
	}
	return templates, nil
}

func (s *SQLiteDB) GetPlayerTemplateByID(id int64) (*PlayerTemplate, error) {
	pt := &PlayerTemplate{}
	err := s.db.QueryRow(
		`SELECT id, name, background_css, control_bar_css, play_button_css,
		 logo_url, logo_position, logo_opacity, watermark_text, show_title, show_live_badge,
		 theme, custom_css, created_at, updated_at FROM player_templates WHERE id = ?`, id,
	).Scan(
		&pt.ID, &pt.Name, &pt.BackgroundCSS, &pt.ControlBarCSS, &pt.PlayButtonCSS,
		&pt.LogoURL, &pt.LogoPosition, &pt.LogoOpacity, &pt.WatermarkText,
		&pt.ShowTitle, &pt.ShowLiveBadge, &pt.Theme, &pt.CustomCSS,
		&pt.CreatedAt, &pt.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return pt, err
}

func (s *SQLiteDB) UpdatePlayerTemplate(pt *PlayerTemplate) error {
	_, err := s.db.Exec(
		`UPDATE player_templates SET name = ?, background_css = ?, control_bar_css = ?,
		 play_button_css = ?, logo_url = ?, logo_position = ?, logo_opacity = ?,
		 watermark_text = ?, show_title = ?, show_live_badge = ?, theme = ?,
		 custom_css = ?, updated_at = ? WHERE id = ?`,
		pt.Name, pt.BackgroundCSS, pt.ControlBarCSS, pt.PlayButtonCSS,
		pt.LogoURL, pt.LogoPosition, pt.LogoOpacity, pt.WatermarkText,
		pt.ShowTitle, pt.ShowLiveBadge, pt.Theme, pt.CustomCSS,
		time.Now(), pt.ID,
	)
	return err
}

func (s *SQLiteDB) DeletePlayerTemplate(id int64) error {
	_, err := s.db.Exec("DELETE FROM player_templates WHERE id = ?", id)
	return err
}

// ─── User Management ─────────────────────────────────────────

func (s *SQLiteDB) GetUserByID(id int64) (*User, error) {
	u := &User{}
	err := s.db.QueryRow(
		"SELECT id, username, password_hash, role, created_at, updated_at FROM users WHERE id = ?", id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (s *SQLiteDB) UpdateUser(id int64, username, role string) error {
	_, err := s.db.Exec(
		"UPDATE users SET username = ?, role = ?, updated_at = ? WHERE id = ?",
		username, role, time.Now(), id,
	)
	return err
}

func (s *SQLiteDB) DeleteUser(id int64) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
