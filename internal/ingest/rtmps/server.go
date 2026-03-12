package rtmps

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/fluxstream/fluxstream/internal/ingest/rtmp"
)

// Server is the RTMPS (TLS) listener
type Server struct {
	port      int
	handler   rtmp.StreamHandler
	certFile  string
	keyFile   string
}

// NewServer creates a new RTMPS server
func NewServer(port int, handler rtmp.StreamHandler, certFile, keyFile string) *Server {
	return &Server{
		port:     port,
		handler:  handler,
		certFile: certFile,
		keyFile:  keyFile,
	}
}

// Start begins listening for RTMPS connections
func (s *Server) Start(ctx context.Context) error {
	if s.certFile == "" || s.keyFile == "" {
		return fmt.Errorf("RTMPS: sertifika dosyaları belirtilmedi")
	}

	cert, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		return fmt.Errorf("RTMPS TLS sertifika yükleme: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	addr := fmt.Sprintf(":%d", s.port)
	listener, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("RTMPS listen %s: %w", addr, err)
	}

	log.Printf("[RTMPS] Dinleniyor: %s", addr)

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				log.Printf("[RTMPS] Accept hatası: %v", err)
				continue
			}
		}

		log.Printf("[RTMPS] Yeni TLS bağlantı: %s", conn.RemoteAddr())
		handler := rtmp.NewHandler(conn, s.handler)
		go handler.Handle()
	}
}

// SetCerts updates the certificate paths (for runtime cert changes)
func (s *Server) SetCerts(certFile, keyFile string) {
	s.certFile = certFile
	s.keyFile = keyFile
}

// HasValidCerts checks if certificate files are configured
func (s *Server) HasValidCerts() bool {
	return s.certFile != "" && s.keyFile != ""
}

// Addr returns the listen address
func (s *Server) Addr() string {
	return fmt.Sprintf(":%d", s.port)
}

// dummyConn wraps net.Conn for the handler
type dummyConn struct {
	net.Conn
}
