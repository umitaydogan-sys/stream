param(
  [string]$BaseURL = "http://127.0.0.1:8844",
  [string]$RTMPURL = "rtmp://127.0.0.1:1935/live"
)

$ErrorActionPreference = 'Stop'
$ffmpeg = (Get-Command ffmpeg -ErrorAction Stop).Source

$body = @{
  name = 'adaptive-smoke'
  description = 'adaptive smoke'
  record_enabled = $false
  record_format = 'ts'
  policy_json = '{"mode":"balanced","enable_abr":true,"profile_set":"balanced"}'
  output_formats = '["hls","ll_hls","dash","flv","mp4","webm","mp3","aac","ogg","wav","flac","icecast"]'
} | ConvertTo-Json -Compress

$created = Invoke-RestMethod -Method Post -Uri "$BaseURL/api/streams" -ContentType 'application/json' -Body $body
$key = $created.stream.stream_key
$id = $created.stream.id
$pub = Start-Process -FilePath $ffmpeg -ArgumentList @(
  '-re',
  '-f', 'lavfi',
  '-i', 'testsrc=size=1280x720:rate=30',
  '-f', 'lavfi',
  '-i', 'sine=frequency=1000:sample_rate=44100',
  '-c:v', 'libopenh264',
  '-b:v', '2000k',
  '-pix_fmt', 'yuv420p',
  '-c:a', 'aac',
  '-b:a', '128k',
  '-f', 'flv',
  "$RTMPURL/$key"
) -PassThru -WindowStyle Hidden

try {
  Start-Sleep -Seconds 10
  $checks = [ordered]@{}
  foreach ($url in @(
    "$BaseURL/hls/$key/master.m3u8",
    "$BaseURL/play/$key",
    "$BaseURL/embed/$key",
    "$BaseURL/embed/$key?format=mp4",
    "$BaseURL/embed/$key?format=ll_hls"
  )) {
    $resp = curl.exe -I -s --max-time 15 $url
    $status = (($resp | Select-String -Pattern '^HTTP/\S+\s+(\d+)').Matches | Select-Object -First 1)
    $checks[$url] = if ($status) { [int]$status.Groups[1].Value } else { 0 }
  }
  $players = Invoke-RestMethod -Uri "$BaseURL/api/players"
  [pscustomobject]@{
    stream_key        = $key
    checks            = $checks
    player_templates  = @($players).Count
  } | ConvertTo-Json -Depth 6
}
finally {
  if ($pub -and -not $pub.HasExited) {
    Stop-Process -Id $pub.Id -Force
  }
  try {
    Invoke-RestMethod -Method Delete -Uri "$BaseURL/api/streams/$id" | Out-Null
  } catch {}
}
