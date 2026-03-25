package archive

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/recording"
	"github.com/fluxstream/fluxstream/internal/storage"
)

type Settings struct {
	Enabled       bool
	Provider      string
	LocalDir      string
	Endpoint      string
	Region        string
	Bucket        string
	AccessKey     string
	SecretKey     string
	Prefix        string
	UsePathStyle  bool
	PublicBaseURL string
	AutoUpload    bool
	DeleteLocal   bool
	ScanMinutes   int
	BatchSize     int
}

type Summary struct {
	Enabled           bool      `json:"enabled"`
	Provider          string    `json:"provider"`
	AutoUpload        bool      `json:"auto_upload"`
	DeleteLocal       bool      `json:"delete_local_after_upload"`
	LocalDir          string    `json:"local_dir,omitempty"`
	Endpoint          string    `json:"endpoint,omitempty"`
	Bucket            string    `json:"bucket,omitempty"`
	Prefix            string    `json:"prefix,omitempty"`
	Items             int       `json:"items"`
	ErrorItems        int       `json:"error_items"`
	LocalDeletedItems int       `json:"local_deleted_items"`
	LastSyncAt        time.Time `json:"last_sync_at,omitempty"`
	LastError         string    `json:"last_error,omitempty"`
}

type Manager struct {
	cfg        *config.Manager
	db         *storage.SQLiteDB
	recordings *recording.Manager
	dataDir    string

	mu         sync.RWMutex
	lastSyncAt time.Time
	lastError  string
}

type storeObject struct {
	Key  string
	URL  string
	ETag string
	Size int64
}

type storeClient interface {
	UploadFile(ctx context.Context, objectKey, localPath, contentType string) (storeObject, error)
	DownloadFile(ctx context.Context, objectKey, destPath string) (int64, error)
}

func NewManager(cfg *config.Manager, db *storage.SQLiteDB, recordings *recording.Manager, dataDir string) *Manager {
	return &Manager{
		cfg:        cfg,
		db:         db,
		recordings: recordings,
		dataDir:    dataDir,
	}
}

func (m *Manager) Settings() Settings {
	settings := Settings{
		Enabled:       m.cfg.GetBool("archive_enabled", false),
		Provider:      strings.ToLower(strings.TrimSpace(m.cfg.Get("archive_provider", "disabled"))),
		LocalDir:      strings.TrimSpace(m.cfg.Get("archive_local_dir", filepath.Join(m.dataDir, "archive"))),
		Endpoint:      strings.TrimRight(strings.TrimSpace(m.cfg.Get("archive_endpoint", "")), "/"),
		Region:        strings.TrimSpace(m.cfg.Get("archive_region", "us-east-1")),
		Bucket:        strings.TrimSpace(m.cfg.Get("archive_bucket", "")),
		AccessKey:     strings.TrimSpace(m.cfg.Get("archive_access_key", "")),
		SecretKey:     strings.TrimSpace(m.cfg.Get("archive_secret_key", "")),
		Prefix:        strings.Trim(strings.TrimSpace(m.cfg.Get("archive_prefix", "fluxstream")), "/"),
		UsePathStyle:  m.cfg.GetBool("archive_use_path_style", true),
		PublicBaseURL: strings.TrimRight(strings.TrimSpace(m.cfg.Get("archive_public_base_url", "")), "/"),
		AutoUpload:    m.cfg.GetBool("archive_auto_upload", false),
		DeleteLocal:   m.cfg.GetBool("archive_delete_local_after_upload", false),
		ScanMinutes:   m.cfg.GetInt("archive_scan_interval_minutes", 10),
		BatchSize:     m.cfg.GetInt("archive_batch_size", 3),
	}
	if settings.Provider == "" {
		settings.Provider = "disabled"
	}
	if settings.ScanMinutes <= 0 {
		settings.ScanMinutes = 10
	}
	if settings.BatchSize <= 0 {
		settings.BatchSize = 3
	}
	if settings.Provider == "local" && settings.LocalDir == "" {
		settings.LocalDir = filepath.Join(m.dataDir, "archive")
	}
	if settings.Provider == "disabled" {
		settings.Enabled = false
	}
	if !settings.Enabled {
		return settings
	}
	switch settings.Provider {
	case "local":
		settings.Enabled = strings.TrimSpace(settings.LocalDir) != ""
	case "s3", "minio":
		settings.Enabled = settings.Endpoint != "" && settings.Bucket != "" && settings.AccessKey != "" && settings.SecretKey != ""
	default:
		settings.Enabled = false
	}
	return settings
}

