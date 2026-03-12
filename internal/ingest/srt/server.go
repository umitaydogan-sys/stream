package srt

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/ingest/rtmp"
	"github.com/fluxstream/fluxstream/internal/media"
)

// SRT Protocol Constants
const (
	SRTHeaderSize = 16

	// Control packet types
	SRTHandshake  = 0x0000
	SRTKeepAlive  = 0x0001
	SRTACK        = 0x0002
	SRTNAK        = 0x0003
	SRTShutdown   = 0x0005

	// Handshake types
	HSTypeInduction  = 1
	HSTypeConclusion = 0xFFFFFFFF

	// SRT version
	SRTVersion = 0x00010401 // 1.4.1
)

// Server is the SRT UDP listener
type Server struct {
	port    int
	handler rtmp.StreamHandler
	latency int // milliseconds
}

// NewServer creates a new SRT server
func NewServer(port int, handler rtmp.StreamHandler, latency int) *Server {
	if latency <= 0 {
		latency = 120
	}
	return &Server{
		port:    port,
		handler: handler,
		latency: latency,
	}
}

// session represents a single SRT connection
type session struct {
	socketID   uint32
	remoteAddr *net.UDPAddr
	streamKey  string
	lastActive time.Time
	seqNum     uint32
	mu         sync.Mutex
	established bool
	tsParser   *tsExtractor
}

// tsExtractor extracts media packets from MPEG-TS data
type tsExtractor struct {
	buffer    []byte
	streamKey string
	handler   rtmp.StreamHandler
	lastTS    uint32
}

func newTSExtractor(streamKey string, handler rtmp.StreamHandler) *tsExtractor {
	return &tsExtractor{
		streamKey: streamKey,
		handler:   handler,
	}
}

func (te *tsExtractor) Feed(data []byte) {
	te.buffer = append(te.buffer, data...)

	// Process complete TS packets (188 bytes each)
	for len(te.buffer) >= 188 {
		// Sync to TS packet boundary
		if te.buffer[0] != 0x47 {
			// Find next sync byte
			idx := -1
			for i := 1; i < len(te.buffer); i++ {
				if te.buffer[i] == 0x47 {
					idx = i
					break
				}
			}
			if idx < 0 {
				te.buffer = nil
				return
			}
			te.buffer = te.buffer[idx:]
			continue
		}

		if len(te.buffer) < 188 {
			break
		}

		tsPkt := te.buffer[:188]
		te.buffer = te.buffer[188:]

		// Parse TS header
		pid := (uint16(tsPkt[1]&0x1F) << 8) | uint16(tsPkt[2])

		// Skip PAT/PMT/null packets
		if pid == 0 || pid == 0x1FFF {
			continue
		}

		hasPayload := (tsPkt[3] & 0x10) != 0
		if !hasPayload {
			continue
		}

		offset := 4
		hasAdaptation := (tsPkt[3] & 0x20) != 0
		if hasAdaptation && offset < 188 {
			adaptLen := int(tsPkt[4])
			offset += 1 + adaptLen
		}

		if offset >= 188 {
			continue
		}

		payload := tsPkt[offset:]

		// Check for PES start
		if len(payload) >= 9 && payload[0] == 0x00 && payload[1] == 0x00 && payload[2] == 0x01 {
			streamID := payload[3]

			// Get PTS from PES header if available
			var pts uint32
			pesHeaderLen := int(payload[8])
			if (payload[7]&0x80) != 0 && pesHeaderLen >= 5 {
				pts = extractPTS(payload[9:])
			}

			pesPayload := payload[9+pesHeaderLen:]
			if len(pesPayload) == 0 {
				continue
			}

			var pktType media.PacketType
			isKeyframe := false

			if streamID >= 0xE0 && streamID <= 0xEF {
				pktType = media.PacketTypeVideo
				if len(pesPayload) > 4 {
					nalType := pesPayload[4] & 0x1F
					isKeyframe = nalType == 5 || nalType == 7
				}
			} else if streamID >= 0xC0 && streamID <= 0xDF {
				pktType = media.PacketTypeAudio
			} else {
				continue
			}

			if pts == 0 {
				te.lastTS += 33 // ~30fps
				pts = te.lastTS
			} else {
				te.lastTS = pts
			}

			pkt := &media.Packet{
				Type:       pktType,
				Timestamp:  pts,
				Data:       pesPayload,
				IsKeyframe: isKeyframe,
				StreamKey:  te.streamKey,
				ReceivedAt: time.Now(),
			}
			te.handler.OnPacket(te.streamKey, pkt)
		}
	}
}

