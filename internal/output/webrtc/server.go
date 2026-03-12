package webrtc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/stream"
)

// WHEPServer implements WebRTC-HTTP Egress Protocol for sub-second playback
type WHEPServer struct {
	manager  *stream.Manager
	sessions map[string]*whepSession
	mu       sync.RWMutex
}

type whepSession struct {
	id        string
	streamKey string
	subID     string
	sdpOffer  string
	sdpAnswer string
	createdAt time.Time
	candidates []string
}

// NewWHEPServer creates a new WHEP output server
func NewWHEPServer(manager *stream.Manager) *WHEPServer {
	return &WHEPServer{
		manager:  manager,
		sessions: make(map[string]*whepSession),
	}
}

// HandleWHEP handles WHEP endpoint requests
func (s *WHEPServer) HandleWHEP(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Expose-Headers", "Location, Link")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := r.URL.Path

	// POST /whep/play/{key} — create session
	if r.Method == http.MethodPost && strings.HasPrefix(path, "/whep/play/") {
		s.handleOffer(w, r)
		return
	}

	// PATCH /whep/ice/{sessionID} — ICE candidate
	if r.Method == http.MethodPatch && strings.HasPrefix(path, "/whep/ice/") {
		s.handleICE(w, r)
		return
	}

	// DELETE /whep/session/{sessionID} — teardown
	if r.Method == http.MethodDelete && strings.HasPrefix(path, "/whep/session/") {
		s.handleTeardown(w, r)
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

func (s *WHEPServer) handleOffer(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/whep/play/")
	key = strings.TrimRight(key, "/")

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	// Read SDP offer
	buf := make([]byte, 65536)
	n, err := r.Body.Read(buf)
	if err != nil && n == 0 {
		http.Error(w, "Invalid SDP offer", http.StatusBadRequest)
		return
	}
	sdpOffer := string(buf[:n])

	// Create session
	sessionID := fmt.Sprintf("whep_%d", time.Now().UnixNano())

	// Generate SDP answer (basic, supports H.264 + Opus)
	sdpAnswer := generateWHEPAnswer(sdpOffer, sessionID)

	// Subscribe to stream for this session
	subID := fmt.Sprintf("whep_%s", sessionID)
	sub := s.manager.Subscribe(key, subID, 128)
	if sub == nil {
		http.Error(w, "Subscribe failed", http.StatusInternalServerError)
		return
	}

	session := &whepSession{
		id:        sessionID,
		streamKey: key,
		subID:     subID,
		sdpOffer:  sdpOffer,
		sdpAnswer: sdpAnswer,
		createdAt: time.Now(),
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	// Start packet relay (would send via RTP over DTLS-SRTP in production)
	go s.relayPackets(session, sub)

	w.Header().Set("Content-Type", "application/sdp")
	w.Header().Set("Location", "/whep/session/"+sessionID)
	w.Header().Set("Link", fmt.Sprintf("</whep/ice/%s>; rel=\"ice\"", sessionID))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(sdpAnswer))

	log.Printf("[WHEP] Oturum oluşturuldu: %s -> %s", sessionID, key)
}

func (s *WHEPServer) handleICE(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimPrefix(r.URL.Path, "/whep/ice/")

	s.mu.RLock()
	sess, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	buf := make([]byte, 4096)
	n, _ := r.Body.Read(buf)
	candidate := string(buf[:n])

	s.mu.Lock()
	sess.candidates = append(sess.candidates, candidate)
	s.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func (s *WHEPServer) handleTeardown(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimPrefix(r.URL.Path, "/whep/session/")

	s.mu.Lock()
	sess, exists := s.sessions[sessionID]
	if exists {
		s.manager.Unsubscribe(sess.streamKey, sess.subID)
		delete(s.sessions, sessionID)
	}
	s.mu.Unlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("[WHEP] Oturum sonlandırıldı: %s", sessionID)
}

func (s *WHEPServer) relayPackets(sess *whepSession, sub *stream.OutputSubscriber) {
	defer s.manager.Unsubscribe(sess.streamKey, sess.subID)

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}
			// In a full implementation, this would:
			// 1. Packetize H.264 into RTP packets
			// 2. Send via DTLS-SRTP over UDP
			// For now, we process the packet to maintain the subscriber
			_ = pkt
		case <-sub.Done:
			return
		}
	}
}

func generateWHEPAnswer(offer string, sessionID string) string {
	// Generate a basic SDP answer for H.264 video + Opus audio
	answer := "v=0\r\n"
	answer += fmt.Sprintf("o=- %d 2 IN IP4 127.0.0.1\r\n", time.Now().UnixNano())
	answer += "s=FluxStream WHEP\r\n"
	answer += "t=0 0\r\n"
	answer += "a=group:BUNDLE 0 1\r\n"
	answer += "a=msid-semantic: WMS stream\r\n"

	// Video
	answer += "m=video 9 UDP/TLS/RTP/SAVPF 96\r\n"
	answer += "c=IN IP4 0.0.0.0\r\n"
	answer += "a=rtcp:9 IN IP4 0.0.0.0\r\n"
	answer += "a=sendonly\r\n"
	answer += "a=mid:0\r\n"
	answer += "a=rtpmap:96 H264/90000\r\n"
	answer += "a=fmtp:96 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f\r\n"
	answer += "a=rtcp-mux\r\n"
	answer += "a=setup:active\r\n"
	answer += fmt.Sprintf("a=ice-ufrag:%s\r\n", sessionID[:8])
	answer += fmt.Sprintf("a=ice-pwd:%s\r\n", sessionID)
	answer += "a=fingerprint:sha-256 00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00\r\n"

	// Audio
	answer += "m=audio 9 UDP/TLS/RTP/SAVPF 111\r\n"
	answer += "c=IN IP4 0.0.0.0\r\n"
	answer += "a=sendonly\r\n"
	answer += "a=mid:1\r\n"
	answer += "a=rtpmap:111 opus/48000/2\r\n"
	answer += "a=rtcp-mux\r\n"
	answer += "a=setup:active\r\n"

	return answer
}

// HandleStatus returns WHEP session stats
func (s *WHEPServer) HandleStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type sessionInfo struct {
		ID        string `json:"id"`
		StreamKey string `json:"stream_key"`
		CreatedAt string `json:"created_at"`
	}

	var sessions []sessionInfo
	for _, sess := range s.sessions {
		sessions = append(sessions, sessionInfo{
			ID:        sess.id,
			StreamKey: sess.streamKey,
			CreatedAt: sess.createdAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"active_sessions": len(sessions),
		"sessions":        sessions,
	})
}

// Cleanup removes expired sessions
func (s *WHEPServer) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, sess := range s.sessions {
		if now.Sub(sess.createdAt) > 2*time.Hour {
			s.manager.Unsubscribe(sess.streamKey, sess.subID)
			delete(s.sessions, id)
		}
	}
}
