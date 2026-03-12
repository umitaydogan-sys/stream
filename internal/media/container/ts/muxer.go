package ts

import (
	"bytes"
	"encoding/binary"

	"github.com/fluxstream/fluxstream/internal/media"
)

// MPEG-TS constants
const (
	TSPacketSize = 188
	SyncByte     = byte(0x47)

	StreamTypeH264 = byte(0x1B)
	StreamTypeH265 = byte(0x24)
	StreamTypeAAC  = byte(0x0F)
	StreamTypeMP3  = byte(0x03)
)

// pidBytes returns high and low bytes of a uint16 PID
func pidBytes(pid uint16) (hi, lo byte) {
	return byte(pid >> 8), byte(pid)
}

// PID values as uint16
const (
	PatPID   uint16 = 0x0000
	PmtPID   uint16 = 0x1000
	VideoPID uint16 = 0x0100
	AudioPID uint16 = 0x0101
)

// Muxer converts media packets to MPEG-TS format
type Muxer struct {
	videoContinuity byte
	audioContinuity byte
	patContinuity   byte
	pmtContinuity   byte
	patPmtInterval  int
	packetCount     int
}

// NewMuxer creates a new MPEG-TS muxer
func NewMuxer() *Muxer {
	return &Muxer{
		patPmtInterval: 50, // send PAT/PMT every 50 packets
	}
}

// MuxPacket converts a media packet to MPEG-TS bytes
func (m *Muxer) MuxPacket(pkt *media.Packet) []byte {
	var buf bytes.Buffer

	// Periodically insert PAT/PMT
	if m.packetCount%m.patPmtInterval == 0 {
		buf.Write(m.buildPAT())
		buf.Write(m.buildPMT())
	}
	m.packetCount++

	// Build PES packet
	var pid uint16
	var continuity *byte

	switch pkt.Type {
	case media.PacketTypeVideo:
		pid = VideoPID
		continuity = &m.videoContinuity
	case media.PacketTypeAudio:
		pid = AudioPID
		continuity = &m.audioContinuity
	default:
		return nil
	}

	pesData := m.buildPES(pkt)
	tsPackets := m.packIntoTS(pid, continuity, pesData, pkt.IsKeyframe, pkt.Timestamp)
	buf.Write(tsPackets)

	return buf.Bytes()
}

// GeneratePatPmt returns PAT+PMT TS packets
func (m *Muxer) GeneratePatPmt() []byte {
	var buf bytes.Buffer
	buf.Write(m.buildPAT())
	buf.Write(m.buildPMT())
	return buf.Bytes()
}

func (m *Muxer) buildPAT() []byte {
	pkt := make([]byte, TSPacketSize)
	pkt[0] = SyncByte
	pkt[1] = 0x40                            // payload unit start
	pkt[2] = 0x00                            // PID = 0
	pkt[3] = 0x10 | (m.patContinuity & 0x0F) // payload only
	m.patContinuity++

	// Pointer field
	pkt[4] = 0x00

	// PAT table
	pmtHi, pmtLo := pidBytes(PmtPID)
	pat := []byte{
		0x00,       // table_id
		0xB0, 0x0D, // section length = 13
		0x00, 0x01, // transport_stream_id
		0xC1,       // version, current
		0x00, 0x00, // section number, last section
		0x00, 0x01, // program_number = 1
		pmtHi | 0xE0, pmtLo, // PMT PID
	}

	// CRC32
	crc := crc32MPEG2(pat)
	pat = append(pat, byte(crc>>24), byte(crc>>16), byte(crc>>8), byte(crc))

	copy(pkt[5:], pat)

	// Fill rest with 0xFF
	for i := 5 + len(pat); i < TSPacketSize; i++ {
		pkt[i] = 0xFF
	}

	return pkt
}

func (m *Muxer) buildPMT() []byte {
	pkt := make([]byte, TSPacketSize)
	pmtHi, pmtLo := pidBytes(PmtPID)
	vidHi, vidLo := pidBytes(VideoPID)
	audHi, audLo := pidBytes(AudioPID)

	pkt[0] = SyncByte
	pkt[1] = pmtHi | 0x40 // payload unit start + PID high
	pkt[2] = pmtLo        // PID low
	pkt[3] = 0x10 | (m.pmtContinuity & 0x0F)
	m.pmtContinuity++

	pkt[4] = 0x00 // pointer

	// PMT table
	pmt := []byte{
		0x02,       // table_id
		0xB0, 0x17, // section length = 23
		0x00, 0x01, // program_number
		0xC1,       // version, current
		0x00, 0x00, // section number, last section
		vidHi | 0xE0, vidLo, // PCR PID = video
		0xF0, 0x00, // program info length = 0
		// Video stream
		StreamTypeH264,
		vidHi | 0xE0, vidLo,
		0xF0, 0x00, // ES info length = 0
		// Audio stream
		StreamTypeAAC,
		audHi | 0xE0, audLo,
		0xF0, 0x00, // ES info length = 0
	}

	crc := crc32MPEG2(pmt)
	pmt = append(pmt, byte(crc>>24), byte(crc>>16), byte(crc>>8), byte(crc))

	copy(pkt[5:], pmt)

	for i := 5 + len(pmt); i < TSPacketSize; i++ {
		pkt[i] = 0xFF
	}

	return pkt
}

