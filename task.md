## Update 2026-03-12 (Hybrid preview and clean uninstall)

Completed in this pass:

- [x] Switched `/play` and `/embed` video playback to a hybrid DASH-first, HLS-fallback flow
- [x] Simplified advanced admin preview so HLS/JS API tabs reuse the same stable player iframe
- [x] Updated preview guidance text from HLS-only to DASH/HLS fallback
- [x] Added uninstall cleanup for `{app}\data`, `{app}\ffmpeg` and the install directory
- [x] Rebuilt installer only for live validation

## Update 2026-03-11 (Preview origin and HTTPS readiness)

Completed in this pass:

- [x] Separated admin preview URLs from public embed/output URLs
- [x] Forced admin preview/player/embed iframes to use the current panel origin
- [x] Limited public HTTPS URL generation to cases where SSL is actually ready
- [x] Added HTTP fallback when public HTTPS is requested but no active SSL listener exists

## Update 2026-03-11 (Playback token passthrough)

Completed in this pass:

- [x] Fixed live preview and copied external links when playback token protection is enabled
- [x] Added token-aware URL generation for admin preview, embed and player links
- [x] Rewrote HLS manifests so nested playlists and segments keep token/password query params
- [x] Rewrote DASH MPDs so init/media segment URLs keep token/password query params
- [x] Preserved query strings on `/audio/hls` and `/audio/dash` redirects
- [x] Verified token-protected HLS playback via `ffprobe`
- [x] Verified tokenized DASH init/media URLs return `200`

## Update 2026-03-11 (Adaptive preview and player template library)

Completed in this pass:

- [x] Fixed adaptive preview/player/embed authorization so live admin preview works on ABR streams
- [x] Switched admin HLS preview to master.m3u8 with media-playlist fallback
- [x] Added Public HTTP/HTTPS port controls to Guided Settings
- [x] Added HTTPS port to first-run setup wizard
- [x] Added built-in editable player template library (6+ presets)
- [x] Added adaptive preview smoke test script

﻿# FluxStream - Live Streaming Media Server

## Update 2026-03-11

Completed in this iteration:

- [x] Live ABR connected to the active live HLS pipeline
- [x] Per-stream policy persistence (policy_json)
- [x] Playback authorization for HLS/DASH/player/embed/raw outputs
- [x] Persistent analytics snapshots in SQLite
- [x] Automatic maintenance loop for retention and cleanup
- [x] Health report API and diagnostics API
- [x] Guided admin pages for easy settings, ABR, health and diagnostics
- [x] Simpler topbar server control flow

Remaining work is now mostly tuning and production validation, not missing control surface.
> **Son GÃ¼ncelleme**: 9 Mart 2026
> **Genel Ä°lerleme**: Phase 1 âœ… %100 | Phase 2 âœ… %100 | Phase 3 âœ… %100 | Phase 4 âœ… %100 | Phase 5 âœ… %100

## Planning
- [x] Architecture design and implementation plan
- [x] User review and approval of updated plan (all protocols + embed port config + install process)

---

## Phase 1: Core Foundation âœ… %100
- [x] Go project init (go.mod, directory structure)
- [x] SQLite database (pure Go, modernc.org/sqlite) + config in DB
  - ğŸ“„ `internal/storage/sqlite.go` â€” 8 tablo, full CRUD, player template + user management
  - ğŸ“„ `internal/storage/models.go` â€” Config, User, Stream, Viewer, BannedIP, LogEntry, PlayerTemplate, EmbedDefaults, ServerStats
  - ğŸ“„ `internal/config/config.go` â€” 120+ config entry, kategori bazlÄ± Get/Set/GetInt/GetBool/GetAll/GetByCategory
- [x] RTMP ingest server (pure Go)
  - ğŸ“„ `internal/ingest/rtmp/server.go` â€” TCP listener, connection handler
  - ğŸ“„ `internal/ingest/rtmp/handler.go` â€” connect, createStream, publish, audio/video data
  - ğŸ“„ `internal/ingest/rtmp/handshake.go` â€” C0-S2 el sÄ±kÄ±ÅŸma
  - ğŸ“„ `internal/ingest/rtmp/chunk.go` â€” 4 format type chunk okuma/yazma, extended timestamp
  - ğŸ“„ `internal/ingest/rtmp/amf.go` â€” AMF0 (Number, Boolean, String, Object, Null, ECMAArray)
