package rtmp

import (
	"context"
	"fmt"
	"log"
	"net"
)

// Server is the RTMP listener
type Server struct {
	port    int
	handler StreamHandler
}

// NewServer creates a new RTMP server
func NewServer(port int, handler StreamHandler) *Server {
	return &Server{
		port:    port,
		handler: handler,
	}
}

// Start begins listening for RTMP connections
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("RTMP listen %s: %w", addr, err)
	}

	log.Printf("[RTMP] Dinleniyor: %s", addr)

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
				log.Printf("[RTMP] Accept hatası: %v", err)
				continue
			}
		}

		log.Printf("[RTMP] Yeni bağlantı: %s", conn.RemoteAddr())
		handler := NewHandler(conn, s.handler)
		go handler.Handle()
	}
}
