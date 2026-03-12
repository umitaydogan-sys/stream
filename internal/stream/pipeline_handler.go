package stream

import (
	"log"
	"net"
	"sync"

	"github.com/fluxstream/fluxstream/internal/media"
	"github.com/fluxstream/fluxstream/internal/transcode"
)

type autoRecordingManager interface {
	StartManagedRecording(streamKey, format string) (string, error)
	StopRecording(recID string) error
}

// PipelineHandler wraps the stream manager with optional live HLS transcoding.
type PipelineHandler struct {
	manager          *Manager
	transcodeManager *transcode.Manager
	recordingManager autoRecordingManager
	liveHLSEnabled   bool
	liveDASHEnabled  bool
	liveOptions      transcode.LiveOptions
	autoRecordings   map[string]string
	mu               sync.Mutex
}

// NewPipelineHandler creates an ingest handler for the live pipeline.
func NewPipelineHandler(manager *Manager, transcodeManager *transcode.Manager, recordingManager autoRecordingManager, liveHLSEnabled, liveDASHEnabled bool, liveOptions transcode.LiveOptions) *PipelineHandler {
	return &PipelineHandler{
		manager:          manager,
		transcodeManager: transcodeManager,
		recordingManager: recordingManager,
		liveHLSEnabled:   liveHLSEnabled,
		liveDASHEnabled:  liveDASHEnabled,
		liveOptions:      liveOptions,
		autoRecordings:   make(map[string]string),
	}
}

// OnPublish registers the stream and starts live HLS transcoding if enabled.
func (h *PipelineHandler) OnPublish(streamKey string, conn net.Conn) error {
	if err := h.manager.OnPublish(streamKey, conn); err != nil {
		return err
	}

	if h.transcodeManager != nil {
		opts := h.liveOptions
		if active := h.manager.GetActiveStream(streamKey); active != nil && active.DBStream != nil {
			policy := ParsePolicyJSON(active.DBStream.PolicyJSON)
			if policy.EnableABR {
				opts.ABREnabled = true
			}
			if policy.ProfileSet != "" {
				opts.ProfileSet = policy.ProfileSet
			}
			opts.Profiles = transcode.ResolveProfiles(opts.ProfileSet, opts.ProfilesJSON)
		}
		h.transcodeManager.SetStreamLiveOptions(streamKey, opts)
	}

	if h.liveHLSEnabled && h.transcodeManager != nil {
		if _, err := h.transcodeManager.StartLiveHLS(streamKey); err != nil {
			log.Printf("[TC] Canli HLS baslatilamadi (%s): %v", streamKey, err)
		}
	}
	if h.liveDASHEnabled && h.transcodeManager != nil {
		if _, err := h.transcodeManager.StartLiveDASH(streamKey); err != nil {
			log.Printf("[TC] Canli DASH baslatilamadi (%s): %v", streamKey, err)
		}
	}

	if h.recordingManager != nil {
		if active := h.manager.GetActiveStream(streamKey); active != nil && active.DBStream != nil && active.DBStream.RecordEnabled {
			recID, err := h.recordingManager.StartManagedRecording(streamKey, active.DBStream.RecordFormat)
			if err != nil {
				log.Printf("[REC] Otomatik kayit baslatilamadi (%s): %v", streamKey, err)
			} else {
				h.mu.Lock()
				h.autoRecordings[streamKey] = recID
				h.mu.Unlock()
			}
		}
	}

	return nil
}

// OnUnpublish tears down live HLS transcoding and the base stream.
func (h *PipelineHandler) OnUnpublish(streamKey string) {
	if h.recordingManager != nil {
		h.mu.Lock()
		recID := h.autoRecordings[streamKey]
		delete(h.autoRecordings, streamKey)
		h.mu.Unlock()

		if recID != "" {
			if err := h.recordingManager.StopRecording(recID); err != nil {
				log.Printf("[REC] Otomatik kayit durdurulamadi (%s): %v", streamKey, err)
			}
		}
	}

	if h.transcodeManager != nil {
		h.transcodeManager.StopLiveDASH(streamKey)
		h.transcodeManager.StopLiveHLS(streamKey)
	}
	h.manager.OnUnpublish(streamKey)
}

// OnPacket fans packets into the live transcode path and base outputs.
func (h *PipelineHandler) OnPacket(streamKey string, pkt *media.Packet) {
	if h.transcodeManager != nil && h.liveHLSEnabled {
		h.transcodeManager.WriteLivePacket(streamKey, pkt)
	}
	h.manager.OnPacket(streamKey, pkt)
}
