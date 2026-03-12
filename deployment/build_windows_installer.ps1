param(
  [string]$ServicePackageDir = ".\dist\fluxstream-windows-amd64-service",
  [string]$OutputDir = ".\dist",
  [string]$InnoCompilerPath = ""
)

$ErrorActionPreference = 'Stop'

Set-Location (Split-Path -Parent $PSScriptRoot)

if (-not $InnoCompilerPath) {
  $registryPath = $null
  try {
    $registryPath = (Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\Inno Setup 6_is1' -ErrorAction Stop).'Inno Setup: App Path'
  } catch {}

  $candidates = @(
    (Get-Command iscc.exe -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source -ErrorAction SilentlyContinue),
    $(if ($registryPath) { Join-Path $registryPath 'ISCC.exe' }),
    "$env:LOCALAPPDATA\Programs\Inno Setup 6\ISCC.exe",
    'C:\Program Files (x86)\Inno Setup 6\ISCC.exe',
    'C:\Program Files\Inno Setup 6\ISCC.exe'
  ) | Where-Object { $_ -and (Test-Path $_) }

  if ($candidates.Count -gt 0) {
    $InnoCompilerPath = $candidates[0]
  }
}

if (-not $InnoCompilerPath) {
  throw 'ISCC.exe bulunamadi. Inno Setup 6 kurun veya -InnoCompilerPath verin.'
}

& $InnoCompilerPath "/DSourceDir=$((Resolve-Path $ServicePackageDir).Path)" "/DOutputDir=$((Resolve-Path $OutputDir).Path)" "/DSetupIconPath=$((Resolve-Path .\deployment\fluxstream.ico).Path)" ".\deployment\fluxstream.iss"

Write-Host "Windows installer hazirlandi:" (Join-Path (Resolve-Path $OutputDir).Path 'FluxStream-Setup.exe')
