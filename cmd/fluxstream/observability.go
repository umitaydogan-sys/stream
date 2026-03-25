package main

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/fluxstream/fluxstream/internal/analytics"
	"github.com/fluxstream/fluxstream/internal/config"
	"github.com/fluxstream/fluxstream/internal/storage"
	"github.com/fluxstream/fluxstream/internal/stream"
	"github.com/fluxstream/fluxstream/internal/transcode"
)

type runtimeObservabilityStream struct {
	StreamKey  string
	StreamName string
	Live       bool
	Telemetry  playerTelemetrySnapshot
	Tracks     transcode.LiveTrackSnapshot
}

func trackTelemetrySamplesFromSnapshot(snapshot transcode.LiveTrackSnapshot) []storage.TrackTelemetrySample {
	samples := make([]storage.TrackTelemetrySample, 0, len(snapshot.VideoTracks)+len(snapshot.AudioTracks))
	for _, track := range snapshot.VideoTracks {
		samples = append(samples, storage.TrackTelemetrySample{
			StreamKey:    snapshot.StreamKey,
			TrackID:      track.TrackID,
			Kind:         track.Kind,
			Codec:        track.Codec,
			Width:        track.Width,
			Height:       track.Height,
			SampleRate:   track.SampleRate,
			Channels:     track.Channels,
			Bitrate:      track.Bitrate,
			Packets:      track.Packets,
			Bytes:        track.Bytes,
			IsDefault:    track.IsDefault,
			IsActive:     track.IsActive,
			DisplayLabel: track.DisplayLabel,
		})
	}
	for _, track := range snapshot.AudioTracks {
		samples = append(samples, storage.TrackTelemetrySample{
			StreamKey:    snapshot.StreamKey,
			TrackID:      track.TrackID,
			Kind:         track.Kind,
			Codec:        track.Codec,
			Width:        track.Width,
			Height:       track.Height,
			SampleRate:   track.SampleRate,
			Channels:     track.Channels,
			Bitrate:      track.Bitrate,
			Packets:      track.Packets,
			Bytes:        track.Bytes,
			IsDefault:    track.IsDefault,
			IsActive:     track.IsActive,
			DisplayLabel: track.DisplayLabel,
		})
	}
	return samples
}

func collectRuntimeObservability(streamMgr *stream.Manager, tcManager *transcode.Manager, playerTelemetry *playerTelemetryCollector) []runtimeObservabilityStream {
	if streamMgr == nil {
		return []runtimeObservabilityStream{}
	}
	activeStreams := streamMgr.GetActiveStreams()
	items := make([]runtimeObservabilityStream, 0, len(activeStreams))
	for _, active := range activeStreams {
		if active == nil {
			continue
		}
		item := runtimeObservabilityStream{
			StreamKey: active.Key,
			Live:      true,
		}
		if active.DBStream != nil {
			item.StreamName = active.DBStream.Name
		}
		if playerTelemetry != nil {
			item.Telemetry = playerTelemetry.Snapshot(active.Key)
		}
		if tcManager != nil {
			item.Tracks = tcManager.GetLiveTrackSnapshot(active.Key)
		}
		items = append(items, item)
	}
	return items
}

