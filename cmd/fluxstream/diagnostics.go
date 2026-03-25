package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/storage"
	"github.com/fluxstream/fluxstream/internal/transcode"
)

func buildStreamDiagnostics(st *storage.Stream, cfg *config.Manager, dataDir string, tcManager *transcode.Manager) map[string]interface{} {
	liveManifest := ""
	liveRoot := ""
	dashRoot := ""
	if tcManager != nil {
		liveManifest = tcManager.GetLiveManifestPath(st.StreamKey)
		liveRoot = filepath.Join(tcManager.GetLiveOutputDir(), st.StreamKey)
		dashRoot = filepath.Join(tcManager.GetLiveDashOutputDir(), st.StreamKey)
	}
	hlsMasterPath := ""
	if liveRoot != "" {
		hlsMasterPath = filepath.Join(liveRoot, "master.m3u8")
	}
	dashManifestPath := ""
	if dashRoot != "" {
		dashManifestPath = filepath.Join(dashRoot, "manifest.mpd")
	}
	hlsVariantCount := countInTextFile(hlsMasterPath, "#EXT-X-STREAM-INF")
	dashRepresentationCount := countInTextFile(dashManifestPath, "<Representation")
	llHLSPath := ""
	if liveRoot != "" {
		llHLSPath = filepath.Join(liveRoot, "ll.m3u8")
	}

	streamLive := strings.EqualFold(strings.TrimSpace(st.Status), "live")
	hlsEnabled := cfg.GetBool("hls_enabled", true)
	abrEnabled := cfg.GetBool("abr_enabled", false)
	abrMasterEnabled := cfg.GetBool("abr_master_enabled", true)
	llHLSEnabled := cfg.GetBool("hls_ll_enabled", false)
	dashEnabled := cfg.GetBool("dash_enabled", false) && cfg.GetBool("transcode_live_dash_enabled", true)
	recordingSystemEnabled := cfg.GetBool("recording_enabled", true)
	recordingExpected := recordingSystemEnabled && st.RecordEnabled

	checks := []map[string]interface{}{
		diagnosticCheck("hls", diagnosticStatusPrimaryHLS(hlsEnabled, streamLive, fileExists(liveManifest)), "Canli HLS playlist", ""),
		diagnosticCheck("hls_master", diagnosticStatusABRMaster(hlsEnabled, abrEnabled, abrMasterEnabled, streamLive, fileExists(hlsMasterPath), hlsVariantCount), "Adaptif HLS master playlist", ""),
		diagnosticCheck("ll_hls", diagnosticStatusOptionalOutput(llHLSEnabled, streamLive, fileExists(llHLSPath)), "Low latency playlist", "LL-HLS kapaliysa bu alan sorun sayilmaz."),
		diagnosticCheck("dash", diagnosticStatusDash(dashEnabled, streamLive, fileExists(dashManifestPath), dashRepresentationCount), "DASH manifest", "DASH kapaliysa veya repack yeni basliyorsa sari/mavi gorunebilir."),
		diagnosticCheck("recordings", diagnosticStatusRecording(recordingExpected, streamLive, folderHasFiles(filepath.Join(dataDir, "recordings", st.StreamKey))), "Kayit dosyalari", "Kayit pasifse ya da yayin yeni basladiysa hemen dosya gormeyebilirsin."),
	}

	outputFormats := strings.TrimSpace(st.OutputFormats)
	return map[string]interface{}{
		"stream_id":                 st.ID,
		"stream_key":                st.StreamKey,
		"stream_name":               st.Name,
		"status":                    st.Status,
		"output_formats":            outputFormats,
		"policy_json":               st.PolicyJSON,
		"checks":                    checks,
		"abr_enabled":               cfg.GetBool("abr_enabled", false),
		"abr_master_enabled":        abrMasterEnabled,
		"abr_profile_set":           cfg.Get("abr_profile_set", "balanced"),
		"hls_enabled":               hlsEnabled,
		"dash_enabled":              dashEnabled,
		"hls_ll_enabled":            llHLSEnabled,
		"recording_system_enabled":  recordingSystemEnabled,
		"recording_stream_enabled":  st.RecordEnabled,
		"delivery_summary":          diagnosticDeliverySummary(abrEnabled, hlsVariantCount, dashEnabled, dashRepresentationCount),
		"hls_variant_count":         hlsVariantCount,
		"dash_representation_count": dashRepresentationCount,
	}
}

