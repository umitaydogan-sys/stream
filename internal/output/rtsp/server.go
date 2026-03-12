package rtsp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/stream"
)

// Server serves RTSP output for live streams
type Server struct {
	port      int
	manager   *stream.Manager
	listener  net.Listener
	sessions  map[string]*rtspSession
	mu        sync.RWMutex
}

type rtspSession struct {
	id         string
	streamKey  string
	subID      string
	conn       net.Conn
	rtpConn    *net.UDPConn
	rtpPort    int
	rtcpPort   int
	clientAddr *net.UDPAddr
	transport  string
	cseq       int
	ssrc       uint32
	seqNum     uint16
}

// NewServer creates a new RTSP output server
func NewServer(port int, manager *stream.Manager) *Server {
	return &Server{
		port:     port,
		manager:  manager,
		sessions: make(map[string]*rtspSession),
	}
}

// Start begins the RTSP output server
func (s *Server) Start(stop <-chan struct{}) error {
	addr := fmt.Sprintf(":%d", s.port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("RTSP output listen: %w", err)
	}
	s.listener = ln
	log.Printf("[RTSP-OUT] Dinleniyor: %s", addr)

	go func() {
		<-stop
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-stop:
				return nil
			default:
				log.Printf("[RTSP-OUT] Accept hatası: %v", err)
				continue
			}
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	sessionID := fmt.Sprintf("rtsp_out_%d", time.Now().UnixNano())

	sess := &rtspSession{
		id:   sessionID,
		conn: conn,
		ssrc: uint32(time.Now().UnixNano() & 0xFFFFFFFF),
	}

	s.mu.Lock()
	s.sessions[sessionID] = sess
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		if sess.subID != "" {
			s.manager.Unsubscribe(sess.streamKey, sess.subID)
		}
		delete(s.sessions, sessionID)
		s.mu.Unlock()
	}()

	log.Printf("[RTSP-OUT] Bağlantı: %s", conn.RemoteAddr())

	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		// Parse RTSP request
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		method := parts[0]
		uri := parts[1]

		// Read headers
		headers := make(map[string]string)
		for {
			hline, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			hline = strings.TrimSpace(hline)
			if hline == "" {
				break
			}
			idx := strings.Index(hline, ":")
			if idx > 0 {
				headers[strings.TrimSpace(hline[:idx])] = strings.TrimSpace(hline[idx+1:])
			}
		}

		cseq := headers["CSeq"]

		switch method {
		case "OPTIONS":
			s.sendResponse(conn, "200 OK", cseq, sessionID,
				"Public: OPTIONS, DESCRIBE, SETUP, PLAY, TEARDOWN\r\n")

		case "DESCRIBE":
			streamKey := extractStreamKey(uri)
			sess.streamKey = streamKey

			if !s.manager.IsLive(streamKey) {
				s.sendResponse(conn, "404 Not Found", cseq, sessionID, "")
				continue
			}

			sdp := generateSDP(streamKey, conn.LocalAddr().String())
			s.sendResponse(conn, "200 OK", cseq, sessionID,
				fmt.Sprintf("Content-Type: application/sdp\r\nContent-Length: %d\r\n\r\n%s", len(sdp), sdp))

		case "SETUP":
			transport := headers["Transport"]
			sess.transport = transport

			// Parse client RTP port
			if idx := strings.Index(transport, "client_port="); idx >= 0 {
				portStr := transport[idx+12:]
				if dashIdx := strings.Index(portStr, "-"); dashIdx > 0 {
					portStr = portStr[:dashIdx]
				}
				var port int
				fmt.Sscanf(portStr, "%d", &port)
				sess.rtpPort = port
				sess.rtcpPort = port + 1
			}

			s.sendResponse(conn, "200 OK", cseq, sessionID,
				fmt.Sprintf("Transport: %s;server_port=20000-20001\r\nSession: %s\r\n", transport, sessionID))

		case "PLAY":
			if sess.streamKey == "" {
				s.sendResponse(conn, "400 Bad Request", cseq, sessionID, "")
				continue
			}

			// Subscribe and start streaming
			sess.subID = fmt.Sprintf("rtsp_out_%s", sessionID)
			sub := s.manager.Subscribe(sess.streamKey, sess.subID, 256)
			if sub == nil {
				s.sendResponse(conn, "500 Internal Server Error", cseq, sessionID, "")
				continue
			}

			s.sendResponse(conn, "200 OK", cseq, sessionID,
				fmt.Sprintf("Session: %s\r\nRTP-Info: url=%s;seq=0;rtptime=0\r\n", sessionID, uri))

			// Start sending RTP packets
			go s.streamRTP(sess, sub)

		case "TEARDOWN":
			s.sendResponse(conn, "200 OK", cseq, sessionID, "")
			return
		}
	}
}

