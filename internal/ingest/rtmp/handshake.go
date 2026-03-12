package rtmp

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	HandshakeSize = 1536
)

// DoHandshake performs the RTMP handshake as server
func DoHandshake(rw io.ReadWriter) error {
	// ─── C0 ───
	c0 := make([]byte, 1)
	if _, err := io.ReadFull(rw, c0); err != nil {
		return fmt.Errorf("read C0: %w", err)
	}
	if c0[0] != 3 {
		return fmt.Errorf("unsupported RTMP version: %d", c0[0])
	}

	// ─── S0 ───
	if _, err := rw.Write([]byte{3}); err != nil {
		return fmt.Errorf("write S0: %w", err)
	}

	// ─── C1 ───
	c1 := make([]byte, HandshakeSize)
	if _, err := io.ReadFull(rw, c1); err != nil {
		return fmt.Errorf("read C1: %w", err)
	}

	// ─── S1 ───
	s1 := make([]byte, HandshakeSize)
	binary.BigEndian.PutUint32(s1[0:4], uint32(time.Now().Unix()))
	// Zero bytes at 4-7 (simple handshake)
	rand.Read(s1[8:])
	if _, err := rw.Write(s1); err != nil {
		return fmt.Errorf("write S1: %w", err)
	}

	// ─── S2 (echo of C1) ───
	s2 := make([]byte, HandshakeSize)
	copy(s2, c1)
	binary.BigEndian.PutUint32(s2[4:8], uint32(time.Now().Unix()))
	if _, err := rw.Write(s2); err != nil {
		return fmt.Errorf("write S2: %w", err)
	}

	// ─── C2 ───
	c2 := make([]byte, HandshakeSize)
	if _, err := io.ReadFull(rw, c2); err != nil {
		return fmt.Errorf("read C2: %w", err)
	}

	return nil
}
