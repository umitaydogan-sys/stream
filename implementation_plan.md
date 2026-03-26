# FluxStream Uygulama Plani

## 1. Urun Hedefi

FluxStream, tek binary calisan, yerelde ve Linux sunucuda kolay kurulabilen,
cok protokollu, urunlesmeye uygun bir canli yayin sunucusudur.

Ana hedefler:

- yayini guvenli sekilde almak
- farkli cikis formatlarinda dagitmak
- dusuk bant genisliginde bile akici oynatma saglamak
- adaptif bitrate ile kaliteyi baglanti kosullarina gore ayarlamak
- lisans, yedek, servis, kurulum ve guncelleme akisini urun seviyesine tasimak

Konferans, chat, sanal sinif ve anlik mesajlasma katmanlari,
streaming cekirdegi yeterince olgunlastiktan sonra eklenecek.

## 2. Bugun Itibariyla Cekirdekte Olanlar

### 2.1 Ingest ve Dagitim

- RTMP, RTMPS, SRT, RTP, RTSP, WebRTC/WHIP, MPEG-TS ve HTTP Push ingest aktif
- HLS, LL-HLS, DASH, HTTP-FLV, WHEP, MP4, WebM ve ses cikislari mevcut
- recording, analytics, subscriber fanout ve stream yasam dongusu calisiyor
- FFmpeg tabanli canli transcode ve ABR merdiveni devrede

### 2.2 Yonetim ve Urunlestirme

- admin paneli, setup wizard, stream CRUD, embed uretici ve player template sistemi hazir
- runtime lisans modeli ABR, RTMPS, recording ve branding tarafina bagli
- backup, restore plani, Linux servis yonetimi ve dagitim akisi var

### 2.3 Oynatici ve Canli Izleme

- preview, direct link ve embed ekranlari yeniden kararlilastirildi
- `play` ve `embed` ekranlari iframe icinde kullanilabilir durumda
- player template preview gercek gomulu player akisina cekildi
- onceki `403`, framing ve sahte `offline` problemleri kapatildi
- player tarafinda QoE debug overlay ve heartbeat tabanli telemetri var
- debug paneli artik aktif kaynak, fallback gecisi ve gecerli hata durumunu daha durust gosteriyor
- QoE verisi SQLite icinde kalici ornekler halinde saklanabiliyor
- admin stream detay ekraninda temel zaman serisi grafikler gosterilebiliyor
- admin panelde MPD/HLS manifestlerini ham metin olarak acabilen kullanim ve tanilama karti eklendi
- VLC icin onerilen HLS linki, tarayici player linki ve ham MP4 farki arayuzde netlestirildi
- menude ayri bir `Operasyon Merkezi` sayfasi eklendi
- tum streamleri secerek sekmeli canli izleme ve tanilama merkezi kullanilabiliyor
- stream arama alani yerine tum streamleri listeleyen merkezi bir selectbox eklendi
- secim katmani ileride on-demand playlistleri de ayni merkezden yonetebilecek sekilde hazirlandi
- `Genel Durum`, `Player ve Teslimat`, `QoE ve Telemetri`, `Track ve ABR`, `Manifest ve Ham Veri`, `OBS ve Ingest`, `Teshis` sekmeleri ayni merkezde bir araya getirildi
- `MP4 Player` ve `WebM` preview davranisi tarayici dostu HLS / DASH oncelikli hale getirildi; ham cikis URL'leri yine arayuzde korunuyor

### 2.4 OBS Cok Kanalli Video ve Adaptif Dagitim

