package rtp

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

// RTP header constants
const (
	RTPHeaderSize    = 12
	RTPVersion       = 2
	PayloadTypeH264  = 96
	PayloadTypeAAC   = 97
	PayloadTypeOpus  = 111

	// H.264 NAL unit types
	NALSingle  = 0  // 1-23: single NAL unit
	NALSTAPA   = 24 // STAP-A
	NALFUA     = 28 // FU-A
)

// Server is the RTP UDP listener
type Server struct {
	port    int
	handler rtmp.StreamHandler
}

// NewServer creates a new RTP server
func NewServer(port int, handler rtmp.StreamHandler) *Server {
	return &Server{
		port:    port,
		handler: handler,
	}
}

// rtpSession tracks an RTP stream session
type rtpSession struct {
	streamKey   string
	remoteAddr  *net.UDPAddr
	ssrc        uint32
	lastActive  time.Time
	lastSeq     uint16
	published   bool

	// H.264 FU-A reassembly
	fuBuffer    []byte
	fuStarted   bool
	fuTimestamp uint32

	// Jitter buffer
	jitterBuf   map[uint16]*rtpPacket
	jitterMu    sync.Mutex
	nextSeq     uint16
}

type rtpPacket struct {
	seq       uint16
	timestamp uint32
	payload   []byte
	marker    bool
	pt        uint8
	receivedAt time.Time
}

// Start begins listening for RTP packets
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("RTP resolve: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("RTP listen %s: %w", addr, err)
	}

	log.Printf("[RTP] Dinleniyor: %s", addr)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	sessions := make(map[uint32]*rtpSession)
	var mu sync.Mutex
	buf := make([]byte, 65536)

	// Session cleanup
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
				for ssrc, sess := range sessions {
					if now.Sub(sess.lastActive) > 30*time.Second {
						log.Printf("[RTP] Oturum zaman aşımı: %s", sess.remoteAddr)
						if sess.published {
							s.handler.OnUnpublish(sess.streamKey)
						}
						delete(sessions, ssrc)
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
				mu.Lock()
				for _, sess := range sessions {
					if sess.published {
						s.handler.OnUnpublish(sess.streamKey)
					}
				}
				mu.Unlock()
				return nil
			default:
				continue
			}
		}

		if n < RTPHeaderSize {
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		s.processRTPPacket(conn, remoteAddr, data, sessions, &mu)
	}
}

func (s *Server) processRTPPacket(conn *net.UDPConn, addr *net.UDPAddr, data []byte, sessions map[uint32]*rtpSession, mu *sync.Mutex) {
	// Parse RTP header
	version := (data[0] >> 6) & 0x03
	if version != RTPVersion {
		return
	}

	// padding := (data[0] >> 5) & 0x01
	// extension := (data[0] >> 4) & 0x01
	csrcCount := data[0] & 0x0F
	marker := (data[1] >> 7) & 0x01
	pt := data[1] & 0x7F
	seq := binary.BigEndian.Uint16(data[2:4])
	timestamp := binary.BigEndian.Uint32(data[4:8])
	ssrc := binary.BigEndian.Uint32(data[8:12])

	headerLen := RTPHeaderSize + int(csrcCount)*4
	if headerLen >= len(data) {
		return
	}

	// Check for header extension
	if (data[0] & 0x10) != 0 {
		if headerLen+4 > len(data) {
			return
		}
		extLen := binary.BigEndian.Uint16(data[headerLen+2 : headerLen+4])
		headerLen += 4 + int(extLen)*4
	}

	if headerLen >= len(data) {
		return
	}

	payload := data[headerLen:]

	mu.Lock()
	sess, exists := sessions[ssrc]
	if !exists {
		streamKey := fmt.Sprintf("rtp_%d", ssrc)
		sess = &rtpSession{
			streamKey:  streamKey,
			remoteAddr: addr,
			ssrc:       ssrc,
			lastActive: time.Now(),
			jitterBuf:  make(map[uint16]*rtpPacket),
		}
		sessions[ssrc] = sess
		mu.Unlock()

		// Publish
		dummyConn := newRTPConn(addr, conn)
		if err := s.handler.OnPublish(streamKey, dummyConn); err != nil {
			log.Printf("[RTP] Publish reddedildi %s: %v", streamKey, err)
			mu.Lock()
			delete(sessions, ssrc)
			mu.Unlock()
			return
		}
		sess.published = true
		log.Printf("[RTP] Yeni akış: %s (SSRC: %d)", addr, ssrc)
	} else {
		mu.Unlock()
	}

	sess.lastActive = time.Now()
	sess.lastSeq = seq

	pkt := &rtpPacket{
		seq:        seq,
		timestamp:  timestamp,
		payload:    payload,
		marker:     marker != 0,
		pt:         pt,
		receivedAt: time.Now(),
	}

	s.processPayload(sess, pkt)
}

