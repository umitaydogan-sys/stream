package web

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/analytics"
	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/storage"
	"github.com/fluxstream/fluxstream/internal/stream"
	"github.com/fluxstream/fluxstream/internal/tlsutil"
)

// Server handles HTTP requests
type Server struct {
	port            int
	db              *storage.SQLiteDB
	cfg             *config.Manager
	streamMgr       *stream.Manager
	hlsOutputDir    string
	hlsOverrideDir  string
	dashOutputDir   string
	dashOverrideDir string
	dataDir         string
	mux             *http.ServeMux
	startTime       time.Time
	sessions        sync.Map // sessionToken -> *session
	analytics       *analytics.Tracker
	playbackAuth    func(r *http.Request, streamKey, format string) (bool, int, string)
}

type session struct {
	UserID   int64
	Username string
	Role     string
	Expiry   time.Time
}

// NewServer creates a new web server
func NewServer(port int, db *storage.SQLiteDB, cfg *config.Manager, streamMgr *stream.Manager, hlsOutputDir, dataDir string) *Server {
	s := &Server{
		port:         port,
		db:           db,
		cfg:          cfg,
		streamMgr:    streamMgr,
		hlsOutputDir: hlsOutputDir,
		dataDir:      dataDir,
		mux:          http.NewServeMux(),
		startTime:    time.Now(),
	}
	// Fix corrupted user records from older versions
	if fixed, err := db.FixCorruptedUsers(); err == nil && fixed > 0 {
		log.Printf("[WEB] %d corrupted user(s) removed â€” setup wizard will run again", fixed)
	}
	if err := s.ensureDefaultPlayerTemplates(); err != nil {
		log.Printf("[WEB] default player templates could not be ensured: %v", err)
	}
	s.setupRoutes()
	return s
}

// Start begins serving HTTP
func (s *Server) Start(ctx context.Context) error {
	webTLSSource, err := tlsutil.NewSource(s.cfg, tlsutil.ProfileWeb, s.dataDir)
	if err != nil {
		return err
	}
	streamTLSSource, err := tlsutil.NewSource(s.cfg, tlsutil.ProfileStream, s.dataDir)
	if err != nil {
		return err
	}

	baseHandler := s.corsMiddleware(s.mux)
	challengeFallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "FluxStream ACME endpoint\n")
	})
	challengeHandler := tlsutil.ChallengeHandler(challengeFallback, webTLSSource, streamTLSSource)
	httpHandler := baseHandler
	if s.port == 80 && (webTLSSource.UsesLetsEncrypt() || streamTLSSource.UsesLetsEncrypt()) {
		httpHandler = tlsutil.ChallengeHandler(baseHandler, webTLSSource, streamTLSSource)
	}

	httpAddr := fmt.Sprintf(":%d", s.port)
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: httpHandler,
	}

	var httpsServer *http.Server
	httpsAddr := ""
	if webTLSSource.Enabled && webTLSSource.Ready {
		httpsAddr = fmt.Sprintf(":%d", s.cfg.GetInt("https_port", 443))
		httpsServer = &http.Server{
			Addr:      httpsAddr,
			Handler:   baseHandler,
			TLSConfig: webTLSSource.TLSConfig(),
		}
	}

	var acmeServer *http.Server
	if (webTLSSource.UsesLetsEncrypt() || streamTLSSource.UsesLetsEncrypt()) && s.port != 80 {
		acmeServer = &http.Server{
			Addr:    ":80",
			Handler: challengeHandler,
		}
	}

	errCh := make(chan error, 3)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
		if httpsServer != nil {
			_ = httpsServer.Shutdown(shutdownCtx)
		}
		if acmeServer != nil {
			_ = acmeServer.Shutdown(shutdownCtx)
		}
	}()

	go func() {
		log.Printf("[HTTP] Dinleniyor: %s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	if httpsServer != nil {
		go func() {
			log.Printf("[HTTPS] Dinleniyor: %s", httpsAddr)
			httpsListener, err := net.Listen("tcp", httpsAddr)
			if err != nil {
				errCh <- err
				return
			}
			if err := httpsServer.Serve(tls.NewListener(httpsListener, webTLSSource.TLSConfig())); err != nil && err != http.ErrServerClosed {
				errCh <- err
			}
		}()
	} else if webTLSSource.Enabled {
		if webTLSSource.UsesLetsEncrypt() {
			log.Printf("[HTTPS] Let's Encrypt etkin fakat domain henuz hazir degil: %s", webTLSSource.Domain)
		} else {
			log.Printf("[HTTPS] SSL etkin fakat cert/key hazir degil, HTTPS baslatilmadi")
		}
	}

	if acmeServer != nil {
		go func() {
			log.Printf("[ACME] HTTP-01 dinleniyor: :80")
			if err := acmeServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- err
			}
		}()
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

// RegisterHandler registers an external HTTP handler on the web server mux
func (s *Server) RegisterHandler(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

// RegisterAdminHandler registers an HTTP handler that requires an authenticated admin session.
func (s *Server) RegisterAdminHandler(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		sess := s.getSession(r)
		if sess == nil || sess.Role != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			jsonResponse(w, map[string]interface{}{
				"success": false,
				"message": "Admin oturumu gerekli",
			})
			return
		}
		handler(w, r)
	})
}