- Enhanced RTMP multitrack ingest calisiyor
- OBS icin panelde kopyalanabilir `Config Override JSON` veriliyor
- stream olusturma ve stream detay ekranlarinda adim adim OBS rehberi var
- OBS multitrack video katmanlari HLS varyantlarina baglanabiliyor
- DASH repack, HLS master icindeki coklu varyantlari representation olarak mapleyebiliyor
- admin/API tarafinda canli video ve audio track listesi gorulebiliyor
- varsayilan video ve audio track secimi policy ve runtime seviyesinde uygulanabiliyor
- direct multitrack HLS master artik alternate audio group uretebiliyor
- video player icinde canli audio track secici gorunebiliyor
- player audio secimi artik tarayici oturumunda kalici tercih olarak saklanabiliyor ve HLS/DASH fallback'lerinde yeniden uygulanabiliyor
- audio-only HLS yonlendirmesi secili audio track playlistine gidebiliyor
- track bitrate ve runtime ornekleri kalici analytics olarak SQLite'a yazilabiliyor
- `/metrics` ve `/api/observability/otel` cikislari hazir
- QoE retention, esik tabanli uyarilar ve saglik raporu entegrasyonu aktif
- Operasyon Merkezi teshis bolumu artik opsiyonel cikislari gereksiz yere kirmizi gostermiyor
- `Hazir / Bekliyor / Kapali / Opsiyonel / Sorunlu` ayrimi ile daha durust teshis mantigi eklendi
- DASH preview tarafinda daha stabil buffer, yeniden deneme ve gec fallback mantigi eklendi
- player telemetrisi kalite gecisi ve audio track degisimi sayaclarini da tasiyor
- stream detay ve Operasyon Merkezi ekranlari artik kalite gecisi, ses izi degisimi ve secili audio track dagilimini gosterebiliyor
- QoE uyari esikleri aktif oturum oranina gore daha akilli hesaplanabiliyor
- saglik ekraninda QoE riskli yayinlar, kalite gecisi ve ses gecisi yogunlugu daha derin raporlanabiliyor
- dusuk bant icin `resilient` ABR profil seti eklendi, `balanced` ve `mobile` merdivenleri daha korumaci hale getirildi
- canli dogrulamada HLS master 2 video katmani, DASH MPD ise 2 video + 1 audio representation ile dogrulandi
- recording kutuphanesi icin object storage / archive yonetimi eklendi
- lokal arsiv klasoru ve S3/MinIO uyumlu archive akisi ayni panelden yonetilebilir hale geldi
- arsivlenen kayitlari panelden geri yukleme ve otomatik senkron mantigi eklendi
- `Depolama ve Arsiv Merkezi` ile kayit yonetimi, depolama ayarlari ve sistem yedekleri ayni ekranda birlestirildi
- harici hedef seceneklerine `SFTP` de eklendi; ayni merkezden `local`, `S3`, `MinIO`, `SFTP` secilebiliyor
- recording tarafinda varsayilan format `mp4` oldu; `mp4` ve `mkv` secildiginde yayin once guvenli `TS capture` olarak alinip kapanista `ffmpeg copy remux` ile son dosyaya donusturuluyor
- tum kayit baslatma ekranlari ve hizli kayit endpoint'leri `mp4` varsayilanina cekildi; kullanici isterse `TS`, `MKV` veya `FLV` secebiliyor
- kayit kutuphanesi artik gecici `.capture.ts` dosyalarini gostermiyor
- mevcut `TS`, `FLV` ve `MKV` kayitlari panelden tek tusla `MP4 Hazirla` akisi ile izlenebilir formata donusturulebiliyor
- sistem yedekleri icin de ayni archive altyapisi kullaniliyor; otomatik yukleme, geri yukleme ve lokal kopyayi silme politikasi ayarlanabiliyor
- `Depolama ve Arsiv Merkezi` ekranindaki buton aksiyonlari artik tam sayfa yeniden cizmeden calisiyor; renderer crash zinciri kapatildi
- `MP4 Hazirla` akisinin durumu arka plan isi olarak izlenebiliyor; sayfa degisse de remux isi devam ediyor
- recording tarafinda TS capture paketleme mantigi HLS ile hizalandi; AVC payload Annex-B, AAC payload ADTS olarak yaziliyor
- yeni kayitlar artik ilk gecerli video keyframe'inden baslatilarak daha guvenilir MP4 remux kaynagi uretiyor
- sistem yedegi silme ve kayit / arsiv aksiyonlari ayni ekranda daha guvenli parcali yenileme ile calisiyor
- MinIO / S3 upload yolunda eksik `Content-Length` basligi giderildi; MinIO artik `411 MissingContentLength` vermeden gercek upload kabul ediyor
- ayni VPS uzerinde kurulan `MinIO` ve ayri `SFTP` kullanicisi ile recording + backup upload / restore saha testi basariyla tamamlandi

