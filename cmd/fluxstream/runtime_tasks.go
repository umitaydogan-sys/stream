package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/analytics"
	"github.com/fluxstream/fluxstream/internal/archive"
	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/recording"
	"github.com/fluxstream/fluxstream/internal/storage"
	"github.com/fluxstream/fluxstream/internal/stream"
	"github.com/fluxstream/fluxstream/internal/tlsutil"
	"github.com/fluxstream/fluxstream/internal/transcode"
)

type systemAlert struct {
	Level       string `json:"level"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action,omitempty"`
}

func buildLiveOptionsFromConfig(cfg *config.Manager) transcode.LiveOptions {
	opts := transcode.DefaultLiveOptions()
	opts.ABREnabled = cfg.GetBool("abr_enabled", false)
	opts.MasterEnabled = cfg.GetBool("abr_master_enabled", true)
	opts.ProfileSet = cfg.Get("abr_profile_set", "balanced")
	opts.ProfilesJSON = cfg.Get("abr_profiles_json", "")
	opts.Profiles = transcode.ResolveProfiles(opts.ProfileSet, opts.ProfilesJSON)
	opts.SegmentDuration = cfg.GetInt("hls_segment_duration", 2)
	opts.PlaylistLength = cfg.GetInt("hls_playlist_length", 10)
	opts.AudioPassthrough = cfg.GetBool("abr_audio_passthrough", false)
	return opts
}

