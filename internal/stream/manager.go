package stream

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/output/dash"
	"github.com/fluxstream/fluxstream/internal/output/hls"
	"github.com/fluxstream/fluxstream/internal/storage"
)

// OutputSubscriber receives packets from an active stream
type OutputSubscriber struct {
	ID      string
	PacketC chan *media.Packet
	Done    chan struct{}
}

// Manager manages all active streams
type Manager struct {
	db         *storage.SQLiteDB
	hlsMuxer   *hls.Muxer
	dashMuxer  *dash.Muxer
	llhlsMuxer *hls.LLMuxer
	streams    map[string]*ActiveStream
	mu         sync.RWMutex
}

// ActiveStream represents a currently live stream
type ActiveStream struct {
	Key         string
	DBStream    *storage.Stream
	HLSMuxer    *hls.StreamMuxer
	DASHMuxer   *dash.StreamMuxer
	LLHLSMuxer  *hls.LLStreamMuxer
	StartedAt   time.Time
	PacketCount int64
	BytesIn     int64
	Conn        net.Conn

	subscribers  map[string]*OutputSubscriber
	videoSeq     *media.Packet
	audioSeq     *media.Packet
	lastKeyframe *media.Packet
	subMu        sync.RWMutex
}

// Subscribe adds an output subscriber to a live stream
func (m *Manager) Subscribe(streamKey, subscriberID string, bufSize int) *OutputSubscriber {
	m.mu.RLock()
	active, exists := m.streams[streamKey]
	m.mu.RUnlock()
	if !exists {
		return nil
	}

	sub := &OutputSubscriber{
		ID:      subscriberID,
		PacketC: make(chan *media.Packet, bufSize),
		Done:    make(chan struct{}),
	}

	active.subMu.Lock()
	active.subscribers[subscriberID] = sub
	var cached []*media.Packet
	if active.videoSeq != nil {
		cached = append(cached, active.videoSeq.Clone())
	}
	if active.audioSeq != nil {
		cached = append(cached, active.audioSeq.Clone())
	}
	if active.lastKeyframe != nil {
		cached = append(cached, active.lastKeyframe.Clone())
	}
	active.subMu.Unlock()

	for _, pkt := range cached {
		select {
		case sub.PacketC <- pkt:
		default:
		}
	}

	log.Printf("[Stream] Subscriber eklendi: %s -> %s", subscriberID, streamKey)
	return sub
}

// Unsubscribe removes an output subscriber
func (m *Manager) Unsubscribe(streamKey, subscriberID string) {
	m.mu.RLock()
	active, exists := m.streams[streamKey]
	m.mu.RUnlock()
	if !exists {
		return
	}

	active.subMu.Lock()
	if sub, ok := active.subscribers[subscriberID]; ok {
		close(sub.Done)
		delete(active.subscribers, subscriberID)
	}
	active.subMu.Unlock()
}

// NewManager creates a new stream manager
func NewManager(db *storage.SQLiteDB, hlsMuxer *hls.Muxer) *Manager {
	return &Manager{
		db:       db,
		hlsMuxer: hlsMuxer,
		streams:  make(map[string]*ActiveStream),
	}
}

// SetOutputMuxers attaches optional output muxers to the stream manager.
func (m *Manager) SetOutputMuxers(dashMuxer *dash.Muxer, llhlsMuxer *hls.LLMuxer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dashMuxer = dashMuxer
	m.llhlsMuxer = llhlsMuxer
}

// OnPublish handles a new publish event from RTMP
func (m *Manager) OnPublish(streamKey string, conn net.Conn) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if stream key exists in DB
	dbStream, err := m.db.GetStreamByKey(streamKey)
	if err != nil {
		return fmt.Errorf("db lookup: %w", err)
	}
	if dbStream == nil {
		return fmt.Errorf("stream key bulunamadı: %s", streamKey)
	}

	// Check if already live
	if _, exists := m.streams[streamKey]; exists {
		return fmt.Errorf("stream zaten yayında: %s", streamKey)
	}

	// Create HLS output
	hlsStream := m.hlsMuxer.AddStream(streamKey)
	var dashStream *dash.StreamMuxer
	if m.dashMuxer != nil {
		dashStream = m.dashMuxer.AddStream(streamKey)
	}
	var llhlsStream *hls.LLStreamMuxer
	if m.llhlsMuxer != nil {
		llhlsStream = m.llhlsMuxer.AddStream(streamKey)
	}

	// Update DB status
	m.db.UpdateStreamStatus(streamKey, "live", "rtmp")
	m.db.AddLog("INFO", "stream", fmt.Sprintf("Yayın başladı: %s (%s)", dbStream.Name, streamKey))

	m.streams[streamKey] = &ActiveStream{
		Key:         streamKey,
		DBStream:    dbStream,
		HLSMuxer:    hlsStream,
		DASHMuxer:   dashStream,
		LLHLSMuxer:  llhlsStream,
		StartedAt:   time.Now(),
		Conn:        conn,
		subscribers: make(map[string]*OutputSubscriber),
	}

	log.Printf("[Stream] Yayın başladı: %s (%s)", dbStream.Name, streamKey)
	return nil
}

