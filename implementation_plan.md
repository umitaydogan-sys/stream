# FluxStream Uygulama Plani

Tarih: 26 Mart 2026

## 1. Urun Vizyonu

FluxStream, tek binary ile calisan, yerelde ve Linux sunucuda kolay
kurulabilen, cok protokollu ingest alip HLS/DASH merkezli dagitim
yapabilen, urunlesmeye uygun bir canli yayin sunucusudur.

Ana hedefler:

- yayini guvenli ve kararlı sekilde almak
- dusuk bant genisliginde dahi akici izleme saglamak
- adaptif bitrate ile kaliteyi baglanti kosullarina gore evirmek
- kayit, arsiv, yedek ve operasyon akislarini urun seviyesine tasimak
- playback guvenligi, lisans ve Linux urunlestirmesini cekirdege entegre etmek

Konferans, chat, sanal sinif ve mesajlasma katmanlari cekirdek
streaming omurgasi yeterince olgunlastiktan sonra eklenecek.

## 2. Bugun Itibariyla Cekirdekte Olanlar

### 2.1 Ingest ve Dagitim

- RTMP, RTMPS, SRT, RTP, RTSP, WebRTC/WHIP, MPEG-TS ve HTTP Push ingest
- HLS, LL-HLS, DASH, HTTP-FLV, MP4, WebM ve ses cikislari
- FFmpeg tabanli live transcode ve ABR merdiveni
- recording, analytics, subscriber fanout ve stream yasam dongusu
- OBS Enhanced RTMP / multitrack ingest

### 2.2 Player, Embed ve Operasyon

- player, embed, iframe ve direct link akisi kararlilasti
- template preview gercek gomulu player akisi ile hizalandi
- QoE debug overlay, heartbeat telemetrisi ve kalici SQLite telemetry
- `Operasyon Merkezi` ve sekmeli `Canli Izleme ve Tanilama Merkezi`
- ham HLS / MPD manifest inceleme ve kullanim rehberi kartlari

### 2.3 Multitrack Video ve Audio

- OBS multitrack video katmanlari HLS varyantlarina baglanabiliyor
- DASH repack HLS varyantlarini representation olarak mapleyebiliyor
- audio alternate group ve player tarafinda ses izi secimi var
- kalite gecisi ve audio switch verisi telemetry / rapora yaziliyor
- `audio-only DASH` icin ayri `audio.mpd` ve init segment uretilebiliyor

### 2.4 Gozlemlenebilirlik ve Tanilama

- Prometheus `/metrics`
- OTel-benzeri `/api/observability/otel`
- QoE riskli yayinlar, uyarilar ve retention temelli housekeeping
- track bazli bitrate / runtime analytics
- `Hazir / Bekliyor / Kapali / Opsiyonel / Sorunlu` mantigina sahip teshis ekrani

### 2.5 Kayıt, Arsiv ve Yedek

- `Depolama ve Arsiv Merkezi`
- kayit, arsiv ve sistem yedeklerini tek merkezden yonetme
- kayit icin varsayilan `mp4`
- guvenli `TS capture + finalize/remux` modeli
- `MP4 Hazirla` arka plan isi
- lokal, S3/MinIO, SFTP ve rclone tabanli bulut hedefleri
- ayri kayit hedefi ve ayri sistem yedegi hedefi
- zamanlama, hedef seviyesi, soguk katman hazirligi
- ayni VPS uzerinde MinIO + SFTP saha testi

### 2.6 Urunlestirme

- runtime lisans modeli
- Linux servis yonetimi
- backup / restore / deploy omurgasi
- saglik endpoint ve tekrar edilebilir deploy akisi

## 3. Bu Fazda Kapanan Teknik Paket

Bu son fazda depolama ve bulut tarafi genisletildi.

Kapananlar:

- basit ve gelismis modlu yeni depolama akisi
- kayitlar ve yedekler icin ayri hedef tanimlama
- `Yerel Disk`, `AWS S3`, `MinIO`, `Cloudflare R2`, `Backblaze B2`, `Wasabi`,
  `DigitalOcean Spaces`, `Linode Object Storage`, `Scaleway Object Storage`,
  `IDrive e2`, `SFTP` kartlari
- rclone profili uzerinden `Google Drive`, `OneDrive`, `Dropbox`,
  `Google Cloud Storage`, `Azure Blob`, `Box`, `pCloud`, `MEGA`,
  `Nextcloud`, `WebDAV` gibi hedefleri kullanma altyapisi