func buildQoEAlerts(cfg *config.Manager, streamName string, snapshot playerTelemetrySnapshot) []systemAlert {
	if cfg == nil || snapshot.ActiveSessions <= 0 {
		return nil
	}
	label := strings.TrimSpace(streamName)
	if label == "" {
		label = strings.TrimSpace(snapshot.StreamKey)
	}
	maxStalls := 0
	for _, session := range snapshot.Sessions {
		if session.StallCount > maxStalls {
			maxStalls = session.StallCount
		}
	}

	stallThreshold := cfg.GetInt("alerts_qoe_stalls_threshold", 4)
	bufferWarn := cfg.GetFloat("alerts_qoe_buffer_warn_seconds", float64(cfg.GetInt("alerts_qoe_buffer_seconds", 1)))
	bufferCritical := cfg.GetFloat("alerts_qoe_buffer_critical_seconds", math.Max(0.6, bufferWarn*0.5))
	waitingBase := cfg.GetInt("alerts_qoe_waiting_sessions", 2)
	waitingRatio := cfg.GetInt("alerts_qoe_waiting_ratio_percent", 35)
	offlineBase := cfg.GetInt("alerts_qoe_offline_sessions", 1)
	offlineRatio := cfg.GetInt("alerts_qoe_offline_ratio_percent", 20)
	transitionRatioThreshold := cfg.GetInt("alerts_qoe_transition_ratio_threshold", 4)

	waitingThreshold := maxInt(waitingBase, int(math.Ceil(float64(snapshot.ActiveSessions)*float64(waitingRatio)/100.0)))
	offlineThreshold := maxInt(offlineBase, int(math.Ceil(float64(snapshot.ActiveSessions)*float64(offlineRatio)/100.0)))

	alerts := make([]systemAlert, 0, 5)
	if stallThreshold > 0 && maxStalls >= stallThreshold {
		level := "warning"
		if maxStalls >= stallThreshold*2 {
			level = "critical"
		}
		alerts = append(alerts, systemAlert{
			Level:       level,
			Code:        "qoe_stalls_high",
			Title:       "QoE stall sayisi yuksek",
			Description: fmt.Sprintf("%s yayininda tek bir oturumda stall sayisi %d degerine ulasti.", label, maxStalls),
			Action:      "ABR profilini hafifletin, HLS/DASH cikislarini ve istemci tarafindaki aktif kaliteyi kontrol edin.",
		})
	}
	if bufferWarn > 0 && snapshot.AverageBufferSeconds > 0 && snapshot.AverageBufferSeconds <= bufferWarn {
		level := "warning"
		if snapshot.AverageBufferSeconds <= bufferCritical {
			level = "critical"
		}
		alerts = append(alerts, systemAlert{
			Level:       level,
			Code:        "qoe_buffer_low",
			Title:       "QoE buffer seviyesi dusuk",
			Description: fmt.Sprintf("%s yayininda ortalama buffer %.1fs seviyesine dustu.", label, snapshot.AverageBufferSeconds),
			Action:      "Segment suresi, ABR merdiveni ve kaynak CPU kullanimi kontrol edilmeli.",
		})
	}
	if waitingThreshold > 0 && snapshot.WaitingSessions >= waitingThreshold {
		alerts = append(alerts, systemAlert{
			Level:       "warning",
			Code:        "qoe_waiting_sessions",
			Title:       "Bekleyen player oturumlari artti",
			Description: fmt.Sprintf("%s yayininda %d player bekleme durumunda gorunuyor. Esik %d oturum olarak hesaplandi.", label, snapshot.WaitingSessions, waitingThreshold),
			Action:      "Player QoE paneli ve stream teslimat katmanlari kontrol edilmelidir.",
		})
	}
	if offlineThreshold > 0 && snapshot.OfflineSessions >= offlineThreshold {
		alerts = append(alerts, systemAlert{
			Level:       "warning",
			Code:        "qoe_offline_sessions",
			Title:       "Offline gorunen player oturumlari var",
			Description: fmt.Sprintf("%s yayininda %d player offline ekranina dustu. Esik %d oturum olarak hesaplandi.", label, snapshot.OfflineSessions, offlineThreshold),
			Action:      "Manifest surekliligi ve istemci fallback davranisi incelenmelidir.",
		})
	}
	if transitionRatioThreshold > 0 && snapshot.ActiveSessions > 0 && snapshot.TotalQualityTransitions >= int64(snapshot.ActiveSessions*transitionRatioThreshold) && snapshot.AverageBufferSeconds > 0 && snapshot.AverageBufferSeconds <= (bufferWarn+0.8) {
		alerts = append(alerts, systemAlert{
			Level:       "warning",
			Code:        "qoe_quality_flapping",
			Title:       "Kalite gecisleri siklasti",
			Description: fmt.Sprintf("%s yayininda kalite gecisi sayisi %d oldu. Bu deger aktif oturum sayisina gore yuksek gorunuyor.", label, snapshot.TotalQualityTransitions),
			Action:      "Dusuk bant profiline gecin, ABR merdivenini hafifletin ve istemcinin hizli yukselmesini sinirlayin.",
		})
	}
	return alerts
}

