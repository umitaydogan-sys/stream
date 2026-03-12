package rtp

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/stream"
)

// Sender sends RTP packets to configured destinations
type Sender struct {
	manager    *stream.Manager
	targets    map[string]*RTPTarget
	mu         sync.RWMutex
}

// RTPTarget represents an RTP output destination
type RTPTarget struct {
	ID          string
	StreamKey   string
	DestAddr    string // host:port
	DestPort    int
	Status      string // "idle", "active", "error"
	SSRC        uint32
	SeqNum      uint16
	PacketsSent int64
	BytesSent   int64
	conn        *net.UDPConn
	subID       string
	stopCh      chan struct{}
}

// NewSender creates a new RTP output sender
func NewSender(manager *stream.Manager) *Sender {
	return &Sender{
		manager: manager,
		targets: make(map[string]*RTPTarget),
	}
}

// AddTarget adds a new RTP output target
func (s *Sender) AddTarget(id, streamKey, destAddr string, destPort int) *RTPTarget {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := &RTPTarget{
		ID:        id,
		StreamKey: streamKey,
		DestAddr:  destAddr,
		DestPort:  destPort,
		Status:    "idle",
		SSRC:      uint32(time.Now().UnixNano() & 0xFFFFFFFF),
		stopCh:    make(chan struct{}),
	}
	s.targets[id] = t
	return t
}

// RemoveTarget removes an RTP output target
func (s *Sender) RemoveTarget(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.targets[id]; ok {
		t.Stop()
		delete(s.targets, id)
	}
}

// StartTarget begins sending RTP packets to the target
func (s *Sender) StartTarget(id string) error {
	s.mu.RLock()
	target, exists := s.targets[id]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("RTP target not found: %s", id)
	}

	if !s.manager.IsLive(target.StreamKey) {
		return fmt.Errorf("source stream not live: %s", target.StreamKey)
	}

	go s.sendLoop(target)
	return nil
}

// Stop stops sending RTP packets
func (t *RTPTarget) Stop() {
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

func (s *Sender) sendLoop(target *RTPTarget) {
	// Subscribe to stream
	target.subID = fmt.Sprintf("rtp_out_%s", target.ID)
	sub := s.manager.Subscribe(target.StreamKey, target.subID, 256)
	if sub == nil {
		target.Status = "error"
		return
	}
	defer s.manager.Unsubscribe(target.StreamKey, target.subID)

	// Resolve destination
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", target.DestAddr, target.DestPort))
	if err != nil {
		target.Status = "error"
		log.Printf("[RTP-OUT] Adres çözüleme hatası: %v", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		target.Status = "error"
		log.Printf("[RTP-OUT] UDP bağlantı hatası: %v", err)
		return
	}
	target.conn = conn
	defer conn.Close()

	target.Status = "active"
	log.Printf("[RTP-OUT] Başlatıldı: %s -> %s:%d", target.StreamKey, target.DestAddr, target.DestPort)

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				target.Status = "idle"
				return
			}

			rtpPkt := s.buildPacket(pkt, target)
			if rtpPkt == nil {
				continue
			}

			n, err := conn.Write(rtpPkt)
			if err != nil {
				target.Status = "error"
				return
			}
			target.PacketsSent++
			target.BytesSent += int64(n)

		case <-target.stopCh:
			target.Status = "idle"
			return
		case <-sub.Done:
			target.Status = "idle"
			return
		}
	}
}

func (s *Sender) buildPacket(pkt *media.Packet, target *RTPTarget) []byte {
	if pkt.IsSequenceHeader {
		return nil
	}

	var payloadType byte
	var clockRate uint32
	if pkt.Type == media.PacketTypeVideo {
		payloadType = 96 // H.264 dynamic
		clockRate = 90000
	} else if pkt.Type == media.PacketTypeAudio {
		payloadType = 97 // AAC dynamic
		clockRate = 44100
	} else {
		return nil
	}

	// Strip FLV headers
	data := pkt.Data
	if pkt.Type == media.PacketTypeVideo && len(data) > 5 {
		data = data[5:]
	} else if pkt.Type == media.PacketTypeAudio && len(data) > 2 {
		data = data[2:]
	}

	if len(data) == 0 {
		return nil
	}

	target.SeqNum++
	timestamp := uint32(float64(pkt.Timestamp) / 1000.0 * float64(clockRate))

	// RTP header (12 bytes)
	header := make([]byte, 12)
	header[0] = 0x80 // Version 2
	header[1] = payloadType
	if pkt.IsKeyframe {
		header[1] |= 0x80 // marker bit
	}
	header[2] = byte(target.SeqNum >> 8)
	header[3] = byte(target.SeqNum)
	header[4] = byte(timestamp >> 24)
	header[5] = byte(timestamp >> 16)
	header[6] = byte(timestamp >> 8)
	header[7] = byte(timestamp)
	header[8] = byte(target.SSRC >> 24)
	header[9] = byte(target.SSRC >> 16)
	header[10] = byte(target.SSRC >> 8)
	header[11] = byte(target.SSRC)

	return append(header, data...)
}

// GetTargets returns all RTP output targets
func (s *Sender) GetTargets() []*RTPTarget {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*RTPTarget
	for _, t := range s.targets {
		result = append(result, t)
	}
	return result
}

// StopAll stops all RTP output targets
func (s *Sender) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.targets {
		t.Stop()
	}
}