- [x] FLV demux â†’ MPEG-TS transmux (pure Go)
  - ğŸ“„ `internal/media/container/flv/reader.go` â€” FLV tag parser, frame type detect, sequence header
  - ğŸ“„ `internal/media/container/ts/muxer.go` â€” PAT/PMT generation, PES packaging, CRC32
  - ğŸ“„ `internal/media/packet.go` â€” Packet struct, codec ID'leri (H.264/H.265/VP8/VP9/AV1, AAC/MP3/Opus)
- [x] HLS muxer (M3U8 + TS segments)
  - ğŸ“„ `internal/output/hls/muxer.go` â€” 2sn segment, 6 segment rolling window, keyframe rotation
  - âœ… writePlaylist() â†’ valid M3U8
  - âœ… writeEndPlaylist() â†’ #EXT-X-ENDLIST on close
  - âœ… Proper CORS headers, Content-Type, Cache-Control
- [x] Basic HTTP/HTTPS web server
  - ğŸ“„ `internal/web/server.go` â€” net/http :8844, CORS middleware, tÃ¼m API route'larÄ±
  - âœ… HLS serving with proper headers (`/hls/`)
  - âœ… Player page (`/play/:key`)
  - âœ… Embed page (`/embed/:key`)
  - âœ… crypto/rand stream key generation
  - âœ… Real stats (runtime.MemStats, uptime)
- [x] Setup Wizard (first-run)
  - âœ… 3 adÄ±mlÄ± wizard: HoÅŸgeldin â†’ Admin HesabÄ± â†’ Port AyarlarÄ±
  - âœ… `/api/setup/status` + `/api/setup/complete`
- [x] Dashboard screen
  - âœ… Stat cards (aktif yayÄ±n, izleyici, uptime, bellek)
  - âœ… Aktif yayÄ±nlar listesi
  - âœ… 5sn auto-refresh
  - âœ… Sunucu bilgileri panel
- [x] Custom HTML5 player + preview screen
  - ğŸ“„ `internal/web/player_html.go` â€” Full-page HLS.js player + embed player
  - âœ… Auto-retry (3sn), canlÄ± badge, responsive
  - âœ… Safari native HLS fallback
- [x] Stream manager (lifecycle, fanout)
  - ğŸ“„ `internal/stream/manager.go` â€” OnPublish, OnUnpublish, OnPacket, GetActiveStreams, IsLive, StopAll, GetStats
  - âœ… ActiveStream struct: HLS muxer, packet count, bytes tracking
- [x] Entry point graceful shutdown polish
  - ğŸ“„ `cmd/fluxstream/main.go` â€” âœ… ASCII banner, health check, startup timing, tÃ¼m protokol wiring, padRight helper
  - âœ… Conditional protocol startup (config-based)
  - âœ… Service status display with active protocols

---

## Phase 2: Full Web UI & Management âœ… %100
- [x] All Settings screens (General, Protocols, Outputs, SSL, Security, Storage, Transcode, Users)
  - âœ… Genel Ayarlar â€” server_name, http_port, https_port, language, timezone
  - âœ… Protokoller â€” 8 protokol toggle kartÄ± (RTMP, RTMPS, SRT, RTP, RTSP, WebRTC, MPEG-TS, HTTP Push) + port + ekstra ayarlar
  - âœ… Ã‡Ä±kÄ±ÅŸ FormatlarÄ± â€” 8 Ã§Ä±kÄ±ÅŸ toggle kartÄ± (HLS, LL-HLS, DASH, HTTP-FLV, WebRTC/WHEP, MP3, AAC, Icecast)
  - âœ… SSL/TLS â€” Sertifika yolu + Let's Encrypt ayarlarÄ± + Dosya upload
  - âœ… GÃ¼venlik â€” Stream key zorunlu, token doÄŸrulama, rate limit
  - âœ… Depolama â€” Max GB, otomatik temizlik
  - âœ… Transkod â€” FFmpeg yolu, GPU hÄ±zlandÄ±rma (NVENC/QSV/AMF)