func (m *Muxer) buildPES(pkt *media.Packet) []byte {
	var streamID byte
	switch pkt.Type {
	case media.PacketTypeVideo:
		streamID = 0xE0 // video
	case media.PacketTypeAudio:
		streamID = 0xC0 // audio
	default:
		return nil
	}

	// PTS in 90kHz clock
	pts := uint64(pkt.Timestamp) * 90 // ms to 90kHz

	// PES header
	pesHeader := make([]byte, 0, 19)
	pesHeader = append(pesHeader, 0x00, 0x00, 0x01) // start code
	pesHeader = append(pesHeader, streamID)

	// PES packet length (0 = unbounded for video, or calculated)
	dataLen := len(pkt.Data) + 8 // 8 = PES header extension
	if dataLen > 65535 {
		pesHeader = append(pesHeader, 0x00, 0x00) // unbounded
	} else {
		pesHeader = append(pesHeader, byte(dataLen>>8), byte(dataLen))
	}

	// PES header extension
	pesHeader = append(pesHeader, 0x80) // marker bits
	pesHeader = append(pesHeader, 0x80) // PTS present
	pesHeader = append(pesHeader, 0x05) // PES header data length

	// PTS
	pesHeader = append(pesHeader, encodePTS(pts)...)

	return append(pesHeader, pkt.Data...)
}

func encodePTS(pts uint64) []byte {
	b := make([]byte, 5)
	b[0] = byte(((pts>>30)&0x07)<<1) | 0x21
	v := uint16((pts >> 15) & 0x7FFF)
	b[1] = byte(v >> 7)
	b[2] = byte(v<<1) | 0x01
	v = uint16(pts & 0x7FFF)
	b[3] = byte(v >> 7)
	b[4] = byte(v<<1) | 0x01
	return b
}

func (m *Muxer) packIntoTS(pid uint16, continuity *byte, pesData []byte, isKeyframe bool, pcrMS uint32) []byte {
	var result bytes.Buffer
	first := true
	offset := 0

	for offset < len(pesData) {
		pkt := make([]byte, TSPacketSize)
		pkt[0] = SyncByte

		// PID
		pkt[1] = byte(pid >> 8)
		pkt[2] = byte(pid)

		if first {
			pkt[1] |= 0x40 // payload unit start indicator
		}

		headerSize := 4
		adaptationLen := 0

		// If first video packet and keyframe, add adaptation field with PCR
		if first && isKeyframe && pid == VideoPID {
			adaptationLen = 8                    // minimum for PCR
			pkt[3] = 0x30 | (*continuity & 0x0F) // adaptation + payload
			pkt[4] = byte(adaptationLen - 1)     // adaptation field length
			pkt[5] = 0x10                        // PCR flag
			pcr := encodePCR90k(pcrMS)
			copy(pkt[6:12], pcr[:])
			headerSize = 4 + adaptationLen
		} else {
			pkt[3] = 0x10 | (*continuity & 0x0F) // payload only
		}

		*continuity++

		payloadSize := TSPacketSize - headerSize
		remaining := len(pesData) - offset

		if remaining < payloadSize {
			// Need stuffing via adaptation field
			stuffLen := payloadSize - remaining
			if pkt[3]&0x20 != 0 {
				// Already have adaptation field, extend it
				currentAdaptLen := int(pkt[4])
				pkt[4] = byte(currentAdaptLen + stuffLen)
				for i := 0; i < stuffLen; i++ {
					pkt[headerSize+i] = 0xFF
				}
				headerSize += stuffLen
			} else {
				// Add adaptation field for stuffing
				pkt[3] |= 0x20 // set adaptation flag
				if stuffLen == 1 {
					// shift everything
					pkt[4] = 0x00 // adaptation length = 0
					headerSize = 5
				} else {
					pkt[4] = byte(stuffLen - 1)
					pkt[5] = 0x00 // flags
					for i := 6; i < 4+stuffLen; i++ {
						pkt[i] = 0xFF
					}
					headerSize = 4 + stuffLen
				}
			}
			payloadSize = remaining
		}

		copy(pkt[headerSize:], pesData[offset:offset+payloadSize])
		offset += payloadSize
		first = false

		result.Write(pkt)
	}

	return result.Bytes()
}

func encodePCR90k(ms uint32) [6]byte {
	// PCR base uses 90kHz clock. Keep extension at 0.
	base := uint64(ms) * 90
	var b [6]byte
	b[0] = byte(base >> 25)
	b[1] = byte(base >> 17)
	b[2] = byte(base >> 9)
	b[3] = byte(base >> 1)
	b[4] = byte((base&0x01)<<7) | 0x7E
	b[5] = 0x00
	return b
}

// CRC32 for MPEG-2 (polynomial 0x04C11DB7)
func crc32MPEG2(data []byte) uint32 {
	crc := uint32(0xFFFFFFFF)
	for _, b := range data {
		for i := 0; i < 8; i++ {
			if (crc>>31)^(uint32(b>>(7-uint(i)))&1) != 0 {
				crc = (crc << 1) ^ 0x04C11DB7
			} else {
				crc = crc << 1
			}
		}
	}
	return crc
}

// MPEG-TS Timestamp helpers
func TSTimestamp(ms uint32) uint64 {
	return uint64(ms) * 90
}

func WritePCR(buf []byte, pcr uint64) {
	base := pcr
	ext := uint64(0)
	binary.BigEndian.PutUint32(buf[0:4], uint32(base>>1))
	buf[4] = byte((base&1)<<7) | 0x7E | byte(ext>>8&1)
	buf[5] = byte(ext)
}
