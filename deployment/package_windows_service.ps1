param(
  [string]$OutputDir = ".\dist\fluxstream-windows-amd64-service",
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

@'
param(
  [string]$InstallDir = "$env:ProgramFiles\FluxStream"
)

$ErrorActionPreference = 'Stop'
$PackageRoot = Split-Path -Parent $MyInvocation.MyCommand.Path

if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
  $proc = Start-Process powershell -Verb RunAs -ArgumentList "-ExecutionPolicy Bypass -File `"$PSCommandPath`" -InstallDir `"$InstallDir`"" -Wait -PassThru
  exit $proc.ExitCode
}

if (Test-Path $InstallDir) {
  Write-Host "Guncelleniyor:" $InstallDir
} else {
  New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Copy-Item (Join-Path $PackageRoot '*') $InstallDir -Recurse -Force

& (Join-Path $InstallDir 'fluxstream.exe') service install
& (Join-Path $InstallDir 'fluxstream.exe') service start

Write-Host "FluxStream service kurulumu tamamlandi:" $InstallDir
'@ | Set-Content -Path (Join-Path $OutputDir 'install_service.ps1') -Encoding UTF8

@'
param(
  [string]$InstallDir = "$env:ProgramFiles\FluxStream",
  [switch]$KeepFiles
)

$ErrorActionPreference = 'Stop'

if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
  $keepArg = ''
  if ($KeepFiles) { $keepArg = ' -KeepFiles' }
  $proc = Start-Process powershell -Verb RunAs -ArgumentList "-ExecutionPolicy Bypass -File `"$PSCommandPath`" -InstallDir `"$InstallDir`"$keepArg" -Wait -PassThru
  exit $proc.ExitCode
}

if (Test-Path (Join-Path $InstallDir 'fluxstream.exe')) {
  try { & (Join-Path $InstallDir 'fluxstream.exe') service stop } catch {}
  Start-Sleep -Seconds 2
  try { & (Join-Path $InstallDir 'fluxstream.exe') service uninstall } catch {}
}

if ((-not $KeepFiles) -and (Test-Path $InstallDir)) {
  Remove-Item $InstallDir -Recurse -Force
}

Write-Host "FluxStream service kaldirildi."
'@ | Set-Content -Path (Join-Path $OutputDir 'uninstall_service.ps1') -Encoding UTF8

@"
FluxStream Windows Service Package

Icerik:
- fluxstream.exe
- ffmpeg\ffmpeg.exe
- ffmpeg\*.dll
- install_service.ps1
- uninstall_service.ps1

Kurulum:
1. PowerShell'i yonetici olarak acin.
2. install_service.ps1 komutunu calistirin.
3. Servis adi: FluxStream

Kaldirma:
1. uninstall_service.ps1 komutunu calistirin.
"@ | Set-Content -Path (Join-Path $OutputDir 'README.txt') -Encoding UTF8

Write-Host "Windows service package hazir:" $OutputDir
