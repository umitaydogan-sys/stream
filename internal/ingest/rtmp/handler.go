package rtmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	flv "github.com/fluxstream/fluxstream/internal/media/container/flv"
)

// StreamHandler defines the interface for handling stream events
type StreamHandler interface {
	OnPublish(streamKey string, conn net.Conn) error
	OnUnpublish(streamKey string)
	OnPacket(streamKey string, pkt *media.Packet)
}

// Handler handles a single RTMP connection
type Handler struct {
	conn      net.Conn
	reader    *ChunkReader
	writer    *ChunkWriter
	handler   StreamHandler
	streamKey string
	appName   string
	chunkSize uint32
}

// NewHandler creates a new RTMP connection handler
func NewHandler(conn net.Conn, handler StreamHandler) *Handler {
	return &Handler{
		conn:      conn,
		reader:    NewChunkReader(conn),
		writer:    NewChunkWriter(conn),
		handler:   handler,
		chunkSize: 4096,
	}
}

// Handle processes the RTMP connection
func (h *Handler) Handle() {
	defer h.conn.Close()
	defer func() {
		if h.streamKey != "" {
			h.handler.OnUnpublish(h.streamKey)
			log.Printf("[RTMP] Yayın sonlandı: %s", h.streamKey)
		}
	}()

	// Set read deadline for handshake
	h.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	// Perform handshake
	if err := DoHandshake(h.conn); err != nil {
		log.Printf("[RTMP] Handshake hatası: %v", err)
		return
	}

	// Clear deadline after handshake
	h.conn.SetReadDeadline(time.Time{})

	// Send server configuration
	h.sendSetChunkSize(h.chunkSize)
	h.sendAckSize(2500000)
	h.sendPeerBandwidth(2500000)

	// Set our writer chunk size
	h.writer.SetChunkSize(h.chunkSize)

	// Process messages
	for {
		msg, err := h.reader.ReadMessage()
		if err != nil {
			if err != io.EOF && !isConnectionClosed(err) {
				log.Printf("[RTMP] Mesaj okuma hatası: %v", err)
			}
			return
		}

		if err := h.processMessage(msg); err != nil {
			log.Printf("[RTMP] Mesaj işleme hatası: %v", err)
			return
		}
	}
}

func (h *Handler) processMessage(msg *Message) error {
	switch msg.TypeID {
	case MsgSetChunkSize:
		if len(msg.Data) >= 4 {
			size := binary.BigEndian.Uint32(msg.Data)
			h.reader.SetChunkSize(size)
		}
	case MsgAbort:
		// Ignore
	case MsgAck:
		// Ignore
	case MsgUserControl:
		// Ignore for now
	case MsgAckSize:
		// Acknowledge
	case MsgPeerBandwidth:
		// Ignore
	case MsgAMF0Command:
		return h.handleCommand(msg)
	case MsgAMF0Data:
		return h.handleDataMessage(msg)
	case MsgVideoData:
		return h.handleVideoData(msg)
	case MsgAudioData:
		return h.handleAudioData(msg)
	}
	return nil
}

func (h *Handler) handleCommand(msg *Message) error {
	r := bytes.NewReader(msg.Data)
	values, err := ReadAMF0(r)
	if err != nil && len(values) == 0 {
		return fmt.Errorf("read AMF0 command: %w", err)
	}

	if len(values) == 0 {
		return nil
	}

	cmdName := values[0].Str

	switch cmdName {
	case "connect":
		return h.handleConnect(values, msg)
	case "releaseStream":
		return h.sendResult(values, msg.StreamID)
	case "FCPublish":
		return h.sendResult(values, msg.StreamID)
	case "createStream":
		return h.handleCreateStream(values, msg)
	case "publish":
		return h.handlePublish(values, msg)
	case "FCUnpublish":
		return nil
	case "deleteStream":
		return nil
	case "_checkbw":
		return nil
	default:
		log.Printf("[RTMP] Bilinmeyen komut: %s", cmdName)
	}
	return nil
}

