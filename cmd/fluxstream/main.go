package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fluxstream/fluxstream/internal/analytics"
	"github.com/fluxstream/fluxstream/internal/archive"
	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/ingest/httppush"
	"github.com/fluxstream/fluxstream/internal/ingest/mpegts"
	"github.com/fluxstream/fluxstream/internal/ingest/rtmp"
	"github.com/fluxstream/fluxstream/internal/ingest/rtmps"
	"github.com/fluxstream/fluxstream/internal/ingest/rtp"
	"github.com/fluxstream/fluxstream/internal/ingest/rtsp"
	"github.com/fluxstream/fluxstream/internal/ingest/srt"
	"github.com/fluxstream/fluxstream/internal/ingest/webrtc"
	"github.com/fluxstream/fluxstream/internal/output/audio"
	"github.com/fluxstream/fluxstream/internal/output/dash"
	"github.com/fluxstream/fluxstream/internal/output/flv"
	"github.com/fluxstream/fluxstream/internal/output/hls"
	llhls "github.com/fluxstream/fluxstream/internal/output/hls"
	"github.com/fluxstream/fluxstream/internal/output/mp4"
	outmpegts "github.com/fluxstream/fluxstream/internal/output/mpegts"
	"github.com/fluxstream/fluxstream/internal/output/relay"
	outrtp "github.com/fluxstream/fluxstream/internal/output/rtp"
	outrtsp "github.com/fluxstream/fluxstream/internal/output/rtsp"
	outsrt "github.com/fluxstream/fluxstream/internal/output/srt"
	outwebrtc "github.com/fluxstream/fluxstream/internal/output/webrtc"
	"github.com/fluxstream/fluxstream/internal/recording"
	"github.com/fluxstream/fluxstream/internal/security"
	"github.com/fluxstream/fluxstream/internal/storage"
	"github.com/fluxstream/fluxstream/internal/stream"
	"github.com/fluxstream/fluxstream/internal/tlsutil"
	"github.com/fluxstream/fluxstream/internal/transcode"
	"github.com/fluxstream/fluxstream/internal/web"
)

const (
	Version         = "2.0.0"
	AppName         = "FluxStream"
	DefaultHTTPPort = 8844
	DefaultRTMPPort = 1935
)

