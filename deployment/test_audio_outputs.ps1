param(
  [string]$StreamKey = "",
  [string]$BaseURL = "http://127.0.0.1:8844",
  [string]$RTMPURL = "rtmp://127.0.0.1:1935/live"
)

$ErrorActionPreference = 'Stop'

$ffmpeg = (Get-Command ffmpeg -ErrorAction Stop).Source
$ffprobe = (Get-Command ffprobe -ErrorAction Stop).Source
$createdStreamId = $null

if ([string]::IsNullOrWhiteSpace($StreamKey)) {
  $createBody = @{ name = "audio-test"; description = "audio output test"; record_enabled = $false; record_format = "ts" } | ConvertTo-Json
  $created = Invoke-RestMethod -Method Post -Uri "$BaseURL/api/streams" -ContentType "application/json" -Body $createBody
  $StreamKey = $created.stream.stream_key
  $createdStreamId = $created.stream.id
}

$publish = Start-Process -FilePath $ffmpeg -ArgumentList @(
  '-re',
  '-f', 'lavfi',
  '-i', 'testsrc=size=1280x720:rate=30',
  '-f', 'lavfi',
  '-i', 'sine=frequency=1000:sample_rate=44100',
  '-c:v', 'libopenh264',
  '-b:v', '1500k',
  '-pix_fmt', 'yuv420p',
  '-c:a', 'aac',
  '-b:a', '128k',
  '-f', 'flv',
  "$RTMPURL/$StreamKey"
) -PassThru -WindowStyle Hidden

try {
  $manifest = Join-Path (Resolve-Path '.\data\transcode\hls').Path "$StreamKey\index.m3u8"
  for ($i = 0; $i -lt 40; $i++) {
    if (Test-Path $manifest) { break }
    Start-Sleep -Milliseconds 500
  }
  if (-not (Test-Path $manifest)) {
    throw "live HLS manifest olusmadi: $manifest"
  }
  for ($i = 0; $i -lt 20; $i++) {
    $streamState = Invoke-RestMethod -Method Get -Uri "$BaseURL/api/streams"
    $liveItem = $streamState | Where-Object { $_.stream_key -eq $StreamKey -and $_.status -eq 'live' }
    if ($liveItem) { break }
    Start-Sleep -Milliseconds 500
  }

  $urls = [ordered]@{
    mp3  = "$BaseURL/audio/mp3/$StreamKey/test.mp3"
    aac  = "$BaseURL/audio/aac/$StreamKey/test.aac"
    ogg  = "$BaseURL/audio/ogg/$StreamKey/test.ogg"
    wav  = "$BaseURL/audio/wav/$StreamKey/test.wav"
    flac = "$BaseURL/audio/flac/$StreamKey/test.flac"
  }

  $results = [ordered]@{}
  foreach ($name in $urls.Keys) {
    $probe = & $ffprobe -v error -show_entries format=format_name -of default=nokey=1:noprint_wrappers=1 $urls[$name] 2>$null
    $headerText = & curl.exe -I -s --max-time 15 $urls[$name]
    $statusMatch = $headerText | Select-String -Pattern '^HTTP/\S+\s+(\d{3})'
    $contentType = (($headerText | Select-String -Pattern '^Content-Type:' -CaseSensitive).Line -replace '^Content-Type:\s*', '').Trim()
    $disposition = (($headerText | Select-String -Pattern '^Content-Disposition:' -CaseSensitive).Line -replace '^Content-Disposition:\s*', '').Trim()
    $results[$name] = [pscustomobject]@{
      status       = if ($statusMatch -and $statusMatch.Matches.Count -gt 0) { [int]$statusMatch.Matches[0].Groups[1].Value } else { 0 }
      content_type = $contentType
      disposition  = $disposition
      format       = ($probe -join '')
    }
  }

  $results | ConvertTo-Json -Depth 4 -Compress
}
finally {
  if ($publish -and -not $publish.HasExited) {
    Stop-Process -Id $publish.Id -Force
  }
  if ($createdStreamId) {
    try {
      Invoke-RestMethod -Method Delete -Uri "$BaseURL/api/streams/$createdStreamId" | Out-Null
    } catch {}
  }
}
