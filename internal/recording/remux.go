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
	targetExt := strings.ToLower(filepath.Ext(targetPath))
	tmpPath := strings.TrimSuffix(targetPath, targetExt) + ".partial" + targetExt
	_ = os.Remove(tmpPath)
	runFFmpeg := func(timeout time.Duration, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		cmd := exec.CommandContext(ctx, m.ffmpegPath, args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			msg := strings.TrimSpace(string(out))
			if msg == "" {
				msg = err.Error()
			}
			return fmt.Errorf("ffmpeg remux hatasi: %s", msg)
		}
		return nil
	}
	baseArgs := []string{
		"-hide_banner", "-loglevel", "error", "-y",
		"-fflags", "+genpts+discardcorrupt",
		"-err_detect", "ignore_err",
		"-analyzeduration", "100M",
		"-probesize", "100M",
		"-i", sourcePath,
		"-map", "0:v?",
		"-map", "0:a?",
		"-map", "0:s?",
	}
	copyArgs := append([]string{}, baseArgs...)
	copyArgs = append(copyArgs, "-c", "copy")
	switch format {
	case FormatMP4:
		copyArgs = append(copyArgs, "-movflags", "+faststart", "-f", "mp4")
	case FormatMKV:
		copyArgs = append(copyArgs, "-f", "matroska")
	}
	copyArgs = append(copyArgs, tmpPath)
	copyErr := runFFmpeg(10*time.Minute, copyArgs)
	if copyErr != nil && format == FormatMP4 {
		_ = os.Remove(tmpPath)
		fallbackArgs := append([]string{}, baseArgs...)
		fallbackArgs = append(fallbackArgs,
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-crf", "23",
			"-c:a", "aac",
			"-b:a", "160k",
			"-movflags", "+faststart",
			"-f", "mp4",
			tmpPath,
		)
		if fallbackErr := runFFmpeg(30*time.Minute, fallbackArgs); fallbackErr != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("%v | fallback: %v", copyErr, fallbackErr)
		}
	} else if copyErr != nil {
		_ = os.Remove(tmpPath)
		return copyErr
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

func (m *Manager) runRemuxJob(jobID, streamKey, filename string, targetFormat Format) {
	m.mu.Lock()
	job, ok := m.remuxJobs[jobID]
	if ok {
		job.Status = "running"
	}
	m.mu.Unlock()
	if !ok {
		return
	}

	saved, err := m.RemuxSavedRecording(streamKey, filename, targetFormat)

	m.mu.Lock()
	defer m.mu.Unlock()
	job, ok = m.remuxJobs[jobID]
	if !ok {
		return
	}
	job.FinishedAt = time.Now()
	if err != nil {
		job.Status = "error"
		job.LastError = err.Error()
		return
	}
	job.Status = "completed"
	if saved != nil && strings.TrimSpace(saved.Name) != "" {
		job.TargetName = saved.Name
	}
}