// RegisterOutputDir registers a static file handler for an output directory
func (s *Server) RegisterOutputDir(pattern string, dir string) {
	s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		filePath := strings.TrimPrefix(r.URL.Path, pattern)
		http.ServeFile(w, r, filepath.Join(dir, filepath.Clean("/"+filePath)))
	})
}

// SetPlaybackAuthorizer sets a shared playback authorization hook.
func (s *Server) SetPlaybackAuthorizer(fn func(r *http.Request, streamKey, format string) (bool, int, string)) {
	s.playbackAuth = fn
}

// SetHLSOverrideDir registers a preferred directory for HLS assets.
func (s *Server) SetHLSOverrideDir(dir string) {
	s.hlsOverrideDir = dir
}

// SetDashOutputDir registers the native DASH output directory.
func (s *Server) SetDashOutputDir(dir string) {
	s.dashOutputDir = dir
}

// SetDashOverrideDir registers a preferred directory for DASH assets.
func (s *Server) SetDashOverrideDir(dir string) {
	s.dashOverrideDir = dir
}

func (s *Server) setupRoutes() {
	// â”€â”€ API Routes â”€â”€
	s.mux.HandleFunc("/api/auth/login", s.handleLogin)
	s.mux.HandleFunc("/api/auth/me", s.handleMe)

	s.mux.HandleFunc("/api/setup/status", s.handleSetupStatus)
	s.mux.HandleFunc("/api/setup/complete", s.handleSetupComplete)

	s.mux.HandleFunc("/api/streams", s.handleStreams)
	s.mux.HandleFunc("/api/streams/", s.handleStreamByID)

	s.mux.HandleFunc("/api/settings", s.handleSettings)
	s.mux.HandleFunc("/api/settings/", s.handleSettingsSection)

	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/logs", s.handleLogs)

	s.mux.HandleFunc("/api/embed/defaults", s.handleEmbedDefaults)
	s.mux.HandleFunc("/api/embed/", s.handleEmbedCodes)

	// â”€â”€ Player Template Routes â”€â”€
	s.mux.HandleFunc("/api/players", s.handlePlayerTemplates)
	s.mux.HandleFunc("/api/players/", s.handlePlayerTemplateByID)

	// â”€â”€ User Management Routes â”€â”€
	s.mux.HandleFunc("/api/users", s.handleUsers)
	s.mux.HandleFunc("/api/users/", s.handleUserByID)

	// â”€â”€ SSL Certificate Upload â”€â”€
	s.mux.HandleFunc("/api/ssl/upload", s.handleSSLUpload)
	s.mux.HandleFunc("/api/ssl/status", s.handleSSLStatus)

	// â”€â”€ Health Check â”€â”€
	s.mux.HandleFunc("/api/health", s.handleHealthCheck)

	// â”€â”€ Media Routes (HLS with proper headers) â”€â”€
	s.mux.HandleFunc("/hls/", s.handleHLS)
	s.mux.HandleFunc("/dash/", s.handleDASH)
	s.mux.Handle("/static/", http.StripPrefix("/static/", embeddedStaticHandler()))

	// â”€â”€ Player Routes â”€â”€
	s.mux.HandleFunc("/play/", s.handlePlayer)
	s.mux.HandleFunc("/embed/", s.handleEmbedPage)

	// â”€â”€ Static Admin UI â”€â”€
	s.mux.HandleFunc("/", s.handleAdmin)
}

// â”€â”€ CORS Middleware â”€â”€

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// â”€â”€ Auth Routes â”€â”€

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request", 400)
		return
	}

	user, err := s.db.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		jsonError(w, "GeÃ§ersiz kullanÄ±cÄ± adÄ± veya ÅŸifre", 401)
		return
	}

	// Simple password check (in production, use bcrypt)
	if user.PasswordHash != hashPassword(req.Password) {
		jsonError(w, "GeÃ§ersiz kullanÄ±cÄ± adÄ± veya ÅŸifre", 401)
		return
	}

	s.db.UpdateUserLogin(user.ID)

	// Create a secure session token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	s.sessions.Store(token, &session{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Expiry:   time.Now().Add(24 * time.Hour),
	})

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
		"token": token,
	})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	sess := s.getSession(r)
	if sess == nil {
		jsonResponse(w, map[string]interface{}{"authenticated": false})
		return
	}
	jsonResponse(w, map[string]interface{}{
		"authenticated": true,
		"user": map[string]interface{}{
			"username": sess.Username,
			"role":     sess.Role,
		},
	})
}

// getSession extracts and validates the session from the Authorization header
func (s *Server) getSession(r *http.Request) *session {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil
	}
	val, ok := s.sessions.Load(token)
	if !ok {
		return nil
	}
	sess := val.(*session)
	if time.Now().After(sess.Expiry) {
		s.sessions.Delete(token)
		return nil
	}
	return sess
}

// â”€â”€ Setup Routes â”€â”€

func (s *Server) handleSetupStatus(w http.ResponseWriter, r *http.Request) {
	completed := s.cfg.GetBool("setup_completed", false)
	jsonResponse(w, map[string]interface{}{
		"setup_completed": completed,
	})
}

