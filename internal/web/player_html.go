package web

import (
	"fmt"
	"strings"
)

func normalizePlayerFormat(format string) string {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "", "iframe", "player", "jsapi":
		return "auto"
	case "auto", "hls", "ll_hls", "dash", "flv", "mp4", "webm", "mp3", "aac", "ogg", "wav", "flac", "icecast":
		return format
	default:
		return "auto"
	}
}

func isAudioFormat(format string) bool {
	switch normalizePlayerFormat(format) {
	case "mp3", "aac", "ogg", "wav", "flac", "icecast":
		return true
	default:
		return false
	}
}

func renderPlayerHTML(streamKey, format string, autoplay, muted bool) string {
	format = normalizePlayerFormat(format)
	if isAudioFormat(format) {
		return renderAudioHTML(streamKey, format, autoplay, muted, false)
	}
	return renderVideoHTML(streamKey, format, autoplay, muted, false)
}

func renderEmbedHTML(streamKey, format string, autoplay, muted bool) string {
	format = normalizePlayerFormat(format)
	if isAudioFormat(format) {
		return renderAudioHTML(streamKey, format, autoplay, muted, true)
	}
	return renderVideoHTML(streamKey, format, autoplay, muted, true)
}

func renderVideoHTML(streamKey, format string, autoplay, muted bool, embedded bool) string {
	stageClass := "player-stage player-shell"
	bodyClass := "player-body"
	header := ""
	footer := ""
	if embedded {
		stageClass += " embedded"
		bodyClass += " embedded"
	} else {
		header = `<div class="topbar">
  <div class="topbar-brand">
    <div class="topbar-title"><strong id="player-title">` + streamKey + `</strong><small>FluxStream Player</small></div>
  </div>
  <div id="live-badge"></div>
</div>`
		footer = `<div class="info-bar">
  <span id="stream-info">` + streamKey + `</span>
  <span id="viewer-count"></span>
</div>`
	}

	return fmt.Sprintf(`<!DOCTYPE html>
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
html,body{width:100%%;height:100%%}
body.player-body{background:#020617;color:#fff;font-family:'Segoe UI',Tahoma,Geneva,Verdana,sans-serif;display:flex;flex-direction:column;min-height:100vh}
body.player-body.embedded{overflow:hidden}
.topbar{height:52px;background:#0f172a;display:flex;align-items:center;justify-content:space-between;padding:0 16px;border-bottom:1px solid rgba(148,163,184,.18)}
.topbar-brand{display:flex;align-items:center;gap:10px;min-width:0}
.topbar-title{display:flex;flex-direction:column;min-width:0}
.topbar-title strong{font-size:14px;color:#fff;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.topbar-title small{font-size:11px;color:#94a3b8}
.badge-live{display:inline-flex;align-items:center;gap:6px;padding:6px 10px;border-radius:999px;background:rgba(239,68,68,.18);font-size:10px;font-weight:700;letter-spacing:.08em;text-transform:uppercase;color:#fff}
.badge-live::before{content:'';width:6px;height:6px;border-radius:50%%;background:#fb7185;animation:pulse 2s infinite}
.player-stage{position:relative;flex:1;display:flex;align-items:center;justify-content:center;background:#020617;overflow:hidden}
.player-stage.embedded{width:100%%;height:100%%;min-height:100%%}
video{width:100%%;height:100%%;object-fit:contain;background:#000}
.floating-logo{position:absolute;z-index:5;display:none;pointer-events:none}
.floating-logo img{height:36px;max-width:160px;object-fit:contain}
.player-watermark{position:absolute;left:16px;bottom:16px;padding:6px 10px;border-radius:999px;background:rgba(15,23,42,.48);backdrop-filter:blur(8px);font-size:11px;font-weight:700;letter-spacing:.08em;color:rgba(255,255,255,.78);z-index:4;display:none}
.resume-button{position:absolute;left:50%%;top:50%%;transform:translate(-50%%,-50%%);display:none;align-items:center;gap:10px;padding:14px 22px;border:0;border-radius:999px;background:rgba(37,99,235,.88);color:#fff;font-size:14px;font-weight:700;cursor:pointer;z-index:6;box-shadow:0 18px 36px rgba(15,23,42,.28)}
.offline-msg{position:absolute;left:50%%;top:50%%;transform:translate(-50%%,-50%%);text-align:center;color:#cbd5e1;z-index:3;display:none;padding:18px 20px;border-radius:18px;background:rgba(15,23,42,.44);backdrop-filter:blur(10px)}
.offline-msg h2{font-size:20px;margin-bottom:8px}
.offline-msg p{font-size:13px;color:#94a3b8}
.info-bar{height:40px;background:#0f172a;display:flex;align-items:center;justify-content:space-between;padding:0 16px;border-top:1px solid rgba(148,163,184,.18);font-size:12px;color:#94a3b8}
@keyframes pulse{0%%,100%%{opacity:1}50%%{opacity:.5}}
</style>
</head>
<body class="%s">
%s
<div class="%s">
  <div id="player-floating-logo" class="floating-logo"></div>
  <video id="video" controls playsinline></video>
  <button id="resume-button" class="resume-button" type="button"><i class="bi bi-play-fill"></i><span>Play</span></button>
  <div id="player-watermark" class="player-watermark"></div>
  <div id="offline" class="offline-msg">
    <h2>Yayin cevrimdisi</h2>
    <p>Yayin basladiginda otomatik oynatilacak</p>
  </div>
</div>
%s
<script>
const streamKey = %q;
const queryParams = new URLSearchParams(location.search);
const preferredFormat = (queryParams.get('format') || %q).toLowerCase();
const autoplay = queryParams.has('autoplay') ? (queryParams.get('autoplay') === '1' || queryParams.get('autoplay') === 'true') : %t;
const muted = queryParams.has('muted') ? (queryParams.get('muted') === '1' || queryParams.get('muted') === 'true') : %t;
const video = document.getElementById('video');
const offline = document.getElementById('offline');
const badge = document.getElementById('live-badge');
const playerTitle = document.getElementById('player-title');
const logoBox = document.getElementById('player-floating-logo');
const watermark = document.getElementById('player-watermark');
const resumeButton = document.getElementById('resume-button');
let hls;
let dashPlayer;
let retryTimer = null;
let activeSourceKind = '';
let sourceOverride = '';
let watchdogTimer = null;
let lastProgressAt = 0;
let lastCurrentTime = 0;
let stallRecoveries = 0;

video.autoplay = autoplay;
video.muted = muted;

function passthroughURL(url) {
  const next = new URL(url, location.origin);
  ['token','password'].forEach(function(key) {
    const val = queryParams.get(key);
    if (val && !next.searchParams.has(key)) next.searchParams.set(key, val);
  });
  return next.toString();
}

function applyStyleText(node, value, fallbackProp) {
  if (!node || !value) return;
  if (value.indexOf(':') !== -1 || value.indexOf(';') !== -1) {
    node.style.cssText += ';' + value;
    return;
  }
  if (fallbackProp) node.style[fallbackProp] = value;
}

function logoPositionStyle(position) {
  switch ((position || 'top-right').toLowerCase()) {
    case 'top-left':
      return 'top:16px;left:16px;';
    case 'bottom-left':
      return 'bottom:16px;left:16px;';
    case 'bottom-right':
      return 'bottom:16px;right:16px;';
    default:
      return 'top:16px;right:16px;';
  }
}

function showResumeButton() {
  if (resumeButton) resumeButton.style.display = 'inline-flex';
}

function hideResumeButton() {
  if (resumeButton) resumeButton.style.display = 'none';
}

function tryAutoplay() {
  if (!autoplay) {
    showResumeButton();
    return;
  }
  const playPromise = video.play();
  if (playPromise && playPromise.catch) {
    playPromise.catch(function() { showResumeButton(); });
  }
}

function applySkin() {
  const title = queryParams.get('player_title') || streamKey;
  const bg = queryParams.get('player_bg') || '';
  const controls = queryParams.get('player_controls') || '';
  const playStyle = queryParams.get('player_play') || '';
  const logo = queryParams.get('player_logo') || '';
  const logoPosition = queryParams.get('player_logo_position') || 'top-right';
  const logoOpacity = queryParams.get('player_logo_opacity') || '1';
  const watermarkText = queryParams.get('player_watermark') || '';
  const showTitle = queryParams.get('player_show_title') !== '0';
  const showBadge = queryParams.get('player_show_badge') !== '0';
  const customCSS = queryParams.get('player_custom_css') || '';
  if (playerTitle) {
    playerTitle.textContent = title;
    if (!showTitle) {
      const topbar = document.querySelector('.topbar');
      if (topbar) topbar.style.display = 'none';
    }
  }
  if (badge && !showBadge) badge.style.display = 'none';
  if (logo && logoBox) {
    logoBox.style.display = 'block';
    logoBox.style.cssText += ';' + logoPositionStyle(logoPosition);
    logoBox.innerHTML = '<img src="' + logo + '" alt="logo" style="opacity:' + logoOpacity + '">';
  }
  if (watermarkText && watermark) {
    watermark.style.display = 'block';
    watermark.textContent = watermarkText;
  }
  applyStyleText(document.body, bg, 'background');
  applyStyleText(document.querySelector('.player-stage'), bg, 'background');
  applyStyleText(document.querySelector('.topbar'), controls);
  applyStyleText(document.querySelector('.info-bar'), controls);
  applyStyleText(resumeButton, playStyle);
  if (customCSS) {
    const style = document.createElement('style');
    style.textContent = customCSS;
    document.head.appendChild(style);
  }
}

function sourceCatalog() {
  const catalog = {
    auto: [
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/index.m3u8'), marker: '#EXTM3U' },
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' }
    ],
    hls: [
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/index.m3u8'), marker: '#EXTM3U' },
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' }
    ],
    hls_media: [
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/index.m3u8'), marker: '#EXTM3U' },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' }
    ],
    ll_hls: [
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/ll.m3u8'), marker: '#EXTM3U' },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/index.m3u8'), marker: '#EXTM3U' },
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' }
    ],
    dash: [
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/index.m3u8'), marker: '#EXTM3U' }
    ],
    mp4: [
      { kind: 'native', mime: 'video/mp4', url: passthroughURL(location.origin + '/mp4/' + streamKey + '/' + streamKey + '.mp4') },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' }
    ],
    webm: [
      { kind: 'native', mime: 'video/webm', url: passthroughURL(location.origin + '/webm/' + streamKey + '/' + streamKey + '.webm') },
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' }
    ]
  };
  catalog.flv = catalog.hls.slice();
  return catalog[sourceOverride] || catalog[preferredFormat] || catalog.auto;
}

async function probeSource(candidate) {
  try {
    const res = await fetch(candidate.url, { cache: 'no-store' });
    if (!res.ok) return false;
    if (candidate.kind === 'native') {
      if (res.body && res.body.cancel) res.body.cancel().catch(function() {});
      return true;
    }
    const body = await res.text();
    return !candidate.marker || body.indexOf(candidate.marker) !== -1;
  } catch (e) {
    return false;
  }
}

async function resolveSource() {
  const candidates = sourceCatalog();
  for (const candidate of candidates) {
    if (await probeSource(candidate)) return candidate;
  }
  return candidates[0];
}

function cleanupPlayers() {
  try { if (hls) { hls.destroy(); hls = null; } } catch (e) {}
  try { if (dashPlayer) { dashPlayer.reset(); dashPlayer = null; } } catch (e) {}
  activeSourceKind = '';
  hideResumeButton();
  video.pause();
  video.removeAttribute('src');
  video.load();
}

function noteProgress() {
  lastProgressAt = Date.now();
  lastCurrentTime = video.currentTime || 0;
  stallRecoveries = 0;
}

async function trySourceFallback(nextOverride) {
  if (sourceOverride === nextOverride) {
    scheduleRetry();
    return;
  }
  sourceOverride = nextOverride;
  await tryPlay();
}

async function tryHLSMediaFallback() {
  await trySourceFallback('hls_media');
}

async function tryHLSMasterFallback() {
  await trySourceFallback('hls');
}

function ensurePlaybackWatchdog() {
  if (watchdogTimer) return;
  watchdogTimer = setInterval(function() {
    if (!video || video.paused || video.ended) return;
    if (offline && offline.style.display !== 'none') return;
    const now = Date.now();
    const currentTime = video.currentTime || 0;
    if (currentTime > lastCurrentTime + 0.15) {
      noteProgress();
      return;
    }
    if (now-lastProgressAt < 6500) return;
    stallRecoveries += 1;
    lastProgressAt = now;
    if (hls && activeSourceKind === 'hls' && stallRecoveries === 1) {
      try { hls.startLoad(); } catch (e) {}
      try { video.play().catch(function() {}); } catch (e) {}
      return;
    }
    if (activeSourceKind === 'dash' || preferredFormat === 'dash') {
      tryHLSMasterFallback();
      return;
    }
    tryHLSMediaFallback();
  }, 2500);
}

function markReady() {
  if (retryTimer) {
    clearTimeout(retryTimer);
    retryTimer = null;
  }
  video.style.display = 'block';
  if (offline) offline.style.display = 'none';
  if (badge && badge.style.display !== 'none') badge.innerHTML = '<span class="badge-live">CANLI</span>';
  noteProgress();
  ensurePlaybackWatchdog();
  tryAutoplay();
}

function scheduleRetry() {
  if (retryTimer) return;
  video.style.display = 'none';
  if (offline) offline.style.display = 'block';
  if (badge) badge.innerHTML = '';
  cleanupPlayers();
  retryTimer = setTimeout(function() {
    retryTimer = null;
    tryPlay();
  }, 3000);
}

function startNative(url) {
  activeSourceKind = 'native';
  video.src = url;
  video.load();
}

function startHLS(url) {
  if (window.Hls && Hls.isSupported()) {
    activeSourceKind = 'hls';
    hls = new Hls({
      liveSyncDurationCount: preferredFormat === 'll_hls' ? 3 : 4,
      liveMaxLatencyDurationCount: preferredFormat === 'll_hls' ? 6 : 10,
      maxBufferLength: preferredFormat === 'll_hls' ? 12 : 20,
      maxMaxBufferLength: preferredFormat === 'll_hls' ? 20 : 40,
      backBufferLength: 30,
      startLevel: 0,
      abrEwmaDefaultEstimate: preferredFormat === 'll_hls' ? 450000 : 300000,
      abrBandWidthFactor: 0.7,
      abrBandWidthUpFactor: 0.5,
      capLevelToPlayerSize: true,
      capLevelOnFPSDrop: true,
      enableWorker: true,
      lowLatencyMode: preferredFormat === 'll_hls'
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
      tryHLSMediaFallback();
    });
    return;
  }
  if (video.canPlayType('application/vnd.apple.mpegurl')) {
    activeSourceKind = 'hls';
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
  activeSourceKind = 'dash';
  dashPlayer = window.dashjs.MediaPlayer().create();
  dashPlayer.updateSettings({ streaming: { lowLatencyEnabled: true } });
  dashPlayer.on(window.dashjs.MediaPlayer.events.STREAM_INITIALIZED, markReady);
  dashPlayer.initialize(video, url, autoplay);
  dashPlayer.on(window.dashjs.MediaPlayer.events.ERROR, function() {
    tryHLSMasterFallback();
  });
}

async function tryPlay() {
  cleanupPlayers();
  const source = await resolveSource();
  if (source.kind === 'dash') {
    startDASH(source.url);
    return;
  }
  if (source.kind === 'native') {
    startNative(source.url);
    return;
  }
  startHLS(source.url);
}

if (resumeButton) {
  resumeButton.addEventListener('click', function() {
    video.play().then(function() {
      hideResumeButton();
    }).catch(function() {});
  });
}

video.addEventListener('loadedmetadata', markReady);
video.addEventListener('canplay', markReady);
video.addEventListener('playing', function() { hideResumeButton(); markReady(); });
video.addEventListener('timeupdate', noteProgress);
video.addEventListener('stalled', function() {
  if (activeSourceKind === 'dash' || preferredFormat === 'dash') {
    tryHLSMasterFallback();
    return;
  }
  tryHLSMediaFallback();
});
video.addEventListener('waiting', function() {
  if (video.readyState <= 2) {
    lastProgressAt = Math.min(lastProgressAt || Date.now(), Date.now()-5000);
  }
});
video.addEventListener('pause', function() {
  if (offline && offline.style.display === 'none') showResumeButton();
});
video.addEventListener('error', scheduleRetry);

applySkin();
tryPlay();
</script>
</body>
</html>`,
		streamKey,
		bodyClass,
		header,
		stageClass,
		footer,
		streamKey,
		format,
		autoplay,
		muted,
	)
}

