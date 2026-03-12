package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/security"
	"github.com/fluxstream/fluxstream/internal/storage"
	streampolicy "github.com/fluxstream/fluxstream/internal/stream"
)

type playbackAuthorizer func(r *http.Request, streamKey, format string) (bool, int, string)

func makePlaybackAuthorizer(cfg *config.Manager, db *storage.SQLiteDB, tokenMgr *security.TokenManager) playbackAuthorizer {
	return func(r *http.Request, streamKey, format string) (bool, int, string) {
		if streamKey == "" || isInternalPlaybackRequest(r) {
			return true, 0, ""
		}
		st, err := db.GetStreamByKey(streamKey)
		if err != nil || st == nil {
			return true, 0, ""
		}

		policy := streampolicy.ParsePolicyJSON(st.PolicyJSON)
		if !streamAllowsFormat(st.OutputFormats, format) {
			return false, http.StatusForbidden, "Bu yayin icin istenen cikis formati kapali."
		}
		if len(policy.AllowedOutputs) > 0 && !policy.AllowsOutput(format) {
			return false, http.StatusForbidden, "Bu yayin icin istenen cikis formati kapali."
		}
		if !ipAllowed(st.IPWhitelist, requestClientIP(r)) {
			return false, http.StatusForbidden, "IP bu yayin icin yetkili degil."
		}
		if !domainAllowed(st.DomainLock, r) {
			return false, http.StatusForbidden, "Bu yayin yalnizca izinli domainlerden acilabilir."
		}
		if st.Password != "" {
			pw := strings.TrimSpace(r.URL.Query().Get("password"))
			if pw == "" {
				pw = strings.TrimSpace(r.Header.Get("X-Stream-Password"))
			}
			if pw != st.Password {
				return false, http.StatusUnauthorized, "Yayin sifresi gerekli."
			}
		}
		needsToken := cfg.GetBool("token_enabled", false) || policy.RequirePlaybackToken || policy.RequireSignedURL
		if needsToken {
			token := strings.TrimSpace(r.URL.Query().Get("token"))
			if token == "" {
				token = strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
			}
			if token == "" || !tokenMgr.ValidateToken(token, streamKey) {
				return false, http.StatusUnauthorized, "Gecerli playback token gerekli."
			}
		}
		return true, 0, ""
	}
}

func streamAllowsFormat(raw, format string) bool {
	format = normalizeOutputFormat(format)
	if strings.TrimSpace(raw) == "" {
		return true
	}
	var items []string
	if err := json.Unmarshal([]byte(raw), &items); err != nil || len(items) == 0 {
		return true
	}
	if len(items) == 1 && normalizeOutputFormat(items[0]) == "hls" {
		// Older installs stored HLS as a placeholder default rather than an explicit allow-list.
		return true
	}
	for _, item := range items {
		if normalizeOutputFormat(item) == format {
			return true
		}
	}
	return false
}

func normalizeOutputFormat(format string) string {
	format = strings.TrimSpace(strings.ToLower(format))
	switch format {
	case "", "player", "embed", "iframe", "jsapi", "hls_master", "ll_hls", "hls_audio":
		return "hls"
	case "dash_audio":
		return "dash"
	case "http_flv":
		return "flv"
	case "fmp4":
		return "mp4"
	default:
		return format
	}
}

func ipAllowed(raw, ip string) bool {
	raw = strings.TrimSpace(raw)
	ip = strings.TrimSpace(ip)
	if raw == "" || ip == "" {
		return true
	}
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, cidr, err := net.ParseCIDR(item); err == nil {
			if parsed := net.ParseIP(ip); parsed != nil && cidr.Contains(parsed) {
				return true
			}
			continue
		}
		if item == ip {
			return true
		}
	}
	return false
}

func domainAllowed(raw string, r *http.Request) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return true
	}
	originHost := strings.TrimSpace(r.Header.Get("Origin"))
	refererHost := strings.TrimSpace(r.Referer())
	targets := strings.Split(raw, ",")
	for _, target := range targets {
		target = strings.ToLower(strings.TrimSpace(target))
		if target == "" {
			continue
		}
		if strings.Contains(strings.ToLower(originHost), target) || strings.Contains(strings.ToLower(refererHost), target) {
			return true
		}
	}
	return false
}
