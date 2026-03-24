//go:build !linux

package systemutil

import "fmt"

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
	return ServiceStatus{Platform: "non-linux", Manager: "unsupported", Unit: unit, Message: "service status unsupported on this platform"}
}

func RunServiceAction(unit, action string) error {
	return fmt.Errorf("service action unsupported on this platform")
}