- [x] SSL/TLS file upload (CRT/KEY upload via UI)
  - âœ… `/api/ssl/upload` â€” multipart form POST
  - âœ… `/api/ssl/status` â€” sertifika durumu
  - âœ… UI'da dosya seÃ§me + yÃ¼kleme + durum gÃ¶stergesi
- [x] Stream CRUD UI (create/edit/delete/list)
  - âœ… Streams list â€” tablo, durum badge, delete
  - âœ… Create stream â€” ad + aÃ§Ä±klama, OBS talimatlarÄ±
  - âœ… Stream detail â€” baÄŸlantÄ± bilgileri, embed kodu, yayin bilgileri, canlÄ± Ã¶nizleme
- [x] Embed code generator with configurable port/domain
  - âœ… Embed KodlarÄ± sayfasÄ± â€” tÃ¼m yayÄ±nlar iÃ§in iframe, HLS, Player, RTMP URL
  - âœ… Stream detail sayfasÄ±nda da embed kodu
  - âœ… `/api/embed/:id` API â€” domain/port ayarlÄ± tÃ¼m URL'ler
- [x] Advanced embed generator
  - âœ… GeliÅŸmiÅŸ Embed sayfasÄ± â€” boyut, autoplay, muted, tema seÃ§imi
  - âœ… 7 format sekmesi: iframe, HLS URL, RTMP, RTSP, SRT, Player URL, JS API
  - âœ… CanlÄ± Ã¶nizleme iframe
  - âœ… Kopyala butonu
- [x] Player template editor
  - âœ… `/api/players` CRUD â€” oluÅŸtur, listele, gÃ¼ncelle, sil
  - âœ… Player ÅablonlarÄ± sayfasÄ± â€” kart grid, gÃ¶rsel Ã¶nizleme
  - âœ… Tema, logo (URL/pozisyon/opacity), watermark, show_title, show_live_badge
  - âœ… Arka plan/kontrol/oynat rengi, Ã¶zel CSS
- [x] User management UI
  - âœ… `/api/users` CRUD â€” oluÅŸtur, listele, gÃ¼ncelle, sil
  - âœ… KullanÄ±cÄ±lar sayfasÄ± â€” tablo, ekle/dÃ¼zenle modal
  - âœ… Rol seÃ§imi (admin/editor/viewer), ÅŸifre deÄŸiÅŸtirme

---

## Phase 3: All Ingest Protocols âœ… %100
- [x] RTMPS (TLS wrapped RTMP)
  - ğŸ“„ `internal/ingest/rtmps/server.go` â€” crypto/tls, reuses rtmp.NewHandler
  - âœ… HasValidCerts(), SetCerts() runtime updates
- [x] SRT ingest (pure Go)
  - ğŸ“„ `internal/ingest/srt/server.go` â€” UDP handshake (Induction/Conclusion), session management
  - âœ… Stream ID extraction, MPEG-TS demux, PES parsing, PTS extraction
- [x] RTP ingest (UDP, H.264 NAL assembly)
  - ğŸ“„ `internal/ingest/rtp/server.go` â€” SSRC-based sessions, FU-A/STAP-A, header extensions
  - âœ… AAC (PT 97) + Opus (PT 111) audio
- [x] RTSP ingest (push from IP cameras)
  - ğŸ“„ `internal/ingest/rtsp/server.go` â€” ANNOUNCE/SETUP/RECORD/TEARDOWN
  - âœ… SDP parsing, TCP interleaved + UDP, H.264 reassembly
- [x] WebRTC/WHIP ingest (browser-based)
  - ğŸ“„ `internal/ingest/webrtc/server.go` â€” WHIP HTTP endpoint, SDP offer/answer
  - âœ… ICE trickle (PATCH), session management, CORS