func (m *Manager) Enabled() bool {
	return m.Settings().Enabled
}

func (m *Manager) Summary() Summary {
	settings := m.Settings()
	items, _ := m.db.ListRecordingArchives("", 0)
	summary := Summary{
		Enabled:     settings.Enabled,
		Provider:    settings.Provider,
		AutoUpload:  settings.AutoUpload,
		DeleteLocal: settings.DeleteLocal,
		LocalDir:    settings.LocalDir,
		Endpoint:    settings.Endpoint,
		Bucket:      settings.Bucket,
		Prefix:      settings.Prefix,
		Items:       len(items),
	}
	for _, item := range items {
		if strings.EqualFold(item.Status, "error") {
			summary.ErrorItems++
		}
		if item.LocalDeleted {
			summary.LocalDeletedItems++
		}
	}
	m.mu.RLock()
	summary.LastSyncAt = m.lastSyncAt
	summary.LastError = m.lastError
	m.mu.RUnlock()
	return summary
}

func (m *Manager) ListArchives() ([]storage.RecordingArchive, error) {
	return m.db.ListRecordingArchives("", 0)
}

func (m *Manager) ArchiveRecording(ctx context.Context, streamKey, filename string) (*storage.RecordingArchive, error) {
	settings := m.Settings()
	if !settings.Enabled {
		return nil, fmt.Errorf("arsivleme etkin degil")
	}
	if m.recordings == nil {
		return nil, fmt.Errorf("kayit yoneticisi hazir degil")
	}
	streamKey = strings.TrimSpace(streamKey)
	filename = filepath.Base(strings.TrimSpace(filename))
	if streamKey == "" || filename == "" {
		return nil, fmt.Errorf("stream_key ve filename gerekli")
	}
	localPath := m.recordings.RecordingFilePath(streamKey, filename)
	info, err := os.Stat(localPath)
	if err != nil {
		return nil, err
	}
	client, err := m.newClient(settings)
	if err != nil {
		return nil, err
	}
	objectKey := buildObjectKey(settings.Prefix, streamKey, filename)
	item := &storage.RecordingArchive{
		StreamKey: streamKey,
		Filename:  filename,
		Format:    strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), "."),
		Provider:  providerLabel(settings.Provider),
		Bucket:    settings.Bucket,
		Endpoint:  settings.Endpoint,
		ObjectKey: objectKey,
		Size:      info.Size(),
		Status:    "archived",
	}
	obj, err := client.UploadFile(ctx, objectKey, localPath, detectArchiveContentType(filename))
	if err != nil {
		item.Status = "error"
		item.LastError = err.Error()
		_ = m.db.UpsertRecordingArchive(item)
		m.setLastError(err)
		return nil, err
	}
	item.ObjectURL = obj.URL
	item.ETag = obj.ETag
	item.Size = obj.Size
	item.ArchivedAt = time.Now()
	if settings.DeleteLocal {
		if err := os.Remove(localPath); err == nil {
			item.LocalDeleted = true
		}
	}
	if err := m.db.UpsertRecordingArchive(item); err != nil {
		return nil, err
	}
	m.markSyncSuccess()
	return m.db.GetRecordingArchive(streamKey, filename)
}