func renderAudioHTML(streamKey, format string, autoplay, muted bool, embedded bool) string {
	bodyClass := "audio-page"
	if embedded {
		bodyClass += " embedded"
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="tr">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s - FluxStream Audio</title>
<link rel="stylesheet" href="/static/vendor/bootstrap-icons.css">
<script src="/static/vendor/hls.min.js"></script>
<script src="/static/vendor/dash.all.min.js"></script>
<style>
*{margin:0;padding:0;box-sizing:border-box}
html,body{width:100%%;height:100%%}
body.audio-page{background:#07111f;color:#fff;font-family:'Segoe UI',Tahoma,Geneva,Verdana,sans-serif;display:flex;align-items:center;justify-content:center;padding:18px}
body.audio-page.embedded{padding:10px;background:transparent}
.audio-shell{position:relative;width:100%%;max-width:720px;display:flex;flex-direction:column;gap:14px;padding:18px;border-radius:22px;background:rgba(15,23,42,.76);backdrop-filter:blur(14px);box-shadow:0 24px 60px rgba(2,6,23,.28)}
.audio-shell.player-shell.audio-only{min-height:180px}
.audio-head{display:flex;justify-content:space-between;align-items:center;gap:12px}
.audio-title{display:flex;align-items:center;gap:10px;min-width:0}
.audio-title strong{font-size:16px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.audio-badge{display:inline-flex;align-items:center;gap:6px;padding:6px 10px;border-radius:999px;background:rgba(16,185,129,.18);font-size:10px;font-weight:700;letter-spacing:.08em;text-transform:uppercase;color:#d1fae5}
.audio-badge::before{content:'';width:6px;height:6px;border-radius:50%%;background:#34d399}
.floating-logo{position:absolute;z-index:4;display:none;pointer-events:none}
.floating-logo img{height:34px;max-width:150px;object-fit:contain}
.audio-watermark{display:none;position:absolute;left:18px;bottom:18px;padding:6px 10px;border-radius:999px;background:rgba(15,23,42,.44);font-size:11px;font-weight:700;letter-spacing:.08em;color:rgba(255,255,255,.78)}
audio{width:100%%}
.audio-note{font-size:12px;color:#94a3b8;line-height:1.6;min-height:18px}
.resume-button{align-self:flex-start;display:none;align-items:center;gap:8px;padding:12px 18px;border:0;border-radius:999px;background:rgba(37,99,235,.88);color:#fff;font-size:13px;font-weight:700;cursor:pointer}
</style>
</head>
<body class="%s">
<div class="audio-shell player-shell audio-only">
  <div id="audio-floating-logo" class="floating-logo"></div>
  <div class="audio-head">
    <div class="audio-title" id="audio-title"><strong>%s</strong></div>
    <div class="audio-badge" id="audio-badge">Live</div>
  </div>
  <audio id="audio" controls playsinline></audio>
  <button id="audio-resume" class="resume-button" type="button"><i class="bi bi-play-fill"></i><span>Play</span></button>
  <div id="audio-watermark" class="audio-watermark"></div>
  <div id="audio-note" class="audio-note"></div>
</div>
<script>
const streamKey = %q;
const preferredFormat = (new URLSearchParams(location.search).get('format') || %q).toLowerCase();
const queryParams = new URLSearchParams(location.search);
const autoplay = queryParams.has('autoplay') ? (queryParams.get('autoplay') === '1' || queryParams.get('autoplay') === 'true') : %t;
const muted = queryParams.has('muted') ? (queryParams.get('muted') === '1' || queryParams.get('muted') === 'true') : %t;
const audio = document.getElementById('audio');
const note = document.getElementById('audio-note');
const titleBox = document.getElementById('audio-title');
const badgeBox = document.getElementById('audio-badge');
const logoBox = document.getElementById('audio-floating-logo');
const watermark = document.getElementById('audio-watermark');
const resumeButton = document.getElementById('audio-resume');
let hls;
let dashPlayer;
let primaryReady = false;
let fallbackTimer = null;

audio.autoplay = autoplay;
audio.muted = muted;

function passthroughURL(url) {
  const next = new URL(url, location.origin);
  ['token','password'].forEach(function(key) {
    const val = queryParams.get(key);
    if (val && !next.searchParams.has(key)) next.searchParams.set(key, val);
  });
  return next.toString();
}

function applyStyleText(node, value, fallbackProp) {
  if (!node || !value) return;
  if (value.indexOf(':') !== -1 || value.indexOf(';') !== -1) {
    node.style.cssText += ';' + value;
    return;
  }
  if (fallbackProp) node.style[fallbackProp] = value;
}

function logoPositionStyle(position) {
  switch ((position || 'top-right').toLowerCase()) {
    case 'top-left':
      return 'top:16px;left:16px;';
    case 'bottom-left':
      return 'bottom:16px;left:16px;';
    case 'bottom-right':
      return 'bottom:16px;right:16px;';
    default:
      return 'top:16px;right:16px;';
  }
}

function setNote(message) {
  if (note) note.textContent = message || '';
}

function showResumeButton() {
  if (resumeButton) resumeButton.style.display = 'inline-flex';
}

function hideResumeButton() {
  if (resumeButton) resumeButton.style.display = 'none';
}

function tryAutoplay() {
  if (!autoplay) {
    showResumeButton();
    return;
  }
  const playPromise = audio.play();
  if (playPromise && playPromise.catch) {
    playPromise.catch(function() { showResumeButton(); });
  }
}

function applySkin() {
  const title = queryParams.get('player_title') || streamKey;
  const bg = queryParams.get('player_bg') || '';
  const controls = queryParams.get('player_controls') || '';
  const playStyle = queryParams.get('player_play') || '';
  const logo = queryParams.get('player_logo') || '';
  const logoPosition = queryParams.get('player_logo_position') || 'top-right';
  const logoOpacity = queryParams.get('player_logo_opacity') || '1';
  const watermarkText = queryParams.get('player_watermark') || '';
  const showTitle = queryParams.get('player_show_title') !== '0';
  const showBadge = queryParams.get('player_show_badge') !== '0';
  const customCSS = queryParams.get('player_custom_css') || '';
  if (titleBox && showTitle) {
    titleBox.innerHTML = '<strong>' + title + '</strong>';
  } else if (titleBox) {
    titleBox.style.display = 'none';
  }
  if (badgeBox && !showBadge) badgeBox.style.display = 'none';
  if (logo && logoBox) {
    logoBox.style.display = 'block';
    logoBox.style.cssText += ';' + logoPositionStyle(logoPosition);
    logoBox.innerHTML = '<img src="' + logo + '" alt="logo" style="opacity:' + logoOpacity + '">';
  }
  if (watermarkText && watermark) {
    watermark.style.display = 'block';
    watermark.textContent = watermarkText;
  }
  applyStyleText(document.body, bg, 'background');
  applyStyleText(document.querySelector('.audio-shell'), bg, 'background');
  applyStyleText(document.querySelector('.audio-shell'), controls);
  applyStyleText(resumeButton, playStyle);
  if (customCSS) {
    const style = document.createElement('style');
    style.textContent = customCSS;
    document.head.appendChild(style);
  }
}

function primaryURLForFormat(fmt) {
  switch (fmt) {
    case 'mp3':
      return passthroughURL(location.origin + '/audio/mp3/' + streamKey + '/' + streamKey + '.mp3');
    case 'aac':
      return passthroughURL(location.origin + '/audio/aac/' + streamKey + '/' + streamKey + '.aac');
    case 'ogg':
      return passthroughURL(location.origin + '/audio/ogg/' + streamKey + '/' + streamKey + '.ogg');
    case 'wav':
      return passthroughURL(location.origin + '/audio/wav/' + streamKey + '/' + streamKey + '.wav');
    case 'flac':
      return passthroughURL(location.origin + '/audio/flac/' + streamKey + '/' + streamKey + '.flac');
    case 'icecast':
      return passthroughURL(location.origin + '/icecast/' + streamKey);
    default:
      return passthroughURL(location.origin + '/audio/aac/' + streamKey + '/' + streamKey + '.aac');
  }
}

const primaryUrl = primaryURLForFormat(preferredFormat);
const hlsFallbackUrl = passthroughURL(location.origin + '/audio/hls/' + streamKey);
const dashFallbackUrl = passthroughURL(location.origin + '/audio/dash/' + streamKey);

function cleanupPlayers() {
  try { if (hls) { hls.destroy(); hls = null; } } catch (e) {}
  try { if (dashPlayer) { dashPlayer.reset(); dashPlayer = null; } } catch (e) {}
  audio.pause();
  audio.removeAttribute('src');
  audio.load();
}

function markReady(message) {
  primaryReady = true;
  if (fallbackTimer) {
    clearTimeout(fallbackTimer);
    fallbackTimer = null;
  }
  setNote(message || 'Canli ses baglandi.');
  hideResumeButton();
  if (badgeBox && badgeBox.style.display !== 'none') badgeBox.textContent = 'Live';
  tryAutoplay();
}

function startNative(url, message) {
  cleanupPlayers();
  audio.src = url;
  audio.load();
  setNote(message);
}

function startHLS(url) {
  cleanupPlayers();
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
    hls.attachMedia(audio);
    hls.on(Hls.Events.MANIFEST_PARSED, function() {
      markReady('Canli HLS ses baglandi.');
    });
    hls.on(Hls.Events.ERROR, function(event, data) {
      if (data && data.fatal) startDASHFallback();
    });
    return;
  }
  if (audio.canPlayType && audio.canPlayType('application/vnd.apple.mpegurl')) {
    audio.src = url;
    audio.load();
    setNote('Canli HLS ses baglandi.');
    return;
  }
  startDASHFallback();
}

function startDASHFallback() {
  if (!window.dashjs || !window.dashjs.MediaPlayer) {
    setNote('Ses akisi simdilik hazir degil.');
    return;
  }
  cleanupPlayers();
  dashPlayer = window.dashjs.MediaPlayer().create();
  dashPlayer.initialize(audio, dashFallbackUrl, autoplay);
  dashPlayer.on(window.dashjs.MediaPlayer.events.STREAM_INITIALIZED, function() {
    markReady('Canli DASH ses baglandi.');
  });
  dashPlayer.on(window.dashjs.MediaPlayer.events.ERROR, function() {
    setNote('Ses akisi simdilik hazir degil.');
  });
}

function startFallback() {
  setNote('Dogrudan ses cikisi hazir degil, HLS/DASH fallback deneniyor.');
  startHLS(hlsFallbackUrl);
}

if (resumeButton) {
  resumeButton.addEventListener('click', function() {
    audio.play().then(function() { hideResumeButton(); }).catch(function() {});
  });
}

audio.addEventListener('canplay', function() {
  markReady(preferredFormat.toUpperCase() + ' ses akisi baglandi.');
});
audio.addEventListener('playing', hideResumeButton);
audio.addEventListener('pause', function() {
  if (primaryReady) showResumeButton();
});
audio.addEventListener('error', function() {
  startFallback();
});

fallbackTimer = setTimeout(function() {
  if (!primaryReady) startFallback();
}, 2200);

applySkin();
startNative(primaryUrl, preferredFormat.toUpperCase() + ' ses akisi bekleniyor...');
</script>
</body>
</html>`,
		streamKey,
		bodyClass,
		streamKey,
		streamKey,
		format,
		autoplay,
		muted,
	)
}
