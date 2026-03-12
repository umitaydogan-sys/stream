package mpegts

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

const (
	TSPacketSize = 188
	SyncByte     = 0x47

	// Stream type constants
	StreamTypeH264 = 0x1B
	StreamTypeH265 = 0x24
	StreamTypeAAC  = 0x0F
	StreamTypeMP3  = 0x03

	// PID constants
	PIDNull = 0x1FFF
	PIDPAT  = 0x0000
)

// Server is the MPEG-TS UDP listener (supports unicast and multicast)
type Server struct {
	port    int
	handler rtmp.StreamHandler
}

// NewServer creates a new MPEG-TS UDP server
func NewServer(port int, handler rtmp.StreamHandler) *Server {
	return &Server{
		port:    port,
		handler: handler,
	}
}

// tsSession tracks an MPEG-TS stream from a specific source
type tsSession struct {
	streamKey    string
	remoteAddr   *net.UDPAddr
	published    bool
	lastActive   time.Time
	pmtPID       uint16
	videoPID     uint16
	audioPID     uint16
	videoType    uint8
	audioType    uint8
	demuxer      *demuxer
}

// demuxer reassembles PES packets from TS packets
type demuxer struct {
	pesBuffers map[uint16]*pesBuffer
	streamKey  string
	handler    rtmp.StreamHandler
	mu         sync.Mutex
}

type pesBuffer struct {
	data        []byte
	pts         uint64
	dts         uint64
	isVideo     bool
	continuity  uint8
}

func newDemuxer(streamKey string, handler rtmp.StreamHandler) *demuxer {
	return &demuxer{
		pesBuffers: make(map[uint16]*pesBuffer),
		streamKey:  streamKey,
		handler:    handler,
	}
}

