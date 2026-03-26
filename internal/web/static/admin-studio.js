(function(){
  if(!window.api || !window.escHtml) return;

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

  function toNumber(value,fallback){
    const n=Number(value);
    return Number.isFinite(n)?n:fallback;
  }

  function studioRerender(page){
    if(typeof loadPage==='function') return loadPage(page);
  }

  function studioField(label,control,hint){
    return '<div><label class="form-label">'+escHtml(label)+'</label>'+control+(hint?'<div class="form-hint">'+escHtml(hint)+'</div>':'')+'</div>';
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
        preset.security.applyStreamPolicy=true;
        preset.options.sharePackage='private';
        break;
      case 'token':
        preset.outputType='iframe';
        preset.security.tokenRequired=true;
        preset.security.signedURL=true;
        preset.security.applyStreamPolicy=true;
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
        signed_url:!!state.security.signedURL,
        token_required:!!state.security.tokenRequired,
        session_bound:!!state.security.sessionBound,
        expiry_minutes:toNumber(state.security.expiryMinutes,60),
        domain_restriction:state.security.domainRestriction||'',
        ip_restriction:state.security.ipRestriction||'',
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
          '<div class="studio-card"><div><h2 class="studio-section-title">Playback guvenligi v1</h2><div class="studio-section-sub">Signed URL, sureli token, domain ve IP kisiti, iframe baglamasi ve watermark tek merkezde yonetilir.</div></div><div class="studio-option-grid">'+
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

  function renderABRLayerRow(layer,index,total){
    layer=normalizeABRLayer(layer,index);
    return '<div class="studio-layer" data-layer-index="'+index+'"><div class="studio-layer-head"><strong>'+escHtml(layer.name||('Katman '+(index+1)))+'</strong><div style="display:flex;gap:8px;flex-wrap:wrap"><button class="btn btn-secondary btn-sm" data-layer-up="'+index+'"'+(index===0?' disabled':'')+'>Yukari</button><button class="btn btn-secondary btn-sm" data-layer-down="'+index+'"'+(index===total-1?' disabled':'')+'>Asagi</button><button class="btn btn-danger btn-sm" data-layer-delete="'+index+'">Sil</button></div></div><div class="studio-grid studio-grid-3">'+
      studioField('Ad','<input class="input" data-layer-field="name" data-layer-index="'+index+'" value="'+escHtml(layer.name||'')+'">','Katman etiketi.')+
      studioField('Genislik','<input class="input" type="number" min="0" data-layer-field="width" data-layer-index="'+index+'" value="'+escHtml(String(layer.width||0))+'">','0 ise audio-only kabul edilir.')+
      studioField('Yukseklik','<input class="input" type="number" min="0" data-layer-field="height" data-layer-index="'+index+'" value="'+escHtml(String(layer.height||0))+'">','0 ise audio-only kabul edilir.')+
      studioField('Bitrate','<input class="input" data-layer-field="bitrate" data-layer-index="'+index+'" value="'+escHtml(layer.bitrate||'0')+'">','Ornek: 2500k')+
      studioField('Max bitrate','<input class="input" data-layer-field="max_bitrate" data-layer-index="'+index+'" value="'+escHtml(layer.max_bitrate||layer.bitrate||'0')+'">','ABR ust siniri.')+
      studioField('Buffer','<input class="input" data-layer-field="buf_size" data-layer-index="'+index+'" value="'+escHtml(layer.buf_size||'0')+'">','Encoder buffer boyutu.')+
      studioField('FPS','<input class="input" type="number" min="0" data-layer-field="fps" data-layer-index="'+index+'" value="'+escHtml(String(layer.fps||0))+'">','Audio-only icin 0 olabilir.')+
      studioField('Preset','<select class="input" data-layer-field="preset" data-layer-index="'+index+'">'+studioSelectOptions([['copy','Copy'],['superfast','superfast'],['veryfast','veryfast'],['fast','fast'],['faster','faster']],layer.preset||'fast')+'</select>','Encoder hiz profili.')+
      studioField('Audio bitrate','<input class="input" data-layer-field="audio_rate" data-layer-index="'+index+'" value="'+escHtml(layer.audio_rate||'128k')+'">','Katmana eslik eden ses hizi.')+
    '</div></div>';
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
        '<div class="studio-card"><div><h2 class="studio-section-title">Katman olusturucu</h2><div class="studio-section-sub">Katman ekle, sil, yukari/asagi tasi. Surukle birak yerine daha guvenli kontrol butonlari kullaniliyor.</div></div><div style="display:flex;gap:10px;flex-wrap:wrap"><button class="btn btn-primary" id="studio-abr-add-layer">Katman Ekle</button><button class="btn btn-secondary" id="studio-abr-apply-preset">Preseti Yukle</button><button class="btn btn-secondary" id="studio-abr-reset">Sifirla</button></div><div id="studio-abr-layer-list">'+(state.layers||[]).map(function(layer,index,arr){ return renderABRLayerRow(layer,index,arr.length); }).join('')+'</div>'+(state.mode==='advanced'?'<div class="studio-card soft" style="padding:14px"><div class="studio-section-title" style="font-size:16px">Gelistirilmis JSON gorunumu</div><textarea id="studio-abr-json" class="input" style="min-height:220px;font-family:Consolas,monospace">'+escHtml(state.json||'[]')+'</textarea><div class="studio-chip-row"><button class="btn btn-secondary" id="studio-abr-sync-json">JSON\'dan katmanlari guncelle</button></div></div>':'')+'</div>'+
        '<div class="studio-grid studio-grid-2"><div class="studio-card"><div class="studio-section-title">Canli test ve cikti tahmini</div><div class="studio-section-sub">Secilen profil ile beklenen teslimat yapisi.</div><div class="studio-chip-row"><span class="studio-chip active">HLS varyant: '+fmtInt(summary.variants)+'</span><span class="studio-chip">'+escHtml(summary.audioOnly?'Audio-only MPD':'Video + Audio MPD')+'</span><span class="studio-chip">'+escHtml((stream&&stream.output_formats)||'Tum cikislar')+'</span></div><div class="studio-code-block">'+escHtml(JSON.stringify({profile_set:state.profileSet,layers:state.layers,hls_master_variants:summary.variants,dash_representations:summary.audioOnly?summary.variants:(summary.variants+1),audio_only:summary.audioOnly},null,2))+'</div></div><div class="studio-card"><div class="studio-section-title">Audio-only DASH sertlestirme</div><div class="studio-section-sub">Radyo ve podcast senaryolari icin ses odakli teslimat tanisi.</div><div class="studio-chip-row"><span class="studio-chip '+((summary.audioOnly||(diagnostics.output_formats||'').indexOf('mp3')>=0)?'active':'')+'">Audio preset</span><span class="studio-chip '+(((diagnostics.dash_enabled)?'active':''))+'">DASH cikisi</span><span class="studio-chip '+(((diagnostics.checks||[]).find(function(item){ return item.code==='dash' && item.status==='ready'; })?'active':''))+'">MPD hazir</span></div><div class="form-hint">Sadece ses yayini icin Audio-only veya Radyo presetlerini sec; ardindan Embed Studyosu icinden DASH Ses veya HLS Ses linklerini kullan.</div>'+(stream?('<div style="margin-top:12px">'+copyField('DASH Ses',location.origin+'/audio/dash/'+stream.stream_key)+copyField('HLS Ses',location.origin+'/audio/hls/'+stream.stream_key)+'</div>'):'')+'</div></div>'+
      '</div>';

    document.querySelectorAll('.studio-option-card[data-studio-key]').forEach(function(btn){ if(btn.closest('#studio-embed-usecases')||btn.closest('#studio-embed-outputs')) return; });
    document.querySelectorAll('.studio-option-grid .studio-option-card').forEach(function(btn){ if(btn.closest('.studio-card.soft') && btn.closest('.studio-card.soft').querySelector('.studio-section-title') && btn.closest('.studio-card.soft').querySelector('.studio-section-title').textContent.indexOf('Hazir preset')>=0){ btn.onclick=function(){ const key=btn.dataset.studioKey||'balanced'; const preset=ABR_PRESETS[key]; if(!preset) return; state.profileSet=key; state.name=preset.title; state.description=preset.desc; state.layers=cloneJSON(preset.layers); state.mode=state.mode||'simple'; studioRerender('settings-abr'); }; } });
    document.querySelectorAll('[data-abr-mode]').forEach(function(btn){ btn.onclick=function(){ state.mode=btn.getAttribute('data-abr-mode')||'simple'; studioRerender('settings-abr'); }; });
    const streamSelect=document.getElementById('studio-abr-stream'); if(streamSelect) streamSelect.onchange=function(){ state.streamKey=streamSelect.value||''; studioRerender('settings-abr'); };
    const savedSelect=document.getElementById('studio-abr-saved'); if(savedSelect) savedSelect.onchange=function(){ state.profileId=toNumber(savedSelect.value,0); };
    const loadBtn=document.getElementById('studio-abr-load'); if(loadBtn) loadBtn.onclick=async function(){ if(!state.profileId){ toast('Once kayitli profil secin','warning'); return; } const item=await api('/api/admin/abr-profiles/'+state.profileId); if(!item || item.error){ toast('Profil yuklenemedi','error'); return; } state.profileSet=item.profile_set||'custom-profile'; state.name=item.name||''; state.description=item.description||''; state.scope=item.scope||'global'; state.streamKey=item.stream_key||state.streamKey; state.layers=parseJSONSafeStudio(item.profiles_json,[]).map(normalizeABRLayer); studioRerender('settings-abr'); };
    const duplicateBtn=document.getElementById('studio-abr-duplicate'); if(duplicateBtn) duplicateBtn.onclick=function(){ state.profileId=0; state.name=(state.name||'Profil')+' Kopya'; studioRerender('settings-abr'); };
    const addBtn=document.getElementById('studio-abr-add-layer'); if(addBtn) addBtn.onclick=function(){ state.layers.push(normalizeABRLayer({name:'Yeni katman',width:640,height:360,bitrate:'600k',max_bitrate:'700k',buf_size:'1200k',fps:24,preset:'fast',audio_rate:'64k'},state.layers.length)); studioRerender('settings-abr'); };
    const applyPreset=document.getElementById('studio-abr-apply-preset'); if(applyPreset) applyPreset.onclick=function(){ const preset=ABR_PRESETS[state.profileSet] || ABR_PRESETS.balanced; state.layers=cloneJSON(preset.layers); state.name=state.name||preset.title; state.description=state.description||preset.desc; studioRerender('settings-abr'); };
    const resetBtn=document.getElementById('studio-abr-reset'); if(resetBtn) resetBtn.onclick=function(){ window.abrStudioState=defaultABRState(); if(streamSelect) window.abrStudioState.streamKey=streamSelect.value||''; studioRerender('settings-abr'); };
    const exportBtn=document.getElementById('studio-abr-export'); if(exportBtn) exportBtn.onclick=function(){ studioDownloadFile('fluxstream-abr-'+(state.profileSet||'profile')+'.json',abrLayersToJSON(state.layers),'application/json;charset=utf-8'); };
    const importBtn=document.getElementById('studio-abr-import'); if(importBtn) importBtn.onclick=function(){ const raw=prompt('ABR JSON yapistirin'); if(!raw) return; try{ state.layers=parseJSONSafeStudio(raw,[]).map(normalizeABRLayer); studioRerender('settings-abr'); }catch(e){ toast('JSON okunamadi','error'); } };
    const syncJSON=document.getElementById('studio-abr-sync-json'); if(syncJSON) syncJSON.onclick=function(){ const raw=(document.getElementById('studio-abr-json')||{}).value||'[]'; state.layers=parseJSONSafeStudio(raw,[]).map(normalizeABRLayer); studioRerender('settings-abr'); };
    ['studio-abr-profile-set','studio-abr-name','studio-abr-description','studio-abr-scope'].forEach(function(id){ const el=document.getElementById(id); if(!el) return; el.onchange=function(){ if(id==='studio-abr-profile-set') state.profileSet=el.value||'custom-profile'; else if(id==='studio-abr-name') state.name=el.value||''; else if(id==='studio-abr-description') state.description=el.value||''; else if(id==='studio-abr-scope') state.scope=el.value||'global'; }; });
    document.querySelectorAll('[data-layer-field]').forEach(function(el){ el.onchange=function(){ const index=toNumber(el.getAttribute('data-layer-index'),0); const field=el.getAttribute('data-layer-field'); state.layers[index]=normalizeABRLayer(state.layers[index],index); state.layers[index][field]=el.value; state.json=abrLayersToJSON(state.layers); }; });
    document.querySelectorAll('[data-layer-up]').forEach(function(btn){ btn.onclick=function(){ const index=toNumber(btn.getAttribute('data-layer-up'),0); if(index<=0) return; const temp=state.layers[index-1]; state.layers[index-1]=state.layers[index]; state.layers[index]=temp; studioRerender('settings-abr'); }; });
    document.querySelectorAll('[data-layer-down]').forEach(function(btn){ btn.onclick=function(){ const index=toNumber(btn.getAttribute('data-layer-down'),0); if(index>=state.layers.length-1) return; const temp=state.layers[index+1]; state.layers[index+1]=state.layers[index]; state.layers[index]=temp; studioRerender('settings-abr'); }; });
    document.querySelectorAll('[data-layer-delete]').forEach(function(btn){ btn.onclick=function(){ const index=toNumber(btn.getAttribute('data-layer-delete'),0); state.layers.splice(index,1); if(!state.layers.length) state.layers=cloneJSON(ABR_PRESETS.balanced.layers); studioRerender('settings-abr'); }; });
    const applyCurrent=document.getElementById('studio-abr-apply-current'); if(applyCurrent) applyCurrent.onclick=async function(){ const payload={profile_set:state.profileSet||'custom-profile',profiles_json:abrLayersToJSON(state.layers),stream_key:state.streamKey||'',scope:state.scope||'global'}; const res=await api('/api/admin/abr-profiles/direct-apply',{method:'POST',body:payload}); if(!res || res.error){ toast('Profil uygulanamadi','error'); return; } toast('ABR profili uygulandi'); studioRerender('settings-abr'); };
    const saveBtn=document.getElementById('studio-abr-save'); if(saveBtn) saveBtn.onclick=async function(){ const payload={profile_set:state.profileSet||'custom-profile',name:state.name||state.profileSet||'Yeni profil',scope:state.scope||'global',stream_key:state.scope==='stream'?(state.streamKey||''):'',description:state.description||'',preset:state.profileSet||'',profiles_json:abrLayersToJSON(state.layers),summary_json:JSON.stringify(summarizeABRLayers(state.layers))}; const path=state.profileId?('/api/admin/abr-profiles/'+state.profileId):'/api/admin/abr-profiles'; const method=state.profileId?'PUT':'POST'; const res=await api(path,{method:method,body:payload}); if(!res || res.error){ toast('ABR profili kaydedilemedi','error'); return; } if(res.item && res.item.id) state.profileId=res.item.id; toast('ABR profili kaydedildi'); studioRerender('settings-abr'); };
    const deleteBtn=document.getElementById('studio-abr-delete'); if(deleteBtn) deleteBtn.onclick=async function(){ if(!state.profileId || !confirm('Secili kayitli ABR profili silinsin mi?')) return; const res=await api('/api/admin/abr-profiles/'+state.profileId,{method:'DELETE'}); if(!res || res.error){ toast('ABR profili silinemedi','error'); return; } state.profileId=0; toast('ABR profili silindi'); studioRerender('settings-abr'); };
  };
})();
