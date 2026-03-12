param(
  [string]$OutputDir = ".\dist\fluxstream-linux-amd64-systemd",
  [string]$LinuxFFmpegPath = ""
)

$ErrorActionPreference = 'Stop'

Set-Location (Split-Path -Parent $PSScriptRoot)

if (Test-Path $OutputDir) {
  Remove-Item $OutputDir -Recurse -Force
}

New-Item -ItemType Directory -Path $OutputDir | Out-Null
New-Item -ItemType Directory -Path (Join-Path $OutputDir 'systemd') | Out-Null

$env:GOOS = 'linux'
$env:GOARCH = 'amd64'
$env:CGO_ENABLED = '0'
go build -o (Join-Path $OutputDir 'fluxstream') .\cmd\fluxstream\
Remove-Item Env:GOOS, Env:GOARCH, Env:CGO_ENABLED

if ($LinuxFFmpegPath) {
  New-Item -ItemType Directory -Path (Join-Path $OutputDir 'ffmpeg') | Out-Null
  Copy-Item $LinuxFFmpegPath (Join-Path $OutputDir 'ffmpeg\ffmpeg')
}

@'
[Unit]
Description=FluxStream Live Streaming Server
After=network.target

[Service]
Type=simple
User=fluxstream
Group=fluxstream
WorkingDirectory=/opt/fluxstream
Environment=FLUXSTREAM_NO_BROWSER=1
ExecStart=/opt/fluxstream/fluxstream
Restart=always
RestartSec=2
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
'@ | Set-Content -Path (Join-Path $OutputDir 'systemd\fluxstream.service') -Encoding UTF8

@'
#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR=/opt/fluxstream
PACKAGE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

sudo useradd --system --home "$INSTALL_DIR" --shell /usr/sbin/nologin fluxstream 2>/dev/null || true
sudo mkdir -p "$INSTALL_DIR"
sudo cp -R "$PACKAGE_DIR"/. "$INSTALL_DIR"/
sudo chown -R fluxstream:fluxstream "$INSTALL_DIR"
sudo chmod +x "$INSTALL_DIR/fluxstream"
if [ -f "$INSTALL_DIR/ffmpeg/ffmpeg" ]; then
  sudo chmod +x "$INSTALL_DIR/ffmpeg/ffmpeg"
fi
sudo cp "$INSTALL_DIR/systemd/fluxstream.service" /etc/systemd/system/fluxstream.service
sudo systemctl daemon-reload
sudo systemctl enable --now fluxstream
echo "FluxStream systemd kurulumu tamamlandi."
'@ | Set-Content -Path (Join-Path $OutputDir 'install.sh') -Encoding UTF8

@'
#!/usr/bin/env bash
set -euo pipefail

sudo systemctl disable --now fluxstream || true
sudo rm -f /etc/systemd/system/fluxstream.service
sudo systemctl daemon-reload
sudo rm -rf /opt/fluxstream
echo "FluxStream systemd kurulumu kaldirildi."
'@ | Set-Content -Path (Join-Path $OutputDir 'uninstall.sh') -Encoding UTF8

@"
FluxStream Linux systemd Package

Icerik:
- fluxstream
- systemd/fluxstream.service
- install.sh
- uninstall.sh

Not:
- Linux ffmpeg binary saglanirsa ./ffmpeg/ffmpeg olarak pakete eklenir.
- Saglanmazsa sistem PATH icindeki ffmpeg kullanilir.
"@ | Set-Content -Path (Join-Path $OutputDir 'README.txt') -Encoding UTF8

Write-Host "Linux systemd package hazir:" $OutputDir
