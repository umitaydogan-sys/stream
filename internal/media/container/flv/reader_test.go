package flv

import "testing"

func TestReadTagEnhancedVideoMultitrackSequenceStart(t *testing.T) {
	reader := NewReader(nil)
	config := []byte{0x01, 0x64, 0x00, 0x1F}
	data := append([]byte{
		0x86,
		0x00,
		'a', 'v', 'c', '1',
		0x01,
	}, config...)

	pkt, err := reader.ReadTag(TagVideo, uint32(len(data)), 42, data)
	if err != nil {
		t.Fatalf("ReadTag() error = %v", err)
	}
	if pkt == nil {
		t.Fatal("ReadTag() returned nil packet")
	}
	if !pkt.IsSequenceHeader {
		t.Fatal("expected sequence header")
	}
	if !pkt.IsKeyframe {
		t.Fatal("expected keyframe sequence header")
	}
	if pkt.TrackID != 1 {
		t.Fatalf("expected track 1, got %d", pkt.TrackID)
	}
	if pkt.FourCC != "avc1" {
		t.Fatalf("expected avc1, got %q", pkt.FourCC)
	}
	if len(pkt.Data) != 5+len(config) {
		t.Fatalf("unexpected converted payload size: %d", len(pkt.Data))
	}
	if pkt.Data[0] != 0x17 || pkt.Data[1] != 0x00 {
		t.Fatalf("unexpected converted video header: %x %x", pkt.Data[0], pkt.Data[1])
	}
}

func TestReadTagEnhancedVideoMultitrackCodedFrames(t *testing.T) {
	reader := NewReader(nil)
	avcc := []byte{0x00, 0x00, 0x00, 0x01, 0x65}
	data := append([]byte{
		0x86,
		0x01,
		'a', 'v', 'c', '1',
		0x01,
		0x00, 0x00, 0x02,
	}, avcc...)

	pkt, err := reader.ReadTag(TagVideo, uint32(len(data)), 42, data)
	if err != nil {
		t.Fatalf("ReadTag() error = %v", err)
	}
	if pkt == nil {
		t.Fatal("ReadTag() returned nil packet")
	}
	if !pkt.IsKeyframe {
		t.Fatal("expected IDR frame to be marked as keyframe")
	}
	if pkt.TrackID != 1 {
		t.Fatalf("expected track 1, got %d", pkt.TrackID)
	}
	if pkt.Data[0] != 0x17 || pkt.Data[1] != 0x01 {
		t.Fatalf("unexpected converted video header: %x %x", pkt.Data[0], pkt.Data[1])
	}
	if pkt.Data[4] != 0x02 {
		t.Fatalf("unexpected CTS byte: %x", pkt.Data[4])
	}
}

func TestReadTagEnhancedAudioMultitrackAAC(t *testing.T) {
	reader := NewReader(nil)

	seq := []byte{
		0x95,
		0x00,
		'm', 'p', '4', 'a',
		0x01,
		0x12, 0x10,
	}
	pkt, err := reader.ReadTag(TagAudio, uint32(len(seq)), 7, seq)
	if err != nil {
		t.Fatalf("ReadTag() error = %v", err)
	}
	if pkt == nil {
		t.Fatal("ReadTag() returned nil packet")
	}
	if !pkt.IsSequenceHeader {
		t.Fatal("expected AAC sequence header")
	}
	if pkt.TrackID != 1 {
		t.Fatalf("expected track 1, got %d", pkt.TrackID)
	}
	if pkt.Data[0] != 0xAF || pkt.Data[1] != 0x00 {
		t.Fatalf("unexpected converted audio header: %x %x", pkt.Data[0], pkt.Data[1])
	}

	raw := []byte{
		0x95,
		0x01,
		'm', 'p', '4', 'a',
		0x01,
		0x11, 0x22, 0x33,
	}
	pkt, err = reader.ReadTag(TagAudio, uint32(len(raw)), 8, raw)
	if err != nil {
		t.Fatalf("ReadTag() error = %v", err)
	}
	if pkt == nil {
		t.Fatal("ReadTag() returned nil packet")
	}
	if pkt.IsSequenceHeader {
		t.Fatal("did not expect AAC raw frame to be sequence header")
	}
	if pkt.Data[0] != 0xAF || pkt.Data[1] != 0x01 {
		t.Fatalf("unexpected converted audio header: %x %x", pkt.Data[0], pkt.Data[1])
	}
}

func TestReadTagEnhancedMultitrackUnsupportedModeIsIgnored(t *testing.T) {
	reader := NewReader(nil)
	data := []byte{
		0x86,
		0x10,
		'a', 'v', 'c', '1',
		0x01,
		0x01, 0x02,
	}

	pkt, err := reader.ReadTag(TagVideo, uint32(len(data)), 1, data)
	if err != nil {
		t.Fatalf("ReadTag() error = %v", err)
	}
	if pkt != nil {
		t.Fatalf("expected packet to be ignored, got %+v", pkt)
	}
}
