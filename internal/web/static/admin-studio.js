(function(){
  if(!window.api || !window.escHtml) return;

  const legacyStudioRenders = {
    dashboard: typeof window.renderDashboard === 'function' ? window.renderDashboard : null,
    streams: typeof window.renderStreams === 'function' ? window.renderStreams : null,
    guidedSettings: typeof window.renderGuidedSettings === 'function' ? window.renderGuidedSettings : null,
    settingsGeneral: typeof window.renderSettingsGeneral === 'function' ? window.renderSettingsGeneral : null,
    settingsEmbed: typeof window.renderSettingsEmbed === 'function' ? window.renderSettingsEmbed : null,
    settingsProtocols: typeof window.renderSettingsProtocols === 'function' ? window.renderSettingsProtocols : null,
    settingsOutputs: typeof window.renderSettingsOutputs === 'function' ? window.renderSettingsOutputs : null,
    settingsSecurity: typeof window.renderSettingsSecurity === 'function' ? window.renderSettingsSecurity : null,
    diagnostics: typeof window.renderDiagnostics === 'function' ? window.renderDiagnostics : null,
    maintenanceCenter: typeof window.renderMaintenanceCenter === 'function' ? window.renderMaintenanceCenter : null,
    securityTokens: typeof window.renderSecurityTokens === 'function' ? window.renderSecurityTokens : null,
    playerTemplates: typeof window.renderPlayerTemplates === 'function' ? window.renderPlayerTemplates : null,
    advancedEmbed: typeof window.renderAdvancedEmbed === 'function' ? window.renderAdvancedEmbed : null,
    analytics: typeof window.renderAnalytics === 'function' ? window.renderAnalytics : null,
    settingsTranscode: typeof window.renderSettingsTranscode === 'function' ? window.renderSettingsTranscode : null,
    viewers: typeof window.renderViewers === 'function' ? window.renderViewers : null,
    transcodeJobs: typeof window.renderTranscodeJobs === 'function' ? window.renderTranscodeJobs : null
  };

  const EMBED_USE_CASES = [
    {key:'website', title:'Web Sitesi', desc:'Genel web yayini ve kurumsal sayfa icin dengeli secim.', output:'iframe', format:'player', badge:'Genel'},
    {key:'newsroom', title:'Haber Portali', desc:'Haber siteleri ve hizli gomulu player kullanimlari.', output:'script', format:'player', badge:'Editoryal'},
    {key:'corporate', title:'Kurumsal Sayfa', desc:'Daha temiz ve markali embed deneyimi.', output:'iframe', format:'hls', badge:'Markali'},
    {key:'mobile', title:'Mobil Uyumlu', desc:'Dusuk bant ve mobil oncelikli dagitim.', output:'iframe', format:'hls', badge:'Mobil'},
    {key:'audio', title:'Sadece Ses', desc:'Radyo, podcast ve audio-only oynatim.', output:'audio', format:'dash_audio', badge:'Audio'},
    {key:'private', title:'Gizli Yayin', desc:'Paylasimi kisitli, korumali oynatim senaryosu.', output:'player', format:'player', badge:'Guvenli'},
    {key:'token', title:'Token Korumali', desc:'Imzali veya sureli playback baglantisi uretir.', output:'iframe', format:'player', badge:'Token'},
    {key:'low_latency', title:'Dusuk Gecikme', desc:'Canli olay ve hizli geri bildirim icin.', output:'iframe', format:'ll_hls', badge:'Canli'},
    {key:'dash', title:'DASH', desc:'MPEG-DASH oncelikli teslimat ve istemci testi.', output:'manifest', format:'dash', badge:'DASH'},
    {key:'hls', title:'HLS', desc:'Tarayici ve mobil oynatici uyumlulugu icin HLS.', output:'manifest', format:'hls', badge:'HLS'},
    {key:'mp4_fallback', title:'MP4 Fallback', desc:'Eski istemciler ve dogrudan dosya linki icin.', output:'player', format:'mp4', badge:'Fallback'}
  ];

  const EMBED_OUTPUTS = [
    {key:'iframe', title:'Iframe Embed', desc:'Hazir sitelere hizli gomme.'},
    {key:'script', title:'Script Embed', desc:'JS ile kolay yerlestirme.'},
    {key:'player', title:'Player URL', desc:'Direkt player sayfasi.'},
    {key:'audio', title:'Audio Player', desc:'Sadece ses yayini icin oynatici.'},
    {key:'popup', title:'Popup Player', desc:'Tiklaninca ayri pencere acar.'},
    {key:'manifest', title:'Direct Manifest', desc:'HLS, DASH veya audio manifest linki.'},
    {key:'vlc', title:'VLC Linki', desc:'Harici player ile test ve izleme.'}
  ];

  const EMBED_FORMAT_OPTIONS = [
    ['player','Akilli Player'],['hls','HLS'],['dash','DASH'],['ll_hls','Low Latency HLS'],['mp4','MP4 Fallback'],
    ['hls_audio','HLS Ses'],['dash_audio','DASH Ses'],['mp3','MP3'],['aac','AAC'],['icecast','Icecast']
  ];

  const EMBED_THEME_OPTIONS = [
    ['clean','Clean'],['newsroom','Newsroom'],['corporate','Corporate'],['radio','Radio'],['minimal','Minimal']
  ];

  function parseJSONSafeStudio(raw,fallback){
    try{
      const parsed=JSON.parse(raw||'');
      return parsed==null?fallback:parsed;
    }catch(e){
      return fallback;
    }
  }

  function cloneJSON(value){
    return JSON.parse(JSON.stringify(value==null?{}:value));
  }

  async function studioUploadFile(category,file){
    if(!file) return {error:true,message:'Dosya secilmedi'};
    const form=new FormData();
    form.append('category',category||'branding');
    form.append('file',file,file.name||'asset');
    const headers={};
    if(typeof authToken!=='undefined' && authToken) headers.Authorization='Bearer '+authToken;
    try{
      const res=await fetch(API+'/api/admin/assets',{method:'POST',headers:headers,body:form,cache:'no-store'});
      return await res.json();
    }catch(err){
      return {error:true,message:(err&&err.message)||'Yukleme hatasi'};
    }
  }

  async function studioListAssets(category){
    const result=await api('/api/admin/assets?category='+encodeURIComponent(category||'branding'));
    return Array.isArray(result&&result.items)?result.items:[];
  }

  async function studioDeleteAsset(path){
    return api('/api/admin/assets',{method:'DELETE',body:{path:path}});
  }

  function toNumber(value,fallback){
    const n=Number(value);
    return Number.isFinite(n)?n:fallback;
  }

  function fmtStudioNumber(value,digits){
    const n=Number(value||0);
    if(!Number.isFinite(n)) return '-';
    return typeof digits==='number' ? n.toFixed(digits) : n.toLocaleString(localeForLang());
  }

  function studioRerender(page){
    if(typeof loadPage==='function') return loadPage(page);
  }

  function studioAuditControls(root){
    const scope=root||document;
    scope.querySelectorAll('textarea, .form-textarea').forEach(function(el){
      el.classList.add('studio-textarea');
      if(!el.style.minHeight) el.style.minHeight=el.rows && Number(el.rows)>=4 ? '160px' : '120px';
    });
    scope.querySelectorAll('input.form-input, select.form-select, input.input, select.input, textarea.input').forEach(function(el){
      el.classList.add('studio-control');
    });
    scope.querySelectorAll('.page-header').forEach(function(header){
      if(header.classList.contains('studio-headerized')) return;
      header.classList.add('studio-headerized');
    });
    scope.querySelectorAll('.copy-input, .mono-wrap, code').forEach(function(el){
      el.classList.add('studio-mono');
    });
  }

  function studioWrapLegacyPage(container,title,subtitle,pills,actionsHTML){
    if(!container) return;
    const temp=document.createElement('div');
    temp.innerHTML=container.innerHTML;
    temp.querySelectorAll(':scope > .page-header').forEach(function(node){ node.remove(); });
    const body=temp.innerHTML;
    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero">'+
          '<div style="display:flex;justify-content:space-between;align-items:flex-start;gap:14px;flex-wrap:wrap">'+
            '<div><h1 class="studio-hero-title">'+escHtml(title||'Studio')+'</h1><div class="studio-hero-sub">'+escHtml(subtitle||'')+'</div></div>'+
            (actionsHTML?'<div class="studio-toolbar-group">'+actionsHTML+'</div>':'')+
          '</div>'+
          (Array.isArray(pills)&&pills.length?'<div class="studio-pill-row" style="margin-top:14px">'+pills.map(function(item){ return '<span class="studio-pill'+(item&&item.active?' active':'')+'">'+escHtml(String((item&&item.label)||item||''))+'</span>'; }).join('')+'</div>':'')+
        '</section>'+
        body+
      '</div>';
    studioAuditControls(container);
  }

  async function studioRenderLegacy(container,key,hero){
    const fn=legacyStudioRenders[key];
    if(typeof fn!=='function') return false;
    await fn(container);
    if(hero){
      studioWrapLegacyPage(container,hero.title,hero.subtitle,hero.pills,hero.actionsHTML);
    }else{
      studioAuditControls(container);
    }
    return true;
  }

  function studioField(label,control,hint){
    return '<div class="studio-field"><label class="form-label">'+escHtml(label)+'</label>'+control+(hint?'<div class="form-hint">'+escHtml(hint)+'</div>':'')+'</div>';
  }

  function studioSelectOptions(items,selected,emptyLabel){
    const list=Array.isArray(items)?items:[];
    const prefix=emptyLabel?'<option value="">'+escHtml(emptyLabel)+'</option>':'';
    return prefix+list.map(function(item){
      const value=Array.isArray(item)?item[0]:item.value;
      const label=Array.isArray(item)?item[1]:item.label;
      return '<option value="'+escHtml(String(value))+'"'+(String(value)===String(selected)?' selected':'')+'>'+escHtml(String(label))+'</option>';
    }).join('');
  }

  function studioOptionCard(item,active){
    return '<button type="button" class="studio-option-card'+(active?' active':'')+'" data-studio-key="'+escHtml(item.key)+'">'+
      '<div style="display:flex;align-items:center;justify-content:space-between;gap:8px;margin-bottom:8px"><h4 class="studio-option-title">'+escHtml(item.title)+'</h4>'+(item.badge?'<span class="studio-chip'+(active?' active':'')+'">'+escHtml(item.badge)+'</span>':'')+'</div>'+
      '<div class="studio-option-meta">'+escHtml(item.desc||'')+'</div>'+
    '</button>';
  }

  function studioProfileSelectOptions(items,selected){
    const list=Array.isArray(items)?items:[];
    return '<option value="">Kayitli profil sec</option>'+list.map(function(item){
      const scope=item.scope==='stream'?'Bu stream':'Global';
      return '<option value="'+Number(item.id)+'"'+(Number(item.id)===Number(selected)?' selected':'')+'>'+escHtml(item.name)+' • '+escHtml(scope)+'</option>';
    }).join('');
  }

  function defaultEmbedState(){
    return {
      mode:'simple',
      profileId:0,
      profileName:'',
      streamKey:'',
      useCase:'website',
      outputType:'iframe',
      primaryFormat:'player',
      width:1280,
      height:720,
      theme:'clean',
      templateId:0,
      notes:'',
      options:{
        responsive:true,
        autoplay:false,
        muted:false,
        debug:false,
        audioOnly:false,
        startQuality:'auto',
        referrerPolicy:'strict-origin-when-cross-origin',
        posterURL:'',
        sharePackage:'general'
      },
      branding:{
        watermarkText:'',
        logoURL:'',
        posterURL:'',
        brandColor:'#2563eb'
      },
      security:{
        signedURL:false,
        tokenRequired:false,
        sessionBound:false,
        applyStreamPolicy:false,
        expiryMinutes:60,
        domainRestriction:'',
        ipRestriction:'',
        viewerID:'',
        watermark:''
      }
    };
  }

  function ensureEmbedState(){
    if(!window.embedStudioState) window.embedStudioState=defaultEmbedState();
    const state=window.embedStudioState;
    if(!state.options) state.options=defaultEmbedState().options;
    if(!state.branding) state.branding=defaultEmbedState().branding;
    if(!state.security) state.security=defaultEmbedState().security;
    return state;
  }

  function studioPresetForUseCase(key){
    const preset=cloneJSON(defaultEmbedState());
    switch(String(key||'')){
      case 'newsroom':
        preset.outputType='script';
        preset.theme='newsroom';
        preset.options.sharePackage='newsroom';
        break;
      case 'corporate':
        preset.outputType='iframe';
        preset.primaryFormat='hls';
        preset.theme='corporate';
        break;
      case 'mobile':
        preset.outputType='iframe';
        preset.primaryFormat='hls';
        preset.width=960;
        preset.height=540;
        preset.options.startQuality='480p';
        break;
      case 'audio':
        preset.outputType='audio';
        preset.primaryFormat='dash_audio';
        preset.height=180;
        preset.options.audioOnly=true;
        preset.theme='radio';
        preset.options.sharePackage='audio';
        break;
      case 'private':
        preset.outputType='player';
        preset.security.tokenRequired=true;
        preset.security.signedURL=true;
        preset.security.applyStreamPolicy=false;
        preset.options.sharePackage='private';
        break;
      case 'token':
        preset.outputType='iframe';
        preset.security.tokenRequired=true;
        preset.security.signedURL=true;
        preset.security.applyStreamPolicy=false;
        break;
      case 'low_latency':
        preset.outputType='iframe';
        preset.primaryFormat='ll_hls';
        preset.options.autoplay=true;
        preset.options.muted=true;
        break;
      case 'dash':
        preset.outputType='manifest';
        preset.primaryFormat='dash';
        break;
      case 'hls':
        preset.outputType='manifest';
        preset.primaryFormat='hls';
        break;
      case 'mp4_fallback':
        preset.outputType='player';
        preset.primaryFormat='mp4';
        break;
      default:
        preset.outputType='iframe';
        preset.primaryFormat='player';
    }
    return preset;
  }

  function studioMergeEmbedState(target,patch){
    Object.assign(target,patch||{});
    if(patch && patch.options) target.options=Object.assign({},target.options||{},patch.options);
    if(patch && patch.branding) target.branding=Object.assign({},target.branding||{},patch.branding);
    if(patch && patch.security) target.security=Object.assign({},target.security||{},patch.security);
    return target;
  }

  function studioStreamSecurityState(stream){
    const policy=parseJSONSafeStudio((stream&&stream.policy_json)||'{}',{});
    return {
      requireToken:!!policy.require_playback_token,
      requireSignedURL:!!policy.require_signed_url,
      domainLock:String((stream&&stream.domain_lock)||'').trim(),
      ipWhitelist:String((stream&&stream.ip_whitelist)||'').trim(),
      active:!!(policy.require_playback_token || policy.require_signed_url || String((stream&&stream.domain_lock)||'').trim() || String((stream&&stream.ip_whitelist)||'').trim())
    };
  }

  function studioBuildScriptEmbed(url,width,height){
    const w=Math.max(320,parseInt(width||1280,10)||1280);
    const h=Math.max(180,parseInt(height||720,10)||720);
    return [
      '<div id="fluxstream-embed-root"></div>',
      '<script>',
      '(function(){',
      'var root=document.getElementById("fluxstream-embed-root");',
      'if(!root)return;',
      'root.style.position="relative";',
      'root.style.width="100%";',
      'root.style.maxWidth="'+w+'px";',
      'root.style.aspectRatio="'+w+' / '+h+'";',
      'var iframe=document.createElement("iframe");',
      'iframe.src='+JSON.stringify(url)+';',
      'iframe.width="'+w+'";',
      'iframe.height="'+h+'";',
      'iframe.allow="autoplay; fullscreen";',
      'iframe.allowFullscreen=true;',
      'iframe.style.border="0";',
      'iframe.style.width="100%";',
      'iframe.style.height="100%";',
      'root.appendChild(iframe);',
      '})();',
      '</script>'
    ].join('\n');
  }

  function studioBuildPopupCode(url,width,height){
    const w=Math.max(640,parseInt(width||1280,10)||1280);
    const h=Math.max(360,parseInt(height||720,10)||720);
    return '<button type="button" onclick="window.open('+JSON.stringify(url)+',\'fluxstream_popup\',\'width='+w+',height='+h+',noopener=yes,resizable=yes\')">FluxStream Player Ac</button>';
  }

  function studioApplyTemplateToURLs(urls,template,streamName){
    if(!template || typeof appendTemplateQuery!=='function') return urls;
    const next=Object.assign({},urls||{});
    Object.keys(next).forEach(function(key){
      if(next[key]) next[key]=appendTemplateQuery(next[key],template,streamName);
    });
    return next;
  }

  function studioSelectedTemplate(templates,id){
    const list=Array.isArray(templates)?templates:[];
    return list.find(function(item){ return Number(item.id)===Number(id); }) || null;
  }

  async function studioFetchPlaybackBundle(state,stream,template){
    const streamSecurity=studioStreamSecurityState(stream);
    const payload={
      stream_key:stream.stream_key,
      page:(state.outputType==='player'||state.outputType==='popup'||state.outputType==='audio')?'player':'embed',
      format:state.outputType==='audio'?(state.primaryFormat||'dash_audio'):(state.primaryFormat||'player'),
      width:state.width,
      height:state.height,
      autoplay:!!state.options.autoplay,
      muted:!!state.options.muted,
      options:cloneJSON(state.options),
      security:{
        signed_url:!!(state.security.signedURL || streamSecurity.requireSignedURL),
        token_required:!!(state.security.tokenRequired || streamSecurity.requireToken),
        session_bound:!!state.security.sessionBound,
        expiry_minutes:toNumber(state.security.expiryMinutes,60),
        domain_restriction:state.security.domainRestriction||streamSecurity.domainLock||'',
        ip_restriction:state.security.ipRestriction||streamSecurity.ipWhitelist||'',
        viewer_id:state.security.viewerID||'',
        watermark:(state.security.watermark||state.branding.watermarkText||'').trim()
      },
      apply_stream_policy:!!state.security.applyStreamPolicy
    };
    const result=await api('/api/admin/security/playback-link',{method:'POST',body:payload});
    if(result && !result.error){
      const urls=studioApplyTemplateToURLs({
        player_url:result.player_url,
        embed_url:result.embed_url,
        manifest_url:result.manifest_url,
        audio_url:result.audio_url,
        vlc_url:result.vlc_url
      },template,stream.name);
      result.player_url=urls.player_url;
      result.embed_url=urls.embed_url;
      result.manifest_url=urls.manifest_url;
      result.audio_url=urls.audio_url;
      result.vlc_url=urls.vlc_url;
    }
    return result||{error:true,message:'Playback baglantisi olusturulamadi'};
  }

  function studioBuildCodeBundle(result,state){
    const output=state.outputType;
    const width=state.width||1280;
    const height=state.height||720;
    const playerURL=result.player_url||'';
    const embedURL=result.embed_url||playerURL;
    const manifestURL=result.manifest_url||'';
    const audioURL=result.audio_url||'';
    const vlcURL=result.vlc_url||manifestURL;
    let code='', label='Embed Kodu', primaryURL=embedURL;
    if(output==='iframe'){
      code='<iframe src="'+embedURL+'" width="'+width+'" height="'+height+'" frameborder="0" allow="autoplay; fullscreen" allowfullscreen></iframe>';
    }else if(output==='script'){
      label='Script Embed';
      code=studioBuildScriptEmbed(embedURL,width,height);
    }else if(output==='player'){
      label='Player URL';
      code=playerURL;
      primaryURL=playerURL;
    }else if(output==='audio'){
      label='Audio Player URL';
      code=audioURL||playerURL;
      primaryURL=playerURL;
    }else if(output==='popup'){
      label='Popup Kodu';
      code=studioBuildPopupCode(playerURL,width,height);
      primaryURL=playerURL;
    }else if(output==='manifest'){
      label='Manifest URL';
      code=manifestURL||audioURL||playerURL;
      primaryURL=manifestURL||audioURL||playerURL;
    }else if(output==='vlc'){
      label='VLC URL';
      code=vlcURL;
      primaryURL=vlcURL;
    }
    return {label:label,code:code,primaryURL:primaryURL,playerURL:playerURL,embedURL:embedURL,manifestURL:manifestURL,audioURL:audioURL,vlcURL:vlcURL};
  }

  function studioRenderWarnings(state,bundle){
    const warnings=[];
    if(state.outputType==='audio' && !String(state.primaryFormat||'').includes('audio')) warnings.push('Audio player secili ama format audio-only degil. HLS Ses veya DASH Ses tercih et.');
    if((state.security.domainRestriction||state.security.ipRestriction) && !state.security.signedURL && !state.security.tokenRequired) warnings.push('Domain veya IP kisiti icin sureli token veya signed URL secmek gerekir.');
    if(state.security.signedURL && !state.security.applyStreamPolicy) warnings.push('Signed URL olusturuldu. Kalici enforcement icin stream policy uygula secenegini de ac.');
    if((state.useCase==='private' || state.useCase==='token') && !state.security.applyStreamPolicy) warnings.push('Gizli ve token korumali presetler varsayilan olarak sadece bu ekranda gecici link uretir. Kalici hale gelmesi icin Stream policy uygula secenegini bilerek acman gerekir.');
    if(bundle && !bundle.primaryURL) warnings.push('Uretilecek cikti henuz hazir degil.');
    if(!warnings.length) return '';
    return '<div class="studio-alert warning"><strong>Dikkat edilmesi gerekenler</strong><ul style="margin:10px 0 0 18px;padding:0">'+warnings.map(function(item){return '<li>'+escHtml(item)+'</li>';}).join('')+'</ul></div>';
  }

  async function renderEmbedStudio(container){
    const state=ensureEmbedState();
    const [streams,templates,profiles]=await Promise.all([
      api('/api/streams'),
      api('/api/players'),
      api('/api/admin/embed-profiles'+(state.streamKey?('?stream_key='+encodeURIComponent(state.streamKey)):''))
    ]);
    const streamList=Array.isArray(streams)?streams:[];
    if(!state.streamKey && streamList[0]) state.streamKey=streamList[0].stream_key;
    const stream=streamList.find(function(item){return item.stream_key===state.streamKey;}) || streamList[0] || null;
    if(stream && !state.streamKey) state.streamKey=stream.stream_key;
    const streamSecurity=studioStreamSecurityState(stream);
    const template=studioSelectedTemplate(Array.isArray(templates)?templates:[],state.templateId);
    const playbackResult=stream?await studioFetchPlaybackBundle(state,stream,template):null;
    const codeBundle=studioBuildCodeBundle(playbackResult||{},state);
    const selectedUseCase=EMBED_USE_CASES.find(function(item){return item.key===state.useCase;}) || EMBED_USE_CASES[0];
    const previewURL=(state.outputType==='player'||state.outputType==='audio'?codeBundle.playerURL:codeBundle.embedURL) || '';
    const profileItems=Array.isArray(profiles)?profiles:[];

    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero">'+
          '<h1 class="studio-hero-title">Embed Studyosu</h1>'+
          '<div class="studio-hero-sub">Kullanim tipine gore embed kodu, markali player, canli onizleme, guvenli paylasim baglantisi ve kaydedilebilir embed profilleri artik tek merkezde.</div>'+
          '<div class="studio-pill-row" style="margin-top:14px">'+
            '<span class="studio-pill active">'+escHtml((selectedUseCase&&selectedUseCase.title)||'Kullanim tipi')+'</span>'+
            '<span class="studio-pill">'+escHtml((stream&&stream.name)||'Stream sec')+'</span>'+
            '<span class="studio-pill">'+escHtml(codeBundle.label||'Cikti')+'</span>'+
            '<span class="studio-pill">'+escHtml((state.primaryFormat||'player').toUpperCase())+'</span>'+
          '</div>'+
        '</section>'+
        '<section class="studio-toolbar">'+
          '<div class="studio-toolbar-group">'+
            '<select id="studio-embed-stream" class="input">'+streamList.map(function(item){return '<option value="'+escHtml(item.stream_key)+'"'+(item.stream_key===state.streamKey?' selected':'')+'>'+escHtml(item.name)+' • '+escHtml(item.stream_key)+'</option>';}).join('')+'</select>'+
            '<select id="studio-embed-profile" class="input">'+studioProfileSelectOptions(profileItems,state.profileId)+'</select>'+
            '<button class="btn btn-secondary" id="studio-embed-load-profile">Profili Yukle</button>'+
          '</div>'+
          '<div class="studio-toolbar-group">'+
            '<div class="segmented"><button class="segment'+(state.mode==='simple'?' active':'')+'" data-mode="simple" type="button">Basit Mod</button><button class="segment'+(state.mode==='advanced'?' active':'')+'" data-mode="advanced" type="button">Gelismis Mod</button></div>'+
            '<button class="btn btn-secondary" id="studio-embed-new-profile">Temizle</button>'+
            '<button class="btn btn-primary" id="studio-embed-save-profile">Profili Kaydet</button>'+
          '</div>'+
        '</section>'+
        '<div class="studio-grid studio-grid-2">'+
          '<div class="studio-card soft"><div><h2 class="studio-section-title">Hazir kullanim tipleri</h2><div class="studio-section-sub">Hazir paket secildiginde format, embed tipi ve guvenlik secimleri buna gore onerilir.</div></div><div class="studio-option-grid" id="studio-embed-usecases">'+EMBED_USE_CASES.map(function(item){return studioOptionCard(item,item.key===state.useCase);}).join('')+'</div><div><h2 class="studio-section-title">Cikti tipi</h2><div class="studio-section-sub">Tek textarea yerine secilebilir, onizlenebilir ve test edilebilir cikti kartlari.</div></div><div class="studio-option-grid" id="studio-embed-outputs">'+EMBED_OUTPUTS.map(function(item){return studioOptionCard(item,item.key===state.outputType);}).join('')+'</div></div>'+
          '<div class="studio-card"><div><h2 class="studio-section-title">Canli onizleme ve paylasim</h2><div class="studio-section-sub">Secili stream ve cikti turune gore uretilen player burada gorunur.</div></div><div class="studio-preview-shell">'+(previewURL?'<iframe src="'+escHtml(withQueryParam(previewURL,'debug',state.options.debug?'1':'0'))+'" allow="autoplay; fullscreen" allowfullscreen></iframe>':'<div class="empty-state" style="padding:34px"><div class="icon"><i class="bi bi-broadcast"></i></div><h3>Onizleme icin stream secin</h3><p style="color:var(--text-muted)">Kod ve player burada gosterilecek.</p></div>')+'</div><div class="studio-chip-row"><button class="btn btn-secondary" id="studio-embed-copy-code">Kopyala</button><button class="btn btn-secondary" id="studio-embed-copy-url">Baglantiyi Kopyala</button><button class="btn btn-secondary" id="studio-embed-open">Yeni Sekmede Ac</button><button class="btn btn-secondary" id="studio-embed-debug">Debug Ile Ac</button><button class="btn btn-secondary" id="studio-embed-test">Test Et</button></div><div class="studio-card soft" style="padding:14px"><div class="studio-section-title" style="font-size:16px">Kullanim ozeti</div><div class="studio-section-sub">'+escHtml((selectedUseCase&&selectedUseCase.desc)||'')+'</div><div class="studio-chip-row" style="margin-top:10px"><span class="studio-chip active">'+escHtml(codeBundle.label||'Cikti')+'</span><span class="studio-chip">'+escHtml(state.options.responsive?'Responsive':'Sabit')+'</span><span class="studio-chip">'+escHtml(state.security.signedURL||state.security.tokenRequired?'Korumali':'Acik')+'</span></div></div></div>'+
        '</div>'+
        '<div class="studio-grid studio-grid-2">'+
          '<div class="studio-card"><div><h2 class="studio-section-title">Embed ayarlari</h2><div class="studio-section-sub">Boyut, tema, player sablonu, poster, marka ve davranis secimleri burada yonetilir.</div></div><div class="studio-grid studio-grid-2">'+
            studioField('Profil adi','<input id="studio-embed-name" class="input" value="'+escHtml(state.profileName||'')+'" placeholder="Ornek: Haber sitesi tokenli embed">','Kaydederken profil basligi olarak kullanilir.')+
            studioField('Format','<select id="studio-embed-format" class="input">'+studioSelectOptions(EMBED_FORMAT_OPTIONS,state.primaryFormat)+'</select>','HLS, DASH, audio-only ve fallback secimleri.')+
            studioField('Genislik','<input id="studio-embed-width" class="input" type="number" min="320" value="'+escHtml(String(state.width||1280))+'">','Iframe ve popup ciktisi icin.')+
            studioField('Yukseklik','<input id="studio-embed-height" class="input" type="number" min="120" value="'+escHtml(String(state.height||720))+'">','Audio paketleri daha dusuk olabilir.')+
            studioField('Tema','<select id="studio-embed-theme" class="input">'+studioSelectOptions(EMBED_THEME_OPTIONS,state.theme)+'</select>','Marka ve paket etiketleri icin.')+
            studioField('Player sablonu','<select id="studio-embed-template" class="input"><option value="0">Varsayilan player</option>'+((Array.isArray(templates)?templates:[]).map(function(item){return '<option value="'+Number(item.id)+'"'+(Number(item.id)===Number(state.templateId)?' selected':'')+'>'+escHtml(item.name)+'</option>';}).join(''))+'</select>','Kayitli player sablonu varsa player ve embed linkine eklenir.')+
            studioField('Poster URL','<input id="studio-embed-poster" class="input" value="'+escHtml(state.branding.posterURL||state.options.posterURL||'')+'" placeholder="https://ornek.com/poster.jpg">','Poster veya bos ekran resmi.')+
            studioField('Watermark yazisi','<input id="studio-embed-watermark" class="input" value="'+escHtml(state.branding.watermarkText||'')+'" placeholder="CANLI • FluxStream">','Player icinde gorunecek markalama metni.')+
          '</div><div class="studio-option-grid">'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Responsive</strong><div class="form-hint">Genislige gore otomatik uyarlanir.</div></div><input type="checkbox" id="studio-embed-responsive" '+(state.options.responsive?'checked':'')+'></div></label>'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Autoplay</strong><div class="form-hint">Uygun ortamlarda otomatik oynatir.</div></div><input type="checkbox" id="studio-embed-autoplay" '+(state.options.autoplay?'checked':'')+'></div></label>'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Muted</strong><div class="form-hint">Autoplay ile uyumluluk icin.</div></div><input type="checkbox" id="studio-embed-muted" '+(state.options.muted?'checked':'')+'></div></label>'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Debug gorunumu</strong><div class="form-hint">Test linkine QoE debug ekler.</div></div><input type="checkbox" id="studio-embed-debug-flag" '+(state.options.debug?'checked':'')+'></div></label>'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Audio-only davranis</strong><div class="form-hint">Radyo ve podcast paketleri icin.</div></div><input type="checkbox" id="studio-embed-audio-only" '+(state.options.audioOnly?'checked':'')+'></div></label>'+
            studioField('Baslangic kalitesi','<select id="studio-embed-start-quality" class="input">'+studioSelectOptions([['auto','Otomatik'],['1080p','1080p'],['720p','720p'],['480p','480p'],['360p','360p']],state.options.startQuality)+'</select>','Oncelikli kalite ipucu olarak saklanir.')+
            studioField('Referrer policy','<select id="studio-embed-referrer" class="input">'+studioSelectOptions([['strict-origin-when-cross-origin','Strict origin when cross origin'],['origin','Origin'],['same-origin','Same origin'],['no-referrer','No referrer']],state.options.referrerPolicy)+'</select>','Gomme sayfasina iliskin tavsiye edilen policy.')+
            studioField('Paylasim paketi','<select id="studio-embed-share-package" class="input">'+studioSelectOptions([['general','Genel'],['newsroom','Haber sitesi'],['corporate','Kurumsal duyuru'],['audio','Sadece ses'],['private','Gizli yayin']],state.options.sharePackage)+'</select>','Profil ile birlikte saklanir.')+
          '</div></div>'+
          '<div class="studio-card"><div><h2 class="studio-section-title">Playback guvenligi v1</h2><div class="studio-section-sub">Signed URL, sureli token, domain ve IP kisiti, iframe baglamasi ve watermark tek merkezde yonetilir.</div></div>'+
          '<div class="studio-alert info"><strong>Gizli preset nasil calisir?</strong><div style="margin-top:8px" class="form-hint">Gizli ve Token Korumali secenekleri varsayilan olarak sadece bu ekranda gecici link uretir. Tum paneli veya tum izleyicileri kilitlemez. Kalici bir politika istiyorsan ayrica <strong>Stream policy uygula</strong> secenegini bilerek acmalisin.</div></div>'+
          (streamSecurity.active?'<div class="studio-alert warning"><strong>Bu streamde kalici playback korumasi aktif</strong><div style="margin-top:8px" class="form-hint">Bu nedenle onizleme ve diger ekranlar token veya signed URL bekliyor olabilir. Eski normale donmek icin alttaki dugme ile korumayi temizleyebilirsin.</div><div style="margin-top:12px"><button class="btn btn-danger btn-sm" id="studio-sec-reset-policy">Kalici korumayi kaldir</button></div></div>':'')+
          '<div class="studio-option-grid">'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Signed URL</strong><div class="form-hint">Imzali player ve manifest baglantisi uretir.</div></div><input type="checkbox" id="studio-sec-signed" '+(state.security.signedURL?'checked':'')+'></div></label>'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Token zorunlu</strong><div class="form-hint">Playback token uretir.</div></div><input type="checkbox" id="studio-sec-token" '+(state.security.tokenRequired?'checked':'')+'></div></label>'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Oturuma bagli</strong><div class="form-hint">Viewer ID ile bagli izleme baglantisi hazirlar.</div></div><input type="checkbox" id="studio-sec-session" '+(state.security.sessionBound?'checked':'')+'></div></label>'+
            '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Stream policy uygula</strong><div class="form-hint">Ayni guvenligi stream policy tarafina da yazar.</div></div><input type="checkbox" id="studio-sec-apply" '+(state.security.applyStreamPolicy?'checked':'')+'></div></label>'+
          '</div><div class="studio-grid studio-grid-2">'+
            studioField('Gecerlilik suresi (dk)','<input id="studio-sec-expiry" class="input" type="number" min="1" value="'+escHtml(String(state.security.expiryMinutes||60))+'">','Token ve signed URL suresi.')+
            studioField('Viewer ID','<input id="studio-sec-viewer" class="input" value="'+escHtml(state.security.viewerID||'')+'" placeholder="izleyici-42">','Oturuma bagli ve watermarkli oynatim icin.')+
            studioField('Domain kisiti','<input id="studio-sec-domain" class="input" value="'+escHtml(state.security.domainRestriction||'')+'" placeholder="haber.example.com">','Tek domain veya referrer kisiti.')+
            studioField('IP kisiti','<input id="studio-sec-ip" class="input" value="'+escHtml(state.security.ipRestriction||'')+'" placeholder="203.0.113.10">','Belirli istemci veya ofis agi icin.')+
            studioField('Guvenlik watermark','<input id="studio-sec-watermark" class="input" value="'+escHtml(state.security.watermark||'')+'" placeholder="Viewer #42">','Token icine yazilan gorunur iz.')+
            studioField('Profil notu','<input id="studio-sec-notes" class="input" value="'+escHtml(state.notes||'')+'" placeholder="Bu paket gizli yayin icin">','Kaydedilen embed profiline not duser.')+
          '</div>'+(playbackResult&&playbackResult.token?'<div class="studio-alert info"><strong>Son uretilen guvenli baglanti</strong><div style="margin-top:8px" class="form-hint">Token: '+escHtml(shortKey(playbackResult.token))+' • Gecerlilik: '+escHtml(new Date(playbackResult.expires_at).toLocaleString(localeForLang()))+'</div></div>':'')+studioRenderWarnings(state,codeBundle)+'</div>'+
        '</div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-card"><div><h2 class="studio-section-title">'+escHtml(codeBundle.label||'Uretilen kod')+'</h2><div class="studio-section-sub">Guncel secimlere gore uretilen kod ve baglantilar burada gorunur.</div></div><pre class="studio-code-block">'+escHtml(codeBundle.code||'Oncelikle stream secin ve cikti tipi belirleyin.')+'</pre><div class="studio-grid studio-grid-2">'+studioField('Player URL','<input class="input" readonly value="'+escHtml(codeBundle.playerURL||'')+'">','Player sayfasina direkt gider.')+studioField('Embed URL','<input class="input" readonly value="'+escHtml(codeBundle.embedURL||'')+'">','Iframe ve site embed ciktilari icin.')+studioField('Manifest URL','<input class="input" readonly value="'+escHtml(codeBundle.manifestURL||'')+'">','HLS veya DASH istemci testleri icin.')+studioField('VLC URL','<input class="input" readonly value="'+escHtml(codeBundle.vlcURL||'')+'">','Harici player ile test etmek icin.')+'</div></div><div class="studio-card"><div><h2 class="studio-section-title">Profil kutuphanesi</h2><div class="studio-section-sub">Her stream icin kaydedilebilir embed profilleri, marka bilgisi ve guvenlik profili saklanir.</div></div>'+(profileItems.length?'<div style="overflow:auto"><table class="studio-table"><thead><tr><th>Profil</th><th>Kullanim</th><th>Format</th><th>Mod</th><th>Guncel</th><th>Islem</th></tr></thead><tbody>'+profileItems.map(function(item){return '<tr><td><strong>'+escHtml(item.name)+'</strong><div class="form-hint">'+escHtml(item.notes||'')+'</div></td><td>'+escHtml(item.use_case||'-')+'</td><td>'+escHtml(item.primary_format||'-')+'</td><td>'+escHtml(item.mode||'-')+'</td><td>'+escHtml(new Date(item.updated_at).toLocaleString(localeForLang()))+'</td><td><div style="display:flex;gap:8px;flex-wrap:wrap"><button class="btn btn-secondary btn-sm" data-profile-load="'+Number(item.id)+'">Yukle</button><button class="btn btn-danger btn-sm" data-profile-delete="'+Number(item.id)+'">Sil</button></div></td></tr>';}).join('')+'</tbody></table></div>':'<div class="empty-state" style="padding:30px"><div class="icon"><i class="bi bi-bookmark-star"></i></div><h3>Kayitli embed profili yok</h3><p style="color:var(--text-muted)">Bu stream icin ilk embed profilini kaydederek tekrar kullanim, marka ve guvenlik ayarlarini tek tik hale getirebilirsin.</p></div>')+'</div></div>'+
      '</div>';

    bindEmbedStudioEvents(profileItems,codeBundle);
  }

  function bindEmbedStudioEvents(profileItems,codeBundle){
    const state=ensureEmbedState();
    document.querySelectorAll('#studio-embed-usecases .studio-option-card').forEach(function(btn){ btn.onclick=function(){ const preset=studioPresetForUseCase(btn.dataset.studioKey); const streamKey=state.streamKey; window.embedStudioState=studioMergeEmbedState(defaultEmbedState(),preset); window.embedStudioState.streamKey=streamKey; window.embedStudioState.useCase=btn.dataset.studioKey; studioRerender('embed-codes'); }; });
    document.querySelectorAll('#studio-embed-outputs .studio-option-card').forEach(function(btn){ btn.onclick=function(){ state.outputType=btn.dataset.studioKey||'iframe'; if(state.outputType==='audio' && !String(state.primaryFormat||'').includes('audio')) state.primaryFormat='dash_audio'; studioRerender('embed-codes'); }; });
    document.querySelectorAll('.segmented .segment').forEach(function(btn){ btn.onclick=function(){ state.mode=btn.dataset.mode||'simple'; studioRerender('embed-codes'); }; });
    [['studio-embed-stream',function(v){state.streamKey=v;state.profileId=0;}],['studio-embed-format',function(v){state.primaryFormat=v;}],['studio-embed-width',function(v){state.width=toNumber(v,1280);}],['studio-embed-height',function(v){state.height=toNumber(v,720);}],['studio-embed-theme',function(v){state.theme=v;}],['studio-embed-template',function(v){state.templateId=toNumber(v,0);}],['studio-embed-poster',function(v){state.branding.posterURL=v;state.options.posterURL=v;}],['studio-embed-watermark',function(v){state.branding.watermarkText=v;}],['studio-embed-start-quality',function(v){state.options.startQuality=v;}],['studio-embed-referrer',function(v){state.options.referrerPolicy=v;}],['studio-embed-share-package',function(v){state.options.sharePackage=v;}],['studio-sec-expiry',function(v){state.security.expiryMinutes=toNumber(v,60);}],['studio-sec-viewer',function(v){state.security.viewerID=v;}],['studio-sec-domain',function(v){state.security.domainRestriction=v;}],['studio-sec-ip',function(v){state.security.ipRestriction=v;}],['studio-sec-watermark',function(v){state.security.watermark=v;}],['studio-sec-notes',function(v){state.notes=v;}],['studio-embed-name',function(v){state.profileName=v;}],['studio-embed-profile',function(v){state.profileId=toNumber(v,0);}]].forEach(function(entry){ const el=document.getElementById(entry[0]); if(!el) return; el.onchange=function(){ entry[1](el.value); studioRerender('embed-codes'); }; });
    [['studio-embed-responsive',function(v){state.options.responsive=v;}],['studio-embed-autoplay',function(v){state.options.autoplay=v;}],['studio-embed-muted',function(v){state.options.muted=v;}],['studio-embed-debug-flag',function(v){state.options.debug=v;}],['studio-embed-audio-only',function(v){state.options.audioOnly=v;}],['studio-sec-signed',function(v){state.security.signedURL=v;}],['studio-sec-token',function(v){state.security.tokenRequired=v;}],['studio-sec-session',function(v){state.security.sessionBound=v;}],['studio-sec-apply',function(v){state.security.applyStreamPolicy=v;}]].forEach(function(entry){ const el=document.getElementById(entry[0]); if(!el) return; el.onchange=function(){ entry[1](!!el.checked); studioRerender('embed-codes'); }; });
    function openPrimary(debug){ const url=(state.outputType==='vlc'?codeBundle.vlcURL:(state.outputType==='manifest'?codeBundle.manifestURL:(state.outputType==='audio'?codeBundle.playerURL:(state.outputType==='popup'?codeBundle.playerURL:(state.outputType==='player'?codeBundle.playerURL:codeBundle.embedURL))))) || ''; if(!url){ toast('Acilacak baglanti bulunamadi','error'); return; } window.open(debug?withQueryParam(url,'debug','1'):url,'_blank','noopener'); }
    const copyCode=document.getElementById('studio-embed-copy-code'); if(copyCode) copyCode.onclick=function(){ copyText(codeBundle.code||''); };
    const copyURL=document.getElementById('studio-embed-copy-url'); if(copyURL) copyURL.onclick=function(){ copyText(codeBundle.primaryURL||codeBundle.playerURL||''); };
    const openBtn=document.getElementById('studio-embed-open'); if(openBtn) openBtn.onclick=function(){ openPrimary(false); };
    const debugBtn=document.getElementById('studio-embed-debug'); if(debugBtn) debugBtn.onclick=function(){ openPrimary(true); };
    const testBtn=document.getElementById('studio-embed-test'); if(testBtn) testBtn.onclick=function(){ openPrimary(!!state.options.debug); };
    const loadBtn=document.getElementById('studio-embed-load-profile'); if(loadBtn) loadBtn.onclick=async function(){ if(!state.profileId){ toast('Yuklemek icin once profil secin','warning'); return; } const item=await api('/api/admin/embed-profiles/'+state.profileId); if(!item || item.error){ toast('Profil yuklenemedi','error'); return; } const next=defaultEmbedState(); next.profileId=item.id||0; next.profileName=item.name||''; next.streamKey=item.stream_key||state.streamKey; next.useCase=item.use_case||'website'; next.mode=item.mode||'simple'; next.primaryFormat=item.primary_format||'player'; next.width=item.width||1280; next.height=item.height||720; next.theme=item.theme||'clean'; next.notes=item.notes||''; next.options=Object.assign(next.options,parseJSONSafeStudio(item.options_json,{})); next.branding=Object.assign(next.branding,parseJSONSafeStudio(item.branding_json,{})); next.security=Object.assign(next.security,parseJSONSafeStudio(item.security_json,{})); window.embedStudioState=next; studioRerender('embed-codes'); };
    const newBtn=document.getElementById('studio-embed-new-profile'); if(newBtn) newBtn.onclick=function(){ const streamKey=state.streamKey; window.embedStudioState=defaultEmbedState(); window.embedStudioState.streamKey=streamKey; studioRerender('embed-codes'); };
    const saveBtn=document.getElementById('studio-embed-save-profile'); if(saveBtn) saveBtn.onclick=async function(){ if(!state.streamKey){ toast('Once stream secin','warning'); return; } const name=(document.getElementById('studio-embed-name')||{}).value || state.profileName || 'Yeni embed profili'; const payload={stream_key:state.streamKey,name:name,use_case:state.useCase,mode:state.mode,primary_format:state.primaryFormat,width:state.width,height:state.height,theme:state.theme,options_json:JSON.stringify(state.options||{}),branding_json:JSON.stringify(state.branding||{}),security_json:JSON.stringify(state.security||{}),notes:state.notes||''}; const path=state.profileId?('/api/admin/embed-profiles/'+state.profileId):'/api/admin/embed-profiles'; const method=state.profileId?'PUT':'POST'; const res=await api(path,{method:method,body:payload}); if(!res || res.error){ toast('Embed profili kaydedilemedi','error'); return; } state.profileName=name; if(res.item && res.item.id) state.profileId=res.item.id; toast('Embed profili kaydedildi'); studioRerender('embed-codes'); };
    document.querySelectorAll('[data-profile-load]').forEach(function(btn){ btn.onclick=function(){ state.profileId=toNumber(btn.getAttribute('data-profile-load'),0); const load=document.getElementById('studio-embed-load-profile'); if(load) load.click(); }; });
    document.querySelectorAll('[data-profile-delete]').forEach(function(btn){ btn.onclick=async function(){ const id=toNumber(btn.getAttribute('data-profile-delete'),0); if(!id || !confirm('Bu embed profili silinsin mi?')) return; const res=await api('/api/admin/embed-profiles/'+id,{method:'DELETE'}); if(!res || res.error){ toast('Profil silinemedi','error'); return; } if(state.profileId===id) state.profileId=0; toast('Profil silindi'); studioRerender('embed-codes'); }; });
    const resetPolicy=document.getElementById('studio-sec-reset-policy');
    if(resetPolicy) resetPolicy.onclick=async function(){
      if(!state.streamKey || !confirm('Bu stream icin kalici playback guvenligini kaldirmak istiyor musunuz?')) return;
      const res=await api('/api/admin/security/stream-policy/reset',{method:'POST',body:{stream_key:state.streamKey,clear_domain_lock:true,clear_ip_whitelist:false}});
      if(!res || res.error){
        toast((res&&res.message)||'Kalici koruma kaldirilamadi','error');
        return;
      }
      state.security.applyStreamPolicy=false;
      toast('Kalici playback korumasi kaldirildi');
      studioRerender('embed-codes');
    };
  }

  window.renderEmbedCodes = async function(container){
    return renderEmbedStudio(container);
  };

  function defaultAnalyticsState(){
    return {period:'24h',streamKey:'',mode:'live'};
  }

  function ensureAnalyticsState(){
    if(!window.analyticsCenterState2) window.analyticsCenterState2=defaultAnalyticsState();
    return window.analyticsCenterState2;
  }

  function studioDownloadFile(filename,content,type){
    const blob=new Blob([content],{type:type||'text/plain;charset=utf-8'});
    const url=URL.createObjectURL(blob);
    const a=document.createElement('a');
    a.href=url;
    a.download=filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  }

  function analyticsCSV(data){
    const rows=[['stream','viewers','health','stalls','audio_switches','quality_transitions','status']];
    (Array.isArray(data.streams)?data.streams:[]).forEach(function(item){
      rows.push([
        item.name||'',
        item.viewer_count||0,
        item.health_score||0,
        ((item.telemetry||{}).total_stalls)||0,
        ((item.telemetry||{}).total_audio_switches)||0,
        ((item.telemetry||{}).total_quality_transitions)||0,
        item.status||''
      ]);
    });
    return rows.map(function(row){
      return row.map(function(cell){
        const value=String(cell==null?'':cell);
        return /[",\n]/.test(value)?('"'+value.replace(/"/g,'""')+'"'):value;
      }).join(',');
    }).join('\n');
  }

  function renderAnalyticsKPI(label,value,sub){
    return '<div class="studio-kpi"><div class="studio-kpi-label">'+escHtml(label)+'</div><div class="studio-kpi-value">'+escHtml(String(value))+'</div><div class="studio-kpi-sub">'+escHtml(sub||'')+'</div></div>';
  }

  function renderAnalyticsCenter(container,data,state){
    const history=data.history||{};
    const selected=(data.selected||{});
    const selectedStream=(selected.stream||null);
    const selectedHistory=Array.isArray(selected.history)?selected.history:[];
    const viewerPoints=(history.viewers||[]).map(function(item){ return {timestamp:item.timestamp,value:Number(item.value||0)}; });
    const bufferPoints=selectedHistory.map(function(item){ return {timestamp:item.created_at,value:Number(item.average_buffer_seconds||0)}; });
    const stallPoints=selectedHistory.map(function(item){ return {timestamp:item.created_at,value:Number(item.total_stalls||0)}; });
    const selectedTelemetry=selected.telemetry||{};
    const selectedTrackHistory=Array.isArray(selected.track_history)?selected.track_history:[];
    const qualityJSON=parseJSONSafeStudio((selectedHistory[selectedHistory.length-1]||{}).qualities_json||'{}',{});
    const audioJSON=parseJSONSafeStudio((selectedHistory[selectedHistory.length-1]||{}).audio_tracks_json||'{}',{});
    const sourceJSON=parseJSONSafeStudio((selectedHistory[selectedHistory.length-1]||{}).sources_json||'{}',{});
    const pageJSON=parseJSONSafeStudio((selectedHistory[selectedHistory.length-1]||{}).pages_json||'{}',{});
    const trackAgg={};
    selectedTrackHistory.forEach(function(item){
      const key=(item.display_label||item.kind||('Track '+item.track_id));
      trackAgg[key]=(trackAgg[key]||0)+Number(item.bitrate||0);
    });
    const trackBars=Object.keys(trackAgg).map(function(key){ return {label:key,value:Math.round(trackAgg[key]/1000)}; }).sort(function(a,b){ return b.value-a.value; });
    const risky=Array.isArray(data.risky_streams)?data.risky_streams:[];
    const sessions=Array.isArray(data.viewer_sessions)?data.viewer_sessions:[];
    const topStream=(data.kpis||{}).top_stream || '-';

    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero"><h1 class="studio-hero-title">Analitik Merkezi</h1><div class="studio-hero-sub">Canli ve gecmis telemetri, QoE sinyalleri, kalite gecisleri, audio switch davranisi ve sorunlu yayinlari tek merkezde izleyebilirsin.</div><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active">'+escHtml((history.label||'Canli pencere'))+'</span><span class="studio-pill">'+escHtml((selectedStream&&selectedStream.name)||'Tum streamler')+'</span><span class="studio-pill">'+escHtml((state.mode||'live')==='live'?'Canli izleme':'Gecmis rapor')+'</span></div></section>'+
        '<section class="studio-toolbar"><div class="studio-toolbar-group"><select id="studio-analytics-period" class="input">'+studioSelectOptions([['24h','Son 24 saat'],['7d','Son 7 gun'],['30d','Son 30 gun']],state.period)+'</select><select id="studio-analytics-stream" class="input"><option value="">Tum streamler</option>'+(Array.isArray(data.streams)?data.streams:[]).map(function(item){ return '<option value="'+escHtml(item.stream_key)+'"'+(item.stream_key===state.streamKey?' selected':'')+'>'+escHtml(item.name)+'</option>'; }).join('')+'</select><div class="segmented"><button class="segment'+(state.mode==='live'?' active':'')+'" data-analytics-mode="live" type="button">Canli</button><button class="segment'+(state.mode==='history'?' active':'')+'" data-analytics-mode="history" type="button">Gecmis</button></div></div><div class="studio-toolbar-group"><button class="btn btn-secondary" id="studio-analytics-export-json">JSON Disa Aktar</button><button class="btn btn-secondary" id="studio-analytics-export-csv">CSV Disa Aktar</button><button class="btn btn-primary" id="studio-analytics-refresh">Yenile</button></div></section>'+
        '<div class="studio-kpi-grid">'+
          renderAnalyticsKPI('Aktif Izleyici',fmtInt((data.kpis||{}).active_viewers||0),'Canli oturum sayisi')+
          renderAnalyticsKPI('Tepe Izleyici',fmtInt((data.kpis||{}).peak_viewers||0),'Secilen penceredeki en yuksek eszamanli izleyici')+
          renderAnalyticsKPI('Ortalama Buffer',fmtStudioNumber((data.kpis||{}).average_buffer,1)+' sn','Player tarafindaki ortalama buffer')+
          renderAnalyticsKPI('Toplam Stall',fmtInt((data.kpis||{}).stalls||0),'Takilma olaylari')+
          renderAnalyticsKPI('Kalite Gecisi',fmtInt((data.kpis||{}).quality_transitions||0),'ABR kalite degisimi sayisi')+
          renderAnalyticsKPI('Audio Switch',fmtInt((data.kpis||{}).audio_switches||0),'Ses izi degisimi sayisi')+
          renderAnalyticsKPI('Hata Orani',fmtStudioNumber((data.kpis||{}).error_rate,1)+'%','Waiting ve offline oturum orani')+
          renderAnalyticsKPI('En Cok Izlenen',topStream,'O anki lider stream')+
        '</div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-chart-card"><div class="studio-section-title">Izleyici zaman serisi</div><div class="studio-section-sub">Dashboard snapshot verileriyle olusan ozet trend.</div><div class="studio-chart-wrap">'+renderTimelineChart(viewerPoints,'Izleyici verisi yok',function(v){ return fmtInt(v); },{meta:history.label||'Secili pencere'})+'</div></div><div class="studio-chart-card"><div class="studio-section-title">Buffer trendi</div><div class="studio-section-sub">Secilen stream icin son telemetri ornekleri.</div><div class="studio-chart-wrap">'+renderTimelineChart(bufferPoints,'Buffer ornegi yok',function(v){ return Number(v||0).toFixed(1)+' sn'; },{meta:selectedStream?(selectedStream.name+' buffer'):history.label||''})+'</div></div></div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-chart-card"><div class="studio-section-title">Stall trendi</div><div class="studio-section-sub">Secilen stream icin birikimli stall davranisi.</div><div class="studio-chart-wrap">'+renderTimelineChart(stallPoints,'Stall verisi yok',function(v){ return fmtInt(v); },{meta:selectedStream?(selectedStream.name+' stall'):history.label||''})+'</div></div><div class="studio-chart-card"><div class="studio-section-title">ABR ve audio dagilimi</div><div class="studio-section-sub">Son telemetri ornegindeki kalite, ses ve player kaynagi dagilimi.</div><div class="studio-grid" style="gap:14px"><div>'+renderBarList(Object.keys(qualityJSON).map(function(key){ return {label:key,value:Number(qualityJSON[key]||0)}; }),'Kalite verisi yok',function(v){ return fmtInt(v); })+'</div><div>'+renderBarList(Object.keys(audioJSON).map(function(key){ return {label:key,value:Number(audioJSON[key]||0)}; }),'Ses izi verisi yok',function(v){ return fmtInt(v); })+'</div></div></div></div>'+
        '<div class="studio-grid studio-grid-3"><div class="studio-card"><div class="studio-section-title">Kaynak dagilimi</div>'+renderBarList(Object.keys(sourceJSON).map(function(key){ return {label:key,value:Number(sourceJSON[key]||0)}; }),'Kaynak verisi yok',function(v){ return fmtInt(v); })+'</div><div class="studio-card"><div class="studio-section-title">Sayfa dagilimi</div>'+renderBarList(Object.keys(pageJSON).map(function(key){ return {label:key,value:Number(pageJSON[key]||0)}; }),'Sayfa verisi yok',function(v){ return fmtInt(v); })+'</div><div class="studio-card"><div class="studio-section-title">ABR katman dagilimi</div>'+renderBarList(trackBars,'Track ornegi yok',function(v){ return fmtInt(v)+' kbps'; })+'</div></div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-card"><div class="studio-section-title">Sorunlu yayinlar</div><div class="studio-section-sub">QoE uyari esiklerini asan yayinlar burada one cikar.</div>'+(risky.length?'<div style="overflow:auto"><table class="studio-table"><thead><tr><th>Yayin</th><th>Durum</th><th>Izleyici</th><th>Uyarilar</th><th>Islem</th></tr></thead><tbody>'+risky.map(function(item){ return '<tr><td><strong>'+escHtml(item.name)+'</strong><div class="form-hint">'+escHtml(item.stream_key)+'</div></td><td>'+escHtml(item.status||'-')+'</td><td>'+fmtInt(item.viewer_count||0)+'</td><td>'+escHtml((item.alerts||[]).map(function(a){ return a.message||a.title||a.code; }).join(' • ')||'-')+'</td><td><button class="btn btn-secondary btn-sm" data-open-ops="'+escHtml(item.stream_key)+'">Operasyon Merkezi</button></td></tr>'; }).join('')+'</tbody></table></div>':'<div class="empty-state" style="padding:28px"><div class="icon"><i class="bi bi-shield-check"></i></div><h3>Su an kirmizi alarm yok</h3><p style="color:var(--text-muted)">Secilen pencere icinde uyari sinirlarini asan yayin bulunmadi.</p></div>')+'</div><div class="studio-card"><div class="studio-section-title">Kalite ve audio degisim raporu</div><div class="studio-section-sub">Secilen stream icin telemetri ozetleri.</div><div class="studio-kpi-grid" style="grid-template-columns:repeat(2,minmax(0,1fr))">'+renderAnalyticsKPI('Kalite Gecisi',fmtInt((selectedTelemetry.total_quality_transitions)||0),'ABR katmanlari arasi gecis')+renderAnalyticsKPI('Audio Switch',fmtInt((selectedTelemetry.total_audio_switches)||0),'Ses izi degisimi')+renderAnalyticsKPI('Toplam Recoveries',fmtInt((selectedTelemetry.total_recoveries)||0),'Player toparlanma sayisi')+renderAnalyticsKPI('Son Hata',((selectedTelemetry.last_error)||'-').slice(0,44),selectedStream?selectedStream.name:'Secili stream yok')+'</div></div></div>'+
        '<div class="studio-card"><div class="studio-section-title">Canli player oturumlari</div><div class="studio-section-sub">Aktif oturumlar, kalite ve audio secimi ile birlikte listelenir.</div>'+renderTelemetrySessionsTable(sessions)+'</div>'+
      '</div>';

    document.querySelectorAll('[data-analytics-mode]').forEach(function(btn){ btn.onclick=function(){ const next=ensureAnalyticsState(); next.mode=btn.getAttribute('data-analytics-mode')||'live'; studioRerender('analytics'); }; });
    const period=document.getElementById('studio-analytics-period');
    if(period) period.onchange=function(){ const next=ensureAnalyticsState(); next.period=period.value||'24h'; studioRerender('analytics'); };
    const stream=document.getElementById('studio-analytics-stream');
    if(stream) stream.onchange=function(){ const next=ensureAnalyticsState(); next.streamKey=stream.value||''; studioRerender('analytics'); };
    const exportJSON=document.getElementById('studio-analytics-export-json');
    if(exportJSON) exportJSON.onclick=function(){ studioDownloadFile('fluxstream-analytics-'+(state.period||'24h')+'.json',JSON.stringify(data,null,2),'application/json;charset=utf-8'); };
    const exportCSV=document.getElementById('studio-analytics-export-csv');
    if(exportCSV) exportCSV.onclick=function(){ studioDownloadFile('fluxstream-analytics-'+(state.period||'24h')+'.csv',analyticsCSV(data),'text/csv;charset=utf-8'); };
    const refresh=document.getElementById('studio-analytics-refresh');
    if(refresh) refresh.onclick=function(){ studioRerender('analytics'); };
    document.querySelectorAll('[data-open-ops]').forEach(function(btn){ btn.onclick=function(){ const streamKey=btn.getAttribute('data-open-ops')||''; const match=(Array.isArray(data.streams)?data.streams:[]).find(function(item){ return item.stream_key===streamKey; }); if(match && typeof selectOperationsStream==='function'){ selectOperationsStream(match.id||0); if(typeof setOperationsCenterTab==='function') setOperationsCenterTab('qoe'); } navigate('operations-center'); }; });
    if(typeof schedulePageRefresh==='function') schedulePageRefresh('analytics',15000);
  }

  window.renderAnalytics = async function(container){
    const state=ensureAnalyticsState();
    const params='?period='+encodeURIComponent(state.period||'24h')+'&mode='+encodeURIComponent(state.mode||'live')+(state.streamKey?('&stream_key='+encodeURIComponent(state.streamKey)):'');
    const data=await api('/api/admin/analytics/center'+params);
    return renderAnalyticsCenter(container,data||{},state);
  };

  const ABR_PRESETS = {
    balanced:{title:'Dengeli', desc:'Genel amacli TV ve webinar yayini.', scope:'global', layers:[{name:'1080p',width:1920,height:1080,bitrate:'4500k',max_bitrate:'5000k',buf_size:'9000k',fps:30,preset:'fast',audio_rate:'192k'},{name:'720p',width:1280,height:720,bitrate:'2500k',max_bitrate:'3000k',buf_size:'5000k',fps:30,preset:'fast',audio_rate:'128k'},{name:'480p',width:854,height:480,bitrate:'1100k',max_bitrate:'1200k',buf_size:'2000k',fps:30,preset:'fast',audio_rate:'96k'},{name:'360p',width:640,height:360,bitrate:'600k',max_bitrate:'700k',buf_size:'1200k',fps:25,preset:'fast',audio_rate:'64k'}]},
    mobile:{title:'Mobil', desc:'Dusuk bantta daha dayanikli dagitim.', scope:'global', layers:[{name:'720p',width:1280,height:720,bitrate:'1800k',max_bitrate:'2200k',buf_size:'3600k',fps:30,preset:'fast',audio_rate:'128k'},{name:'480p',width:854,height:480,bitrate:'900k',max_bitrate:'1100k',buf_size:'1800k',fps:30,preset:'fast',audio_rate:'96k'},{name:'360p',width:640,height:360,bitrate:'480k',max_bitrate:'600k',buf_size:'1000k',fps:24,preset:'veryfast',audio_rate:'64k'}]},
    resilient:{title:'Dayanikli', desc:'Sorunlu aglarda donmayi azaltmaya odaklanir.', scope:'global', layers:[{name:'720p',width:1280,height:720,bitrate:'1600k',max_bitrate:'1900k',buf_size:'3200k',fps:25,preset:'veryfast',audio_rate:'128k'},{name:'480p',width:854,height:480,bitrate:'700k',max_bitrate:'900k',buf_size:'1500k',fps:25,preset:'veryfast',audio_rate:'96k'},{name:'360p',width:640,height:360,bitrate:'420k',max_bitrate:'520k',buf_size:'900k',fps:24,preset:'superfast',audio_rate:'64k'}]},
    tv:{title:'TV', desc:'Yuksek kalite ekranlar icin daha genis merdiven.', scope:'global', layers:[{name:'1080p',width:1920,height:1080,bitrate:'5500k',max_bitrate:'6200k',buf_size:'11000k',fps:30,preset:'fast',audio_rate:'192k'},{name:'720p',width:1280,height:720,bitrate:'3200k',max_bitrate:'3600k',buf_size:'6400k',fps:30,preset:'fast',audio_rate:'160k'},{name:'540p',width:960,height:540,bitrate:'1800k',max_bitrate:'2100k',buf_size:'3600k',fps:30,preset:'fast',audio_rate:'128k'}]},
    high_quality:{title:'Yuksek Kalite', desc:'Goruntu kalitesini onceler, CPU yuksek olabilir.', scope:'global', layers:[{name:'1440p',width:2560,height:1440,bitrate:'9000k',max_bitrate:'9800k',buf_size:'18000k',fps:30,preset:'faster',audio_rate:'192k'},{name:'1080p',width:1920,height:1080,bitrate:'5200k',max_bitrate:'5800k',buf_size:'10400k',fps:30,preset:'fast',audio_rate:'192k'},{name:'720p',width:1280,height:720,bitrate:'2600k',max_bitrate:'3200k',buf_size:'5200k',fps:30,preset:'fast',audio_rate:'128k'}]},
    audio_only:{title:'Audio-only', desc:'Radyo, podcast ve DASH ses dagitimi icin.', scope:'global', layers:[{name:'audio-192k',width:0,height:0,bitrate:'0',max_bitrate:'0',buf_size:'0',fps:0,preset:'copy',audio_rate:'192k'},{name:'audio-96k',width:0,height:0,bitrate:'0',max_bitrate:'0',buf_size:'0',fps:0,preset:'copy',audio_rate:'96k'}]},
    radio:{title:'Radyo', desc:'Dusuk maliyetli surekli ses yayini icin.', scope:'global', layers:[{name:'audio-128k',width:0,height:0,bitrate:'0',max_bitrate:'0',buf_size:'0',fps:0,preset:'copy',audio_rate:'128k'},{name:'audio-64k',width:0,height:0,bitrate:'0',max_bitrate:'0',buf_size:'0',fps:0,preset:'copy',audio_rate:'64k'}]},
    low_band:{title:'Sadece dusuk bant', desc:'Dar bant ve saha baglantilari icin minimum merdiven.', scope:'global', layers:[{name:'480p',width:854,height:480,bitrate:'650k',max_bitrate:'800k',buf_size:'1300k',fps:25,preset:'veryfast',audio_rate:'96k'},{name:'360p',width:640,height:360,bitrate:'360k',max_bitrate:'450k',buf_size:'800k',fps:24,preset:'superfast',audio_rate:'64k'},{name:'240p',width:426,height:240,bitrate:'220k',max_bitrate:'280k',buf_size:'500k',fps:20,preset:'superfast',audio_rate:'48k'}]}
  };

  const ABR_RESOLUTION_OPTIONS = [
    {key:'1440p',label:'1440p / 2560x1440',width:2560,height:1440,default_name:'1440p',fps:30},
    {key:'1080p',label:'1080p / 1920x1080',width:1920,height:1080,default_name:'1080p',fps:30},
    {key:'720p',label:'720p / 1280x720',width:1280,height:720,default_name:'720p',fps:30},
    {key:'540p',label:'540p / 960x540',width:960,height:540,default_name:'540p',fps:30},
    {key:'480p',label:'480p / 854x480',width:854,height:480,default_name:'480p',fps:30},
    {key:'360p',label:'360p / 640x360',width:640,height:360,default_name:'360p',fps:25},
    {key:'240p',label:'240p / 426x240',width:426,height:240,default_name:'240p',fps:20},
    {key:'audio',label:'Audio-only katman',width:0,height:0,default_name:'audio',fps:0}
  ];

  const ABR_BITRATE_PACKS = [
    {key:'ultra_1440',label:'Ultra 1440p',bitrate:'9000k',max_bitrate:'9800k',buf_size:'18000k'},
    {key:'strong_1080',label:'Guclu 1080p',bitrate:'5500k',max_bitrate:'6200k',buf_size:'11000k'},
    {key:'balanced_1080',label:'Dengeli 1080p',bitrate:'4500k',max_bitrate:'5000k',buf_size:'9000k'},
    {key:'balanced_720',label:'Dengeli 720p',bitrate:'2500k',max_bitrate:'3000k',buf_size:'5000k'},
    {key:'mobile_720',label:'Mobil 720p',bitrate:'1800k',max_bitrate:'2200k',buf_size:'3600k'},
    {key:'balanced_540',label:'Dengeli 540p',bitrate:'1800k',max_bitrate:'2100k',buf_size:'3600k'},
    {key:'balanced_480',label:'Dengeli 480p',bitrate:'1100k',max_bitrate:'1200k',buf_size:'2000k'},
    {key:'mobile_480',label:'Mobil 480p',bitrate:'900k',max_bitrate:'1100k',buf_size:'1800k'},
    {key:'resilient_480',label:'Dayanikli 480p',bitrate:'700k',max_bitrate:'900k',buf_size:'1500k'},
    {key:'balanced_360',label:'Dengeli 360p',bitrate:'600k',max_bitrate:'700k',buf_size:'1200k'},
    {key:'mobile_360',label:'Mobil 360p',bitrate:'480k',max_bitrate:'600k',buf_size:'1000k'},
    {key:'resilient_360',label:'Dayanikli 360p',bitrate:'420k',max_bitrate:'520k',buf_size:'900k'},
    {key:'low_240',label:'Dusuk 240p',bitrate:'220k',max_bitrate:'280k',buf_size:'500k'},
    {key:'audio_passthrough',label:'Audio-only',bitrate:'0',max_bitrate:'0',buf_size:'0'}
  ];

  const ABR_AUDIO_OPTIONS = [['192k','192 kbps'],['160k','160 kbps'],['128k','128 kbps'],['96k','96 kbps'],['64k','64 kbps'],['48k','48 kbps']];
  const ABR_FPS_OPTIONS = [['0','0 / audio-only'],['20','20'],['24','24'],['25','25'],['30','30'],['50','50'],['60','60']];
  const ABR_PRESET_OPTIONS = [['copy','copy'],['superfast','superfast'],['veryfast','veryfast'],['fast','fast'],['faster','faster']];

  function defaultABRState(){
    return {mode:'simple',profileId:0,profileSet:'balanced',name:'',description:'',scope:'global',streamKey:'',compareKey:'',layers:cloneJSON(ABR_PRESETS.balanced.layers),json:'',importJSON:''};
  }

  function ensureABRState(){
    if(!window.abrStudioState) window.abrStudioState=defaultABRState();
    return window.abrStudioState;
  }

  function normalizeABRLayer(layer,index){
    const next=Object.assign({name:'Katman '+(index+1),width:0,height:0,bitrate:'0',max_bitrate:'0',buf_size:'0',fps:0,preset:'fast',audio_rate:'128k'},layer||{});
    return next;
  }

  function parseABRProfilesJSON(raw){
    const parsed=parseJSONSafeStudio(raw,{});
    const out={};
    Object.keys(parsed||{}).forEach(function(key){
      out[key]=(parsed[key]||[]).map(function(item,index){ return normalizeABRLayer(item,index); });
    });
    return out;
  }

  function bitrateToNumber(raw){
    const value=String(raw||'0').trim().toLowerCase();
    const num=parseFloat(value)||0;
    if(value.endsWith('k')) return num*1000;
    if(value.endsWith('m')) return num*1000000;
    return num;
  }

  function summarizeABRLayers(layers){
    const list=(Array.isArray(layers)?layers:[]).map(normalizeABRLayer);
    let upload=0, peak=0, cpu=0, lowBand=0, audioOnly=true;
    list.forEach(function(layer){
      const bitrate=bitrateToNumber(layer.bitrate);
      const maxBitrate=bitrateToNumber(layer.max_bitrate||layer.bitrate);
      const pixels=Math.max(1,(Number(layer.width||0)*Number(layer.height||0))||1);
      const fps=Math.max(1,Number(layer.fps||24)||24);
      upload+=Math.max(bitrate,maxBitrate);
      peak+=maxBitrate||bitrate;
      cpu+=Math.round((pixels/921600)*(fps/30)*22);
      if((Number(layer.width||0)<=640 && bitrate<=700000) || bitrate<=500000) lowBand++;
      if(Number(layer.width||0)>0 && Number(layer.height||0)>0) audioOnly=false;
    });
    return {
      variants:list.length,
      uploadMbps:(upload/1000000),
      peakMbps:(peak/1000000),
      cpuScore:cpu,
      lowBandScore:list.length?Math.round((lowBand/list.length)*100):0,
      audioOnly:audioOnly
    };
  }

  function abrLayersToJSON(layers){
    return JSON.stringify((Array.isArray(layers)?layers:[]).map(function(layer,index){ return normalizeABRLayer(layer,index); }),null,2);
  }

  function abrRecommendation(stream,diagnostics){
    if(!stream) return {preset:'balanced',why:'Varsayilan yayin davranisi icin dengeli profil iyi baslangictir.'};
    if((diagnostics&&diagnostics.delivery_summary&&diagnostics.delivery_summary.label==='ABR bekliyor') || Number(stream.input_width||0)>=1920) return {preset:'resilient',why:'Canli varyant sertlestirmesi ve yuksek giris cozunurlugu icin dayaniksiz degil, dayanikli merdiven daha guvenli.'};
    if(Number(stream.input_width||0)>0 && Number(stream.input_width||0)<=960) return {preset:'mobile',why:'Kaynak zaten dusuk/mobil odakli gorunuyor.'};
    return {preset:'balanced',why:'Genel kullanim ve coklu ekran dagitimi icin dengeli profil uygun.'};
  }

  function abrResolutionKey(layer){
    const item=ABR_RESOLUTION_OPTIONS.find(function(option){
      return Number(option.width)===Number(layer.width||0) && Number(option.height)===Number(layer.height||0);
    });
    return item?item.key:(Number(layer.width||0)===0 && Number(layer.height||0)===0?'audio':'custom');
  }

  function abrBitratePackKey(layer){
    const item=ABR_BITRATE_PACKS.find(function(pack){
      return String(pack.bitrate)===String(layer.bitrate||'0')
        && String(pack.max_bitrate)===String(layer.max_bitrate||layer.bitrate||'0')
        && String(pack.buf_size)===String(layer.buf_size||'0');
    });
    if(item) return item.key;
    return Number(layer.width||0)===0 && Number(layer.height||0)===0 ? 'audio_passthrough' : 'balanced_720';
  }

  function applyABRResolutionChoice(layer,key,index){
    const picked=ABR_RESOLUTION_OPTIONS.find(function(option){ return option.key===key; }) || ABR_RESOLUTION_OPTIONS[2];
    layer.width=picked.width;
    layer.height=picked.height;
    layer.fps=picked.fps;
    if(!layer.name || /^Katman \d+$/i.test(layer.name) || /^audio/i.test(layer.name) || /\d+p/i.test(layer.name)) layer.name=picked.default_name;
    if(picked.key==='audio'){
      layer.bitrate='0';
      layer.max_bitrate='0';
      layer.buf_size='0';
      layer.preset='copy';
      layer.audio_rate=layer.audio_rate||'128k';
    }else if(String(layer.bitrate||'0')==='0'){
      const fallbackPack=ABR_BITRATE_PACKS.find(function(pack){ return pack.key.indexOf(picked.key)>=0; }) || ABR_BITRATE_PACKS.find(function(pack){ return pack.key==='balanced_720'; });
      if(fallbackPack) applyABRBitratePack(layer,fallbackPack.key);
    }
    return normalizeABRLayer(layer,index);
  }

  function applyABRBitratePack(layer,key){
    const picked=ABR_BITRATE_PACKS.find(function(pack){ return pack.key===key; }) || ABR_BITRATE_PACKS[0];
    layer.bitrate=picked.bitrate;
    layer.max_bitrate=picked.max_bitrate;
    layer.buf_size=picked.buf_size;
    if(key==='audio_passthrough'){
      layer.width=0;
      layer.height=0;
      layer.fps=0;
      layer.preset='copy';
    }
    return layer;
  }

  function renderABRLayerRow(layer,index,total,mode){
    layer=normalizeABRLayer(layer,index);
    const resolutionKey=abrResolutionKey(layer);
    const bitrateKey=abrBitratePackKey(layer);
    const simpleFields=
      '<div class="studio-layer-simple">'+
        studioField('Katman adi','<input class="input" data-layer-field="name" data-layer-index="'+index+'" value="'+escHtml(layer.name||'')+'">','Profil kutuphanesinde gorunecek ad.')+
        studioField('Cozunurluk','<select class="input" data-layer-select="resolution" data-layer-index="'+index+'">'+studioSelectOptions(ABR_RESOLUTION_OPTIONS.map(function(item){ return [item.key,item.label]; }),resolutionKey)+'</select>','Hazir cozumlerden birini sec.')+
        studioField('Bitrate paketi','<select class="input" data-layer-select="bitrate_pack" data-layer-index="'+index+'">'+studioSelectOptions(ABR_BITRATE_PACKS.map(function(item){ return [item.key,item.label]; }),bitrateKey)+'</select>','Video bitrate, max ve buffer birlikte secilir.')+
        studioField('FPS','<select class="input" data-layer-field="fps" data-layer-index="'+index+'">'+studioSelectOptions(ABR_FPS_OPTIONS,String(layer.fps||0))+'</select>','Kare hizi secimi.')+
        studioField('Encoder hizi','<select class="input" data-layer-field="preset" data-layer-index="'+index+'">'+studioSelectOptions(ABR_PRESET_OPTIONS,layer.preset||'fast')+'</select>','FFmpeg encoder profili.')+
        studioField('Ses bitrate','<select class="input" data-layer-field="audio_rate" data-layer-index="'+index+'">'+studioSelectOptions(ABR_AUDIO_OPTIONS,layer.audio_rate||'128k')+'</select>','Katmana eslik eden ses hizi.')+
      '</div>';
    const advancedFields=
      '<div class="studio-layer-manual">'+
        studioField('Genislik','<input class="input" type="number" min="0" data-layer-field="width" data-layer-index="'+index+'" value="'+escHtml(String(layer.width||0))+'">','0 ise audio-only kabul edilir.')+
        studioField('Yukseklik','<input class="input" type="number" min="0" data-layer-field="height" data-layer-index="'+index+'" value="'+escHtml(String(layer.height||0))+'">','0 ise audio-only kabul edilir.')+
        studioField('Bitrate','<input class="input" data-layer-field="bitrate" data-layer-index="'+index+'" value="'+escHtml(layer.bitrate||'0')+'">','Ornek: 2500k')+
        studioField('Max bitrate','<input class="input" data-layer-field="max_bitrate" data-layer-index="'+index+'" value="'+escHtml(layer.max_bitrate||layer.bitrate||'0')+'">','ABR ust siniri.')+
        studioField('Buffer','<input class="input" data-layer-field="buf_size" data-layer-index="'+index+'" value="'+escHtml(layer.buf_size||'0')+'">','Encoder buffer boyutu.')+
      '</div>';
    return '<div class="studio-layer" data-layer-index="'+index+'"><div class="studio-layer-head"><strong>'+escHtml(layer.name||('Katman '+(index+1)))+'</strong><div style="display:flex;gap:8px;flex-wrap:wrap"><button class="btn btn-secondary btn-sm" data-layer-up="'+index+'"'+(index===0?' disabled':'')+'>Yukari</button><button class="btn btn-secondary btn-sm" data-layer-down="'+index+'"'+(index===total-1?' disabled':'')+'>Asagi</button><button class="btn btn-danger btn-sm" data-layer-delete="'+index+'">Sil</button></div></div><div class="studio-layer-meta"><span class="studio-chip active">'+escHtml(layer.width&&layer.height?(layer.width+'x'+layer.height):'Audio-only')+'</span><span class="studio-chip">'+escHtml((layer.bitrate||'0')+' / '+(layer.max_bitrate||layer.bitrate||'0'))+'</span><span class="studio-chip">'+escHtml((layer.audio_rate||'128k')+' ses')+'</span></div>'+simpleFields+(mode==='advanced'?advancedFields:'')+'</div>';
  }

  window.renderSettingsABR = async function(container){
    const state=ensureABRState();
    const [settings,streams,saved]=await Promise.all([api('/api/settings'),api('/api/streams'),api('/api/admin/abr-profiles')]);
    const streamList=Array.isArray(streams)?streams:[];
    if(!state.streamKey && streamList[0]) state.streamKey=streamList[0].stream_key;
    const stream=streamList.find(function(item){ return item.stream_key===state.streamKey; }) || streamList[0] || null;
    const diagnostics=stream?await api('/api/diagnostics/stream/'+stream.id):{};
    const savedProfiles=Array.isArray(saved)?saved:[];
    const globalProfiles=parseABRProfilesJSON((settings&&settings.abr_profiles_json)||'{}');
    if(!state.layers || !state.layers.length){
      const baseKey=state.profileSet || (settings&&settings.abr_profile_set) || 'balanced';
      state.layers=cloneJSON(globalProfiles[baseKey] || (ABR_PRESETS[baseKey]&&ABR_PRESETS[baseKey].layers) || ABR_PRESETS.balanced.layers);
      state.profileSet=baseKey;
    }
    state.json=abrLayersToJSON(state.layers);
    const summary=summarizeABRLayers(state.layers);
    const recommendation=abrRecommendation(stream,diagnostics);
    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero"><h1 class="studio-hero-title">ABR Profilleri ve Teslimat Merkezi</h1><div class="studio-hero-sub">Ham JSON duzenleme yerine secilebilir, kaydedilebilir ve tekrar kullanilabilir ABR profil studyosu. HLS, DASH ve audio-only teslimat davranisini burada modelleyebilir ve uygulayabilirsin.</div><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active">'+escHtml((state.profileSet||'balanced'))+'</span><span class="studio-pill">'+escHtml((stream&&stream.name)||'Stream sec')+'</span><span class="studio-pill">'+escHtml(summary.audioOnly?'Audio-only':'Video + Audio')+'</span><span class="studio-pill">'+escHtml((diagnostics.delivery_summary&&diagnostics.delivery_summary.label)||'Teslimat ozeti')+'</span></div></section>'+
        '<section class="studio-toolbar"><div class="studio-toolbar-group"><select id="studio-abr-stream" class="input">'+streamList.map(function(item){ return '<option value="'+escHtml(item.stream_key)+'"'+(item.stream_key===state.streamKey?' selected':'')+'>'+escHtml(item.name)+'</option>'; }).join('')+'</select><select id="studio-abr-saved" class="input">'+studioProfileSelectOptions(savedProfiles,state.profileId)+'</select><button class="btn btn-secondary" id="studio-abr-load">Profili Yukle</button></div><div class="studio-toolbar-group"><div class="segmented"><button class="segment'+(state.mode==='simple'?' active':'')+'" data-abr-mode="simple" type="button">Basit Mod</button><button class="segment'+(state.mode==='advanced'?' active':'')+'" data-abr-mode="advanced" type="button">Gelismis Mod</button></div><button class="btn btn-secondary" id="studio-abr-duplicate">Cogalt</button><button class="btn btn-primary" id="studio-abr-save">Profili Kaydet</button></div></section>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-card soft"><div><h2 class="studio-section-title">Hazir preset kutuphanesi</h2><div class="studio-section-sub">Kullanici JSON yazmak zorunda kalmadan secilebilir profil setleri.</div></div><div class="studio-option-grid">'+Object.keys(ABR_PRESETS).map(function(key){ const item=ABR_PRESETS[key]; return studioOptionCard({key:key,title:item.title,desc:item.desc,badge:key===recommendation.preset?'Oneri':'Preset'},key===state.profileSet); }).join('')+'</div><div class="studio-alert info"><strong>Yayin bazli onerim</strong><div style="margin-top:8px" class="form-hint">'+escHtml((ABR_PRESETS[recommendation.preset]&&ABR_PRESETS[recommendation.preset].title)||recommendation.preset)+' • '+escHtml(recommendation.why)+'</div></div></div><div class="studio-card"><div><h2 class="studio-section-title">Profil ozeti</h2><div class="studio-section-sub">Secili merdivenin tahmini kaynak kullanimi ve cikti yapisi.</div></div><div class="studio-kpi-grid" style="grid-template-columns:repeat(2,minmax(0,1fr))">'+renderAnalyticsKPI('Varyant',fmtInt(summary.variants),'Tahmini HLS / DASH katman sayisi')+renderAnalyticsKPI('Tahmini Upload',summary.uploadMbps.toFixed(2)+' Mbps','Bitrate toplami')+renderAnalyticsKPI('Tepe Upload',summary.peakMbps.toFixed(2)+' Mbps','Max bitrate toplami')+renderAnalyticsKPI('CPU Skoru',fmtInt(summary.cpuScore),'Yaklasik encoder maliyeti')+renderAnalyticsKPI('Dusuk Bant Uyumu',fmtInt(summary.lowBandScore)+'%','Daha dusuk katman yogunlugu')+renderAnalyticsKPI('Teslimat Sagligi',(diagnostics.delivery_summary&&diagnostics.delivery_summary.label)||'-',(diagnostics.delivery_summary&&diagnostics.delivery_summary.description)||'')+'</div><div class="studio-chip-row">'+(diagnostics.checks||[]).map(function(item){ return '<span class="studio-chip'+(item.tone==='green'?' active':'')+'">'+escHtml(item.description)+' • '+escHtml(item.label)+'</span>'; }).join('')+'</div><div class="studio-grid studio-grid-2">'+studioField('Profil anahtari','<input id="studio-abr-profile-set" class="input" value="'+escHtml(state.profileSet||'balanced')+'">','Kayitli profillerde benzersiz anahtar.')+studioField('Profil adi','<input id="studio-abr-name" class="input" value="'+escHtml(state.name||'')+'" placeholder="Ornek: Mobil Dayanikli">','Kutuphanede gosterilecek ad.')+studioField('Aciklama','<input id="studio-abr-description" class="input" value="'+escHtml(state.description||'')+'" placeholder="Dusuk bant saha yayini">','Profilin nerede kullanilacagi.')+studioField('Uygulama kapsami','<select id="studio-abr-scope" class="input">'+studioSelectOptions([['global','Global varsayilan'],['stream','Sadece secili stream']],state.scope)+'</select>','Global veya stream bazli uygula.')+'</div><div class="studio-chip-row"><button class="btn btn-secondary" id="studio-abr-import">JSON Icice Al</button><button class="btn btn-secondary" id="studio-abr-export">JSON Disari Aktar</button><button class="btn btn-secondary" id="studio-abr-apply-current">Mevcut Profili Uygula</button><button class="btn btn-secondary" id="studio-abr-delete">Secili Kayitli Profili Sil</button></div></div></div>'+
        '<div class="studio-card"><div><h2 class="studio-section-title">Katman olusturucu</h2><div class="studio-section-sub">Manuel sayi girmek yerine hazir cozum ve bitrate paketleri secilir. Istersen gelismis modda ayrintilari yine elle duzenleyebilirsin.</div></div><div style="display:flex;gap:10px;flex-wrap:wrap"><button class="btn btn-primary" id="studio-abr-add-layer">Katman Ekle</button><button class="btn btn-secondary" id="studio-abr-apply-preset">Preseti Yukle</button><button class="btn btn-secondary" id="studio-abr-reset">Sifirla</button></div><div id="studio-abr-layer-list">'+(state.layers||[]).map(function(layer,index,arr){ return renderABRLayerRow(layer,index,arr.length,state.mode); }).join('')+'</div>'+(state.mode==='advanced'?'<div class="studio-card soft" style="padding:14px"><div class="studio-section-title" style="font-size:16px">Gelistirilmis JSON gorunumu</div><textarea id="studio-abr-json" class="input" style="min-height:220px;font-family:Consolas,monospace">'+escHtml(state.json||'[]')+'</textarea><div class="studio-chip-row"><button class="btn btn-secondary" id="studio-abr-sync-json">JSON\'dan katmanlari guncelle</button></div></div>':'')+'</div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-card"><div class="studio-section-title">Canli test ve cikti tahmini</div><div class="studio-section-sub">Secilen profil ile beklenen teslimat yapisi.</div><div class="studio-chip-row"><span class="studio-chip active">HLS varyant: '+fmtInt(summary.variants)+'</span><span class="studio-chip">'+escHtml(summary.audioOnly?'Audio-only MPD':'Video + Audio MPD')+'</span><span class="studio-chip">'+escHtml((stream&&stream.output_formats)||'Tum cikislar')+'</span></div><div class="studio-code-block">'+escHtml(JSON.stringify({profile_set:state.profileSet,layers:state.layers,hls_master_variants:summary.variants,dash_representations:summary.audioOnly?summary.variants:(summary.variants+1),audio_only:summary.audioOnly},null,2))+'</div></div><div class="studio-card"><div class="studio-section-title">Audio-only DASH sertlestirme</div><div class="studio-section-sub">Radyo ve podcast senaryolari icin ses odakli teslimat tanisi.</div><div class="studio-chip-row"><span class="studio-chip '+((summary.audioOnly||(diagnostics.output_formats||'').indexOf('mp3')>=0)?'active':'')+'">Audio preset</span><span class="studio-chip '+(((diagnostics.dash_enabled)?'active':''))+'">DASH cikisi</span><span class="studio-chip '+(((diagnostics.checks||[]).find(function(item){ return item.code==='dash' && item.status==='ready'; })?'active':''))+'">MPD hazir</span></div><div class="form-hint">Sadece ses yayini icin Audio-only veya Radyo presetlerini sec; ardindan Embed Studyosu icinden DASH Ses veya HLS Ses linklerini kullan.</div>'+(stream?('<div style="margin-top:12px">'+copyField('DASH Ses',location.origin+'/audio/dash/'+stream.stream_key)+copyField('HLS Ses',location.origin+'/audio/hls/'+stream.stream_key)+'</div>'):'')+'</div></div>'+
      '</div>';

    document.querySelectorAll('.studio-option-card[data-studio-key]').forEach(function(btn){ if(btn.closest('#studio-embed-usecases')||btn.closest('#studio-embed-outputs')) return; });
    document.querySelectorAll('.studio-option-grid .studio-option-card').forEach(function(btn){ if(btn.closest('.studio-card.soft') && btn.closest('.studio-card.soft').querySelector('.studio-section-title') && btn.closest('.studio-card.soft').querySelector('.studio-section-title').textContent.indexOf('Hazir preset')>=0){ btn.onclick=function(){ const key=btn.dataset.studioKey||'balanced'; const preset=ABR_PRESETS[key]; if(!preset) return; state.profileSet=key; state.name=preset.title; state.description=preset.desc; state.layers=cloneJSON(preset.layers); state.mode=state.mode||'simple'; studioRerender('settings-abr'); }; } });
    document.querySelectorAll('[data-abr-mode]').forEach(function(btn){ btn.onclick=function(){ state.mode=btn.getAttribute('data-abr-mode')||'simple'; studioRerender('settings-abr'); }; });
    const streamSelect=document.getElementById('studio-abr-stream'); if(streamSelect) streamSelect.onchange=function(){ state.streamKey=streamSelect.value||''; studioRerender('settings-abr'); };
    const savedSelect=document.getElementById('studio-abr-saved'); if(savedSelect) savedSelect.onchange=function(){ state.profileId=toNumber(savedSelect.value,0); };
    const loadBtn=document.getElementById('studio-abr-load'); if(loadBtn) loadBtn.onclick=async function(){ if(!state.profileId){ toast('Once kayitli profil secin','warning'); return; } const item=await api('/api/admin/abr-profiles/'+state.profileId); if(!item || item.error){ toast('Profil yuklenemedi','error'); return; } state.profileSet=item.profile_set||'custom-profile'; state.name=item.name||''; state.description=item.description||''; state.scope=item.scope||'global'; state.streamKey=item.stream_key||state.streamKey; state.layers=parseJSONSafeStudio(item.profiles_json,[]).map(normalizeABRLayer); studioRerender('settings-abr'); };
    const duplicateBtn=document.getElementById('studio-abr-duplicate'); if(duplicateBtn) duplicateBtn.onclick=function(){ state.profileId=0; state.name=(state.name||'Profil')+' Kopya'; studioRerender('settings-abr'); };
    const addBtn=document.getElementById('studio-abr-add-layer'); if(addBtn) addBtn.onclick=function(){ var fresh=normalizeABRLayer({name:'Yeni katman',width:1280,height:720,bitrate:'2500k',max_bitrate:'3000k',buf_size:'5000k',fps:30,preset:'fast',audio_rate:'128k'},state.layers.length); fresh=applyABRResolutionChoice(fresh,'720p',state.layers.length); fresh=applyABRBitratePack(fresh,'balanced_720'); state.layers.push(fresh); studioRerender('settings-abr'); };
    const applyPreset=document.getElementById('studio-abr-apply-preset'); if(applyPreset) applyPreset.onclick=function(){ const preset=ABR_PRESETS[state.profileSet] || ABR_PRESETS.balanced; state.layers=cloneJSON(preset.layers); state.name=state.name||preset.title; state.description=state.description||preset.desc; studioRerender('settings-abr'); };
    const resetBtn=document.getElementById('studio-abr-reset'); if(resetBtn) resetBtn.onclick=function(){ window.abrStudioState=defaultABRState(); if(streamSelect) window.abrStudioState.streamKey=streamSelect.value||''; studioRerender('settings-abr'); };
    const exportBtn=document.getElementById('studio-abr-export'); if(exportBtn) exportBtn.onclick=function(){ studioDownloadFile('fluxstream-abr-'+(state.profileSet||'profile')+'.json',abrLayersToJSON(state.layers),'application/json;charset=utf-8'); };
    const importBtn=document.getElementById('studio-abr-import'); if(importBtn) importBtn.onclick=function(){ const raw=prompt('ABR JSON yapistirin'); if(!raw) return; try{ state.layers=parseJSONSafeStudio(raw,[]).map(normalizeABRLayer); studioRerender('settings-abr'); }catch(e){ toast('JSON okunamadi','error'); } };
    const syncJSON=document.getElementById('studio-abr-sync-json'); if(syncJSON) syncJSON.onclick=function(){ const raw=(document.getElementById('studio-abr-json')||{}).value||'[]'; state.layers=parseJSONSafeStudio(raw,[]).map(normalizeABRLayer); studioRerender('settings-abr'); };
    ['studio-abr-profile-set','studio-abr-name','studio-abr-description','studio-abr-scope'].forEach(function(id){ const el=document.getElementById(id); if(!el) return; el.onchange=function(){ if(id==='studio-abr-profile-set') state.profileSet=el.value||'custom-profile'; else if(id==='studio-abr-name') state.name=el.value||''; else if(id==='studio-abr-description') state.description=el.value||''; else if(id==='studio-abr-scope') state.scope=el.value||'global'; }; });
    document.querySelectorAll('[data-layer-field]').forEach(function(el){ el.onchange=function(){ const index=toNumber(el.getAttribute('data-layer-index'),0); const field=el.getAttribute('data-layer-field'); state.layers[index]=normalizeABRLayer(state.layers[index],index); state.layers[index][field]=el.value; state.json=abrLayersToJSON(state.layers); if(state.mode==='advanced') return; studioRerender('settings-abr'); }; });
    document.querySelectorAll('[data-layer-select="resolution"]').forEach(function(el){ el.onchange=function(){ const index=toNumber(el.getAttribute('data-layer-index'),0); state.layers[index]=normalizeABRLayer(state.layers[index],index); state.layers[index]=applyABRResolutionChoice(state.layers[index],el.value,index); state.json=abrLayersToJSON(state.layers); studioRerender('settings-abr'); }; });
    document.querySelectorAll('[data-layer-select="bitrate_pack"]').forEach(function(el){ el.onchange=function(){ const index=toNumber(el.getAttribute('data-layer-index'),0); state.layers[index]=normalizeABRLayer(state.layers[index],index); state.layers[index]=applyABRBitratePack(state.layers[index],el.value); state.layers[index]=normalizeABRLayer(state.layers[index],index); state.json=abrLayersToJSON(state.layers); studioRerender('settings-abr'); }; });
    document.querySelectorAll('[data-layer-up]').forEach(function(btn){ btn.onclick=function(){ const index=toNumber(btn.getAttribute('data-layer-up'),0); if(index<=0) return; const temp=state.layers[index-1]; state.layers[index-1]=state.layers[index]; state.layers[index]=temp; studioRerender('settings-abr'); }; });
    document.querySelectorAll('[data-layer-down]').forEach(function(btn){ btn.onclick=function(){ const index=toNumber(btn.getAttribute('data-layer-down'),0); if(index>=state.layers.length-1) return; const temp=state.layers[index+1]; state.layers[index+1]=state.layers[index]; state.layers[index]=temp; studioRerender('settings-abr'); }; });
    document.querySelectorAll('[data-layer-delete]').forEach(function(btn){ btn.onclick=function(){ const index=toNumber(btn.getAttribute('data-layer-delete'),0); state.layers.splice(index,1); if(!state.layers.length) state.layers=cloneJSON(ABR_PRESETS.balanced.layers); studioRerender('settings-abr'); }; });
    const applyCurrent=document.getElementById('studio-abr-apply-current'); if(applyCurrent) applyCurrent.onclick=async function(){ const payload={profile_set:state.profileSet||'custom-profile',profiles_json:abrLayersToJSON(state.layers),stream_key:state.streamKey||'',scope:state.scope||'global'}; const res=await api('/api/admin/abr-profiles/direct-apply',{method:'POST',body:payload}); if(!res || res.error){ toast('Profil uygulanamadi','error'); return; } toast('ABR profili uygulandi'); studioRerender('settings-abr'); };
    const saveBtn=document.getElementById('studio-abr-save'); if(saveBtn) saveBtn.onclick=async function(){ const payload={profile_set:state.profileSet||'custom-profile',name:state.name||state.profileSet||'Yeni profil',scope:state.scope||'global',stream_key:state.scope==='stream'?(state.streamKey||''):'',description:state.description||'',preset:state.profileSet||'',profiles_json:abrLayersToJSON(state.layers),summary_json:JSON.stringify(summarizeABRLayers(state.layers))}; const path=state.profileId?('/api/admin/abr-profiles/'+state.profileId):'/api/admin/abr-profiles'; const method=state.profileId?'PUT':'POST'; const res=await api(path,{method:method,body:payload}); if(!res || res.error){ toast('ABR profili kaydedilemedi','error'); return; } if(res.item && res.item.id) state.profileId=res.item.id; toast('ABR profili kaydedildi'); studioRerender('settings-abr'); };
    const deleteBtn=document.getElementById('studio-abr-delete'); if(deleteBtn) deleteBtn.onclick=async function(){ if(!state.profileId || !confirm('Secili kayitli ABR profili silinsin mi?')) return; const res=await api('/api/admin/abr-profiles/'+state.profileId,{method:'DELETE'}); if(!res || res.error){ toast('ABR profili silinemedi','error'); return; } state.profileId=0; toast('ABR profili silindi'); studioRerender('settings-abr'); };
  };

  function transcodePresetCards(current){
    const cards=[
      {key:'none',title:'CPU Modu',desc:'Ek GPU yoksa veya stabiliteyi onceliyorsan klasik encoder yolu.',badge:'Guvenli'},
      {key:'nvenc',title:'NVIDIA NVENC',desc:'Yuksek yogunluklu canli yayinlar icin GPU uzerinden hizli encode.',badge:'NVIDIA'},
      {key:'qsv',title:'Intel Quick Sync',desc:'Ofis ve mini PC sinifinda dusuk watt ile hizli donanim encode.',badge:'Intel'},
      {key:'amf',title:'AMD AMF',desc:'AMD GPU bulunan sistemlerde donanim hizlandirma.',badge:'AMD'}
    ];
    return cards.map(function(item){
      return '<button type="button" class="studio-option-card'+(item.key===current?' active':'')+'" data-transcode-gpu="'+escHtml(item.key)+'">'+
        '<div style="display:flex;align-items:center;justify-content:space-between;gap:8px;margin-bottom:8px"><h4 class="studio-option-title">'+escHtml(item.title)+'</h4><span class="studio-chip'+(item.key===current?' active':'')+'">'+escHtml(item.badge)+'</span></div>'+
        '<div class="studio-option-meta">'+escHtml(item.desc)+'</div></button>';
    }).join('');
  }

  window.renderSettingsTranscode = async function(container){
    const [settings,status,jobs]=await Promise.all([
      api('/api/settings'),
      api('/api/transcode/status'),
      api('/api/transcode/jobs')
    ]);
    const s=settings||{};
    const st=status||{};
    const list=Array.isArray(jobs)?jobs:[];
    const live=list.filter(function(item){ return item.status==='running'; });
    const gpu=String((s.gpu_accel||st.gpu_accel||'none')).toLowerCase();
    const ffmpegVersion=String(st.ffmpeg_version||'bilinmiyor');
    const liveOptions=st.live_options||{};
    const presets=(liveOptions.profiles||[]).map(function(item){ return {label:item.name||((item.width&&item.height)?(item.height+'p'):'audio'),value:item.bitrate||item.audio_rate||'-'}; });
    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero"><h1 class="studio-hero-title">Transkod / FFmpeg Merkezi</h1><div class="studio-hero-sub">FFmpeg yolu, GPU hizlandirma, canli ABR merdiveni ve calisan isler artik ayni ekranda. Teknik ayrintilar gorunur, ama kararlar kartlar ve oneriler uzerinden verilir.</div><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active">GPU: '+escHtml(gpu.toUpperCase())+'</span><span class="studio-pill">Aktif is: '+fmtInt(st.active_jobs||0)+'</span><span class="studio-pill">Profil: '+escHtml(liveOptions.profile_set||s.abr_profile_set||'balanced')+'</span></div></section>'+
        '<div class="studio-kpi-grid">'+
          renderAnalyticsKPI('Aktif Is',fmtInt(st.active_jobs||0),'FFmpeg ve canli repack islemleri')+
          renderAnalyticsKPI('Toplam Kayitli Is',fmtInt(st.total_jobs||list.length||0),'Durum kaydi tutulan tum isler')+
          renderAnalyticsKPI('Calisan HLS / DASH',fmtInt(live.length),'Su an calisan transcode isleri')+
          renderAnalyticsKPI('GPU Modu',gpu.toUpperCase(),'Secili hizlandirma katmani')+
          renderAnalyticsKPI('Segment',fmtStudioNumber(liveOptions.segment_duration||2,0)+' sn','Canli cikis parcasi')+
          renderAnalyticsKPI('Playlist',fmtInt(liveOptions.playlist_length||10),'Tutulan medya pencere boyu')+
          renderAnalyticsKPI('Audio Passthrough',liveOptions.audio_passthrough?'Acik':'Kapali','Ses dogrudan kopyalaniyor mu')+
          renderAnalyticsKPI('FFmpeg',ffmpegVersion.split(' ').slice(0,2).join(' ')||ffmpegVersion,'Calisan encoder surumu')+
        '</div>'+
        '<div class="studio-grid studio-grid-2">'+
          '<div class="studio-card soft"><div><h2 class="studio-section-title">Hizlandirma secimi</h2><div class="studio-section-sub">Sunucudaki encoder yolunu kartlardan sec, sonra alttaki ayari kaydet.</div></div><div class="studio-option-grid">'+transcodePresetCards(gpu)+'</div><div class="studio-grid studio-grid-2">'+
            studioField('FFmpeg yolu','<input id="studio-transcode-path" class="input" value="'+escHtml(s.ffmpeg_path||st.ffmpeg_path||'ffmpeg')+'">','Tam yol veya PATH icindeki ad yeterlidir.')+
            studioField('GPU modu','<select id="studio-transcode-gpu" class="input">'+studioSelectOptions([['none','CPU / Yok'],['nvenc','NVIDIA NVENC'],['qsv','Intel Quick Sync'],['amf','AMD AMF']],gpu)+'</select>','Donanim hizlandirma secimi.')+
          '</div><div class="studio-chip-row"><button class="btn btn-primary" id="studio-transcode-save">FFmpeg Ayarlarini Kaydet</button><button class="btn btn-secondary" id="studio-transcode-open-jobs">Transcode Islerini Ac</button></div></div>'+
          '<div class="studio-card"><div><h2 class="studio-section-title">Canli ABR ozet karti</h2><div class="studio-section-sub">Mevcut profile gore sunucunun uretecegi katmanlar ve maliyet ozeti.</div></div><div class="studio-chip-row">'+(presets.length?presets.map(function(item){ return '<span class="studio-chip active">'+escHtml(item.label)+' • '+escHtml(item.value)+'</span>'; }).join(''):'<span class="studio-chip">Canli katman bilgisi yok</span>')+'</div><div class="metric-list">'+
            '<div class="metric-row"><span>Profil seti</span><strong>'+escHtml(liveOptions.profile_set||s.abr_profile_set||'balanced')+'</strong></div>'+
            '<div class="metric-row"><span>ABR</span><strong>'+(liveOptions.abr_enabled?'Acik':'Kapali')+'</strong></div>'+
            '<div class="metric-row"><span>Master playlist</span><strong>'+(liveOptions.master_enabled?'Acik':'Kapali')+'</strong></div>'+
            '<div class="metric-row"><span>FFmpeg yolu</span><span class="mono-wrap">'+escHtml(st.ffmpeg_path||s.ffmpeg_path||'ffmpeg')+'</span></div>'+
          '</div><div class="studio-alert info"><strong>Oneri</strong><div style="margin-top:8px" class="form-hint">'+(gpu==='none'?'GPU yoksa Dayanikli veya Mobil profil daha guvenli olur.':'Donanim hizlandirma acik. Canli yogunlukte yine de segment ve bitrate merdivenini izlemek iyi olur.')+'</div></div></div>'+
        '</div>'+
        '<div class="studio-card"><div class="studio-section-title">Canli ve son isler</div><div class="studio-section-sub">Bu sayfa ayar merkezidir; ayrintili job izleme icin Transcode Isleri sayfasina gecebilirsin.</div>'+(list.length?'<div style="overflow:auto"><table class="studio-table"><thead><tr><th>ID</th><th>Stream</th><th>Tip</th><th>Durum</th><th>Baslangic</th><th>Cikti</th></tr></thead><tbody>'+list.slice(0,8).map(function(job){ return '<tr><td><code>'+escHtml(shortKey(job.id||'-'))+'</code></td><td><code>'+escHtml(shortKey(job.stream_key||'-'))+'</code></td><td>'+escHtml(job.type||'-')+'</td><td><span class="studio-chip'+(job.status==='running'?' active':'')+'">'+escHtml(job.status||'-')+'</span></td><td>'+escHtml(fmtLocaleDateTime(job.started_at))+'</td><td class="mono-wrap">'+escHtml(job.output_dir||'-')+'</td></tr>'; }).join('')+'</tbody></table></div>':'<div class="empty-state" style="padding:26px"><div class="icon"><i class="bi bi-list-task"></i></div><h3>Henuz transcode isi yok</h3><p style="color:var(--text-muted)">Canli yayin basladiginda burada FFmpeg ve cikis islerini goreceksin.</p></div>')+'</div>'+
      '</div>';
    document.querySelectorAll('[data-transcode-gpu]').forEach(function(btn){ btn.onclick=function(){ const key=btn.getAttribute('data-transcode-gpu')||'none'; const select=document.getElementById('studio-transcode-gpu'); if(select) select.value=key; }; });
    const save=document.getElementById('studio-transcode-save');
    if(save) save.onclick=async function(){
      const path=(document.getElementById('studio-transcode-path')||{}).value||'ffmpeg';
      const gpuValue=(document.getElementById('studio-transcode-gpu')||{}).value||'none';
      const res=await saveSettingsValues('transcode',{ffmpeg_path:path,gpu_accel:gpuValue},true);
      if(res && res.success!==false){
        toast('Transcode ayarlari kaydedildi');
        studioRerender('settings-transcode');
      }else{
        toast((res&&res.message)||'Transcode ayarlari kaydedilemedi','error');
      }
    };
    const openJobs=document.getElementById('studio-transcode-open-jobs');
    if(openJobs) openJobs.onclick=function(){ navigate('transcode-jobs'); };
    if(typeof schedulePageRefresh==='function') schedulePageRefresh('settings-transcode',15000);
  };

  window.renderViewers = async function(container){
    const [viewers,stats,bans]=await Promise.all([
      api('/api/viewers'),
      api('/api/stats/viewers'),
      api('/api/security/bans')
    ]);
    const data=viewers||{};
    const stat=stats||{};
    const sessions=Array.isArray(data.sessions)?data.sessions:[];
    const banList=Array.isArray(bans)?bans:[];
    const viewerTimeline=(Array.isArray(stat.timeline)?stat.timeline:[]).map(function(item){ return {timestamp:item.timestamp,value:Number(item.value||0)}; });
    const formatBars=Object.keys(stat.by_format||{}).map(function(key){ return {label:key.toUpperCase(),value:Number(stat.by_format[key]||0)}; }).sort(function(a,b){ return b.value-a.value; });
    const countryBars=Object.keys(stat.by_country||{}).map(function(key){ return {label:key,value:Number(stat.by_country[key]||0)}; }).sort(function(a,b){ return b.value-a.value; }).slice(0,8);
    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero"><h1 class="studio-hero-title">Izleyici Merkezi</h1><div class="studio-hero-sub">Canli izleyici sayisi, oturum dagilimi, format kullanimi ve IP yasaklari ayni panelde. Destek ve operasyon akisi icin hizli aksiyonlarla sade tutuldu.</div><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active">Canli: '+fmtInt(data.active||0)+'</span><span class="studio-pill">Toplam: '+fmtInt(data.total||0)+'</span><span class="studio-pill">Peak: '+fmtInt(stat.peak||0)+'</span></div></section>'+
        '<div class="studio-kpi-grid">'+
          renderAnalyticsKPI('Aktif Izleyici',fmtInt(data.active||0),'Su an acik baglantilar')+
          renderAnalyticsKPI('Toplam Izleyici',fmtInt(data.total||0),'Tum zamanlar toplami')+
          renderAnalyticsKPI('Tepe Izleyici',fmtInt(stat.peak||0),'Simdiye kadar gorulen en yuksek an')+
          renderAnalyticsKPI('Yasakli IP',fmtInt(data.banned||0),'Aktif ban kayitlari')+
          renderAnalyticsKPI('Format Cesidi',fmtInt(formatBars.length),'Ayni anda kullanilan teslimatlar')+
          renderAnalyticsKPI('Ulke Cesidi',fmtInt(countryBars.length),'Son pencere ulke dagilimi')+
        '</div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-chart-card"><div class="studio-section-title">Izleyici trendi</div><div class="studio-section-sub">Son snapshot penceresindeki izleyici zaman serisi.</div><div class="studio-chart-wrap">'+renderTimelineChart(viewerTimeline,'Izleyici verisi yok',function(v){ return fmtInt(v); },{meta:'Izleyici timeline'})+'</div></div><div class="studio-chart-card"><div class="studio-section-title">Format ve ulke dagilimi</div><div class="studio-section-sub">En cok kullanilan oynatim tipleri ve ulkeler.</div><div class="studio-grid" style="gap:14px"><div>'+renderBarList(formatBars,'Format verisi yok',function(v){ return fmtInt(v); })+'</div><div>'+renderBarList(countryBars,'Ulke verisi yok',function(v){ return fmtInt(v); })+'</div></div></div></div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-card"><div class="studio-section-title">Aktif oturumlar</div><div class="studio-section-sub">Canli izleyici oturumlari, format ve trafik bilgileriyle listelenir.</div>'+(sessions.length?'<div style="overflow:auto"><table class="studio-table"><thead><tr><th>Yayin</th><th>Format</th><th>IP</th><th>Ulke</th><th>Sure</th><th>Trafik</th><th>Son gorulme</th></tr></thead><tbody>'+sessions.map(function(sess){ return '<tr><td><strong>'+escHtml(sess.stream_name||shortKey(sess.stream_key||'-'))+'</strong><div class="form-hint"><code>'+escHtml(sess.stream_key||'-')+'</code></div></td><td><span class="studio-chip active">'+escHtml(String(sess.format||'-').toUpperCase())+'</span></td><td><code>'+escHtml(sess.ip||'-')+'</code></td><td>'+escHtml(sess.country||'-')+'</td><td>'+escHtml(formatDurationSeconds(sess.duration_seconds||0))+'</td><td>'+escHtml(fmtBytes(sess.bytes_sent||0))+'</td><td>'+escHtml(fmtLocaleDateTime(sess.last_seen))+'</td></tr>'; }).join('')+'</tbody></table></div>':'<div class="empty-state" style="padding:26px"><div class="icon"><i class="bi bi-people"></i></div><h3>Aktif oturum yok</h3><p style="color:var(--text-muted)">Su an canli izleyici baglantisi gorunmuyor.</p></div>')+'</div>'+
        '<div class="studio-card"><div class="studio-section-title">IP yasaklama</div><div class="studio-section-sub">Sorunlu istemcileri gecici veya kalici olarak engelle.</div><div class="studio-grid studio-grid-2">'+
          studioField('IP adresi','<input id="studio-ban-ip" class="input" placeholder="203.0.113.10">','Yasaklanacak istemci.')+
          studioField('Sure (dk)','<input id="studio-ban-dur" class="input" type="number" min="0" placeholder="60">','0 ise kalici sayilir.')+
          studioField('Neden','<input id="studio-ban-reason" class="input" placeholder="Yuksek trafik / kotuye kullanim">','Loglarda gorunecek not.')+
          studioField('Hizli not','<select id="studio-ban-preset" class="input">'+studioSelectOptions([['','Ozel neden'],['abuse','Kotuye kullanim'],['load','Asiri istek'],['geo','Bolge kisiti'],['ops','Operasyon karari']],'')+'</select>','Hizli neden secimi.')+
        '</div><div class="studio-chip-row"><button class="btn btn-danger" id="studio-ban-save">IP Yasakla</button><button class="btn btn-secondary" id="studio-ban-refresh">Listeyi Yenile</button></div>'+(banList.length?'<div style="overflow:auto;margin-top:14px"><table class="studio-table"><thead><tr><th>IP</th><th>Neden</th><th>Tarih</th><th>Islem</th></tr></thead><tbody>'+banList.map(function(item){ return '<tr><td><code>'+escHtml(item.IP||item.ip||'-')+'</code></td><td>'+escHtml(item.Reason||item.reason||'-')+'</td><td>'+escHtml(fmtLocaleDateTime(item.BannedAt||item.banned_at))+'</td><td><button class="btn btn-secondary btn-sm" data-unban-ip="'+escHtml(item.IP||item.ip||'')+'">Kaldir</button></td></tr>'; }).join('')+'</tbody></table></div>':'<div class="form-hint" style="margin-top:12px">Aktif yasakli IP kaydi yok.</div>')+'</div></div>'+
      '</div>';
    const preset=document.getElementById('studio-ban-preset');
    if(preset) preset.onchange=function(){
      const reason=document.getElementById('studio-ban-reason');
      const labels={abuse:'Kotuye kullanim',load:'Asiri istek',geo:'Bolge kisiti',ops:'Operasyon karari'};
      if(reason && labels[preset.value]) reason.value=labels[preset.value];
    };
    const save=document.getElementById('studio-ban-save');
    if(save) save.onclick=async function(){
      const ip=(document.getElementById('studio-ban-ip')||{}).value||'';
      const reason=(document.getElementById('studio-ban-reason')||{}).value||'Manuel';
      const duration=toNumber((document.getElementById('studio-ban-dur')||{}).value,0);
      if(!ip){ toast('Yasaklamak icin IP girin','warning'); return; }
      const res=await api('/api/security/bans',{method:'POST',body:{ip:ip,reason:reason,duration_minutes:duration}});
      if(res && res.status==='banned'){
        toast('IP yasaklandi');
        studioRerender('viewers');
      }else{
        toast((res&&res.message)||'IP yasaklanamadi','error');
      }
    };
    const refresh=document.getElementById('studio-ban-refresh');
    if(refresh) refresh.onclick=function(){ studioRerender('viewers'); };
    document.querySelectorAll('[data-unban-ip]').forEach(function(btn){ btn.onclick=async function(){ const ip=btn.getAttribute('data-unban-ip')||''; if(!ip) return; const res=await api('/api/security/bans',{method:'DELETE',body:{ip:ip}}); if(res && res.status==='unbanned'){ toast('IP engeli kaldirildi'); studioRerender('viewers'); } else { toast((res&&res.message)||'IP engeli kaldirilamadi','error'); } }; });
    if(typeof schedulePageRefresh==='function') schedulePageRefresh('viewers',10000);
  };

  window.renderTranscodeJobs = async function(container){
    const [status,jobs,streams]=await Promise.all([
      api('/api/transcode/status'),
      api('/api/transcode/jobs'),
      api('/api/streams')
    ]);
    const st=status||{};
    const list=Array.isArray(jobs)?jobs:[];
    const streamList=Array.isArray(streams)?streams:[];
    const streamNames={};
    streamList.forEach(function(item){ streamNames[item.stream_key]=item.name||item.stream_key; });
    const running=list.filter(function(item){ return item.status==='running'; });
    const completed=list.filter(function(item){ return item.status==='completed'; });
    const errored=list.filter(function(item){ return item.status==='error'; });
    const byType={};
    list.forEach(function(item){ byType[item.type||'other']=(byType[item.type||'other']||0)+1; });
    const typeBars=Object.keys(byType).map(function(key){ return {label:key,value:byType[key]}; }).sort(function(a,b){ return b.value-a.value; });
    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero"><h1 class="studio-hero-title">Transcode Isleri Merkezi</h1><div class="studio-hero-sub">Canli HLS, DASH ve remux islerinin durumunu, hata satirlarini ve cikti dizinlerini tek bakista izle. Teknik ayrinti gerekirken de okunakli kalacak sekilde duzenlendi.</div><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active">Aktif: '+fmtInt(st.active_jobs||running.length)+'</span><span class="studio-pill">GPU: '+escHtml(String(st.gpu_accel||'none').toUpperCase())+'</span><span class="studio-pill">FFmpeg: '+escHtml((String(st.ffmpeg_version||'bilinmiyor').split(' ').slice(0,2).join(' ')||String(st.ffmpeg_version||'bilinmiyor')) )+'</span></div></section>'+
        '<div class="studio-kpi-grid">'+
          renderAnalyticsKPI('Calisan Is',fmtInt(running.length),'Su an devam eden job')+
          renderAnalyticsKPI('Tamamlanan',fmtInt(completed.length),'Sorunsuz biten job')+
          renderAnalyticsKPI('Hatali',fmtInt(errored.length),'Elle bakilmasi gereken isler')+
          renderAnalyticsKPI('Toplam',fmtInt(list.length),'Kayitli transcode job satiri')+
          renderAnalyticsKPI('HLS / DASH',fmtInt(list.filter(function(item){ return item.type==='live_hls'||item.type==='live_dash'; }).length),'Canli cikis turleri')+
          renderAnalyticsKPI('Profil',escHtml((st.live_options&&st.live_options.profile_set)||'balanced'),'Varsayilan canli merdiven')+
        '</div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-chart-card"><div class="studio-section-title">Is tipi dagilimi</div><div class="studio-section-sub">Canli cikis turleri ve job tipi yogunlugu.</div>'+renderBarList(typeBars,'Job verisi yok',function(v){ return fmtInt(v); })+'</div><div class="studio-card"><div class="studio-section-title">FFmpeg ve ortam</div><div class="metric-list">'+
          '<div class="metric-row"><span>FFmpeg yolu</span><span class="mono-wrap">'+escHtml(st.ffmpeg_path||'-')+'</span></div>'+
          '<div class="metric-row"><span>OS / Mimari</span><strong>'+escHtml((st.os||'-')+' / '+(st.arch||'-'))+'</strong></div>'+
          '<div class="metric-row"><span>ABR profil seti</span><strong>'+escHtml((st.live_options&&st.live_options.profile_set)||'balanced')+'</strong></div>'+
          '<div class="metric-row"><span>Audio passthrough</span><strong>'+((st.live_options&&st.live_options.audio_passthrough)?'Acik':'Kapali')+'</strong></div>'+
          '<div class="metric-row"><span>Segment / Playlist</span><strong>'+fmtInt((st.live_options&&st.live_options.segment_duration)||2)+' sn / '+fmtInt((st.live_options&&st.live_options.playlist_length)||10)+'</strong></div>'+
        '</div><div class="studio-chip-row"><button class="btn btn-secondary" id="studio-jobs-open-settings">Transcode Ayarlarina Git</button><button class="btn btn-secondary" id="studio-jobs-refresh">Yenile</button></div></div></div>'+
        '<div class="studio-card"><div class="studio-section-title">Job listesi</div><div class="studio-section-sub">Hata varsa en basta gorunur. Ayrinti satirinda cikti dizini ve manifest yolu korunur.</div>'+(list.length?'<div style="overflow:auto"><table class="studio-table"><thead><tr><th>ID</th><th>Yayin</th><th>Tip</th><th>Durum</th><th>Baslangic</th><th>PID</th><th>Cikti</th><th>Hata</th></tr></thead><tbody>'+list.map(function(job){ const tone=job.status==='running'?' active':''; return '<tr><td><code>'+escHtml(shortKey(job.id||'-'))+'</code></td><td><strong>'+escHtml(streamNames[job.stream_key]||shortKey(job.stream_key||'-'))+'</strong><div class="form-hint"><code>'+escHtml(job.stream_key||'-')+'</code></div></td><td>'+escHtml(job.type||'-')+'</td><td><span class="studio-chip'+tone+'">'+escHtml(job.status||'-')+'</span></td><td>'+escHtml(fmtLocaleDateTime(job.started_at))+'</td><td>'+escHtml(String(job.pid||'-'))+'</td><td><div class="mono-wrap">'+escHtml(job.output_dir||'-')+'</div>'+(job.manifest_path?'<div class="form-hint mono-wrap">'+escHtml(job.manifest_path)+'</div>':'')+'</td><td style="max-width:320px;color:'+(job.error?'#dc2626':'var(--text-muted)')+'">'+escHtml(job.error||'-')+'</td></tr>'; }).join('')+'</tbody></table></div>':'<div class="empty-state" style="padding:26px"><div class="icon"><i class="bi bi-list-task"></i></div><h3>Job bulunmuyor</h3><p style="color:var(--text-muted)">Canli yayin veya donusum baslayinca bu liste dolacak.</p></div>')+'</div>'+
      '</div>';
    const openSettings=document.getElementById('studio-jobs-open-settings');
    if(openSettings) openSettings.onclick=function(){ navigate('settings-transcode'); };
    const refresh=document.getElementById('studio-jobs-refresh');
    if(refresh) refresh.onclick=function(){ studioRerender('transcode-jobs'); };
    if(typeof schedulePageRefresh==='function') schedulePageRefresh('transcode-jobs',10000);
  };

  function studioInsertAfter(referenceNode, html){
    if(!referenceNode || !referenceNode.parentNode || !html) return null;
    const wrap=document.createElement('div');
    wrap.innerHTML=html;
    const node=wrap.firstElementChild;
    if(!node) return null;
    referenceNode.parentNode.insertBefore(node, referenceNode.nextSibling);
    return node;
  }

  function renderBrandAssetTiles(items,currentURL){
    const list=Array.isArray(items)?items:[];
    if(!list.length) return '<div class="empty-state" style="padding:24px"><div class="icon"><i class="bi bi-images"></i></div><h3>Henüz asset yok</h3><p style="color:var(--text-muted)">Logo, poster veya marka görseli yüklediğinde burada görünecek.</p></div>';
    return '<div class="studio-asset-grid">'+list.map(function(item){
      const active=String(item.url||'')===String(currentURL||'');
      return '<div class="studio-asset-card'+(active?' active':'')+'">'+
        '<div class="studio-asset-thumb">'+(item.url?'<img src="'+escHtml(item.url)+'" alt="'+escHtml(item.name||'asset')+'">':'<span class="studio-chip">Asset</span>')+'</div>'+
        '<div><div style="font-weight:700;font-size:13px;word-break:break-word">'+escHtml(item.name||'asset')+'</div><div class="form-hint">'+escHtml((item.category||'branding').toUpperCase())+' • '+escHtml(fmtBytes(item.size||0))+'</div></div>'+
        '<div class="studio-inline-actions"><button class="btn btn-secondary btn-sm" data-asset-use="'+escHtml(item.url||'')+'">Kullan</button><button class="btn btn-secondary btn-sm" data-asset-copy="'+escHtml(item.url||'')+'">Kopyala</button><button class="btn btn-danger btn-sm" data-asset-delete="'+escHtml(item.path||'')+'">Sil</button></div>'+
      '</div>';
    }).join('')+'</div>';
  }

  async function bindBrandAssetTiles(root,onUse){
    if(!root) return;
    root.querySelectorAll('[data-asset-copy]').forEach(function(btn){
      btn.onclick=function(){ copyText(btn.getAttribute('data-asset-copy')||''); };
    });
    root.querySelectorAll('[data-asset-use]').forEach(function(btn){
      btn.onclick=function(){ if(typeof onUse==='function') onUse(btn.getAttribute('data-asset-use')||''); };
    });
    root.querySelectorAll('[data-asset-delete]').forEach(function(btn){
      btn.onclick=async function(){
        const path=btn.getAttribute('data-asset-delete')||'';
        if(!path || !confirm('Bu asset silinsin mi?')) return;
        const res=await studioDeleteAsset(path);
        if(res && !res.error){
          toast('Asset silindi');
          studioRerender(currentPage);
        }else{
          toast((res&&res.message)||'Asset silinemedi','error');
        }
      };
    });
  }

  function playerTemplateModalValues(){
    return {
      name:(document.getElementById('pt-name')||{}).value||'',
      theme:(document.getElementById('pt-theme')||{}).value||'dark',
      logo_url:(document.getElementById('pt-logo-url')||{}).value||'',
      logo_position:(document.getElementById('pt-logo-pos')||{}).value||'top-right',
      logo_opacity:parseFloat((document.getElementById('pt-logo-opacity')||{}).value||'1')||1,
      watermark_text:(document.getElementById('pt-watermark')||{}).value||'',
      show_title:!!((document.getElementById('pt-show-title')||{}).checked),
      show_live_badge:!!((document.getElementById('pt-show-badge')||{}).checked),
      background_css:(document.getElementById('pt-bg-css')||{}).value||'',
      control_bar_css:(document.getElementById('pt-ctrl-css')||{}).value||'',
      play_button_css:(document.getElementById('pt-play-css')||{}).value||'',
      custom_css:(document.getElementById('pt-custom-css')||{}).value||''
    };
  }

  async function refreshPlayerTemplateAssetShelf(currentURL){
    const host=document.getElementById('pt-brand-assets');
    if(!host) return;
    const items=await studioListAssets('branding');
    host.innerHTML=renderBrandAssetTiles(items,currentURL);
    bindBrandAssetTiles(host,function(url){
      const input=document.getElementById('pt-logo-url');
      if(input){
        input.value=url||'';
        updatePlayerTemplateModalPreview();
      }
    });
  }

  async function uploadPlayerTemplateLogo(){
    const input=document.getElementById('pt-logo-file');
    if(!input || !input.files || !input.files[0]){
      toast('Yüklenecek logo dosyasını seçin','warning');
      return;
    }
    const res=await studioUploadFile('branding',input.files[0]);
    if(res && res.item && res.item.url){
      const logo=document.getElementById('pt-logo-url');
      if(logo) logo.value=res.item.url;
      toast('Logo yüklendi');
      input.value='';
      await refreshPlayerTemplateAssetShelf(res.item.url);
      updatePlayerTemplateModalPreview();
    }else{
      toast((res&&res.message)||'Logo yüklenemedi','error');
    }
  }

  window.showPlayerModal = async function(id){
    let pt={name:'',theme:'dark',background_css:'',control_bar_css:'',play_button_css:'',logo_url:'',logo_position:'top-right',logo_opacity:1,watermark_text:'',show_title:true,show_live_badge:true,custom_css:''};
    if(id){
      const data=await api('/api/players/'+id);
      if(data && !data.error) pt=data;
    }
    const streams=await api('/api/streams');
    window._playerTemplateStreams=Array.isArray(streams)?streams:[];
    ensureTemplateStudioState();
    const modalRoot=document.getElementById('player-modal');
    if(!modalRoot) return;
    modalRoot.innerHTML=
      '<div class="modal-overlay" onclick="if(event.target===this)this.remove()">'+
        '<div class="modal" style="max-width:1280px">'+
          '<div class="modal-title">'+(id?'Player Şablonu Stüdyosu':'Yeni Player Şablonu')+'</div>'+
          '<div class="studio-grid studio-grid-2" style="align-items:start">'+
            '<div class="studio-card soft">'+
              '<div><h2 class="studio-section-title">Şablon Ayarları</h2><div class="studio-section-sub">Gerçek player görünümüne yakın önizleme ile renk, logo, watermark ve tema ayarlarını aynı pencerede yönet.</div></div>'+
              '<div class="studio-form-grid">'+
                studioField('Şablon Adı *','<input class="input" id="pt-name" value="'+escHtml(pt.name||'')+'" placeholder="Örn: Haber Merkezi">','Kaydedilecek şablonun görünen adı.')+
                studioField('Tema','<select class="input" id="pt-theme">'+studioSelectOptions([['dark','Dark'],['light','Light'],['minimal','Minimal'],['custom','Custom']],pt.theme||'dark')+'</select>','Hazır temadan başla veya custom kullan.')+
                studioField('Logo URL','<input class="input" id="pt-logo-url" value="'+escHtml(pt.logo_url||'')+'" placeholder="/media-assets/branding/logo.png veya https://...">','Harici URL ya da yüklediğin marka dosyası kullanılabilir.')+
                studioField('Logo Konumu','<select class="input" id="pt-logo-pos">'+studioSelectOptions([['top-right','Sağ Üst'],['top-left','Sol Üst'],['bottom-right','Sağ Alt'],['bottom-left','Sol Alt']],pt.logo_position||'top-right')+'</select>','Player içinde markanın görüleceği alan.')+
                studioField('Logo Şeffaflık','<input class="input" id="pt-logo-opacity" type="number" min="0" max="1" step="0.1" value="'+escHtml(String(pt.logo_opacity||1))+'">','0 ile 1 arasında şeffaflık değeri.')+
                studioField('Watermark Yazı','<input class="input" id="pt-watermark" value="'+escHtml(pt.watermark_text||'')+'" placeholder="CANLI • FluxStream">','Küçük metin watermarkı.')+
              '</div>'+
              '<div class="studio-option-grid">'+
                '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Başlık Göster</strong><div class="form-hint">Player başlığını görünür tutar.</div></div><input type="checkbox" id="pt-show-title" '+(pt.show_title?'checked':'')+'></div></label>'+
                '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>CANLI Rozeti</strong><div class="form-hint">Canlı durum etiketini gösterir.</div></div><input type="checkbox" id="pt-show-badge" '+(pt.show_live_badge?'checked':'')+'></div></label>'+
              '</div>'+
              '<div class="studio-grid">'+
                studioField('Arkaplan CSS','<textarea class="input" id="pt-bg-css" rows="3" placeholder="background:linear-gradient(180deg,#0b1120 0%,#111827 100%);">'+escHtml(pt.background_css||'')+'</textarea>','Ana player yüzeyi.')+
                studioField('Kontrol Çubuğu CSS','<textarea class="input" id="pt-ctrl-css" rows="3" placeholder="background:rgba(8,15,32,.88);">'+escHtml(pt.control_bar_css||'')+'</textarea>','Alt kontrol alanı.')+
                studioField('Play Butonu CSS','<textarea class="input" id="pt-play-css" rows="3" placeholder="background:#2563eb; color:#fff;">'+escHtml(pt.play_button_css||'')+'</textarea>','Ortadaki büyük play butonu.')+
                studioField('Özel CSS','<textarea class="input" id="pt-custom-css" rows="5" placeholder=".player-shell{backdrop-filter:blur(10px);}">'+escHtml(pt.custom_css||'')+'</textarea>','Tema dışında kalan küçük dokunuşlar.')+
              '</div>'+
              '<div class="studio-card soft" style="padding:14px">'+
                '<div class="studio-section-title" style="font-size:16px">Logo yükleme ve kütüphane</div>'+
                '<div class="studio-section-sub">Dış URL zorunluluğunu kaldırır. Yüklenen logo ve poster dosyaları tüm panelde tekrar kullanılabilir.</div>'+
                '<div class="studio-inline-actions" style="margin-top:10px"><input type="file" id="pt-logo-file" accept="image/*"><button class="btn btn-secondary" id="pt-logo-upload">Logo Yükle</button><button class="btn btn-secondary" id="pt-logo-refresh">Kütüphaneyi Yenile</button></div>'+
                '<div id="pt-brand-assets" style="margin-top:14px"></div>'+
              '</div>'+
              '<div class="studio-inline-actions"><button class="btn btn-secondary" onclick="document.getElementById(\'player-modal\').innerHTML=\'\'">Kapat</button><button class="btn btn-secondary" id="pt-apply-preview">Güncelleştir</button><button class="btn btn-primary" id="pt-save-btn">'+(id?'Kaydet ve Açık Kal':'Oluştur ve Açık Kal')+'</button></div>'+
            '</div>'+
            '<div class="studio-sticky-preview">'+
              '<div class="studio-card">'+
                '<div class="studio-section-title">Gerçek player önizlemesi</div>'+
                '<div class="studio-section-sub">Kaydet dediğinde modal kapanmaz; önizleme ve kodlar aynı pencerede güncellenir.</div>'+
                '<div class="studio-grid studio-grid-2">'+
                  studioField('Kaynak stream','<select class="input" id="pt-modal-stream">'+templateStudioStreamOptions()+'</select>','Önizleme kaynağı.')+
                  studioField('Format','<select class="input" id="pt-modal-format">'+templateStudioFormatOptions()+'</select>','Şablonun test formatı.')+
                '</div>'+
                '<input type="hidden" id="pt-current-template-id" value="'+(id||0)+'">'+
                '<div id="pt-live-preview" class="studio-preview-shell player-frame" style="margin-top:14px"></div>'+
              '</div>'+
              '<div class="studio-card" style="margin-top:14px"><div class="studio-section-title">Şablonlu çıkışlar</div><div class="studio-section-sub">Player URL, embed URL ve ana çıktı burada güncellenir.</div><div id="pt-live-embed-code"></div></div>'+
            '</div>'+
          '</div>'+
        '</div>'+
      '</div>';
    applyTranslations(modalRoot);
    ['pt-name','pt-theme','pt-logo-url','pt-logo-pos','pt-logo-opacity','pt-watermark','pt-bg-css','pt-ctrl-css','pt-play-css','pt-custom-css','pt-show-title','pt-show-badge','pt-modal-stream','pt-modal-format'].forEach(function(idValue){
      const el=document.getElementById(idValue);
      if(!el) return;
      el.onchange=function(){ updatePlayerTemplateModalPreview(); };
      el.oninput=function(){ updatePlayerTemplateModalPreview(); };
    });
    const upload=document.getElementById('pt-logo-upload');
    if(upload) upload.onclick=uploadPlayerTemplateLogo;
    const refresh=document.getElementById('pt-logo-refresh');
    if(refresh) refresh.onclick=function(){ refreshPlayerTemplateAssetShelf((document.getElementById('pt-logo-url')||{}).value||''); };
    const applyPreview=document.getElementById('pt-apply-preview');
    if(applyPreview) applyPreview.onclick=function(){ savePlayerTemplate(toNumber((document.getElementById('pt-current-template-id')||{}).value,0)||null); };
    const saveBtn=document.getElementById('pt-save-btn');
    if(saveBtn) saveBtn.onclick=function(){ savePlayerTemplate(toNumber((document.getElementById('pt-current-template-id')||{}).value,0)||null); };
    await refreshPlayerTemplateAssetShelf(pt.logo_url||'');
    updatePlayerTemplateModalPreview();
  };

  window.savePlayerTemplate = async function(id){
    const body=playerTemplateModalValues();
    if(!body.name){
      toast('Şablon adı gerekli','error');
      return;
    }
    let newID=id;
    if(id){
      const res=await api('/api/players/'+id,{method:'PUT',body:body});
      if(!res || res.error){
        toast((res&&res.message)||'Şablon güncellenemedi','error');
        return;
      }
      toast('Şablon güncellendi');
    }else{
      const res=await api('/api/players',{method:'POST',body:body});
      newID=toNumber((res&&res.id)||(res&&res.item&&res.item.id),0);
      if(!newID){
        toast((res&&res.message)||'Şablon oluşturulamadı','error');
        return;
      }
      const current=document.getElementById('pt-current-template-id');
      if(current) current.value=String(newID);
      toast('Şablon oluşturuldu');
    }
    await refreshPlayerTemplateAssetShelf(body.logo_url||'');
    updatePlayerTemplatePreview(newID);
  };

  window.renderPlayerTemplates = async function(container){
    const rendered=await studioRenderLegacy(container,'playerTemplates',{
      title:'Player Şablonları Stüdyosu',
      subtitle:'Gerçek player görünümüne yakın canlı önizleme, marka kütüphanesi, yüklenebilir logo dosyaları ve kapanmayan düzenleme akışı tek merkezde.',
      pills:[{label:'Gerçek önizleme',active:true},{label:'Upload destekli logo akışı'},{label:'Kaydet ve açık kal'}],
      actionsHTML:'<button class="btn btn-secondary btn-sm" onclick="navigate(\'logos\')"><i class="bi bi-images"></i> Marka Kütüphanesi</button><button class="btn btn-primary btn-sm" onclick="showPlayerModal()"><i class="bi bi-plus-circle"></i> Yeni Şablon</button>'
    });
    if(rendered){
      const host=container.querySelector('.studio-page');
      if(host){
        const summary=document.createElement('section');
        summary.className='studio-grid studio-grid-3';
        summary.innerHTML=
          '<div class="studio-summary"><span>Canlı önizleme</span><strong>Gerçek player</strong><div class="form-hint">Önizleme artık gerçek player tasarımını daha net gösterir.</div></div>'+
          '<div class="studio-summary"><span>Güncelleştir davranışı</span><strong>Kapanmadan</strong><div class="form-hint">Kaydet dediğinde modal kapanmaz; sonucu aynı ekranda görürsün.</div></div>'+
          '<div class="studio-summary"><span>Marka varlıkları</span><strong>Upload + URL</strong><div class="form-hint">Dış URL yanında doğrudan yüklenen logo ve posterleri de kullanabilirsin.</div></div>';
        host.insertBefore(summary, host.children[1] || null);
      }
      studioAuditControls(container);
    }
  };

  window.renderAdvancedEmbed = async function(container){
    const state=ensureEmbedState();
    state.mode='advanced';
    await renderEmbedStudio(container);
    const title=container.querySelector('.studio-hero-title');
    if(title) title.textContent='Gelişmiş Embed Stüdyosu';
    const sub=container.querySelector('.studio-hero-sub');
    if(sub) sub.textContent='Hazır paylaşım paketleri, güvenli link laboratuvarı, manifest testleri ve gerçek önizleme ile gelişmiş embed üretimi.';
    const studioPage=container.querySelector('.studio-page');
    const streamList=await api('/api/streams');
    const templates=await api('/api/players');
    const streams=Array.isArray(streamList)?streamList:[];
    const stream=streams.find(function(item){ return item.stream_key===state.streamKey; }) || streams[0] || null;
    const template=studioSelectedTemplate(Array.isArray(templates)?templates:[],state.templateId);
    const playback=stream?await studioFetchPlaybackBundle(state,stream,template):null;
    const extra=document.createElement('section');
    extra.className='studio-grid studio-grid-2';
    extra.innerHTML=
      '<div class="studio-card soft"><div><h2 class="studio-section-title">Embed preset laboratuvarı</h2><div class="studio-section-sub">Gizli yayın, token, audio-only, popup ve manifest dağıtımını ayrı ayrı test etmek için hızlı kısayollar.</div></div><div class="studio-chip-row">'+
      '<button class="btn btn-secondary" data-advanced-open="player">Player Aç</button><button class="btn btn-secondary" data-advanced-open="embed">Embed Aç</button><button class="btn btn-secondary" data-advanced-open="manifest">Manifest Aç</button><button class="btn btn-secondary" data-advanced-open="audio">Ses Çıkışını Aç</button><button class="btn btn-secondary" data-advanced-open="vlc">VLC Linki</button></div><div class="studio-alert info" style="margin-top:12px"><strong>Gizli preset davranışı</strong><div style="margin-top:8px" class="form-hint">Gizli veya tokenli paket seçmek artık varsayılan olarak stream policy’yi kalıcı kilitlemez. Bu ekranda geçici güvenli link üretir; kalıcı uygulama için ayrıca stream policy açmalısın.</div></div></div>'+
      '<div class="studio-card"><div><h2 class="studio-section-title">Doğrudan çıktı kutusu</h2><div class="studio-section-sub">Manifest, player, iframe ve ses URL’lerini burada tek tıkla test et.</div></div>'+
      copyField('Player URL',(playback&&playback.player_url)||'')+
      copyField('Embed URL',(playback&&playback.embed_url)||'')+
      copyField('Manifest URL',(playback&&playback.manifest_url)||'')+
      copyField('Ses URL',(playback&&playback.audio_url)||'')+
      '</div>';
    studioPage.appendChild(extra);
    extra.querySelectorAll('[data-advanced-open]').forEach(function(btn){
      btn.onclick=function(){
        if(!playback){ toast('Önce bir stream seçin','warning'); return; }
        const kind=btn.getAttribute('data-advanced-open');
        const map={player:playback.player_url,embed:playback.embed_url,manifest:playback.manifest_url,audio:playback.audio_url,vlc:playback.vlc_url};
        if(!map[kind]){ toast('Bu çıkış hazır değil','warning'); return; }
        window.open(map[kind],'_blank','noopener');
      };
    });
    studioAuditControls(container);
  };

  function defaultPlayerTemplateDraftStudioV2(){
    return {name:'',theme:'dark',background_css:'',control_bar_css:'',play_button_css:'',logo_url:'',logo_position:'top-right',logo_opacity:1,watermark_text:'',show_title:true,show_live_badge:true,custom_css:''};
  }

  function playerTemplateRecordToDraftStudioV2(item){
    return Object.assign(defaultPlayerTemplateDraftStudioV2(),cloneJSON(item||{}));
  }

  function ensurePlayerTemplateWorkbenchStudioV2(templates){
    if(!window.playerTemplateWorkbenchStateV2) window.playerTemplateWorkbenchStateV2={selectedId:0,draft:defaultPlayerTemplateDraftStudioV2()};
    const state=window.playerTemplateWorkbenchStateV2;
    const list=Array.isArray(templates)?templates:[];
    if(!state.draft) state.draft=defaultPlayerTemplateDraftStudioV2();
    if(state.selectedId){
      const selected=list.find(function(item){ return Number(item.id)===Number(state.selectedId); });
      if(selected && !state.draft.name) state.draft=playerTemplateRecordToDraftStudioV2(selected);
    }else if(list.length && !state.draft.name){
      state.selectedId=Number(list[0].id)||0;
      state.draft=playerTemplateRecordToDraftStudioV2(list[0]);
    }
    return state;
  }

  function loadPlayerTemplateWorkbenchStudioV2(id,templates){
    const state=ensurePlayerTemplateWorkbenchStudioV2(templates);
    const selected=(Array.isArray(templates)?templates:[]).find(function(item){ return Number(item.id)===Number(id); }) || null;
    state.selectedId=selected?Number(selected.id):0;
    state.draft=playerTemplateRecordToDraftStudioV2(selected);
    return state;
  }

  function playerTemplateWorkbenchValuesStudioV2(){
    return {
      name:(document.getElementById('ptw-name')||{}).value||'',
      theme:(document.getElementById('ptw-theme')||{}).value||'dark',
      logo_url:(document.getElementById('ptw-logo-url')||{}).value||'',
      logo_position:(document.getElementById('ptw-logo-pos')||{}).value||'top-right',
      logo_opacity:parseFloat((document.getElementById('ptw-logo-opacity')||{}).value||'1')||1,
      watermark_text:(document.getElementById('ptw-watermark')||{}).value||'',
      show_title:!!((document.getElementById('ptw-show-title')||{}).checked),
      show_live_badge:!!((document.getElementById('ptw-show-badge')||{}).checked),
      background_css:(document.getElementById('ptw-bg-css')||{}).value||'',
      control_bar_css:(document.getElementById('ptw-ctrl-css')||{}).value||'',
      play_button_css:(document.getElementById('ptw-play-css')||{}).value||'',
      custom_css:(document.getElementById('ptw-custom-css')||{}).value||''
    };
  }

  function syncPlayerTemplateWorkbenchStudioV2(){
    const state=window.playerTemplateWorkbenchStateV2||{selectedId:0,draft:defaultPlayerTemplateDraftStudioV2()};
    state.draft=playerTemplateWorkbenchValuesStudioV2();
    const streamSelect=document.getElementById('ptw-stream');
    const formatSelect=document.getElementById('ptw-format');
    if(streamSelect) playerTemplateStudioState.streamKey=streamSelect.value||'';
    if(formatSelect) playerTemplateStudioState.format=formatSelect.value||'player';
    ensureTemplateStudioState();
    return state.draft;
  }

  async function refreshPlayerTemplateWorkbenchAssetShelfStudioV2(currentURL){
    const host=document.getElementById('ptw-brand-assets');
    if(!host) return;
    const items=await studioListAssets('branding');
    host.innerHTML=renderBrandAssetTiles(items,currentURL);
    bindBrandAssetTiles(host,function(url){
      const input=document.getElementById('ptw-logo-url');
      if(input){
        input.value=url||'';
        syncPlayerTemplateWorkbenchStudioV2();
        renderPlayerTemplateWorkbenchPreviewStudioV2();
      }
    });
  }

  async function uploadPlayerTemplateWorkbenchLogoStudioV2(){
    const input=document.getElementById('ptw-logo-file');
    if(!input || !input.files || !input.files[0]){
      toast('Yuklenecek logo dosyasini secin','warning');
      return;
    }
    const res=await studioUploadFile('branding',input.files[0]);
    if(res && res.item && res.item.url){
      const logo=document.getElementById('ptw-logo-url');
      if(logo) logo.value=res.item.url;
      input.value='';
      syncPlayerTemplateWorkbenchStudioV2();
      await refreshPlayerTemplateWorkbenchAssetShelfStudioV2(res.item.url);
      await renderPlayerTemplateWorkbenchPreviewStudioV2();
      toast('Logo yuklendi');
    }else{
      toast((res&&res.message)||'Logo yuklenemedi','error');
    }
  }

  async function renderPlayerTemplateWorkbenchPreviewStudioV2(){
    const preview=document.getElementById('ptw-live-preview');
    const outputs=document.getElementById('ptw-live-output');
    const streamHint=document.getElementById('ptw-stream-hint');
    if(!preview || !outputs) return;
    const draft=syncPlayerTemplateWorkbenchStudioV2();
    const stream=ensureTemplateStudioState();
    if(streamHint) streamHint.innerHTML=stream?('Kaynak: <strong>'+escHtml(stream.name)+'</strong> • '+escHtml(stream.stream_key)):'Kaynak stream secilmedi';
    if(!stream){
      preview.innerHTML='<div class="empty-state" style="padding:34px"><div class="icon"><i class="bi bi-broadcast"></i></div><h3>Kaynak stream yok</h3><p style="color:var(--text-muted)">Canli onizleme icin once bir stream secin.</p></div>';
      outputs.innerHTML='';
      return;
    }
    const settings=await api('/api/settings');
    const access=await getPlaybackAccess(stream.stream_key,settings,stream.policy_json||'');
    const previewRawURLs=getPreviewURLs(stream.stream_key,settings,stream.name,access);
    const publicRawURLs=getAllURLs(stream.stream_key,settings,stream.name,access);
    const previewURLs=templateAwareURLs(previewRawURLs,draft,stream.name);
    const urls=templateAwareURLs(publicRawURLs,draft,stream.name);
    const format=playerTemplateStudioState.format||'player';
    const isAudio=format==='mp3'||format==='aac'||format==='ogg'||format==='wav'||format==='flac'||format==='icecast';
    const previewSrc=buildTemplatePreviewSrc(previewURLs,format);
    const bundle=buildEmbedBundle(format,stream.stream_key,urls,960,isAudio?140:540,true,true);
    preview.innerHTML=previewSrc
      ?'<div style="position:relative;'+(isAudio?'height:150px;':'padding-top:56.25%;')+'background:#05070b;border-radius:18px;overflow:hidden"><iframe src="'+previewSrc+'" style="position:absolute;inset:0;width:100%;height:100%;border:none;background:#000" allow="autoplay;fullscreen" allowfullscreen></iframe></div>'
      :'<div class="empty-state" style="padding:34px"><div class="icon"><i class="bi bi-palette"></i></div><h3>Onizleme hazir degil</h3><p style="color:var(--text-muted)">Secili format icin gecerli bir cikti olusunca burada gorunecek.</p></div>';
    outputs.innerHTML=
      '<div class="studio-grid studio-grid-2">'+
        studioField('Player URL','<input class="input" readonly value="'+escHtml(playerURLForFormat(urls.play||'',format))+'">','Template sorgusu eklenmis player baglantisi.')+
        studioField('Embed URL','<input class="input" readonly value="'+escHtml(urls.embed||'')+'">','Iframe ve sitelere gomulu cikti.')+
        studioField(bundle.primaryLabel||'Birincil cikti','<input class="input" readonly value="'+escHtml(bundle.primary||'')+'">','Secili formatin ana linki.')+
        studioField(bundle.directLabel||'Direkt link','<input class="input" readonly value="'+escHtml(bundle.direct||'')+'">','Direkt manifest veya medya baglantisi.')+
      '</div>'+
      (bundle.note?'<div class="studio-alert info" style="margin-top:12px"><strong>Format notu</strong><div style="margin-top:8px" class="form-hint">'+escHtml(bundle.note)+'</div></div>':'');
  }

  let playerTemplateWorkbenchPreviewTimer=null;
  function schedulePlayerTemplateWorkbenchPreviewStudioV2(){
    if(playerTemplateWorkbenchPreviewTimer) clearTimeout(playerTemplateWorkbenchPreviewTimer);
    playerTemplateWorkbenchPreviewTimer=setTimeout(function(){
      renderPlayerTemplateWorkbenchPreviewStudioV2();
    },180);
  }

  async function savePlayerTemplateWorkbenchStudioV2(){
    const state=window.playerTemplateWorkbenchStateV2||{selectedId:0,draft:defaultPlayerTemplateDraftStudioV2()};
    const body=syncPlayerTemplateWorkbenchStudioV2();
    if(!body.name){
      toast('Sablon adi gerekli','error');
      return;
    }
    if(state.selectedId){
      const res=await api('/api/players/'+state.selectedId,{method:'PUT',body:body});
      if(!res || res.error){
        toast((res&&res.message)||'Sablon guncellenemedi','error');
        return;
      }
      toast('Sablon guncellendi');
    }else{
      const res=await api('/api/players',{method:'POST',body:body});
      const newID=toNumber((res&&res.id)||(res&&res.item&&res.item.id),0);
      if(!newID){
        toast((res&&res.message)||'Sablon olusturulamadi','error');
        return;
      }
      state.selectedId=newID;
      toast('Sablon olusturuldu');
    }
    studioRerender('player-templates');
  }

  window.renderPlayerTemplates = async function(container){
    const [templates,streams]=await Promise.all([api('/api/players'),api('/api/streams')]);
    window._playerTemplateStreams=Array.isArray(streams)?streams:[];
    ensureTemplateStudioState();
    const library=Array.isArray(templates)?templates:[];
    const state=ensurePlayerTemplateWorkbenchStudioV2(library);
    const draft=playerTemplateRecordToDraftStudioV2(state.draft);
    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero"><h1 class="studio-hero-title">Player Sablonlari Studyosu</h1><div class="studio-hero-sub">Modala bagimli duzen yerine sabit kutuphane, canli taslak alani ve gercek player onizlemesiyle daha kullanisli bir sablon merkezi.</div><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active">Canli taslak</span><span class="studio-pill">Upload destekli logo</span><span class="studio-pill">Kapanmayan onizleme</span></div></section>'+
        '<section class="studio-toolbar"><div class="studio-toolbar-group"><select id="ptw-stream" class="input">'+templateStudioStreamOptions()+'</select><select id="ptw-format" class="input">'+templateStudioFormatOptions()+'</select></div><div class="studio-toolbar-group"><button class="btn btn-secondary" id="ptw-new">Yeni Taslak</button><button class="btn btn-secondary" id="ptw-duplicate">Cogalt</button><button class="btn btn-primary" id="ptw-save">Kaydet</button><button class="btn btn-danger" id="ptw-delete"'+(state.selectedId?'':' disabled')+'>Sil</button></div></section>'+
        '<section class="studio-template-workbench">'+
          '<aside class="studio-card soft"><div><h2 class="studio-section-title">Sablon kutuphanesi</h2><div class="studio-section-sub">Hazir taslaklardan birini sec, cogalt veya sifirdan basla.</div></div><div class="studio-template-library">'+
            '<button type="button" class="studio-template-tile'+(!state.selectedId?' active':'')+'" id="ptw-pick-new"><div class="title">Yeni taslak</div><div class="meta">Bos calisma alani acilir.</div></button>'+
            (library.length?library.map(function(item){ return '<button type="button" class="studio-template-tile'+(Number(item.id)===Number(state.selectedId)?' active':'')+'" data-template-select="'+Number(item.id)+'"><div class="title">'+escHtml(item.name||'Sablon')+'</div><div class="meta">'+escHtml((item.theme||'dark').toUpperCase())+' • '+escHtml(item.logo_url?'Logo var':'Logo yok')+'</div><div style="margin-top:10px">'+renderPlayerTemplateThumbnail(item)+'</div></button>'; }).join(''):'<div class="empty-state" style="padding:24px"><div class="icon"><i class="bi bi-pc-display-horizontal"></i></div><h3>Kayitli sablon yok</h3><p style="color:var(--text-muted)">Ilk taslagi sag taraftan olusturup kaydedebilirsin.</p></div>')+
          '</div></aside>'+
          '<section class="studio-card"><div><h2 class="studio-section-title">Canli taslak duzenleyici</h2><div class="studio-section-sub">Alanlar degistikce sagdaki player gercek gorunume yakin sekilde yenilenir. Kaydet demeden once sonucu gorebilirsin.</div></div>'+
            '<div class="studio-grid studio-grid-2">'+
              studioField('Sablon adi *','<input class="input" id="ptw-name" value="'+escHtml(draft.name||'')+'" placeholder="Ornek: Haber Merkezi">','Kutuphane ve embed baglantilarinda gorunen ad.')+
              studioField('Tema','<select class="input" id="ptw-theme">'+studioSelectOptions([['dark','Dark'],['light','Light'],['minimal','Minimal'],['custom','Custom']],draft.theme||'dark')+'</select>','Hazir tema ile baslayip sonradan ince ayar yapabilirsin.')+
              studioField('Logo URL','<input class="input" id="ptw-logo-url" value="'+escHtml(draft.logo_url||'')+'" placeholder="/media-assets/branding/logo.png veya https://...">','Dis URL veya yukledigin marka dosyasi kullanilabilir.')+
              studioField('Logo konumu','<select class="input" id="ptw-logo-pos">'+studioSelectOptions([['top-right','Sag ust'],['top-left','Sol ust'],['bottom-right','Sag alt'],['bottom-left','Sol alt']],draft.logo_position||'top-right')+'</select>','Player icinde logonun gorunecegi alan.')+
              studioField('Logo seffaflik','<input class="input" id="ptw-logo-opacity" type="number" min="0" max="1" step="0.1" value="'+escHtml(String(draft.logo_opacity||1))+'">','0 ile 1 arasinda yumusaklik degeri.')+
              studioField('Watermark yazi','<input class="input" id="ptw-watermark" value="'+escHtml(draft.watermark_text||'')+'" placeholder="CANLI • FluxStream">','Kucuk metin markalama.')+
            '</div>'+
            '<div class="studio-option-grid">'+
              '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Baslik goster</strong><div class="form-hint">Player basligini acik tutar.</div></div><input type="checkbox" id="ptw-show-title" '+(draft.show_title?'checked':'')+'></div></label>'+
              '<label class="card" style="padding:14px"><div style="display:flex;justify-content:space-between;gap:12px"><div><strong>Canli rozeti</strong><div class="form-hint">Canli yayin etiketini gosterir.</div></div><input type="checkbox" id="ptw-show-badge" '+(draft.show_live_badge?'checked':'')+'></div></label>'+
            '</div>'+
            '<div class="studio-grid">'+
              studioField('Arkaplan CSS','<textarea class="input" id="ptw-bg-css" rows="3" placeholder="background:linear-gradient(180deg,#0b1120 0%,#111827 100%);">'+escHtml(draft.background_css||'')+'</textarea>','Ana player yuzeyi.')+
              studioField('Kontrol cubugu CSS','<textarea class="input" id="ptw-ctrl-css" rows="3" placeholder="background:rgba(8,15,32,.88);">'+escHtml(draft.control_bar_css||'')+'</textarea>','Alt kontrol alaninin tonu.')+
              studioField('Play butonu CSS','<textarea class="input" id="ptw-play-css" rows="3" placeholder="background:#2563eb; color:#fff;">'+escHtml(draft.play_button_css||'')+'</textarea>','Ortadaki buyuk oynat butonu.')+
              studioField('Ozel CSS','<textarea class="input" id="ptw-custom-css" rows="5" placeholder=".player-shell{backdrop-filter:blur(10px);}">'+escHtml(draft.custom_css||'')+'</textarea>','Tema disi kucuk dokunuslar.')+
            '</div>'+
            '<div class="studio-card soft" style="padding:14px"><div class="studio-section-title" style="font-size:16px">Logo kutuphanesi</div><div class="studio-section-sub">Dis URL mecburiyeti yok. Dilersen bu sayfadan dosya yukle, sec ve aninda onizlemede gor.</div><div class="studio-inline-actions" style="margin-top:10px"><input type="file" id="ptw-logo-file" accept="image/*"><button class="btn btn-secondary" id="ptw-logo-upload">Logo Yukle</button><button class="btn btn-secondary" id="ptw-logo-refresh">Kutuphane Yenile</button><button class="btn btn-secondary" onclick="navigate(\'logos\')">Marka Merkezi</button></div><div id="ptw-brand-assets" style="margin-top:14px"></div></div>'+
          '</section>'+
          '<aside class="studio-sticky-preview"><div class="studio-card"><div class="studio-section-title">Gercek player onizlemesi</div><div class="studio-section-sub">Guncellestir dugmesine gerek kalmadan secili stream ve format ile canli sonucu gor.</div><div class="form-hint" id="ptw-stream-hint" style="margin-top:6px"></div><div id="ptw-live-preview" class="studio-preview-shell player-frame" style="margin-top:14px"></div></div><div class="studio-card" style="margin-top:14px"><div class="studio-section-title">Sablonlu ciktilar</div><div class="studio-section-sub">Player URL, embed URL ve direkt cikti baglantilari burada tutulur.</div><div id="ptw-live-output" style="margin-top:12px"></div></div></aside>'+
        '</section>'+
      '</div>';

    document.getElementById('ptw-stream').value=playerTemplateStudioState.streamKey||'';
    document.getElementById('ptw-format').value=playerTemplateStudioState.format||'player';
    ['ptw-name','ptw-theme','ptw-logo-url','ptw-logo-pos','ptw-logo-opacity','ptw-watermark','ptw-show-title','ptw-show-badge','ptw-bg-css','ptw-ctrl-css','ptw-play-css','ptw-custom-css','ptw-stream','ptw-format'].forEach(function(id){
      const el=document.getElementById(id);
      if(!el) return;
      const handler=function(){ syncPlayerTemplateWorkbenchStudioV2(); schedulePlayerTemplateWorkbenchPreviewStudioV2(); };
      el.onchange=handler;
      el.oninput=handler;
    });
    document.getElementById('ptw-pick-new').onclick=function(){
      window.playerTemplateWorkbenchStateV2={selectedId:0,draft:defaultPlayerTemplateDraftStudioV2()};
      studioRerender('player-templates');
    };
    document.querySelectorAll('[data-template-select]').forEach(function(btn){
      btn.onclick=function(){
        loadPlayerTemplateWorkbenchStudioV2(toNumber(btn.getAttribute('data-template-select'),0),library);
        studioRerender('player-templates');
      };
    });
    const newBtn=document.getElementById('ptw-new'); if(newBtn) newBtn.onclick=function(){ window.playerTemplateWorkbenchStateV2={selectedId:0,draft:defaultPlayerTemplateDraftStudioV2()}; studioRerender('player-templates'); };
    const duplicateBtn=document.getElementById('ptw-duplicate'); if(duplicateBtn) duplicateBtn.onclick=function(){ const draft=syncPlayerTemplateWorkbenchStudioV2(); window.playerTemplateWorkbenchStateV2={selectedId:0,draft:Object.assign({},draft,{name:(draft.name||'Yeni sablon')+' Kopya'})}; studioRerender('player-templates'); };
    const saveBtn=document.getElementById('ptw-save'); if(saveBtn) saveBtn.onclick=savePlayerTemplateWorkbenchStudioV2;
    const deleteBtn=document.getElementById('ptw-delete'); if(deleteBtn) deleteBtn.onclick=async function(){ const current=window.playerTemplateWorkbenchStateV2||{}; if(!current.selectedId || !confirm('Bu sablonu silmek istediginize emin misiniz?')) return; await api('/api/players/'+current.selectedId,{method:'DELETE'}); toast('Sablon silindi'); window.playerTemplateWorkbenchStateV2={selectedId:0,draft:defaultPlayerTemplateDraftStudioV2()}; studioRerender('player-templates'); };
    const uploadBtn=document.getElementById('ptw-logo-upload'); if(uploadBtn) uploadBtn.onclick=uploadPlayerTemplateWorkbenchLogoStudioV2;
    const refreshBtn=document.getElementById('ptw-logo-refresh'); if(refreshBtn) refreshBtn.onclick=function(){ refreshPlayerTemplateWorkbenchAssetShelfStudioV2((document.getElementById('ptw-logo-url')||{}).value||''); };
    await refreshPlayerTemplateWorkbenchAssetShelfStudioV2(draft.logo_url||'');
    await renderPlayerTemplateWorkbenchPreviewStudioV2();
    studioAuditControls(container);
  };

  window.renderAdvancedEmbed = async function(container){
    const rendered=await studioRenderLegacy(container,'advancedEmbed',{
      title:'Gelismis Embed Studyosu',
      subtitle:'Tum direkt linkler, output sekmeleri ve canli onizleme panelleri burada. HLS, DASH, MP4, ses cikislari ve protokol bazli teslimat linkleri eski guclu davranisiyla geri doner.',
      pills:[{label:'Tum outputlar',active:true},{label:'Direkt linkler'},{label:'Sekmeli onizleme'}],
      actionsHTML:'<button class="btn btn-secondary btn-sm" onclick="navigate(\'embed-codes\')"><i class="bi bi-grid"></i> Embed Studyosuna Don</button>'
    });
    if(rendered){
      const page=container.querySelector('.studio-page');
      const hero=page&&page.querySelector('.studio-hero');
      if(hero){
        studioInsertAfter(hero,
          '<section class="studio-grid studio-grid-3">'+
            '<div class="studio-summary"><span>Direkt linkler</span><strong>Tum teslimat ciktilari</strong><div class="form-hint">HLS, DASH, MP4, ses ve protokol bazli linkler ayni ekranda listelenir.</div></div>'+
            '<div class="studio-summary"><span>Guvenli erisim testi</span><strong>Token ve imzali link</strong><div class="form-hint">Korumali oynatim baglantilarini burada uretebilir ve yayina etkisini test edebilirsin.</div></div>'+
            '<div class="studio-summary"><span>Onizleme sekmeleri</span><strong>Tarayici ve harici oynatici</strong><div class="form-hint">Player, manifest ve medya linklerini ayni yerden dogrulayabilirsin.</div></div>'+
          '</section>'
        );
      }
      studioAuditControls(container);
    }
  };

  window.renderLogos = async function(container){
    const items=await studioListAssets('branding');
    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero"><h1 class="studio-hero-title">Logo ve Marka Merkezi</h1><div class="studio-hero-sub">Player şablonları, embed profilleri ve ileride gelecek marka presetleri için tekrar kullanılabilir logo, poster ve görsel varlık kütüphanesi.</div><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active">'+fmtInt(items.length)+' asset</span><span class="studio-pill">Upload + Kopyala + Sil</span><span class="studio-pill">Player ve Embed ile ortak</span></div></section>'+
        '<section class="studio-grid studio-grid-2">'+
          '<div class="studio-card soft"><div><h2 class="studio-section-title">Yeni asset yükle</h2><div class="studio-section-sub">Logo, poster veya küçük marka görseli yükleyebilirsin. Yüklenen URL tüm panelde kullanılabilir.</div></div>'+
            studioField('Kategori','<select id="studio-brand-category" class="input">'+studioSelectOptions([['branding','Genel Marka'],['logos','Logo'],['posters','Poster'],['players','Player Asset']], 'branding')+'</select>','Yüklenecek varlığın sınıfı.')+
            '<div class="studio-inline-actions"><input type="file" id="studio-brand-file" accept="image/*"><button class="btn btn-primary" id="studio-brand-upload">Yükle</button><button class="btn btn-secondary" id="studio-brand-refresh">Listeyi Yenile</button></div>'+
            '<div class="studio-alert info" style="margin-top:14px"><strong>Nerede kullanılır?</strong><div style="margin-top:8px" class="form-hint">Player Şablonları içindeki logo alanında, Embed Stüdyosu içindeki poster alanında ve ileride marka presetlerinde doğrudan kullanılabilir.</div></div>'+
          '</div>'+
          '<div class="studio-card"><div><h2 class="studio-section-title">Marka kullanım rehberi</h2><div class="studio-section-sub">Teknik URL’ler yerine bu kütüphaneden seçerek daha hızlı ilerleyebilirsin.</div></div><div class="metric-list">'+
            '<div class="metric-row"><span>Player şablonu</span><strong>Logo + watermark</strong></div>'+
            '<div class="metric-row"><span>Embed stüdyosu</span><strong>Poster + paylaşım paketi</strong></div>'+
            '<div class="metric-row"><span>Kısa kullanım</span><strong>Kopyala ve kullan</strong></div>'+
            '<div class="metric-row"><span>Öneri</span><strong>Şeffaf PNG veya SVG</strong></div>'+
          '</div><div class="studio-chip-row" style="margin-top:12px"><button class="btn btn-secondary" onclick="navigate(\'player-templates\')">Player Şablonlarına Git</button><button class="btn btn-secondary" onclick="navigate(\'embed-codes\')">Embed Stüdyosuna Git</button></div></div>'+
        '</section>'+
        '<section class="studio-card"><div><h2 class="studio-section-title">Asset kütüphanesi</h2><div class="studio-section-sub">Kullan, kopyala veya sil. En güncel dosyalar üstte görünür.</div></div><div id="studio-brand-assets">'+renderBrandAssetTiles(items,'')+'</div></section>'+
      '</div>';
    studioAuditControls(container);
    const upload=document.getElementById('studio-brand-upload');
    if(upload) upload.onclick=async function(){
      const file=(document.getElementById('studio-brand-file')||{}).files;
      const category=(document.getElementById('studio-brand-category')||{}).value||'branding';
      if(!file || !file[0]){
        toast('Önce yüklemek için bir dosya seçin','warning');
        return;
      }
      const res=await studioUploadFile(category,file[0]);
      if(res && res.item){
        toast('Asset yüklendi');
        studioRerender('logos');
      }else{
        toast((res&&res.message)||'Asset yüklenemedi','error');
      }
    };
    const refresh=document.getElementById('studio-brand-refresh');
    if(refresh) refresh.onclick=function(){ studioRerender('logos'); };
    await bindBrandAssetTiles(document.getElementById('studio-brand-assets'),function(url){ copyText(url); toast('Asset URL kopyalandi'); });
  };

  window.renderDashboard = async function(container){
    await studioRenderLegacy(container,'dashboard');
  };

  window.renderStreams = async function(container){
    await studioRenderLegacy(container,'streams');
  };

  window.renderGuidedSettings = async function(container){
    await studioRenderLegacy(container,'guidedSettings',{
      title:'Hızlı Ayarlar Stüdyosu',
      subtitle:'Günlük kullanım için risk seviyesi düşük, preset odaklı kısa yol ekranı.',
      pills:[{label:'Preset kartları',active:true},{label:'Hızlı güvenlik'},{label:'Yayın tipi senaryoları'}]
    });
  };

  window.renderSettingsGeneral = async function(container){
    await studioRenderLegacy(container,'settingsGeneral');
    const page=container.querySelector('.studio-page');
    if(page){
      const extra=document.createElement('section');
      extra.className='studio-grid studio-grid-3';
      extra.innerHTML=
        '<div class="studio-card soft"><div class="studio-section-title">Global varsayılanlar</div><div class="form-hint">Ürünün genel davranışını belirleyen ayarlar burada kalır. Marka, güvenlik ve teslimat varsayımlarını bu merkezden yönetebilirsin.</div></div>'+
        '<div class="studio-card"><div class="studio-section-title">Hızlı geçişler</div><div class="studio-chip-row"><button class="btn btn-secondary" onclick="navigate(\'logos\')">Logo ve Marka</button><button class="btn btn-secondary" onclick="navigate(\'settings-embed\')">Alan Adı / Embed</button><button class="btn btn-secondary" onclick="navigate(\'settings-security\')">Güvenlik</button></div></div>'+
        '<div class="studio-card"><div class="studio-section-title">Yönlendirme</div><div class="form-hint">Detaylı yayın davranışları için Protokoller, Çıkış Formatları, ABR ve Depolama merkezleriyle birlikte kullan.</div></div>';
      page.appendChild(extra);
    }
    studioAuditControls(container);
  };

  window.renderSettingsEmbed = async function(container){
    await studioRenderLegacy(container,'settingsEmbed');
    studioAuditControls(container);
  };

  window.renderSettingsProtocols = async function(container){
    await studioRenderLegacy(container,'settingsProtocols');
    const page=container.querySelector('.studio-page');
    if(page){
      const hero=page.querySelector('.studio-hero');
      studioInsertAfter(hero,'<section class="studio-grid studio-grid-3"><div class="studio-summary"><span>OBS / yayın encoder</span><strong>RTMP</strong><div class="form-hint">En kolay ve en yaygın giriş protokolü.</div></div><div class="studio-summary"><span>Kararsız ağlar</span><strong>SRT</strong><div class="form-hint">Düşük gecikme ve daha güvenli iletim.</div></div><div class="studio-summary"><span>Tarayıcı publish</span><strong>WHIP / WebRTC</strong><div class="form-hint">Doğrudan browser tabanlı gönderim için.</div></div></section>');
    }
    studioAuditControls(container);
  };

  window.renderSettingsOutputs = async function(container){
    await studioRenderLegacy(container,'settingsOutputs');
    const page=container.querySelector('.studio-page');
    if(page){
      const hero=page.querySelector('.studio-hero');
      studioInsertAfter(hero,'<section class="studio-grid studio-grid-3"><div class="studio-summary"><span>Genel web dağıtımı</span><strong>HLS + DASH</strong><div class="form-hint">En dengeli başlangıç kombinasyonu.</div></div><div class="studio-summary"><span>Düşük gecikme</span><strong>HTTP-FLV / WHEP</strong><div class="form-hint">Canlı etkileşim için hızlı teslimat.</div></div><div class="studio-summary"><span>Radyo / podcast</span><strong>HLS Ses / DASH Ses</strong><div class="form-hint">Audio-only yayınlar için en temiz yol.</div></div></section>');
    }
    studioAuditControls(container);
  };

  window.renderSettingsSecurity = async function(container){
    await studioRenderLegacy(container,'settingsSecurity');
    studioAuditControls(container);
  };

  window.renderSecurityTokens = async function(container){
    await studioRenderLegacy(container,'securityTokens',{
      title:'Token ve Signed URL Merkezi',
      subtitle:'Süreli playback linkleri, test akışı ve güvenlik profilleri için tek merkez.',
      pills:[{label:'Playback güvenliği',active:true},{label:'Kopyala + test et'},{label:'Embed ile uyumlu'}]
    });
  };

  window.renderMaintenanceCenter = async function(container){
    await studioRenderLegacy(container,'maintenanceCenter',{
      title:'Bakım ve Yedek Merkezi',
      subtitle:'Servis aksiyonları, upgrade planı ve offline restore komutları. Kayıt ve arşiv senkronu için Depolama ve Arşiv Merkezi ile rol ayrımı korunur.',
      pills:[{label:'Servis yönetimi',active:true},{label:'Offline restore'},{label:'Upgrade planı'}]
    });
  };

  window.renderDiagnostics = async function(container){
    const streams=await api('/api/streams')||[];
    const selectedID=(window.diagnosticsStudioState&&window.diagnosticsStudioState.streamID)||((streams[0]&&streams[0].id)||0);
    window.diagnosticsStudioState={streamID:selectedID};
    const selected=streams.find(function(item){ return Number(item.id)===Number(selectedID); }) || streams[0] || null;
    const diag=selected?await api('/api/diagnostics/stream/'+selected.id):null;
    const checks=Array.isArray(diag&&diag.checks)?diag.checks:[];
    const readyCount=checks.filter(function(item){ return item.status==='ready'; }).length;
    container.innerHTML=
      '<div class="studio-page">'+
        '<section class="studio-hero"><h1 class="studio-hero-title">Teşhis ve Tedavi Merkezi</h1><div class="studio-hero-sub">HLS, DASH, audio-only, güvenlik policy ve kayıt/arsiv zincirini denetle; gerektiğinde aynı ekrandan müdahale et.</div><div class="studio-pill-row" style="margin-top:14px"><span class="studio-pill active">'+(selected?escHtml(selected.name):'Yayın seç')+'</span><span class="studio-pill">'+fmtInt(readyCount)+' hazır kontrol</span><span class="studio-pill">'+fmtInt(checks.length-readyCount)+' dikkat</span></div></section>'+
        '<section class="studio-toolbar"><div class="studio-toolbar-group"><select id="studio-diag-stream" class="input">'+(streams.map(function(st){ return '<option value="'+Number(st.id)+'"'+(Number(st.id)===Number(selectedID)?' selected':'')+'>'+escHtml(st.name)+' • '+escHtml(st.stream_key)+'</option>'; }).join(''))+'</select></div><div class="studio-toolbar-group"><button class="btn btn-secondary" id="studio-diag-refresh">Yenile</button><button class="btn btn-secondary" id="studio-diag-open-ops">Operasyon Merkezi</button><button class="btn btn-primary" id="studio-diag-run-maint">Bakımı Çalıştır</button></div></section>'+
        '<div class="studio-kpi-grid">'+
          renderAnalyticsKPI('HLS Varyant',fmtInt((diag&&diag.hls_variant_count)||0),'Master playlist katmanı')+
          renderAnalyticsKPI('DASH Rep.',fmtInt((diag&&diag.dash_representation_count)||0),'MPD representation toplamı')+
          renderAnalyticsKPI('DASH Ses',fmtInt((diag&&diag.dash_audio_representation_count)||0),'Audio-only veya ses adaptation seti')+
          renderAnalyticsKPI('Aktif Player',fmtInt((diag&&diag.telemetry&&diag.telemetry.active_sessions)||0),'Anlık oturum')+
          renderAnalyticsKPI('Toplam Stall',fmtInt((diag&&diag.telemetry&&diag.telemetry.total_stalls)||0),'QoE takılma olayı')+
          renderAnalyticsKPI('Teslimat Özeti',escHtml((diag&&diag.delivery_summary&&diag.delivery_summary.label)||'-'),(diag&&diag.delivery_summary&&diag.delivery_summary.description)||'')+
        '</div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-card"><div class="studio-section-title">Kontrol matrisi</div><div class="studio-section-sub">Hazır, bekliyor, sorunlu veya kapalı durumlarını aynı tabloda gör.</div>'+(checks.length?'<div style="overflow:auto"><table class="studio-table"><thead><tr><th>Kontrol</th><th>Durum</th><th>Açıklama</th></tr></thead><tbody>'+checks.map(function(item){ return '<tr><td><strong>'+escHtml(item.description||item.code||'-')+'</strong></td><td><span class="studio-chip'+(item.status==='ready'?' active':'')+'">'+escHtml(item.label||item.status||'-')+'</span></td><td>'+escHtml(item.detail||item.message||'-')+'</td></tr>'; }).join('')+'</tbody></table></div>':'<div class="empty-state" style="padding:24px"><div class="icon"><i class="bi bi-heart-pulse"></i></div><h3>Henüz teşhis verisi yok</h3></div>')+'</div><div class="studio-card"><div class="studio-section-title">Tedavi aksiyonları</div><div class="studio-section-sub">Sorun çözmek için ilgili merkeze veya onarım akışına tek tık geçiş.</div><div class="studio-inline-actions"><button class="btn btn-secondary" id="studio-diag-open-security">Güvenlik</button><button class="btn btn-secondary" id="studio-diag-open-storage">Depolama</button><button class="btn btn-secondary" id="studio-diag-open-transcode">Transkod</button><button class="btn btn-secondary" id="studio-diag-open-text">Manifesti Aç</button>'+(selected?'<button class="btn btn-danger" id="studio-diag-reset-policy">Playback policy sıfırla</button>':'')+'</div><div class="studio-alert info" style="margin-top:14px"><strong>Tedavi mantığı</strong><div style="margin-top:8px" class="form-hint">Tanı ekranı, en sık kullanılan onarım akışlarını ilgili merkeze bağlar. Kayıt ve arşiv konuları için Depolama Merkezi; token ve signed URL sorunları için Güvenlik; varyant/FFmpeg sorunları için Transkod sayfasına geç.</div></div></div></div>'+
      '</div>';
    const streamSel=document.getElementById('studio-diag-stream');
    if(streamSel) streamSel.onchange=function(){ window.diagnosticsStudioState.streamID=toNumber(streamSel.value,0); studioRerender('diagnostics'); };
    const refresh=document.getElementById('studio-diag-refresh');
    if(refresh) refresh.onclick=function(){ studioRerender('diagnostics'); };
    const openOps=document.getElementById('studio-diag-open-ops');
    if(openOps) openOps.onclick=function(){ if(selected && typeof selectOperationsStream==='function') selectOperationsStream(selected.id||0); navigate('operations-center'); };
    const runMaint=document.getElementById('studio-diag-run-maint');
    if(runMaint) runMaint.onclick=function(){ runMaintenance(); };
    const openSec=document.getElementById('studio-diag-open-security');
    if(openSec) openSec.onclick=function(){ navigate('settings-security'); };
    const openStorage=document.getElementById('studio-diag-open-storage');
    if(openStorage) openStorage.onclick=function(){ navigate('settings-storage'); };
    const openTranscode=document.getElementById('studio-diag-open-transcode');
    if(openTranscode) openTranscode.onclick=function(){ navigate('settings-transcode'); };
    const openText=document.getElementById('studio-diag-open-text');
    if(openText) openText.onclick=function(){
      if(!selected){ toast('Önce bir yayın seçin','warning'); return; }
      openTextInspectModal('HLS Master',location.origin+'/hls/'+selected.stream_key+'/master.m3u8');
    };
    const reset=document.getElementById('studio-diag-reset-policy');
    if(reset) reset.onclick=async function(){
      if(!selected || !confirm('Bu yayın için kalıcı playback koruması sıfırlansın mı?')) return;
      const res=await api('/api/admin/security/stream-policy/reset',{method:'POST',body:{stream_key:selected.stream_key,clear_domain_lock:true,clear_ip_whitelist:false}});
      if(res && !res.error){
        toast('Playback policy sıfırlandı');
        studioRerender('diagnostics');
      }else{
        toast((res&&res.message)||'Policy sıfırlanamadı','error');
      }
    };
  };
})();
