package service

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GenerateWindowsService creates a Windows service installer script (NSIS)
func GenerateWindowsService(exePath, serviceName, displayName, description string) string {
	return fmt.Sprintf(`; FluxStream Windows Service Installer (NSIS)
; Build: makensis fluxstream-installer.nsi

!include "MUI2.nsh"

Name "%s"
OutFile "FluxStream-Setup.exe"
InstallDir "$PROGRAMFILES\FluxStream"
RequestExecutionLevel admin

!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_LANGUAGE "English"

Section "Install"
    SetOutPath $INSTDIR
    File /r "dist\*.*"
    
    ; Create service
    nsExec::ExecToLog 'sc create %s binPath= "$INSTDIR\fluxstream.exe" start= auto DisplayName= "%s"'
    nsExec::ExecToLog 'sc description %s "%s"'
    
    ; Start service
    nsExec::ExecToLog 'sc start %s'
    
    ; Create start menu shortcuts
    CreateDirectory "$SMPROGRAMS\FluxStream"
    CreateShortCut "$SMPROGRAMS\FluxStream\FluxStream Dashboard.lnk" "http://localhost:8844"
    CreateShortCut "$SMPROGRAMS\FluxStream\Uninstall.lnk" "$INSTDIR\uninstall.exe"
    
    ; Write uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
    
    ; Add/Remove Programs
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FluxStream" "DisplayName" "%s"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FluxStream" "UninstallString" '"$INSTDIR\uninstall.exe"'
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FluxStream" "DisplayIcon" "$INSTDIR\fluxstream.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FluxStream" "Publisher" "FluxStream"
SectionEnd

Section "Uninstall"
    nsExec::ExecToLog 'sc stop %s'
    nsExec::ExecToLog 'sc delete %s'
    
    RMDir /r "$INSTDIR"
    RMDir /r "$SMPROGRAMS\FluxStream"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\FluxStream"
SectionEnd
`, displayName, serviceName, displayName, serviceName, description, serviceName, displayName, serviceName, serviceName)
}

// GenerateSystemdUnit creates a Linux systemd service unit file
func GenerateSystemdUnit(exePath, user, group, workDir string) string {
	return fmt.Sprintf(`[Unit]
Description=FluxStream Live Streaming Media Server
Documentation=https://github.com/fluxstream/fluxstream
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=%s
Group=%s
ExecStart=%s
WorkingDirectory=%s
Restart=always
RestartSec=5
LimitNOFILE=65536
StandardOutput=journal
StandardError=journal
SyslogIdentifier=fluxstream

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=%s/data
PrivateTmp=yes
ProtectKernelTunables=yes
ProtectControlGroups=yes

[Install]
WantedBy=multi-user.target
`, user, group, exePath, workDir, workDir)
}

// GenerateDebControl creates DEBIAN/control for .deb package
func GenerateDebControl(version, arch string) string {
	return fmt.Sprintf(`Package: fluxstream
Version: %s
Section: video
Priority: optional
Architecture: %s
Maintainer: FluxStream Team <info@fluxstream.io>
Description: FluxStream Live Streaming Media Server
 Zero-dependency live streaming server with RTMP, HLS, DASH,
 WebRTC, and 20+ output formats. Single binary, pure Go.
Depends: libc6
Homepage: https://github.com/fluxstream/fluxstream
`, version, arch)
}

// GenerateBuildScript creates a cross-platform build script
func GenerateBuildScript(version string) string {
	return fmt.Sprintf(`#!/bin/bash
# FluxStream Build Script
# Usage: ./build.sh [version]

VERSION="${1:-%s}"
APP_NAME="fluxstream"
BUILD_DIR="dist"
LDFLAGS="-s -w -X main.Version=$VERSION"

echo "🔧 Building FluxStream v$VERSION..."
echo ""

# Clean
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# Build for all target platforms
PLATFORMS=(
    "windows/amd64"
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="${PLATFORM%%/*}"
    GOARCH="${PLATFORM##*/}"
    OUTPUT="$BUILD_DIR/${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT="${OUTPUT}.exe"
    fi
    
    echo "  📦 Building $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "$OUTPUT" ./cmd/fluxstream/
    
    if [ $? -ne 0 ]; then
        echo "  ❌ Build failed for $GOOS/$GOARCH"
    else
        SIZE=$(du -h "$OUTPUT" | cut -f1)
        echo "  ✅ $OUTPUT ($SIZE)"
    fi
done

echo ""
echo "✅ Build complete! Outputs in $BUILD_DIR/"
ls -la "$BUILD_DIR/"
`, version)
}

// WriteServiceFiles generates and writes all service/installer files to a directory
func WriteServiceFiles(outputDir, version string) error {
	os.MkdirAll(outputDir, 0755)

	exePath := filepath.Join("/opt/fluxstream", "fluxstream")
	if runtime.GOOS == "windows" {
		exePath = filepath.Join("C:\\Program Files\\FluxStream", "fluxstream.exe")
	}

	// Systemd unit
	systemd := GenerateSystemdUnit(exePath, "fluxstream", "fluxstream", "/opt/fluxstream")
	if err := os.WriteFile(filepath.Join(outputDir, "fluxstream.service"), []byte(systemd), 0644); err != nil {
		return fmt.Errorf("write systemd: %w", err)
	}

	// NSIS installer script
	nsis := GenerateWindowsService(exePath, "FluxStream", "FluxStream Media Server", "Live streaming media server")
	if err := os.WriteFile(filepath.Join(outputDir, "fluxstream-installer.nsi"), []byte(nsis), 0644); err != nil {
		return fmt.Errorf("write nsis: %w", err)
	}

	// Debian control
	arch := "amd64"
	if runtime.GOARCH == "arm64" {
		arch = "arm64"
	}
	deb := GenerateDebControl(version, arch)
	debDir := filepath.Join(outputDir, "DEBIAN")
	os.MkdirAll(debDir, 0755)
	if err := os.WriteFile(filepath.Join(debDir, "control"), []byte(deb), 0644); err != nil {
		return fmt.Errorf("write deb control: %w", err)
	}

	// Build script
	buildScript := GenerateBuildScript(version)
	if err := os.WriteFile(filepath.Join(outputDir, "build.sh"), []byte(buildScript), 0755); err != nil {
		return fmt.Errorf("write build script: %w", err)
	}

	return nil
}
