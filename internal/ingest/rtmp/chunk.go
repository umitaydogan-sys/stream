package rtmp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// RTMP chunk stream constants
const (
	DefaultChunkSize = 128
	MaxChunkSize     = 65536

	// Message type IDs
	MsgSetChunkSize     = 1
	MsgAbort            = 2
	MsgAck              = 3
	MsgUserControl      = 4
	MsgAckSize          = 5
	MsgPeerBandwidth    = 6
	MsgAudioData        = 8
	MsgVideoData        = 9
	MsgAMF3Command      = 17
	MsgAMF0Command      = 20
	MsgAMF0Data         = 18
	MsgAMF3Data         = 15

	// Chunk stream IDs
	ChunkStreamProtocol = 2
	ChunkStreamCommand  = 3
	ChunkStreamAudio    = 4
	ChunkStreamVideo    = 6
)

// ChunkHeader represents an RTMP chunk header
type ChunkHeader struct {
	Fmt         byte   // 0-3
	CSID        uint32 // chunk stream ID
	Timestamp   uint32
	Length      uint32
	TypeID      byte
	StreamID    uint32
	ExtTimestamp bool
}

// Message represents a complete RTMP message reassembled from chunks
type Message struct {
	ChunkStreamID uint32
	Timestamp     uint32
	TypeID        byte
	StreamID      uint32
	Data          []byte
}

// ChunkReader reads RTMP chunks from a connection
type ChunkReader struct {
	r         io.Reader
	chunkSize uint32
	prevHeaders map[uint32]*ChunkHeader
	prevData    map[uint32][]byte
}

// NewChunkReader creates a new chunk reader
func NewChunkReader(r io.Reader) *ChunkReader {
	return &ChunkReader{
		r:           r,
		chunkSize:   DefaultChunkSize,
		prevHeaders: make(map[uint32]*ChunkHeader),
		prevData:    make(map[uint32][]byte),
	}
}

// SetChunkSize sets the incoming chunk size
func (cr *ChunkReader) SetChunkSize(size uint32) {
	cr.chunkSize = size
}

// ReadMessage reads a complete RTMP message (may span multiple chunks)
func (cr *ChunkReader) ReadMessage() (*Message, error) {
	for {
		header, err := cr.readChunkHeader()
		if err != nil {
			return nil, err
		}

		// Get previous header for this CSID
		prev := cr.prevHeaders[header.CSID]

		// Fill in missing fields based on format type
		switch header.Fmt {
		case 0:
			// Full header - all fields present
		case 1:
			// No stream ID
			if prev != nil {
				header.StreamID = prev.StreamID
			}
		case 2:
			// Only timestamp delta
			if prev != nil {
				header.Length = prev.Length
				header.TypeID = prev.TypeID
				header.StreamID = prev.StreamID
			}
		case 3:
			// No header
			if prev != nil {
				header.Timestamp = prev.Timestamp
				header.Length = prev.Length
				header.TypeID = prev.TypeID
				header.StreamID = prev.StreamID
			}
		}

		// Read chunk data
		buf := cr.prevData[header.CSID]
		remaining := int(header.Length) - len(buf)
		if remaining <= 0 {
			// Reset for new message
			buf = nil
			remaining = int(header.Length)
		}

		readSize := int(cr.chunkSize)
		if readSize > remaining {
			readSize = remaining
		}

		chunk := make([]byte, readSize)
		if _, err := io.ReadFull(cr.r, chunk); err != nil {
			return nil, fmt.Errorf("read chunk data: %w", err)
		}

		buf = append(buf, chunk...)
		cr.prevHeaders[header.CSID] = header

		if len(buf) >= int(header.Length) {
			// Complete message
			msg := &Message{
				ChunkStreamID: header.CSID,
				Timestamp:     header.Timestamp,
				TypeID:        header.TypeID,
				StreamID:      header.StreamID,
				Data:          buf[:header.Length],
			}
			delete(cr.prevData, header.CSID)
			return msg, nil
		}

		// Incomplete, store and continue
		cr.prevData[header.CSID] = buf
	}
}

