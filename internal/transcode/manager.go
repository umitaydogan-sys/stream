package transcode

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/media"
)

var ffmpegExecutableNames = []string{"ffmpeg"}

// Profile represents a transcoding profile (quality variant)
type Profile struct {
	Name       string `json:"name"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Bitrate    string `json:"bitrate"`     // e.g., "2500k"
	MaxBitrate string `json:"max_bitrate"` // e.g., "3000k"
	BufSize    string `json:"buf_size"`    // e.g., "5000k"
	Preset     string `json:"preset"`      // e.g., "fast", "medium"
	FPS        int    `json:"fps"`
	AudioRate  string `json:"audio_rate"` // e.g., "128k"
}

// GPUAccel represents GPU acceleration type
type GPUAccel string

const (
	GPUNone  GPUAccel = "none"
	GPUNVENC GPUAccel = "nvenc"
	GPUQSV   GPUAccel = "qsv"
	GPUAMF   GPUAccel = "amf"
	GPUVaapi GPUAccel = "vaapi"
)

// Job represents a transcoding job
type Job struct {
	ID           string    `json:"id"`
	StreamKey    string    `json:"stream_key"`
	Profiles     []Profile `json:"profiles"`
	Status       string    `json:"status"` // "pending", "running", "completed", "error"
	Error        string    `json:"error,omitempty"`
	StartedAt    time.Time `json:"started_at"`
	OutputDir    string    `json:"output_dir"`
	PID          int       `json:"pid,omitempty"`
	Type         string    `json:"type,omitempty"`
	ManifestPath string    `json:"manifest_path,omitempty"`
	cmd          *exec.Cmd
	cancel       context.CancelFunc
	stdin        io.WriteCloser
	packetCh     chan *media.Packet
	closeOnce    sync.Once
	logFile      *os.File
}

type liveTimestampState struct {
	baseTime        time.Time
	hasBase         bool
	lastTS          uint32
	audioSampleRate int
	audioBaseTS     uint32
	audioSamples    int64
	audioStarted    bool
}

// Manager manages FFmpeg transcoding
type Manager struct {
	ffmpegPath string
	gpuAccel   GPUAccel
	outputDir  string
	httpPort   int
	liveOpts   LiveOptions
	jobs       map[string]*Job
	liveJobs   map[string]*Job
	liveDash   map[string]*Job
	streamOpts map[string]LiveOptions
	mu         sync.RWMutex
}

// DefaultProfiles returns standard ABR profiles
func DefaultProfiles() []Profile {
	return []Profile{
		{Name: "1080p", Width: 1920, Height: 1080, Bitrate: "4500k", MaxBitrate: "5000k", BufSize: "9000k", Preset: "fast", FPS: 30, AudioRate: "192k"},
		{Name: "720p", Width: 1280, Height: 720, Bitrate: "2500k", MaxBitrate: "3000k", BufSize: "5000k", Preset: "fast", FPS: 30, AudioRate: "128k"},
		{Name: "480p", Width: 854, Height: 480, Bitrate: "1000k", MaxBitrate: "1200k", BufSize: "2000k", Preset: "fast", FPS: 30, AudioRate: "96k"},
		{Name: "360p", Width: 640, Height: 360, Bitrate: "600k", MaxBitrate: "700k", BufSize: "1200k", Preset: "fast", FPS: 25, AudioRate: "64k"},
	}
}

// NewManager creates a new transcoding manager
func NewManager(ffmpegPath string, gpuAccel GPUAccel, outputDir string) *Manager {
	os.MkdirAll(outputDir, 0755)
	os.MkdirAll(filepath.Join(outputDir, "hls"), 0755)
	if runtime.GOOS == "windows" {
		ffmpegExecutableNames = []string{"ffmpeg.exe", "ffmpeg"}
	}
	return &Manager{
		ffmpegPath: ffmpegPath,
		gpuAccel:   gpuAccel,
		outputDir:  outputDir,
		liveOpts:   DefaultLiveOptions(),
		jobs:       make(map[string]*Job),
		liveJobs:   make(map[string]*Job),
		liveDash:   make(map[string]*Job),
		streamOpts: make(map[string]LiveOptions),
	}
}

// SetHTTPPort stores the local HTTP port for loopback transcode inputs.
func (m *Manager) SetHTTPPort(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.httpPort = port
}

// SetLiveOptions updates default live HLS output behavior.
func (m *Manager) SetLiveOptions(opts LiveOptions) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.liveOpts = normalizeLiveOptions(opts)
}

// SetStreamLiveOptions overrides live output behavior for a single stream.
func (m *Manager) SetStreamLiveOptions(streamKey string, opts LiveOptions) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.streamOpts[streamKey] = normalizeLiveOptions(opts)
}

func (m *Manager) getLiveOptions(streamKey string) LiveOptions {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if opts, ok := m.streamOpts[streamKey]; ok {
		return normalizeLiveOptions(opts)
	}
	return normalizeLiveOptions(m.liveOpts)
}

// DetectFFmpeg checks if FFmpeg is available
func (m *Manager) DetectFFmpeg() (string, error) {
	for _, candidate := range m.ffmpegCandidates() {
		if path, ok := firstExistingExecutable(candidate); ok {
			return path, nil
		}
	}

	for _, name := range ffmpegExecutableNames {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("FFmpeg bulunamadi: fluxstream.exe yaninda, ffmpeg/ klasorunde veya PATH icinde ffmpeg gerekli")
}

func (m *Manager) ffmpegCandidates() []string {
	candidates := make([]string, 0, 16)
	execDir := ""
	workingDir := ""

	addNames := func(base string) {
		if strings.TrimSpace(base) == "" {
			return
		}
		for _, name := range ffmpegExecutableNames {
			candidates = append(candidates,
				filepath.Join(base, name),
				filepath.Join(base, "ffmpeg", name),
				filepath.Join(base, "bin", name),
				filepath.Join(base, "tools", name),
				filepath.Join(base, "tools", "ffmpeg", name),
				filepath.Join(base, "data", "tools", name),
				filepath.Join(base, "data", "tools", "ffmpeg", name),
			)
		}
	}

	addConfigured := func(candidate string) {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			return
		}
		candidates = append(candidates, candidate)
		if !filepath.IsAbs(candidate) {
			if execDir != "" {
				candidates = append(candidates, filepath.Join(execDir, candidate))
			}
			if workingDir != "" {
				candidates = append(candidates, filepath.Join(workingDir, candidate))
			}
		}
	}

	if execPath, err := os.Executable(); err == nil {
		execDir = filepath.Dir(execPath)
	}
	if wd, err := os.Getwd(); err == nil {
		workingDir = wd
	}

	if configured := strings.TrimSpace(m.ffmpegPath); configured != "" && configured != "ffmpeg" && configured != "ffmpeg.exe" {
		addConfigured(configured)
	}
	if envPath := strings.TrimSpace(os.Getenv("FLUXSTREAM_FFMPEG_PATH")); envPath != "" {
		addConfigured(envPath)
	}
	addNames(execDir)
	addNames(workingDir)

	return uniqueStrings(candidates)
}

func firstExistingExecutable(candidate string) (string, bool) {
	if strings.TrimSpace(candidate) == "" {
		return "", false
	}
	if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
		return candidate, true
	}
	return "", false
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		key := filepath.Clean(value)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, key)
	}
	return out
}

// DetectGPU checks available GPU acceleration
func (m *Manager) DetectGPU() GPUAccel {
	ffPath, err := m.DetectFFmpeg()
	if err != nil {
		return GPUNone
	}

	// Check NVENC
	cmd := exec.Command(ffPath, "-hide_banner", "-encoders")
	out, err := cmd.Output()
	if err != nil {
		return GPUNone
	}

	output := string(out)
	if strings.Contains(output, "h264_nvenc") {
		return GPUNVENC
	}
	if strings.Contains(output, "h264_qsv") {
		return GPUQSV
	}
	if strings.Contains(output, "h264_amf") {
		return GPUAMF
	}
	if strings.Contains(output, "h264_vaapi") {
		return GPUVaapi
	}

	return GPUNone
}

// StartTranscoding begins multi-quality ABR transcoding for a stream
func (m *Manager) StartTranscoding(streamKey, inputURL string, profiles []Profile) (*Job, error) {
	ffPath, err := m.DetectFFmpeg()
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	jobID := fmt.Sprintf("tc_%s_%d", streamKey, time.Now().Unix())
	jobOutputDir := filepath.Join(m.outputDir, streamKey)
	os.MkdirAll(jobOutputDir, 0755)

	ctx, cancel := context.WithCancel(context.Background())

	job := &Job{
		ID:        jobID,
		StreamKey: streamKey,
		Profiles:  profiles,
		Status:    "pending",
		StartedAt: time.Now(),
		OutputDir: jobOutputDir,
		cancel:    cancel,
	}

	// Build FFmpeg command
	args := m.buildFFmpegArgs(inputURL, jobOutputDir, profiles)
	cmd := exec.CommandContext(ctx, ffPath, args...)
	job.logFile = openLogFile(jobOutputDir)
	if job.logFile != nil {
		cmd.Stdout = job.logFile
		cmd.Stderr = job.logFile
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}
	job.cmd = cmd

	m.jobs[jobID] = job

	go m.runJob(job)
	return job, nil
}

// StartLiveHLS starts a controlled-GOP live HLS transcoder fed via FLV on stdin.
func (m *Manager) StartLiveHLS(streamKey string) (*Job, error) {
	ffPath, err := m.DetectFFmpeg()
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	if existing, ok := m.liveJobs[streamKey]; ok && (existing.Status == "running" || existing.Status == "pending") {
		m.mu.Unlock()
		return existing, nil
	}
	m.mu.Unlock()

	jobOutputDir := filepath.Join(m.GetLiveOutputDir(), streamKey)
	if err := os.MkdirAll(jobOutputDir, 0755); err != nil {
		return nil, err
	}
	cleanupLiveOutputDir(jobOutputDir)

	ctx, cancel := context.WithCancel(context.Background())
	liveOpts := m.getLiveOptions(streamKey)
	args := m.buildLiveHLSArgs(jobOutputDir, liveOpts)
	cmd := exec.CommandContext(ctx, ffPath, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, err
	}

	job := &Job{
		ID:           fmt.Sprintf("live_hls_%s_%d", streamKey, time.Now().Unix()),
		StreamKey:    streamKey,
		Status:       "pending",
		StartedAt:    time.Now(),
		OutputDir:    jobOutputDir,
		Type:         "live_hls",
		ManifestPath: filepath.Join(jobOutputDir, manifestNameForLiveOptions(liveOpts)),
		cmd:          cmd,
		cancel:       cancel,
		stdin:        stdin,
		packetCh:     make(chan *media.Packet, 4096),
	}
	job.logFile = openLogFile(jobOutputDir)
	if job.logFile != nil {
		cmd.Stdout = job.logFile
		cmd.Stderr = job.logFile
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	m.mu.Lock()
	if existing, ok := m.liveJobs[streamKey]; ok && (existing.Status == "running" || existing.Status == "pending") {
		m.mu.Unlock()
		cancel()
		_ = stdin.Close()
		return existing, nil
	}
	m.jobs[job.ID] = job
	m.liveJobs[streamKey] = job
	m.mu.Unlock()

	go m.runLiveHLSJob(job)
	return job, nil
}

// StartLiveDASH starts a live DASH repack job from the local HLS output.
func (m *Manager) StartLiveDASH(streamKey string) (*Job, error) {
	ffPath, err := m.DetectFFmpeg()
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	if existing, ok := m.liveDash[streamKey]; ok && (existing.Status == "running" || existing.Status == "pending") {
		m.mu.Unlock()
		return existing, nil
	}
	httpPort := m.httpPort
	m.mu.Unlock()

	if httpPort <= 0 {
		return nil, fmt.Errorf("http port not configured")
	}

	jobOutputDir := filepath.Join(m.GetLiveDashOutputDir(), streamKey)
	if err := os.MkdirAll(jobOutputDir, 0755); err != nil {
		return nil, err
	}
	cleanupLiveOutputDir(jobOutputDir)

	ctx, cancel := context.WithCancel(context.Background())
	inputURL := fmt.Sprintf("http://127.0.0.1:%d/hls/%s/%s", httpPort, streamKey, manifestNameForLiveOptions(m.getLiveOptions(streamKey)))
	args := m.buildLiveDASHArgs(inputURL, jobOutputDir)
	cmd := exec.CommandContext(ctx, ffPath, args...)
	cmd.Dir = jobOutputDir

	job := &Job{
		ID:           fmt.Sprintf("live_dash_%s_%d", streamKey, time.Now().Unix()),
		StreamKey:    streamKey,
		Status:       "pending",
		StartedAt:    time.Now(),
		OutputDir:    jobOutputDir,
		Type:         "live_dash",
		ManifestPath: filepath.Join(jobOutputDir, "manifest.mpd"),
		cmd:          cmd,
		cancel:       cancel,
	}
	job.logFile = openLogFile(jobOutputDir)
	if job.logFile != nil {
		cmd.Stdout = job.logFile
		cmd.Stderr = job.logFile
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	m.mu.Lock()
	if existing, ok := m.liveDash[streamKey]; ok && (existing.Status == "running" || existing.Status == "pending") {
		m.mu.Unlock()
		cancel()
		return existing, nil
	}
	m.jobs[job.ID] = job
	m.liveDash[streamKey] = job
	m.mu.Unlock()

	go m.runLiveDASHJob(job)
	return job, nil
}

// StopTranscoding stops a transcoding job
func (m *Manager) StopTranscoding(jobID string) {
	m.mu.RLock()
	job, exists := m.jobs[jobID]
	m.mu.RUnlock()

	if exists && job.cancel != nil {
		job.cancel()
		job.Status = "completed"
	}
}

// StopLiveHLS stops the live HLS transcoder for a stream.
func (m *Manager) StopLiveHLS(streamKey string) {
	m.mu.Lock()
	job, exists := m.liveJobs[streamKey]
	if exists {
		delete(m.liveJobs, streamKey)
	}
	m.mu.Unlock()

	if !exists {
		return
	}

	job.Status = "completed"
	job.closeInput()
	go func() {
		time.Sleep(2 * time.Second)
		if job.cancel != nil {
			job.cancel()
		}
		if job.OutputDir != "" {
			_ = os.RemoveAll(job.OutputDir)
		}
	}()
}

// StopLiveDASH stops the live DASH repack job for a stream.
func (m *Manager) StopLiveDASH(streamKey string) {
	m.mu.Lock()
	job, exists := m.liveDash[streamKey]
	if exists {
		delete(m.liveDash, streamKey)
	}
	m.mu.Unlock()

	if !exists {
		return
	}

	job.Status = "completed"
	go func() {
		time.Sleep(2 * time.Second)
		if job.cancel != nil {
			job.cancel()
		}
		if job.OutputDir != "" {
			_ = os.RemoveAll(job.OutputDir)
		}
	}()
}

// GetJob returns a job by ID
func (m *Manager) GetJob(jobID string) *Job {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.jobs[jobID]
}

// GetJobs returns all jobs
func (m *Manager) GetJobs() []*Job {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Job
	for _, j := range m.jobs {
		result = append(result, j)
	}
	return result
}

// WriteLivePacket feeds a live packet to an active live HLS job.
func (m *Manager) WriteLivePacket(streamKey string, pkt *media.Packet) {
	m.mu.RLock()
	job := m.liveJobs[streamKey]
	m.mu.RUnlock()

	if job == nil || job.packetCh == nil || pkt == nil {
		return
	}

	select {
	case job.packetCh <- pkt.Clone():
	default:
		log.Printf("[TC] CanlÄ± HLS kuyruk dolu, paket atÄ±ldÄ±: %s", streamKey)
	}
}

func (m *Manager) runJob(job *Job) {
	defer job.closeLog()
	job.Status = "running"
	log.Printf("[TC] Transkod baÅŸlatÄ±ldÄ±: %s", job.ID)

	if err := job.cmd.Start(); err != nil {
		job.Status = "error"
		job.Error = err.Error()
		log.Printf("[TC] BaÅŸlatma hatasÄ±: %v", err)
		return
	}

	job.PID = job.cmd.Process.Pid

	if err := job.cmd.Wait(); err != nil {
		if job.Status != "completed" {
			job.Status = "error"
			job.Error = err.Error()
		}
		return
	}

	job.Status = "completed"
	log.Printf("[TC] Transkod tamamlandÄ±: %s", job.ID)
}

func (m *Manager) runLiveHLSJob(job *Job) {
	defer job.closeLog()
	job.Status = "running"
	log.Printf("[TC] CanlÄ± HLS transcode baÅŸlatÄ±ldÄ±: %s", job.StreamKey)

	if err := job.cmd.Start(); err != nil {
		job.Status = "error"
		job.Error = err.Error()
		if job.stdin != nil {
			_ = job.stdin.Close()
		}
		log.Printf("[TC] CanlÄ± HLS baÅŸlatma hatasÄ± (%s): %v", job.StreamKey, err)
		m.removeLiveJob(job)
		return
	}

	job.PID = job.cmd.Process.Pid

	go func() {
		state := liveTimestampState{
			audioSampleRate: 48000,
		}

		defer func() {
			if job.stdin != nil {
				_ = job.stdin.Close()
			}
		}()

		if _, err := job.stdin.Write(buildFLVHeader(true, true)); err != nil {
			log.Printf("[TC] FLV header yazma hatasÄ± (%s): %v", job.StreamKey, err)
			return
		}

		for pkt := range job.packetCh {
			pkt = normalizePacketTimestamp(pkt, &state)
			tag := buildFLVTag(pkt)
			if tag == nil {
				continue
			}
			if _, err := job.stdin.Write(tag); err != nil {
				log.Printf("[TC] FLV packet yazma hatasÄ± (%s): %v", job.StreamKey, err)
				return
			}
		}
	}()

	if err := job.cmd.Wait(); err != nil {
		if job.Status != "completed" {
			job.Status = "error"
			job.Error = err.Error()
			log.Printf("[TC] CanlÄ± HLS hata ile sonlandÄ± (%s): %v", job.StreamKey, err)
		}
	} else if job.Status != "completed" {
		job.Status = "completed"
		log.Printf("[TC] CanlÄ± HLS tamamlandÄ±: %s", job.StreamKey)
	}

	job.closeInput()
	m.removeLiveJob(job)
}

func (m *Manager) runLiveStaticJob(job *Job, cleanup func(*Job), logPrefix string) {
	defer job.closeLog()
	job.Status = "running"
	log.Printf("%s baslatildi: %s", logPrefix, job.StreamKey)

	if err := job.cmd.Start(); err != nil {
		job.Status = "error"
		job.Error = err.Error()
		log.Printf("%s baslatma hatasi (%s): %v", logPrefix, job.StreamKey, err)
		cleanup(job)
		return
	}

	job.PID = job.cmd.Process.Pid

	if err := job.cmd.Wait(); err != nil {
		if job.Status != "completed" {
			job.Status = "error"
			job.Error = err.Error()
			log.Printf("%s hata (%s): %v", logPrefix, job.StreamKey, err)
		}
		cleanup(job)
		return
	}

	job.Status = "completed"
	cleanup(job)
}

func (m *Manager) runLiveDASHJob(job *Job) {
	if err := m.waitForLiveManifest(job.StreamKey, 20*time.Second); err != nil {
		defer job.closeLog()
		if job.Status != "completed" {
			job.Status = "error"
			job.Error = err.Error()
			log.Printf("[TC] Canli DASH baslatilamadi (%s): %v", job.StreamKey, err)
		}
		m.removeLiveDashJob(job)
		return
	}

	m.runLiveStaticJob(job, m.removeLiveDashJob, "[TC] Canli DASH")
}

func (m *Manager) buildFFmpegArgs(inputURL, outputDir string, profiles []Profile) []string {
	args := []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-i", inputURL,
	}

	// Add GPU acceleration flags
	switch m.gpuAccel {
	case GPUNVENC:
		args = append([]string{"-hwaccel", "cuda", "-hwaccel_output_format", "cuda"}, args...)
	case GPUQSV:
		args = append([]string{"-hwaccel", "qsv"}, args...)
	case GPUAMF:
		// AMF doesn't need special input flags
	case GPUVaapi:
		args = append([]string{"-hwaccel", "vaapi", "-vaapi_device", "/dev/dri/renderD128"}, args...)
	}

	encoder := m.selectVideoEncoder()

	// Generate multi-output HLS with master playlist
	for i, p := range profiles {
		args = append(args,
			"-map", "0:v:0", "-map", "0:a:0",
		)

		// Video encoding
		args = append(args,
			fmt.Sprintf("-c:v:%d", i), encoder,
			fmt.Sprintf("-b:v:%d", i), p.Bitrate,
			fmt.Sprintf("-maxrate:v:%d", i), p.MaxBitrate,
			fmt.Sprintf("-bufsize:v:%d", i), p.BufSize,
		)

		if encoder == "libx264" {
			args = append(args,
				fmt.Sprintf("-preset:v:%d", i), p.Preset,
			)
		}

		args = append(args,
			fmt.Sprintf("-s:v:%d", i), fmt.Sprintf("%dx%d", p.Width, p.Height),
			fmt.Sprintf("-r:v:%d", i), fmt.Sprintf("%d", p.FPS),
		)

		// Audio encoding
		args = append(args,
			fmt.Sprintf("-c:a:%d", i), "aac",
			fmt.Sprintf("-b:a:%d", i), p.AudioRate,
		)
	}

	// HLS output
	args = append(args,
		"-f", "hls",
		"-hls_time", "2",
		"-hls_list_size", "6",
		"-hls_flags", "delete_segments+independent_segments",
		"-hls_segment_type", "mpegts",
		"-master_pl_name", "master.m3u8",
	)

	// Variant stream map
	var varStreams []string
	for i := range profiles {
		varStreams = append(varStreams, fmt.Sprintf("v:%d,a:%d", i, i))
	}
	args = append(args, "-var_stream_map", strings.Join(varStreams, " "))

	args = append(args, filepath.Join(outputDir, "stream_%v", "data%d.ts"))

	return args
}

func manifestNameForLiveOptions(opts LiveOptions) string {
	opts = normalizeLiveOptions(opts)
	if opts.ABREnabled && opts.MasterEnabled && len(opts.Profiles) > 1 {
		return "master.m3u8"
	}
	return "index.m3u8"
}

func manifestCandidatesForLiveOptions(outputDir string, opts LiveOptions) []string {
	name := manifestNameForLiveOptions(opts)
	if name == "master.m3u8" {
		return []string{
			filepath.Join(outputDir, "master.m3u8.tmp"),
			filepath.Join(outputDir, "master.m3u8"),
			filepath.Join(outputDir, "index.m3u8.tmp"),
			filepath.Join(outputDir, "index.m3u8"),
		}
	}
	return []string{
		filepath.Join(outputDir, "index.m3u8.tmp"),
		filepath.Join(outputDir, "index.m3u8"),
		filepath.Join(outputDir, "master.m3u8.tmp"),
		filepath.Join(outputDir, "master.m3u8"),
	}
}

func normalizeLiveOptions(opts LiveOptions) LiveOptions {
	if opts.SegmentDuration <= 0 {
		opts.SegmentDuration = 2
	}
	if opts.PlaylistLength <= 0 {
		opts.PlaylistLength = 6
	}
	if strings.TrimSpace(opts.ProfileSet) == "" {
		opts.ProfileSet = "balanced"
	}
	if len(opts.Profiles) == 0 {
		opts.Profiles = ResolveProfiles(opts.ProfileSet, opts.ProfilesJSON)
	}
	return opts
}

func (m *Manager) buildLiveHLSArgs(outputDir string, opts LiveOptions) []string {
	opts = normalizeLiveOptions(opts)
	encoder := m.selectVideoEncoder()
	args := []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-fflags", "+genpts+igndts",
		"-f", "flv",
		"-i", "pipe:0",
		"-sn",
		"-dn",
	}

	if !(opts.ABREnabled && opts.MasterEnabled && len(opts.Profiles) > 1) {
		switch encoder {
		case "h264_nvenc":
			args = append(args,
				"-c:v", "h264_nvenc",
				"-preset", "p4",
				"-cq", "19",
				"-b:v", "0",
			)
		case "h264_qsv":
			args = append(args,
				"-c:v", "h264_qsv",
				"-preset", "medium",
				"-global_quality", "22",
			)
		case "h264_amf":
			args = append(args,
				"-c:v", "h264_amf",
				"-quality", "quality",
				"-rc", "cqp",
				"-qp_i", "19",
				"-qp_p", "21",
			)
		case "h264_vaapi":
			args = append(args,
				"-c:v", "h264_vaapi",
				"-b:v", "5000k",
				"-maxrate", "6500k",
				"-bufsize", "13000k",
			)
		case "libopenh264":
			args = append(args,
				"-c:v", "libopenh264",
				"-b:v", "5000k",
				"-maxrate", "6500k",
				"-bufsize", "13000k",
				"-pix_fmt", "yuv420p",
			)
		case "h264_mf":
			args = append(args,
				"-c:v", "h264_mf",
				"-b:v", "5000k",
				"-maxrate", "6500k",
				"-bufsize", "13000k",
				"-pix_fmt", "yuv420p",
			)
		default:
			args = append(args,
				"-c:v", "libx264",
				"-preset", "veryfast",
				"-tune", "zerolatency",
				"-crf", "21",
				"-maxrate", "6500k",
				"-bufsize", "13000k",
				"-pix_fmt", "yuv420p",
			)
		}
	}

	if opts.ABREnabled && opts.MasterEnabled && len(opts.Profiles) > 1 {
		for i, p := range opts.Profiles {
			args = append(args, "-map", "0:v:0", "-map", "0:a:0?")
			args = append(args,
				fmt.Sprintf("-c:v:%d", i), encoder,
				fmt.Sprintf("-b:v:%d", i), p.Bitrate,
				fmt.Sprintf("-maxrate:v:%d", i), p.MaxBitrate,
				fmt.Sprintf("-bufsize:v:%d", i), p.BufSize,
				fmt.Sprintf("-s:v:%d", i), fmt.Sprintf("%dx%d", p.Width, p.Height),
				fmt.Sprintf("-r:v:%d", i), fmt.Sprintf("%d", p.FPS),
				fmt.Sprintf("-g:v:%d", i), "60",
				fmt.Sprintf("-keyint_min:v:%d", i), "60",
				fmt.Sprintf("-sc_threshold:v:%d", i), "0",
			)
			switch encoder {
			case "h264_nvenc":
				args = append(args, fmt.Sprintf("-preset:v:%d", i), "p4", fmt.Sprintf("-cq:v:%d", i), "19")
			case "h264_qsv":
				args = append(args, fmt.Sprintf("-preset:v:%d", i), "medium", fmt.Sprintf("-global_quality:v:%d", i), "22")
			case "h264_amf":
				args = append(args, fmt.Sprintf("-quality:v:%d", i), "quality", fmt.Sprintf("-rc:v:%d", i), "cqp", fmt.Sprintf("-qp_i:v:%d", i), "19", fmt.Sprintf("-qp_p:v:%d", i), "21")
			case "h264_vaapi", "libopenh264", "h264_mf":
				args = append(args, fmt.Sprintf("-pix_fmt:v:%d", i), "yuv420p")
			default:
				args = append(args, fmt.Sprintf("-preset:v:%d", i), "veryfast", fmt.Sprintf("-crf:v:%d", i), "21", fmt.Sprintf("-tune:v:%d", i), "zerolatency", fmt.Sprintf("-pix_fmt:v:%d", i), "yuv420p")
			}
			if p.AudioRate != "" {
				args = append(args, fmt.Sprintf("-b:a:%d", i), p.AudioRate)
			}
		}
		args = append(args,
			"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", opts.SegmentDuration),
			"-c:a", "aac",
			"-ac", "2",
			"-ar", "48000",
			"-f", "hls",
			"-hls_time", fmt.Sprintf("%d", opts.SegmentDuration),
			"-hls_list_size", fmt.Sprintf("%d", opts.PlaylistLength),
			"-hls_allow_cache", "0",
			"-hls_flags", "delete_segments+independent_segments+append_list",
			"-master_pl_name", "master.m3u8",
		)
		var streamMap []string
		for _, p := range opts.Profiles {
			streamMap = append(streamMap, fmt.Sprintf("v:%d,a:%d,name:%s", len(streamMap), len(streamMap), sanitizeProfileName(p.Name)))
		}
		args = append(args, "-var_stream_map", strings.Join(streamMap, " "))
		args = append(args,
			"-hls_segment_filename", filepath.Join(outputDir, "%v", "seg_%06d.ts"),
			filepath.Join(outputDir, "%v", "index.m3u8"),
		)
		return args
	}

	audioRate := "160k"
	if len(opts.Profiles) > 0 && opts.Profiles[0].AudioRate != "" {
		audioRate = opts.Profiles[0].AudioRate
	}
	args = append(args,
		"-map", "0:v:0",
		"-map", "0:a:0?",
		"-g", "60",
		"-keyint_min", "60",
		"-sc_threshold", "0",
		"-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", opts.SegmentDuration),
		"-c:a", "aac",
		"-b:a", audioRate,
		"-ac", "2",
		"-ar", "48000",
		"-f", "hls",
		"-hls_time", fmt.Sprintf("%d", opts.SegmentDuration),
		"-hls_list_size", fmt.Sprintf("%d", opts.PlaylistLength),
		"-hls_allow_cache", "0",
		"-hls_flags", "delete_segments+independent_segments+append_list",
		"-hls_segment_filename", filepath.Join(outputDir, "seg_%06d.ts"),
		filepath.Join(outputDir, "index.m3u8"),
	)

	return args
}

func sanitizeProfileName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return "stream"
	}
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			continue
		}
		if r == '-' || r == '_' {
			b.WriteRune(r)
			continue
		}
	}
	if b.Len() == 0 {
		return "stream"
	}
	return b.String()
}

func (m *Manager) buildLiveDASHArgs(inputURL, _ string) []string {
	return []string{
		"-hide_banner",
		"-loglevel", "warning",
		"-user_agent", "FluxStreamInternal/2.0",
		"-fflags", "+genpts+igndts+discardcorrupt",
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "2",
		"-i", inputURL,
		"-map", "0:v:0",
		"-map", "0:a:0?",
		"-c:v", "copy",
		"-tag:v", "avc1",
		"-b:v", "2500k",
		"-c:a", "aac",
		"-b:a", "128k",
		"-ar", "48000",
		"-ac", "2",
		"-af", "aresample=async=1:min_hard_comp=0.100:first_pts=0,asetpts=N/SR/TB",
		"-avoid_negative_ts", "make_zero",
		"-f", "dash",
		"-seg_duration", "2",
		"-streaming", "1",
		"-remove_at_exit", "0",
		"-window_size", "6",
		"-extra_window_size", "0",
		"-use_template", "1",
		"-use_timeline", "0",
		"-ldash", "1",
		"-adaptation_sets", "id=0,streams=v id=1,streams=a",
		"-init_seg_name", "init-$RepresentationID$.m4s",
		"-media_seg_name", "chunk-$RepresentationID$-$Number%05d$.m4s",
		"manifest.mpd",
	}
}

func cleanupLiveOutputDir(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		_ = os.RemoveAll(filepath.Join(dir, entry.Name()))
	}
}

func (m *Manager) waitForLiveManifest(streamKey string, timeout time.Duration) error {
	if path := m.GetLiveManifestPath(streamKey); path != "" {
		return nil
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if path := m.GetLiveManifestPath(streamKey); path != "" {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("live hls manifest hazir degil: %s", streamKey)
}

func (m *Manager) selectVideoEncoder() string {
	encoders := m.availableEncoders()

	switch m.gpuAccel {
	case GPUNVENC:
		if strings.Contains(encoders, "h264_nvenc") {
			return "h264_nvenc"
		}
	case GPUQSV:
		if strings.Contains(encoders, "h264_qsv") {
			return "h264_qsv"
		}
	case GPUAMF:
		if strings.Contains(encoders, "h264_amf") {
			return "h264_amf"
		}
	case GPUVaapi:
		if strings.Contains(encoders, "h264_vaapi") {
			return "h264_vaapi"
		}
	}

	switch {
	case strings.Contains(encoders, "libx264"):
		return "libx264"
	case strings.Contains(encoders, "libopenh264"):
		return "libopenh264"
	case strings.Contains(encoders, "h264_mf"):
		return "h264_mf"
	case strings.Contains(encoders, "h264_qsv"):
		return "h264_qsv"
	case strings.Contains(encoders, "h264_amf"):
		return "h264_amf"
	case strings.Contains(encoders, "h264_nvenc"):
		return "h264_nvenc"
	default:
		return "libx264"
	}
}

func (m *Manager) availableEncoders() string {
	ffPath, err := m.DetectFFmpeg()
	if err != nil {
		return ""
	}

	out, err := exec.Command(ffPath, "-hide_banner", "-encoders").Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// GetStreamMasterPlaylist returns path to master.m3u8 for a stream
func (m *Manager) GetStreamMasterPlaylist(streamKey string) string {
	return filepath.Join(m.outputDir, streamKey, "master.m3u8")
}

// GetLiveManifestPath returns the preferred live playlist path for a stream.
func (m *Manager) GetLiveManifestPath(streamKey string) string {
	outputDir := filepath.Join(m.GetLiveOutputDir(), streamKey)
	for _, candidate := range manifestCandidatesForLiveOptions(outputDir, m.getLiveOptions(streamKey)) {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() && info.Size() > 0 {
			return candidate
		}
	}
	return ""
}

// WaitForLiveManifestPath waits until the preferred live playlist exists.
func (m *Manager) WaitForLiveManifestPath(streamKey string, timeout time.Duration) string {
	if path := m.GetLiveManifestPath(streamKey); path != "" {
		return path
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		time.Sleep(200 * time.Millisecond)
		if path := m.GetLiveManifestPath(streamKey); path != "" {
			return path
		}
	}
	return ""
}

// GetLiveManifestURL returns the local loopback HLS playlist URL for a stream.
func (m *Manager) GetLiveManifestURL(streamKey string) string {
	if m.httpPort <= 0 {
		return ""
	}
	path := m.GetLiveManifestPath(streamKey)
	if path == "" {
		return ""
	}
	root := filepath.Join(m.GetLiveOutputDir(), streamKey)
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}
	rel = filepath.ToSlash(rel)
	return fmt.Sprintf("http://127.0.0.1:%d/hls/%s/%s", m.httpPort, streamKey, rel)
}

// WaitForLiveManifestURL waits until a usable local loopback HLS playlist URL is available.
func (m *Manager) WaitForLiveManifestURL(streamKey string, timeout time.Duration) string {
	if m.httpPort <= 0 {
		return ""
	}
	if path := m.WaitForLiveManifestPath(streamKey, timeout); path != "" {
		root := filepath.Join(m.GetLiveOutputDir(), streamKey)
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return ""
		}
		rel = filepath.ToSlash(rel)
		return fmt.Sprintf("http://127.0.0.1:%d/hls/%s/%s", m.httpPort, streamKey, rel)
	}
	return ""
}

// GetLiveOutputDir returns the root directory for live HLS transcode outputs.
func (m *Manager) GetLiveOutputDir() string {
	return filepath.Join(m.outputDir, "hls")
}

// GetLiveDashOutputDir returns the root directory for live DASH repack outputs.
func (m *Manager) GetLiveDashOutputDir() string {
	return filepath.Join(m.outputDir, "dash")
}

// GetFFmpegVersion returns FFmpeg version info
func (m *Manager) GetFFmpegVersion() string {
	ffPath, err := m.DetectFFmpeg()
	if err != nil {
		return "not found"
	}
	cmd := exec.Command(ffPath, "-version")
	out, err := cmd.Output()
	if err != nil {
		return "error"
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return "unknown"
}

// GetStatus returns transcode manager status
func (m *Manager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	active := 0
	for _, j := range m.jobs {
		if j.Status == "running" {
			active++
		}
	}

	ffmpegPath, err := m.DetectFFmpeg()
	if err != nil {
		ffmpegPath = ""
	}

	return map[string]interface{}{
		"ffmpeg_version": m.GetFFmpegVersion(),
		"ffmpeg_path":    ffmpegPath,
		"gpu_accel":      string(m.gpuAccel),
		"active_jobs":    active,
		"total_jobs":     len(m.jobs),
		"os":             runtime.GOOS,
		"arch":           runtime.GOARCH,
		"live_options":   m.liveOpts,
	}
}

// StopAll stops all transcoding jobs
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, j := range m.jobs {
		j.closeInput()
		if j.cancel != nil {
			j.cancel()
		}
	}
	m.liveJobs = make(map[string]*Job)
	m.liveDash = make(map[string]*Job)
	m.streamOpts = make(map[string]LiveOptions)
}

func (m *Manager) removeLiveJob(job *Job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if live, ok := m.liveJobs[job.StreamKey]; ok && live == job {
		delete(m.liveJobs, job.StreamKey)
	}
}

func (m *Manager) removeLiveDashJob(job *Job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if live, ok := m.liveDash[job.StreamKey]; ok && live == job {
		delete(m.liveDash, job.StreamKey)
	}
}

func (j *Job) closeInput() {
	j.closeOnce.Do(func() {
		if j.packetCh != nil {
			close(j.packetCh)
		}
	})
}

func (j *Job) closeLog() {
	if j.logFile != nil {
		_ = j.logFile.Close()
		j.logFile = nil
	}
}

func openLogFile(outputDir string) *os.File {
	path := filepath.Join(outputDir, "ffmpeg.log")
	f, err := os.Create(path)
	if err != nil {
		return nil
	}
	return f
}

func buildFLVHeader(hasVideo, hasAudio bool) []byte {
	header := make([]byte, 13)
	copy(header[0:3], "FLV")
	header[3] = 0x01
	flags := byte(0)
	if hasAudio {
		flags |= 0x04
	}
	if hasVideo {
		flags |= 0x01
	}
	header[4] = flags
	binary.BigEndian.PutUint32(header[5:9], 9)
	return header
}

func buildFLVTag(pkt *media.Packet) []byte {
	if pkt == nil || len(pkt.Data) == 0 {
		return nil
	}

	var tagType byte
	switch pkt.Type {
	case media.PacketTypeVideo:
		tagType = 0x09
	case media.PacketTypeAudio:
		tagType = 0x08
	case media.PacketTypeMeta:
		tagType = 0x12
	default:
		return nil
	}

	dataSize := len(pkt.Data)
	tagSize := 11 + dataSize
	tag := make([]byte, tagSize+4)
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
	binary.BigEndian.PutUint32(tag[tagSize:], uint32(tagSize))
	return tag
}

func normalizePacketTimestamp(pkt *media.Packet, state *liveTimestampState) *media.Packet {
	if pkt == nil {
		return nil
	}

	clone := pkt.Clone()
	packetTime := clone.ReceivedAt
	if packetTime.IsZero() {
		packetTime = time.Now()
	}
	if !state.hasBase {
		state.baseTime = packetTime
		state.hasBase = true
	}

	elapsed := packetTime.Sub(state.baseTime)
	if elapsed < 0 {
		elapsed = 0
	}

	wallTS := uint32(elapsed / time.Millisecond)
	ts := wallTS

	if clone.Type == media.PacketTypeAudio {
		if clone.IsSequenceHeader {
			if rate := parseAACSampleRate(clone.Data); rate > 0 {
				state.audioSampleRate = rate
			}
		} else if state.audioSampleRate > 0 {
			if !state.audioStarted {
				state.audioBaseTS = wallTS
				state.audioSamples = 0
				state.audioStarted = true
			}
			ts = state.audioBaseTS + uint32((state.audioSamples*1000)/int64(state.audioSampleRate))
			state.audioSamples += 1024
			if ts < wallTS {
				ts = wallTS
			}
		}
	}

	if ts <= state.lastTS {
		ts = state.lastTS + 1
	}
	state.lastTS = ts
	clone.Timestamp = ts

	return clone
}

func parseAACSampleRate(data []byte) int {
	if len(data) < 4 {
		return 0
	}
	if len(data) < 2 || ((data[0]>>4)&0x0F) != byte(media.AudioCodecAAC) {
		return 0
	}
	if data[1] != 0 {
		return 0
	}

	asc := data[2:]
	if len(asc) < 2 {
		return 0
	}

	freqIdx := int(((asc[0] & 0x07) << 1) | ((asc[1] >> 7) & 0x01))
	sampleRates := []int{96000, 88200, 64000, 48000, 44100, 32000, 24000, 22050, 16000, 12000, 11025, 8000, 7350}
	if freqIdx < 0 || freqIdx >= len(sampleRates) {
		return 0
	}
	return sampleRates[freqIdx]
}

// MarshalJob returns JSON representation of a job
func MarshalJob(j *Job) ([]byte, error) {
	return json.Marshal(j)
}