## 3. Bu Surecte Neleri Kapatmis Olduk

### 3.1 Player, Preview ve Embed Tarafi

- preview iframe engelleri kapatildi
- direct link ve embed senaryolari yeniden duzeltildi
- player template preview, gercek oynatici gorunumuyle hizalandi
- sahte `offline` ve eski framing problemleri kapatildi

### 3.2 OBS Cok Kanalli Video Destegi

- Enhanced RTMP multitrack paketleri parse edilir hale getirildi
- `trackId` bilgisi ingest tarafinda alinmaya baslandi
- OBS icin panelde hazir JSON ve kurulum rehberi eklendi
- multitrack video katmanlari HLS master playlist tarafina baglandi
- DASH tarafi HLS varyantlarindan representation uretebilir hale geldi

### 3.3 QoE ve Teshis Katmani

- player heartbeat telemetrisi eklendi
- stall, toparlanma, reconnect, buffer ve hata bilgisi runtime olarak toplanmaya baslandi
- admin stream detay ekranina canli `QoE ve Stall Telemetrisi` karti eklendi
- diagnostics ekranina `HLS varyant sayisi` ve `DASH representation sayisi` eklendi
- telemetri ornekleri SQLite icinde kalici olarak saklanmaya baslandi
- admin stream detay ekranina temel zaman serisi grafikler eklendi
- stream detay ekraninda canli track listesi ve varsayilan secim karti eklendi
- kalite gecisi ve ses izi degisimi artik player oturum bazinda raporlanabiliyor
- telemetry history icinde kalite gecisi ve audio switch trendleri de birikiyor

### 3.4 Log ve Metin Tarafi

- log ekraninda yeni kayitlar icin metin normalize etme katmani eklendi
- API hata mesajlarinda metin normalize etme uygulandi
- eski DASH hata metninin fallback sonrasi gereksiz yere kalmasi buyuk olcude temizlendi
- debug panelinde artik aktif kaynak ve fallback gecisi ayri alanlarda izlenebiliyor

## 4. Canli Testte Bulunan Kok Neden ve Kalici Cozum

Canli VPS testlerinde `1080p` OBS multitrack varyantinda gorulen
mikro segment probleminin asil kok nedeni RTMP chunk reader tarafinda bulundu.

Kok neden:

- RTMP Type 1 ve Type 2 chunk header'larindaki `timestamp delta`
  degeri mutlak zaman gibi ele aliniyordu
- ayni CSID uzerinden akan yeni mesajlarda delta birikmedigi icin
  bazi paketler `11ms`, `16ms` gibi yanlis zaman damgalari aliyordu
- OBS multitrack yayininda butun video katmanlari ayni video CSID
  uzerinden geldigi icin sorun ozellikle ust kalite izlerde buyuyordu
- HLS ve DASH segment sureleri bu bozuk timestamp farklarindan
  hesaplandiginda `0.010`, `0.011`, `0.016` gibi mikro segmentler olusuyordu

Yapilan kalici duzeltmeler:

- `internal/ingest/rtmp/chunk.go`
  icinde CSID bazli mutlak zaman ve delta birikimi eklendi
- Type 1 ve Type 2 chunk'lar artik onceki mutlak zamana delta ekleyerek
  gercek timestamp uretiyor
