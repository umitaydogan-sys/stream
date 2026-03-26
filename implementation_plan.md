# FluxStream Uygulama Plani

Tarih: 26 Mart 2026

## 0. Bugun Kapanan Yeni Faz

Bu turda urun hissini buyuten ikinci buyuk admin panel fazi kapatildi.
Ana paket:

- `Embed Studyosu`
- `Analitik Merkezi`
- `ABR Profilleri ve Teslimat Merkezi`
- `audio-only DASH` gorunurlugu ve istemci hazirlik katmani
- `Playback Guvenligi V1`

Bu faz ile birlikte:

- embed kodlari tek textarea ekranindan cikti studyosuna donustu
- analitik ekrani KPI, grafik ve sorunlu yayin merkezi haline geldi
- ABR ayarlari ham JSON alani olmaktan cikti, preset ve katman studyosuna donustu
- signed playback URL ve token temelli guvenlik ayarlari embed akisina baglandi
- audio-only teslimat linkleri daha gorunur ve daha uygulanabilir hale geldi

## 0.1 Fazin Canli Durumu

Yerelde:

- `node --check internal/web/static/admin-studio.js`
- `go build ./cmd/fluxstream/`
- `go test ./...`

temiz gecti.

VPS:

- servis: `active`
- health: `http://127.0.0.1:8844/api/health`
- Linux binary SHA256:
  `44E5A33FD771B67934CF0187F2A0F9ABCE5A6DE544F77F3F7B00A97321DA9E70`

## 0.2 Son Cekirdek Sertlestirme Paketi

Bu faz kapandiktan sonra cekirdek tarafta cakismaz bir sertlestirme turu daha uygulandi.

Bu turda:

- `Analitik Merkezi` acilisini kiran eksik istemci yardimcisi kapatildi
- `require_signed_url` aktif streamlerde sadece sorgu parametreli `v2` signed URL kabul edilir hale getirildi
- domain / referrer eslesmesi alt dizi mantigindan cikarilip gercek host ve subdomain siniri ile guvenli hale getirildi
- tokenli HLS / DASH teslimatta cache davranisi `private, no-store` seviyesine cekildi
- `audio.mpd`, `audio_init.mp4` ve `audio_*.m4s` icin daha net audio MIME ve baslik davranisi eklendi
- teshis tarafina `Audio-only DASH manifest` ve `DASH ses representation` gorunurlugu eklendi
- `Bakim ve Yedek` ile `Depolama ve Arsiv Merkezi` arasindaki rol ayrimi sayfa icinden netlestirildi

Yeni statik dosyalar:

- `internal/web/static/admin-studio.js`
- `internal/web/static/admin-studio.css`

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

## 6. Siradaki Buyuk Faz

### 6.1 Embed + Analitik + ABR Stüdyosu + Playback Guvenligi Fazı

Bu faz, birbirinden kopuk duran ama ayni urun hissi eksigini tasiyan
uc ekrani ayni paket icinde urun seviyesine tasir:

- `Embed Kodlari` menusu
- `Analitik` sayfasi
- `Teslimat / ABR` sayfasi

Ayni faz icinde iki teknik sertlestirme hatti da kapatilir:

- `audio-only DASH` istemci sertlestirmesi
- `playback guvenligi v1`

### 6.2 Embed Stüdyosu

Hedef:

- mevcut embed ekranini `Embed Stüdyosu` seviyesine tasimak
- kod uretimini tek textarea yerine secilebilir, profilli ve canli onizlemeli bir merkez haline getirmek

Planlanan alt basliklar:

- `Basit Mod` ve `Gelismis Mod`
- hazir kullanim tipleri:
  `Web sitesi`, `Haber portalı`, `Kurumsal sayfa`, `Mobil uyumlu`,
  `Sadece ses`, `Gizli yayın`, `Token korumalı`, `Düşük gecikme`,
  `DASH`, `HLS`, `MP4 fallback`
- kartli cikislar:
  `Iframe`, `Script embed`, `Player URL`, `Audio player`,
  `Popup player`, `Direct manifest`, `VLC linki`
- canli onizleme, kullanim aciklamasi ve `nerede kullanilir` kutusu
- secilebilir opsiyonlar:
  `responsive`, `autoplay`, `muted`, `poster`, `branding`,
  `watermark`, `audio-only`, `start quality`, `token`,
  `signed URL`, `referrer policy`
- stream bazli kaydedilebilir `Embed Profili`
- `Kopyala`, `Paylaş`, `Test Et`, `Yeni sekmede aç`, `Debug ile aç`
- eksik veya gecersiz parametrelerde korumali uyari akisi

### 6.3 Analitik Merkezi

Hedef:

- mevcut basit analitik ekranini tek merkezli bir `Analitik Merkezi` haline getirmek

Planlanan alt basliklar:

- ust sabit filtre blogu:
  `tarih aralığı`, `stream seçimi`, `canlı görünüm`, `geçmiş rapor`
- KPI kartlari:
  `aktif izleyici`, `tepe izleyici`, `ortalama buffer`, `stall`,
  `kalite geçişi`, `audio switch`, `hata oranı`, `en çok izlenen stream`
