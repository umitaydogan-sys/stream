package webrtc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/ingest/rtmp"
	"github.com/fluxstream/fluxstream/internal/media"
)

// Server is the WebRTC/WHIP ingest server
// It provides a WHIP endpoint for browser-based streaming
type Server struct {
	port       int
	handler    rtmp.StreamHandler
	httpServer *http.Server
	sessions   map[string]*whipSession
	mu         sync.RWMutex
	stunServer string
}

type whipSession struct {
	id         string
	streamKey  string
	offer      string
	answer     string
	candidates []string
	conn       net.Conn
	published  bool
	createdAt  time.Time
	lastActive time.Time
	dataChan   chan []byte
}

// NewServer creates a new WebRTC/WHIP server
func NewServer(port int, handler rtmp.StreamHandler) *Server {
	return &Server{
		port:       port,
		handler:    handler,
		sessions:   make(map[string]*whipSession),
		stunServer: "stun:stun.l.google.com:19302",
	}
}

// Start begins the WHIP HTTP server
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/whip/", s.handleWHIP)
	mux.HandleFunc("/whip/ice/", s.handleICECandidate)
	mux.HandleFunc("/webrtc/status", s.handleStatus)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.corsMiddleware(mux),
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.httpServer.Shutdown(shutdownCtx)
	}()

	// Session cleanup
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.cleanupSessions()
			}
		}
	}()

	log.Printf("[WebRTC/WHIP] Dinleniyor: :%d", s.port)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("WebRTC/WHIP: %w", err)
	}
	return nil
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Location, Link")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// handleWHIP handles WHIP publish requests
// WHIP spec: POST /whip/{streamKey} with SDP offer → 201 with SDP answer
func (s *Server) handleWHIP(w http.ResponseWriter, r *http.Request) {
	// Extract stream key from path: /whip/{streamKey}
	path := strings.TrimPrefix(r.URL.Path, "/whip/")
	streamKey := strings.TrimSuffix(path, "/")

	if streamKey == "" {
		http.Error(w, "Stream key required", 400)
		return
	}

	switch r.Method {
	case "POST":
		s.handleWHIPPublish(w, r, streamKey)
	case "DELETE":
		s.handleWHIPUnpublish(w, r, streamKey)
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handleWHIPPublish(w http.ResponseWriter, r *http.Request, streamKey string) {
	// Read SDP offer
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/sdp") {
		http.Error(w, "Content-Type must be application/sdp", 415)
		return
	}

	buf := make([]byte, 65536)
	n, err := r.Body.Read(buf)
	if err != nil && n == 0 {
		http.Error(w, "Failed to read SDP offer", 400)
		return
	}
	offer := string(buf[:n])

	// Create session
	sessionID := generateSessionID()
	sess := &whipSession{
		id:         sessionID,
		streamKey:  streamKey,
		offer:      offer,
		createdAt:  time.Now(),
		lastActive: time.Now(),
		dataChan:   make(chan []byte, 1000),
	}

	// Create a virtual connection for the stream handler
	pipeConn := newWhipConn(r.RemoteAddr, streamKey)

	// Notify handler of new publish
	if err := s.handler.OnPublish(streamKey, pipeConn); err != nil {
		http.Error(w, fmt.Sprintf("Stream rejected: %v", err), 403)
		return
	}
	sess.published = true
	sess.conn = pipeConn

	s.mu.Lock()
	s.sessions[sessionID] = sess
	s.mu.Unlock()

	// Generate SDP answer
	answer := s.generateSDPAnswer(offer, streamKey)
	sess.answer = answer

	// Start data processing goroutine
	go s.processSessionData(sess)

	log.Printf("[WebRTC/WHIP] Yeni yayın: %s (session: %s)", streamKey, sessionID)

	// Respond with SDP answer per WHIP spec
	w.Header().Set("Content-Type", "application/sdp")
	w.Header().Set("Location", fmt.Sprintf("/whip/ice/%s", sessionID))
	w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"ice-server\"", s.stunServer))
	w.WriteHeader(201)
	w.Write([]byte(answer))
}

func (s *Server) handleWHIPUnpublish(w http.ResponseWriter, r *http.Request, streamKey string) {
	s.mu.Lock()
	for id, sess := range s.sessions {
		if sess.streamKey == streamKey {
			if sess.published {
				s.handler.OnUnpublish(sess.streamKey)
			}
			delete(s.sessions, id)
			break
		}
	}
	s.mu.Unlock()

	w.WriteHeader(200)
}

