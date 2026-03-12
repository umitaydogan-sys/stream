package relay

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/stream"
)

// Target represents an RTMP relay target (YouTube, Twitch, etc.)
type Target struct {
	ID        string
	StreamKey string    // source stream key
	URL       string    // rtmp://a.rtmp.youtube.com/live2/xxxx
	Name      string    // "YouTube Relay"
	Status    string    // "idle", "connecting", "live", "error"
	Error     string
	StartedAt time.Time
	BytesSent int64
	conn      net.Conn
	stopCh    chan struct{}
}

// Manager manages RTMP relay targets
type Manager struct {
	streamMgr *stream.Manager
	targets   map[string]*Target
	mu        sync.RWMutex
}

// NewManager creates a new relay manager
func NewManager(streamMgr *stream.Manager) *Manager {
	return &Manager{
		streamMgr: streamMgr,
		targets:   make(map[string]*Target),
	}
}

// AddTarget adds a new relay target
func (m *Manager) AddTarget(id, streamKey, url, name string) *Target {
	m.mu.Lock()
	defer m.mu.Unlock()

	t := &Target{
		ID:        id,
		StreamKey: streamKey,
		URL:       url,
		Name:      name,
		Status:    "idle",
		stopCh:    make(chan struct{}),
	}
	m.targets[id] = t
	return t
}

// RemoveTarget removes a relay target
func (m *Manager) RemoveTarget(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t, ok := m.targets[id]; ok {
		t.Stop()
		delete(m.targets, id)
	}
}

// StartTarget begins relaying to the target
func (m *Manager) StartTarget(id string) error {
	m.mu.RLock()
	target, exists := m.targets[id]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("relay target not found: %s", id)
	}

	if !m.streamMgr.IsLive(target.StreamKey) {
		return fmt.Errorf("source stream not live: %s", target.StreamKey)
	}

	go m.relayLoop(target)
	return nil
}

// StopTarget stops relaying to the target
func (m *Manager) StopTarget(id string) {
	m.mu.RLock()
	target, exists := m.targets[id]
	m.mu.RUnlock()

	if exists {
		target.Stop()
	}
}

// GetTargets returns all relay targets
func (m *Manager) GetTargets() []*Target {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Target
	for _, t := range m.targets {
		result = append(result, t)
	}
	return result
}

// Stop stops the relay target
func (t *Target) Stop() {
	select {
	case <-t.stopCh:
		// already stopped
	default:
		close(t.stopCh)
	}
	if t.conn != nil {
		t.conn.Close()
	}
	t.Status = "idle"
}

func (m *Manager) relayLoop(target *Target) {
	target.Status = "connecting"
	target.StartedAt = time.Now()

	// Subscribe to source stream
	subID := fmt.Sprintf("relay_%s", target.ID)
	sub := m.streamMgr.Subscribe(target.StreamKey, subID, 512)
	if sub == nil {
		target.Status = "error"
		target.Error = "subscribe failed"
		return
	}
	defer m.streamMgr.Unsubscribe(target.StreamKey, subID)

	// Connect to target RTMP server
	conn, err := net.DialTimeout("tcp", parseRTMPHost(target.URL), 10*time.Second)
	if err != nil {
		target.Status = "error"
		target.Error = fmt.Sprintf("connect failed: %v", err)
		log.Printf("[Relay] Bağlantı hatası %s: %v", target.Name, err)
		return
	}
	target.conn = conn
	defer conn.Close()

	// Perform RTMP handshake
	if err := performRelayHandshake(conn); err != nil {
		target.Status = "error"
		target.Error = fmt.Sprintf("handshake failed: %v", err)
		return
	}

	// Send connect + publish commands
	app, stream := parseRTMPURL(target.URL)
	if err := sendRelayConnect(conn, app); err != nil {
		target.Status = "error"
		target.Error = fmt.Sprintf("connect cmd failed: %v", err)
		return
	}

	if err := sendRelayPublish(conn, stream); err != nil {
		target.Status = "error"
		target.Error = fmt.Sprintf("publish cmd failed: %v", err)
		return
	}

	target.Status = "live"
	log.Printf("[Relay] Başlatıldı: %s -> %s", target.StreamKey, target.Name)

	// Relay packets
	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				target.Status = "idle"
				return
			}
			data := buildRTMPMessage(pkt)
			if _, err := conn.Write(data); err != nil {
				target.Status = "error"
				target.Error = fmt.Sprintf("write error: %v", err)
				return
			}
			target.BytesSent += int64(len(data))

		case <-target.stopCh:
			target.Status = "idle"
			return

		case <-sub.Done:
			target.Status = "idle"
			return
		}
	}
}

func performRelayHandshake(conn net.Conn) error {
	// C0 + C1
	c0c1 := make([]byte, 1537)
	c0c1[0] = 3 // RTMP version
	binary.BigEndian.PutUint32(c0c1[1:5], uint32(time.Now().Unix()))
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	if _, err := conn.Write(c0c1); err != nil {
		return err
	}

	// Read S0 + S1
	s0s1 := make([]byte, 1537)
	if _, err := readFull(conn, s0s1); err != nil {
		return err
	}

	// C2 (echo S1)
	c2 := make([]byte, 1536)
	copy(c2, s0s1[1:])
	if _, err := conn.Write(c2); err != nil {
		return err
	}

	// Read S2
	s2 := make([]byte, 1536)
	if _, err := readFull(conn, s2); err != nil {
		return err
	}

	conn.SetDeadline(time.Time{})
	return nil
}