func (s *Server) handleSetupComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		HTTPPort    int    `json:"http_port"`
		HTTPSPort   int    `json:"https_port"`
		RTMPPort    int    `json:"rtmp_port"`
		EmbedDomain string `json:"embed_domain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request", 400)
		return
	}

	// Create admin user
	if _, err := s.db.CreateUser(req.Username, hashPassword(req.Password), "admin"); err != nil {
		jsonError(w, fmt.Sprintf("KullanÄ±cÄ± oluÅŸturulamadÄ±: %v", err), 500)
		return
	}

	// Update ports if changed
	if req.HTTPPort > 0 {
		s.cfg.Set("http_port", fmt.Sprintf("%d", req.HTTPPort), "general")
		s.cfg.Set("embed_http_port", fmt.Sprintf("%d", req.HTTPPort), "embed")
	}
	if req.HTTPSPort > 0 {
		s.cfg.Set("https_port", fmt.Sprintf("%d", req.HTTPSPort), "general")
		s.cfg.Set("embed_https_port", fmt.Sprintf("%d", req.HTTPSPort), "embed")
	}
	if req.RTMPPort > 0 {
		s.cfg.Set("rtmp_port", fmt.Sprintf("%d", req.RTMPPort), "general")
	}
	if strings.TrimSpace(req.EmbedDomain) != "" {
		s.cfg.Set("embed_domain", strings.TrimSpace(req.EmbedDomain), "embed")
	}

	// Mark setup as completed
	s.cfg.Set("setup_completed", "true", "general")

	s.db.AddLog("INFO", "setup", "Ä°lk kurulum tamamlandÄ±")

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"message": "Kurulum tamamlandÄ±!",
	})
}

// â”€â”€ Stream Routes â”€â”€

func (s *Server) handleStreams(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		streams, err := s.db.GetAllStreams()
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		if streams == nil {
			streams = []storage.Stream{}
		}
		viewerCounts := map[string]int{}
		if s.analytics != nil {
			viewerCounts = s.analytics.CurrentViewersByStream()
		}
		// Enrich with live status
		for i := range streams {
			if s.analytics != nil {
				s.analytics.RegisterStreamName(streams[i].StreamKey, streams[i].Name)
			}
			if s.streamMgr.IsLive(streams[i].StreamKey) {
				streams[i].Status = "live"
			}
			streams[i].ViewerCount = viewerCounts[streams[i].StreamKey]
		}
		jsonResponse(w, streams)

	case "POST":
		var req struct {
			Name          string `json:"name"`
			Description   string `json:"description"`
			OutputFormats string `json:"output_formats"`
			PolicyJSON    string `json:"policy_json"`
			MaxViewers    int    `json:"max_viewers"`
			MaxBitrate    int    `json:"max_bitrate"`
			RecordEnabled bool   `json:"record_enabled"`
			RecordFormat  string `json:"record_format"`
			Password      string `json:"password"`
			DomainLock    string `json:"domain_lock"`
			IPWhitelist   string `json:"ip_whitelist"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "Invalid request", 400)
			return
		}

		if req.Name == "" {
			jsonError(w, "YayÄ±n adÄ± gerekli", 400)
			return
		}

		if req.OutputFormats == "" {
			req.OutputFormats = `["hls","ll_hls","dash","flv","whep","mp4","webm","mp3","aac","ogg","wav","flac","icecast"]`
		}
		if req.RecordFormat == "" {
			req.RecordFormat = "ts"
		}

		streamKey := generateStreamKey()
		st := &storage.Stream{
			Name:          req.Name,
			Description:   req.Description,
			StreamKey:     streamKey,
			OutputFormats: req.OutputFormats,
			PolicyJSON:    req.PolicyJSON,
			MaxViewers:    req.MaxViewers,
			MaxBitrate:    req.MaxBitrate,
			RecordEnabled: req.RecordEnabled,
			RecordFormat:  req.RecordFormat,
			Password:      req.Password,
			DomainLock:    req.DomainLock,
			IPWhitelist:   req.IPWhitelist,
		}

		id, err := s.db.CreateStream(st)
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		st.ID = id
		st.StreamKey = streamKey
		if s.analytics != nil {
			s.analytics.RegisterStreamName(streamKey, st.Name)
		}

		s.db.AddLog("INFO", "stream", fmt.Sprintf("Yeni yayÄ±n oluÅŸturuldu: %s (key: %s)", req.Name, streamKey))

		base := s.publicBaseURL(r)
		rtmpBase := s.publicRTMPBase(r, false)

		jsonResponse(w, map[string]interface{}{
			"stream":           st,
			"rtmp_url":         rtmpBase + "/live",
			"rtmp_publish_url": rtmpBase + "/live/" + streamKey,
			"hls_url":          base + "/hls/" + streamKey + "/master.m3u8",
			"play_url":         base + "/play/" + streamKey,
		})

	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handleStreamByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL: /api/streams/{id}
	idStr := r.URL.Path[len("/api/streams/"):]
	if idStr == "" {
		jsonError(w, "Stream ID gerekli", 400)
		return
	}

	var id int64
	fmt.Sscanf(idStr, "%d", &id)

	switch r.Method {
	case "GET":
		st, err := s.db.GetStreamByID(id)
		if err != nil || st == nil {
			jsonError(w, "Stream bulunamadÄ±", 404)
			return
		}
		if s.streamMgr.IsLive(st.StreamKey) {
			st.Status = "live"
		}
		if s.analytics != nil {
			s.analytics.RegisterStreamName(st.StreamKey, st.Name)
			st.ViewerCount = s.analytics.CurrentViewersByStream()[st.StreamKey]
		}
		jsonResponse(w, st)

	case "PUT":
		var req storage.Stream
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		if req.RecordFormat == "" {
			req.RecordFormat = "ts"
		}
		req.ID = id
		if err := s.db.UpdateStream(&req); err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		if s.analytics != nil && req.StreamKey != "" && req.Name != "" {
			s.analytics.RegisterStreamName(req.StreamKey, req.Name)
		}
		jsonResponse(w, map[string]interface{}{"success": true})

	case "DELETE":
		if err := s.db.DeleteStream(id); err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		jsonResponse(w, map[string]interface{}{"success": true})

	default:
		http.Error(w, "Method not allowed", 405)
	}
}

