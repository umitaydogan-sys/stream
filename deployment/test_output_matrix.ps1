param(
  [string]$BaseURL = "http://127.0.0.1:8844",
  [string]$RTMPURL = "rtmp://127.0.0.1:1935/live"
)
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $false
$ffmpeg = (Get-Command ffmpeg -ErrorAction Stop).Source
$ffprobe = (Get-Command ffprobe -ErrorAction Stop).Source
$body = @{ name = 'matrix-test'; description = 'matrix'; record_enabled = $false; record_format = 'ts' } | ConvertTo-Json
$created = Invoke-RestMethod -Method Post -Uri "$BaseURL/api/streams" -ContentType 'application/json' -Body $body
$key = $created.stream.stream_key
$id = $created.stream.id
$pub = Start-Process -FilePath $ffmpeg -ArgumentList @('-re','-f','lavfi','-i','testsrc=size=1280x720:rate=30','-f','lavfi','-i','sine=frequency=1000:sample_rate=44100','-c:v','libopenh264','-b:v','1800k','-pix_fmt','yuv420p','-c:a','aac','-b:a','128k','-f','flv',"$RTMPURL/$key") -PassThru -WindowStyle Hidden
try {
  Start-Sleep -Seconds 8
  $headers = @{ 'User-Agent' = 'Mozilla/5.0 FluxStreamMatrixTest' }
  $null = Invoke-WebRequest -Headers $headers -Uri "$BaseURL/hls/$key/index.m3u8" -UseBasicParsing
  $seg = (Invoke-WebRequest -Headers $headers -Uri "$BaseURL/hls/$key/index.m3u8" -UseBasicParsing).Content | Select-String -Pattern '^[^#].+\.ts$' | Select-Object -First 1
  if ($seg) { $segUrl = "$BaseURL/hls/$key/" + $seg.Matches[0].Value.Trim(); $null = Invoke-WebRequest -Headers $headers -Uri $segUrl -UseBasicParsing }
  for ($i = 0; $i -lt 20; $i++) {
    try {
      $null = Invoke-WebRequest -Headers $headers -Uri "$BaseURL/dash/$key/manifest.mpd" -UseBasicParsing
      break
    } catch {
      Start-Sleep -Milliseconds 500
    }
  }
  try { $null = Invoke-WebRequest -Headers $headers -Uri "$BaseURL/embed/$key?format=mp4" -UseBasicParsing } catch {}
  $viewer = Invoke-RestMethod -Uri "$BaseURL/api/viewers"
  $analytics = Invoke-RestMethod -Uri "$BaseURL/api/analytics"
  $targets = [ordered]@{
    hls = "$BaseURL/hls/$key/index.m3u8"
    ll_hls = "$BaseURL/hls/$key/ll.m3u8"
    dash = "$BaseURL/dash/$key/manifest.mpd"
    flv = "$BaseURL/flv/$key"
    mp4 = "$BaseURL/mp4/$key/test.mp4"
    webm = "$BaseURL/webm/$key/test.webm"
    mp3 = "$BaseURL/audio/mp3/$key/test.mp3"
    aac = "$BaseURL/audio/aac/$key/test.aac"
    ogg = "$BaseURL/audio/ogg/$key/test.ogg"
    wav = "$BaseURL/audio/wav/$key/test.wav"
    flac = "$BaseURL/audio/flac/$key/test.flac"
    icecast = "$BaseURL/icecast/$key"
  }
  $results = [ordered]@{}
  foreach ($name in $targets.Keys) {
    $probe = cmd.exe /c "`"$ffprobe`" -v error -show_entries format=format_name -of default=nokey=1:noprint_wrappers=1 `"$($targets[$name])`" 2>nul"
    $probeExit = $LASTEXITCODE
    $head = & curl.exe -I -s --max-time 15 $targets[$name]
    $status = (($head | Select-String -Pattern '^HTTP/\S+\s+(\d+)').Matches | Select-Object -First 1)
    $type = (($head | Select-String -Pattern '^Content-Type:' -CaseSensitive).Line -replace '^Content-Type:\s*','').Trim()
    $results[$name] = [pscustomobject]@{
      status = if($status){[int]$status.Groups[1].Value}else{0}
      content_type = $type
      format = ($probe -join '')
      probe_ok = ($probeExit -eq 0)
    }
  }
  [pscustomobject]@{ viewers=$viewer; analytics=$analytics; outputs=$results } | ConvertTo-Json -Depth 8
}
finally {
  if ($pub -and -not $pub.HasExited) { Stop-Process -Id $pub.Id -Force }
  try { Invoke-RestMethod -Method Delete -Uri "$BaseURL/api/streams/$id" | Out-Null } catch {}
}
