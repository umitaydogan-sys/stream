package tlsutil

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxstream/fluxstream/internal/config"
	"golang.org/x/crypto/acme/autocert"
)

type Profile string

const (
	ProfileWeb    Profile = "web"
	ProfileStream Profile = "stream"
)

type Source struct {
	Profile  Profile
	Enabled  bool
	Mode     string
	Domain   string
	Email    string
	CertPath string
	KeyPath  string
	Ready    bool

	manager    *autocert.Manager
	manualCert *tls.Certificate
}

func NewSource(cfg *config.Manager, profile Profile, dataDir string) (*Source, error) {
	source := &Source{
		Profile:  profile,
		Enabled:  enabled(cfg, profile),
		Mode:     mode(cfg, profile),
		Domain:   strings.TrimSpace(domain(cfg, profile)),
		Email:    strings.TrimSpace(email(cfg, profile)),
		CertPath: manualCertPath(cfg, profile, dataDir),
		KeyPath:  manualKeyPath(cfg, profile, dataDir),
	}
	if !source.Enabled {
		return source, nil
	}

	if source.Mode == "letsencrypt" {
		if source.Domain == "" {
			return source, nil
		}
		cacheDir := filepath.Join(dataDir, "certs", "acme", string(profile))
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			return nil, fmt.Errorf("acme cache directory: %w", err)
		}
		source.manager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(cacheDir),
			HostPolicy: autocert.HostWhitelist(source.Domain),
			Email:      source.Email,
		}
		source.Ready = true
		return source, nil
	}

	if source.CertPath == "" || source.KeyPath == "" {
		return source, nil
	}
	cert, err := tls.LoadX509KeyPair(source.CertPath, source.KeyPath)
	if err != nil {
		return source, nil
	}
	source.manualCert = &cert
	source.Ready = true
	return source, nil
}

func (s *Source) UsesLetsEncrypt() bool {
	return s != nil && s.manager != nil
}

func (s *Source) Manager() *autocert.Manager {
	if s == nil {
		return nil
	}
	return s.manager
}

func (s *Source) TLSConfig() *tls.Config {
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if s == nil || !s.Ready {
		return cfg
	}
	if s.manager != nil {
		tlsCfg := s.manager.TLSConfig()
		tlsCfg.MinVersion = tls.VersionTLS12
		return tlsCfg
	}
	if s.manualCert != nil {
		cfg.Certificates = []tls.Certificate{*s.manualCert}
	}
	return cfg
}

func ChallengeHandler(fallback http.Handler, sources ...*Source) http.Handler {
	if fallback == nil {
		fallback = http.NotFoundHandler()
	}
	handlers := map[string]http.Handler{}
	for _, source := range sources {
		if source == nil || source.manager == nil || source.Domain == "" {
			continue
		}
		handlers[strings.ToLower(source.Domain)] = source.manager.HTTPHandler(fallback)
	}
	if len(handlers) == 0 {
		return fallback
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := strings.ToLower(requestHostName(r.Host))
		if handler, ok := handlers[host]; ok {
			handler.ServeHTTP(w, r)
			return
		}
		fallback.ServeHTTP(w, r)
	})
}

func enabled(cfg *config.Manager, profile Profile) bool {
	switch profile {
	case ProfileStream:
		return cfg.GetBool("rtmps_enabled", false)
	default:
		return cfg.GetBool("ssl_enabled", false)
	}
}

func mode(cfg *config.Manager, profile Profile) string {
	switch profile {
	case ProfileStream:
		return strings.ToLower(strings.TrimSpace(cfg.Get("stream_ssl_mode", "file")))
	default:
		return strings.ToLower(strings.TrimSpace(cfg.Get("ssl_mode", "file")))
	}
}

func domain(cfg *config.Manager, profile Profile) string {
	switch profile {
	case ProfileStream:
		return cfg.Get("stream_ssl_le_domain", "")
	default:
		return cfg.Get("ssl_le_domain", "")
	}
}

func email(cfg *config.Manager, profile Profile) string {
	switch profile {
	case ProfileStream:
		return cfg.Get("stream_ssl_le_email", "")
	default:
		return cfg.Get("ssl_le_email", "")
	}
}

func manualCertPath(cfg *config.Manager, profile Profile, dataDir string) string {
	switch profile {
	case ProfileStream:
		path := strings.TrimSpace(cfg.Get("stream_ssl_cert_path", ""))
		if path != "" {
			return path
		}
		fallbacks := []string{
			filepath.Join(dataDir, "certs", "stream", "server.crt"),
			filepath.Join(dataDir, "certs", "server.crt"),
		}
		for _, candidate := range fallbacks {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
		return fallbacks[0]
	default:
		path := strings.TrimSpace(cfg.Get("ssl_cert_path", ""))
		if path != "" {
			return path
		}
		fallbacks := []string{
			filepath.Join(dataDir, "certs", "web", "server.crt"),
			filepath.Join(dataDir, "certs", "server.crt"),
		}
		for _, candidate := range fallbacks {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
		return fallbacks[0]
	}
}

func manualKeyPath(cfg *config.Manager, profile Profile, dataDir string) string {
	switch profile {
	case ProfileStream:
		path := strings.TrimSpace(cfg.Get("stream_ssl_key_path", ""))
		if path != "" {
			return path
		}
		fallbacks := []string{
			filepath.Join(dataDir, "certs", "stream", "server.key"),
			filepath.Join(dataDir, "certs", "server.key"),
		}
		for _, candidate := range fallbacks {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
		return fallbacks[0]
	default:
		path := strings.TrimSpace(cfg.Get("ssl_key_path", ""))
		if path != "" {
			return path
		}
		fallbacks := []string{
			filepath.Join(dataDir, "certs", "web", "server.key"),
			filepath.Join(dataDir, "certs", "server.key"),
		}
		for _, candidate := range fallbacks {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
		return fallbacks[0]
	}
}

func requestHostName(hostPort string) string {
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
	return strings.Trim(hostPort, "[]")
}