func startMaintenanceLoops(ctxDone <-chan struct{}, cfg *config.Manager, db *storage.SQLiteDB, tracker *analytics.Tracker, recManager *recording.Manager, archiveManager *archive.Manager, tcManager *transcode.Manager, dataDir string) {
	if cfg.GetBool("analytics_persist_enabled", true) && tracker != nil {
		interval := time.Duration(cfg.GetInt("analytics_snapshot_interval", 60)) * time.Second
		if interval < 15*time.Second {
			interval = 15 * time.Second
		}
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				select {
				case <-ctxDone:
					return
				case <-ticker.C:
					dash := tracker.GetDashboard()
					byFormat, _ := json.Marshal(dash.ViewersByFormat)
					byCountry, _ := json.Marshal(dash.ViewersByCountry)
					_ = db.SaveAnalyticsSnapshot(&storage.AnalyticsSnapshot{
						Timestamp:        time.Now(),
						TotalStreams:     dash.TotalStreams,
						TotalViewers:     dash.TotalViewers,
						CurrentViewers:   dash.CurrentViewers,
						PeakConcurrent:   dash.PeakConcurrent,
						TotalBandwidth:   dash.TotalBandwidth,
						ViewersByFormat:  string(byFormat),
						ViewersByCountry: string(byCountry),
					})
				}
			}
		}()
	}

	if cfg.GetBool("track_analytics_enabled", true) && tcManager != nil {
		interval := time.Duration(cfg.GetInt("track_analytics_interval", 20)) * time.Second
		if interval < 10*time.Second {
			interval = 10 * time.Second
		}
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			persistTrackSamples := func() {
				streams, err := db.GetAllStreams()
				if err != nil {
					return
				}
				for _, st := range streams {
					snapshot := tcManager.GetLiveTrackSnapshot(st.StreamKey)
					if len(snapshot.VideoTracks) == 0 && len(snapshot.AudioTracks) == 0 {
						continue
					}
					if err := db.SaveTrackTelemetrySamples(trackTelemetrySamplesFromSnapshot(snapshot)); err != nil {
						continue
					}
				}
			}
			persistTrackSamples()
			for {
				select {
				case <-ctxDone:
					return
				case <-ticker.C:
					persistTrackSamples()
				}
			}
		}()
	}

	if archiveManager != nil {
		interval := time.Duration(cfg.GetInt("archive_scan_interval_minutes", 10)) * time.Minute
		if interval < 2*time.Minute {
			interval = 2 * time.Minute
		}
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			runArchiveSync := func() {
				if !cfg.GetBool("archive_enabled", false) || !cfg.GetBool("archive_auto_upload", false) {
					return
				}
				if !archiveManager.ShouldRunRecordingSchedule(time.Now()) {
					return
				}
				uploaded, err := archiveManager.SyncPending(context.Background(), cfg.GetInt("archive_batch_size", 3))
				if err != nil {
					_ = db.AddLog("WARN", "archive", fmt.Sprintf("Arsiv senkronizasyonu basarisiz: %v", err))
					return
				}
				if uploaded > 0 {
					_ = db.AddLog("INFO", "archive", fmt.Sprintf("Yeni arsivlenen kayit sayisi: %d", uploaded))
				}
			}
			runArchiveSync()
			for {
				select {
				case <-ctxDone:
					return
				case <-ticker.C:
					runArchiveSync()
				}
			}
		}()
	}

	if archiveManager != nil {
		interval := time.Duration(cfg.GetInt("backup_archive_scan_interval_minutes", 30)) * time.Minute
		if interval < 5*time.Minute {
			interval = 5 * time.Minute
		}
		go func() {
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			runBackupArchiveSync := func() {
				if !cfg.GetBool("backup_archive_enabled", false) || !cfg.GetBool("backup_archive_auto_upload", false) {
					return
				}
				if !archiveManager.ShouldRunBackupSchedule(time.Now()) {
					return
				}
				uploaded, err := archiveManager.SyncPendingBackups(context.Background(), cfg.GetInt("backup_archive_batch_size", 2))
				if err != nil {
					_ = db.AddLog("WARN", "backup", fmt.Sprintf("Yedek arsiv senkronizasyonu basarisiz: %v", err))
					return
				}
				if uploaded > 0 {
					_ = db.AddLog("INFO", "backup", fmt.Sprintf("Yeni arsivlenen sistem yedegi sayisi: %d", uploaded))
				}
			}
			runBackupArchiveSync()
			for {
				select {
				case <-ctxDone:
					return
				case <-ticker.C:
					runBackupArchiveSync()
				}
			}
		}()
	}

	if cfg.GetBool("maintenance_auto_cleanup", true) {
		intervalHours := cfg.GetInt("maintenance_cleanup_interval", 6)
		if intervalHours <= 0 {
			intervalHours = 6
		}
		go func() {
			ticker := time.NewTicker(time.Duration(intervalHours) * time.Hour)
			defer ticker.Stop()
			runMaintenance := func() {
				if recManager != nil {
					retentionDays := cfg.GetInt("recordings_retention_days", cfg.GetInt("storage_auto_clean", 30))
					if retentionDays > 0 {
						if deleted, err := recManager.CleanupOld(time.Duration(retentionDays) * 24 * time.Hour); err == nil && deleted > 0 {
							_ = db.AddLog("INFO", "maintenance", fmt.Sprintf("Eski kayitlar temizlendi: %d dosya", deleted))
						}
					}
					keepLatest := cfg.GetInt("recordings_keep_latest", 10)
					if keepLatest > 0 {
						if deleted, err := recManager.TrimLatestPerStream(keepLatest); err == nil && deleted > 0 {
							_ = db.AddLog("INFO", "maintenance", fmt.Sprintf("Kayit trim uygulandi: %d dosya", deleted))
						}
					}
				}
				if days := cfg.GetInt("analytics_retention_days", 30); days > 0 {
					if deleted, err := db.CleanupAnalyticsSnapshots(time.Duration(days) * 24 * time.Hour); err == nil && deleted > 0 {
						_ = db.AddLog("INFO", "maintenance", fmt.Sprintf("Eski analytics snapshot temizlendi: %d", deleted))
					}
				}
				if days := cfg.GetInt("player_telemetry_retention_days", 30); days > 0 {
					if deleted, err := db.CleanupPlayerTelemetrySamples(time.Duration(days) * 24 * time.Hour); err == nil && deleted > 0 {
						_ = db.AddLog("INFO", "maintenance", fmt.Sprintf("Eski QoE telemetry ornekleri temizlendi: %d", deleted))
					}
				}
				if days := cfg.GetInt("track_analytics_retention_days", 30); days > 0 {
					if deleted, err := db.CleanupTrackTelemetrySamples(time.Duration(days) * 24 * time.Hour); err == nil && deleted > 0 {
						_ = db.AddLog("INFO", "maintenance", fmt.Sprintf("Eski track analytics ornekleri temizlendi: %d", deleted))
					}
				}
			}

			runMaintenance()
			for {
				select {
				case <-ctxDone:
					return
				case <-ticker.C:
					runMaintenance()
				}
			}
		}()
	}

	_ = dataDir
}

