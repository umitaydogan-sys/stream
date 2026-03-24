//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func handleServiceMode(args []string) (bool, error) {
	if len(args) < 2 || !strings.EqualFold(args[0], "service") {
		return false, nil
	}

	switch strings.ToLower(strings.TrimSpace(args[1])) {
	case "install":
		return true, installLinuxService()
	case "uninstall":
		return true, uninstallLinuxService()
	case "start":
		return true, runLinuxSystemctl("start", linuxServiceUnit)
	case "stop":
		return true, runLinuxSystemctl("stop", linuxServiceUnit)
	case "restart":
		return true, runLinuxSystemctl("restart", linuxServiceUnit)
	case "status":
		return true, linuxServiceStatus()
	default:
		return true, fmt.Errorf("bilinmeyen service komutu: %s", args[1])
	}
}

func installLinuxService() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("linux service install root yetkisi gerektirir")
	}
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath = filepath.Clean(exePath)
	workDir := filepath.Dir(exePath)
	serviceUser := envOrDefault("FLUXSTREAM_SERVICE_USER", "fluxstream")
	serviceGroup := envOrDefault("FLUXSTREAM_SERVICE_GROUP", serviceUser)
	unitPath := filepath.Join("/etc/systemd/system", linuxServiceUnit+".service")

	unitBody := fmt.Sprintf(`[Unit]
Description=FluxStream live streaming server
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=%s
Group=%s
WorkingDirectory=%s
ExecStart=%s
Restart=always
RestartSec=3
AmbientCapabilities=CAP_NET_BIND_SERVICE
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
NoNewPrivileges=true
LimitNOFILE=65535
Environment=FLUXSTREAM_NO_BROWSER=1

[Install]
WantedBy=multi-user.target
`, serviceUser, serviceGroup, workDir, exePath)

	if err := os.WriteFile(unitPath, []byte(unitBody), 0644); err != nil {
		return err
	}
	if err := runLinuxSystemctl("daemon-reload"); err != nil {
		return err
	}
	if err := runLinuxSystemctl("enable", linuxServiceUnit); err != nil {
		return err
	}
	_ = exec.Command("setcap", "cap_net_bind_service=+ep", exePath).Run()
	if err := runLinuxSystemctl("restart", linuxServiceUnit); err != nil {
		return err
	}

	fmt.Printf("%s servisi kuruldu: %s\n", linuxServiceUnit, unitPath)
	return nil
}

func uninstallLinuxService() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("linux service uninstall root yetkisi gerektirir")
	}
	unitPath := filepath.Join("/etc/systemd/system", linuxServiceUnit+".service")
	_ = runLinuxSystemctl("stop", linuxServiceUnit)
	_ = runLinuxSystemctl("disable", linuxServiceUnit)
	if err := os.Remove(unitPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := runLinuxSystemctl("daemon-reload"); err != nil {
		return err
	}
	fmt.Printf("%s servisi silindi.\n", linuxServiceUnit)
	return nil
}

func linuxServiceStatus() error {
	cmd := exec.Command("systemctl", "status", linuxServiceUnit, "--no-pager")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runLinuxSystemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		text := strings.TrimSpace(string(out))
		if text == "" {
			return err
		}
		return fmt.Errorf(text)
	}
	if text := strings.TrimSpace(string(out)); text != "" {
		fmt.Println(text)
	}
	return nil
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
