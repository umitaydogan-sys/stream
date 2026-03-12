param(
  [string]$PackageDir = ".\dist\fluxstream-windows-amd64-service",
  [string]$InstallDir = "$env:ProgramFiles\FluxStream-ServiceTest",
  [string]$ReportFile = ".\dist\windows-service-selftest.log"
)

$ErrorActionPreference = 'Stop'

function Wait-ServiceState($name, $desired, $timeoutSec = 30) {
  $sw = [Diagnostics.Stopwatch]::StartNew()
  while ($sw.Elapsed.TotalSeconds -lt $timeoutSec) {
    try {
      $svc = Get-Service -Name $name -ErrorAction Stop
      if ($svc.Status -eq $desired) { return $true }
    } catch {}
    Start-Sleep -Seconds 1
  }
  return $false
}

function Wait-Health($url, $timeoutSec = 30) {
  $sw = [Diagnostics.Stopwatch]::StartNew()
  while ($sw.Elapsed.TotalSeconds -lt $timeoutSec) {
    try {
      $resp = Invoke-WebRequest -UseBasicParsing $url -TimeoutSec 3
      if ($resp.StatusCode -eq 200) { return $true }
    } catch {}
    Start-Sleep -Milliseconds 750
  }
  return $false
}

function Write-Stage($message) {
  $line = "$(Get-Date -Format s) $message"
  Add-Content -Path $ReportFile -Value $line
}

if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
  $resolvedPackageDir = (Resolve-Path $PackageDir).Path
  $resolvedReportFile = $ReportFile
  if (-not [System.IO.Path]::IsPathRooted($resolvedReportFile)) {
    $resolvedReportFile = Join-Path (Get-Location) $resolvedReportFile
  }
  $proc = Start-Process powershell -Verb RunAs -ArgumentList "-ExecutionPolicy Bypass -File `"$PSCommandPath`" -PackageDir `"$resolvedPackageDir`" -InstallDir `"$InstallDir`" -ReportFile `"$resolvedReportFile`"" -Wait -PassThru
  exit $proc.ExitCode
}

$packageRoot = Resolve-Path $PackageDir
$serviceName = 'FluxStream'
if (-not [System.IO.Path]::IsPathRooted($ReportFile)) {
  $ReportFile = Join-Path (Get-Location) $ReportFile
}
Set-Content -Path $ReportFile -Value ""
Write-Stage "START"

Get-Process fluxstream -ErrorAction SilentlyContinue | Stop-Process -Force
Write-Stage "STOPPED_EXISTING_PROCESS"
try {
  if (Get-Service -Name $serviceName -ErrorAction Stop) {
    throw "FluxStream servisi zaten kurulu. Self-test mevcut kurulumu ezmemek icin durduruldu."
  }
} catch {
  if ($_.Exception.Message -notlike '*cannot find*' -and $_.Exception.Message -notlike '*Bulunamadi*' -and $_.Exception.Message -notlike '*No service*') {
    throw
  }
}
Write-Stage "SERVICE_NOT_PRESENT"

if (Test-Path $InstallDir) {
  Remove-Item $InstallDir -Recurse -Force
}
New-Item -ItemType Directory -Path $InstallDir | Out-Null
Copy-Item (Join-Path $packageRoot '*') $InstallDir -Recurse -Force
Write-Stage "COPIED_PACKAGE"

& (Join-Path $InstallDir 'fluxstream.exe') service install
Write-Stage "SERVICE_INSTALLED"
& (Join-Path $InstallDir 'fluxstream.exe') service start
Write-Stage "SERVICE_START_SENT"

if (-not (Wait-ServiceState $serviceName 'Running' 30)) {
  throw 'Servis RUNNING durumuna gelmedi.'
}
Write-Stage "SERVICE_RUNNING"
if (-not (Wait-Health 'http://localhost:8844/api/health' 30)) {
  throw 'Servis acildi ancak health endpoint cevap vermedi.'
}
Write-Stage "HEALTH_OK"

& (Join-Path $InstallDir 'fluxstream.exe') service stop
Write-Stage "SERVICE_STOP_SENT"
if (-not (Wait-ServiceState $serviceName 'Stopped' 30)) {
  throw 'Servis STOPPED durumuna gelmedi.'
}
Write-Stage "SERVICE_STOPPED"

& (Join-Path $InstallDir 'fluxstream.exe') service uninstall
Write-Stage "SERVICE_UNINSTALLED"
Start-Sleep -Seconds 2
if (Get-Service -Name $serviceName -ErrorAction SilentlyContinue) {
  throw 'Servis uninstall sonrasinda sistemde kaldi.'
}

Remove-Item $InstallDir -Recurse -Force
Write-Stage "REMOVED_INSTALL_DIR"
Write-Host 'Windows service self-test basarili.'
