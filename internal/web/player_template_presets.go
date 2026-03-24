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
		{
			Name:          "Aurora Radio",
			Theme:         "custom",
			BackgroundCSS: "background:linear-gradient(135deg,#041c32 0%,#064663 52%,#0c7b93 100%);",
			ControlBarCSS: "background:rgba(4,28,50,.82);",
			PlayButtonCSS: "background:rgba(34,211,238,.18);color:#cffafe;",
			WatermarkText: "RADIO",
			ShowTitle:     true,
			ShowLiveBadge: false,
			CustomCSS:     ".audio-shell.player-shell.audio-only{max-width:560px;border:1px solid rgba(103,232,249,.25);}.audio-note{color:#dbeafe;}",
		},
		{
			Name:          "Podcast Copper",
			Theme:         "custom",
			BackgroundCSS: "background:linear-gradient(145deg,#1c1917 0%,#292524 48%,#7c2d12 100%);",
			ControlBarCSS: "background:rgba(28,25,23,.84);",
			PlayButtonCSS: "background:rgba(249,115,22,.18);color:#fed7aa;",
			WatermarkText: "PODCAST",
			ShowTitle:     true,
			ShowLiveBadge: false,
			CustomCSS:     ".audio-shell.player-shell.audio-only{box-shadow:0 28px 70px rgba(120,53,15,.22);}.audio-badge{background:rgba(251,146,60,.18);color:#ffedd5;}",
		},
		{
			Name:          "Talkback Mono",
			Theme:         "minimal",
			BackgroundCSS: "background:linear-gradient(180deg,#111111 0%,#27272a 100%);",
			ControlBarCSS: "background:rgba(17,17,17,.88);",
			PlayButtonCSS: "background:rgba(244,244,245,.12);color:#fafafa;",
			WatermarkText: "VOICE",
			ShowTitle:     true,
			ShowLiveBadge: false,
			CustomCSS:     ".audio-shell.player-shell.audio-only{max-width:520px;border-radius:20px;}.audio-title strong{text-transform:uppercase;letter-spacing:.08em;font-size:14px;}",
		},
		{
			Name:          "Signal Amber",
			Theme:         "custom",
			BackgroundCSS: "background:linear-gradient(135deg,#111827 0%,#451a03 50%,#f59e0b 100%);",
			ControlBarCSS: "background:rgba(17,24,39,.82);",
			PlayButtonCSS: "background:rgba(245,158,11,.22);color:#fef3c7;",
			WatermarkText: "SIGNAL",
			ShowTitle:     true,
			ShowLiveBadge: true,
			CustomCSS:     ".player-shell{box-shadow:0 26px 80px rgba(146,64,14,.22);}.resume-button{letter-spacing:.05em;}",
		},
		{
			Name:          "Ocean Voice",
			Theme:         "custom",
			BackgroundCSS: "background:linear-gradient(135deg,#082f49 0%,#0f766e 55%,#d1fae5 100%);",
			ControlBarCSS: "background:rgba(8,47,73,.78);",
			PlayButtonCSS: "background:rgba(16,185,129,.22);color:#ecfdf5;",
			WatermarkText: "AUDIO",
			ShowTitle:     true,
			ShowLiveBadge: false,
			CustomCSS:     ".audio-shell.player-shell.audio-only{border:1px solid rgba(153,246,228,.26);}.audio-watermark{background:rgba(6,78,59,.42);}",
		},
		{
			Name:          "Newsroom Slate",
			Theme:         "custom",
			BackgroundCSS: "background:linear-gradient(180deg,#0b1120 0%,#1f2937 100%);",
			ControlBarCSS: "background:rgba(11,17,32,.9);",
			PlayButtonCSS: "background:rgba(59,130,246,.2);color:#dbeafe;",
			WatermarkText: "NEWS",
			ShowTitle:     true,
			ShowLiveBadge: true,
			CustomCSS:     ".player-shell{border-top:4px solid rgba(59,130,246,.72);}.topbar-title small{text-transform:uppercase;letter-spacing:.08em;}",
		},
	}
}
