package recording

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
	ts "github.com/fluxstream/fluxstream/internal/media/container/ts"
	"github.com/fluxstream/fluxstream/internal/stream"
)

// Format represents the recording format
type Format string

const (
	FormatTS  Format = "ts"
	FormatMP4 Format = "mp4"
	FormatMKV Format = "mkv"
	FormatFLV Format = "flv"
)

// Recording represents an active recording
type Recording struct {
	ID          string
	StreamKey   string
	Format      Format
	FilePath    string
	StartedAt   time.Time
	Duration    time.Duration
	Size        int64
	Status      string // "recording", "completed", "error"
	file        *os.File
	tsMuxer     *ts.Muxer
	subID       string
	stopCh      chan struct{}
	mu          sync.Mutex
	capturePath string
	finalPath   string
	finalizeErr string
	hasVideo    bool
	hasAudio    bool
	started     bool

	videoConfigNALU []byte
	aacProfile      int
	aacFreqIndex    int
	aacChannelCfg   int
}

// SavedRecording represents a completed recording file on disk.
type SavedRecording struct {
	StreamKey string    `json:"stream_key"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	Format    string    `json:"format"`
	Path      string    `json:"-"`
}

// RemuxJob represents a background conversion job for a saved recording.
type RemuxJob struct {
	ID           string    `json:"id"`
	StreamKey    string    `json:"stream_key"`
	SourceName   string    `json:"source_name"`
	TargetName   string    `json:"target_name"`
	TargetFormat string    `json:"target_format"`
	Status       string    `json:"status"`
	LastError    string    `json:"last_error,omitempty"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at,omitempty"`
}

// Manager manages stream recordings and DVR
type Manager struct {
	streamMgr     *stream.Manager
	recordingsDir string
	ffmpegPath    string
	recordings    map[string]*Recording
	remuxJobs     map[string]*RemuxJob
	maxDuration   time.Duration
	mu            sync.RWMutex
}

// NewManager creates a new recording manager
func NewManager(streamMgr *stream.Manager, recordingsDir, ffmpegPath string) *Manager {
	os.MkdirAll(recordingsDir, 0755)
	if strings.TrimSpace(ffmpegPath) == "" {
		ffmpegPath = "ffmpeg"
	}
	return &Manager{
		streamMgr:     streamMgr,
		recordingsDir: recordingsDir,
		ffmpegPath:    ffmpegPath,
		recordings:    make(map[string]*Recording),
		remuxJobs:     make(map[string]*RemuxJob),
		maxDuration:   24 * time.Hour,
	}
}

// StartRecording begins recording a live stream
func (m *Manager) StartRecording(streamKey string, format Format) (*Recording, error) {
	if !m.streamMgr.IsLive(streamKey) {
		return nil, fmt.Errorf("stream not live: %s", streamKey)
	}

	format = normalizeFormat(format)

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for existing recording
	recID := fmt.Sprintf("%s_%d", streamKey, time.Now().Unix())
	if _, exists := m.recordings[recID]; exists {
		return nil, fmt.Errorf("recording already exists: %s", recID)
	}

	// Create recording directory
	streamDir := filepath.Join(m.recordingsDir, streamKey)
	os.MkdirAll(streamDir, 0755)

	// Build filename
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.%s", streamKey, timestamp, format)
	filePath := filepath.Join(streamDir, filename)
	capturePath := filePath
	if format == FormatMP4 || format == FormatMKV {
		capturePath = filepath.Join(streamDir, fmt.Sprintf("%s_%s.capture.ts", streamKey, timestamp))
	}

	file, err := os.Create(capturePath)
	if err != nil {
		return nil, fmt.Errorf("create recording file: %w", err)
	}

	if format == FormatFLV {
		if _, err := file.Write([]byte{'F', 'L', 'V', 0x01, 0x05, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00}); err != nil {
			file.Close()
			return nil, fmt.Errorf("write flv header: %w", err)
		}
	}

	rec := &Recording{
		ID:            recID,
		StreamKey:     streamKey,
		Format:        format,
		FilePath:      filePath,
		StartedAt:     time.Now(),
		Status:        "recording",
		file:          file,
		tsMuxer:       ts.NewMuxer(),
		stopCh:        make(chan struct{}),
		capturePath:   capturePath,
		finalPath:     filePath,
		aacProfile:    1,
		aacFreqIndex:  4,
		aacChannelCfg: 2,
	}

	m.recordings[recID] = rec

	// Subscribe and start recording goroutine
	go m.recordLoop(rec)

	log.Printf("[REC] Kayıt başlatıldı: %s -> %s", streamKey, filePath)
	return rec, nil
}

