package systemutil

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/storage"
)

type BackupInfo struct {
	Name              string    `json:"name"`
	Size              int64     `json:"size"`
	ModTime           time.Time `json:"mod_time"`
	IncludeRecordings bool      `json:"include_recordings"`
}

type backupManifest struct {
	CreatedAt         time.Time `json:"created_at"`
	IncludeRecordings bool      `json:"include_recordings"`
	Items             []string  `json:"items"`
}

func BackupDir(dataDir string) string {
	return filepath.Join(dataDir, "backups")
}

func CreateBackup(dataDir string, db *storage.SQLiteDB, includeRecordings bool) (BackupInfo, error) {
	if err := os.MkdirAll(BackupDir(dataDir), 0755); err != nil {
		return BackupInfo{}, err
	}
	stamp := time.Now().UTC().Format("20060102-150405")
	name := fmt.Sprintf("fluxstream-backup-%s%s.tar.gz", stamp, ternary(includeRecordings, "-with-recordings", ""))
	path := filepath.Join(BackupDir(dataDir), name)
	tmpDB := filepath.Join(os.TempDir(), fmt.Sprintf("fluxstream-snapshot-%s.db", stamp))
	defer os.Remove(tmpDB)
	if err := db.ExportBackupSnapshot(tmpDB); err != nil {
		return BackupInfo{}, err
	}
	file, err := os.Create(path)
	if err != nil {
		return BackupInfo{}, err
	}
	defer file.Close()
	gz := gzip.NewWriter(file)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	manifest := backupManifest{CreatedAt: time.Now().UTC(), IncludeRecordings: includeRecordings}
	if err := addFileToTar(tw, tmpDB, "data/fluxstream.db"); err != nil {
		return BackupInfo{}, err
	}
	manifest.Items = append(manifest.Items, "data/fluxstream.db")
	for _, dir := range []string{"certs", "players", "license"} {
		root := filepath.Join(dataDir, dir)
		items, err := addDirTreeToTar(tw, root, filepath.ToSlash(filepath.Join("data", dir)))
		if err != nil {
			return BackupInfo{}, err
		}
		manifest.Items = append(manifest.Items, items...)
	}
	if includeRecordings {
		items, err := addDirTreeToTar(tw, filepath.Join(dataDir, "recordings"), "data/recordings")
		if err != nil {
			return BackupInfo{}, err
		}
		manifest.Items = append(manifest.Items, items...)
	}
	sort.Strings(manifest.Items)
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return BackupInfo{}, err
	}
	if err := writeTarEntry(tw, "backup-manifest.json", manifestBytes, manifest.CreatedAt); err != nil {
		return BackupInfo{}, err
	}
	if err := tw.Close(); err != nil {
		return BackupInfo{}, err
	}
	if err := gz.Close(); err != nil {
		return BackupInfo{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return BackupInfo{}, err
	}
	return BackupInfo{Name: info.Name(), Size: info.Size(), ModTime: info.ModTime(), IncludeRecordings: includeRecordings}, nil
}

func ListBackups(dataDir string) ([]BackupInfo, error) {
	entries, err := os.ReadDir(BackupDir(dataDir))
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupInfo{}, nil
		}
		return nil, err
	}
	items := make([]BackupInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tar.gz") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		items = append(items, BackupInfo{
			Name:              entry.Name(),
			Size:              info.Size(),
			ModTime:           info.ModTime(),
			IncludeRecordings: strings.Contains(entry.Name(), "with-recordings"),
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ModTime.After(items[j].ModTime) })
	return items, nil
}

func DeleteBackup(dataDir, name string) error {
	clean := filepath.Base(name)
	if clean == "." || clean == "" {
		return fmt.Errorf("gecersiz backup adi")
	}
	return os.Remove(filepath.Join(BackupDir(dataDir), clean))
}

func BackupFilePath(dataDir, name string) string {
	return filepath.Join(BackupDir(dataDir), filepath.Base(name))
}

func RestoreBackup(archivePath, dataDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	restoreRoots := map[string]bool{
		"data/fluxstream.db": true,
		"data/certs":         true,
		"data/players":       true,
		"data/license":       true,
		"data/recordings":    true,
	}
	prepareTargets := func() error {
		for rel := range restoreRoots {
			target := filepath.Join(dataDir, strings.TrimPrefix(rel, "data/"))
			if strings.HasSuffix(rel, ".db") {
				if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
					return err
				}
				continue
			}
			if err := os.RemoveAll(target); err != nil {
				return err
			}
		}
		return nil
	}
	if err := prepareTargets(); err != nil {
		return err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := filepath.ToSlash(strings.TrimSpace(hdr.Name))
		if name == "" || name == "." || strings.HasPrefix(name, "../") {
			continue
		}
		if name == "backup-manifest.json" {
			continue
		}
		if !strings.HasPrefix(name, "data/") {
			continue
		}

		target := filepath.Join(dataDir, strings.TrimPrefix(name, "data/"))
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			out.Close()
			_ = os.Chtimes(target, time.Now(), hdr.ModTime)
		}
	}
	return nil
}

func addDirTreeToTar(tw *tar.Writer, root, prefix string) ([]string, error) {
	items := make([]string, 0, 16)
	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return items, nil
		}
		return nil, err
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		tarPath := filepath.ToSlash(filepath.Join(prefix, rel))
		if err := addFileToTar(tw, path, tarPath); err != nil {
			return err
		}
		items = append(items, tarPath)
		return nil
	})
	return items, err
}

func addFileToTar(tw *tar.Writer, src, tarPath string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	hdr.Name = filepath.ToSlash(tarPath)
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	return err
}

func writeTarEntry(tw *tar.Writer, name string, data []byte, modTime time.Time) error {
	hdr := &tar.Header{Name: name, Mode: 0644, Size: int64(len(data)), ModTime: modTime}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := tw.Write(data)
	return err
}

func ternary(ok bool, yes, no string) string {
	if ok {
		return yes
	}
	return no
}
