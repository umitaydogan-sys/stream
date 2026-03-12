package httppush

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/ingest/rtmp"
	"github.com/fluxstream/fluxstream/internal/media"
)

// Server accepts HTTP PUT/POST with MPEG-TS data for ingest
type Server struct {
	port      int
	handler   rtmp.StreamHandler
	authToken string
}

// NewServer creates a new HTTP Push ingest server
func NewServer(port int, handler rtmp.StreamHandler, authToken string) *Server {
	return &Server{
		port:      port,
		handler:   handler,
		authToken: authToken,
	}
}

type pushSession struct {
	streamKey  string
	published  bool
	mu         sync.Mutex
}

// Start begins listening for HTTP Push requests
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	sessions := &sync.Map{}

	mux.HandleFunc("/push/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Auth check
		if s.authToken != "" {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") || auth[7:] != s.authToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Extract stream key from path
		path := strings.TrimPrefix(r.URL.Path, "/push/")
		streamKey := strings.TrimRight(path, "/")
		if streamKey == "" {
			http.Error(w, "Stream key required", http.StatusBadRequest)
			return
		}

		// Sanitize stream key
		if strings.Contains(streamKey, "..") || strings.Contains(streamKey, "/") {
			http.Error(w, "Invalid stream key", http.StatusBadRequest)
			return
		}

		// Get or create session
		var sess *pushSession
		existing, loaded := sessions.LoadOrStore(streamKey, &pushSession{
			streamKey: streamKey,
		})
		sess = existing.(*pushSession)

		if !loaded {
			conn := newHTTPConn(r.RemoteAddr)
			if err := s.handler.OnPublish(streamKey, conn); err != nil {
				sessions.Delete(streamKey)
				http.Error(w, "Publish rejected", http.StatusForbidden)
				return
			}
			sess.published = true
			log.Printf("[HTTP-PUSH] Akış başladı: %s (%s)", streamKey, r.RemoteAddr)
		}

		contentType := r.Header.Get("Content-Type")

		switch {
		case strings.Contains(contentType, "video/mp2t"), contentType == "":
			s.handleTSPush(sess, r.Body)
		default:
			// Accept anything, try to treat as TS
			s.handleTSPush(sess, r.Body)
		}

		// For non-chunked, unpublish when request ends
		isChunked := r.Header.Get("Transfer-Encoding") == "chunked"
		if !isChunked {
			sess.mu.Lock()
			if sess.published {
				s.handler.OnUnpublish(streamKey)
				sess.published = false
				sessions.Delete(streamKey)
				log.Printf("[HTTP-PUSH] Akış bitti: %s", streamKey)
			}
			sess.mu.Unlock()
		}

		w.WriteHeader(http.StatusOK)
	})

	// Status endpoint
	mux.HandleFunc("/push/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		activeCount := 0
		sessions.Range(func(_, _ interface{}) bool {
			activeCount++
			return true
		})
		fmt.Fprintf(w, `{"active_streams":%d}`, activeCount)
	})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("HTTP-Push listen: %w", err)
	}

	server := &http.Server{
		Handler:      mux,
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("[HTTP-PUSH] Dinleniyor: :%d", s.port)

	go func() {
		<-ctx.Done()
		sessions.Range(func(key, value interface{}) bool {
			sess := value.(*pushSession)
			sess.mu.Lock()
			if sess.published {
				s.handler.OnUnpublish(sess.streamKey)
			}
			sess.mu.Unlock()
			return true
		})
		server.Close()
	}()

	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP-Push serve: %w", err)
	}
	return nil
}

func (s *Server) handleTSPush(sess *pushSession, body io.ReadCloser) {
	defer body.Close()

	buf := make([]byte, 188*7) // Read multiple TS packets at once
	var residual []byte

	for {
		n, err := body.Read(buf)
		if n > 0 {
			data := append(residual, buf[:n]...)
			residual = nil

			// Process complete TS packets
			offset := 0
			// Find first sync byte
			for offset < len(data) && data[offset] != 0x47 {
				offset++
			}

			for offset+188 <= len(data) {
				if data[offset] != 0x47 {
					offset++
					continue
				}

				tsPkt := data[offset : offset+188]
				offset += 188

				s.processTSPacket(sess, tsPkt)
			}

			// Keep remaining data
			if offset < len(data) {
				residual = make([]byte, len(data)-offset)
				copy(residual, data[offset:])
			}
		}

		if err != nil {
			break
		}
	}
}

func (s *Server) processTSPacket(sess *pushSession, pkt []byte) {
	// Basic TS packet parsing - extract PID and payload
	pid := (uint16(pkt[1]&0x1F) << 8) | uint16(pkt[2])
	payloadStart := (pkt[1] & 0x40) != 0
	hasPayload := (pkt[3] & 0x10) != 0
	hasAdapt := (pkt[3] & 0x20) != 0

	if pid == 0x1FFF || !hasPayload {
		return
	}

	offset := 4
	if hasAdapt {
		adaptLen := int(pkt[4])
		offset += 1 + adaptLen
	}

	if offset >= 188 {
		return
	}

	payload := pkt[offset:]

	// For HTTP Push, we send raw TS data as video packets
	// The downstream (HLS muxer) handles actual demuxing
	if payloadStart && len(payload) >= 4 {
		// Check for PES start code
		if payload[0] == 0x00 && payload[1] == 0x00 && payload[2] == 0x01 {
			streamID := payload[3]
			isVideo := (streamID >= 0xE0 && streamID <= 0xEF)
			isAudio := (streamID >= 0xC0 && streamID <= 0xDF)

			var pktType media.PacketType
			if isVideo {
				pktType = media.PacketTypeVideo
			} else if isAudio {
				pktType = media.PacketTypeAudio
			} else {
				return
			}

			mediaPkt := &media.Packet{
				Type:       pktType,
				Timestamp:  uint32(time.Now().UnixMilli()),
				Data:       payload,
				IsKeyframe: false,
				StreamKey:  sess.streamKey,
				ReceivedAt: time.Now(),
			}
			s.handler.OnPacket(sess.streamKey, mediaPkt)
		}
	}
}

// httpConn implements net.Conn for the handler interface
type httpConn struct {
	remoteAddr string
}

func newHTTPConn(remoteAddr string) *httpConn {
	return &httpConn{remoteAddr: remoteAddr}
}

func (c *httpConn) Read(b []byte) (n int, err error)  { return 0, nil }
func (c *httpConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (c *httpConn) Close() error                       { return nil }
func (c *httpConn) LocalAddr() net.Addr                { return &net.TCPAddr{} }
func (c *httpConn) RemoteAddr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp", c.remoteAddr)
	if addr == nil {
		return &net.TCPAddr{}
	}
	return addr
}
func (c *httpConn) SetDeadline(t time.Time) error      { return nil }
func (c *httpConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *httpConn) SetWriteDeadline(t time.Time) error { return nil }