func (m *Manager) RestoreRecording(ctx context.Context, streamKey, filename string) (*storage.RecordingArchive, error) {
	settings := m.Settings()
	if !settings.Enabled {
		return nil, fmt.Errorf("arsivleme etkin degil")
	}
	streamKey = strings.TrimSpace(streamKey)
	filename = filepath.Base(strings.TrimSpace(filename))
	if streamKey == "" || filename == "" {
		return nil, fmt.Errorf("stream_key ve filename gerekli")
	}
	item, err := m.db.GetRecordingArchive(streamKey, filename)
	if err != nil {
		return nil, err
	}
	if item == nil || strings.TrimSpace(item.ObjectKey) == "" {
		return nil, fmt.Errorf("arsiv kaydi bulunamadi")
	}
	client, err := m.newClient(settings)
	if err != nil {
		return nil, err
	}
	if m.recordings == nil {
		return nil, fmt.Errorf("kayit yoneticisi hazir degil")
	}
	targetPath := m.recordings.RecordingFilePath(streamKey, filename)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return nil, err
	}
	size, err := client.DownloadFile(ctx, item.ObjectKey, targetPath)
	if err != nil {
		item.Status = "error"
		item.LastError = err.Error()
		_ = m.db.UpsertRecordingArchive(item)
		m.setLastError(err)
		return nil, err
	}
	item.Size = size
	item.Status = "archived"
	item.LastError = ""
	item.LocalDeleted = false
	item.RestoredAt = time.Now()
	if err := m.db.UpsertRecordingArchive(item); err != nil {
		return nil, err
	}
	m.markSyncSuccess()
	return m.db.GetRecordingArchive(streamKey, filename)
}

func (m *Manager) SyncPending(ctx context.Context, limit int) (int, error) {
	settings := m.Settings()
	if !settings.Enabled || !settings.AutoUpload || m.recordings == nil {
		return 0, nil
	}
	if limit <= 0 {
		limit = settings.BatchSize
	}
	files, err := m.recordings.ListAllRecordingFiles()
	if err != nil {
		m.setLastError(err)
		return 0, err
	}
	archives, err := m.db.ListRecordingArchives("", 0)
	if err != nil {
		m.setLastError(err)
		return 0, err
	}
	known := make(map[string]storage.RecordingArchive, len(archives))
	for _, item := range archives {
		known[item.StreamKey+"::"+item.Filename] = item
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.Before(files[j].ModTime)
	})
	uploaded := 0
	for _, file := range files {
		key := file.StreamKey + "::" + file.Name
		if existing, ok := known[key]; ok && strings.EqualFold(existing.Status, "archived") && strings.TrimSpace(existing.ObjectKey) != "" {
			continue
		}
		if _, err := m.ArchiveRecording(ctx, file.StreamKey, file.Name); err != nil {
			return uploaded, err
		}
		uploaded++
		if uploaded >= limit {
			break
		}
	}
	m.markSyncSuccess()
	return uploaded, nil
}

func (m *Manager) newClient(settings Settings) (storeClient, error) {
	switch settings.Provider {
	case "local":
		return &localStore{
			rootDir:       settings.LocalDir,
			publicBaseURL: settings.PublicBaseURL,
		}, nil
	case "s3", "minio":
		return &s3Store{
			endpoint:      settings.Endpoint,
			region:        settings.Region,
			bucket:        settings.Bucket,
			accessKey:     settings.AccessKey,
			secretKey:     settings.SecretKey,
			usePathStyle:  settings.UsePathStyle,
			publicBaseURL: settings.PublicBaseURL,
			httpClient:    &http.Client{Timeout: 2 * time.Minute},
		}, nil
	default:
		return nil, fmt.Errorf("desteklenmeyen arsiv saglayicisi: %s", settings.Provider)
	}
}

func (m *Manager) markSyncSuccess() {
	m.mu.Lock()
	m.lastSyncAt = time.Now()
	m.lastError = ""
	m.mu.Unlock()
}

