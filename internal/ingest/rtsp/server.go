package rtsp

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/ingest/rtmp"
	"github.com/fluxstream/fluxstream/internal/media"
)

// RTSP Methods
const (
	MethodOPTIONS  = "OPTIONS"
	MethodDESCRIBE = "DESCRIBE"
	MethodANNOUNCE = "ANNOUNCE"
	MethodSETUP    = "SETUP"
	MethodPLAY     = "PLAY"
	MethodRECORD   = "RECORD"
	MethodTEARDOWN = "TEARDOWN"
)

// Server is the RTSP listener
type Server struct {
	port    int
	handler rtmp.StreamHandler
}

// NewServer creates a new RTSP server
func NewServer(port int, handler rtmp.StreamHandler) *Server {
	return &Server{
		port:    port,
		handler: handler,
	}
}

// session represents an RTSP session
type session struct {
	id          string
	streamKey   string
	conn        net.Conn
	reader      *bufio.Reader
	cseq        int
	transport   string // "TCP" or "UDP"
	published   bool
	videoTrack  int
	audioTrack  int
	udpVideoPort int
	udpAudioPort int
	lastActive  time.Time

	// H.264 FU-A reassembly
	fuBuffer    []byte
	fuStarted   bool
}

// Start begins listening for RTSP connections
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("RTSP listen %s: %w", addr, err)
	}

	log.Printf("[RTSP] Dinleniyor: %s", addr)

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
				log.Printf("[RTSP] Accept hatası: %v", err)
				continue
			}
		}

		log.Printf("[RTSP] Yeni bağlantı: %s", conn.RemoteAddr())
		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	sess := &session{
		id:         fmt.Sprintf("rtsp_%d", time.Now().UnixNano()),
		conn:       conn,
		reader:     bufio.NewReaderSize(conn, 65536),
		videoTrack: -1,
		audioTrack: -1,
		lastActive: time.Now(),
	}

	defer func() {
		if sess.published && sess.streamKey != "" {
			s.handler.OnUnpublish(sess.streamKey)
			log.Printf("[RTSP] Yayın bitti: %s", sess.streamKey)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		// Check for interleaved RTP data (TCP mode)
		firstByte, err := sess.reader.ReadByte()
		if err != nil {
			return
		}

		if firstByte == '$' {
			// Interleaved RTP/RTCP data
			s.handleInterleavedData(sess)
			continue
		}

		// RTSP request
		sess.reader.UnreadByte()
		s.handleRTSPRequest(sess)
	}
}

func (s *Server) handleInterleavedData(sess *session) {
	header := make([]byte, 3)
	if _, err := sess.reader.Read(header); err != nil {
		return
	}

	channel := header[0]
	length := binary.BigEndian.Uint16(header[1:3])

	if length > 65535 || length == 0 {
		return
	}

	data := make([]byte, length)
	total := 0
	for total < int(length) {
		n, err := sess.reader.Read(data[total:])
		if err != nil {
			return
		}
		total += n
	}

	// Process RTP data
	if len(data) < 12 {
		return
	}

	isVideo := int(channel) == sess.videoTrack*2
	isAudio := int(channel) == sess.audioTrack*2

	if !isVideo && !isAudio {
		return // RTCP or unknown
	}

	timestamp := binary.BigEndian.Uint32(data[4:8])

	// payload starts after RTP header
	csrcCount := data[0] & 0x0F
	headerLen := 12 + int(csrcCount)*4

	// Extension header
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

	if isVideo {
		s.processVideoPayload(sess, payload, timestamp)
	} else if isAudio {
		pkt := &media.Packet{
			Type:       media.PacketTypeAudio,
			Timestamp:  timestamp / 48, // Approximate
			Data:       payload,
			StreamKey:  sess.streamKey,
			ReceivedAt: time.Now(),
		}
		s.handler.OnPacket(sess.streamKey, pkt)
	}
}

