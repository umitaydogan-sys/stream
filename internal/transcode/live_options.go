package transcode

import (
	"encoding/json"
	"strings"
)

// LiveOptions controls the live HLS transcode path.
type LiveOptions struct {
	ABREnabled          bool      `json:"abr_enabled"`
	MasterEnabled       bool      `json:"master_enabled"`
	ProfileSet          string    `json:"profile_set"`
	ProfilesJSON        string    `json:"profiles_json,omitempty"`
	Profiles            []Profile `json:"profiles"`
	SegmentDuration     int       `json:"segment_duration"`
	PlaylistLength      int       `json:"playlist_length"`
	AudioPassthrough    bool      `json:"audio_passthrough"`
	DefaultVideoTrackID uint8     `json:"default_video_track_id,omitempty"`
	DefaultAudioTrackID uint8     `json:"default_audio_track_id,omitempty"`
}

func DefaultLiveOptions() LiveOptions {
	return LiveOptions{
		ABREnabled:      false,
		MasterEnabled:   true,
		ProfileSet:      "balanced",
		Profiles:        DefaultProfiles(),
		SegmentDuration: 2,
		PlaylistLength:  10,
	}
}

func ResolveProfiles(profileSet, rawJSON string) []Profile {
	profileSet = strings.TrimSpace(strings.ToLower(profileSet))
	if profileSet == "" {
		profileSet = "balanced"
	}
	if strings.TrimSpace(rawJSON) == "" {
		return DefaultProfiles()
	}
	sets := map[string][]Profile{}
	if err := json.Unmarshal([]byte(rawJSON), &sets); err != nil {
		return DefaultProfiles()
	}
	if profiles := sets[profileSet]; len(profiles) > 0 {
		return profiles
	}
	if profiles := sets["balanced"]; len(profiles) > 0 {
		return profiles
	}
	return DefaultProfiles()
}