func buildPrometheusMetrics(version string, tracker *analytics.Tracker, streamMgr *stream.Manager, tcManager *transcode.Manager, playerTelemetry *playerTelemetryCollector) string {
	items := collectRuntimeObservability(streamMgr, tcManager, playerTelemetry)
	dashboard := analytics.Dashboard{}
	if tracker != nil {
		dashboard = tracker.GetDashboard()
	}
	var b strings.Builder
	now := time.Now().Unix()
	b.WriteString("# HELP fluxstream_info FluxStream surum bilgisi.\n")
	b.WriteString("# TYPE fluxstream_info gauge\n")
	b.WriteString(fmt.Sprintf("fluxstream_info{version=%q} 1\n", version))
	b.WriteString("# HELP fluxstream_streams_active Aktif canli yayin sayisi.\n")
	b.WriteString("# TYPE fluxstream_streams_active gauge\n")
	b.WriteString(fmt.Sprintf("fluxstream_streams_active %d\n", len(items)))
	b.WriteString("# HELP fluxstream_viewers_current Anlik izleyici sayisi.\n")
	b.WriteString("# TYPE fluxstream_viewers_current gauge\n")
	b.WriteString(fmt.Sprintf("fluxstream_viewers_current %d\n", dashboard.CurrentViewers))
	b.WriteString("# HELP fluxstream_metrics_timestamp_seconds Metrics uretim zamani.\n")
	b.WriteString("# TYPE fluxstream_metrics_timestamp_seconds gauge\n")
	b.WriteString(fmt.Sprintf("fluxstream_metrics_timestamp_seconds %d\n", now))
	b.WriteString("# HELP fluxstream_stream_live Yayin canli mi.\n")
	b.WriteString("# TYPE fluxstream_stream_live gauge\n")
	b.WriteString("# HELP fluxstream_player_active_sessions Aktif player oturumlari.\n")
	b.WriteString("# TYPE fluxstream_player_active_sessions gauge\n")
	b.WriteString("# HELP fluxstream_player_waiting_sessions Bekleyen player oturumlari.\n")
	b.WriteString("# TYPE fluxstream_player_waiting_sessions gauge\n")
	b.WriteString("# HELP fluxstream_player_offline_sessions Offline player oturumlari.\n")
	b.WriteString("# TYPE fluxstream_player_offline_sessions gauge\n")
	b.WriteString("# HELP fluxstream_player_average_buffer_seconds Ortalama buffer suresi.\n")
	b.WriteString("# TYPE fluxstream_player_average_buffer_seconds gauge\n")
	b.WriteString("# HELP fluxstream_player_total_stalls Toplam stall sayisi.\n")
	b.WriteString("# TYPE fluxstream_player_total_stalls gauge\n")
	b.WriteString("# HELP fluxstream_player_total_quality_transitions Toplam kalite gecisi sayisi.\n")
	b.WriteString("# TYPE fluxstream_player_total_quality_transitions gauge\n")
	b.WriteString("# HELP fluxstream_player_total_audio_switches Toplam audio track degisimi sayisi.\n")
	b.WriteString("# TYPE fluxstream_player_total_audio_switches gauge\n")
	b.WriteString("# HELP fluxstream_track_bitrate_bps Canli track bitrate degeri.\n")
	b.WriteString("# TYPE fluxstream_track_bitrate_bps gauge\n")
	b.WriteString("# HELP fluxstream_track_packets_total Track paket sayisi.\n")
	b.WriteString("# TYPE fluxstream_track_packets_total gauge\n")
	b.WriteString("# HELP fluxstream_track_bytes_total Track byte sayisi.\n")
	b.WriteString("# TYPE fluxstream_track_bytes_total gauge\n")
	for _, item := range items {
		labels := fmt.Sprintf("stream_key=%q,stream_name=%q", item.StreamKey, item.StreamName)
		b.WriteString(fmt.Sprintf("fluxstream_stream_live{%s} 1\n", labels))
		b.WriteString(fmt.Sprintf("fluxstream_player_active_sessions{%s} %d\n", labels, item.Telemetry.ActiveSessions))
		b.WriteString(fmt.Sprintf("fluxstream_player_waiting_sessions{%s} %d\n", labels, item.Telemetry.WaitingSessions))
		b.WriteString(fmt.Sprintf("fluxstream_player_offline_sessions{%s} %d\n", labels, item.Telemetry.OfflineSessions))
		b.WriteString(fmt.Sprintf("fluxstream_player_average_buffer_seconds{%s} %.3f\n", labels, item.Telemetry.AverageBufferSeconds))
		b.WriteString(fmt.Sprintf("fluxstream_player_total_stalls{%s} %d\n", labels, item.Telemetry.TotalStalls))
		b.WriteString(fmt.Sprintf("fluxstream_player_total_quality_transitions{%s} %d\n", labels, item.Telemetry.TotalQualityTransitions))
		b.WriteString(fmt.Sprintf("fluxstream_player_total_audio_switches{%s} %d\n", labels, item.Telemetry.TotalAudioSwitches))
		for _, track := range item.Tracks.VideoTracks {
			trackLabels := fmt.Sprintf("%s,kind=%q,track_id=%q,display_label=%q", labels, track.Kind, fmt.Sprintf("%d", track.TrackID), track.DisplayLabel)
			b.WriteString(fmt.Sprintf("fluxstream_track_bitrate_bps{%s} %d\n", trackLabels, track.Bitrate))
			b.WriteString(fmt.Sprintf("fluxstream_track_packets_total{%s} %d\n", trackLabels, track.Packets))
			b.WriteString(fmt.Sprintf("fluxstream_track_bytes_total{%s} %d\n", trackLabels, track.Bytes))
		}
		for _, track := range item.Tracks.AudioTracks {
			trackLabels := fmt.Sprintf("%s,kind=%q,track_id=%q,display_label=%q", labels, track.Kind, fmt.Sprintf("%d", track.TrackID), track.DisplayLabel)
			b.WriteString(fmt.Sprintf("fluxstream_track_bitrate_bps{%s} %d\n", trackLabels, track.Bitrate))
			b.WriteString(fmt.Sprintf("fluxstream_track_packets_total{%s} %d\n", trackLabels, track.Packets))
			b.WriteString(fmt.Sprintf("fluxstream_track_bytes_total{%s} %d\n", trackLabels, track.Bytes))
		}
	}
	return b.String()
}

