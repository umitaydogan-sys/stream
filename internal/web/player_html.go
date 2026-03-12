package web

import (
	"fmt"
	"strings"
)

// playerHTML is the full-page HLS player
// Format args: streamKey, streamKey, streamKey
const playerHTML = `<!DOCTYPE html>
<html lang="tr">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s - FluxStream Player</title>
<link rel="stylesheet" href="/static/vendor/bootstrap-icons.css">
<script src="/static/vendor/hls.min.js"></script>
<script src="/static/vendor/dash.all.min.js"></script>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{background:#000;color:#fff;font-family:'Inter',-apple-system,sans-serif;display:flex;flex-direction:column;min-height:100vh}
.player-wrap{flex:1;display:flex;align-items:center;justify-content:center;position:relative;background:#111}
video{width:100%%;height:100%%;max-height:100vh;object-fit:contain;background:#000}
.topbar{height:48px;background:#111;display:flex;align-items:center;justify-content:space-between;padding:0 16px;border-bottom:1px solid #222}
.topbar-brand{display:flex;align-items:center;gap:8px;font-weight:600;font-size:14px}
.topbar-brand i{color:#6366f1;font-size:15px}
.topbar-brand span{color:#6366f1}
.badge-live{background:rgba(255,59,59,0.15);color:#ff3b3b;padding:4px 10px;border-radius:20px;font-size:11px;font-weight:700;display:inline-flex;align-items:center;gap:5px}
.badge-live::before{content:'';width:6px;height:6px;border-radius:50%%;background:#ff3b3b;animation:pulse 2s infinite}
@keyframes pulse{0%%,100%%{opacity:1}50%%{opacity:.5}}
.info-bar{height:40px;background:#111;display:flex;align-items:center;justify-content:space-between;padding:0 16px;border-top:1px solid #222;font-size:12px;color:#888}
.offline-msg{text-align:center;color:#666}
.offline-msg h2{font-size:20px;color:#999;margin-bottom:8px}
.offline-msg p{font-size:14px}
</style>
</head>
<body>
<div class="topbar">
  <div class="topbar-brand"><i class="bi bi-lightning-charge-fill"></i><span>FluxStream</span></div>
  <div id="live-badge"></div>
</div>
<div class="player-wrap">
  <video id="video" controls autoplay muted playsinline></video>
  <div id="offline" style="display:none" class="offline-msg">
    <h2>Yayin cevrimdisi</h2>
    <p>Yayin basladiginda otomatik oynatilacak</p>
  </div>
</div>
<div class="info-bar">
  <span id="stream-info">%s</span>
  <span id="viewer-count"></span>
</div>
<script>
const streamKey = '%s';
const queryParams = new URLSearchParams(location.search);
const video = document.getElementById('video');
const offline = document.getElementById('offline');
const badge = document.getElementById('live-badge');
let hls;
let dashPlayer;
let retryTimer = null;

function passthroughURL(url) {
  const next = new URL(url, location.origin);
  ['token','password'].forEach(function(key){
    const val = queryParams.get(key);
    if (val && !next.searchParams.has(key)) next.searchParams.set(key, val);
  });
  return next.toString();
}

function cleanupPlayers() {
  try { if (hls) { hls.destroy(); hls = null; } } catch (e) {}
  try { if (dashPlayer) { dashPlayer.reset(); dashPlayer = null; } } catch (e) {}
  video.pause();
  video.removeAttribute('src');
  video.load();
}

function sourceCandidates() {
  return [
    { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' },
    { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
    { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/index.m3u8'), marker: '#EXTM3U' }
  ];
}

async function probeSource(candidate) {
  try {
    const res = await fetch(candidate.url, { cache: 'no-store' });
    if (!res.ok) return false;
    const body = await res.text();
    return !candidate.marker || body.indexOf(candidate.marker) !== -1;
  } catch (e) {
    return false;
  }
}

async function resolveSource() {
  const candidates = sourceCandidates();
  for (const candidate of candidates) {
    if (await probeSource(candidate)) return candidate;
  }
  return candidates[1];
}

function markReady() {
  if (retryTimer) {
    clearTimeout(retryTimer);
    retryTimer = null;
  }
  video.style.display = 'block';
  offline.style.display = 'none';
  badge.innerHTML = '<span class="badge-live">CANLI</span>';
  video.play().catch(()=>{});
}

function scheduleRetry() {
  if (retryTimer) return;
  cleanupPlayers();
  video.style.display = 'none';
  offline.style.display = 'block';
  badge.innerHTML = '';
  retryTimer = setTimeout(function() {
    retryTimer = null;
    tryPlay();
  }, 3000);
}

function startHLS(url) {
  if (window.Hls && Hls.isSupported()) {
    hls = new Hls({
      liveSyncDurationCount: 4,
      liveMaxLatencyDurationCount: 10,
      maxBufferLength: 20,
      maxMaxBufferLength: 40,
      backBufferLength: 30,
      enableWorker: true,
      lowLatencyMode: false
    });
    hls.loadSource(url);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, markReady);
    hls.on(Hls.Events.ERROR, function(event, data) {
      if (!data || !data.fatal) return;
      if (data.type === Hls.ErrorTypes.NETWORK_ERROR) {
        hls.startLoad();
        return;
      }
      if (data.type === Hls.ErrorTypes.MEDIA_ERROR) {
        hls.recoverMediaError();
        return;
      }
      scheduleRetry();
    });
    return;
  }
  if (video.canPlayType('application/vnd.apple.mpegurl')) {
    video.src = url;
    video.load();
    return;
  }
  scheduleRetry();
}

function startDASH(url) {
  if (!window.dashjs || !window.dashjs.MediaPlayer) {
    startHLS(passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'));
    return;
  }
  dashPlayer = window.dashjs.MediaPlayer().create();
  dashPlayer.updateSettings({ streaming: { lowLatencyEnabled: true } });
  dashPlayer.initialize(video, url, true);
  dashPlayer.on(window.dashjs.MediaPlayer.events.ERROR, function() {
    scheduleRetry();
  });
}

video.addEventListener('loadedmetadata', markReady);
video.addEventListener('canplay', markReady);
video.addEventListener('playing', markReady);
video.addEventListener('error', scheduleRetry);

async function tryPlay() {
  cleanupPlayers();
  const source = await resolveSource();
  if (source.kind === 'dash') {
    startDASH(source.url);
    return;
  }
  startHLS(source.url);
}

tryPlay();
</script>
</body>
</html>`

