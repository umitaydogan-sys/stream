package srt

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	ts "github.com/fluxstream/fluxstream/internal/media/container/ts"
	"github.com/fluxstream/fluxstream/internal/stream"
)

// Server serves SRT output for live streams
type Server struct {
	port    int
	manager *stream.Manager
	clients map[string]*srtClient
	mu      sync.RWMutex
}

type srtClient struct {
	id        string
	streamKey string
	subID     string
	conn      *net.UDPConn
	addr      *net.UDPAddr
	tsMuxer   *ts.Muxer
	stopCh    chan struct{}
	seqNum    uint32
}

// NewServer creates a new SRT output server
func NewServer(port int, manager *stream.Manager) *Server {
	return &Server{
		port:    port,
		manager: manager,
		clients: make(map[string]*srtClient),
	}
}

// Start begins the SRT output server
func (s *Server) Start(stop <-chan struct{}) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("SRT output resolve: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("SRT output listen: %w", err)
	}
	defer conn.Close()

	log.Printf("[SRT-OUT] Dinleniyor: :%d", s.port)

	go func() {
		<-stop
		conn.Close()
	}()

	buf := make([]byte, 1500)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-stop:
				return nil
			default:
				continue
			}
		}
		go s.handlePacket(conn, remoteAddr, buf[:n])
	}
}

func (s *Server) handlePacket(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	if len(data) < 16 {
		return
	}

	// Check if it's a SRT handshake
	isControl := (data[0] & 0x80) != 0
	if !isControl {
		return
	}

	msgType := binary.BigEndian.Uint16(data[0:2]) & 0x7FFF

	switch msgType {
	case 0x0000: // Handshake
		s.handleHandshake(conn, addr, data)
	case 0x0001: // Keep-alive
		// Echo back
		conn.WriteToUDP(data, addr)
	case 0x0005: // Shutdown
		clientID := addr.String()
		s.removeClient(clientID)
	}
}

func (s *Server) handleHandshake(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	if len(data) < 64 {
		return
	}

	clientID := addr.String()

	// Build handshake response
	resp := make([]byte, 64)
	copy(resp, data) // echo back with modifications
	resp[0] = 0x80   // control packet

	conn.WriteToUDP(resp, addr)

	// Extract stream key from extension (simplified)
	streamKey := ""
	if len(data) > 64 {
		// Stream ID extension
		extData := data[64:]
		streamKey = string(extData)
	}

	if streamKey == "" {
		streamKey = "default"
	}

	// Create client and start sending
	s.mu.Lock()
	client := &srtClient{
		id:        clientID,
		streamKey: streamKey,
		conn:      conn,
		addr:      addr,
		tsMuxer:   ts.NewMuxer(),
		stopCh:    make(chan struct{}),
	}
	s.clients[clientID] = client
	s.mu.Unlock()

	go s.sendLoop(client)
}

func (s *Server) sendLoop(client *srtClient) {
	// Subscribe to stream
	client.subID = fmt.Sprintf("srt_out_%s", client.id)
	sub := s.manager.Subscribe(client.streamKey, client.subID, 256)
	if sub == nil {
		s.removeClient(client.id)
		return
	}
	defer s.manager.Unsubscribe(client.streamKey, client.subID)

	log.Printf("[SRT-OUT] İstemci bağlandı: %s -> %s", client.addr, client.streamKey)

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}

			if pkt.IsSequenceHeader {
				continue
			}

			// Mux to TS
			mediaPkt := pkt.Clone()
			if pkt.Type == media.PacketTypeVideo && len(pkt.Data) > 5 {
				mediaPkt.Data = pkt.Data[5:]
			} else if pkt.Type == media.PacketTypeAudio && len(pkt.Data) > 2 {
				mediaPkt.Data = pkt.Data[2:]
			}

			tsData := client.tsMuxer.MuxPacket(mediaPkt)
			if tsData == nil {
				continue
			}

			// Send as SRT data packet (simplified)
			srtPkt := buildSRTDataPacket(client.seqNum, tsData)
			client.seqNum++
			client.conn.WriteToUDP(srtPkt, client.addr)

		case <-client.stopCh:
			return
		case <-sub.Done:
			return
		}
	}
}

func buildSRTDataPacket(seqNum uint32, payload []byte) []byte {
	// SRT data packet header (16 bytes) + payload
	header := make([]byte, 16)
	// Bit 0 = 0 (data packet)
	binary.BigEndian.PutUint32(header[0:4], seqNum)
	// Position flags (bits 0-1 of byte 4): 11 = solo packet
	header[4] = 0xC0
	// Message number
	binary.BigEndian.PutUint32(header[4:8], 0xC0000000|seqNum)
	// Timestamp
	binary.BigEndian.PutUint32(header[8:12], uint32(time.Now().UnixMilli()&0xFFFFFFFF))
	// Destination socket ID (0 for now)
	binary.BigEndian.PutUint32(header[12:16], 0)

	return append(header, payload...)
}

func (s *Server) removeClient(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, ok := s.clients[id]; ok {
		select {
		case <-client.stopCh:
		default:
			close(client.stopCh)
		}
		if client.subID != "" {
			s.manager.Unsubscribe(client.streamKey, client.subID)
		}
		delete(s.clients, id)
	}
}

// Stop stops the SRT output server
func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, client := range s.clients {
		select {
		case <-client.stopCh:
		default:
			close(client.stopCh)
		}
		delete(s.clients, id)
	}
}