// â”€â”€ Settings Routes â”€â”€

func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	configs, err := s.cfg.GetAll()
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonResponse(w, configs)
}

func (s *Server) handleSettingsSection(w http.ResponseWriter, r *http.Request) {
	section := r.URL.Path[len("/api/settings/"):]

	if r.Method == "PUT" {
		var updates map[string]string
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		for key, value := range updates {
			if err := s.cfg.Set(key, value, section); err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
		}
		s.db.AddLog("INFO", "settings", fmt.Sprintf("Ayarlar gÃ¼ncellendi: %s", section))
		jsonResponse(w, map[string]interface{}{"success": true})
		return
	}

	configs, err := s.cfg.GetByCategory(section)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonResponse(w, configs)
}

// â”€â”€ Stats Route â”€â”€

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := s.streamMgr.GetStats()
	if s.analytics != nil {
		dash := s.analytics.GetDashboard()
		stats.TotalViewers = dash.CurrentViewers
		stats.BandwidthOut = dash.TotalBandwidth
	}

	// Real memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	stats.MemoryUsedMB = int64(memStats.Alloc / 1024 / 1024)
	stats.MemoryTotalMB = int64(memStats.Sys / 1024 / 1024)
	stats.UptimeSeconds = int64(time.Since(s.startTime).Seconds())

	jsonResponse(w, stats)
}

// â”€â”€ Logs Route â”€â”€

func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		s.db.ClearLogs()
		jsonResponse(w, map[string]interface{}{"success": true})
		return
	}

	level := r.URL.Query().Get("level")
	component := r.URL.Query().Get("component")
	logs, err := s.db.GetLogs(200, level, component)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	if logs == nil {
		logs = []storage.LogEntry{}
	}
	jsonResponse(w, logs)
}

// â”€â”€ Embed Defaults Route â”€â”€

func (s *Server) handleEmbedDefaults(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var defaults storage.EmbedDefaults
		if err := json.NewDecoder(r.Body).Decode(&defaults); err != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		s.cfg.Set("embed_domain", defaults.Domain, "embed")
		s.cfg.Set("embed_http_port", fmt.Sprintf("%d", defaults.HTTPPort), "embed")
		s.cfg.Set("embed_https_port", fmt.Sprintf("%d", defaults.HTTPSPort), "embed")
		s.cfg.Set("embed_use_https", fmt.Sprintf("%t", defaults.UseHTTPS), "embed")
		jsonResponse(w, map[string]interface{}{"success": true})
		return
	}

	jsonResponse(w, storage.EmbedDefaults{
		Domain:      s.cfg.Get("embed_domain", "localhost"),
		HTTPPort:    s.cfg.GetInt("embed_http_port", 8844),
		HTTPSPort:   s.cfg.GetInt("embed_https_port", 443),
		RTMPPort:    s.cfg.GetInt("rtmp_port", 1935),
		RTSPPort:    s.cfg.GetInt("rtsp_port", 8554),
		SRTPort:     s.cfg.GetInt("srt_port", 9000),
		WebRTCPort:  s.cfg.GetInt("webrtc_port", 8855),
		IcecastPort: s.cfg.GetInt("icecast_port", 8000),
		UseHTTPS:    s.cfg.GetBool("embed_use_https", false),
	})
}

// â”€â”€ Admin UI Handler â”€â”€

func (s *Server) handleAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(adminHTML))
}

// â”€â”€ Helpers â”€â”€

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

// hashPasswordSalt is a fixed salt â€” in production rotate via config.
var hashPasswordSalt = []byte("fluxstream-v1-salt")

func hashPassword(password string) string {
	mac := hmac.New(sha256.New, hashPasswordSalt)
	mac.Write([]byte(password))
	return hex.EncodeToString(mac.Sum(nil))
}

func generateStreamKey() string {
	b := make([]byte, 12)
	rand.Read(b)
	return "live_" + hex.EncodeToString(b)
}

// â”€â”€ HLS Handler with proper headers â”€â”€