- [x] MPEG-TS UDP ingest (multicast/unicast)
  - ğŸ“„ `internal/ingest/mpegts/server.go` â€” PAT/PMT parsing, PES demux, PID tracking
  - âœ… H.264/H.265/AAC/MP3 stream type detection, keyframe detection
- [x] HTTP Push ingest
  - ğŸ“„ `internal/ingest/httppush/server.go` â€” PUT/POST with Bearer auth
  - âœ… Chunked + non-chunked, TS packet processing, status endpoint
- [x] Shared StreamHandler interface
  - ğŸ“„ `internal/ingest/handler.go` â€” OnPublish, OnUnpublish, OnPacket

---

## Phase 4: All Output Formats âœ… %100
- [x] DASH + CMAF output
  - ğŸ“„ `internal/output/dash/muxer.go` â€” fMP4 segment, MPD manifest, init.mp4, SegmentTimeline
  - âœ… Video + Audio AdaptationSet, rolling window, styp/moof/mdat
- [x] LL-HLS (Low Latency HLS)
  - ğŸ“„ `internal/output/hls/ll_muxer.go` â€” Partial segments (~200ms), ll.m3u8
  - âœ… EXT-X-PART, EXT-X-PRELOAD-HINT, EXT-X-SERVER-CONTROL, PART-HOLD-BACK
  - âœ… Segment concatenation, rolling window
- [x] HTTP-FLV output
  - ğŸ“„ `internal/output/flv/server.go` â€” Chunked transfer encoding, FLV tag builder
  - âœ… FLV header, tag header (video/audio/meta), subscriber-based streaming
- [x] WebRTC/WHEP output (sub-second latency)
  - ğŸ“„ `internal/output/webrtc/server.go` â€” WHEPServer, SDP offer/answer, ICE trickle
  - âœ… POST /whep/play/:key, PATCH /whep/ice/:id, DELETE /whep/session/:id
  - âœ… Session management, cleanup, status API
- [x] RTMP Relay (re-stream to YouTube/Twitch/etc.)
  - ğŸ“„ `internal/output/relay/manager.go` â€” Target management, RTMP handshake, AMF0 connect/publish
  - âœ… Add/Remove/Start/Stop targets, relay loop, URL parsing
- [x] RTSP output server
  - ğŸ“„ `internal/output/rtsp/server.go` â€” OPTIONS/DESCRIBE/SETUP/PLAY/TEARDOWN
  - âœ… SDP generation, RTP over UDP, session management
- [x] RTP output sender
  - ğŸ“„ `internal/output/rtp/sender.go` â€” Target management, RTP packet builder
  - âœ… H.264 (PT 96) + AAC (PT 97), SSRC, seq num, clock rate conversion
- [x] SRT output
  - ğŸ“„ `internal/output/srt/server.go` â€” UDP handshake, MPEG-TS mux, SRT data packet
  - âœ… Client management, stream ID, keep-alive
- [x] MPEG-TS UDP output
  - ğŸ“„ `internal/output/mpegts/sender.go` â€” UDP unicast/multicast, 7 TS packet bundling
  - âœ… Target management, multicast TTL, auto-target
- [x] fMP4/MP4 progressive streaming
  - ğŸ“„ `internal/output/mp4/muxer.go` â€” HandleFMP4, ftyp/moov/moof/mdat, GOP fragmentation
  - âœ… Chunked transfer, init segment, subscriber-based
- [x] WebM streaming output
  - âœ… HandleWebM in mp4/muxer.go â€” EBML header, Segment, Info, Tracks (VP9+Opus)
  - âœ… Cluster on keyframe, SimpleBlock, live unknown-size segments
- [x] Audio: MP3/Icecast, AAC, Opus, OGG, HLS Audio, DASH Audio, FLAC, WAV
  - ğŸ“„ `internal/output/audio/server.go` â€” 9 audio format handler
  - âœ… HandleMP3 (Icecast headers), HandleAAC (ADTS), HandleOpus (OGG), HandleOGG
  - âœ… HandleWAV (PCM), HandleFLAC, HandleHLSAudio, HandleDASHAudio
  - âœ… Common streamAudio() helper, extractKey(), subscriber-based
