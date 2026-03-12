package flv

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/stream"
)

// Server serves HTTP-FLV streams via chunked transfer encoding
type Server struct {
	manager   *stream.Manager
	gopCache  bool
}

// NewServer creates a new HTTP-FLV server
func NewServer(manager *stream.Manager, gopCache bool) *Server {
	return &Server{
		manager:  manager,
		gopCache: gopCache,
	}
}

// HandleFLV is the HTTP handler for /flv/{key}
func (s *Server) HandleFLV(w http.ResponseWriter, r *http.Request) {
	// Extract stream key
	path := r.URL.Path
	key := strings.TrimPrefix(path, "/flv/")
	key = strings.TrimRight(key, "/")
	if key == "" {
		http.Error(w, "Stream key required", http.StatusBadRequest)
		return
	}

	if !s.manager.IsLive(key) {
		http.Error(w, "Stream not live", http.StatusNotFound)
		return
	}

	// Set headers for chunked FLV
	w.Header().Set("Content-Type", "video/x-flv")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Cache-Control", "no-cache, no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Write FLV header
	flvHeader := buildFLVHeader(true, true)
	w.Write(flvHeader)
	flusher.Flush()

	// Subscribe to stream
	subID := fmt.Sprintf("flv_%s_%d", r.RemoteAddr, time.Now().UnixNano())
	sub := s.manager.Subscribe(key, subID, 256)
	if sub == nil {
		http.Error(w, "Subscribe failed", http.StatusInternalServerError)
		return
	}
	defer s.manager.Unsubscribe(key, subID)

	log.Printf("[HTTP-FLV] İzleyici bağlandı: %s -> %s", r.RemoteAddr, key)

	// Stream packets as FLV tags
	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}
			tag := buildFLVTag(pkt)
			if tag == nil {
				continue
			}
			if _, err := w.Write(tag); err != nil {
				return
			}
			flusher.Flush()

		case <-sub.Done:
			return

		case <-r.Context().Done():
			return
		}
	}
}

func buildFLVHeader(hasVideo, hasAudio bool) []byte {
	header := make([]byte, 13) // 9 byte header + 4 byte prev tag size
	copy(header[0:3], "FLV")
	header[3] = 0x01 // version
	flags := byte(0)
	if hasAudio {
		flags |= 0x04
	}
	if hasVideo {
		flags |= 0x01
	}
	header[4] = flags
	binary.BigEndian.PutUint32(header[5:9], 9) // header size
	// previous tag size 0
	return header
}

func buildFLVTag(pkt *media.Packet) []byte {
	if pkt == nil || len(pkt.Data) == 0 {
		return nil
	}

	var tagType byte
	switch pkt.Type {
	case media.PacketTypeVideo:
		tagType = 0x09
	case media.PacketTypeAudio:
		tagType = 0x08
	case media.PacketTypeMeta:
		tagType = 0x12
	default:
		return nil
	}

	dataSize := len(pkt.Data)
	tagSize := 11 + dataSize // 11 byte tag header + data

	tag := make([]byte, tagSize+4) // +4 for previous tag size

	// Tag header
	tag[0] = tagType

	// Data size (24-bit)
	tag[1] = byte(dataSize >> 16)
	tag[2] = byte(dataSize >> 8)
	tag[3] = byte(dataSize)

	// Timestamp (24-bit + 8-bit extension)
	ts := pkt.Timestamp
	tag[4] = byte(ts >> 16)
	tag[5] = byte(ts >> 8)
	tag[6] = byte(ts)
	tag[7] = byte(ts >> 24) // timestamp extended

	// Stream ID (always 0)
	tag[8] = 0
	tag[9] = 0
	tag[10] = 0

	// Data
	copy(tag[11:], pkt.Data)

	// Previous tag size
	prevSize := uint32(tagSize)
	binary.BigEndian.PutUint32(tag[tagSize:], prevSize)

	return tag
}

// FLVClientCount returns the number of active FLV viewers (for stats)
type FLVStats struct {
	mu      sync.RWMutex
	clients int
}