- Type 3 yeni mesaj algisi, ayni delta'yi yeni mesaja dogru sekilde tasiyor
- `internal/output/hls/muxer.go`
  icinde wall-clock savunma katmani eklendi
- `internal/output/dash/muxer.go`
  icinde ayni duration fallback mantigi uygulandi
- `internal/transcode/live_multitrack.go`
  icinde master playlist tekrar tum saglikli varyantlari ilan edecek hale getirildi

Sonuc:

- mikro `EXTINF` segmentleri ortadan kalkti
- DASH `SegmentTimeline` degerleri tutarli hale geldi
- `master.m3u8` icinde `360p` ve `1080p` katmanlari yeniden saglikli adaptif bitrate
  olarak gorunebilir hale geldi
- onceki gecici guvenli mod, hata ayiklama asamasinda kullanilan ara adim olarak kaldi

## 5. Bugun Acik Kalan Gercek Eksikler

Bugun artik asagidaki maddeler kapanmis sayilmali:

- OBS multitrack baglantisini kabul etme
- multitrack video katmanlarini HLS varyantlarina baglama
- mikro segment kok nedenini kapatma
- HLS ve DASH tarafinda temel adaptif dagitimi ayaga kaldirma
- admin tarafinda QoE ve stall gorunurlugu saglama
- admin/API tarafinda track metadata ve track listesi acma
- varsayilan video/audio track secimini policy ve runtime seviyesinde urunlestirme
- QoE telemetrisini kalici depoya alma ve temel grafiklerle gosterebilme

Bugun hala acik olan gercek eksikler ise sunlar:

- `audio-only DASH` icin eklenen `audio.mpd` ve audio-only init segment akisinin farkli oynaticilarla saha testi
- kalite gecisi ve audio switch verisini daha ileri alarm otomasyonu ve uzun periyot rapor katmanina baglama
- ABR profil merdivenlerini gercek trafik ve cihaz verisine gore tekrar ince ayarlama
- dusuk bant sahalarinda uzun sureli soak test ve canli benchmark calistirma
- harici bir bucket ile gercek AWS S3 saha testi
- ayni VPS uzerindeki MinIO + SFTP laboratuvar hedeflerinde yapilan testleri daha uzun sureli senaryolarla sertlestirme
- kayit finalize/remux akisinda buyuk dosya, uzun sureli kayit ve beklenmeyen servis restart senaryolarini sertlestirme
- onceki bozuk `TS` kayitlar icin kullaniciya acik kurtarma / uyari akislarini tasarlama
- `Depolama ve Arsiv Merkezi` ekranini teknik terimleri azaltarak daha sade, daha anlasilir hale getirme
- `Google Drive` ve `OneDrive` gibi populer cloud hedefleri icin daha basit arsiv baglanti secenekleri ekleme
- playback guvenligini signed URL, oturum bagli token, hotlink korumasi ve watermark ile urunlestirme
- dusuk butceye uygun playback security fazi ile tam DRM fazini ayri planlama

## 6. Rakiplere Gore Bugunku Konum

FluxStream'in bugun guclu oldugu taraflar:

- tek binary ile kolay kurulum
- ayni urunde admin paneli + setup wizard + stream CRUD + embed + template akisi
- zengin ingest ve output matrisi
- runtime lisans, backup ve Linux servis omurgasi
- OBS cok kanalli video icin panel destekli kullanim rehberi
- artik gercek OBS multitrack to ABR omurgasinin calisiyor olmasi

FluxStream'in bugun zayif veya eksik oldugu taraflar:

- multi-node origin-edge cluster ve autoscaling yok
- archive/object storage ve SFTP hedefleri artik var ama saha sertlestirmesi henuz yeni
- Prometheus / OpenTelemetry / alarm omurgasi artik var ama daha derin entegrasyon acik
- RBAC, audit log ve SSO eksik
- DRM, SSAI ve gelismis monetizasyon eksik
- otomatik test, yuk testi ve uzun sureli soak test kapsami dar