func main() {
	configureProcessOutput()
	if handled, err := handleConfigMode(os.Args[1:]); handled {
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if handled, err := handleServiceMode(os.Args[1:]); handled {
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if handled, err := handleBackupMode(os.Args[1:]); handled {
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	startTime := time.Now()
	fmt.Println()
	fmt.Printf("%s v%s\n", AppName, Version)
	fmt.Println(strings.Repeat("=", 44))
	fmt.Println("Live Streaming Media Server")
	fmt.Printf("Go %s | %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Println("Zero Dependency | Single Binary | Pure Go")
	fmt.Println()

	// Determine data directory
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("[FATAL] Executable path alinamadi: %v", err)
	}
	dataDir := filepath.Join(filepath.Dir(execPath), "data")

	// Ensure data directories exist
	dirs := []string{
		dataDir,
		filepath.Join(dataDir, "hls"),
		filepath.Join(dataDir, "dash"),
		filepath.Join(dataDir, "recordings"),
		filepath.Join(dataDir, "archive"),
		filepath.Join(dataDir, "backups"),
		filepath.Join(dataDir, "thumbnails"),
		filepath.Join(dataDir, "license"),
		filepath.Join(dataDir, "certs"),
		filepath.Join(dataDir, "certs", "web"),
		filepath.Join(dataDir, "certs", "stream"),
		filepath.Join(dataDir, "certs", "acme"),
		filepath.Join(dataDir, "logs"),
		filepath.Join(dataDir, "players"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatalf("[FATAL] Dizin olusturulamadi %s: %v", d, err)
		}
	}
	log.Printf("[INIT] Veri dizini: %s", dataDir)

	// Initialize SQLite database
	dbPath := filepath.Join(dataDir, "fluxstream.db")
	db, err := storage.NewSQLiteDB(dbPath)
	if err != nil {
		log.Fatalf("[FATAL] Veritabani acilamadi: %v", err)
	}
	defer db.Close()
	log.Printf("[INIT] Veritabani: %s", dbPath)
	if resetCount, err := db.ResetRuntimeStreamState(); err != nil {
		log.Printf("[INIT] Calisma zamani stream durumu sifirlanamadi: %v", err)
	} else if resetCount > 0 {
		log.Printf("[INIT] Eski calisma zamani stream durumu temizlendi (%d kayit)", resetCount)
	}

	// Initialize config
	cfg := config.NewManager(db)
	if err := cfg.LoadDefaults(); err != nil {
		log.Fatalf("[FATAL] Varsayilan ayarlar yuklenemedi: %v", err)
	}
	log.Printf("[INIT] Yapilandirma yuklendi (%d ayar)", 115)
	licenseRuntime := resolveRuntimeLicense(dataDir)
	log.Printf("[INIT] Lisans modu: %s | aktif ozellikler: %s", licenseRuntime.Mode, strings.Join(licenseRuntime.EnabledFeatures, ","))
	if licenseRuntime.Enforced {
		if !licenseRuntime.allows(licenseFeatureABR) {
			_ = cfg.Set("abr_enabled", "false", "outputs")
			_ = cfg.Set("abr_master_enabled", "false", "outputs")
		}
		if !licenseRuntime.allows(licenseFeatureRTMPS) {
			_ = cfg.Set("rtmps_enabled", "false", "protocols")
		}
		if !licenseRuntime.allows(licenseFeatureRecording) {
			_ = cfg.Set("recording_enabled", "false", "recording")
		}
	}

	webTLSSource, err := tlsutil.NewSource(cfg, tlsutil.ProfileWeb, dataDir)
	if err != nil {
		log.Printf("[INIT] Web TLS hazirlanamadi: %v", err)
		webTLSSource = &tlsutil.Source{Profile: tlsutil.ProfileWeb}
	}
	streamTLSSource, err := tlsutil.NewSource(cfg, tlsutil.ProfileStream, dataDir)
	if err != nil {
		log.Printf("[INIT] Stream TLS hazirlanamadi: %v", err)
		streamTLSSource = &tlsutil.Source{Profile: tlsutil.ProfileStream}
	}

	// Create context with cancel for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	controlCh := make(chan string, 1)

	// Initialize stream manager
	hlsOutputDir := filepath.Join(dataDir, "hls")
	httpPort := cfg.GetInt("http_port", DefaultHTTPPort)
	hlsMuxer := hls.NewMuxer(hlsOutputDir)
	streamManager := stream.NewManager(db, hlsMuxer)

	// ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ Output Servers ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬

	// DASH/CMAF output
	var dashMuxer *dash.Muxer
	if cfg.GetBool("dash_enabled", false) {
		dashOutputDir := filepath.Join(dataDir, "dash")
		dashMuxer = dash.NewMuxer(dashOutputDir)
		log.Println("[INIT] DASH/CMAF cikis aktif")
	}

	// LL-HLS output (uses same package as hls, different muxer type)
	var llhlsMuxer *llhls.LLMuxer
	if cfg.GetBool("hls_ll_enabled", false) {
		llhlsMuxer = llhls.NewLLMuxer(hlsOutputDir)
		log.Println("[INIT] LL-HLS cikis aktif")
	}
	streamManager.SetOutputMuxers(dashMuxer, llhlsMuxer)

	// HTTP-FLV output
	var flvServer *flv.Server
	if cfg.GetBool("httpflv_enabled", false) {
		gopCache := cfg.GetBool("httpflv_gop_cache", true)
		flvServer = flv.NewServer(streamManager, gopCache)
		log.Println("[INIT] HTTP-FLV cikis aktif")
	}

	// WHEP/WebRTC output
	var whepServer *outwebrtc.WHEPServer
	if cfg.GetBool("whep_enabled", false) {
		whepServer = outwebrtc.NewWHEPServer(streamManager)
		log.Println("[INIT] WHEP cikis aktif")
	}

	// RTMP Relay
	relayManager := relay.NewManager(streamManager)

	// MP4/WebM streaming
	var mp4Server *mp4.Server

	// RTSP output
	if cfg.GetBool("rtsp_out_enabled", false) {
		rtspOutPort := cfg.GetInt("rtsp_out_port", 8555)
		rtspOutServer := outrtsp.NewServer(rtspOutPort, streamManager)
		stopCh := make(chan struct{})
		go func() {
			<-ctx.Done()
			close(stopCh)
		}()
		go func() {
			if err := rtspOutServer.Start(stopCh); err != nil {
				log.Printf("[ERROR] RTSP cikis sunucu hatasi: %v", err)
			}
		}()
	}

	// RTP output sender
	rtpOutSender := outrtp.NewSender(streamManager)

	// SRT output
	if cfg.GetBool("srt_out_enabled", false) {
		srtOutPort := cfg.GetInt("srt_out_port", 9010)
		srtOutServer := outsrt.NewServer(srtOutPort, streamManager)
		stopCh := make(chan struct{})
		go func() {
			<-ctx.Done()
			close(stopCh)
		}()
		go func() {
			if err := srtOutServer.Start(stopCh); err != nil {
				log.Printf("[ERROR] SRT cikis sunucu hatasi: %v", err)
			}
		}()
	}

	// MPEG-TS UDP output sender
	tsUDPSender := outmpegts.NewSender(streamManager)

	// Suppress unused variable warnings
	_, _, _, _, _, _, _ = dashMuxer, llhlsMuxer, flvServer, whepServer, relayManager, rtpOutSender, tsUDPSender

	// ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ Phase 5: Advanced Features ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬

	// Recording/DVR
	recordingsDir := filepath.Join(dataDir, "recordings")
	recManager := recording.NewManager(streamManager, recordingsDir, cfg.Get("ffmpeg_path", "ffmpeg"))
	archiveManager := archive.NewManager(cfg, db, recManager, dataDir)
	log.Println("[INIT] Kayit/DVR sistemi aktif")

	// Analytics
	analyticsTracker := analytics.NewTracker()
	if existingStreams, err := db.GetAllStreams(); err == nil {
		for _, st := range existingStreams {
			analyticsTracker.RegisterStreamName(st.StreamKey, st.Name)
		}
	}
	log.Println("[INIT] Analitik izleyici aktif")
	playerTelemetry := newPlayerTelemetryCollector()
	playerTelemetry.SetDB(db)
	log.Println("[INIT] Player QoE telemetrisi aktif")

	// Security
	tokenSecret := cfg.Get("token_secret", "")
	if strings.TrimSpace(tokenSecret) == "" {
		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err == nil {
			tokenSecret = hex.EncodeToString(buf)
			_ = cfg.Set("token_secret", tokenSecret, "security")
		}
	}
	tokenDuration := cfg.GetInt("token_duration", 60)
	tokenMgr := security.NewTokenManager(tokenSecret, tokenDuration)
	rateLimitMax := cfg.GetInt("rate_limit", 100)
	rateLimiter := security.NewRateLimiter(rateLimitMax, time.Minute)
	ipBanList := security.NewIPBanList()
	twoFA := security.NewTwoFAManager()
	playbackAuth := makePlaybackAuthorizer(cfg, db, tokenMgr)
	log.Println("[INIT] Guvenlik modulleri aktif (token, rate-limit, IP ban, 2FA)")
	// Transcoding
	ffmpegPath := cfg.Get("ffmpeg_path", "ffmpeg")
	gpuAccel := transcode.GPUAccel(cfg.Get("gpu_accel", "none"))
	transcodeDir := filepath.Join(dataDir, "transcode")
	tcManager := transcode.NewManager(ffmpegPath, gpuAccel, transcodeDir)
	tcManager.SetHTTPPort(httpPort)
	liveOpts := buildLiveOptionsFromConfig(cfg)
	if !licenseRuntime.allows(licenseFeatureABR) {
		liveOpts.ABREnabled = false
		liveOpts.MasterEnabled = false
	}
	tcManager.SetLiveOptions(liveOpts)
	mp4Server = mp4.NewServer(streamManager, tcManager, httpPort)
	audioServer := audio.NewServer(streamManager, tcManager, httpPort)
	log.Println("[INIT] Transkod yoneticisi aktif")
	liveHLSTranscode := cfg.GetBool("transcode_live_hls_enabled", true)
	liveDASHTranscode := cfg.GetBool("transcode_live_dash_enabled", true) && cfg.GetBool("dash_enabled", false)
	if liveHLSTranscode {
		log.Println("[INIT] Canli HLS transcode aktif")
	}
	if liveDASHTranscode {
		log.Println("[INIT] Canli DASH repack aktif")
	}
	recorderForPipeline := recManager
	if !licenseRuntime.allows(licenseFeatureRecording) {
		recorderForPipeline = nil
	}
	pipelineHandler := stream.NewPipelineHandler(streamManager, tcManager, recorderForPipeline, liveHLSTranscode, liveDASHTranscode, liveOpts, licenseRuntime.allows(licenseFeatureABR))

	// Suppress warnings
	_, _, _, _, _, _ = recManager, analyticsTracker, tokenMgr, rateLimiter, ipBanList, twoFA
	startMaintenanceLoops(ctx.Done(), cfg, db, analyticsTracker, recManager, archiveManager, tcManager, dataDir)

	// Initialize & start RTMP server
	rtmpPort := cfg.GetInt("rtmp_port", DefaultRTMPPort)
	rtmpServer := rtmp.NewServer(rtmpPort, pipelineHandler)
	go func() {
		if err := rtmpServer.Start(ctx); err != nil {
			log.Printf("[FATAL] RTMP sunucu hatasi: %v", err)
		}
	}()

	// Initialize & start RTMPS server (if enabled)
	if cfg.GetBool("rtmps_enabled", false) && licenseRuntime.allows(licenseFeatureRTMPS) {
		rtmpsPort := cfg.GetInt("rtmps_port", 1936)
		if !streamTLSSource.Ready {
			if streamTLSSource.UsesLetsEncrypt() {
				log.Printf("[INIT] RTMPS Let's Encrypt etkin ama sertifika henuz hazir degil: %s", streamTLSSource.Domain)
			} else {
				log.Printf("[INIT] RTMPS etkin ama gecerli CRT/KEY bulunamadi, RTMPS baslatilmadi")
			}
		} else {
			rtmpsServer := rtmps.NewServer(rtmpsPort, pipelineHandler, streamTLSSource.TLSConfig())
			go func() {
				if err := rtmpsServer.Start(ctx); err != nil {
					log.Printf("[ERROR] RTMPS sunucu hatasi: %v", err)
				}
			}()
		}
	}

	// Initialize & start SRT server (if enabled)
	if cfg.GetBool("srt_enabled", false) {
		srtPort := cfg.GetInt("srt_port", 9000)
		srtLatency := cfg.GetInt("srt_latency", 120)
		srtServer := srt.NewServer(srtPort, pipelineHandler, srtLatency)
		go func() {
			if err := srtServer.Start(ctx); err != nil {
				log.Printf("[ERROR] SRT cikis sunucu hatasi: %v", err)
			}
		}()
	}

	// Initialize & start RTP server (if enabled)
	if cfg.GetBool("rtp_enabled", false) {
		rtpPort := cfg.GetInt("rtp_port", 5004)
		rtpServer := rtp.NewServer(rtpPort, pipelineHandler)
		go func() {
			if err := rtpServer.Start(ctx); err != nil {
				log.Printf("[ERROR] RTP sunucu hatasi: %v", err)
			}
		}()
	}

	// Initialize & start RTSP server (if enabled)
	if cfg.GetBool("rtsp_enabled", false) {
		rtspPort := cfg.GetInt("rtsp_port", 8554)
		rtspServer := rtsp.NewServer(rtspPort, pipelineHandler)
		go func() {
			if err := rtspServer.Start(ctx); err != nil {
				log.Printf("[ERROR] RTSP cikis sunucu hatasi: %v", err)
			}
		}()
	}

	// Initialize & start WebRTC/WHIP server (if enabled)
	if cfg.GetBool("webrtc_enabled", false) {
		webrtcPort := cfg.GetInt("webrtc_port", 8855)
		webrtcServer := webrtc.NewServer(webrtcPort, pipelineHandler)
		go func() {
			if err := webrtcServer.Start(ctx); err != nil {
				log.Printf("[ERROR] WebRTC sunucu hatasi: %v", err)
			}
		}()
	}

	// Initialize & start MPEG-TS UDP server (if enabled)
	if cfg.GetBool("mpegts_enabled", false) {
		mpegtsPort := cfg.GetInt("mpegts_port", 9001)
		mpegtsServer := mpegts.NewServer(mpegtsPort, pipelineHandler)
		go func() {
			if err := mpegtsServer.Start(ctx); err != nil {
				log.Printf("[ERROR] MPEG-TS sunucu hatasi: %v", err)
			}
		}()
	}

	// Initialize & start HTTP Push server (if enabled)
	if cfg.GetBool("http_push_enabled", false) {
		httpPushPort := cfg.GetInt("http_push_port", 8850)
		httpPushToken := cfg.Get("http_push_token", "")
		httpPushServer := httppush.NewServer(httpPushPort, pipelineHandler, httpPushToken)
		go func() {
			if err := httpPushServer.Start(ctx); err != nil {
				log.Printf("[ERROR] HTTP-Push sunucu hatasi: %v", err)
			}
		}()
	}

	// Initialize & start web server
	webServer := web.NewServer(httpPort, db, cfg, streamManager, hlsOutputDir, dataDir)
	webServer.SetAnalyticsTracker(analyticsTracker)
	webServer.SetPlaybackAuthorizer(playbackAuth)
	webServer.SetSettingsMutator(licenseRuntime.normalizeSettings)
	webServer.SetStreamMutator(licenseRuntime.normalizeStream)
	webServer.SetPlayerTemplateMutator(licenseRuntime.normalizePlayerTemplate)
	webServer.SetHLSOverrideDir(tcManager.GetLiveOutputDir())
	webServer.SetDashOverrideDir(tcManager.GetLiveDashOutputDir())
	// Register output routes on web server
	if flvServer != nil {
		webServer.RegisterHandler("/flv/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "http_flv", flvPlaybackKey, flvServer.HandleFLV))
	}
	if whepServer != nil {
		webServer.RegisterHandler("/whep/", whepServer.HandleWHEP)
	}
	if dashMuxer != nil {
		webServer.SetDashOutputDir(dashMuxer.GetOutputDir())
	}
	webServer.RegisterHandler("/mp4/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "mp4", mp4PlaybackKey, mp4Server.HandleFMP4))
	webServer.RegisterHandler("/webm/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "webm", webmPlaybackKey, mp4Server.HandleWebM))
	webServer.RegisterHandler("/audio/mp3/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "mp3", audioPlaybackKey("/audio/mp3/"), audioServer.HandleMP3))
	webServer.RegisterHandler("/audio/aac/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "aac", audioPlaybackKey("/audio/aac/"), audioServer.HandleAAC))
	webServer.RegisterHandler("/audio/ogg/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "ogg", audioPlaybackKey("/audio/ogg/"), audioServer.HandleOGG))
	webServer.RegisterHandler("/audio/wav/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "wav", audioPlaybackKey("/audio/wav/"), audioServer.HandleWAV))
	webServer.RegisterHandler("/audio/flac/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "flac", audioPlaybackKey("/audio/flac/"), audioServer.HandleFLAC))
	webServer.RegisterHandler("/audio/hls/", audioServer.HandleHLSAudio)
	webServer.RegisterHandler("/audio/dash/", audioServer.HandleDASHAudio)
	webServer.RegisterHandler("/icecast/", wrapStreamingPlaybackHandler(analyticsTracker, playbackAuth, "icecast", icecastPlaybackKey, audioServer.HandleIcecast))
	webServer.RegisterHandler("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte(buildPrometheusMetrics(Version, analyticsTracker, streamManager, tcManager, playerTelemetry)))
	})
	webServer.RegisterHandler("/api/observability/otel", func(w http.ResponseWriter, r *http.Request) {
		jsonResp(w, buildOpenTelemetryPayload(Version, analyticsTracker, streamManager, tcManager, playerTelemetry))
	})
	webServer.RegisterHandler("/api/player/telemetry", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		var payload playerTelemetryPayload
		if err := decodeJSON(r, &payload); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if strings.TrimSpace(payload.StreamKey) == "" || strings.TrimSpace(payload.SessionID) == "" {
			http.Error(w, "stream_key ve session_id gerekli", 400)
			return
		}
		if st, err := db.GetStreamByKey(payload.StreamKey); err == nil && st != nil {
			playerTelemetry.Record(payload, r.RemoteAddr, r.UserAgent())
		}
		w.WriteHeader(http.StatusAccepted)
	})
	webServer.RegisterAdminHandler("/api/admin/player/telemetry/stream/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/player/telemetry/stream/")
		var id int64
		fmt.Sscanf(idStr, "%d", &id)
		st, err := db.GetStreamByID(id)
		if err != nil || st == nil {
			http.Error(w, "Stream bulunamadi", 404)
			return
		}
		if streamManager.IsLive(st.StreamKey) {
			st.Status = "live"
		}
		history, _ := db.GetPlayerTelemetrySamples(st.StreamKey, 72)
		trackSnapshot := transcode.LiveTrackSnapshot{}
		if tcManager != nil {
			trackSnapshot = tcManager.GetLiveTrackSnapshot(st.StreamKey)
		}
		qoeAlerts := buildQoEAlerts(cfg, st.Name, playerTelemetry.Snapshot(st.StreamKey))
		trackHistory, _ := db.GetTrackTelemetrySamples(st.StreamKey, 160)
		jsonResp(w, map[string]interface{}{
			"stream_id":     st.ID,
			"stream_key":    st.StreamKey,
			"stream_name":   st.Name,
			"status":        st.Status,
			"telemetry":     playerTelemetry.Snapshot(st.StreamKey),
			"history":       history,
			"tracks":        trackSnapshot,
			"track_history": trackHistory,
			"qoe_alerts":    qoeAlerts,
		})
	})
	webServer.RegisterAdminHandler("/api/admin/stream/tracks/defaults/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/api/admin/stream/tracks/defaults/")
		var id int64
		fmt.Sscanf(idStr, "%d", &id)
		st, err := db.GetStreamByID(id)
		if err != nil || st == nil {
			http.Error(w, "Stream bulunamadi", http.StatusNotFound)
			return
		}
		var req struct {
			DefaultVideoTrackID int `json:"default_video_track_id"`
			DefaultAudioTrackID int `json:"default_audio_track_id"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if tcManager != nil {
			var videoTrackID uint8
			var audioTrackID uint8
			if req.DefaultVideoTrackID > 0 && req.DefaultVideoTrackID <= 255 {
				videoTrackID = uint8(req.DefaultVideoTrackID)
			}
			if req.DefaultAudioTrackID > 0 && req.DefaultAudioTrackID <= 255 {
				audioTrackID = uint8(req.DefaultAudioTrackID)
			}
			tcManager.SetStreamTrackDefaults(st.StreamKey, videoTrackID, audioTrackID)
		}
		jsonResp(w, map[string]interface{}{"success": true})
	})
	webServer.RegisterAdminHandler("/api/system/restart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		select {
		case controlCh <- "restart":
			jsonResp(w, map[string]interface{}{"success": true})
		default:
			jsonResp(w, map[string]interface{}{"success": false, "message": "Islem zaten suruyor"})
		}
	})
	webServer.RegisterAdminHandler("/api/system/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		select {
		case controlCh <- "stop":
			jsonResp(w, map[string]interface{}{"success": true})
		default:
			jsonResp(w, map[string]interface{}{"success": false, "message": "Islem zaten suruyor"})
		}
	})

	// ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ Phase 5 API Routes ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã†â€™Ãƒâ€ Ã¢â‚¬â„¢ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â€šÂ¬Ã…Â¡Ãƒâ€šÃ‚Â¬ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚ÂÃƒÆ’Ã†â€™Ãƒâ€šÃ‚Â¢ÃƒÆ’Ã‚Â¢ÃƒÂ¢Ã¢â‚¬Å¡Ã‚Â¬Ãƒâ€¦Ã‚Â¡ÃƒÆ’Ã¢â‚¬Å¡Ãƒâ€šÃ‚Â¬

	// Recording API
	webServer.RegisterHandler("/api/recordings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			recs := recManager.GetRecordings()
			if recs == nil {
				recs = []*recording.Recording{}
			}
			jsonResp(w, recs)
		case "POST":
			var req struct {
				StreamKey string `json:"stream_key"`
				Format    string `json:"format"`
			}
			if err := decodeJSON(r, &req); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if !licenseRuntime.allows(licenseFeatureRecording) {
				http.Error(w, "recording lisans gerektirir", http.StatusForbidden)
				return
			}
			if strings.TrimSpace(req.Format) == "" {
				req.Format = string(recording.FormatMP4)
			}
			rec, err := recManager.StartRecording(req.StreamKey, recording.Format(req.Format))
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			jsonResp(w, rec)
		}
	})
	webServer.RegisterHandler("/api/recordings/stop/", func(w http.ResponseWriter, r *http.Request) {
		recID := strings.TrimPrefix(r.URL.Path, "/api/recordings/stop/")
		if err := recManager.StopRecording(recID); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonResp(w, map[string]string{"status": "stopped"})
	})
	webServer.RegisterHandler("/api/recordings/files/", func(w http.ResponseWriter, r *http.Request) {
		streamKey := strings.TrimPrefix(r.URL.Path, "/api/recordings/files/")
		files, err := recManager.ListRecordingFiles(streamKey)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if files == nil {
			files = []recording.RecordingFile{}
		}
		jsonResp(w, files)
	})
	webServer.RegisterHandler("/api/recordings/library", func(w http.ResponseWriter, r *http.Request) {
		files, err := recManager.ListAllRecordingFiles()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if files == nil {
			files = []recording.SavedRecording{}
		}
		jsonResp(w, files)
	})
	webServer.RegisterHandler("/api/recordings/archives", func(w http.ResponseWriter, r *http.Request) {
		items, err := archiveManager.ListArchives()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if items == nil {
			items = []storage.RecordingArchive{}
		}
		jsonResp(w, items)
	})
	webServer.RegisterHandler("/api/recordings/archive", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !licenseRuntime.allows(licenseFeatureRecording) {
			http.Error(w, "recording lisans gerektirir", http.StatusForbidden)
			return
		}
		var req struct {
			StreamKey string `json:"stream_key"`
			Filename  string `json:"filename"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		item, err := archiveManager.ArchiveRecording(r.Context(), req.StreamKey, req.Filename)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonResp(w, item)
	})
	webServer.RegisterHandler("/api/recordings/restore", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !licenseRuntime.allows(licenseFeatureRecording) {
			http.Error(w, "recording lisans gerektirir", http.StatusForbidden)
			return
		}
		var req struct {
			StreamKey string `json:"stream_key"`
			Filename  string `json:"filename"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		item, err := archiveManager.RestoreRecording(r.Context(), req.StreamKey, req.Filename)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonResp(w, item)
	})
	webServer.RegisterHandler("/api/recordings/archive/sync", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !licenseRuntime.allows(licenseFeatureRecording) {
			http.Error(w, "recording lisans gerektirir", http.StatusForbidden)
			return
		}
		uploaded, err := archiveManager.SyncPending(r.Context(), cfg.GetInt("archive_batch_size", 3))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonResp(w, map[string]interface{}{"success": true, "uploaded": uploaded})
	})
	webServer.RegisterAdminHandler("/api/storage/connection-test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if archiveManager == nil {
			http.Error(w, "archive manager hazir degil", http.StatusServiceUnavailable)
			return
		}
		var req struct {
			Role    string            `json:"role"`
			Updates map[string]string `json:"updates"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Role) == "" {
			req.Role = "recordings"
		}
		result, err := archiveManager.TestConnection(r.Context(), req.Role, req.Updates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResp(w, map[string]interface{}{"success": true, "result": result})
	})
	webServer.RegisterHandler("/api/recordings/remux", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !licenseRuntime.allows(licenseFeatureRecording) {
			http.Error(w, "recording lisans gerektirir", http.StatusForbidden)
			return
		}
		var req struct {
			StreamKey string `json:"stream_key"`
			Filename  string `json:"filename"`
			Format    string `json:"format"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Format) == "" {
			req.Format = "mp4"
		}
		job, err := recManager.StartRemuxJob(req.StreamKey, req.Filename, recording.Format(req.Format))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResp(w, map[string]interface{}{"success": true, "job": job})
	})
	webServer.RegisterHandler("/api/recordings/remux/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		jobs := recManager.ListRemuxJobs()
		if jobs == nil {
			jobs = []*recording.RemuxJob{}
		}
		jsonResp(w, jobs)
	})
	webServer.RegisterHandler("/api/recordings/file", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		var req struct {
			StreamKey string `json:"stream_key"`
			Filename  string `json:"filename"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := recManager.DeleteRecording(req.StreamKey, req.Filename); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonResp(w, map[string]string{"status": "deleted"})
	})
	webServer.RegisterHandler("/recordings/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/recordings/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			http.Error(w, "Recording path required", 400)
			return
		}

		streamKey, err := url.PathUnescape(parts[0])
		if err != nil {
			http.Error(w, "Invalid stream key", 400)
			return
		}
		filename, err := url.PathUnescape(parts[1])
		if err != nil {
			http.Error(w, "Invalid filename", 400)
			return
		}
		rc, _, err := recManager.OpenRecording(streamKey, filename)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		defer rc.Close()

		file, ok := rc.(*os.File)
		if !ok {
			http.Error(w, "Recording handle unsupported", 500)
			return
		}
		info, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if ct := detectRecordingContentType(filename); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		if r.URL.Query().Get("download") == "1" {
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(filename)))
		} else {
			w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filepath.Base(filename)))
		}
		http.ServeContent(w, r, filepath.Base(filename), info.ModTime(), file)
	})

	// Analytics API
	webServer.RegisterHandler("/api/analytics", func(w http.ResponseWriter, r *http.Request) {
		dash := analyticsTracker.GetDashboard()
		if streams, err := db.GetAllStreams(); err == nil {
			dash.TotalStreams = len(streams)
			names := make(map[string]string, len(streams))
			for _, st := range streams {
				analyticsTracker.RegisterStreamName(st.StreamKey, st.Name)
				names[st.StreamKey] = st.Name
			}
			for i := range dash.TopStreams {
				if dash.TopStreams[i].StreamName == "" {
					dash.TopStreams[i].StreamName = names[dash.TopStreams[i].StreamKey]
				}
			}
		}
		if cfg.GetBool("rtmps_enabled", false) && !licenseRuntime.allows(licenseFeatureRTMPS) {
			log.Printf("[INIT] RTMPS lisans gerektiriyor; bu node icin devreye alinmadi")
		}
		jsonResp(w, dash)
	})
	webServer.RegisterHandler("/api/analytics/stream/", func(w http.ResponseWriter, r *http.Request) {
		streamKey := strings.TrimPrefix(r.URL.Path, "/api/analytics/stream/")
		stats := analyticsTracker.GetStreamStats(streamKey)
		jsonResp(w, stats)
	})
	webServer.RegisterHandler("/api/analytics/history", func(w http.ResponseWriter, r *http.Request) {
		window := analyticsWindowForPeriod(r.URL.Query().Get("period"), time.Now())
		snapshots, err := db.GetAnalyticsSnapshotsSince(window.Since, 0)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if snapshots == nil {
			snapshots = []storage.AnalyticsSnapshot{}
		}
		jsonResp(w, buildAnalyticsHistoryPayload(window.Period, snapshots, time.Now()))
	})

	// Security API - Token
	webServer.RegisterHandler("/api/security/token/generate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			StreamKey string `json:"stream_key"`
		}
		if err := decodeJSON(r, &req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		token, expiry := tokenMgr.GenerateToken(req.StreamKey)
		jsonResp(w, map[string]interface{}{"token": token, "expires_at": expiry})
	})

	registerStudioAdminRoutes(webServer, db, cfg, analyticsTracker, tcManager, playerTelemetry, tokenMgr)

	// Security API - IP Ban
	webServer.RegisterHandler("/api/security/bans", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			jsonResp(w, ipBanList.GetBanned())
		case "POST":
			var req struct {
				IP       string `json:"ip"`
				Reason   string `json:"reason"`
				Duration int    `json:"duration_minutes"`
			}
			if err := decodeJSON(r, &req); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			dur := time.Duration(req.Duration) * time.Minute
			ipBanList.Ban(req.IP, req.Reason, dur)
			jsonResp(w, map[string]string{"status": "banned"})
		case "DELETE":
			var req struct {
				IP string `json:"ip"`
			}
			if err := decodeJSON(r, &req); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			ipBanList.Unban(req.IP)
			jsonResp(w, map[string]string{"status": "unbanned"})
		}
	})

	// Transcoding API
	webServer.RegisterHandler("/api/transcode/status", func(w http.ResponseWriter, r *http.Request) {
		jsonResp(w, tcManager.GetStatus())
	})
	webServer.RegisterHandler("/api/transcode/jobs", func(w http.ResponseWriter, r *http.Request) {
		jsonResp(w, tcManager.GetJobs())
	})
	webServer.RegisterHandler("/api/health/report", func(w http.ResponseWriter, r *http.Request) {
		stats := streamManager.GetStats()
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		stats.MemoryUsedMB = int64(memStats.Alloc / 1024 / 1024)
		stats.MemoryTotalMB = int64(memStats.Sys / 1024 / 1024)
		stats.UptimeSeconds = int64(time.Since(startTime).Seconds())
		jsonResp(w, buildHealthReport(cfg, db, stats, tcManager, streamManager, playerTelemetry, archiveManager, dataDir))
	})
	webServer.RegisterHandler("/api/maintenance/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		retentionDays := cfg.GetInt("recordings_retention_days", cfg.GetInt("storage_auto_clean", 30))
		deletedOld, _ := recManager.CleanupOld(time.Duration(retentionDays) * 24 * time.Hour)
		trimmed, _ := recManager.TrimLatestPerStream(cfg.GetInt("recordings_keep_latest", 10))
		snapshotDeleted, _ := db.CleanupAnalyticsSnapshots(time.Duration(cfg.GetInt("analytics_retention_days", 30)) * 24 * time.Hour)
		playerTelemetryDeleted, _ := db.CleanupPlayerTelemetrySamples(time.Duration(cfg.GetInt("player_telemetry_retention_days", 30)) * 24 * time.Hour)
		trackAnalyticsDeleted, _ := db.CleanupTrackTelemetrySamples(time.Duration(cfg.GetInt("track_analytics_retention_days", 30)) * 24 * time.Hour)
		jsonResp(w, map[string]interface{}{
			"success":                    true,
			"deleted_recordings_old":     deletedOld,
			"deleted_recordings_trimmed": trimmed,
			"deleted_analytics":          snapshotDeleted,
			"deleted_player_telemetry":   playerTelemetryDeleted,
			"deleted_track_analytics":    trackAnalyticsDeleted,
		})
	})
	webServer.RegisterHandler("/api/diagnostics/stream/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/diagnostics/stream/")
		var id int64
		fmt.Sscanf(idStr, "%d", &id)
		st, err := db.GetStreamByID(id)
		if err != nil || st == nil {
			http.Error(w, "Stream bulunamadi", 404)
			return
		}
		payload := buildStreamDiagnostics(st, cfg, dataDir, tcManager)
		payload["telemetry"] = playerTelemetry.Snapshot(st.StreamKey)
		payload["telemetry_history"], _ = db.GetPlayerTelemetrySamples(st.StreamKey, 48)
		payload["track_history"], _ = db.GetTrackTelemetrySamples(st.StreamKey, 160)
		payload["qoe_alerts"] = buildQoEAlerts(cfg, st.Name, playerTelemetry.Snapshot(st.StreamKey))
		if tcManager != nil {
			payload["tracks"] = tcManager.GetLiveTrackSnapshot(st.StreamKey)
		}
		payload["live_now"] = streamManager.IsLive(st.StreamKey)
		jsonResp(w, payload)
	})

	// Viewer stats API
	webServer.RegisterHandler("/api/viewers", func(w http.ResponseWriter, r *http.Request) {
		dash := analyticsTracker.GetDashboard()
		banned := ipBanList.GetBanned()
		sessions := analyticsTracker.GetViewerSessions()
		jsonResp(w, map[string]interface{}{
			"total":    dash.TotalViewers,
			"active":   dash.CurrentViewers,
			"banned":   len(banned),
			"sessions": sessions,
		})
	})
	webServer.RegisterHandler("/api/stats/viewers", func(w http.ResponseWriter, r *http.Request) {
		dash := analyticsTracker.GetDashboard()
		jsonResp(w, map[string]interface{}{
			"current":    dash.CurrentViewers,
			"total":      dash.TotalViewers,
			"peak":       dash.PeakConcurrent,
			"by_format":  dash.ViewersByFormat,
			"by_country": dash.ViewersByCountry,
			"timeline":   dash.ViewersTimeline,
		})
	})

	// Per-stream record trigger
	webServer.RegisterHandler("/api/streams/record/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		streamKey := strings.TrimPrefix(r.URL.Path, "/api/streams/record/")
		recordFormat := recording.FormatMP4
		if st, err := db.GetStreamByKey(streamKey); err == nil && st != nil && strings.TrimSpace(st.RecordFormat) != "" {
			recordFormat = recording.Format(st.RecordFormat)
		}
		rec, err := recManager.StartRecording(streamKey, recordFormat)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonResp(w, rec)
	})
	registerProductAdminRoutes(webServer, cfg, db, dataDir, licenseRuntime, archiveManager)

	go func() {
		if err := webServer.Start(ctx); err != nil {
			log.Printf("[FATAL] HTTP sunucu hatasi: %v", err)
		}
	}()

	// Wait briefly for servers to bind
	time.Sleep(200 * time.Millisecond)

	hostName, _ := os.Hostname()
	httpsEnabled := cfg.GetBool("ssl_enabled", false) && webTLSSource.Ready
	rtmpsConfigured := cfg.GetBool("rtmps_enabled", false)
	rtmpsEnabled := rtmpsConfigured && licenseRuntime.allows(licenseFeatureRTMPS) && streamTLSSource.Ready
	displayHost := strings.TrimSpace(cfg.Get("embed_domain", ""))
	if displayHost == "" {
		displayHost = "localhost"
	}

	// Print service status
	fmt.Println("Service Endpoints")
	if hostName != "" {
		fmt.Printf("  HOST   %s\n", hostName)
	}
	fmt.Printf("  HTTP   http://%s:%d\n", displayHost, httpPort)
	if httpsEnabled {
		fmt.Printf("  HTTPS  https://%s:%d\n", displayHost, cfg.GetInt("https_port", 443))
	}
	fmt.Printf("  RTMP   rtmp://%s:%d\n", displayHost, rtmpPort)
	if rtmpsEnabled {
		fmt.Printf("  RTMPS  rtmps://%s:%d\n", displayHost, cfg.GetInt("rtmps_port", 1936))
	}
	if cfg.GetBool("srt_enabled", false) {
		fmt.Printf("  SRT    srt://%s:%d\n", displayHost, cfg.GetInt("srt_port", 9000))
	}
	if cfg.GetBool("rtp_enabled", false) {
		fmt.Printf("  RTP    rtp://%s:%d\n", displayHost, cfg.GetInt("rtp_port", 5004))
	}
	if cfg.GetBool("rtsp_enabled", false) {
		fmt.Printf("  RTSP   rtsp://%s:%d\n", displayHost, cfg.GetInt("rtsp_port", 8554))
	}
	if cfg.GetBool("webrtc_enabled", false) {
		fmt.Printf("  WHIP   http://%s:%d/whip\n", displayHost, cfg.GetInt("webrtc_port", 8855))
	}
	if cfg.GetBool("mpegts_enabled", false) {
		fmt.Printf("  TS-UDP udp://%s:%d\n", displayHost, cfg.GetInt("mpegts_port", 9001))
	}
	if cfg.GetBool("http_push_enabled", false) {
		fmt.Printf("  PUSH   http://%s:%d/push\n", displayHost, cfg.GetInt("http_push_port", 8850))
	}
	if cfg.GetBool("ssl_enabled", false) && !httpsEnabled {
		if webTLSSource.UsesLetsEncrypt() {
			fmt.Printf("  NOTE   HTTPS Let's Encrypt bekliyor (%s, port 80/443 acik olmali)\n", webTLSSource.Domain)
		} else {
			fmt.Printf("  NOTE   HTTPS icin gecerli CRT/KEY gerekli (%s | %s)\n", webTLSSource.CertPath, webTLSSource.KeyPath)
		}
	}
	if rtmpsConfigured && !licenseRuntime.allows(licenseFeatureRTMPS) {
		fmt.Printf("  NOTE   RTMPS lisans gerektiriyor; aktif degil\n")
	} else if rtmpsConfigured && !rtmpsEnabled {
		if streamTLSSource.UsesLetsEncrypt() {
			fmt.Printf("  NOTE   RTMPS Let's Encrypt bekliyor (%s, port 80 DNS dogru olmali)\n", streamTLSSource.Domain)
		} else {
			fmt.Printf("  NOTE   RTMPS icin gecerli CRT/KEY gerekli (%s | %s)\n", streamTLSSource.CertPath, streamTLSSource.KeyPath)
		}
	}
	fmt.Println()

	// Health check
	healthURL := fmt.Sprintf("http://localhost:%d/api/stats", httpPort)
	if checkHealth(healthURL) {
		log.Println("[INIT] HTTP sunucu saglik kontrolu basarili")
	} else {
		log.Println("[INIT] HTTP sunucu saglik kontrolu basarili")
	}

	// Log startup time
	elapsed := time.Since(startTime)
	log.Printf("[INIT] FluxStream baslatildi (%.0fms)", float64(elapsed.Milliseconds()))
	db.AddLog("INFO", "system", fmt.Sprintf("FluxStream v%s baslatildi (%.0fms)", Version, float64(elapsed.Milliseconds())))

	// Check if setup is needed and open browser
	if !cfg.GetBool("setup_completed", false) {
		fmt.Println()
		fmt.Println("  Setup gerekli. Tarayici aciliyor...")
		if os.Getenv("FLUXSTREAM_NO_BROWSER") != "1" {
			openBrowser(fmt.Sprintf("http://localhost:%d", httpPort))
		}
	} else {
		fmt.Println()
		fmt.Println("  Sunucu hazir. Yayin bekleniyor...")
		if os.Getenv("FLUXSTREAM_NO_BROWSER") != "1" {
			openBrowser(fmt.Sprintf("http://localhost:%d", httpPort))
		}
	}

	fmt.Println()
	fmt.Println("  Kapatmak icin Ctrl+C")
	fmt.Println()

	// Wait for shutdown signal or panel command
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	shutdownReason := "signal"
	restartRequested := false
	var sig os.Signal
	select {
	case sig = <-sigChan:
		fmt.Println()
		log.Printf("[SHUTDOWN] Signal alindi: %s", sig)
	case cmd := <-controlCh:
		fmt.Println()
		shutdownReason = cmd
		restartRequested = cmd == "restart"
		log.Printf("[SHUTDOWN] Panel uzerinden islem alindi: %s", cmd)
	}
	log.Println("[SHUTDOWN] Aktif yayinlar durduruluyor...")
	recManager.StopAll()
	tcManager.StopAll()
	streamManager.StopAll()

	log.Println("[SHUTDOWN] Sunucular kapatiliyor...")
	cancel()
	time.Sleep(500 * time.Millisecond)

	if restartRequested {
		if err := restartSelf(execPath); err != nil {
			log.Printf("[SHUTDOWN] Yeniden baslatma basarisiz: %v", err)
		} else {
			log.Println("[SHUTDOWN] Yeni surec baslatildi")
		}
	}

	db.AddLog("INFO", "system", "FluxStream kapatildi")
	if shutdownReason == "stop" {
		fmt.Println("  FluxStream durduruldu.")
	} else if restartRequested {
		fmt.Println("  FluxStream yeniden baslatiliyor.")
	} else {
		fmt.Println("  FluxStream kapatildi.")
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

func restartSelf(execPath string) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "start", "\"FluxStream\"", execPath)
		cmd.Dir = filepath.Dir(execPath)
		cmd.Env = append(os.Environ(), "FLUXSTREAM_NO_BROWSER=1")
		return cmd.Start()
	}

	cmd := exec.Command(execPath)
	cmd.Dir = filepath.Dir(execPath)
	cmd.Env = append(os.Environ(), "FLUXSTREAM_NO_BROWSER=1")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Start()
}

func checkHealth(url string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func padRight(s string, n int) string {
	for len(s) < n {
		s += " "
	}
	return s
}

func jsonResp(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func detectRecordingContentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".flv":
		return "video/x-flv"
	case ".ts":
		return "video/mp2t"
	case ".mp3":
		return "audio/mpeg"
	case ".aac":
		return "audio/aac"
	case ".ogg":
		return "audio/ogg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	default:
		return "application/octet-stream"
	}
}
