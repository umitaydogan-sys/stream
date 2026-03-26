package web

// adminHTML contains the complete admin SPA
// In production, this would be served via go:embed from web/admin/ files
// For Phase 1-2, we embed the full HTML directly for simplicity
const adminHTML = `<!DOCTYPE html>
<html lang="tr">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>FluxStream</title>
<!-- Bootstrap Icons -->
<link rel="stylesheet" href="/static/vendor/bootstrap-icons.css">
<link rel="stylesheet" href="/static/admin-studio.css">
<style>
/* Light theme variables */
:root {
  --bg-primary: #eef3f8;
  --bg-secondary: #ffffff;
  --bg-card: #ffffff;
  --bg-card-hover: #f7faff;
  --bg-input: #fbfdff;
  --border: #d5deea;
  --border-focus: #2563eb;
  --text-primary: #1f2a3a;
  --text-secondary: #4f5d71;
  --text-muted: #7d8aa0;
  --accent: #2563eb;
  --accent-hover: #1d4ed8;
  --accent-glow: rgba(37, 99, 235, 0.14);
  --success: #10b981;
  --success-bg: rgba(16, 185, 129, 0.10);
  --danger: #ef4444;
  --danger-bg: rgba(239, 68, 68, 0.10);
  --warning: #f59e0b;
  --warning-bg: rgba(245, 158, 11, 0.12);
  --live-red: #ef4444;
  --live-glow: rgba(239, 68, 68, 0.2);
  --gradient-1: linear-gradient(135deg, #2563eb 0%, #3b82f6 100%);
  --gradient-2: linear-gradient(135deg, #0ea5e9 0%, #2563eb 100%);
  --gradient-3: linear-gradient(135deg, #10b981 0%, #06b6d4 100%);
  --shadow-sm: 0 4px 14px rgba(15, 23, 42, 0.06);
  --shadow-md: 0 10px 24px rgba(15, 23, 42, 0.10);
  --shadow-lg: 0 16px 36px rgba(15, 23, 42, 0.14);
  --radius: 9px;
  --radius-sm: 7px;
  --radius-xs: 5px;
}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'Inter',-apple-system,BlinkMacSystemFont,sans-serif;background:radial-gradient(circle at 12% -8%,#dbeafe 0%,transparent 32%),var(--bg-primary);color:var(--text-primary);min-height:100vh;overflow-x:hidden}
.app{display:flex;min-height:100vh}
.sidebar{width:260px;background:var(--bg-secondary);border-right:1px solid var(--border);display:flex;flex-direction:column;position:fixed;top:0;left:0;bottom:0;z-index:100;transition:transform .3s;box-shadow:var(--shadow-md)}
.main{flex:1;margin-left:260px;background:var(--bg-primary);}
.topbar{height:64px;background:rgba(255,255,255,.95);border-bottom:1px solid var(--border);display:flex;align-items:center;justify-content:space-between;padding:0 24px;position:sticky;top:0;z-index:50;backdrop-filter:blur(12px)}
.topbar-actions{display:flex;align-items:center;gap:10px}
.content{padding:24px;max-width:1400px;margin:0 auto}
.logo{padding:20px;display:flex;align-items:center;gap:12px;border-bottom:1px solid var(--border)}
.logo-icon{width:40px;height:40px;background:var(--gradient-1);border-radius:var(--radius-sm);display:flex;align-items:center;justify-content:center;font-size:20px;box-shadow:var(--shadow-sm)}
.logo-text{font-size:20px;font-weight:700;letter-spacing:-.5px}
.logo-version{font-size:11px;color:var(--text-muted);margin-top:2px}
.nav{flex:1;padding:12px;overflow-y:auto}
.nav-section{margin-bottom:8px}
.nav-section-title{font-size:10px;font-weight:600;text-transform:uppercase;letter-spacing:1.2px;color:var(--text-muted);padding:8px 12px}
.nav-item{display:flex;align-items:center;gap:10px;padding:10px 12px;border-radius:var(--radius-sm);cursor:pointer;transition:all .15s;color:var(--text-secondary);font-size:14px;font-weight:500;text-decoration:none}
.nav-item:hover{background:var(--bg-card);color:var(--text-primary)}
.nav-item.active{background:var(--accent-glow);color:var(--accent-hover);border:1px solid rgba(37,99,235,.25);box-shadow:inset 0 0 0 1px rgba(37,99,235,.08)}
.nav-item .icon{font-size:18px;width:22px;text-align:center}
.card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius);padding:24px;transition:all .2s;box-shadow:var(--shadow-sm)}
.card:hover{border-color:rgba(37,99,235,.35);box-shadow:var(--shadow-md)}
.card-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:16px}
.card-title{font-size:16px;font-weight:600}
.card-grid{display:grid;gap:16px}
.card-grid-4{grid-template-columns:repeat(auto-fit,minmax(200px,1fr))}
.card-grid-3{grid-template-columns:repeat(auto-fit,minmax(280px,1fr))}
.card-grid-2{grid-template-columns:repeat(auto-fit,minmax(380px,1fr))}
.stat-card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius);padding:20px;position:relative;overflow:hidden;box-shadow:var(--shadow-sm)}
.stat-card.clickable{cursor:pointer}
.stat-card.clickable:hover{transform:translateY(-2px);box-shadow:var(--shadow-md);border-color:rgba(37,99,235,.35)}
.stat-card::before{content:'';position:absolute;top:0;left:0;right:0;height:3px}
.stat-card.purple::before{background:var(--gradient-1)}
.stat-card.blue::before{background:var(--gradient-2)}
.stat-card.green::before{background:var(--gradient-3)}
.stat-card.red::before{background:linear-gradient(135deg,#ef4444,#f97316)}
.stat-card.orange::before{background:linear-gradient(135deg,#f59e0b,#f97316)}
.stat-value{font-size:32px;font-weight:800;letter-spacing:-1px;margin:8px 0 4px}
.stat-label{font-size:13px;color:var(--text-muted);font-weight:500}
.stat-subtext{font-size:12px;color:var(--text-muted);margin-top:6px;line-height:1.5}
.stat-icon{font-size:22px;color:var(--accent);display:inline-flex;align-items:center}
.btn{display:inline-flex;align-items:center;gap:8px;padding:10px 20px;border-radius:var(--radius-sm);font-size:14px;font-weight:600;cursor:pointer;border:none;transition:all .2s;font-family:inherit}
.btn-primary{background:var(--gradient-1);color:#fff;box-shadow:var(--shadow-sm)}
.btn-primary:hover{transform:translateY(-1px);box-shadow:var(--shadow-md)}
.btn-secondary{background:var(--bg-card);color:var(--text-primary);border:1px solid var(--border)}
.btn-secondary:hover{background:var(--bg-card-hover)}
.btn-danger{background:var(--danger);color:#fff}
.btn-success{background:var(--success);color:#fff}
.btn-sm{padding:6px 14px;font-size:13px}
.btn-icon{padding:8px;width:36px;height:36px;justify-content:center}
.form-group{margin-bottom:16px}
.form-label{display:block;font-size:13px;font-weight:600;color:var(--text-secondary);margin-bottom:6px}
.form-input,.form-select,.form-textarea{width:100%;padding:9px 13px;background:linear-gradient(180deg,#ffffff 0%,#f7fbff 100%);border:1px solid rgba(148,163,184,.22);border-radius:10px;color:var(--text-primary);font-size:13px;line-height:1.35;font-family:inherit;transition:border-color .2s, box-shadow .2s, transform .2s;outline:none;box-shadow:inset 0 1px 0 rgba(255,255,255,.92),0 6px 16px rgba(15,23,42,.04)}
.form-input:focus,.form-select:focus,.form-textarea:focus{border-color:rgba(37,99,235,.34);box-shadow:0 0 0 3px rgba(37,99,235,.10),0 12px 22px rgba(37,99,235,.08);transform:translateY(-1px)}
.form-input::placeholder{color:var(--text-muted)}
.form-select{appearance:none;-webkit-appearance:none;background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='20' height='20' viewBox='0 0 20 20' fill='none'%3E%3Cpath d='M5 7.5L10 12.5L15 7.5' stroke='%2364758B' stroke-width='1.8' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");background-repeat:no-repeat;background-position:right 11px center;padding-right:36px}
.form-textarea{min-height:110px;resize:vertical;border-radius:12px;line-height:1.45}
.form-hint{font-size:12px;color:var(--text-muted);margin-top:4px}
table{width:100%;border-collapse:collapse}
th{text-align:left;padding:12px 16px;font-size:12px;font-weight:600;text-transform:uppercase;letter-spacing:.5px;color:var(--text-muted);border-bottom:1px solid var(--border)}
td{padding:14px 16px;border-bottom:1px solid rgba(213,222,234,.8);font-size:14px}
tr:hover td{background:rgba(37,99,235,.04)}
.badge{display:inline-flex;align-items:center;gap:6px;padding:4px 10px;border-radius:20px;font-size:12px;font-weight:600}
.badge-live{background:var(--danger-bg);color:var(--live-red);animation:pulse-live 2s infinite}
.badge-offline{background:rgba(100,116,139,.15);color:var(--text-muted)}
.badge-live::before{content:'';width:7px;height:7px;border-radius:50%;background:var(--live-red);box-shadow:0 0 6px var(--live-glow)}
@keyframes pulse-live{0%,100%{opacity:1}50%{opacity:.7}}
.stream-thumb{width:120px;height:68px;background:var(--bg-primary);border-radius:var(--radius-xs);display:flex;align-items:center;justify-content:center;color:var(--text-muted);font-size:24px;position:relative;overflow:hidden;border:1px solid var(--border)}
.stream-thumb.live{border:2px solid var(--live-red)}
.copy-group{display:flex;gap:8px;align-items:center}
.copy-input{flex:1;padding:8px 12px;background:var(--bg-primary);border:1px solid var(--border);border-radius:var(--radius-xs);color:var(--accent-hover);font-size:13px;font-family:'Consolas',monospace}
.copy-btn{padding:8px 12px;background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-xs);color:var(--text-primary);cursor:pointer;font-size:13px;transition:all .2s}
.copy-btn:hover{background:var(--bg-card-hover)}
.modal-overlay{position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,.7);backdrop-filter:blur(4px);z-index:200;display:flex;align-items:center;justify-content:center}
.modal{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius);padding:32px;max-width:600px;width:90%;max-height:80vh;overflow-y:auto;box-shadow:var(--shadow-lg)}
.modal-title{font-size:20px;font-weight:700;margin-bottom:20px}
.wizard-container{min-height:100vh;display:flex;align-items:center;justify-content:center;padding:24px;background:radial-gradient(ellipse at top,rgba(99,102,241,.08) 0%,transparent 60%)}
.wizard-card{background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius);padding:48px;max-width:520px;width:100%;box-shadow:var(--shadow-lg),0 0 60px rgba(37,99,235,.06)}
.wizard-title{font-size:28px;font-weight:800;text-align:center;margin-bottom:8px;background:var(--gradient-1);-webkit-background-clip:text;-webkit-text-fill-color:transparent}
.wizard-subtitle{text-align:center;color:var(--text-muted);font-size:14px;margin-bottom:32px}
.wizard-dot{width:10px;height:10px;border-radius:50%;background:var(--border);transition:all .3s;display:inline-block;margin:0 4px}
.wizard-dot.active{background:var(--accent);box-shadow:0 0 8px var(--accent-glow);width:28px;border-radius:5px}
.wizard-dot.done{background:var(--success)}
.proto-status{display:flex;gap:12px;flex-wrap:wrap}
.proto-dot{display:flex;align-items:center;gap:6px;font-size:12px;font-weight:600;color:var(--text-secondary)}
.proto-dot::before{content:'';width:8px;height:8px;border-radius:50%}
.proto-dot.on::before{background:var(--success);box-shadow:0 0 6px rgba(16,185,129,.4)}
.proto-dot.off::before{background:var(--text-muted)}
.empty-state{text-align:center;padding:60px 20px;color:var(--text-muted)}
.empty-state .icon{font-size:42px;margin-bottom:16px;opacity:.55;display:inline-flex}
.empty-state h3{font-size:18px;color:var(--text-secondary);margin-bottom:8px}
.page-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:24px}
.page-title{font-size:24px;font-weight:700;letter-spacing:-.5px}
.toast{position:fixed;bottom:24px;right:24px;background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-sm);padding:14px 20px;box-shadow:var(--shadow-lg);z-index:300;font-size:14px;font-weight:500;transform:translateY(100px);opacity:0;transition:all .3s}
.toast.show{transform:translateY(0);opacity:1}
.toast.success{border-left:4px solid var(--success)}
.toast.error{border-left:4px solid var(--danger)}
.tabs{display:flex;gap:4px;margin-bottom:20px;border-bottom:1px solid var(--border);padding-bottom:0;overflow-x:auto}
.tab{padding:10px 18px;cursor:pointer;font-size:13px;font-weight:600;color:var(--text-muted);border-bottom:2px solid transparent;transition:all .2s;white-space:nowrap}
.tab:hover{color:var(--text-primary)}
.tab.active{color:var(--accent);border-bottom-color:var(--accent)}
.toggle{position:relative;width:44px;height:24px;display:inline-block}
.toggle input{display:none}
.toggle-slider{position:absolute;top:0;left:0;right:0;bottom:0;background:var(--border);border-radius:12px;cursor:pointer;transition:.3s}
.toggle-slider::before{content:'';position:absolute;height:18px;width:18px;left:3px;bottom:3px;background:#fff;border-radius:50%;transition:.3s}
.toggle input:checked+.toggle-slider{background:var(--accent)}
.toggle input:checked+.toggle-slider::before{transform:translateX(20px)}
.setting-row{display:flex;align-items:center;justify-content:space-between;padding:14px 0;border-bottom:1px solid rgba(45,53,72,.5)}
.setting-row:last-child{border-bottom:none}
.setting-label{font-size:14px;font-weight:500}
.setting-desc{font-size:12px;color:var(--text-muted);margin-top:2px}
.tag{display:inline-block;padding:2px 8px;border-radius:4px;font-size:11px;font-weight:600;margin:2px}
.tag-green{background:var(--success-bg);color:var(--success)}
.tag-yellow{background:var(--warning-bg);color:var(--warning)}
.tag-red{background:var(--danger-bg);color:var(--danger)}
.tag-blue{background:rgba(99,102,241,.1);color:var(--accent)}
.title-icon{margin-right:8px;color:var(--accent);font-size:15px}
.quick-grid{display:grid;gap:16px;grid-template-columns:2fr 1fr}
.insight-grid{display:grid;gap:16px;grid-template-columns:repeat(auto-fit,minmax(260px,1fr))}
.metric-list{display:flex;flex-direction:column;gap:10px}
.metric-row{display:flex;justify-content:space-between;align-items:center;gap:12px;padding:10px 0;border-bottom:1px solid rgba(213,222,234,.8)}
.metric-row:last-child{border-bottom:none}
.bar-list{display:flex;flex-direction:column;gap:12px}
.bar-item{display:grid;grid-template-columns:140px 1fr 56px;gap:12px;align-items:center}
.bar-track{height:10px;background:var(--bg-primary);border-radius:999px;overflow:hidden;border:1px solid var(--border)}
.bar-fill{height:100%;background:var(--gradient-2);border-radius:999px}
.timeline-meta{font-size:11px;color:var(--text-muted);margin-bottom:12px}
.sparkline-shell{position:relative}
.sparkline-frame{position:relative;border-radius:14px;background:linear-gradient(180deg,#f8fbff 0%,#eef4fb 100%);border:1px solid rgba(37,99,235,.08);padding:14px 14px 10px;min-height:178px;overflow:hidden}
.sparkline-svg{display:block;width:100%;height:118px}
.sparkline-grid line{stroke:rgba(37,99,235,.08);stroke-width:1}
.sparkline-area{fill:url(#sparkline-fill)}
.sparkline-line{fill:none;stroke:#14b8a6;stroke-width:3;stroke-linecap:round;stroke-linejoin:round;filter:drop-shadow(0 8px 18px rgba(20,184,166,.18))}
.sparkline-point{fill:#ffffff;stroke:#14b8a6;stroke-width:2}
.sparkline-hitmap{position:absolute;left:14px;right:14px;top:14px;height:118px;display:grid;gap:0}
.sparkline-hit{position:relative;height:100%;cursor:default}
.sparkline-hit::after{content:attr(data-tooltip);position:absolute;left:50%;top:8px;transform:translate(-50%,-110%);background:rgba(15,23,42,.96);color:#fff;padding:7px 9px;border-radius:8px;font-size:11px;line-height:1.35;white-space:nowrap;box-shadow:var(--shadow-md);pointer-events:none;opacity:0;transition:opacity .15s ease;z-index:3}
.sparkline-hit::before{content:'';position:absolute;left:50%;top:8px;transform:translate(-50%,-40%);border:6px solid transparent;border-top-color:rgba(15,23,42,.96);pointer-events:none;opacity:0;transition:opacity .15s ease;z-index:2}
.sparkline-hit:hover::after,.sparkline-hit:hover::before{opacity:1}
.sparkline-footer{display:flex;justify-content:space-between;align-items:flex-start;gap:12px;margin-top:12px}
.sparkline-axis{display:grid;gap:8px;font-size:11px;color:var(--text-muted);flex:1}
.sparkline-axis span{white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.sparkline-summary{display:flex;gap:16px;flex-wrap:wrap;justify-content:flex-end}
.sparkline-chip{min-width:82px}
.sparkline-chip strong{display:block;font-size:15px;color:var(--text-primary)}
.sparkline-chip span{display:block;font-size:11px;color:var(--text-muted)}
.template-thumb{height:150px;border-radius:12px;display:flex;align-items:stretch;justify-content:center;position:relative;overflow:hidden;margin-bottom:12px;border:1px solid rgba(255,255,255,.1);box-shadow:inset 0 0 0 1px rgba(255,255,255,.06)}
.template-thumb-shell{position:relative;display:flex;flex-direction:column;justify-content:space-between;width:100%;height:100%;padding:14px}
.template-thumb-shell::before{content:'';position:absolute;inset:0;background:radial-gradient(circle at top right,rgba(255,255,255,.12),transparent 35%)}
.template-thumb-header,.template-thumb-footer,.template-thumb-center{position:relative;z-index:1}
.template-thumb-header{display:flex;justify-content:space-between;align-items:flex-start;gap:10px}
.template-thumb-title{font-size:13px;font-weight:700;color:#fff;max-width:70%;text-shadow:0 2px 10px rgba(0,0,0,.35)}
.template-thumb-logo{font-size:11px;font-weight:700;padding:5px 8px;border-radius:999px;background:rgba(255,255,255,.18);color:#fff;backdrop-filter:blur(6px)}
.template-thumb-badge{display:inline-flex;align-items:center;gap:6px;padding:4px 10px;border-radius:999px;background:rgba(255,255,255,.14);font-size:10px;font-weight:700;color:#fff;text-transform:uppercase;letter-spacing:.08em}
.template-thumb-badge::before{content:'';width:6px;height:6px;border-radius:50%;background:#fb7185;box-shadow:0 0 10px rgba(251,113,133,.7)}
.template-thumb-center{display:flex;align-items:center;justify-content:center}
.template-thumb-play{width:58px;height:58px;border-radius:50%;display:flex;align-items:center;justify-content:center;background:rgba(255,255,255,.18);backdrop-filter:blur(10px);color:#fff;font-size:26px;box-shadow:0 10px 22px rgba(2,6,23,.28)}
.template-thumb-footer{display:flex;flex-direction:column;gap:10px}
.template-thumb-progress{height:5px;border-radius:999px;background:rgba(255,255,255,.14);overflow:hidden}
.template-thumb-progress span{display:block;height:100%;width:58%;border-radius:999px;background:linear-gradient(90deg,#60a5fa,#22d3ee)}
.template-thumb-controls{display:flex;justify-content:space-between;align-items:center;gap:10px;padding:10px 12px;border-radius:12px;color:#fff}
.template-thumb-controls .left,.template-thumb-controls .right{display:flex;align-items:center;gap:10px;font-size:12px}
.template-thumb-watermark{position:absolute;left:14px;bottom:56px;font-size:11px;letter-spacing:.08em;font-weight:700;color:rgba(255,255,255,.7)}
.segment-control{display:inline-flex;gap:6px;padding:4px;background:var(--bg-primary);border:1px solid var(--border);border-radius:999px}
.segment-btn{border:none;background:transparent;color:var(--text-muted);font-size:12px;font-weight:700;padding:7px 12px;border-radius:999px;cursor:pointer;transition:all .2s}
.segment-btn.active{background:#fff;color:var(--text-primary);box-shadow:var(--shadow-sm)}
.storage-shell{display:flex;flex-direction:column;gap:16px}
.storage-hero{display:grid;gap:16px;grid-template-columns:1.2fr .8fr}
.storage-choice-grid{display:grid;gap:12px;grid-template-columns:repeat(auto-fit,minmax(170px,1fr))}
.storage-choice-card{border:1px solid var(--border);border-radius:16px;background:linear-gradient(180deg,#fff 0%,#f7fbff 100%);padding:16px;cursor:pointer;transition:transform .18s ease,box-shadow .18s ease,border-color .18s ease}
.storage-choice-card:hover{transform:translateY(-1px);box-shadow:var(--shadow-md)}
.storage-choice-card.active{border-color:rgba(37,99,235,.35);box-shadow:0 10px 28px rgba(37,99,235,.12)}
.storage-choice-card .icon{width:40px;height:40px;border-radius:12px;display:flex;align-items:center;justify-content:center;font-size:18px;margin-bottom:10px;background:rgba(37,99,235,.08);color:var(--accent)}
.storage-choice-card .title{font-size:14px;font-weight:700;color:var(--text-primary);margin-bottom:4px}
.storage-choice-card .desc{font-size:12px;color:var(--text-muted);line-height:1.6}
.storage-target-shell{display:flex;flex-direction:column;gap:14px}
.storage-target-top{display:flex;justify-content:space-between;align-items:flex-start;gap:16px;flex-wrap:wrap}
.storage-target-meta{display:flex;flex-wrap:wrap;gap:8px}
.storage-pill{display:inline-flex;align-items:center;gap:6px;padding:6px 10px;border-radius:999px;background:var(--bg-primary);border:1px solid var(--border);font-size:12px;font-weight:700;color:var(--text-secondary)}
.storage-provider-grid{display:grid;gap:12px;grid-template-columns:repeat(auto-fit,minmax(170px,1fr))}
.storage-form-grid{display:grid;gap:14px;grid-template-columns:repeat(2,minmax(0,1fr))}
.storage-subtle{padding:14px 16px;border-radius:14px;background:linear-gradient(180deg,#fbfdff 0%,#f5f9ff 100%);border:1px solid rgba(37,99,235,.08)}
.storage-kpi-list{display:grid;gap:10px;grid-template-columns:repeat(auto-fit,minmax(160px,1fr))}
.storage-kpi{padding:12px 14px;border-radius:14px;background:var(--bg-primary);border:1px solid var(--border)}
.storage-kpi strong{display:block;font-size:18px;color:var(--text-primary)}
.storage-kpi span{display:block;font-size:12px;color:var(--text-muted);margin-top:4px}
.storage-test-row{display:flex;gap:10px;flex-wrap:wrap;align-items:center}
.storage-provider-note{font-size:12px;color:var(--text-muted);line-height:1.7}
.viewer-table td,.viewer-table th{font-size:12px}
.mono-wrap{font-family:Consolas,monospace;font-size:12px;word-break:break-all;color:var(--text-secondary)}
@media(max-width:980px){.quick-grid{grid-template-columns:1fr}.storage-hero,.storage-form-grid{grid-template-columns:1fr}}
@media(max-width:768px){.sidebar{transform:translateX(-100%)}.sidebar.open{transform:translateX(0)}.main{margin-left:0}.card-grid-4,.card-grid-3,.card-grid-2{grid-template-columns:1fr}}
::-webkit-scrollbar{width:6px}::-webkit-scrollbar-track{background:var(--bg-primary)}::-webkit-scrollbar-thumb{background:var(--border);border-radius:3px}
.hidden{display:none!important}
</style>
</head>
<body>
<div id="app"></div>
<div id="toast" class="toast"></div>
<script>
const API='';
let currentPage='dashboard';
let setupCompleted=false;
let authToken=sessionStorage.getItem('fluxstream_token')||'';
let pageRefreshTimer=null;
let streamTelemetryTimer=null;
let currentLang=localStorage.getItem('fluxstream_lang')||'tr';
let runtimeSettings={};
const operationsCenterState={sourceType:'streams',streamID:0,tab:'general',filter:'all'};
const LANGUAGE_META={
  tr:{label:'Turkce',locale:'tr-TR'},
  en:{label:'English',locale:'en-US'},
  de:{label:'Deutsch',locale:'de-DE'},
  es:{label:'Espanol',locale:'es-ES'},
  fr:{label:'Francais',locale:'fr-FR'}
};
const I18N={
  en:{
    'Ana Menu':'Main Menu','Yayin':'Streaming','Ayarlar':'Settings','Izleme':'Monitoring','Sistem':'System',
    'Yayinlar':'Streams','Yeni Yayin':'New Stream','Embed Kodlari':'Embed Codes','Gelismis Embed':'Advanced Embed','Player Sablonlari':'Player Templates',
    'Kolay Ayarlar':'Quick Settings','Genel':'General','Alan Adi / Embed':'Domain / Embed','Protokoller':'Protocols','Cikis Formatlari':'Output Formats','Teslimat / ABR':'Delivery / ABR','SSL/TLS':'SSL/TLS','Guvenlik':'Security','Depolama':'Storage','Saglik ve Uyari':'Health & Alerts','Transkod':'Transcode',
    'Analitik':'Analytics','Kayitlar':'Recordings','Izleyiciler':'Viewers','Transcode Isleri':'Transcode Jobs','Teshis':'Diagnostics','Bakim ve Yedek':'Maintenance & Backups','Lisans':'License','Tokenlar':'Tokens','Kullanicilar':'Users','Loglar':'Logs',
    'Yonetim paneline giris yapin':'Sign in to the admin panel','Kullanici Adi':'Username','Sifre':'Password','Sifreniz':'Your password','Giris Yap':'Sign In','Kullanici adi ve sifre gerekli':'Username and password are required','Giris basarili!':'Login successful!','Giris hatasi':'Login error',
    'Live Streaming Media Server':'Live Streaming Media Server',"FluxStream'e hos geldiniz!":'Welcome to FluxStream!','Canli yayin sunucunuzu birkac adimda kuralim.':'Let’s set up your live streaming server in a few steps.','Baslayalim':'Let’s Start','Admin Hesabi':'Admin Account','Yonetim paneli icin giris bilgileri':'Sign-in details for the admin panel','Sifre Tekrar':'Repeat Password','En az 4 karakter':'At least 4 characters','Ileri':'Next','Geri':'Back','Port ve Domain':'Ports and Domain','Sunucu portlarini ve public alan adini yapilandirin':'Configure server ports and the public domain','HTTP Port (Web Arayuzu)':'HTTP Port (Web UI)','HTTPS Port (SSL aktifse)':'HTTPS Port (if SSL is enabled)','RTMP Port (OBS Yayin)':'RTMP Port (OBS ingest)','Public Domain / IP':'Public Domain / IP','Bos birakirsaniz panelin acildigi host kullanilir. HTTP ve HTTPS public portlari kurulumdan sonra Kolay Ayarlar veya Alan Adi / Embed ekranindan degistirilebilir.':'If left empty, the host used to open the panel is used. Public HTTP and HTTPS ports can be changed later from Quick Settings or the Domain / Embed page.','Kurulumu Tamamla':'Finish Setup','Kurulum tamamlandi!':'Setup completed!','Kurulum hatasi':'Setup error','Sifre en az 4 karakter olmali':'Password must be at least 4 characters','Sifreler eslesiyor!':'Passwords do not match!',
    'Kaydet':'Save','Iptal':'Cancel','Guncelle':'Update','Olustur':'Create','Sil':'Delete','Duzenle':'Edit','Onizle':'Preview','Indir':'Download','Direkt Link':'Direct Link',
    'Toplam Yayin':'Total Streams','Aktif Izleyici':'Active Viewers','Tepe Esz.':'Peak Concurrency','Toplam Bant':'Total Bandwidth','Izleyici Trendi':'Viewer Trend','Bant Trendi':'Bandwidth Trend','Format Dagilimi':'Format Distribution','Ulke Dagilimi':'Country Distribution','En Populer Yayinlar':'Top Streams','Secili periyot':'Selected period','Ayni periyotta toplam cikis':'Total output in the same period','Henuz timeline verisi yok':'No timeline data yet','Henuz bant snapshot yok':'No bandwidth snapshots yet','Henuz format verisi yok':'No format data yet','Henuz ulke verisi yok':'No country data yet','Aktif yayin yok':'No active stream','izleyici':'viewer',
    'Sunucu Kontrol':'Server Control','Yeniden Baslat':'Restart','Durdur':'Stop','Kapat':'Close','Kopyalandi!':'Copied!','Kopyalama basarisiz':'Copy failed','Son':'Latest','Tepe':'Peak','Minimum':'Minimum','Ayarlar kaydedildi!':'Settings saved!','Kayit hatasi':'Save error',
    'Genel Ayarlar':'General Settings','Kimlik ve Yerellesme':'Identity and Localization','Sunucu Adi':'Server Name','Sunucu goruntuleme adi':'Display name of the server','Dil':'Language','Kurulumda secilen dil burada degistirilebilir. Login, setup ve panel kabugu bu secime gore acilir.':'The installation language can be changed here. Login, setup, and the admin shell follow this choice.','Saat Dilimi':'Time Zone','Tarih ve saat gosterimleri bu timezone ile yorumlanir.':'Dates and times are interpreted using this time zone.','Tema':'Theme','Admin panelinin gorsel yonunu belirler. Su an acik tema varsayilandir.':'Defines the visual style of the admin panel. Light theme is the default for now.','Kolay mod acik':'Guided mode enabled','Yeni kurulumlarda rehber odakli ayarlari one cikarir.':'Highlights guidance-first settings during new installations.',
    'Sunucu ve Panel Varsayilanlari':'Server and Panel Defaults','Web arayuzu portu':'Web interface port','SSL portu':'SSL port','Varsayilan Public Domain':'Default Public Domain','Link uretiminde kullanilan ilk alan adi. Bossa mevcut host kullanilir.':'Primary domain used while generating links. If empty, the current host is used.','Varsayilan Public HTTP Port':'Default Public HTTP Port','Embed ve player linkleri icin':'Used for player and embed links','Varsayilan Public HTTPS Port':'Default Public HTTPS Port','SSL ile uretilen linkler icin':'Used for SSL-generated links','Player kalite secici':'Player quality selector','ABR yayinlarda kullanici kaliteyi elle de secebilir.':'Lets viewers manually choose quality on ABR streams.','Otomatik bakim':'Automatic maintenance','Temizleme ve bakim islerini zamanli calistirir.':'Runs cleanup and maintenance tasks on schedule.','Kayit Saklama Suresi (gun)':'Recording Retention (days)','0 verilirse otomatik silme yapilmaz.':'Set to 0 to disable automatic deletion.',
    'Baglanti Rehberi':'Connection Guide','OBS RTMP URL':'OBS RTMP URL','RTP URL':'RTP URL','HLS Izleme URL':'HLS Playback URL',
    'Kurulu gelen hazir sablonlari temel alip duzenleyebilir veya sifirdan yeni sablon olusturabilirsiniz.':'Use the bundled starter templates as a base, or create a new one from scratch.','+ Yeni Sablon':'+ New Template','Onizleme Kaynagi':'Preview Source','Onizleme Formati':'Preview Format','Kaydedilen template icin bu formatta embed kodu ve preview olusur.':'Preview and embed code are generated in this format for the saved template.','Hazir baslangic sablonu':'Starter templates','Kullanim':'Usage','Duzenle -> Kaydet -> Embed tarafinda kullan':'Edit -> Save -> Use on the embed side','Amac':'Purpose','Canli TV, radyo, minimal player, cam tasarim ve parlak vitrini hizla baslatmak':'Quick-start templates for live TV, radio, minimal player, glass style, and showcase layouts.',
    'Kaynak stream yok':'No source stream','Template preview icin en az bir stream olusturun.':'Create at least one stream to preview templates.','Kaydedin ve deneyin':'Save and try','Yeni bir template icin once kaydet, sonra secili stream ile player preview ve embed kodunu gor.':'For a new template, save it first, then review the player preview and embed code with the selected stream.','Secili kaynak:':'Selected source:','Sablon Duzenle':'Edit Template','Yeni Player Sablonu':'New Player Template','Sablon Adi *':'Template Name *','Logo URL':'Logo URL','Logo Konum':'Logo Position','Sag Ust':'Top Right','Sol Ust':'Top Left','Sag Alt':'Bottom Right','Sol Alt':'Bottom Left','Logo Seffaflik':'Logo Opacity','Watermark Yazi':'Watermark Text','Baslik Goster':'Show Title','CANLI Badge':'LIVE Badge','Arkaplan CSS':'Background CSS','Kontrol Cubugu CSS':'Control Bar CSS','Play Butonu CSS':'Play Button CSS','Ozel CSS':'Custom CSS','Kaynak stream':'Source stream','Format':'Format','Canli Player Onizleme':'Live Player Preview','Secili stream ve format ile':'Using the selected stream and format','Template + stream birlesik cikti':'Combined template + stream output','Sablon adi gerekli':'Template name is required','Sablon guncellendi!':'Template updated!','Sablon olusturuldu!':'Template created!','Sablon silindi':'Template deleted',
    'Servis durumu, tek tikla yedek alma ve temiz geri donus komutlari burada toplanir.':'Service status, one-click backups, and clean recovery actions are gathered here.','Offline imzali lisans dosyasi burada saklanir. Internet baglantisi olmadan dogrulama yapilir.':'The offline signed license file is stored here. Validation works without internet access.','Mevcut Lisans Durumu':'Current License Status','Bekleniyor':'Pending','Lisans ID':'License ID','Lisans yuklenince aktif ozellikler burada gorunur.':'Active licensed features appear here after a license is loaded.','Lisans Dosyasi Yukle':'Upload License File','Lisans JSON':'License JSON','Imzali lisans JSONunu buraya yapistirin':'Paste the signed license JSON here','Lisansi Kaydet':'Save License','Ornek JSON Yukle':'Load Sample JSON','Lisans kaydedildi':'License saved','Lisans kaydedilemedi':'License could not be saved'
  },
  de:{
    'Ana Menu':'Hauptmenu','Yayin':'Streaming','Ayarlar':'Einstellungen','Izleme':'Ueberwachung','Sistem':'System',
    'Yayinlar':'Streams','Yeni Yayin':'Neuer Stream','Embed Kodlari':'Embed-Codes','Gelismis Embed':'Erweitertes Embed','Player Sablonlari':'Player-Vorlagen',
    'Kolay Ayarlar':'Schnelleinstellungen','Genel':'Allgemein','Alan Adi / Embed':'Domain / Embed','Protokoller':'Protokolle','Cikis Formatlari':'Ausgabeformate','Teslimat / ABR':'Auslieferung / ABR','SSL/TLS':'SSL/TLS','Guvenlik':'Sicherheit','Depolama':'Speicher','Saglik ve Uyari':'Status und Warnungen','Transkod':'Transkodierung',
    'Analitik':'Analytik','Kayitlar':'Aufnahmen','Izleyiciler':'Zuschauer','Transcode Isleri':'Transcode-Jobs','Teshis':'Diagnose','Bakim ve Yedek':'Wartung und Backups','Lisans':'Lizenz','Tokenlar':'Token','Kullanicilar':'Benutzer','Loglar':'Protokolle',
    'Yonetim paneline giris yapin':'Am Admin-Panel anmelden','Kullanici Adi':'Benutzername','Sifre':'Passwort','Sifreniz':'Ihr Passwort','Giris Yap':'Anmelden','Kullanici adi ve sifre gerekli':'Benutzername und Passwort sind erforderlich','Giris basarili!':'Anmeldung erfolgreich!','Giris hatasi':'Anmeldefehler',
    'Live Streaming Media Server':'Live-Streaming-Medienserver',"FluxStream'e hos geldiniz!":'Willkommen bei FluxStream!','Canli yayin sunucunuzu birkac adimda kuralim.':'Richten wir Ihren Live-Streaming-Server in wenigen Schritten ein.','Baslayalim':'Los geht’s','Admin Hesabi':'Admin-Konto','Yonetim paneli icin giris bilgileri':'Anmeldedaten fuer das Admin-Panel','Sifre Tekrar':'Passwort wiederholen','En az 4 karakter':'Mindestens 4 Zeichen','Ileri':'Weiter','Geri':'Zurueck','Port ve Domain':'Ports und Domain','Sunucu portlarini ve public alan adini yapilandirin':'Server-Ports und die oeffentliche Domain konfigurieren','HTTP Port (Web Arayuzu)':'HTTP-Port (Weboberflaeche)','HTTPS Port (SSL aktifse)':'HTTPS-Port (wenn SSL aktiv ist)','RTMP Port (OBS Yayin)':'RTMP-Port (OBS-Ingest)','Public Domain / IP':'Oeffentliche Domain / IP','Bos birakirsaniz panelin acildigi host kullanilir. HTTP ve HTTPS public portlari kurulumdan sonra Kolay Ayarlar veya Alan Adi / Embed ekranindan degistirilebilir.':'Wenn leer, wird der Host der aktuellen Panel-URL verwendet. Oeffentliche HTTP- und HTTPS-Ports koennen spaeter im Bereich Schnelleinstellungen oder Domain / Embed geaendert werden.','Kurulumu Tamamla':'Einrichtung abschliessen','Kurulum tamamlandi!':'Einrichtung abgeschlossen!','Kurulum hatasi':'Einrichtungsfehler','Sifre en az 4 karakter olmali':'Das Passwort muss mindestens 4 Zeichen lang sein','Sifreler eslesiyor!':'Die Passwoerter stimmen nicht ueberein!',
    'Kaydet':'Speichern','Iptal':'Abbrechen','Guncelle':'Aktualisieren','Olustur':'Erstellen','Sil':'Loeschen','Duzenle':'Bearbeiten','Onizle':'Vorschau','Indir':'Herunterladen','Direkt Link':'Direktlink',
    'Toplam Yayin':'Gesamtzahl Streams','Aktif Izleyici':'Aktive Zuschauer','Tepe Esz.':'Spitzenwert','Toplam Bant':'Gesamtbandbreite','Izleyici Trendi':'Zuschauertrend','Bant Trendi':'Bandbreitentrend','Format Dagilimi':'Formatverteilung','Ulke Dagilimi':'Laenderverteilung','En Populer Yayinlar':'Beliebteste Streams','Secili periyot':'Ausgewaehlter Zeitraum','Ayni periyotta toplam cikis':'Gesamtausgabe im selben Zeitraum','Henuz timeline verisi yok':'Noch keine Zeitreihendaten','Henuz bant snapshot yok':'Noch keine Bandbreiten-Snapshots','Henuz format verisi yok':'Noch keine Formatdaten','Henuz ulke verisi yok':'Noch keine Laenderdaten','Aktif yayin yok':'Kein aktiver Stream','izleyici':'Zuschauer',
    'Sunucu Kontrol':'Serversteuerung','Yeniden Baslat':'Neu starten','Durdur':'Beenden','Kapat':'Schliessen','Kopyalandi!':'Kopiert!','Kopyalama basarisiz':'Kopieren fehlgeschlagen','Son':'Aktuell','Tepe':'Spitze','Minimum':'Minimum','Ayarlar kaydedildi!':'Einstellungen gespeichert!','Kayit hatasi':'Speicherfehler'
  },
  es:{
    'Ana Menu':'Menu principal','Yayin':'Streaming','Ayarlar':'Configuracion','Izleme':'Monitoreo','Sistem':'Sistema',
    'Yayinlar':'Streams','Yeni Yayin':'Nuevo stream','Embed Kodlari':'Codigos embed','Gelismis Embed':'Embed avanzado','Player Sablonlari':'Plantillas de reproductor',
    'Kolay Ayarlar':'Ajustes rapidos','Genel':'General','Alan Adi / Embed':'Dominio / Embed','Protokoller':'Protocolos','Cikis Formatlari':'Formatos de salida','Teslimat / ABR':'Entrega / ABR','SSL/TLS':'SSL/TLS','Guvenlik':'Seguridad','Depolama':'Almacenamiento','Saglik ve Uyari':'Salud y alertas','Transkod':'Transcodificacion',
    'Analitik':'Analitica','Kayitlar':'Grabaciones','Izleyiciler':'Espectadores','Transcode Isleri':'Tareas de transcodificacion','Teshis':'Diagnostico','Bakim ve Yedek':'Mantenimiento y copias','Lisans':'Licencia','Tokenlar':'Tokens','Kullanicilar':'Usuarios','Loglar':'Registros',
    'Yonetim paneline giris yapin':'Inicia sesion en el panel de administracion','Kullanici Adi':'Usuario','Sifre':'Contrasena','Sifreniz':'Tu contrasena','Giris Yap':'Entrar','Kullanici adi ve sifre gerekli':'Se requieren usuario y contrasena','Giris basarili!':'Inicio de sesion correcto','Giris hatasi':'Error de inicio de sesion',
    'Live Streaming Media Server':'Servidor de streaming en vivo',"FluxStream'e hos geldiniz!":'Bienvenido a FluxStream!','Canli yayin sunucunuzu birkac adimda kuralim.':'Configuremos tu servidor de streaming en pocos pasos.','Baslayalim':'Empecemos','Admin Hesabi':'Cuenta de administrador','Yonetim paneli icin giris bilgileri':'Credenciales para el panel de administracion','Sifre Tekrar':'Repetir contrasena','En az 4 karakter':'Al menos 4 caracteres','Ileri':'Siguiente','Geri':'Atras','Port ve Domain':'Puertos y dominio','Sunucu portlarini ve public alan adini yapilandirin':'Configura los puertos del servidor y el dominio publico','HTTP Port (Web Arayuzu)':'Puerto HTTP (interfaz web)','HTTPS Port (SSL aktifse)':'Puerto HTTPS (si SSL esta activo)','RTMP Port (OBS Yayin)':'Puerto RTMP (ingesta OBS)','Public Domain / IP':'Dominio publico / IP','Bos birakirsaniz panelin acildigi host kullanilir. HTTP ve HTTPS public portlari kurulumdan sonra Kolay Ayarlar veya Alan Adi / Embed ekranindan degistirilebilir.':'Si se deja vacio, se usa el host con el que se abre el panel. Los puertos publicos HTTP y HTTPS pueden cambiarse despues desde Ajustes rapidos o Dominio / Embed.','Kurulumu Tamamla':'Finalizar instalacion','Kurulum tamamlandi!':'Instalacion completada!','Kurulum hatasi':'Error de instalacion','Sifre en az 4 karakter olmali':'La contrasena debe tener al menos 4 caracteres','Sifreler eslesiyor!':'Las contrasenas no coinciden!',
    'Kaydet':'Guardar','Iptal':'Cancelar','Guncelle':'Actualizar','Olustur':'Crear','Sil':'Eliminar','Duzenle':'Editar','Onizle':'Vista previa','Indir':'Descargar','Direkt Link':'Enlace directo',
    'Toplam Yayin':'Streams totales','Aktif Izleyici':'Espectadores activos','Tepe Esz.':'Pico concurrente','Toplam Bant':'Ancho de banda total','Izleyici Trendi':'Tendencia de audiencia','Bant Trendi':'Tendencia de ancho de banda','Format Dagilimi':'Distribucion por formato','Ulke Dagilimi':'Distribucion por pais','En Populer Yayinlar':'Streams mas populares','Secili periyot':'Periodo seleccionado','Ayni periyotta toplam cikis':'Salida total del mismo periodo','Henuz timeline verisi yok':'Aun no hay datos de la linea temporal','Henuz bant snapshot yok':'Aun no hay capturas de ancho de banda','Henuz format verisi yok':'Aun no hay datos por formato','Henuz ulke verisi yok':'Aun no hay datos por pais','Aktif yayin yok':'No hay stream activo','izleyici':'espectador',
    'Sunucu Kontrol':'Control del servidor','Yeniden Baslat':'Reiniciar','Durdur':'Detener','Kapat':'Cerrar','Kopyalandi!':'Copiado!','Kopyalama basarisiz':'Error al copiar','Son':'Actual','Tepe':'Pico','Minimum':'Minimo','Ayarlar kaydedildi!':'Configuracion guardada!','Kayit hatasi':'Error al guardar'
  },
  fr:{
    'Ana Menu':'Menu principal','Yayin':'Streaming','Ayarlar':'Parametres','Izleme':'Supervision','Sistem':'Systeme',
    'Yayinlar':'Streams','Yeni Yayin':'Nouveau stream','Embed Kodlari':'Codes embed','Gelismis Embed':'Embed avance','Player Sablonlari':'Modeles de lecteur',
    'Kolay Ayarlar':'Reglages rapides','Genel':'General','Alan Adi / Embed':'Domaine / Embed','Protokoller':'Protocoles','Cikis Formatlari':'Formats de sortie','Teslimat / ABR':'Distribution / ABR','SSL/TLS':'SSL/TLS','Guvenlik':'Securite','Depolama':'Stockage','Saglik ve Uyari':'Sante et alertes','Transkod':'Transcodage',
    'Analitik':'Analytique','Kayitlar':'Enregistrements','Izleyiciler':'Spectateurs','Transcode Isleri':'Taches de transcodage','Teshis':'Diagnostic','Bakim ve Yedek':'Maintenance et sauvegardes','Lisans':'Licence','Tokenlar':'Jetons','Kullanicilar':'Utilisateurs','Loglar':'Journaux',
    'Yonetim paneline giris yapin':'Connectez-vous au panneau d’administration','Kullanici Adi':'Nom d’utilisateur','Sifre':'Mot de passe','Sifreniz':'Votre mot de passe','Giris Yap':'Connexion','Kullanici adi ve sifre gerekli':'Nom d’utilisateur et mot de passe requis','Giris basarili!':'Connexion reussie !','Giris hatasi':'Erreur de connexion',
    'Live Streaming Media Server':'Serveur multimedia de streaming en direct',"FluxStream'e hos geldiniz!":'Bienvenue sur FluxStream !','Canli yayin sunucunuzu birkac adimda kuralim.':'Configurons votre serveur de streaming en quelques etapes.','Baslayalim':'Commencons','Admin Hesabi':'Compte administrateur','Yonetim paneli icin giris bilgileri':'Identifiants du panneau d’administration','Sifre Tekrar':'Repeter le mot de passe','En az 4 karakter':'Au moins 4 caracteres','Ileri':'Suivant','Geri':'Retour','Port ve Domain':'Ports et domaine','Sunucu portlarini ve public alan adini yapilandirin':'Configurez les ports du serveur et le domaine public','HTTP Port (Web Arayuzu)':'Port HTTP (interface web)','HTTPS Port (SSL aktifse)':'Port HTTPS (si SSL est actif)','RTMP Port (OBS Yayin)':'Port RTMP (ingest OBS)','Public Domain / IP':'Domaine public / IP','Bos birakirsaniz panelin acildigi host kullanilir. HTTP ve HTTPS public portlari kurulumdan sonra Kolay Ayarlar veya Alan Adi / Embed ekranindan degistirilebilir.':'Si ce champ est vide, l’hote utilise pour ouvrir le panneau sera repris. Les ports publics HTTP et HTTPS peuvent etre modifies apres l’installation depuis Reglages rapides ou Domaine / Embed.','Kurulumu Tamamla':'Terminer l’installation','Kurulum tamamlandi!':'Installation terminee !','Kurulum hatasi':'Erreur d’installation','Sifre en az 4 karakter olmali':'Le mot de passe doit contenir au moins 4 caracteres','Sifreler eslesiyor!':'Les mots de passe ne correspondent pas !',
    'Kaydet':'Enregistrer','Iptal':'Annuler','Guncelle':'Mettre a jour','Olustur':'Creer','Sil':'Supprimer','Duzenle':'Modifier','Onizle':'Apercu','Indir':'Telecharger','Direkt Link':'Lien direct',
    'Toplam Yayin':'Total des streams','Aktif Izleyici':'Spectateurs actifs','Tepe Esz.':'Pic de simultaneite','Toplam Bant':'Bande passante totale','Izleyici Trendi':'Tendance des spectateurs','Bant Trendi':'Tendance de bande passante','Format Dagilimi':'Repartition par format','Ulke Dagilimi':'Repartition par pays','En Populer Yayinlar':'Streams les plus populaires','Secili periyot':'Periode selectionnee','Ayni periyotta toplam cikis':'Sortie totale sur la meme periode','Henuz timeline verisi yok':'Pas encore de donnees de chronologie','Henuz bant snapshot yok':'Aucun instantane de bande passante','Henuz format verisi yok':'Pas encore de donnees de format','Henuz ulke verisi yok':'Pas encore de donnees par pays','Aktif yayin yok':'Aucun stream actif','izleyici':'spectateur',
    'Sunucu Kontrol':'Controle du serveur','Yeniden Baslat':'Redemarrer','Durdur':'Arreter','Kapat':'Fermer','Kopyalandi!':'Copie !','Kopyalama basarisiz':'La copie a echoue','Son':'Actuel','Tepe':'Pic','Minimum':'Minimum','Ayarlar kaydedildi!':'Parametres enregistres !','Kayit hatasi':'Erreur d’enregistrement'
  }
};

Object.assign(I18N.en,{
  'Runtime Modu':'Runtime Mode',
  'Feature enforcement':'Feature Enforcement',
  'Gelistirme':'Development',
  'Servisi Yeniden Baslat':'Restart Service',
  'Servisi Baslat':'Start Service',
  'Servisi Durdur':'Stop Service',
  'Kurulum Dizini':'Install Directory',
  'Atomic Upgrade Komutu':'Atomic Upgrade Command',
  'Servis aksiyonu gonderildi':'Service action sent',
  'Servis aksiyonu basarisiz':'Service action failed',
  'Embedded development key aktif; production icin imzali lisans yukleyin.':'Embedded development key is active; upload a signed license for production.',
  'Yeni binary once *.next olarak yuklenir, servis durdurulur, atomik rename yapilip servis yeniden baslatilir.':'Upload the new binary as *.next, stop the service, perform an atomic rename, then start the service again.'
});

function normalizeLang(lang){
  return LANGUAGE_META[lang]?lang:'tr';
}
function localeForLang(){
  return (LANGUAGE_META[normalizeLang(currentLang)]||LANGUAGE_META.tr).locale;
}
function fmtLocaleDateTime(value){
  if(!value)return '-';
  const d=new Date(value);
  return Number.isNaN(d.getTime())?'-':d.toLocaleString(localeForLang());
}
function fmtLocaleDate(value){
  if(!value)return '-';
  const d=new Date(value);
  return Number.isNaN(d.getTime())?'-':d.toLocaleDateString(localeForLang());
}
function fmtLocaleTime(value){
  if(!value)return '-';
  const d=new Date(value);
  return Number.isNaN(d.getTime())?'-':d.toLocaleTimeString(localeForLang());
}
function languageOptions(selected){
  selected=normalizeLang(selected||currentLang);
  return Object.keys(LANGUAGE_META).map(function(code){
    return '<option value="'+code+'" '+(code===selected?'selected':'')+'>'+LANGUAGE_META[code].label+'</option>';
  }).join('');
}
function t(key,fallback,vars){
  const lang=normalizeLang(currentLang);
  let text=(I18N[lang]||{})[key];
  if((text===undefined||text===null) && lang!=='en') text=(I18N.en||{})[key];
  if(text===undefined||text===null) text=fallback||key;
  if(vars){
    Object.keys(vars).forEach(function(name){
      text=text.replace(new RegExp('\\{'+name+'\\}','g'),String(vars[name]));
    });
  }
  return text;
}
function translateLiteral(value){
  if(currentLang==='tr'||value==null)return value;
  const raw=String(value);
  const trimmed=raw.trim();
  if(!trimmed)return raw;
  const translated=t(trimmed,trimmed);
  return translated===trimmed?raw:raw.replace(trimmed,translated);
}
function applyTranslations(root){
  if(!root||currentLang==='tr')return;
  root.querySelectorAll('[placeholder],[title],[aria-label]').forEach(function(el){
    ['placeholder','title','aria-label'].forEach(function(attr){
      const val=el.getAttribute(attr);
      if(val)el.setAttribute(attr,translateLiteral(val));
    });
  });
  const walker=document.createTreeWalker(root,NodeFilter.SHOW_TEXT,{
    acceptNode:function(node){
      if(!node.parentElement)return NodeFilter.FILTER_REJECT;
      const tag=node.parentElement.tagName;
      if(['SCRIPT','STYLE','TEXTAREA','CODE','PRE'].indexOf(tag)!==-1)return NodeFilter.FILTER_REJECT;
      return String(node.nodeValue||'').trim()?NodeFilter.FILTER_ACCEPT:NodeFilter.FILTER_REJECT;
    }
  });
  const nodes=[];
  while(walker.nextNode())nodes.push(walker.currentNode);
  nodes.forEach(function(node){
    node.nodeValue=translateLiteral(node.nodeValue);
  });
}
function setCurrentLanguage(lang,silent){
  currentLang=normalizeLang(lang);
  document.documentElement.lang=currentLang;
  localStorage.setItem('fluxstream_lang',currentLang);
  if(typeof wizardData!=='undefined'&&wizardData)wizardData.language=currentLang;
  if(!silent)applyTranslations(document.getElementById('app'));
}

async function api(path,opts={}){
  try{
    const hdrs={'Content-Type':'application/json',...opts.headers};
    if(authToken) hdrs['Authorization']='Bearer '+authToken;
    const res=await fetch(API+path,{
      cache:opts.cache||'no-store',
      headers:hdrs,
      ...opts,
      body:opts.body?JSON.stringify(opts.body):undefined,
    });
    return res.json();
  }catch(e){return {error:true,message:e.message}}
}

function toast(msg,type='success'){
  const el=document.getElementById('toast');
  el.textContent=msg;el.className='toast '+type+' show';
  setTimeout(()=>el.classList.remove('show'),3000);
}

async function copyText(text){
  const value=String(text==null?'':text);
  try{
    if(navigator.clipboard&&window.isSecureContext){
      await navigator.clipboard.writeText(value);
      toast('Kopyalandi!');
      return;
    }
  }catch(e){}
  const ta=document.createElement('textarea');
  ta.value=value;
  ta.setAttribute('readonly','readonly');
  ta.style.position='fixed';
  ta.style.opacity='0';
  ta.style.pointerEvents='none';
  document.body.appendChild(ta);
  ta.focus();
  ta.select();
  try{
    if(document.execCommand('copy'))toast('Kopyalandi!');
    else toast('Kopyalama basarisiz','error');
  }catch(e){
    toast('Kopyalama basarisiz','error');
  }finally{
    document.body.removeChild(ta);
  }
}
function escHtml(s){if(!s)return '';return String(s).replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;')}
function formatBytes(b){if(!b||b===0)return '0 B';const k=1024,s=['B','KB','MB','GB','TB'],i=Math.floor(Math.log(b)/Math.log(k));return parseFloat((b/Math.pow(k,i)).toFixed(1))+' '+s[i]}
function formatUptime(sec){if(!sec)return '0s';const d=Math.floor(sec/86400),h=Math.floor((sec%86400)/3600),m=Math.floor((sec%3600)/60);if(d>0)return d+'g '+h+'s '+m+'dk';if(h>0)return h+'s '+m+'dk';return m+'dk'}

async function init(){
  const [status,settings]=await Promise.all([api('/api/setup/status'),api('/api/settings')]);
  if(settings&&!settings.error)runtimeSettings=settings;
  setCurrentLanguage((runtimeSettings&&runtimeSettings.language)||currentLang,true);
  setupCompleted=status.setup_completed;
  if(!setupCompleted){renderWizard();return}
  if(authToken){
    const me=await api('/api/auth/me');
    if(me.authenticated){renderApp();return}
    authToken='';sessionStorage.removeItem('fluxstream_token');
  }
  renderLogin();
}

function renderLogin(){
  document.getElementById('app').innerHTML=
  '<div class="wizard-container"><div class="wizard-card"><div style="text-align:center;font-size:48px;margin-bottom:16px;color:var(--accent)"><i class="bi bi-lightning-charge-fill"></i></div>'+
  '<h1 class="wizard-title">FluxStream</h1><p class="wizard-subtitle">'+t('Yonetim paneline giris yapin')+'</p>'+
  '<div class="form-group"><label class="form-label">'+t('Kullanici Adi')+'</label><input class="form-input" id="login-user" value="admin"></div>'+
  '<div class="form-group"><label class="form-label">'+t('Sifre')+'</label><input class="form-input" id="login-pass" type="password" placeholder="'+t('Sifreniz')+'"></div>'+
  '<button class="btn btn-primary" style="width:100%" onclick="doLogin()">'+t('Giris Yap')+'</button></div></div>';
  applyTranslations(document.getElementById('app'));
}
async function doLogin(){
  const u=document.getElementById('login-user').value;
  const p=document.getElementById('login-pass').value;
  if(!u||!p){toast(t('Kullanici adi ve sifre gerekli'),'error');return}
  const res=await api('/api/auth/login',{method:'POST',body:{username:u,password:p}});
  if(res.success){authToken=res.token;sessionStorage.setItem('fluxstream_token',authToken);toast(t('Giris basarili!'));renderApp()}
  else{toast(res.message||t('Giris hatasi'),'error')}
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETUP WIZARD ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
let wizardStep=1;
const wizardData={username:'admin',password:'',http_port:8844,https_port:443,rtmp_port:1935,embed_domain:'',language:'tr'};

function renderWizard(){
  document.getElementById('app').innerHTML='<div class="wizard-container">'+getWizardContent()+'</div>';
  applyTranslations(document.getElementById('app'));
}
function getWizardContent(){
  const steps=[
    '<div class="wizard-card"><div style="display:flex;justify-content:flex-end;margin-bottom:12px"><select class="form-select" style="max-width:150px" onchange="setCurrentLanguage(this.value,true);wizardData.language=this.value;renderWizard()">'+languageOptions(wizardData.language||currentLang)+'</select></div><div style="text-align:center;font-size:48px;margin-bottom:16px;color:var(--accent)"><i class="bi bi-lightning-charge-fill"></i></div><h1 class="wizard-title">FluxStream</h1><p class="wizard-subtitle">'+t('Live Streaming Media Server')+'</p><div style="text-align:center;margin-bottom:24px">'+stepDots(1)+'</div><p style="text-align:center;color:var(--text-secondary);margin-bottom:32px;line-height:1.7">'+t("FluxStream'e hos geldiniz!")+'<br>'+t('Canli yayin sunucunuzu birkac adimda kuralim.')+'</p><button class="btn btn-primary" style="width:100%" onclick="wizardNext()">'+t('Baslayalim')+' <i class="bi bi-arrow-right"></i></button></div>',
    '<div class="wizard-card"><h1 class="wizard-title">'+t('Admin Hesabi')+'</h1><p class="wizard-subtitle">'+t('Yonetim paneli icin giris bilgileri')+'</p><div style="text-align:center;margin-bottom:24px">'+stepDots(2)+'</div><div class="form-group"><label class="form-label">'+t('Kullanici Adi')+'</label><input class="form-input" id="w-username" value="'+escHtml(wizardData.username||'admin')+'"></div><div class="form-group"><label class="form-label">'+t('Sifre')+'</label><input class="form-input" id="w-password" type="password" placeholder="'+t('En az 4 karakter')+'"></div><div class="form-group"><label class="form-label">'+t('Sifre Tekrar')+'</label><input class="form-input" id="w-password2" type="password"></div><div style="display:flex;gap:12px"><button class="btn btn-secondary" style="flex:1" onclick="wizardPrev()"><i class="bi bi-arrow-left"></i> '+t('Geri')+'</button><button class="btn btn-primary" style="flex:1" onclick="wizardNext()">'+t('Ileri')+' <i class="bi bi-arrow-right"></i></button></div></div>',
    '<div class="wizard-card"><h1 class="wizard-title">'+t('Port ve Domain')+'</h1><p class="wizard-subtitle">'+t('Sunucu portlarini ve public alan adini yapilandirin')+'</p><div style="text-align:center;margin-bottom:24px">'+stepDots(3)+'</div><div class="form-group"><label class="form-label">'+t('HTTP Port (Web Arayuzu)')+'</label><input class="form-input" id="w-http-port" type="number" value="'+wizardData.http_port+'"></div><div class="form-group"><label class="form-label">'+t('HTTPS Port (SSL aktifse)')+'</label><input class="form-input" id="w-https-port" type="number" value="'+wizardData.https_port+'"></div><div class="form-group"><label class="form-label">'+t('RTMP Port (OBS Yayin)')+'</label><input class="form-input" id="w-rtmp-port" type="number" value="'+wizardData.rtmp_port+'"></div><div class="form-group"><label class="form-label">'+t('Public Domain / IP')+'</label><input class="form-input" id="w-embed-domain" placeholder="Orn: stream.ornek.com veya 203.0.113.10" value="'+escHtml(wizardData.embed_domain||'')+'"></div><div style="background:var(--bg-primary);border-radius:var(--radius-sm);padding:14px;margin-bottom:20px"><div style="font-size:13px;color:var(--text-muted)">'+t('Bos birakirsaniz panelin acildigi host kullanilir. HTTP ve HTTPS public portlari kurulumdan sonra Kolay Ayarlar veya Alan Adi / Embed ekranindan degistirilebilir.')+'</div></div><div style="display:flex;gap:12px"><button class="btn btn-secondary" style="flex:1" onclick="wizardPrev()"><i class="bi bi-arrow-left"></i> '+t('Geri')+'</button><button class="btn btn-primary" style="flex:1" onclick="wizardFinish()">'+t('Kurulumu Tamamla')+'</button></div></div>'
  ];
  return steps[wizardStep-1]||steps[0];
}
function stepDots(c){let d='';for(let i=1;i<=3;i++){d+='<span class="wizard-dot'+(i===c?' active':i<c?' done':'')+'"></span>'}return d}
function wizardNext(){
  if(wizardStep===2){
    const pw=document.getElementById('w-password').value;
    const pw2=document.getElementById('w-password2').value;
    const user=document.getElementById('w-username').value;
    if(!pw||pw.length<4){toast(t('Sifre en az 4 karakter olmali'),'error');return}
    if(pw!==pw2){toast(t('Sifreler eslesiyor!'),'error');return}
    wizardData.username=user||'admin';wizardData.password=pw;
  }
  wizardStep++;renderWizard();
}
function wizardPrev(){wizardStep--;renderWizard()}
async function wizardFinish(){
  wizardData.http_port=parseInt(document.getElementById('w-http-port').value)||8844;
  wizardData.https_port=parseInt(document.getElementById('w-https-port').value)||443;
  wizardData.rtmp_port=parseInt(document.getElementById('w-rtmp-port').value)||1935;
  wizardData.embed_domain=(document.getElementById('w-embed-domain').value||'').trim();
  const res=await api('/api/setup/complete',{method:'POST',body:wizardData});
  if(res.success){
    setupCompleted=true;toast(t('Kurulum tamamlandi!'));
    // Auto-login after setup
    const lr=await api('/api/auth/login',{method:'POST',body:{username:wizardData.username,password:wizardData.password}});
    if(lr.success){authToken=lr.token;sessionStorage.setItem('fluxstream_token',authToken)}
    setTimeout(()=>renderApp(),500);
  }
  else{toast(res.message||t('Kurulum hatasi'),'error')}
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â MAIN APP ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
function renderApp(){
  document.getElementById('app').innerHTML=
  '<div class="app">'+
    '<nav class="sidebar" id="sidebar">'+
      '<div class="logo"><div class="logo-icon"><i class="bi bi-lightning-charge-fill"></i></div><div><div class="logo-text">FluxStream</div><div class="logo-version">v2.0.0</div></div></div>'+
      '<div class="nav">'+
        '<div class="nav-section"><div class="nav-section-title">'+t('Ana Menu')+'</div>'+
          navItem('dashboard','bi-bar-chart-fill',t('Dashboard','Dashboard'))+
        '</div>'+
        '<div class="nav-section"><div class="nav-section-title">'+t('Yayin')+'</div>'+
          navItem('streams','bi-collection-play-fill',t('Yayinlar'))+
          navItem('create-stream','bi-plus-circle-fill',t('Yeni Yayin'))+
          navItem('embed-codes','bi-code-slash',t('Embed Kodlari'))+
          navItem('embed-advanced','bi-sliders',t('Gelismis Embed'))+
          navItem('player-templates','bi-pc-display',t('Player Sablonlari'))+
        '</div>'+
        '<div class="nav-section"><div class="nav-section-title">'+t('Ayarlar')+'</div>'+
          navItem('guided-settings','bi-magic',t('Kolay Ayarlar'))+
          navItem('settings-general','bi-gear-fill',t('Genel'))+
          navItem('logos','bi-images',t('Logo ve Marka'))+
          navItem('settings-embed','bi-globe2',t('Alan Adi / Embed'))+
          navItem('settings-protocols','bi-diagram-3-fill',t('Protokoller'))+
          navItem('settings-outputs','bi-boxes',t('Cikis Formatlari'))+
          navItem('settings-abr','bi-badge-hd',t('Teslimat / ABR'))+
          navItem('settings-ssl','bi-shield-lock-fill',t('SSL/TLS'))+
          navItem('settings-security','bi-shield-shaded',t('Guvenlik'))+
          navItem('settings-storage','bi-hdd-fill',t('Depolama ve Arsiv'))+
          navItem('settings-health','bi-heart-pulse-fill',t('Saglik ve Uyari'))+
          navItem('settings-transcode','bi-cpu-fill',t('Transkod'))+
        '</div>'+
        '<div class="nav-section"><div class="nav-section-title">'+t('Izleme')+'</div>'+
          navItem('operations-center','bi-broadcast-pin',t('Operasyon Merkezi'))+
          navItem('analytics','bi-graph-up',t('Analitik'))+
          navItem('viewers','bi-people-fill',t('Izleyiciler'))+
          navItem('transcode-jobs','bi-cpu',t('Transcode Isleri'))+
          navItem('diagnostics','bi-activity',t('Teshis'))+
        '</div>'+
        '<div class="nav-section"><div class="nav-section-title">'+t('Sistem')+'</div>'+
          navItem('maintenance-center','bi-safe2-fill',t('Bakim ve Yedek'))+
          navItem('license','bi-patch-check-fill',t('Lisans'))+
          navItem('security-tokens','bi-key-fill',t('Tokenlar'))+
          navItem('users','bi-person-fill',t('Kullanicilar'))+
          navItem('logs','bi-journal-text',t('Loglar'))+
        '</div>'+
      '</div>'+
    '</nav>'+
    '<div class="main">'+
      '<div class="topbar">'+
        '<div id="proto-status" class="proto-status"></div>'+
        '<div class="topbar-actions">'+
          '<button class="btn btn-secondary btn-sm" onclick="openSystemControl()"><i class="bi bi-power"></i> '+t('Sunucu Kontrol')+'</button>'+
          '<span style="font-size:13px;color:var(--text-muted)" id="clock"></span>'+
        '</div>'+
      '</div>'+
      '<div class="content" id="page-content"></div>'+
    '</div>'+
  '</div>';
  applyTranslations(document.getElementById('app'));
  navigate('dashboard');startClock();loadProtoStatus();
}

// navItem now uses Bootstrap Icons
function navItem(page,icon,label){
  return '<a class="nav-item'+(currentPage===page?' active':'')+'" onclick="navigate(\''+page+'\')" data-page="'+page+'"><span class="icon"><i class="bi '+icon+'"></i></span>'+label+'</a>';
}
function navigate(page){
  currentPage=page;
  document.querySelectorAll('.nav-item').forEach(el=>{
    el.classList.toggle('active',el.dataset.page===page);
  });
  if(pageRefreshTimer){
    clearTimeout(pageRefreshTimer);
    pageRefreshTimer=null;
  }
  if(streamTelemetryTimer){
    clearTimeout(streamTelemetryTimer);
    streamTelemetryTimer=null;
  }
  loadPage(page);
}
function startClock(){setInterval(()=>{const el=document.getElementById('clock');if(el)el.textContent=new Date().toLocaleTimeString(localeForLang())},1000)}
function schedulePageRefresh(page,ms){
  if(currentPage!==page)return;
  if(pageRefreshTimer)clearTimeout(pageRefreshTimer);
  pageRefreshTimer=setTimeout(()=>{if(currentPage===page)loadPage(page)},ms);
}

async function waitForServerBack(retries=30){
  for(let i=0;i<retries;i++){
    try{
      const res=await fetch(API+'/api/health',{cache:'no-store'});
      if(res.ok){location.reload();return}
    }catch(e){}
    await new Promise(r=>setTimeout(r,1000));
  }
  toast('Sunucu geri donmedi','error');
}

async function restartServer(){
  if(!confirm('Sunucu yeniden baslatilsin mi?'))return;
  const res=await api('/api/system/restart',{method:'POST'});
  if(res&&res.success){
    toast('Sunucu yeniden baslatiliyor...');
    setTimeout(()=>waitForServerBack(),1500);
    return;
  }
  toast((res&&res.message)||'Yeniden baslatma baslatilamadi','error');
}

async function stopServer(){
  if(!confirm('Sunucu durdurulsun mu?'))return;
  const res=await api('/api/system/stop',{method:'POST'});
  if(res&&res.success){
    toast('Sunucu durduruluyor...');
    return;
  }
  toast((res&&res.message)||'Durdurma baslatilamadi','error');
}

function openSystemControl(){
  const html=
    '<div class="modal-overlay" id="system-control-modal">'+
      '<div class="modal">'+
        '<div class="modal-title">Sunucu Kontrol</div>'+
        '<div style="color:var(--text-secondary);font-size:13px;line-height:1.7;margin-bottom:18px">Bu ekran servis islemlerini basitlestirir. Yeniden baslat, ayar degisikliklerini uygulamak icin kullanilir. Durdur secenegi sunucuyu tamamen kapatir.</div>'+
        '<div class="card-grid card-grid-2">'+
          '<div class="card"><div class="card-title" style="margin-bottom:8px">Yeniden Baslat</div><div class="form-hint" style="margin-bottom:14px">Port, SSL, transcode ve cikis ayarlari degistiginde kullanin.</div><button class="btn btn-primary" style="width:100%" onclick="closeModal(\'system-control-modal\');restartServer()">Yeniden Baslat</button></div>'+
          '<div class="card"><div class="card-title" style="margin-bottom:8px">Durdur</div><div class="form-hint" style="margin-bottom:14px">Exe kapanir. Tekrar calistirmak icin masaustu simgesine veya service aracina donmeniz gerekir.</div><button class="btn btn-danger" style="width:100%" onclick="closeModal(\'system-control-modal\');stopServer()">Durdur</button></div>'+
        '</div>'+
        '<div style="display:flex;justify-content:flex-end;margin-top:16px"><button class="btn btn-secondary" onclick="closeModal(\'system-control-modal\')">Kapat</button></div>'+
      '</div>'+
    '</div>';
  document.body.insertAdjacentHTML('beforeend',html);
  applyTranslations(document.getElementById('system-control-modal'));
}

function closeModal(id){
  const el=document.getElementById(id);
  if(el&&el.parentNode)el.parentNode.removeChild(el);
}

function scrollToElementId(id){
  const el=document.getElementById(id);
  if(el&&typeof el.scrollIntoView==='function'){
    el.scrollIntoView({behavior:'smooth',block:'start'});
  }
}

async function openTextInspectModal(title,url){
  const modalId='text-inspect-modal';
  closeModal(modalId);
  const html=
    '<div class="modal-overlay" id="'+modalId+'" onclick="if(event.target===this)closeModal(\''+modalId+'\')">'+
      '<div class="modal" style="max-width:980px">'+
        '<div class="modal-title">'+escHtml(title||'Metin Onizleme')+'</div>'+
        '<div class="form-hint" style="margin-bottom:12px">'+escHtml(url||'-')+'</div>'+
        '<pre id="text-inspect-body" style="margin:0;white-space:pre-wrap;word-break:break-word;max-height:60vh;overflow:auto;background:var(--bg-primary);border:1px solid var(--border);border-radius:12px;padding:16px;font-size:12px;line-height:1.55">Yukleniyor...</pre>'+
        '<div style="display:flex;justify-content:flex-end;gap:10px;margin-top:16px"><a class="btn btn-secondary" href="'+escHtml(url||'#')+'" target="_blank" rel="noopener">Yeni Sekmede Ac</a><button class="btn btn-primary" onclick="closeModal(\''+modalId+'\')">Kapat</button></div>'+
      '</div>'+
    '</div>';
  document.body.insertAdjacentHTML('beforeend',html);
  const body=document.getElementById('text-inspect-body');
  if(!body)return;
  try{
    const resp=await fetch(url,{cache:'no-store'});
    const text=await resp.text();
    body.textContent=text||'(Bos yanit)';
  }catch(e){
    body.textContent='Icerik yuklenemedi: '+String((e&&e.message)||e||'Bilinmeyen hata');
  }
}

async function loadProtoStatus(){
  const s=await api('/api/settings');
  const el=document.getElementById('proto-status');
  if(!el)return;
  const protos=[
    {n:'RTMP',k:'rtmp_enabled'},{n:'SRT',k:'srt_enabled'},{n:'RTSP',k:'rtsp_enabled'},
    {n:'WebRTC',k:'webrtc_enabled'},{n:'HLS',k:'hls_enabled'},{n:'DASH',k:'dash_enabled'}
  ];
  el.innerHTML=protos.map(p=>'<div class="proto-dot '+(s[p.k]==='true'?'on':'off')+'">'+p.n+'</div>').join('');
}

async function loadPage(page){
  const c=document.getElementById('page-content');
  if(!c)return;
  if(page==='dashboard')await renderDashboard(c);
  else if(page==='streams')await renderStreams(c);
  else if(page==='create-stream')renderCreateStream(c);
  else if(page==='guided-settings')await renderGuidedSettings(c);
  else if(page==='embed-codes')await renderEmbedCodes(c);
  else if(page.startsWith('stream-detail-'))await renderStreamDetail(c,page.replace('stream-detail-',''));
  else if(page==='operations-center')await renderOperationsCenter(c);
  else if(page==='logos')await renderLogos(c);
  else if(page==='settings-general')await renderSettingsGeneral(c);
  else if(page==='settings-embed')await renderSettingsEmbed(c);
  else if(page==='settings-protocols')await renderSettingsProtocols(c);
  else if(page==='settings-outputs')await renderSettingsOutputs(c);
  else if(page==='settings-abr')await renderSettingsABR(c);
  else if(page==='settings-ssl')await renderSettingsSSL(c);
  else if(page==='settings-security')await renderSettingsSecurity(c);
  else if(page==='settings-storage')await renderSettingsStorage(c);
  else if(page==='settings-health')await renderSettingsHealth(c);
  else if(page==='settings-transcode')await renderSettingsTranscode(c);
  else if(page==='logs')await renderLogs(c);
  else if(page==='users')await renderUsers(c);
  else if(page==='player-templates')await renderPlayerTemplates(c);
  else if(page==='embed-advanced')await renderAdvancedEmbed(c);
  else if(page==='analytics')await renderAnalytics(c);
  else if(page==='recordings')await renderSettingsStorage(c);
  else if(page==='viewers')await renderViewers(c);
  else if(page==='maintenance-center')await renderMaintenanceCenter(c);
  else if(page==='license')await renderLicensePage(c);
  else if(page==='security-tokens')await renderSecurityTokens(c);
  else if(page==='transcode-jobs')await renderTranscodeJobs(c);
  else if(page==='diagnostics')await renderDiagnostics(c);
  else c.innerHTML='<div class="empty-state"><div class="icon"><i class="bi bi-cone-striped"></i></div><h3>Yakinda</h3></div>';
  applyTranslations(c);
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â DASHBOARD ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderDashboard(c){
  const [stats,streams,analytics,health]=await Promise.all([api('/api/stats'),api('/api/streams'),api('/api/analytics'),api('/api/health/report')]);
  const live=(streams||[]).filter(s=>s.status==='live');
  const fmtItems=Object.entries((analytics&&analytics.viewers_by_format)||{}).sort((a,b)=>b[1]-a[1]).map(([label,value])=>({label:label,value:value}));
  const topStreams=((analytics&&analytics.top_streams)||[]).slice(0,5);
  const alerts=Array.isArray(health&&health.alerts)?health.alerts:[];
  const healthStatus=String((health&&health.status)||'ok').toUpperCase();
  const healthTone=healthStatus==='OK'?'green':(healthStatus==='WARNING'?'orange':'red');
  const bwIn=stats.bandwidth_in||0;const bwOut=stats.bandwidth_out||0;
  const memUsed=stats.memory_used_mb||0;const memTotal=stats.memory_total_mb||0;
  const memPct=memTotal>0?Math.round((memUsed/memTotal)*100):0;
  c.innerHTML=
    '<div class="studio-page">'+
      '<div class="studio-hero">'+
        '<div style="display:flex;justify-content:space-between;align-items:flex-start;flex-wrap:wrap;gap:14px">'+
          '<div><h1 class="studio-hero-title">Komuta Merkezi</h1>'+
          '<p class="studio-hero-sub">Yayinlar, izleyiciler, sistem sagligi ve performans metrikleri tek bakista. Kartlara tiklayarak hizlica ilgili ekrana gecin.</p></div>'+
          '<div style="display:flex;gap:10px;flex-wrap:wrap">'+
            '<button class="btn btn-secondary btn-sm" onclick="navigate(\'analytics\')"><i class="bi bi-graph-up"></i> Analitik</button>'+
            '<button class="btn btn-secondary btn-sm" onclick="navigate(\'operations-center\')"><i class="bi bi-sliders"></i> Operasyon</button>'+
            '<button class="btn btn-primary btn-sm" onclick="navigate(\'create-stream\')"><i class="bi bi-plus-circle"></i> Yeni Yayin</button>'+
          '</div>'+
        '</div>'+
        '<div class="studio-pill-row" style="margin-top:14px">'+
          '<span class="studio-pill '+(healthTone==='green'?'active':'')+'"><i class="bi bi-heart-pulse-fill"></i> Sistem: '+healthStatus+'</span>'+
          '<span class="studio-pill '+(live.length>0?'active':'')+'"><i class="bi bi-broadcast"></i> '+fmtInt(live.length)+' Canli Yayin</span>'+
          '<span class="studio-pill"><i class="bi bi-people-fill"></i> '+fmtInt(stats.total_viewers||0)+' Izleyici</span>'+
          '<span class="studio-pill"><i class="bi bi-clock-history"></i> '+formatUptime(stats.uptime_seconds)+'</span>'+
        '</div>'+
      '</div>'+
      '<div class="studio-kpi-grid">'+
        '<div class="studio-kpi" onclick="navigate(\'streams\')" style="cursor:pointer"><div class="studio-kpi-label"><i class="bi bi-broadcast" style="margin-right:6px"></i>Aktif Yayin</div><div class="studio-kpi-value">'+fmtInt(stats.active_streams||0)+'</div><div class="studio-kpi-sub">'+fmtInt((streams||[]).length)+' toplam tanimli</div></div>'+
        '<div class="studio-kpi" onclick="navigate(\'viewers\')" style="cursor:pointer"><div class="studio-kpi-label"><i class="bi bi-people-fill" style="margin-right:6px"></i>Aktif Izleyici</div><div class="studio-kpi-value">'+fmtInt(stats.total_viewers||0)+'</div><div class="studio-kpi-sub">Tum formatlardaki anlik oturumlar</div></div>'+
        '<div class="studio-kpi" onclick="navigate(\'transcode-jobs\')" style="cursor:pointer"><div class="studio-kpi-label"><i class="bi bi-memory" style="margin-right:6px"></i>Bellek</div><div class="studio-kpi-value">'+memUsed+' <span style="font-size:16px;color:var(--text-muted)">MB</span></div><div class="studio-kpi-sub">'+memPct+'% kullanim · '+memTotal+' MB toplam</div></div>'+
        '<div class="studio-kpi"><div class="studio-kpi-label"><i class="bi bi-arrow-down-up" style="margin-right:6px"></i>Bant Genisligi</div><div class="studio-kpi-value" style="font-size:22px">'+formatBytes(bwIn)+'</div><div class="studio-kpi-sub">'+formatBytes(bwOut)+' cikis</div></div>'+
        '<div class="studio-kpi" onclick="navigate(\'settings-health\')" style="cursor:pointer"><div class="studio-kpi-label"><i class="bi bi-bell-fill" style="margin-right:6px"></i>Uyarilar</div><div class="studio-kpi-value" style="color:'+(alerts.length>0?'var(--danger)':'var(--success)')+'">'+fmtInt(alerts.length)+'</div><div class="studio-kpi-sub">'+(alerts.length?escHtml(alerts[0].title||'Aktif uyari var'):'Sorun yok')+'</div></div>'+
      '</div>'+
      '<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span><i class="bi bi-broadcast-pin" style="margin-right:8px;color:var(--accent)"></i>Aktif Yayinlar</span><button class="btn btn-sm btn-secondary" onclick="navigate(\'streams\')">Tumu</button></div>'+
          (live.length===0
            ?'<div class="studio-empty"><i class="bi bi-broadcast" style="font-size:32px;display:block;margin-bottom:10px;opacity:.4"></i>Aktif yayin yok<br><span style="font-size:12px">Yeni bir yayin olusturun ve OBS ile baglanin</span></div>'
            :'<div style="display:grid;gap:10px">'+live.slice(0,6).map(function(s){
              return '<div style="display:flex;align-items:center;gap:12px;padding:10px 14px;border:1px solid var(--border);border-radius:14px;cursor:pointer;transition:all .15s;background:linear-gradient(180deg,#fff,#f8fbff)" onclick="navigate(\'stream-detail-'+s.id+'\')">'+
                '<span class="badge badge-live">CANLI</span>'+
                '<div style="flex:1;min-width:0"><div style="font-weight:700;font-size:14px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis">'+escHtml(s.name)+'</div><div style="font-size:12px;color:var(--text-muted)">'+escHtml(s.input_codec||'-')+' · '+fmtInt(s.viewer_count||0)+' izleyici</div></div>'+
              '</div>';
            }).join('')+'</div>')+
        '</div>'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span><i class="bi bi-lightning-charge" style="margin-right:8px;color:var(--warning)"></i>Hizli Operasyon</span></div>'+
          '<div class="metric-list">'+
            '<div class="metric-row"><span>Toplam giris</span><strong>'+formatBytes(bwIn)+'</strong></div>'+
            '<div class="metric-row"><span>Toplam cikis</span><strong>'+formatBytes(bwOut)+'</strong></div>'+
            '<div class="metric-row"><span>Top format</span><strong>'+(fmtItems[0]?escHtml(fmtItems[0].label):'Yok')+'</strong></div>'+
            '<div class="metric-row"><span>Saglik</span><strong><span class="tag tag-'+(healthTone==='green'?'green':(healthTone==='orange'?'yellow':'red'))+'">'+healthStatus+'</span></strong></div>'+
          '</div>'+
          '<div style="display:flex;gap:8px;flex-wrap:wrap;margin-top:14px">'+
            '<button class="btn btn-sm btn-secondary" onclick="openSystemControl()"><i class="bi bi-gear"></i> Sunucu</button>'+
            '<button class="btn btn-sm btn-secondary" onclick="navigate(\'settings-health\')"><i class="bi bi-heart-pulse"></i> Saglik</button>'+
            '<button class="btn btn-sm btn-secondary" onclick="navigate(\'guided-settings\')"><i class="bi bi-magic"></i> Hizli Ayar</button>'+
            '<button class="btn btn-sm btn-secondary" onclick="navigate(\'maintenance-center\')"><i class="bi bi-tools"></i> Bakim</button>'+
          '</div>'+
        '</div>'+
      '</div>'+
      '<div style="display:grid;gap:18px;grid-template-columns:repeat(auto-fit,minmax(300px,1fr))">'+
        '<div class="studio-card"><div class="studio-section-title" style="margin-bottom:8px"><span>24 Saat Izleyici</span><button class="btn btn-sm btn-secondary" onclick="navigate(\'analytics\')">Detay</button></div>'+renderTimelineChart((analytics&&analytics.viewers_timeline)||[],'Henuz izleyici verisi yok',function(v){return String(v)})+'</div>'+
        '<div class="studio-card"><div class="studio-section-title" style="margin-bottom:8px"><span>Format Dagilimi</span><button class="btn btn-sm btn-secondary" onclick="navigate(\'embed-advanced\')">Embed</button></div>'+renderBarList(fmtItems,'Henuz format verisi yok',function(v){return String(v)})+'</div>'+
        '<div class="studio-card"><div class="studio-section-title" style="margin-bottom:8px"><span>Populer Yayinlar</span><button class="btn btn-sm btn-secondary" onclick="navigate(\'analytics\')">Analitik</button></div>'+
          (topStreams.length?topStreams.map(function(s){
            const label=escHtml(s.stream_name||shortKey(s.stream_key));
            const sid=findStreamIdByKey((streams||[]),s.stream_key);
            const action=sid?' onclick="navigate(\'stream-detail-'+sid+'\')" style="cursor:pointer"':'';
            return '<div class="metric-row"'+action+'><span>'+label+'</span><span class="badge">'+fmtInt(s.viewers||0)+' izleyici</span></div>';
          }).join(''):'<div style="color:var(--text-muted)">Henuz veri yok</div>')+
        '</div>'+
        '<div class="studio-card"><div class="studio-section-title" style="margin-bottom:8px"><span>Uyarilar</span><button class="btn btn-sm btn-secondary" onclick="navigate(\'settings-health\')">Yonet</button></div>'+
          (alerts.length?alerts.slice(0,5).map(function(item){
            var tone=item.level==='critical'?'critical':item.level==='warning'?'warning':'info';
            return '<div class="studio-alert '+tone+'" style="margin-bottom:6px"><div style="display:flex;justify-content:space-between;gap:8px"><strong style="font-size:13px">'+escHtml(item.title||item.code||'Uyari')+'</strong><span class="tag '+(item.level==='critical'?'tag-red':item.level==='warning'?'tag-yellow':'tag-blue')+'">'+escHtml(String(item.level||'info').toUpperCase())+'</span></div>'+(item.description?'<div style="font-size:12px;color:var(--text-muted);margin-top:4px">'+escHtml(item.description)+'</div>':'')+'</div>';
          }).join(''):'<div class="studio-empty" style="padding:18px"><i class="bi bi-check-circle" style="font-size:24px;color:var(--success);display:block;margin-bottom:6px"></i>Aktif uyari yok</div>')+
        '</div>'+
      '</div>'+
    '</div>';
  schedulePageRefresh('dashboard',5000);
}
function statCard(color,iconClass,value,label,route,subtext){
  const clickable=route?' clickable':'';
  const action=route?' onclick="navigate(\''+route+'\')"':'';
  return '<div class="stat-card '+color+clickable+'"'+action+'><div class="stat-icon"><i class="bi '+iconClass+'"></i></div><div class="stat-value">'+value+'</div><div class="stat-label">'+label+'</div>'+(subtext?'<div class="stat-subtext">'+subtext+'</div>':'')+'</div>';
}
function fmtInt(n){return Number(n||0).toLocaleString(localeForLang())}
function shortKey(value){value=String(value||'');return value.length>18?value.slice(0,8)+'...'+value.slice(-6):value}
function renderBarList(items,emptyText,formatter){
  const list=Array.isArray(items)?items:[];
  if(!list.length)return '<div style="color:var(--text-muted)">'+(emptyText||'Henuz veri yok')+'</div>';
  const max=Math.max.apply(null,list.map(function(item){return Number(item.value||0)}).concat([1]));
  return '<div class="bar-list">'+list.map(function(item){
    const value=Number(item.value||0);
    const width=Math.max(6,Math.round((value/max)*100));
    return '<div class="bar-item"><div>'+escHtml(item.label||'-')+'</div><div class="bar-track"><div class="bar-fill" style="width:'+width+'%"></div></div><div style="text-align:right;font-weight:600">'+escHtml(formatter?formatter(value):String(value))+'</div></div>';
  }).join('')+'</div>';
}
function renderTimelineChart(points,emptyText,formatter,options){
  options=options||{};
  const source=Array.isArray(points)?points:[];
  if(!source.length)return '<div style="color:var(--text-muted)">'+(emptyText||'Henuz veri yok')+'</div>';
  const maxPoints=parseInt(options.maxPoints||20,10)||20;
  const list=source.slice(-maxPoints);
  const max=Math.max.apply(null,list.map(function(p){return Number(p.value||0)}).concat([1]));
  const min=Math.min.apply(null,list.map(function(p){return Number(p.value||0)}).concat([0]));
  const lastIndex=list.length-1;
  let meta=options.meta||('Son '+list.length+' nokta gosteriliyor');
  const firstDate=list[0]&&list[0].timestamp?new Date(list[0].timestamp):null;
  const lastDate=list[lastIndex]&&list[lastIndex].timestamp?new Date(list[lastIndex].timestamp):null;
  if(!options.meta&&firstDate&&lastDate&&!Number.isNaN(firstDate.getTime())&&!Number.isNaN(lastDate.getTime())){
    meta+=' • '+firstDate.toLocaleTimeString(localeForLang(),{hour:'2-digit',minute:'2-digit'})+' - '+lastDate.toLocaleTimeString(localeForLang(),{hour:'2-digit',minute:'2-digit'});
  }
  const width=640;
  const height=118;
  const baseY=height-8;
  const step=list.length===1?width:(width/(list.length-1));
  const pointToCoord=function(point,index){
    const value=Number(point.value||0);
    const normalized=max<=0?0:(value/max);
    const x=Math.round(index*step*100)/100;
    const y=Math.round((baseY-(normalized*(height-26)))*100)/100;
    return {x:x,y:y,value:value};
  };
  const coords=list.map(pointToCoord);
  const linePath=coords.map(function(coord,index){
    return (index===0?'M':'L')+coord.x+' '+coord.y;
  }).join(' ');
  const areaPath=linePath+' L '+coords[lastIndex].x+' '+baseY+' L 0 '+baseY+' Z';
  const labelIndices=[];
  const labelSlots=Math.max(3,Number(options.labelSlots||4));
  const stepSize=Math.max(1,Math.floor((list.length-1)/Math.max(1,labelSlots-1)));
  for(let i=0;i<list.length;i+=stepSize)labelIndices.push(i);
  if(labelIndices[labelIndices.length-1]!==lastIndex)labelIndices.push(lastIndex);
  const axisLabels=labelIndices.map(function(index){
    const point=list[index];
    const date=point&&point.timestamp?new Date(point.timestamp):null;
    const label=date?(options.labelFormatter?options.labelFormatter(date,index,list.length):date.toLocaleTimeString(localeForLang(),{hour:'2-digit',minute:'2-digit'})):'';
    return '<span title="'+escHtml(label)+'">'+escHtml(label||' ')+'</span>';
  }).join('');
  const currentValue=Number(list[lastIndex].value||0);
  const currentText=formatter?formatter(currentValue):String(currentValue);
  const peakText=formatter?formatter(max):String(max);
  const minText=formatter?formatter(min):String(min);
  const showPoints=list.length<=14;
  return '<div class="timeline-meta">'+escHtml(meta)+'</div>'+
    '<div class="sparkline-shell">'+
      '<div class="sparkline-frame">'+
        '<svg class="sparkline-svg" viewBox="0 0 '+width+' '+height+'" preserveAspectRatio="none" aria-hidden="true">'+
          '<defs><linearGradient id="sparkline-fill" x1="0" y1="0" x2="0" y2="1"><stop offset="0%" stop-color="rgba(20,184,166,.30)"></stop><stop offset="100%" stop-color="rgba(20,184,166,0)"></stop></linearGradient></defs>'+
          '<g class="sparkline-grid"><line x1="0" y1="'+(height*0.2).toFixed(1)+'" x2="'+width+'" y2="'+(height*0.2).toFixed(1)+'"></line><line x1="0" y1="'+(height*0.5).toFixed(1)+'" x2="'+width+'" y2="'+(height*0.5).toFixed(1)+'"></line><line x1="0" y1="'+(height*0.8).toFixed(1)+'" x2="'+width+'" y2="'+(height*0.8).toFixed(1)+'"></line></g>'+
          '<path class="sparkline-area" d="'+areaPath+'"></path>'+
          '<path class="sparkline-line" d="'+linePath+'"></path>'+
          (showPoints?coords.map(function(coord,index){
            const point=list[index];
            const date=point.timestamp?new Date(point.timestamp):null;
            const label=date?(options.labelFormatter?options.labelFormatter(date,index,list.length):date.toLocaleTimeString(localeForLang(),{hour:'2-digit',minute:'2-digit'})):'';
            const valueText=formatter?formatter(coord.value):String(coord.value);
            const tooltip=(date&&!Number.isNaN(date.getTime())?date.toLocaleString(localeForLang())+' • ':'')+valueText;
            return '<circle class="sparkline-point" cx="'+coord.x+'" cy="'+coord.y+'" r="4"><title>'+escHtml(tooltip)+' - '+escHtml(label)+'</title></circle>';
          }).join(''):'')+
        '</svg>'+
        '<div class="sparkline-hitmap" style="grid-template-columns:repeat('+list.length+',1fr)">'+list.map(function(point,index){
          const date=point.timestamp?new Date(point.timestamp):null;
          const label=date?(options.labelFormatter?options.labelFormatter(date,index,list.length):date.toLocaleTimeString(localeForLang(),{hour:'2-digit',minute:'2-digit'})):'';
          const value=Number(point.value||0);
          const valueText=formatter?formatter(value):String(value);
          const tooltip=(label?label+' • ':'')+valueText;
          return '<span class="sparkline-hit" title="'+escHtml(tooltip)+'" data-tooltip="'+escHtml(tooltip)+'"></span>';
        }).join('')+'</div>'+
      '</div>'+
      '<div class="sparkline-footer"><div class="sparkline-axis" style="grid-template-columns:repeat('+labelIndices.length+',minmax(0,1fr))">'+axisLabels+'</div><div class="sparkline-summary"><div class="sparkline-chip"><strong>'+escHtml(currentText)+'</strong><span>Son</span></div><div class="sparkline-chip"><strong>'+escHtml(peakText)+'</strong><span>Tepe</span></div><div class="sparkline-chip"><strong>'+escHtml(minText)+'</strong><span>Minimum</span></div></div></div>'+
    '</div>';
}
function streamCard(s){
  return '<div class="card" style="padding:16px;cursor:pointer" onclick="navigate(\'stream-detail-'+s.id+'\')">'+
    '<div style="display:flex;gap:14px;align-items:center;margin-bottom:12px">'+
      '<div class="stream-thumb '+(s.status==='live'?'live':'')+'"><i class="bi bi-play-btn-fill"></i></div>'+
      '<div><div style="font-weight:600;margin-bottom:4px">'+escHtml(s.name)+'</div>'+
        '<span class="badge badge-'+s.status+'">'+(s.status==='live'?'CANLI':'Cevrimdisi')+'</span></div>'+
    '</div>'+
    (s.status==='live'?'<div style="font-size:13px;color:var(--text-muted);display:flex;align-items:center;gap:6px;flex-wrap:wrap"><span><i class="bi bi-eye"></i> '+(s.viewer_count||0)+' izleyici</span><span>'+escHtml(s.input_codec||'Bilinmiyor')+'</span><span>'+(s.input_width&&s.input_height?(s.input_width+'x'+s.input_height):'cozunurluk yok')+'</span></div>':'')+
  '</div>';
}
function findStreamIdByKey(streams,key){
  const match=(streams||[]).find(function(s){return s.stream_key===key;});
  return match?match.id:0;
}

function withQueryParam(url,key,value){
  try{
    const next=new URL(String(url||''),window.location.origin);
    next.searchParams.set(key,String(value));
    return next.toString();
  }catch(e){
    return String(url||'');
  }
}
function formatShortSeconds(value){
  const num=Number(value||0);
  return Number.isFinite(num)?num.toFixed(1)+' sn':'-';
}
function formatAgoSeconds(value){
  const sec=Math.max(0,parseInt(value||0,10)||0);
  if(sec<60)return sec+' sn once';
  if(sec<3600)return Math.round(sec/60)+' dk once';
  return Math.round(sec/3600)+' sa once';
}
function renderTelemetryPills(values,emptyText){
  const entries=Object.entries(values||{}).sort(function(a,b){return Number(b[1]||0)-Number(a[1]||0)});
  if(!entries.length)return '<div class="form-hint">'+(emptyText||'Veri yok')+'</div>';
  return '<div style="display:flex;gap:8px;flex-wrap:wrap">'+entries.map(function(entry){
    return '<span class="tag tag-blue">'+escHtml(String(entry[0]||'-'))+' '+fmtInt(entry[1]||0)+'</span>';
  }).join('')+'</div>';
}
function renderTelemetrySessionsTable(items){
  const sessions=Array.isArray(items)?items:[];
  if(!sessions.length)return '<div class="form-hint">Henuz aktif player telemetrisi gelmedi.</div>';
  return '<div style="overflow:auto"><table class="viewer-table"><thead><tr><th>Oturum</th><th>Sayfa</th><th>Kaynak</th><th>Kalite</th><th>Kalite Gecisi</th><th>Ses</th><th>Ses Gecisi</th><th>Buffer</th><th>Stall</th><th>Durum</th><th>Son Gorus</th></tr></thead><tbody>'+
    sessions.map(function(item){
      return '<tr>'+
        '<td>'+escHtml(shortKey(item.session_id||'-'))+'</td>'+
        '<td>'+escHtml(item.page||'-')+'</td>'+
        '<td>'+escHtml(item.active_source_kind||'-')+'</td>'+
        '<td>'+escHtml(item.quality||'-')+'</td>'+
        '<td>'+escHtml(String(item.quality_transitions||0))+'</td>'+
        '<td>'+escHtml(item.selected_audio_label||item.selected_audio_track||'-')+'</td>'+
        '<td>'+escHtml(String(item.audio_switches||0))+'</td>'+
        '<td>'+escHtml(formatShortSeconds(item.buffer_seconds))+'</td>'+
        '<td>'+escHtml(String(item.stall_count||0))+'</td>'+
        '<td>'+escHtml(item.reconnect||'-')+(item.offline?' / offline':'')+'</td>'+
        '<td>'+escHtml(formatAgoSeconds(item.last_seen_ago_sec))+'</td>'+
      '</tr>';
    }).join('')+
  '</tbody></table></div>';
}
function renderTelemetryTrendChart(history,key,color,label,formatter){
  const items=Array.isArray(history)?history:[];
  if(!items.length){
    return '<div class="card" style="padding:16px"><div class="card-title" style="margin-bottom:12px">'+escHtml(label)+'</div><div class="form-hint">Kalici zaman serisi henuz olusmadi.</div></div>';
  }
  const values=items.map(function(item){
    const raw=Number(item&&item[key]);
    return Number.isFinite(raw)?raw:0;
  });
  const latest=values.length?values[values.length-1]:0;
  const min=Math.min.apply(null,values);
  const max=Math.max.apply(null,values);
  const width=320;
  const height=84;
  const pad=10;
  const span=Math.max(1,max-min);
  const points=values.map(function(value,index){
    const x=pad+((width-(pad*2))*index/Math.max(1,values.length-1));
    const normalized=(value-min)/span;
    const y=(height-pad)-((height-(pad*2))*normalized);
    return x.toFixed(1)+','+y.toFixed(1);
  }).join(' ');
  return '<div class="card" style="padding:16px">'+
    '<div style="display:flex;align-items:center;justify-content:space-between;gap:12px;margin-bottom:10px">'+
      '<div class="card-title">'+escHtml(label)+'</div>'+
      '<strong style="color:'+escHtml(color)+'">'+escHtml((formatter||fmtInt)(latest))+'</strong>'+
    '</div>'+
    '<svg viewBox="0 0 '+width+' '+height+'" preserveAspectRatio="none" style="width:100%;height:84px;display:block;background:var(--bg-primary);border-radius:10px;border:1px solid var(--border-color)">'+
      '<polyline fill="none" stroke="'+escHtml(color)+'" stroke-width="3" points="'+points+'"></polyline>'+
    '</svg>'+
    '<div class="form-hint" style="margin-top:8px">Son '+fmtInt(items.length)+' kalici ornek</div>'+
  '</div>';
}
function trackBitrateLabel(value){
  const num=Number(value||0);
  if(!Number.isFinite(num)||num<=0)return '-';
  return Math.round(num/1000)+' kbps';
}
function renderAlertList(items,emptyText){
  const alerts=Array.isArray(items)?items:[];
  if(!alerts.length)return '<div class="form-hint">'+(emptyText||'Aktif uyari yok.')+'</div>';
  return '<div style="display:grid;gap:10px">'+alerts.map(function(alert){
    const tone=alert.level==='critical'?'tag-red':(alert.level==='warning'?'tag-yellow':'tag-blue');
    return '<div class="card" style="padding:14px;border:1px solid var(--border-color)">'+
      '<div style="display:flex;align-items:flex-start;justify-content:space-between;gap:12px;margin-bottom:8px">'+
        '<div class="card-title" style="font-size:14px">'+escHtml(alert.title||alert.code||'Uyari')+'</div>'+
        '<span class="tag '+tone+'">'+escHtml((alert.level||'info').toUpperCase())+'</span>'+
      '</div>'+
      '<div class="form-hint" style="line-height:1.7">'+escHtml(alert.description||'-')+'</div>'+
      (alert.action?'<div class="form-hint" style="margin-top:8px;color:var(--text-primary)"><strong>Oneri:</strong> '+escHtml(alert.action)+'</div>':'')+
    '</div>';
  }).join('')+'</div>';
}
function groupTrackHistory(items,kind){
  const groups={};
  (Array.isArray(items)?items:[]).forEach(function(sample){
    if(kind&&sample.kind!==kind)return;
    const id=String(Number(sample.track_id||0));
    if(!groups[id]){
      groups[id]={
        track_id:Number(sample.track_id||0),
        kind:sample.kind||kind||'video',
        display_label:sample.display_label||('Track '+id),
        items:[]
      };
    }
    groups[id].items.push(sample);
  });
  return Object.values(groups).map(function(group){
    group.items.sort(function(a,b){
      return new Date(a.created_at||0).getTime()-new Date(b.created_at||0).getTime();
    });
    return group;
  }).sort(function(a,b){return a.track_id-b.track_id;});
}
function renderTrackAnalyticsGroups(items,kind,color){
  const groups=groupTrackHistory(items,kind);
  if(!groups.length)return '<div class="form-hint">'+(kind==='audio'?'Audio':'Video')+' track analytics verisi henuz birikmedi.</div>';
  return '<div class="card-grid card-grid-2">'+groups.map(function(group){
    return renderTelemetryTrendChart(group.items,'bitrate',color,group.display_label+' bitrate',trackBitrateLabel);
  }).join('')+'</div>';
}
function trackMetaLabel(track){
  if(!track)return '-';
  if(track.kind==='video'){
    if(track.width&&track.height)return track.width+'x'+track.height;
    if(track.height)return track.height+'p';
    return '-';
  }
  const parts=[];
  if(track.sample_rate)parts.push(track.sample_rate+' Hz');
  if(track.channels)parts.push(track.channels+' ch');
  return parts.length?parts.join(' / '):'-';
}
function renderTrackSelector(id,items,selectedID,emptyText,disabled){
  const tracks=Array.isArray(items)?items:[];
  return '<select class="form-select" id="'+id+'"'+(disabled?' disabled':'')+'>'+
    '<option value="0">Otomatik</option>'+
    tracks.map(function(track){
      const selected=Number(selectedID||0)===Number(track.track_id||0)?' selected':'';
      return '<option value="'+String(Number(track.track_id||0))+'"'+selected+'>'+escHtml(track.display_label||('Track '+track.track_id))+'</option>';
    }).join('')+
  '</select>'+(tracks.length?'':'<div class="form-hint" style="margin-top:8px">'+(emptyText||'Canli track verisi bekleniyor')+'</div>');
}
function renderTrackTable(items){
  const tracks=Array.isArray(items)?items:[];
  if(!tracks.length)return '<div class="form-hint">Track bilgisi henuz gelmedi.</div>';
  return '<div style="overflow:auto"><table class="viewer-table"><thead><tr><th>Track</th><th>Codec</th><th>Meta</th><th>Bitrate</th><th>Durum</th><th>Son Gorus</th></tr></thead><tbody>'+
    tracks.map(function(track){
      const tags=[];
      if(track.is_default)tags.push('<span class="tag tag-blue">Varsayilan</span>');
      if(track.is_active)tags.push('<span class="tag tag-green">Aktif</span>');
      if(track.enhanced)tags.push('<span class="tag tag-yellow">Enhanced</span>');
      return '<tr>'+
        '<td><strong>'+escHtml(track.display_label||('Track '+track.track_id))+'</strong><div class="form-hint">ID '+fmtInt(track.track_id||0)+'</div></td>'+
        '<td>'+escHtml(track.codec||'-')+'</td>'+
        '<td>'+escHtml(trackMetaLabel(track))+'</td>'+
        '<td>'+escHtml(trackBitrateLabel(track.bitrate))+'</td>'+
        '<td>'+(tags.join(' ')||'<span class="form-hint">-</span>')+'</td>'+
        '<td>'+escHtml(formatAgoSeconds(track.last_seen_ago_sec))+'</td>'+
      '</tr>';
    }).join('')+
  '</tbody></table></div>';
}
function renderTrackRuntimeBody(payload,policy,options){
  const opts=options||{};
  const tracks=payload&&payload.tracks?payload.tracks:{};
  const trackHistory=Array.isArray(payload&&payload.track_history)?payload.track_history:[];
  const videoTracks=Array.isArray(tracks.video_tracks)?tracks.video_tracks:[];
  const audioTracks=Array.isArray(tracks.audio_tracks)?tracks.audio_tracks:[];
  const defaultVideoID=Number((tracks.default_video_track_id!=null?tracks.default_video_track_id:policy&&policy.default_video_track_id)||0);
  const defaultAudioID=Number((tracks.default_audio_track_id!=null?tracks.default_audio_track_id:policy&&policy.default_audio_track_id)||0);
  const directMode=!!tracks.direct_mode;
  const readOnly=!!opts.readOnly;
  const footerHint=readOnly
    ?'Varsayilan secimleri kalici degistirmek icin yayin detay ekranindaki politika kartini kullan.'
    :'Video secimi yeni publish ile kokten etkili olur. Audio secimi mevcut canli oturuma da uygulanabilir.';
  return '<div class="card-grid card-grid-2" style="margin-bottom:16px">'+
      '<div class="card" style="padding:16px">'+
        '<div class="card-title" style="margin-bottom:12px">Varsayilan Track Secimi</div>'+
        '<div class="form-group"><label class="form-label">Varsayilan Video Track</label>'+renderTrackSelector('sd-default-video-track',videoTracks,defaultVideoID,'Video track secimi yeni yayinda tam olarak uygulanir.',readOnly)+'</div>'+
        '<div class="form-group"><label class="form-label">Varsayilan Audio Track</label>'+renderTrackSelector('sd-default-audio-track',audioTracks,defaultAudioID,'Audio track secimi canli yayinda uygulanabilir.',readOnly)+'</div>'+
        '<div class="form-hint">Durum: '+(directMode?'Direct multitrack HLS aktif':'Tek track veya klasik pipeline modu')+'</div>'+
        '<div class="form-hint" style="margin-top:8px">'+footerHint+'</div>'+
      '</div>'+
      '<div class="card" style="padding:16px">'+
        '<div class="card-title" style="margin-bottom:12px">Canli Runtime Ozet</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Aktif video track</span><strong>'+(tracks.active_video_track_id?fmtInt(tracks.active_video_track_id):'-')+'</strong></div>'+
          '<div class="metric-row"><span>Aktif audio track</span><strong>'+(tracks.active_audio_track_id?fmtInt(tracks.active_audio_track_id):'-')+'</strong></div>'+
          '<div class="metric-row"><span>Video track sayisi</span><strong>'+fmtInt(videoTracks.length)+'</strong></div>'+
          '<div class="metric-row"><span>Audio track sayisi</span><strong>'+fmtInt(audioTracks.length)+'</strong></div>'+
          '<div class="metric-row"><span>Son guncelleme</span><strong>'+escHtml(tracks.updated_at?fmtLocaleDateTime(tracks.updated_at):'-')+'</strong></div>'+
        '</div>'+
      '</div>'+
    '</div>'+
    '<div class="card-grid card-grid-2" style="margin-bottom:16px">'+
      '<div class="card" style="padding:16px"><div class="card-title" style="margin-bottom:12px">Video Trackleri</div>'+renderTrackTable(videoTracks)+'</div>'+
      '<div class="card" style="padding:16px"><div class="card-title" style="margin-bottom:12px">Audio Trackleri</div>'+renderTrackTable(audioTracks)+'</div>'+
    '</div>'+
    '<div class="card" style="padding:16px">'+
      '<div class="card-title" style="margin-bottom:12px">Track Analytics</div>'+
      '<div class="form-hint" style="margin-bottom:14px">Kalici bitrate ve track runtime ornekleri burada zaman serisi olarak birikir.</div>'+
      '<div class="form-hint" style="margin-bottom:8px">Video Trackleri</div>'+
      renderTrackAnalyticsGroups(trackHistory,'video','var(--accent)')+
      '<div class="form-hint" style="margin:16px 0 8px">Audio Trackleri</div>'+
      renderTrackAnalyticsGroups(trackHistory,'audio','var(--success)')+
    '</div>';
}
function renderStreamTelemetryBody(payload){
  const telemetry=payload&&payload.telemetry?payload.telemetry:{};
  const sessions=Array.isArray(telemetry.sessions)?telemetry.sessions:[];
  const history=Array.isArray(payload&&payload.history)?payload.history:[];
  const alerts=Array.isArray(payload&&payload.qoe_alerts)?payload.qoe_alerts:[];
  const lastUpdate=telemetry.last_update?fmtLocaleDateTime(telemetry.last_update):'-';
  const lastError=telemetry.last_error&&telemetry.last_error!=='-'?telemetry.last_error:'Yok';
  return (alerts.length?'<div class="card" style="padding:16px;margin-bottom:16px"><div class="card-title" style="margin-bottom:12px">QoE Uyarilari</div>'+renderAlertList(alerts,'Aktif QoE uyarisi yok.')+'</div>':'')+
    '<div class="card-grid card-grid-4" style="margin-bottom:16px">'+
      statCard('blue','bi-play-circle-fill',fmtInt(telemetry.active_sessions||0),'Aktif Player')+
      statCard('orange','bi-hourglass-split',fmtInt(telemetry.waiting_sessions||0),'Bekleyen')+
      statCard('red','bi-exclamation-triangle-fill',fmtInt(telemetry.total_stalls||0),'Toplam Stall')+
      statCard('green','bi-arrow-repeat',fmtInt(telemetry.total_recoveries||0),'Toparlanma')+
      statCard('purple','bi-shuffle',fmtInt(telemetry.total_quality_transitions||0),'Kalite Gecisi')+
      statCard('blue','bi-music-note-list',fmtInt(telemetry.total_audio_switches||0),'Ses Gecisi')+
    '</div>'+
    '<div class="card-grid card-grid-2" style="margin-bottom:16px">'+
      '<div class="card" style="padding:16px">'+
        '<div class="card-title" style="margin-bottom:12px">QoE Ozet</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Son guncelleme</span><strong>'+escHtml(lastUpdate)+'</strong></div>'+
          '<div class="metric-row"><span>Ortalama buffer</span><strong>'+escHtml(formatShortSeconds(telemetry.average_buffer_seconds))+'</strong></div>'+
          '<div class="metric-row"><span>Ortalama oynatma suresi</span><strong>'+escHtml(formatShortSeconds(telemetry.average_playback_seconds))+'</strong></div>'+
          '<div class="metric-row"><span>Offline oturum</span><strong>'+fmtInt(telemetry.offline_sessions||0)+'</strong></div>'+
          '<div class="metric-row"><span>Debug acik</span><strong>'+fmtInt(telemetry.debug_sessions||0)+'</strong></div>'+
        '<div class="metric-row"><span>Son hata</span><strong style="text-align:right">'+escHtml(lastError)+'</strong></div>'+
      '</div>'+
      '</div>'+
      '<div class="card" style="padding:16px">'+
        '<div class="card-title" style="margin-bottom:12px">Dagilim</div>'+
        '<div class="form-hint" style="margin-bottom:8px">Kaynak</div>'+renderTelemetryPills(telemetry.sources,'Kaynak verisi yok')+
        '<div class="form-hint" style="margin:14px 0 8px">Format</div>'+renderTelemetryPills(telemetry.formats,'Format verisi yok')+
        '<div class="form-hint" style="margin:14px 0 8px">Sayfa</div>'+renderTelemetryPills(telemetry.pages,'Sayfa verisi yok')+
        '<div class="form-hint" style="margin:14px 0 8px">Kalite</div>'+renderTelemetryPills(telemetry.qualities,'Kalite verisi yok')+
        '<div class="form-hint" style="margin:14px 0 8px">Ses Track</div>'+renderTelemetryPills(telemetry.audio_tracks,'Ses track verisi yok')+
      '</div>'+
    '</div>'+
    '<div class="card-grid card-grid-3" style="margin-bottom:16px">'+
      renderTelemetryTrendChart(history,'active_sessions','var(--accent)','Aktif Player Trendi',fmtInt)+
      renderTelemetryTrendChart(history,'average_buffer_seconds','var(--warning)','Buffer Trendi',formatShortSeconds)+
      renderTelemetryTrendChart(history,'total_stalls','var(--danger)','Stall Birikimi',fmtInt)+
      renderTelemetryTrendChart(history,'total_quality_transitions','var(--accent-hover)','Kalite Gecisi Trendi',fmtInt)+
      renderTelemetryTrendChart(history,'total_audio_switches','var(--success)','Ses Gecisi Trendi',fmtInt)+
    '</div>'+
    '<div class="card" style="padding:16px">'+
      '<div class="card-title" style="margin-bottom:12px">Son Aktif Oturumlar</div>'+
      renderTelemetrySessionsTable(sessions)+
    '</div>';
}
async function loadStreamTelemetry(id){
  const body=document.getElementById('stream-qoe-body');
  const trackBody=document.getElementById('stream-track-body');
  if(!body&&!trackBody)return;
  const currentVideoSelection=document.getElementById('sd-default-video-track')?.value||'';
  const currentAudioSelection=document.getElementById('sd-default-audio-track')?.value||'';
  const data=await api('/api/admin/player/telemetry/stream/'+id);
  if(!data||data.error){
    if(body)body.innerHTML='<div class="form-hint" style="color:var(--danger)">QoE telemetrisi alinamadi.</div>';
    if(trackBody)trackBody.innerHTML='<div class="form-hint" style="color:var(--danger)">Track runtime verisi alinamadi.</div>';
    return;
  }
  if(body)body.innerHTML=renderStreamTelemetryBody(data);
  if(trackBody){
    trackBody.innerHTML=renderTrackRuntimeBody(data,parseStreamPolicy(window._streamDetailData&&window._streamDetailData.policy_json));
    if(currentVideoSelection&&document.getElementById('sd-default-video-track'))document.getElementById('sd-default-video-track').value=currentVideoSelection;
    if(currentAudioSelection&&document.getElementById('sd-default-audio-track'))document.getElementById('sd-default-audio-track').value=currentAudioSelection;
  }
}
function startStreamTelemetryLoop(id){
  if(streamTelemetryTimer){
    clearTimeout(streamTelemetryTimer);
    streamTelemetryTimer=null;
  }
  const tick=async function(){
    if(currentPage!=='stream-detail-'+id)return;
    await loadStreamTelemetry(id);
    if(currentPage==='stream-detail-'+id)streamTelemetryTimer=setTimeout(tick,5000);
  };
  tick();
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â STREAMS LIST ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderStreams(c){
  const streams=await api('/api/streams')||[];
  const live=streams.filter(s=>s.status==='live');
  const offline=streams.filter(s=>s.status!=='live');
  const totalViewers=streams.reduce((a,s)=>a+(s.viewer_count||0),0);
  window._streamsFilter=window._streamsFilter||'all';
  const filter=window._streamsFilter;
  const filtered=filter==='live'?live:(filter==='offline'?offline:streams);
  c.innerHTML=
    '<div class="studio-page">'+
      '<div class="studio-hero">'+
        '<div style="display:flex;justify-content:space-between;align-items:flex-start;flex-wrap:wrap;gap:14px">'+
          '<div><h1 class="studio-hero-title">Yayin Yonetimi</h1>'+
          '<p class="studio-hero-sub">Tum yayinlarinizi buradan olusturun, izleyin ve yonetin. Canli yayinlara tiklayarak detay ekranina gecin.</p></div>'+
          '<div style="display:flex;gap:10px;flex-wrap:wrap">'+
            '<button class="btn btn-secondary btn-sm" onclick="navigate(\'operations-center\')"><i class="bi bi-sliders"></i> Operasyon</button>'+
            '<button class="btn btn-primary" onclick="navigate(\'create-stream\')"><i class="bi bi-plus-circle"></i> Yeni Yayin</button>'+
          '</div>'+
        '</div>'+
        '<div class="studio-pill-row" style="margin-top:14px">'+
          '<span class="studio-pill active"><i class="bi bi-collection"></i> '+fmtInt(streams.length)+' Toplam</span>'+
          '<span class="studio-pill '+(live.length>0?'active':'')+'"><i class="bi bi-broadcast"></i> '+fmtInt(live.length)+' Canli</span>'+
          '<span class="studio-pill"><i class="bi bi-people-fill"></i> '+fmtInt(totalViewers)+' Izleyici</span>'+
        '</div>'+
      '</div>'+
      '<div class="studio-toolbar">'+
        '<div class="studio-toolbar-group">'+
          '<div class="segmented">'+
            '<button class="segment '+(filter==='all'?'active':'')+'" onclick="window._streamsFilter=\'all\';navigate(\'streams\')">Tumu ('+streams.length+')</button>'+
            '<button class="segment '+(filter==='live'?'active':'')+'" onclick="window._streamsFilter=\'live\';navigate(\'streams\')">Canli ('+live.length+')</button>'+
            '<button class="segment '+(filter==='offline'?'active':'')+'" onclick="window._streamsFilter=\'offline\';navigate(\'streams\')">Cevrimdisi ('+offline.length+')</button>'+
          '</div>'+
        '</div>'+
        '<div class="studio-toolbar-group">'+
          '<button class="btn btn-sm btn-secondary" onclick="navigate(\'embed-codes\')"><i class="bi bi-code-slash"></i> Embed Kodlari</button>'+
        '</div>'+
      '</div>'+
      (filtered.length===0
        ?'<div class="studio-empty" style="padding:48px"><i class="bi bi-broadcast" style="font-size:42px;display:block;margin-bottom:14px;opacity:.4"></i><div style="font-size:16px;font-weight:700;margin-bottom:6px">'+(filter==='all'?'Henuz yayin yok':'Bu filtrede yayin yok')+'</div><div style="color:var(--text-muted);font-size:13px">'+(filter==='all'?'Ilk yayininizi olusturmak icin yukardaki butonu kullanin':'Filtreyi degistirerek diger yayinlari gorebilirsiniz')+'</div></div>'
        :'<div style="display:grid;gap:12px">'+filtered.map(function(s){
          const policy=parseStreamPolicy(s.policy_json);
          const isLive=s.status==='live';
          const borderColor=isLive?'rgba(239,68,68,.2)':'var(--border)';
          const bgGrad=isLive?'linear-gradient(180deg,#fff 0%,#fef8f8 100%)':'linear-gradient(180deg,#fff 0%,#f8fbff 100%)';
          return '<div style="display:flex;align-items:center;gap:16px;padding:16px 20px;border:1px solid '+borderColor+';border-radius:18px;cursor:pointer;transition:all .18s;background:'+bgGrad+'" onclick="navigate(\'stream-detail-'+s.id+'\')" onmouseover="this.style.transform=\'translateY(-2px)\';this.style.boxShadow=\'0 12px 28px rgba(37,99,235,.08)\'" onmouseout="this.style.transform=\'none\';this.style.boxShadow=\'none\'">'+
            '<div style="flex-shrink:0"><span class="badge badge-'+(isLive?'live':'offline')+'">'+(isLive?'CANLI':'Cevrimdisi')+'</span></div>'+
            '<div style="flex:1;min-width:0">'+
              '<div style="font-weight:700;font-size:15px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis">'+escHtml(s.name)+'</div>'+
              '<div style="font-size:12px;color:var(--text-muted);margin-top:2px"><code style="color:var(--accent)">'+escHtml(s.stream_key)+'</code></div>'+
              '<div style="display:flex;gap:8px;flex-wrap:wrap;margin-top:6px">'+
                '<span class="tag '+(policy.enable_abr?'tag-green':'tag-blue')+'">'+(policy.enable_abr?('Adaptive acik · '+escHtml(policy.profile_set||'balanced')):'Tek kalite teslimat')+'</span>'+
                (isLive?'<span class="tag tag-red">Canli yayin</span>':'<span class="tag tag-gray">Sonraki publish hazir</span>')+
              '</div>'+
            '</div>'+
            '<div style="display:flex;align-items:center;gap:16px;flex-shrink:0">'+
              (isLive?'<div style="text-align:center"><div style="font-size:18px;font-weight:800;color:var(--text-primary)">'+fmtInt(s.viewer_count||0)+'</div><div style="font-size:11px;color:var(--text-muted)">izleyici</div></div>':'')+
              '<div style="text-align:center;min-width:60px"><div style="font-size:12px;font-weight:600;color:var(--text-secondary)">'+escHtml(s.input_codec||'-')+'</div><div style="font-size:11px;color:var(--text-muted)">codec</div></div>'+
              '<div style="display:flex;gap:6px">'+
                '<button class="btn btn-sm '+(policy.enable_abr?'btn-secondary':'btn-primary')+'" onclick="event.stopPropagation();showAdaptiveModeModal('+s.id+')" title="Adaptive teslimat">'+(policy.enable_abr?'Adaptive':'Adaptiveye Al')+'</button>'+
                '<button class="btn btn-sm btn-secondary btn-icon" onclick="event.stopPropagation();navigate(\'operations-center\');setTimeout(function(){if(window.operationsCenterState)window.operationsCenterState.selectedStreamId=\''+s.id+'\'},100)" title="Operasyon"><i class="bi bi-sliders"></i></button>'+
                '<button class="btn btn-sm btn-danger btn-icon" onclick="event.stopPropagation();deleteStream('+s.id+')" title="Sil"><i class="bi bi-trash"></i></button>'+
              '</div>'+
            '</div>'+
          '</div>';
        }).join('')+'</div>')+
    '</div>';
  schedulePageRefresh('streams',5000);
}
async function deleteStream(id){
  if(!confirm('Bu yayini silmek istediginize emin misiniz?'))return;
  await api('/api/streams/'+id,{method:'DELETE'});toast('Yayin silindi');navigate('streams');
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â CREATE STREAM ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
function renderCreateStream(c){
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Yeni Yayin Olustur</h1></div>'+
    '<div class="card-grid card-grid-2">'+
      '<div class="card">'+
        '<div class="form-group"><label class="form-label">Yayin Adi *</label><input class="form-input" id="cs-name" placeholder="Orn: Konser Canli Yayin"></div>'+
        '<div class="form-group"><label class="form-label">Aciklama</label><input class="form-input" id="cs-desc" placeholder="Kisa aciklama"></div>'+
        '<div class="form-group"><label class="form-label">Yayin Modu</label><select class="form-select" id="cs-mode" onchange="updateCreateStreamGuide()"><option value="balanced">TV / Dengeli</option><option value="mobile">Mobil / Hafif</option><option value="resilient">Dusuk Bant / Dayanikli</option><option value="radio">Radyo / Audio</option></select><div class="form-hint">Bu secim ABR, cikis ve kaynak kullanimini belirleyen baslangic davranisini tanimlar.</div></div>'+
        '<div class="setting-row"><div><div class="setting-label">Adaptif Bitrate</div><div class="setting-desc">Acilirsa izleyicinin baglantisina gore kalite otomatik degisir.</div></div>'+
          '<label class="toggle"><input type="checkbox" id="cs-abr-enabled"><span class="toggle-slider"></span></label></div>'+
        '<div class="form-group" style="margin-top:16px"><label class="form-label">ABR Profil Seti</label><select class="form-select" id="cs-profile-set"><option value="balanced">Dengeli</option><option value="mobile">Mobil</option><option value="resilient">Dayanikli</option><option value="radio">Radyo</option></select></div>'+
        '<div class="setting-row"><div><div class="setting-label">Playback Token Gerekli</div><div class="setting-desc">Bu yayini izlemek icin token aranir.</div></div>'+
          '<label class="toggle"><input type="checkbox" id="cs-token-required"><span class="toggle-slider"></span></label></div>'+
        '<div class="form-group" style="margin-top:16px"><label class="form-label">Domain Kilidi</label><input class="form-input" id="cs-domain-lock" placeholder="Orn: mysite.com, embed.partner.com"><div class="form-hint">Bossa her yerde acilir. Doluysa sadece bu domainlerden gelen embed/referer kabul edilir.</div></div>'+
        '<div class="form-group"><label class="form-label">IP Beyaz Liste</label><input class="form-input" id="cs-ip-whitelist" placeholder="Orn: 203.0.113.20, 10.0.0.0/24"></div>'+
        '<div class="card-grid card-grid-2">'+
          '<div class="form-group"><label class="form-label">Maks Izleyici</label><input class="form-input" id="cs-max-viewers" type="number" value="0"><div class="form-hint">0 sinirsiz anlamina gelir.</div></div>'+
          '<div class="form-group"><label class="form-label">Maks Bitrate (kbps)</label><input class="form-input" id="cs-max-bitrate" type="number" value="0"><div class="form-hint">Kaynak kontrolu icin opsiyoneldir.</div></div>'+
        '</div>'+
        '<div class="form-group"><label class="form-label">Acik Cikis Formatlari</label><div class="form-hint" style="margin-bottom:10px">Bu yayinin disariya hangi formatlarda servis edilecegini secin.</div>'+renderOutputSelector(defaultStreamOutputs(),'cs')+'</div>'+
        '<div class="setting-row"><div><div class="setting-label">Yayin kaydedilsin mi?</div><div class="setting-desc">Varsayilan olarak kapali. Kalici kayitlar data/recordings altina yazilir.</div></div>'+
          '<label class="toggle"><input type="checkbox" id="cs-record-enabled" onchange="toggleCreateRecordFormat()"><span class="toggle-slider"></span></label></div>'+
        '<div class="form-group" style="margin-top:16px"><label class="form-label">Kayit Formati</label><select class="form-select" id="cs-record-format" disabled>'+recordingFormatOptions('mp4')+'</select><div class="form-hint">MP4 secildiginde yayin once guvenli bicimde yakalanir, yayin bitince izlenebilir dosyaya finalize edilir.</div></div>'+
        '<button class="btn btn-primary" onclick="createStream()">Yayin Olustur</button>'+
      '</div>'+
      '<div class="card" id="cs-side-panel">'+
        '<div class="card-title" style="margin-bottom:12px">Baglanti Rehberi</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>1. Kaynagi hazirla</span><strong>OBS / encoder</strong></div>'+
          '<div class="metric-row"><span>2. Yayin olustur</span><strong>Bu formdan</strong></div>'+
          '<div class="metric-row"><span>3. RTMP veya RTP ile baglan</span><strong>Bilgiler sag panelde kalir</strong></div>'+
          '<div class="metric-row"><span>4. Player ve embed kontrol et</span><strong>Olusur olusmaz kopyalanir</strong></div>'+
        '</div>'+
        '<div class="form-hint" style="margin-top:14px">Yayin olusturulduktan sonra OBS, RTP ve izleme baglanti bilgileri burada sabit kalir. Ekranin altina kaymaz.</div>'+
        '<div id="cs-guide" style="margin-top:18px"></div>'+
        '<div id="cs-result" style="margin-top:18px"></div>'+
      '</div>'+
    '</div>';
  updateCreateStreamGuide();
}
async function createStream(){
  const name=document.getElementById('cs-name').value;
  const desc=document.getElementById('cs-desc').value;
  const recordEnabled=document.getElementById('cs-record-enabled').checked;
  const recordFormat=document.getElementById('cs-record-format').value||'mp4';
  const outputFormats=collectOutputSelector('cs');
  const policy={
    mode:document.getElementById('cs-mode')?.value||'balanced',
    enable_abr:document.getElementById('cs-abr-enabled')?.checked||false,
    profile_set:document.getElementById('cs-profile-set')?.value||'balanced',
    require_playback_token:document.getElementById('cs-token-required')?.checked||false
  };
  if(!name){toast('Yayin adi gerekli','error');return}
  const res=await api('/api/streams',{method:'POST',body:{
    name,
    description:desc,
    output_formats:JSON.stringify(outputFormats),
    policy_json:JSON.stringify(policy),
    max_viewers:parseInt(document.getElementById('cs-max-viewers')?.value||'0')||0,
    max_bitrate:parseInt(document.getElementById('cs-max-bitrate')?.value||'0')||0,
    domain_lock:document.getElementById('cs-domain-lock')?.value||'',
    ip_whitelist:document.getElementById('cs-ip-whitelist')?.value||'',
    record_enabled:recordEnabled,
    record_format:recordFormat
  }});
  if(res.stream){
    toast('Yayin olusturuldu!');
    const settings=await api('/api/settings');
    const access=await getPlaybackAccess(res.stream.stream_key,settings,JSON.stringify(policy));
    const urls=getAllURLs(res.stream.stream_key,settings,name,access);
    updateCreateStreamGuide({mode:policy.mode,rtmp_url:res.rtmp_url,stream_key:res.stream.stream_key,stream_name:name});
    const r=document.getElementById('cs-result');
    r.innerHTML='<div class="card" style="padding:18px;background:var(--bg-primary)">'+
      '<div class="card-title" style="margin-bottom:16px">Yayin Hazir!</div>'+
      copyField('Stream Key',res.stream.stream_key)+
      copyField('OBS RTMP URL',res.rtmp_url)+
      copyField('RTP URL',urls.rtp)+
      copyField('HLS Izleme URL',urls.hls)+
      (access&&access.needs_token?'<div class="form-hint" style="margin-bottom:10px;color:var(--warning)">Bu yayinda playback token gerekli. Izleme linkine gecici token eklendi.</div>':'')+
      '<div style="margin-top:12px"><button class="btn btn-sm btn-primary" onclick="navigate(\'stream-detail-'+res.stream.id+'\')">Yayin Detaylarina Git <i class="bi bi-arrow-right"></i></button></div>'+
      '<div style="background:var(--bg-primary);border-radius:var(--radius-sm);padding:16px;margin-top:12px">'+
        '<div style="font-size:13px;color:var(--text-muted);line-height:1.6">OBS Studio\'da:<br>1. Ayarlar -> Yayin -> Hizmet: Ozel<br>2. Sunucu: <strong>'+escHtml(res.rtmp_url||'')+'</strong><br>3. Yayin Anahtari: <strong>'+res.stream.stream_key+'</strong><br>4. Yayina Baslat butonuna basin</div>'+
      '</div>'+
      '<div class="form-hint" style="margin-top:12px">Cok kanalli video kullanacaksaniz sagdaki rehberdeki JSON alanini tek tusla kopyalayabilirsiniz.</div>'+
      '<div style="background:var(--bg-primary);border-radius:var(--radius-sm);padding:16px;margin-top:12px">'+
        '<div style="font-size:13px;color:var(--text-muted);line-height:1.6">RTP push kullaniyorsaniz encoder hedefini <strong>'+escHtml(urls.rtp||'')+'</strong> olarak girebilirsiniz. MPEG-TS ve diger cikislar stream detay ekraninda hazirdir.</div>'+
      '</div></div>';
  }else{toast(res.message||'Hata','error')}
}
function toggleCreateRecordFormat(){
  const enabled=document.getElementById('cs-record-enabled')?.checked;
  const format=document.getElementById('cs-record-format');
  if(format)format.disabled=!enabled;
}
function getCreateStreamMode(){
  return document.getElementById('cs-mode')?.value||'balanced';
}
function getOBSMultitrackOverrideObject(mode){
  const presets={
    balanced:[
      {type:'obs_x264',width:1920,height:1080,framerate:{numerator:30,denominator:1},settings:{rate_control:'CBR',bitrate:6000,keyint_sec:2,preset:'veryfast',profile:'high',tune:'zerolatency'},canvas_index:0},
      {type:'obs_x264',width:854,height:480,framerate:{numerator:30,denominator:1},settings:{rate_control:'CBR',bitrate:1800,keyint_sec:2,preset:'veryfast',profile:'main',tune:'zerolatency'},canvas_index:0}
    ],
    mobile:[
      {type:'obs_x264',width:1280,height:720,framerate:{numerator:30,denominator:1},settings:{rate_control:'CBR',bitrate:2800,keyint_sec:2,preset:'veryfast',profile:'high',tune:'zerolatency'},canvas_index:0},
      {type:'obs_x264',width:640,height:360,framerate:{numerator:30,denominator:1},settings:{rate_control:'CBR',bitrate:900,keyint_sec:2,preset:'veryfast',profile:'main',tune:'zerolatency'},canvas_index:0}
    ],
    resilient:[
      {type:'obs_x264',width:960,height:540,framerate:{numerator:25,denominator:1},settings:{rate_control:'CBR',bitrate:1500,keyint_sec:2,preset:'veryfast',profile:'main',tune:'zerolatency'},canvas_index:0},
      {type:'obs_x264',width:640,height:360,framerate:{numerator:24,denominator:1},settings:{rate_control:'CBR',bitrate:650,keyint_sec:2,preset:'veryfast',profile:'baseline',tune:'zerolatency'},canvas_index:0},
      {type:'obs_x264',width:426,height:240,framerate:{numerator:20,denominator:1},settings:{rate_control:'CBR',bitrate:320,keyint_sec:2,preset:'veryfast',profile:'baseline',tune:'zerolatency'},canvas_index:0}
    ],
    radio:[
      {type:'obs_x264',width:1280,height:720,framerate:{numerator:25,denominator:1},settings:{rate_control:'CBR',bitrate:1800,keyint_sec:2,preset:'veryfast',profile:'main',tune:'zerolatency'},canvas_index:0},
      {type:'obs_x264',width:640,height:360,framerate:{numerator:25,denominator:1},settings:{rate_control:'CBR',bitrate:700,keyint_sec:2,preset:'veryfast',profile:'main',tune:'zerolatency'},canvas_index:0}
    ]
  };
  return {
    encoder_configurations:(presets[mode]||presets.balanced),
    audio_configurations:{
      live:[
        {codec:'ffmpeg_aac',track_id:1,channels:2,settings:{bitrate:160}}
      ]
    }
  };
}
function getOBSMultitrackOverrideJSON(mode){
  return JSON.stringify(getOBSMultitrackOverrideObject(mode),null,2);
}
function copyCodeField(label,value,rows){
  var raw=String(value==null?'':value);
  var id='copy_'+(++copyFieldSeq);
  copyValues[id]=raw;
  var minHeight=Math.max((rows||12)*18,120);
  return '<div class="form-group"><label class="form-label">'+label+'</label><div style="display:grid;gap:8px"><textarea class="form-textarea" readonly spellcheck="false" style="min-height:'+minHeight+'px;font-size:12px;line-height:1.55;font-family:Consolas,monospace;white-space:pre">'+escHtml(raw)+'</textarea><button type="button" class="copy-btn" style="justify-self:start" onclick="copyStoredValue(\''+id+'\')"><i class="bi bi-clipboard"></i> JSON Kopyala</button></div></div>';
}
function renderCreateStreamGuide(data){
  data=data||{};
  const settings=runtimeSettings||{};
  const mode=data.mode||getCreateStreamMode();
  const rtmpURL=String(data.rtmp_url||getOBSRTMPServerURL(settings));
  const streamKey=String(data.stream_key||'Yayin olusturunca burada gorunecek');
  const hasRealStream=!!data.stream_key;
  const json=getOBSMultitrackOverrideJSON(mode);
  const intro=hasRealStream
    ?'Bu yayin icin gerekli alanlar hazir. Asagidaki URL, stream key ve JSON\'u kopyalayip OBS\'e yapistirabilirsiniz.'
    :'Cok kanalli video normal RTMP gibi sadece URL ve key ile calismaz. Bu ozellikte Config Override JSON zorunludur. Once yayini olusturun, sonra bu alandaki bilgiler gercek degerlerle dolar.';
  return '<div class="card" style="padding:18px;background:var(--bg-primary)">'+
    '<div class="card-title" style="margin-bottom:12px">OBS Cok Kanalli Video Rehberi</div>'+
    '<div class="form-hint" style="margin-bottom:14px;line-height:1.7">'+intro+'</div>'+
    (hasRealStream?copyField('OBS RTMP URL',rtmpURL)+copyField('OBS Yayin Anahtari',streamKey):'')+
    copyCodeField('Config Override JSON',json,18)+
    '<div style="background:var(--bg-card);border:1px solid var(--border);border-radius:var(--radius-sm);padding:16px">'+
      '<div style="font-size:13px;font-weight:700;margin-bottom:10px">Adim adim kurulum rehberi:</div>'+
      '<ol style="margin:0;padding-left:18px;font-size:13px;line-height:1.85;color:var(--text-secondary)">'+
        '<li>OBS programini ac.</li>'+
        '<li>Alttaki bu sayfadan once yayini olustur. Yayin olusturunca burada <strong>OBS RTMP URL</strong> ve <strong>OBS Yayin Anahtari</strong> gorunur.</li>'+
        '<li>OBS icinde <strong>Ayarlar</strong> menusu ac.</li>'+
        '<li><strong>Yayin</strong> sekmesine gir.</li>'+
        '<li><strong>Service / Hizmet</strong> alanini <strong>Custom / Ozel</strong> yap.</li>'+
        '<li><strong>Server / Sunucu</strong> alanina burada gordugun <strong>OBS RTMP URL</strong> degerini yapistir.</li>'+
        '<li><strong>Stream Key / Yayin Anahtari</strong> alanina burada gordugun <strong>OBS Yayin Anahtari</strong> degerini yapistir.</li>'+
        '<li><strong>Enable Multitrack Video / Cok Kanalli Video</strong> secenegini ac.</li>'+
        '<li><strong>Enable Config Override</strong> secenegini ac.</li>'+
        '<li>Bu sayfadaki <strong>Config Override JSON</strong> alanini kopyala ve OBS icindeki kutuya komple yapistir.</li>'+
        '<li>OBS\'i tamamen kapatip tekrar ac. Bu adim onemli; bazi OBS surumlerinde yeniden acmadan multitrack baslamaz.</li>'+
        '<li>Son olarak <strong>Yayina Baslat</strong> dugmesine bas.</li>'+
      '</ol>'+
    '</div>'+
    '<div class="form-hint" style="margin-top:12px;line-height:1.7">Not: Normal baglanti calisip cok kanalli video calismiyorsa sebep genelde bu JSON\'un yapistirilmamasi veya OBS\'in yeniden acilmamasidir.</div>'+
  '</div>';
}
function updateCreateStreamGuide(data){
  const el=document.getElementById('cs-guide');
  if(!el)return;
  el.innerHTML=renderCreateStreamGuide(data||{mode:getCreateStreamMode()});
}
const copyValues={};
let copyFieldSeq=0;
const streamAccessCache={};
function copyField(label,value){
  var raw=String(value==null?'':value);
  var id='copy_'+(++copyFieldSeq);
  copyValues[id]=raw;
  return '<div class="form-group"><label class="form-label">'+label+'</label><div class="copy-group"><input class="copy-input" readonly value="'+escHtml(raw)+'"><button type="button" class="copy-btn" onclick="copyStoredValue(\''+id+'\')"><i class="bi bi-clipboard"></i></button></div></div>';
}
function copyStoredValue(id){copyText(copyValues[id]||'')}
function isTruthy(v){return v===true||v==='true'||v===1||v==='1'}
function getFallbackHost(){
  return (window.location&&window.location.hostname)||'localhost';
}
function getConfiguredDomain(s){
  var configured=String((s&&s.embed_domain)||'').trim();
  if(!configured||configured.toLowerCase()==='localhost'){
    return getFallbackHost();
  }
  return configured;
}
function hasConfiguredSSL(s){
  if(!isTruthy(s&&s.ssl_enabled))return false;
  const mode=String((s&&s.ssl_mode)||'file').toLowerCase();
  if(mode==='letsencrypt')return !!String((s&&s.ssl_le_domain)||'').trim();
  return !!(String((s&&s.ssl_cert_path)||'').trim()&&String((s&&s.ssl_key_path)||'').trim());
}
function shouldUsePublicHTTPS(s){
  return !!(isTruthy(s&&s.embed_use_https)&&hasConfiguredSSL(s));
}
function appendURLQuery(url,key,value){
  if(!value)return url;
  try{
    const next=new URL(url,window.location.origin);
    if(!next.searchParams.has(key))next.searchParams.set(key,value);
    return next.toString();
  }catch(e){
    return url+(url.indexOf('?')===-1?'?':'&')+encodeURIComponent(key)+'='+encodeURIComponent(value);
  }
}
function policyRequiresToken(raw){
  const policy=parseStreamPolicy(raw);
  return !!(policy.require_playback_token||policy.require_signed_url);
}
function cachedAccessValid(entry){
  if(!entry||!entry.token)return false;
  if(!entry.expires_at)return true;
  const expiry=new Date(entry.expires_at).getTime();
  return Number.isFinite(expiry)&&expiry>(Date.now()+15000);
}
async function getPlaybackAccess(streamKey,settings,policyRaw){
  const needsToken=isTruthy(settings&&settings.token_enabled)||policyRequiresToken(policyRaw);
  if(!streamKey||!needsToken)return {token:'',expires_at:'',needs_token:false};
  const cached=streamAccessCache[streamKey];
  if(cachedAccessValid(cached))return cached;
  const res=await api('/api/security/token/generate',{method:'POST',body:{stream_key:streamKey}});
  const access={token:(res&&res.token)||'',expires_at:(res&&res.expires_at)||'',needs_token:true};
  streamAccessCache[streamKey]=access;
  return access;
}
function getPublicBase(s){
  var domain=getConfiguredDomain(s);
  var useHTTPS=shouldUsePublicHTTPS(s);
  var scheme=useHTTPS?'https':'http';
  var port=useHTTPS?(s.embed_https_port||s.https_port||'443'):(s.embed_http_port||s.http_port||'8844');
  var defaultPort=useHTTPS?'443':'80';
  var portPart=port&&String(port)!==defaultPort?':'+port:'';
  return {domain:domain,useHTTPS:useHTTPS,scheme:scheme,port:String(port||''),base:scheme+'://'+domain+portPart};
}
function getOBSRTMPServerURL(s){
  return 'rtmp://'+getConfiguredDomain(s||{})+':'+(((s||{}).rtmp_port)||'1935')+'/live';
}
function recordingFormatOptions(selected){
  selected=selected||'mp4';
  return '<option value="mp4"'+(selected==='mp4'?' selected':'')+'>MP4 (.mp4) - Onerilen</option>'+
    '<option value="mkv"'+(selected==='mkv'?' selected':'')+'>Matroska (.mkv)</option>'+
    '<option value="ts"'+(selected==='ts'?' selected':'')+'>MPEG-TS (.ts) - Ham capture</option>'+
    '<option value="flv"'+(selected==='flv'?' selected':'')+'>FLV (.flv)</option>';
}
const streamOutputChoices=[
  ['hls','HLS'],['ll_hls','LL-HLS'],['dash','DASH'],['flv','HTTP-FLV'],['whep','WHEP'],
  ['mp4','MP4'],['webm','WebM'],['mp3','MP3'],['aac','AAC'],['ogg','OGG'],['wav','WAV'],['flac','FLAC'],['icecast','Icecast']
];
function defaultStreamOutputs(){
  return ['hls','ll_hls','dash','flv','whep','mp4','webm','mp3','aac','ogg','wav','flac','icecast'];
}
function parseJSONSafe(raw,fallback){
  try{return JSON.parse(raw||'')}catch(e){return fallback}
}
function parseStreamPolicy(raw){
  const policy=parseJSONSafe(raw,{})||{};
  if(!policy.profile_set)policy.profile_set='balanced';
  return policy;
}
function renderOutputSelector(selected,prefix){
  const active=Array.isArray(selected)&&selected.length?selected:defaultStreamOutputs();
  return '<div class="card-grid card-grid-3">'+streamOutputChoices.map(function(item){
    const key=item[0],label=item[1];
    const checked=active.indexOf(key)>=0;
    return '<label class="card" style="padding:14px;cursor:pointer"><div style="display:flex;align-items:center;justify-content:space-between;gap:10px"><div><div style="font-weight:600">'+label+'</div><div class="form-hint">'+key+'</div></div><input type="checkbox" class="'+prefix+'-output" value="'+key+'" '+(checked?'checked':'')+'></div></label>';
  }).join('')+'</div>';
}
function collectOutputSelector(prefix){
  return Array.from(document.querySelectorAll('.'+prefix+'-output:checked')).map(function(el){return el.value});
}

// URL HELPERS
function slugifyFileName(name){
  return String(name||'stream').toLowerCase().replace(/[^a-z0-9]+/g,'-').replace(/^-+|-+$/g,'')||'stream';
}
function namedMediaURL(basePath,key,name,ext){
  return basePath+'/'+key+'/'+slugifyFileName(name)+'.'+ext;
}
function buildURLSet(base,domain,key,s,name,access){
  var fileBase=slugifyFileName(name||key);
  var token=(access&&access.token)||'';
  function withToken(url){
    return token?appendURLQuery(url,'token',token):url;
  }
  return {
    hls:withToken(base+'/hls/'+key+'/master.m3u8'),
    hls_media:withToken(base+'/hls/'+key+'/index.m3u8'),
    ll_hls:withToken(base+'/hls/'+key+'/ll.m3u8'),
    dash:withToken(base+'/dash/'+key+'/manifest.mpd'),
    http_flv:withToken(base+'/flv/'+key),
    whep:withToken(base+'/whep/play/'+key),
    fmp4:withToken(base+'/mp4/'+key+'/'+fileBase+'.mp4'),
    webm:withToken(base+'/webm/'+key+'/'+fileBase+'.webm'),
    rtmp:'rtmp://'+domain+':'+(s.rtmp_port||'1935')+'/live/'+key,
    rtsp:'rtsp://'+domain+':'+(s.rtsp_port||'8554')+'/live/'+key,
    srt:'srt://'+domain+':'+(s.srt_port||'9000')+'?streamid='+key,
    rtp:'rtp://'+domain+':'+(s.rtp_port||'5004'),
    mpegts:'udp://'+domain+':'+(s.mpegts_port||'9001'),
    rtsp_out:'rtsp://'+domain+':'+(s.rtsp_out_port||'8555')+'/live/'+key,
    srt_out:'srt://'+domain+':'+(s.srt_out_port||'9010')+'?streamid='+key,
    mp3:withToken(base+'/audio/mp3/'+key+'/'+fileBase+'.mp3'),
    aac:withToken(base+'/audio/aac/'+key+'/'+fileBase+'.aac'),
    ogg:withToken(base+'/audio/ogg/'+key+'/'+fileBase+'.ogg'),
    wav:withToken(base+'/audio/wav/'+key+'/'+fileBase+'.wav'),
    flac:withToken(base+'/audio/flac/'+key+'/'+fileBase+'.flac'),
    hls_audio:withToken(base+'/audio/hls/'+key),
    dash_audio:withToken(base+'/audio/dash/'+key),
    icecast:withToken(base+'/icecast/'+key),
    asset_hls:base+'/static/vendor/hls.min.js',
    asset_dash:base+'/static/vendor/dash.all.min.js',
    asset_mpegts:base+'/static/vendor/mpegts.min.js',
    play:withToken(base+'/play/'+key),
    embed:withToken(base+'/embed/'+key)
  };
}
function getAllURLs(key,s,name,access){
  var publicConfig=getPublicBase(s||{});
  return buildURLSet(publicConfig.base,publicConfig.domain,key,s||{},name,access);
}
function getPreviewURLs(key,s,name,access){
  var base=(window.location&&window.location.origin)?window.location.origin:('http://'+getFallbackHost());
  return buildURLSet(base,getFallbackHost(),key,s||{},name,access);
}
function urlSection(title,pairs){
  return '<div class="card" style="margin-bottom:16px"><div class="card-title" style="margin-bottom:12px">'+title+'</div>'+
    pairs.map(function(p){return copyField(p[0],p[1])}).join('')+'</div>';
}

function renderDeliveryUsageCard(streamID,urls,options){
  const u=urls||{};
  const opts=options||{};
  const dashPlayer=withQueryParam(u.play||'','format','dash');
  const mp4Player=withQueryParam(u.play||'','format','mp4');
  let telemetryAction='<button class="btn btn-secondary btn-sm" onclick="scrollToElementId(\'stream-qoe-card\')">Telemetri Kartina Git</button>';
  if(opts.telemetryMode==='navigate'){
    telemetryAction='<button class="btn btn-secondary btn-sm" onclick="navigate(\'stream-detail-'+Number(streamID||0)+'\')">Telemetri Ekrani</button>';
  }else if(opts.telemetryMode==='operations'){
    telemetryAction='<button class="btn btn-secondary btn-sm" onclick="setOperationsCenterTab(\'qoe\')">Telemetri Sekmesini Ac</button>';
  }
  return '<div class="card" style="margin-bottom:16px">'+
    '<div class="card-header"><div class="card-title">Kullanim ve Tanilama Rehberi</div><div style="display:flex;gap:10px;flex-wrap:wrap">'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(u.play||'#')+'" target="_blank" rel="noopener">Tarayici Player</a>'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(dashPlayer||'#')+'" target="_blank" rel="noopener">DASH Player</a>'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(mp4Player||'#')+'" target="_blank" rel="noopener">MP4 Player</a>'+
      telemetryAction+
    '</div></div>'+
    '<div class="card-grid card-grid-2">'+
      '<div class="card" style="padding:16px">'+
        '<div class="card-title" style="margin-bottom:12px">Hangi link nerede kullanilir?</div>'+
        '<div class="form-hint" style="line-height:1.8;margin-bottom:12px">VLC icin en guvenli secim HLS URL\'dir. DASH MPD genelde teshis ve DASH uyumlu player icindir. Ham MP4 cikisi tarayicida dogrudan sekmede her zaman en iyi deneyimi vermez; tarayicida MP4 izlemek icin ustteki <strong>MP4 Player</strong> dugmesini kullan.</div>'+
        copyField('VLC icin onerilen HLS URL',u.hls||'')+
        copyField('DASH MPD URL',u.dash||'')+
        copyField('Ham MP4 URL',u.fmp4||'')+
      '</div>'+
      '<div class="card" style="padding:16px">'+
        '<div class="card-title" style="margin-bottom:12px">Manifest ve Telemetri</div>'+
        '<div class="form-hint" style="line-height:1.8;margin-bottom:12px">MPD veya HLS manifest dosyasini ham metin olarak acip kontrol edebilirsin. Canli player telemetrisi stream detay ekranindaki <strong>QoE ve Stall Telemetrisi</strong> kartinda gorunur.</div>'+
        '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-bottom:12px">'+
          '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("DASH MPD",'+JSON.stringify(u.dash||'')+')\'>MPD XML Goster</button>'+
          '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("HLS Master",'+JSON.stringify(u.hls||'')+')\'>HLS Master Goster</button>'+
          '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("HLS Media",'+JSON.stringify(u.hls_media||'')+')\'>HLS Media Goster</button>'+
        '</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Tarayicida izleme</span><strong>Player URL kullan</strong></div>'+
          '<div class="metric-row"><span>VLC / harici oynatici</span><strong>HLS URL kullan</strong></div>'+
          '<div class="metric-row"><span>Manifest kontrolu</span><strong>MPD XML / HLS Master</strong></div>'+
          '<div class="metric-row"><span>Canli telemetri</span><strong>QoE karti</strong></div>'+
        '</div>'+
      '</div>'+
    '</div>'+
  '</div>';
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â STREAM DETAIL Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
function setOperationsCenterFilter(filter){
  operationsCenterState.filter=String(filter||'all');
  const page=document.getElementById('page-content');
  if(page)renderOperationsCenter(page);
}
function setOperationsCenterSourceType(value){
  operationsCenterState.sourceType=String(value||'streams');
  const page=document.getElementById('page-content');
  if(page)renderOperationsCenter(page);
}
function selectOperationsStream(id){
  operationsCenterState.streamID=Number(id||0);
  const page=document.getElementById('page-content');
  if(page)renderOperationsCenter(page);
}
function setOperationsCenterStream(value){
  selectOperationsStream(value);
}
function setOperationsCenterTab(tab){
  operationsCenterState.tab=String(tab||'general');
  const page=document.getElementById('page-content');
  if(page)renderOperationsCenter(page);
}
function operationsCenterFilterMatches(stream,filter){
  const st=stream||{};
  switch(String(filter||'all')){
    case 'live':
      return st.status==='live';
    case 'offline':
      return st.status!=='live';
    case 'watched':
      return Number(st.viewer_count||0)>0;
    default:
      return true;
  }
}
function renderOperationsFilterButton(filter,label){
  const active=operationsCenterState.filter===filter;
  return '<button class="segment-btn '+(active?'active':'')+'" onclick="setOperationsCenterFilter(\''+filter+'\')">'+label+'</button>';
}
function renderOperationsSourceTypeSelect(){
  return '<select class="form-select" onchange="setOperationsCenterSourceType(this.value)">'+
    '<option value="streams"'+(operationsCenterState.sourceType==='streams'?' selected':'')+'>Canli Yayinlar / Streamler</option>'+
    '<option value="playlists" disabled>On-Demand Playlistler (yakinda)</option>'+
  '</select>';
}
function renderOperationsStreamSelect(streams){
  const items=Array.isArray(streams)?streams:[];
  if(!items.length){
    return '<select class="form-select" disabled><option>Gorunur stream yok</option></select>';
  }
  return '<select class="form-select" onchange="setOperationsCenterStream(this.value)">'+
    items.map(function(stream){
      const selected=Number(stream.id||0)===Number(operationsCenterState.streamID||0)?' selected':'';
      const label=(stream.name||'Yayin')+' ['+(stream.status==='live'?'CANLI':'Cevrimdisi')+'] - '+(stream.stream_key||'-');
      return '<option value="'+String(Number(stream.id||0))+'"'+selected+'>'+escHtml(label)+'</option>';
    }).join('')+
  '</select>';
}
function renderOperationsStreamListItem(stream,selected){
  const resolution=stream.input_width&&stream.input_height?(stream.input_width+'x'+stream.input_height):'cozunurluk yok';
  return '<button type="button" class="card" style="width:100%;text-align:left;padding:12px;border:'+(selected?'1px solid var(--accent)':'1px solid var(--border)')+';background:'+(selected?'rgba(59,130,246,.08)':'var(--bg-card)')+';cursor:pointer;box-shadow:none" onclick="selectOperationsStream('+Number(stream.id||0)+')">'+
    '<div style="display:flex;align-items:flex-start;justify-content:space-between;gap:12px;margin-bottom:8px">'+
      '<div><div style="font-weight:700;margin-bottom:4px">'+escHtml(stream.name||'Yayin')+'</div><div class="form-hint"><code>'+escHtml(shortKey(stream.stream_key||'-'))+'</code></div></div>'+
      '<span class="badge badge-'+escHtml(stream.status||'offline')+'">'+(stream.status==='live'?'CANLI':'Cevrimdisi')+'</span>'+
    '</div>'+
    '<div class="metric-list">'+
      '<div class="metric-row"><span>Izleyici</span><strong>'+fmtInt(stream.viewer_count||0)+'</strong></div>'+
      '<div class="metric-row"><span>Codec</span><strong>'+escHtml(stream.input_codec||'-')+'</strong></div>'+
      '<div class="metric-row"><span>Cozunurluk</span><strong>'+escHtml(resolution)+'</strong></div>'+
    '</div>'+
  '</button>';
}
function renderOperationsTabButton(tab,label){
  return '<button class="segment-btn '+(operationsCenterState.tab===tab?'active':'')+'" onclick="setOperationsCenterTab(\''+tab+'\')">'+label+'</button>';
}
function renderOperationsQuickActions(stream,urls,previewURLs){
  const u=urls||{};
  const preview=previewURLs||{};
  const playerDebug=withQueryParam(u.play||'','debug','1');
  const embedDebug=withQueryParam(u.embed||'','debug','1');
  const dashPlayer=withQueryParam(u.play||'','format','dash');
  const mp4Player=withQueryParam(u.play||'','format','mp4');
  const previewDebug=withQueryParam(preview.play||u.play||'','debug','1');
  return '<div class="card" style="margin-bottom:16px">'+
    '<div class="card-header"><div><div class="card-title">Hizli Eylemler</div><div class="form-hint">Secili yayin icin hizli linkler ve tanilama gecisleri.</div></div><div><button class="btn btn-secondary btn-sm" onclick="navigate(\'stream-detail-'+Number(stream.id||0)+'\')">Yayin Detayi</button></div></div>'+
    '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-bottom:14px">'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(u.play||'#')+'" target="_blank" rel="noopener">Player</a>'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(u.embed||'#')+'" target="_blank" rel="noopener">Embed</a>'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(dashPlayer||'#')+'" target="_blank" rel="noopener">DASH Player</a>'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(mp4Player||'#')+'" target="_blank" rel="noopener">MP4 Player</a>'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(playerDebug||'#')+'" target="_blank" rel="noopener">Debug Player</a>'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(embedDebug||'#')+'" target="_blank" rel="noopener">Debug Embed</a>'+
      '<a class="btn btn-secondary btn-sm" href="'+escHtml(previewDebug||'#')+'" target="_blank" rel="noopener">Canli Preview</a>'+
    '</div>'+
    '<div class="metric-list">'+
      '<div class="metric-row"><span>Player URL</span><strong class="mono-wrap">'+escHtml(u.play||'-')+'</strong></div>'+
      '<div class="metric-row"><span>Embed URL</span><strong class="mono-wrap">'+escHtml(u.embed||'-')+'</strong></div>'+
      '<div class="metric-row"><span>HLS Master</span><strong class="mono-wrap">'+escHtml(u.hls||'-')+'</strong></div>'+
      '<div class="metric-row"><span>DASH MPD</span><strong class="mono-wrap">'+escHtml(u.dash||'-')+'</strong></div>'+
    '</div>'+
  '</div>';
}
function renderOperationsDiagnosticsBody(data,urls){
  const diag=data||{};
  const checks=Array.isArray(diag.checks)?diag.checks:[];
  const telemetry=diag.telemetry||{};
  const hlsVariants=Number(diag.hls_variant_count||0);
  const dashRepresentations=Number(diag.dash_representation_count||0);
  const deliverySummary=diag.delivery_summary||{};
  const summaryTone='tag-'+(deliverySummary.tone||'yellow');
  return '<div class="card-grid card-grid-4" style="margin-bottom:16px">'+
      statCard('blue','bi-collection-play',fmtInt(hlsVariants),'HLS Varyant')+
      statCard('purple','bi-diagram-3',fmtInt(dashRepresentations),'DASH Representation')+
      statCard('orange','bi-people-fill',fmtInt(telemetry.active_sessions||0),'Aktif Player')+
      statCard('red','bi-exclamation-triangle-fill',fmtInt(telemetry.total_stalls||0),'Toplam Stall')+
    '</div>'+
    '<div class="card" style="margin-bottom:16px">'+
      '<div class="card-header"><div><div class="card-title">Teshis Ozeti</div><div class="form-hint">Manifest, cikis ve telemetry sagligi bu bolumde gorunur.</div></div><div><span class="tag '+summaryTone+'">'+escHtml(deliverySummary.label||'Durum bekleniyor')+'</span></div></div>'+
      '<div class="metric-list" style="margin-bottom:16px">'+
        '<div class="metric-row"><span>ABR Profil</span><strong>'+escHtml(diag.abr_profile_set||'balanced')+'</strong></div>'+
        '<div class="metric-row"><span>Player telemetrisi</span><strong>'+fmtInt(telemetry.active_sessions||0)+' aktif / '+fmtInt(telemetry.total_stalls||0)+' stall</strong></div>'+
        '<div class="metric-row"><span>Teslimat Ozeti</span><strong style="text-align:right">'+escHtml(deliverySummary.description||'-')+'</strong></div>'+
        '<div class="metric-row"><span>Policy JSON</span><strong class="mono-wrap">'+escHtml(diag.policy_json||'{}')+'</strong></div>'+
      '</div>'+
      '<div class="bar-list">'+checks.map(function(check){
        const tone='tag-'+(check.tone||'red');
        return '<div class="metric-row"><div><div>'+escHtml(check.description||check.code||'-')+'</div>'+(check.detail?'<div class="form-hint" style="margin-top:4px">'+escHtml(check.detail)+'</div>':'')+'</div><span class="tag '+tone+'">'+escHtml(check.label||'Sorunlu')+'</span></div>';
      }).join('')+'</div>'+
      '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-top:16px">'+
        '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("DASH MPD",'+JSON.stringify((urls&&urls.dash)||'')+')\'>MPD XML Goster</button>'+
        '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("DASH Ses MPD",'+JSON.stringify((urls&&urls.dash_audio)||'')+')\'>DASH Ses MPD</button>'+
        '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("HLS Master",'+JSON.stringify((urls&&urls.hls)||'')+')\'>HLS Master Goster</button>'+
        '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("HLS Media",'+JSON.stringify((urls&&urls.hls_media)||'')+')\'>HLS Media Goster</button>'+
        '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("Prometheus Metrics","/metrics")\'>Prometheus Metrics</button>'+
        '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("OpenTelemetry JSON","/api/observability/otel")\'>OpenTelemetry JSON</button>'+
      '</div>'+
    '</div>';
}
async function renderOperationsCenter(c){
  const streamsRes=await api('/api/streams');
  const streams=Array.isArray(streamsRes)?streamsRes:[];
  const filtered=streams.filter(function(stream){
    return operationsCenterFilterMatches(stream,operationsCenterState.filter);
  });
  if(filtered.length&&filtered.findIndex(function(stream){return Number(stream.id||0)===Number(operationsCenterState.streamID||0)})===-1){
    operationsCenterState.streamID=Number(filtered[0].id||0);
  }
  if(!filtered.length)operationsCenterState.streamID=0;
  let st=null;
  let settings={};
  let access={};
  let urls={};
  let previewURLs={};
  let telemetryData=null;
  let diagnosticsData=null;
  let policy={};
  if(operationsCenterState.streamID){
    const selectedID=Number(operationsCenterState.streamID||0);
    st=await api('/api/streams/'+selectedID);
    if(st&&!st.error){
      settings=await api('/api/settings')||{};
      access=await getPlaybackAccess(st.stream_key,settings,st.policy_json);
      urls=getAllURLs(st.stream_key,settings,st.name,access);
      previewURLs=getPreviewURLs(st.stream_key,settings,st.name,access);
      const diagAndTelemetry=await Promise.all([
        api('/api/admin/player/telemetry/stream/'+selectedID),
        api('/api/diagnostics/stream/'+selectedID)
      ]);
      telemetryData=diagAndTelemetry[0];
      diagnosticsData=diagAndTelemetry[1];
      policy=parseStreamPolicy(st.policy_json);
    }
  }
  let tabBody='<div class="card"><div class="empty-state"><div class="icon"><i class="bi bi-broadcast"></i></div><h3>Bir yayin secin</h3><p style="color:var(--text-muted)">Listedeki veya secim kutusundaki bir stream secerek operasyon verilerini gorebilirsiniz.</p></div></div>';
  if(operationsCenterState.sourceType!=='streams'){
    tabBody='<div class="card"><div class="empty-state"><div class="icon"><i class="bi bi-list-stars"></i></div><h3>On-demand playlist hazirligi tamamlandi</h3><p style="color:var(--text-muted)">Bu secim alani ileride on-demand playlistleri de ayni merkezden yonetmek icin kullanilacak. Bu fazda yalnizca canli streamler aktif durumda.</p></div></div>';
  }else if(st&&!st.error){
    const previewDebugURL=withQueryParam(previewURLs.play||urls.play||'','debug','1');
    const deliveryBody=
      renderDeliveryUsageCard(st.id,urls,{telemetryMode:'operations'})+
      urlSection('Video Akis URLleri',[
        ['HLS',urls.hls],['LL-HLS',urls.ll_hls],['DASH',urls.dash],['HTTP-FLV',urls.http_flv],['fMP4',urls.fmp4],['WebM',urls.webm]
      ])+
      urlSection('Ses ve Harici Oynatici',[
        ['AAC',urls.aac],['MP3',urls.mp3],['HLS Ses',urls.hls_audio],['DASH Ses',urls.dash_audio],['Icecast',urls.icecast]
      ]);
    const generalBody=
      '<div class="card-grid card-grid-4" style="margin-bottom:16px">'+
        statCard(st.status==='live'?'green':'red','bi-broadcast',st.status==='live'?'CANLI':'Cevrimdisi','Durum')+
        statCard('blue','bi-people-fill',fmtInt(st.viewer_count||0),'Izleyici')+
        statCard('purple','bi-badge-hd',st.input_width&&st.input_height?(st.input_width+'x'+st.input_height):'-','Cozunurluk')+
        statCard('orange','bi-camera-video-fill',st.input_codec||'-','Codec')+
      '</div>'+
      '<div class="card" style="margin-bottom:16px">'+
        '<div class="card-header"><div><div class="card-title">Canli Preview</div><div class="form-hint">Secili stream icin debug destekli canli player onizlemesi.</div></div><div><a class="btn btn-secondary btn-sm" href="'+escHtml(previewDebugURL||'#')+'" target="_blank" rel="noopener">Ayrica Ac</a></div></div>'+
        '<div style="position:relative;padding-top:56.25%;background:#000;border-radius:12px;overflow:hidden">'+
          '<iframe src="'+escHtml(previewDebugURL||'#')+'" style="position:absolute;top:0;left:0;width:100%;height:100%;border:none" allowfullscreen></iframe>'+
        '</div>'+
      '</div>'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Giris ve Yayin Ozeti</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Yayin adi</span><strong>'+escHtml(st.name||'-')+'</strong></div>'+
          '<div class="metric-row"><span>Stream key</span><strong class="mono-wrap">'+escHtml(st.stream_key||'-')+'</strong></div>'+
          '<div class="metric-row"><span>FPS</span><strong>'+escHtml(String(st.input_fps||'-'))+'</strong></div>'+
          '<div class="metric-row"><span>Bitrate</span><strong>'+(st.input_bitrate?formatBytes(st.input_bitrate)+'/s':'-')+'</strong></div>'+
          '<div class="metric-row"><span>OBS RTMP URL</span><strong class="mono-wrap">'+escHtml(getOBSRTMPServerURL(settings)||'-')+'</strong></div>'+
        '</div>'+
      '</div>';
    const trackBody=(telemetryData&&!telemetryData.error)
      ?renderTrackRuntimeBody(telemetryData,policy,{readOnly:true})+
        '<div class="card"><div class="form-hint">Varsayilan video ve audio track secimlerini kalici degistirmek icin yayin detay ekranindaki politika kartini kullanabilirsiniz.</div></div>'
      :'<div class="card"><div class="form-hint">Track runtime verisi henuz gelmedi.</div></div>';
    const manifestBody=
      '<div class="card" style="margin-bottom:16px">'+
        '<div class="card-header"><div><div class="card-title">Manifest ve Ham Veri</div><div class="form-hint">Ham manifest, harici oynatici URLleri ve dosya inceleme dugmeleri burada toplanir.</div></div></div>'+
        '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-bottom:16px">'+
          '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("DASH MPD",'+JSON.stringify(urls.dash||'')+')\'>MPD XML Goster</button>'+
          '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("HLS Master",'+JSON.stringify(urls.hls||'')+')\'>HLS Master Goster</button>'+
          '<button class="btn btn-secondary btn-sm" onclick=\'openTextInspectModal("HLS Media",'+JSON.stringify(urls.hls_media||'')+')\'>HLS Media Goster</button>'+
          '<a class="btn btn-secondary btn-sm" href="'+escHtml(urls.hls||'#')+'" target="_blank" rel="noopener">HLS Master Ac</a>'+
          '<a class="btn btn-secondary btn-sm" href="'+escHtml(urls.dash||'#')+'" target="_blank" rel="noopener">MPD Ac</a>'+
        '</div>'+
        copyField('HLS Master',urls.hls||'')+
        copyField('HLS Media',urls.hls_media||'')+
        copyField('DASH MPD',urls.dash||'')+
        copyField('Ham MP4',urls.fmp4||'')+
        copyField('Ham WebM',urls.webm||'')+
      '</div>';
    const obsBody=
      renderCreateStreamGuide({mode:policy.mode||'balanced',rtmp_url:getOBSRTMPServerURL(settings),stream_key:st.stream_key,stream_name:st.name})+
      '<div class="card"><div class="card-title" style="margin-bottom:12px">Ingest Ozeti</div><div class="metric-list">'+
        '<div class="metric-row"><span>Yayin modu</span><strong>'+escHtml(policy.mode||'balanced')+'</strong></div>'+
        '<div class="metric-row"><span>ABR</span><strong>'+(policy.enable_abr?'Acik':'Kapali')+'</strong></div>'+
        '<div class="metric-row"><span>Profil seti</span><strong>'+escHtml(policy.profile_set||'balanced')+'</strong></div>'+
        '<div class="metric-row"><span>Token gerekli</span><strong>'+(policy.require_playback_token?'Evet':'Hayir')+'</strong></div>'+
      '</div></div>';
    const qoeBody=(telemetryData&&!telemetryData.error)
      ?renderStreamTelemetryBody(telemetryData)
      :'<div class="card"><div class="form-hint">QoE telemetrisi henuz gelmedi.</div></div>';
    const diagnosisBody=(diagnosticsData&&!diagnosticsData.error)
      ?renderOperationsDiagnosticsBody(diagnosticsData,urls)
      :'<div class="card"><div class="form-hint">Teshis verisi henuz gelmedi.</div></div>';
    switch(operationsCenterState.tab){
      case 'delivery':
        tabBody=deliveryBody;
        break;
      case 'qoe':
        tabBody=qoeBody;
        break;
      case 'tracks':
        tabBody=trackBody+diagnosisBody;
        break;
      case 'manifests':
        tabBody=manifestBody;
        break;
      case 'obs':
        tabBody=obsBody;
        break;
      case 'diagnostics':
        tabBody=diagnosisBody;
        break;
      default:
        tabBody=generalBody;
        break;
    }
  }
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Canli Izleme ve Tanilama Merkezi</h1><div style="display:flex;gap:10px;flex-wrap:wrap"><button class="btn btn-secondary btn-sm" onclick="loadPage(\'operations-center\')"><i class="bi bi-arrow-clockwise"></i> Yenile</button><button class="btn btn-primary btn-sm" onclick="navigate(\'streams\')"><i class="bi bi-collection-play"></i> Yayinlar</button></div></div>'+
    '<div class="card" style="margin-bottom:16px">'+
      '<div class="card-header"><div><div class="card-title">Kaynak Secimi</div><div class="form-hint">Bu alan gelecekte canli streamlere ek olarak on-demand playlistleri de ayni merkezden secebilecegin bir yapi olacak.</div></div></div>'+
      '<div class="card-grid card-grid-2" style="margin-bottom:16px">'+
        '<div class="form-group" style="margin:0"><label class="form-label">Kaynak Turu</label>'+renderOperationsSourceTypeSelect()+'</div>'+
        '<div class="form-group" style="margin:0"><label class="form-label">Tum Streamler</label>'+renderOperationsStreamSelect(streams)+'</div>'+
      '</div>'+
      '<div class="segment-control">'+
          renderOperationsFilterButton('all','Tum')+
          renderOperationsFilterButton('live','Canli')+
          renderOperationsFilterButton('offline','Cevrimdisi')+
          renderOperationsFilterButton('watched','Izleyicili')+
      '</div>'+
    '</div>'+
    '<div style="display:grid;grid-template-columns:minmax(260px,300px) minmax(0,1fr);gap:16px;align-items:start">'+
      '<div class="card" style="position:sticky;top:16px">'+
        '<div class="card-header"><div><div class="card-title">Hizli Liste</div><div class="form-hint">Filtreye uyan tum streamler burada kalir. Selectbox ise her zaman tum streamleri listeler.</div></div></div>'+
        (filtered.length
          ?'<div style="display:grid;gap:10px;max-height:72vh;overflow:auto;padding-right:4px">'+filtered.map(function(stream){
            return renderOperationsStreamListItem(stream,Number(stream.id||0)===Number(operationsCenterState.streamID||0));
          }).join('')+'</div>'
          :'<div class="empty-state"><div class="icon"><i class="bi bi-search"></i></div><h3>Uygun yayin yok</h3><p style="color:var(--text-muted)">Secili filtrede gorunecek bir stream bulunmuyor.</p></div>')+
      '</div>'+
      '<div>'+
        (st&&!st.error?renderOperationsQuickActions(st,urls,previewURLs):'')+
        '<div class="card" style="margin-bottom:16px"><div class="segment-control" style="justify-content:flex-start;flex-wrap:wrap">'+
          renderOperationsTabButton('general','Genel Durum')+
          renderOperationsTabButton('delivery','Player ve Teslimat')+
          renderOperationsTabButton('qoe','QoE ve Telemetri')+
          renderOperationsTabButton('tracks','Track ve ABR')+
          renderOperationsTabButton('manifests','Manifest ve Ham Veri')+
          renderOperationsTabButton('obs','OBS ve Ingest')+
          renderOperationsTabButton('diagnostics','Teshis')+
        '</div></div>'+
        tabBody+
      '</div>'+
    '</div>';
  if(currentPage==='operations-center'){
    schedulePageRefresh('operations-center',8000);
  }
}

async function renderStreamDetail(c,id){
  const st=await api('/api/streams/'+id);
  if(!st||st.error){c.innerHTML='<div class="empty-state"><h3>Yayin bulunamadi</h3></div>';return}
  window._streamDetailData=st;
  const [settings,recsRes]=await Promise.all([api('/api/settings'),api('/api/recordings')]);
  const activeRecordings=(Array.isArray(recsRes)?recsRes:[]).filter(function(item){
    return item&&item.Status==='recording'&&item.StreamKey===st.stream_key;
  });
  const access=await getPlaybackAccess(st.stream_key,settings,st.policy_json);
  const u=getAllURLs(st.stream_key,settings,st.name,access);
  const previewURLs=getPreviewURLs(st.stream_key,settings,st.name,access);
  const policy=parseStreamPolicy(st.policy_json);
  const outputFormats=parseJSONSafe(st.output_formats,defaultStreamOutputs());
  const previewDebugURL=withQueryParam(previewURLs.play,'debug','1');
  const playerDebugURL=withQueryParam(u.play,'debug','1');
  const embedDebugURL=withQueryParam(u.embed,'debug','1');

  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">'+escHtml(st.name)+'</h1>'+
      '<div><span class="badge badge-'+st.status+'" style="margin-right:8px">'+(st.status==='live'?'CANLI':'Cevrimdisi')+'</span>'+
      '<button class="btn btn-sm btn-danger" onclick="deleteStream('+st.id+');navigate(\'streams\')">Sil</button></div></div>'+

    '<div class="card" style="margin-bottom:16px">'+
      '<div class="card-title" style="margin-bottom:12px">Yayin Bilgileri</div>'+
      copyField('Stream Key',st.stream_key)+
      '<div class="card-grid card-grid-2">'+
        '<div>'+
          '<div class="setting-row"><div class="setting-label">Durum</div><span class="badge badge-'+st.status+'">'+(st.status==='live'?'CANLI':'Cevrimdisi')+'</span></div>'+
          '<div class="setting-row"><div class="setting-label">Codec</div><div>'+(st.input_codec||'-')+'</div></div>'+
          '<div class="setting-row"><div class="setting-label">Cozunurluk</div><div>'+(st.input_width?st.input_width+'x'+st.input_height:'-')+'</div></div>'+
        '</div>'+
        '<div>'+
          '<div class="setting-row"><div class="setting-label">FPS</div><div>'+(st.input_fps||'-')+'</div></div>'+
          '<div class="setting-row"><div class="setting-label">Bitrate</div><div>'+(st.input_bitrate?formatBytes(st.input_bitrate)+'/s':'-')+'</div></div>'+
          '<div class="setting-row"><div class="setting-label">Izleyici</div><div>'+(st.viewer_count||0)+'</div></div>'+
        '</div>'+
      '</div>'+
    '</div>'+

    renderCreateStreamGuide({mode:policy.mode||'balanced',rtmp_url:getOBSRTMPServerURL(settings),stream_key:st.stream_key,stream_name:st.name})+

    urlSection('Video Akis URL\'leri',[
      ['HLS',u.hls],['LL-HLS',u.ll_hls],['DASH',u.dash],['HTTP-FLV',u.http_flv],
      ['WHEP (WebRTC)',u.whep],['fMP4',u.fmp4],['WebM',u.webm]
    ])+

    urlSection('Protokol URL\'leri',[
      ['RTMP',u.rtmp],['RTSP',u.rtsp],['SRT',u.srt],
      ['RTP',u.rtp],['MPEG-TS',u.mpegts],['RTSP Cikis',u.rtsp_out],['SRT Cikis',u.srt_out]
    ])+

    urlSection('Ses URL\'leri',[
      ['MP3',u.mp3],['AAC',u.aac],['OGG',u.ogg],
      ['WAV',u.wav],['FLAC',u.flac],['HLS Ses',u.hls_audio],['DASH Ses',u.dash_audio],['Icecast',u.icecast]
    ])+

    '<div class="card" style="margin-bottom:16px;border-color:'+(policy.enable_abr?'rgba(37,99,235,.18)':'rgba(59,130,246,.16)')+';background:'+(policy.enable_abr?'linear-gradient(180deg,#ffffff 0%,#f8fbff 100%)':'linear-gradient(180deg,#ffffff 0%,#f9fbff 100%)')+'">'+
      '<div class="card-header"><div><div class="card-title">Adaptive Teslimat</div><div class="form-hint">Kaynak yayin tek kalite olsa bile sunucu bunu coklu kalite HLS / DASH teslimatina cevirebilir.</div></div><span class="tag '+(policy.enable_abr?'tag-green':'tag-blue')+'">'+(policy.enable_abr?('Aktif · '+escHtml(policy.profile_set||'balanced')):'Kapali')+'</span></div>'+
      '<div class="metric-list">'+
        '<div class="metric-row"><span>Davranis</span><strong>'+(policy.enable_abr?'Izleyiciye baglanti hizina gore kalite sunulur':'Yayin su an tek kalite teslim ediliyor')+'</strong></div>'+
        '<div class="metric-row"><span>Onerilen profil</span><strong>'+escHtml(recommendedAdaptiveProfileSet(st))+'</strong></div>'+
        '<div class="metric-row"><span>Canli uygulama</span><strong>'+(st.status==='live'?'Mumkun · kisa yeniden kurulum etkisi olabilir':'Bir sonraki publishte risksiz uygulanir')+'</strong></div>'+
      '</div>'+
      '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-top:14px">'+
        '<button class="btn '+(policy.enable_abr?'btn-secondary':'btn-primary')+'" onclick="showAdaptiveModeModal('+st.id+')">'+(policy.enable_abr?'Adaptive Ayarini Guncelle':'Adaptiveye Al')+'</button>'+
        '<button class="btn btn-secondary" onclick="navigate(\'settings-abr\')">ABR Studyosunu Ac</button>'+
      '</div>'+
    '</div>'+

    '<div class="card" style="margin-bottom:16px"><div class="card-title" style="margin-bottom:12px">Teslimat ve Guvenlik Politikasi</div>'+
      '<div class="form-group"><label class="form-label">Yayin Modu</label><select class="form-select" id="sd-policy-mode"><option value="balanced" '+((policy.mode||'balanced')==='balanced'?'selected':'')+'>TV / Dengeli</option><option value="mobile" '+((policy.mode||'')==='mobile'?'selected':'')+'>Mobil / Hafif</option><option value="resilient" '+((policy.mode||'')==='resilient'?'selected':'')+'>Dusuk Bant / Dayanikli</option><option value="radio" '+((policy.mode||'')==='radio'?'selected':'')+'>Radyo / Audio</option></select><div class="form-hint">Bu, yayin icin secilen genel davranis profilidir.</div></div>'+
      '<div class="setting-row"><div><div class="setting-label">Adaptif Bitrate</div><div class="setting-desc">Acik oldugunda izleyiciye baglanti hizina gore farkli kalite katmanlari sunulur.</div></div>'+
      '<label class="toggle"><input type="checkbox" id="sd-abr-enabled" '+(policy.enable_abr?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      '<div class="form-group" style="margin-top:16px"><label class="form-label">ABR Profil Seti</label><select class="form-select" id="sd-profile-set"><option value="balanced" '+((policy.profile_set||'balanced')==='balanced'?'selected':'')+'>Dengeli</option><option value="mobile" '+((policy.profile_set||'')==='mobile'?'selected':'')+'>Mobil</option><option value="resilient" '+((policy.profile_set||'')==='resilient'?'selected':'')+'>Dayanikli</option><option value="radio" '+((policy.profile_set||'')==='radio'?'selected':'')+'>Radyo</option></select></div>'+
      '<div class="setting-row"><div><div class="setting-label">Playback Token Gerekli</div><div class="setting-desc">Acilirsa bu yayin icin token olmadan izleme baslamaz.</div></div>'+
      '<label class="toggle"><input type="checkbox" id="sd-token-required" '+(policy.require_playback_token?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      '<div class="form-group" style="margin-top:16px"><label class="form-label">Domain Kilidi</label><input class="form-input" id="sd-domain-lock" value="'+escHtml(st.domain_lock||'')+'" placeholder="mysite.com, partner.com"><div class="form-hint">Bos ise tum domainlerden acilir.</div></div>'+
      '<div class="form-group"><label class="form-label">IP Beyaz Liste</label><input class="form-input" id="sd-ip-whitelist" value="'+escHtml(st.ip_whitelist||'')+'" placeholder="203.0.113.20, 10.0.0.0/24"></div>'+
      '<div class="form-group"><label class="form-label">Maks Izleyici</label><input class="form-input" id="sd-max-viewers" type="number" value="'+escHtml(String(st.max_viewers||0))+'"></div>'+
      '<div class="form-group"><label class="form-label">Maks Bitrate (kbps)</label><input class="form-input" id="sd-max-bitrate" type="number" value="'+escHtml(String(st.max_bitrate||0))+'"></div>'+
      '<div class="form-group"><label class="form-label">Izinli Cikis Formatlari</label><div class="form-hint" style="margin-bottom:10px">Secilmeyen formatlar playback tarafinda reddedilir.</div>'+renderOutputSelector(outputFormats,'sd')+'</div>'+
      '<div style="margin-top:16px"><button class="btn btn-primary" onclick="saveStreamPolicySettings('+st.id+')">Politikayi Kaydet</button></div>'+
    '</div>'+

    '<div class="card" style="margin-bottom:16px"><div class="card-title" style="margin-bottom:12px">Kayit Politikasi</div>'+
      '<div class="setting-row"><div><div class="setting-label">Otomatik kayit</div><div class="setting-desc">Canli yayin basladiginda secili formatta kalici kayit baslatilir. HLS segmentleri kayit sayilmaz.</div></div>'+
      '<label class="toggle"><input type="checkbox" id="sd-record-enabled" '+(st.record_enabled?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      '<div class="form-group" style="margin-top:16px"><label class="form-label">Kayit Formati</label><select class="form-select" id="sd-record-format">'+recordingFormatOptions(st.record_format||'mp4')+'</select></div>'+
      '<div class="form-hint">Kalici kayitlar <code>data/recordings</code> altina yazilir. MP4 ve MKV secenekleri yayin kapaninca finalize edilir; canli cache dizinleri kayit sayilmaz.</div>'+
      '<div style="margin-top:16px"><button class="btn btn-primary" onclick="saveStreamRecordSettings('+st.id+')">Kayit Ayarlarini Kaydet</button></div>'+
    '</div>'+
    (activeRecordings.length?'<div class="card" style="margin-bottom:16px;border-color:rgba(239,68,68,.28);box-shadow:0 8px 22px rgba(239,68,68,.08)"><div class="card-header"><div><div class="card-title">Aktif Kayit</div><div class="form-hint">Bu yayin icin calisan kayit oturumlari burada gorunur.</div></div><span class="badge badge-live">'+fmtInt(activeRecordings.length)+' aktif</span></div><div style="display:flex;gap:10px;flex-wrap:wrap">'+activeRecordings.map(function(r){return '<div class="tag tag-red" style="display:flex;align-items:center;gap:10px;padding:8px 12px"><span>'+escHtml(String(r.Format||'').toUpperCase())+' · '+fmtBytes(r.Size||0)+'</span><button class="btn btn-sm btn-danger" onclick="stopRec(\''+r.ID+'\')">Kaydi Durdur</button></div>';}).join('')+'</div></div>':'')+

    '<div class="card" style="margin-bottom:16px"><div class="card-title" style="margin-bottom:12px">Embed Kodlari</div>'+
      (access&&access.needs_token?'<div class="form-hint" style="margin-bottom:10px;color:var(--warning)">Bu yayinda playback token gerekli. Aasagidaki preview ve linkler gecici token ile uretildi.</div>':'')+
      copyField('iframe','<iframe src="'+u.embed+'" width="1280" height="720" frameborder="0" allowfullscreen></iframe>')+
      copyField('Player URL',u.play)+
      copyField('Audio Player URL',playerURLForFormat(u.play,'aac'))+
      copyField('Embed URL',u.embed)+
    '</div>'+

    renderDeliveryUsageCard(st.id,u,{telemetryMode:'scroll'})+

    '<div class="card" id="stream-qoe-card" style="margin-bottom:16px"><div class="card-header"><div><div class="card-title">QoE ve Stall Telemetrisi</div><div class="form-hint">Canli player oturumlari, buffer, stall ve hata verileri burada gorunur.</div></div><div style="display:flex;gap:10px;flex-wrap:wrap"><a class="btn btn-secondary btn-sm" href="'+playerDebugURL+'" target="_blank" rel="noopener">Debug Player</a><a class="btn btn-secondary btn-sm" href="'+embedDebugURL+'" target="_blank" rel="noopener">Debug Embed</a></div></div><div id="stream-qoe-body" style="color:var(--text-muted)">QoE verisi bekleniyor...</div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><div class="card-title">Canli Track ve Varsayilan Secim</div><div class="form-hint">Multitrack video ve audio yapisi burada gorunur.</div></div><div id="stream-track-body" style="color:var(--text-muted)">Track verisi bekleniyor...</div></div>'+

    (st.status==='live'?
      '<div class="card"><div class="card-title" style="margin-bottom:12px">Onizleme (QoE Debug)</div>'+
        '<div style="position:relative;padding-top:56.25%;background:#000;border-radius:8px;overflow:hidden">'+
          '<iframe src="'+previewDebugURL+'" style="position:absolute;top:0;left:0;width:100%;height:100%;border:none" allowfullscreen></iframe>'+
        '</div></div>':'');

  startStreamTelemetryLoop(String(st.id));

}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â EMBED CODES Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function saveStreamRecordSettings(id){
  const st=window._streamDetailData;
  if(!st)return;
  const payload=Object.assign({},st,{
    record_enabled:document.getElementById('sd-record-enabled')?.checked||false,
    record_format:document.getElementById('sd-record-format')?.value||'mp4'
  });
  const res=await api('/api/streams/'+id,{method:'PUT',body:payload});
  if(res&&res.success){
    toast('Kayit ayarlari kaydedildi');
    navigate('stream-detail-'+id);
  }else{
    toast((res&&res.message)||'Kaydedilemedi','error');
  }
}
async function saveStreamPolicySettings(id){
  const st=window._streamDetailData;
  if(!st)return;
  const defaultVideoTrackID=parseInt(document.getElementById('sd-default-video-track')?.value||'0',10)||0;
  const defaultAudioTrackID=parseInt(document.getElementById('sd-default-audio-track')?.value||'0',10)||0;
  const policy={
    mode:document.getElementById('sd-policy-mode')?.value||'balanced',
    enable_abr:document.getElementById('sd-abr-enabled')?.checked||false,
    profile_set:document.getElementById('sd-profile-set')?.value||'balanced',
    require_playback_token:document.getElementById('sd-token-required')?.checked||false,
    default_video_track_id:defaultVideoTrackID,
    default_audio_track_id:defaultAudioTrackID
  };
  const payload=Object.assign({},st,{
    max_viewers:parseInt(document.getElementById('sd-max-viewers')?.value||'0')||0,
    max_bitrate:parseInt(document.getElementById('sd-max-bitrate')?.value||'0')||0,
    domain_lock:document.getElementById('sd-domain-lock')?.value||'',
    ip_whitelist:document.getElementById('sd-ip-whitelist')?.value||'',
    output_formats:JSON.stringify(collectOutputSelector('sd')),
    policy_json:JSON.stringify(policy)
  });
  const res=await api('/api/streams/'+id,{method:'PUT',body:payload});
  if(res&&res.success){
    await api('/api/admin/stream/tracks/defaults/'+id,{method:'POST',body:{
      default_video_track_id:defaultVideoTrackID,
      default_audio_track_id:defaultAudioTrackID
    }});
    toast('Politika kaydedildi');
    navigate('stream-detail-'+id);
  }else{
    toast((res&&res.message)||'Kaydedilemedi','error');
  }
}
function recommendedAdaptiveProfileSet(stream){
  const st=stream||{};
  const policy=parseStreamPolicy(st.policy_json);
  if(policy.profile_set)return policy.profile_set;
  if((policy.mode||'')==='radio')return 'radio';
  if((st.input_width||0)>0 && (st.input_width||0)<=960)return 'mobile';
  if((st.input_width||0)>=1920 || (st.input_bitrate||0)>=3500000)return 'balanced';
  return policy.mode||'balanced';
}
function adaptiveProfileSelectOptions(selected){
  const current=String(selected||'balanced');
  return [
    ['balanced','TV / Dengeli'],
    ['mobile','Mobil / Hafif'],
    ['resilient','Dusuk Bant / Dayanikli'],
    ['radio','Radyo / Audio']
  ].map(function(item){
    return '<option value="'+item[0]+'"'+(item[0]===current?' selected':'')+'>'+item[1]+'</option>';
  }).join('');
}
async function showAdaptiveModeModal(streamId){
  const st=await api('/api/streams/'+streamId);
  if(!st||st.error){toast('Yayin bilgisi alinamadi','error');return;}
  const policy=parseStreamPolicy(st.policy_json);
  const recommended=recommendedAdaptiveProfileSet(st);
  closeModal('adaptive-mode-modal');
  const isLive=st.status==='live';
  const currentProfile=policy.profile_set||recommended||'balanced';
  document.body.insertAdjacentHTML('beforeend',
    '<div class="modal-overlay" id="adaptive-mode-modal" onclick="if(event.target===this)closeModal(\'adaptive-mode-modal\')">'+
      '<div class="modal" style="max-width:760px">'+
        '<div class="modal-title">Adaptive teslimata gec</div>'+
        '<div class="form-hint" style="margin-bottom:16px">Kaynak yayin tek kalite olsa bile sunucu bunu coklu kalite HLS / DASH teslimatina cevirebilir. Varsayilan guvenli akis, ayari simdi kaydedip bir sonraki yayinda devreye almaktir.</div>'+
        '<div class="card-grid card-grid-2" style="margin-bottom:16px">'+
          '<div class="card" style="padding:16px"><div class="card-title" style="font-size:16px;margin-bottom:10px">Yayin ozeti</div><div class="metric-list">'+
            '<div class="metric-row"><span>Yayin</span><strong>'+escHtml(st.name||st.stream_key)+'</strong></div>'+
            '<div class="metric-row"><span>Durum</span><strong>'+(isLive?'Canli':'Cevrimdisi')+'</strong></div>'+
            '<div class="metric-row"><span>Mevcut teslimat</span><strong>'+(policy.enable_abr?('Adaptive · '+escHtml(currentProfile)):'Tek kalite')+'</strong></div>'+
            '<div class="metric-row"><span>Oneri</span><strong>'+escHtml(recommended)+'</strong></div>'+
          '</div></div>'+
          '<div class="card" style="padding:16px"><div class="card-title" style="font-size:16px;margin-bottom:10px">Ne olacak?</div><div class="metric-list">'+
            '<div class="metric-row"><span>Sonraki yayinda</span><strong>Kesintisiz ve onerilen</strong></div>'+
            '<div class="metric-row"><span>Canli uygula</span><strong>'+(isLive?'Kisa yeniden kurulum etkisi olabilir':'Yayin canli degilse otomatik olarak sonraki publishe kalir')+'</strong></div>'+
            '<div class="metric-row"><span>Kaynak tipi</span><strong>Tek bitrate kaynak da uygundur</strong></div>'+
            '<div class="metric-row"><span>Sunucu maliyeti</span><strong>FFmpeg ve CPU/GPU kullanimi artar</strong></div>'+
          '</div></div>'+
        '</div>'+
        '<div class="form-group"><label class="form-label">ABR profil seti</label><select class="form-select" id="adaptive-profile-set">'+adaptiveProfileSelectOptions(currentProfile)+'</select><div class="form-hint">TV / Dengeli genel amacli baslangictir. Dusuk uplink veya zayif sunucuda Dayanikli daha guvenlidir.</div></div>'+
        '<div class="form-group"><label class="form-label">Uygulama modu</label><select class="form-select" id="adaptive-apply-mode">'+
          '<option value="next_publish">Sonraki yayinda etkinlestir (onerilen)</option>'+
          (isLive?'<option value="live_now">Canli yayina simdi uygula</option>':'')+
        '</select><div class="form-hint">'+(isLive?'Canli uygula secenegi HLS/DASH zincirini kisa sureli yeniden kurar.':'Yayin su an canli olmadigi icin ayar risksiz sekilde bir sonraki publishte devreye girer.')+'</div></div>'+
        '<div class="setting-row"><div><div class="setting-label">Yayin modu ile hizala</div><div class="setting-desc">Secilen profil seti, yayin modunu da ayni aileye ceker.</div></div><label class="toggle"><input type="checkbox" id="adaptive-sync-mode" checked><span class="toggle-slider"></span></label></div>'+
        '<div style="display:flex;justify-content:flex-end;gap:10px;margin-top:18px"><button class="btn btn-secondary" onclick="closeModal(\'adaptive-mode-modal\')">Vazgec</button><button class="btn btn-primary" onclick="applyAdaptiveModeForStream('+streamId+')">Adaptive Olarak Isaretle</button></div>'+
      '</div>'+
    '</div>');
}
async function applyAdaptiveModeForStream(streamId){
  const profileSet=document.getElementById('adaptive-profile-set')?.value||'balanced';
  const applyMode=document.getElementById('adaptive-apply-mode')?.value||'next_publish';
  const syncMode=!!document.getElementById('adaptive-sync-mode')?.checked;
  const res=await api('/api/admin/streams/adaptive-mode',{method:'POST',body:{
    stream_id:streamId,
    profile_set:profileSet,
    apply_mode:applyMode,
    sync_mode:syncMode
  }});
  if(!res||res.error){
    toast((res&&res.message)||'Adaptive teslimat uygulanamadi','error');
    return;
  }
  if(Array.isArray(res.warnings)&&res.warnings.length){
    toast(res.warnings[0],'warning');
  }else{
    toast(res.message||'Adaptive teslimat guncellendi');
  }
  closeModal('adaptive-mode-modal');
  if(String(currentPage||'').startsWith('stream-detail-')){
    navigate('stream-detail-'+streamId);
  }else{
    navigate('streams');
  }
}
async function renderEmbedCodes(c){
  try{
    const streamsRes=await api('/api/streams');
    const streams=Array.isArray(streamsRes)?streamsRes:[];
    const settings=await api('/api/settings');
    const entries=await Promise.all(streams.map(async function(s){
      const access=await getPlaybackAccess(s.stream_key,settings,s.policy_json);
      return {stream:s,urls:getAllURLs(s.stream_key,settings,s.name,access),access:access};
    }));

    c.innerHTML=
      '<div class="page-header"><h1 class="page-title">Embed Kodlari</h1></div>'+
      (streams.length===0?'<div class="card"><div class="empty-state"><h3>Henuz yayin yok</h3><p style="color:var(--text-muted)">Once bir yayin olusturun</p></div></div>'
      :entries.map(function(entry){
        var s=entry.stream;
        var u=entry.urls;
        return '<div class="card" style="margin-bottom:20px">'+
          '<div class="card-header"><div class="card-title">'+escHtml(s.name)+' <span class="badge badge-'+s.status+'" style="margin-left:8px">'+(s.status==='live'?'CANLI':'Cevrimdisi')+'</span></div></div>'+
          (entry.access&&entry.access.needs_token?'<div class="form-hint" style="margin-bottom:10px;color:var(--warning)">Token korumasi aktif. Aasagidaki player ve URL alanlari gecici playback token ile uretildi.</div>':'')+
          copyField('iframe Embed','<iframe src="'+u.embed+'" width="1280" height="720" frameborder="0" allowfullscreen></iframe>')+
          '<details style="margin-top:8px"><summary style="cursor:pointer;font-weight:600;padding:8px 0;color:var(--text-secondary)">Video Akis URL\'leri (7)</summary>'+
            copyField('HLS',u.hls)+copyField('LL-HLS',u.ll_hls)+copyField('DASH',u.dash)+
            copyField('HTTP-FLV',u.http_flv)+copyField('WHEP',u.whep)+copyField('fMP4',u.fmp4)+copyField('WebM',u.webm)+
          '</details>'+
          '<details style="margin-top:4px"><summary style="cursor:pointer;font-weight:600;padding:8px 0;color:var(--text-secondary)">Protokol URL\'leri (7)</summary>'+
            copyField('RTMP',u.rtmp)+copyField('RTSP',u.rtsp)+copyField('SRT',u.srt)+
            copyField('RTP',u.rtp)+copyField('MPEG-TS',u.mpegts)+copyField('RTSP Cikis',u.rtsp_out)+copyField('SRT Cikis',u.srt_out)+
          '</details>'+
          '<details style="margin-top:4px"><summary style="cursor:pointer;font-weight:600;padding:8px 0;color:var(--text-secondary)">Ses URL\'leri (8)</summary>'+
            copyField('MP3',u.mp3)+copyField('AAC',u.aac)+copyField('OGG',u.ogg)+
            copyField('WAV',u.wav)+copyField('FLAC',u.flac)+copyField('HLS Ses',u.hls_audio)+copyField('DASH Ses',u.dash_audio)+copyField('Icecast',u.icecast)+
          '</details>'+
          '<details style="margin-top:4px"><summary style="cursor:pointer;font-weight:600;padding:8px 0;color:var(--text-secondary)">Player & Embed (3)</summary>'+
            copyField('Player URL',u.play)+copyField('Audio Player URL',playerURLForFormat(u.play,'aac'))+copyField('Embed URL',u.embed)+
          '</details>'+
          '<details style="margin-top:4px"><summary style="cursor:pointer;font-weight:600;padding:8px 0;color:var(--text-secondary)">Kullanim ve Tanilama</summary>'+
            renderDeliveryUsageCard(s.id,u,{telemetryMode:'navigate'})+
          '</details>'+
        '</div>';
      }).join(''));
  }catch(e){
    c.innerHTML='<div class="card"><div class="empty-state"><h3>Embed kodlari yuklenemedi</h3><p style="color:var(--text-muted)">'+escHtml(e.message||'Bilinmeyen hata')+'</p></div></div>';
  }
}
// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETTINGS - GENERAL ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderSettingsGeneral(c){
  const s=await api('/api/settings');
  window._generalAdvanced=window._generalAdvanced||false;
  const adv=window._generalAdvanced;
  c.innerHTML=
    '<div class="studio-page">'+
      '<div class="studio-hero">'+
        '<h1 class="studio-hero-title"><i class="bi bi-gear" style="margin-right:10px"></i>Genel Ayarlar</h1>'+
        '<p class="studio-hero-sub">Sunucunun omurgasini olusturan temel yapilandirma ayarlari. Basit mod en sik kullanilan ayarlari gosterir; gelismis mod tum teknik alanlari acar.</p>'+
      '</div>'+
      '<div class="studio-toolbar">'+
        '<div class="studio-toolbar-group"><div class="segmented"><button class="segment '+(!adv?'active':'')+'" onclick="window._generalAdvanced=false;navigate(\'settings-general\')">Basit Mod</button><button class="segment '+(adv?'active':'')+'" onclick="window._generalAdvanced=true;navigate(\'settings-general\')">Gelismis Mod</button></div></div>'+
        '<div class="studio-toolbar-group"><button class="btn btn-primary btn-sm" onclick="saveGeneralSettingsExtended()"><i class="bi bi-check-lg"></i> Tum Ayarlari Kaydet</button></div>'+
      '</div>'+
      '<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span><i class="bi bi-person-badge" style="margin-right:8px;color:var(--accent)"></i>Kimlik ve Yerellesme</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:12px">Sunucu adi, dil ve saat dilimi ayarlari.</div>'+
          settingInput('server_name','Sunucu Adi',s.server_name||'FluxStream','text','Panel basliginda ve embed kodlarinda goruntulenir.')+
          settingSelect('language','Dil',s.language||'tr',[{value:'tr',label:'Turkce'},{value:'en',label:'English'},{value:'de',label:'Deutsch'},{value:'es',label:'Espanol'},{value:'fr',label:'Francais'}],'Login, setup ve panel arayuzu bu dilde acilir.')+
          settingInput('timezone','Saat Dilimi',s.timezone||'Europe/Istanbul','text','Tarih ve saat gosterimleri bu timezone ile yorumlanir.')+
          '<div class="form-group"><label class="form-label">Tema</label><select class="form-select setting-input" data-key="theme"><option value="light" '+((s.theme||'light')==='light'?'selected':'')+'>Acik (Light)</option><option value="dark" '+((s.theme||'')==='dark'?'selected':'')+'>Koyu (Dark)</option><option value="minimal" '+((s.theme||'')==='minimal'?'selected':'')+'>Minimal</option></select><div class="form-hint">Admin panelinin gorsel temasini belirler.</div></div>'+
          '<div class="setting-row"><div><div class="setting-label">Kolay Mod</div><div class="setting-desc">Acikken Hizli Ayarlar ekrani varsayilan olarak one cikar.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="guided_mode_enabled" '+(s.guided_mode_enabled!=='false'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
        '</div>'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span><i class="bi bi-globe" style="margin-right:8px;color:var(--accent)"></i>Public Erisim ve Portlar</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:12px">Embed linkleri ve player URL\'leri bu ayarlara gore uretilir.</div>'+
          settingInput('embed_domain','Public Domain / IP',s.embed_domain||'','text','Bossa mevcut host kullanilir. Ornek: stream.ornek.com')+
          '<div class="studio-form-grid">'+
            settingInput('embed_http_port','Public HTTP Port',s.embed_http_port||s.http_port||'8844','number','Embed linklerinde')+
            settingInput('embed_https_port','Public HTTPS Port',s.embed_https_port||s.https_port||'443','number','SSL linklerinde')+
          '</div>'+
          '<div class="setting-row"><div><div class="setting-label">Player Kalite Secici</div><div class="setting-desc">ABR yayinlarda izleyici kaliteyi elle secebilir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="player_quality_selector" '+(s.player_quality_selector!=='false'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
        '</div>'+
      '</div>'+
      (adv?'<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span><i class="bi bi-hdd-rack" style="margin-right:8px;color:var(--text-muted)"></i>Sunucu Portlari</span><span class="tag tag-blue">Gelismis</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:12px">Degistirirseniz sunucu restart gerektirir.</div>'+
          settingInput('http_port','Web HTTP Port',s.http_port||'8844','number','Admin paneli portu')+
          settingInput('https_port','Web HTTPS Port',s.https_port||'443','number','SSL portu')+
        '</div>'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span><i class="bi bi-tools" style="margin-right:8px;color:var(--text-muted)"></i>Bakim ve Saklama</span><span class="tag tag-blue">Gelismis</span></div>'+
          '<div class="setting-row"><div><div class="setting-label">Otomatik Bakim</div><div class="setting-desc">Temizleme ve bakim islerini zamanli calistirir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="maintenance_auto_cleanup" '+(s.maintenance_auto_cleanup!=='false'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          settingInput('recordings_retention_days','Kayit Saklama (gun)',s.recordings_retention_days||'30','number','0 = sinirsiz')+
        '</div>'+
      '</div>':'')+
    '</div>';
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETTINGS - PROTOCOLS ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderSettingsEmbed(c){
  const s=await api('/api/settings');
  const publicConfig=getPublicBase(s);
  const streamHost=getConfiguredDomain(s);
  c.innerHTML=
    '<div class="studio-page">'+
      '<div class="studio-hero">'+
        '<h1 class="studio-hero-title"><i class="bi bi-globe2" style="margin-right:10px"></i>Alan Adi ve Embed</h1>'+
        '<p class="studio-hero-sub">Embed kodlari ve izleme linklerinde kullanilan public alan adini buradan belirleyin. Web erisimi HTTP/HTTPS portlarindan gider, yayin gonderimi RTMP/RTMPS portlarindan devam eder.</p>'+
      '</div>'+
      '<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span>Public Erisim Ayarlari</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:12px">Bu alanlar uretilen tum embed ve player linklerini etkiler.</div>'+
          settingInput('embed_domain','Public Domain / IP',s.embed_domain||'','text','Bos birakirsaniz mevcut host kullanilir. Ornek: stream.ornek.com')+
          '<div class="studio-form-grid">'+
            settingInput('embed_http_port','Public HTTP Port',s.embed_http_port||s.http_port||'8844','number','HTTP linkleri icin')+
            settingInput('embed_https_port','Public HTTPS Port',s.embed_https_port||s.https_port||'443','number','HTTPS linkleri icin')+
          '</div>'+
          '<div class="setting-row"><div><div class="setting-label">HTTPS Link Uret</div><div class="setting-desc">SSL etkin ve sertifika hazirsa embed linkleri HTTPS olur.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="embed_use_https" '+(shouldUsePublicHTTPS(s)?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<button class="btn btn-primary" style="margin-top:12px" onclick="saveSettingsCategory(\'embed\')"><i class="bi bi-check-lg"></i> Kaydet</button>'+
        '</div>'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span>Canli Onizleme</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:12px">Mevcut ayarlarla uretilen ornek URL\'ler. Kopyalayip kullanabilirsiniz.</div>'+
          copyField('Player Base URL',publicConfig.base)+
          copyField('RTMP Sunucu','rtmp://'+streamHost+':'+(s.rtmp_port||'1935')+'/live')+
          copyField('RTMPS Sunucu','rtmps://'+streamHost+':'+(s.rtmps_port||'1936')+'/live')+
          '<div class="studio-alert info" style="margin-top:12px"><i class="bi bi-info-circle" style="margin-right:6px"></i>Bu URL\'ler sol paneldeki ayarlara gore otomatik uretilir. Ayarlar degisince buradaki ornekler de guncellenir.</div>'+
        '</div>'+
      '</div>'+
    '</div>';
}

async function renderSettingsProtocols(c){
  const s=await api('/api/settings');
  const protos=[
    {name:'RTMP',desc:'En yaygin - OBS, Wirecast, vMix',icon:'bi-broadcast',key:'rtmp_enabled',port:'rtmp_port',def:'1935',extra:[settingInput('rtmp_chunk_size','Chunk Size',s.rtmp_chunk_size||'4096','number',''),settingInput('rtmp_max_conns','Maks Baglanti',s.rtmp_max_conns||'100','number','')]},
    {name:'RTMPS',desc:'Sifreli RTMP - guvenli aglar',icon:'bi-lock',key:'rtmps_enabled',port:'rtmps_port',def:'1936',extra:[]},
    {name:'SRT',desc:'Dusuk gecikme + guvenilmez ag',icon:'bi-lightning',key:'srt_enabled',port:'srt_port',def:'9000',extra:[settingInput('srt_latency','Latency (ms)',s.srt_latency||'120','number','')]},
    {name:'RTP',desc:'Profesyonel encoder push',icon:'bi-arrow-right-circle',key:'rtp_enabled',port:'rtp_port',def:'5004',extra:[]},
    {name:'RTSP',desc:'IP kameralar, profesyonel encoderlar',icon:'bi-camera-video',key:'rtsp_enabled',port:'rtsp_port',def:'8554',extra:[]},
    {name:'WebRTC/WHIP',desc:'Tarayicidan dogrudan yayin',icon:'bi-browser-chrome',key:'webrtc_enabled',port:'webrtc_port',def:'8855',extra:[]},
    {name:'MPEG-TS',desc:'Uydu alici, profesyonel broadcast',icon:'bi-reception-4',key:'mpegts_enabled',port:'mpegts_port',def:'9001',extra:[]}
  ];
  const activeCount=protos.filter(function(p){return s[p.key]==='true'}).length;
  c.innerHTML=
    '<div class="studio-page">'+
      '<div class="studio-hero">'+
        '<h1 class="studio-hero-title"><i class="bi bi-arrow-down-circle" style="margin-right:10px"></i>Giris Protokolleri</h1>'+
        '<p class="studio-hero-sub">Encoder\'lardan kabul edilen yayin protokollerini buradan yonetin. Kullanmadiginiz protokolleri kapatarak guvenlik yuzeyini daraltabilirsiniz.</p>'+
        '<div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active"><i class="bi bi-check-circle"></i> '+activeCount+'/'+protos.length+' Aktif</span></div>'+
      '</div>'+
      '<div style="display:grid;gap:14px">'+
        protos.map(function(p){
          var enabled=s[p.key]==='true';
          return '<div class="studio-card" style="border-left:4px solid '+(enabled?'var(--success)':'var(--border)')+'">'+
            '<div style="display:flex;align-items:center;justify-content:space-between;gap:14px;flex-wrap:wrap">'+
              '<div style="display:flex;align-items:center;gap:12px"><i class="bi '+p.icon+'" style="font-size:22px;color:'+(enabled?'var(--accent)':'var(--text-muted)')+'"></i><div><div style="font-weight:700;font-size:15px">'+p.name+' <span class="tag '+(enabled?'tag-green':'tag-red')+'">'+(enabled?'Aktif':'Kapali')+'</span></div><div style="font-size:12px;color:var(--text-muted)">'+p.desc+'</div></div></div>'+
              '<div style="display:flex;align-items:center;gap:14px"><div style="font-size:13px;color:var(--text-secondary)">Port: <strong>'+(s[p.port]||p.def)+'</strong></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="'+p.key+'" '+(enabled?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
            '</div>'+
            '<div class="studio-form-grid" style="margin-top:10px">'+settingInput(p.port,'Port',s[p.port]||p.def,'number','')+p.extra.join('')+'</div>'+
          '</div>';
        }).join('')+
      '</div>'+
      '<div style="margin-top:16px"><button class="btn btn-primary" onclick="saveSettingsCategory(\'protocols\')"><i class="bi bi-check-lg"></i> Tum Protokolleri Kaydet</button></div>'+
    '</div>';
}
function protoCard(name,desc,enableKey,s,portKey,portVal,extra){
  const enabled=s[enableKey]==='true';
  return '<div class="card" style="margin-bottom:16px">'+
    '<div class="card-header"><div><div class="card-title">'+name+' <span class="tag '+(enabled?'tag-green':'tag-red')+'">'+(enabled?'Aktif':'Kapali')+'</span></div>'+
      '<div class="setting-desc" style="margin-top:4px">'+desc+'</div></div>'+
      '<label class="toggle"><input type="checkbox" class="setting-input" data-key="'+enableKey+'" '+(enabled?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
    '<div class="card-grid card-grid-2" style="margin-top:12px">'+
      settingInput(portKey,'Port',portVal,'number','')+
      (extra||[]).join('')+
    '</div></div>';
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETTINGS - OUTPUTS ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderSettingsOutputs(c){
  const s=await api('/api/settings');
  const outs=[
    {name:'HLS',desc:'Apple HLS - en uyumlu, adaptif',icon:'bi-play-circle',key:'hls_enabled',cat:'video',extra:[settingInput('hls_segment_duration','Segment (sn)',s.hls_segment_duration||'2','number',''),settingInput('hls_playlist_length','Playlist Uzunlugu',s.hls_playlist_length||'6','number','')]},
    {name:'Low Latency HLS',desc:'2 saniye alti gecikme',icon:'bi-lightning',key:'hls_ll_enabled',cat:'video',extra:[]},
    {name:'DASH',desc:'MPEG-DASH adaptive streaming',icon:'bi-speedometer2',key:'dash_enabled',cat:'video',extra:[settingInput('dash_segment_duration','Segment (sn)',s.dash_segment_duration||'2','number','')]},
    {name:'HTTP-FLV',desc:'Ultra dusuk gecikme (~1sn)',icon:'bi-film',key:'httpflv_enabled',cat:'video',extra:[]},
    {name:'WebRTC/WHEP',desc:'Sub-second gecikme (<500ms)',icon:'bi-webcam',key:'whep_enabled',cat:'video',extra:[]},
    {name:'MP3',desc:'MP3 ses akisi',icon:'bi-music-note-beamed',key:'mp3_enabled',cat:'audio',extra:[settingInput('mp3_bitrate','Bitrate (kbps)',s.mp3_bitrate||'128','number','')]},
    {name:'AAC',desc:'Yuksek kalite ses',icon:'bi-soundwave',key:'aac_out_enabled',cat:'audio',extra:[]},
    {name:'Icecast',desc:'Icecast uyumlu akis',icon:'bi-broadcast',key:'icecast_enabled',cat:'audio',extra:[settingInput('icecast_port','Port',s.icecast_port||'8000','number','')]}
  ];
  const videoOuts=outs.filter(function(o){return o.cat==='video'});
  const audioOuts=outs.filter(function(o){return o.cat==='audio'});
  const activeCount=outs.filter(function(o){return s[o.key]==='true'}).length;
  function renderOutputItem(o){
    var enabled=s[o.key]==='true';
    return '<div style="display:flex;align-items:center;justify-content:space-between;gap:12px;padding:14px 16px;border:1px solid '+(enabled?'rgba(16,185,129,.2)':'var(--border)')+';border-radius:16px;background:'+(enabled?'linear-gradient(180deg,#f0fdf4,#f8fbff)':'#fff')+'">'+
      '<div style="display:flex;align-items:center;gap:12px"><i class="bi '+o.icon+'" style="font-size:20px;color:'+(enabled?'var(--success)':'var(--text-muted)')+'"></i><div><div style="font-weight:700;font-size:14px">'+o.name+' <span class="tag '+(enabled?'tag-green':'tag-red')+'" style="margin-left:4px">'+(enabled?'Acik':'Kapali')+'</span></div><div style="font-size:12px;color:var(--text-muted)">'+o.desc+'</div></div></div>'+
      '<label class="toggle"><input type="checkbox" class="setting-input" data-key="'+o.key+'" '+(enabled?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      (o.extra.length?'<div class="studio-form-grid" style="margin-top:6px;padding-left:44px">'+o.extra.join('')+'</div>':'');
  }
  c.innerHTML=
    '<div class="studio-page">'+
      '<div class="studio-hero">'+
        '<h1 class="studio-hero-title"><i class="bi bi-arrow-up-circle" style="margin-right:10px"></i>Cikis Formatlari</h1>'+
        '<p class="studio-hero-sub">Izleyicilere sunulan video ve ses formatlarini buradan yonetin. Kullanmadiginiz formatlari kapatarak sunucu yukunu azaltabilirsiniz.</p>'+
        '<div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active"><i class="bi bi-check-circle"></i> '+activeCount+'/'+outs.length+' Aktif Format</span></div>'+
      '</div>'+
      '<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card"><div class="studio-section-title"><span><i class="bi bi-camera-video" style="margin-right:8px;color:var(--accent)"></i>Video Formatlari</span></div><div style="display:grid;gap:10px;margin-top:10px">'+videoOuts.map(renderOutputItem).join('')+'</div></div>'+
        '<div class="studio-card"><div class="studio-section-title"><span><i class="bi bi-music-note-beamed" style="margin-right:8px;color:var(--accent)"></i>Ses Formatlari</span></div><div style="display:grid;gap:10px;margin-top:10px">'+audioOuts.map(renderOutputItem).join('')+'</div></div>'+
      '</div>'+
      '<div style="margin-top:16px"><button class="btn btn-primary" onclick="saveSettingsCategory(\'outputs\')"><i class="bi bi-check-lg"></i> Tum Cikislari Kaydet</button></div>'+
    '</div>';
}
function outputCard(name,desc,enableKey,s,extra){
  const enabled=s[enableKey]==='true';
  return '<div class="card" style="margin-bottom:16px">'+
    '<div class="card-header"><div><div class="card-title">'+name+' <span class="tag '+(enabled?'tag-green':'tag-red')+'">'+(enabled?'Aktif':'Kapali')+'</span></div>'+
      '<div class="setting-desc" style="margin-top:4px">'+desc+'</div></div>'+
      '<label class="toggle"><input type="checkbox" class="setting-input" data-key="'+enableKey+'" '+(enabled?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
    (extra.length?'<div class="card-grid card-grid-2" style="margin-top:12px">'+extra.join('')+'</div>':'')+
  '</div>';
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETTINGS - SSL ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderSettingsSSL(c){
  const s=await api('/api/settings');
  const sslStatus=await api('/api/ssl/status');
  const webSSL=(sslStatus&&sslStatus.web)||{};
  const streamSSL=(sslStatus&&sslStatus.stream)||{};
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">SSL/TLS Sertifika</h1></div>'+
    '<div class="card" style="max-width:920px;margin-bottom:16px">'+
      '<div class="card-title" style="margin-bottom:12px">Kullanim Mantigi</div>'+
      '<div class="form-hint" style="line-height:1.8">Web HTTPS sertifikasi admin paneli ve embed/player sayfalari icin kullanilir. Stream SSL ise yalnizca RTMPS ingest tarafini korur. Isterseniz ayni domaini, isterseniz ayri domain ve ayri sertifika kullanabilirsiniz. Let\' Encrypt icin alan adlarinin bu VPS\'e yonlenmis olmasi ve 80/443 portlarinin acik olmasi gerekir.</div>'+
    '</div>'+
    '<div class="card-grid card-grid-2">'+
      renderSSLProfileCard('web','Web HTTPS',webSSL,s,'ssl_enabled','https_port','ssl_mode','ssl_cert_path','ssl_key_path','ssl_le_domain','ssl_le_email','Admin paneli, embed ve player linkleri bu sertifikayi kullanir.')+
      renderSSLProfileCard('stream','Stream RTMPS',streamSSL,s,'rtmps_enabled','rtmps_port','stream_ssl_mode','stream_ssl_cert_path','stream_ssl_key_path','stream_ssl_le_domain','stream_ssl_le_email','OBS veya baska encoder RTMPS ile baglanacaksa bu sertifika kullanilir.')+
    '</div>'+
    '<div style="margin-top:16px"><button class="btn btn-primary" onclick="saveSSLSettings()">SSL Ayarlarini Kaydet</button></div>';
}
function renderSSLProfileCard(target,title,status,s,enableKey,portKey,modeKey,certKey,keyKey,domainKey,emailKey,desc){
  const mode=String((s&&s[modeKey])||'file').toLowerCase();
  const ready=!!(status&&status.ready);
  const enabled=isTruthy(s&&s[enableKey]);
  return '<div class="card">'+
    '<div class="card-title" style="margin-bottom:10px">'+title+'</div>'+
    '<div class="form-hint" style="line-height:1.7;margin-bottom:14px">'+desc+'</div>'+
    '<div class="setting-row"><div><div class="setting-label">Ozellik Acik</div><div class="setting-desc">Kapaliysa bu profil hic kullanilmaz.</div></div>'+
      '<label class="toggle"><input type="checkbox" class="setting-input" data-key="'+enableKey+'" '+(enabled?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
    '<div style="padding:14px;background:var(--bg-primary);border-radius:var(--radius-sm);margin-bottom:14px">'+
      '<div style="display:flex;justify-content:space-between;gap:10px;align-items:center;margin-bottom:8px"><strong>Durum</strong><span class="tag '+(ready?'tag-green':'tag-red')+'">'+(ready?'Hazir':'Hazir Degil')+'</span></div>'+
      '<div class="form-hint">Port: <b>'+escHtml(String((s&&s[portKey])||(status&&status[portKey])||''))+'</b></div>'+
      (mode==='letsencrypt'
        ?'<div class="form-hint" style="margin-top:6px">Domain: <b>'+escHtml(String((s&&s[domainKey])||(status&&status.domain)||'-'))+'</b></div>'
        :'<div class="form-hint" style="margin-top:6px">CRT: <code>'+escHtml(String((status&&status.cert_path)||(s&&s[certKey])||'-'))+'</code></div><div class="form-hint">KEY: <code>'+escHtml(String((status&&status.key_path)||(s&&s[keyKey])||'-'))+'</code></div>')+
    '</div>'+
    '<div class="form-group"><label class="form-label">Sertifika Modu</label><select class="form-select setting-input" data-key="'+modeKey+'"><option value="file" '+(mode==='file'?'selected':'')+'>Manuel CRT/KEY</option><option value="letsencrypt" '+(mode==='letsencrypt'?'selected':'')+'>Let\' Encrypt</option></select><div class="form-hint">Manuel modda dosya yuklersiniz. Let\' Encrypt modunda domain ve e-posta yeterlidir.</div></div>'+
    '<div class="form-group"><label class="form-label">CRT / PEM Yukle</label><input type="file" id="ssl-cert-file-'+target+'" accept=".crt,.pem,.cert" class="form-input" style="padding:8px"></div>'+
    '<div class="form-group"><label class="form-label">KEY / PEM Yukle</label><input type="file" id="ssl-key-file-'+target+'" accept=".key,.pem" class="form-input" style="padding:8px"></div>'+
    '<div style="margin-bottom:16px"><button class="btn btn-secondary" onclick="uploadSSL(\''+target+'\')">Bu Profil Icin Sertifika Yukle</button></div>'+
    settingInput(certKey,'Sertifika Dosyasi (.crt)',s[certKey]||'','text','Orn: /opt/fluxstream/data/certs/'+target+'/server.crt')+
    settingInput(keyKey,'Ozel Anahtar (.key)',s[keyKey]||'','text','Orn: /opt/fluxstream/data/certs/'+target+'/server.key')+
    settingInput(domainKey,'Let\' Encrypt Domain',s[domainKey]||'','text','Orn: '+(target==='web'?'panel.example.com':'stream.example.com'))+
    settingInput(emailKey,'Let\' Encrypt E-posta',s[emailKey]||'','text','Bildirim ve yenileme icin kullanilir.')+
  '</div>';
}
async function uploadSSL(target){
  const certInput=document.getElementById('ssl-cert-file-'+target);
  const keyInput=document.getElementById('ssl-key-file-'+target);
  if(!certInput.files[0]||!keyInput.files[0]){toast('Her iki dosyayi da secin','error');return}
  const fd=new FormData();
  fd.append('cert',certInput.files[0]);
  fd.append('key',keyInput.files[0]);
  fd.append('target',target);
  try{
    const res=await fetch('/api/ssl/upload',{method:'POST',body:fd});
    const data=await res.json();
    if(data.success){toast('SSL sertifikalari yuklendi. Uygulamak icin restart gerekli.');navigate('settings-ssl')}
    else{toast(data.message||'Yukleme hatasi','error')}
  }catch(e){toast('Yukleme hatasi: '+e.message,'error')}
}
// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETTINGS - SECURITY ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderSettingsSecurity(c){
  const s=await api('/api/settings');
  const streamKeyOn=s.stream_key_required==='true';
  const tokenOn=s.token_enabled==='true';
  const riskLevel=(!streamKeyOn&&!tokenOn)?'Dusuk':(streamKeyOn&&tokenOn?'Yuksek':'Orta');
  const riskTone=riskLevel==='Yuksek'?'green':(riskLevel==='Orta'?'yellow':'red');
  c.innerHTML=
    '<div class="studio-page">'+
      '<div class="studio-hero">'+
        '<h1 class="studio-hero-title"><i class="bi bi-shield-lock" style="margin-right:10px"></i>Guvenlik Ayarlari</h1>'+
        '<p class="studio-hero-sub">Yayin gonderimine ve izleme erisimini koruma altina alan ayarlar. Koruma seviyesi arttirilirsa yetkisiz erisim engellenir ancak entegrasyon karmasikligi artar.</p>'+
        '<div class="studio-pill-row" style="margin-top:14px">'+
          '<span class="studio-pill '+(riskTone==='green'?'active':'')+'"><i class="bi bi-shield-check"></i> Koruma: '+riskLevel+'</span>'+
          '<span class="studio-pill '+(streamKeyOn?'active':'')+'"><i class="bi bi-key"></i> Stream Key: '+(streamKeyOn?'Zorunlu':'Serbest')+'</span>'+
          '<span class="studio-pill '+(tokenOn?'active':'')+'"><i class="bi bi-lock"></i> Token: '+(tokenOn?'Aktif':'Kapali')+'</span>'+
        '</div>'+
      '</div>'+
      '<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span><i class="bi bi-key" style="margin-right:8px;color:var(--accent)"></i>Yayin Korumasi</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:14px">Encoder\'larin yayin gonderirken dogrulama gereksinimi.</div>'+
          '<div class="setting-row"><div><div class="setting-label">Stream Key Zorunlu</div><div class="setting-desc">Kapatilirsa herhangi bir key ile yayin gonderilebilir. <strong>Onerilir: Acik</strong></div></div>'+
            '<label class="toggle"><input type="checkbox" class="setting-input" data-key="stream_key_required" '+(streamKeyOn?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          (!streamKeyOn?'<div class="studio-alert warning" style="margin-top:10px"><i class="bi bi-exclamation-triangle" style="margin-right:6px"></i><strong>Risk:</strong> Stream key kontrolu kapali. Herkes RTMP ile yayin gonderebilir.</div>':'')+
        '</div>'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span><i class="bi bi-lock" style="margin-right:8px;color:var(--accent)"></i>Izleme Korumasi</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:14px">Izleyicilerin yayin izlerken token sunmasini zorunlu kilar.</div>'+
          '<div class="setting-row"><div><div class="setting-label">Token Dogrulama</div><div class="setting-desc">Acilirsa izleme linkleri gecerli token olmadan calismaz.</div></div>'+
            '<label class="toggle"><input type="checkbox" class="setting-input" data-key="token_enabled" '+(tokenOn?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          settingInput('token_duration','Token Gecerlilik Suresi (sn)',s.token_duration||'60','number','Token uretildikten sonra bu kadar sure gecerli kalir. Onerilen: 60-3600')+
          (tokenOn?'<div class="studio-alert info" style="margin-top:10px"><i class="bi bi-info-circle" style="margin-right:6px"></i>Token uretimi icin <a href="#" onclick="navigate(\'security-tokens\');return false" style="color:var(--accent);font-weight:600">Tokenlar</a> ekranini kullanin.</div>':'')+
        '</div>'+
      '</div>'+
      '<div class="studio-card">'+
        '<div class="studio-section-title"><span><i class="bi bi-speedometer" style="margin-right:8px;color:var(--accent)"></i>Trafik Sinirlamasi</span></div>'+
        '<div class="studio-form-grid" style="margin-top:10px">'+
          settingInput('rate_limit','Rate Limit (istek/sn)',s.rate_limit||'100','number','Ani ve asiri istekleri sinirlar. Cok dusuk degerler normal izleyicileri de etkileyebilir. Onerilen: 50-200')+
        '</div>'+
      '</div>'+
      '<div style="margin-top:16px"><button class="btn btn-primary" onclick="saveSettingsCategory(\'security\')"><i class="bi bi-check-lg"></i> Guvenlik Ayarlarini Kaydet</button></div>'+
    '</div>';
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETTINGS - STORAGE ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderSettingsStorage(c){
  teardownRecordingPreview();
  const [s,report,archivesRes,recsRes,streamsRes,savedRes,backupsRes,backupArchivesRes,upgradeRes,remuxJobsRes]=await Promise.all([
    api('/api/settings'),
    api('/api/health/report'),
    api('/api/recordings/archives'),
    api('/api/recordings'),
    api('/api/streams'),
    api('/api/recordings/library'),
    api('/api/system/backups'),
    api('/api/system/backups/archives'),
    api('/api/system/upgrade/plan'),
    api('/api/recordings/remux/jobs')
  ]);
  {
    const commands=(upgradeRes&&upgradeRes.commands)||{};
    const restoreCmd=commands.backup_restore||'fluxstream backup restore fluxstream-backup-YYYYMMDD-HHMMSS.tar.gz';
    const storageData=normalizeStorageSnapshot(s,report,archivesRes,recsRes,streamsRes,savedRes,backupsRes,backupArchivesRes,remuxJobsRes);
    c.innerHTML=renderStorageCenter(storageData,restoreCmd);
    window._recordingPreviewSelection=null;
    resetRecordingPreviewPanel();
    applyStorageSnapshot(storageData,{resetPreview:true});
    updateStorageProviderUI();
    return;
  }
  const archiveSummary=report&&report.storage&&report.storage.archive?report.storage.archive:{};
  const archives=Array.isArray(archivesRes)?archivesRes:[];
  const recs=Array.isArray(recsRes)?recsRes:[];
  const activeRecordings=recs.filter(function(item){return item&&item.Status==='recording';});
  const streams=Array.isArray(streamsRes)?streamsRes:[];
  const saved=Array.isArray(savedRes)?savedRes:[];
  const backups=(backupsRes&&Array.isArray(backupsRes.items))?backupsRes.items:[];
  const backupArchives=Array.isArray(backupArchivesRes)?backupArchivesRes:[];
  const remuxJobs=Array.isArray(remuxJobsRes)?remuxJobsRes:[];
  const archiveMap={};
  const backupArchiveMap={};
  archives.forEach(function(item){archiveMap[item.stream_key+'::'+item.filename]=item;});
  backupArchives.forEach(function(item){backupArchiveMap[item.name]=item;});
  const archiveEnabled=s&&s.archive_enabled==='true';
  const backupArchiveEnabled=s&&s.backup_archive_enabled==='true';
  const commands=(upgradeRes&&upgradeRes.commands)||{};
  const restoreCmd=commands.backup_restore||'fluxstream backup restore fluxstream-backup-YYYYMMDD-HHMMSS.tar.gz';
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Depolama ve Arsiv Merkezi</h1><div style="display:flex;gap:10px;flex-wrap:wrap"><button class="btn btn-primary" onclick="showRecordModal()">Kayit Baslat</button><button class="btn btn-secondary" onclick="createSystemBackupFromStorage(false)">Hafif Yedek Al</button><button class="btn btn-secondary" onclick="createSystemBackupFromStorage(true)">Kayitlarla Yedek Al</button></div></div>'+
    '<div id="storage-active-banner">'+(activeRecordings.length?'<div class="card" style="margin-bottom:16px;border-color:rgba(239,68,68,.28);box-shadow:0 8px 22px rgba(239,68,68,.08)"><div class="card-header"><div><div class="card-title">Aktif Kayit Uyarisi</div><div class="form-hint">Calisan kayit oturumlari burada sabit tutulur. Durdur dugmesine buradan da erisebilirsiniz.</div></div><span class="badge badge-live">'+fmtInt(activeRecordings.length)+' aktif kayit</span></div><div style="display:flex;gap:10px;flex-wrap:wrap">'+activeRecordings.map(function(r){return '<div class="tag tag-red" style="display:flex;align-items:center;gap:10px;padding:8px 12px"><span><strong>'+escHtml(String(r.StreamKey||'-'))+'</strong> · '+escHtml(String(r.Format||'').toUpperCase())+'</span>'+(r.Status==='recording'?'<button class="btn btn-sm btn-danger" onclick="stopRec(\''+r.ID+'\')">Durdur</button>':'')+'</div>';}).join('')+'</div></div>':'')+'</div>'+
    '<div id="storage-remux-jobs">'+(remuxJobs.length?'<div class="card" style="margin-bottom:16px;background:linear-gradient(180deg,#f8fbff 0%,#f2f8ff 100%)"><div class="card-header"><div><div class="card-title">Donusum ve Senkron Isleri</div><div class="form-hint">MP4 hazirlama ve benzeri uzun isler arka planda devam eder.</div></div><button class="btn btn-secondary btn-sm" onclick="refreshStorageSnapshot({resetPreview:false})">Yenile</button></div><div style="display:flex;gap:10px;flex-wrap:wrap">'+remuxJobs.slice(0,6).map(function(job){var tone=job.status==='completed'?'green':(job.status==='error'?'red':'yellow'); var label=job.status==='completed'?'Hazir':(job.status==='error'?'Hata':'Calisiyor'); return '<div class="tag tag-'+tone+'" style="display:flex;align-items:center;gap:8px;padding:8px 12px"><span><strong>'+escHtml(job.source_name||'-')+'</strong> → '+escHtml((job.target_format||'mp4').toUpperCase())+'</span><span>'+label+'</span></div>';}).join('')+'</div></div>':'')+'</div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><div><div class="card-title">Ne yapmak istiyorsunuz?</div><div class="form-hint">Bu ekran uc seyi birlikte yonetir: kayitlar, yedekler ve harici arsiv kopyalari.</div></div></div><div class="card-grid card-grid-3">'+
      '<div style="padding:14px;border:1px solid var(--border);border-radius:12px;background:var(--bg-primary)"><div style="font-weight:700;margin-bottom:6px">1. Yerelde tut</div><div class="form-hint">En kolay baslangic. Kayitlar bu sunucuda kalir.</div></div>'+
      '<div style="padding:14px;border:1px solid var(--border);border-radius:12px;background:var(--bg-primary)"><div style="font-weight:700;margin-bottom:6px">2. Dis kopya ekle</div><div class="form-hint">Ayni kayitlari ikinci bir hedefe de gonderirsin. Yedek icin en guvenli yoldur.</div></div>'+
      '<div style="padding:14px;border:1px solid var(--border);border-radius:12px;background:var(--bg-primary)"><div style="font-weight:700;margin-bottom:6px">3. Geri yukle ve indir</div><div class="form-hint">Arsivden geri getir, yerelden indir veya gerekmiyorsa sil.</div></div>'+
    '</div><div class="form-hint" style="margin-top:12px">MP4 secilen kayitlar yayin boyunca guvenli bicimde yakalanir, yayin bitince izlenebilir dosyaya finalize edilir.</div></div>'+
    '<div id="storage-stats-grid" class="card-grid card-grid-4" style="margin-bottom:16px">'+
      statCard('blue','bi-hdd-fill',formatBytes((report&&report.storage&&report.storage.recordings_bytes)||0),'Yerel Kayitlar')+
      statCard('purple','bi-archive-fill',fmtInt(backups.length),'Yerel Yedekler')+
      statCard('orange','bi-cloud-arrow-up-fill',fmtInt(archives.length),'Kayit Arsivi')+
      statCard('green','bi-safe2-fill',fmtInt(backupArchives.length),'Yedek Arsivi')+
    '</div>'+
    '<div class="card-grid card-grid-2" style="margin-bottom:16px">'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Yerel Depolama ve Temizlik</div>'+
        settingInput('storage_max_gb','Maksimum Depolama (GB)',s.storage_max_gb||'50','number','Toplam kayit ve yedek alanini izlemek icin uyarilarda kullanilir.')+
        settingInput('storage_auto_clean','Otomatik Temizlik (gun)',s.storage_auto_clean||'30','number','Gecmis davranis uyumlulugu icin korunur.')+
        settingInput('recordings_retention_days','Kayit Saklama Suresi (gun)',s.recordings_retention_days||'30','number','0 verilirse otomatik silme yapilmaz.')+
        settingInput('recordings_keep_latest','Yayin Basina Sakla',s.recordings_keep_latest||'10','number','Her yayinda tutulacak son kayit sayisi.')+
        '<div class="setting-row"><div><div class="setting-label">Otomatik Bakim</div><div class="setting-desc">Kayit, telemetry ve trim bakimlarini periyodik olarak calistirir.</div></div>'+
          '<label class="toggle"><input type="checkbox" class="setting-input" data-key="maintenance_auto_cleanup" '+(s.maintenance_auto_cleanup!=='false'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
        '<button class="btn btn-primary" style="margin-top:8px" onclick="saveSettingsCategory(\'storage\')">Yerel Ayarlari Kaydet</button>'+
      '</div>'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Dis Kopya ve Arsiv Hedefi</div>'+
        '<div class="card" style="margin-bottom:14px;background:var(--bg-primary)"><div class="card-title" style="margin-bottom:10px">Kolay secim rehberi</div><div class="form-hint" style="line-height:1.8"><strong>Yerel:</strong> Kayitlar ayni sunucuda ikinci bir klasore kopyalanir. En kolay secenektir.<br><strong>S3 / MinIO:</strong> Dosyalar bulut benzeri bir depoya gider. Daha duzenli backup icin uygundur.<br><strong>SFTP:</strong> Kayitlari baska bir sunucuya klasor gibi kopyalar. Dusuk butcede en pratik dis hedeflerden biridir.<br><strong>Google Drive / OneDrive:</strong> Yol haritasinda. Bu turda dogrudan entegrasyon yok.</div></div>'+
        '<div class="form-group"><label class="form-label">Arsiv hedefi</label><select id="storage-provider-select" class="form-select setting-input" data-key="archive_provider" onchange="updateStorageProviderUI()"><option value="local" '+((s.archive_provider||'local')==='local'?'selected':'')+'>Bu sunucuda sakla</option><option value="s3" '+((s.archive_provider||'')==='s3'?'selected':'')+'>S3 bulut deposu</option><option value="minio" '+((s.archive_provider||'')==='minio'?'selected':'')+'>MinIO sunucusu</option><option value="sftp" '+((s.archive_provider||'')==='sftp'?'selected':'')+'>SFTP ile baska sunucu</option></select><div class="form-hint">Yerel ile baslamak en kolayidir. Sonra istersen dis hedefe gecebilirsin.</div></div>'+
        '<div id="storage-provider-guide" class="card" style="margin-bottom:14px;background:var(--bg-primary)"></div>'+
        settingInput('archive_prefix','Arsiv klasor adi',s.archive_prefix||'fluxstream','text','Kayitlar ve yedekler bu isim altinda toplanir.')+
        settingInput('archive_public_base_url','Genel erisim link tabani',s.archive_public_base_url||'','text','Varsa panel tiklanabilir baglanti uretir.')+
        settingInput('archive_local_dir','Bu sunucudaki arsiv klasoru',s.archive_local_dir||'','text','Yerel hedef secildiginde dosyalar bu klasore kopyalanir.')+
        settingInput('archive_endpoint','Baglanti adresi',s.archive_endpoint||'','text','Ornek: https://s3.eu-central-1.amazonaws.com')+
        settingInput('archive_region','Bolge',s.archive_region||'us-east-1','text','S3 / MinIO imzalama bolgesi')+
        settingInput('archive_bucket','Depo adi (bucket)',s.archive_bucket||'','text','Dosyalarin yazilacagi depo adi')+
        settingInput('archive_access_key','Kullanici anahtari',s.archive_access_key||'','text','S3 / MinIO erisim anahtari')+
        settingInput('archive_secret_key','Gizli anahtar',s.archive_secret_key||'','password','S3 / MinIO gizli anahtari')+
        settingInput('archive_sftp_host','SFTP sunucu adresi',s.archive_sftp_host||'','text','Host adi veya IP')+
        settingInput('archive_sftp_port','SFTP portu',s.archive_sftp_port||'22','number','Genelde 22')+
        settingInput('archive_sftp_user','Kullanici adi',s.archive_sftp_user||'','text','Sunucuda baglanacak kullanici')+
        settingInput('archive_sftp_remote_dir','Sunucudaki hedef klasor',s.archive_sftp_remote_dir||'','text','Ornek: /srv/fluxstream-archive')+
        settingInput('archive_sftp_key_path','Anahtar dosyasi yolu',s.archive_sftp_key_path||'','text','Bos ise varsayilan SSH anahtari / agent denenir.')+
        '<div class="setting-row"><div><div class="setting-label">SFTP Host Key Kontrolunu Gevset</div><div class="setting-desc">Ilk testte kolaylik saglar; production icin kapali tutmak daha guvenlidir.</div></div>'+
          '<label class="toggle"><input type="checkbox" class="setting-input" data-key="archive_sftp_disable_host_key_check" '+(s.archive_sftp_disable_host_key_check==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
        '<div class="setting-row"><div><div class="setting-label">MinIO uyum modu</div><div class="setting-desc">MinIO ve bazi S3 uyumlu servislerde acik olmali.</div></div>'+
          '<label class="toggle"><input type="checkbox" class="setting-input" data-key="archive_use_path_style" '+(s.archive_use_path_style!=='false'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
        '<div class="card" style="margin-top:14px;background:var(--bg-primary)">'+
          '<div class="card-title" style="margin-bottom:12px">Kayit Arsivi</div>'+
          '<div class="setting-row"><div><div class="setting-label">Kayit arsivi etkin</div><div class="setting-desc">Yerel kayit kutuphanesini bu hedefe tasir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="archive_enabled" '+(s.archive_enabled==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<div class="setting-row"><div><div class="setting-label">Otomatik yukle</div><div class="setting-desc">Yeni kayitlar periyodik olarak arsivlenir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="archive_auto_upload" '+(s.archive_auto_upload==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<div class="setting-row"><div><div class="setting-label">Yukleme sonrasi yereli sil</div><div class="setting-desc">Basarili upload sonrasi diski bosaltir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="archive_delete_local_after_upload" '+(s.archive_delete_local_after_upload==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          settingInput('archive_scan_interval_minutes','Ne kadar sik kontrol edilsin (dk)',s.archive_scan_interval_minutes||'10','number','Yeni kayitlarin ne kadar sik gonderilecegini belirler')+
          settingInput('archive_batch_size','Tek Seferde Maksimum Oge',s.archive_batch_size||'3','number','Bir turda yuklenecek kayit sayisi')+
          '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-top:8px"><button class="btn btn-primary" onclick="saveSettingsCategory(\'storage\')">Hedefi Kaydet</button><button class="btn btn-secondary" onclick="runArchiveSync()">Kayitlari Simdi Gonder</button></div>'+
        '</div>'+
        '<div class="card" style="margin-top:14px;background:var(--bg-primary)">'+
          '<div class="card-title" style="margin-bottom:12px">Sistem Yedegi Arsivi</div>'+
          '<div class="setting-row"><div><div class="setting-label">Yedek arsivi etkin</div><div class="setting-desc">Olusan sistem yedekleri ayni hedefe aktarilir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="backup_archive_enabled" '+(s.backup_archive_enabled==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<div class="setting-row"><div><div class="setting-label">Otomatik yukle</div><div class="setting-desc">Yeni backup dosyalari periyodik olarak yuklenir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="backup_archive_auto_upload" '+(s.backup_archive_auto_upload==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<div class="setting-row"><div><div class="setting-label">Yukleme sonrasi yereli sil</div><div class="setting-desc">Basarili upload sonrasi backup dosyasini yerelden kaldirir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="backup_archive_delete_local_after_upload" '+(s.backup_archive_delete_local_after_upload==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          settingInput('backup_archive_scan_interval_minutes','Ne kadar sik kontrol edilsin (dk)',s.backup_archive_scan_interval_minutes||'30','number','Yeni yedeklerin ne kadar sik gonderilecegini belirler')+
          settingInput('backup_archive_batch_size','Tek Seferde Maksimum Oge',s.backup_archive_batch_size||'2','number','Bir turda yuklenecek backup sayisi')+
          '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-top:8px"><button class="btn btn-primary" onclick="saveSettingsCategory(\'storage\')">Yedek Ayarlarini Kaydet</button><button class="btn btn-secondary" onclick="runBackupArchiveSync()">Yedekleri Simdi Gonder</button></div>'+
        '</div>'+
      '</div>'+
    '</div>'+
    '<div class="card" style="margin-top:16px;margin-bottom:16px"><div class="card-header"><h3 class="card-title">Aktif Kayitlar</h3><span class="form-hint" id="storage-active-count">'+fmtInt(recs.length)+' aktif oturum</span></div><div class="card-body"><table class="table"><thead><tr><th>ID</th><th>Yayin</th><th>Format</th><th>Durum</th><th>Boyut</th><th style="white-space:nowrap">Islem</th></tr></thead><tbody id="rec-list"></tbody></table></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Secili Kayit Onizleme</h3><span class="form-hint">TS / FLV eski dosyalarda gerekirse MP4 donusumu baslatabilirsiniz.</span></div><div class="card-body"><div id="recording-preview-panel"><div class="empty-state"><div class="icon"><i class="bi bi-film"></i></div><h3>Kayit secin</h3><p style="color:var(--text-muted)">Panel secili kaydi ayni sayfada oynatir.</p></div></div></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Kayit Kutuphanesi</h3><span class="form-hint">Yerelde bulunan dosyalar ve izlenebilir kopyalar</span></div><div class="card-body"><table class="table"><thead><tr><th>Yayin</th><th>Dosya</th><th>Format</th><th>Tarih</th><th>Boyut</th><th>Arsiv</th><th>Islem</th></tr></thead><tbody id="saved-rec-list"></tbody></table></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Kayit Arsiv Kutuphanesi</h3><span class="form-hint">Object storage, MinIO veya SFTP hedefindeki kayitlar</span></div><div class="card-body"><table class="table"><thead><tr><th>Yayin</th><th>Dosya</th><th>Saglayici</th><th>Tarih</th><th>Yerel Durum</th><th>Sonuc</th><th>Islem</th></tr></thead><tbody id="archive-rec-list"></tbody></table></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Sistem Yedekleri</h3><span class="form-hint">Restore komutu: '+escHtml(restoreCmd)+'</span></div><div class="card-body"><table class="table"><thead><tr><th>Dosya</th><th>Boyut</th><th>Tarih</th><th>Tur</th><th>Arsiv</th><th>Islem</th></tr></thead><tbody id="system-backup-list"></tbody></table></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Yedek Arsiv Kutuphanesi</h3><span class="form-hint">Harici hedefte saklanan sistem yedekleri</span></div><div class="card-body"><table class="table"><thead><tr><th>Dosya</th><th>Saglayici</th><th>Tarih</th><th>Yerel Durum</th><th>Sonuc</th><th>Islem</th></tr></thead><tbody id="backup-archive-list"></tbody></table></div></div>'+
    '<div id="rec-modal" style="display:none"></div>';

  const rl=document.getElementById('rec-list');
  if(rl){
    rl.innerHTML=recs.length?recs.map(function(r){
      const recID=String(r.ID||'');
      const streamKey=String(r.StreamKey||'');
      const shortID=recID.length>28?recID.slice(0,28)+'…':recID;
      const shortStream=streamKey.length>22?streamKey.slice(0,22)+'…':streamKey;
      return '<tr>'+
        '<td><code title="'+escHtml(recID)+'" style="display:inline-block;max-width:260px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;vertical-align:bottom">'+escHtml(shortID)+'</code></td>'+
        '<td><code title="'+escHtml(streamKey)+'" style="display:inline-block;max-width:220px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;vertical-align:bottom">'+escHtml(shortStream)+'</code></td>'+
        '<td style="white-space:nowrap">'+escHtml(String(r.Format||'').toUpperCase())+'</td>'+
        '<td style="white-space:nowrap"><span class="badge badge-'+(r.Status==='recording'?'green':(r.Status==='error'?'red':'gray'))+'">'+escHtml(String(r.Status||'-'))+'</span></td>'+
        '<td style="white-space:nowrap">'+fmtBytes(r.Size||0)+'</td>'+
        '<td style="white-space:nowrap">'+(r.Status==='recording'?'<button class="btn btn-sm btn-danger" onclick="stopRec(\''+r.ID+'\')">Durdur</button>':'—')+'</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Aktif kayit yok</td></tr>';
  }

  const srl=document.getElementById('saved-rec-list');
  if(srl){
    srl.innerHTML=saved.length?saved.map(function(r){
      const archiveInfo=archiveMap[r.stream_key+'::'+r.name];
      const archiveBadge=archiveInfo?renderArchiveStatusBadge(archiveInfo):'<span class="tag tag-blue">Yerelde</span>';
      const format=String(r.format||'').toLowerCase();
      const canRemux=format==='ts'||format==='flv'||format==='mkv';
      return '<tr>'+
        '<td><code>'+escHtml(r.stream_key)+'</code></td>'+
        '<td>'+escHtml(r.name)+'</td>'+
        '<td>'+(r.format||'-').toUpperCase()+'</td>'+
        '<td>'+fmtLocaleDateTime(r.mod_time)+'</td>'+
        '<td>'+fmtBytes(r.size||0)+'</td>'+
        '<td>'+archiveBadge+'</td>'+
        '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
          '<button class="btn btn-sm btn-secondary" onclick=\'previewRecordingPanel('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+','+JSON.stringify(r.format||'')+','+JSON.stringify(r.mod_time||'')+','+(r.size||0)+')\'>Onizle</button>'+
          '<button class="btn btn-sm btn-secondary" onclick=\'downloadRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Indir</button>'+
          (canRemux?'<button class="btn btn-sm btn-secondary" onclick=\'remuxRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+','+JSON.stringify('mp4')+')\'>MP4 Hazirla</button>':'')+
          (archiveEnabled?'<button class="btn btn-sm btn-secondary" onclick=\'archiveRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>'+(archiveInfo&&archiveInfo.status==='archived'?'Yeniden Arsivle':'Arsive Gonder')+'</button>':'')+
          (archiveInfo&&archiveInfo.object_url?'<button class="btn btn-sm btn-secondary" onclick=\'window.open('+JSON.stringify(archiveInfo.object_url)+',"_blank")\'>Arsiv Linki</button>':'')+
          '<button class="btn btn-sm btn-danger" onclick=\'deleteRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Sil</button>'+
        '</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="7" style="text-align:center;color:var(--text-muted);padding:24px">Kaydedilmis dosya yok</td></tr>';
  }

  const arl=document.getElementById('archive-rec-list');
  if(arl){
    arl.innerHTML=archives.length?archives.map(function(item){
      const localState=item.local_deleted?'<span class="tag tag-yellow">Yerelde yok</span>':'<span class="tag tag-green">Yerelde var</span>';
      const statusBadge=renderArchiveStatusBadge(item);
      return '<tr>'+
        '<td><code>'+escHtml(item.stream_key)+'</code></td>'+
        '<td>'+escHtml(item.filename)+'</td>'+
        '<td>'+escHtml(String(item.provider||'-').toUpperCase())+'</td>'+
        '<td>'+fmtLocaleDateTime(item.archived_at||item.updated_at||item.created_at)+'</td>'+
        '<td>'+localState+'</td>'+
        '<td>'+statusBadge+(item.last_error?'<div class="setting-desc" style="max-width:320px">'+escHtml(item.last_error)+'</div>':'')+'</td>'+
        '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
          '<button class="btn btn-sm btn-secondary" onclick=\'restoreRecordingArchive('+JSON.stringify(item.stream_key)+','+JSON.stringify(item.filename)+')\'>Geri Yukle</button>'+
          (item.object_url?'<button class="btn btn-sm btn-secondary" onclick=\'window.open('+JSON.stringify(item.object_url)+',"_blank")\'>Arsiv Linki</button>':'')+
        '</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="7" style="text-align:center;color:var(--text-muted);padding:24px">Arsivlenmis kayit yok</td></tr>';
  }

  const bl=document.getElementById('system-backup-list');
  if(bl){
    bl.innerHTML=backups.length?backups.map(function(item){
      const archiveInfo=backupArchiveMap[item.name];
      const archiveBadge=archiveInfo?renderBackupArchiveStatusBadge(archiveInfo):'<span class="tag tag-blue">Yerelde</span>';
      return '<tr>'+
        '<td class="mono-wrap">'+escHtml(item.name)+'</td>'+
        '<td>'+formatBytes(item.size||0)+'</td>'+
        '<td>'+escHtml(fmtLocaleDateTime(item.mod_time))+'</td>'+
        '<td>'+(item.include_recordings?'<span class="tag tag-blue">Kayitlar dahil</span>':'<span class="tag tag-green">Hafif</span>')+'</td>'+
        '<td>'+archiveBadge+'</td>'+
        '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
          '<a class="btn btn-sm btn-secondary" href="/api/system/backups/download/'+encodeURIComponent(item.name)+'" target="_blank" rel="noopener">Indir</a>'+
          (backupArchiveEnabled?'<button class="btn btn-sm btn-secondary" onclick=\'archiveSystemBackup('+JSON.stringify(item.name)+')\'>'+(archiveInfo&&archiveInfo.status==='archived'?'Yeniden Arsivle':'Arsive Gonder')+'</button>':'')+
          '<button class="btn btn-sm btn-danger" onclick=\'deleteSystemBackup('+JSON.stringify(item.name)+')\'>Sil</button>'+
        '</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Yerel sistem yedegi yok</td></tr>';
  }

  const bal=document.getElementById('backup-archive-list');
  if(bal){
    bal.innerHTML=backupArchives.length?backupArchives.map(function(item){
      const localState=item.local_deleted?'<span class="tag tag-yellow">Yerelde yok</span>':'<span class="tag tag-green">Yerelde var</span>';
      const statusBadge=renderBackupArchiveStatusBadge(item);
      return '<tr>'+
        '<td class="mono-wrap">'+escHtml(item.name)+'</td>'+
        '<td>'+escHtml(String(item.provider||'-').toUpperCase())+'</td>'+
        '<td>'+fmtLocaleDateTime(item.archived_at||item.updated_at||item.created_at)+'</td>'+
        '<td>'+localState+'</td>'+
        '<td>'+statusBadge+(item.last_error?'<div class="setting-desc" style="max-width:320px">'+escHtml(item.last_error)+'</div>':'')+'</td>'+
        '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
          '<button class="btn btn-sm btn-secondary" onclick=\'restoreSystemBackupArchive('+JSON.stringify(item.name)+')\'>Geri Getir</button>'+
          (item.object_url?'<button class="btn btn-sm btn-secondary" onclick=\'window.open('+JSON.stringify(item.object_url)+',"_blank")\'>Arsiv Linki</button>':'')+
        '</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Arsivlenmis sistem yedegi yok</td></tr>';
  }

  window._recStreams=streams;
  window._savedRecordings=saved;
  window._recordingArchives=archives;
  window._systemBackups=backups;
  window._backupArchives=backupArchives;
  window._recordingPreviewSelection=null;
  resetRecordingPreviewPanel();
  applyStorageSnapshot(normalizeStorageSnapshot(s,report,archivesRes,recsRes,streamsRes,savedRes,backupsRes,backupArchivesRes,remuxJobsRes),{resetPreview:true});
  updateStorageProviderUI();
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETTINGS - TRANSCODE ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderSettingsTranscode(c){
  const s=await api('/api/settings');
  const gpuMode=(s.gpu_accel||'none').toUpperCase();
  const gpuIcons={NONE:'bi-cpu',NVENC:'bi-gpu-card',QSV:'bi-gpu-card',AMF:'bi-gpu-card'};
  c.innerHTML=
    '<div class="studio-page">'+
      '<div class="studio-hero">'+
        '<h1 class="studio-hero-title"><i class="bi bi-cpu" style="margin-right:10px"></i>Transkod / FFmpeg</h1>'+
        '<p class="studio-hero-sub">Canli transkod ve ABR katman uretimi icin FFmpeg yapilandirmasi. GPU hizlandirma aktifse performans onemli olcude artar.</p>'+
        '<div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active"><i class="bi '+(gpuIcons[gpuMode]||'bi-cpu')+'"></i> GPU: '+gpuMode+'</span></div>'+
      '</div>'+
      '<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span>FFmpeg Ayarlari</span></div>'+
          settingInput('ffmpeg_path','FFmpeg Yolu',s.ffmpeg_path||'ffmpeg','text','FFmpeg calistirilabilir dosya yolu. Sistem PATH\'inde ise sadece "ffmpeg" yeterli.')+
          '<div class="form-group"><label class="form-label">GPU Hizlandirma</label>'+
            '<select class="form-select setting-input" data-key="gpu_accel">'+
              '<option value="none" '+(s.gpu_accel==='none'||!s.gpu_accel?'selected':'')+'>Yok (CPU)</option>'+
              '<option value="nvenc" '+(s.gpu_accel==='nvenc'?'selected':'')+'>NVIDIA NVENC</option>'+
              '<option value="qsv" '+(s.gpu_accel==='qsv'?'selected':'')+'>Intel Quick Sync</option>'+
              '<option value="amf" '+(s.gpu_accel==='amf'?'selected':'')+'>AMD AMF</option>'+
            '</select>'+
            '<div class="form-hint">GPU varsa encoding islemleri donanim uzerinden yapilir. CPU yukunu onemli olcude azaltir.</div>'+
          '</div>'+
          '<button class="btn btn-primary" style="margin-top:8px" onclick="saveSettingsCategory(\'transcode\')"><i class="bi bi-check-lg"></i> Kaydet</button>'+
        '</div>'+
        '<div class="studio-card">'+
          '<div class="studio-section-title"><span>Oneri ve Bilgilendirme</span></div>'+
          '<div class="studio-alert info" style="margin-bottom:10px"><i class="bi bi-lightbulb" style="margin-right:6px"></i><strong>Onerilen:</strong> ffmpeg 5.x veya 6.x, GPU mevcut ise nvenc secimi CPU yukunu %80 azaltir.</div>'+
          '<div class="metric-list">'+
            '<div class="metric-row"><span>Aktif GPU modu</span><strong>'+gpuMode+'</strong></div>'+
            '<div class="metric-row"><span>Mevcut yol</span><span class="mono-wrap" style="max-width:300px">'+escHtml(s.ffmpeg_path||'ffmpeg')+'</span></div>'+
          '</div>'+
          '<div style="margin-top:14px"><button class="btn btn-sm btn-secondary" onclick="navigate(\'transcode-jobs\')"><i class="bi bi-list-task"></i> Transcode Islerini Gor</button></div>'+
        '</div>'+
      '</div>'+
    '</div>';
}

async function renderGuidedSettings(c){
  const [s,health]=await Promise.all([api('/api/settings'),api('/api/health/report')]);
  const alerts=Array.isArray(health&&health.alerts)?health.alerts:[];
  const publicConfig=getPublicBase(s);
  const healthStatus=String((health&&health.status)||'ok').toUpperCase();
  const currentPreset=s.abr_profile_set||'balanced';
  c.innerHTML=
    '<div class="studio-page"><div class="studio-hero"><h1 class="studio-hero-title"><i class="bi bi-magic" style="margin-right:10px"></i>Hizli Ayarlar</h1><p class="studio-hero-sub">Ilk kurulum veya gunluk kullanim icin en sik ihtiyac duyulan ayarlar tek ekranda. Her blogu bagimsiz kaydedebilirsiniz.</p><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active"><i class="bi bi-speedometer2"></i> Profil: '+escHtml(currentPreset.charAt(0).toUpperCase()+currentPreset.slice(1))+'</span><span class="studio-pill"><i class="bi bi-heart-pulse-fill"></i> '+healthStatus+'</span><span class="studio-pill"><i class="bi bi-bell"></i> '+fmtInt(alerts.length)+' Uyari</span></div></div>'+
    '<div class="studio-card"><div class="studio-section-title"><span><i class="bi bi-broadcast-pin" style="margin-right:8px;color:var(--accent)"></i>1. Yayin Profili</span><span class="tag tag-blue">Tek tikla uygula</span></div>'+
      '<div class="studio-section-sub" style="margin-bottom:14px">Hazir profiller onerilen en iyi ayar kombinasyonlarini icinde tasir.</div>'+
      '<div class="studio-option-grid">'+
        '<div class="studio-option-card '+(currentPreset==='balanced'?'active':'')+'" onclick="applyDeliveryPreset(\'balanced\')"><div style="font-size:28px;margin-bottom:8px"><i class="bi bi-broadcast-pin"></i></div><div class="studio-option-title">TV / Dengeli</div><div class="studio-option-meta">Canli TV, YouTube benzeri. Adaptif kalite merdiveni.</div></div>'+
        '<div class="studio-option-card '+(currentPreset==='mobile'?'active':'')+'" onclick="applyDeliveryPreset(\'mobile\')"><div style="font-size:28px;margin-bottom:8px"><i class="bi bi-phone"></i></div><div class="studio-option-title">Mobil / Hafif</div><div class="studio-option-meta">Zayif baglanti ve dusuk CPU icin optimize.</div></div>'+
        '<div class="studio-option-card '+(currentPreset==='resilient'?'active':'')+'" onclick="applyDeliveryPreset(\'resilient\')"><div style="font-size:28px;margin-bottom:8px"><i class="bi bi-wifi-off"></i></div><div class="studio-option-title">Dayanikli</div><div class="studio-option-meta">Kesinti riskini azaltir ve bitrate sinirlar.</div></div>'+
        '<div class="studio-option-card '+(currentPreset==='radio'?'active':'')+'" onclick="applyDeliveryPreset(\'radio\')"><div style="font-size:28px;margin-bottom:8px"><i class="bi bi-mic-fill"></i></div><div class="studio-option-title">Radyo / Audio</div><div class="studio-option-meta">Ses yayinlari icin sade ayar.</div></div>'+
      '</div></div>'+
      '<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card"><div class="studio-section-title"><span><i class="bi bi-globe" style="margin-right:8px;color:var(--accent)"></i>2. Public Erisim</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:12px">Embed ve player linkleri bu ayarlara gore uretilir.</div>'+
          settingInput('embed_domain','Public Domain / IP',s.embed_domain||'','text','Bos birakirsaniz mevcut host kullanilir.')+
          '<div class="studio-form-grid">'+settingInput('embed_http_port','HTTP Port',s.embed_http_port||s.http_port||'8844','number','')+settingInput('embed_https_port','HTTPS Port',s.embed_https_port||s.https_port||'443','number','')+'</div>'+
          '<div class="setting-row"><div><div class="setting-label">HTTPS Link Uret</div><div class="setting-desc">SSL etkin ve sertifika hazirsa linkleri HTTPS yapar.</div></div><label class="toggle"><input type="checkbox" class="guided-input" data-key="embed_use_https" '+(shouldUsePublicHTTPS(s)?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<div class="studio-alert info" style="margin-top:10px"><i class="bi bi-link-45deg" style="margin-right:6px"></i>Ornek: <code>'+escHtml(publicConfig.base)+'</code></div>'+
          '<div style="margin-top:12px"><button class="btn btn-primary" onclick="saveGuidedPublic()"><i class="bi bi-check-lg"></i> Kaydet</button></div>'+
        '</div>'+
        '<div class="studio-card"><div class="studio-section-title"><span><i class="bi bi-hdd" style="margin-right:8px;color:var(--accent)"></i>3. Kayit ve Temizlik</span></div>'+
          '<div class="studio-section-sub" style="margin-bottom:12px">Kayitlarin ne kadar saklanacagini ve otomatik bakimi ayarlayin.</div>'+
          '<div class="form-group"><label class="form-label">Kayit Saklama (gun)</label><input class="form-input guided-input" data-key="recordings_retention_days" type="number" value="'+escHtml(s.recordings_retention_days||'30')+'"><div class="form-hint">0 = sinirsiz. Onerilen: 30</div></div>'+
          '<div class="form-group"><label class="form-label">Yayin Basina Maks Kayit</label><input class="form-input guided-input" data-key="recordings_keep_latest" type="number" value="'+escHtml(s.recordings_keep_latest||'10')+'"><div class="form-hint">Eski kayitlar otomatik budanir.</div></div>'+
          '<div class="setting-row"><div><div class="setting-label">Otomatik Bakim</div><div class="setting-desc">Kayit ve analytics temizligi periyodik calisir.</div></div><label class="toggle"><input type="checkbox" class="guided-input" data-key="maintenance_auto_cleanup" '+(s.maintenance_auto_cleanup==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<div style="display:flex;gap:10px;margin-top:12px"><button class="btn btn-primary" onclick="saveGuidedStorage()"><i class="bi bi-check-lg"></i> Kaydet</button><button class="btn btn-secondary" onclick="runMaintenance()"><i class="bi bi-arrow-repeat"></i> Bakimi Calistir</button></div>'+
        '</div>'+
      '</div>'+
      '<div class="studio-grid studio-grid-2">'+
        '<div class="studio-card"><div class="studio-section-title"><span><i class="bi bi-shield-lock" style="margin-right:8px;color:var(--accent)"></i>4. Guvenlik</span></div>'+
          '<div class="setting-row"><div><div class="setting-label">Playback Token Zorunlu</div><div class="setting-desc">Acilirsa izleyicilerde gecerli token aranir.</div></div><label class="toggle"><input type="checkbox" class="guided-input" data-key="token_enabled" '+(s.token_enabled==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<div class="form-group" style="margin-top:16px"><label class="form-label">Rate Limit (istek/sn)</label><input class="form-input guided-input" data-key="rate_limit" type="number" value="'+escHtml(s.rate_limit||'100')+'"><div class="form-hint">Onerilen: 100</div></div>'+
          '<div style="margin-top:12px"><button class="btn btn-primary" onclick="saveGuidedSecurity()"><i class="bi bi-check-lg"></i> Kaydet</button></div>'+
        '</div>'+
        '<div class="studio-card"><div class="studio-section-title"><span><i class="bi bi-heart-pulse" style="margin-right:8px;color:'+(healthStatus==='OK'?'var(--success)':'var(--warning)')+'"></i>5. Sistem Durumu</span></div>'+
          '<div class="metric-list" style="margin-top:8px"><div class="metric-row"><span>Durum</span><strong><span class="tag tag-'+(healthStatus==='OK'?'green':'yellow')+'">'+healthStatus+'</span></strong></div><div class="metric-row"><span>Uyarilar</span><strong>'+fmtInt(alerts.length)+'</strong></div><div class="metric-row"><span>ABR</span><strong><span class="tag '+(s.abr_enabled==='true'?'tag-green':'tag-yellow')+'">'+(s.abr_enabled==='true'?'Acik':'Kapali')+'</span></strong></div><div class="metric-row"><span>Analitik</span><strong>'+(s.analytics_persist_enabled==='true'?'Acik':'Kapali')+'</strong></div></div>'+
          (alerts.length?'<div class="studio-alert warning" style="margin-top:12px"><strong>'+escHtml(alerts[0].title||'Uyari')+'</strong></div>':'')+
          '<div style="display:flex;gap:10px;margin-top:14px"><button class="btn btn-secondary" onclick="navigate(\'settings-health\')"><i class="bi bi-heart-pulse"></i> Saglik</button><button class="btn btn-secondary" onclick="openSystemControl()"><i class="bi bi-gear"></i> Sunucu</button></div>'+
        '</div>'+
      '</div>'+
    '</div>';
}

function presetCard(title,desc,preset,icon){
  return '<div class="card" style="cursor:pointer" onclick="applyDeliveryPreset(\''+preset+'\')"><div class="card-title" style="margin-bottom:8px"><i class="bi '+icon+' title-icon"></i>'+title+'</div><div class="form-hint" style="line-height:1.7">'+desc+'</div></div>';
}

async function applyDeliveryPreset(preset){
  const updates={abr_enabled:'true',abr_master_enabled:'true',abr_profile_set:preset,hls_enabled:'true',transcode_live_hls_enabled:'true'};
  if(preset==='radio'){
    updates.abr_enabled='false';
    updates.mp3_enabled='true';
    updates.aac_out_enabled='true';
  }else if(preset==='mobile' || preset==='resilient'){
    updates.hls_ll_enabled='true';
  }else{
    updates.dash_enabled='true';
  }
  await saveSettingsValues('outputs',updates,true);
  toast('Teslimat profili uygulandi');
  if(currentPage==='guided-settings')loadPage('guided-settings');
}

async function renderSettingsABR(c){
  const s=await api('/api/settings');
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Teslimat / ABR</h1><div style="color:var(--text-muted);font-size:13px">Adaptif bitrate, izleyicinin baglantisina gore kaliteyi otomatik yukselten veya dusuren HLS master playlist uretilmesini saglar.</div></div>'+
    '<div class="card" style="max-width:880px;margin-bottom:16px">'+
      '<div class="setting-row"><div><div class="setting-label">Adaptif Bitrate (ABR)</div><div class="setting-desc">Acik oldugunda canli HLS icin coklu kalite katmanlari ve master.m3u8 uretilir.</div></div>'+
      '<label class="toggle"><input type="checkbox" class="setting-input" data-key="abr_enabled" '+(s.abr_enabled==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      '<div class="setting-row"><div><div class="setting-label">Master Playlist Uret</div><div class="setting-desc">Player once master playlist arar. Yoksa index.m3u8 fallback kullanilir.</div></div>'+
      '<label class="toggle"><input type="checkbox" class="setting-input" data-key="abr_master_enabled" '+(s.abr_master_enabled!=='false'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      '<div class="form-group" style="margin-top:16px"><label class="form-label">Hazir Profil Seti</label><select class="form-select setting-input" data-key="abr_profile_set"><option value="balanced" '+((s.abr_profile_set||'balanced')==='balanced'?'selected':'')+'>TV / Dengeli</option><option value="mobile" '+((s.abr_profile_set||'')==='mobile'?'selected':'')+'>Mobil / Hafif</option><option value="resilient" '+((s.abr_profile_set||'')==='resilient'?'selected':'')+'>Dusuk Bant / Dayanikli</option><option value="radio" '+((s.abr_profile_set||'')==='radio'?'selected':'')+'>Radyo / Audio</option></select><div class="form-hint">Dengeli cogu video yayin icin en iyi baslangic noktasidir.</div></div>'+
      '<div class="form-group"><label class="form-label">ABR Profil JSON</label><textarea class="form-textarea setting-input" data-key="abr_profiles_json" style="min-height:220px">'+escHtml(s.abr_profiles_json||'')+'</textarea><div class="form-hint">Gelistirilmis kullanim icin. Hazir setler yukaridan secilebilir.</div></div>'+
      '<button class="btn btn-primary" onclick="saveSettingsCategory(\'outputs\')">ABR Ayarlarini Kaydet</button>'+
    '</div>';
}

async function renderSettingsHealth(c){
  const [s,report]=await Promise.all([api('/api/settings'),api('/api/health/report')]);
  const alerts=Array.isArray(report&&report.alerts)?report.alerts:[];
  const qoeStreams=(Array.isArray(report&&report.qoe_streams)?report.qoe_streams:[]).slice().sort(function(a,b){
    return (Number(b.qoe_alert_count||0)-Number(a.qoe_alert_count||0))||
      (Number(b.total_stalls||0)-Number(a.total_stalls||0))||
      (Number(b.total_quality_transitions||0)-Number(a.total_quality_transitions||0));
  });
  const archiveSummary=report&&report.storage&&report.storage.archive?report.storage.archive:{};
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Saglik ve Uyarilar</h1><div style="color:var(--text-muted);font-size:13px">Sunucu sagligi, sertifika, bellek, depolama ve bakim isleri bu ekrandan izlenir.</div></div>'+
    '<div class="card-grid card-grid-4" style="margin-bottom:16px">'+
      statCard('green','bi-heart-pulse-fill',String((report&&report.status)||'ok').toUpperCase(),'Genel Durum')+
      statCard('orange','bi-bell-fill',fmtInt(alerts.length),'Aktif Uyari')+
      statCard('blue','bi-database-fill',formatBytes((report&&report.storage&&report.storage.recordings_bytes)||0),'Kayit Depolama')+
      statCard('purple','bi-clock-history',fmtInt(((report&&report.snapshots)||[]).length),'Kalici Snapshot')+
    '</div>'+
    '<div class="card-grid card-grid-2">'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Uyarilar</div>'+
        (alerts.length?alerts.map(function(item){
          var tone=item.level==='critical'?'tag-red':item.level==='warning'?'tag-yellow':'tag-blue';
          return '<div style="padding:12px 0;border-bottom:1px solid var(--border)"><div style="display:flex;justify-content:space-between;gap:12px"><strong>'+escHtml(item.title||item.code||'Uyari')+'</strong><span class="tag '+tone+'">'+escHtml(String(item.level||'info').toUpperCase())+'</span></div><div class="form-hint" style="margin-top:6px">'+escHtml(item.description||'')+'</div>'+(item.action?'<div class="form-hint" style="margin-top:6px;color:var(--text-secondary)">'+escHtml(item.action)+'</div>':'')+'</div>';
        }).join(''):'<div style="color:var(--text-muted)">Aktif uyari yok.</div>')+
      '</div>'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Esikler</div>'+
        '<div class="form-group"><label class="form-label">Depolama Uyari Esigi (GB)</label><input class="form-input setting-input" data-key="alerts_disk_threshold_gb" type="number" value="'+escHtml(s.alerts_disk_threshold_gb||'5')+'"><div class="form-hint">Toplam depolama limitine bu kadar kala uyari olusur.</div></div>'+
        '<div class="form-group"><label class="form-label">Bellek Uyari Esigi (MB)</label><input class="form-input setting-input" data-key="alerts_memory_threshold_mb" type="number" value="'+escHtml(s.alerts_memory_threshold_mb||'2048')+'"><div class="form-hint">Asilirsa panel warning uretir.</div></div>'+
        '<div class="form-group"><label class="form-label">Sertifika Uyari Esigi (gun)</label><input class="form-input setting-input" data-key="alerts_cert_days" type="number" value="'+escHtml(s.alerts_cert_days||'21')+'"><div class="form-hint">Bu sureye girince sertifika yenileme uyarisi verilir.</div></div>'+
        '<div class="form-group"><label class="form-label">Kalite Gecisi Uyari Carpani</label><input class="form-input setting-input" data-key="alerts_qoe_transition_ratio_threshold" type="number" value="'+escHtml(s.alerts_qoe_transition_ratio_threshold||'4')+'"><div class="form-hint">Aktif player basina kabul edilen kalite gecisi carpani.</div></div>'+
        '<div class="form-group"><label class="form-label">Ses Gecisi Uyari Carpani</label><input class="form-input setting-input" data-key="alerts_qoe_audio_ratio_threshold" type="number" value="'+escHtml(s.alerts_qoe_audio_ratio_threshold||'3')+'"><div class="form-hint">Aktif player basina kabul edilen ses izi degisimi carpani.</div></div>'+
        '<div style="display:flex;gap:10px"><button class="btn btn-primary" onclick="saveSettingsCategory(\'health\')">Esikleri Kaydet</button><button class="btn btn-secondary" onclick="runMaintenance()">Bakimi Simdi Calistir</button></div>'+
      '</div>'+
    '</div>'+
    '<div class="card-grid card-grid-2" style="margin-top:16px">'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">QoE Riskli Yayinlar</div>'+
        (qoeStreams.length
          ?'<table class="table"><thead><tr><th>Yayin</th><th>Player</th><th>Stall</th><th>Kalite</th><th>Ses</th><th>Baskin</th><th>Durum</th></tr></thead><tbody>'+
            qoeStreams.map(function(item){
              const tone=Number(item.qoe_alert_count||0)>0?'tag-red':'tag-blue';
              return '<tr>'+
                '<td><div style="font-weight:600">'+escHtml(item.stream_name||item.stream_key||'-')+'</div><div class="setting-desc"><code>'+escHtml(item.stream_key||'-')+'</code></div></td>'+
                '<td>'+fmtInt(item.active_sessions||0)+'</td>'+
                '<td>'+fmtInt(item.total_stalls||0)+'</td>'+
                '<td>'+fmtInt(item.total_quality_transitions||0)+'</td>'+
                '<td>'+fmtInt(item.total_audio_switches||0)+'</td>'+
                '<td><div class="setting-desc">'+escHtml(item.dominant_quality||'-')+'</div><div class="setting-desc">'+escHtml(item.dominant_audio||'-')+'</div></td>'+
                '<td><span class="tag '+tone+'">'+fmtInt(item.qoe_alert_count||0)+' uyari</span></td>'+
              '</tr>';
            }).join('')+
           '</tbody></table>'
          :'<div style="color:var(--text-muted)">Canli QoE riski gorunen yayin yok.</div>')+
      '</div>'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Arsiv Ozeti</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Arsiv etkin</span><strong>'+(archiveSummary&&archiveSummary.enabled?'Evet':'Hayir')+'</strong></div>'+
          '<div class="metric-row"><span>Saglayici</span><strong>'+escHtml(String(archiveSummary&&archiveSummary.provider||'kapali').toUpperCase())+'</strong></div>'+
          '<div class="metric-row"><span>Arsivlenen oge</span><strong>'+fmtInt(archiveSummary&&archiveSummary.items||0)+'</strong></div>'+
          '<div class="metric-row"><span>Hata durumundaki oge</span><strong>'+fmtInt(archiveSummary&&archiveSummary.error_items||0)+'</strong></div>'+
          '<div class="metric-row"><span>Yerelden silinmis oge</span><strong>'+fmtInt(archiveSummary&&archiveSummary.local_deleted_items||0)+'</strong></div>'+
          '<div class="metric-row"><span>Yedek arsiv ogesi</span><strong>'+fmtInt(archiveSummary&&archiveSummary.backup_items||0)+'</strong></div>'+
          '<div class="metric-row"><span>Yedek arsiv hatasi</span><strong>'+fmtInt(archiveSummary&&archiveSummary.backup_error_items||0)+'</strong></div>'+
          '<div class="metric-row"><span>Son senkron</span><strong>'+escHtml(archiveSummary&&archiveSummary.last_sync_at?fmtLocaleDateTime(archiveSummary.last_sync_at):'-')+'</strong></div>'+
        '</div>'+
        '<div class="form-hint" style="margin-top:10px;line-height:1.7">'+escHtml((archiveSummary&&archiveSummary.last_error)||'Object storage akisinda yeni hata gorunmuyor.')+'</div>'+
      '</div>'+
    '</div>';
}

async function renderDiagnostics(c){
  const streams=await api('/api/streams')||[];
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Teshis</h1><div style="color:var(--text-muted);font-size:13px">Bir yayinin HLS, DASH, kayit ve ABR ciktilarinin dosya seviyesinde hazir olup olmadigini hizli kontrol eder.</div></div>'+
    '<div class="card" style="max-width:900px">'+
      '<div class="form-group"><label class="form-label">Yayin Sec</label><select class="form-select" id="diag-stream">'+
        (streams||[]).map(function(st){return '<option value="'+st.id+'">'+escHtml(st.name)+' ('+escHtml(st.stream_key)+')</option>'}).join('')+
      '</select></div>'+
      '<div style="display:flex;gap:10px;margin-bottom:16px"><button class="btn btn-primary" onclick="loadDiagnostics()">Kontrol Et</button><button class="btn btn-secondary" onclick="navigate(\'stream-detail-\'+(document.getElementById(\'diag-stream\')?.value||\'\'))">Yayin Detayini Ac</button></div>'+
      '<div id="diag-output" style="color:var(--text-muted)">Bir yayin secip kontrol edebilirsiniz.</div>'+
    '</div>';
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â SETTINGS HELPERS ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
function settingInput(key,label,value,type,hint){
  return '<div class="form-group"><label class="form-label">'+label+'</label>'+
    '<input class="form-input setting-input" data-key="'+key+'" type="'+(type||'text')+'" value="'+escHtml(String(value||''))+'">'+
    (hint?'<div class="form-hint">'+hint+'</div>':'')+
  '</div>';
}
async function saveSettingsValues(category,updates,silent){
  const res=await api('/api/settings/'+category,{method:'PUT',body:updates});
  if(res&&res.success!==false&&updates&&updates.language){
    runtimeSettings=runtimeSettings||{};
    runtimeSettings.language=updates.language;
    setCurrentLanguage(updates.language,true);
    if(setupCompleted&&authToken){
      renderApp();
      navigate(currentPage);
    }else if(!setupCompleted){
      renderWizard();
    }else{
      renderLogin();
    }
  }
  if(!silent){
    if(res&&res.success!==false)toast(t('Ayarlar kaydedildi!'));
    else toast((res&&res.message)||t('Kayit hatasi'),'error');
  }
  loadProtoStatus();
  return res;
}
async function saveSettingsCategory(category){
  const inputs=document.querySelectorAll('.setting-input');
  const updates={};
  inputs.forEach(el=>{
    const key=el.dataset.key;
    if(!key)return;
    if(el.type==='checkbox')updates[key]=el.checked?'true':'false';
    else updates[key]=el.value;
  });
  await saveSettingsValues(category,updates,false);
}
async function saveSSLSettings(){
  const inputs=document.querySelectorAll('.setting-input');
  const generalUpdates={};
  const protocolUpdates={};
  const sslUpdates={};
  inputs.forEach(el=>{
    const key=el.dataset.key;
    if(!key)return;
    const value=el.type==='checkbox'?(el.checked?'true':'false'):el.value;
    if(key==='https_port')generalUpdates[key]=value;
    else if(key==='rtmps_enabled'||key==='rtmps_port')protocolUpdates[key]=value;
    else sslUpdates[key]=value;
  });
  if(Object.keys(generalUpdates).length)await saveSettingsValues('general',generalUpdates,true);
  if(Object.keys(protocolUpdates).length)await saveSettingsValues('protocols',protocolUpdates,true);
  await saveSettingsValues('ssl',sslUpdates,false);
}
async function saveGuidedPublic(){
  const httpsToggle=document.querySelector('.guided-input[data-key="embed_use_https"]');
  await saveSettingsValues('embed',{
    embed_domain:document.querySelector('.setting-input[data-key="embed_domain"]')?.value||'',
    embed_http_port:document.querySelector('.setting-input[data-key="embed_http_port"]')?.value||'',
    embed_https_port:document.querySelector('.setting-input[data-key="embed_https_port"]')?.value||'',
    embed_use_https:httpsToggle&&httpsToggle.checked?'true':'false'
  },false);
}
async function saveGeneralSettingsExtended(){
  await saveSettingsValues('general',{
    server_name:document.querySelector('.setting-input[data-key="server_name"]')?.value||'FluxStream',
    language:document.querySelector('.setting-input[data-key="language"]')?.value||'tr',
    timezone:document.querySelector('.setting-input[data-key="timezone"]')?.value||'Europe/Istanbul',
    theme:document.querySelector('.setting-input[data-key="theme"]')?.value||'light',
    guided_mode_enabled:document.querySelector('.setting-input[data-key="guided_mode_enabled"]')?.checked?'true':'false',
    http_port:document.querySelector('.setting-input[data-key="http_port"]')?.value||'8844',
    https_port:document.querySelector('.setting-input[data-key="https_port"]')?.value||'443'
  },true);
  await saveSettingsValues('embed',{
    embed_domain:document.querySelector('.setting-input[data-key="embed_domain"]')?.value||'',
    embed_http_port:document.querySelector('.setting-input[data-key="embed_http_port"]')?.value||'8844',
    embed_https_port:document.querySelector('.setting-input[data-key="embed_https_port"]')?.value||'443'
  },true);
  await saveSettingsValues('outputs',{
    player_quality_selector:document.querySelector('.setting-input[data-key="player_quality_selector"]')?.checked?'true':'false'
  },true);
  await saveSettingsValues('storage',{
    maintenance_auto_cleanup:document.querySelector('.setting-input[data-key="maintenance_auto_cleanup"]')?.checked?'true':'false',
    recordings_retention_days:document.querySelector('.setting-input[data-key="recordings_retention_days"]')?.value||'30'
  },false);
}
async function saveGuidedStorage(){
  const updates={};
  document.querySelectorAll('.guided-input').forEach(function(el){
    const key=el.dataset.key;
    if(!key)return;
    if(['recordings_retention_days','recordings_keep_latest','maintenance_auto_cleanup'].indexOf(key)===-1)return;
    updates[key]=el.type==='checkbox'?(el.checked?'true':'false'):el.value;
  });
  await saveSettingsValues('storage',updates,false);
}
async function saveGuidedSecurity(){
  const updates={};
  document.querySelectorAll('.guided-input').forEach(function(el){
    const key=el.dataset.key;
    if(!key)return;
    if(['token_enabled','rate_limit'].indexOf(key)===-1)return;
    updates[key]=el.type==='checkbox'?(el.checked?'true':'false'):el.value;
  });
  await saveSettingsValues('security',updates,false);
}
async function runMaintenance(){
  const res=await api('/api/maintenance/run',{method:'POST'});
  if(res&&res.success){
    toast('Bakim tamamlandi');
    if(currentPage==='settings-health'||currentPage==='guided-settings')loadPage(currentPage);
  }else{
    toast((res&&res.message)||'Bakim basarisiz','error');
  }
}
async function runArchiveSync(showToast=true){
  const res=await api('/api/recordings/archive/sync',{method:'POST'});
  if(res&&res.success){
    if(showToast)toast('Arsiv senkronu tamamlandi');
    if(currentPage==='recordings'||currentPage==='settings-storage')await refreshStorageSnapshot({resetPreview:false});
    else if(currentPage==='maintenance-center')await loadPage(currentPage);
  }else{
    toast((res&&res.message)||'Arsiv senkronu basarisiz','error');
  }
}
async function runBackupArchiveSync(showToast=true){
  const res=await api('/api/system/backups/archive/sync',{method:'POST'});
  if(res&&res.success){
    if(showToast)toast('Yedek arsiv senkronu tamamlandi');
    if(currentPage==='recordings'||currentPage==='settings-storage')await refreshStorageSnapshot({resetPreview:false});
    else if(currentPage==='maintenance-center')await loadPage(currentPage);
  }else{
    toast((res&&res.message)||'Yedek arsiv senkronu basarisiz','error');
  }
}
async function loadDiagnostics(){
  const id=document.getElementById('diag-stream')?.value;
  const out=document.getElementById('diag-output');
  if(!id||!out){return}
  out.innerHTML='Kontrol ediliyor...';
  const data=await api('/api/diagnostics/stream/'+id);
  if(!data||data.error){
    out.innerHTML='<div style="color:var(--danger)">Teshis verisi alinamadi.</div>';
    return;
  }
  const checks=Array.isArray(data.checks)?data.checks:[];
  const telemetry=data.telemetry||{};
  const hlsVariants=Number(data.hls_variant_count||0);
  const dashRepresentations=Number(data.dash_representation_count||0);
  const deliverySummary=data.delivery_summary||{};
  const dominantQuality=(telemetry.qualities&&Object.keys(telemetry.qualities).length)?Object.entries(telemetry.qualities).sort(function(a,b){return b[1]-a[1]})[0][0]:'-';
  const dominantAudio=(telemetry.audio_tracks&&Object.keys(telemetry.audio_tracks).length)?Object.entries(telemetry.audio_tracks).sort(function(a,b){return b[1]-a[1]})[0][0]:'-';
  out.innerHTML=
    '<div class="metric-list">'+
      '<div class="metric-row"><span>Yayin</span><strong>'+escHtml(data.stream_name||data.stream_key||'-')+'</strong></div>'+
      '<div class="metric-row"><span>ABR Profil</span><strong>'+escHtml(data.abr_profile_set||'balanced')+'</strong></div>'+
      '<div class="metric-row"><span>HLS varyant sayisi</span><strong>'+fmtInt(hlsVariants)+'</strong></div>'+
      '<div class="metric-row"><span>DASH representation sayisi</span><strong>'+fmtInt(dashRepresentations)+'</strong></div>'+
      '<div class="metric-row"><span>DASH ses representation</span><strong>'+fmtInt(data.dash_audio_representation_count||0)+'</strong></div>'+
      '<div class="metric-row"><span>Player telemetrisi</span><strong>'+fmtInt(telemetry.active_sessions||0)+' aktif / '+fmtInt(telemetry.total_stalls||0)+' stall</strong></div>'+
      '<div class="metric-row"><span>Kalite gecisi</span><strong>'+fmtInt(telemetry.total_quality_transitions||0)+'</strong></div>'+
      '<div class="metric-row"><span>Ses gecisi</span><strong>'+fmtInt(telemetry.total_audio_switches||0)+'</strong></div>'+
      '<div class="metric-row"><span>Baskin kalite</span><strong>'+escHtml(dominantQuality)+'</strong></div>'+
      '<div class="metric-row"><span>Baskin ses</span><strong>'+escHtml(dominantAudio)+'</strong></div>'+
      '<div class="metric-row"><span>Teslimat Ozeti</span><strong style="text-align:right">'+escHtml(deliverySummary.description||'-')+'</strong></div>'+
      '<div class="metric-row"><span>Policy JSON</span><span class="mono-wrap">'+escHtml(data.policy_json||'{}')+'</span></div>'+
    '</div>'+
    '<div style="margin-top:12px"><span class="tag tag-'+escHtml(deliverySummary.tone||'yellow')+'">'+escHtml(deliverySummary.label||'Durum bekleniyor')+'</span></div>'+
    '<div class="bar-list" style="margin-top:16px">'+checks.map(function(check){
      const tone='tag-'+(check.tone||'red');
      return '<div class="metric-row"><div><div>'+escHtml(check.description||check.code)+'</div>'+(check.detail?'<div class="form-hint" style="margin-top:4px">'+escHtml(check.detail)+'</div>':'')+'</div><span class="tag '+tone+'">'+escHtml(check.label||'Sorunlu')+'</span></div>';
    }).join('')+'</div>';
}

// ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â LOGS ÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚ÂÃƒÂ¢Ã¢â‚¬Â¢Ã‚Â
async function renderLogs(c){
  const logs=await api('/api/logs')||[];
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Loglar</h1>'+
      '<button class="btn btn-sm btn-danger" onclick="clearLogs()"><i class="bi bi-trash"></i> Temizle</button></div>'+
    '<div class="card">'+(logs.length===0
      ?'<div class="empty-state"><div class="icon"><i class="bi bi-file-earmark-text"></i></div><h3>Log kaydi yok</h3></div>'
      :'<div style="max-height:600px;overflow-y:auto"><table><thead><tr><th>Zaman</th><th>Seviye</th><th>Bilesen</th><th>Mesaj</th></tr></thead><tbody>'+
        logs.map(l=>{
          const colors={INFO:'var(--accent)',WARN:'var(--warning)',ERROR:'var(--danger)'};
          const time=fmtLocaleDateTime(l.created_at);
          return '<tr><td style="font-size:12px;color:var(--text-muted);white-space:nowrap">'+time+'</td>'+
            '<td><span style="color:'+(colors[l.level]||'var(--text-secondary)')+';font-weight:600;font-size:12px">'+l.level+'</span></td>'+
            '<td style="font-size:13px">'+escHtml(l.component)+'</td>'+
            '<td style="font-size:13px">'+escHtml(l.message)+'</td></tr>';
        }).join('')+
      '</tbody></table></div>')+
    '</div>';
}
async function clearLogs(){
  await api('/api/logs',{method:'DELETE'});toast('Loglar temizlendi');navigate('logs');
}

// ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ USERS ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢
async function renderUsers(c){
  const users=await api('/api/users')||[];
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Kullanici Yonetimi</h1>'+
      '<button class="btn btn-primary" onclick="showAddUserModal()">+ Yeni Kullanici</button></div>'+
    '<div class="card">'+(users.length===0
      ?'<div class="empty-state"><div class="icon">&#128100;</div><h3>Kullanici yok</h3></div>'
      :'<table><thead><tr><th>ID</th><th>Kullanici Adi</th><th>Rol</th><th>Olusturma</th><th></th></tr></thead><tbody>'+
        users.map(u=>'<tr>'+
          '<td>'+u.id+'</td>'+
          '<td><strong>'+escHtml(u.username)+'</strong></td>'+
          '<td><span class="tag '+(u.role==='admin'?'tag-blue':'tag-green')+'">'+escHtml(u.role)+'</span></td>'+
          '<td style="font-size:12px;color:var(--text-muted)">'+fmtLocaleDate(u.created_at)+'</td>'+
          '<td><button class="btn btn-sm btn-secondary" onclick="showEditUserModal('+u.id+',\''+escHtml(u.username)+'\',\''+escHtml(u.role)+'\')">Duzenle</button> '+
            '<button class="btn btn-sm btn-danger" onclick="deleteUser('+u.id+')">Sil</button></td>'+
        '</tr>').join('')+
      '</tbody></table>')+
    '</div><div id="user-modal"></div>';
}
function showAddUserModal(){
  document.getElementById('user-modal').innerHTML=
    '<div class="modal-overlay" onclick="if(event.target===this)this.remove()">'+
      '<div class="modal"><div class="modal-title">Yeni Kullanici</div>'+
        '<div class="form-group"><label class="form-label">Kullanici Adi</label><input class="form-input" id="mu-username"></div>'+
        '<div class="form-group"><label class="form-label">Sifre</label><input class="form-input" id="mu-password" type="password"></div>'+
        '<div class="form-group"><label class="form-label">Rol</label>'+
          '<select class="form-select" id="mu-role"><option value="admin">Admin</option><option value="editor">Editor</option><option value="viewer">Viewer</option></select></div>'+
        '<div style="display:flex;gap:12px;margin-top:20px">'+
          '<button class="btn btn-secondary" onclick="document.getElementById(\'user-modal\').innerHTML=\'\'">Iptal</button>'+
          '<button class="btn btn-primary" onclick="addUser()">Olustur</button></div>'+
      '</div></div>';
}
async function addUser(){
  const username=document.getElementById('mu-username').value;
  const password=document.getElementById('mu-password').value;
  const role=document.getElementById('mu-role').value;
  if(!username||!password){toast('Kullanici adi ve sifre gerekli','error');return}
  const res=await api('/api/users',{method:'POST',body:{username,password,role}});
  if(res.success){toast('Kullanici olusturuldu!');navigate('users')}
  else{toast(res.message||'Hata','error')}
}
function showEditUserModal(id,username,role){
  document.getElementById('user-modal').innerHTML=
    '<div class="modal-overlay" onclick="if(event.target===this)this.remove()">'+
      '<div class="modal"><div class="modal-title">Kullanici Duzenle</div>'+
        '<div class="form-group"><label class="form-label">Kullanici Adi</label><input class="form-input" id="eu-username" value="'+escHtml(username)+'"></div>'+
        '<div class="form-group"><label class="form-label">Yeni Sifre (bos birakin degistirmemek icin)</label><input class="form-input" id="eu-password" type="password"></div>'+
        '<div class="form-group"><label class="form-label">Rol</label>'+
          '<select class="form-select" id="eu-role"><option value="admin" '+(role==='admin'?'selected':'')+'>Admin</option><option value="editor" '+(role==='editor'?'selected':'')+'>Editor</option><option value="viewer" '+(role==='viewer'?'selected':'')+'>Viewer</option></select></div>'+
        '<div style="display:flex;gap:12px;margin-top:20px">'+
          '<button class="btn btn-secondary" onclick="document.getElementById(\'user-modal\').innerHTML=\'\'">Iptal</button>'+
          '<button class="btn btn-primary" onclick="editUser('+id+')">Kaydet</button></div>'+
      '</div></div>';
}
async function editUser(id){
  const username=document.getElementById('eu-username').value;
  const password=document.getElementById('eu-password').value;
  const role=document.getElementById('eu-role').value;
  const body={username,role};
  if(password)body.password=password;
  const res=await api('/api/users/'+id,{method:'PUT',body});
  if(res.success){toast('Kullanici guncellendi!');navigate('users')}
  else{toast(res.message||'Hata','error')}
}
async function deleteUser(id){
  if(!confirm('Bu kullaniciyi silmek istediginize emin misiniz?'))return;
  await api('/api/users/'+id,{method:'DELETE'});toast('Kullanici silindi');navigate('users');
}

// ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ PLAYER TEMPLATES ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢
function templateBackgroundStyle(t){
  if(t&&t.background_css)return t.background_css;
  if(t&&t.theme==='light')return 'background:linear-gradient(180deg,#f8fbff 0%,#dbeafe 100%);';
  if(t&&t.theme==='minimal')return 'background:#0f172a;';
  return 'background:linear-gradient(135deg,#030712 0%,#1d4ed8 100%);';
}
function templateControlStyle(t){
  return (t&&t.control_bar_css)||'background:rgba(15,23,42,.72);backdrop-filter:blur(10px);';
}
function templatePlayStyle(t){
  return (t&&t.play_button_css)||'color:#ffffff;';
}
function templateLogoPositionStyle(position){
  switch(position){
    case 'top-left': return 'top:14px;left:14px;';
    case 'bottom-left': return 'left:14px;bottom:14px;';
    case 'bottom-right': return 'right:14px;bottom:14px;';
    default: return 'top:14px;right:14px;';
  }
}
function renderPlayerTemplateThumbnail(t){
  const logoStyle=templateLogoPositionStyle(t&&t.logo_position);
  const logo=t&&t.logo_url?'<img src="'+escHtml(t.logo_url)+'" alt="" style="position:absolute;'+logoStyle+'height:26px;max-width:96px;object-fit:contain;opacity:'+(Number(t.logo_opacity||1))+';z-index:2">':'';
  const fallbackLogo=!logo?'<span class="template-thumb-logo" style="position:absolute;'+logoStyle+'">'+escHtml(((t&&t.watermark_text)||'FluxStream').slice(0,14))+'</span>':'';
  const title=(t&&t.show_title!==false)?'<div class="template-thumb-title">'+escHtml((t&&t.name)||'Player')+'</div>':'<div></div>';
  const badge=(t&&t.show_live_badge!==false)?'<span class="template-thumb-badge">Live</span>':'<span></span>';
  const watermark=(t&&t.watermark_text)?'<div class="template-thumb-watermark">'+escHtml(t.watermark_text)+'</div>':'';
  return '<div class="template-thumb" style="'+escHtml(templateBackgroundStyle(t))+'">'+
    '<div class="template-thumb-shell">'+
      logo+fallbackLogo+watermark+
      '<div class="template-thumb-header">'+title+badge+'</div>'+
      '<div class="template-thumb-center"><div class="template-thumb-play" style="'+escHtml(templatePlayStyle(t))+'"><i class="bi bi-play-fill"></i></div></div>'+
      '<div class="template-thumb-footer">'+
        '<div class="template-thumb-progress"><span></span></div>'+
        '<div class="template-thumb-controls" style="'+escHtml(templateControlStyle(t))+'">'+
          '<div class="left"><i class="bi bi-play-fill"></i><span>00:18</span></div>'+
          '<div class="right"><i class="bi bi-volume-up"></i><i class="bi bi-badge-hd"></i><i class="bi bi-fullscreen"></i></div>'+
        '</div>'+
      '</div>'+
    '</div>'+
  '</div>';
}
let playerTemplateStudioState={streamKey:'',streamName:'',streamPolicy:'',format:'player'};
function getTemplateStudioStreams(){
  return Array.isArray(window._playerTemplateStreams)?window._playerTemplateStreams:[];
}
function templateStudioCurrentStream(){
  const streams=getTemplateStudioStreams();
  return streams.find(function(s){return s.stream_key===playerTemplateStudioState.streamKey;})||streams[0]||null;
}
function ensureTemplateStudioState(){
  const streams=getTemplateStudioStreams();
  if(!streams.length){
    playerTemplateStudioState.streamKey='';
    playerTemplateStudioState.streamName='';
    playerTemplateStudioState.streamPolicy='';
    return null;
  }
  let current=streams.find(function(s){return s.stream_key===playerTemplateStudioState.streamKey;});
  if(!current){
    current=streams.find(function(s){return s.status==='live';})||streams[0];
    playerTemplateStudioState.streamKey=current.stream_key;
  }
  playerTemplateStudioState.streamName=current.name||current.stream_key;
  playerTemplateStudioState.streamPolicy=current.policy_json||'';
  if(!playerTemplateStudioState.format)playerTemplateStudioState.format='player';
  return current;
}
function templateStudioStreamOptions(){
  const streams=getTemplateStudioStreams();
  if(!streams.length)return '<option value="">-- Stream yok --</option>';
  return streams.map(function(s){
    return '<option value="'+escHtml(s.stream_key)+'" '+(s.stream_key===playerTemplateStudioState.streamKey?'selected':'')+'>'+escHtml(s.name)+' ('+escHtml(s.stream_key)+')</option>';
  }).join('');
}
function templateStudioFormatOptions(){
  const formats=[
    {value:'player',label:'Player'},
    {value:'iframe',label:'iframe'},
    {value:'hls',label:'HLS'},
    {value:'ll_hls',label:'LL-HLS'},
    {value:'dash',label:'DASH'},
    {value:'flv',label:'HTTP-FLV'},
    {value:'mp4',label:'MP4'},
    {value:'webm',label:'WebM'},
    {value:'mp3',label:'MP3'},
    {value:'aac',label:'AAC'},
    {value:'ogg',label:'OGG'},
    {value:'wav',label:'WAV'},
    {value:'flac',label:'FLAC'},
    {value:'icecast',label:'Icecast'}
  ];
  return formats.map(function(item){
    return '<option value="'+item.value+'" '+(item.value===playerTemplateStudioState.format?'selected':'')+'>'+item.label+'</option>';
  }).join('');
}
function buildTemplateQuery(t,streamName){
  const params=new URLSearchParams();
  params.set('player_title',streamName||t.name||'FluxStream');
  params.set('player_theme',t.theme||'dark');
  params.set('player_bg',t.background_css||'');
  params.set('player_controls',t.control_bar_css||'');
  params.set('player_play',t.play_button_css||'');
  params.set('player_logo',t.logo_url||'');
  params.set('player_logo_position',t.logo_position||'top-right');
  params.set('player_logo_opacity',String(Number(t.logo_opacity||1)));
  params.set('player_watermark',t.watermark_text||'');
  params.set('player_show_title',t.show_title===false?'0':'1');
  params.set('player_show_badge',t.show_live_badge===false?'0':'1');
  params.set('player_custom_css',t.custom_css||'');
  return params.toString();
}
function appendTemplateQuery(url,t,streamName){
  if(!url)return url;
  const query=buildTemplateQuery(t,streamName);
  if(!query)return url;
  return url+(url.indexOf('?')===-1?'?':'&')+query;
}
function templateAwareURLs(urls,t,streamName){
  const next=Object.assign({},urls||{});
  next.play=appendTemplateQuery(next.play,t,streamName);
  next.embed=appendTemplateQuery(next.embed,t,streamName);
  return next;
}
function playerURLForFormat(url,format){
  if(!url)return url;
  format=String(format||'').toLowerCase();
  if(!format || format==='player' || format==='iframe' || format==='jsapi')return url;
  return appendURLQuery(url,'format',format);
}
function updatePlayerTemplateStudioControls(){
  const stream=ensureTemplateStudioState();
  const streamSelect=document.getElementById('pt-stream-select');
  const formatSelect=document.getElementById('pt-format-select');
  if(streamSelect)streamSelect.innerHTML=templateStudioStreamOptions();
  if(formatSelect)formatSelect.innerHTML=templateStudioFormatOptions();
  const hint=document.getElementById('pt-stream-hint');
  if(hint){
    hint.innerHTML=stream?('Secili kaynak: <strong>'+escHtml(stream.name)+'</strong> • '+escHtml(stream.stream_key)):'Kaynak stream secilmedi';
  }
  const modalStream=document.getElementById('pt-modal-stream');
  const modalFormat=document.getElementById('pt-modal-format');
  if(modalStream)modalStream.innerHTML=templateStudioStreamOptions();
  if(modalFormat)modalFormat.innerHTML=templateStudioFormatOptions();
}
function updatePlayerTemplateModalPreview(){
  const holder=document.getElementById('pt-current-template-id');
  if(!holder)return;
  const id=parseInt(holder.value||'0',10)||0;
  updatePlayerTemplatePreview(id);
}
function buildTemplatePreviewSrc(previewURLs,format){
  var src=(previewURLs&&previewURLs.embed)||'';
  if(!src)return '';
  src=playerURLForFormat(src,format);
  src=appendURLQuery(src,'debug','1');
  src=appendURLQuery(src,'autoplay','1');
  src=appendURLQuery(src,'muted','1');
  return src;
}
async function updatePlayerTemplatePreview(id){
  const prev=document.getElementById('pt-live-preview');
  const code=document.getElementById('pt-live-embed-code');
  if(!prev||!code)return;
  const stream=ensureTemplateStudioState();
  if(!stream){
    prev.innerHTML='<div class="empty-state"><div class="icon"><i class="bi bi-broadcast"></i></div><h3>Kaynak stream yok</h3><p style="color:var(--text-muted)">Template preview icin en az bir stream olusturun.</p></div>';
    code.innerHTML='';
    return;
  }
  if(!id){
    prev.innerHTML='<div class="empty-state"><div class="icon"><i class="bi bi-palette"></i></div><h3>Kaydedin ve deneyin</h3><p style="color:var(--text-muted)">Yeni bir template icin once kaydet, sonra secili stream ile player preview ve embed kodunu gor.</p></div>';
    code.innerHTML='';
    return;
  }
  const template=await api('/api/players/'+id);
  const settings=await api('/api/settings');
  const access=await getPlaybackAccess(stream.stream_key,settings,stream.policy_json||'');
  const previewRawURLs=getPreviewURLs(stream.stream_key,settings,stream.name,access);
  const publicRawURLs=getAllURLs(stream.stream_key,settings,stream.name,access);
  const previewURLs=templateAwareURLs(previewRawURLs,template,stream.name);
  const urls=templateAwareURLs(publicRawURLs,template,stream.name);
  const isAudioPreview=playerTemplateStudioState.format==='mp3'||playerTemplateStudioState.format==='aac'||playerTemplateStudioState.format==='ogg'||playerTemplateStudioState.format==='wav'||playerTemplateStudioState.format==='flac'||playerTemplateStudioState.format==='icecast';
  const previewSrc=buildTemplatePreviewSrc(previewURLs,playerTemplateStudioState.format);
  const previewBundle=buildEmbedBundle(playerTemplateStudioState.format,stream.stream_key,previewURLs,960,playerTemplateStudioState.format==='mp3'||playerTemplateStudioState.format==='aac'||playerTemplateStudioState.format==='ogg'||playerTemplateStudioState.format==='wav'||playerTemplateStudioState.format==='flac'||playerTemplateStudioState.format==='icecast'?120:540,true,true);
  const bundle=buildEmbedBundle(playerTemplateStudioState.format,stream.stream_key,urls,960,playerTemplateStudioState.format==='mp3'||playerTemplateStudioState.format==='aac'||playerTemplateStudioState.format==='ogg'||playerTemplateStudioState.format==='wav'||playerTemplateStudioState.format==='flac'||playerTemplateStudioState.format==='icecast'?120:540,true,true);
  prev.innerHTML='<div style="position:relative;'+(isAudioPreview?'height:140px;':'padding-top:56.25%;')+'background:#05070b;border-radius:12px;overflow:hidden">'+
    (previewSrc?'<iframe src="'+previewSrc+'" style="position:absolute;inset:0;width:100%;height:100%;border:none;background:#000" allow="autoplay;fullscreen" allowfullscreen></iframe>':'')+
    '</div>';
  code.innerHTML=
    '<div class="metric-row"><div><div class="setting-label">'+escHtml(stream.name)+'</div><div class="setting-desc">'+escHtml(playerTemplateStudioState.format.toUpperCase())+' • Template preview</div></div><span class="tag tag-blue">'+escHtml((template.name||'Template'))+'</span></div>'+
    copyField(bundle.primaryLabel||'Embed',bundle.primary||'')+
    (bundle.direct?copyField(bundle.directLabel||'URL',bundle.direct):'')+
    copyField('Player URL',playerURLForFormat(urls.play||'',playerTemplateStudioState.format))+
    copyField('Embed URL',urls.embed||'')+
    (bundle.note?'<div class="form-hint" style="margin-top:10px">'+escHtml(bundle.note)+'</div>':'');
}
async function renderPlayerTemplates(c){
  const [templates,streams]=await Promise.all([api('/api/players'),api('/api/streams')]);
  window._playerTemplateStreams=Array.isArray(streams)?streams:[];
  ensureTemplateStudioState();
  c.innerHTML=
    '<div class="page-header"><div><h1 class="page-title">Player Sablonlari</h1><div style="color:var(--text-muted);font-size:13px;margin-top:6px">Kurulu gelen hazir sablonlari temel alip duzenleyebilir veya sifirdan yeni sablon olusturabilirsiniz.</div></div>'+
      '<button class="btn btn-primary" onclick="showPlayerModal()">+ Yeni Sablon</button></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-grid card-grid-2">'+
      '<div class="form-group"><label class="form-label">Onizleme Kaynagi</label><select class="form-select" id="pt-stream-select" onchange="playerTemplateStudioState.streamKey=this.value;updatePlayerTemplateStudioControls();updatePlayerTemplateModalPreview()">'+templateStudioStreamOptions()+'</select><div class="form-hint" id="pt-stream-hint"></div></div>'+
      '<div class="form-group"><label class="form-label">Onizleme Formati</label><select class="form-select" id="pt-format-select" onchange="playerTemplateStudioState.format=this.value;updatePlayerTemplateModalPreview()">'+templateStudioFormatOptions()+'</select><div class="form-hint">Kaydedilen template icin bu formatta embed kodu ve preview olusur.</div></div>'+
    '</div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="metric-list">'+
      '<div class="metric-row"><span>Hazir baslangic sablonu</span><strong>12+</strong></div>'+
      '<div class="metric-row"><span>Kullanim</span><strong>Duzenle -> Kaydet -> Embed tarafinda kullan</strong></div>'+
      '<div class="metric-row"><span>Amac</span><strong>Canli TV, radyo, minimal player, cam tasarim ve parlak vitrini hizla baslatmak</strong></div>'+
    '</div></div>'+
    (templates.length===0
      ?'<div class="card"><div class="empty-state"><div class="icon"><i class="bi bi-pc-display-horizontal"></i></div><h3>Henuz sablon yok</h3><p style="color:var(--text-muted)">Ozel player sablonu olusturun</p></div></div>'
      :'<div class="card-grid card-grid-3">'+templates.map(t=>
        '<div class="card" style="cursor:pointer" onclick="showPlayerModal('+t.id+')">'+
          '<div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">'+
            '<div class="card-title">'+escHtml(t.name)+'</div>'+
            '<span class="tag tag-blue">'+escHtml(t.theme||'dark')+'</span>'+
          '</div>'+
          renderPlayerTemplateThumbnail(t)+
          '<div class="form-hint" style="margin:-2px 0 12px">Kaynak: '+escHtml(playerTemplateStudioState.streamName||'Stream secin')+'</div>'+
          '<div style="display:flex;gap:8px">'+
            '<button class="btn btn-sm btn-secondary" onclick="event.stopPropagation();showPlayerModal('+t.id+')">Duzenle</button>'+
            '<button class="btn btn-sm btn-primary" onclick="event.stopPropagation();showPlayerModal('+t.id+')">Onizle ve Kodlar</button>'+
            '<button class="btn btn-sm btn-danger" onclick="event.stopPropagation();deletePlayerTemplate('+t.id+')">Sil</button>'+
          '</div>'+
        '</div>'
      ).join('')+'</div>')+
    '<div id="player-modal"></div>';
  updatePlayerTemplateStudioControls();
}
async function showPlayerModal(id){
  let pt={name:'',theme:'dark',background_css:'',control_bar_css:'',play_button_css:'',logo_url:'',logo_position:'top-right',logo_opacity:1,watermark_text:'',show_title:true,show_live_badge:true,custom_css:''};
  if(id){const data=await api('/api/players/'+id);if(data&&!data.error)pt=data}
  const isEdit=!!id;
  document.getElementById('player-modal').innerHTML=
    '<div class="modal-overlay" onclick="if(event.target===this)this.remove()">'+
      '<div class="modal" style="max-width:1080px">'+
        '<div class="modal-title">'+(isEdit?'Sablon Duzenle':'Yeni Player Sablonu')+'</div>'+
        '<div class="card-grid card-grid-2" style="align-items:start">'+
          '<div>'+
            '<div style="margin-bottom:18px">'+renderPlayerTemplateThumbnail(pt)+'</div>'+
            '<div class="card-grid card-grid-2">'+
          '<div class="form-group"><label class="form-label">Sablon Adi *</label><input class="form-input" id="pt-name" value="'+escHtml(pt.name)+'"></div>'+
          '<div class="form-group"><label class="form-label">Tema</label>'+
            '<select class="form-select" id="pt-theme"><option value="dark" '+(pt.theme==='dark'?'selected':'')+'>Dark</option><option value="light" '+(pt.theme==='light'?'selected':'')+'>Light</option><option value="minimal" '+(pt.theme==='minimal'?'selected':'')+'>Minimal</option><option value="custom" '+(pt.theme==='custom'?'selected':'')+'>Custom</option></select></div>'+
            '</div>'+
            '<div class="card-grid card-grid-2">'+
          '<div class="form-group"><label class="form-label">Logo URL</label><input class="form-input" id="pt-logo-url" value="'+escHtml(pt.logo_url)+'" placeholder="https://..."></div>'+
          '<div class="form-group"><label class="form-label">Logo Konum</label>'+
            '<select class="form-select" id="pt-logo-pos"><option value="top-right" '+(pt.logo_position==='top-right'?'selected':'')+'>Sag Ust</option><option value="top-left" '+(pt.logo_position==='top-left'?'selected':'')+'>Sol Ust</option><option value="bottom-right" '+(pt.logo_position==='bottom-right'?'selected':'')+'>Sag Alt</option><option value="bottom-left" '+(pt.logo_position==='bottom-left'?'selected':'')+'>Sol Alt</option></select></div>'+
            '</div>'+
            '<div class="card-grid card-grid-2">'+
          '<div class="form-group"><label class="form-label">Logo Seffaflik</label><input class="form-input" id="pt-logo-opacity" type="number" min="0" max="1" step="0.1" value="'+(pt.logo_opacity||1)+'"></div>'+
          '<div class="form-group"><label class="form-label">Watermark Yazi</label><input class="form-input" id="pt-watermark" value="'+escHtml(pt.watermark_text)+'"></div>'+
            '</div>'+
            '<div class="card-grid card-grid-2">'+
          '<div class="setting-row" style="border:none"><div><div class="setting-label">Baslik Goster</div></div>'+
            '<label class="toggle"><input type="checkbox" id="pt-show-title" '+(pt.show_title?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
          '<div class="setting-row" style="border:none"><div><div class="setting-label">CANLI Badge</div></div>'+
            '<label class="toggle"><input type="checkbox" id="pt-show-badge" '+(pt.show_live_badge?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
            '</div>'+
            '<div class="form-group"><label class="form-label">Arkaplan CSS</label><input class="form-input" id="pt-bg-css" value="'+escHtml(pt.background_css)+'" placeholder="background: #000;"></div>'+
            '<div class="form-group"><label class="form-label">Kontrol Cubugu CSS</label><input class="form-input" id="pt-ctrl-css" value="'+escHtml(pt.control_bar_css)+'"></div>'+
            '<div class="form-group"><label class="form-label">Play Butonu CSS</label><input class="form-input" id="pt-play-css" value="'+escHtml(pt.play_button_css)+'"></div>'+
            '<div class="form-group"><label class="form-label">Ozel CSS</label><textarea class="form-textarea" id="pt-custom-css" rows="4">'+escHtml(pt.custom_css)+'</textarea></div>'+
            '<div style="display:flex;gap:12px;margin-top:20px">'+
              '<button class="btn btn-secondary" onclick="document.getElementById(\'player-modal\').innerHTML=\'\'">Iptal</button>'+
              '<button class="btn btn-primary" onclick="savePlayerTemplate('+(id||'null')+')">'+(isEdit?'Guncelle':'Olustur')+'</button></div>'+
          '</div>'+
          '<div>'+
            '<div class="card" style="padding:18px;margin-bottom:16px"><div class="card-grid card-grid-2"><div class="form-group"><label class="form-label">Kaynak stream</label><select class="form-select" id="pt-modal-stream" onchange="playerTemplateStudioState.streamKey=this.value;updatePlayerTemplateStudioControls();updatePlayerTemplateModalPreview()">'+templateStudioStreamOptions()+'</select></div><div class="form-group"><label class="form-label">Format</label><select class="form-select" id="pt-modal-format" onchange="playerTemplateStudioState.format=this.value;updatePlayerTemplateModalPreview()">'+templateStudioFormatOptions()+'</select></div></div><input type="hidden" id="pt-current-template-id" value="'+(id||0)+'"></div>'+
            '<div class="card" style="padding:18px;margin-bottom:16px"><div class="card-header"><div class="card-title">Canli Player Onizleme</div><span class="form-hint">Secili stream ve format ile</span></div><div class="card-body"><div id="pt-live-preview"></div></div></div>'+
            '<div class="card" style="padding:18px"><div class="card-header"><div class="card-title">Embed Kodlari</div><span class="form-hint">Template + stream birlesik cikti</span></div><div class="card-body" id="pt-live-embed-code"></div></div>'+
          '</div>'+
        '</div>'+
      '</div></div>';
  applyTranslations(document.getElementById('player-modal'));
  updatePlayerTemplatePreview(id);
}
async function savePlayerTemplate(id){
  const body={
    name:document.getElementById('pt-name').value,
    theme:document.getElementById('pt-theme').value,
    logo_url:document.getElementById('pt-logo-url').value,
    logo_position:document.getElementById('pt-logo-pos').value,
    logo_opacity:parseFloat(document.getElementById('pt-logo-opacity').value)||1,
    watermark_text:document.getElementById('pt-watermark').value,
    show_title:document.getElementById('pt-show-title').checked,
    show_live_badge:document.getElementById('pt-show-badge').checked,
    background_css:document.getElementById('pt-bg-css').value,
    control_bar_css:document.getElementById('pt-ctrl-css').value,
    play_button_css:document.getElementById('pt-play-css').value,
    custom_css:document.getElementById('pt-custom-css').value,
  };
  if(!body.name){toast('Sablon adi gerekli','error');return}
  if(id){
    const res=await api('/api/players/'+id,{method:'PUT',body});
    if(res.success)toast('Sablon guncellendi!');
    else{toast(res.message||'Hata','error');return}
  }else{
    const res=await api('/api/players',{method:'POST',body});
    if(res.id)toast('Sablon olusturuldu!');
    else{toast(res.message||'Hata','error');return}
  }
  navigate('player-templates');
}
async function deletePlayerTemplate(id){
  if(!confirm('Bu sablonu silmek istediginize emin misiniz?'))return;
  await api('/api/players/'+id,{method:'DELETE'});toast('Sablon silindi');navigate('player-templates');
}

// ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ ADVANCED EMBED GENERATOR ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢ÃƒÂ¢Ã¢â‚¬Â¢
let embedTab='iframe';
async function renderAdvancedEmbed(c){
  embedTab='iframe';
  const streams=await api('/api/streams')||[];
  const settings=await api('/api/settings');
  var settingsJSON=JSON.stringify(settings).replace(/</g,'\\u003c');

  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Gelismis Embed Olusturucu</h1>'+
      '<p style="color:var(--text-muted);font-size:13px">19 format icin embed kodu uretici</p></div>'+
    '<div class="card" style="margin-bottom:16px">'+
      '<div class="card-grid card-grid-2">'+
        '<div class="form-group"><label class="form-label">Yayin Sec</label>'+
          '<select class="form-select" id="ae-stream" onchange="updateEmbedPreview()" data-settings=\''+settingsJSON+'\'>'+
            '<option value="">-- Yayin Secin --</option>'+
            streams.map(function(s){
              return '<option value="'+s.stream_key+'" data-stream-name="'+escHtml(s.name)+'" data-policy-json="'+escHtml(String(s.policy_json||''))+'">'+escHtml(s.name)+' ('+s.stream_key+')</option>';
            }).join('')+
          '</select></div>'+
        '<div class="form-group"><label class="form-label">Tema</label>'+
          '<select class="form-select" id="ae-theme" onchange="updateEmbedPreview()">'+
            '<option value="dark">Dark</option><option value="light">Light</option><option value="minimal">Minimal</option></select></div>'+
      '</div>'+
      '<div class="card-grid card-grid-4">'+
        '<div class="form-group"><label class="form-label">Genislik</label><input class="form-input" id="ae-width" value="1280" type="number" onchange="updateEmbedPreview()"></div>'+
        '<div class="form-group"><label class="form-label">Yukseklik</label><input class="form-input" id="ae-height" value="720" type="number" onchange="updateEmbedPreview()"></div>'+
        '<div class="setting-row" style="border:none;padding:0;margin-top:16px"><div><div class="setting-label">Otomatik Oynat</div></div>'+
          '<label class="toggle"><input type="checkbox" id="ae-autoplay" checked onchange="updateEmbedPreview()"><span class="toggle-slider"></span></label></div>'+
        '<div class="setting-row" style="border:none;padding:0;margin-top:16px"><div><div class="setting-label">Sessiz</div></div>'+
          '<label class="toggle"><input type="checkbox" id="ae-muted" checked onchange="updateEmbedPreview()"><span class="toggle-slider"></span></label></div>'+
      '</div>'+
    '</div>'+
    '<div class="card">'+
      '<div class="tabs" id="embed-tabs" style="flex-wrap:wrap">'+
        '<div class="tab active" onclick="switchEmbedTab(\'iframe\')">iframe</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'hls\')">HLS</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'ll_hls\')">LL-HLS</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'dash\')">DASH</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'flv\')">HTTP-FLV</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'whep\')">WHEP</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'mp4\')">MP4</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'webm\')">WebM</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'rtmp\')">RTMP</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'rtsp\')">RTSP</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'srt\')">SRT</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'mp3\')">MP3</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'aac\')">AAC</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'ogg\')">OGG</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'wav\')">WAV</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'flac\')">FLAC</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'icecast\')">Icecast</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'player\')">Player</div>'+
        '<div class="tab" onclick="switchEmbedTab(\'jsapi\')">JS API</div>'+
      '</div>'+
      '<div id="embed-output" style="margin-top:16px"></div>'+
      '<div style="margin-top:16px;padding:16px;background:#000;border-radius:var(--radius-sm)">'+
        '<div style="font-size:12px;color:var(--text-muted);margin-bottom:8px">Onizleme</div>'+
        '<div id="embed-preview" style="background:#111;border-radius:8px;overflow:hidden;position:relative;padding-top:56.25%">'+
          '<div style="position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);color:var(--text-muted)">Yayin secin</div>'+
        '</div>'+
      '</div>'+
    '</div>';
  updateEmbedPreview();
}
function switchEmbedTab(tab){
  embedTab=tab;
  document.querySelectorAll('#embed-tabs .tab').forEach(t=>t.classList.remove('active'));
  document.querySelectorAll('#embed-tabs .tab').forEach(t=>{
    const oc=t.getAttribute('onclick')||'';
    if(oc.indexOf('\''+tab+'\'')!==-1)t.classList.add('active');
  });
  updateEmbedPreview();
}
let previewHls=null;
let previewDash=null;
let previewFlv=null;
const embedScriptCache={};
function destroyEmbedPreviewPlayers(){
  try{if(previewHls){previewHls.destroy();previewHls=null}}catch(e){}
  try{if(previewDash){previewDash.reset();previewDash=null}}catch(e){}
  try{if(previewFlv){previewFlv.destroy();previewFlv=null}}catch(e){}
}
function loadEmbedScript(url){
  if(embedScriptCache[url])return embedScriptCache[url];
  embedScriptCache[url]=new Promise(function(resolve,reject){
    var ex=document.querySelector('script[data-embed-src="'+url+'"]');
    if(ex){
      if(ex.dataset.loaded==='1'){resolve();return}
      ex.addEventListener('load',function(){ex.dataset.loaded='1';resolve()},{once:true});
      ex.addEventListener('error',function(){reject(new Error('Script yuklenemedi'))},{once:true});
      return;
    }
    var s=document.createElement('script');
    s.src=url;
    s.async=true;
    s.dataset.embedSrc=url;
    s.onload=function(){s.dataset.loaded='1';resolve()};
    s.onerror=function(){reject(new Error('Script yuklenemedi'))};
    document.head.appendChild(s);
  });
  return embedScriptCache[url];
}
function setPreviewFrame(prev,src){
  prev.innerHTML='<iframe src="'+src+'" style="position:absolute;top:0;left:0;width:100%;height:100%;border:none;background:#000" allow="autoplay;fullscreen" allowfullscreen></iframe>';
}
function setPreviewFallback(prev,embedURL,msg){
  setPreviewFrame(prev,embedURL);
  prev.innerHTML+=
    '<div style="position:absolute;left:12px;bottom:12px;right:12px;padding:8px 10px;border-radius:8px;background:rgba(0,0,0,.55);color:#e8eefc;font-size:12px;line-height:1.4">'+
    msg+'</div>';
}
function setPreviewVideo(prev,src,autoplay,muted){
  var a=autoplay?' autoplay':'';
  var m=muted?' muted':'';
  prev.innerHTML='<video id="embed-preview-media" controls playsinline'+a+m+' style="position:absolute;top:0;left:0;width:100%;height:100%;background:#000;object-fit:contain"></video>';
  var v=document.getElementById('embed-preview-media');
  if(v&&src)v.src=src;
  return v;
}
function setPreviewAudio(prev,src,autoplay,muted){
  var a=autoplay?' autoplay':'';
  var m=muted?' muted':'';
  prev.innerHTML=
    '<div style="position:absolute;inset:0;display:flex;align-items:center;justify-content:center;background:#070b14">'+
      '<audio id="embed-preview-audio" controls'+a+m+' style="width:min(92%,460px)"></audio>'+
    '</div>';
  var au=document.getElementById('embed-preview-audio');
  if(au&&src)au.src=src;
}
function mountDirectStreamPreview(video,url,autoplay){
  if(!video||!url)return;
  var retryTimer=null;
  function load(){
    video.src=url+(url.indexOf('?')===-1?'?':'&')+'ts='+Date.now();
    video.load();
    if(autoplay)video.play().catch(function(){});
  }
  video.addEventListener('error',function(){
    if(retryTimer)clearTimeout(retryTimer);
    retryTimer=setTimeout(load,3000);
  });
  load();
}
function mountHLSPreview(video,url,autoplay,isLowLatency){
  if(!video)return;
  var candidates=Array.isArray(url)?url.filter(Boolean):[url];
  if(!candidates.length)return;
  function startWith(resolvedURL){
    if(video.canPlayType&&video.canPlayType('application/vnd.apple.mpegurl')){
      video.src=resolvedURL;
      if(autoplay)video.play().catch(function(){});
      return;
    }
    loadEmbedScript('/static/vendor/hls.min.js').then(function(){
      if(!window.Hls)return;
      if(window.Hls.isSupported()){
        previewHls=new Hls({
          liveSyncDurationCount:isLowLatency?3:4,
          liveMaxLatencyDurationCount:isLowLatency?6:10,
          lowLatencyMode:!!isLowLatency,
          maxBufferLength:isLowLatency?12:20,
          maxMaxBufferLength:isLowLatency?20:40,
          backBufferLength:30,
          enableWorker:true
        });
        previewHls.loadSource(resolvedURL);
        previewHls.attachMedia(video);
        previewHls.on(Hls.Events.MANIFEST_PARSED,function(){
          if(autoplay)video.play().catch(function(){});
        });
      }else{
        video.src=resolvedURL;
        if(autoplay)video.play().catch(function(){});
      }
    }).catch(function(){});
  }
  (async function(){
    for(const candidate of candidates){
      try{
        const res=await fetch(candidate,{method:'HEAD',cache:'no-store'});
        if(res.ok){startWith(candidate);return;}
      }catch(e){}
    }
    startWith(candidates[candidates.length-1]);
  })();
}
function mountDashPreview(video,url,autoplay){
  if(!video)return;
  loadEmbedScript('/static/vendor/dash.all.min.js').then(function(){
    if(!window.dashjs)return;
    previewDash=window.dashjs.MediaPlayer().create();
    previewDash.initialize(video,url,autoplay);
    previewDash.updateSettings({streaming:{lowLatencyEnabled:true}});
  }).catch(function(){});
}
function mountFLVPreview(video,url,autoplay){
  if(!video)return;
  loadEmbedScript('/static/vendor/mpegts.min.js').then(function(){
    if(!window.mpegts||!window.mpegts.getFeatureList||!window.mpegts.getFeatureList().mseLivePlayback)return;
    previewFlv=window.mpegts.createPlayer({type:'flv',isLive:true,url:url});
    previewFlv.attachMediaElement(video);
    previewFlv.load();
    if(autoplay)previewFlv.play().catch(function(){});
  }).catch(function(){});
}
function buildHTML5MediaSnippet(tag,url,w,h,autoplay,muted){
  var attrs=' controls playsinline'+(autoplay?' autoplay':'')+(muted?' muted':'');
  var style=tag==='audio'?'width:100%;max-width:'+w+'px':'width:'+w+'px;height:'+h+'px;max-width:100%;background:#000';
  return '<'+tag+' src="'+url+'"'+attrs+' style="'+style+'"></'+tag+'>';
}
function buildScriptedVideoSnippet(playerVar,videoID,scriptURL,scriptBody,w,h,autoplay,muted){
  return '<video id="'+videoID+'" controls playsinline'+(autoplay?' autoplay':'')+(muted?' muted':'')+' style="width:'+w+'px;height:'+h+'px;max-width:100%;background:#000"></video>\\n'+
    '<script src="'+scriptURL+'"><\\/script>\\n<script>\\n'+scriptBody+'\\n<\\/script>';
}
function buildDirectLiveVideoSnippet(videoID,url,w,h,autoplay,muted){
  return '<video id="'+videoID+'" controls playsinline'+(autoplay?' autoplay':'')+(muted?' muted':'')+' style="width:'+w+'px;height:'+h+'px;max-width:100%;background:#000"></video>\\n'+
    '<script>\\n(function(){var v=document.getElementById("'+videoID+'");if(!v)return;var timer=null;function load(){v.src="'+url+'"+(("'+url+'".indexOf("?")===-1)?"?":"&")+"ts="+Date.now();v.load();'+(autoplay?'v.play().catch(function(){});':'')+'}v.addEventListener("error",function(){if(timer)clearTimeout(timer);timer=setTimeout(load,3000);});load();})();\\n<\\/script>';
}
function buildFormatEmbedFrame(key,format,w,h,autoplay,muted,embedBase){
  var src=(embedBase||('/embed/'+key))+(String(embedBase||'').indexOf('?')===-1?'?':'&')+'format='+format+'&autoplay='+(autoplay?'1':'0')+'&muted='+(muted?'1':'0');
  return '<iframe src="'+src+'" width="'+w+'" height="'+h+'" frameborder="0" allow="autoplay;fullscreen" allowfullscreen></iframe>';
}
function directURLForEmbedTab(tab,urls){
  switch(tab){
    case 'iframe': return urls.embed;
    case 'hls': return urls.hls;
    case 'll_hls': return urls.ll_hls;
    case 'dash': return urls.dash;
    case 'flv': return urls.http_flv;
    case 'whep': return urls.whep;
    case 'mp4': return urls.fmp4;
    case 'webm': return urls.webm;
    case 'rtmp': return urls.rtmp;
    case 'rtsp': return urls.rtsp;
    case 'srt': return urls.srt;
    case 'mp3': return urls.mp3;
    case 'aac': return urls.aac;
    case 'ogg': return urls.ogg;
    case 'wav': return urls.wav;
    case 'flac': return urls.flac;
    case 'icecast': return urls.icecast;
    case 'player': return urls.play;
    case 'jsapi': return urls.hls;
    default: return urls.embed;
  }
}
function buildEmbedBundle(tab,key,urls,w,h,autoplay,muted){
  var embedBase=urls.embed||('/embed/'+key);
  var fallbackFrame='<iframe src="'+embedBase+(embedBase.indexOf('?')===-1?'?':'&')+'autoplay='+(autoplay?'1':'0')+'&muted='+(muted?'1':'0')+'" width="'+w+'" height="'+h+'" frameborder="0" allow="autoplay;fullscreen" allowfullscreen></iframe>';
  var audioFrameHeight=Math.min(parseInt(h,10)||720,160);
  var audioFrame=function(fmt){
    var base=urls.embed||('/embed/'+key);
    return '<iframe src="'+base+(base.indexOf('?')===-1?'?':'&')+'format='+fmt+'&autoplay='+(autoplay?'1':'0')+'&muted='+(muted?'1':'0')+'" width="'+w+'" height="'+audioFrameHeight+'" frameborder="0" allow="autoplay" allowfullscreen></iframe>';
  };
  switch(tab){
    case 'iframe':
      return {primaryLabel:'Tarayici Embed Kodu',primary:fallbackFrame,directLabel:'Embed URL',direct:urls.embed};
    case 'hls':
      return {
        primaryLabel:'Tarayici Embed Kodu',
        primary:buildScriptedVideoSnippet('hls','fs-hls-player',urls.asset_hls,'var video=document.getElementById(\"fs-hls-player\");if(window.Hls&&Hls.isSupported()){var hls=new Hls();hls.loadSource(\"'+urls.hls+'\");hls.attachMedia(video);}else{video.src=\"'+urls.hls+'\";}',w,h,autoplay,muted),
        directLabel:'HLS URL',
        direct:urls.hls,
      };
    case 'll_hls':
      return {
        primaryLabel:'Tarayici Embed Kodu',
        primary:buildFormatEmbedFrame(key,'ll_hls',w,h,autoplay,muted,urls.embed),
        directLabel:'LL-HLS URL',
        direct:urls.ll_hls,
        note:'Tarayici uyumlu iframe player stabil DASH/HLS onizleme kullanir. Ham LL-HLS cikisi alttaki URL alanindadir.'
      };
    case 'dash':
      return {
        primaryLabel:'Tarayici Embed Kodu',
        primary:buildFormatEmbedFrame(key,'dash',w,h,autoplay,muted,urls.embed),
        directLabel:'DASH URL',
        direct:urls.dash,
        note:'Tarayici uyumlu iframe player stabil DASH/HLS onizleme kullanir. Ham DASH cikisi alttaki URL alanindadir.'
      };
    case 'flv':
      return {
        primaryLabel:'Tarayici Embed Kodu',
        primary:buildFormatEmbedFrame(key,'flv',w,h,autoplay,muted,urls.embed),
        directLabel:'HTTP-FLV URL',
        direct:urls.http_flv,
        note:'Tarayici uyumlu iframe player stabil DASH/HLS onizleme kullanir. Ham HTTP-FLV cikisi alttaki URL alanindadir.'
      };
    case 'mp4':
      return {primaryLabel:'Tarayici Embed Kodu',primary:buildFormatEmbedFrame(key,'mp4',w,h,autoplay,muted,urls.embed),directLabel:'MP4 URL',direct:urls.fmp4,note:'Tarayici uyumlu iframe player stabil DASH/HLS onizleme kullanir. Ham MP4 cikisi alttaki URL alanindadir.'};
    case 'webm':
      return {primaryLabel:'Tarayici Embed Kodu',primary:buildFormatEmbedFrame(key,'webm',w,h,autoplay,muted,urls.embed),directLabel:'WebM URL',direct:urls.webm,note:'Tarayici uyumlu iframe player stabil DASH/HLS onizleme kullanir. Ham WebM cikisi alttaki URL alanindadir.'};
    case 'player':
      return {
        primaryLabel:'Tarayici Embed Kodu',
        primary:fallbackFrame,
        directLabel:'Player URL',
        direct:urls.play,
        note:'Gomulu kullanimda tam sayfa player yerine iframe dostu embed gorunumu kullanilir. Direkt link yine Player URL alanindadir.'
      };
    case 'jsapi':
      return {
        primaryLabel:'JS API Kodu',
        primary:'<div id="player"></div>\\n<script src="'+urls.asset_hls+'"><\\/script>\\n<script>\\nvar video=document.createElement("video");\\nvideo.controls=true;\\nvideo.autoplay='+(autoplay?'true':'false')+';\\nvideo.muted='+(muted?'true':'false')+';\\nvideo.style.width="100%";\\ndocument.getElementById("player").appendChild(video);\\nif(window.Hls&&Hls.isSupported()){var hls=new Hls();hls.loadSource("'+urls.hls+'");hls.attachMedia(video);}else{video.src="'+urls.hls+'";}\\n<\\/script>',
        directLabel:'HLS URL',
        direct:urls.hls,
      };
    case 'mp3':
    case 'aac':
    case 'ogg':
    case 'wav':
    case 'flac':
    case 'icecast':
      return {
        primaryLabel:'Tarayici Embed Kodu',
        primary:audioFrame(tab),
        directLabel:'Dogrudan Cikis URL',
        direct:directURLForEmbedTab(tab,urls),
        note:'Bu sekme icin iframe player yalnizca ses oynatir. Alttaki alan dogrudan cikis URL\'sidir.'
      };
    case 'whep':
    case 'rtmp':
    case 'rtsp':
    case 'srt':
      return {
        primaryLabel:'Tarayici Embed Kodu',
        primary:fallbackFrame,
        directLabel:'Protokol URL',
        direct:directURLForEmbedTab(tab,urls),
        note:'Bu protokol tarayicida dogrudan oynatilamaz. Tarayici embed kodu DASH/HLS tabanli player fallback kullanir.'
      };
    default:
      return {primaryLabel:'Tarayici Embed Kodu',primary:fallbackFrame,directLabel:'Embed URL',direct:urls.embed};
  }
}
function renderAdvancedPreview(prev,key,autoplay,muted,urls){
  if(!prev||!key)return;
  destroyEmbedPreviewPlayers();
  const previewURLs=urls||{};
  const hlsCandidates=[previewURLs.hls,previewURLs.hls_media].filter(Boolean);
  const defaultFrame=(previewURLs.embed||('/embed/'+key))+(String(previewURLs.embed||'').indexOf('?')===-1?'?':'&')+'autoplay='+(autoplay?'1':'0')+'&muted='+(muted?'1':'0')+'&debug=1';
  const formatFrame=function(fmt){
    const base=previewURLs.embed||('/embed/'+key);
    return base+(String(base).indexOf('?')===-1?'?':'&')+'format='+fmt+'&autoplay='+(autoplay?'1':'0')+'&muted='+(muted?'1':'0')+'&debug=1';
  };

  switch(embedTab){
    case 'iframe':
      setPreviewFrame(prev,defaultFrame);
      break;
    case 'player':
      setPreviewFrame(prev,defaultFrame);
      break;
    case 'hls':
    case 'jsapi':{
      setPreviewFrame(prev,formatFrame('hls'));
      break;
    }
    case 'll_hls':{
      setPreviewFrame(prev,formatFrame('ll_hls'));
      break;
    }
    case 'dash':{
      setPreviewFrame(prev,formatFrame('dash'));
      break;
    }
    case 'flv':{
      setPreviewFrame(prev,formatFrame('flv'));
      break;
    }
    case 'mp4':
    {
      setPreviewFrame(prev,formatFrame('mp4'));
      break;
    }
    case 'webm':
    {
      setPreviewFrame(prev,formatFrame('webm'));
      break;
    }
    case 'mp3':
    case 'aac':
    case 'ogg':
    case 'wav':
    case 'flac':
    case 'icecast':
      setPreviewFrame(prev,formatFrame(embedTab));
      break;
    case 'whep':
      setPreviewFallback(prev,previewURLs.embed||('/embed/'+key),'WHEP icin tarayici onizleme yerine standart player gosteriliyor.');
      break;
    case 'rtmp':
    case 'rtsp':
    case 'srt':
      setPreviewFallback(prev,previewURLs.embed||('/embed/'+key),'Bu protokol tarayicida dogrudan oynatilamaz. HLS onizleme gosteriliyor.');
      break;
    default:
      setPreviewFallback(prev,previewURLs.embed||('/embed/'+key),'Onizleme varsayilan player ile gosteriliyor.');
      break;
  }
}
async function updateEmbedPreview(){
  const key=document.getElementById('ae-stream')?document.getElementById('ae-stream').value:'';
  const w=document.getElementById('ae-width')?document.getElementById('ae-width').value:'1280';
  const h=document.getElementById('ae-height')?document.getElementById('ae-height').value:'720';
  const autoplay=document.getElementById('ae-autoplay')?document.getElementById('ae-autoplay').checked:true;
  const muted=document.getElementById('ae-muted')?document.getElementById('ae-muted').checked:true;
  const out=document.getElementById('embed-output');
  const prev=document.getElementById('embed-preview');
  if(!out||!key){if(out)out.innerHTML='<div style="color:var(--text-muted)">Once bir yayin secin</div>';return}

  var se=document.getElementById('ae-stream');
  var sData=se?JSON.parse(se.dataset.settings||'{}'):{};
  var streamName='';
  var policyRaw='';
  if(se&&se.selectedOptions&&se.selectedOptions[0]){
    streamName=se.selectedOptions[0].dataset.streamName||'';
    policyRaw=se.selectedOptions[0].dataset.policyJson||'';
  }
  try{
    const access=await getPlaybackAccess(key,sData,policyRaw);
    var urls=getAllURLs(key,sData,streamName||key,access);
    var previewURLs=getPreviewURLs(key,sData,streamName||key,access);
    var bundle=buildEmbedBundle(embedTab,key,urls,w,h,autoplay,muted);
    out.innerHTML=
      copyField(bundle.primaryLabel||'Embed Kodu',bundle.primary||'')+
      (bundle.direct?copyField(bundle.directLabel||'URL',bundle.direct):'')+
      (access&&access.needs_token?'<div class="form-hint" style="color:var(--warning);margin-top:8px">Token korumasi aktif. Uretilen preview ve linklere gecici playback token eklendi.</div>':'')+
      (bundle.note?'<div class="form-hint">'+bundle.note+'</div>':'');
    renderAdvancedPreview(prev,key,autoplay,muted,previewURLs);
  }catch(e){
    out.innerHTML='<div style="color:var(--danger)">Embed bilgileri yuklenemedi</div>';
    if(prev)prev.innerHTML='<div style="position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);color:var(--text-muted)">Onizleme yuklenemedi</div>';
  }
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â ANALYTICS Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
let analyticsPeriod='24h';
function analyticsLabelFormatter(period){
  return function(date,index,total){
    if(!(date instanceof Date)||Number.isNaN(date.getTime()))return '';
    if(period==='30d')return date.toLocaleDateString(localeForLang(),{day:'2-digit',month:'2-digit'});
    if(period==='7d')return date.toLocaleDateString(localeForLang(),{day:'2-digit',month:'2-digit'})+' '+date.toLocaleTimeString(localeForLang(),{hour:'2-digit',minute:'2-digit'});
    return date.toLocaleTimeString(localeForLang(),{hour:'2-digit',minute:'2-digit'});
  };
}
function analyticsMeta(period,history){
  if(history&&history.label){
    const points=history.points||((history.viewers&&history.viewers.length)||0);
    return history.label+(points?(' - '+points+' nokta'):'');
  }
  if(period==='30d')return 'Son 30 gun';
  if(period==='7d')return 'Son 7 gun';
  return 'Son 24 saat';
}
async function fetchAnalyticsHistory(period){
  return api('/api/analytics/history?period='+encodeURIComponent(period||'24h'));
}
function renderAnalyticsPeriodSelector(){
  return '<div class="segment-control">'+['24h','7d','30d'].map(function(period){
    const label=period==='24h'?'24 Saat':(period==='7d'?'7 Gun':'30 Gun');
    return '<button class="segment-btn '+(analyticsPeriod===period?'active':'')+'" onclick="setAnalyticsPeriod(\''+period+'\')">'+label+'</button>';
  }).join('')+'</div>';
}
function setAnalyticsPeriod(period){
  if(!period||analyticsPeriod===period)return;
  analyticsPeriod=period;
  navigate('analytics');
}
async function renderAnalytics(c){
  const [data,history]=await Promise.all([api('/api/analytics'),fetchAnalyticsHistory(analyticsPeriod)]);
  if(!data){c.innerHTML='<div class="empty-state"><div class="icon"><i class="bi bi-bar-chart-line"></i></div><h3>Analitik verisi yok</h3></div>';return}
  const fmtItems=Object.entries(data.viewers_by_format||{}).sort((a,b)=>b[1]-a[1]).map(([label,value])=>({label:label,value:value}));
  const countryItems=Object.entries(data.viewers_by_country||{}).sort((a,b)=>b[1]-a[1]).slice(0,8).map(([label,value])=>({label:label,value:value}));
  const historyViewers=(history&&Array.isArray(history.viewers)?history.viewers:[]);
  const historyBandwidth=(history&&Array.isArray(history.bandwidth)?history.bandwidth:[]);
  const viewerTimeline=historyViewers.length?historyViewers:(data.viewers_timeline||[]);
  const labelFormatter=analyticsLabelFormatter(analyticsPeriod);
  const maxPoints=analyticsPeriod==='30d'?30:(analyticsPeriod==='7d'?28:24);
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Analitik</h1><div>'+renderAnalyticsPeriodSelector()+'</div></div>'+
    '<div class="card-grid card-grid-4" style="margin-bottom:24px">'+
      statCard('purple','bi-collection-play',fmtInt(data.total_streams||0),'Toplam Yayin','streams','Olusturulan tum streamler')+
      statCard('green','bi-people-fill',fmtInt(data.current_viewers||0),'Aktif Izleyici','viewers','Su an acik oturumlar')+
      statCard('orange','bi-graph-up-arrow',fmtInt(data.peak_concurrent||0),'Tepe Esz.','viewers','Kaydedilen en yuksek eszamanli izleyici')+
      statCard('blue','bi-diagram-3',fmtBytes(data.total_bandwidth||0),'Toplam Bant','transcode-jobs','Sunucudan cikan toplam trafik')+
    '</div>'+
    '<div class="insight-grid">'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Izleyici Trendi</h3><span class="form-hint">'+escHtml((history&&history.label)||'Secili periyot')+'</span></div><div class="card-body">'+renderTimelineChart(viewerTimeline,'Henuz timeline verisi yok',function(v){return String(v)},{meta:analyticsMeta(analyticsPeriod,history),labelFormatter:labelFormatter,maxPoints:maxPoints,labelSlots:6,valueSlots:6})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Bant Trendi</h3><span class="form-hint">Ayni periyotta toplam cikis</span></div><div class="card-body">'+renderTimelineChart(historyBandwidth,'Henuz bant snapshot yok',function(v){return fmtBytes(v)},{meta:analyticsMeta(analyticsPeriod,history),labelFormatter:labelFormatter,maxPoints:maxPoints,labelSlots:6,valueSlots:5})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Format Dagilimi</h3></div><div class="card-body">'+renderBarList(fmtItems,'Henuz format verisi yok',function(v){return String(v)})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Ulke Dagilimi</h3></div><div class="card-body">'+renderBarList(countryItems,'Henuz ulke verisi yok',function(v){return String(v)})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">En Populer Yayinlar</h3></div><div class="card-body">'+
        ((data.top_streams||[]).length?(data.top_streams||[]).map(function(s){
          return '<div class="metric-row"><div><div class="setting-label">'+escHtml(s.stream_name||shortKey(s.stream_key))+'</div><div class="setting-desc"><code>'+escHtml(s.stream_key)+'</code></div></div><span class="badge">'+fmtInt(s.viewers||0)+' izleyici</span></div>';
        }).join(''):'<div style="color:var(--text-muted)">Aktif yayin yok</div>')+
      '</div></div>'+
    '</div>';
  schedulePageRefresh('analytics',15000);
}
function fmtBytes(b){if(!b||b===0)return '0 B';const k=1024,s=['B','KB','MB','GB','TB'];const i=Math.floor(Math.log(b)/Math.log(k));return (b/Math.pow(k,i)).toFixed(1)+' '+s[i]}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â RECORDINGS Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderRecordings(c){
  const [recsRes,streamsRes,savedRes,archivesRes,settings]=await Promise.all([
    api('/api/recordings'),
    api('/api/streams'),
    api('/api/recordings/library'),
    api('/api/recordings/archives'),
    api('/api/settings')
  ]);
  const recs=Array.isArray(recsRes)?recsRes:[];
  const streams=Array.isArray(streamsRes)?streamsRes:[];
  const saved=Array.isArray(savedRes)?savedRes:[];
  const archives=Array.isArray(archivesRes)?archivesRes:[];
  const archiveEnabled=settings&&settings.archive_enabled==='true';
  const archiveMap={};
  archives.forEach(function(item){
    archiveMap[item.stream_key+'::'+item.filename]=item;
  });
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Kayitlar</h1>'+
      '<div style="display:flex;gap:10px;flex-wrap:wrap"><button class="btn btn-primary" onclick="showRecordModal()">Kayit Baslat</button>'+(archiveEnabled?'<button class="btn btn-secondary" onclick="runArchiveSync()">Arsiv Senkronu</button>':'')+'</div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-title" style="margin-bottom:8px">Depolama Notu</div>'+
      '<div class="form-hint">Kalici kayitlar <code>data/recordings</code> altindadir. <code>data/hls</code> ve <code>data/transcode/hls</code> dizinleri canli yayin cache alanidir; kayit olarak listelenmez. Object storage aciksa bu ekrandan arsive gonderip geri yukleyebilirsiniz.</div>'+
    '</div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Aktif Kayitlar</h3></div>'+
    '<div class="card-body"><table class="table"><thead><tr><th>ID</th><th>Yayin</th><th>Format</th><th>Durum</th><th>Boyut</th><th>Islem</th></tr></thead><tbody id="rec-list"></tbody></table></div></div>'+
    '<div class="card-grid card-grid-2" style="margin-bottom:16px">'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Kayit Kutuphanesi</h3><span class="form-hint">Yerelde bulunan dosyalar</span></div>'+
      '<div class="card-body"><table class="table"><thead><tr><th>Yayin</th><th>Dosya</th><th>Format</th><th>Tarih</th><th>Boyut</th><th>Arsiv</th><th>Islem</th></tr></thead><tbody id="saved-rec-list"></tbody></table></div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Secili Kayit Onizleme</h3><span class="form-hint">Listeden bir kayit secin</span></div>'+
      '<div class="card-body"><div id="recording-preview-panel"><div class="empty-state"><div class="icon"><i class="bi bi-film"></i></div><h3>Kayit secin</h3><p style="color:var(--text-muted)">Panel ayni sayfada secili kaydi oynatir.</p></div></div></div></div>'+
    '</div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Arsiv Kutuphanesi</h3><span class="form-hint">Object storage veya lokal arsivdeki kayitlar</span></div><div class="card-body"><table class="table"><thead><tr><th>Yayin</th><th>Dosya</th><th>Saglayici</th><th>Tarih</th><th>Yerel Durum</th><th>Sonuc</th><th>Islem</th></tr></thead><tbody id="archive-rec-list"></tbody></table></div></div>'+
    '<div id="rec-modal" style="display:none"></div>';
  const rl=document.getElementById('rec-list');
  if(rl){
    rl.innerHTML=recs.length?recs.map(r=>'<tr><td style="font-size:12px">'+r.ID+'</td><td>'+r.StreamKey+'</td><td>'+r.Format+'</td><td><span class="badge badge-'+(r.Status==='recording'?'green':'gray')+'">'+r.Status+'</span></td><td>'+fmtBytes(r.Size||0)+'</td><td>'+(r.Status==='recording'?'<button class="btn btn-sm btn-danger" onclick="stopRec(\''+r.ID+'\')">Durdur</button>':'\u2014')+'</td></tr>').join(''):'<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Aktif kayit yok</td></tr>';
  }
  const srl=document.getElementById('saved-rec-list');
  if(srl){
    srl.innerHTML=saved.length?saved.map(function(r,i){
      const archiveInfo=archiveMap[r.stream_key+'::'+r.name];
      const archiveBadge=archiveInfo?renderArchiveStatusBadge(archiveInfo):'<span class="tag tag-blue">Yerelde</span>';
      return '<tr>'+
        '<td><code>'+escHtml(r.stream_key)+'</code></td>'+
        '<td>'+escHtml(r.name)+'</td>'+
        '<td>'+(r.format||'-').toUpperCase()+'</td>'+
        '<td>'+fmtLocaleDateTime(r.mod_time)+'</td>'+
        '<td>'+fmtBytes(r.size||0)+'</td>'+
        '<td>'+archiveBadge+'</td>'+
        '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
          '<button class="btn btn-sm btn-secondary" onclick=\'previewRecordingPanel('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+','+JSON.stringify(r.format||'')+','+JSON.stringify(r.mod_time||'')+','+(r.size||0)+')\'>Onizle</button>'+
          '<button class="btn btn-sm btn-secondary" onclick=\'downloadRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Indir</button>'+
          (archiveEnabled?'<button class="btn btn-sm btn-secondary" onclick=\'archiveRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>'+(archiveInfo&&archiveInfo.status==='archived'?'Yeniden Arsivle':'Arsive Gonder')+'</button>':'')+
          (archiveInfo&&archiveInfo.object_url?'<button class="btn btn-sm btn-secondary" onclick=\'window.open('+JSON.stringify(archiveInfo.object_url)+',"_blank")\'>Arsiv Linki</button>':'')+
          '<button class="btn btn-sm btn-danger" onclick=\'deleteRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Sil</button>'+
        '</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="7" style="text-align:center;color:var(--text-muted);padding:24px">Kaydedilmis dosya yok</td></tr>';
  }
  const arl=document.getElementById('archive-rec-list');
  if(arl){
    arl.innerHTML=archives.length?archives.map(function(item){
      const localState=item.local_deleted?'<span class="tag tag-yellow">Yerelde yok</span>':'<span class="tag tag-green">Yerelde var</span>';
      const statusBadge=renderArchiveStatusBadge(item);
      return '<tr>'+
        '<td><code>'+escHtml(item.stream_key)+'</code></td>'+
        '<td>'+escHtml(item.filename)+'</td>'+
        '<td>'+escHtml(String(item.provider||'-').toUpperCase())+'</td>'+
        '<td>'+fmtLocaleDateTime(item.archived_at||item.updated_at||item.created_at)+'</td>'+
        '<td>'+localState+'</td>'+
        '<td>'+statusBadge+(item.last_error?'<div class="setting-desc" style="max-width:320px">'+escHtml(item.last_error)+'</div>':'')+'</td>'+
        '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
          '<button class="btn btn-sm btn-secondary" onclick=\'restoreRecordingArchive('+JSON.stringify(item.stream_key)+','+JSON.stringify(item.filename)+')\'>Geri Yukle</button>'+
          (item.object_url?'<button class="btn btn-sm btn-secondary" onclick=\'window.open('+JSON.stringify(item.object_url)+',"_blank")\'>Arsiv Linki</button>':'')+
          '</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="7" style="text-align:center;color:var(--text-muted);padding:24px">Arsivlenmis kayit yok</td></tr>';
  }
  window._recStreams=streams;
  window._savedRecordings=saved;
  window._recordingArchives=archives;
  if(saved.length){
    previewRecordingPanel(saved[0].stream_key,saved[0].name,saved[0].format||'',saved[0].mod_time||'',saved[0].size||0);
  }
}
function renderArchiveStatusBadge(item){
  if(!item)return '<span class="tag tag-blue">Yerelde</span>';
  const status=String(item.status||'archived').toLowerCase();
  if(status==='error')return '<span class="tag tag-red">Hata</span>';
  if(item.local_deleted)return '<span class="tag tag-yellow">Arsivde</span>';
  return '<span class="tag tag-green">Arsivlendi</span>';
}
function renderBackupArchiveStatusBadge(item){
  if(!item)return '<span class="tag tag-blue">Yerelde</span>';
  const status=String(item.status||'archived').toLowerCase();
  if(status==='error')return '<span class="tag tag-red">Hata</span>';
  if(item.local_deleted)return '<span class="tag tag-yellow">Arsivde</span>';
  return '<span class="tag tag-green">Arsivlendi</span>';
}
function normalizeStorageSnapshot(settings,report,archivesRes,recsRes,streamsRes,savedRes,backupsRes,backupArchivesRes,remuxJobsRes){
  const archives=Array.isArray(archivesRes)?archivesRes:[];
  const recs=Array.isArray(recsRes)?recsRes:[];
  const streams=Array.isArray(streamsRes)?streamsRes:[];
  const saved=Array.isArray(savedRes)?savedRes:[];
  const backups=(backupsRes&&Array.isArray(backupsRes.items))?backupsRes.items:[];
  const backupArchives=Array.isArray(backupArchivesRes)?backupArchivesRes:[];
  const remuxJobs=Array.isArray(remuxJobsRes)?remuxJobsRes:[];
  const archiveMap={};
  const backupArchiveMap={};
  archives.forEach(function(item){archiveMap[item.stream_key+'::'+item.filename]=item;});
  backupArchives.forEach(function(item){backupArchiveMap[item.name]=item;});
  return {
    settings:settings||{},
    report:report||{},
    archiveSummary:(report&&report.storage&&report.storage.archive)?report.storage.archive:{},
    archives:archives,
    recs:recs,
    activeRecordings:recs.filter(function(item){return item&&item.Status==='recording';}),
    streams:streams,
    saved:saved,
    backups:backups,
    backupArchives:backupArchives,
    remuxJobs:remuxJobs,
    archiveMap:archiveMap,
    backupArchiveMap:backupArchiveMap,
    archiveEnabled:settings&&settings.archive_enabled==='true',
    backupArchiveEnabled:settings&&settings.backup_archive_enabled==='true'
  };
}
function storageProviderCatalog(){
  return [
    {engine:'local',variant:'local',title:'Yerel Disk',icon:'bi-hdd-stack',desc:'Dosyalari bu sunucuda ikinci bir klasorde tutar.',hint:'Tek sunuculu ve en kolay baslangic secenegi.',family:'local'},
    {engine:'s3',variant:'aws_s3',title:'AWS S3',icon:'bi-cloud-fill',desc:'Amazon S3 bucket hedefine yazar.',hint:'Gercek dis felaket yedegi icin en guclu seceneklerden biri.',family:'s3'},
    {engine:'minio',variant:'minio',title:'MinIO',icon:'bi-server',desc:'Kendi S3 uyumlu MinIO sunucuna yazar.',hint:'Kendi objeni barindiriyorsan idealdir.',family:'s3'},
    {engine:'s3',variant:'cloudflare_r2',title:'Cloudflare R2',icon:'bi-cloud-fill',desc:'R2 bucket hedefine S3 uyumlu yazim yapar.',hint:'Egress maliyetini dusurmek isteyenlerde cok populer.',family:'s3'},
    {engine:'s3',variant:'backblaze_b2',title:'Backblaze B2',icon:'bi-cloud-fill',desc:'B2 S3 uyumlu bucket hedefine yazar.',hint:'Dusuk maliyetli dis backup icin sik kullanilir.',family:'s3'},
    {engine:'s3',variant:'wasabi',title:'Wasabi',icon:'bi-cloud-fill',desc:'Wasabi bucket hedefine yazar.',hint:'Uzun sureli arsiv icin populer bir secenek.',family:'s3'},
    {engine:'s3',variant:'digitalocean_spaces',title:'DO Spaces',icon:'bi-cloud-fill',desc:'DigitalOcean Spaces bucket hedefine yazar.',hint:'Kucuk ekipler icin pratik bir S3 uyumlu secim.',family:'s3'},
    {engine:'s3',variant:'linode_object_storage',title:'Linode Object',icon:'bi-cloud-fill',desc:'Linode object storage hedefine yazar.',hint:'S3 uyumlu baska bir uygun maliyetli secenek.',family:'s3'},
    {engine:'s3',variant:'scaleway_object_storage',title:'Scaleway Object',icon:'bi-cloud-fill',desc:'Scaleway object storage hedefine yazar.',hint:'Avrupa odakli bir S3 uyumlu alternatiftir.',family:'s3'},
    {engine:'s3',variant:'idrive_e2',title:'IDrive e2',icon:'bi-cloud-fill',desc:'IDrive e2 bucket hedefine yazar.',hint:'S3 uyumlu uygun maliyetli bir depolama secenegi.',family:'s3'},
    {engine:'sftp',variant:'sftp',title:'SFTP Sunucu',icon:'bi-folder-symlink',desc:'Dosyalari baska bir Linux sunucusuna klasor gibi kopyalar.',hint:'Dusuk butcede en pratik uzak hedeflerden biridir.',family:'sftp'},
    {engine:'rclone',variant:'google_drive',title:'Google Drive',icon:'bi-cloud-fill',desc:'Google Drive klasorune yukler.',hint:'Rclone baglanti profiliyle calisir.',family:'rclone'},
    {engine:'rclone',variant:'onedrive',title:'OneDrive',icon:'bi-cloud-fill',desc:'Microsoft OneDrive hedefine yukler.',hint:'Rclone baglanti profiliyle calisir.',family:'rclone'},
    {engine:'rclone',variant:'dropbox',title:'Dropbox',icon:'bi-cloud-fill',desc:'Dropbox klasorune yukler.',hint:'Rclone baglanti profiliyle calisir.',family:'rclone'},
    {engine:'rclone',variant:'google_cloud_storage',title:'Google Cloud Storage',icon:'bi-cloud-fill',desc:'GCS bucket benzeri hedefine yazar.',hint:'Rclone baglantisi ile calisir.',family:'rclone'},
    {engine:'rclone',variant:'azure_blob',title:'Azure Blob',icon:'bi-cloud-fill',desc:'Azure Blob container hedefine yazar.',hint:'Rclone baglantisi ile calisir.',family:'rclone'},
    {engine:'rclone',variant:'box',title:'Box',icon:'bi-cloud-fill',desc:'Box klasorune yukler.',hint:'Rclone baglanti profili ile calisir.',family:'rclone'},
    {engine:'rclone',variant:'pcloud',title:'pCloud',icon:'bi-cloud-fill',desc:'pCloud klasorune yukler.',hint:'Rclone baglanti profili ile calisir.',family:'rclone'},
    {engine:'rclone',variant:'mega',title:'MEGA',icon:'bi-cloud-fill',desc:'MEGA klasorune yukler.',hint:'Rclone baglanti profili ile calisir.',family:'rclone'},
    {engine:'rclone',variant:'nextcloud',title:'Nextcloud',icon:'bi-cloud-fill',desc:'Nextcloud / ownCloud deposuna yazar.',hint:'Rclone veya WebDAV profili ile calisir.',family:'rclone'},
    {engine:'rclone',variant:'webdav',title:'WebDAV',icon:'bi-cloud-fill',desc:'WebDAV tabanli genel bulut hedeflerine yazar.',hint:'Kurumsal dosya servisleri icin genis uyumluluk saglar.',family:'rclone'}
  ];
}
function storageResolveProvider(engine, variant){
  engine=String(engine||'').toLowerCase().trim();
  variant=String(variant||'').toLowerCase().trim();
  if(engine)return engine;
  if(['aws_s3','cloudflare_r2','backblaze_b2','wasabi','digitalocean_spaces','linode_object_storage','scaleway_object_storage','idrive_e2'].indexOf(variant)>=0)return 's3';
  if(variant==='minio')return 'minio';
  if(variant==='sftp')return 'sftp';
  if(['google_drive','onedrive','dropbox','google_cloud_storage','azure_blob','box','pcloud','mega','nextcloud','webdav'].indexOf(variant)>=0)return 'rclone';
  return 'local';
}
function storageDefaultVariant(engine){
  switch(storageResolveProvider(engine,'')){
    case 's3': return 'aws_s3';
    case 'minio': return 'minio';
    case 'sftp': return 'sftp';
    case 'rclone': return 'google_drive';
    default: return 'local';
  }
}
function storageProviderInfo(engine, variant){
  const resolvedEngine=storageResolveProvider(engine,variant);
  const resolvedVariant=String(variant||storageDefaultVariant(resolvedEngine)).toLowerCase();
  const found=storageProviderCatalog().find(function(item){
    return item.engine===resolvedEngine && item.variant===resolvedVariant;
  });
  if(found)return found;
  return storageProviderCatalog().find(function(item){return item.engine===resolvedEngine;})||storageProviderCatalog()[0];
}
function storagePrettyLabel(raw){
  const value=String(raw||'').trim().toLowerCase();
  if(!value)return '-';
  const info=storageProviderCatalog().find(function(item){
    return item.variant===value || (item.engine===value && item.variant===storageDefaultVariant(item.engine));
  });
  if(info)return info.title;
  return value.replace(/_/g,' ').replace(/\b\w/g,function(ch){return ch.toUpperCase();});
}
function storageRoleMeta(role){
  return role==='backups'
    ?{prefix:'backup_archive',title:'Sistem Yedegi Hedefi',syncAction:'runBackupArchiveSync()',syncLabel:'Yedekleri Simdi Gonder',testLabel:'Yedek hedefini test et'}
    :{prefix:'archive',title:'Kayit Arsiv Hedefi',syncAction:'runArchiveSync()',syncLabel:'Kayitlari Simdi Gonder',testLabel:'Kayit hedefini test et'};
}
function readStorageDraftSettings(){
  const base=Object.assign({},(window._storageData&&window._storageData.settings)||{});
  document.querySelectorAll('.setting-input').forEach(function(el){
    const key=el.dataset.key;
    if(!key)return;
    if(el.type==='checkbox')base[key]=el.checked?'true':'false';
    else base[key]=el.value;
  });
  return base;
}
function storageTargetState(role, settings){
  const meta=storageRoleMeta(role);
  const sameTarget=role==='backups' && String(settings.backup_archive_use_same_target||'true')!=='false';
  const sourcePrefix=sameTarget?'archive':meta.prefix;
  const savedEngine=String(settings[meta.prefix+'_provider']||'').toLowerCase();
  const savedVariant=String(settings[meta.prefix+'_provider_variant']||'').toLowerCase();
  const engine=storageResolveProvider(settings[sourcePrefix+'_provider']||'', settings[sourcePrefix+'_provider_variant']||'');
  const variant=String(settings[sourcePrefix+'_provider_variant']||storageDefaultVariant(engine)).toLowerCase();
  return {
    role:role,
    meta:meta,
    sameTarget:sameTarget,
    sourcePrefix:sourcePrefix,
    engine:engine,
    variant:variant,
    savedEngine:savedEngine||engine,
    savedVariant:savedVariant||variant,
    provider:storageProviderInfo(engine,variant),
    mode:String(settings.archive_ui_mode||'simple').toLowerCase()==='advanced'?'advanced':'simple'
  };
}
function renderStorageHiddenField(key,value){
  return '<input type="hidden" class="setting-input" data-key="'+escHtml(key)+'" value="'+escHtml(String(value==null?'':value))+'">';
}
function renderStorageCenter(data,restoreCmd){
  const settings=(data&&data.settings)||{};
  return '<div class="storage-shell">'+
    '<div class="page-header"><h1 class="page-title">Depolama ve Arsiv Merkezi</h1><div style="display:flex;gap:10px;flex-wrap:wrap"><button class="btn btn-primary" onclick="showRecordModal()">Kayit Baslat</button><button class="btn btn-secondary" onclick="createSystemBackupFromStorage(false)">Hafif Yedek Al</button><button class="btn btn-secondary" onclick="createSystemBackupFromStorage(true)">Kayitlarla Yedek Al</button><button class="btn btn-secondary" onclick="loadPage(\'maintenance-center\')">Bakim ve Servis</button></div></div>'+
    '<div class="studio-alert info" style="margin-bottom:16px"><strong>Rol ayrimi</strong><div style="margin-top:8px" class="form-hint">Bu merkez kayit, arsiv ve harici hedefleri yonetir. Servis durumu, upgrade ve offline restore komutlari icin <strong>Bakim ve Yedek</strong> sayfasini kullan.</div></div>'+
    '<div id="storage-active-banner"></div>'+
    '<div id="storage-remux-jobs"></div>'+
    '<div class="storage-hero">'+
      '<div class="card"><div class="card-title" style="margin-bottom:12px">Kolay secim rehberi</div><div class="storage-choice-grid">'+
        '<div class="storage-choice-card active"><div class="icon"><i class="bi bi-1-circle"></i></div><div class="title">Yerelde basla</div><div class="desc">Kayitlar once bu sunucuda dursun. En hizli ve risksiz kurulum boyle olur.</div></div>'+
        '<div class="storage-choice-card active"><div class="icon"><i class="bi bi-2-circle"></i></div><div class="title">Ikinci kopya ekle</div><div class="desc">Kayitlari ve sistem yedeklerini uzak hedefe gondererek disk ve felaket riskini dusur.</div></div>'+
        '<div class="storage-choice-card active"><div class="icon"><i class="bi bi-3-circle"></i></div><div class="title">Geri yukle</div><div class="desc">Yerelde olmayan dosyayi arsivden geri getir, gerekirse indir ve tekrar kullan.</div></div>'+
      '</div><div class="storage-provider-note" style="margin-top:12px">MP4 kayitlar varsayilan olarak onerilir. Tarayici onizleme, indirme ve arsiv akisi icin en sorunsuz secim budur.</div></div>'+
      '<div class="card"><div class="card-title" style="margin-bottom:12px">Durum Ozeti</div><div class="storage-kpi-list">'+
        '<div class="storage-kpi"><strong>'+fmtInt((data&&data.saved||[]).length)+'</strong><span>Yerel kayit dosyasi</span></div>'+
        '<div class="storage-kpi"><strong>'+fmtInt((data&&data.archives||[]).length)+'</strong><span>Kayit arsiv kopyasi</span></div>'+
        '<div class="storage-kpi"><strong>'+fmtInt((data&&data.backups||[]).length)+'</strong><span>Yerel sistem yedegi</span></div>'+
        '<div class="storage-kpi"><strong>'+fmtInt((data&&data.backupArchives||[]).length)+'</strong><span>Yedek arsiv kopyasi</span></div>'+
      '</div><div class="storage-provider-note" style="margin-top:14px">Restore komutu: <code>'+escHtml(restoreCmd)+'</code></div></div>'+
    '</div>'+
    '<div id="storage-stats-grid" class="card-grid card-grid-4" style="margin-bottom:0"></div>'+
    '<div id="storage-settings-shell">'+renderStorageSettingsShell(settings,data)+'</div>'+
    '<div class="card" style="margin-top:16px;margin-bottom:16px"><div class="card-header"><h3 class="card-title">Aktif Kayitlar</h3><span class="form-hint" id="storage-active-count">0 aktif oturum</span></div><div class="card-body"><table class="table"><thead><tr><th>ID</th><th>Yayin</th><th>Format</th><th>Durum</th><th>Boyut</th><th style="white-space:nowrap">Islem</th></tr></thead><tbody id="rec-list"></tbody></table></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Secili Kayit Onizleme</h3><span class="form-hint">Ham TS/FLV/MKV kayitlarda once MP4 hazirlamak daha kararlidir.</span></div><div class="card-body"><div id="recording-preview-panel"><div class="empty-state"><div class="icon"><i class="bi bi-film"></i></div><h3>Kayit secin</h3><p style="color:var(--text-muted)">Panel secili kaydi ayni sayfada oynatir.</p></div></div></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Kayit Kutuphanesi</h3><span class="form-hint">Yereldeki kayitlar, indirme ve donusum aksiyonlari</span></div><div class="card-body"><table class="table"><thead><tr><th>Yayin</th><th>Dosya</th><th>Format</th><th>Tarih</th><th>Boyut</th><th>Arsiv</th><th>Islem</th></tr></thead><tbody id="saved-rec-list"></tbody></table></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Kayit Arsiv Kutuphanesi</h3><span class="form-hint">Harici hedefte bulunan kayitlar</span></div><div class="card-body"><table class="table"><thead><tr><th>Yayin</th><th>Dosya</th><th>Saglayici</th><th>Tarih</th><th>Yerel Durum</th><th>Sonuc</th><th>Islem</th></tr></thead><tbody id="archive-rec-list"></tbody></table></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Sistem Yedekleri</h3><span class="form-hint">Yerelde duran sistem yedekleri</span></div><div class="card-body"><table class="table"><thead><tr><th>Dosya</th><th>Boyut</th><th>Tarih</th><th>Tur</th><th>Arsiv</th><th>Islem</th></tr></thead><tbody id="system-backup-list"></tbody></table></div></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Yedek Arsiv Kutuphanesi</h3><span class="form-hint">Harici hedefte saklanan sistem yedekleri</span></div><div class="card-body"><table class="table"><thead><tr><th>Dosya</th><th>Saglayici</th><th>Tarih</th><th>Yerel Durum</th><th>Sonuc</th><th>Islem</th></tr></thead><tbody id="backup-archive-list"></tbody></table></div></div>'+
    '<div id="rec-modal" style="display:none"></div>'+
  '</div>';
}
function renderStorageSettingsShell(settings,data){
  return '<div class="card-grid card-grid-2" style="margin-top:16px">'+
    '<div class="card">'+renderStorageModeCard(settings,data)+'</div>'+
    '<div class="card">'+renderStorageLocalMaintenanceCard(settings)+'</div>'+
  '</div>'+
  '<div class="card-grid card-grid-2" style="margin-top:16px">'+
    '<div id="storage-recordings-panel">'+renderStorageTargetPanel('recordings',settings,data)+'</div>'+
    '<div id="storage-backups-panel">'+renderStorageTargetPanel('backups',settings,data)+'</div>'+
  '</div>';
}
function renderStorageModeCard(settings,data){
  const mode=String(settings.archive_ui_mode||'simple').toLowerCase()==='advanced'?'advanced':'simple';
  const sameTarget=String(settings.backup_archive_use_same_target||'true')!=='false';
  return renderStorageHiddenField('archive_ui_mode',mode)+renderStorageHiddenField('backup_archive_use_same_target',sameTarget?'true':'false')+
    '<div class="storage-target-shell"><div class="storage-target-top"><div><div class="card-title">Kullanim Modu ve Mimari</div><div class="storage-provider-note">Basit mod gundelik isler icin sade alanlar gosterir. Gelismis mod, endpoint, path-style, rclone config ve yasam dongusu gibi teknik ayrintilari acar.</div></div><div class="segment-control"><button class="segment-btn '+(mode==='simple'?'active':'')+'" onclick="switchStorageUIMode(\'simple\')">Basit Mod</button><button class="segment-btn '+(mode==='advanced'?'active':'')+'" onclick="switchStorageUIMode(\'advanced\')">Gelismis Mod</button></div></div>'+
      '<div class="storage-subtle"><div class="setting-row"><div><div class="setting-label">Kayıt ve yedek ayni hedefi kullansin</div><div class="setting-desc">Acik olursa sistem yedekleri kayitlarla ayni cloud veya uzak klasore gider. Kaparsaniz yedekler icin ayri hedef secebilirsiniz.</div></div><label class="toggle"><input type="checkbox" '+(sameTarget?'checked':'')+' onchange="toggleBackupTargetMode(this.checked)"><span class="toggle-slider"></span></label></div></div>'+
      '<div class="storage-provider-note">S3 ailesi: AWS S3, MinIO, R2, Backblaze B2, Wasabi, DigitalOcean Spaces, Linode, Scaleway, IDrive e2. Rclone ailesi: Google Drive, OneDrive, Dropbox, GCS, Azure Blob, Box, pCloud, MEGA, Nextcloud ve WebDAV.</div>'+
      '<div class="storage-test-row"><button class="btn btn-primary" onclick="saveStorageCenter()">Tum Ayarlari Kaydet</button><button class="btn btn-secondary" onclick="refreshStorageSnapshot({resetPreview:false})">Yeniden Tara</button></div></div>';
}
function renderStorageLocalMaintenanceCard(settings){
  return '<div class="storage-target-shell"><div class="card-title">Yerel Depolama ve Temizlik</div><div class="storage-provider-note">Yerel kopya her zaman korunur. Uzak hedef ekleseniz bile burada ne kadar veri saklanacagini ve bakimin nasil calisacagini belirleyebilirsiniz.</div><div class="storage-form-grid">'+
    settingInput('storage_max_gb','Maksimum Depolama (GB)',settings.storage_max_gb||'50','number','Disk ve saglik uyari esiklerinde kullanilir.')+
    settingInput('recordings_retention_days','Kayit Saklama Suresi (gun)',settings.recordings_retention_days||'30','number','0 verilirse otomatik silinmez.')+
    settingInput('recordings_keep_latest','Yayin Basina Sakla',settings.recordings_keep_latest||'10','number','Her yayindan tutulacak son kayit sayisi.')+
    settingInput('storage_auto_clean','Uyumluluk Temizlik Gun Sayisi',settings.storage_auto_clean||'30','number','Eski davranis uyumlulugu icin korunur.')+
  '</div><div class="setting-row"><div><div class="setting-label">Otomatik Bakim</div><div class="setting-desc">Temizlik, telemetry ve trim bakimlarini periyodik olarak calistirir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="maintenance_auto_cleanup" '+(settings.maintenance_auto_cleanup!=='false'?'checked':'')+'><span class="toggle-slider"></span></label></div></div>';
}
function renderStorageTargetPanel(role,settings,data){
  const target=storageTargetState(role,settings);
  const mode=target.mode;
  const meta=target.meta;
  const summary=(data&&data.archiveSummary)||{};
  const sameTargetNote=role==='backups'&&target.sameTarget?'<div class="storage-subtle">Sistem yedekleri su an kayitlarla ayni bulut hedefini kullanir. Ayri hedef istiyorsaniz yukaridan bu secenegi kapatin.</div>':'';
  return '<div class="card storage-target-shell">'+
    renderStorageHiddenField(meta.prefix+'_provider',target.savedEngine)+
    renderStorageHiddenField(meta.prefix+'_provider_variant',target.savedVariant)+
    '<div class="storage-target-top"><div><div class="card-title">'+meta.title+'</div><div class="storage-provider-note">'+escHtml(target.provider.desc+' '+target.provider.hint)+'</div></div><div class="storage-target-meta">'+renderStorageTargetBadges(role,target,summary)+'</div></div>'+
    sameTargetNote+
    (role==='backups'&&target.sameTarget?'':'<div><div class="storage-provider-note" style="margin-bottom:10px">Bir hedef karti secin. Kullanici bulut markasini gorur; sistem alt tarafta dogru motoru kendisi kullanir.</div><div class="storage-provider-grid">'+renderStorageProviderCards(role,target)+'</div></div>')+
    renderStorageTargetFields(role,target,settings,mode)+
    '<div class="storage-test-row"><button class="btn btn-primary" onclick="saveStorageCenter()">Kaydet</button><button class="btn btn-secondary" onclick="testStorageTarget(\''+role+'\')">'+meta.testLabel+'</button><button class="btn btn-secondary" onclick="'+meta.syncAction+'">'+meta.syncLabel+'</button></div>'+
    '<div id="storage-test-result-'+role+'">'+renderStorageTestResult(role)+'</div>'+
  '</div>';
}
function renderStorageProviderCards(role,target){
  return storageProviderCatalog().map(function(item){
    const active=item.engine===target.engine && item.variant===target.variant;
    return '<div class="storage-choice-card '+(active?'active':'')+'" onclick="selectStorageProvider(\''+role+'\',\''+item.engine+'\',\''+item.variant+'\')"><div class="icon"><i class="bi '+item.icon+'"></i></div><div class="title">'+escHtml(item.title)+'</div><div class="desc">'+escHtml(item.desc)+'</div></div>';
  }).join('');
}
function renderStorageTargetBadges(role,target,summary){
  const configured=role==='backups'?!!summary.backup_configured:!!summary.recording_configured;
  const provider=role==='backups'?(summary.backup_provider||target.provider.title):(summary.recording_provider||target.provider.title);
  const schedule=role==='backups'?(summary.backup_schedule||'weekly'):(summary.recording_schedule||'immediate');
  const lastSync=role==='backups'?summary.last_backup_sync_at:summary.last_sync_at;
  return '<span class="storage-pill">'+escHtml(storagePrettyLabel(provider||target.provider.title))+'</span>'+
    '<span class="storage-pill">'+(configured?'Hazir':'Eksik ayar')+'</span>'+
    '<span class="storage-pill">'+escHtml(String(schedule||'-')).toUpperCase()+'</span>'+
    '<span class="storage-pill">'+(lastSync?escHtml(fmtLocaleDateTime(lastSync)):'Henuz senkron yok')+'</span>';
}
function renderStorageTargetFields(role,target,settings,mode){
  const prefix=target.meta.prefix;
  const sourcePrefix=target.sourcePrefix;
  const advanced=mode==='advanced';
  const isBackup=role==='backups';
  const scheduleValue=settings[prefix+'_schedule']||(isBackup?'weekly':'immediate');
  const tierValue=settings[prefix+'_target_tier']||(isBackup?'cold':'standard');
  const common='<div class="storage-form-grid">'+
    '<div class="setting-row"><div><div class="setting-label">'+(isBackup?'Yedek arsivi etkin':'Kayit arsivi etkin')+'</div><div class="setting-desc">'+(isBackup?'Yeni sistem yedekleri secilen hedefe gidebilir.':'Yeni kayitlar secilen hedefe gidebilir.')+'</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="'+prefix+'_enabled" '+(settings[prefix+'_enabled']==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
    '<div class="setting-row"><div><div class="setting-label">Otomatik yukle</div><div class="setting-desc">Yeni dosyalar zamanlama kuralina gore otomatik gonderilir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="'+prefix+'_auto_upload" '+(settings[prefix+'_auto_upload']==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
    '<div class="setting-row"><div><div class="setting-label">Yukleme sonrasi yereli sil</div><div class="setting-desc">Basarili kopya sonrasi yerel dosyayi temizler.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="'+prefix+'_delete_local_after_upload" '+(settings[prefix+'_delete_local_after_upload']==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
    settingSelect(prefix+'_schedule','Gonderim sikligi',scheduleValue,[{value:'manual',label:'Elle'},{value:'immediate',label:'Yeni dosya geldikce'},{value:'hourly',label:'Saatlik'},{value:'daily',label:'Gunluk'},{value:'weekly',label:'Haftalik'}],'Otomatik yukleme aciksa hangi ritimde calisacagini belirler.')+
    settingInput(prefix+'_scan_interval_minutes','Tarama araligi (dk)',settings[prefix+'_scan_interval_minutes']||(isBackup?'30':'10'),'number','Yeni dosyalari hangi aralikla kontrol edecegi.')+
    settingInput(prefix+'_batch_size','Tek turda islenecek oge',settings[prefix+'_batch_size']||(isBackup?'2':'3'),'number','Bir turda kac dosya gonderilecegi.')+
    settingSelect(prefix+'_target_tier','Arsiv seviyesi',tierValue,[{value:'standard',label:'Standart'},{value:'cold',label:'Soguk arsiv'},{value:'hot',label:'Hizli erisim'}],'Hedefin kullanim tipini operator tarafinda isaretlemek icindir.')+
    settingInput(prefix+'_cold_after_days','Soguk arsive gecis gunu',settings[prefix+'_cold_after_days']||(isBackup?'7':'30'),'number','Yasam dongusu planlamasi icin saklanir.')+
  '</div>';
  let providerFields='';
  if(target.provider.family==='local'){
    providerFields=settingInput(sourcePrefix+'_local_dir','Yerel arsiv klasoru',settings[sourcePrefix+'_local_dir']||'','text','Dosyalar bu sunucuda bu klasore kopyalanir.');
  }else if(target.provider.family==='s3'){
    providerFields='<div class="storage-form-grid">'+
      settingInput(sourcePrefix+'_endpoint','Baglanti adresi',settings[sourcePrefix+'_endpoint']||'','text','Ornek: https://s3.eu-central-1.amazonaws.com veya servisinizin endpoint adresi.')+
      settingInput(sourcePrefix+'_region','Bolge',settings[sourcePrefix+'_region']||'us-east-1','text','S3 uyumlu servisler icin imzalama bolgesi.')+
      settingInput(sourcePrefix+'_bucket','Bucket / alan adi',settings[sourcePrefix+'_bucket']||'','text','Yazma yapilacak container veya bucket adi.')+
      settingInput(sourcePrefix+'_access_key','Erisim anahtari',settings[sourcePrefix+'_access_key']||'','text','Servisin verdigi access key degeri.')+
      settingInput(sourcePrefix+'_secret_key','Gizli anahtar',settings[sourcePrefix+'_secret_key']||'','password','Servisin verdigi secret key degeri.')+
      (advanced?'<div class="setting-row"><div><div class="setting-label">Path style kullan</div><div class="setting-desc">MinIO, R2 ve bazi S3 uyumlu servislerde acik olmasi gerekir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="'+sourcePrefix+'_use_path_style" '+(settings[sourcePrefix+'_use_path_style']!=='false'?'checked':'')+'><span class="toggle-slider"></span></label></div>':'')+
    '</div>';
  }else if(target.provider.family==='sftp'){
    providerFields='<div class="storage-form-grid">'+
      settingInput(sourcePrefix+'_sftp_host','Sunucu adresi',settings[sourcePrefix+'_sftp_host']||'','text','Host adi veya IP adresi.')+
      settingInput(sourcePrefix+'_sftp_port','Port',settings[sourcePrefix+'_sftp_port']||'22','number','Genelde 22 olur.')+
      settingInput(sourcePrefix+'_sftp_user','Kullanici adi',settings[sourcePrefix+'_sftp_user']||'','text','Uzak sunucuda baglanacak hesap.')+
      settingInput(sourcePrefix+'_sftp_remote_dir','Uzak klasor',settings[sourcePrefix+'_sftp_remote_dir']||'','text','Ornek: /srv/fluxstream-archive')+
      settingInput(sourcePrefix+'_sftp_key_path','Anahtar dosya yolu',settings[sourcePrefix+'_sftp_key_path']||'','text','Bos ise sunucudaki varsayilan SSH anahtari denenir.')+
      (advanced?'<div class="setting-row"><div><div class="setting-label">Host key kontrolunu gevset</div><div class="setting-desc">Testte kolaylik saglar; production icin kapali tutmak daha guvenlidir.</div></div><label class="toggle"><input type="checkbox" class="setting-input" data-key="'+sourcePrefix+'_sftp_disable_host_key_check" '+(settings[sourcePrefix+'_sftp_disable_host_key_check']==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>':'')+
    '</div>';
  }else{
    providerFields='<div class="storage-form-grid">'+
      settingInput(sourcePrefix+'_rclone_remote','Bulut baglanti profili',settings[sourcePrefix+'_rclone_remote']||'','text','Sunucuda tanimli rclone profil adi. Ornek: gdrive-ana veya onedrive-media')+
      settingInput(sourcePrefix+'_rclone_path','Hedef klasor / yol',settings[sourcePrefix+'_rclone_path']||'','text','Bulut hedefinde yazilacak klasor. Ornek: fluxstream/recordings')+
      (advanced?settingInput(sourcePrefix+'_rclone_config_path','Rclone config yolu',settings[sourcePrefix+'_rclone_config_path']||'','text','Bos ise varsayilan rclone.conf kullanilir.'):'' )+
    '</div>';
  }
  const extras=(advanced?'<div class="storage-form-grid">'+
    settingInput(prefix+'_prefix','Object key kok klasoru',settings[prefix+'_prefix']||'fluxstream','text','Kayitlar ve yedekler bu kok klasor altina yazilir.')+
    settingInput(prefix+'_public_base_url','Public link tabani',settings[prefix+'_public_base_url']||'','text','Varsa panelde tiklanabilir public link uretir.')+
  '</div>':'');
  return '<div class="storage-subtle">'+providerFields+extras+common+'</div>';
}
function renderStorageTestResult(role){
  window._storageTestResults=window._storageTestResults||{};
  const result=window._storageTestResults[role];
  if(!result)return '';
  return '<div class="storage-subtle" style="margin-top:12px;border-color:'+(result.ok?'rgba(16,185,129,.25)':'rgba(239,68,68,.22)')+'"><div class="setting-label" style="margin-bottom:6px">'+(result.ok?'Baglanti testi basarili':'Baglanti testi hatali')+'</div><div class="storage-provider-note">'+escHtml(result.message||'')+'</div></div>';
}
async function saveStorageCenter(){
  await saveSettingsCategory('storage');
  await loadPage('settings-storage');
}
function switchStorageUIMode(mode){
  setStorageInputValue('archive_ui_mode',mode);
  updateStorageProviderUI();
}
function toggleBackupTargetMode(useSame){
  setStorageInputValue('backup_archive_use_same_target',useSame?'true':'false');
  updateStorageProviderUI();
}
function setStorageInputValue(key,value){
  const el=document.querySelector('.setting-input[data-key="'+key+'"]');
  if(!el)return;
  if(el.type==='checkbox')el.checked=String(value)==='true';
  else el.value=value;
}
function selectStorageProvider(role,engine,variant){
  const prefix=storageRoleMeta(role).prefix;
  setStorageInputValue(prefix+'_provider',engine);
  setStorageInputValue(prefix+'_provider_variant',variant);
  applyStorageProviderPreset(role,engine,variant);
  updateStorageProviderUI();
}
function applyStorageProviderPreset(role,engine,variant){
  const prefix=storageRoleMeta(role).prefix;
  const sourcePrefix=(role==='backups' && String(readStorageDraftSettings().backup_archive_use_same_target||'true')!=='false')?'archive':prefix;
  const presets={
    aws_s3:{endpoint:'https://s3.eu-central-1.amazonaws.com',region:'eu-central-1',pathStyle:'false'},
    minio:{endpoint:'http://127.0.0.1:9002',region:'us-east-1',pathStyle:'true'},
    cloudflare_r2:{endpoint:'https://ACCOUNT_ID.r2.cloudflarestorage.com',region:'auto',pathStyle:'true'},
    backblaze_b2:{endpoint:'https://s3.us-west-004.backblazeb2.com',region:'us-west-004',pathStyle:'true'},
    wasabi:{endpoint:'https://s3.eu-central-1.wasabisys.com',region:'eu-central-1',pathStyle:'true'},
    digitalocean_spaces:{endpoint:'https://fra1.digitaloceanspaces.com',region:'fra1',pathStyle:'true'},
    linode_object_storage:{endpoint:'https://eu-central-1.linodeobjects.com',region:'eu-central-1',pathStyle:'true'},
    scaleway_object_storage:{endpoint:'https://s3.fr-par.scw.cloud',region:'fr-par',pathStyle:'true'},
    idrive_e2:{endpoint:'https://storage.wdc.idrivee2-23.com',region:'us-east-1',pathStyle:'true'}
  };
  const preset=presets[variant];
  if(!preset)return;
  if(!document.querySelector('.setting-input[data-key="'+sourcePrefix+'_endpoint"]')?.value)setStorageInputValue(sourcePrefix+'_endpoint',preset.endpoint);
  if(!document.querySelector('.setting-input[data-key="'+sourcePrefix+'_region"]')?.value)setStorageInputValue(sourcePrefix+'_region',preset.region);
  setStorageInputValue(sourcePrefix+'_use_path_style',preset.pathStyle);
}
async function testStorageTarget(role){
  window._storageTestResults=window._storageTestResults||{};
  const updates=readStorageDraftSettings();
  try{
    const res=await api('/api/storage/connection-test',{method:'POST',body:{role:role,updates:updates}});
    if(res&&res.success===false){
      window._storageTestResults[role]={ok:false,message:res.message||'Baglanti testi basarisiz'};
      toast(res.message||'Baglanti testi basarisiz','error');
    }else{
      const result=(res&&res.result)||{};
      window._storageTestResults[role]={ok:true,message:'Test dosyasi yazildi ve silindi. Saglayici: '+String(result.provider||'-')};
      toast('Baglanti testi basarili');
    }
  }catch(e){
    window._storageTestResults[role]={ok:false,message:e.message||'Baglanti testi basarisiz'};
    toast('Baglanti testi basarisiz: '+e.message,'error');
  }
  updateStorageProviderUI();
}
function renderStorageActiveBanner(data){
  const activeRecordings=Array.isArray(data&&data.activeRecordings)?data.activeRecordings:[];
  if(!activeRecordings.length)return '';
  return '<div class="card" style="margin-bottom:16px;border-color:rgba(239,68,68,.28);box-shadow:0 8px 22px rgba(239,68,68,.08)"><div class="card-header"><div><div class="card-title">Aktif Kayit Uyarisi</div><div class="form-hint">Calisan kayit oturumlari burada sabit tutulur. Durdur dugmesine buradan da erisebilirsiniz.</div></div><span class="badge badge-live">'+fmtInt(activeRecordings.length)+' aktif kayit</span></div><div style="display:flex;gap:10px;flex-wrap:wrap">'+activeRecordings.map(function(r){return '<div class="tag tag-red" style="display:flex;align-items:center;gap:10px;padding:8px 12px"><span><strong>'+escHtml(String(r.StreamKey||'-'))+'</strong> · '+escHtml(String(r.Format||'').toUpperCase())+'</span>'+(r.Status==='recording'?'<button class="btn btn-sm btn-danger" onclick="stopRec(\''+r.ID+'\')">Durdur</button>':'')+'</div>';}).join('')+'</div></div>';
}
function renderStorageRemuxJobs(data){
  const remuxJobs=Array.isArray(data&&data.remuxJobs)?data.remuxJobs:[];
  if(!remuxJobs.length)return '';
  return '<div class="card" style="margin-bottom:16px;background:linear-gradient(180deg,#f8fbff 0%,#f2f8ff 100%)"><div class="card-header"><div><div class="card-title">Donusum ve Senkron Isleri</div><div class="form-hint">MP4 hazirlama ve benzeri uzun isler arka planda devam eder.</div></div><button class="btn btn-secondary btn-sm" onclick="refreshStorageSnapshot({resetPreview:false})">Yenile</button></div><div style="display:flex;gap:10px;flex-wrap:wrap">'+remuxJobs.slice(0,8).map(function(job){var tone=job.status==='completed'?'green':(job.status==='error'?'red':'yellow'); var label=job.status==='completed'?'Hazir':(job.status==='error'?'Hata':'Calisiyor'); return '<div class="tag tag-'+tone+'" style="display:flex;align-items:center;gap:8px;padding:8px 12px"><span><strong>'+escHtml(job.source_name||'-')+'</strong> &rarr; '+escHtml((job.target_format||'mp4').toUpperCase())+'</span><span>'+label+'</span></div>';}).join('')+'</div></div>';
}
function renderStorageStatsGrid(data){
  const report=data&&data.report?data.report:{};
  return statCard('blue','bi-hdd-fill',formatBytes((report&&report.storage&&report.storage.recordings_bytes)||0),'Yerel Kayitlar')+
    statCard('purple','bi-archive-fill',fmtInt((data&&data.backups||[]).length),'Yerel Yedekler')+
    statCard('orange','bi-cloud-arrow-up-fill',fmtInt((data&&data.archives||[]).length),'Kayit Arsivi')+
    statCard('green','bi-safe2-fill',fmtInt((data&&data.backupArchives||[]).length),'Yedek Arsivi');
}
function renderStorageActiveRecordingRows(data){
  const recs=Array.isArray(data&&data.recs)?data.recs:[];
  if(!recs.length)return '<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Aktif kayit yok</td></tr>';
  return recs.map(function(r){
    const recID=String(r.ID||'');
    const streamKey=String(r.StreamKey||'');
    const shortID=recID.length>28?recID.slice(0,28)+'...':recID;
    const shortStream=streamKey.length>22?streamKey.slice(0,22)+'...':streamKey;
    return '<tr>'+
      '<td><code title="'+escHtml(recID)+'" style="display:inline-block;max-width:260px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;vertical-align:bottom">'+escHtml(shortID)+'</code></td>'+
      '<td><code title="'+escHtml(streamKey)+'" style="display:inline-block;max-width:220px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;vertical-align:bottom">'+escHtml(shortStream)+'</code></td>'+
      '<td style="white-space:nowrap">'+escHtml(String(r.Format||'').toUpperCase())+'</td>'+
      '<td style="white-space:nowrap"><span class="badge badge-'+(r.Status==='recording'?'green':(r.Status==='error'?'red':'gray'))+'">'+escHtml(String(r.Status||'-'))+'</span></td>'+
      '<td style="white-space:nowrap">'+fmtBytes(r.Size||0)+'</td>'+
      '<td style="white-space:nowrap">'+(r.Status==='recording'?'<button class="btn btn-sm btn-danger" onclick="stopRec(\''+r.ID+'\')">Durdur</button>':'-')+'</td>'+
    '</tr>';
  }).join('');
}
function renderStorageSavedRecordingRows(data){
  const saved=Array.isArray(data&&data.saved)?data.saved:[];
  const archiveMap=data&&data.archiveMap?data.archiveMap:{};
  const archiveEnabled=!!(data&&data.archiveEnabled);
  if(!saved.length)return '<tr><td colspan="7" style="text-align:center;color:var(--text-muted);padding:24px">Kaydedilmis dosya yok</td></tr>';
  return saved.map(function(r){
    const archiveInfo=archiveMap[r.stream_key+'::'+r.name];
    const archiveBadge=archiveInfo?renderArchiveStatusBadge(archiveInfo):'<span class="tag tag-blue">Yerelde</span>';
    const format=String(r.format||'').toLowerCase();
    const canRemux=format==='ts'||format==='flv'||format==='mkv';
    return '<tr>'+
      '<td><code>'+escHtml(r.stream_key)+'</code></td>'+
      '<td>'+escHtml(r.name)+'</td>'+
      '<td>'+(r.format||'-').toUpperCase()+'</td>'+
      '<td>'+fmtLocaleDateTime(r.mod_time)+'</td>'+
      '<td>'+fmtBytes(r.size||0)+'</td>'+
      '<td>'+archiveBadge+'</td>'+
      '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
        '<button class="btn btn-sm btn-secondary" onclick=\'previewRecordingPanel('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+','+JSON.stringify(r.format||'')+','+JSON.stringify(r.mod_time||'')+','+(r.size||0)+')\'>Onizle</button>'+
        '<button class="btn btn-sm btn-secondary" onclick=\'downloadRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Indir</button>'+
        (canRemux?'<button class="btn btn-sm btn-secondary" onclick=\'remuxRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+','+JSON.stringify('mp4')+')\'>MP4 Hazirla</button>':'')+
        (archiveEnabled?'<button class="btn btn-sm btn-secondary" onclick=\'archiveRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>'+(archiveInfo&&archiveInfo.status==='archived'?'Yeniden Arsivle':'Arsive Gonder')+'</button>':'')+
        (archiveInfo&&archiveInfo.object_url?'<button class="btn btn-sm btn-secondary" onclick=\'window.open('+JSON.stringify(archiveInfo.object_url)+',"_blank")\'>Arsiv Linki</button>':'')+
        '<button class="btn btn-sm btn-danger" onclick=\'deleteRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Sil</button>'+
      '</td>'+
    '</tr>';
  }).join('');
}
function renderStorageArchiveRows(data){
  const archives=Array.isArray(data&&data.archives)?data.archives:[];
  if(!archives.length)return '<tr><td colspan="7" style="text-align:center;color:var(--text-muted);padding:24px">Arsivlenmis kayit yok</td></tr>';
  return archives.map(function(item){
    const localState=item.local_deleted?'<span class="tag tag-yellow">Yerelde yok</span>':'<span class="tag tag-green">Yerelde var</span>';
    const statusBadge=renderArchiveStatusBadge(item);
    return '<tr>'+
      '<td><code>'+escHtml(item.stream_key)+'</code></td>'+
      '<td>'+escHtml(item.filename)+'</td>'+
      '<td>'+escHtml(String(item.provider||'-').toUpperCase())+'</td>'+
      '<td>'+fmtLocaleDateTime(item.archived_at||item.updated_at||item.created_at)+'</td>'+
      '<td>'+localState+'</td>'+
      '<td>'+statusBadge+(item.last_error?'<div class="setting-desc" style="max-width:320px">'+escHtml(item.last_error)+'</div>':'')+'</td>'+
      '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
        '<button class="btn btn-sm btn-secondary" onclick=\'restoreRecordingArchive('+JSON.stringify(item.stream_key)+','+JSON.stringify(item.filename)+')\'>Geri Yukle</button>'+
        (item.object_url?'<button class="btn btn-sm btn-secondary" onclick=\'window.open('+JSON.stringify(item.object_url)+',"_blank")\'>Arsiv Linki</button>':'')+
      '</td>'+
    '</tr>';
  }).join('');
}
function renderStorageBackupRows(data){
  const backups=Array.isArray(data&&data.backups)?data.backups:[];
  const backupArchiveMap=data&&data.backupArchiveMap?data.backupArchiveMap:{};
  const backupArchiveEnabled=!!(data&&data.backupArchiveEnabled);
  if(!backups.length)return '<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Yerel sistem yedegi yok</td></tr>';
  return backups.map(function(item){
    const archiveInfo=backupArchiveMap[item.name];
    const archiveBadge=archiveInfo?renderBackupArchiveStatusBadge(archiveInfo):'<span class="tag tag-blue">Yerelde</span>';
    return '<tr data-backup-name="'+escHtml(item.name)+'">'+
      '<td class="mono-wrap">'+escHtml(item.name)+'</td>'+
      '<td>'+formatBytes(item.size||0)+'</td>'+
      '<td>'+escHtml(fmtLocaleDateTime(item.mod_time))+'</td>'+
      '<td>'+(item.include_recordings?'<span class="tag tag-blue">Kayitlar dahil</span>':'<span class="tag tag-green">Hafif</span>')+'</td>'+
      '<td>'+archiveBadge+'</td>'+
      '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
        '<button class="btn btn-sm btn-secondary" onclick=\'downloadSystemBackup('+JSON.stringify(item.name)+')\'>Indir</button>'+
        (backupArchiveEnabled?'<button class="btn btn-sm btn-secondary" onclick=\'archiveSystemBackup('+JSON.stringify(item.name)+')\'>'+(archiveInfo&&archiveInfo.status==='archived'?'Yeniden Arsivle':'Arsive Gonder')+'</button>':'')+
        '<button class="btn btn-sm btn-danger" onclick=\'deleteSystemBackup('+JSON.stringify(item.name)+')\'>Sil</button>'+
      '</td>'+
    '</tr>';
  }).join('');
}
function renderStorageBackupArchiveRows(data){
  const backupArchives=Array.isArray(data&&data.backupArchives)?data.backupArchives:[];
  if(!backupArchives.length)return '<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Arsivlenmis sistem yedegi yok</td></tr>';
  return backupArchives.map(function(item){
    const localState=item.local_deleted?'<span class="tag tag-yellow">Yerelde yok</span>':'<span class="tag tag-green">Yerelde var</span>';
    const statusBadge=renderBackupArchiveStatusBadge(item);
    return '<tr>'+
      '<td class="mono-wrap">'+escHtml(item.name)+'</td>'+
      '<td>'+escHtml(String(item.provider||'-').toUpperCase())+'</td>'+
      '<td>'+fmtLocaleDateTime(item.archived_at||item.updated_at||item.created_at)+'</td>'+
      '<td>'+localState+'</td>'+
      '<td>'+statusBadge+(item.last_error?'<div class="setting-desc" style="max-width:320px">'+escHtml(item.last_error)+'</div>':'')+'</td>'+
      '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
        '<button class="btn btn-sm btn-secondary" onclick=\'restoreSystemBackupArchive('+JSON.stringify(item.name)+')\'>Geri Getir</button>'+
        (item.object_url?'<button class="btn btn-sm btn-secondary" onclick=\'window.open('+JSON.stringify(item.object_url)+',"_blank")\'>Arsiv Linki</button>':'')+
      '</td>'+
    '</tr>';
  }).join('');
}
function applyStorageSnapshot(data,opts){
  const options=opts||{};
  window._storageData=data;
  window._recStreams=data.streams;
  window._savedRecordings=data.saved;
  window._recordingArchives=data.archives;
  window._systemBackups=data.backups;
  window._backupArchives=data.backupArchives;
  const activeBanner=document.getElementById('storage-active-banner');
  if(activeBanner)activeBanner.innerHTML=renderStorageActiveBanner(data);
  const jobs=document.getElementById('storage-remux-jobs');
  if(jobs)jobs.innerHTML=renderStorageRemuxJobs(data);
  const stats=document.getElementById('storage-stats-grid');
  if(stats)stats.innerHTML=renderStorageStatsGrid(data);
  const recCount=document.getElementById('storage-active-count');
  if(recCount)recCount.textContent=fmtInt((data.recs||[]).length)+' aktif oturum';
  const recList=document.getElementById('rec-list');
  if(recList)recList.innerHTML=renderStorageActiveRecordingRows(data);
  const savedList=document.getElementById('saved-rec-list');
  if(savedList)savedList.innerHTML=renderStorageSavedRecordingRows(data);
  const archiveList=document.getElementById('archive-rec-list');
  if(archiveList)archiveList.innerHTML=renderStorageArchiveRows(data);
  const backupList=document.getElementById('system-backup-list');
  if(backupList)backupList.innerHTML=renderStorageBackupRows(data);
  const backupArchiveList=document.getElementById('backup-archive-list');
  if(backupArchiveList)backupArchiveList.innerHTML=renderStorageBackupArchiveRows(data);
  const selection=window._recordingPreviewSelection;
  const selectedStillExists=!!(selection&&Array.isArray(data.saved)&&data.saved.some(function(item){
    return String(item.stream_key||'')===String(selection.stream_key||'') && String(item.name||'')===String(selection.name||'');
  }));
  if(options.resetPreview || !selectedStillExists){
    window._recordingPreviewSelection=null;
    teardownRecordingPreview();
    resetRecordingPreviewPanel();
  }
  updateStorageProviderUI();
}
async function fetchStorageSnapshot(){
  const [s,report,archivesRes,recsRes,streamsRes,savedRes,backupsRes,backupArchivesRes,remuxJobsRes]=await Promise.all([
    api('/api/settings'),
    api('/api/health/report'),
    api('/api/recordings/archives'),
    api('/api/recordings'),
    api('/api/streams'),
    api('/api/recordings/library'),
    api('/api/system/backups'),
    api('/api/system/backups/archives'),
    api('/api/recordings/remux/jobs')
  ]);
  return normalizeStorageSnapshot(s,report,archivesRes,recsRes,streamsRes,savedRes,backupsRes,backupArchivesRes,remuxJobsRes);
}
async function refreshStorageSnapshot(opts){
  const data=await fetchStorageSnapshot();
  applyStorageSnapshot(data,opts||{});
  return data;
}
function setStorageFieldVisible(key, visible){
  const input=document.querySelector('.setting-input[data-key="'+key+'"]');
  if(!input)return;
  const row=input.closest('.form-group, .setting-row');
  if(row)row.style.display=visible?'':'none';
}
function updateStorageProviderUI(){
  const shell=document.getElementById('storage-settings-shell');
  if(!shell)return;
  const draft=readStorageDraftSettings();
  shell.innerHTML=renderStorageSettingsShell(draft,window._storageData||{});
}
let recordingPreviewPlayer=null;
window._recordingPreviewSelection=null;
function recordingFileURL(streamKey,name,download){
  return '/recordings/'+encodeURIComponent(streamKey)+'/'+encodeURIComponent(name)+(download?'?download=1':'');
}
function destroyRecordingPreviewPlayer(){
  try{
    if(recordingPreviewPlayer){
      recordingPreviewPlayer.destroy();
      recordingPreviewPlayer=null;
    }
  }catch(e){}
}
function teardownRecordingPreview(){
  destroyRecordingPreviewPlayer();
  try{
    const panel=document.getElementById('recording-preview-panel');
    if(!panel)return;
    panel.querySelectorAll('video,audio').forEach(function(media){
      try{
        media.pause();
        media.removeAttribute('src');
        media.load();
      }catch(e){}
    });
  }catch(e){}
}
function resetRecordingPreviewPanel(){
  const panel=document.getElementById('recording-preview-panel');
  if(!panel)return;
  panel.innerHTML='<div class="empty-state"><div class="icon"><i class="bi bi-film"></i></div><h3>Kayit secin</h3><p style="color:var(--text-muted)">Panel secili kaydi ayni sayfada oynatir.</p></div>';
}
async function refreshStorageSurface(targetPage){
  teardownRecordingPreview();
  await new Promise(function(resolve){setTimeout(resolve,0);});
  const page=targetPage||currentPage;
  if(page==='settings-storage'||page==='recordings'||page==='maintenance-center'){
    await loadPage(page);
    return;
  }
  navigate('settings-storage');
}
async function prepareStorageAction(){
  window._recordingPreviewSelection=null;
  teardownRecordingPreview();
  resetRecordingPreviewPanel();
  await new Promise(function(resolve){
    if(typeof requestAnimationFrame==='function'){
      requestAnimationFrame(function(){setTimeout(resolve,0);});
      return;
    }
    setTimeout(resolve,0);
  });
}
async function stopRec(id){
  const res=await api('/api/recordings/stop/'+id);
  if(res&&res.error){
    toast(res.message||'Kayit durdurulamadi','error');
    return;
  }
  toast('Kayit durduruldu');
  if(currentPage==='recordings'||currentPage==='settings-storage'){
    await refreshStorageSnapshot({resetPreview:false});
  }else if(currentPage==='maintenance-center'){
    await loadPage(currentPage);
  }else if(String(currentPage||'').indexOf('stream-detail-')===0){
    await loadPage(currentPage);
  }else{
    navigate('settings-storage');
  }
}
function showRecordModal(){
  const streams=(window._recStreams||[]).filter(s=>s.status==='live');
  const modal=document.getElementById('rec-modal');
  if(!modal)return;
  modal.style.display='block';
  modal.innerHTML='<div style="position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,.5);z-index:999;display:flex;align-items:center;justify-content:center" onclick="if(event.target===this)this.parentElement.style.display=\'none\'">'+
    '<div class="card" style="width:400px"><div class="card-header"><h3 class="card-title">Kayit Baslat</h3></div><div class="card-body">'+
    '<div class="form-group"><label class="form-label">Yayin</label><select class="form-select" id="rec-key">'+
    (streams.length?streams.map(s=>'<option value="'+s.stream_key+'">'+s.name+'</option>').join(''):'<option>Canli yayin yok</option>')+
    '</select></div>'+ 
    '<div class="form-group"><label class="form-label">Format</label><select class="form-select" id="rec-fmt">'+recordingFormatOptions('mp4')+'</select><div class="form-hint">MP4 tarayici ve panel onizlemesi icin onerilir.</div></div>'+ 
    '<button class="btn btn-primary" onclick="startNewRec()" style="width:100%">Kaydi Baslat</button>'+ 
    '</div></div></div>';
}
async function startNewRec(){
  const key=document.getElementById('rec-key')?.value;
  const fmt=document.getElementById('rec-fmt')?.value||'mp4';
  if(!key)return;
  await api('/api/recordings',{method:'POST',body:{stream_key:key,format:fmt}});
  document.getElementById('rec-modal').style.display='none';
  navigate('recordings');
}
// Yeni: Kayıt önizlemesini panelde gösteren fonksiyon
async function previewRecordingPanel(streamKey,name,format,mod_time,size){
  destroyRecordingPreviewPlayer();
  window._recordingPreviewSelection={stream_key:String(streamKey||''),name:String(name||'')};
  const panel=document.getElementById('recording-preview-panel');
  if(!panel)return;
  const url=recordingFileURL(streamKey,name,false);
  const ext=(name.split('.').pop()||'').toLowerCase();
  const canRemux=ext==='ts'||ext==='flv'||ext==='mkv';
  const header='<div style="margin-bottom:10px"><strong>'+escHtml(name)+'</strong> ('+fmtBytes(size||0)+')</div>';
  const actions='<div style="margin-top:12px;display:flex;gap:12px;flex-wrap:wrap">'+
    '<a class="btn btn-sm btn-secondary" href="'+url+'" target="_blank">Direkt Link</a>'+
    '<a class="btn btn-sm btn-secondary" href="'+recordingFileURL(streamKey,name,true)+'" target="_blank">Indir</a>'+
    (canRemux?'<button class="btn btn-sm btn-secondary" onclick=\'remuxRecordingFile('+JSON.stringify(streamKey)+','+JSON.stringify(name)+','+JSON.stringify('mp4')+')\'>MP4 Hazirla</button>':'')+
  '</div>';
  if(ext==='mp4'||ext==='webm'||ext==='ogg'){
    panel.innerHTML=header+'<div style="position:relative;width:100%;aspect-ratio:16/9;min-height:280px;background:#000;border-radius:14px;overflow:hidden"><video controls playsinline src="'+url+'" style="position:absolute;inset:0;width:100%;height:100%;background:#000;object-fit:contain"></video></div>'+actions;
    return;
  }
  if(ext==='mp3'||ext==='aac'||ext==='wav'||ext==='flac'){
    panel.innerHTML=header+'<div style="padding:24px"><audio controls src="'+url+'" style="width:100%"></audio></div>'+actions;
    return;
  }
  if(ext==='flv'||ext==='ts'||ext==='mkv'){
    panel.innerHTML=header+'<div class="empty-state"><div class="icon"><i class="bi bi-magic"></i></div><h3>Guvenli onizleme icin MP4 onerilir</h3><p style="color:var(--text-muted)">Bu kayit '+escHtml(ext.toUpperCase())+' olarak saklandi. Tarayici ici onizleme yerine once <strong>MP4 Hazirla</strong> kullanmaniz daha kararlidir.</p></div>'+actions;
    return;
  }
  panel.innerHTML=header+'<div class="empty-state"><h3>Onizleme yok</h3><p style="color:var(--text-muted)">Bu format panelde dogrudan oynatilamiyor.</p></div>'+actions;
}
function downloadRecordingFile(streamKey,name){
  window.open(recordingFileURL(streamKey,name,true),'_blank');
}
function downloadSystemBackup(name){
  const link=document.createElement('a');
  link.href='/api/system/backups/download/'+encodeURIComponent(name);
  link.target='_blank';
  link.rel='noopener';
  link.download=name||'backup.tar.gz';
  document.body.appendChild(link);
  link.click();
  setTimeout(function(){ if(link.parentNode)link.parentNode.removeChild(link); },0);
}
function removeStorageBackupRow(name){
  const tbody=document.getElementById('system-backup-list');
  if(!tbody)return;
  tbody.querySelectorAll('tr[data-backup-name]').forEach(function(row){
    if(String(row.getAttribute('data-backup-name')||'')===String(name||'')){
      row.remove();
    }
  });
  if(!tbody.children.length){
    tbody.innerHTML='<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Yerel sistem yedegi yok</td></tr>';
  }
}
async function remuxRecordingFile(streamKey,name,format){
  const res=await api('/api/recordings/remux',{method:'POST',body:{stream_key:streamKey,filename:name,format:format||'mp4'}});
  if(res&&res.success&&res.job){
    toast('MP4 donusumu arka planda basladi');
    if(currentPage==='recordings'||currentPage==='settings-storage')await refreshStorageSnapshot({resetPreview:false});
    else if(currentPage==='maintenance-center')await loadPage(currentPage);
  }else{
    toast((res&&res.message)||'Donusum basarisiz','error');
  }
}
async function archiveRecordingFile(streamKey,name){
  const res=await api('/api/recordings/archive',{method:'POST',body:{stream_key:streamKey,filename:name}});
  if(res&&res.stream_key){
    toast('Kayit arsive gonderildi');
    if(currentPage==='recordings'||currentPage==='settings-storage')await refreshStorageSnapshot({resetPreview:false});
    else if(currentPage==='maintenance-center')await loadPage(currentPage);
  }else{
    toast((res&&res.message)||'Arsivleme basarisiz','error');
  }
}
async function restoreRecordingArchive(streamKey,name){
  const res=await api('/api/recordings/restore',{method:'POST',body:{stream_key:streamKey,filename:name}});
  if(res&&res.stream_key){
    toast('Kayit geri yuklendi');
    if(currentPage==='recordings'||currentPage==='settings-storage')await refreshStorageSnapshot({resetPreview:false});
    else if(currentPage==='maintenance-center')await loadPage(currentPage);
  }else{
    toast((res&&res.message)||'Geri yukleme basarisiz','error');
  }
}
async function deleteRecordingFile(streamKey,name){
  if(!confirm('Bu kayit dosyasini silmek istediginize emin misiniz?'))return;
  const deletingSelected=!!(window._recordingPreviewSelection&&String(window._recordingPreviewSelection.stream_key||'')===String(streamKey||'')&&String(window._recordingPreviewSelection.name||'')===String(name||''));
  const res=await api('/api/recordings/file',{method:'DELETE',body:{stream_key:streamKey,filename:name}});
  if(res&&res.status==='deleted'){
    if(deletingSelected){
      window._recordingPreviewSelection=null;
      teardownRecordingPreview();
      resetRecordingPreviewPanel();
    }
    toast('Kayit silindi');
    if(currentPage==='recordings'||currentPage==='settings-storage')await refreshStorageSnapshot({resetPreview:false});
    else if(currentPage==='maintenance-center')await loadPage(currentPage);
    else navigate('settings-storage');
  }else{
    toast((res&&res.message)||'Kayit silinemedi','error');
  }
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â VIEWERS Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderViewers(c){
  const data=await api('/api/viewers');
  const sessions=Array.isArray(data&&data.sessions)?data.sessions:[];
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Izleyiciler</h1></div>'+
    '<div class="card-grid card-grid-3" style="margin-bottom:24px">'+
      statCard('green','bi-people',fmtInt(data?data.total:0),'Toplam Izleyici','analytics','Toplam oturum sayisi')+
      statCard('blue','bi-eye-fill',fmtInt(data?data.active:0),'Aktif Izleyici','dashboard','Su an acik baglantilar')+
      statCard('orange','bi-shield-fill-x',fmtInt(data?data.banned:0),'Yasakli IP','settings-security','Aktif IP ban kayitlari')+
    '</div>'+
    '<div class="card" style="margin-bottom:24px"><div class="card-header"><h3 class="card-title">Aktif Oturumlar</h3></div>'+
    '<div class="card-body"><table class="table viewer-table"><thead><tr><th>Yayin</th><th>Format</th><th>IP</th><th>Ulke</th><th>Sure</th><th>Trafik</th><th>Son Gorulme</th></tr></thead><tbody id="viewer-session-list"></tbody></table></div></div>'+
    '<div class="card"><div class="card-header"><h3 class="card-title">IP Yasaklama</h3></div>'+
    '<div class="card-body"><div style="display:flex;gap:8px;margin-bottom:16px;flex-wrap:wrap"><input class="input" id="ban-ip" placeholder="IP adresi" style="flex:1;min-width:140px"><input class="input" id="ban-reason" placeholder="Neden" style="flex:1;min-width:140px"><input class="input" id="ban-dur" type="number" placeholder="Sure (dk)" style="width:120px"><button class="btn btn-primary" onclick="banIP()">Yasakla</button></div>'+
    '<table class="table"><thead><tr><th>IP</th><th>Neden</th><th>Tarih</th><th>Islem</th></tr></thead><tbody id="ban-list"></tbody></table></div></div>';
  const sl=document.getElementById('viewer-session-list');
  if(sl){
    sl.innerHTML=sessions.length?sessions.map(function(sess){
      return '<tr>'+
        '<td><div style="font-weight:600">'+escHtml(sess.stream_name||shortKey(sess.stream_key))+'</div><div class="setting-desc"><code>'+escHtml(sess.stream_key)+'</code></div></td>'+
        '<td><span class="badge">'+escHtml((sess.format||'-').toUpperCase())+'</span></td>'+
        '<td><code>'+escHtml(sess.ip||'-')+'</code></td>'+
        '<td>'+escHtml(sess.country||'-')+'</td>'+
        '<td>'+formatDurationSeconds(sess.duration_seconds||0)+'</td>'+
        '<td>'+fmtBytes(sess.bytes_sent||0)+'</td>'+
        '<td>'+fmtLocaleTime(sess.last_seen)+'</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="7" style="text-align:center;color:var(--text-muted);padding:24px">Aktif izleyici oturumu yok</td></tr>';
  }
  const bans=await api('/api/security/bans');
  const bl=document.getElementById('ban-list');
  if(bl&&bans){
    bl.innerHTML=bans.length?bans.map(b=>'<tr><td><code>'+b.IP+'</code></td><td>'+b.Reason+'</td><td>'+fmtLocaleDateTime(b.BannedAt).replace(/^-$/,'\u2014')+'</td><td><button onclick="unbanIP(\''+b.IP+'\')" style="background:#e74c3c;color:#fff;padding:4px 12px;border:none;border-radius:6px;cursor:pointer;font-size:12px">Kaldir</button></td></tr>').join(''):'<tr><td colspan="4" style="text-align:center;color:var(--text-muted);padding:24px">Yasakli IP yok</td></tr>';
  }
  schedulePageRefresh('viewers',5000);
}
function formatDurationSeconds(total){
  total=Math.max(0,parseInt(total||0,10));
  const h=Math.floor(total/3600);
  const m=Math.floor((total%3600)/60);
  const s=total%60;
  if(h>0)return h+'s '+m+'d';
  if(m>0)return m+'d '+s+'sn';
  return s+'sn';
}
async function banIP(){
  const ip=document.getElementById('ban-ip')?.value;
  const reason=document.getElementById('ban-reason')?.value||'Manuel';
  const dur=parseInt(document.getElementById('ban-dur')?.value)||0;
  if(!ip)return;
  await api('/api/security/bans',{method:'POST',body:{ip:ip,reason:reason,duration_minutes:dur}});
  navigate('viewers');
}
async function unbanIP(ip){
  await api('/api/security/bans',{method:'DELETE',body:{ip:ip}});
  navigate('viewers');
}

async function renderSecurityTokens(c){
  const streamsRes=await api('/api/streams');
  const settings=await api('/api/settings');
  const streams=Array.isArray(streamsRes)?streamsRes:[];
  const tokenEnabled=(settings&&settings.token_enabled==='true')?'Acik':'Kapali';
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Tokenlar</h1></div>'+
    '<div class="card-grid card-grid-3" style="margin-bottom:24px">'+
      statCard('blue','bi-key-fill',streams.length,'Toplam Yayin')+
      statCard('green','bi-broadcast-pin',(streams.filter(function(s){return s.status==='live'}).length||0),'Canli Yayin')+
      statCard('orange','bi-shield-lock-fill',tokenEnabled,'Token Durumu')+
    '</div>'+
    '<div class="card" style="margin-bottom:24px">'+
      '<div class="card-header"><h3 class="card-title">Izleme Tokeni Uret</h3></div>'+
      '<div class="card-body">'+
        '<div class="card-grid card-grid-2">'+
          '<div class="form-group"><label class="form-label">Yayin</label><select class="form-select" id="token-stream">'+
            '<option value="">-- Yayin Secin --</option>'+
            streams.map(function(s){return '<option value="'+s.stream_key+'">'+escHtml(s.name)+' ('+escHtml(s.stream_key)+')</option>'}).join('')+
          '</select></div>'+
          '<div class="form-group"><label class="form-label">Not</label><div class="form-hint" style="padding-top:12px">Token uretimi hazir. Zorunluluk ayari Guvenlik ekranindaki <code>Token Dogrulama</code> alanindan yonetilir.</div></div>'+
        '</div>'+
        '<button class="btn btn-primary" onclick="generateStreamToken()">Token Uret</button>'+
        '<div id="token-output" style="margin-top:16px"></div>'+
      '</div>'+
    '</div>'+
    '<div class="card"><div class="card-header"><h3 class="card-title">Kullanim Notu</h3></div>'+
      '<div class="card-body"><div class="form-hint">Uretilen degeri uygulama katmaninda <code>token</code> parametresi veya yetkilendirme basligi olarak tasiyabilirsiniz. Panel bu ekranda yalnizca token uretir; hangi client tarafinda nasil eklenecegi entegrasyona gore belirlenir.</div></div>'+
    '</div>';
}

async function renderLogos(c){
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Logo ve Marka</h1></div>'+
    '<div class="card"><div class="empty-state"><div class="icon"><i class="bi bi-images"></i></div><h3>Marka studyosu yukleniyor</h3><p style="color:var(--text-muted)">Bu sayfa yeni studio katmani ile yuklenir. Sayfayi yenileyip tekrar deneyin.</p></div></div>';
}
async function generateStreamToken(){
  const key=document.getElementById('token-stream')?.value||'';
  const out=document.getElementById('token-output');
  if(!key){
    toast('Once bir yayin secin','error');
    return;
  }
  if(out)out.innerHTML='<div class="form-hint">Token uretiliyor...</div>';
  const res=await api('/api/security/token/generate',{method:'POST',body:{stream_key:key}});
  if(!res||!res.token){
    if(out)out.innerHTML='<div class="form-hint" style="color:#ef4444">Token uretilemedi</div>';
    toast((res&&res.message)||'Token uretilemedi','error');
    return;
  }
  if(out){
    out.innerHTML=
      copyField('Token',res.token)+
      copyField('Gecerlilik',fmtLocaleDateTime(res.expires_at));
  }
  toast('Token olusturuldu');
}

async function renderTranscodeJobs(c){
  const statusRes=await api('/api/transcode/status');
  const jobsRes=await api('/api/transcode/jobs');
  const status=(statusRes&&typeof statusRes==='object')?statusRes:{};
  const jobs=Array.isArray(jobsRes)?jobsRes:[];
  const liveHLS=jobs.filter(function(job){return job.type==='live_hls'&&job.status==='running'}).length;
  const liveDASH=jobs.filter(function(job){return job.type==='live_dash'&&job.status==='running'}).length;
  const ffmpegVersion=String(status.ffmpeg_version||'bilinmiyor');
  const liveOptions=status.live_options||{};
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Transcode Isleri</h1></div>'+
    '<div class="card-grid card-grid-4" style="margin-bottom:24px">'+
      statCard('blue','bi-cpu-fill',fmtInt(status.active_jobs||0),'Aktif Is')+
      statCard('purple','bi-badge-hd-fill',fmtInt(liveHLS),'Canli HLS Isleri')+
      statCard('green','bi-badge-4k-fill',fmtInt(liveDASH),'Canli DASH Isleri')+
      statCard('orange','bi-gpu-card',escHtml((status.gpu_accel||'none').toUpperCase()),'GPU Hizlandirma')+
    '</div>'+
    '<div class="insight-grid" style="margin-bottom:24px">'+
      '<div class="card"><div class="card-header"><h3 class="card-title">FFmpeg</h3></div>'+
        '<div class="card-body"><div class="metric-list">'+
          '<div class="metric-row"><span>Surum</span><strong>'+escHtml(ffmpegVersion.split(' ').slice(0,2).join(' ')||ffmpegVersion)+'</strong></div>'+
          '<div class="metric-row"><span>Calisma yolu</span><span class="mono-wrap">'+escHtml(status.ffmpeg_path||'-')+'</span></div>'+
          '<div class="metric-row"><span>Toplam is</span><strong>'+fmtInt(status.total_jobs||0)+'</strong></div>'+
        '</div></div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Calisma Ortami</h3></div>'+
        '<div class="card-body"><div class="metric-list">'+
          '<div class="metric-row"><span>Isletim Sistemi</span><strong>'+escHtml(status.os||'-')+'</strong></div>'+
          '<div class="metric-row"><span>Mimari</span><strong>'+escHtml(status.arch||'-')+'</strong></div>'+
          '<div class="metric-row"><span>GPU modu</span><strong>'+escHtml((status.gpu_accel||'none').toUpperCase())+'</strong></div>'+
          '<div class="metric-row"><span>ABR</span><strong>'+(liveOptions.abr_enabled?'Acik':'Kapali')+'</strong></div>'+
          '<div class="metric-row"><span>Profil Seti</span><strong>'+escHtml(liveOptions.profile_set||'balanced')+'</strong></div>'+
        '</div></div></div>'+
    '</div>'+
    '<div class="card"><div class="card-header"><h3 class="card-title">Job Listesi</h3></div>'+
      '<div class="card-body"><table class="table"><thead><tr><th>ID</th><th>Yayin</th><th>Tip</th><th>Durum</th><th>Baslangic</th><th>PID</th><th>Cikti</th><th>Hata</th></tr></thead><tbody id="tc-job-list"></tbody></table></div></div>';
  const list=document.getElementById('tc-job-list');
  if(list){
    list.innerHTML=jobs.length?jobs.map(function(job){
      const started=fmtLocaleDateTime(job.started_at);
      const type=job.type||'abr';
      const statusClass=job.status==='running'?'green':(job.status==='error'?'red':'gray');
      return '<tr>'+
        '<td style="font-size:12px"><code>'+escHtml(shortKey(job.id||'-'))+'</code></td>'+
        '<td><code>'+escHtml(shortKey(job.stream_key||'-'))+'</code></td>'+
        '<td>'+escHtml(type)+'</td>'+
        '<td><span class="badge badge-'+statusClass+'">'+escHtml(job.status||'-')+'</span></td>'+
        '<td>'+started+'</td>'+
        '<td>'+(job.pid||'-')+'</td>'+
        '<td class="mono-wrap">'+escHtml(job.output_dir||'-')+'</td>'+
        '<td style="max-width:260px;color:'+(job.error?'#ef4444':'var(--text-muted)')+'">'+escHtml(job.error||'-')+'</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="8" style="text-align:center;color:var(--text-muted);padding:24px">Henuz transcode isi yok</td></tr>';
  }
  schedulePageRefresh('transcode-jobs',5000);
}

async function renderMaintenanceCenter(c){
  const [serviceRes,backupsRes,healthRes,upgradeRes]=await Promise.all([
    api('/api/system/service/status'),
    api('/api/system/backups'),
    api('/api/health/report'),
    api('/api/system/upgrade/plan')
  ]);
  const status=(serviceRes&&serviceRes.status)||{};
  const backups=(backupsRes&&backupsRes.items)||[];
  const upgrade=(upgradeRes&&typeof upgradeRes==='object')?upgradeRes:{};
  const commands=(upgrade&&upgrade.commands)||{};
  const platform=String((serviceRes&&serviceRes.platform)||status.platform||'unknown');
  const unit=String((serviceRes&&serviceRes.unit)||status.unit||'-');
  const restoreCmd=commands.backup_restore||(platform==='linux'
    ? 'sudo systemctl stop '+unit+' && sudo /opt/fluxstream/fluxstream backup restore fluxstream-backup-YYYYMMDD-HHMMSS.tar.gz && sudo systemctl start '+unit
    : 'FluxStream.exe backup restore fluxstream-backup-YYYYMMDD-HHMMSS.tar.gz');
  const upgradeCmd=commands.atomic_upgrade||'-';
  const serviceButtons=platform==='linux'
    ? '<button class="btn btn-primary" onclick="serviceAction(\'restart\')">Servisi Yeniden Baslat</button>'+
      '<button class="btn btn-secondary" onclick="serviceAction(\'start\')">Servisi Baslat</button>'+
      '<button class="btn btn-danger" onclick="serviceAction(\'stop\')">Servisi Durdur</button>'
    : '<button class="btn btn-primary" onclick="restartServer()">Yeniden Baslat</button>'+
      '<button class="btn btn-danger" onclick="stopServer()">Durdur</button>';
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Bakim ve Yedek</h1><div style="display:flex;align-items:center;justify-content:space-between;gap:12px;flex-wrap:wrap"><div style="color:var(--text-muted);font-size:13px">Servis durumu, tek tikla yedek alma ve temiz geri donus komutlari burada toplanir.</div><button class="btn btn-secondary" onclick="loadPage(\'settings-storage\')">Depolama ve Arsiv Merkezine Git</button></div></div>'+
    '<div class="studio-alert info" style="margin-bottom:16px"><strong>Bu sayfa ne icin?</strong><div style="margin-top:8px" class="form-hint">Servis aksiyonlari, upgrade plani ve geri donus komutlari burada kalir. Kayit arsivi, bulut hedefleri ve senkron kuyruklari ise <strong>Depolama ve Arsiv Merkezi</strong> icinde yonetilir.</div></div>'+
    '<div class="card-grid card-grid-4" style="margin-bottom:16px">'+
      statCard(status.active?'green':'red','bi-hdd-network',status.active?'AKTIF':'DURDU','Servis Durumu')+
      statCard(status.enabled?'blue':'orange','bi-toggle-on',status.enabled?'ACIK':'KAPALI','Otomatik Baslangic')+
      statCard('purple','bi-archive-fill',fmtInt(backups.length),'Toplam Yedek')+
      statCard('orange','bi-heart-pulse-fill',String((healthRes&&healthRes.status)||'ok').toUpperCase(),'Saglik')+
    '</div>'+
    '<div class="card-grid card-grid-2">'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Servis Bilgisi</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Platform</span><strong>'+escHtml(platform)+'</strong></div>'+
          '<div class="metric-row"><span>Servis Yonetici</span><strong>'+escHtml(String(status.manager||'-'))+'</strong></div>'+
          '<div class="metric-row"><span>Unit</span><span class="mono-wrap">'+escHtml(unit)+'</span></div>'+
          '<div class="metric-row"><span>Main PID</span><strong>'+escHtml(String(status.main_pid||0))+'</strong></div>'+
          '<div class="metric-row"><span>Aktif Oldugu Zaman</span><strong>'+escHtml(String(status.since||'-'))+'</strong></div>'+
          '<div class="metric-row"><span>Kurulum Dizini</span><span class="mono-wrap">'+escHtml(String(upgrade.install_dir||'-'))+'</span></div>'+
        '</div>'+
        '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-top:14px">'+
          serviceButtons+
          '<button class="btn btn-secondary" onclick="loadPage(\'maintenance-center\')">Durumu Yenile</button>'+
        '</div>'+
        (status.message?'<div class="form-hint" style="margin-top:12px">'+escHtml(String(status.message))+'</div>':'')+
      '</div>'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Backup / Restore Akisi</div>'+
        '<div class="form-hint" style="line-height:1.8;margin-bottom:12px">Panelden yedek alabilirsiniz. Geri yukleme bilerek web panelinden yapilmiyor; calisan servis uzerinde restore riskli oldugu icin offline komutla ilerliyoruz.</div>'+
        '<div style="display:flex;gap:10px;flex-wrap:wrap;margin-bottom:12px">'+
          '<button class="btn btn-primary" onclick="createSystemBackup(false)">Yedek Al</button>'+
          '<button class="btn btn-secondary" onclick="createSystemBackup(true)">Kayitlarla Birlikte Al</button>'+
        '</div>'+
        '<div class="form-group" style="margin-bottom:14px"><label class="form-label">Geri Yukleme Komutu</label><textarea class="form-textarea" readonly style="min-height:96px">'+escHtml(restoreCmd)+'</textarea><div class="form-hint">Linux kurulumu icin servis once durdurulur, backup geri yuklenir, sonra servis tekrar baslatilir.</div></div>'+
        '<div class="form-group" style="margin-bottom:0"><label class="form-label">Atomic Upgrade Komutu</label><textarea class="form-textarea" readonly style="min-height:96px">'+escHtml(upgradeCmd)+'</textarea><div class="form-hint">Yeni binary once *.next olarak yuklenir, servis durdurulur, atomik rename yapilip servis yeniden baslatilir.</div></div>'+
      '</div>'+
    '</div>'+
    '<div class="card" style="margin-top:16px">'+
      '<div class="card-header"><div class="card-title">Hazir Yedekler</div><div class="form-hint">'+escHtml(backups.length?('Son yedek: '+fmtLocaleDateTime(backups[0].mod_time)):'Henuz backup yok')+'</div></div>'+
      (backups.length
        ?'<table><thead><tr><th>Dosya</th><th>Boyut</th><th>Tarih</th><th>Tur</th><th>Islem</th></tr></thead><tbody>'+
          backups.map(function(item){
            return '<tr><td class="mono-wrap">'+escHtml(item.name)+'</td><td>'+formatBytes(item.size||0)+'</td><td>'+escHtml(fmtLocaleDateTime(item.mod_time))+'</td><td>'+(item.include_recordings?'<span class="tag tag-blue">Kayitlar dahil</span>':'<span class="tag tag-green">Hafif</span>')+'</td><td style="white-space:nowrap"><a class="btn btn-sm btn-secondary" href="/api/system/backups/download/'+encodeURIComponent(item.name)+'" target="_blank" rel="noopener">Indir</a> <button class="btn btn-sm btn-danger" onclick=\'deleteSystemBackup('+JSON.stringify(item.name)+')\'>Sil</button></td></tr>';
          }).join('')+
        '</tbody></table>'
        :'<div class="empty-state"><div class="icon"><i class="bi bi-archive"></i></div><h3>Henuz backup yok</h3><p style="color:var(--text-muted)">Ilk yedegi bu ekrandan tek tikla alabilirsiniz.</p></div>')+
    '</div>';
}

async function createSystemBackup(includeRecordings){
  const res=await api('/api/system/backups',{method:'POST',body:{include_recordings:!!includeRecordings}});
  if(res&&res.success){
    toast('Backup hazirlandi');
    if(currentPage==='settings-storage'||currentPage==='recordings')await refreshStorageSnapshot({resetPreview:false});
    else loadPage('maintenance-center');
  }else{
    toast((res&&res.message)||'Backup olusturulamadi','error');
  }
}
async function createSystemBackupFromStorage(includeRecordings){
  await createSystemBackup(includeRecordings);
}

async function archiveSystemBackup(name){
  const res=await api('/api/system/backups/archive',{method:'POST',body:{name:name}});
  if(res&&res.success){
    toast('Sistem yedegi arsive gonderildi');
    if(currentPage==='settings-storage'||currentPage==='recordings')await refreshStorageSnapshot({resetPreview:false});
    else await loadPage('maintenance-center');
  }else{
    toast((res&&res.message)||'Yedek arsivleme basarisiz','error');
  }
}

async function restoreSystemBackupArchive(name){
  const res=await api('/api/system/backups/archive/restore',{method:'POST',body:{name:name}});
  if(res&&res.success){
    toast('Arsiv yedegi yerel backup klasorune geri getirildi');
    if(currentPage==='settings-storage'||currentPage==='recordings')await refreshStorageSnapshot({resetPreview:false});
    else await loadPage('maintenance-center');
  }else{
    toast((res&&res.message)||'Arsiv yedegi geri getirilemedi','error');
  }
}

async function deleteSystemBackup(name){
  if(!confirm('Bu backup silinsin mi?'))return;
  const res=await api('/api/system/backups/delete',{method:'POST',body:{name:name}});
  if(res&&res.success){
    toast('Backup silindi');
    removeStorageBackupRow(name);
    if(currentPage==='settings-storage'||currentPage==='recordings')await refreshStorageSnapshot({resetPreview:false});
    else loadPage('maintenance-center');
  }else{
    toast((res&&res.message)||'Backup silinemedi','error');
  }
}

async function serviceAction(action){
  const res=await api('/api/system/service/action',{method:'POST',body:{action:action}});
  if(res&&res.success){
    toast('Servis aksiyonu gonderildi');
    loadPage('maintenance-center');
  }else{
    toast((res&&res.message)||'Servis aksiyonu basarisiz','error');
  }
}

async function renderLicensePage(c){
  const [statusRes,sampleRes]=await Promise.all([api('/api/license/status'),api('/api/license/sample')]);
  const status=(statusRes&&statusRes.status)||{};
  const runtimeInfo=(statusRes&&statusRes.runtime)||{};
  const sample=(sampleRes&&sampleRes.sample)||{};
  const features=Array.isArray(runtimeInfo.enabled_features)&&runtimeInfo.enabled_features.length?runtimeInfo.enabled_features:(Array.isArray(status.features)?status.features:[]);
  const mode=String(runtimeInfo.mode||status.mode||'unlicensed');
  const warnings=Array.isArray(runtimeInfo.warnings)?runtimeInfo.warnings:[];
  const tone=status.valid?'tag-green':(runtimeInfo.development?'tag-blue':mode==='unlicensed'?'tag-yellow':'tag-red');
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Lisans</h1><div style="color:var(--text-muted);font-size:13px">Offline imzali lisans dosyasi burada saklanir. Internet baglantisi olmadan dogrulama yapilir.</div></div>'+
    '<div class="card-grid card-grid-4" style="margin-bottom:16px">'+
      statCard(status.valid?'green':(runtimeInfo.development?'blue':'orange'),'bi-patch-check-fill',escHtml(mode.toUpperCase()),'Durum')+
      statCard('blue','bi-building',escHtml(status.customer||'-'),'Musteri')+
      statCard('purple','bi-calendar-check',escHtml(status.valid_until||'-'),'Gecerlilik')+
      statCard('orange','bi-grid-1x2-fill',fmtInt(features.length),'Ozellik')+
    '</div>'+
    '<div class="card-grid card-grid-2">'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Mevcut Lisans Durumu <span class="tag '+tone+'" style="margin-left:8px">'+escHtml(String(status.message||'Bekleniyor'))+'</span></div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Runtime Modu</span><strong>'+escHtml(mode)+'</strong></div>'+
          '<div class="metric-row"><span>Feature enforcement</span><strong>'+(runtimeInfo.enforced===false?'Gelistirme':'Aktif')+'</strong></div>'+
          '<div class="metric-row"><span>Public key kaynagi</span><span class="mono-wrap">'+escHtml(String(status.public_key_source||'-'))+'</span></div>'+
          '<div class="metric-row"><span>Embedded key kullaniyor mu?</span><strong>'+(status.using_embedded_key?'Evet':'Hayir')+'</strong></div>'+
          '<div class="metric-row"><span>Lisans ID</span><span class="mono-wrap">'+escHtml(String(status.license_id||'-'))+'</span></div>'+
          '<div class="metric-row"><span>Bakim Bitisi</span><strong>'+escHtml(String(status.maintenance_until||'-'))+'</strong></div>'+
          '<div class="metric-row"><span>Maksimum Node</span><strong>'+escHtml(String(status.max_nodes||1))+'</strong></div>'+
        '</div>'+
        '<div style="margin-top:14px">'+
          (features.length?features.map(function(item){return '<span class="tag tag-blue">'+escHtml(String(item))+'</span>'}).join(''):'<div class="form-hint">Lisans yuklenince aktif ozellikler burada gorunur.</div>')+
        '</div>'+
        (warnings.length?'<div class="form-hint" style="margin-top:14px">'+warnings.map(function(item){return escHtml(String(item));}).join('<br>')+'</div>':'')+
      '</div>'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">Lisans Dosyasi Yukle</div>'+
        '<div class="form-group"><label class="form-label">Lisans JSON</label><textarea class="form-textarea" id="license-json-input" style="min-height:220px" placeholder="Imzali lisans JSONunu buraya yapistirin"></textarea></div>'+
        '<div class="form-group"><label class="form-label">Public Key PEM (opsiyonel)</label><textarea class="form-textarea" id="license-public-key-input" style="min-height:140px" placeholder="Ozel bir public key kullanacaksaniz buraya yapistirin"></textarea><div class="form-hint">Bos birakirsaniz uygulamanin icindeki development public key kullanilir.</div></div>'+
        '<div style="display:flex;gap:10px;flex-wrap:wrap"><button class="btn btn-primary" onclick="saveLicenseConfig()">Lisansi Kaydet</button><button class="btn btn-secondary" onclick=\'loadSampleLicense('+JSON.stringify(JSON.stringify(sample,null,2))+')\'>Ornek JSON Yukle</button></div>'+
      '</div>'+
    '</div>';
}

function loadSampleLicense(sampleJSON){
  const el=document.getElementById('license-json-input');
  if(el)el.value=sampleJSON;
}

async function saveLicenseConfig(){
  const licenseJSON=document.getElementById('license-json-input')?.value||'';
  const publicKeyPEM=document.getElementById('license-public-key-input')?.value||'';
  const res=await api('/api/license/upload',{method:'POST',body:{license_json:licenseJSON,public_key_pem:publicKeyPEM}});
  if(res&&res.success){
    toast('Lisans kaydedildi');
    loadPage('license');
  }else{
    toast((res&&res.message)||'Lisans kaydedilemedi','error');
  }
}

function settingSelect(key,label,value,options,hint){
  return '<div class="form-group"><label class="form-label">'+label+'</label><select class="form-select setting-input" data-key="'+key+'">'+
    (options||[]).map(function(opt){
      return '<option value="'+escHtml(String(opt.value))+'" '+(String(opt.value)===String(value)?'selected':'')+'>'+escHtml(String(opt.label))+'</option>';
    }).join('')+
    '</select>'+(hint?'<div class="form-hint">'+hint+'</div>':'')+'</div>';
}

</script>
<script src="/static/admin-studio.js"></script>
<script>
init();
</script>
</body>
</html>` + "`"
