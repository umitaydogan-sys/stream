package security

import (
	"net/http/httptest"
	"testing"
)

func TestRequestMatchesDomainUsesHostBoundaries(t *testing.T) {
	req := httptest.NewRequest("GET", "https://player.example.com/embed/live_demo", nil)
	req.Host = "player.example.com"
	req.Header.Set("Origin", "https://player.example.com")
	req.Header.Set("Referer", "https://portal.example.com/watch/live_demo")

	if !requestMatchesDomain(req, "example.com") {
		t.Fatalf("expected example.com to match subdomains")
	}
	if !requestMatchesDomain(req, "*.example.com") {
		t.Fatalf("expected wildcard rule to match subdomains")
	}
	if requestMatchesDomain(req, "evil-example.com") {
		t.Fatalf("did not expect substring-like domain to match")
	}
}
