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

const (
	enhancedAudioFormat = 9

	videoExTypeSequenceStart = 0
	videoExTypeCodedFrames   = 1
	videoExTypeSequenceEnd   = 2
	videoExTypeFramesX       = 3
	videoExTypeMetadata      = 4
	videoExTypeMultitrack    = 6

	audioExTypeSequenceStart = 0
	audioExTypeCodedFrames   = 1
	audioExTypeSequenceEnd   = 2
	audioExTypeMultitrack    = 5
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
	switch tagType {
	case TagVideo:
		if isEnhancedVideoTag(data) {
			return f.readEnhancedVideoTag(timestamp, data)
		}
		return f.readLegacyVideoTag(timestamp, data), nil

	case TagAudio:
		if isEnhancedAudioTag(data) {
			return f.readEnhancedAudioTag(timestamp, data)
		}
		return f.readLegacyAudioTag(timestamp, data), nil

	case TagScript:
		return &media.Packet{
			Type:      media.PacketTypeMeta,
			Timestamp: timestamp,
			Data:      data,
		}, nil
	}

	return &media.Packet{
		Timestamp: timestamp,
		Data:      data,
	}, nil
}

func (f *Reader) readLegacyVideoTag(timestamp uint32, data []byte) *media.Packet {
	pkt := &media.Packet{
		Type:      media.PacketTypeVideo,
		Timestamp: timestamp,
		Data:      data,
	}
	if len(data) == 0 {
		return pkt
	}

	frameType := (data[0] >> 4) & 0x0F
	pkt.IsKeyframe = (frameType == FrameKeyframe)

	codecID := data[0] & 0x0F
	if codecID == byte(media.VideoCodecH264) && len(data) > 1 {
		pkt.IsSequenceHeader = (data[1] == 0) // AVC sequence header
	}
	return pkt
}

func (f *Reader) readLegacyAudioTag(timestamp uint32, data []byte) *media.Packet {
	pkt := &media.Packet{
		Type:      media.PacketTypeAudio,
		Timestamp: timestamp,
		Data:      data,
	}
	if len(data) == 0 {
		return pkt
	}

	codecID := (data[0] >> 4) & 0x0F
	if codecID == byte(media.AudioCodecAAC) && len(data) > 1 {
		pkt.IsSequenceHeader = (data[1] == 0) // AAC sequence header
	}
	return pkt
}

func (f *Reader) readEnhancedVideoTag(timestamp uint32, data []byte) (*media.Packet, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("enhanced video tag too short")
	}

	packetType := data[0] & 0x0F
	trackID := uint8(0)
	fourCCOffset := 1
	bodyOffset := 5

	if packetType == videoExTypeMultitrack {
		if len(data) < 7 {
			return nil, fmt.Errorf("enhanced multitrack video tag too short")
		}
		multitrackType := data[1] >> 4
		if multitrackType != 0 {
			return nil, nil
		}
		packetType = data[1] & 0x0F
		fourCCOffset = 2
		bodyOffset = 7
		trackID = data[6]
	}

	fourCC := readFourCC(data, fourCCOffset)
	switch packetType {
	case videoExTypeSequenceStart:
		if fourCC != "avc1" {
			return nil, nil
		}
		body := data[bodyOffset:]
		converted := make([]byte, 5+len(body))
		converted[0] = 0x17
		converted[1] = 0x00
		copy(converted[5:], body)
		return &media.Packet{
			Type:             media.PacketTypeVideo,
			Timestamp:        timestamp,
			Data:             converted,
			IsKeyframe:       true,
			IsSequenceHeader: true,
			TrackID:          trackID,
			IsEnhanced:       true,
			FourCC:           fourCC,
		}, nil

	case videoExTypeCodedFrames:
		if fourCC != "avc1" {
			return nil, nil
		}
		if len(data) < bodyOffset+3 {
			return nil, fmt.Errorf("enhanced AVC coded-frames tag too short")
		}
		cts := data[bodyOffset : bodyOffset+3]
		payload := data[bodyOffset+3:]
		return buildEnhancedAVCPacket(timestamp, trackID, fourCC, cts, payload), nil

	case videoExTypeFramesX:
		if fourCC != "avc1" {
			return nil, nil
		}
		payload := data[bodyOffset:]
		return buildEnhancedAVCPacket(timestamp, trackID, fourCC, []byte{0x00, 0x00, 0x00}, payload), nil

	case videoExTypeSequenceEnd, videoExTypeMetadata:
		return nil, nil

	default:
		return nil, nil
	}
}

