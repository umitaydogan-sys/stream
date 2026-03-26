package storage

import (
	"database/sql"
	"strings"
	"time"
)

func (s *SQLiteDB) ListEmbedProfiles(streamKey string) ([]EmbedProfile, error) {
	streamKey = strings.TrimSpace(streamKey)
	query := `SELECT id, stream_key, name, use_case, mode, primary_format, width, height, theme,
		options_json, branding_json, security_json, notes, created_at, updated_at
		FROM embed_profiles`
	args := []interface{}{}
	if streamKey != "" {
		query += ` WHERE stream_key = ? OR stream_key = ''`
		args = append(args, streamKey)
	}
	query += ` ORDER BY CASE WHEN stream_key = '' THEN 1 ELSE 0 END, updated_at DESC, id DESC`
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]EmbedProfile, 0, 32)
	for rows.Next() {
		var item EmbedProfile
		if err := rows.Scan(
			&item.ID,
			&item.StreamKey,
			&item.Name,
			&item.UseCase,
			&item.Mode,
			&item.PrimaryFormat,
			&item.Width,
			&item.Height,
			&item.Theme,
			&item.OptionsJSON,
			&item.BrandingJSON,
			&item.SecurityJSON,
			&item.Notes,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SQLiteDB) GetEmbedProfile(id int64) (*EmbedProfile, error) {
	var item EmbedProfile
	err := s.db.QueryRow(
		`SELECT id, stream_key, name, use_case, mode, primary_format, width, height, theme,
		 options_json, branding_json, security_json, notes, created_at, updated_at
		 FROM embed_profiles WHERE id = ?`, id,
	).Scan(
		&item.ID,
		&item.StreamKey,
		&item.Name,
		&item.UseCase,
		&item.Mode,
		&item.PrimaryFormat,
		&item.Width,
		&item.Height,
		&item.Theme,
		&item.OptionsJSON,
		&item.BrandingJSON,
		&item.SecurityJSON,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *SQLiteDB) CreateEmbedProfile(item *EmbedProfile) (int64, error) {
	now := time.Now()
	res, err := s.db.Exec(
		`INSERT INTO embed_profiles
		(stream_key, name, use_case, mode, primary_format, width, height, theme, options_json, branding_json, security_json, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		strings.TrimSpace(item.StreamKey),
		strings.TrimSpace(item.Name),
		strings.TrimSpace(item.UseCase),
		strings.TrimSpace(item.Mode),
		strings.TrimSpace(item.PrimaryFormat),
		item.Width,
		item.Height,
		strings.TrimSpace(item.Theme),
		coalesceJSONString(item.OptionsJSON),
		coalesceJSONString(item.BrandingJSON),
		coalesceJSONString(item.SecurityJSON),
		strings.TrimSpace(item.Notes),
		now,
		now,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *SQLiteDB) UpdateEmbedProfile(item *EmbedProfile) error {
	_, err := s.db.Exec(
		`UPDATE embed_profiles SET
		stream_key = ?, name = ?, use_case = ?, mode = ?, primary_format = ?, width = ?, height = ?, theme = ?,
		options_json = ?, branding_json = ?, security_json = ?, notes = ?, updated_at = ?
		WHERE id = ?`,
		strings.TrimSpace(item.StreamKey),
		strings.TrimSpace(item.Name),
		strings.TrimSpace(item.UseCase),
		strings.TrimSpace(item.Mode),
		strings.TrimSpace(item.PrimaryFormat),
		item.Width,
		item.Height,
		strings.TrimSpace(item.Theme),
		coalesceJSONString(item.OptionsJSON),
		coalesceJSONString(item.BrandingJSON),
		coalesceJSONString(item.SecurityJSON),
		strings.TrimSpace(item.Notes),
		time.Now(),
		item.ID,
	)
	return err
}

func (s *SQLiteDB) DeleteEmbedProfile(id int64) error {
	_, err := s.db.Exec(`DELETE FROM embed_profiles WHERE id = ?`, id)
	return err
}

func (s *SQLiteDB) ListABRProfiles(streamKey string) ([]ABRProfile, error) {
	streamKey = strings.TrimSpace(streamKey)
	query := `SELECT id, profile_set, name, scope, stream_key, description, preset, profiles_json, summary_json, created_at, updated_at
		FROM abr_profiles`
	args := []interface{}{}
	if streamKey != "" {
		query += ` WHERE scope = 'global' OR stream_key = ?`
		args = append(args, streamKey)
	}
	query += ` ORDER BY CASE WHEN scope = 'global' THEN 0 ELSE 1 END, updated_at DESC, id DESC`
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ABRProfile, 0, 32)
	for rows.Next() {
		var item ABRProfile
		if err := rows.Scan(
			&item.ID,
			&item.ProfileSet,
			&item.Name,
			&item.Scope,
			&item.StreamKey,
			&item.Description,
			&item.Preset,
			&item.ProfilesJSON,
			&item.SummaryJSON,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *SQLiteDB) GetABRProfile(id int64) (*ABRProfile, error) {
	var item ABRProfile
	err := s.db.QueryRow(
		`SELECT id, profile_set, name, scope, stream_key, description, preset, profiles_json, summary_json, created_at, updated_at
		 FROM abr_profiles WHERE id = ?`, id,
	).Scan(
		&item.ID,
		&item.ProfileSet,
		&item.Name,
		&item.Scope,
		&item.StreamKey,
		&item.Description,
		&item.Preset,
		&item.ProfilesJSON,
		&item.SummaryJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *SQLiteDB) CreateABRProfile(item *ABRProfile) (int64, error) {
	now := time.Now()
	res, err := s.db.Exec(
		`INSERT INTO abr_profiles
		(profile_set, name, scope, stream_key, description, preset, profiles_json, summary_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		strings.TrimSpace(item.ProfileSet),
		strings.TrimSpace(item.Name),
		strings.TrimSpace(item.Scope),
		strings.TrimSpace(item.StreamKey),
		strings.TrimSpace(item.Description),
		strings.TrimSpace(item.Preset),
		coalesceJSONString(item.ProfilesJSON),
		coalesceJSONString(item.SummaryJSON),
		now,
		now,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *SQLiteDB) UpdateABRProfile(item *ABRProfile) error {
	_, err := s.db.Exec(
		`UPDATE abr_profiles SET
		profile_set = ?, name = ?, scope = ?, stream_key = ?, description = ?, preset = ?, profiles_json = ?, summary_json = ?, updated_at = ?
		WHERE id = ?`,
		strings.TrimSpace(item.ProfileSet),
		strings.TrimSpace(item.Name),
		strings.TrimSpace(item.Scope),
		strings.TrimSpace(item.StreamKey),
		strings.TrimSpace(item.Description),
		strings.TrimSpace(item.Preset),
		coalesceJSONString(item.ProfilesJSON),
		coalesceJSONString(item.SummaryJSON),
		time.Now(),
		item.ID,
	)
	return err
}

func (s *SQLiteDB) DeleteABRProfile(id int64) error {
	_, err := s.db.Exec(`DELETE FROM abr_profiles WHERE id = ?`, id)
	return err
}

func coalesceJSONString(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "{}"
	}
	return raw
}
