package main

import (
	"sort"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/license"
	"github.com/fluxstream/fluxstream/internal/storage"
	streamcfg "github.com/fluxstream/fluxstream/internal/stream"
)

const (
	licenseFeatureABR       = "abr"
	licenseFeatureRTMPS     = "rtmps"
	licenseFeatureRecording = "recording"
	licenseFeatureBranding  = "branding"
)

var developmentFeatureSet = []string{
	licenseFeatureABR,
	licenseFeatureRTMPS,
	licenseFeatureRecording,
	licenseFeatureBranding,
}

type runtimeLicense struct {
	Status          license.Status  `json:"status"`
	Mode            string          `json:"mode"`
	Development     bool            `json:"development"`
	Enforced        bool            `json:"enforced"`
	EnabledFeatures []string        `json:"enabled_features"`
	FeatureMap      map[string]bool `json:"feature_map"`
	Warnings        []string        `json:"warnings,omitempty"`
}

func resolveRuntimeLicense(dataDir string) *runtimeLicense {
	status := license.NewManager(dataDir).Status(time.Now().UTC())
	rt := &runtimeLicense{
		Status:     status,
		Mode:       status.Mode,
		Enforced:   true,
		FeatureMap: map[string]bool{},
	}
	for _, feature := range developmentFeatureSet {
		rt.FeatureMap[feature] = false
	}

	switch {
	case status.Valid:
		rt.Mode = "licensed"
		rt.enableFeatures(status.Features)
	case status.UsingEmbeddedKey && status.Mode == "unlicensed":
		rt.Mode = "development"
		rt.Development = true
		rt.Enforced = false
		rt.enableFeatures(developmentFeatureSet)
		rt.Warnings = append(rt.Warnings, "Embedded development key aktif; production icin imzali lisans yukleyin.")
	default:
		rt.Mode = status.Mode
		if strings.TrimSpace(status.Message) != "" {
			rt.Warnings = append(rt.Warnings, status.Message)
		}
	}

	rt.EnabledFeatures = rt.enabledFeatureList()
	return rt
}

func (r *runtimeLicense) enableFeatures(features []string) {
	for _, feature := range features {
		key := strings.ToLower(strings.TrimSpace(feature))
		if key == "" {
			continue
		}
		r.FeatureMap[key] = true
	}
}

func (r *runtimeLicense) enabledFeatureList() []string {
	list := make([]string, 0, len(r.FeatureMap))
	for feature, enabled := range r.FeatureMap {
		if enabled {
			list = append(list, feature)
		}
	}
	sort.Strings(list)
	return list
}

func (r *runtimeLicense) allows(feature string) bool {
	if r == nil {
		return false
	}
	return r.FeatureMap[strings.ToLower(strings.TrimSpace(feature))]
}

func (r *runtimeLicense) normalizeSettings(section string, updates map[string]string) {
	if r == nil || updates == nil {
		return
	}
	switch strings.TrimSpace(strings.ToLower(section)) {
	case "outputs":
		if !r.allows(licenseFeatureABR) {
			updates["abr_enabled"] = "false"
			updates["abr_master_enabled"] = "false"
		}
	case "protocols":
		if !r.allows(licenseFeatureRTMPS) {
			updates["rtmps_enabled"] = "false"
		}
	case "recording":
		if !r.allows(licenseFeatureRecording) {
			updates["recording_enabled"] = "false"
		}
	}
}

func (r *runtimeLicense) normalizeStream(st *storage.Stream) {
	if r == nil || st == nil {
		return
	}
	if !r.allows(licenseFeatureRecording) {
		st.RecordEnabled = false
	}
	policy := streamcfg.ParsePolicyJSON(st.PolicyJSON)
	if !r.allows(licenseFeatureABR) && policy.EnableABR {
		policy.EnableABR = false
	}
	st.PolicyJSON = streamcfg.EncodePolicyJSON(policy)
}

func (r *runtimeLicense) normalizePlayerTemplate(pt *storage.PlayerTemplate) {
	if r == nil || pt == nil || r.allows(licenseFeatureBranding) {
		return
	}
	pt.LogoURL = ""
	pt.LogoPosition = "top-right"
	pt.LogoOpacity = 1
	pt.WatermarkText = ""
	pt.CustomCSS = ""
}
