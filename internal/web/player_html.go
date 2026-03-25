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
.audio-track-box{position:absolute;top:16px;right:16px;z-index:6;display:none;align-items:center;gap:8px;padding:10px 12px;border-radius:14px;background:rgba(15,23,42,.56);backdrop-filter:blur(12px);border:1px solid rgba(148,163,184,.18)}
.audio-track-box label{font-size:11px;font-weight:600;letter-spacing:.06em;text-transform:uppercase;color:#cbd5e1}
.audio-track-box select{min-width:160px;padding:8px 10px;border-radius:10px;border:1px solid rgba(148,163,184,.22);background:#0f172a;color:#fff;font-size:12px}
.qoe-debug{position:absolute;right:14px;bottom:14px;z-index:7;display:none;min-width:260px;max-width:min(360px,calc(100%% - 28px));padding:12px 14px;border-radius:14px;background:rgba(2,6,23,.78);backdrop-filter:blur(12px);border:1px solid rgba(148,163,184,.18);box-shadow:0 16px 40px rgba(2,6,23,.34);font:12px/1.5 Consolas,Monaco,'Courier New',monospace;color:#dbeafe}
.qoe-debug.visible{display:block}
.qoe-debug-title{font:600 11px/1.2 'Segoe UI',Tahoma,Geneva,Verdana,sans-serif;letter-spacing:.08em;text-transform:uppercase;color:#93c5fd;margin-bottom:8px}
.qoe-debug-row{display:flex;justify-content:space-between;gap:12px;padding:2px 0}
.qoe-debug-row span:first-child{color:#94a3b8}
.qoe-debug-row span:last-child{text-align:right;word-break:break-word}
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
  <div id="audio-track-box" class="audio-track-box">
    <label for="audio-track-select">Ses</label>
    <select id="audio-track-select"></select>
  </div>
  <div id="offline" class="offline-msg">
    <h2>Yayin cevrimdisi</h2>
    <p>Yayin basladiginda otomatik oynatilacak</p>
  </div>
  <div id="qoe-debug" class="qoe-debug"></div>
</div>
%s
<script>
const streamKey = %q;
const queryParams = new URLSearchParams(location.search);
const preferredFormat = (queryParams.get('format') || %q).toLowerCase();
const embedded = %t;
const debugEnabled = queryParams.get('debug') === '1' || queryParams.get('debug') === 'true';
const autoplay = queryParams.has('autoplay') ? (queryParams.get('autoplay') === '1' || queryParams.get('autoplay') === 'true') : %t;
const muted = queryParams.has('muted') ? (queryParams.get('muted') === '1' || queryParams.get('muted') === 'true') : %t;
const telemetryEndpoint = location.origin + '/api/player/telemetry';
const telemetrySessionKey = 'fluxstream_qoe_' + streamKey;
const audioPreferenceKey = 'fluxstream_audio_pref_' + streamKey;
const video = document.getElementById('video');
const offline = document.getElementById('offline');
const qoeDebug = document.getElementById('qoe-debug');
const badge = document.getElementById('live-badge');
const playerTitle = document.getElementById('player-title');
const logoBox = document.getElementById('player-floating-logo');
const watermark = document.getElementById('player-watermark');
const resumeButton = document.getElementById('resume-button');
const audioTrackBox = document.getElementById('audio-track-box');
const audioTrackSelect = document.getElementById('audio-track-select');
let hls;
let dashPlayer;
let retryTimer = null;
let activeSourceKind = '';
let sourceOverride = '';
let watchdogTimer = null;
let lastProgressAt = 0;
let lastCurrentTime = 0;
let stallRecoveries = 0;
let lastErrorMessage = '-';
let lastErrorSourceKind = '';
let fallbackNote = '-';
let reconnectState = 'idle';
let retryAt = 0;
let qoeTimer = null;
let telemetryTimer = null;
let telemetryInflight = false;
let telemetryDirty = false;
let hlsAudioTracks = [];
let dashAudioTracks = [];
let preferredAudioTrack = (queryParams.get('audio_track') || '').trim();
let preferredAudioTrackLabel = '';
let preferredAudioTrackApplied = false;
let qualityTransitionCount = 0;
let audioSwitchCount = 0;
let lastQualitySignature = '';
let lastSelectedAudioSignature = '';
let selectedAudioTrackID = '';
let selectedAudioTrackLabel = '-';
let dashRetryCount = 0;
const qoeState = {
  preferredFormat: preferredFormat || 'auto',
  sourceOverride: 'auto',
  activeSourceKind: '-',
  quality: '-',
  stallCount: 0,
  recoveries: 0,
  lastError: '-',
  fallback: '-',
  reconnect: 'idle',
  offline: 'hidden'
};

function loadSavedAudioPreference() {
  try {
    const raw = localStorage.getItem(audioPreferenceKey);
    if (!raw) return { id: '', label: '' };
    const parsed = JSON.parse(raw);
    return {
      id: String((parsed && parsed.id) || '').trim(),
      label: String((parsed && parsed.label) || '').trim()
    };
  } catch (e) {
    return { id: '', label: '' };
  }
}

function saveAudioPreference(id, label) {
  try {
    localStorage.setItem(audioPreferenceKey, JSON.stringify({
      id: String(id || '').trim(),
      label: String(label || '').trim()
    }));
  } catch (e) {}
}

const savedAudioPreference = loadSavedAudioPreference();
if (!preferredAudioTrack && savedAudioPreference.id) preferredAudioTrack = savedAudioPreference.id;
preferredAudioTrackLabel = savedAudioPreference.label || '';

video.autoplay = autoplay;
video.muted = muted;

function createTelemetrySessionID() {
  try {
    if (window.crypto && typeof window.crypto.randomUUID === 'function') return window.crypto.randomUUID();
  } catch (e) {}
  return 'qoe_' + Math.random().toString(36).slice(2) + Date.now().toString(36);
}

function getTelemetrySessionID() {
  try {
    const current = sessionStorage.getItem(telemetrySessionKey);
    if (current) return current;
    const created = createTelemetrySessionID();
    sessionStorage.setItem(telemetrySessionKey, created);
    return created;
  } catch (e) {
    return createTelemetrySessionID();
  }
}

const telemetrySessionID = getTelemetrySessionID();

function formatSeconds(value) {
  if (!Number.isFinite(value)) return '-';
  return value.toFixed(1) + 's';
}

function getBufferedSeconds() {
  try {
    if (video.buffered && video.buffered.length) {
      return Math.max(0, video.buffered.end(video.buffered.length - 1) - (video.currentTime || 0));
    }
  } catch (e) {}
  return 0;
}

function updateQualityLabel() {
  const hlsInfo = getHLSQualityInfo();
  if (hlsInfo) {
    qoeState.quality = hlsInfo.label;
    return;
  }
  const dashInfo = getDashQualityInfo();
  if (dashInfo) {
    qoeState.quality = dashInfo.label;
    return;
  }
  qoeState.quality = activeSourceKind || '-';
}

function getHLSQualityInfo() {
  if (!hls || !hls.levels || !hls.levels.length) return null;
  const levelIndex = hls.currentLevel >= 0 ? hls.currentLevel : hls.loadLevel;
  const level = levelIndex >= 0 ? hls.levels[levelIndex] : null;
  if (!level) return null;
  return {
    id: String(levelIndex),
    label: (level.height || '?') + 'p @ ' + Math.round((level.bitrate || 0) / 1000) + ' kbps'
  };
}

function getDashQualityInfo() {
  if (!dashPlayer || typeof dashPlayer.getQualityFor !== 'function' || typeof dashPlayer.getBitrateInfoListFor !== 'function') {
    return null;
  }
  const idx = dashPlayer.getQualityFor('video');
  const list = dashPlayer.getBitrateInfoListFor('video') || [];
  const info = idx >= 0 ? list[idx] : null;
  if (!info) return null;
  return {
    id: String(info.id != null ? info.id : idx),
    label: (info.height || '?') + 'p @ ' + Math.round((info.bitrate || 0) / 1000) + ' kbps'
  };
}

function rememberQualityInfo(info, allowTransition) {
  if (!info) return;
  qoeState.quality = info.label || 'auto';
  const nextSignature = String(info.id || info.label || '');
  if (!nextSignature) return;
  if (!lastQualitySignature) {
    lastQualitySignature = nextSignature;
    return;
  }
  if (allowTransition && nextSignature !== lastQualitySignature) {
    qualityTransitionCount += 1;
  }
  lastQualitySignature = nextSignature;
}

function rememberSelectedAudioTrack(id, label, allowTransition) {
  const nextID = String(id || '').trim();
  const nextLabel = String(label || nextID || '-').trim() || '-';
  selectedAudioTrackID = nextID;
  selectedAudioTrackLabel = nextLabel;
  if (nextID || nextLabel !== '-') {
    preferredAudioTrack = nextID || preferredAudioTrack;
    preferredAudioTrackLabel = nextLabel;
    saveAudioPreference(preferredAudioTrack, preferredAudioTrackLabel);
  }
  const nextSignature = nextID + '|' + nextLabel;
  if (!nextSignature.trim()) return;
  if (!lastSelectedAudioSignature) {
    lastSelectedAudioSignature = nextSignature;
    return;
  }
  if (allowTransition && nextSignature !== lastSelectedAudioSignature) {
    audioSwitchCount += 1;
  }
  lastSelectedAudioSignature = nextSignature;
}

function renderQOEDebug() {
  if (!debugEnabled || !qoeDebug) return;
  qoeState.preferredFormat = preferredFormat || 'auto';
  qoeState.sourceOverride = sourceOverride || 'auto';
  qoeState.activeSourceKind = activeSourceKind || '-';
  qoeState.recoveries = stallRecoveries;
  qoeState.lastError = getVisibleError();
  qoeState.fallback = fallbackNote || '-';
  qoeState.reconnect = reconnectState;
  qoeState.offline = offline && offline.style.display !== 'none' ? 'visible' : 'hidden';
  updateQualityLabel();
  qoeDebug.classList.add('visible');
  qoeDebug.innerHTML =
    '<div class="qoe-debug-title">QoE Debug</div>' +
    '<div class="qoe-debug-row"><span>Format</span><span>' + qoeState.preferredFormat + '</span></div>' +
    '<div class="qoe-debug-row"><span>Kaynak</span><span>' + qoeState.activeSourceKind + '</span></div>' +
    '<div class="qoe-debug-row"><span>Override</span><span>' + qoeState.sourceOverride + '</span></div>' +
    '<div class="qoe-debug-row"><span>Kalite</span><span>' + qoeState.quality + '</span></div>' +
    '<div class="qoe-debug-row"><span>Ses</span><span>' + (selectedAudioTrackLabel || '-') + '</span></div>' +
    '<div class="qoe-debug-row"><span>Sure</span><span>' + formatSeconds(video.currentTime || 0) + '</span></div>' +
    '<div class="qoe-debug-row"><span>Buffer</span><span>' + formatSeconds((video.buffered && video.buffered.length) ? (video.buffered.end(video.buffered.length - 1) - (video.currentTime || 0)) : 0) + '</span></div>' +
    '<div class="qoe-debug-row"><span>Stall</span><span>' + qoeState.stallCount + '</span></div>' +
    '<div class="qoe-debug-row"><span>Toparlanma</span><span>' + qoeState.recoveries + '</span></div>' +
    '<div class="qoe-debug-row"><span>Kalite Gecisi</span><span>' + qualityTransitionCount + '</span></div>' +
    '<div class="qoe-debug-row"><span>Audio Gecisi</span><span>' + audioSwitchCount + '</span></div>' +
    '<div class="qoe-debug-row"><span>Reconnect</span><span>' + qoeState.reconnect + '</span></div>' +
    '<div class="qoe-debug-row"><span>Gecis</span><span>' + qoeState.fallback + '</span></div>' +
    '<div class="qoe-debug-row"><span>Offline</span><span>' + qoeState.offline + '</span></div>' +
    '<div class="qoe-debug-row"><span>Hata</span><span>' + qoeState.lastError + '</span></div>';
}

function buildTelemetryPayload() {
  updateQualityLabel();
  return {
    stream_key: streamKey,
    session_id: telemetrySessionID,
    page: embedded ? 'embed' : 'player',
    preferred_format: preferredFormat || 'auto',
    active_source_kind: activeSourceKind || '-',
    source_override: sourceOverride || 'auto',
    quality: qoeState.quality || '-',
    selected_audio_track: selectedAudioTrackID || '-',
    selected_audio_label: selectedAudioTrackLabel || '-',
    playback_seconds: Number(video.currentTime || 0),
    buffer_seconds: Number(getBufferedSeconds()),
    stall_count: Number(qoeState.stallCount || 0),
    recoveries: Number(stallRecoveries || 0),
    quality_transitions: Number(qualityTransitionCount || 0),
    audio_switches: Number(audioSwitchCount || 0),
    last_error: getVisibleError(),
    reconnect: reconnectState || 'idle',
    offline: !!(offline && offline.style.display !== 'none'),
    waiting: reconnectState === 'waiting' || reconnectState === 'stalled' || reconnectState === 'retrying' || reconnectState === 'recovering',
    debug_enabled: debugEnabled
  };
}

function sendTelemetry() {
  if (!streamKey) return;
  if (telemetryInflight) {
    telemetryDirty = true;
    return;
  }
  telemetryInflight = true;
  telemetryDirty = false;
  fetch(telemetryEndpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    cache: 'no-store',
    keepalive: true,
    body: JSON.stringify(buildTelemetryPayload())
  }).catch(function() {}).finally(function() {
    telemetryInflight = false;
    if (telemetryDirty) {
      telemetryDirty = false;
      sendTelemetry();
    }
  });
}

function startTelemetryLoop() {
  if (telemetryTimer) return;
  sendTelemetry();
  telemetryTimer = setInterval(sendTelemetry, 5000);
}

function startQOEDebugLoop() {
  if (!debugEnabled || qoeTimer) return;
  renderQOEDebug();
  qoeTimer = setInterval(renderQOEDebug, 1000);
}

function getVisibleError() {
  if (reconnectState === 'playing' && activeSourceKind && lastErrorSourceKind && activeSourceKind !== lastErrorSourceKind) {
    return '-';
  }
  return lastErrorMessage || '-';
}

function setLastError(message, sourceKind) {
  lastErrorMessage = message || '-';
  lastErrorSourceKind = sourceKind || activeSourceKind || '';
  renderQOEDebug();
  sendTelemetry();
}

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

function hideAudioTrackSelector() {
  if (audioTrackBox) audioTrackBox.style.display = 'none';
  if (audioTrackSelect) audioTrackSelect.innerHTML = '';
}

function normalizeAudioToken(value) {
  return String(value || '').trim().toLowerCase();
}

function buildAudioTrackLabel(track, fallbackIndex) {
  const index = Number(fallbackIndex || 0);
  const attrs = track && track.attrs ? track.attrs : {};
  const parts = [];
  const primary = String(
    (track && (track.label || track.name || track.displayName || track.lang || track.groupId)) ||
    attrs.NAME ||
    attrs.LANGUAGE ||
    ''
  ).trim();
  if (primary) parts.push(primary);
  const lang = String((track && track.lang) || attrs.LANGUAGE || '').trim();
  if (lang && parts.every(function(item) { return normalizeAudioToken(item) !== normalizeAudioToken(lang); })) {
    parts.push(lang.toUpperCase());
  }
  const channels = String((track && track.channels) || attrs.CHANNELS || '').trim();
  if (channels) parts.push(channels.indexOf('ch') !== -1 ? channels : channels + 'ch');
  const role = String((track && track.role) || attrs.CHARACTERISTICS || '').trim();
  if (role) parts.push(role);
  if (!parts.length) parts.push('Ses Track ' + (index + 1));
  return parts.join(' • ');
}

function audioTrackAliases(track) {
  const aliases = [];
  const pushToken = function(value) {
    const token = normalizeAudioToken(value);
    if (!token || aliases.indexOf(token) !== -1) return;
    aliases.push(token);
  };
  pushToken(track && track.id);
  pushToken(track && track.name);
  pushToken(track && track.lang);
  pushToken(track && track.label);
  pushToken(track && track.displayName);
  pushToken(track && track.groupId);
  pushToken(track && track.role);
  pushToken(buildAudioTrackLabel(track, Number(track && track.id) || 0));
  if (track && track.id != null) {
    pushToken('track ' + track.id);
    pushToken('track-' + track.id);
    pushToken('audio track ' + track.id);
  }
  return aliases;
}

function trackMatchesPreferred(track, preferred) {
  const target = normalizeAudioToken(preferred);
  const labelTarget = normalizeAudioToken(preferredAudioTrackLabel);
  if (!target && !labelTarget) return false;
  const aliases = audioTrackAliases(track);
  if (target && aliases.indexOf(target) !== -1) return true;
  if (labelTarget && aliases.indexOf(labelTarget) !== -1) return true;
  return aliases.some(function(alias) {
    return (target && alias.indexOf(target) !== -1) || (labelTarget && alias.indexOf(labelTarget) !== -1);
  });
}

function renderAudioTrackSelector(items, selectedIndex, applyFn) {
  if (!audioTrackBox || !audioTrackSelect || !Array.isArray(items) || items.length <= 1) {
    hideAudioTrackSelector();
    return;
  }
  audioTrackBox.style.display = 'inline-flex';
  audioTrackSelect.innerHTML = items.map(function(item, index) {
    const label = item.label || item.name || item.lang || ('Track ' + (index + 1));
    const selected = index === selectedIndex ? ' selected' : '';
    return '<option value="' + index + '"' + selected + '>' + label + '</option>';
  }).join('');
  audioTrackSelect.onchange = function() {
    const nextIndex = parseInt(audioTrackSelect.value || '0', 10) || 0;
    applyFn(nextIndex);
  };
}

function syncHLSAudioTracks(allowTransition) {
  if (!hls || !Array.isArray(hls.audioTracks)) {
    hideAudioTrackSelector();
    return;
  }
  hlsAudioTracks = hls.audioTracks.map(function(track, index) {
    return {
      id: track && track.id != null ? track.id : index,
      name: (track && (track.name || track.lang || track.groupId)) || ('Track ' + (index + 1)),
      label: buildAudioTrackLabel({
        id: track && track.id != null ? track.id : index,
        name: track && track.name,
        lang: track && track.lang,
        label: track && (track.name || track.lang || track.groupId),
        groupId: track && track.groupId,
        attrs: track && track.attrs
      }, index)
    };
  });
  if (!preferredAudioTrackApplied && preferredAudioTrack) {
    const requestedIndex = hlsAudioTracks.findIndex(function(track) {
      return trackMatchesPreferred(track, preferredAudioTrack);
    });
    if (requestedIndex >= 0) {
      try { hls.audioTrack = requestedIndex; } catch (e) {}
      preferredAudioTrackApplied = true;
    }
  }
  const selectedIndex = hls.audioTrack >= 0 ? hls.audioTrack : hlsAudioTracks.findIndex(function(track) { return trackMatchesPreferred(track, preferredAudioTrack); });
  const currentIndex = selectedIndex >= 0 ? selectedIndex : 0;
  if (hlsAudioTracks[currentIndex]) {
    rememberSelectedAudioTrack(hlsAudioTracks[currentIndex].id, hlsAudioTracks[currentIndex].label, !!allowTransition);
  }
  renderAudioTrackSelector(hlsAudioTracks, currentIndex, function(nextIndex) {
    if (!hls) return;
    try { hls.audioTrack = nextIndex; } catch (e) {}
    const selected = hlsAudioTracks[nextIndex];
    if (selected) {
      preferredAudioTrack = String(selected.id || nextIndex + 1);
      rememberSelectedAudioTrack(selected.id, selected.label, true);
    }
    renderQOEDebug();
    sendTelemetry();
  });
}

function syncDashAudioTracks(allowTransition) {
  if (!dashPlayer || typeof dashPlayer.getTracksFor !== 'function') {
    hideAudioTrackSelector();
    return;
  }
  dashAudioTracks = (dashPlayer.getTracksFor('audio') || []).map(function(track, index) {
    return {
      raw: track,
      id: track && track.id != null ? track.id : index,
      name: (track && (track.lang || track.id || track.labels && track.labels[0])) || ('Track ' + (index + 1)),
      label: buildAudioTrackLabel({
        id: track && track.id != null ? track.id : index,
        lang: track && track.lang,
        label: track && ((track.labels && track.labels[0]) || track.lang || track.id),
        role: track && track.roles && track.roles[0],
        channels: track && track.audioChannelConfiguration && track.audioChannelConfiguration[0]
      }, index)
    };
  });
  if (!dashAudioTracks.length || dashAudioTracks.length === 1) {
    hideAudioTrackSelector();
    return;
  }
  let selectedIndex = 0;
  if (typeof dashPlayer.getCurrentTrackFor === 'function') {
    const current = dashPlayer.getCurrentTrackFor('audio');
    const idx = dashAudioTracks.findIndex(function(track) { return track.raw === current; });
    if (idx >= 0) selectedIndex = idx;
  }
  if (!preferredAudioTrackApplied && preferredAudioTrack) {
    const requestedIndex = dashAudioTracks.findIndex(function(track) {
      return trackMatchesPreferred(track, preferredAudioTrack);
    });
    if (requestedIndex >= 0) {
      selectedIndex = requestedIndex;
      if (typeof dashPlayer.setCurrentTrack === 'function') {
        try { dashPlayer.setCurrentTrack(dashAudioTracks[requestedIndex].raw); } catch (e) {}
      }
      preferredAudioTrackApplied = true;
    }
  }
  if (dashAudioTracks[selectedIndex]) {
    rememberSelectedAudioTrack(dashAudioTracks[selectedIndex].id, dashAudioTracks[selectedIndex].label, !!allowTransition);
  }
  renderAudioTrackSelector(dashAudioTracks, selectedIndex, function(nextIndex) {
    if (!dashPlayer || typeof dashPlayer.setCurrentTrack !== 'function' || !dashAudioTracks[nextIndex]) return;
    try { dashPlayer.setCurrentTrack(dashAudioTracks[nextIndex].raw); } catch (e) {}
    preferredAudioTrack = String(dashAudioTracks[nextIndex].id || nextIndex + 1);
    rememberSelectedAudioTrack(dashAudioTracks[nextIndex].id, dashAudioTracks[nextIndex].label, true);
    renderQOEDebug();
    sendTelemetry();
  });
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
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' },
      { kind: 'native', mime: 'video/mp4', url: passthroughURL(location.origin + '/mp4/' + streamKey + '/' + streamKey + '.mp4') }
    ],
    webm: [
      { kind: 'hls', url: passthroughURL(location.origin + '/hls/' + streamKey + '/master.m3u8'), marker: '#EXTM3U' },
      { kind: 'dash', url: passthroughURL(location.origin + '/dash/' + streamKey + '/manifest.mpd'), marker: '<MPD' },
      { kind: 'native', mime: 'video/webm', url: passthroughURL(location.origin + '/webm/' + streamKey + '/' + streamKey + '.webm') }
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
  hlsAudioTracks = [];
  dashAudioTracks = [];
  preferredAudioTrackApplied = false;
  selectedAudioTrackID = '';
  selectedAudioTrackLabel = '-';
  hideAudioTrackSelector();
  activeSourceKind = '';
  reconnectState = 'idle';
  fallbackNote = '-';
  hideResumeButton();
  video.pause();
  video.removeAttribute('src');
  video.load();
  renderQOEDebug();
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
  reconnectState = 'fallback:' + nextOverride;
  renderQOEDebug();
  await tryPlay();
}

async function tryHLSMediaFallback() {
  await trySourceFallback('hls_media');
}

async function tryHLSMasterFallback() {
  await trySourceFallback('hls');
}

function handleHLSLevelChange(allowTransition) {
  rememberQualityInfo(getHLSQualityInfo(), !!allowTransition);
  renderQOEDebug();
  sendTelemetry();
}

function handleDashQualityChange(allowTransition) {
  rememberQualityInfo(getDashQualityInfo(), !!allowTransition);
  renderQOEDebug();
  sendTelemetry();
}

function retryDASHSource(url, reason) {
  if (dashRetryCount >= 2) {
    tryHLSMasterFallback();
    return;
  }
  dashRetryCount += 1;
  reconnectState = 'dash-retry-' + dashRetryCount;
  fallbackNote = reason || 'dash retry';
  renderQOEDebug();
  sendTelemetry();
  setTimeout(function() {
    cleanupPlayers();
    reconnectState = 'dash-retrying';
    renderQOEDebug();
    startDASH(url);
  }, 900);
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
    qoeState.stallCount += 1;
    lastProgressAt = now;
    reconnectState = 'recovering';
    renderQOEDebug();
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
  fallbackNote = (lastErrorSourceKind && activeSourceKind && lastErrorSourceKind !== activeSourceKind) ? (lastErrorSourceKind + ' -> ' + activeSourceKind) : '-';
  lastErrorMessage = '-';
  lastErrorSourceKind = '';
  retryAt = 0;
  video.style.display = 'block';
  if (offline) offline.style.display = 'none';
  if (badge && badge.style.display !== 'none') badge.innerHTML = '<span class="badge-live">CANLI</span>';
  reconnectState = 'playing';
  noteProgress();
  renderQOEDebug();
  ensurePlaybackWatchdog();
  startQOEDebugLoop();
  startTelemetryLoop();
  sendTelemetry();
  tryAutoplay();
}

function scheduleRetry() {
  if (retryTimer) return;
  video.style.display = 'none';
  if (offline) offline.style.display = 'block';
  if (badge) badge.innerHTML = '';
  retryAt = Date.now() + 3000;
  reconnectState = 'retrying';
  cleanupPlayers();
  renderQOEDebug();
  sendTelemetry();
  retryTimer = setTimeout(function() {
    retryTimer = null;
    reconnectState = 'retry-now';
    renderQOEDebug();
    tryPlay();
  }, 3000);
}

function startNative(url) {
  activeSourceKind = 'native';
  reconnectState = 'native';
  renderQOEDebug();
  sendTelemetry();
  video.src = url;
  video.load();
}

function startHLS(url) {
  if (window.Hls && Hls.isSupported()) {
    activeSourceKind = 'hls';
    preferredAudioTrackApplied = false;
    hls = new Hls({
      liveSyncDurationCount: preferredFormat === 'll_hls' ? 3 : 4,
      liveMaxLatencyDurationCount: preferredFormat === 'll_hls' ? 6 : 10,
      maxBufferLength: preferredFormat === 'll_hls' ? 14 : 24,
      maxMaxBufferLength: preferredFormat === 'll_hls' ? 24 : 48,
      backBufferLength: 45,
      startLevel: 0,
      abrEwmaDefaultEstimate: preferredFormat === 'll_hls' ? 380000 : 240000,
      abrBandWidthFactor: 0.65,
      abrBandWidthUpFactor: 0.45,
      capLevelToPlayerSize: true,
      capLevelOnFPSDrop: true,
      enableWorker: true,
      lowLatencyMode: preferredFormat === 'll_hls'
    });
    hls.loadSource(url);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, function() {
      syncHLSAudioTracks(false);
      handleHLSLevelChange(false);
      markReady();
    });
    hls.on(Hls.Events.LEVEL_SWITCHED, function() { handleHLSLevelChange(true); });
    hls.on(Hls.Events.AUDIO_TRACKS_UPDATED, function() { syncHLSAudioTracks(false); });
    hls.on(Hls.Events.AUDIO_TRACK_SWITCHED, function() { syncHLSAudioTracks(true); });
    hls.on(Hls.Events.ERROR, function(event, data) {
      setLastError('hls:' + ((data && data.details) || (data && data.type) || 'unknown'), 'hls');
      if (!data || !data.fatal) return;
      if (data.type === Hls.ErrorTypes.NETWORK_ERROR) {
        reconnectState = 'hls-network-retry';
        renderQOEDebug();
        sendTelemetry();
        hls.startLoad();
        return;
      }
      if (data.type === Hls.ErrorTypes.MEDIA_ERROR) {
        reconnectState = 'hls-media-recover';
        renderQOEDebug();
        sendTelemetry();
        hls.recoverMediaError();
        return;
      }
      tryHLSMediaFallback();
    });
    return;
  }
  if (video.canPlayType('application/vnd.apple.mpegurl')) {
    activeSourceKind = 'hls';
    reconnectState = 'native-hls';
    renderQOEDebug();
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
  preferredAudioTrackApplied = false;
  dashPlayer = window.dashjs.MediaPlayer().create();
  dashPlayer.updateSettings({
    streaming: {
      lowLatencyEnabled: false,
      liveDelay: 10,
      liveCatchup: { enabled: true },
      buffer: {
        fastSwitchEnabled: true,
        stableBufferTime: 10,
        bufferTimeAtTopQuality: 12,
        bufferTimeAtTopQualityLongForm: 20
      },
      gaps: {
        jumpGaps: true,
        smallGapLimit: 1.5
      },
      retryIntervals: {
        MPD: 700,
        InitializationSegment: 1000,
        MediaSegment: 1000
      },
      retryAttempts: {
        MPD: 5,
        InitializationSegment: 4,
        MediaSegment: 6
      },
      abr: {
        autoSwitchBitrate: { video: true, audio: true },
        initialBitrate: { video: 350, audio: 64 }
      }
    }
  });
  dashPlayer.on(window.dashjs.MediaPlayer.events.STREAM_INITIALIZED, function() {
    dashRetryCount = 0;
    syncDashAudioTracks(false);
    handleDashQualityChange(false);
    markReady();
  });
  dashPlayer.on(window.dashjs.MediaPlayer.events.QUALITY_CHANGE_RENDERED, function() { handleDashQualityChange(true); });
  if (window.dashjs.MediaPlayer.events.TRACK_CHANGE_RENDERED) {
    dashPlayer.on(window.dashjs.MediaPlayer.events.TRACK_CHANGE_RENDERED, function() { syncDashAudioTracks(true); });
  }
  dashPlayer.initialize(video, url, autoplay);
  renderQOEDebug();
  dashPlayer.on(window.dashjs.MediaPlayer.events.ERROR, function(evt) {
    const message = ((evt && evt.error && evt.error.message) || (evt && evt.event && evt.event.message) || (evt && evt.message) || 'unknown');
    const normalized = String(message || '').toLowerCase();
    setLastError('dash:' + message, 'dash');
    if (normalized.indexOf('nostreamscomposed') !== -1 || normalized.indexOf('no periods') !== -1 || normalized.indexOf('buffer') !== -1 || normalized.indexOf('segment') !== -1 || normalized.indexOf('manifest') !== -1) {
      retryDASHSource(url, 'dash recovery');
      return;
    }
    if (dashRetryCount < 1) {
      retryDASHSource(url, 'dash retry');
      return;
    }
    tryHLSMasterFallback();
  });
}

async function tryPlay() {
  cleanupPlayers();
  dashRetryCount = 0;
  reconnectState = 'probing';
  fallbackNote = '-';
  renderQOEDebug();
  startTelemetryLoop();
  sendTelemetry();
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
  qoeState.stallCount += 1;
  reconnectState = 'stalled';
  renderQOEDebug();
  sendTelemetry();
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
  reconnectState = 'waiting';
  renderQOEDebug();
  sendTelemetry();
});
video.addEventListener('pause', function() {
  if (offline && offline.style.display === 'none') showResumeButton();
  reconnectState = 'paused';
  renderQOEDebug();
  sendTelemetry();
});
video.addEventListener('error', function() {
  setLastError('video:' + ((video.error && video.error.message) || (video.error && video.error.code) || 'unknown'), activeSourceKind || 'video');
  sendTelemetry();
  scheduleRetry();
});

window.addEventListener('pagehide', function() {
  if (!navigator.sendBeacon) return;
  try {
    navigator.sendBeacon(telemetryEndpoint, new Blob([JSON.stringify(buildTelemetryPayload())], { type: 'application/json' }));
  } catch (e) {}
});

applySkin();
startQOEDebugLoop();
startTelemetryLoop();
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
		embedded,
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