// Start begins listening for MPEG-TS UDP packets
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("MPEG-TS resolve: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("MPEG-TS listen %s: %w", addr, err)
	}

	log.Printf("[MPEG-TS] Dinleniyor: %s", addr)

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	sessions := make(map[string]*tsSession)
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
				for key, sess := range sessions {
					if now.Sub(sess.lastActive) > 30*time.Second {
						log.Printf("[MPEG-TS] Oturum zaman aşımı: %s", key)
						if sess.published {
							s.handler.OnUnpublish(sess.streamKey)
						}
						delete(sessions, key)
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

		if n < TSPacketSize {
			continue
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		addrKey := remoteAddr.String()

		mu.Lock()
		sess, exists := sessions[addrKey]
		if !exists {
			streamKey := fmt.Sprintf("ts_%s_%d", remoteAddr.IP, remoteAddr.Port)
			sess = &tsSession{
				streamKey:  streamKey,
				remoteAddr: remoteAddr,
				lastActive: time.Now(),
				demuxer:    newDemuxer(streamKey, s.handler),
			}
			sessions[addrKey] = sess
			mu.Unlock()

			// Publish
			dummyConn := newTSConn(remoteAddr, conn)
			if err := s.handler.OnPublish(streamKey, dummyConn); err != nil {
				log.Printf("[MPEG-TS] Publish reddedildi %s: %v", streamKey, err)
				mu.Lock()
				delete(sessions, addrKey)
				mu.Unlock()
				continue
			}
			sess.published = true
			log.Printf("[MPEG-TS] Yeni akış: %s", streamKey)
		} else {
			mu.Unlock()
		}

		sess.lastActive = time.Now()

		// Process all TS packets in the UDP datagram
		s.processTSPackets(sess, data)
	}
}

func (s *Server) processTSPackets(sess *tsSession, data []byte) {
	offset := 0

	// Find sync byte
	for offset < len(data) && data[offset] != SyncByte {
		offset++
	}

	for offset+TSPacketSize <= len(data) {
		if data[offset] != SyncByte {
			offset++
			continue
		}

		pkt := data[offset : offset+TSPacketSize]
		offset += TSPacketSize

		s.processOneTSPacket(sess, pkt)
	}
}

func (s *Server) processOneTSPacket(sess *tsSession, pkt []byte) {
	// Parse TS header
	pid := (uint16(pkt[1]&0x1F) << 8) | uint16(pkt[2])
	payloadStart := (pkt[1] & 0x40) != 0
	hasPayload := (pkt[3] & 0x10) != 0
	hasAdapt := (pkt[3] & 0x20) != 0

	if pid == PIDNull || !hasPayload {
		return
	}

	offset := 4
	if hasAdapt && offset < TSPacketSize {
		adaptLen := int(pkt[4])
		offset += 1 + adaptLen
	}

	if offset >= TSPacketSize {
		return
	}

	payload := pkt[offset:]

	// PAT (PID 0)
	if pid == PIDPAT {
		s.parsePAT(sess, payload, payloadStart)
		return
	}

	// PMT
	if pid == sess.pmtPID && sess.pmtPID != 0 {
		s.parsePMT(sess, payload, payloadStart)
		return
	}

	// Media PES data
	if pid == sess.videoPID || pid == sess.audioPID {
		sess.demuxer.feed(pid, payload, payloadStart, pid == sess.videoPID)
	}
}

func (s *Server) parsePAT(sess *tsSession, payload []byte, payloadStart bool) {
	offset := 0
	if payloadStart {
		if len(payload) == 0 {
			return
		}
		pointerField := int(payload[0])
		offset = 1 + pointerField
	}

	if offset+8 > len(payload) {
		return
	}

	// Skip table header (8 bytes minimum)
	// table_id(1) + section_syntax(2) + transport_stream_id(2) + version(1) + section(1) + last_section(1)
	offset += 8

	// Each program entry is 4 bytes
	for offset+4 <= len(payload) {
		programNum := binary.BigEndian.Uint16(payload[offset : offset+2])
		pmtPID := (uint16(payload[offset+2]&0x1F) << 8) | uint16(payload[offset+3])
		offset += 4

		if programNum != 0 {
			sess.pmtPID = pmtPID
			break
		}
	}
}

func (s *Server) parsePMT(sess *tsSession, payload []byte, payloadStart bool) {
	offset := 0
	if payloadStart {
		if len(payload) == 0 {
			return
		}
		pointerField := int(payload[0])
		offset = 1 + pointerField
	}

	if offset+12 > len(payload) {
		return
	}

	// PMT header
	// table_id(1) + section_length(2) + program_number(2) + version(1) + section_number(1) + last_section(1) + PCR_PID(2) + program_info_length(2)
	sectionLen := (int(payload[offset+1]&0x0F) << 8) | int(payload[offset+2])
	progInfoLen := (int(payload[offset+10]&0x0F) << 8) | int(payload[offset+11])
	offset += 12 + progInfoLen

	endOffset := offset + sectionLen - 12 - progInfoLen - 4 // minus CRC
	if endOffset > len(payload) {
		endOffset = len(payload)
	}

	for offset+5 <= endOffset {
		streamType := payload[offset]
		elementaryPID := (uint16(payload[offset+1]&0x1F) << 8) | uint16(payload[offset+2])
		esInfoLen := (int(payload[offset+3]&0x0F) << 8) | int(payload[offset+4])
		offset += 5 + esInfoLen

		switch streamType {
		case StreamTypeH264, StreamTypeH265:
			sess.videoPID = elementaryPID
			sess.videoType = streamType
		case StreamTypeAAC, StreamTypeMP3:
			sess.audioPID = elementaryPID
			sess.audioType = streamType
		}
	}
}

func (d *demuxer) feed(pid uint16, payload []byte, payloadStart bool, isVideo bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if payloadStart {
		// Flush any existing buffer for this PID
		if buf, ok := d.pesBuffers[pid]; ok && len(buf.data) > 0 {
			d.emitPacket(buf)
		}

		// Parse PES header
		if len(payload) < 9 || payload[0] != 0x00 || payload[1] != 0x00 || payload[2] != 0x01 {
			d.pesBuffers[pid] = &pesBuffer{data: payload, isVideo: isVideo}
			return
		}

		var pts, dts uint64
		ptsFlag := (payload[7] & 0x80) != 0
		dtsFlag := (payload[7] & 0x40) != 0
		pesHeaderLen := int(payload[8])

		if ptsFlag && pesHeaderLen >= 5 {
			pts = extractTimestamp(payload[9:])
		}
		if dtsFlag && pesHeaderLen >= 10 {
			dts = extractTimestamp(payload[14:])
		}

		pesPayloadStart := 9 + pesHeaderLen
		if pesPayloadStart > len(payload) {
			pesPayloadStart = len(payload)
		}

		d.pesBuffers[pid] = &pesBuffer{
			data:    payload[pesPayloadStart:],
			pts:     pts,
			dts:     dts,
			isVideo: isVideo,
		}
	} else {
		if buf, ok := d.pesBuffers[pid]; ok {
			buf.data = append(buf.data, payload...)
		}
	}
}

func (d *demuxer) emitPacket(buf *pesBuffer) {
	if len(buf.data) == 0 {
		return
	}

	var pktType media.PacketType
	isKeyframe := false

	if buf.isVideo {
		pktType = media.PacketTypeVideo
		// Check for H.264 IDR
		if len(buf.data) > 4 {
			nalType := buf.data[0] & 0x1F
			if nalType == 0 && len(buf.data) > 5 {
				// Check after start code
				for i := 0; i < len(buf.data)-4; i++ {
					if buf.data[i] == 0 && buf.data[i+1] == 0 && buf.data[i+2] == 1 {
						nalType = buf.data[i+3] & 0x1F
						break
					}
					if buf.data[i] == 0 && buf.data[i+1] == 0 && buf.data[i+2] == 0 && i+3 < len(buf.data) && buf.data[i+3] == 1 {
						if i+4 < len(buf.data) {
							nalType = buf.data[i+4] & 0x1F
						}
						break
					}
				}
			}
			isKeyframe = nalType == 5 || nalType == 7
		}
	} else {
		pktType = media.PacketTypeAudio
	}

	timestamp := uint32(buf.pts / 90) // Convert from 90kHz
	if timestamp == 0 {
		timestamp = uint32(buf.dts / 90)
	}

	pkt := &media.Packet{
		Type:       pktType,
		Timestamp:  timestamp,
		Data:       buf.data,
		IsKeyframe: isKeyframe,
		StreamKey:  d.streamKey,
		ReceivedAt: time.Now(),
	}
	d.handler.OnPacket(d.streamKey, pkt)
}

func extractTimestamp(data []byte) uint64 {
	if len(data) < 5 {
		return 0
	}
	return (uint64(data[0]&0x0E) << 29) |
		(uint64(data[1]) << 22) |
		(uint64(data[2]&0xFE) << 14) |
		(uint64(data[3]) << 7) |
		(uint64(data[4]) >> 1)
}

// tsConn wraps UDP connection for the handler interface
type tsConn struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func newTSConn(addr *net.UDPAddr, conn *net.UDPConn) *tsConn {
	return &tsConn{addr: addr, conn: conn}
}

func (c *tsConn) Read(b []byte) (n int, err error)  { return 0, nil }
func (c *tsConn) Write(b []byte) (n int, err error)  { return c.conn.WriteToUDP(b, c.addr) }
func (c *tsConn) Close() error                       { return nil }
func (c *tsConn) LocalAddr() net.Addr                { return c.conn.LocalAddr() }
func (c *tsConn) RemoteAddr() net.Addr                { return c.addr }
func (c *tsConn) SetDeadline(t time.Time) error      { return nil }
func (c *tsConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *tsConn) SetWriteDeadline(t time.Time) error { return nil }