## 7. Bugunku Uretim Degerlendirmesi

25 Mart 2026 itibariyla FluxStream artik yalnizca prototip degildir.

Bugun icin en dogru tanim:

- iyi bir tek-node medya sunucusu
- urunlesmis bir yayin cekirdegi
- canli testte kendini gostermis bir HLS merkezli dagitim sistemi

Bu haliyle su alanlar icin ciddi bicimde kullanilabilir durumdadir:

- kurum ici TV
- radyo
- webcast
- webinar
- markali player ve embed dagitimi

Ama su haliyle henuz "ust seviye enterprise rakiplerle ayni ligde"
demek icin erken:

- cluster
- storage
- telemetry
- security
- audit
- monetization

## 8. Simdi Ne Yapmaliyiz

Bir sonraki dogru kapatma sirasi bence su:

1. `Depolama ve Arsiv Merkezi` ekranini sade, daha az teknik ve daha kullanici dostu hale getir
2. harici bucket ile gercek S3 saha testi al
3. ayni VPS uzerindeki MinIO + SFTP laboratuvar hedeflerini UI/UX ve uzun sureli senaryolarla yeniden dogrula
4. audio-only DASH ve uzun sureli recording/finalize davranisini canli saha testiyle sertlestir
5. dusuk butceye uygun playback guvenligi fazini ekle:
   signed URL, signed manifest/segment, oturum tokeni, hotlink korumasi, watermark
6. daha sonra tam DRM fazini tasarla:
   AES-128 HLS key servis, DRM abstraction, Widevine/FairPlay/PlayReady hazirligi
7. sonra multi-node origin-edge mimarisini tasarla
8. RBAC, audit log, SSO ve lisans enforcement tarafini sertlestir
9. ABR profil merdivenini gercek saha benchmarklari ile tekrar optimize et

## 9. Operasyon Merkezi Fazinin Uygulama Taslagi

Bu faz, stream detay ekranindaki yeni kartlari kaldirmadan daha guclu bir
operasyon yuzeyi uretmeyi hedefler.

Menu ve sayfa:

- menu adi: `Operasyon Merkezi`
- sayfa basligi: `Canli Izleme ve Tanilama Merkezi`
- amac: tum streamleri tek merkezden secip, canli teslimat, telemetry,
  track, manifest ve OBS bilgisini ayni yerden gormek

Sayfa mimarisi:

- sol kolon: stream listesi, arama ve filtreler
- orta kolon: secilen stream icin sekmeli detay paneli
- sag kolon: hizli linkler ve operator eylemleri

Sekmeler:

- `Genel Durum`
- `Player ve Teslimat`
- `QoE ve Telemetri`
- `Track ve ABR`
- `Manifest ve Ham Veri`
- `OBS ve Ingest`
- `Teshis`

Bu fazda kullanilacak mevcut veri kaynaklari:

- `/api/streams`
- `/api/streams/{id}`
- `/api/admin/player/telemetry/stream/{id}`
- `/api/diagnostics/stream/{id}`
- `/api/settings`
- `/hls/{key}/master.m3u8`
- `/hls/{key}/index.m3u8`
- `/dash/{key}/manifest.mpd`

Bu fazin ayni teslimat paketi icindeki ikinci zorunlu isi:

- `MP4 preview fix`
- ham MP4 cikisini korurken, tarayici preview ve panel butonlarini
  kullaniciya daha durust ve kararli sekilde sunmak
- `MP4 Player` davranisini operasyon paneli, gelismis embed ve rehber
  metinleriyle tutarli hale getirmek

## 10. Cekirdek Sonrasi Buyuk Faz

Streaming cekirdegi yeterince olgunlastiginda siradaki urun katmanlari:

- konferans odalari
- canli chat
- moderasyonlu soru-cevap
- sanal sinif rolleri
- yoklama ve katilim
- breakout room mantigi
- takim ici mesajlasma