// StartManagedRecording starts a recording and returns only the recording ID.
func (m *Manager) StartManagedRecording(streamKey, format string) (string, error) {
	rec, err := m.StartRecording(streamKey, Format(format))
	if err != nil {
		return "", err
	}
	return rec.ID, nil
}

// StopRecording stops an active recording
func (m *Manager) StopRecording(recID string) error {
	m.mu.RLock()
	rec, exists := m.recordings[recID]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("recording not found: %s", recID)
	}

	select {
	case <-rec.stopCh:
		// already stopped
	default:
		close(rec.stopCh)
	}
	return nil
}

// GetRecording returns a recording by ID
func (m *Manager) GetRecording(recID string) *Recording {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.recordings[recID]
}

// GetRecordings returns all recordings
func (m *Manager) GetRecordings() []*Recording {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Recording
	for _, r := range m.recordings {
		result = append(result, r)
	}
	return result
}

// GetStreamRecordings returns recordings for a specific stream
func (m *Manager) GetStreamRecordings(streamKey string) []*Recording {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Recording
	for _, r := range m.recordings {
		if r.StreamKey == streamKey {
			result = append(result, r)
		}
	}
	return result
}

// StartRemuxJob starts a background remux operation for a saved recording.
func (m *Manager) StartRemuxJob(streamKey, filename string, targetFormat Format) (*RemuxJob, error) {
	targetFormat = normalizeRemuxTarget(targetFormat)
	streamKey = strings.TrimSpace(streamKey)
	filename = filepath.Base(strings.TrimSpace(filename))
	if streamKey == "" || filename == "" {
		return nil, fmt.Errorf("stream_key ve filename gerekli")
	}
	targetName := strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + string(targetFormat)
	m.mu.Lock()
	for _, job := range m.remuxJobs {
		if job.StreamKey == streamKey && job.SourceName == filename && job.TargetName == targetName && (job.Status == "queued" || job.Status == "running") {
			existing := *job
			m.mu.Unlock()
			return &existing, nil
		}
	}
	jobID := fmt.Sprintf("remux_%d", time.Now().UnixNano())
	job := &RemuxJob{
		ID:           jobID,
		StreamKey:    streamKey,
		SourceName:   filename,
		TargetName:   targetName,
		TargetFormat: string(targetFormat),
		Status:       "queued",
		StartedAt:    time.Now(),
	}
	m.remuxJobs[jobID] = job
	m.mu.Unlock()

	go m.runRemuxJob(jobID, streamKey, filename, targetFormat)

	copyJob := *job
	return &copyJob, nil
}

// ListRemuxJobs returns all background remux jobs.
func (m *Manager) ListRemuxJobs() []*RemuxJob {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]*RemuxJob, 0, len(m.remuxJobs))
	for _, job := range m.remuxJobs {
		copyJob := *job
		result = append(result, &copyJob)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].StartedAt.After(result[j].StartedAt)
	})
	return result
}

// ListRecordingFiles returns saved recording files
func (m *Manager) ListRecordingFiles(streamKey string) ([]RecordingFile, error) {
	dir := filepath.Join(m.recordingsDir, streamKey)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []RecordingFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(e.Name()), ".capture.ts") || strings.HasSuffix(strings.ToLower(e.Name()), ".tmp") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, RecordingFile{
			StreamKey: streamKey,
			Name:      e.Name(),
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			Path:      filepath.Join(dir, e.Name()),
			Format:    strings.TrimPrefix(strings.ToLower(filepath.Ext(e.Name())), "."),
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})
	return files, nil
}

// RecordingFile represents a saved recording file
type RecordingFile struct {
	StreamKey string
	Name      string
	Size      int64
	ModTime   time.Time
	Path      string
	Format    string
}

