package main

import (
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/security"
	"github.com/fluxstream/fluxstream/internal/storage"
	streampolicy "github.com/fluxstream/fluxstream/internal/stream"
)

func TestDomainAllowedMatchesHostSafely(t *testing.T) {
	req := httptest.NewRequest("GET", "https://player.example.com/play/live_demo", nil)
	req.Host = "player.example.com"

	if !domainAllowed("example.com", req) {
		t.Fatalf("expected parent domain to match subdomain host")
	}
	if domainAllowed("evil-example.com", req) {
		t.Fatalf("substring domain should not match")
	}
}

func TestPlaybackAuthorizerRequiresSignedURLQueryToken(t *testing.T) {
	tempDir := t.TempDir()
	db, err := storage.NewSQLiteDB(filepath.Join(tempDir, "test.db"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()
	cfg := config.NewManager(db)
	if err := cfg.LoadDefaults(); err != nil {
		t.Fatalf("load defaults: %v", err)
	}

	policy := streampolicy.DefaultPolicy()
	policy.RequireSignedURL = true
	streamID, err := db.CreateStream(&storage.Stream{
		Name:       "Signed Stream",
		StreamKey:  "live_signed_demo",
		Status:     "live",
		PolicyJSON: streampolicy.EncodePolicyJSON(policy),
	})
	if err != nil {
		t.Fatalf("create stream: %v", err)
	}
	if streamID <= 0 {
		t.Fatalf("expected stream to be created")
	}

	tokenMgr := security.NewTokenManager("test-secret", 60)
	auth := makePlaybackAuthorizer(cfg, db, tokenMgr)

	queryToken, _, err := tokenMgr.GeneratePlaybackToken("live_signed_demo", security.PlaybackTokenClaims{
		AllowedFormat: "hls",
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	noTokenReq := httptest.NewRequest("GET", "https://example.com/hls/live_signed_demo/master.m3u8", nil)
	if ok, status, _ := auth(noTokenReq, "live_signed_demo", "hls"); ok || status == 0 {
		t.Fatalf("expected missing query token to be rejected")
	}

	headerTokenReq := httptest.NewRequest("GET", "https://example.com/hls/live_signed_demo/master.m3u8", nil)
	headerTokenReq.Header.Set("Authorization", "Bearer "+queryToken)
	if ok, status, _ := auth(headerTokenReq, "live_signed_demo", "hls"); ok || status == 0 {
		t.Fatalf("expected header-only token to be rejected for signed url policy")
	}

	queryTokenReq := httptest.NewRequest("GET", "https://example.com/hls/live_signed_demo/master.m3u8?token="+queryToken, nil)
	if ok, status, msg := auth(queryTokenReq, "live_signed_demo", "hls"); !ok || status != 0 || msg != "" {
		t.Fatalf("expected query token to pass, got ok=%v status=%d msg=%q", ok, status, msg)
	}
}