func extractPTS(data []byte) uint32 {
	if len(data) < 5 {
		return 0
	}
	pts := (uint64(data[0]&0x0E) << 29) |
		(uint64(data[1]) << 22) |
		(uint64(data[2]&0xFE) << 14) |
		(uint64(data[3]) << 7) |
		(uint64(data[4]) >> 1)
	return uint32(pts / 90) // Convert to milliseconds
}

// Start begins listening for SRT connections
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("SRT resolve: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("SRT listen %s: %w", addr, err)
	}

	log.Printf("[SRT] Dinleniyor: %s (latency: %dms)", addr, s.latency)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	sessions := make(map[uint32]*session)
	var mu sync.Mutex
	buf := make([]byte, 65536)

	// Session cleanup goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				now := time.Now()
				for id, sess := range sessions {
					if now.Sub(sess.lastActive) > 30*time.Second {
						log.Printf("[SRT] Oturum zaman aşımı: %s", sess.remoteAddr)
						if sess.streamKey != "" {
							s.handler.OnUnpublish(sess.streamKey)
						}
						delete(sessions, id)
					}
				}
				mu.Unlock()
			}
		}
	}()

	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-ctx.Done():
				// Cleanup all sessions
				mu.Lock()
				for _, sess := range sessions {
					if sess.streamKey != "" {
						s.handler.OnUnpublish(sess.streamKey)
					}
				}
				mu.Unlock()
				return nil
			default:
				continue
			}
		}

		if n < SRTHeaderSize {
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		isControl := (data[0] & 0x80) != 0

		if isControl {
			s.handleControlPacket(conn, remoteAddr, data[:n], sessions, &mu)
		} else {
			s.handleDataPacket(data[:n], sessions, &mu)
		}
	}
}

func (s *Server) handleControlPacket(conn *net.UDPConn, addr *net.UDPAddr, data []byte, sessions map[uint32]*session, mu *sync.Mutex) {
	if len(data) < 16 {
		return
	}

	pktType := binary.BigEndian.Uint16(data[0:2]) & 0x7FFF

	switch pktType {
	case SRTHandshake:
		s.handleHandshake(conn, addr, data, sessions, mu)
	case SRTKeepAlive:
		socketID := binary.BigEndian.Uint32(data[12:16])
		mu.Lock()
		if sess, ok := sessions[socketID]; ok {
			sess.lastActive = time.Now()
		}
		mu.Unlock()
		// Send keepalive response
		resp := make([]byte, 16)
		resp[0] = 0x80
		binary.BigEndian.PutUint16(resp[0:2], 0x8001) // Keepalive
		binary.BigEndian.PutUint32(resp[4:8], binary.BigEndian.Uint32(data[4:8]))
		conn.WriteToUDP(resp, addr)

	case SRTACK:
		// ACK received
		socketID := binary.BigEndian.Uint32(data[12:16])
		mu.Lock()
		if sess, ok := sessions[socketID]; ok {
			sess.lastActive = time.Now()
		}
		mu.Unlock()

	case SRTShutdown:
		socketID := binary.BigEndian.Uint32(data[12:16])
		mu.Lock()
		if sess, ok := sessions[socketID]; ok {
			log.Printf("[SRT] Bağlantı kapatıldı: %s (key: %s)", addr, sess.streamKey)
			if sess.streamKey != "" {
				s.handler.OnUnpublish(sess.streamKey)
			}
			delete(sessions, socketID)
		}
		mu.Unlock()
	}
}

