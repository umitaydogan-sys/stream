package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TokenManager manages viewer authentication tokens
type TokenManager struct {
	secret   []byte
	duration time.Duration
}

// PlaybackTokenClaims stores optional restrictions for signed playback links.
type PlaybackTokenClaims struct {
	StreamKey     string `json:"stream_key"`
	ExpiresAtUnix int64  `json:"exp"`
	AllowedIP     string `json:"allowed_ip,omitempty"`
	AllowedDomain string `json:"allowed_domain,omitempty"`
	ViewerID      string `json:"viewer_id,omitempty"`
	Watermark     string `json:"watermark,omitempty"`
	AllowedFormat string `json:"allowed_format,omitempty"`
}

// NewTokenManager creates a new token manager
func NewTokenManager(secret string, durationMinutes int) *TokenManager {
	key := []byte(secret)
	if len(key) == 0 {
		key = make([]byte, 32)
		rand.Read(key)
	}
	return &TokenManager{
		secret:   key,
		duration: time.Duration(durationMinutes) * time.Minute,
	}
}

// GenerateToken creates a signed token for a stream key
func (tm *TokenManager) GenerateToken(streamKey string) (string, time.Time) {
	expiry := time.Now().Add(tm.duration)
	payload := fmt.Sprintf("%s:%d", streamKey, expiry.Unix())
	mac := hmac.New(sha256.New, tm.secret)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	token := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", payload, sig)))
	return token, expiry
}

// ValidateToken verifies a token
func (tm *TokenManager) ValidateToken(token, streamKey string) bool {
	if strings.HasPrefix(strings.TrimSpace(token), "v2.") {
		return tm.ValidatePlaybackToken(token, streamKey, "", nil)
	}
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	parts := strings.SplitN(string(decoded), ":", 3)
	if len(parts) != 3 {
		return false
	}

	tokenStreamKey := parts[0]
	expiryStr := parts[1]
	providedSig := parts[2]

	if tokenStreamKey != streamKey {
		return false
	}

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		return false
	}

	if time.Now().Unix() > expiry {
		return false
	}

	// Verify signature
	payload := fmt.Sprintf("%s:%s", tokenStreamKey, expiryStr)
	mac := hmac.New(sha256.New, tm.secret)
	mac.Write([]byte(payload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(providedSig), []byte(expectedSig))
}

// GeneratePlaybackToken creates a scoped signed playback token.
func (tm *TokenManager) GeneratePlaybackToken(streamKey string, claims PlaybackTokenClaims) (string, time.Time, error) {
	streamKey = strings.TrimSpace(streamKey)
	if streamKey == "" {
		return "", time.Time{}, fmt.Errorf("stream key gerekli")
	}
	expiry := time.Now().Add(tm.duration)
	if claims.ExpiresAtUnix > 0 {
		expiry = time.Unix(claims.ExpiresAtUnix, 0)
	}
	claims.StreamKey = streamKey
	claims.ExpiresAtUnix = expiry.Unix()
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	payloadEnc := base64.RawURLEncoding.EncodeToString(payload)
	sig := tm.sign(payloadEnc)
	return "v2." + payloadEnc + "." + sig, expiry, nil
}

// ValidatePlaybackToken verifies both legacy and scoped playback tokens.
func (tm *TokenManager) ValidatePlaybackToken(token, streamKey, format string, r *http.Request) bool {
	token = strings.TrimSpace(token)
	streamKey = strings.TrimSpace(streamKey)
	if token == "" || streamKey == "" {
		return false
	}
	if !strings.HasPrefix(token, "v2.") {
		return tm.ValidateToken(token, streamKey)
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	payloadEnc := parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(tm.sign(payloadEnc))) {
		return false
	}
	raw, err := base64.RawURLEncoding.DecodeString(payloadEnc)
	if err != nil {
		return false
	}
	var claims PlaybackTokenClaims
	if err := json.Unmarshal(raw, &claims); err != nil {
		return false
	}
	if strings.TrimSpace(claims.StreamKey) != streamKey {
		return false
	}
	if claims.ExpiresAtUnix > 0 && time.Now().Unix() > claims.ExpiresAtUnix {
		return false
	}
	if claims.AllowedIP != "" && r != nil {
		if extractIP(r) != strings.TrimSpace(claims.AllowedIP) {
			return false
		}
	}
	if claims.AllowedDomain != "" && r != nil {
		if !requestMatchesDomain(r, claims.AllowedDomain) {
			return false
		}
	}
	if claims.ViewerID != "" && r != nil {
		if strings.TrimSpace(r.URL.Query().Get("viewer_id")) != strings.TrimSpace(claims.ViewerID) {
			return false
		}
	}
	if claims.AllowedFormat != "" && format != "" {
		if normalizePlaybackFormat(claims.AllowedFormat) != normalizePlaybackFormat(format) {
			return false
		}
	}
	return true
}

func (tm *TokenManager) sign(payloadEnc string) string {
	mac := hmac.New(sha256.New, tm.secret)
	mac.Write([]byte(payloadEnc))
	return hex.EncodeToString(mac.Sum(nil))
}

func requestMatchesDomain(r *http.Request, domain string) bool {
	domain = strings.TrimSpace(domain)
	if domain == "" || r == nil {
		return true
	}
	candidates := []string{
		strings.TrimSpace(r.Header.Get("Origin")),
		strings.TrimSpace(r.Referer()),
		strings.TrimSpace(r.Host),
	}
	for _, rule := range strings.Split(domain, ",") {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}
		for _, candidate := range candidates {
			if requestHostMatchesRule(candidate, rule) {
				return true
			}
		}
	}
	return false
}