func readFull(conn net.Conn, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := conn.Read(buf[total:])
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}

func sendRelayConnect(conn net.Conn, app string) error {
	// AMF0 connect command
	payload := amf0String("connect")
	payload = append(payload, amf0Number(1)...) // transaction ID

	// Command object
	payload = append(payload, 0x03) // Object marker
	payload = append(payload, amf0Property("app", app)...)
	payload = append(payload, amf0Property("type", "nonprivate")...)
	payload = append(payload, amf0Property("flashVer", "FMLE/3.0")...)
	payload = append(payload, amf0Property("tcUrl", "rtmp://localhost/"+app)...)
	payload = append(payload, 0x00, 0x00, 0x09) // Object end

	return writeRTMPChunk(conn, 3, 0x14, 0, payload)
}

func sendRelayPublish(conn net.Conn, streamName string) error {
	// createStream
	cs := amf0String("createStream")
	cs = append(cs, amf0Number(2)...)
	cs = append(cs, 0x05) // null
	writeRTMPChunk(conn, 3, 0x14, 0, cs)

	// Read response (simplified - skip)
	time.Sleep(100 * time.Millisecond)

	// publish
	pub := amf0String("publish")
	pub = append(pub, amf0Number(0)...)
	pub = append(pub, 0x05)                          // null
	pub = append(pub, amf0String(streamName)...)
	pub = append(pub, amf0String("live")...)

	return writeRTMPChunk(conn, 8, 0x14, 1, pub)
}

func writeRTMPChunk(conn net.Conn, csid byte, msgType byte, streamID uint32, payload []byte) error {
	// Format 0 chunk header (12 bytes)
	header := make([]byte, 12)
	header[0] = csid // fmt=0, csid
	// timestamp = 0 (3 bytes)
	header[4] = byte(len(payload) >> 16)
	header[5] = byte(len(payload) >> 8)
	header[6] = byte(len(payload))
	header[7] = msgType
	binary.LittleEndian.PutUint32(header[8:12], streamID)

	data := append(header, payload...)
	_, err := conn.Write(data)
	return err
}

func buildRTMPMessage(pkt *media.Packet) []byte {
	var msgType byte
	switch pkt.Type {
	case media.PacketTypeVideo:
		msgType = 0x09
	case media.PacketTypeAudio:
		msgType = 0x08
	case media.PacketTypeMeta:
		msgType = 0x12
	default:
		return nil
	}

	payload := pkt.Data
	header := make([]byte, 12)
	header[0] = 0x06 // fmt=0, csid=6

	// Timestamp (3 bytes)
	ts := pkt.Timestamp
	header[1] = byte(ts >> 16)
	header[2] = byte(ts >> 8)
	header[3] = byte(ts)

	// Message length (3 bytes)
	header[4] = byte(len(payload) >> 16)
	header[5] = byte(len(payload) >> 8)
	header[6] = byte(len(payload))

	// Message type
	header[7] = msgType

	// Stream ID (4 bytes LE) = 1
	binary.LittleEndian.PutUint32(header[8:12], 1)

	return append(header, payload...)
}

// AMF0 encoding helpers
func amf0String(s string) []byte {
	buf := make([]byte, 3+len(s))
	buf[0] = 0x02
	binary.BigEndian.PutUint16(buf[1:3], uint16(len(s)))
	copy(buf[3:], s)
	return buf
}

func amf0Number(n float64) []byte {
	buf := make([]byte, 9)
	buf[0] = 0x00
	bits := math.Float64bits(n)
	binary.BigEndian.PutUint64(buf[1:9], bits)
	return buf
}

func amf0Property(key, value string) []byte {
	// Property name (without type marker)
	buf := make([]byte, 2+len(key))
	binary.BigEndian.PutUint16(buf[0:2], uint16(len(key)))
	copy(buf[2:], key)
	// Property value (string)
	buf = append(buf, amf0String(value)...)
	return buf
}

// URL parsing helpers
func parseRTMPHost(url string) string {
	// rtmp://host:port/app/key -> host:port
	s := url
	if len(s) > 7 && s[:7] == "rtmp://" {
		s = s[7:]
	}
	// Find first /
	for i, c := range s {
		if c == '/' {
			host := s[:i]
			if !hasPort(host) {
				return host + ":1935"
			}
			return host
		}
	}
	if !hasPort(s) {
		return s + ":1935"
	}
	return s
}

func hasPort(s string) bool {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ':' {
			return true
		}
		if s[i] == ']' {
			return false
		}
	}
	return false
}

func parseRTMPURL(url string) (app, streamName string) {
	// rtmp://host/app/stream_key
	s := url
	if len(s) > 7 && s[:7] == "rtmp://" {
		s = s[7:]
	}
	// Skip host
	idx := 0
	for i, c := range s {
		if c == '/' {
			idx = i + 1
			break
		}
	}
	rest := s[idx:]
	// Split app/streamKey
	for i, c := range rest {
		if c == '/' {
			return rest[:i], rest[i+1:]
		}
	}
	return rest, ""
}

// StopAll stops all relay targets
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, t := range m.targets {
		t.Stop()
	}
}
