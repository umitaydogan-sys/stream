package archive

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type rcloneStore struct {
	remote        string
	remotePath    string
	configPath    string
	publicBaseURL string
}

func (s *rcloneStore) UploadFile(ctx context.Context, objectKey, localPath, _ string) (storeObject, error) {
	if err := s.ensureTool(); err != nil {
		return storeObject{}, err
	}
	info, err := os.Stat(localPath)
	if err != nil {
		return storeObject{}, err
	}
	file, err := os.Open(localPath)
	if err != nil {
		return storeObject{}, err
	}
	defer file.Close()
	hash, err := hashSeeker(file)
	if err != nil {
		return storeObject{}, err
	}
	cmd := s.command(ctx, "copyto", localPath, s.remoteSpec(objectKey))
	if out, err := cmd.CombinedOutput(); err != nil {
		return storeObject{}, fmt.Errorf("rclone upload hatasi: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return storeObject{
		Key:  objectKey,
		URL:  s.objectURL(objectKey),
		ETag: hash,
		Size: info.Size(),
	}, nil
}

func (s *rcloneStore) DownloadFile(ctx context.Context, objectKey, destPath string) (int64, error) {
	if err := s.ensureTool(); err != nil {
		return 0, err
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return 0, err
	}
	tmpPath := destPath + ".restore"
	cmd := s.command(ctx, "copyto", s.remoteSpec(objectKey), tmpPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.Remove(tmpPath)
		return 0, fmt.Errorf("rclone indirme hatasi: %v: %s", err, strings.TrimSpace(string(out)))
	}
	info, err := os.Stat(tmpPath)
	if err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return 0, err
	}
	return info.Size(), nil
}

func (s *rcloneStore) DeleteObject(ctx context.Context, objectKey string) error {
	if err := s.ensureTool(); err != nil {
		return err
	}
	cmd := s.command(ctx, "deletefile", s.remoteSpec(objectKey))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("rclone silme hatasi: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (s *rcloneStore) ensureTool() error {
	if _, err := exec.LookPath("rclone"); err != nil {
		return fmt.Errorf("rclone bulunamadi")
	}
	if strings.TrimSpace(s.remote) == "" {
		return fmt.Errorf("rclone baglanti profili gerekli")
	}
	return nil
}

func (s *rcloneStore) command(ctx context.Context, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "rclone", args...)
	if strings.TrimSpace(s.configPath) != "" {
		cmd.Env = append(os.Environ(), "RCLONE_CONFIG="+s.configPath)
	}
	return cmd
}

func (s *rcloneStore) remoteSpec(objectKey string) string {
	base := strings.TrimSuffix(strings.TrimSpace(s.remote), ":")
	fullPath := strings.Trim(strings.TrimSpace(s.remotePath), "/")
	if fullPath != "" {
		fullPath = path.Join(fullPath, strings.TrimLeft(objectKey, "/"))
	} else {
		fullPath = strings.TrimLeft(objectKey, "/")
	}
	return base + ":" + fullPath
}

func (s *rcloneStore) objectURL(objectKey string) string {
	if strings.TrimSpace(s.publicBaseURL) != "" {
		return strings.TrimRight(s.publicBaseURL, "/") + "/" + strings.TrimLeft(objectKey, "/")
	}
	return "rclone://" + s.remoteSpec(objectKey)
}