func requestHostMatchesRule(candidate, rule string) bool {
	candidateHost := normalizeRequestHost(candidate)
	ruleHost := normalizeRequestHost(rule)
	if candidateHost == "" || ruleHost == "" {
		return false
	}
	if strings.HasPrefix(ruleHost, "*.") {
		ruleHost = strings.TrimPrefix(ruleHost, "*.")
	}
	return candidateHost == ruleHost || strings.HasSuffix(candidateHost, "."+ruleHost)
}

func normalizeRequestHost(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") && strings.Contains(raw, "/") {
		raw = strings.Split(raw, "/")[0]
	}
	if strings.Contains(raw, "://") {
		if parsed, err := url.Parse(raw); err == nil {
			raw = parsed.Host
		}
	}
	if strings.Contains(raw, "/") {
		raw = strings.Split(raw, "/")[0]
	}
	if host, _, err := net.SplitHostPort(raw); err == nil {
		raw = host
	}
	return strings.Trim(raw, "[]")
}

func normalizePlaybackFormat(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "", "auto", "player", "embed", "iframe", "jsapi", "hls_master", "ll_hls", "hls_audio":
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

// RateLimiter implements IP-based rate limiting
type RateLimiter struct {
	requests    map[string]*ipState
	maxRequests int
	window      time.Duration
	mu          sync.RWMutex
}

type ipState struct {
	count   int
	resetAt time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:    make(map[string]*ipState),
		maxRequests: maxRequests,
		window:      window,
	}
	go rl.cleanup()
	return rl
}

// Allow checks if a request from IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	state, exists := rl.requests[ip]
	if !exists || now.After(state.resetAt) {
		rl.requests[ip] = &ipState{count: 1, resetAt: now.Add(rl.window)}
		return true
	}

	state.count++
	return state.count <= rl.maxRequests
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, state := range rl.requests {
			if now.After(state.resetAt) {
				delete(rl.requests, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Middleware returns an HTTP middleware for rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		if !rl.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// IPBanList manages banned IPs
type IPBanList struct {
	banned map[string]banEntry
	mu     sync.RWMutex
}

type banEntry struct {
	IP        string
	Reason    string
	BannedAt  time.Time
	ExpiresAt time.Time // zero = permanent
}

// NewIPBanList creates a new IP ban list
func NewIPBanList() *IPBanList {
	return &IPBanList{
		banned: make(map[string]banEntry),
	}
}

// Ban adds an IP to the ban list
func (bl *IPBanList) Ban(ip, reason string, duration time.Duration) {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	entry := banEntry{
		IP:       ip,
		Reason:   reason,
		BannedAt: time.Now(),
	}
	if duration > 0 {
		entry.ExpiresAt = time.Now().Add(duration)
	}
	bl.banned[ip] = entry
}

// Unban removes an IP from the ban list
func (bl *IPBanList) Unban(ip string) {
	bl.mu.Lock()
	defer bl.mu.Unlock()
	delete(bl.banned, ip)
}

// IsBanned checks if an IP is banned
func (bl *IPBanList) IsBanned(ip string) (bool, string) {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	entry, exists := bl.banned[ip]
	if !exists {
		return false, ""
	}

	// Check expiration
	if !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
		return false, ""
	}

	return true, entry.Reason
}

// GetBanned returns all banned IPs
func (bl *IPBanList) GetBanned() []banEntry {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	var entries []banEntry
	for _, e := range bl.banned {
		entries = append(entries, e)
	}
	return entries
}

// Middleware returns HTTP middleware that blocks banned IPs
func (bl *IPBanList) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		if banned, reason := bl.IsBanned(ip); banned {
			http.Error(w, fmt.Sprintf("Access denied: %s", reason), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// TwoFAManager manages TOTP-based two-factor authentication
type TwoFAManager struct {
	secrets map[string]string // userID -> base32 secret
	mu      sync.RWMutex
}

// NewTwoFAManager creates a 2FA manager
func NewTwoFAManager() *TwoFAManager {
	return &TwoFAManager{
		secrets: make(map[string]string),
	}
}

// GenerateSecret creates a new TOTP secret
func (m *TwoFAManager) GenerateSecret(userID string) string {
	secret := make([]byte, 20)
	rand.Read(secret)
	encoded := base64.StdEncoding.EncodeToString(secret)[:20]
	m.mu.Lock()
	m.secrets[userID] = encoded
	m.mu.Unlock()
	return encoded
}

// VerifyCode verifies a TOTP code (simplified - time-based)
func (m *TwoFAManager) VerifyCode(userID, code string) bool {
	m.mu.RLock()
	secret, exists := m.secrets[userID]
	m.mu.RUnlock()

	if !exists || len(code) != 6 {
		return false
	}

	// TOTP: HMAC-SHA1(secret, time/30)
	timeStep := time.Now().Unix() / 30
	for offset := int64(-1); offset <= 1; offset++ {
		expected := generateTOTP(secret, timeStep+offset)
		if expected == code {
			return true
		}
	}
	return false
}

// IsEnabled checks if 2FA is enabled for a user
func (m *TwoFAManager) IsEnabled(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.secrets[userID]
	return exists
}

// Disable removes 2FA for a user
func (m *TwoFAManager) Disable(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.secrets, userID)
}

func generateTOTP(secret string, timeStep int64) string {
	key := []byte(secret)
	msg := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		msg[i] = byte(timeStep & 0xFF)
		timeStep >>= 8
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0F
	truncated := (uint32(hash[offset])&0x7F)<<24 |
		uint32(hash[offset+1])<<16 |
		uint32(hash[offset+2])<<8 |
		uint32(hash[offset+3])
	code := truncated % 1000000
	return fmt.Sprintf("%06d", code)
}

func extractIP(r *http.Request) string {
	// Check X-Forwarded-For (trust proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