// embedPlayerHTML is the iframe-embeddable player
// Format args: title, dash key, hls master key, hls media key, hls fallback key
const embedPlayerHTML = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s - FluxStream</title>
<link rel="stylesheet" href="/static/vendor/bootstrap-icons.css">
<script src="/static/vendor/hls.min.js"></script>
<script src="/static/vendor/dash.all.min.js"></script>
<style>
*{margin:0;padding:0;box-sizing:border-box}
html,body{width:100%%;height:100%%;overflow:hidden;background:#000}
video{width:100%%;height:100%%;object-fit:contain}
.offline{position:absolute;top:50%%;left:50%%;transform:translate(-50%%,-50%%);text-align:center;color:#666;font-family:sans-serif}
</style>
</head>
<body>
<video id="video" controls autoplay muted playsinline></video>
<div id="offline" style="display:none" class="offline">
  <p style="font-size:14px">Yayin cevrimdisi</p>
</div>
<script>
const queryParams = new URLSearchParams(location.search);
const video = document.getElementById('video');
const offline = document.getElementById('offline');
let hls;
let dashPlayer;
let retryTimer = null;
const autoplay = !queryParams.has('autoplay') || queryParams.get('autoplay') === '1' || queryParams.get('autoplay') === 'true';
const muted = !queryParams.has('muted') || queryParams.get('muted') === '1' || queryParams.get('muted') === 'true';

function passthroughURL(url) {
  const next = new URL(url, location.origin);
  ['token','password'].forEach(function(key){
    const val = queryParams.get(key);
    if (val && !next.searchParams.has(key)) next.searchParams.set(key, val);
  });
  return next.toString();
}

video.autoplay = autoplay;
video.muted = muted;

function cleanupPlayers() {
  try { if (hls) { hls.destroy(); hls = null; } } catch (e) {}
  try { if (dashPlayer) { dashPlayer.reset(); dashPlayer = null; } } catch (e) {}
  video.pause();
  video.removeAttribute('src');
  video.load();
}

function sourceCandidates() {
  return [
    { kind: 'dash', url: passthroughURL(location.origin + '/dash/%s/manifest.mpd'), marker: '<MPD' },
    { kind: 'hls', url: passthroughURL(location.origin + '/hls/%s/master.m3u8'), marker: '#EXTM3U' },
    { kind: 'hls', url: passthroughURL(location.origin + '/hls/%s/index.m3u8'), marker: '#EXTM3U' }
  ];
}

async function probeSource(candidate) {
  try {
    const res = await fetch(candidate.url, { cache: 'no-store' });
    if (!res.ok) return false;
    const body = await res.text();
    return !candidate.marker || body.indexOf(candidate.marker) !== -1;
  } catch (e) {
    return false;
  }
}

async function resolveSource() {
  const candidates = sourceCandidates();
  for (const candidate of candidates) {
    if (await probeSource(candidate)) return candidate;
  }
  return candidates[1];
}

function markReady() {
  if (retryTimer) {
    clearTimeout(retryTimer);
    retryTimer = null;
  }
  offline.style.display = 'none';
  if (autoplay) video.play().catch(()=>{});
}

function scheduleRetry() {
  if (retryTimer) return;
  cleanupPlayers();
  offline.style.display = 'block';
  retryTimer = setTimeout(function() {
    retryTimer = null;
    tryPlay();
  }, 3000);
}

function startHLS(url) {
  if (window.Hls && Hls.isSupported()) {
    hls = new Hls({
      liveSyncDurationCount: 4,
      liveMaxLatencyDurationCount: 10,
      maxBufferLength: 20,
      maxMaxBufferLength: 40,
      backBufferLength: 30,
      lowLatencyMode: false
    });
    hls.loadSource(url);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, markReady);
    hls.on(Hls.Events.ERROR, function(event, data) {
      if (!data || !data.fatal) return;
      if (data.type === Hls.ErrorTypes.NETWORK_ERROR) {
        hls.startLoad();
        return;
      }
      if (data.type === Hls.ErrorTypes.MEDIA_ERROR) {
        hls.recoverMediaError();
        return;
      }
      scheduleRetry();
    });
    return;
  }
  if (video.canPlayType('application/vnd.apple.mpegurl')) {
    video.src = url;
    video.load();
    return;
  }
  scheduleRetry();
}

function startDASH(url) {
  if (!window.dashjs || !window.dashjs.MediaPlayer) {
    startHLS(passthroughURL(location.origin + '/hls/%s/master.m3u8'));
    return;
  }
  dashPlayer = window.dashjs.MediaPlayer().create();
  dashPlayer.updateSettings({ streaming: { lowLatencyEnabled: true } });
  dashPlayer.initialize(video, url, autoplay);
  dashPlayer.on(window.dashjs.MediaPlayer.events.ERROR, function() {
    scheduleRetry();
  });
}

video.addEventListener('loadedmetadata', markReady);
video.addEventListener('canplay', markReady);
video.addEventListener('playing', markReady);
video.addEventListener('error', scheduleRetry);

async function tryPlay() {
  cleanupPlayers();
  const source = await resolveSource();
  if (source.kind === 'dash') {
    startDASH(source.url);
    return;
  }
  startHLS(source.url);
}
tryPlay();
</script>
</body>
</html>`

const embedAudioHTML = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s</title>
<link rel="stylesheet" href="/static/vendor/bootstrap-icons.css">
<script src="/static/vendor/hls.min.js"></script>
<style>
*{margin:0;padding:0;box-sizing:border-box}
html,body{width:100%%;height:100%%;overflow:hidden;background:transparent}
body{display:flex;align-items:center;justify-content:center;padding:0;margin:0}
audio{width:100%%;height:100%%;max-width:100%%}
#audio-note{display:none}
</style>
</head>
<body>
<audio id="audio" controls playsinline></audio>
<div id="audio-note">%s</div>
<script>
const format = '%s';
const queryParams = new URLSearchParams(location.search);
const autoplay = %t;
const muted = %t;
const audio = document.getElementById('audio');
const note = document.getElementById('audio-note');
let hls;
let switchedToHLS = false;

function passthroughURL(url) {
  const next = new URL(url, location.origin);
  ['token','password'].forEach(function(key){
    const val = queryParams.get(key);
    if (val && !next.searchParams.has(key)) next.searchParams.set(key, val);
  });
  return next.toString();
}

const primaryUrl = passthroughURL('%s');
const hlsCandidates = [
  passthroughURL(location.origin + '/hls/%s/master.m3u8'),
  passthroughURL(location.origin + '/hls/%s/index.m3u8')
];

audio.autoplay = autoplay;
audio.muted = muted;

function setNote(msg){ note.textContent = msg; }

async function resolveHLSFallback(){
  for (const url of hlsCandidates) {
    try {
      const res = await fetch(url, { method: 'HEAD', cache: 'no-store' });
      if (res.ok) return url;
    } catch (e) {}
  }
  return hlsCandidates[hlsCandidates.length - 1];
}

function startHLSFallback(){
  if(switchedToHLS) return;
  switchedToHLS = true;
  setNote('Dogrudan ses cikisi acilamadi, uyumlu HLS ses akisi deneniyor.');
  resolveHLSFallback().then(function(hlsUrl){
    if(window.Hls && Hls.isSupported()){
      if(hls) hls.destroy();
      hls = new Hls({
        liveSyncDurationCount: 4,
        liveMaxLatencyDurationCount: 10,
        maxBufferLength: 20,
        maxMaxBufferLength: 40,
        backBufferLength: 30,
        lowLatencyMode: false
      });
      hls.loadSource(hlsUrl);
      hls.attachMedia(audio);
      hls.on(Hls.Events.MANIFEST_PARSED, function(){
        setNote('Canli ses baglandi.');
        if(autoplay) audio.play().catch(function(){});
      });
      hls.on(Hls.Events.ERROR, function(event,data){
        if(data && data.fatal){
          setNote('Ses akisi simdilik hazir degil.');
        }
      });
      return;
    }
    if(audio.canPlayType && audio.canPlayType('application/vnd.apple.mpegurl')){
      audio.src = hlsUrl;
      audio.load();
      if(autoplay) audio.play().catch(function(){});
      setNote('Canli ses baglandi.');
      return;
    }
    setNote('Bu tarayici ses fallback akisini da oynatamiyor.');
  });
}

let primaryReady = false;
const fallbackTimer = setTimeout(function(){
  if(!primaryReady) startHLSFallback();
}, 2500);

audio.addEventListener('canplay', function(){
  primaryReady = true;
  clearTimeout(fallbackTimer);
  setNote(format.toUpperCase() + ' ses akisi baglandi.');
  if(autoplay) audio.play().catch(function(){});
});
audio.addEventListener('error', function(){
  clearTimeout(fallbackTimer);
  startHLSFallback();
});

audio.src = primaryUrl;
audio.load();
</script>
</body>
</html>`

func renderEmbedHTML(streamKey, format string, autoplay, muted bool) string {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "mp3", "aac", "ogg", "wav", "flac", "icecast":
		return fmt.Sprintf(embedAudioHTML,
			streamKey,
			"",
			format,
			autoplay,
			muted,
			embedAudioPrimaryURL(streamKey, format),
			streamKey,
			streamKey,
		)
	case "ll_hls", "dash", "flv", "mp4", "webm", "hls", "player", "jsapi":
		return renderVideoEmbedHTML(streamKey, format, autoplay, muted)
	default:
		return renderVideoEmbedHTML(streamKey, "hls", autoplay, muted)
	}
}