func buildHealthReport(cfg *config.Manager, db *storage.SQLiteDB, stats storage.ServerStats, tcManager *transcode.Manager, streamMgr *stream.Manager, playerTelemetry *playerTelemetryCollector, archiveManager *archive.Manager, dataDir string) map[string]interface{} {
	alerts := make([]systemAlert, 0, 8)
	status := "ok"
	qoeStreams := make([]map[string]interface{}, 0, 8)
	archiveSummary := archive.Summary{}
	if archiveManager != nil {
		archiveSummary = archiveManager.Summary()
	}

	if cfg.GetBool("alerts_enabled", true) {
		if cfg.GetBool("ssl_enabled", false) {
			webTLSSource, _ := tlsutil.NewSource(cfg, tlsutil.ProfileWeb, dataDir)
			if webTLSSource == nil || !webTLSSource.Ready {
				alerts = append(alerts, systemAlert{
					Level:       "warning",
					Code:        "ssl_missing",
					Title:       "SSL etkin ama sertifika hazir degil",
					Description: "Web HTTPS acik gorunuyor ancak sertifika profili tamamlanmadi.",
					Action:      "Ayarlar > SSL/TLS ekranindan web sertifikasini veya Let's Encrypt domainini tamamlayin ve restart edin.",
				})
			} else if webTLSSource.CertPath != "" {
				if expiry, ok := readCertificateExpiry(webTLSSource.CertPath); ok {
					daysLeft := int(time.Until(expiry).Hours() / 24)
					if daysLeft <= cfg.GetInt("alerts_cert_days", 21) {
						level := "warning"
						if daysLeft <= 7 {
							level = "critical"
						}
						alerts = append(alerts, systemAlert{
							Level:       level,
							Code:        "ssl_expiry",
							Title:       "SSL sertifikasi yakinda bitiyor",
							Description: fmt.Sprintf("Sertifika son gecerlilik tarihi: %s (%d gun kaldi)", expiry.Format("02.01.2006"), daysLeft),
							Action:      "Yeni sertifikayi yukleyin veya Let's Encrypt ayarlarini yenileyin.",
						})
					}
				}
			}
		}
		if cfg.GetBool("rtmps_enabled", false) {
			streamTLSSource, _ := tlsutil.NewSource(cfg, tlsutil.ProfileStream, dataDir)
			if streamTLSSource == nil || !streamTLSSource.Ready {
				alerts = append(alerts, systemAlert{
					Level:       "warning",
					Code:        "rtmps_ssl_missing",
					Title:       "RTMPS etkin ama stream sertifikasi hazir degil",
					Description: "OBS gibi encoder'lar RTMPS kullanacaksa stream SSL profili tamamlanmalidir.",
					Action:      "Ayarlar > SSL/TLS ekranindan stream sertifikasini veya stream Let's Encrypt domainini tamamlayin.",
				})
			} else if streamTLSSource.CertPath != "" {
				if expiry, ok := readCertificateExpiry(streamTLSSource.CertPath); ok {
					daysLeft := int(time.Until(expiry).Hours() / 24)
					if daysLeft <= cfg.GetInt("alerts_cert_days", 21) {
						level := "warning"
						if daysLeft <= 7 {
							level = "critical"
						}
						alerts = append(alerts, systemAlert{
							Level:       level,
							Code:        "rtmps_ssl_expiry",
							Title:       "RTMPS sertifikasi yakinda bitiyor",
							Description: fmt.Sprintf("Stream sertifika son gecerlilik tarihi: %s (%d gun kaldi)", expiry.Format("02.01.2006"), daysLeft),
							Action:      "Yeni stream sertifikasini yukleyin veya stream Let's Encrypt ayarlarini yenileyin.",
						})
					}
				}
			}
		}

		if threshold := cfg.GetInt("alerts_memory_threshold_mb", 2048); threshold > 0 && stats.MemoryUsedMB >= int64(threshold) {
			alerts = append(alerts, systemAlert{
				Level:       "warning",
				Code:        "memory_high",
				Title:       "Bellek kullanimi yuksek",
				Description: fmt.Sprintf("Sunucu %d MB bellek kullaniyor, esik %d MB.", stats.MemoryUsedMB, threshold),
				Action:      "Transcode profil sayisini veya kalite merdivenini azaltin.",
			})
		}

		if maxGB := cfg.GetInt("storage_max_gb", 50); maxGB > 0 {
			recordingsSizeGB := float64(folderSize(filepath.Join(dataDir, "recordings"))) / (1024 * 1024 * 1024)
			thresholdGB := float64(maxGB - cfg.GetInt("alerts_disk_threshold_gb", 5))
			if recordingsSizeGB >= thresholdGB {
				alerts = append(alerts, systemAlert{
					Level:       "warning",
					Code:        "storage_high",
					Title:       "Kayit depolamasi limite yaklasiyor",
					Description: fmt.Sprintf("Kayit klasoru %.1f GB kullaniyor. Yapilandirilmis limit %d GB.", recordingsSizeGB, maxGB),
					Action:      "Retention suresini kisaltin veya eski kayitlari silin.",
				})
			}
		}

		if ffmpegPath, err := tcManager.DetectFFmpeg(); err != nil || ffmpegPath == "" {
			alerts = append(alerts, systemAlert{
				Level:       "critical",
				Code:        "ffmpeg_missing",
				Title:       "FFmpeg bulunamadi",
				Description: "Transcode, MP4/WebM, audio output ve ABR zinciri FFmpeg olmadan tam calismaz.",
				Action:      "Portable paketteki ffmpeg klasorunun exe ile ayni dizinde oldugunu dogrulayin.",
			})
		}

		if !cfg.GetBool("abr_enabled", false) {
			alerts = append(alerts, systemAlert{
				Level:       "info",
				Code:        "abr_disabled",
				Title:       "Adaptif bitrate kapali",
				Description: "Su an canli HLS cikisi tek kalite olarak uretiliyor.",
				Action:      "Teslimat / ABR ekranindan adaptif bitrate'i acabilirsiniz.",
			})
		}
		if cfg.GetBool("archive_enabled", false) {
			switch {
			case !archiveSummary.RecordingConfigured:
				alerts = append(alerts, systemAlert{
					Level:       "warning",
					Code:        "archive_config_invalid",
					Title:       "Arsivleme acik ama nesne depolama ayari eksik",
					Description: "Archive / object storage akisi etkinlestirilmis gorunuyor ancak saglayici ayarlari tamamlanmadi.",
					Action:      "Depolama ekranindan provider, endpoint veya bucket bilgilerini dogrulayin.",
				})
			case archiveSummary.ErrorItems > 0:
				alerts = append(alerts, systemAlert{
					Level:       "warning",
					Code:        "archive_errors_present",
					Title:       "Arsiv hatalari var",
					Description: fmt.Sprintf("Kayit arsiv kutugunde %d hata durumundaki oge gorunuyor.", archiveSummary.ErrorItems),
					Action:      "Kayitlar ekranindaki Arsiv Kutuphanesi tablosunu acip hatali ogeleri yeniden deneyin.",
				})
			}
		}
		if cfg.GetBool("backup_archive_enabled", false) {
			switch {
			case !archiveSummary.BackupConfigured:
				alerts = append(alerts, systemAlert{
					Level:       "warning",
					Code:        "backup_archive_config_invalid",
					Title:       "Yedek arsivi acik ama hedef ayari eksik",
					Description: "Sistem yedeklerini object storage veya SFTP hedefine tasimak etkin gorunuyor ancak hedef ayarlari tamamlanmadi.",
					Action:      "Depolama ve Arsiv Merkezi ekranindan provider ve hedef bilgilerini dogrulayin.",
				})
			case archiveSummary.BackupErrorItems > 0:
				alerts = append(alerts, systemAlert{
					Level:       "warning",
					Code:        "backup_archive_errors_present",
					Title:       "Yedek arsiv hatalari var",
					Description: fmt.Sprintf("Yedek arsiv kutugunde %d hata durumundaki oge gorunuyor.", archiveSummary.BackupErrorItems),
					Action:      "Depolama ve Arsiv Merkezi ekranindaki Yedek Arsiv Kutuphanesi tablosundan hatali ogeleri yeniden deneyin.",
				})
			}
		}

		for _, item := range collectRuntimeObservability(streamMgr, tcManager, playerTelemetry) {
			streamAlerts := buildQoEAlerts(cfg, item.StreamName, item.Telemetry)
			if len(streamAlerts) > 0 {
				alerts = append(alerts, streamAlerts...)
			}
			qoeStreams = append(qoeStreams, map[string]interface{}{
				"stream_key":                item.StreamKey,
				"stream_name":               item.StreamName,
				"active_sessions":           item.Telemetry.ActiveSessions,
				"waiting_sessions":          item.Telemetry.WaitingSessions,
				"offline_sessions":          item.Telemetry.OfflineSessions,
				"average_buffer_seconds":    item.Telemetry.AverageBufferSeconds,
				"total_stalls":              item.Telemetry.TotalStalls,
				"total_quality_transitions": item.Telemetry.TotalQualityTransitions,
				"total_audio_switches":      item.Telemetry.TotalAudioSwitches,
				"dominant_quality":          dominantTelemetryLabel(item.Telemetry.Qualities),
				"dominant_audio":            dominantTelemetryLabel(item.Telemetry.AudioTracks),
				"qoe_alert_count":           len(streamAlerts),
				"active_audio_track_id":     item.Tracks.ActiveAudioTrackID,
				"active_video_track_id":     item.Tracks.ActiveVideoTrackID,
				"video_track_count":         len(item.Tracks.VideoTracks),
				"audio_track_count":         len(item.Tracks.AudioTracks),
			})
		}
	}

	for _, alert := range alerts {
		switch alert.Level {
		case "critical":
			status = "critical"
		case "warning":
			if status != "critical" {
				status = "warning"
			}
		}
	}

	services := map[string]bool{
		"http":     true,
		"https":    cfg.GetBool("ssl_enabled", false),
		"rtmp":     cfg.GetBool("rtmp_enabled", true),
		"dash":     cfg.GetBool("dash_enabled", false),
		"hls":      cfg.GetBool("hls_enabled", true),
		"ll_hls":   cfg.GetBool("hls_ll_enabled", false),
		"http_flv": cfg.GetBool("httpflv_enabled", false),
		"whep":     cfg.GetBool("whep_enabled", false),
		"abr":      cfg.GetBool("abr_enabled", false),
	}

	snapshots, _ := db.GetAnalyticsSnapshots(12)

	return map[string]interface{}{
		"status":       status,
		"generated_at": time.Now(),
		"alerts":       alerts,
		"services":     services,
		"storage": map[string]interface{}{
			"recordings_bytes": folderSize(filepath.Join(dataDir, "recordings")),
			"backups_bytes":    folderSize(filepath.Join(dataDir, "backups")),
			"hls_bytes":        folderSize(filepath.Join(dataDir, "hls")),
			"dash_bytes":       folderSize(filepath.Join(dataDir, "dash")),
			"archive":          archiveSummary,
		},
		"snapshots":       snapshots,
		"qoe_streams":     qoeStreams,
		"recommendations": buildRecommendations(cfg, alerts),
	}
}