func buildOpenTelemetryPayload(version string, tracker *analytics.Tracker, streamMgr *stream.Manager, tcManager *transcode.Manager, playerTelemetry *playerTelemetryCollector) map[string]interface{} {
	items := collectRuntimeObservability(streamMgr, tcManager, playerTelemetry)
	dashboard := analytics.Dashboard{}
	if tracker != nil {
		dashboard = tracker.GetDashboard()
	}
	nowUnixNano := time.Now().UnixNano()
	metrics := make([]map[string]interface{}, 0, 8)
	metrics = append(metrics,
		buildOTELGaugeMetric("fluxstream.streams.active", "1", []map[string]interface{}{
			otelNumberPoint(nil, nowUnixNano, float64(len(items))),
		}),
		buildOTELGaugeMetric("fluxstream.viewers.current", "1", []map[string]interface{}{
			otelNumberPoint(nil, nowUnixNano, float64(dashboard.CurrentViewers)),
		}),
	)

	livePoints := make([]map[string]interface{}, 0, len(items))
	activeSessionPoints := make([]map[string]interface{}, 0, len(items))
	waitingPoints := make([]map[string]interface{}, 0, len(items))
	bufferPoints := make([]map[string]interface{}, 0, len(items))
	stallPoints := make([]map[string]interface{}, 0, len(items))
	qualityTransitionPoints := make([]map[string]interface{}, 0, len(items))
	audioSwitchPoints := make([]map[string]interface{}, 0, len(items))
	trackBitratePoints := make([]map[string]interface{}, 0, len(items)*4)
	for _, item := range items {
		attrs := otelAttributes(
			"stream.key", item.StreamKey,
			"stream.name", item.StreamName,
		)
		livePoints = append(livePoints, otelNumberPoint(attrs, nowUnixNano, 1))
		activeSessionPoints = append(activeSessionPoints, otelNumberPoint(attrs, nowUnixNano, float64(item.Telemetry.ActiveSessions)))
		waitingPoints = append(waitingPoints, otelNumberPoint(attrs, nowUnixNano, float64(item.Telemetry.WaitingSessions)))
		bufferPoints = append(bufferPoints, otelNumberPoint(attrs, nowUnixNano, item.Telemetry.AverageBufferSeconds))
		stallPoints = append(stallPoints, otelNumberPoint(attrs, nowUnixNano, float64(item.Telemetry.TotalStalls)))
		qualityTransitionPoints = append(qualityTransitionPoints, otelNumberPoint(attrs, nowUnixNano, float64(item.Telemetry.TotalQualityTransitions)))
		audioSwitchPoints = append(audioSwitchPoints, otelNumberPoint(attrs, nowUnixNano, float64(item.Telemetry.TotalAudioSwitches)))
		for _, track := range item.Tracks.VideoTracks {
			trackBitratePoints = append(trackBitratePoints, otelNumberPoint(otelAttributes(
				"stream.key", item.StreamKey,
				"stream.name", item.StreamName,
				"track.id", fmt.Sprintf("%d", track.TrackID),
				"track.kind", track.Kind,
				"track.label", track.DisplayLabel,
			), nowUnixNano, float64(track.Bitrate)))
		}
		for _, track := range item.Tracks.AudioTracks {
			trackBitratePoints = append(trackBitratePoints, otelNumberPoint(otelAttributes(
				"stream.key", item.StreamKey,
				"stream.name", item.StreamName,
				"track.id", fmt.Sprintf("%d", track.TrackID),
				"track.kind", track.Kind,
				"track.label", track.DisplayLabel,
			), nowUnixNano, float64(track.Bitrate)))
		}
	}
	if len(livePoints) > 0 {
		metrics = append(metrics,
			buildOTELGaugeMetric("fluxstream.stream.live", "1", livePoints),
			buildOTELGaugeMetric("fluxstream.player.active_sessions", "1", activeSessionPoints),
			buildOTELGaugeMetric("fluxstream.player.waiting_sessions", "1", waitingPoints),
			buildOTELGaugeMetric("fluxstream.player.average_buffer_seconds", "s", bufferPoints),
			buildOTELGaugeMetric("fluxstream.player.total_stalls", "1", stallPoints),
			buildOTELGaugeMetric("fluxstream.player.total_quality_transitions", "1", qualityTransitionPoints),
			buildOTELGaugeMetric("fluxstream.player.total_audio_switches", "1", audioSwitchPoints),
		)
	}
	if len(trackBitratePoints) > 0 {
		metrics = append(metrics, buildOTELGaugeMetric("fluxstream.track.bitrate", "bit/s", trackBitratePoints))
	}
	return map[string]interface{}{
		"resourceMetrics": []map[string]interface{}{
			{
				"resource": map[string]interface{}{
					"attributes": otelAttributes(
						"service.name", AppName,
						"service.version", version,
						"telemetry.sdk.name", "fluxstream",
						"telemetry.sdk.language", "go",
					),
				},
				"scopeMetrics": []map[string]interface{}{
					{
						"scope": map[string]interface{}{
							"name":    "fluxstream.observability",
							"version": version,
						},
						"metrics": metrics,
					},
				},
			},
		},
	}
}

func buildOTELGaugeMetric(name, unit string, points []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name": name,
		"unit": unit,
		"gauge": map[string]interface{}{
			"data_points": points,
		},
	}
}

func otelAttributes(keyValues ...string) []map[string]interface{} {
	attrs := make([]map[string]interface{}, 0, len(keyValues)/2)
	for i := 0; i+1 < len(keyValues); i += 2 {
		if strings.TrimSpace(keyValues[i+1]) == "" {
			continue
		}
		attrs = append(attrs, map[string]interface{}{
			"key": keyValues[i],
			"value": map[string]interface{}{
				"stringValue": keyValues[i+1],
			},
		})
	}
	return attrs
}

func otelNumberPoint(attrs []map[string]interface{}, ts int64, value float64) map[string]interface{} {
	point := map[string]interface{}{
		"timeUnixNano": ts,
		"asDouble":     value,
	}
	if len(attrs) > 0 {
		point["attributes"] = attrs
	}
	return point
}

func prometheusEscape(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", "\"", "\\\"", "\n", " ")
	return replacer.Replace(value)
}

func marshalJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
