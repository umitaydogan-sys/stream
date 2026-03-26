package archive

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type sftpStore struct {
	host                string
	port                int
	user                string
	remoteDir           string
	keyPath             string
	disableHostKeyCheck bool
	publicBaseURL       string
}

func (s *sftpStore) UploadFile(ctx context.Context, objectKey, localPath, _ string) (storeObject, error) {
	if err := s.ensureTools(); err != nil {
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
	remotePath := s.remoteObjectPath(objectKey)
	if err := s.ensureRemoteDir(ctx, path.Dir(remotePath)); err != nil {
		return storeObject{}, err
	}
	args := append(s.baseSCPArgs(), localPath, s.remoteSpec(remotePath))
	cmd := exec.CommandContext(ctx, "scp", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return storeObject{}, fmt.Errorf("sftp upload hatasi: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return storeObject{
		Key:  objectKey,
		URL:  s.objectURL(objectKey),
		ETag: hash,
		Size: info.Size(),
	}, nil
}

func (s *sftpStore) DownloadFile(ctx context.Context, objectKey, destPath string) (int64, error) {
	if err := s.ensureTools(); err != nil {
		return 0, err
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return 0, err
	}
	tmpPath := destPath + ".restore"
	args := append(s.baseSCPArgs(), s.remoteSpec(s.remoteObjectPath(objectKey)), tmpPath)
	cmd := exec.CommandContext(ctx, "scp", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.Remove(tmpPath)
		return 0, fmt.Errorf("sftp indirme hatasi: %v: %s", err, strings.TrimSpace(string(out)))
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

func (s *sftpStore) DeleteObject(ctx context.Context, objectKey string) error {
	if err := s.ensureTools(); err != nil {
		return err
	}
	remotePath := s.remoteObjectPath(objectKey)
	args := append(s.baseSSHArgs(), s.remoteHost(), "rm -f -- "+shellQuote(remotePath))
	cmd := exec.CommandContext(ctx, "ssh", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("sftp silme hatasi: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (s *sftpStore) ensureTools() error {
	if _, err := exec.LookPath("ssh"); err != nil {
		return fmt.Errorf("ssh araci bulunamadi")
	}
	if _, err := exec.LookPath("scp"); err != nil {
		return fmt.Errorf("scp araci bulunamadi")
	}
	return nil
}

func (s *sftpStore) ensureRemoteDir(ctx context.Context, dir string) error {
	dir = strings.TrimSpace(dir)
	if dir == "" || dir == "." || dir == "/" {
		return nil
	}
	args := append(s.baseSSHArgs(), s.remoteHost(), "mkdir -p -- "+shellQuote(dir))
	cmd := exec.CommandContext(ctx, "ssh", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("sftp uzak dizin hazirlanamadi: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (s *sftpStore) baseSSHArgs() []string {
	args := []string{"-o", "BatchMode=yes"}
	if s.port > 0 {
		args = append(args, "-p", strconv.Itoa(s.port))
	}
	if strings.TrimSpace(s.keyPath) != "" {
		args = append(args, "-i", s.keyPath)
	}
	if s.disableHostKeyCheck {
		args = append(args, "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile="+os.DevNull)
	}
	return args
}

func (s *sftpStore) baseSCPArgs() []string {
	args := []string{"-q"}
	if s.port > 0 {
		args = append(args, "-P", strconv.Itoa(s.port))
	}
	if strings.TrimSpace(s.keyPath) != "" {
		args = append(args, "-i", s.keyPath)
	}
	args = append(args, "-o", "BatchMode=yes")
	if s.disableHostKeyCheck {
		args = append(args, "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile="+os.DevNull)
	}
	return args
}

func (s *sftpStore) remoteHost() string {
	if strings.TrimSpace(s.user) == "" {
		return s.host
	}
	return s.user + "@" + s.host
}

func (s *sftpStore) remoteSpec(remotePath string) string {
	return s.remoteHost() + ":" + remotePath
}

func (s *sftpStore) remoteObjectPath(objectKey string) string {
	root := strings.TrimRight(strings.TrimSpace(s.remoteDir), "/")
	clean := strings.TrimLeft(objectKey, "/")
	if root == "" {
		return "/" + clean
	}
	return path.Join(root, clean)
}

func (s *sftpStore) objectURL(objectKey string) string {
	if strings.TrimSpace(s.publicBaseURL) == "" {
		return ""
	}
	return strings.TrimRight(s.publicBaseURL, "/") + "/" + strings.TrimLeft(objectKey, "/")
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}
