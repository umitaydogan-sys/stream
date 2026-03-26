package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/analytics"
	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/security"
	"github.com/fluxstream/fluxstream/internal/storage"
	streampolicy "github.com/fluxstream/fluxstream/internal/stream"
	"github.com/fluxstream/fluxstream/internal/transcode"
	"github.com/fluxstream/fluxstream/internal/web"
)

func registerStudioAdminRoutes(
	webServer *web.Server,
	db *storage.SQLiteDB,
	cfg *config.Manager,
	analyticsTracker *analytics.Tracker,
	tcManager *transcode.Manager,
	playerTelemetry *playerTelemetryCollector,
	tokenMgr *security.TokenManager,
	dataDir string,
) {
	if webServer == nil || db == nil || cfg == nil {
		return
	}

	webServer.RegisterAdminHandler("/api/admin/assets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			category := studioAssetCategory(r.URL.Query().Get("category"))
			items, err := listStudioAssets(dataDir, category)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonResp(w, map[string]interface{}{"items": items})
		case http.MethodPost:
			if err := r.ParseMultipartForm(16 << 20); err != nil {
				http.Error(w, "Dosya yukleme istegi okunamadi", 400)
				return
			}
			category := studioAssetCategory(r.FormValue("category"))
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "Yuklenecek dosya bulunamadi", 400)
				return
			}
			defer file.Close()
			item, err := saveStudioAsset(dataDir, category, header.Filename, file)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonResp(w, map[string]interface{}{"success": true, "item": item})
		case http.MethodDelete:
			var req struct {
				Path string `json:"path"`
			}
			if err := decodeJSON(r, &req); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if err := deleteStudioAsset(dataDir, req.Path); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			jsonResp(w, map[string]interface{}{"success": true})
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	webServer.RegisterAdminHandler("/api/admin/embed-profiles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			items, err := db.ListEmbedProfiles(strings.TrimSpace(r.URL.Query().Get("stream_key")))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			if items == nil {
				items = []storage.EmbedProfile{}
			}
			jsonResp(w, items)
		case http.MethodPost:
			var item storage.EmbedProfile
			if err := decodeJSON(r, &item); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if strings.TrimSpace(item.Name) == "" {
				http.Error(w, "Profil adi gerekli", 400)
				return
			}
			if item.Width <= 0 {
				item.Width = 1280
			}
			if item.Height <= 0 {
				item.Height = 720
			}
			if strings.TrimSpace(item.Mode) == "" {
				item.Mode = "simple"
			}
			if strings.TrimSpace(item.PrimaryFormat) == "" {
				item.PrimaryFormat = "player"
			}
			if strings.TrimSpace(item.Theme) == "" {
				item.Theme = "clean"
			}
			id, err := db.CreateEmbedProfile(&item)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			created, _ := db.GetEmbedProfile(id)
			jsonResp(w, map[string]interface{}{"success": true, "item": created})
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	webServer.RegisterAdminHandler("/api/admin/embed-profiles/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/embed-profiles/")
		id, _ := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if id <= 0 {
			http.Error(w, "Profil bulunamadi", 404)
			return
		}
		switch r.Method {
		case http.MethodGet:
			item, err := db.GetEmbedProfile(id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			if item == nil {
				http.Error(w, "Profil bulunamadi", 404)
				return
			}
			jsonResp(w, item)
		case http.MethodPut:
			item, err := db.GetEmbedProfile(id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			if item == nil {
				http.Error(w, "Profil bulunamadi", 404)
				return
			}
			var payload storage.EmbedProfile
			if err := decodeJSON(r, &payload); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			payload.ID = id
			if payload.Width <= 0 {
				payload.Width = item.Width
			}
			if payload.Height <= 0 {
				payload.Height = item.Height
			}
			if err := db.UpdateEmbedProfile(&payload); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			updated, _ := db.GetEmbedProfile(id)
			jsonResp(w, map[string]interface{}{"success": true, "item": updated})
		case http.MethodDelete:
			if err := db.DeleteEmbedProfile(id); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonResp(w, map[string]interface{}{"success": true})
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	webServer.RegisterAdminHandler("/api/admin/abr-profiles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			items, err := db.ListABRProfiles(strings.TrimSpace(r.URL.Query().Get("stream_key")))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			if items == nil {
				items = []storage.ABRProfile{}
			}
			jsonResp(w, items)
		case http.MethodPost:
			var item storage.ABRProfile
			if err := decodeJSON(r, &item); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.ProfileSet) == "" {
				http.Error(w, "Profil adi ve anahtari gerekli", 400)
				return
			}
			if strings.TrimSpace(item.Scope) == "" {
				item.Scope = "global"
			}
			item.ProfileSet = studioSlug(item.ProfileSet, "custom-profile")
			id, err := db.CreateABRProfile(&item)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			created, _ := db.GetABRProfile(id)
			jsonResp(w, map[string]interface{}{"success": true, "item": created})
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	webServer.RegisterAdminHandler("/api/admin/abr-profiles/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/abr-profiles/")
		id, _ := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if id <= 0 {
			http.Error(w, "Profil bulunamadi", 404)
			return
		}
		switch r.Method {
		case http.MethodGet:
			item, err := db.GetABRProfile(id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			if item == nil {
				http.Error(w, "Profil bulunamadi", 404)
				return
			}
			jsonResp(w, item)
		case http.MethodPut:
			existing, err := db.GetABRProfile(id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			if existing == nil {
				http.Error(w, "Profil bulunamadi", 404)
				return
			}
			var payload storage.ABRProfile
			if err := decodeJSON(r, &payload); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			payload.ID = id
			if strings.TrimSpace(payload.ProfileSet) == "" {
				payload.ProfileSet = existing.ProfileSet
			}
			payload.ProfileSet = studioSlug(payload.ProfileSet, existing.ProfileSet)
			if strings.TrimSpace(payload.Scope) == "" {
				payload.Scope = existing.Scope
			}
			if err := db.UpdateABRProfile(&payload); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			updated, _ := db.GetABRProfile(id)
			jsonResp(w, map[string]interface{}{"success": true, "item": updated})
		case http.MethodDelete:
			if err := db.DeleteABRProfile(id); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonResp(w, map[string]interface{}{"success": true})
		default:
			http.Error(w, "Method not allowed", 405)
		}
	})

	webServer.RegisterAdminHandler("/api/admin/abr-profiles/apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		var req struct {
			ProfileID int64  `json:"profile_id"`
			StreamKey string `json:"stream_key"`
			Scope     string `json:"scope"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		profile, err := db.GetABRProfile(req.ProfileID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if profile == nil {
			http.Error(w, "Profil bulunamadi", 404)
			return
		}
		scope := strings.ToLower(strings.TrimSpace(req.Scope))
		if scope == "" {
			scope = strings.ToLower(strings.TrimSpace(profile.Scope))
		}
		if scope == "" {
			scope = "global"
		}
		if scope == "stream" {
			streamKey := strings.TrimSpace(req.StreamKey)
			if streamKey == "" {
				streamKey = strings.TrimSpace(profile.StreamKey)
			}
			st, err := db.GetStreamByKey(streamKey)
			if err != nil || st == nil {
				http.Error(w, "Stream bulunamadi", 404)
				return
			}
			policy := streampolicy.ParsePolicyJSON(st.PolicyJSON)
			policy.EnableABR = true
			policy.ProfileSet = profile.ProfileSet
			st.PolicyJSON = streampolicy.EncodePolicyJSON(policy)
			if err := db.UpdateStream(st); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonResp(w, map[string]interface{}{"success": true, "mode": "stream", "stream_key": st.StreamKey, "profile_set": profile.ProfileSet})
			return
		}

		allProfiles := parseABRProfilesMap(cfg.Get("abr_profiles_json", "{}"))
		layers := parseABRProfileLayers(profile.ProfilesJSON)
		allProfiles[profile.ProfileSet] = layers
		encoded, _ := json.Marshal(allProfiles)
		if err := cfg.Set("abr_profiles_json", string(encoded), "outputs"); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := cfg.Set("abr_profile_set", profile.ProfileSet, "outputs"); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := cfg.Set("abr_enabled", "true", "outputs"); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonResp(w, map[string]interface{}{"success": true, "mode": "global", "profile_set": profile.ProfileSet})
	})

	webServer.RegisterAdminHandler("/api/admin/abr-profiles/direct-apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		var req struct {
			ProfileSet  string `json:"profile_set"`
			ProfilesJSON string `json:"profiles_json"`
			StreamKey   string `json:"stream_key"`
			Scope       string `json:"scope"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		req.ProfileSet = studioSlug(req.ProfileSet, "custom-profile")
		layers := parseABRProfileLayers(req.ProfilesJSON)
		scope := strings.ToLower(strings.TrimSpace(req.Scope))
		if scope == "stream" {
			st, err := db.GetStreamByKey(strings.TrimSpace(req.StreamKey))
			if err != nil || st == nil {
				http.Error(w, "Stream bulunamadi", 404)
				return
			}
			policy := streampolicy.ParsePolicyJSON(st.PolicyJSON)
			policy.EnableABR = true
			policy.ProfileSet = req.ProfileSet
			st.PolicyJSON = streampolicy.EncodePolicyJSON(policy)
			if err := db.UpdateStream(st); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			jsonResp(w, map[string]interface{}{"success": true, "mode": "stream", "profile_set": req.ProfileSet, "stream_key": st.StreamKey, "layers": layers})
			return
		}
		allProfiles := parseABRProfilesMap(cfg.Get("abr_profiles_json", "{}"))
		allProfiles[req.ProfileSet] = layers
		encoded, _ := json.Marshal(allProfiles)
		if err := cfg.Set("abr_profiles_json", string(encoded), "outputs"); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := cfg.Set("abr_profile_set", req.ProfileSet, "outputs"); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := cfg.Set("abr_enabled", "true", "outputs"); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonResp(w, map[string]interface{}{"success": true, "mode": "global", "profile_set": req.ProfileSet, "layers": layers})
	})

	webServer.RegisterAdminHandler("/api/admin/streams/adaptive-mode", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		var req struct {
			StreamID   int64  `json:"stream_id"`
			StreamKey  string `json:"stream_key"`
			ProfileSet string `json:"profile_set"`
			ApplyMode  string `json:"apply_mode"`
			SyncMode   *bool  `json:"sync_mode"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		var (
			st  *storage.Stream
			err error
		)
		if req.StreamID > 0 {
			st, err = db.GetStreamByID(req.StreamID)
		} else {
			st, err = db.GetStreamByKey(strings.TrimSpace(req.StreamKey))
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if st == nil {
			http.Error(w, "Stream bulunamadi", 404)
			return
		}

		policy := streampolicy.ParsePolicyJSON(st.PolicyJSON)
		profileSet := strings.TrimSpace(strings.ToLower(req.ProfileSet))
		if profileSet == "" {
			profileSet = strings.TrimSpace(strings.ToLower(policy.ProfileSet))
		}
		if profileSet == "" {
			profileSet = "balanced"
		}
		applyMode := strings.TrimSpace(strings.ToLower(req.ApplyMode))
		if applyMode != "live_now" {
			applyMode = "next_publish"
		}
		syncMode := true
		if req.SyncMode != nil {
			syncMode = *req.SyncMode
		}

		policy.EnableABR = true
		policy.ProfileSet = profileSet
		if syncMode {
			switch profileSet {
			case "balanced", "mobile", "resilient", "radio":
				policy.Mode = profileSet
			}
		}
		st.PolicyJSON = streampolicy.EncodePolicyJSON(policy)
		if err := db.UpdateStream(st); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		liveApplied := false
		warnings := []string{}
		message := "Adaptive teslimat sonraki yayin icin isaretlendi."
		if applyMode == "live_now" {
			if !strings.EqualFold(st.Status, "live") {
				warnings = append(warnings, "Yayin su an canli degil; ayar bir sonraki publish aninda devreye girecek.")
				applyMode = "next_publish"
			} else if tcManager == nil {
				warnings = append(warnings, "Transcode yoneticisi hazir degil; ayar bir sonraki publish aninda devreye girecek.")
				applyMode = "next_publish"
			} else {
				opts := buildLiveOptionsFromConfig(cfg)
				opts.ABREnabled = true
				opts.ProfileSet = profileSet
				opts.ProfilesJSON = cfg.Get("abr_profiles_json", "")
				opts.Profiles = transcode.ResolveProfiles(opts.ProfileSet, opts.ProfilesJSON)
				if policy.DefaultVideoTrackID > 0 && policy.DefaultVideoTrackID <= 255 {
					opts.DefaultVideoTrackID = uint8(policy.DefaultVideoTrackID)
				}
				if policy.DefaultAudioTrackID > 0 && policy.DefaultAudioTrackID <= 255 {
					opts.DefaultAudioTrackID = uint8(policy.DefaultAudioTrackID)
				}
				tcManager.SetStreamLiveOptions(st.StreamKey, opts)
				tcManager.StopLiveDASH(st.StreamKey)
				tcManager.StopLiveHLS(st.StreamKey)
				hlsStarted := false
				if _, err := tcManager.StartLiveHLS(st.StreamKey); err != nil {
					warnings = append(warnings, "Canli HLS yeniden baslatilamadi: "+err.Error())
				} else {
					hlsStarted = true
				}
				if cfg.GetBool("dash_enabled", false) && cfg.GetBool("transcode_live_dash_enabled", true) {
					if _, err := tcManager.StartLiveDASH(st.StreamKey); err != nil {
						warnings = append(warnings, "Canli DASH yeniden baslatilamadi: "+err.Error())
					}
				}
				liveApplied = hlsStarted
				if liveApplied {
					message = "Adaptive teslimat canli yayina uygulandi. Oynaticida kisa bir yeniden kurulum etkisi olabilir."
				}
			}
		}

		jsonResp(w, map[string]interface{}{
			"success":      true,
			"stream_id":    st.ID,
			"stream_key":   st.StreamKey,
			"profile_set":  profileSet,
			"apply_mode":   applyMode,
			"live_applied": liveApplied,
			"warnings":     warnings,
			"policy_json":  st.PolicyJSON,
			"message":      message,
		})
	})

	webServer.RegisterAdminHandler("/api/admin/analytics/center", func(w http.ResponseWriter, r *http.Request) {
		streamKey := strings.TrimSpace(r.URL.Query().Get("stream_key"))
		period := strings.TrimSpace(r.URL.Query().Get("period"))
		if period == "" {
			period = "24h"
		}
		mode := strings.TrimSpace(r.URL.Query().Get("mode"))
		if mode == "" {
			mode = "live"
		}

		streams, _ := db.GetAllStreams()
		dashboard := analytics.Dashboard{}
		if analyticsTracker != nil {
			dashboard = analyticsTracker.GetDashboard()
		}
		window := analyticsWindowForPeriod(period, time.Now())
		snapshots, _ := db.GetAnalyticsSnapshotsSince(window.Since, 0)
		if snapshots == nil {
			snapshots = []storage.AnalyticsSnapshot{}
		}
		history := buildAnalyticsHistoryPayload(window.Period, snapshots, time.Now())
		viewerSessions := []analytics.ViewerSession{}
		if analyticsTracker != nil {
			viewerSessions = analyticsTracker.GetViewerSessions()
		}

		streamSummaries := make([]map[string]interface{}, 0, len(streams))
		risky := make([]map[string]interface{}, 0, len(streams))
		var topStreamName string
		var topStreamViewers int

		for _, st := range streams {
			if analyticsTracker != nil {
				analyticsTracker.RegisterStreamName(st.StreamKey, st.Name)
			}
			telemetry := playerTelemetry.Snapshot(st.StreamKey)
			qoeAlerts := buildQoEAlerts(cfg, st.Name, telemetry)
			trackSnapshot := transcode.LiveTrackSnapshot{}
			if tcManager != nil {
				trackSnapshot = tcManager.GetLiveTrackSnapshot(st.StreamKey)
			}
			currentViewers := st.ViewerCount
			streamStats := (*analytics.StreamStats)(nil)
			if analyticsTracker != nil {
				streamStats = analyticsTracker.GetStreamStats(st.StreamKey)
				if streamStats != nil && int(streamStats.CurrentViewers) > currentViewers {
					currentViewers = streamStats.CurrentViewers
				}
			}
			summary := map[string]interface{}{
				"id":               st.ID,
				"stream_key":       st.StreamKey,
				"name":             st.Name,
				"status":           st.Status,
				"viewer_count":     currentViewers,
				"input_bitrate":    st.InputBitrate,
				"input_codec":      st.InputCodec,
				"input_width":      st.InputWidth,
				"input_height":     st.InputHeight,
				"input_fps":        st.InputFPS,
				"telemetry":        telemetry,
				"qoe_alerts":       qoeAlerts,
				"tracks":           trackSnapshot,
				"analytics":        streamStats,
				"health_score":     studioHealthScore(telemetry, qoeAlerts),
				"recommended_mode": studioRecommendedMode(st, telemetry),
			}
			streamSummaries = append(streamSummaries, summary)
			if currentViewers > topStreamViewers {
				topStreamViewers = currentViewers
				topStreamName = st.Name
			}
			if len(qoeAlerts) > 0 {
				risky = append(risky, map[string]interface{}{
					"id":           st.ID,
					"stream_key":   st.StreamKey,
					"name":         st.Name,
					"status":       st.Status,
					"viewer_count": currentViewers,
					"alerts":       qoeAlerts,
				})
			}
		}

		selected := map[string]interface{}{}
		if streamKey != "" {
			if st, err := db.GetStreamByKey(streamKey); err == nil && st != nil {
				playerHistory, _ := db.GetPlayerTelemetrySamples(st.StreamKey, 96)
				trackHistory, _ := db.GetTrackTelemetrySamples(st.StreamKey, 200)
				selected["stream"] = st
				if analyticsTracker != nil {
					selected["analytics"] = analyticsTracker.GetStreamStats(st.StreamKey)
				}
				selected["telemetry"] = playerTelemetry.Snapshot(st.StreamKey)
				selected["history"] = playerHistory
				selected["track_history"] = trackHistory
				if tcManager != nil {
					selected["tracks"] = tcManager.GetLiveTrackSnapshot(st.StreamKey)
				} else {
					selected["tracks"] = transcode.LiveTrackSnapshot{}
				}
				selected["qoe_alerts"] = buildQoEAlerts(cfg, st.Name, playerTelemetry.Snapshot(st.StreamKey))
			}
		}

		aggregate := aggregateTelemetrySummaries(streamSummaries, streamKey)
		errorRate := 0.0
		if aggregate.ActiveSessions > 0 {
			errorRate = float64(aggregate.OfflineSessions+aggregate.WaitingSessions) / float64(aggregate.ActiveSessions) * 100.0
		}
		jsonResp(w, map[string]interface{}{
			"mode":            mode,
			"period":          period,
			"streams":         streamSummaries,
			"dashboard":       dashboard,
			"history":         history,
			"viewer_sessions": viewerSessions,
			"risky_streams":   risky,
			"selected":        selected,
			"kpis": map[string]interface{}{
				"active_viewers":      aggregate.ActiveSessions,
				"peak_viewers":        dashboard.PeakConcurrent,
				"average_buffer":      aggregate.AverageBufferSeconds,
				"stalls":              aggregate.TotalStalls,
				"quality_transitions": aggregate.TotalQualityTransitions,
				"audio_switches":      aggregate.TotalAudioSwitches,
				"error_rate":          math.Round(errorRate*10) / 10,
				"top_stream":          topStreamName,
			},
		})
	})

	webServer.RegisterAdminHandler("/api/admin/security/playback-link", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		var req struct {
			StreamKey        string                 `json:"stream_key"`
			Page             string                 `json:"page"`
			Format           string                 `json:"format"`
			Width            int                    `json:"width"`
			Height           int                    `json:"height"`
			Autoplay         bool                   `json:"autoplay"`
			Muted            bool                   `json:"muted"`
			Options          map[string]interface{} `json:"options"`
			Security         map[string]interface{} `json:"security"`
			ApplyStreamPolicy bool                  `json:"apply_stream_policy"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		streamKey := strings.TrimSpace(req.StreamKey)
		if streamKey == "" {
			http.Error(w, "Stream key gerekli", 400)
			return
		}
		st, err := db.GetStreamByKey(streamKey)
		if err != nil || st == nil {
			http.Error(w, "Stream bulunamadi", 404)
			return
		}
		page := strings.TrimSpace(strings.ToLower(req.Page))
		if page == "" {
			page = "embed"
		}
		format := strings.TrimSpace(strings.ToLower(req.Format))
		if format == "" {
			format = "player"
		}
		base := studioPublicBaseURL(r, cfg)
		previewBase := studioPreviewBaseURL(r)

		secureEnabled := asBool(req.Security["signed_url"]) || asBool(req.Security["token_required"]) || asBool(req.Security["session_bound"])
		expiryMinutes := asInt(req.Security["expiry_minutes"], cfg.GetInt("token_duration", 60))
		allowedIP := strings.TrimSpace(asString(req.Security["ip_restriction"]))
		allowedDomain := strings.TrimSpace(asString(req.Security["domain_restriction"]))
		watermark := strings.TrimSpace(asString(req.Security["watermark"]))
		viewerID := strings.TrimSpace(asString(req.Security["viewer_id"]))
		if asBool(req.Security["session_bound"]) && viewerID == "" {
			viewerID = randomTokenFragment(8)
		}
		claims := security.PlaybackTokenClaims{
			AllowedIP:     allowedIP,
			AllowedDomain: allowedDomain,
			ViewerID:      viewerID,
			Watermark:     watermark,
			AllowedFormat: format,
			ExpiresAtUnix: time.Now().Add(time.Duration(maxInt(expiryMinutes, 1)) * time.Minute).Unix(),
		}
		token := ""
		expiry := time.Time{}
		if secureEnabled {
			generated, exp, err := tokenMgr.GeneratePlaybackToken(streamKey, claims)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			token = generated
			expiry = exp
		}

		playerURL := fmt.Sprintf("%s/play/%s", base, streamKey)
		embedURL := fmt.Sprintf("%s/embed/%s", base, streamKey)
		if page == "player" {
			embedURL = playerURL
		}
		manifestURL := studioManifestURL(base, streamKey, format)
		audioURL := studioManifestURL(base, streamKey, "dash_audio")
		query := map[string]string{}
		if format != "" && format != "player" && format != "embed" && format != "iframe" {
			query["format"] = format
		}
		if token != "" {
			query["token"] = token
		}
		if viewerID != "" {
			query["viewer_id"] = viewerID
		}
		if watermark != "" {
			query["player_watermark"] = watermark
		}
		if req.Autoplay {
			query["autoplay"] = "1"
		}
		if req.Muted {
			query["muted"] = "1"
		}
		playerURL = withQueryMap(playerURL, query)
		embedURL = withQueryMap(embedURL, query)
		if manifestURL != "" {
			manifestURL = withQueryMap(manifestURL, map[string]string{
				"token":     token,
				"viewer_id": viewerID,
			})
		}
		if audioURL != "" {
			audioURL = withQueryMap(audioURL, map[string]string{
				"token":     token,
				"viewer_id": viewerID,
			})
		}

		if req.ApplyStreamPolicy {
			policy := streampolicy.ParsePolicyJSON(st.PolicyJSON)
			policy.RequirePlaybackToken = secureEnabled || policy.RequirePlaybackToken
			policy.RequireSignedURL = asBool(req.Security["signed_url"]) || policy.RequireSignedURL
			st.PolicyJSON = streampolicy.EncodePolicyJSON(policy)
			if allowedDomain != "" {
				st.DomainLock = allowedDomain
			}
			if err := db.UpdateStream(st); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		width := req.Width
		height := req.Height
		if width <= 0 {
			width = 1280
		}
		if height <= 0 {
			height = 720
		}
		iframeCode := fmt.Sprintf(`<iframe src="%s" width="%d" height="%d" frameborder="0" allowfullscreen></iframe>`, embedURL, width, height)
		scriptCode := fmt.Sprintf(`<script src="%s/static/vendor/hls.min.js"></script><div data-fluxstream="%s"></div><script>window.open(%q,'_blank','width=%d,height=%d');</script>`, base, streamKey, playerURL, maxInt(width, 640), maxInt(height, 360))

		jsonResp(w, map[string]interface{}{
			"success":     true,
			"stream_key":  streamKey,
			"token":       token,
			"viewer_id":   viewerID,
			"expires_at":  expiry,
			"player_url":  playerURL,
			"embed_url":   embedURL,
			"manifest_url": manifestURL,
			"audio_url":   audioURL,
			"vlc_url":     studioManifestURL(base, streamKey, "hls"),
			"preview_player_url": withQueryMap(fmt.Sprintf("%s/play/%s", previewBase, streamKey), query),
			"preview_embed_url": withQueryMap(func() string {
				if page == "player" {
					return fmt.Sprintf("%s/play/%s", previewBase, streamKey)
				}
				return fmt.Sprintf("%s/embed/%s", previewBase, streamKey)
			}(), query),
			"preview_manifest_url": func() string {
				url := studioManifestURL(previewBase, streamKey, format)
				if url == "" {
					return ""
				}
				return withQueryMap(url, map[string]string{
					"token":     token,
					"viewer_id": viewerID,
				})
			}(),
			"preview_audio_url": func() string {
				url := studioManifestURL(previewBase, streamKey, "dash_audio")
				if url == "" {
					return ""
				}
				return withQueryMap(url, map[string]string{
					"token":     token,
					"viewer_id": viewerID,
				})
			}(),
			"iframe_code": iframeCode,
			"script_code": scriptCode,
			"applied":     req.ApplyStreamPolicy,
		})
	})

	webServer.RegisterAdminHandler("/api/admin/security/stream-policy/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", 405)
			return
		}
		var req struct {
			StreamKey         string `json:"stream_key"`
			ClearDomainLock   bool   `json:"clear_domain_lock"`
			ClearIPWhitelist  bool   `json:"clear_ip_whitelist"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		streamKey := strings.TrimSpace(req.StreamKey)
		if streamKey == "" {
			http.Error(w, "Stream key gerekli", 400)
			return
		}
		st, err := db.GetStreamByKey(streamKey)
		if err != nil || st == nil {
			http.Error(w, "Stream bulunamadi", 404)
			return
		}
		policy := streampolicy.ParsePolicyJSON(st.PolicyJSON)
		policy.RequirePlaybackToken = false
		policy.RequireSignedURL = false
		st.PolicyJSON = streampolicy.EncodePolicyJSON(policy)
		if req.ClearDomainLock {
			st.DomainLock = ""
		}
		if req.ClearIPWhitelist {
			st.IPWhitelist = ""
		}
		if err := db.UpdateStream(st); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		jsonResp(w, map[string]interface{}{
			"success":        true,
			"stream_key":     st.StreamKey,
			"policy_json":    st.PolicyJSON,
			"domain_lock":    st.DomainLock,
			"ip_whitelist":   st.IPWhitelist,
		})
	})
}

type aggregatedTelemetrySummary struct {
	ActiveSessions          int
	WaitingSessions         int
	OfflineSessions         int
	TotalStalls             int64
	TotalQualityTransitions int64
	TotalAudioSwitches      int64
	AverageBufferSeconds    float64
}

func aggregateTelemetrySummaries(items []map[string]interface{}, selectedStreamKey string) aggregatedTelemetrySummary {
	out := aggregatedTelemetrySummary{}
	bufferSum := 0.0
	bufferCount := 0
	for _, item := range items {
		if selectedStreamKey != "" && asString(item["stream_key"]) != selectedStreamKey {
			continue
		}
		raw, _ := item["telemetry"].(playerTelemetrySnapshot)
		out.ActiveSessions += raw.ActiveSessions
		out.WaitingSessions += raw.WaitingSessions
		out.OfflineSessions += raw.OfflineSessions
		out.TotalStalls += raw.TotalStalls
		out.TotalQualityTransitions += raw.TotalQualityTransitions
		out.TotalAudioSwitches += raw.TotalAudioSwitches
		if raw.ActiveSessions > 0 {
			bufferSum += raw.AverageBufferSeconds
			bufferCount++
		}
	}
	if bufferCount > 0 {
		out.AverageBufferSeconds = bufferSum / float64(bufferCount)
	}
	return out
}

func studioHealthScore(snapshot playerTelemetrySnapshot, alerts []systemAlert) int {
	score := 100
	score -= int(snapshot.TotalStalls) * 3
	score -= snapshot.WaitingSessions * 4
	score -= snapshot.OfflineSessions * 6
	score -= len(alerts) * 8
	if snapshot.AverageBufferSeconds > 0 && snapshot.AverageBufferSeconds < 1.2 {
		score -= 10
	}
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func studioRecommendedMode(st storage.Stream, snapshot playerTelemetrySnapshot) string {
	if snapshot.TotalStalls >= 5 || snapshot.WaitingSessions >= 2 {
		return "resilient"
	}
	if st.InputWidth >= 1920 && st.InputBitrate > 3500000 {
		return "balanced"
	}
	if st.InputWidth > 0 && st.InputWidth <= 960 {
		return "mobile"
	}
	return "balanced"
}

func studioPublicBaseURL(r *http.Request, cfg *config.Manager) string {
	if cfg == nil {
		return strings.TrimRight("http://localhost:8844", "/")
	}
	host := strings.TrimSpace(cfg.Get("embed_domain", ""))
	if host == "" || strings.EqualFold(host, "localhost") {
		host = studioRequestHostName(r.Host)
	}
	if host == "" {
		host = "localhost"
	}
	useHTTPS := cfg.GetBool("embed_use_https", false)
	scheme := "http"
	port := cfg.GetInt("embed_http_port", cfg.GetInt("http_port", 8844))
	if useHTTPS {
		scheme = "https"
		port = cfg.GetInt("embed_https_port", cfg.GetInt("https_port", 443))
	}
	defaultPort := 80
	if useHTTPS {
		defaultPort = 443
	}
	if port == defaultPort || port <= 0 {
		return fmt.Sprintf("%s://%s", scheme, host)
	}
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

func studioPreviewBaseURL(r *http.Request) string {
	if r == nil {
		return "http://localhost:8844"
	}
	scheme := "http"
	if xf := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); xf != "" {
		scheme = strings.ToLower(xf)
	} else if r.TLS != nil {
		scheme = "https"
	}
	host := studioRequestHostName(r.Host)
	if host == "" {
		host = "localhost"
	}
	port := 80
	if scheme == "https" {
		port = 443
	}
	if h, p, err := net.SplitHostPort(r.Host); err == nil {
		host = strings.Trim(h, "[]")
		if parsed, convErr := strconv.Atoi(p); convErr == nil && parsed > 0 {
			port = parsed
		}
	} else if strings.Count(r.Host, ":") == 1 {
		if idx := strings.LastIndex(r.Host, ":"); idx > 0 {
			if parsed, convErr := strconv.Atoi(r.Host[idx+1:]); convErr == nil && parsed > 0 {
				port = parsed
			}
		}
	}
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", scheme, host)
	}
	return fmt.Sprintf("%s://%s:%d", scheme, host, port)
}

func studioManifestURL(base, streamKey, format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "dash", "dash_audio":
		if format == "dash_audio" {
			return fmt.Sprintf("%s/audio/dash/%s", base, streamKey)
		}
		return fmt.Sprintf("%s/dash/%s/manifest.mpd", base, streamKey)
	case "hls", "hls_audio":
		if format == "hls_audio" {
			return fmt.Sprintf("%s/audio/hls/%s", base, streamKey)
		}
		return fmt.Sprintf("%s/hls/%s/master.m3u8", base, streamKey)
	case "mp4":
		return fmt.Sprintf("%s/mp4/%s/%s.mp4", base, streamKey, streamKey)
	case "audio", "aac", "mp3", "ogg", "wav", "flac":
		return fmt.Sprintf("%s/audio/%s/%s/%s.%s", base, format, streamKey, streamKey, strings.TrimPrefix(format, "audio/"))
	default:
		return fmt.Sprintf("%s/hls/%s/master.m3u8", base, streamKey)
	}
}

func studioRequestHostName(hostPort string) string {
	hostPort = strings.TrimSpace(hostPort)
	if hostPort == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(hostPort); err == nil {
		return strings.Trim(host, "[]")
	}
	if strings.HasPrefix(hostPort, "[") {
		if end := strings.Index(hostPort, "]"); end > 1 {
			return strings.Trim(hostPort[1:end], "[]")
		}
	}
	if strings.Count(hostPort, ":") == 1 {
		if idx := strings.LastIndex(hostPort, ":"); idx > 0 {
			return hostPort[:idx]
		}
	}
	return hostPort
}

func withQueryMap(raw string, params map[string]string) string {
	if strings.TrimSpace(raw) == "" {
		return raw
	}
	if len(params) == 0 {
		return raw
	}
	separator := "?"
	if strings.Contains(raw, "?") {
		separator = "&"
	}
	parts := make([]string, 0, len(params))
	for _, key := range []string{"format", "token", "viewer_id", "player_watermark", "autoplay", "muted"} {
		value := strings.TrimSpace(params[key])
		if value == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", key, urlQueryEscape(value)))
	}
	if len(parts) == 0 {
		return raw
	}
	return raw + separator + strings.Join(parts, "&")
}

func urlQueryEscape(value string) string {
	replacer := strings.NewReplacer(
		"%", "%25",
		" ", "%20",
		"+", "%2B",
		"&", "%26",
		"=", "%3D",
		"?", "%3F",
		"#", "%23",
		"/", "%2F",
	)
	return replacer.Replace(value)
}

func parseABRProfilesMap(raw string) map[string][]map[string]interface{} {
	out := map[string][]map[string]interface{}{}
	if strings.TrimSpace(raw) == "" {
		return out
	}
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}

func parseABRProfileLayers(raw string) []map[string]interface{} {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []map[string]interface{}{}
	}
	var layers []map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &layers); err == nil {
		return layers
	}
	wrapped := map[string][]map[string]interface{}{}
	if err := json.Unmarshal([]byte(raw), &wrapped); err == nil {
		for _, items := range wrapped {
			return items
		}
	}
	return []map[string]interface{}{}
}

func randomTokenFragment(size int) string {
	if size <= 0 {
		size = 8
	}
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(buf)[:size]
}

func studioSlug(raw, fallback string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		raw = fallback
	}
	var b strings.Builder
	lastDash := false
	for _, r := range raw {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	result := strings.Trim(b.String(), "-")
	if result == "" {
		return fallback
	}
	return result
}

func studioAssetCategory(raw string) string {
	value := studioSlug(raw, "branding")
	switch value {
	case "branding", "logos", "posters", "players", "embed":
		return value
	default:
		return "branding"
	}
}

func studioAssetsDir(dataDir, category string) string {
	return filepath.Join(dataDir, "assets", studioAssetCategory(category))
}

func studioAssetURL(category, name string) string {
	return "/media-assets/" + studioAssetCategory(category) + "/" + urlQueryEscape(name)
}

func listStudioAssets(dataDir, category string) ([]map[string]interface{}, error) {
	dir := studioAssetsDir(dataDir, category)
	_ = os.MkdirAll(dir, 0755)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	items := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		name := entry.Name()
		items = append(items, map[string]interface{}{
			"name":      name,
			"category":  studioAssetCategory(category),
			"path":      studioAssetCategory(category) + "/" + name,
			"url":       studioAssetURL(category, name),
			"size":      info.Size(),
			"mod_time":  info.ModTime().UTC().Format(time.RFC3339),
			"extension": strings.ToLower(filepath.Ext(name)),
		})
	}
	sort.Slice(items, func(i, j int) bool {
		left, _ := time.Parse(time.RFC3339, asString(items[i]["mod_time"]))
		right, _ := time.Parse(time.RFC3339, asString(items[j]["mod_time"]))
		return right.Before(left)
	})
	return items, nil
}

func saveStudioAsset(dataDir, category, originalName string, src io.Reader) (map[string]interface{}, error) {
	dir := studioAssetsDir(dataDir, category)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	base := strings.TrimSuffix(filepath.Base(strings.TrimSpace(originalName)), filepath.Ext(strings.TrimSpace(originalName)))
	base = studioSlug(base, "asset")
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(originalName)))
	if ext == "" {
		ext = ".bin"
	}
	filename := fmt.Sprintf("%s-%s%s", base, time.Now().UTC().Format("20060102-150405"), ext)
	target := filepath.Join(dir, filename)
	out, err := os.Create(target)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	size, err := io.Copy(out, src)
	if err != nil {
		return nil, err
	}
	info, err := out.Stat()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"name":      filename,
		"category":  studioAssetCategory(category),
		"path":      studioAssetCategory(category) + "/" + filename,
		"url":       studioAssetURL(category, filename),
		"size":      size,
		"mod_time":  info.ModTime().UTC().Format(time.RFC3339),
		"extension": ext,
	}, nil
}

func deleteStudioAsset(dataDir, relPath string) error {
	relPath = strings.Trim(strings.TrimSpace(relPath), "/\\")
	if relPath == "" {
		return fmt.Errorf("Silinecek dosya yolu gerekli")
	}
	parts := strings.FieldsFunc(relPath, func(r rune) bool { return r == '/' || r == '\\' })
	if len(parts) != 2 {
		return fmt.Errorf("Gecersiz asset yolu")
	}
	category := studioAssetCategory(parts[0])
	name := filepath.Base(parts[1])
	if name == "." || name == "" || strings.Contains(name, "..") {
		return fmt.Errorf("Gecersiz asset adi")
	}
	target := filepath.Join(studioAssetsDir(dataDir, category), name)
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func asString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		return ""
	}
}

func asInt(value interface{}, fallback int) int {
	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case json.Number:
		if parsed, err := v.Int64(); err == nil {
			return int(parsed)
		}
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return parsed
		}
	}
	return fallback
}

func asBool(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case int:
		return v != 0
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "1", "true", "yes", "on":
			return true
		}
	}
	return false
}
