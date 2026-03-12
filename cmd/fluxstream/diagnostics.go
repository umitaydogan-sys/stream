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
	if tcManager != nil {
		liveManifest = tcManager.GetLiveManifestPath(st.StreamKey)
	}
	checks := []map[string]interface{}{
		diagnosticCheck("hls", fileExists(liveManifest), "Canli HLS playlist"),
		diagnosticCheck("hls_master", fileExists(filepath.Join(tcManager.GetLiveOutputDir(), st.StreamKey, "master.m3u8")), "Adaptif HLS master playlist"),
		diagnosticCheck("ll_hls", fileExists(filepath.Join(tcManager.GetLiveOutputDir(), st.StreamKey, "ll.m3u8")), "Low latency playlist"),
		diagnosticCheck("dash", fileExists(filepath.Join(tcManager.GetLiveDashOutputDir(), st.StreamKey, "manifest.mpd")), "DASH manifest"),
		diagnosticCheck("recordings", folderHasFiles(filepath.Join(dataDir, "recordings", st.StreamKey)), "Kayit dosyalari"),
	}

	outputFormats := strings.TrimSpace(st.OutputFormats)
	return map[string]interface{}{
		"stream_id":       st.ID,
		"stream_key":      st.StreamKey,
		"stream_name":     st.Name,
		"status":          st.Status,
		"output_formats":  outputFormats,
		"policy_json":     st.PolicyJSON,
		"checks":          checks,
		"abr_enabled":     cfg.GetBool("abr_enabled", false),
		"abr_profile_set": cfg.Get("abr_profile_set", "balanced"),
	}
}

func diagnosticCheck(code string, ok bool, description string) map[string]interface{} {
	status := "ready"
	if !ok {
		status = "missing"
	}
	return map[string]interface{}{
		"code":        code,
		"ok":          ok,
		"status":      status,
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
