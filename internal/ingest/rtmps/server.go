package rtmps

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/fluxstream/fluxstream/internal/ingest/rtmp"
)

// Server is the RTMPS (TLS) listener.
type Server struct {
	port      int
	handler   rtmp.StreamHandler
	tlsConfig *tls.Config
}

// NewServer creates a new RTMPS server.
func NewServer(port int, handler rtmp.StreamHandler, tlsConfig *tls.Config) *Server {
	return &Server{
		port:      port,
		handler:   handler,
		tlsConfig: tlsConfig,
	}
}

// Start begins listening for RTMPS connections.
func (s *Server) Start(ctx context.Context) error {
	if s.tlsConfig == nil {
		return fmt.Errorf("RTMPS TLS config belirtilmedi")
	}

	addr := fmt.Sprintf(":%d", s.port)
	listener, err := tls.Listen("tcp", addr, s.tlsConfig.Clone())
	if err != nil {
		return fmt.Errorf("RTMPS listen %s: %w", addr, err)
	}

	log.Printf("[RTMPS] Dinleniyor: %s", addr)

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Printf("[RTMPS] Accept hatasi: %v", err)
				continue
			}
		}

		log.Printf("[RTMPS] Yeni TLS baglanti: %s", conn.RemoteAddr())
		handler := rtmp.NewHandler(conn, s.handler)
		go handler.Handle()
	}
}

// Addr returns the listen address.
func (s *Server) Addr() string {
	return fmt.Sprintf(":%d", s.port)
}

// dummyConn wraps net.Conn for the handler.
type dummyConn struct {
	net.Conn
}
