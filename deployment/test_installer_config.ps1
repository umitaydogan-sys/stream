param(
  [string]$PackageDir = ".\dist\fluxstream-windows-amd64-service",
  [string]$TestRoot = ".\dist\installer-config-test",
  [int]$HTTPPort = 9911,
  [int]$RTMPPort = 2935,
  [int]$HTTPSPort = 9443,
  [int]$RTMPSPort = 2940
)

$ErrorActionPreference = 'Stop'

$packagePath = (Resolve-Path $PackageDir).Path
$testPath = Join-Path (Resolve-Path '.\dist').Path (Split-Path $TestRoot -Leaf)

if (Test-Path $testPath) {
  Remove-Item $testPath -Recurse -Force
}

Copy-Item $packagePath $testPath -Recurse
Push-Location $testPath

try {
  $configOutput = & .\fluxstream.exe config set `
    "http_port=$HTTPPort" `
    "embed_http_port=$HTTPPort" `
    "rtmp_port=$RTMPPort" `
    "https_port=$HTTPSPort" `
    "embed_https_port=$HTTPSPort" `
    "rtmps_port=$RTMPSPort" `
    "server_name=InstallerTest" `
    "setup_completed=true" 2>&1

  $stdoutLog = Join-Path $testPath 'config-test.out.log'
  $stderrLog = Join-Path $testPath 'config-test.err.log'
  $proc = Start-Process -FilePath .\fluxstream.exe -WorkingDirectory $testPath -RedirectStandardOutput $stdoutLog -RedirectStandardError $stderrLog -PassThru

  Start-Sleep -Seconds 6

  $health = Invoke-WebRequest -UseBasicParsing -Uri "http://localhost:$HTTPPort/api/health" -TimeoutSec 10

  [pscustomobject]@{
    ConfigOutput = ($configOutput -join "`n")
    HealthStatus = $health.StatusCode
    HealthBody = $health.Content
    DataDirExists = Test-Path (Join-Path $testPath 'data')
    DBExists = Test-Path (Join-Path $testPath 'data\fluxstream.db')
  } | ConvertTo-Json -Compress
}
finally {
  if ($proc -and -not $proc.HasExited) {
    Stop-Process -Id $proc.Id -Force
    Start-Sleep -Seconds 1
  }
  Pop-Location
  if (Test-Path $testPath) {
    Remove-Item $testPath -Recurse -Force
  }
}