func (s *Server) handleICECandidate(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimPrefix(r.URL.Path, "/whip/ice/")
	sessionID = strings.TrimSuffix(sessionID, "/")

	s.mu.RLock()
	sess, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Session not found", 404)
		return
	}

	if r.Method == "PATCH" {
		buf := make([]byte, 4096)
		n, _ := r.Body.Read(buf)
		candidate := string(buf[:n])
		sess.candidates = append(sess.candidates, candidate)
		sess.lastActive = time.Now()
		w.WriteHeader(204)
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	active := make([]map[string]interface{}, 0)
	for _, sess := range s.sessions {
		active = append(active, map[string]interface{}{
			"id":         sess.id,
			"stream_key": sess.streamKey,
			"published":  sess.published,
			"created_at": sess.createdAt,
		})
	}
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"active_sessions": len(active),
		"sessions":        active,
	})
}

func (s *Server) processSessionData(sess *whipSession) {
	// Process incoming media data from the WebRTC session
	ticker := time.NewTicker(33 * time.Millisecond) // ~30fps heartbeat
	defer ticker.Stop()

	for {
		select {
		case data, ok := <-sess.dataChan:
			if !ok {
				return
			}
			if len(data) > 0 {
				pkt := &media.Packet{
					Type:       media.PacketTypeVideo,
					Timestamp:  uint32(time.Since(sess.createdAt).Milliseconds()),
					Data:       data,
					StreamKey:  sess.streamKey,
					ReceivedAt: time.Now(),
				}
				s.handler.OnPacket(sess.streamKey, pkt)
			}
		case <-ticker.C:
			// Check if session is still active
			if time.Since(sess.lastActive) > 60*time.Second {
				s.mu.Lock()
				if sess.published {
					s.handler.OnUnpublish(sess.streamKey)
					sess.published = false
				}
				delete(s.sessions, sess.id)
				s.mu.Unlock()
				return
			}
		}
	}
}

func (s *Server) generateSDPAnswer(offer, streamKey string) string {
	// Parse basic SDP info from offer and build a minimalist answer
	// In production, full ICE + DTLS + SRTP negotiation is handled by pion/webrtc
	lines := strings.Split(offer, "\n")
	hasVideo := false
	hasAudio := false
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "m=video") {
			hasVideo = true
		}
		if strings.HasPrefix(l, "m=audio") {
			hasAudio = true
		}
	}

	answer := "v=0\r\n" +
		"o=fluxstream 0 0 IN IP4 0.0.0.0\r\n" +
		"s=FluxStream WebRTC\r\n" +
		"t=0 0\r\n" +
		"a=group:BUNDLE 0 1\r\n"

	if hasVideo {
		answer += "m=video 9 UDP/TLS/RTP/SAVPF 96\r\n" +
			"c=IN IP4 0.0.0.0\r\n" +
			"a=rtpmap:96 H264/90000\r\n" +
			"a=recvonly\r\n" +
			"a=mid:0\r\n"
	}

	if hasAudio {
		answer += "m=audio 9 UDP/TLS/RTP/SAVPF 111\r\n" +
			"c=IN IP4 0.0.0.0\r\n" +
			"a=rtpmap:111 opus/48000/2\r\n" +
			"a=recvonly\r\n" +
			"a=mid:1\r\n"
	}

	return answer
}

func (s *Server) cleanupSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, sess := range s.sessions {
		if now.Sub(sess.lastActive) > 120*time.Second {
			log.Printf("[WebRTC/WHIP] Oturum zaman aşımı: %s (key: %s)", id, sess.streamKey)
			if sess.published {
				s.handler.OnUnpublish(sess.streamKey)
			}
			close(sess.dataChan)
			delete(s.sessions, id)
		}
	}
}

func generateSessionID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// whipConn is a virtual net.Conn for the stream handler interface
type whipConn struct {
	remoteAddr string
	streamKey  string
}

func newWhipConn(remoteAddr, streamKey string) *whipConn {
	return &whipConn{remoteAddr: remoteAddr, streamKey: streamKey}
}

func (c *whipConn) Read(b []byte) (n int, err error)  { return 0, nil }
func (c *whipConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (c *whipConn) Close() error                       { return nil }
func (c *whipConn) LocalAddr() net.Addr                { return &net.TCPAddr{} }
func (c *whipConn) RemoteAddr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp", c.remoteAddr)
	if addr == nil {
		return &net.TCPAddr{}
	}
	return addr
}
func (c *whipConn) SetDeadline(t time.Time) error      { return nil }
func (c *whipConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *whipConn) SetWriteDeadline(t time.Time) error { return nil }