func (s *Server) handleHLS(w http.ResponseWriter, r *http.Request) {
	// Set proper CORS and cache headers for HLS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/hls/")
	cleanPath := filepath.Clean(path)
	if cleanPath == "." || strings.HasPrefix(cleanPath, "..") {
		http.Error(w, "Forbidden", 403)
		return
	}
	streamKey := ""
	if parts := strings.Split(strings.Trim(cleanPath, "/"), string(filepath.Separator)); len(parts) > 0 {
		streamKey = parts[0]
	}

	filePath := filepath.Join(s.hlsOutputDir, cleanPath)
	if s.hlsOverrideDir != "" {
		overridePath := filepath.Join(s.hlsOverrideDir, cleanPath)
		if override := newestPlaylistVariant(overridePath); override != "" {
			filePath = override
		} else if info, err := os.Stat(overridePath); err == nil && !info.IsDir() {
			filePath = overridePath
		} else if strings.HasSuffix(cleanPath, string(filepath.Separator)+"master.m3u8") || strings.HasSuffix(cleanPath, "/master.m3u8") {
			fallbackPath := filepath.Join(s.hlsOverrideDir, strings.TrimSuffix(cleanPath, "master.m3u8")+"index.m3u8")
			if override := newestPlaylistVariant(fallbackPath); override != "" {
				filePath = override
			} else if info, err := os.Stat(fallbackPath); err == nil && !info.IsDir() {
				filePath = fallbackPath
			}
		} else if strings.HasSuffix(cleanPath, string(filepath.Separator)+"ll.m3u8") {
			fallbackPath := filepath.Join(s.hlsOverrideDir, strings.TrimSuffix(cleanPath, "ll.m3u8")+"index.m3u8")
			if override := newestPlaylistVariant(fallbackPath); override != "" {
				filePath = override
			} else if info, err := os.Stat(fallbackPath); err == nil && !info.IsDir() {
				filePath = fallbackPath
			}
		}
	}

	if strings.HasSuffix(path, ".m3u8") {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		w.Header().Set("Cache-Control", "no-cache, no-store")
	} else if strings.HasSuffix(path, ".ts") {
		w.Header().Set("Content-Type", "video/mp2t")
		w.Header().Set("Cache-Control", "public, max-age=30")
	}
	format := "hls"
	switch {
	case strings.HasSuffix(path, "ll.m3u8"), strings.Contains(path, "/ll_"), strings.Contains(path, "\\ll_"):
		format = "ll_hls"
	case strings.HasSuffix(path, "audio.m3u8"):
		format = "hls_audio"
	}
	if !s.authorizePlayback(w, r, streamKey, format) {
		return
	}

	s.trackPlaybackHeartbeat(w, r, streamKey, format, func(rw http.ResponseWriter) {
		if strings.HasSuffix(path, ".m3u8") && serveManifestWithPassthrough(rw, r, filePath, "application/vnd.apple.mpegurl", rewriteHLSManifest) {
			return
		}
		http.ServeFile(rw, r, filePath)
	})
}

// â”€â”€ Player Page Handler â”€â”€

func (s *Server) handleDASH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/dash/")
	cleanPath := filepath.Clean(path)
	if cleanPath == "." || strings.HasPrefix(cleanPath, "..") {
		http.Error(w, "Forbidden", 403)
		return
	}
	streamKey := ""
	if parts := strings.Split(strings.Trim(cleanPath, "/"), string(filepath.Separator)); len(parts) > 0 {
		streamKey = parts[0]
	}

	filePath := ""
	if s.dashOverrideDir != "" {
		overridePath := filepath.Join(s.dashOverrideDir, cleanPath)
		if override := newestManifestVariant(overridePath); override != "" {
			filePath = override
		} else if info, err := os.Stat(overridePath); err == nil && !info.IsDir() {
			filePath = overridePath
		}
	}
	if filePath == "" && s.dashOutputDir != "" {
		nativePath := filepath.Join(s.dashOutputDir, cleanPath)
		if override := newestManifestVariant(nativePath); override != "" {
			filePath = override
		} else {
			filePath = nativePath
		}
	}

	switch {
	case strings.HasSuffix(cleanPath, ".mpd"):
		w.Header().Set("Content-Type", "application/dash+xml")
		w.Header().Set("Cache-Control", "no-cache, no-store")
	case strings.HasSuffix(cleanPath, ".m4s"), strings.HasSuffix(cleanPath, ".mp4"):
		w.Header().Set("Content-Type", "video/iso.segment")
		w.Header().Set("Cache-Control", "public, max-age=30")
	}
	format := "dash"
	if strings.HasSuffix(cleanPath, "audio.mpd") {
		format = "dash_audio"
	}
	if !s.authorizePlayback(w, r, streamKey, format) {
		return
	}
	s.trackPlaybackHeartbeat(w, r, streamKey, format, func(rw http.ResponseWriter) {
		if strings.HasSuffix(cleanPath, ".mpd") && serveManifestWithPassthrough(rw, r, filePath, "application/dash+xml", rewriteDASHManifest) {
			return
		}
		http.ServeFile(rw, r, filePath)
	})
}

func newestManifestVariant(path string) string {
	if !strings.HasSuffix(path, ".m3u8") && !strings.HasSuffix(path, ".mpd") {
		return ""
	}
	newest := ""
	var newestTime time.Time
	for _, candidate := range []string{path + ".tmp", path} {
		info, err := os.Stat(candidate)
		if err != nil || info.IsDir() {
			continue
		}
		if newest == "" || info.ModTime().After(newestTime) {
			newest = candidate
			newestTime = info.ModTime()
		}
	}
	return newest
}

func newestPlaylistVariant(path string) string {
	return newestManifestVariant(path)
}

func (s *Server) handlePlayer(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/play/"):]
	if key == "" {
		http.Error(w, "Stream key gerekli", 400)
		return
	}
	if !s.authorizePlayback(w, r, key, "hls") {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, playerHTML, key, key, key)
}

// â”€â”€ Embed Page Handler â”€â”€

