package mpegts

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	ts "github.com/fluxstream/fluxstream/internal/media/container/ts"
	"github.com/fluxstream/fluxstream/internal/stream"
)

// Sender sends MPEG-TS over UDP (multicast/unicast)
type Sender struct {
	manager *stream.Manager
	targets map[string]*UDPTarget
	mu      sync.RWMutex
}

// UDPTarget represents an MPEG-TS UDP output
type UDPTarget struct {
	ID          string
	StreamKey   string
	DestAddr    string // e.g. "239.1.1.1:5000" for multicast or "192.168.1.100:5000"
	Status      string // "idle", "active", "error"
	PacketsSent int64
	BytesSent   int64
	subID       string
	conn        *net.UDPConn
	tsMuxer     *ts.Muxer
	stopCh      chan struct{}
}

// NewSender creates a new MPEG-TS UDP sender
func NewSender(manager *stream.Manager) *Sender {
	return &Sender{
		manager: manager,
		targets: make(map[string]*UDPTarget),
	}
}

// AddTarget adds a new UDP output target
func (s *Sender) AddTarget(id, streamKey, destAddr string) *UDPTarget {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := &UDPTarget{
		ID:        id,
		StreamKey: streamKey,
		DestAddr:  destAddr,
		Status:    "idle",
		tsMuxer:   ts.NewMuxer(),
		stopCh:    make(chan struct{}),
	}
	s.targets[id] = t
	return t
}

// RemoveTarget removes a UDP output target
func (s *Sender) RemoveTarget(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.targets[id]; ok {
		t.Stop()
		delete(s.targets, id)
	}
}

// StartTarget begins sending MPEG-TS over UDP
func (s *Sender) StartTarget(id string) error {
	s.mu.RLock()
	target, exists := s.targets[id]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("TS-UDP target not found: %s", id)
	}

	if !s.manager.IsLive(target.StreamKey) {
		return fmt.Errorf("source stream not live: %s", target.StreamKey)
	}

	go s.sendLoop(target)
	return nil
}

// Stop stops the UDP output
func (t *UDPTarget) Stop() {
	select {
	case <-t.stopCh:
	default:
		close(t.stopCh)
	}
	if t.conn != nil {
		t.conn.Close()
	}
	t.Status = "idle"
}

func (s *Sender) sendLoop(target *UDPTarget) {
	// Subscribe to stream
	target.subID = fmt.Sprintf("tsudp_out_%s", target.ID)
	sub := s.manager.Subscribe(target.StreamKey, target.subID, 256)
	if sub == nil {
		target.Status = "error"
		return
	}
	defer s.manager.Unsubscribe(target.StreamKey, target.subID)

	// Resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", target.DestAddr)
	if err != nil {
		target.Status = "error"
		log.Printf("[TS-UDP] Adres çözüleme hatası: %v", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		target.Status = "error"
		log.Printf("[TS-UDP] UDP bağlantı hatası: %v", err)
		return
	}
	target.conn = conn
	defer conn.Close()

	// Set multicast TTL if multicast address
	if addr.IP.IsMulticast() {
		conn.SetWriteBuffer(1024 * 1024) // 1MB write buffer
	}

	target.Status = "active"
	log.Printf("[TS-UDP] Başlatıldı: %s -> %s", target.StreamKey, target.DestAddr)

	// Buffer for TS packet bundling (7 TS packets per UDP = 1316 bytes)
	const tsPerUDP = 7
	tsBuf := make([]byte, 0, ts.TSPacketSize*tsPerUDP)

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				target.Status = "idle"
				return
			}

			if pkt.IsSequenceHeader {
				continue
			}

			// Mux to MPEG-TS
			mediaPkt := pkt.Clone()
			if pkt.Type == media.PacketTypeVideo && len(pkt.Data) > 5 {
				mediaPkt.Data = pkt.Data[5:]
			} else if pkt.Type == media.PacketTypeAudio && len(pkt.Data) > 2 {
				mediaPkt.Data = pkt.Data[2:]
			}

			tsData := target.tsMuxer.MuxPacket(mediaPkt)
			if tsData == nil {
				continue
			}

			// Bundle TS packets for efficient UDP sending
			tsBuf = append(tsBuf, tsData...)
			for len(tsBuf) >= ts.TSPacketSize*tsPerUDP {
				chunk := tsBuf[:ts.TSPacketSize*tsPerUDP]
				n, err := conn.Write(chunk)
				if err != nil {
					target.Status = "error"
					return
				}
				target.PacketsSent++
				target.BytesSent += int64(n)
				tsBuf = tsBuf[ts.TSPacketSize*tsPerUDP:]
			}

		case <-target.stopCh:
			// Flush remaining
			if len(tsBuf) > 0 {
				conn.Write(tsBuf)
			}
			target.Status = "idle"
			return
		case <-sub.Done:
			target.Status = "idle"
			return
		}
	}
}

// GetTargets returns all UDP output targets
func (s *Sender) GetTargets() []*UDPTarget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*UDPTarget
	for _, t := range s.targets {
		result = append(result, t)
	}
	return result
}

// StopAll stops all UDP output targets
func (s *Sender) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.targets {
		t.Stop()
	}
}

// StartAutoTarget creates and starts a default UDP target for a stream
func (s *Sender) StartAutoTarget(streamKey, destAddr string) error {
	id := fmt.Sprintf("auto_%s_%d", streamKey, time.Now().UnixNano())
	s.AddTarget(id, streamKey, destAddr)
	return s.StartTarget(id)
}
