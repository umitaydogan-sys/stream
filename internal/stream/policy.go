package stream

import (
	"encoding/json"
	"strings"
)

// Policy represents optional per-stream delivery and security policy.
type Policy struct {
	Mode                 string   `json:"mode,omitempty"`
	EnableABR            bool     `json:"enable_abr"`
	ProfileSet           string   `json:"profile_set,omitempty"`
	AllowedOutputs       []string `json:"allowed_outputs,omitempty"`
	RequirePlaybackToken bool     `json:"require_playback_token"`
	RequireSignedURL     bool     `json:"require_signed_url"`
	RetentionDays        int      `json:"retention_days,omitempty"`
	Notes                string   `json:"notes,omitempty"`
}

func DefaultPolicy() Policy {
	return Policy{
		Mode:       "balanced",
		EnableABR:  false,
		ProfileSet: "balanced",
	}
}

func ParsePolicyJSON(raw string) Policy {
	policy := DefaultPolicy()
	if strings.TrimSpace(raw) == "" {
		return policy
	}
	_ = json.Unmarshal([]byte(raw), &policy)
	if strings.TrimSpace(policy.ProfileSet) == "" {
		policy.ProfileSet = "balanced"
	}
	if strings.TrimSpace(policy.Mode) == "" {
		policy.Mode = "balanced"
	}
	return policy
}

func EncodePolicyJSON(policy Policy) string {
	if strings.TrimSpace(policy.ProfileSet) == "" {
		policy.ProfileSet = "balanced"
	}
	if strings.TrimSpace(policy.Mode) == "" {
		policy.Mode = "balanced"
	}
	data, err := json.Marshal(policy)
	if err != nil {
		return ""
	}
	return string(data)
}

func (p Policy) AllowsOutput(name string) bool {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return false
	}
	if len(p.AllowedOutputs) == 0 {
		return true
	}
	for _, item := range p.AllowedOutputs {
		if strings.EqualFold(strings.TrimSpace(item), name) {
			return true
		}
	}
	return false
}