func (s *Server) handleEmbedPage(w http.ResponseWriter, r *http.Request) {
	// /embed/{key} or /embed/{key}/audio
	path := r.URL.Path[len("/embed/"):]
	key := strings.Split(path, "/")[0]
	if key == "" {
		http.Error(w, "Stream key gerekli", 400)
		return
	}

	format := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("format")))
	if format == "" || format == "player" || format == "iframe" || format == "jsapi" {
		format = "hls"
	}
	autoplay := r.URL.Query().Get("autoplay") != "0"
	muted := r.URL.Query().Get("muted") != "0"
	if !s.authorizePlayback(w, r, key, format) {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, renderEmbedHTML(key, format, autoplay, muted))
}

// â”€â”€ Embed Codes API â”€â”€

func (s *Server) handleEmbedCodes(w http.ResponseWriter, r *http.Request) {
	// /api/embed/{id}
	idStr := r.URL.Path[len("/api/embed/"):]
	if idStr == "" || idStr == "defaults" {
		return
	}

	var id int64
	fmt.Sscanf(idStr, "%d", &id)

	st, err := s.db.GetStreamByID(id)
	if err != nil || st == nil {
		jsonError(w, "Stream bulunamadÄ±", 404)
		return
	}

	domain := s.publicHost(r)
	rtspPort := s.cfg.GetInt("rtsp_port", 8554)
	srtPort := s.cfg.GetInt("srt_port", 9000)
	base := s.publicBaseURL(r)

	rtpPort := s.cfg.GetInt("rtp_port", 5004)
	mpegtsPort := s.cfg.GetInt("mpegts_port", 9001)
	rtspOutPort := s.cfg.GetInt("rtsp_out_port", 8555)
	srtOutPort := s.cfg.GetInt("srt_out_port", 9010)
	codes := map[string]interface{}{
		"stream": st,
		"urls": map[string]string{
			"hls":        base + "/hls/" + st.StreamKey + "/master.m3u8",
			"hls_media":  base + "/hls/" + st.StreamKey + "/index.m3u8",
			"ll_hls":     base + "/hls/" + st.StreamKey + "/ll.m3u8",
			"dash":       base + "/dash/" + st.StreamKey + "/manifest.mpd",
			"http_flv":   base + "/flv/" + st.StreamKey,
			"whep":       base + "/whep/play/" + st.StreamKey,
			"fmp4":       mediaNamedURL(base, "/mp4", st.StreamKey, st.Name, "mp4"),
			"webm":       mediaNamedURL(base, "/webm", st.StreamKey, st.Name, "webm"),
			"play":       base + "/play/" + st.StreamKey,
			"embed":      base + "/embed/" + st.StreamKey,
			"rtmp":       fmt.Sprintf("%s/live/%s", s.publicRTMPBase(r, false), st.StreamKey),
			"rtsp":       fmt.Sprintf("rtsp://%s:%d/live/%s", domain, rtspPort, st.StreamKey),
			"srt":        fmt.Sprintf("srt://%s:%d?streamid=%s", domain, srtPort, st.StreamKey),
			"rtp":        fmt.Sprintf("rtp://%s:%d", domain, rtpPort),
			"mpegts":     fmt.Sprintf("udp://%s:%d/%s", domain, mpegtsPort, st.StreamKey),
			"rtsp_out":   fmt.Sprintf("rtsp://%s:%d/live/%s", domain, rtspOutPort, st.StreamKey),
			"srt_out":    fmt.Sprintf("srt://%s:%d?streamid=%s", domain, srtOutPort, st.StreamKey),
			"mp3":        mediaNamedURL(base, "/audio/mp3", st.StreamKey, st.Name, "mp3"),
			"aac":        mediaNamedURL(base, "/audio/aac", st.StreamKey, st.Name, "aac"),
			"ogg":        mediaNamedURL(base, "/audio/ogg", st.StreamKey, st.Name, "ogg"),
			"wav":        mediaNamedURL(base, "/audio/wav", st.StreamKey, st.Name, "wav"),
			"flac":       mediaNamedURL(base, "/audio/flac", st.StreamKey, st.Name, "flac"),
			"hls_audio":  base + "/audio/hls/" + st.StreamKey,
			"dash_audio": base + "/audio/dash/" + st.StreamKey,
			"icecast":    base + "/icecast/" + st.StreamKey,
			"asset_hls":  base + "/static/vendor/hls.min.js",
			"asset_dash": base + "/static/vendor/dash.all.min.js",
			"asset_flv":  base + "/static/vendor/mpegts.min.js",
		},
		"embed_code": fmt.Sprintf(`<iframe src="%s/embed/%s" width="1280" height="720" frameborder="0" allowfullscreen></iframe>`, base, st.StreamKey),
	}

	jsonResponse(w, codes)
}

// â”€â”€ Health Check â”€â”€

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, map[string]interface{}{
		"status":  "ok",
		"version": "2.0.0",
		"uptime":  int64(time.Since(s.startTime).Seconds()),
	})
}

// â”€â”€ Player Template Routes â”€â”€

