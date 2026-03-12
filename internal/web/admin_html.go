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
.form-input,.form-select,.form-textarea{width:100%;padding:10px 14px;background:var(--bg-input);border:1px solid var(--border);border-radius:var(--radius-sm);color:var(--text-primary);font-size:14px;font-family:inherit;transition:border-color .2s;outline:none}
.form-input:focus,.form-select:focus,.form-textarea:focus{border-color:var(--border-focus);box-shadow:0 0 0 3px var(--accent-glow)}
.form-input::placeholder{color:var(--text-muted)}
.form-textarea{min-height:80px;resize:vertical}
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
.timeline-bars{display:grid;grid-template-columns:repeat(auto-fit,minmax(12px,1fr));align-items:end;gap:6px;height:150px}
.timeline-col{display:flex;flex-direction:column;justify-content:flex-end;align-items:center;gap:8px}
.timeline-bar{width:100%;border-radius:999px 999px 4px 4px;background:var(--gradient-3);min-height:6px;box-shadow:inset 0 -1px 0 rgba(255,255,255,.2)}
.timeline-label{font-size:10px;color:var(--text-muted)}
.timeline-value{font-size:10px;color:var(--text-secondary)}
.viewer-table td,.viewer-table th{font-size:12px}
.mono-wrap{font-family:Consolas,monospace;font-size:12px;word-break:break-all;color:var(--text-secondary)}
@media(max-width:980px){.quick-grid{grid-template-columns:1fr}}
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