func (f *Reader) readEnhancedAudioTag(timestamp uint32, data []byte) (*media.Packet, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("enhanced audio tag too short")
	}

	packetType := data[0] & 0x0F
	trackID := uint8(0)
	fourCCOffset := 1
	bodyOffset := 5

	if packetType == audioExTypeMultitrack {
		if len(data) < 7 {
			return nil, fmt.Errorf("enhanced multitrack audio tag too short")
		}
		multitrackType := data[1] >> 4
		if multitrackType != 0 {
			return nil, nil
		}
		packetType = data[1] & 0x0F
		fourCCOffset = 2
		bodyOffset = 7
		trackID = data[6]
	}

	fourCC := readFourCC(data, fourCCOffset)
	if fourCC != "mp4a" {
		if packetType == audioExTypeSequenceEnd {
			return nil, nil
		}
		return nil, nil
	}

	switch packetType {
	case audioExTypeSequenceStart:
		body := data[bodyOffset:]
		converted := make([]byte, 2+len(body))
		converted[0] = 0xAF
		converted[1] = 0x00
		copy(converted[2:], body)
		return &media.Packet{
			Type:             media.PacketTypeAudio,
			Timestamp:        timestamp,
			Data:             converted,
			IsSequenceHeader: true,
			TrackID:          trackID,
			IsEnhanced:       true,
			FourCC:           fourCC,
		}, nil

	case audioExTypeCodedFrames:
		body := data[bodyOffset:]
		converted := make([]byte, 2+len(body))
		converted[0] = 0xAF
		converted[1] = 0x01
		copy(converted[2:], body)
		return &media.Packet{
			Type:       media.PacketTypeAudio,
			Timestamp:  timestamp,
			Data:       converted,
			TrackID:    trackID,
			IsEnhanced: true,
			FourCC:     fourCC,
		}, nil

	case audioExTypeSequenceEnd:
		return nil, nil

	default:
		return nil, nil
	}
}

func buildEnhancedAVCPacket(timestamp uint32, trackID uint8, fourCC string, cts []byte, payload []byte) *media.Packet {
	isKeyframe := isAVCCKeyframe(payload)
	firstByte := byte(0x27)
	if isKeyframe {
		firstByte = 0x17
	}

	converted := make([]byte, 5+len(payload))
	converted[0] = firstByte
	converted[1] = 0x01
	copy(converted[2:5], cts)
	copy(converted[5:], payload)

	return &media.Packet{
		Type:       media.PacketTypeVideo,
		Timestamp:  timestamp,
		Data:       converted,
		IsKeyframe: isKeyframe,
		TrackID:    trackID,
		IsEnhanced: true,
		FourCC:     fourCC,
	}
}

func isEnhancedVideoTag(data []byte) bool {
	return len(data) > 0 && (data[0]&0x80) != 0
}

func isEnhancedAudioTag(data []byte) bool {
	return len(data) > 0 && (data[0]>>4) == enhancedAudioFormat
}

func readFourCC(data []byte, offset int) string {
	if offset < 0 || offset+4 > len(data) {
		return ""
	}
	return string(data[offset : offset+4])
}

func isAVCCKeyframe(payload []byte) bool {
	pos := 0
	for pos+4 <= len(payload) {
		naluLen := int(binary.BigEndian.Uint32(payload[pos : pos+4]))
		pos += 4
		if naluLen <= 0 || pos+naluLen > len(payload) {
			return false
		}
		naluType := payload[pos] & 0x1F
		if naluType == 5 {
			return true
		}
		pos += naluLen
	}
	return false
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
