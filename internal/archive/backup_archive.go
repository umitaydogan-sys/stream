package archive

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/storage"
	"github.com/fluxstream/fluxstream/internal/systemutil"
)

func (m *Manager) ListBackupArchives() ([]storage.BackupArchive, error) {
	return m.db.ListBackupArchives(0)
}

func (m *Manager) ArchiveBackup(ctx context.Context, name string) (*storage.BackupArchive, error) {
	settings := m.Settings()
	target := settings.Backup
	if !target.Configured || !settings.BackupsEnabled {
		return nil, fmt.Errorf("yedek arsivi etkin degil")
	}
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" {
		return nil, fmt.Errorf("backup name gerekli")
	}
	localPath := systemutil.BackupFilePath(m.dataDir, name)
	info, err := os.Stat(localPath)
	if err != nil {
		return nil, err
	}
	client, err := m.newClient(target)
	if err != nil {
		return nil, err
	}
	objectKey := buildBackupObjectKey(target.Prefix, name)
	item := &storage.BackupArchive{
		Name:              name,
		Provider:          providerLabel(target.Provider, target.ProviderVariant),
		Bucket:            target.Bucket,
		Endpoint:          target.Endpoint,
		ObjectKey:         objectKey,
		Size:              info.Size(),
		IncludeRecordings: strings.Contains(strings.ToLower(name), "with-recordings"),
		Status:            "archived",
	}
	obj, err := client.UploadFile(ctx, objectKey, localPath, "application/gzip")
	if err != nil {
		item.Status = "error"
		item.LastError = err.Error()
		_ = m.db.UpsertBackupArchive(item)
		m.setLastBackupError(err)
		return nil, err
	}
	item.ObjectURL = obj.URL
	item.ETag = obj.ETag
	item.Size = obj.Size
	item.ArchivedAt = time.Now()
	if target.DeleteLocal {
		if err := os.Remove(localPath); err == nil {
			item.LocalDeleted = true
		}
	}
	if err := m.db.UpsertBackupArchive(item); err != nil {
		return nil, err
	}
	m.markBackupSyncSuccess()
	return m.db.GetBackupArchive(name)
}

func (m *Manager) RestoreBackupArchive(ctx context.Context, name string) (*storage.BackupArchive, error) {
	settings := m.Settings()
	target := settings.Backup
	if !target.Configured || !settings.BackupsEnabled {
		return nil, fmt.Errorf("yedek arsivi etkin degil")
	}
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" {
		return nil, fmt.Errorf("backup name gerekli")
	}
	item, err := m.db.GetBackupArchive(name)
	if err != nil {
		return nil, err
	}
	if item == nil || strings.TrimSpace(item.ObjectKey) == "" {
		return nil, fmt.Errorf("arsiv yedegi bulunamadi")
	}
	client, err := m.newClient(target)
	if err != nil {
		return nil, err
	}
	targetPath := systemutil.BackupFilePath(m.dataDir, name)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return nil, err
	}
	size, err := client.DownloadFile(ctx, item.ObjectKey, targetPath)
	if err != nil {
		item.Status = "error"
		item.LastError = err.Error()
		_ = m.db.UpsertBackupArchive(item)
		m.setLastBackupError(err)
		return nil, err
	}
	item.Size = size
	item.Status = "archived"
	item.LastError = ""
	item.LocalDeleted = false
	item.RestoredAt = time.Now()
	if err := m.db.UpsertBackupArchive(item); err != nil {
		return nil, err
	}
	m.markBackupSyncSuccess()
	return m.db.GetBackupArchive(name)
}

func (m *Manager) MarkBackupLocalDeleted(name string, deleted bool) error {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" {
		return fmt.Errorf("backup name gerekli")
	}
	item, err := m.db.GetBackupArchive(name)
	if err != nil {
		return err
	}
	if item == nil {
		return nil
	}
	item.LocalDeleted = deleted
	item.UpdatedAt = time.Now()
	if err := m.db.UpsertBackupArchive(item); err != nil {
		return err
	}
	return nil
}

func (m *Manager) SyncPendingBackups(ctx context.Context, limit int) (int, error) {
	settings := m.Settings()
	target := settings.Backup
	if !target.Configured || !settings.BackupsEnabled || !target.AutoUpload {
		return 0, nil
	}
	if limit <= 0 {
		limit = target.BatchSize
	}
	files, err := systemutil.ListBackups(m.dataDir)
	if err != nil {
		m.setLastBackupError(err)
		return 0, err
	}
	archives, err := m.db.ListBackupArchives(0)
	if err != nil {
		m.setLastBackupError(err)
		return 0, err
	}
	known := make(map[string]storage.BackupArchive, len(archives))
	for _, item := range archives {
		known[item.Name] = item
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.Before(files[j].ModTime)
	})
	uploaded := 0
	for _, file := range files {
		if existing, ok := known[file.Name]; ok && strings.EqualFold(existing.Status, "archived") && strings.TrimSpace(existing.ObjectKey) != "" {
			continue
		}
		if _, err := m.ArchiveBackup(ctx, file.Name); err != nil {
			return uploaded, err
		}
		uploaded++
		if uploaded >= limit {
			break
		}
	}
	m.markBackupSyncSuccess()
	return uploaded, nil
}

func buildBackupObjectKey(prefix, name string) string {
	parts := []string{}
	if prefix = strings.Trim(prefix, "/"); prefix != "" {
		parts = append(parts, prefix)
	}
	parts = append(parts, "backups", sanitizeObjectPart(filepath.Base(name)))
	return path.Join(parts...)
}
