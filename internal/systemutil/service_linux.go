//go:build linux

package systemutil

import (
	"os/exec"
	"strconv"
	"strings"
)

type ServiceStatus struct {
	Platform string `json:"platform"`
	Manager  string `json:"manager"`
	Unit     string `json:"unit"`
	Active   bool   `json:"active"`
	Enabled  bool   `json:"enabled"`
	MainPID  int    `json:"main_pid"`
	Since    string `json:"since,omitempty"`
	Message  string `json:"message,omitempty"`
}

func GetServiceStatus(unit string) ServiceStatus {
	status := ServiceStatus{Platform: "linux", Manager: "systemd", Unit: unit}
	active, _ := exec.Command("systemctl", "is-active", unit).Output()
	enabled, _ := exec.Command("systemctl", "is-enabled", unit).Output()
	status.Active = strings.TrimSpace(string(active)) == "active"
	status.Enabled = strings.TrimSpace(string(enabled)) == "enabled"
	showOut, err := exec.Command("systemctl", "show", unit, "-p", "MainPID", "-p", "ActiveEnterTimestamp", "--value").Output()
	if err == nil {
		parts := strings.Split(strings.TrimSpace(string(showOut)), "\n")
		if len(parts) > 0 {
			status.MainPID, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
		}
		if len(parts) > 1 {
			status.Since = strings.TrimSpace(parts[1])
		}
	}
	if !status.Active {
		status.Message = strings.TrimSpace(string(active))
	}
	return status
}

func RunServiceAction(unit, action string) error {
	switch action {
	case "start", "stop", "restart":
		return exec.Command("systemctl", action, unit).Run()
	default:
		return exec.ErrNotFound
	}
}