func (s *Server) handleHandshake(conn *net.UDPConn, addr *net.UDPAddr, data []byte, sessions map[uint32]*session, mu *sync.Mutex) {
	if len(data) < 64 {
		return
	}

	// Parse handshake
	hsType := binary.BigEndian.Uint32(data[16:20])
	peerSocketID := binary.BigEndian.Uint32(data[48:52])

	if hsType == HSTypeInduction {
		// Induction phase: respond with server socket ID
		serverSocketID := uint32(time.Now().UnixNano() & 0xFFFFFFFF)

		mu.Lock()
		sessions[serverSocketID] = &session{
			socketID:   serverSocketID,
			remoteAddr: addr,
			lastActive: time.Now(),
		}
		mu.Unlock()

		resp := s.buildHandshakeResponse(data, serverSocketID, peerSocketID, HSTypeInduction)
		conn.WriteToUDP(resp, addr)
		log.Printf("[SRT] Handshake induction: %s", addr)

	} else if hsType == HSTypeConclusion {
		// Conclusion phase: extract stream ID (stream key)
		streamKey := ""

		// Try to extract streamid from extension block
		if len(data) > 64 {
			extData := data[64:]
			streamKey = extractStreamID(extData)
		}

		if streamKey == "" {
			streamKey = fmt.Sprintf("srt_%d", peerSocketID)
		}

		// Find or create session
		mu.Lock()
		var sess *session
		for _, s := range sessions {
			if s.remoteAddr.String() == addr.String() && !s.established {
				sess = s
				break
			}
		}
		if sess == nil {
			serverSocketID := uint32(time.Now().UnixNano() & 0xFFFFFFFF)
			sess = &session{
				socketID:   serverSocketID,
				remoteAddr: addr,
				lastActive: time.Now(),
			}
			sessions[serverSocketID] = sess
		}
		sess.streamKey = streamKey
		sess.established = true
		sess.tsParser = newTSExtractor(streamKey, s.handler)
		mu.Unlock()

		// Create a dummy connection for OnPublish
		dummyConn := newSRTConn(addr, conn)
		if err := s.handler.OnPublish(streamKey, dummyConn); err != nil {
			log.Printf("[SRT] Publish reddedildi %s: %v", streamKey, err)
			return
		}

		resp := s.buildHandshakeResponse(data, sess.socketID, peerSocketID, HSTypeConclusion)
		conn.WriteToUDP(resp, addr)
		log.Printf("[SRT] Bağlantı kuruldu: %s (key: %s)", addr, streamKey)
	}
}

func extractStreamID(data []byte) string {
	// SRT extension blocks: type(2) + size(2) + data
	offset := 0
	for offset+4 <= len(data) {
		extType := binary.BigEndian.Uint16(data[offset : offset+2])
		extSize := binary.BigEndian.Uint16(data[offset+2 : offset+4])
		offset += 4

		blockSize := int(extSize) * 4
		if offset+blockSize > len(data) {
			break
		}

		// SRT_CMD_SID = 5 (Stream ID)
		if extType == 5 {
			sid := data[offset : offset+blockSize]
			// Trim null bytes
			end := len(sid)
			for end > 0 && sid[end-1] == 0 {
				end--
			}
			if end > 0 {
				return string(sid[:end])
			}
		}

		offset += blockSize
	}
	return ""
}

func (s *Server) buildHandshakeResponse(reqData []byte, serverSocketID, peerSocketID uint32, hsType uint32) []byte {
	resp := make([]byte, 64)
	copy(resp, reqData[:64])

	resp[0] = 0x80 // Control packet
	binary.BigEndian.PutUint16(resp[0:2], 0x8000) // Handshake
	binary.BigEndian.PutUint32(resp[16:20], hsType)     // HS type
	binary.BigEndian.PutUint32(resp[20:24], SRTVersion)  // Version
	binary.BigEndian.PutUint32(resp[48:52], serverSocketID) // syn cookie / socket ID

	return resp
}

func (s *Server) handleDataPacket(data []byte, sessions map[uint32]*session, mu *sync.Mutex) {
	if len(data) < 16 {
		return
	}

	socketID := binary.BigEndian.Uint32(data[12:16])

	mu.Lock()
	sess, ok := sessions[socketID]
	mu.Unlock()

	if !ok || !sess.established {
		return
	}

	sess.mu.Lock()
	sess.lastActive = time.Now()
	sess.seqNum++
	sess.mu.Unlock()

	// Payload starts after 16-byte header
	payload := data[16:]
	if len(payload) == 0 {
		return
	}

	// Feed to TS parser
	sess.tsParser.Feed(payload)
}

// srtConn wraps UDP address info to satisfy net.Conn for OnPublish
type srtConn struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func newSRTConn(addr *net.UDPAddr, conn *net.UDPConn) *srtConn {
	return &srtConn{addr: addr, conn: conn}
}

func (c *srtConn) Read(b []byte) (n int, err error)  { return 0, nil }
func (c *srtConn) Write(b []byte) (n int, err error)  { return c.conn.WriteToUDP(b, c.addr) }
func (c *srtConn) Close() error                       { return nil }
func (c *srtConn) LocalAddr() net.Addr                { return c.conn.LocalAddr() }
func (c *srtConn) RemoteAddr() net.Addr                { return c.addr }
func (c *srtConn) SetDeadline(t time.Time) error      { return nil }
func (c *srtConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *srtConn) SetWriteDeadline(t time.Time) error { return nil }
