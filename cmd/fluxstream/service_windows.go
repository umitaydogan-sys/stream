//go:build windows

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const windowsServiceName = "FluxStream"

func handleServiceMode(args []string) (bool, error) {
	if len(args) < 2 || !strings.EqualFold(args[0], "service") {
		return false, nil
	}

	switch strings.ToLower(args[1]) {
	case "install":
		return true, installWindowsService()
	case "uninstall":
		return true, uninstallWindowsService()
	case "start":
		return true, startWindowsService()
	case "stop":
		return true, stopWindowsService()
	case "status":
		return true, statusWindowsService()
	case "run":
		return true, runWindowsService()
	default:
		return true, fmt.Errorf("bilinmeyen service komutu: %s", args[1])
	}
}

func installWindowsService() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	if existing, err := m.OpenService(windowsServiceName); err == nil {
		existing.Close()
		return fmt.Errorf("%s servisi zaten kurulu", windowsServiceName)
	}

	s, err := m.CreateService(windowsServiceName, exePath, mgr.Config{
		DisplayName: windowsServiceName,
		StartType:   mgr.StartAutomatic,
		Description: "FluxStream live streaming server",
	}, "service", "run")
	if err != nil {
		return err
	}
	defer s.Close()

	fmt.Printf("%s servisi kuruldu.\n", windowsServiceName)
	return nil
}

func uninstallWindowsService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(windowsServiceName)
	if err != nil {
		return err
	}
	defer s.Close()

	_, _ = s.Control(svc.Stop)
	time.Sleep(2 * time.Second)
	if err := s.Delete(); err != nil {
		return err
	}

	fmt.Printf("%s servisi silindi.\n", windowsServiceName)
	return nil
}

func startWindowsService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(windowsServiceName)
	if err != nil {
		return err
	}
	defer s.Close()

	if err := s.Start(); err != nil {
		return err
	}

	fmt.Printf("%s servisi baslatildi.\n", windowsServiceName)
	return nil
}

func stopWindowsService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(windowsServiceName)
	if err != nil {
		return err
	}
	defer s.Close()

	if _, err := s.Control(svc.Stop); err != nil {
		return err
	}

	fmt.Printf("%s servisine durdurma komutu gonderildi.\n", windowsServiceName)
	return nil
}

func statusWindowsService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(windowsServiceName)
	if err != nil {
		return err
	}
	defer s.Close()

	st, err := s.Query()
	if err != nil {
		return err
	}

	fmt.Printf("%s durum: %d\n", windowsServiceName, st.State)
	return nil
}

func runWindowsService() error {
	return svc.Run(windowsServiceName, &fluxstreamService{})
}

type fluxstreamService struct{}

func (s *fluxstreamService) Execute(_ []string, req <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	const accepted = svc.AcceptStop | svc.AcceptShutdown

	logFile, err := openWindowsServiceLog()
	if err != nil {
		return false, 1
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(logFile))

	status <- svc.Status{State: svc.StartPending}
	child, err := startWindowsServiceChild(logFile)
	if err != nil {
		log.Printf("[SERVICE] child start failed: %v", err)
		status <- svc.Status{State: svc.Stopped}
		return false, 1
	}

	childDone := make(chan error, 1)
	go func() {
		childDone <- child.Wait()
	}()

	status <- svc.Status{State: svc.Running, Accepts: accepted}

	for {
		select {
		case err := <-childDone:
			if err != nil {
				log.Printf("[SERVICE] child exited with error: %v", err)
				return false, 1
			}
			log.Printf("[SERVICE] child exited")
			return false, 0
		case c := <-req:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				status <- svc.Status{State: svc.StopPending}
				if child.Process != nil {
					_ = child.Process.Kill()
				}
				select {
				case <-childDone:
				case <-time.After(5 * time.Second):
				}
				return false, 0
			default:
			}
		}
	}
}

func openWindowsServiceLog() (*os.File, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	logDir := filepath.Join(filepath.Dir(exePath), "data", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}
	return os.OpenFile(filepath.Join(logDir, "service.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}

func startWindowsServiceChild(logFile *os.File) (*exec.Cmd, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(exePath)
	cmd.Dir = filepath.Dir(exePath)
	cmd.Env = append(os.Environ(), "FLUXSTREAM_NO_BROWSER=1")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return cmd, nil
}
