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
	checks := []map[string]interface{}{
		diagnosticCheck("hls", fileExists(liveManifest), "Canli HLS playlist"),
		diagnosticCheck("hls_master", fileExists(hlsMasterPath), "Adaptif HLS master playlist"),
		diagnosticCheck("ll_hls", fileExists(llHLSPath), "Low latency playlist"),
		diagnosticCheck("dash", fileExists(dashManifestPath), "DASH manifest"),
		diagnosticCheck("recordings", folderHasFiles(filepath.Join(dataDir, "recordings", st.StreamKey)), "Kayit dosyalari"),
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
		"abr_profile_set":           cfg.Get("abr_profile_set", "balanced"),
		"hls_variant_count":         hlsVariantCount,
		"dash_representation_count": dashRepresentationCount,
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