- tum streamler ve tek stream gorunumu ayni sayfada
- gelismis grafikler:
  `izleyici zaman serisi`, `buffer trendi`, `stall trendi`,
  `kalite dağılımı`, `cihaz/oynatıcı kaynağı`,
  `audio track kullanımı`, `ABR katman dağılımı`
- `Sorunlu yayınlar` bolumu
- ayri kartlar:
  `Kalite geçiş raporu`, `Audio track değişim raporu`
- `CSV` ve `JSON` disa aktarma
- ilgili streamin `Operasyon Merkezi` sayfasina gecis
- `Analitik alarm merkezi` ve esik asimi kartlari

### 6.4 ABR Profilleri ve Teslimat Merkezi

Hedef:

- mevcut JSON odakli ekranı `ABR Profilleri ve Teslimat Merkezi` seviyesine tasimak

Planlanan alt basliklar:

- hazir profil secimi korunacak ama yanina form tabanli profil olusturucu gelecek
- kullanici JSON yazmak zorunda kalmayacak
- `katman ekle`, `katman sil`, `surukle sirala`
- alanlar:
  `çözünürlük`, `bitrate`, `max bitrate`, `buffer`, `fps`,
  `preset`, `audio bitrate`
- hazir kartli presetler:
  `Mobil`, `Dengeli`, `Dayanıklı`, `TV`, `Yüksek kalite`,
  `Audio-only`, `Radyo`, `Sadece düşük bant`
- `Profili kaydet`, `çoğalt`, `içe al`, `dışa aktar`
- `JSON görünümü` sadece gelismis modda
- her profil icin:
  `tahmini CPU yükü`, `tahmini upload`, `düşük bant uyumu`,
  `önerilen kullanım`
- secilen profil icin beklenen HLS / DASH cikisini gosteren canli test kutusu
- `varsayılan profil`, `stream bazlı özel profil`, `global profil kütüphanesi`
- `Yayin bazli öneri motoru`

### 6.5 Audio-only DASH Sertlestirme

Planlanan alt basliklar:

- tarayici, dash.js ve VLC tarafinda `audio-only DASH` dogrulamasi
- `audio.mpd`, `manifest.mpd`, `init segment`, codec ve MIME basliklarini gercek istemcilerle test
- `Sadece ses oynatici` icin daha net UI
- `audio-only embed` ve `audio-only direct link` gorunurlugu
- DASH ses cikisi icin `hazır / bekliyor / sorunlu` tanisini daha netlestirme
- radyo ve podcast presetleri

### 6.6 Playback Guvenligi V1

Planlanan alt basliklar:

- `signed playback URL`
- `signed manifest ve segment erişimi`
- `süreli token`
- `tek domain / referrer kısıtı`
- `iframe domain pinning`
- `IP kısıtı`
- `tek kullanımlık token` veya `oturum bağlı token`
- `görünür watermark`
- `oturuma özel izleme izi`
- `embed güvenlik profilleri`

### 6.7 Bu Fazda Ayni Anda Eklenebilecek Guzel Parcalar

- `Embed Şablon Kütüphanesi`
- `Paylaşım Paketleri`
- `A/B kalite testi`
- `Teslimat sağlık özeti`
- `Yayın bazlı öneri motoru`
- `Analitik alarm merkezi`
- `Preset import/export`
- `Stream’e profil bağla / profili miras al`
- `Gömme kodları için marka profili`
- `Kısa link ve paylaşım linki üretimi`

## 7. Sonraki Kisa Vade Sertlestirme Basliklari

- `Depolama ve Arsiv Merkezi` teknik terimlerini daha da azalt
- kullaniciya secim yardimi ve hazir preset sihirbazi ekle
- rclone tabanli Drive / OneDrive / Dropbox akisini gercek saha testleriyle dogrula
- harici AWS S3 bucket ile gercek saha testi al
- ayni VPS uzerindeki MinIO + SFTP laboratuvar hedeflerini uzun sureli testlerle sertlestir
- buyuk dosya, uzun sureli kayit, servis restart ve gec finalize senaryolari
- eski bozuk `TS` kayitlar icin kurtarma / uyari akisi

## 8. Tam DRM Hazirligi

- AES-128 HLS key servisi
- DRM abstraction
- Widevine / FairPlay / PlayReady hazirligi

## 9. Orta Vade Buyuk Fazlar

- multi-node origin-edge mimarisi
- RBAC, audit log, SSO
- SSAI ve monetizasyon
- uzun sureli soak test / yuk testi

## 10. Cekirdek Sonrasi Fazlar

- konferans odalari
- canli chat
- moderasyonlu soru-cevap
- sanal sinif rolleri
- yoklama ve breakout room
- takim ici mesajlasma

## 11. Sonraki En Dogru Sira

1. `Embed Stüdyosu` ekranini kur
2. `Analitik Merkezi` ekranini urun seviyesine tası
3. `ABR Profilleri ve Teslimat Merkezi` profil mantigini devreye al
4. ayni faz icinde `audio-only DASH` istemci sertlestirmesini kapat
5. ayni faz icinde `Playback Guvenligi V1` katmanini ekle
6. sonra depolama sertlestirmesi ve harici AWS S3 saha testine don
7. sonra DRM ve origin-edge mimarisi tasarimina gec
