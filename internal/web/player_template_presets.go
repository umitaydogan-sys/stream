package web

import (
	"strings"

	"github.com/fluxstream/fluxstream/internal/storage"
)

func (s *Server) ensureDefaultPlayerTemplates() error {
	templates, err := s.db.GetPlayerTemplates()
	if err != nil {
		return err
	}
	existing := make(map[string]struct{}, len(templates))
	for _, tpl := range templates {
		existing[strings.ToLower(strings.TrimSpace(tpl.Name))] = struct{}{}
	}
	for _, tpl := range defaultPlayerTemplates() {
		key := strings.ToLower(strings.TrimSpace(tpl.Name))
		if _, ok := existing[key]; ok {
			continue
		}
		if _, err := s.db.CreatePlayerTemplate(&tpl); err != nil {
			return err
		}
	}
	return nil
}

func defaultPlayerTemplates() []storage.PlayerTemplate {
	return []storage.PlayerTemplate{
		{
			Name:          "Broadcast Dark",
			Theme:         "dark",
			BackgroundCSS: "background:linear-gradient(180deg,#030712 0%,#0f172a 100%);",
			ControlBarCSS: "background:rgba(3,7,18,.78);backdrop-filter:blur(8px);",
			PlayButtonCSS: "color:#ffffff;filter:drop-shadow(0 10px 24px rgba(37,99,235,.35));",
			WatermarkText: "LIVE",
			ShowTitle:     true,
			ShowLiveBadge: true,
		},
		{
			Name:          "Clean Light",
			Theme:         "light",
			BackgroundCSS: "background:linear-gradient(180deg,#f8fbff 0%,#e9f0f8 100%);",
			ControlBarCSS: "background:rgba(255,255,255,.92);backdrop-filter:blur(10px);color:#1f2937;",
			PlayButtonCSS: "color:#2563eb;",
			WatermarkText: "FluxStream",
			ShowTitle:     true,
			ShowLiveBadge: true,
		},
		{
			Name:          "Minimal Slate",
			Theme:         "minimal",
			BackgroundCSS: "background:#0f172a;",
			ControlBarCSS: "background:rgba(15,23,42,.84);border-top:1px solid rgba(148,163,184,.2);",
			PlayButtonCSS: "color:#e2e8f0;",
			ShowTitle:     false,
			ShowLiveBadge: true,
		},
		{
			Name:          "Glass Studio",
			Theme:         "custom",
			BackgroundCSS: "background:radial-gradient(circle at top left,#1d4ed8 0%,#020617 58%);",
			ControlBarCSS: "background:rgba(255,255,255,.12);backdrop-filter:blur(14px);border:1px solid rgba(255,255,255,.14);",
			PlayButtonCSS: "color:#93c5fd;",
			WatermarkText: "STUDIO",
			ShowTitle:     true,
			ShowLiveBadge: true,
			CustomCSS:     ".player-shell{border-radius:18px;overflow:hidden;box-shadow:0 30px 80px rgba(2,6,23,.35);}",
		},
		{
			Name:          "Radio Compact",
			Theme:         "custom",
			BackgroundCSS: "background:linear-gradient(135deg,#111827 0%,#1f2937 100%);",
			ControlBarCSS: "background:rgba(17,24,39,.86);",
			PlayButtonCSS: "color:#34d399;",
			WatermarkText: "RADIO",
			ShowTitle:     true,
			ShowLiveBadge: false,
			CustomCSS:     ".player-shell.audio-only{max-width:480px;margin:0 auto;border-radius:16px;}",
		},
		{
			Name:          "Neon Arena",
			Theme:         "custom",
			BackgroundCSS: "background:linear-gradient(135deg,#020617 0%,#111827 40%,#1d4ed8 100%);",
			ControlBarCSS: "background:rgba(2,6,23,.8);box-shadow:0 -1px 0 rgba(96,165,250,.25) inset;",
			PlayButtonCSS: "color:#60a5fa;filter:drop-shadow(0 0 16px rgba(96,165,250,.45));",
			WatermarkText: "ARENA",
			ShowTitle:     true,
			ShowLiveBadge: true,
			CustomCSS:     ".player-shell{outline:1px solid rgba(96,165,250,.24);}",
		},
	}
}