func (s *Server) handlePlayerTemplates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := s.ensureDefaultPlayerTemplates(); err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		templates, err := s.db.GetPlayerTemplates()
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		if templates == nil {
			templates = []storage.PlayerTemplate{}
		}
		jsonResponse(w, templates)

	case "POST":
		var pt storage.PlayerTemplate
		if err := json.NewDecoder(r.Body).Decode(&pt); err != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		if pt.Name == "" {
			jsonError(w, "Template adÄ± gerekli", 400)
			return
		}
		id, err := s.db.CreatePlayerTemplate(&pt)
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		pt.ID = id
		s.db.AddLog("INFO", "player", fmt.Sprintf("Yeni player template: %s", pt.Name))
		jsonResponse(w, pt)

	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handlePlayerTemplateByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/players/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonError(w, "GeÃ§ersiz ID", 400)
		return
	}

	switch r.Method {
	case "GET":
		pt, err := s.db.GetPlayerTemplateByID(id)
		if err != nil || pt == nil {
			jsonError(w, "Template bulunamadÄ±", 404)
			return
		}
		jsonResponse(w, pt)

	case "PUT":
		var pt storage.PlayerTemplate
		if err := json.NewDecoder(r.Body).Decode(&pt); err != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		pt.ID = id
		if err := s.db.UpdatePlayerTemplate(&pt); err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		jsonResponse(w, map[string]interface{}{"success": true})

	case "DELETE":
		if err := s.db.DeletePlayerTemplate(id); err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		jsonResponse(w, map[string]interface{}{"success": true})

	default:
		http.Error(w, "Method not allowed", 405)
	}
}

// â”€â”€ User Management Routes â”€â”€

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, err := s.db.GetUsers()
		if err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		if users == nil {
			users = []storage.User{}
		}
		jsonResponse(w, users)

	case "POST":
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		if req.Username == "" || req.Password == "" {
			jsonError(w, "KullanÄ±cÄ± adÄ± ve ÅŸifre gerekli", 400)
			return
		}
		if req.Role == "" {
			req.Role = "viewer"
		}
		id, err := s.db.CreateUser(req.Username, hashPassword(req.Password), req.Role)
		if err != nil {
			jsonError(w, fmt.Sprintf("KullanÄ±cÄ± oluÅŸturulamadÄ±: %v", err), 500)
			return
		}
		s.db.AddLog("INFO", "users", fmt.Sprintf("Yeni kullanÄ±cÄ±: %s (%s)", req.Username, req.Role))
		jsonResponse(w, map[string]interface{}{
			"success": true,
			"id":      id,
		})

	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/users/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonError(w, "GeÃ§ersiz ID", 400)
		return
	}

	switch r.Method {
	case "GET":
		user, err := s.db.GetUserByID(id)
		if err != nil || user == nil {
			jsonError(w, "KullanÄ±cÄ± bulunamadÄ±", 404)
			return
		}
		jsonResponse(w, map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})

	case "PUT":
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "Invalid request", 400)
			return
		}
		if req.Username != "" || req.Role != "" {
			if err := s.db.UpdateUser(id, req.Username, req.Role); err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
		}
		if req.Password != "" {
			if err := s.db.UpdateUserPassword(id, hashPassword(req.Password)); err != nil {
				jsonError(w, err.Error(), 500)
				return
			}
		}
		jsonResponse(w, map[string]interface{}{"success": true})

	case "DELETE":
		if err := s.db.DeleteUser(id); err != nil {
			jsonError(w, err.Error(), 500)
			return
		}
		s.db.AddLog("INFO", "users", fmt.Sprintf("KullanÄ±cÄ± silindi: ID %d", id))
		jsonResponse(w, map[string]interface{}{"success": true})

	default:
		http.Error(w, "Method not allowed", 405)
	}
}

// â”€â”€ SSL Certificate Upload â”€â”€

func (s *Server) handleSSLUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// Max 10MB
	r.ParseMultipartForm(10 << 20)

	certFile, _, err := r.FormFile("cert")
	if err != nil {
		jsonError(w, "Sertifika dosyasi gerekli", 400)
		return
	}
	defer certFile.Close()
	certBytes, err := io.ReadAll(certFile)
	if err != nil || len(certBytes) == 0 {
		jsonError(w, "Sertifika okunamadi", 400)
		return
	}

	keyFile, _, err := r.FormFile("key")
	if err != nil {
		jsonError(w, "Anahtar dosyasi gerekli", 400)
		return
	}
	defer keyFile.Close()
	keyBytes, err := io.ReadAll(keyFile)
	if err != nil || len(keyBytes) == 0 {
		jsonError(w, "Anahtar okunamadi", 400)
		return
	}

	if _, err := tls.X509KeyPair(certBytes, keyBytes); err != nil {
		jsonError(w, fmt.Sprintf("CRT/KEY cifti gecersiz: %v", err), 400)
		return
	}

	target := strings.ToLower(strings.TrimSpace(r.FormValue("target")))
	profile := tlsutil.ProfileWeb
	certKeyPrefix := "ssl"
	certsDir := filepath.Join(s.dataDir, "certs", "web")
	if target == "stream" {
		profile = tlsutil.ProfileStream
		certKeyPrefix = "stream_ssl"
		certsDir = filepath.Join(s.dataDir, "certs", "stream")
	}
	if err := os.MkdirAll(certsDir, 0755); err != nil {
		jsonError(w, fmt.Sprintf("Sertifika klasoru olusturulamadi: %v", err), 500)
		return
	}
	certPath := filepath.Join(certsDir, "server.crt")
	keyPath := filepath.Join(certsDir, "server.key")

	if err := os.WriteFile(certPath, certBytes, 0644); err != nil {
		jsonError(w, fmt.Sprintf("Sertifika yazilamadi: %v", err), 500)
		return
	}
	if err := os.WriteFile(keyPath, keyBytes, 0600); err != nil {
		jsonError(w, fmt.Sprintf("Anahtar yazilamadi: %v", err), 500)
		return
	}

	_ = s.cfg.Set(certKeyPrefix+"_cert_path", certPath, "ssl")
	_ = s.cfg.Set(certKeyPrefix+"_key_path", keyPath, "ssl")
	if profile == tlsutil.ProfileWeb {
		_ = s.cfg.Set("ssl_mode", "file", "ssl")
	} else {
		_ = s.cfg.Set("stream_ssl_mode", "file", "ssl")
	}

	s.db.AddLog("INFO", "ssl", fmt.Sprintf("%s SSL sertifikalari yuklendi", strings.ToUpper(string(profile))))

	jsonResponse(w, map[string]interface{}{
		"success":          true,
		"target":           string(profile),
		"cert_path":        certPath,
		"key_path":         keyPath,
		"requires_restart": true,
	})
}