async function api(path,opts={}){
  try{
    const hdrs={'Content-Type':'application/json',...opts.headers};
    if(authToken) hdrs['Authorization']='Bearer '+authToken;
    const res=await fetch(API+path,{
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
  const status=await api('/api/setup/status');
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
  '<h1 class="wizard-title">FluxStream</h1><p class="wizard-subtitle">Yonetim paneline giris yapin</p>'+
  '<div class="form-group"><label class="form-label">Kullanici Adi</label><input class="form-input" id="login-user" value="admin"></div>'+
  '<div class="form-group"><label class="form-label">Sifre</label><input class="form-input" id="login-pass" type="password" placeholder="Sifreniz"></div>'+
  '<button class="btn btn-primary" style="width:100%" onclick="doLogin()">Giris Yap</button></div></div>';
}
async function doLogin(){
  const u=document.getElementById('login-user').value;
  const p=document.getElementById('login-pass').value;
  if(!u||!p){toast('Kullanici adi ve sifre gerekli','error');return}
  const res=await api('/api/auth/login',{method:'POST',body:{username:u,password:p}});
  if(res.success){authToken=res.token;sessionStorage.setItem('fluxstream_token',authToken);toast('Giris basarili!');renderApp()}
  else{toast(res.message||'Giris hatasi','error')}
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETUP WIZARD Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
let wizardStep=1;
const wizardData={username:'admin',password:'',http_port:8844,https_port:443,rtmp_port:1935,embed_domain:''};

function renderWizard(){
  document.getElementById('app').innerHTML='<div class="wizard-container">'+getWizardContent()+'</div>';
}
function getWizardContent(){
  const steps=[
    '<div class="wizard-card"><div style="text-align:center;font-size:48px;margin-bottom:16px;color:var(--accent)"><i class="bi bi-lightning-charge-fill"></i></div><h1 class="wizard-title">FluxStream</h1><p class="wizard-subtitle">Live Streaming Media Server</p><div style="text-align:center;margin-bottom:24px">'+stepDots(1)+'</div><p style="text-align:center;color:var(--text-secondary);margin-bottom:32px;line-height:1.7">FluxStream\'e hos geldiniz!<br>Canli yayin sunucunuzu birkac adimda kuralim.</p><button class="btn btn-primary" style="width:100%" onclick="wizardNext()">Baslayalim <i class="bi bi-arrow-right"></i></button></div>',
    '<div class="wizard-card"><h1 class="wizard-title">Admin Hesabi</h1><p class="wizard-subtitle">Yonetim paneli icin giris bilgileri</p><div style="text-align:center;margin-bottom:24px">'+stepDots(2)+'</div><div class="form-group"><label class="form-label">Kullanici Adi</label><input class="form-input" id="w-username" value="admin"></div><div class="form-group"><label class="form-label">Sifre</label><input class="form-input" id="w-password" type="password" placeholder="En az 4 karakter"></div><div class="form-group"><label class="form-label">Sifre Tekrar</label><input class="form-input" id="w-password2" type="password"></div><div style="display:flex;gap:12px"><button class="btn btn-secondary" style="flex:1" onclick="wizardPrev()"><i class="bi bi-arrow-left"></i> Geri</button><button class="btn btn-primary" style="flex:1" onclick="wizardNext()">Ileri <i class="bi bi-arrow-right"></i></button></div></div>',
    '<div class="wizard-card"><h1 class="wizard-title">Port ve Domain</h1><p class="wizard-subtitle">Sunucu portlarini ve public alan adini yapilandirin</p><div style="text-align:center;margin-bottom:24px">'+stepDots(3)+'</div><div class="form-group"><label class="form-label">HTTP Port (Web Arayuzu)</label><input class="form-input" id="w-http-port" type="number" value="8844"></div><div class="form-group"><label class="form-label">HTTPS Port (SSL aktifse)</label><input class="form-input" id="w-https-port" type="number" value="443"></div><div class="form-group"><label class="form-label">RTMP Port (OBS Yayin)</label><input class="form-input" id="w-rtmp-port" type="number" value="1935"></div><div class="form-group"><label class="form-label">Public Domain / IP</label><input class="form-input" id="w-embed-domain" placeholder="Orn: stream.ornek.com veya 203.0.113.10"></div><div style="background:var(--bg-primary);border-radius:var(--radius-sm);padding:14px;margin-bottom:20px"><div style="font-size:13px;color:var(--text-muted)">Bos birakirsaniz panelin acildigi host kullanilir. HTTP ve HTTPS public portlari kurulumdan sonra Kolay Ayarlar veya Alan Adi / Embed ekranindan degistirilebilir.</div></div><div style="display:flex;gap:12px"><button class="btn btn-secondary" style="flex:1" onclick="wizardPrev()"><i class="bi bi-arrow-left"></i> Geri</button><button class="btn btn-primary" style="flex:1" onclick="wizardFinish()">Kurulumu Tamamla</button></div></div>'
  ];
  return steps[wizardStep-1]||steps[0];
}
function stepDots(c){let d='';for(let i=1;i<=3;i++){d+='<span class="wizard-dot'+(i===c?' active':i<c?' done':'')+'"></span>'}return d}
function wizardNext(){
  if(wizardStep===2){
    const pw=document.getElementById('w-password').value;
    const pw2=document.getElementById('w-password2').value;
    const user=document.getElementById('w-username').value;
    if(!pw||pw.length<4){toast('Sifre en az 4 karakter olmali','error');return}
    if(pw!==pw2){toast('Sifreler eslesiyor!','error');return}
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
    setupCompleted=true;toast('Kurulum tamamlandi!');
    // Auto-login after setup
    const lr=await api('/api/auth/login',{method:'POST',body:{username:wizardData.username,password:wizardData.password}});
    if(lr.success){authToken=lr.token;sessionStorage.setItem('fluxstream_token',authToken)}
    setTimeout(()=>renderApp(),500);
  }
  else{toast(res.message||'Kurulum hatasi','error')}
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â MAIN APP Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
function renderApp(){
  document.getElementById('app').innerHTML=
  '<div class="app">'+
    '<nav class="sidebar" id="sidebar">'+
      '<div class="logo"><div class="logo-icon"><i class="bi bi-lightning-charge-fill"></i></div><div><div class="logo-text">FluxStream</div><div class="logo-version">v2.0.0</div></div></div>'+
      '<div class="nav">'+
        '<div class="nav-section"><div class="nav-section-title">Ana Menu</div>'+
          navItem('dashboard','bi-bar-chart-fill','Dashboard')+
        '</div>'+
        '<div class="nav-section"><div class="nav-section-title">Yayin</div>'+
          navItem('streams','bi-collection-play-fill','Yayinlar')+
          navItem('create-stream','bi-plus-circle-fill','Yeni Yayin')+
          navItem('embed-codes','bi-code-slash','Embed Kodlari')+
          navItem('embed-advanced','bi-sliders','Gelismis Embed')+
          navItem('player-templates','bi-pc-display','Player Sablonlari')+
        '</div>'+
        '<div class="nav-section"><div class="nav-section-title">Ayarlar</div>'+
          navItem('guided-settings','bi-magic','Kolay Ayarlar')+
          navItem('settings-general','bi-gear-fill','Genel')+
          navItem('settings-embed','bi-globe2','Alan Adi / Embed')+
          navItem('settings-protocols','bi-diagram-3-fill','Protokoller')+
          navItem('settings-outputs','bi-boxes','Cikis Formatlari')+
          navItem('settings-abr','bi-badge-hd','Teslimat / ABR')+
          navItem('settings-ssl','bi-shield-lock-fill','SSL/TLS')+
          navItem('settings-security','bi-shield-shaded','Guvenlik')+
          navItem('settings-storage','bi-hdd-fill','Depolama')+
          navItem('settings-health','bi-heart-pulse-fill','Saglik ve Uyari')+
          navItem('settings-transcode','bi-cpu-fill','Transkod')+
        '</div>'+
        '<div class="nav-section"><div class="nav-section-title">Izleme</div>'+
          navItem('analytics','bi-graph-up','Analitik')+
          navItem('recordings','bi-camera-reels-fill','Kayitlar')+
          navItem('viewers','bi-people-fill','Izleyiciler')+
          navItem('transcode-jobs','bi-cpu','Transcode Isleri')+
          navItem('diagnostics','bi-activity','Teshis')+
        '</div>'+
        '<div class="nav-section"><div class="nav-section-title">Sistem</div>'+
          navItem('security-tokens','bi-key-fill','Tokenlar')+
          navItem('users','bi-person-fill','Kullanicilar')+
          navItem('logs','bi-journal-text','Loglar')+
        '</div>'+
      '</div>'+
    '</nav>'+
    '<div class="main">'+
      '<div class="topbar">'+
        '<div id="proto-status" class="proto-status"></div>'+
        '<div class="topbar-actions">'+
          '<button class="btn btn-secondary btn-sm" onclick="openSystemControl()"><i class="bi bi-power"></i> Sunucu Kontrol</button>'+
          '<span style="font-size:13px;color:var(--text-muted)" id="clock"></span>'+
        '</div>'+
      '</div>'+
      '<div class="content" id="page-content"></div>'+
    '</div>'+
  '</div>';
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
  loadPage(page);
}
function startClock(){setInterval(()=>{const el=document.getElementById('clock');if(el)el.textContent=new Date().toLocaleTimeString('tr-TR')},1000)}
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
}

function closeModal(id){
  const el=document.getElementById(id);
  if(el&&el.parentNode)el.parentNode.removeChild(el);
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
  else if(page==='recordings')await renderRecordings(c);
  else if(page==='viewers')await renderViewers(c);
  else if(page==='security-tokens')await renderSecurityTokens(c);
  else if(page==='transcode-jobs')await renderTranscodeJobs(c);
  else if(page==='diagnostics')await renderDiagnostics(c);
  else c.innerHTML='<div class="empty-state"><div class="icon"><i class="bi bi-cone-striped"></i></div><h3>Yakinda</h3></div>';
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â DASHBOARD Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderDashboard(c){
  const [stats,streams,analytics,health]=await Promise.all([api('/api/stats'),api('/api/streams'),api('/api/analytics'),api('/api/health/report')]);
  const live=(streams||[]).filter(s=>s.status==='live');
  const fmtItems=Object.entries((analytics&&analytics.viewers_by_format)||{}).sort((a,b)=>b[1]-a[1]).map(([label,value])=>({label:label,value:value}));
  const topStreams=((analytics&&analytics.top_streams)||[]).slice(0,5);
  const alerts=Array.isArray(health&&health.alerts)?health.alerts:[];
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Dashboard</h1><div style="display:flex;gap:10px;flex-wrap:wrap"><button class="btn btn-secondary btn-sm" onclick="navigate(\'analytics\')"><i class="bi bi-graph-up"></i> Analitik</button><button class="btn btn-primary btn-sm" onclick="navigate(\'create-stream\')"><i class="bi bi-plus-circle"></i> Yeni Yayin</button></div></div>'+
    '<div class="card-grid card-grid-4" style="margin-bottom:24px">'+
      statCard('purple','bi-broadcast',fmtInt(stats.active_streams||0),'Aktif Yayin','streams','Canli stream listesini ac')+
      statCard('blue','bi-people-fill',fmtInt(stats.total_viewers||0),'Aktif Izleyici','viewers','Anlik izleyici oturumlarini ac')+
      statCard('green','bi-clock-history',formatUptime(stats.uptime_seconds),'Calisma Suresi','analytics','24 saatlik hareketleri ac')+
      statCard('red','bi-memory',(stats.memory_used_mb||0)+' MB','Bellek Kullanimi','transcode-jobs','Transcode ve sistem detaylarini ac')+
    '</div>'+
    '<div class="quick-grid" style="margin-bottom:24px">'+
      '<div class="card">'+
        '<div class="card-header"><div class="card-title"><i class="bi bi-broadcast-pin title-icon"></i>Aktif Yayinlar</div><button class="btn btn-sm btn-primary" onclick="navigate(\'streams\')">Tumunu Ac</button></div>'+
        (live.length===0
          ?'<div class="empty-state"><div class="icon"><i class="bi bi-broadcast"></i></div><h3>Aktif yayin yok</h3><p style="color:var(--text-muted)">Yeni bir yayin olusturun ve OBS ile baglanin</p></div>'
          :'<div class="card-grid card-grid-2">'+live.map(streamCard).join('')+'</div>')+
      '</div>'+
      '<div class="card"><div class="card-title" style="margin-bottom:12px"><i class="bi bi-lightning-charge title-icon"></i>Hizli Bakis</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Toplam giris veri</span><strong>'+formatBytes(stats.bandwidth_in||0)+'</strong></div>'+
          '<div class="metric-row"><span>Toplam cikis veri</span><strong>'+formatBytes(stats.bandwidth_out||0)+'</strong></div>'+
          '<div class="metric-row"><span>Bellek ayagi</span><strong>'+(stats.memory_total_mb||0)+' MB</strong></div>'+
          '<div class="metric-row"><span>Top format</span><strong>'+(fmtItems[0]?escHtml(fmtItems[0].label):'Yok')+'</strong></div>'+
          '<div class="metric-row"><span>Saglik durumu</span><strong>'+escHtml(String((health&&health.status)||'ok').toUpperCase())+'</strong></div>'+
        '</div>'+
        '<div style="display:flex;gap:10px;margin-top:14px"><button class="btn btn-secondary btn-sm" onclick="openSystemControl()">Sunucu Kontrol</button><button class="btn btn-secondary btn-sm" onclick="navigate(\'settings-health\')">Saglik</button></div>'+
      '</div>'+
    '</div>'+
    '<div class="insight-grid">'+
      '<div class="card"><div class="card-header"><h3 class="card-title">24 Saat Izleyici Akisi</h3><button class="btn btn-sm btn-secondary" onclick="navigate(\'analytics\')">Detay</button></div><div class="card-body">'+renderTimelineChart((analytics&&analytics.viewers_timeline)||[],'Henuz izleyici verisi yok',function(v){return String(v)})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Format Dagilimi</h3><button class="btn btn-sm btn-secondary" onclick="navigate(\'embed-advanced\')">Embed</button></div><div class="card-body">'+renderBarList(fmtItems,'Henuz format verisi yok',function(v){return String(v)})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Populer Yayinlar</h3><button class="btn btn-sm btn-secondary" onclick="navigate(\'analytics\')">Analitik</button></div><div class="card-body">'+
        (topStreams.length?topStreams.map(function(s){
          const label=escHtml(s.stream_name||shortKey(s.stream_key));
          const sid=findStreamIdByKey((streams||[]),s.stream_key);
          const action=sid?' onclick="navigate(\'stream-detail-'+sid+'\')" style="cursor:pointer"':'';
          return '<div class="metric-row"'+action+'><span>'+label+'</span><span class="badge">'+fmtInt(s.viewers||0)+' izleyici</span></div>';
        }).join(''):'<div style="color:var(--text-muted)">Henuz populer yayin verisi yok</div>')+
      '</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Uyarilar</h3><button class="btn btn-sm btn-secondary" onclick="navigate(\'settings-health\')">Yonet</button></div><div class="card-body">'+
        (alerts.length?alerts.slice(0,4).map(function(item){
          return '<div class="metric-row"><span>'+escHtml(item.title||item.code||'Uyari')+'</span><span class="tag '+(item.level==='critical'?'tag-red':item.level==='warning'?'tag-yellow':'tag-blue')+'">'+escHtml(String(item.level||'info').toUpperCase())+'</span></div>';
        }).join(''):'<div style="color:var(--text-muted)">Aktif uyari yok</div>')+
      '</div></div>'+
    '</div>';
  schedulePageRefresh('dashboard',5000);
}
function statCard(color,iconClass,value,label,route,subtext){
  const clickable=route?' clickable':'';
  const action=route?' onclick="navigate(\''+route+'\')"':'';
  return '<div class="stat-card '+color+clickable+'"'+action+'><div class="stat-icon"><i class="bi '+iconClass+'"></i></div><div class="stat-value">'+value+'</div><div class="stat-label">'+label+'</div>'+(subtext?'<div class="stat-subtext">'+subtext+'</div>':'')+'</div>';
}
function fmtInt(n){return Number(n||0).toLocaleString('tr-TR')}
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
function renderTimelineChart(points,emptyText,formatter){
  const list=Array.isArray(points)?points:[];
  if(!list.length)return '<div style="color:var(--text-muted)">'+(emptyText||'Henuz veri yok')+'</div>';
  const max=Math.max.apply(null,list.map(function(p){return Number(p.value||0)}).concat([1]));
  return '<div class="timeline-bars">'+list.map(function(point){
    const value=Number(point.value||0);
    const date=point.timestamp?new Date(point.timestamp):null;
    const label=date?date.toLocaleTimeString('tr-TR',{hour:'2-digit',minute:'2-digit'}):'';
    const height=Math.max(6,Math.round((value/max)*100));
    return '<div class="timeline-col"><div class="timeline-value">'+escHtml(formatter?formatter(value):String(value))+'</div><div class="timeline-bar" style="height:'+height+'%"></div><div class="timeline-label">'+escHtml(label)+'</div></div>';
  }).join('')+'</div>';
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

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â STREAMS LIST Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderStreams(c){
  const streams=await api('/api/streams')||[];
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Yayinlar</h1>'+
      '<button class="btn btn-primary" onclick="navigate(\'create-stream\')">+ Yeni Yayin</button></div>'+
    '<div class="card">'+(streams.length===0
      ?'<div class="empty-state"><div class="icon"><i class="bi bi-broadcast"></i></div><h3>Henuz yayin yok</h3><p style="color:var(--text-muted)">Ilk yayininizi olusturun</p></div>'
      :'<table><thead><tr><th>Yayin</th><th>Durum</th><th>Stream Key</th><th>Izleyici</th><th>Codec</th><th></th></tr></thead><tbody>'+
        streams.map(s=>'<tr onclick="navigate(\'stream-detail-'+s.id+'\')" style="cursor:pointer">'+
          '<td><strong>'+escHtml(s.name)+'</strong></td>'+
          '<td><span class="badge badge-'+s.status+'">'+(s.status==='live'?'CANLI':'Cevrimdisi')+'</span></td>'+
          '<td><code style="font-size:12px;color:var(--accent)">'+s.stream_key+'</code></td>'+
          '<td>'+(s.viewer_count||0)+'</td>'+
          '<td style="font-size:12px;color:var(--text-muted)">'+(s.input_codec||'-')+'</td>'+
          '<td><button class="btn btn-sm btn-danger" onclick="event.stopPropagation();deleteStream('+s.id+')"><i class="bi bi-trash"></i></button></td>'+
        '</tr>').join('')+
      '</tbody></table>')+
    '</div>';
}
async function deleteStream(id){
  if(!confirm('Bu yayini silmek istediginize emin misiniz?'))return;
  await api('/api/streams/'+id,{method:'DELETE'});toast('Yayin silindi');navigate('streams');
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â CREATE STREAM Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
function renderCreateStream(c){
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Yeni Yayin Olustur</h1></div>'+
    '<div class="card" style="max-width:640px">'+
      '<div class="form-group"><label class="form-label">Yayin Adi *</label><input class="form-input" id="cs-name" placeholder="Orn: Konser Canli Yayin"></div>'+
      '<div class="form-group"><label class="form-label">Aciklama</label><input class="form-input" id="cs-desc" placeholder="Kisa aciklama"></div>'+
      '<div class="form-group"><label class="form-label">Yayin Modu</label><select class="form-select" id="cs-mode"><option value="balanced">TV / Dengeli</option><option value="mobile">Mobil / Hafif</option><option value="radio">Radyo / Audio</option></select><div class="form-hint">Bu secim ABR, cikis ve kaynak kullanimini belirleyen baslangic davranisini tanimlar.</div></div>'+
      '<div class="setting-row"><div><div class="setting-label">Adaptif Bitrate</div><div class="setting-desc">Acilirsa izleyicinin baglantisina gore kalite otomatik degisir.</div></div>'+
        '<label class="toggle"><input type="checkbox" id="cs-abr-enabled"><span class="toggle-slider"></span></label></div>'+
      '<div class="form-group" style="margin-top:16px"><label class="form-label">ABR Profil Seti</label><select class="form-select" id="cs-profile-set"><option value="balanced">Dengeli</option><option value="mobile">Mobil</option><option value="radio">Radyo</option></select></div>'+
      '<div class="setting-row"><div><div class="setting-label">Playback Token Gerekli</div><div class="setting-desc">Bu yayini izlemek icin token aranir.</div></div>'+
        '<label class="toggle"><input type="checkbox" id="cs-token-required"><span class="toggle-slider"></span></label></div>'+
      '<div class="form-group" style="margin-top:16px"><label class="form-label">Domain Kilidi</label><input class="form-input" id="cs-domain-lock" placeholder="Orn: mysite.com, embed.partner.com"><div class="form-hint">Bossa her yerde acilir. Doluysa sadece bu domainlerden gelen embed/referer kabul edilir.</div></div>'+
      '<div class="form-group"><label class="form-label">IP Beyaz Liste</label><input class="form-input" id="cs-ip-whitelist" placeholder="Orn: 203.0.113.20, 10.0.0.0/24"></div>'+
      '<div class="form-group"><label class="form-label">Maks Izleyici</label><input class="form-input" id="cs-max-viewers" type="number" value="0"><div class="form-hint">0 sinirsiz anlamina gelir.</div></div>'+
      '<div class="form-group"><label class="form-label">Maks Bitrate (kbps)</label><input class="form-input" id="cs-max-bitrate" type="number" value="0"><div class="form-hint">Kaynak kontrolu icin opsiyoneldir.</div></div>'+
      '<div class="form-group"><label class="form-label">Acik Cikis Formatlari</label><div class="form-hint" style="margin-bottom:10px">Bu yayinin disariya hangi formatlarda servis edilecegini secin.</div>'+renderOutputSelector(defaultStreamOutputs(),'cs')+'</div>'+
      '<div class="setting-row"><div><div class="setting-label">Yayin kaydedilsin mi?</div><div class="setting-desc">Varsayilan olarak kapali. Kalici kayitlar data/recordings altina yazilir.</div></div>'+
        '<label class="toggle"><input type="checkbox" id="cs-record-enabled" onchange="toggleCreateRecordFormat()"><span class="toggle-slider"></span></label></div>'+
      '<div class="form-group" style="margin-top:16px"><label class="form-label">Kayit Formati</label><select class="form-select" id="cs-record-format" disabled>'+recordingFormatOptions('ts')+'</select></div>'+
      '<button class="btn btn-primary" onclick="createStream()">Yayin Olustur</button>'+
    '</div><div id="cs-result" class="hidden"></div>';
}
async function createStream(){
  const name=document.getElementById('cs-name').value;
  const desc=document.getElementById('cs-desc').value;
  const recordEnabled=document.getElementById('cs-record-enabled').checked;
  const recordFormat=document.getElementById('cs-record-format').value||'ts';
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
    const r=document.getElementById('cs-result');r.classList.remove('hidden');
    r.innerHTML='<div class="card" style="max-width:640px;margin-top:20px">'+
      '<div class="card-title" style="margin-bottom:16px">Yayin Hazir!</div>'+
      copyField('Stream Key',res.stream.stream_key)+
      copyField('OBS RTMP URL',res.rtmp_url)+
      copyField('HLS Izleme URL',urls.hls)+
      (access&&access.needs_token?'<div class="form-hint" style="margin-bottom:10px;color:var(--warning)">Bu yayinda playback token gerekli. Izleme linkine gecici token eklendi.</div>':'')+
      '<div style="margin-top:12px"><button class="btn btn-sm btn-primary" onclick="navigate(\'stream-detail-'+res.stream.id+'\')">Yayin Detaylarina Git <i class="bi bi-arrow-right"></i></button></div>'+
      '<div style="background:var(--bg-primary);border-radius:var(--radius-sm);padding:16px;margin-top:12px">'+
        '<div style="font-size:13px;color:var(--text-muted);line-height:1.6">OBS Studio\'da:<br>1. Ayarlar -> Yayin -> Hizmet: Ozel<br>2. Sunucu: <strong>'+escHtml(res.rtmp_url||'')+'</strong><br>3. Yayin Anahtari: <strong>'+res.stream.stream_key+'</strong><br>4. Yayina Baslat butonuna basin</div>'+
      '</div></div>';
  }else{toast(res.message||'Hata','error')}
}
function toggleCreateRecordFormat(){
  const enabled=document.getElementById('cs-record-enabled')?.checked;
  const format=document.getElementById('cs-record-format');
  if(format)format.disabled=!enabled;
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
function recordingFormatOptions(selected){
  selected=selected||'ts';
  return '<option value="ts"'+(selected==='ts'?' selected':'')+'>MPEG-TS (.ts)</option>'+
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

// â•â•â• STREAM DETAIL â•â•â•
async function renderStreamDetail(c,id){
  const st=await api('/api/streams/'+id);
  if(!st||st.error){c.innerHTML='<div class="empty-state"><h3>Yayin bulunamadi</h3></div>';return}
  window._streamDetailData=st;
  const settings=await api('/api/settings');
  const access=await getPlaybackAccess(st.stream_key,settings,st.policy_json);
  const u=getAllURLs(st.stream_key,settings,st.name,access);
  const previewURLs=getPreviewURLs(st.stream_key,settings,st.name,access);
  const policy=parseStreamPolicy(st.policy_json);
  const outputFormats=parseJSONSafe(st.output_formats,defaultStreamOutputs());

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

    '<div class="card" style="margin-bottom:16px"><div class="card-title" style="margin-bottom:12px">Teslimat ve Guvenlik Politikasi</div>'+
      '<div class="form-group"><label class="form-label">Yayin Modu</label><select class="form-select" id="sd-policy-mode"><option value="balanced" '+((policy.mode||'balanced')==='balanced'?'selected':'')+'>TV / Dengeli</option><option value="mobile" '+((policy.mode||'')==='mobile'?'selected':'')+'>Mobil / Hafif</option><option value="radio" '+((policy.mode||'')==='radio'?'selected':'')+'>Radyo / Audio</option></select><div class="form-hint">Bu, yayin icin secilen genel davranis profilidir.</div></div>'+
      '<div class="setting-row"><div><div class="setting-label">Adaptif Bitrate</div><div class="setting-desc">Acik oldugunda izleyiciye baglanti hizina gore farkli kalite katmanlari sunulur.</div></div>'+
      '<label class="toggle"><input type="checkbox" id="sd-abr-enabled" '+(policy.enable_abr?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      '<div class="form-group" style="margin-top:16px"><label class="form-label">ABR Profil Seti</label><select class="form-select" id="sd-profile-set"><option value="balanced" '+((policy.profile_set||'balanced')==='balanced'?'selected':'')+'>Dengeli</option><option value="mobile" '+((policy.profile_set||'')==='mobile'?'selected':'')+'>Mobil</option><option value="radio" '+((policy.profile_set||'')==='radio'?'selected':'')+'>Radyo</option></select></div>'+
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
      '<div class="form-group" style="margin-top:16px"><label class="form-label">Kayit Formati</label><select class="form-select" id="sd-record-format">'+recordingFormatOptions(st.record_format||'ts')+'</select></div>'+
      '<div class="form-hint">Kalici kayitlar <code>data/recordings</code> altina yazilir. Canli cache dizinleri yayin bitince temizlenir.</div>'+
      '<div style="margin-top:16px"><button class="btn btn-primary" onclick="saveStreamRecordSettings('+st.id+')">Kayit Ayarlarini Kaydet</button></div>'+
    '</div>'+

    '<div class="card" style="margin-bottom:16px"><div class="card-title" style="margin-bottom:12px">Embed Kodlari</div>'+
      (access&&access.needs_token?'<div class="form-hint" style="margin-bottom:10px;color:var(--warning)">Bu yayinda playback token gerekli. Aasagidaki preview ve linkler gecici token ile uretildi.</div>':'')+
      copyField('iframe','<iframe src="'+u.embed+'" width="1280" height="720" frameborder="0" allowfullscreen></iframe>')+
      copyField('Player URL',u.play)+
      copyField('Embed URL',u.embed)+
    '</div>'+

    (st.status==='live'?
      '<div class="card"><div class="card-title" style="margin-bottom:12px">Onizleme</div>'+
        '<div style="position:relative;padding-top:56.25%;background:#000;border-radius:8px;overflow:hidden">'+
          '<iframe src="'+previewURLs.play+'" style="position:absolute;top:0;left:0;width:100%;height:100%;border:none" allowfullscreen></iframe>'+
        '</div></div>':'');

}

// â•â•â• EMBED CODES â•â•â•
async function saveStreamRecordSettings(id){
  const st=window._streamDetailData;
  if(!st)return;
  const payload=Object.assign({},st,{
    record_enabled:document.getElementById('sd-record-enabled')?.checked||false,
    record_format:document.getElementById('sd-record-format')?.value||'ts'
  });
  const res=await api('/api/streams/'+id,{method:'PUT',body:payload});
  if(res&&res.success){
    window._streamDetailData=payload;
    toast('Kayit ayarlari kaydedildi');
  }else{
    toast((res&&res.message)||'Kaydedilemedi','error');
  }
}
async function saveStreamPolicySettings(id){
  const st=window._streamDetailData;
  if(!st)return;
  const policy={
    mode:document.getElementById('sd-policy-mode')?.value||'balanced',
    enable_abr:document.getElementById('sd-abr-enabled')?.checked||false,
    profile_set:document.getElementById('sd-profile-set')?.value||'balanced',
    require_playback_token:document.getElementById('sd-token-required')?.checked||false
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
    window._streamDetailData=payload;
    toast('Politika kaydedildi');
  }else{
    toast((res&&res.message)||'Kaydedilemedi','error');
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
          '<details style="margin-top:4px"><summary style="cursor:pointer;font-weight:600;padding:8px 0;color:var(--text-secondary)">Player & Embed (2)</summary>'+
            copyField('Player URL',u.play)+copyField('Embed URL',u.embed)+
          '</details>'+
        '</div>';
      }).join(''));
  }catch(e){
    c.innerHTML='<div class="card"><div class="empty-state"><h3>Embed kodlari yuklenemedi</h3><p style="color:var(--text-muted)">'+escHtml(e.message||'Bilinmeyen hata')+'</p></div></div>';
  }
}
// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETTINGS - GENERAL Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderSettingsGeneral(c){
  const s=await api('/api/settings');
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Genel Ayarlar</h1></div>'+
    '<div class="card" style="max-width:700px">'+
      settingInput('server_name','Sunucu Adi',s.server_name||'FluxStream','text','Sunucu goruntuleme adi')+
      settingInput('http_port','HTTP Port',s.http_port||'8844','number','Web arayuzu portu')+
      settingInput('https_port','HTTPS Port',s.https_port||'443','number','SSL portu')+
      settingInput('language','Dil',s.language||'tr','text','Arayuz dili')+
      settingInput('timezone','Saat Dilimi',s.timezone||'Europe/Istanbul','text','')+
      '<button class="btn btn-primary" style="margin-top:8px" onclick="saveSettingsCategory(\'general\')">Kaydet</button>'+
    '</div>';
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETTINGS - PROTOCOLS Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderSettingsEmbed(c){
  const s=await api('/api/settings');
  const publicConfig=getPublicBase(s);
  const streamHost=getConfiguredDomain(s);
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Alan Adi ve Embed</h1></div>'+
    '<div class="card" style="max-width:760px;margin-bottom:16px">'+
      '<div class="card-title" style="margin-bottom:12px">Public Erisim Ayarlari</div>'+
      '<div style="font-size:13px;color:var(--text-muted);line-height:1.7;margin-bottom:16px">Bu sayfa, embed kodlari ve izleme linklerinde kullanilan public alan adini belirler. Web erisimi HTTP/HTTPS portlarindan gider, yayin gonderimi ise RTMP/RTMPS gibi kendi portlarindan devam eder.</div>'+
      settingInput('embed_domain','Public Domain / IP',s.embed_domain||'','text','Bos birakirsaniz panelin acildigi host kullanilir. Ornek: stream.ornek.com')+
      settingInput('embed_http_port','Public HTTP Port',s.embed_http_port||s.http_port||'8844','number','HTTP uzerinden uretilen player ve embed linkleri icin')+
      settingInput('embed_https_port','Public HTTPS Port',s.embed_https_port||s.https_port||'443','number','SSL aktifse genelde 443 kullanilir')+
      '<div class="setting-row"><div><div class="setting-label">Embed Linklerinde HTTPS Kullan</div><div class="setting-desc">Yalnizca SSL etkin ve sertifika hazirsa player ve embed linkleri HTTPS olarak uretilir. Hazir degilse sistem HTTP kullanir.</div></div>'+
        '<label class="toggle"><input type="checkbox" class="setting-input" data-key="embed_use_https" '+(shouldUsePublicHTTPS(s)?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      '<button class="btn btn-primary" style="margin-top:8px" onclick="saveSettingsCategory(\'embed\')">Kaydet</button>'+
    '</div>'+
    '<div class="card" style="max-width:760px">'+
      '<div class="card-title" style="margin-bottom:12px">Canli Ornek</div>'+
      copyField('Player Base URL',publicConfig.base)+
      copyField('RTMP Sunucu','rtmp://'+streamHost+':'+(s.rtmp_port||'1935')+'/live')+
      copyField('RTMPS Sunucu','rtmps://'+streamHost+':'+(s.rtmps_port||'1936')+'/live')+
    '</div>';
}

async function renderSettingsProtocols(c){
  const s=await api('/api/settings');
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Giris Protokolleri</h1>'+
      '<p style="color:var(--text-muted);font-size:13px">Encoder\'lardan kabul edilen protokoller</p></div>'+
    protoCard('RTMP','En yaygin - OBS, Wirecast, vMix','rtmp_enabled',s,'rtmp_port',s.rtmp_port||'1935',[
      settingInput('rtmp_chunk_size','Chunk Size',s.rtmp_chunk_size||'4096','number',''),
      settingInput('rtmp_max_conns','Maks Baglanti',s.rtmp_max_conns||'100','number',''),
    ])+
    protoCard('RTMPS','Sifreli RTMP - guvenli aglar','rtmps_enabled',s,'rtmps_port',s.rtmps_port||'1936',[])+
    protoCard('SRT','Dusuk gecikme + guvenilmez aglarda guclu','srt_enabled',s,'srt_port',s.srt_port||'9000',[
      settingInput('srt_latency','Latency (ms)',s.srt_latency||'120','number',''),
    ])+
    protoCard('RTP','Profesyonel encoder push','rtp_enabled',s,'rtp_port',s.rtp_port||'5004',[])+
    protoCard('RTSP','IP kameralar, profesyonel encoderlar','rtsp_enabled',s,'rtsp_port',s.rtsp_port||'8554',[])+
    protoCard('WebRTC/WHIP','Tarayicidan dogrudan yayin','webrtc_enabled',s,'webrtc_port',s.webrtc_port||'8855',[])+
    protoCard('MPEG-TS','Uydu alici, profesyonel broadcast','mpegts_enabled',s,'mpegts_port',s.mpegts_port||'9001',[])+
    '<button class="btn btn-primary" onclick="saveSettingsCategory(\'protocols\')">Tum Protokolleri Kaydet</button>';
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

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETTINGS - OUTPUTS Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderSettingsOutputs(c){
  const s=await api('/api/settings');
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Cikis Formatlari</h1>'+
      '<p style="color:var(--text-muted);font-size:13px">Izleyicilere sunulan formatlar</p></div>'+
    outputCard('HLS','Apple HLS - en uyumlu','hls_enabled',s,[
      settingInput('hls_segment_duration','Segment Suresi (sn)',s.hls_segment_duration||'2','number',''),
      settingInput('hls_playlist_length','Playlist Uzunlugu',s.hls_playlist_length||'6','number',''),
    ])+
    outputCard('Low Latency HLS','2 saniye alti gecikme','hls_ll_enabled',s,[])+
    outputCard('DASH','MPEG-DASH adaptive','dash_enabled',s,[
      settingInput('dash_segment_duration','Segment Suresi (sn)',s.dash_segment_duration||'2','number',''),
    ])+
    outputCard('HTTP-FLV','Ultra dusuk gecikme (~1sn)','httpflv_enabled',s,[])+
    outputCard('WebRTC/WHEP','Sub-second gecikme (<500ms)','whep_enabled',s,[])+
    outputCard('MP3/Icecast','MP3 ses akisi','mp3_enabled',s,[
      settingInput('mp3_bitrate','Bitrate (kbps)',s.mp3_bitrate||'128','number',''),
    ])+
    outputCard('AAC','Yuksek kalite ses','aac_out_enabled',s,[])+
    outputCard('Icecast','Icecast uyumlu akis','icecast_enabled',s,[
      settingInput('icecast_port','Icecast Port',s.icecast_port||'8000','number',''),
    ])+
    '<button class="btn btn-primary" onclick="saveSettingsCategory(\'outputs\')">Tum Cikislari Kaydet</button>';
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

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETTINGS - SSL Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderSettingsSSL(c){
  const s=await api('/api/settings');
  const sslStatus=await api('/api/ssl/status');
  const webSSL=(sslStatus&&sslStatus.web)||{};
  const streamSSL=(sslStatus&&sslStatus.stream)||{};
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">SSL/TLS Sertifika</h1></div>'+
    '<div class="card" style="max-width:920px;margin-bottom:16px">'+
      '<div class="card-title" style="margin-bottom:12px">Kullanim Mantigi</div>'+
      '<div class="form-hint" style="line-height:1.8">Web HTTPS sertifikasi admin paneli ve embed/player sayfalari icin kullanilir. Stream SSL ise yalnizca RTMPS ingest tarafini korur. Isterseniz ayni domaini, isterseniz ayri domain ve ayri sertifika kullanabilirsiniz. Let\\'s Encrypt icin alan adlarinin bu VPS\\'e yonlenmis olmasi ve 80/443 portlarinin acik olmasi gerekir.</div>'+
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
    '<div class="form-group"><label class="form-label">Sertifika Modu</label><select class="form-select setting-input" data-key="'+modeKey+'"><option value="file" '+(mode==='file'?'selected':'')+'>Manuel CRT/KEY</option><option value="letsencrypt" '+(mode==='letsencrypt'?'selected':'')+'>Let\\'s Encrypt</option></select><div class="form-hint">Manuel modda dosya yuklersiniz. Let\\'s Encrypt modunda domain ve e-posta yeterlidir.</div></div>'+
    '<div class="form-group"><label class="form-label">CRT / PEM Yukle</label><input type="file" id="ssl-cert-file-'+target+'" accept=".crt,.pem,.cert" class="form-input" style="padding:8px"></div>'+
    '<div class="form-group"><label class="form-label">KEY / PEM Yukle</label><input type="file" id="ssl-key-file-'+target+'" accept=".key,.pem" class="form-input" style="padding:8px"></div>'+
    '<div style="margin-bottom:16px"><button class="btn btn-secondary" onclick="uploadSSL(\''+target+'\')">Bu Profil Icin Sertifika Yukle</button></div>'+
    settingInput(certKey,'Sertifika Dosyasi (.crt)',s[certKey]||'','text','Orn: /opt/fluxstream/data/certs/'+target+'/server.crt')+
    settingInput(keyKey,'Ozel Anahtar (.key)',s[keyKey]||'','text','Orn: /opt/fluxstream/data/certs/'+target+'/server.key')+
    settingInput(domainKey,'Let\\'s Encrypt Domain',s[domainKey]||'','text','Orn: '+(target==='web'?'panel.example.com':'stream.example.com'))+
    settingInput(emailKey,'Let\\'s Encrypt E-posta',s[emailKey]||'','text','Bildirim ve yenileme icin kullanilir.')+
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
// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETTINGS - SECURITY Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderSettingsSecurity(c){
  const s=await api('/api/settings');
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Guvenlik Ayarlari</h1></div>'+
    '<div class="card" style="max-width:700px">'+
      '<div class="setting-row"><div><div class="setting-label">Stream Key Zorunlu</div><div class="setting-desc">Yayin icin stream key gerektirir</div></div>'+
        '<label class="toggle"><input type="checkbox" class="setting-input" data-key="stream_key_required" '+(s.stream_key_required==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      '<div class="setting-row"><div><div class="setting-label">Token Dogrulama</div><div class="setting-desc">Izleme icin token gerektirir</div></div>'+
        '<label class="toggle"><input type="checkbox" class="setting-input" data-key="token_enabled" '+(s.token_enabled==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
      settingInput('token_duration','Token Suresi (sn)',s.token_duration||'60','number','')+
      settingInput('rate_limit','Rate Limit (istek/sn)',s.rate_limit||'100','number','')+
      '<button class="btn btn-primary" style="margin-top:8px" onclick="saveSettingsCategory(\'security\')">Kaydet</button>'+
    '</div>';
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETTINGS - STORAGE Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderSettingsStorage(c){
  const s=await api('/api/settings');
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Depolama Ayarlari</h1></div>'+
    '<div class="card" style="max-width:700px">'+
      settingInput('storage_max_gb','Maksimum Depolama (GB)',s.storage_max_gb||'50','number','Toplam disk kullanim limiti')+
      settingInput('storage_auto_clean','Otomatik Temizlik (gun)',s.storage_auto_clean||'30','number','Bu sureden eski dosyalar silinir')+
      '<button class="btn btn-primary" style="margin-top:8px" onclick="saveSettingsCategory(\'storage\')">Kaydet</button>'+
    '</div>';
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETTINGS - TRANSCODE Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
async function renderSettingsTranscode(c){
  const s=await api('/api/settings');
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Transkod / FFmpeg</h1></div>'+
    '<div class="card" style="max-width:700px">'+
      settingInput('ffmpeg_path','FFmpeg Yolu',s.ffmpeg_path||'ffmpeg','text','FFmpeg calistirilabilir dosya yolu')+
      '<div class="form-group"><label class="form-label">GPU Hizlandirma</label>'+
        '<select class="form-select setting-input" data-key="gpu_accel">'+
          '<option value="none" '+(s.gpu_accel==='none'?'selected':'')+'>Yok (CPU)</option>'+
          '<option value="nvenc" '+(s.gpu_accel==='nvenc'?'selected':'')+'>NVIDIA NVENC</option>'+
          '<option value="qsv" '+(s.gpu_accel==='qsv'?'selected':'')+'>Intel Quick Sync</option>'+
          '<option value="amf" '+(s.gpu_accel==='amf'?'selected':'')+'>AMD AMF</option>'+
        '</select></div>'+
      '<button class="btn btn-primary" style="margin-top:8px" onclick="saveSettingsCategory(\'transcode\')">Kaydet</button>'+
    '</div>';
}

async function renderGuidedSettings(c){
  const [s,health]=await Promise.all([api('/api/settings'),api('/api/health/report')]);
  const alerts=Array.isArray(health&&health.alerts)?health.alerts:[];
  const publicConfig=getPublicBase(s);
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Kolay Ayarlar</h1><div style="color:var(--text-muted);font-size:13px">En sik kullanilan tum ayarlar aciklamalariyla burada toplanir.</div></div>'+
    '<div class="card" style="margin-bottom:16px">'+
      '<div class="card-title" style="margin-bottom:10px">1. Yayin Profili Sec</div>'+
      '<div class="form-hint" style="margin-bottom:14px">Hazir profiller bir cocugun bile anlayabilecegi kadar sade tutuldu. Isterseniz sonra detay ekranlarindan ince ayar yapabilirsiniz.</div>'+
      '<div class="card-grid card-grid-3">'+
        presetCard('TV / Dengeli','Canli TV, YouTube benzeri genel kullanim. Adaptif kalite merdiveni acilir, 1080p-360p katmanlari kullanilir.','balanced','bi-broadcast-pin')+
        presetCard('Mobil / Hafif','Zayif baglantilar ve daha dusuk CPU kullanimi icin. Daha kucuk bitrate merdiveni kullanilir.','mobile','bi-phone')+
        presetCard('Radyo / Audio','Ses yayinlarina uygun sade ayar. Video gerekmiyorsa depolama ve transcode yukunu azaltir.','radio','bi-mic-fill')+
      '</div>'+
    '</div>'+
    '<div class="card-grid card-grid-2">'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">2. Public Erisim</div>'+
        settingInput('embed_domain','Public Domain / IP',s.embed_domain||'','text','Canli sunucuda localhost yerine uretilen tum embed ve player linkleri burada kullanilir.')+
        settingInput('embed_http_port','Public HTTP Port',s.embed_http_port||s.http_port||'8844','number','HTTP linklerinde kullanilan public port. Kurulumda verdiginiz web portu genelde dogru baslangic degeridir.')+
        settingInput('embed_https_port','Public HTTPS Port',s.embed_https_port||s.https_port||'443','number','HTTPS linklerinde kullanilan public port. SSL reverse proxy veya farkli bir port kullaniyorsaniz burayi degistirin.')+
        '<div class="setting-row"><div><div class="setting-label">HTTPS Link Uret</div><div class="setting-desc">Yalnizca SSL etkin ve sertifika hazirsa embed linkleri HTTPS olarak uretilir. Hazir degilse HTTP kullanilir. 443 sabit degildir.</div></div>'+
        '<label class="toggle"><input type="checkbox" class="guided-input" data-key="embed_use_https" '+(shouldUsePublicHTTPS(s)?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
        '<div class="form-hint" style="margin-top:12px;line-height:1.7">Ornek player tabani: <code>'+escHtml(publicConfig.base)+'</code></div>'+
        '<div style="margin-top:12px"><button class="btn btn-primary" onclick="saveGuidedPublic()">Public Ayarlari Kaydet</button></div>'+
      '</div>'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">3. Kayit ve Temizlik</div>'+
        '<div class="form-group"><label class="form-label">Kayitlari Kac Gun Tut</label><input class="form-input guided-input" data-key="recordings_retention_days" type="number" value="'+escHtml(s.recordings_retention_days||'30')+'"><div class="form-hint">Bu sureyi asan kayitlar otomatik silinir.</div></div>'+
        '<div class="form-group"><label class="form-label">Her Yayinda En Fazla Kac Kayit Sakla</label><input class="form-input guided-input" data-key="recordings_keep_latest" type="number" value="'+escHtml(s.recordings_keep_latest||'10')+'"><div class="form-hint">Disk kontrolu icin eski kayitlar trim edilir.</div></div>'+
        '<div class="setting-row"><div><div class="setting-label">Otomatik Bakim</div><div class="setting-desc">Kayit ve analytics temizlik islemleri belirli araliklarla otomatik calisir.</div></div>'+
        '<label class="toggle"><input type="checkbox" class="guided-input" data-key="maintenance_auto_cleanup" '+(s.maintenance_auto_cleanup==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
        '<div style="display:flex;gap:10px;margin-top:12px"><button class="btn btn-primary" onclick="saveGuidedStorage()">Depolama Ayarlarini Kaydet</button><button class="btn btn-secondary" onclick="runMaintenance()">Bakimi Simdi Calistir</button></div>'+
      '</div>'+
    '</div>'+
    '<div class="card-grid card-grid-2" style="margin-top:16px">'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">4. Guvenlik</div>'+
        '<div class="setting-row"><div><div class="setting-label">Playback Token</div><div class="setting-desc">Acilirsa izleyici linklerinde gecerli token aranir.</div></div><label class="toggle"><input type="checkbox" class="guided-input" data-key="token_enabled" '+(s.token_enabled==='true'?'checked':'')+'><span class="toggle-slider"></span></label></div>'+
        '<div class="form-group" style="margin-top:16px"><label class="form-label">Rate Limit</label><input class="form-input guided-input" data-key="rate_limit" type="number" value="'+escHtml(s.rate_limit||'100')+'"><div class="form-hint">Ani ve asiri istekleri sinirlar.</div></div>'+
        '<div style="margin-top:12px"><button class="btn btn-primary" onclick="saveGuidedSecurity()">Guvenligi Kaydet</button></div>'+
      '</div>'+
      '<div class="card">'+
        '<div class="card-title" style="margin-bottom:12px">5. Sistem Durumu</div>'+
        '<div class="metric-list">'+
          '<div class="metric-row"><span>Genel durum</span><strong>'+escHtml(String((health&&health.status)||'ok').toUpperCase())+'</strong></div>'+
          '<div class="metric-row"><span>Aktif uyari</span><strong>'+fmtInt(alerts.length)+'</strong></div>'+
          '<div class="metric-row"><span>ABR</span><strong>'+(s.abr_enabled==='true'?'Acik':'Kapali')+'</strong></div>'+
          '<div class="metric-row"><span>Kalici analitik</span><strong>'+(s.analytics_persist_enabled==='true'?'Acik':'Kapali')+'</strong></div>'+
        '</div>'+
        '<div style="display:flex;gap:10px;margin-top:14px"><button class="btn btn-secondary" onclick="navigate(\'settings-health\')">Saglik Ekranini Ac</button><button class="btn btn-secondary" onclick="openSystemControl()">Sunucu Kontrol</button></div>'+
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
  }else if(preset==='mobile'){
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
      '<div class="form-group" style="margin-top:16px"><label class="form-label">Hazir Profil Seti</label><select class="form-select setting-input" data-key="abr_profile_set"><option value="balanced" '+((s.abr_profile_set||'balanced')==='balanced'?'selected':'')+'>TV / Dengeli</option><option value="mobile" '+((s.abr_profile_set||'')==='mobile'?'selected':'')+'>Mobil / Hafif</option><option value="radio" '+((s.abr_profile_set||'')==='radio'?'selected':'')+'>Radyo / Audio</option></select><div class="form-hint">Dengeli cogu video yayin icin en iyi baslangic noktasidir.</div></div>'+
      '<div class="form-group"><label class="form-label">ABR Profil JSON</label><textarea class="form-textarea setting-input" data-key="abr_profiles_json" style="min-height:220px">'+escHtml(s.abr_profiles_json||'')+'</textarea><div class="form-hint">Gelistirilmis kullanim icin. Hazir setler yukaridan secilebilir.</div></div>'+
      '<button class="btn btn-primary" onclick="saveSettingsCategory(\'outputs\')">ABR Ayarlarini Kaydet</button>'+
    '</div>';
}

async function renderSettingsHealth(c){
  const [s,report]=await Promise.all([api('/api/settings'),api('/api/health/report')]);
  const alerts=Array.isArray(report&&report.alerts)?report.alerts:[];
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
        '<div style="display:flex;gap:10px"><button class="btn btn-primary" onclick="saveSettingsCategory(\'health\')">Esikleri Kaydet</button><button class="btn btn-secondary" onclick="runMaintenance()">Bakimi Simdi Calistir</button></div>'+
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

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â SETTINGS HELPERS Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
function settingInput(key,label,value,type,hint){
  return '<div class="form-group"><label class="form-label">'+label+'</label>'+
    '<input class="form-input setting-input" data-key="'+key+'" type="'+(type||'text')+'" value="'+escHtml(String(value||''))+'">'+
    (hint?'<div class="form-hint">'+hint+'</div>':'')+
  '</div>';
}
async function saveSettingsValues(category,updates,silent){
  const res=await api('/api/settings/'+category,{method:'PUT',body:updates});
  if(!silent){
    if(res&&res.success!==false)toast('Ayarlar kaydedildi!');
    else toast((res&&res.message)||'Kayit hatasi','error');
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
  out.innerHTML=
    '<div class="metric-list">'+
      '<div class="metric-row"><span>Yayin</span><strong>'+escHtml(data.stream_name||data.stream_key||'-')+'</strong></div>'+
      '<div class="metric-row"><span>ABR Profil</span><strong>'+escHtml(data.abr_profile_set||'balanced')+'</strong></div>'+
      '<div class="metric-row"><span>Policy JSON</span><span class="mono-wrap">'+escHtml(data.policy_json||'{}')+'</span></div>'+
    '</div>'+
    '<div class="bar-list" style="margin-top:16px">'+checks.map(function(check){
      const tone=check.ok?'tag-green':'tag-red';
      return '<div class="metric-row"><span>'+escHtml(check.description||check.code)+'</span><span class="tag '+tone+'">'+(check.ok?'Hazir':'Yok')+'</span></div>';
    }).join('')+'</div>';
}

// Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â LOGS Ã¢â€¢ÂÃ¢â€¢ÂÃ¢â€¢Â
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
          const time=new Date(l.created_at).toLocaleString('tr-TR');
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

// Ã¢â€¢Ã¢â€¢Ã¢â€¢ USERS Ã¢â€¢Ã¢â€¢Ã¢â€¢
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
          '<td style="font-size:12px;color:var(--text-muted)">'+new Date(u.created_at).toLocaleDateString('tr-TR')+'</td>'+
          '<td><button class="btn btn-sm btn-secondary" onclick="showEditUserModal('+u.id+',\''+escHtml(u.username)+'\',\''+escHtml(u.role)+'\')">Duzenle</button> '+
            '<button class="btn btn-sm btn-danger" onclick="deleteUser('+u.id+')">Sil</button></td>'+
        '</tr>').join('')+
      '</tbody></table>')+
    '</div><div id="user-modal"></div>';
}
function showAddUserModal(){
  document.getElementById('user-modal').innerHTML=
    '<div class="modal-overlay" onclick="if(event.target===this)this.innerHTML=\'\'">'+
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
    '<div class="modal-overlay" onclick="if(event.target===this)this.innerHTML=\'\'">'+
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

// Ã¢â€¢Ã¢â€¢Ã¢â€¢ PLAYER TEMPLATES Ã¢â€¢Ã¢â€¢Ã¢â€¢
async function renderPlayerTemplates(c){
  const templates=await api('/api/players')||[];
  c.innerHTML=
    '<div class="page-header"><div><h1 class="page-title">Player Sablonlari</h1><div style="color:var(--text-muted);font-size:13px;margin-top:6px">Kurulu gelen hazir sablonlari temel alip duzenleyebilir veya sifirdan yeni sablon olusturabilirsiniz.</div></div>'+
      '<button class="btn btn-primary" onclick="showPlayerModal()">+ Yeni Sablon</button></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="metric-list">'+
      '<div class="metric-row"><span>Hazir baslangic sablonu</span><strong>6+</strong></div>'+
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
          '<div style="height:120px;background:#000;border-radius:8px;display:flex;align-items:center;justify-content:center;position:relative;overflow:hidden;margin-bottom:12px">'+
            (t.logo_url?'<img src="'+escHtml(t.logo_url)+'" style="position:absolute;'+
              (t.logo_position==='top-left'?'top:8px;left:8px':t.logo_position==='top-right'?'top:8px;right:8px':t.logo_position==='bottom-left'?'bottom:8px;left:8px':'bottom:8px;right:8px')+
              ';height:24px;opacity:'+(t.logo_opacity||1)+'">':'')+
            '<div style="font-size:32px;opacity:.35"><i class="bi bi-play-circle-fill"></i></div>'+
            (t.watermark_text?'<div style="position:absolute;bottom:8px;left:50%;transform:translateX(-50%);font-size:11px;opacity:.4;color:#fff">'+escHtml(t.watermark_text)+'</div>':'')+
          '</div>'+
          '<div style="display:flex;gap:8px">'+
            '<button class="btn btn-sm btn-secondary" onclick="event.stopPropagation();showPlayerModal('+t.id+')">Duzenle</button>'+
            '<button class="btn btn-sm btn-danger" onclick="event.stopPropagation();deletePlayerTemplate('+t.id+')">Sil</button>'+
          '</div>'+
        '</div>'
      ).join('')+'</div>')+
    '<div id="player-modal"></div>';
}
async function showPlayerModal(id){
  let pt={name:'',theme:'dark',background_css:'',control_bar_css:'',play_button_css:'',logo_url:'',logo_position:'top-right',logo_opacity:1,watermark_text:'',show_title:true,show_live_badge:true,custom_css:''};
  if(id){const data=await api('/api/players/'+id);if(data&&!data.error)pt=data}
  const isEdit=!!id;
  document.getElementById('player-modal').innerHTML=
    '<div class="modal-overlay" onclick="if(event.target===this)this.innerHTML=\'\'">'+
      '<div class="modal" style="max-width:700px">'+
        '<div class="modal-title">'+(isEdit?'Sablon Duzenle':'Yeni Player Sablonu')+'</div>'+
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
      '</div></div>';
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

// Ã¢â€¢Ã¢â€¢Ã¢â€¢ ADVANCED EMBED GENERATOR Ã¢â€¢Ã¢â€¢Ã¢â€¢
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
      return {primaryLabel:'Tarayici Embed Kodu',primary:'<iframe src="'+urls.play+'" width="'+w+'" height="'+h+'" frameborder="0" allow="autoplay;fullscreen" allowfullscreen></iframe>',directLabel:'Player URL',direct:urls.play};
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
  const defaultFrame=(previewURLs.embed||('/embed/'+key))+(String(previewURLs.embed||'').indexOf('?')===-1?'?':'&')+'autoplay='+(autoplay?'1':'0')+'&muted='+(muted?'1':'0');
  const formatFrame=function(fmt){
    const base=previewURLs.embed||('/embed/'+key);
    return base+(String(base).indexOf('?')===-1?'?':'&')+'format='+fmt+'&autoplay='+(autoplay?'1':'0')+'&muted='+(muted?'1':'0');
  };

  switch(embedTab){
    case 'iframe':
      setPreviewFrame(prev,defaultFrame);
      break;
    case 'player':
      setPreviewFrame(prev,previewURLs.play||('/play/'+key));
      break;
    case 'hls':
    case 'jsapi':{
      setPreviewFrame(prev,previewURLs.play||('/play/'+key));
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

// â•â•â• ANALYTICS â•â•â•
async function renderAnalytics(c){
  const [data,history]=await Promise.all([api('/api/analytics'),api('/api/analytics/history')]);
  if(!data){c.innerHTML='<div class="empty-state"><div class="icon"><i class="bi bi-bar-chart-line"></i></div><h3>Analitik verisi yok</h3></div>';return}
  const fmtItems=Object.entries(data.viewers_by_format||{}).sort((a,b)=>b[1]-a[1]).map(([label,value])=>({label:label,value:value}));
  const countryItems=Object.entries(data.viewers_by_country||{}).sort((a,b)=>b[1]-a[1]).slice(0,8).map(([label,value])=>({label:label,value:value}));
  const persistedTimeline=(Array.isArray(history)?history:[]).slice().reverse().map(function(item){return {timestamp:item.timestamp,value:item.current_viewers||0}});
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Analitik</h1></div>'+
    '<div class="card-grid card-grid-4" style="margin-bottom:24px">'+
      statCard('purple','bi-collection-play',fmtInt(data.total_streams||0),'Toplam Yayin','streams','Olusturulan tum streamler')+
      statCard('green','bi-people-fill',fmtInt(data.current_viewers||0),'Aktif Izleyici','viewers','Su an acik oturumlar')+
      statCard('orange','bi-graph-up-arrow',fmtInt(data.peak_concurrent||0),'Tepe Esz.','viewers','Kaydedilen en yuksek eszamanli izleyici')+
      statCard('blue','bi-diagram-3',fmtBytes(data.total_bandwidth||0),'Toplam Bant','transcode-jobs','Sunucudan cikan toplam trafik')+
    '</div>'+
    '<div class="insight-grid">'+
      '<div class="card"><div class="card-header"><h3 class="card-title">24 Saat Izleyici Trendi</h3></div><div class="card-body">'+renderTimelineChart(data.viewers_timeline||[],'Henuz timeline verisi yok',function(v){return String(v)})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Kalici Snapshot Trendi</h3></div><div class="card-body">'+renderTimelineChart(persistedTimeline,'Henuz kalici snapshot yok',function(v){return String(v)})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Format Dagilimi</h3></div><div class="card-body">'+renderBarList(fmtItems,'Henuz format verisi yok',function(v){return String(v)})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">Ulke Dagilimi</h3></div><div class="card-body">'+renderBarList(countryItems,'Henuz ulke verisi yok',function(v){return String(v)})+'</div></div>'+
      '<div class="card"><div class="card-header"><h3 class="card-title">En Populer Yayinlar</h3></div><div class="card-body">'+
        ((data.top_streams||[]).length?(data.top_streams||[]).map(function(s){
          return '<div class="metric-row"><div><div class="setting-label">'+escHtml(s.stream_name||shortKey(s.stream_key))+'</div><div class="setting-desc"><code>'+escHtml(s.stream_key)+'</code></div></div><span class="badge">'+fmtInt(s.viewers||0)+' izleyici</span></div>';
        }).join(''):'<div style="color:var(--text-muted)">Aktif yayin yok</div>')+
      '</div></div>'+
    '</div>';
  schedulePageRefresh('analytics',5000);
}
function fmtBytes(b){if(!b||b===0)return '0 B';const k=1024,s=['B','KB','MB','GB','TB'];const i=Math.floor(Math.log(b)/Math.log(k));return (b/Math.pow(k,i)).toFixed(1)+' '+s[i]}

// â•â•â• RECORDINGS â•â•â•
async function renderRecordings(c){
  const recsRes=await api('/api/recordings');
  const streamsRes=await api('/api/streams');
  const savedRes=await api('/api/recordings/library');
  const recs=Array.isArray(recsRes)?recsRes:[];
  const streams=Array.isArray(streamsRes)?streamsRes:[];
  const saved=Array.isArray(savedRes)?savedRes:[];
  c.innerHTML=
    '<div class="page-header"><h1 class="page-title">Kayitlar</h1>'+
    '<button class="btn btn-primary" onclick="showRecordModal()">Kayit Baslat</button></div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-title" style="margin-bottom:8px">Depolama Notu</div>'+
      '<div class="form-hint">Kalici kayitlar <code>data/recordings</code> altindadir. <code>data/hls</code> ve <code>data/transcode/hls</code> dizinleri canli yayin cache alanidir; kayit olarak listelenmez.</div>'+
    '</div>'+
    '<div class="card" style="margin-bottom:16px"><div class="card-header"><h3 class="card-title">Aktif Kayitlar</h3></div>'+
    '<div class="card-body"><table class="table"><thead><tr><th>ID</th><th>Yayin</th><th>Format</th><th>Durum</th><th>Boyut</th><th>Islem</th></tr></thead><tbody id="rec-list"></tbody></table></div></div>'+
    '<div class="card"><div class="card-header"><h3 class="card-title">Kayit Kutuphanesi</h3></div>'+
    '<div class="card-body"><table class="table"><thead><tr><th>Yayin</th><th>Dosya</th><th>Format</th><th>Tarih</th><th>Boyut</th><th>Islem</th></tr></thead><tbody id="saved-rec-list"></tbody></table></div></div>'+
    '<div id="rec-modal" style="display:none"></div><div id="rec-preview-modal" style="display:none"></div>';
  const rl=document.getElementById('rec-list');
  if(rl){
    rl.innerHTML=recs.length?recs.map(r=>'<tr><td style="font-size:12px">'+r.ID+'</td><td>'+r.StreamKey+'</td><td>'+r.Format+'</td><td><span class="badge badge-'+(r.Status==='recording'?'green':'gray')+'">'+r.Status+'</span></td><td>'+fmtBytes(r.Size||0)+'</td><td>'+(r.Status==='recording'?'<button class="btn btn-sm btn-danger" onclick="stopRec(\''+r.ID+'\')">Durdur</button>':'\u2014')+'</td></tr>').join(''):'<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Aktif kayit yok</td></tr>';
  }
  const srl=document.getElementById('saved-rec-list');
  if(srl){
    srl.innerHTML=saved.length?saved.map(function(r){
      return '<tr>'+
        '<td><code>'+escHtml(r.stream_key)+'</code></td>'+
        '<td>'+escHtml(r.name)+'</td>'+
        '<td>'+(r.format||'-').toUpperCase()+'</td>'+
        '<td>'+(r.mod_time?new Date(r.mod_time).toLocaleString('tr-TR'):'-')+'</td>'+
        '<td>'+fmtBytes(r.size||0)+'</td>'+
        '<td style="display:flex;gap:8px;flex-wrap:wrap">'+
          '<button class="btn btn-sm btn-secondary" onclick=\'previewRecording('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Onizle</button>'+
          '<button class="btn btn-sm btn-secondary" onclick=\'downloadRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Indir</button>'+
          '<button class="btn btn-sm btn-danger" onclick=\'deleteRecordingFile('+JSON.stringify(r.stream_key)+','+JSON.stringify(r.name)+')\'>Sil</button>'+
        '</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="6" style="text-align:center;color:var(--text-muted);padding:24px">Kaydedilmis dosya yok</td></tr>';
  }
  window._recStreams=streams;
}
let recordingPreviewPlayer=null;
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
async function stopRec(id){await api('/api/recordings/stop/'+id);navigate('recordings')}
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
    '<div class="form-group"><label class="form-label">Format</label><select class="form-select" id="rec-fmt">'+recordingFormatOptions('ts')+'</select></div>'+
    '<button class="btn btn-primary" onclick="startNewRec()" style="width:100%">Kaydi Baslat</button>'+
    '</div></div></div>';
}
async function startNewRec(){
  const key=document.getElementById('rec-key')?.value;
  const fmt=document.getElementById('rec-fmt')?.value||'ts';
  if(!key)return;
  await api('/api/recordings',{method:'POST',body:{stream_key:key,format:fmt}});
  document.getElementById('rec-modal').style.display='none';
  navigate('recordings');
}
async function previewRecording(streamKey,name){
  destroyRecordingPreviewPlayer();
  const modal=document.getElementById('rec-preview-modal');
  if(!modal)return;
  const url=recordingFileURL(streamKey,name,false);
  const ext=(name.split('.').pop()||'').toLowerCase();
  modal.style.display='block';
  modal.innerHTML='<div style="position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,.55);z-index:1000;display:flex;align-items:center;justify-content:center;padding:20px" onclick="if(event.target===this)closeRecordingPreview()">'+
    '<div class="card" style="width:min(960px,96vw);max-height:92vh;overflow:auto"><div class="card-header"><h3 class="card-title">'+escHtml(name)+'</h3><button class="btn btn-sm btn-secondary" onclick="closeRecordingPreview()">Kapat</button></div>'+
    '<div id="rec-preview-body"></div></div></div>';
  const body=document.getElementById('rec-preview-body');
  if(!body)return;
  if(ext==='flv'||ext==='ts'){
    body.innerHTML='<video id="rec-preview-video" controls playsinline style="width:100%;max-height:70vh;background:#000"></video>';
    try{
      await loadEmbedScript('/static/vendor/mpegts.min.js');
      if(window.mpegts&&window.mpegts.isSupported&&window.mpegts.isSupported()){
        recordingPreviewPlayer=window.mpegts.createPlayer({type:ext==='flv'?'flv':'mpegts',isLive:false,url:url});
        recordingPreviewPlayer.attachMediaElement(document.getElementById('rec-preview-video'));
        recordingPreviewPlayer.load();
      }else{
        body.innerHTML='<div class="empty-state"><h3>Onizleme desteklenmiyor</h3><p style="color:var(--text-muted)">Bu tarayici kaydi dogrudan oynatamiyor.</p></div>';
      }
    }catch(e){
      body.innerHTML='<div class="empty-state"><h3>Onizleme hazirlanamadi</h3><p style="color:var(--text-muted)">'+escHtml(e.message||'Bilinmeyen hata')+'</p></div>';
    }
    return;
  }
  if(ext==='mp4'||ext==='webm'||ext==='ogg'){
    body.innerHTML='<video controls playsinline src="'+url+'" style="width:100%;max-height:70vh;background:#000"></video>';
    return;
  }
  if(ext==='mp3'||ext==='aac'||ext==='wav'||ext==='flac'){
    body.innerHTML='<div style="padding:24px"><audio controls src="'+url+'" style="width:100%"></audio></div>';
    return;
  }
  body.innerHTML='<div class="empty-state"><h3>Onizleme yok</h3><p style="color:var(--text-muted)">Bu format panelde dogrudan oynatilamiyor. Dosyayi indirebilirsiniz.</p></div>';
}
function closeRecordingPreview(){
  destroyRecordingPreviewPlayer();
  const modal=document.getElementById('rec-preview-modal');
  if(!modal)return;
  modal.style.display='none';
  modal.innerHTML='';
}
function downloadRecordingFile(streamKey,name){
  window.open(recordingFileURL(streamKey,name,true),'_blank');
}
async function deleteRecordingFile(streamKey,name){
  if(!confirm('Bu kayit dosyasini silmek istediginize emin misiniz?'))return;
  const res=await api('/api/recordings/file',{method:'DELETE',body:{stream_key:streamKey,filename:name}});
  if(res&&res.status==='deleted'){
    toast('Kayit silindi');
    navigate('recordings');
  }else{
    toast((res&&res.message)||'Kayit silinemedi','error');
  }
}

// â•â•â• VIEWERS â•â•â•
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
        '<td>'+(sess.last_seen?new Date(sess.last_seen).toLocaleTimeString('tr-TR'):'-')+'</td>'+
      '</tr>';
    }).join(''):'<tr><td colspan="7" style="text-align:center;color:var(--text-muted);padding:24px">Aktif izleyici oturumu yok</td></tr>';
  }
  const bans=await api('/api/security/bans');
  const bl=document.getElementById('ban-list');
  if(bl&&bans){
    bl.innerHTML=bans.length?bans.map(b=>'<tr><td><code>'+b.IP+'</code></td><td>'+b.Reason+'</td><td>'+(b.BannedAt?new Date(b.BannedAt).toLocaleString('tr-TR'):'\u2014')+'</td><td><button onclick="unbanIP(\''+b.IP+'\')" style="background:#e74c3c;color:#fff;padding:4px 12px;border:none;border-radius:6px;cursor:pointer;font-size:12px">Kaldir</button></td></tr>').join(''):'<tr><td colspan="4" style="text-align:center;color:var(--text-muted);padding:24px">Yasakli IP yok</td></tr>';
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
      copyField('Gecerlilik',res.expires_at?new Date(res.expires_at).toLocaleString('tr-TR'):'-');
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
      const started=job.started_at?new Date(job.started_at).toLocaleString('tr-TR'):'-';
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

init();
</script>
</body>
</html>` + "`"
