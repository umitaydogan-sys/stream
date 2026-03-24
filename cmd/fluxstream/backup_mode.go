package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxstream/fluxstream/internal/storage"
	"github.com/fluxstream/fluxstream/internal/systemutil"
)

func handleBackupMode(args []string) (bool, error) {
	if len(args) < 1 || !strings.EqualFold(args[0], "backup") {
		return false, nil
	}
	if len(args) < 2 {
		return true, fmt.Errorf("kullanim: fluxstream backup <create|list|restore>")
	}

	execPath, err := os.Executable()
	if err != nil {
		return true, fmt.Errorf("executable path alinamadi: %w", err)
	}
	dataDir := filepath.Join(filepath.Dir(execPath), "data")
	if err := ensureDataDirs(dataDir); err != nil {
		return true, err
	}

	switch strings.ToLower(args[1]) {
	case "create":
		fs := flag.NewFlagSet("backup create", flag.ExitOnError)
		includeRecordings := fs.Bool("with-recordings", false, "kayit dosyalarini da dahil et")
		fs.Parse(args[2:])
		db, err := storage.NewSQLiteDB(filepath.Join(dataDir, "fluxstream.db"))
		if err != nil {
			return true, err
		}
		defer db.Close()
		item, err := systemutil.CreateBackup(dataDir, db, *includeRecordings)
		if err != nil {
			return true, err
		}
		fmt.Printf("backup olusturuldu: %s (%d bytes)\n", item.Name, item.Size)
		return true, nil
	case "list":
		items, err := systemutil.ListBackups(dataDir)
		if err != nil {
			return true, err
		}
		if len(items) == 0 {
			fmt.Println("backup bulunamadi")
			return true, nil
		}
		for _, item := range items {
			fmt.Printf("%s | %s | %d bytes\n", item.ModTime.Format("2006-01-02 15:04:05"), item.Name, item.Size)
		}
		return true, nil
	case "restore":
		if len(args) < 3 {
			return true, fmt.Errorf("kullanim: fluxstream backup restore <backup.tar.gz>")
		}
		archivePath := args[2]
		if !filepath.IsAbs(archivePath) {
			archivePath = filepath.Join(systemutil.BackupDir(dataDir), archivePath)
		}
		if err := systemutil.RestoreBackup(archivePath, dataDir); err != nil {
			return true, err
		}
		fmt.Printf("backup geri yuklendi: %s\n", archivePath)
		return true, nil
	default:
		return true, fmt.Errorf("bilinmeyen backup komutu: %s", args[1])
	}
}