- [x] Stream Manager subscriber/fanout system
  - âœ… OutputSubscriber struct: ID, PacketC, Done channels
  - âœ… Subscribe/Unsubscribe methods, non-blocking fan-out, slow subscriber drop
- [x] Output routes registered on HTTP server
  - âœ… /flv/, /whep/, /dash/, /mp4/, /webm/, /audio/mp3|aac|opus|ogg|wav|flac|hls|dash/
  - âœ… RegisterHandler() + RegisterOutputDir() methods on web.Server

---

## Phase 5: Advanced Features & Polish âœ… %100
- [x] FFmpeg transcoding (multi-quality ABR)
  - ğŸ“„ `internal/transcode/manager.go` â€” FFmpeg entegrasyonu, Ã§oklu profil, GPU tespit
  - âœ… DefaultProfiles: 1080p/720p/480p/360p ABR
  - âœ… GPU hÄ±zlandÄ±rma: NVENC, QSV, AMF, VAAPI otomatik tespit
  - âœ… HLS master playlist Ã§Ä±kÄ±ÅŸÄ±, var_stream_map
  - âœ… Job yÃ¶netimi: Start/Stop/GetJobs/GetStatus
- [x] Stream recording / DVR (MP4/TS/MKV/FLV)
  - ğŸ“„ `internal/recording/manager.go` â€” KayÄ±t baÅŸlatma/durdurma, dosya yÃ¶netimi
  - âœ… TS + FLV format kayÄ±t, subscriber tabanlÄ± paket yazma
  - âœ… Maks sÃ¼re limiti (24h), otomatik durdurma
  - âœ… DVR: dosya listeleme, aÃ§ma/playback, silme
  - âœ… Eski kayÄ±t temizleme (CleanupOld)
- [x] Analytics dashboard
  - ğŸ“„ `internal/analytics/tracker.go` â€” Event izleme, dashboard verisi
  - âœ… Viewer join/leave tracking, format/Ã¼lke bazlÄ± istatistik
  - âœ… Saatlik timeline, peak concurrent, toplam bandwidth
  - âœ… Stream bazlÄ± istatistikler, top streams sÄ±ralamasÄ±
  - âœ… 24 saatlik event pruning
- [x] Security features (2FA, tokens, rate limiting, IP ban)
  - ğŸ“„ `internal/security/security.go` â€” 4 gÃ¼venlik modÃ¼lÃ¼
  - âœ… HMAC-SHA256 token oluÅŸturma/doÄŸrulama, sÃ¼re bazlÄ±
  - âœ… IP tabanlÄ± rate limiting middleware
  - âœ… IP ban listesi (sÃ¼reli/kalÄ±cÄ±) middleware
  - âœ… TOTP 2FA (secret oluÅŸturma, kod doÄŸrulama, time-step tolerans)
- [x] Windows service installer + Linux systemd + Build scripts
  - ğŸ“„ `internal/service/installer.go` â€” TÃ¼m platform desteÄŸi
  - âœ… NSIS installer script (Windows service create/delete)
  - âœ… systemd unit file (gÃ¼venlik sÄ±kÄ±laÅŸtÄ±rma: NoNewPrivileges, ProtectSystem)
  - âœ… Debian .deb control dosyasÄ±
  - âœ… Cross-platform build script (windows/linux/darwin, amd64/arm64)
- [x] Phase 5 API entegrasyonu
  - âœ… /api/recordings (GET/POST), /api/recordings/stop/:id, /api/recordings/files/:key
  - âœ… /api/analytics, /api/analytics/stream/:key
  - âœ… /api/security/token/generate, /api/security/bans (GET/POST/DELETE)
  - âœ… /api/transcode/status, /api/transcode/jobs
  - âœ… Graceful shutdown: recManager.StopAll(), tcManager.StopAll()

---

## Dosya Listesi (GÃ¼ncel)

