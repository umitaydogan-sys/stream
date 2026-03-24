package media

import "time"

// PacketType identifies the type of media packet
type PacketType byte

const (
	PacketTypeAudio PacketType = 0x08
	PacketTypeVideo PacketType = 0x09
	PacketTypeMeta  PacketType = 0x12
)

// Packet represents a generic media packet flowing through FluxStream
type Packet struct {
	Type             PacketType
	Timestamp        uint32 // milliseconds
	Data             []byte
	IsKeyframe       bool
	IsSequenceHeader bool // codec config (SPS/PPS for H.264, AudioSpecificConfig for AAC)
	TrackID          uint8
	IsEnhanced       bool
	FourCC           string
	StreamKey        string
	ReceivedAt       time.Time
}

// Clone creates a deep copy of the packet
func (p *Packet) Clone() *Packet {
	data := make([]byte, len(p.Data))
	copy(data, p.Data)
	return &Packet{
		Type:             p.Type,
		Timestamp:        p.Timestamp,
		Data:             data,
		IsKeyframe:       p.IsKeyframe,
		IsSequenceHeader: p.IsSequenceHeader,
		TrackID:          p.TrackID,
		IsEnhanced:       p.IsEnhanced,
		FourCC:           p.FourCC,
		StreamKey:        p.StreamKey,
		ReceivedAt:       p.ReceivedAt,
	}
}

// VideoCodecID identifies video codec from FLV tag
type VideoCodecID byte

const (
	VideoCodecH264 VideoCodecID = 7
	VideoCodecH265 VideoCodecID = 12
	VideoCodecVP8  VideoCodecID = 13
	VideoCodecVP9  VideoCodecID = 14
	VideoCodecAV1  VideoCodecID = 15
)

// AudioCodecID identifies audio codec from FLV tag
type AudioCodecID byte

const (
	AudioCodecAAC  AudioCodecID = 10
	AudioCodecMP3  AudioCodecID = 2
	AudioCodecOpus AudioCodecID = 13
)

// CodecString returns human-readable codec name
func (v VideoCodecID) String() string {
	switch v {
	case VideoCodecH264:
		return "H.264"
	case VideoCodecH265:
		return "H.265"
	case VideoCodecVP8:
		return "VP8"
	case VideoCodecVP9:
		return "VP9"
	case VideoCodecAV1:
		return "AV1"
	default:
		return "unknown"
	}
}

func (a AudioCodecID) String() string {
	switch a {
	case AudioCodecAAC:
		return "AAC"
	case AudioCodecMP3:
		return "MP3"
	case AudioCodecOpus:
		return "Opus"
	default:
		return "unknown"
	}
}
