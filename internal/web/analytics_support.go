package web

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	"github.com/fluxstream/fluxstream/internal/analytics"
)

type countingResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func newCountingResponseWriter(w http.ResponseWriter) *countingResponseWriter {
	return &countingResponseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (w *countingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *countingResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}

func (w *countingResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (s *Server) SetAnalyticsTracker(tracker *analytics.Tracker) {
	s.analytics = tracker
}

func (s *Server) trackPlaybackHeartbeat(w http.ResponseWriter, r *http.Request, streamKey, format string, serve func(http.ResponseWriter)) {
	if s.analytics == nil || streamKey == "" {
		serve(w)
		return
	}
	if isInternalPlaybackRequest(r) {
		serve(w)
		return
	}

	cw := newCountingResponseWriter(w)
	serve(cw)

	if cw.status >= 400 {
		return
	}

	s.analytics.TrackPlayback(
		streamKey,
		ensureViewerID(cw, r),
		format,
		clientIP(r),
		clientCountry(r),
		r.UserAgent(),
		cw.bytes,
	)
}

func isInternalPlaybackRequest(r *http.Request) bool {
	ua := strings.ToLower(strings.TrimSpace(r.UserAgent()))
	if ua == "" {
		return false
	}
	if strings.Contains(ua, "fluxstreaminternal") || strings.Contains(ua, "lavf") || strings.Contains(ua, "ffmpeg") {
		return true
	}
	return false
}

func ensureViewerID(w http.ResponseWriter, r *http.Request) string {
	if id := strings.TrimSpace(r.URL.Query().Get("viewer_id")); id != "" {
		return id
	}
	if cookie, err := r.Cookie("fluxstream_viewer"); err == nil {
		if value := strings.TrimSpace(cookie.Value); value != "" {
			return value
		}
	}

	sum := sha256.Sum256([]byte(clientIP(r) + "|" + strings.TrimSpace(r.UserAgent())))
	id := hex.EncodeToString(sum[:12])
	http.SetCookie(w, &http.Cookie{
		Name:     "fluxstream_viewer",
		Value:    id,
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
	return id
}

func clientIP(r *http.Request) string {
	for _, header := range []string{"CF-Connecting-IP", "X-Real-IP", "X-Forwarded-For"} {
		if value := strings.TrimSpace(r.Header.Get(header)); value != "" {
			if header == "X-Forwarded-For" && strings.Contains(value, ",") {
				return strings.TrimSpace(strings.Split(value, ",")[0])
			}
			return value
		}
	}
	host := strings.TrimSpace(r.RemoteAddr)
	if ip, _, err := net.SplitHostPort(host); err == nil {
		return ip
	}
	return host
}

func clientCountry(r *http.Request) string {
	for _, header := range []string{"CF-IPCountry", "X-Country-Code", "X-Country"} {
		if value := strings.TrimSpace(r.Header.Get(header)); value != "" {
			return value
		}
	}
	return ""
}