func (s *Server) handleSSLStatus(w http.ResponseWriter, r *http.Request) {
	webSource, err := tlsutil.NewSource(s.cfg, tlsutil.ProfileWeb, s.dataDir)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	streamSource, err := tlsutil.NewSource(s.cfg, tlsutil.ProfileStream, s.dataDir)
	if err != nil {
		jsonError(w, err.Error(), 500)
		return
	}
	jsonResponse(w, map[string]interface{}{
		"web": map[string]interface{}{
			"enabled":          s.cfg.GetBool("ssl_enabled", false),
			"mode":             s.cfg.Get("ssl_mode", "file"),
			"domain":           s.cfg.Get("ssl_le_domain", ""),
			"email":            s.cfg.Get("ssl_le_email", ""),
			"cert_path":        webSource.CertPath,
			"key_path":         webSource.KeyPath,
			"has_cert":         fileExists(webSource.CertPath),
			"has_key":          fileExists(webSource.KeyPath),
			"ready":            webSource.Ready,
			"https_port":       s.cfg.GetInt("https_port", 443),
			"requires_restart": true,
		},
		"stream": map[string]interface{}{
			"enabled":          s.cfg.GetBool("rtmps_enabled", false),
			"mode":             s.cfg.Get("stream_ssl_mode", "file"),
			"domain":           s.cfg.Get("stream_ssl_le_domain", ""),
			"email":            s.cfg.Get("stream_ssl_le_email", ""),
			"cert_path":        streamSource.CertPath,
			"key_path":         streamSource.KeyPath,
			"has_cert":         fileExists(streamSource.CertPath),
			"has_key":          fileExists(streamSource.KeyPath),
			"ready":            streamSource.Ready,
			"rtmps_port":       s.cfg.GetInt("rtmps_port", 1936),
			"requires_restart": true,
		},
		"requires_restart": true,
	})
}

// â”€â”€ Disk Usage Helper â”€â”€

func mediaNamedURL(base, prefix, streamKey, streamName, ext string) string {
	fileBase := slugifyFileName(streamName)
	if fileBase == "" {
		fileBase = streamKey
	}
	return fmt.Sprintf("%s%s/%s/%s.%s", base, prefix, streamKey, fileBase, ext)
}

func (s *Server) publicHost(r *http.Request) string {
	configured := strings.TrimSpace(s.cfg.Get("embed_domain", ""))
	if configured != "" && !strings.EqualFold(configured, "localhost") {
		return configured
	}
	if r != nil {
		if host := requestHostName(r.Host); host != "" && !strings.EqualFold(host, "localhost") {
			return host
		}
	}
	return "localhost"
}

func (s *Server) publicUseHTTPS() bool {
	if !s.cfg.GetBool("embed_use_https", false) {
		return false
	}
	webSource, err := tlsutil.NewSource(s.cfg, tlsutil.ProfileWeb, s.dataDir)
	if err != nil {
		return false
	}
	return webSource.Ready
}

func (s *Server) publicBaseURL(r *http.Request) string {
	scheme := "http"
	port := s.cfg.GetInt("embed_http_port", s.cfg.GetInt("http_port", 8844))
	if s.publicUseHTTPS() {
		scheme = "https"
		port = s.cfg.GetInt("embed_https_port", s.cfg.GetInt("https_port", 443))
	}
	portStr := ""
	if (scheme == "http" && port != 80) || (scheme == "https" && port != 443) {
		portStr = fmt.Sprintf(":%d", port)
	}
	return fmt.Sprintf("%s://%s%s", scheme, s.publicHost(r), portStr)
}

func (s *Server) publicRTMPBase(r *http.Request, secure bool) string {
	scheme := "rtmp"
	port := s.cfg.GetInt("rtmp_port", 1935)
	if secure {
		scheme = "rtmps"
		port = s.cfg.GetInt("rtmps_port", 1936)
	}
	return fmt.Sprintf("%s://%s:%d", scheme, s.publicHost(r), port)
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
	return hostPort
}

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func slugifyFileName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return ""
	}
	var b strings.Builder
	lastDash := false
	for _, r := range name {
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
	return strings.Trim(b.String(), "-")
}

func dirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

func (s *Server) authorizePlayback(w http.ResponseWriter, r *http.Request, streamKey, format string) bool {
	if s.playbackAuth == nil {
		return true
	}
	ok, status, message := s.playbackAuth(r, streamKey, format)
	if ok {
		return true
	}
	if status <= 0 {
		status = http.StatusForbidden
	}
	w.WriteHeader(status)
	jsonResponse(w, map[string]interface{}{
		"success": false,
		"message": message,
	})
	return false
}
