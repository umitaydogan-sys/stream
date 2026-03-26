package dash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fluxstream/fluxstream/internal/media"
)

func TestAudioOnlyDASHProducesAudioManifest(t *testing.T) {
	root := t.TempDir()
	muxer := NewMuxer(root)
	sm := muxer.AddStream("audio_only")

	if err := sm.WritePacket(&media.Packet{
		Type:             media.PacketTypeAudio,
		IsSequenceHeader: true,
		Data:             []byte{0xAF, 0x00, 0x12, 0x10},
	}); err != nil {
		t.Fatalf("write audio sequence header: %v", err)
	}
	if err := sm.WritePacket(&media.Packet{
		Type:      media.PacketTypeAudio,
		Timestamp: 0,
		Data:      []byte{0xAF, 0x01, 0x11, 0x22, 0x33, 0x44},
	}); err != nil {
		t.Fatalf("write first audio packet: %v", err)
	}
	if err := sm.WritePacket(&media.Packet{
		Type:      media.PacketTypeAudio,
		Timestamp: 2200,
		Data:      []byte{0xAF, 0x01, 0x55, 0x66, 0x77, 0x88},
	}); err != nil {
		t.Fatalf("write second audio packet: %v", err)
	}
	sm.Close()

	streamDir := filepath.Join(root, "audio_only")
	manifest := mustReadFile(t, filepath.Join(streamDir, "manifest.mpd"))
	audioManifest := mustReadFile(t, filepath.Join(streamDir, "audio.mpd"))

	if strings.Contains(manifest, `mimeType="video/mp4"`) {
		t.Fatal("audio-only manifest unexpectedly contains a video adaptation set")
	}
	if !strings.Contains(manifest, `mimeType="audio/mp4"`) {
		t.Fatal("audio-only manifest is missing the audio adaptation set")
	}
	if !strings.Contains(audioManifest, `initialization="audio_init.mp4"`) {
		t.Fatal("audio-only dedicated manifest is missing audio_init.mp4")
	}
	if _, err := os.Stat(filepath.Join(streamDir, "audio_init.mp4")); err != nil {
		t.Fatalf("audio_init.mp4 missing: %v", err)
	}
}

func TestAVDASHAlsoProducesAudioOnlyManifest(t *testing.T) {
	root := t.TempDir()
	muxer := NewMuxer(root)
	sm := muxer.AddStream("av")

	if err := sm.WritePacket(&media.Packet{
		Type:             media.PacketTypeVideo,
		IsSequenceHeader: true,
		Data:             []byte{0x17, 0x00, 0x00, 0x00, 0x00, 0x01, 0x64},
	}); err != nil {
		t.Fatalf("write video sequence header: %v", err)
	}
	if err := sm.WritePacket(&media.Packet{
		Type:             media.PacketTypeAudio,
		IsSequenceHeader: true,
		Data:             []byte{0xAF, 0x00, 0x12, 0x10},
	}); err != nil {
		t.Fatalf("write audio sequence header: %v", err)
	}
	if err := sm.WritePacket(&media.Packet{
		Type:       media.PacketTypeVideo,
		Timestamp:  0,
		IsKeyframe: true,
		Data:       []byte{0x17, 0x01, 0x00, 0x00, 0x00, 0x65, 0x88, 0x99},
	}); err != nil {
		t.Fatalf("write first video packet: %v", err)
	}
	if err := sm.WritePacket(&media.Packet{
		Type:      media.PacketTypeAudio,
		Timestamp: 0,
		Data:      []byte{0xAF, 0x01, 0x11, 0x22, 0x33, 0x44},
	}); err != nil {
		t.Fatalf("write first audio packet: %v", err)
	}
	if err := sm.WritePacket(&media.Packet{
		Type:       media.PacketTypeVideo,
		Timestamp:  2200,
		IsKeyframe: true,
		Data:       []byte{0x17, 0x01, 0x00, 0x00, 0x00, 0x65, 0xAA, 0xBB},
	}); err != nil {
		t.Fatalf("write second video packet: %v", err)
	}
	sm.Close()

	streamDir := filepath.Join(root, "av")
	manifest := mustReadFile(t, filepath.Join(streamDir, "manifest.mpd"))
	audioManifest := mustReadFile(t, filepath.Join(streamDir, "audio.mpd"))

	if !strings.Contains(manifest, `mimeType="video/mp4"`) {
		t.Fatal("main manifest should keep the video adaptation set")
	}
	if !strings.Contains(manifest, `mimeType="audio/mp4"`) {
		t.Fatal("main manifest should keep the audio adaptation set")
	}
	if strings.Contains(audioManifest, `mimeType="video/mp4"`) {
		t.Fatal("audio-only manifest should not contain a video adaptation set")
	}
	if !strings.Contains(audioManifest, `mimeType="audio/mp4"`) {
		t.Fatal("audio-only manifest should contain the audio adaptation set")
	}
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