func buildRecommendations(cfg *config.Manager, alerts []systemAlert) []string {
	recs := []string{}
	if !cfg.GetBool("abr_enabled", false) {
		recs = append(recs, "Canli izleyicide otomatik kalite degisimi icin ABR'yi acin.")
	}
	if cfg.GetBool("token_enabled", false) {
		recs = append(recs, "Public embed kullaniyorsaniz token dagitimini istemci uygulamanizla birlikte planlayin.")
	}
	if cfg.GetBool("maintenance_auto_cleanup", true) {
		recs = append(recs, "Kayit retention suresi aktif; disk dolmadan once temizlik otomatik yapilir.")
	}
	if len(alerts) == 0 {
		recs = append(recs, "Kritik bir sorun gorunmuyor. Sistem su an canli dagitim icin hazir.")
	}
	return recs
}

func dominantTelemetryLabel(values map[string]int) string {
	bestLabel := "-"
	bestValue := 0
	for label, value := range values {
		if value > bestValue && strings.TrimSpace(label) != "" {
			bestLabel = label
			bestValue = value
		}
	}
	return bestLabel
}

func readCertificateExpiry(certPath string) (time.Time, bool) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return time.Time{}, false
	}
	for len(data) > 0 {
		var block *pem.Block
		block, data = pem.Decode(data)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		return cert.NotAfter, true
	}
	return time.Time{}, false
}

func folderSize(path string) int64 {
	var size int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}