func (s *Server) processVideoPayload(sess *session, payload []byte, timestamp uint32) {
	if len(payload) == 0 {
		return
	}

	nalType := payload[0] & 0x1F

	switch {
	case nalType >= 1 && nalType <= 23:
		// Single NAL
		isKeyframe := nalType == 5 || nalType == 7
		pkt := &media.Packet{
			Type:       media.PacketTypeVideo,
			Timestamp:  timestamp / 90,
			Data:       payload,
			IsKeyframe: isKeyframe,
			StreamKey:  sess.streamKey,
			ReceivedAt: time.Now(),
		}
		if nalType == 7 {
			pkt.IsSequenceHeader = true
		}
		s.handler.OnPacket(sess.streamKey, pkt)

	case nalType == 24:
		// STAP-A
		offset := 1
		for offset+2 < len(payload) {
			nalSize := binary.BigEndian.Uint16(payload[offset : offset+2])
			offset += 2
			if offset+int(nalSize) > len(payload) {
				break
			}
			nalData := payload[offset : offset+int(nalSize)]
			nt := nalData[0] & 0x1F
			pkt := &media.Packet{
				Type:       media.PacketTypeVideo,
				Timestamp:  timestamp / 90,
				Data:       nalData,
				IsKeyframe: nt == 5 || nt == 7,
				StreamKey:  sess.streamKey,
				ReceivedAt: time.Now(),
			}
			s.handler.OnPacket(sess.streamKey, pkt)
			offset += int(nalSize)
		}

	case nalType == 28:
		// FU-A
		if len(payload) < 2 {
			return
		}
		fuHeader := payload[1]
		startBit := (fuHeader >> 7) & 0x01
		endBit := (fuHeader >> 6) & 0x01
		origNalType := fuHeader & 0x1F

		if startBit == 1 {
			nalHeader := (payload[0] & 0xE0) | origNalType
			sess.fuBuffer = []byte{nalHeader}
			sess.fuBuffer = append(sess.fuBuffer, payload[2:]...)
			sess.fuStarted = true
		} else if sess.fuStarted {
			sess.fuBuffer = append(sess.fuBuffer, payload[2:]...)
		}

		if endBit == 1 && sess.fuStarted {
			pkt := &media.Packet{
				Type:       media.PacketTypeVideo,
				Timestamp:  timestamp / 90,
				Data:       sess.fuBuffer,
				IsKeyframe: origNalType == 5,
				StreamKey:  sess.streamKey,
				ReceivedAt: time.Now(),
			}
			s.handler.OnPacket(sess.streamKey, pkt)
			sess.fuBuffer = nil
			sess.fuStarted = false
		}
	}
}

