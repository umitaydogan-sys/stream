package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	"github.com/fluxstream/fluxstream/internal/analytics"
)

type streamingTrackingWriter struct {
	http.ResponseWriter
	status  int
	bytes   int64
	started bool
	onStart func()
}

func newStreamingTrackingWriter(w http.ResponseWriter, onStart func()) *streamingTrackingWriter {
	return &streamingTrackingWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
		onStart:        onStart,
	}
}

func (w *streamingTrackingWriter) ensureStarted(status int) {
	if w.started || status >= 400 {
		return
	}
	w.started = true
	if w.onStart != nil {
		w.onStart()
	}
}

func (w *streamingTrackingWriter) WriteHeader(status int) {
	w.status = status
	w.ensureStarted(status)
	w.ResponseWriter.WriteHeader(status)
}

func (w *streamingTrackingWriter) Write(b []byte) (int, error) {
	w.ensureStarted(http.StatusOK)
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}

func (w *streamingTrackingWriter) Flush() {
	w.ensureStarted(w.status)
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func wrapStreamingPlaybackHandler(tracker *analytics.Tracker, authorizer playbackAuthorizer, format string, extractKey func(*http.Request) string, next http.HandlerFunc) http.HandlerFunc {
	if tracker == nil {
		return func(w http.ResponseWriter, r *http.Request) {
			if authorizer != nil {
				streamKey := strings.TrimSpace(extractKey(r))
				if ok, status, message := authorizer(r, streamKey, format); !ok {
					if status <= 0 {
						status = http.StatusForbidden
					}
					w.WriteHeader(status)
					_, _ = w.Write([]byte(message))
					return
				}
			}
			next(w, r)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		streamKey := strings.TrimSpace(extractKey(r))
		if authorizer != nil {
			if ok, status, message := authorizer(r, streamKey, format); !ok {
				if status <= 0 {
					status = http.StatusForbidden
				}
				w.WriteHeader(status)
				_, _ = w.Write([]byte(message))
				return
			}
		}
		if isInternalPlaybackRequest(r) {
			next(w, r)
			return
		}

		if streamKey == "" {
			next(w, r)
			return
		}

		viewerID := requestViewerID(w, r)
		writer := newStreamingTrackingWriter(w, func() {
			tracker.TrackPlayback(streamKey, viewerID, format, requestClientIP(r), requestCountry(r), r.UserAgent(), 0)
		})
		defer func() {
			if writer.started {
				tracker.EndPlayback(streamKey, viewerID, writer.bytes)
			}
		}()

		next(writer, r)
	}
}

func requestViewerID(w http.ResponseWriter, r *http.Request) string {
	if id := strings.TrimSpace(r.URL.Query().Get("viewer_id")); id != "" {
		return id
	}
	if cookie, err := r.Cookie("fluxstream_viewer"); err == nil {
		if value := strings.TrimSpace(cookie.Value); value != "" {
			return value
		}
	}

	sum := sha256.Sum256([]byte(requestClientIP(r) + "|" + strings.TrimSpace(r.UserAgent())))
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

func requestClientIP(r *http.Request) string {
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

func requestCountry(r *http.Request) string {
	for _, header := range []string{"CF-IPCountry", "X-Country-Code", "X-Country"} {
		if value := strings.TrimSpace(r.Header.Get(header)); value != "" {
			return value
		}
	}
	return ""
}

func isInternalPlaybackRequest(r *http.Request) bool {
	ua := strings.ToLower(strings.TrimSpace(r.UserAgent()))
	return strings.Contains(ua, "fluxstreaminternal") || strings.Contains(ua, "lavf") || strings.Contains(ua, "ffmpeg")
}

func flvPlaybackKey(r *http.Request) string {
	return firstPathSegment(strings.TrimPrefix(r.URL.Path, "/flv/"))
}

func mp4PlaybackKey(r *http.Request) string {
	return mediaKeyFromPath(strings.TrimPrefix(r.URL.Path, "/mp4/"), ".mp4")
}

func webmPlaybackKey(r *http.Request) string {
	return mediaKeyFromPath(strings.TrimPrefix(r.URL.Path, "/webm/"), ".webm")
}

func audioPlaybackKey(prefix string) func(*http.Request) string {
	return func(r *http.Request) string {
		key := firstPathSegment(strings.TrimPrefix(r.URL.Path, prefix))
		for _, ext := range []string{".mp3", ".aac", ".ogg", ".wav", ".flac"} {
			key = strings.TrimSuffix(key, ext)
		}
		return key
	}
}

func icecastPlaybackKey(r *http.Request) string {
	return firstPathSegment(strings.TrimPrefix(r.URL.Path, "/icecast/"))
}

func firstPathSegment(value string) string {
	value = strings.Trim(value, "/")
	if value == "" {
		return ""
	}
	return strings.Split(value, "/")[0]
}

func mediaKeyFromPath(value, ext string) string {
	key := firstPathSegment(value)
	return strings.TrimSuffix(key, ext)
}
