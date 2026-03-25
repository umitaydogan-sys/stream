package recording

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func (m *Manager) finalizeRecording(rec *Recording) error {
	if rec == nil {
		return nil
	}
	if rec.capturePath == "" || rec.capturePath == rec.finalPath {
		rec.FilePath = rec.finalPath
		return nil
	}
	if rec.Format != FormatMP4 && rec.Format != FormatMKV {
		rec.FilePath = rec.capturePath
		return nil
	}
	if err := m.remuxFile(rec.capturePath, rec.finalPath, rec.Format); err != nil {
		fallbackPath := strings.TrimSuffix(rec.finalPath, filepath.Ext(rec.finalPath)) + ".ts"
		_ = os.Remove(fallbackPath)
		if renameErr := os.Rename(rec.capturePath, fallbackPath); renameErr == nil {
			rec.FilePath = fallbackPath
		} else {
			rec.FilePath = rec.capturePath
		}
		return err
	}
	_ = os.Remove(rec.capturePath)
	rec.FilePath = rec.finalPath
	return nil
}

func (m *Manager) RemuxSavedRecording(streamKey, filename string, targetFormat Format) (*SavedRecording, error) {
	targetFormat = normalizeRemuxTarget(targetFormat)
	streamKey = strings.TrimSpace(streamKey)
	filename = filepath.Base(strings.TrimSpace(filename))
	if streamKey == "" || filename == "" {
		return nil, fmt.Errorf("stream_key ve filename gerekli")
	}
	sourcePath := m.RecordingFilePath(streamKey, filename)
	if _, err := os.Stat(sourcePath); err != nil {
		return nil, err
	}
	targetName := strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + string(targetFormat)
	targetPath := m.RecordingFilePath(streamKey, targetName)
	if strings.EqualFold(sourcePath, targetPath) {
		info, err := os.Stat(sourcePath)
		if err != nil {
			return nil, err
		}
		return &SavedRecording{
			StreamKey: streamKey,
			Name:      targetName,
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			Format:    strings.TrimPrefix(strings.ToLower(filepath.Ext(targetName)), "."),
			Path:      targetPath,
		}, nil
	}
	if err := m.remuxFile(sourcePath, targetPath, targetFormat); err != nil {
		return nil, err
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		return nil, err
	}
	return &SavedRecording{
		StreamKey: streamKey,
		Name:      targetName,
		Size:      info.Size(),
		ModTime:   info.ModTime(),
		Format:    strings.TrimPrefix(strings.ToLower(filepath.Ext(targetName)), "."),
		Path:      targetPath,
	}, nil
}

func (m *Manager) remuxFile(sourcePath, targetPath string, format Format) error {
	if strings.TrimSpace(m.ffmpegPath) == "" {
		m.ffmpegPath = "ffmpeg"
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}
	tmpPath := targetPath + ".tmp"
	_ = os.Remove(tmpPath)
	args := []string{"-hide_banner", "-loglevel", "error", "-y", "-fflags", "+genpts", "-i", sourcePath, "-map", "0", "-c", "copy"}
	if format == FormatMP4 {
		args = append(args, "-movflags", "+faststart")
	}
	args = append(args, tmpPath)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, m.ffmpegPath, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.Remove(tmpPath)
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("ffmpeg remux hatasi: %s", msg)
	}
	_ = os.Remove(targetPath)
	if err := os.Rename(tmpPath, targetPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}

func normalizeRemuxTarget(format Format) Format {
	switch normalizeFormat(format) {
	case FormatMKV:
		return FormatMKV
	case FormatMP4:
		fallthrough
	default:
		return FormatMP4
	}
}