| Dosya | SatÄ±r | Durum |
|-------|-------|-------|
| `cmd/fluxstream/main.go` | ~350 | âœ… TamamlandÄ± |
| `internal/config/config.go` | ~200 | âœ… TamamlandÄ± |
| `internal/storage/sqlite.go` | ~485 | âœ… TamamlandÄ± |
| `internal/storage/models.go` | ~114 | âœ… TamamlandÄ± |
| `internal/stream/manager.go` | ~230 | âœ… TamamlandÄ± |
| `internal/ingest/handler.go` | ~11 | âœ… TamamlandÄ± |
| `internal/ingest/rtmp/server.go` | ~47 | âœ… TamamlandÄ± |
| `internal/ingest/rtmp/handler.go` | ~329 | âœ… TamamlandÄ± |
| `internal/ingest/rtmp/handshake.go` | ~52 | âœ… TamamlandÄ± |
| `internal/ingest/rtmp/chunk.go` | ~326 | âœ… TamamlandÄ± |
| `internal/ingest/rtmp/amf.go` | ~202 | âœ… TamamlandÄ± |
| `internal/ingest/rtmps/server.go` | ~81 | âœ… TamamlandÄ± |
| `internal/ingest/srt/server.go` | ~425 | âœ… TamamlandÄ± |
| `internal/ingest/rtp/server.go` | ~323 | âœ… TamamlandÄ± |
| `internal/ingest/rtsp/server.go` | ~443 | âœ… TamamlandÄ± |
| `internal/ingest/webrtc/server.go` | ~327 | âœ… TamamlandÄ± |
| `internal/ingest/mpegts/server.go` | ~388 | âœ… TamamlandÄ± |
| `internal/ingest/httppush/server.go` | ~241 | âœ… TamamlandÄ± |
| `internal/media/packet.go` | ~78 | âœ… TamamlandÄ± |
| `internal/media/container/flv/reader.go` | ~130 | âœ… TamamlandÄ± |
| `internal/media/container/ts/muxer.go` | ~274 | âœ… TamamlandÄ± |
| `internal/output/hls/muxer.go` | ~246 | âœ… TamamlandÄ± |
| `internal/output/hls/ll_muxer.go` | ~280 | âœ… TamamlandÄ± |
| `internal/output/dash/muxer.go` | ~310 | âœ… TamamlandÄ± |
| `internal/output/flv/server.go` | ~155 | âœ… TamamlandÄ± |
| `internal/output/webrtc/server.go` | ~220 | âœ… TamamlandÄ± |
| `internal/output/relay/manager.go` | ~340 | âœ… TamamlandÄ± |
| `internal/output/rtsp/server.go` | ~310 | âœ… TamamlandÄ± |
| `internal/output/rtp/sender.go` | ~215 | âœ… TamamlandÄ± |
| `internal/output/srt/server.go` | ~215 | âœ… TamamlandÄ± |
| `internal/output/mpegts/sender.go` | ~200 | âœ… TamamlandÄ± |
| `internal/output/mp4/muxer.go` | ~350 | âœ… TamamlandÄ± |
| `internal/output/audio/server.go` | ~410 | âœ… TamamlandÄ± |
| `internal/recording/manager.go` | ~295 | âœ… TamamlandÄ± |
| `internal/transcode/manager.go` | ~295 | âœ… TamamlandÄ± |
| `internal/analytics/tracker.go` | ~210 | âœ… TamamlandÄ± |
| `internal/security/security.go` | ~295 | âœ… TamamlandÄ± |
| `internal/service/installer.go` | ~190 | âœ… TamamlandÄ± |
| `internal/web/server.go` | ~820 | âœ… TamamlandÄ± |
| `internal/web/admin_html.go` | ~963 | âœ… TamamlandÄ± (Full SPA) |
| `internal/web/player_html.go` | ~137 | âœ… TamamlandÄ± |

**Toplam**: ~11,300+ satÄ±r Go kodu, 41 dosya
**Build**: `go build ./cmd/fluxstream/` âœ… HatasÄ±z | `go vet ./...` âœ… HatasÄ±z
**BaÄŸÄ±mlÄ±lÄ±k**: Sadece `modernc.org/sqlite` (pure Go, CGo yok)

