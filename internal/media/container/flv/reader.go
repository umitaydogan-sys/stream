package flv

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/fluxstream/fluxstream/internal/media"
)

// FLV Tag types
const (
	TagAudio  = 0x08
	TagVideo  = 0x09
	TagScript = 0x12
)

// Video frame types
const (
	FrameKeyframe   = 1
	FrameInterframe = 2
)

// Reader reads FLV tags from an io.Reader
type Reader struct {
	r io.Reader
}

// NewReader creates a new FLV reader
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// ReadTag reads the next FLV tag and returns a media.Packet
func (f *Reader) ReadTag(tagType byte, dataSize uint32, timestamp uint32, data []byte) (*media.Packet, error) {
	pkt := &media.Packet{
		Timestamp: timestamp,
		Data:      data,
	}

	switch tagType {
	case TagVideo:
		pkt.Type = media.PacketTypeVideo
		if len(data) > 0 {
			frameType := (data[0] >> 4) & 0x0F
			pkt.IsKeyframe = (frameType == FrameKeyframe)

			codecID := data[0] & 0x0F
			if codecID == byte(media.VideoCodecH264) && len(data) > 1 {
				avcPacketType := data[1]
				pkt.IsSequenceHeader = (avcPacketType == 0) // AVC sequence header
			}
		}

	case TagAudio:
		pkt.Type = media.PacketTypeAudio
		if len(data) > 0 {
			codecID := (data[0] >> 4) & 0x0F
			if codecID == byte(media.AudioCodecAAC) && len(data) > 1 {
				aacPacketType := data[1]
				pkt.IsSequenceHeader = (aacPacketType == 0) // AAC sequence header
			}
		}

	case TagScript:
		pkt.Type = media.PacketTypeMeta
	}

	return pkt, nil
}

// ParseFLVHeader reads and validates an FLV header
func ParseFLVHeader(r io.Reader) error {
	header := make([]byte, 9)
	if _, err := io.ReadFull(r, header); err != nil {
		return fmt.Errorf("read flv header: %w", err)
	}

	if header[0] != 'F' || header[1] != 'L' || header[2] != 'V' {
		return fmt.Errorf("invalid FLV signature: %x%x%x", header[0], header[1], header[2])
	}

	// Read first PreviousTagSize (4 bytes, should be 0)
	prev := make([]byte, 4)
	if _, err := io.ReadFull(r, prev); err != nil {
		return fmt.Errorf("read first prev tag size: %w", err)
	}

	return nil
}

// ReadTagHeader reads an FLV tag header (11 bytes) and returns type, data size, timestamp
func ReadTagHeader(r io.Reader) (tagType byte, dataSize uint32, timestamp uint32, err error) {
	header := make([]byte, 11)
	if _, err = io.ReadFull(r, header); err != nil {
		return
	}

	tagType = header[0]
	dataSize = uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3])

	// Timestamp: 3 bytes + 1 extended byte
	timestamp = uint32(header[4])<<16 | uint32(header[5])<<8 | uint32(header[6])
	timestamp |= uint32(header[7]) << 24 // timestamp extended

	return
}

// ReadTagData reads the tag body and trailing PreviousTagSize
func ReadTagData(r io.Reader, dataSize uint32) ([]byte, error) {
	data := make([]byte, dataSize)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("read tag data: %w", err)
	}

	// Read PreviousTagSize (4 bytes)
	prev := make([]byte, 4)
	if _, err := io.ReadFull(r, prev); err != nil {
		return nil, fmt.Errorf("read prev tag size: %w", err)
	}

	return data, nil
}

// PutUint24 encodes a uint32 as 3 bytes big-endian
func PutUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

// Uint24 decodes 3 bytes big-endian to uint32
func Uint24(b []byte) uint32 {
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// BuildFLVTag constructs an FLV tag from parts
func BuildFLVTag(tagType byte, timestamp uint32, data []byte) []byte {
	dataSize := len(data)
	tag := make([]byte, 11+dataSize+4)

	tag[0] = tagType
	PutUint24(tag[1:4], uint32(dataSize))

	// Timestamp
	tag[4] = byte(timestamp >> 16)
	tag[5] = byte(timestamp >> 8)
	tag[6] = byte(timestamp)
	tag[7] = byte(timestamp >> 24) // extended

	// StreamID = 0
	tag[8] = 0
	tag[9] = 0
	tag[10] = 0

	copy(tag[11:], data)

	// PreviousTagSize
	binary.BigEndian.PutUint32(tag[11+dataSize:], uint32(11+dataSize))

	return tag
}