func (s *Server) sendResponse(conn net.Conn, status, cseq, session, extra string) {
	resp := fmt.Sprintf("RTSP/1.0 %s\r\nCSeq: %s\r\nSession: %s\r\n%s\r\n",
		status, cseq, session, extra)
	conn.Write([]byte(resp))
}

func (s *Server) streamRTP(sess *rtspSession, sub *stream.OutputSubscriber) {
	// Resolve client UDP address for RTP
	host, _, _ := net.SplitHostPort(sess.conn.RemoteAddr().String())
	clientAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, sess.rtpPort))
	if err != nil {
		log.Printf("[RTSP-OUT] UDP resolve hatası: %v", err)
		return
	}

	udpConn, err := net.DialUDP("udp", nil, clientAddr)
	if err != nil {
		log.Printf("[RTSP-OUT] UDP bağlantı hatası: %v", err)
		return
	}
	defer udpConn.Close()

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}
			if pkt.Type == media.PacketTypeVideo || pkt.Type == media.PacketTypeAudio {
				rtpPkt := buildRTPPacket(pkt, sess.ssrc, &sess.seqNum, pkt.Timestamp)
				udpConn.Write(rtpPkt)
			}

		case <-sub.Done:
			return
		}
	}
}

func buildRTPPacket(pkt *media.Packet, ssrc uint32, seqNum *uint16, timestamp uint32) []byte {
	var payloadType byte
	if pkt.Type == media.PacketTypeVideo {
		payloadType = 96 // H.264
	} else {
		payloadType = 97 // AAC
	}

	// Strip FLV headers
	data := pkt.Data
	if pkt.Type == media.PacketTypeVideo && len(data) > 5 {
		data = data[5:]
	} else if pkt.Type == media.PacketTypeAudio && len(data) > 2 {
		data = data[2:]
	}

	*seqNum++
	header := make([]byte, 12)
	header[0] = 0x80                                     // Version 2
	header[1] = payloadType | 0x80                       // marker bit
	header[2] = byte(*seqNum >> 8)                       // sequence number
	header[3] = byte(*seqNum)
	header[4] = byte(timestamp >> 24)                    // timestamp
	header[5] = byte(timestamp >> 16)
	header[6] = byte(timestamp >> 8)
	header[7] = byte(timestamp)
	header[8] = byte(ssrc >> 24)                         // SSRC
	header[9] = byte(ssrc >> 16)
	header[10] = byte(ssrc >> 8)
	header[11] = byte(ssrc)

	return append(header, data...)
}

func extractStreamKey(uri string) string {
	// rtsp://host:port/live/streamkey
	parts := strings.Split(uri, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func generateSDP(streamKey, localAddr string) string {
	host, _, _ := net.SplitHostPort(localAddr)
	sdp := "v=0\r\n"
	sdp += fmt.Sprintf("o=- %d %d IN IP4 %s\r\n", time.Now().Unix(), time.Now().Unix(), host)
	sdp += fmt.Sprintf("s=%s\r\n", streamKey)
	sdp += "t=0 0\r\n"
	sdp += "a=control:*\r\n"

	// Video
	sdp += "m=video 0 RTP/AVP 96\r\n"
	sdp += "a=rtpmap:96 H264/90000\r\n"
	sdp += "a=fmtp:96 packetization-mode=1\r\n"
	sdp += "a=control:trackID=0\r\n"

	// Audio
	sdp += "m=audio 0 RTP/AVP 97\r\n"
	sdp += "a=rtpmap:97 MPEG4-GENERIC/44100/2\r\n"
	sdp += "a=control:trackID=1\r\n"

	return sdp
}

// Stop gracefully stops the server
func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sess := range s.sessions {
		if sess.subID != "" {
			s.manager.Unsubscribe(sess.streamKey, sess.subID)
		}
		sess.conn.Close()
	}
}