func (cr *ChunkReader) readChunkHeader() (*ChunkHeader, error) {
	// Read basic header (1-3 bytes)
	b := make([]byte, 1)
	if _, err := io.ReadFull(cr.r, b); err != nil {
		return nil, err
	}

	header := &ChunkHeader{
		Fmt: (b[0] >> 6) & 0x03,
	}

	csid := uint32(b[0] & 0x3F)
	switch csid {
	case 0:
		// 2-byte form
		b2 := make([]byte, 1)
		if _, err := io.ReadFull(cr.r, b2); err != nil {
			return nil, err
		}
		header.CSID = uint32(b2[0]) + 64
	case 1:
		// 3-byte form
		b2 := make([]byte, 2)
		if _, err := io.ReadFull(cr.r, b2); err != nil {
			return nil, err
		}
		header.CSID = uint32(b2[1])*256 + uint32(b2[0]) + 64
	default:
		header.CSID = csid
	}

	// Read message header based on fmt
	switch header.Fmt {
	case 0:
		// Type 0: 11 bytes
		mh := make([]byte, 11)
		if _, err := io.ReadFull(cr.r, mh); err != nil {
			return nil, err
		}
		header.Timestamp = uint32(mh[0])<<16 | uint32(mh[1])<<8 | uint32(mh[2])
		header.Length = uint32(mh[3])<<16 | uint32(mh[4])<<8 | uint32(mh[5])
		header.TypeID = mh[6]
		header.StreamID = binary.LittleEndian.Uint32(mh[7:11])

		if header.Timestamp == 0xFFFFFF {
			header.ExtTimestamp = true
		}

	case 1:
		// Type 1: 7 bytes
		mh := make([]byte, 7)
		if _, err := io.ReadFull(cr.r, mh); err != nil {
			return nil, err
		}
		header.Timestamp = uint32(mh[0])<<16 | uint32(mh[1])<<8 | uint32(mh[2])
		header.Length = uint32(mh[3])<<16 | uint32(mh[4])<<8 | uint32(mh[5])
		header.TypeID = mh[6]

		if header.Timestamp == 0xFFFFFF {
			header.ExtTimestamp = true
		}

	case 2:
		// Type 2: 3 bytes
		mh := make([]byte, 3)
		if _, err := io.ReadFull(cr.r, mh); err != nil {
			return nil, err
		}
		header.Timestamp = uint32(mh[0])<<16 | uint32(mh[1])<<8 | uint32(mh[2])

		if header.Timestamp == 0xFFFFFF {
			header.ExtTimestamp = true
		}

	case 3:
		// Type 3: 0 bytes
		prev := cr.prevHeaders[header.CSID]
		if prev != nil && prev.ExtTimestamp {
			header.ExtTimestamp = true
		}
	}

	// Read extended timestamp if needed
	if header.ExtTimestamp {
		ext := make([]byte, 4)
		if _, err := io.ReadFull(cr.r, ext); err != nil {
			return nil, err
		}
		header.Timestamp = binary.BigEndian.Uint32(ext)
	}

	return header, nil
}

// ChunkWriter writes RTMP chunks to a connection
type ChunkWriter struct {
	w         io.Writer
	chunkSize uint32
}

// NewChunkWriter creates a new chunk writer
func NewChunkWriter(w io.Writer) *ChunkWriter {
	return &ChunkWriter{
		w:         w,
		chunkSize: DefaultChunkSize,
	}
}

// SetChunkSize sets the outgoing chunk size
func (cw *ChunkWriter) SetChunkSize(size uint32) {
	cw.chunkSize = size
}

// WriteMessage writes a complete RTMP message as chunks
func (cw *ChunkWriter) WriteMessage(csid uint32, timestamp uint32, typeID byte, streamID uint32, data []byte) error {
	// First chunk: Type 0 header
	header := cw.buildChunkHeader(0, csid, timestamp, uint32(len(data)), typeID, streamID)
	if _, err := cw.w.Write(header); err != nil {
		return err
	}

	offset := 0
	first := true
	for offset < len(data) {
		if !first {
			// Continuation chunks: Type 3 header
			contHeader := cw.buildChunkHeader(3, csid, 0, 0, 0, 0)
			if _, err := cw.w.Write(contHeader); err != nil {
				return err
			}
		}

		end := offset + int(cw.chunkSize)
		if end > len(data) {
			end = len(data)
		}

		if _, err := cw.w.Write(data[offset:end]); err != nil {
			return err
		}

		offset = end
		first = false
	}

	return nil
}

func (cw *ChunkWriter) buildChunkHeader(fmt byte, csid uint32, timestamp uint32, length uint32, typeID byte, streamID uint32) []byte {
	var header []byte

	// Basic header
	if csid < 64 {
		header = append(header, (fmt<<6)|byte(csid))
	} else if csid < 320 {
		header = append(header, fmt<<6, byte(csid-64))
	} else {
		header = append(header, (fmt<<6)|1)
		v := csid - 64
		header = append(header, byte(v), byte(v>>8))
	}

	// Message header
	switch fmt {
	case 0:
		ts := timestamp
		if ts >= 0xFFFFFF {
			ts = 0xFFFFFF
		}
		header = append(header,
			byte(ts>>16), byte(ts>>8), byte(ts),
			byte(length>>16), byte(length>>8), byte(length),
			typeID,
		)
		// Stream ID (little-endian)
		sid := make([]byte, 4)
		binary.LittleEndian.PutUint32(sid, streamID)
		header = append(header, sid...)

		if timestamp >= 0xFFFFFF {
			ext := make([]byte, 4)
			binary.BigEndian.PutUint32(ext, timestamp)
			header = append(header, ext...)
		}

	case 1:
		ts := timestamp
		if ts >= 0xFFFFFF {
			ts = 0xFFFFFF
		}
		header = append(header,
			byte(ts>>16), byte(ts>>8), byte(ts),
			byte(length>>16), byte(length>>8), byte(length),
			typeID,
		)
		if timestamp >= 0xFFFFFF {
			ext := make([]byte, 4)
			binary.BigEndian.PutUint32(ext, timestamp)
			header = append(header, ext...)
		}

	case 2:
		ts := timestamp
		if ts >= 0xFFFFFF {
			ts = 0xFFFFFF
		}
		header = append(header, byte(ts>>16), byte(ts>>8), byte(ts))
		if timestamp >= 0xFFFFFF {
			ext := make([]byte, 4)
			binary.BigEndian.PutUint32(ext, timestamp)
			header = append(header, ext...)
		}

	case 3:
		// No additional header
	}

	return header
}