func (s *Server) processPayload(sess *rtpSession, pkt *rtpPacket) {
	if len(pkt.payload) == 0 {
		return
	}

	switch {
	case pkt.pt == PayloadTypeH264 || pkt.pt == PayloadTypeH264+1:
		s.processH264(sess, pkt)
	case pkt.pt == PayloadTypeAAC:
		s.processAudio(sess, pkt, media.PacketTypeAudio)
	default:
		s.processGenericPayload(sess, pkt)
	}
}

func (s *Server) processH264(sess *rtpSession, pkt *rtpPacket) {
	nalType := pkt.payload[0] & 0x1F

	switch {
	case nalType >= 1 && nalType <= 23:
		// Single NAL unit
		isKeyframe := nalType == 5 || nalType == 7
		mediaPkt := &media.Packet{
			Type:       media.PacketTypeVideo,
			Timestamp:  pkt.timestamp / 90, // Convert from 90kHz clock
			Data:       pkt.payload,
			IsKeyframe: isKeyframe,
			StreamKey:  sess.streamKey,
			ReceivedAt: time.Now(),
		}
		if nalType == 7 {
			mediaPkt.IsSequenceHeader = true
		}
		s.handler.OnPacket(sess.streamKey, mediaPkt)

	case nalType == NALSTAPA:
		// STAP-A: multiple NAL units in one packet
		offset := 1
		for offset+2 < len(pkt.payload) {
			nalSize := binary.BigEndian.Uint16(pkt.payload[offset : offset+2])
			offset += 2
			if offset+int(nalSize) > len(pkt.payload) {
				break
			}
			nalData := pkt.payload[offset : offset+int(nalSize)]
			nt := nalData[0] & 0x1F
			isKeyframe := nt == 5 || nt == 7
			mediaPkt := &media.Packet{
				Type:       media.PacketTypeVideo,
				Timestamp:  pkt.timestamp / 90,
				Data:       nalData,
				IsKeyframe: isKeyframe,
				StreamKey:  sess.streamKey,
				ReceivedAt: time.Now(),
			}
			s.handler.OnPacket(sess.streamKey, mediaPkt)
			offset += int(nalSize)
		}

	case nalType == NALFUA:
		// FU-A: fragmented NAL unit
		if len(pkt.payload) < 2 {
			return
		}
		fuHeader := pkt.payload[1]
		startBit := (fuHeader >> 7) & 0x01
		endBit := (fuHeader >> 6) & 0x01
		origNalType := fuHeader & 0x1F

		if startBit == 1 {
			// Start of fragmented NAL
			nalHeader := (pkt.payload[0] & 0xE0) | origNalType
			sess.fuBuffer = []byte{nalHeader}
			sess.fuBuffer = append(sess.fuBuffer, pkt.payload[2:]...)
			sess.fuStarted = true
			sess.fuTimestamp = pkt.timestamp
		} else if sess.fuStarted {
			sess.fuBuffer = append(sess.fuBuffer, pkt.payload[2:]...)
		}

		if endBit == 1 && sess.fuStarted {
			isKeyframe := origNalType == 5 || origNalType == 7
			mediaPkt := &media.Packet{
				Type:       media.PacketTypeVideo,
				Timestamp:  sess.fuTimestamp / 90,
				Data:       sess.fuBuffer,
				IsKeyframe: isKeyframe,
				StreamKey:  sess.streamKey,
				ReceivedAt: time.Now(),
			}
			s.handler.OnPacket(sess.streamKey, mediaPkt)
			sess.fuBuffer = nil
			sess.fuStarted = false
		}
	}
}

func (s *Server) processAudio(sess *rtpSession, pkt *rtpPacket, pktType media.PacketType) {
	mediaPkt := &media.Packet{
		Type:       pktType,
		Timestamp:  pkt.timestamp / (48000 / 1000), // Convert from audio clock
		Data:       pkt.payload,
		StreamKey:  sess.streamKey,
		ReceivedAt: time.Now(),
	}
	s.handler.OnPacket(sess.streamKey, mediaPkt)
}

func (s *Server) processGenericPayload(sess *rtpSession, pkt *rtpPacket) {
	mediaPkt := &media.Packet{
		Type:       media.PacketTypeVideo,
		Timestamp:  pkt.timestamp / 90,
		Data:       pkt.payload,
		StreamKey:  sess.streamKey,
		ReceivedAt: time.Now(),
	}
	s.handler.OnPacket(sess.streamKey, mediaPkt)
}

// rtpConn wraps UDP connection for the handler interface
type rtpConn struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func newRTPConn(addr *net.UDPAddr, conn *net.UDPConn) *rtpConn {
	return &rtpConn{addr: addr, conn: conn}
}

func (c *rtpConn) Read(b []byte) (n int, err error)  { return 0, nil }
func (c *rtpConn) Write(b []byte) (n int, err error)  { return c.conn.WriteToUDP(b, c.addr) }
func (c *rtpConn) Close() error                       { return nil }
func (c *rtpConn) LocalAddr() net.Addr                { return c.conn.LocalAddr() }
func (c *rtpConn) RemoteAddr() net.Addr                { return c.addr }
func (c *rtpConn) SetDeadline(t time.Time) error      { return nil }
func (c *rtpConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *rtpConn) SetWriteDeadline(t time.Time) error { return nil }