func (h *Handler) handleConnect(values []AMFValue, msg *Message) error {
	// Extract app name from connect object
	if len(values) > 2 {
		if obj := values[2]; obj.Type == AMF0Object {
			if app, ok := obj.Obj["app"]; ok {
				h.appName = app.Str
			}
		}
	}

	log.Printf("[RTMP] Connect: app=%s from=%s", h.appName, h.conn.RemoteAddr())

	// Send Window Ack Size
	h.sendAckSize(2500000)

	// Send Set Peer Bandwidth
	h.sendPeerBandwidth(2500000)

	// Send StreamBegin
	h.sendUserControl(0, 0)

	// Send _result for connect
	var buf bytes.Buffer
	WriteAMF0String(&buf, "_result")
	WriteAMF0Number(&buf, 1) // transaction ID

	// Properties object
	WriteAMF0ObjectStart(&buf)
	WriteAMF0Property(&buf, "fmsVer")
	buf.Write([]byte{AMF0String})
	writeShortString(&buf, "FMS/3,5,3,824")
	WriteAMF0Property(&buf, "capabilities")
	buf.Write([]byte{AMF0Number})
	numBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(numBuf, uint64(0x403F000000000000)) // 31.0
	buf.Write(numBuf)
	WriteAMF0ObjectEnd(&buf)

	// Information object
	WriteAMF0ObjectStart(&buf)
	WriteAMF0Property(&buf, "level")
	buf.Write([]byte{AMF0String})
	writeShortString(&buf, "status")
	WriteAMF0Property(&buf, "code")
	buf.Write([]byte{AMF0String})
	writeShortString(&buf, "NetConnection.Connect.Success")
	WriteAMF0Property(&buf, "description")
	buf.Write([]byte{AMF0String})
	writeShortString(&buf, "Connection succeeded.")
	WriteAMF0Property(&buf, "objectEncoding")
	WriteAMF0Number(&buf, 0)
	WriteAMF0ObjectEnd(&buf)

	return h.writer.WriteMessage(ChunkStreamCommand, 0, MsgAMF0Command, 0, buf.Bytes())
}

func (h *Handler) handleCreateStream(values []AMFValue, msg *Message) error {
	transID := float64(0)
	if len(values) > 1 {
		transID = values[1].Num
	}

	var buf bytes.Buffer
	WriteAMF0String(&buf, "_result")
	WriteAMF0Number(&buf, transID)
	WriteAMF0Null(&buf)
	WriteAMF0Number(&buf, 1) // stream ID = 1

	return h.writer.WriteMessage(ChunkStreamCommand, 0, MsgAMF0Command, 0, buf.Bytes())
}

func (h *Handler) handlePublish(values []AMFValue, msg *Message) error {
	if len(values) < 4 {
		return fmt.Errorf("publish: insufficient arguments")
	}

	streamKey := values[3].Str

	// Remove any query parameters from stream key
	if idx := strings.Index(streamKey, "?"); idx != -1 {
		streamKey = streamKey[:idx]
	}

	h.streamKey = streamKey

	log.Printf("[RTMP] Publish: key=%s from=%s", streamKey, h.conn.RemoteAddr())

	// Notify stream manager
	if err := h.handler.OnPublish(streamKey, h.conn); err != nil {
		// Send error
		h.sendStatus("error", "NetStream.Publish.BadName", err.Error(), msg.StreamID)
		return err
	}

	// Send onStatus
	h.sendStatus("status", "NetStream.Publish.Start", fmt.Sprintf("Publishing %s", streamKey), msg.StreamID)
	return nil
}

func (h *Handler) handleDataMessage(msg *Message) error {
	if h.streamKey == "" {
		return nil
	}

	pkt := &media.Packet{
		Type:       media.PacketTypeMeta,
		Timestamp:  msg.Timestamp,
		Data:       msg.Data,
		StreamKey:  h.streamKey,
		ReceivedAt: time.Now(),
	}
	h.handler.OnPacket(h.streamKey, pkt)
	return nil
}