// ListAllRecordingFiles returns every saved recording grouped across stream folders.
func (m *Manager) ListAllRecordingFiles() ([]SavedRecording, error) {
	entries, err := os.ReadDir(m.recordingsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []SavedRecording
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		streamKey := entry.Name()
		streamFiles, err := m.ListRecordingFiles(streamKey)
		if err != nil {
			return nil, err
		}
		for _, file := range streamFiles {
			files = append(files, SavedRecording{
				StreamKey: streamKey,
				Name:      file.Name,
				Size:      file.Size,
				ModTime:   file.ModTime,
				Format:    file.Format,
				Path:      file.Path,
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})
	return files, nil
}

// RecordingsDir returns the root directory for saved recordings.
func (m *Manager) RecordingsDir() string {
	return m.recordingsDir
}

// RecordingFilePath returns the absolute path for a saved recording file.
func (m *Manager) RecordingFilePath(streamKey, filename string) string {
	return filepath.Join(m.recordingsDir, streamKey, filepath.Base(filename))
}

func (m *Manager) recordLoop(rec *Recording) {
	rec.subID = fmt.Sprintf("rec_%s", rec.ID)
	sub := m.streamMgr.Subscribe(rec.StreamKey, rec.subID, 512)
	if sub == nil {
		rec.Status = "error"
		rec.file.Close()
		return
	}
	defer m.streamMgr.Unsubscribe(rec.StreamKey, rec.subID)
	defer func() {
		rec.mu.Lock()
		rec.Duration = time.Since(rec.StartedAt)
		rec.file.Close()
		rec.mu.Unlock()
		if rec.Status != "error" {
			if err := m.finalizeRecording(rec); err != nil {
				rec.mu.Lock()
				rec.finalizeErr = err.Error()
				rec.Status = "error"
				rec.mu.Unlock()
				log.Printf("[REC] Kayit finalize edilemedi: %s (%v)", rec.ID, err)
			} else {
				rec.mu.Lock()
				rec.Status = "completed"
				rec.mu.Unlock()
			}
		}
		m.mu.Lock()
		delete(m.recordings, rec.ID)
		m.mu.Unlock()
		log.Printf("[REC] Kayıt tamamlandı: %s (süre: %s, boyut: %s)",
			rec.ID, rec.Duration.Round(time.Second), formatBytes(rec.Size))
	}()

	maxTimer := time.NewTimer(m.maxDuration)
	defer maxTimer.Stop()

	for {
		select {
		case pkt, ok := <-sub.PacketC:
			if !ok {
				return
			}
			data := m.encodePacket(rec, pkt)
			if data == nil {
				continue
			}
			n, err := rec.file.Write(data)
			if err != nil {
				rec.mu.Lock()
				rec.Status = "error"
				rec.mu.Unlock()
				return
			}
			rec.mu.Lock()
			rec.Size += int64(n)
			rec.mu.Unlock()

		case <-rec.stopCh:
			return
		case <-sub.Done:
			return
		case <-maxTimer.C:
			log.Printf("[REC] Maksimum süre aşıldı: %s", rec.ID)
			return
		}
	}
}

func (m *Manager) encodePacket(rec *Recording, pkt *media.Packet) []byte {
	switch rec.Format {
	case FormatTS:
		return m.encodeTSPacket(rec, pkt)
	case FormatFLV:
		return m.encodeFLVPacket(pkt)
	default:
		return m.encodeTSPacket(rec, pkt)
	}
}

func (m *Manager) encodeTSPacket(rec *Recording, pkt *media.Packet) []byte {
	if pkt == nil {
		return nil
	}
	if pkt.IsSequenceHeader {
		switch pkt.Type {
		case media.PacketTypeVideo:
			rec.hasVideo = true
			rec.videoConfigNALU = recordingParseAVCConfigToAnnexB(pkt.Data)
		case media.PacketTypeAudio:
			rec.hasAudio = true
			recordingParseAACAudioSpecificConfig(rec, pkt.Data)
		}
		return nil
	}
	mediaPkt := pkt.Clone()
	switch pkt.Type {
	case media.PacketTypeVideo:
		rec.hasVideo = true
		if len(pkt.Data) <= 5 {
			return nil
		}
		if !rec.started {
			if !pkt.IsKeyframe {
				return nil
			}
			rec.started = true
		}
		annexB := recordingAVCCToAnnexB(pkt.Data[5:])
		if pkt.IsKeyframe && len(rec.videoConfigNALU) > 0 {
			mediaPkt.Data = append(append([]byte{}, rec.videoConfigNALU...), annexB...)
		} else {
			mediaPkt.Data = annexB
		}
	case media.PacketTypeAudio:
		rec.hasAudio = true
		if rec.hasVideo && !rec.started {
			return nil
		}
		if len(pkt.Data) < 2 {
			return nil
		}
		codecID := (pkt.Data[0] >> 4) & 0x0F
		if codecID == byte(media.AudioCodecAAC) {
			if len(pkt.Data) <= 2 {
				return nil
			}
			mediaPkt.Data = recordingAddADTSHeader(pkt.Data[2:], rec.aacProfile, rec.aacFreqIndex, rec.aacChannelCfg)
		} else {
			mediaPkt.Data = pkt.Data[1:]
		}
	default:
		return nil
	}
	return rec.tsMuxer.MuxPacket(mediaPkt)
}

func (m *Manager) encodeFLVPacket(pkt *media.Packet) []byte {
	if pkt == nil || len(pkt.Data) == 0 {
		return nil
	}

	var tagType byte
	switch pkt.Type {
	case media.PacketTypeVideo:
		tagType = 0x09
	case media.PacketTypeAudio:
		tagType = 0x08
	default:
		return nil
	}

	dataSize := len(pkt.Data)
	tag := make([]byte, 11+dataSize+4)
	tag[0] = tagType
	tag[1] = byte(dataSize >> 16)
	tag[2] = byte(dataSize >> 8)
	tag[3] = byte(dataSize)
	ts := pkt.Timestamp
	tag[4] = byte(ts >> 16)
	tag[5] = byte(ts >> 8)
	tag[6] = byte(ts)
	tag[7] = byte(ts >> 24)
	copy(tag[11:], pkt.Data)

	prevSize := uint32(11 + dataSize)
	tag[11+dataSize] = byte(prevSize >> 24)
	tag[12+dataSize] = byte(prevSize >> 16)
	tag[13+dataSize] = byte(prevSize >> 8)
	tag[14+dataSize] = byte(prevSize)

	return tag
}

// DeleteRecording removes a recording file
func (m *Manager) DeleteRecording(streamKey, filename string) error {
	path := filepath.Join(m.recordingsDir, streamKey, filepath.Base(filename))
	return os.Remove(path)
}

// OpenRecording opens a recording file for reading (DVR playback)
func (m *Manager) OpenRecording(streamKey, filename string) (io.ReadCloser, int64, error) {
	path := filepath.Join(m.recordingsDir, streamKey, filepath.Base(filename))
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, err
	}
	return f, info.Size(), nil
}

// CleanupOld removes recordings older than the given duration
func (m *Manager) CleanupOld(maxAge time.Duration) (int, error) {
	count := 0
	entries, err := os.ReadDir(m.recordingsDir)
	if err != nil {
		return 0, err
	}

	cutoff := time.Now().Add(-maxAge)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		files, _ := os.ReadDir(filepath.Join(m.recordingsDir, e.Name()))
		for _, f := range files {
			info, err := f.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				os.Remove(filepath.Join(m.recordingsDir, e.Name(), f.Name()))
				count++
			}
		}
	}
	return count, nil
}

// TrimLatestPerStream keeps only the newest N recordings per stream directory.
func (m *Manager) TrimLatestPerStream(keep int) (int, error) {
	if keep <= 0 {
		return 0, nil
	}
	totalDeleted := 0
	entries, err := os.ReadDir(m.recordingsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		streamKey := entry.Name()
		files, err := m.ListRecordingFiles(streamKey)
		if err != nil {
			return totalDeleted, err
		}
		for idx, file := range files {
			if idx < keep {
				continue
			}
			if err := os.Remove(file.Path); err == nil {
				totalDeleted++
			}
		}
	}
	return totalDeleted, nil
}

// StopAll stops all active recordings
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, rec := range m.recordings {
		select {
		case <-rec.stopCh:
		default:
			close(rec.stopCh)
		}
	}
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func normalizeFormat(format Format) Format {
	switch Format(strings.ToLower(string(format))) {
	case FormatMP4:
		return FormatMP4
	case FormatMKV:
		return FormatMKV
	case FormatFLV:
		return FormatFLV
	case FormatTS:
		return FormatTS
	default:
		return FormatMP4
	}
}