func diagnosticCheck(code, status, description, detail string) map[string]interface{} {
	ok := status == "ready"
	label, tone := diagnosticBadge(status)
	return map[string]interface{}{
		"code":        code,
		"ok":          ok,
		"status":      status,
		"label":       label,
		"tone":        tone,
		"description": description,
		"detail":      detail,
	}
}

func diagnosticBadge(status string) (string, string) {
	switch strings.TrimSpace(strings.ToLower(status)) {
	case "ready":
		return "Hazir", "green"
	case "waiting":
		return "Bekliyor", "yellow"
	case "disabled":
		return "Kapali", "blue"
	case "optional":
		return "Opsiyonel", "blue"
	default:
		return "Sorunlu", "red"
	}
}

func diagnosticStatusPrimaryHLS(enabled, live, exists bool) string {
	if !enabled {
		return "disabled"
	}
	if exists {
		return "ready"
	}
	if live {
		return "waiting"
	}
	return "optional"
}

func diagnosticStatusABRMaster(hlsEnabled, abrEnabled, abrMasterEnabled, live, exists bool, variantCount int) string {
	if !hlsEnabled || !abrEnabled || !abrMasterEnabled {
		return "disabled"
	}
	if exists && variantCount > 1 {
		return "ready"
	}
	if exists {
		return "waiting"
	}
	if live {
		return "waiting"
	}
	return "optional"
}

func diagnosticStatusOptionalOutput(enabled, live, exists bool) string {
	if !enabled {
		return "disabled"
	}
	if exists {
		return "ready"
	}
	if live {
		return "waiting"
	}
	return "optional"
}

func diagnosticStatusDash(enabled, live, exists bool, representations int) string {
	if !enabled {
		return "disabled"
	}
	if exists && representations > 0 {
		return "ready"
	}
	if live {
		return "waiting"
	}
	return "optional"
}

func diagnosticStatusRecording(expected, live, exists bool) string {
	if !expected {
		return "disabled"
	}
	if exists {
		return "ready"
	}
	if live {
		return "waiting"
	}
	return "optional"
}

func diagnosticDeliverySummary(abrEnabled bool, hlsVariants int, dashEnabled bool, dashRepresentations int) map[string]interface{} {
	label := "Tek kalite"
	tone := "blue"
	description := "ABR kapali; yayin tek kalite teslim ediliyor."
	if abrEnabled {
		label = "ABR bekliyor"
		tone = "yellow"
		description = "ABR acik fakat coklu katman henuz tam gorunmuyor."
		if hlsVariants > 1 && (!dashEnabled || dashRepresentations > 0) {
			label = "ABR hazir"
			tone = "green"
			if dashEnabled {
				description = "HLS ve DASH tarafinda adaptif katmanlar hazir gorunuyor."
			} else {
				description = "HLS tarafinda adaptif katmanlar hazir; DASH kapali veya opsiyonel."
			}
		} else if hlsVariants > 1 {
			label = "HLS ABR hazir"
			tone = "green"
			description = "HLS tarafinda adaptif katmanlar hazir; DASH tarafi bekliyor veya kapali."
		}
	}
	return map[string]interface{}{
		"label":       label,
		"tone":        tone,
		"description": description,
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func folderHasFiles(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			return true
		}
	}
	return false
}

func countInTextFile(path, needle string) int {
	if path == "" || needle == "" {
		return 0
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return strings.Count(string(data), needle)
}