func (s *Server) handleRTSPRequest(sess *session) {
	// Read request line
	line, err := sess.reader.ReadString('\n')
	if err != nil {
		return
	}
	line = strings.TrimSpace(line)
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 3 {
		return
	}

	method := parts[0]
	uri := parts[1]

	// Read headers
	headers := make(map[string]string)
	for {
		hline, err := sess.reader.ReadString('\n')
		if err != nil {
			return
		}
		hline = strings.TrimSpace(hline)
		if hline == "" {
			break
		}
		idx := strings.Index(hline, ":")
		if idx > 0 {
			key := strings.TrimSpace(hline[:idx])
			val := strings.TrimSpace(hline[idx+1:])
			headers[key] = val
		}
	}

	if cs, ok := headers["CSeq"]; ok {
		sess.cseq, _ = strconv.Atoi(cs)
	}

	log.Printf("[RTSP] %s %s (CSeq: %d)", method, uri, sess.cseq)

	switch method {
	case MethodOPTIONS:
		s.sendResponse(sess, 200, "OK", map[string]string{
			"Public": "OPTIONS, DESCRIBE, ANNOUNCE, SETUP, PLAY, RECORD, TEARDOWN",
		}, "")

	case MethodANNOUNCE:
		// Read content body (SDP)
		contentLen := 0
		if cl, ok := headers["Content-Length"]; ok {
			contentLen, _ = strconv.Atoi(cl)
		}
		sdpBody := ""
		if contentLen > 0 {
			buf := make([]byte, contentLen)
			total := 0
			for total < contentLen {
				n, err := sess.reader.Read(buf[total:])
				if err != nil {
					return
				}
				total += n
			}
			sdpBody = string(buf)
		}

		// Extract stream key from URI: rtsp://host:port/live/{key}
		streamKey := extractStreamKey(uri)
		sess.streamKey = streamKey

		// Parse SDP for track info
		s.parseSDP(sess, sdpBody)

		s.sendResponse(sess, 200, "OK", nil, "")

	case MethodDESCRIBE:
		streamKey := extractStreamKey(uri)
		sdp := fmt.Sprintf("v=0\r\no=- 0 0 IN IP4 0.0.0.0\r\ns=%s\r\nc=IN IP4 0.0.0.0\r\nt=0 0\r\n"+
			"m=video 0 RTP/AVP 96\r\na=rtpmap:96 H264/90000\r\n"+
			"m=audio 0 RTP/AVP 97\r\na=rtpmap:97 MPEG4-GENERIC/48000/2\r\n",
			streamKey)
		s.sendResponse(sess, 200, "OK", map[string]string{
			"Content-Type": "application/sdp",
		}, sdp)

	case MethodSETUP:
		transport := headers["Transport"]
		sess.transport = "TCP"

		var respTransport string
		if strings.Contains(transport, "TCP") || strings.Contains(transport, "interleaved") {
			// TCP interleaved
			sess.transport = "TCP"
			// Figure out track
			if strings.Contains(uri, "track") {
				trackStr := uri[strings.LastIndex(uri, "track")+5:]
				trackNum, _ := strconv.Atoi(strings.TrimLeft(trackStr, "="))
				if sess.videoTrack < 0 {
					sess.videoTrack = trackNum
				} else {
					sess.audioTrack = trackNum
				}
			} else {
				if sess.videoTrack < 0 {
					sess.videoTrack = 0
				} else {
					sess.audioTrack = 1
				}
			}
			respTransport = fmt.Sprintf("RTP/AVP/TCP;unicast;interleaved=%d-%d", sess.videoTrack*2, sess.videoTrack*2+1)
		} else {
			// UDP
			sess.transport = "UDP"
			respTransport = transport + ";server_port=6970-6971"
		}

		s.sendResponse(sess, 200, "OK", map[string]string{
			"Transport": respTransport,
			"Session":   sess.id,
		}, "")

	case MethodRECORD:
		// Start receiving stream (push mode)
		if sess.streamKey == "" {
			sess.streamKey = extractStreamKey(uri)
		}
		if err := s.handler.OnPublish(sess.streamKey, sess.conn); err != nil {
			s.sendResponse(sess, 403, "Forbidden", nil, "")
			return
		}
		sess.published = true
		log.Printf("[RTSP] Yayın başladı (RECORD): %s", sess.streamKey)
		s.sendResponse(sess, 200, "OK", map[string]string{
			"Session": sess.id,
		}, "")

	case MethodPLAY:
		s.sendResponse(sess, 200, "OK", map[string]string{
			"Session": sess.id,
		}, "")

	case MethodTEARDOWN:
		s.sendResponse(sess, 200, "OK", nil, "")
		// Connection will close, defer will handle OnUnpublish
	}
}

func (s *Server) parseSDP(sess *session, sdp string) {
	trackNum := 0
	lines := strings.Split(sdp, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "m=video") {
			sess.videoTrack = trackNum
			trackNum++
		} else if strings.HasPrefix(line, "m=audio") {
			sess.audioTrack = trackNum
			trackNum++
		}
	}
}

func (s *Server) sendResponse(sess *session, statusCode int, statusText string, headers map[string]string, body string) {
	resp := fmt.Sprintf("RTSP/1.0 %d %s\r\nCSeq: %d\r\n", statusCode, statusText, sess.cseq)
	for k, v := range headers {
		resp += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	if body != "" {
		resp += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	}
	resp += "\r\n"
	if body != "" {
		resp += body
	}
	sess.conn.Write([]byte(resp))
}

func extractStreamKey(uri string) string {
	// rtsp://host:port/live/{key}
	parts := strings.Split(uri, "/")
	for i, p := range parts {
		if p == "live" && i+1 < len(parts) {
			key := parts[i+1]
			// Remove query parameters
			if idx := strings.Index(key, "?"); idx > 0 {
				key = key[:idx]
			}
			return key
		}
	}
	// Fallback: use last path segment
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		if idx := strings.Index(last, "?"); idx > 0 {
			last = last[:idx]
		}
		return last
	}
	return "unknown"
}

// SessionManager tracks active RTSP sessions (optional advanced features)
type SessionManager struct {
	sessions map[string]*session
	mu       sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*session),
	}
}