// OnUnpublish handles stream disconnect
func (m *Manager) OnUnpublish(streamKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	active, exists := m.streams[streamKey]
	if !exists {
		return
	}

	// Close all subscribers
	active.subMu.Lock()
	for id, sub := range active.subscribers {
		close(sub.Done)
		delete(active.subscribers, id)
	}
	active.subMu.Unlock()

	// Close HLS muxer
	m.hlsMuxer.RemoveStream(streamKey)
	if m.dashMuxer != nil {
		m.dashMuxer.RemoveStream(streamKey)
	}
	if m.llhlsMuxer != nil {
		m.llhlsMuxer.RemoveStream(streamKey)
	}
	m.cleanupOutputFiles(streamKey)

	// Update DB status
	m.db.UpdateStreamStatus(streamKey, "offline", "")

	duration := time.Since(active.StartedAt)
	m.db.AddLog("INFO", "stream", fmt.Sprintf("Yayın sonlandı: %s (süre: %s)", streamKey, duration.Round(time.Second)))

	delete(m.streams, streamKey)
	log.Printf("[Stream] Yayın sonlandı: %s (süre: %s)", streamKey, duration.Round(time.Second))
}

// OnPacket handles incoming media packets
func (m *Manager) OnPacket(streamKey string, pkt *media.Packet) {
	if pkt == nil || pkt.TrackID != 0 {
		return
	}

	m.mu.RLock()
	active, exists := m.streams[streamKey]
	m.mu.RUnlock()

	if !exists {
		return
	}

	active.PacketCount++
	active.BytesIn += int64(len(pkt.Data))

	active.subMu.Lock()
	if pkt.IsSequenceHeader {
		if pkt.Type == media.PacketTypeVideo {
			active.videoSeq = pkt.Clone()
		} else if pkt.Type == media.PacketTypeAudio {
			active.audioSeq = pkt.Clone()
		}
	} else if pkt.Type == media.PacketTypeVideo && pkt.IsKeyframe {
		active.lastKeyframe = pkt.Clone()
	}
	active.subMu.Unlock()

	// Write to HLS
	if err := active.HLSMuxer.WritePacket(pkt); err != nil {
		log.Printf("[Stream] HLS yazma hatası (%s): %v", streamKey, err)
	}
	if active.DASHMuxer != nil {
		if err := active.DASHMuxer.WritePacket(pkt); err != nil {
			log.Printf("[Stream] DASH yazma hatası (%s): %v", streamKey, err)
		}
	}
	if active.LLHLSMuxer != nil {
		if err := active.LLHLSMuxer.WritePacket(pkt); err != nil {
			log.Printf("[Stream] LL-HLS yazma hatası (%s): %v", streamKey, err)
		}
	}

	// Fan out to subscribers (non-blocking)
	active.subMu.RLock()
	for _, sub := range active.subscribers {
		select {
		case sub.PacketC <- pkt:
		default:
			// subscriber too slow, drop packet
		}
	}
	active.subMu.RUnlock()
}

// GetActiveStreams returns all active streams
func (m *Manager) GetActiveStreams() []*ActiveStream {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*ActiveStream
	for _, s := range m.streams {
		result = append(result, s)
	}
	return result
}

// GetActiveStream returns a single active stream by key
func (m *Manager) GetActiveStream(key string) *ActiveStream {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.streams[key]
}

// IsLive checks if a stream key is currently live
func (m *Manager) IsLive(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.streams[key]
	return exists
}

// StopAll stops all active streams
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, active := range m.streams {
		// Close subscribers
		active.subMu.Lock()
		for id, sub := range active.subscribers {
			close(sub.Done)
			delete(active.subscribers, id)
		}
		active.subMu.Unlock()

		m.hlsMuxer.RemoveStream(key)
		if m.dashMuxer != nil {
			m.dashMuxer.RemoveStream(key)
		}
		if m.llhlsMuxer != nil {
			m.llhlsMuxer.RemoveStream(key)
		}
		m.cleanupOutputFiles(key)
		m.db.UpdateStreamStatus(key, "offline", "")
		if active.Conn != nil {
			active.Conn.Close()
		}
		delete(m.streams, key)
	}
}

// GetStats returns current server stats
func (m *Manager) GetStats() storage.ServerStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := storage.ServerStats{
		ActiveStreams: len(m.streams),
	}

	for _, s := range m.streams {
		stats.BandwidthIn += s.BytesIn
	}

	return stats
}

func (m *Manager) cleanupOutputFiles(streamKey string) {
	if m.hlsMuxer != nil {
		_ = os.RemoveAll(filepath.Join(m.hlsMuxer.GetOutputDir(), streamKey))
	}
	if m.dashMuxer != nil {
		_ = os.RemoveAll(filepath.Join(m.dashMuxer.GetOutputDir(), streamKey))
	}
}
