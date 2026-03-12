param(
  [string]$OutputDir = ".\dist\fluxstream-windows-amd64-portable",
  [string]$FFmpegPath = ""
)

$ErrorActionPreference = 'Stop'

Set-Location (Split-Path -Parent $PSScriptRoot)

if (-not $FFmpegPath) {
  $cmd = Get-Command ffmpeg -ErrorAction Stop
  $FFmpegPath = $cmd.Source
}

$ffmpegItem = Get-Item $FFmpegPath -ErrorAction Stop
$ffmpegDir = $ffmpegItem.Directory.FullName

go build -o .\fluxstream.exe .\cmd\fluxstream\

if (Test-Path $OutputDir) {
  Remove-Item $OutputDir -Recurse -Force
}
New-Item -ItemType Directory -Path $OutputDir | Out-Null
New-Item -ItemType Directory -Path (Join-Path $OutputDir 'ffmpeg') | Out-Null

Copy-Item .\fluxstream.exe $OutputDir
Copy-Item .\deployment\fluxstream.ico $OutputDir
Copy-Item (Join-Path $ffmpegDir 'ffmpeg.exe') (Join-Path $OutputDir 'ffmpeg')

Get-ChildItem $ffmpegDir -Filter '*.dll' | ForEach-Object {
  Copy-Item $_.FullName (Join-Path $OutputDir 'ffmpeg')
}

if (Test-Path (Join-Path $ffmpegDir 'ffprobe.exe')) {
  Copy-Item (Join-Path $ffmpegDir 'ffprobe.exe') (Join-Path $OutputDir 'ffmpeg')
}

@"
FluxStream Portable Release

Icerik:
- fluxstream.exe
- ffmpeg\ffmpeg.exe
- ffmpeg\*.dll

Calisma:
1. fluxstream.exe dosyasini calistirin.
2. Uygulama ffmpeg runtime'ini otomatik olarak ./ffmpeg/ffmpeg.exe altinda bulur.
3. data/ klasoru ilk calistirmada otomatik olusur.
"@ | Set-Content -Path (Join-Path $OutputDir 'README.txt') -Encoding UTF8

Write-Host "Portable release hazir:" $OutputDir