func renderVideoEmbedHTML(streamKey, format string, autoplay, muted bool) string {
	_ = format
	_ = autoplay
	_ = muted
	return fmt.Sprintf(embedPlayerHTML, streamKey, streamKey, streamKey, streamKey, streamKey)
}

func embedVideoPreviewURL(streamKey, format string) string {
	return "/hls/" + streamKey + "/master.m3u8"
}

func embedVideoStatusLabel(format string) string {
	switch format {
	case "ll_hls":
		return "LL-HLS uyumlu onizleme"
	case "dash":
		return "DASH uyumlu onizleme"
	case "flv":
		return "HTTP-FLV uyumlu onizleme"
	case "mp4":
		return "MP4 uyumlu onizleme"
	case "webm":
		return "WebM uyumlu onizleme"
	default:
		return "Canli onizleme"
	}
}

func embedVideoPrimaryURL(streamKey, format string) string {
	switch format {
	case "ll_hls":
		return "/hls/" + streamKey + "/ll.m3u8"
	case "dash":
		return "/dash/" + streamKey + "/manifest.mpd"
	case "flv":
		return "/flv/" + streamKey
	case "mp4":
		return "/mp4/" + streamKey + "/" + streamKey + ".mp4"
	case "webm":
		return "/webm/" + streamKey + "/" + streamKey + ".webm"
	default:
		return "/hls/" + streamKey + "/master.m3u8"
	}
}

func boolAttr(enabled bool, attr string) string {
	if enabled {
		return attr
	}
	return ""
}

func embedAudioPrimaryURL(streamKey, format string) string {
	switch format {
	case "mp3":
		return "/audio/mp3/" + streamKey + "/" + streamKey + ".mp3"
	case "aac":
		return "/audio/aac/" + streamKey + "/" + streamKey + ".aac"
	case "ogg":
		return "/audio/ogg/" + streamKey + "/" + streamKey + ".ogg"
	case "wav":
		return "/audio/wav/" + streamKey + "/" + streamKey + ".wav"
	case "flac":
		return "/audio/flac/" + streamKey + "/" + streamKey + ".flac"
	case "icecast":
		return "/icecast/" + streamKey
	default:
		return "/audio/aac/" + streamKey
	}
}