- hedef bazli `Baglantiyi Test Et`
- ayri kayit senkronu ve ayri yedek senkronu
- senkron / donusum isleri icin ust ozet kartlari
- storage ekraninda tam sayfa yeniden cizim yerine parcali yenileme
- storage ekraninda renderer crash zincirinin kapatilmasi

## 4. Canli Saha Ogrenimleri

### 4.1 Multitrack Mikro Segment Sorunu

OBS multitrack yayininda gorulen mikro segment sorununun kok nedeni
RTMP chunk reader tarafindaki timestamp delta birikimiydi.

Kalici duzeltmeler:

- CSID bazli mutlak timestamp birikimi
- HLS segment duration fallback korumasi
- DASH `SegmentTimeline` fallback korumasi
- HLS master playlistin saglikli varyantlari yeniden ilan etmesi

Sonuc:

- mikro `EXTINF` segmentleri ortadan kalkti
- DASH `SegmentTimeline` tutarli hale geldi
- `360p + 1080p` ABR katmanlari yeniden saglikli sekilde ilan edildi

### 4.2 Recording ve Storage Sorunlari

Saha testinde storage ekranindaki tam sayfa donma / renderer crash
zinciri ve kayittan MP4 hazirlama sorunu goruldu.

Kalici duzeltmeler:

- storage aksiyonlari icin tam rerender kaldirildi
- preview teardown ve parcali yenileme akisi sertlestirildi
- `MP4 Hazirla` arka plan job haline getirildi
- recording TS paketleme mantigi HLS ile hizalandi
- yeni kayitlar icin temiz remux kaynagi uretilmeye baslandi

## 5. Bugunku Uretim Degerlendirmesi

Bugun icin en dogru tanim:

- iyi bir tek-node medya sunucusu
- urunlesmis bir yayin cekirdegi
- operasyon merkezi, telemetry ve depolama omurgasi olan HLS merkezli dagitim urunu

Bu haliyle su alanlar icin ciddi bicimde kullanilabilir:

- kurum ici TV
- webcast
- webinar
- radyo ve audio streaming
- markali player / embed dagitimi

Hala enterprise seviyeye cikarmak icin gerekli ana farklar:

- multi-node origin-edge
- daha derin guvenlik ve playback policy
- RBAC / SSO / audit
- gercek dis ortam storage ve failover testleri
- yuk testi ve soak testi

## 6. Acik Kalan Kisa Vade Fazlar

### 6.1 Depolama ve Arsiv Merkezi

- teknik terimleri daha da azalt
- kullaniciya secim yardimi ve hazir preset sihirbazi ekle
- rclone tabanli Drive / OneDrive / Dropbox akisini gercek saha testleriyle dogrula
- harici AWS S3 bucket ile gercek saha testi al
- ayni VPS uzerindeki MinIO + SFTP laboratuvar hedeflerini uzun sureli testlerle sertlestir

### 6.2 Audio-only DASH ve Recording Sertlestirme

- farkli tarayicilar, VLC ve dash.js tabanli istemcilerde `audio-only DASH` testleri
- buyuk dosya, uzun sureli kayit, servis restart ve gec finalize senaryolari
- eski bozuk `TS` kayitlar icin kurtarma / uyari akisi

### 6.3 Playback Guvenligi

- signed URL
- signed manifest / segment
- oturum bagli playback tokeni
- hotlink korumasi
- watermark
- IP / CIDR / geo policy

### 6.4 Tam DRM

- AES-128 HLS key servisi
- DRM abstraction
- Widevine / FairPlay / PlayReady hazirligi

## 7. Orta Vade Buyuk Fazlar

- multi-node origin-edge mimarisi
- RBAC, audit log, SSO
- SSAI ve monetizasyon
- uzun sureli soak test / yuk testi

## 8. Cekirdek Sonrasi Fazlar

- konferans odalari
- canli chat
- moderasyonlu soru-cevap
- sanal sinif rolleri
- yoklama ve breakout room
- takim ici mesajlasma

## 9. Sonraki En Dogru Sira

1. `Depolama ve Arsiv Merkezi` UX sadeleştirmesini ikinci turda tamamla
2. harici bir bucket ile gercek AWS S3 saha testi yap
3. rclone tabanli populer bulut hedeflerini gercek hesaplarla dogrula
4. `audio-only DASH` davranisini canli istemcilerde sertlestir
5. playback guvenligi fazina gir
6. sonra DRM ve origin-edge mimarisi tasarimina gec