func (m *Manager) setLastError(err error) {
	if err == nil {
		return
	}
	m.mu.Lock()
	m.lastError = err.Error()
	m.mu.Unlock()
}

type localStore struct {
	rootDir       string
	publicBaseURL string
}

func (s *localStore) UploadFile(_ context.Context, objectKey, localPath, _ string) (storeObject, error) {
	if err := os.MkdirAll(s.rootDir, 0755); err != nil {
		return storeObject{}, err
	}
	src, err := os.Open(localPath)
	if err != nil {
		return storeObject{}, err
	}
	defer src.Close()

	destPath := filepath.Join(s.rootDir, filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return storeObject{}, err
	}
	dst, err := os.Create(destPath)
	if err != nil {
		return storeObject{}, err
	}
	defer dst.Close()

	hasher := sha256.New()
	written, err := io.Copy(io.MultiWriter(dst, hasher), src)
	if err != nil {
		return storeObject{}, err
	}
	url := ""
	if strings.TrimSpace(s.publicBaseURL) != "" {
		url = strings.TrimRight(s.publicBaseURL, "/") + "/" + strings.TrimLeft(objectKey, "/")
	}
	return storeObject{
		Key:  objectKey,
		URL:  url,
		ETag: hex.EncodeToString(hasher.Sum(nil)),
		Size: written,
	}, nil
}

func (s *localStore) DownloadFile(_ context.Context, objectKey, destPath string) (int64, error) {
	src, err := os.Open(filepath.Join(s.rootDir, filepath.FromSlash(objectKey)))
	if err != nil {
		return 0, err
	}
	defer src.Close()
	tmpPath := destPath + ".restore"
	dst, err := os.Create(tmpPath)
	if err != nil {
		return 0, err
	}
	written, err := io.Copy(dst, src)
	closeErr := dst.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	return written, nil
}

type s3Store struct {
	endpoint      string
	region        string
	bucket        string
	accessKey     string
	secretKey     string
	usePathStyle  bool
	publicBaseURL string
	httpClient    *http.Client
}

func (s *s3Store) UploadFile(ctx context.Context, objectKey, localPath, contentType string) (storeObject, error) {
	file, err := os.Open(localPath)
	if err != nil {
		return storeObject{}, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return storeObject{}, err
	}
	payloadHash, err := hashSeeker(file)
	if err != nil {
		return storeObject{}, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return storeObject{}, err
	}
	req, canonicalURI, err := s.newSignedRequest(ctx, http.MethodPut, objectKey, nil, file, payloadHash)
	if err != nil {
		return storeObject{}, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return storeObject{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		return storeObject{}, fmt.Errorf("object storage upload hatasi (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return storeObject{
		Key:  objectKey,
		URL:  s.objectURL(objectKey, canonicalURI),
		ETag: strings.Trim(resp.Header.Get("ETag"), `"`),
		Size: info.Size(),
	}, nil
}

func (s *s3Store) DownloadFile(ctx context.Context, objectKey, destPath string) (int64, error) {
	req, _, err := s.newSignedRequest(ctx, http.MethodGet, objectKey, nil, nil, emptySHA256)
	if err != nil {
		return 0, err
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		return 0, fmt.Errorf("object storage indirme hatasi (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	tmpPath := destPath + ".restore"
	dst, err := os.Create(tmpPath)
	if err != nil {
		return 0, err
	}
	written, err := io.Copy(dst, resp.Body)
	closeErr := dst.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	return written, nil
}

func (s *s3Store) newSignedRequest(ctx context.Context, method, objectKey string, query url.Values, body io.Reader, payloadHash string) (*http.Request, string, error) {
	endpointURL, host, canonicalURI, canonicalQuery, err := s.buildObjectURL(objectKey, query)
	if err != nil {
		return nil, "", err
	}
	req, err := http.NewRequestWithContext(ctx, method, endpointURL, body)
	if err != nil {
		return nil, "", err
	}
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	dateStamp := now.Format("20060102")
	req.Host = host
	req.Header.Set("Host", host)
	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)

	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", host, payloadHash, amzDate)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalRequest := strings.Join([]string{
		method,
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")
	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateStamp, s.region)
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		credentialScope,
		hashString(canonicalRequest),
	}, "\n")
	signingKey := deriveAWS4SigningKey(s.secretKey, dateStamp, s.region, "s3")
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))
	req.Header.Set("Authorization", fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		s.accessKey, credentialScope, signedHeaders, signature,
	))
	return req, canonicalURI, nil
}

