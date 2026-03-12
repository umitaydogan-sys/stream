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
	ID        string
	StreamKey string
	Format    Format
	FilePath  string
	StartedAt time.Time
	Duration  time.Duration
	Size      int64
	Status    string // "recording", "completed", "error"
	file      *os.File
	tsMuxer   *ts.Muxer
	subID     string
	stopCh    chan struct{}
	mu        sync.Mutex
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

// Manager manages stream recordings and DVR
type Manager struct {
	streamMgr     *stream.Manager
	recordingsDir string
	recordings    map[string]*Recording
	maxDuration   time.Duration
	mu            sync.RWMutex
}

// NewManager creates a new recording manager
func NewManager(streamMgr *stream.Manager, recordingsDir string) *Manager {
	os.MkdirAll(recordingsDir, 0755)
	return &Manager{
		streamMgr:     streamMgr,
		recordingsDir: recordingsDir,
		recordings:    make(map[string]*Recording),
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

	file, err := os.Create(filePath)
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
		ID:        recID,
		StreamKey: streamKey,
		Format:    format,
		FilePath:  filePath,
		StartedAt: time.Now(),
		Status:    "recording",
		file:      file,
		tsMuxer:   ts.NewMuxer(),
		stopCh:    make(chan struct{}),
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
		rec.Status = "completed"
		rec.file.Close()
		rec.mu.Unlock()
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
				rec.Status = "error"
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
	if pkt.IsSequenceHeader {
		return nil
	}
	mediaPkt := pkt.Clone()
	if pkt.Type == media.PacketTypeVideo && len(pkt.Data) > 5 {
		mediaPkt.Data = pkt.Data[5:]
	} else if pkt.Type == media.PacketTypeAudio && len(pkt.Data) > 2 {
		mediaPkt.Data = pkt.Data[2:]
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
	case FormatFLV:
		return FormatFLV
	case FormatMP4:
		return FormatTS
	case FormatMKV:
		return FormatTS
	case FormatTS:
		fallthrough
	default:
		return FormatTS
	}
}