func (h *Handler) handleVideoData(msg *Message) error {
	if h.streamKey == "" || len(msg.Data) < 2 {
		return nil
	}

	reader := flv.NewReader(nil) // We don't need the reader, just use ReadTag
	pkt, err := reader.ReadTag(0x09, uint32(len(msg.Data)), msg.Timestamp, msg.Data)
	if err != nil {
		return err
	}
	pkt.StreamKey = h.streamKey
	pkt.ReceivedAt = time.Now()
	h.handler.OnPacket(h.streamKey, pkt)
	return nil
}

func (h *Handler) handleAudioData(msg *Message) error {
	if h.streamKey == "" || len(msg.Data) < 1 {
		return nil
	}

	reader := flv.NewReader(nil)
	pkt, err := reader.ReadTag(0x08, uint32(len(msg.Data)), msg.Timestamp, msg.Data)
	if err != nil {
		return err
	}
	pkt.StreamKey = h.streamKey
	pkt.ReceivedAt = time.Now()
	h.handler.OnPacket(h.streamKey, pkt)
	return nil
}

// ─── Send helpers ─────────────────────────────────────────

func (h *Handler) sendSetChunkSize(size uint32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, size)
	h.writer.WriteMessage(ChunkStreamProtocol, 0, MsgSetChunkSize, 0, data)
}

func (h *Handler) sendAckSize(size uint32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, size)
	h.writer.WriteMessage(ChunkStreamProtocol, 0, MsgAckSize, 0, data)
}

func (h *Handler) sendPeerBandwidth(size uint32) {
	data := make([]byte, 5)
	binary.BigEndian.PutUint32(data, size)
	data[4] = 2 // dynamic
	h.writer.WriteMessage(ChunkStreamProtocol, 0, MsgPeerBandwidth, 0, data)
}

func (h *Handler) sendUserControl(eventType uint16, param uint32) {
	data := make([]byte, 6)
	binary.BigEndian.PutUint16(data, eventType)
	binary.BigEndian.PutUint32(data[2:], param)
	h.writer.WriteMessage(ChunkStreamProtocol, 0, MsgUserControl, 0, data)
}

func (h *Handler) sendStatus(level, code, description string, streamID uint32) {
	var buf bytes.Buffer
	WriteAMF0String(&buf, "onStatus")
	WriteAMF0Number(&buf, 0)
	WriteAMF0Null(&buf)

	WriteAMF0ObjectStart(&buf)
	WriteAMF0Property(&buf, "level")
	buf.Write([]byte{AMF0String})
	writeShortString(&buf, level)
	WriteAMF0Property(&buf, "code")
	buf.Write([]byte{AMF0String})
	writeShortString(&buf, code)
	WriteAMF0Property(&buf, "description")
	buf.Write([]byte{AMF0String})
	writeShortString(&buf, description)
	WriteAMF0ObjectEnd(&buf)

	h.writer.WriteMessage(ChunkStreamCommand, 0, MsgAMF0Command, streamID, buf.Bytes())
}

func (h *Handler) sendResult(values []AMFValue, streamID uint32) error {
	transID := float64(0)
	if len(values) > 1 {
		transID = values[1].Num
	}

	var buf bytes.Buffer
	WriteAMF0String(&buf, "_result")
	WriteAMF0Number(&buf, transID)
	WriteAMF0Null(&buf)
	WriteAMF0Null(&buf)

	return h.writer.WriteMessage(ChunkStreamCommand, 0, MsgAMF0Command, streamID, buf.Bytes())
}

func writeShortString(w io.Writer, s string) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(len(s)))
	w.Write(b)
	w.Write([]byte(s))
}

func isConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "closed") ||
		strings.Contains(errStr, "reset") ||
		strings.Contains(errStr, "broken pipe") ||
		strings.Contains(errStr, "EOF")
}