func (s *s3Store) buildObjectURL(objectKey string, query url.Values) (string, string, string, string, error) {
	base, err := url.Parse(s.endpoint)
	if err != nil {
		return "", "", "", "", err
	}
	cleanKey := strings.TrimLeft(objectKey, "/")
	basePath := strings.TrimRight(base.EscapedPath(), "/")
	objectPath := escapeObjectPath(cleanKey)
	canonicalURI := ""
	if s.usePathStyle {
		canonicalURI = joinURLPath(basePath, "/"+url.PathEscape(s.bucket)+"/"+objectPath)
	} else {
		base.Host = s.bucket + "." + base.Host
		canonicalURI = joinURLPath(basePath, "/"+objectPath)
	}
	base.Path = canonicalURI
	if query == nil {
		query = url.Values{}
	}
	base.RawQuery = query.Encode()
	return base.String(), base.Host, canonicalURI, query.Encode(), nil
}

func (s *s3Store) objectURL(objectKey, canonicalURI string) string {
	if strings.TrimSpace(s.publicBaseURL) != "" {
		return strings.TrimRight(s.publicBaseURL, "/") + "/" + strings.TrimLeft(objectKey, "/")
	}
	base, err := url.Parse(s.endpoint)
	if err != nil {
		return ""
	}
	if s.usePathStyle {
		base.Path = canonicalURI
	} else {
		base.Host = s.bucket + "." + base.Host
		base.Path = canonicalURI
	}
	base.RawQuery = ""
	return base.String()
}

func buildObjectKey(prefix, streamKey, filename string) string {
	parts := []string{}
	if prefix = strings.Trim(prefix, "/"); prefix != "" {
		parts = append(parts, prefix)
	}
	parts = append(parts, "recordings", sanitizeObjectPart(streamKey), sanitizeObjectPart(filepath.Base(filename)))
	return path.Join(parts...)
}

func detectArchiveContentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".flv":
		return "video/x-flv"
	case ".ts":
		return "video/mp2t"
	case ".mp3":
		return "audio/mpeg"
	case ".aac":
		return "audio/aac"
	case ".ogg":
		return "audio/ogg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	default:
		return "application/octet-stream"
	}
}

func providerLabel(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "minio":
		return "s3"
	default:
		return strings.ToLower(strings.TrimSpace(raw))
	}
}

func sanitizeObjectPart(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}
	value = strings.ReplaceAll(value, "\\", "_")
	value = strings.ReplaceAll(value, "..", "_")
	value = strings.ReplaceAll(value, " ", "_")
	return strings.Trim(value, "/")
}

func escapeObjectPath(objectKey string) string {
	parts := strings.Split(strings.TrimLeft(objectKey, "/"), "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}

func joinURLPath(basePath, suffix string) string {
	basePath = strings.TrimRight(basePath, "/")
	suffix = strings.TrimLeft(suffix, "/")
	if basePath == "" {
		return "/" + suffix
	}
	if suffix == "" {
		return basePath
	}
	return basePath + "/" + suffix
}

func hashSeeker(rs io.ReadSeeker) (string, error) {
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, rs); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func hashString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(data))
	return mac.Sum(nil)
}

func deriveAWS4SigningKey(secret, dateStamp, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), dateStamp)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	return hmacSHA256(kService, "aws4_request")
}

var emptySHA256 = hashString("")
