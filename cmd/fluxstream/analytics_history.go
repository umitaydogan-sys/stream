package main

import (
	"time"

	"github.com/fluxstream/fluxstream/internal/analytics"
	"github.com/fluxstream/fluxstream/internal/storage"
)

type analyticsHistoryWindow struct {
	Period string
	Since  time.Time
	Bucket time.Duration
	Label  string
}

type analyticsHistoryPayload struct {
	Period    string                    `json:"period"`
	Label     string                    `json:"label"`
	Viewers   []analytics.TimelinePoint `json:"viewers"`
	Bandwidth []analytics.TimelinePoint `json:"bandwidth"`
	Points    int                       `json:"points"`
}

func analyticsWindowForPeriod(period string, now time.Time) analyticsHistoryWindow {
	switch period {
	case "7d":
		return analyticsHistoryWindow{
			Period: "7d",
			Since:  now.Add(-7 * 24 * time.Hour),
			Bucket: 6 * time.Hour,
			Label:  "Son 7 gun",
		}
	case "30d":
		return analyticsHistoryWindow{
			Period: "30d",
			Since:  now.Add(-30 * 24 * time.Hour),
			Bucket: 24 * time.Hour,
			Label:  "Son 30 gun",
		}
	default:
		return analyticsHistoryWindow{
			Period: "24h",
			Since:  now.Add(-24 * time.Hour),
			Bucket: 1 * time.Hour,
			Label:  "Son 24 saat",
		}
	}
}

func buildAnalyticsHistoryPayload(period string, snapshots []storage.AnalyticsSnapshot, now time.Time) analyticsHistoryPayload {
	window := analyticsWindowForPeriod(period, now)
	type bucketValue struct {
		viewers   int64
		bandwidth int64
		hasValue  bool
	}

	buckets := make(map[time.Time]*bucketValue)
	for _, item := range snapshots {
		if item.Timestamp.Before(window.Since) {
			continue
		}
		ts := item.Timestamp.UTC().Truncate(window.Bucket)
		entry := buckets[ts]
		if entry == nil {
			entry = &bucketValue{}
			buckets[ts] = entry
		}
		if !entry.hasValue || int64(item.CurrentViewers) > entry.viewers {
			entry.viewers = int64(item.CurrentViewers)
		}
		if !entry.hasValue || item.TotalBandwidth > entry.bandwidth {
			entry.bandwidth = item.TotalBandwidth
		}
		entry.hasValue = true
	}

	payload := analyticsHistoryPayload{
		Period:    window.Period,
		Label:     window.Label,
		Viewers:   make([]analytics.TimelinePoint, 0, 32),
		Bandwidth: make([]analytics.TimelinePoint, 0, 32),
	}

	for ts := window.Since.UTC().Truncate(window.Bucket); !ts.After(now.UTC()); ts = ts.Add(window.Bucket) {
		entry := buckets[ts]
		if entry == nil {
			payload.Viewers = append(payload.Viewers, analytics.TimelinePoint{Timestamp: ts, Value: 0})
			payload.Bandwidth = append(payload.Bandwidth, analytics.TimelinePoint{Timestamp: ts, Value: 0})
			continue
		}
		payload.Viewers = append(payload.Viewers, analytics.TimelinePoint{Timestamp: ts, Value: entry.viewers})
		payload.Bandwidth = append(payload.Bandwidth, analytics.TimelinePoint{Timestamp: ts, Value: entry.bandwidth})
	}

	payload.Points = len(payload.Viewers)
	return payload
}
