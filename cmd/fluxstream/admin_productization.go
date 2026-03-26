package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/archive"
	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/license"
	"github.com/fluxstream/fluxstream/internal/storage"
	"github.com/fluxstream/fluxstream/internal/systemutil"
	"github.com/fluxstream/fluxstream/internal/web"
)

const linuxServiceUnit = "fluxstream"
const windowsServiceUnit = "FluxStream"

func registerProductAdminRoutes(webServer *web.Server, cfg *config.Manager, db *storage.SQLiteDB, dataDir string, runtimeLicense *runtimeLicense, archiveManager *archive.Manager) {
	licMgr := license.NewManager(dataDir)

	webServer.RegisterAdminHandler("/api/license/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		jsonResp(w, map[string]interface{}{
			"success": true,
			"status":  licMgr.Status(time.Now().UTC()),
			"runtime": runtimeLicense,
		})
	})

	webServer.RegisterAdminHandler("/api/license/sample", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		jsonResp(w, map[string]interface{}{
			"success": true,
			"sample":  license.SampleDocument(),
		})
	})

	webServer.RegisterAdminHandler("/api/license/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			LicenseJSON  string `json:"license_json"`
			PublicKeyPEM string `json:"public_key_pem"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.PublicKeyPEM) != "" {
			if err := licMgr.SavePublicKey(req.PublicKeyPEM); err != nil {
				http.Error(w, fmt.Sprintf("public key kaydedilemedi: %v", err), http.StatusBadRequest)
				return
			}
		}
		if strings.TrimSpace(req.LicenseJSON) == "" {
			http.Error(w, "license_json gerekli", http.StatusBadRequest)
			return
		}
		if err := licMgr.SaveLicense(req.LicenseJSON); err != nil {
			http.Error(w, fmt.Sprintf("lisans kaydedilemedi: %v", err), http.StatusBadRequest)
			return
		}
		jsonResp(w, map[string]interface{}{
			"success": true,
			"status":  licMgr.Status(time.Now().UTC()),
		})
	})

	webServer.RegisterAdminHandler("/api/system/service/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		unit := serviceUnitName()
		jsonResp(w, map[string]interface{}{
			"success":  true,
			"status":   systemutil.GetServiceStatus(unit),
			"unit":     unit,
			"platform": runtime.GOOS,
		})
	})

	webServer.RegisterAdminHandler("/api/system/service/action", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Action string `json:"action"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		action := strings.ToLower(strings.TrimSpace(req.Action))
		if action == "" {
			http.Error(w, "action gerekli", http.StatusBadRequest)
			return
		}
		unit := serviceUnitName()
		if err := systemutil.RunServiceAction(unit, action); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResp(w, map[string]interface{}{
			"success": true,
			"action":  action,
			"status":  systemutil.GetServiceStatus(unit),
			"unit":    unit,
		})
	})

	webServer.RegisterAdminHandler("/api/system/upgrade/plan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		execPath, err := os.Executable()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		execDir := filepath.Dir(execPath)
		unit := serviceUnitName()
		binaryName := filepath.Base(execPath)
		binaryPath := execPath
		nextBinary := filepath.Join(execDir, binaryName+".next")
		if runtime.GOOS == "linux" {
			nextBinary = filepath.Join(execDir, "fluxstream.next")
		}
		owner := "root"
		if out, err := exec.Command("whoami").Output(); err == nil {
			owner = strings.TrimSpace(string(out))
		}
		jsonResp(w, map[string]interface{}{
			"success":      true,
			"platform":     runtime.GOOS,
			"service_unit": unit,
			"install_dir":  execDir,
			"binary_path":  binaryPath,
			"next_binary":  nextBinary,
			"data_dir":     dataDir,
			"backup_dir":   filepath.Join(dataDir, "backups"),
			"owner":        owner,
			"commands": map[string]string{
				"backup_create":   fmt.Sprintf("%s backup create", binaryPath),
				"backup_restore":  backupRestoreCommand(runtime.GOOS, unit, binaryPath),
				"atomic_upgrade":  upgradeCommand(runtime.GOOS, unit, binaryPath, nextBinary),
				"service_restart": serviceCommand(runtime.GOOS, unit, "restart"),
			},
		})
	})

	webServer.RegisterAdminHandler("/api/system/backups", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			items, err := systemutil.ListBackups(dataDir)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jsonResp(w, map[string]interface{}{"success": true, "items": items})
		case http.MethodPost:
			var req struct {
				IncludeRecordings bool `json:"include_recordings"`
			}
			if r.ContentLength > 0 {
				if err := decodeJSON(r, &req); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			item, err := systemutil.CreateBackup(dataDir, db, req.IncludeRecordings)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var archived *storage.BackupArchive
			if archiveManager != nil {
				settings := archiveManager.Settings()
				if settings.Configured && settings.BackupsEnabled && settings.BackupAutoUpload {
					archived, _ = archiveManager.ArchiveBackup(r.Context(), item.Name)
				}
			}
			jsonResp(w, map[string]interface{}{"success": true, "item": item, "archived": archived})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	webServer.RegisterAdminHandler("/api/system/backups/archives", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if archiveManager == nil {
			jsonResp(w, []storage.BackupArchive{})
			return
		}
		items, err := archiveManager.ListBackupArchives()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if items == nil {
			items = []storage.BackupArchive{}
		}
		jsonResp(w, items)
	})

	webServer.RegisterAdminHandler("/api/system/backups/archive", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if archiveManager == nil {
			http.Error(w, "archive manager hazir degil", http.StatusServiceUnavailable)
			return
		}
		var req struct {
			Name string `json:"name"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item, err := archiveManager.ArchiveBackup(r.Context(), req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResp(w, map[string]interface{}{"success": true, "item": item})
	})

	webServer.RegisterAdminHandler("/api/system/backups/archive/restore", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if archiveManager == nil {
			http.Error(w, "archive manager hazir degil", http.StatusServiceUnavailable)
			return
		}
		var req struct {
			Name string `json:"name"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		item, err := archiveManager.RestoreBackupArchive(r.Context(), req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResp(w, map[string]interface{}{"success": true, "item": item})
	})

	webServer.RegisterAdminHandler("/api/system/backups/archive/sync", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if archiveManager == nil {
			http.Error(w, "archive manager hazir degil", http.StatusServiceUnavailable)
			return
		}
		uploaded, err := archiveManager.SyncPendingBackups(context.Background(), cfg.GetInt("backup_archive_batch_size", 2))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResp(w, map[string]interface{}{"success": true, "uploaded": uploaded})
	})

	webServer.RegisterAdminHandler("/api/system/backups/download/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := filepath.Base(strings.TrimPrefix(r.URL.Path, "/api/system/backups/download/"))
		path := systemutil.BackupFilePath(dataDir, name)
		if _, err := filepath.Abs(path); err != nil {
			http.Error(w, "Invalid backup path", http.StatusBadRequest)
			return
		}
		if _, err := os.Stat(path); err != nil {
			http.Error(w, "Backup bulunamadi", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(name)))
		http.ServeFile(w, r, path)
	})

	webServer.RegisterAdminHandler("/api/system/backups/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Name string `json:"name"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name := filepath.Base(strings.TrimSpace(req.Name))
		if name == "" || name == "." {
			http.Error(w, "Backup adi gerekli", http.StatusBadRequest)
			return
		}
		if err := systemutil.DeleteBackup(dataDir, name); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if archiveManager != nil {
			_ = archiveManager.MarkBackupLocalDeleted(name, true)
		}
		jsonResp(w, map[string]interface{}{"success": true})
	})

	webServer.RegisterAdminHandler("/api/system/backups/", func(w http.ResponseWriter, r *http.Request) {
		name := filepath.Base(strings.TrimPrefix(r.URL.Path, "/api/system/backups/"))
		if name == "" || name == "." {
			http.Error(w, "Backup adi gerekli", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodDelete:
			if err := systemutil.DeleteBackup(dataDir, name); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if archiveManager != nil {
				_ = archiveManager.MarkBackupLocalDeleted(name, true)
			}
			jsonResp(w, map[string]interface{}{"success": true})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	_ = cfg
}

func serviceCommand(goos, unit, action string) string {
	if goos == "linux" {
		return fmt.Sprintf("sudo systemctl %s %s", action, unit)
	}
	if goos == "windows" {
		return fmt.Sprintf("%s service %s", windowsServiceUnit, action)
	}
	return fmt.Sprintf("%s %s", unit, action)
}

func backupRestoreCommand(goos, unit, binaryPath string) string {
	if goos == "linux" {
		return fmt.Sprintf("sudo systemctl stop %s && sudo %s backup restore fluxstream-backup-YYYYMMDD-HHMMSS.tar.gz && sudo systemctl start %s", unit, binaryPath, unit)
	}
	if goos == "windows" {
		return fmt.Sprintf("%s backup restore fluxstream-backup-YYYYMMDD-HHMMSS.tar.gz", filepath.Base(binaryPath))
	}
	return fmt.Sprintf("%s backup restore fluxstream-backup-YYYYMMDD-HHMMSS.tar.gz", binaryPath)
}

func upgradeCommand(goos, unit, binaryPath, nextBinary string) string {
	if goos == "linux" {
		return fmt.Sprintf("sudo systemctl stop %s && sudo mv -f %s %s && sudo chmod 755 %s && sudo systemctl start %s", unit, nextBinary, binaryPath, binaryPath, unit)
	}
	if goos == "windows" {
		return fmt.Sprintf("Stop-Service %s; Move-Item -Force %s %s; Start-Service %s", unit, nextBinary, binaryPath, unit)
	}
	return fmt.Sprintf("mv -f %s %s", nextBinary, binaryPath)
}

func serviceUnitName() string {
	if runtime.GOOS == "linux" {
		return linuxServiceUnit
	}
	if runtime.GOOS == "windows" {
		return windowsServiceUnit
	}
	return AppName
}
