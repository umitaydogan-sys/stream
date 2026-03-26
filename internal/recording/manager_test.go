package recording

import (
	"bytes"
	"testing"

	"github.com/fluxstream/fluxstream/internal/media"
	ts "github.com/fluxstream/fluxstream/internal/media/container/ts"
)

func TestEncodeTSPacketWaitsForKeyframeAndWritesAnnexBAndADTS(t *testing.T) {
	rec := &Recording{
		Format:        FormatTS,
		tsMuxer:       ts.NewMuxer(),
		aacProfile:    1,
		aacFreqIndex:  4,
		aacChannelCfg: 2,
	}
	mgr := &Manager{}

	videoSeq := &media.Packet{
		Type:             media.PacketTypeVideo,
		IsSequenceHeader: true,
		Data: []byte{
			0x17, 0x00, 0x00, 0x00, 0x00,
			0x01, 0x64, 0x00, 0x1E, 0xFF, 0xE1,
			0x00, 0x05, 0x67, 0x64, 0x00, 0x1E, 0xAC,
			0x01, 0x00, 0x04, 0x68, 0xEE, 0x3C, 0x80,
		},
	}
	audioSeq := &media.Packet{
		Type:             media.PacketTypeAudio,
		IsSequenceHeader: true,
		Data:             []byte{0xAF, 0x00, 0x12, 0x10},
	}
	if out := mgr.encodeTSPacket(rec, videoSeq); out != nil {
		t.Fatalf("sequence header should not be muxed")
	}
	if out := mgr.encodeTSPacket(rec, audioSeq); out != nil {
		t.Fatalf("sequence header should not be muxed")
	}

	audioBeforeKeyframe := &media.Packet{
		Type:      media.PacketTypeAudio,
		Timestamp: 20,
		Data:      []byte{0xAF, 0x01, 0x11, 0x22, 0x33},
	}
	if out := mgr.encodeTSPacket(rec, audioBeforeKeyframe); out != nil {
		t.Fatalf("audio should wait until first video keyframe when video is present")
	}

	videoKeyframe := &media.Packet{
		Type:       media.PacketTypeVideo,
		Timestamp:  40,
		IsKeyframe: true,
		Data: []byte{
			0x17, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x04, 0x65, 0x88, 0x84, 0x21,
		},
	}
	videoOut := mgr.encodeTSPacket(rec, videoKeyframe)
	if len(videoOut) == 0 {
		t.Fatalf("expected keyframe TS output")
	}
	if !bytes.Contains(videoOut, []byte{0x00, 0x00, 0x00, 0x01, 0x67}) {
		t.Fatalf("expected SPS annex-b config in TS payload")
	}
	if !bytes.Contains(videoOut, []byte{0x00, 0x00, 0x00, 0x01, 0x65}) {
		t.Fatalf("expected IDR annex-b payload in TS payload")
	}

	audioAfterKeyframe := &media.Packet{
		Type:      media.PacketTypeAudio,
		Timestamp: 60,
		Data:      []byte{0xAF, 0x01, 0x11, 0x22, 0x33},
	}
	audioOut := mgr.encodeTSPacket(rec, audioAfterKeyframe)
	if len(audioOut) == 0 {
		t.Fatalf("expected audio TS output after video start")
	}
	if !bytes.Contains(audioOut, []byte{0xFF, 0xF1}) {
		t.Fatalf("expected ADTS header in audio TS payload")
	}
}

func TestEncodeTSPacketAudioOnlyStartsImmediately(t *testing.T) {
	rec := &Recording{
		Format:        FormatTS,
		tsMuxer:       ts.NewMuxer(),
		aacProfile:    1,
		aacFreqIndex:  4,
		aacChannelCfg: 2,
	}
	mgr := &Manager{}
	audioSeq := &media.Packet{
		Type:             media.PacketTypeAudio,
		IsSequenceHeader: true,
		Data:             []byte{0xAF, 0x00, 0x12, 0x10},
	}
	mgr.encodeTSPacket(rec, audioSeq)
	audioPkt := &media.Packet{
		Type:      media.PacketTypeAudio,
		Timestamp: 20,
		Data:      []byte{0xAF, 0x01, 0x11, 0x22, 0x33},
	}
	out := mgr.encodeTSPacket(rec, audioPkt)
	if len(out) == 0 {
		t.Fatalf("expected audio-only TS output")
	}
	if !bytes.Contains(out, []byte{0xFF, 0xF1}) {
		t.Fatalf("expected ADTS header in audio-only TS payload")
	}
}
