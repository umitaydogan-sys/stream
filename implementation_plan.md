# FluxStream Uygulama Plani

Tarih: 26 Mart 2026

## 0. Bugun Kapanan Faz

Bu turda `Admin Studio V2` fazi kapatildi.
Bu faz, daha once urunlesen `Embed Studyosu`, `Analitik Merkezi`,
`ABR Profilleri ve Teslimat Merkezi` ve `Playback Guvenligi V1`
omurgasinin ustune panelin geri kalan kritik ekranlarini de ayni urun
diline tasidi.

Kapanan ana paket:

- `Dashboard V2`
- `Streams V2`
- `Quick Settings V2`
- `Genel Ayarlar Merkezi`
- `Gelişmis Embed` laboratuvari
- `Player Sablonlari Studyosu`
- `Domain ve Embed Merkezi`
- `Giris Protokolleri Studyosu`
- `Cikis Formatlari Studyosu`
- `Security Studio`
- `Health & Alerts` studio katmani
- `Transkod / FFmpeg Merkezi`
- `Izleyici Merkezi`
- `Transcode Isleri Merkezi`
- `Teshis ve Tedavi Merkezi`
- `Bakim ve Yedek Merkezi`
- `Token Merkezi`
- `Logo ve Marka Merkezi`

Bu faz ile birlikte:

- admin panelin buyuk kismi ortak `studio` gorunumu etrafinda toplandi
- textarea, input, select ve teknik metin bloklari ortak stile kavustu
- `Gelişmis Embed` ekrani urun hissi verecek sekilde yenilendi
- `Player Sablonlari` modal akisi kapanmadan calisacak hale geldi
- logo yukleme ve medya varlik kutuphanesi eklendi
- teshis ekrani tani koyan degil, yonlendiren ve tedavi aksiyonu veren merkez haline geldi
- `Bakim ve Yedek` ile `Depolama ve Arsiv Merkezi` arasindaki rol farki netlestirildi

## 0.1 Canli Durum

Yerelde:

- `node --check internal/web/static/admin-studio.js`
- `go build ./cmd/fluxstream/`
- `go build ./cmd/fluxstream-license/`
- `go test ./...`

Canli host:

- servis: `active`
- health: `http://127.0.0.1:8844/api/health`
- Linux binary SHA256:
  `BB92CE8B47EA09884D4367DA96B785DBB6DA01275556A0928008C8B611C9D656`

## 0.2 Bu Fazdaki Cekirdek Sertlestirme

Bu fazla birlikte cekirdek tarafta da cakismaz sertlestirmeler yapildi:

- `Analitik Merkezi` acilisini kiran istemci yardimci eksigi kapatildi
- `require_signed_url` aktif streamlerde sadece sorgu parametreli `v2` signed URL kabul eder hale getirildi
- domain / referrer / host eslesmesi gercek host ve subdomain siniri mantigina cekildi
- tokenli HLS / DASH teslimatta daha guvenli `private, no-store` cache davranisi eklendi
- `audio.mpd`, `audio_init.mp4` ve `audio_*.m4s` icin audio odakli MIME davranisi sertlestirildi
- teshis ekranina `Audio-only DASH manifest` ve `DASH ses representation` gorunurlugu eklendi
- admin asset yukleme / listeleme / silme API'leri eklendi
- `/media-assets/` uzerinden logo ve marka varliklari servis edilir hale geldi

## 0.3 Son UI Polish Turu

Bu turun sonundaki ek iyilestirmeler:

- tum panelde input, select ve textarea gorunumu daha kompakt ve daha kosegen hale getirildi
- `Gelişmis Embed` ekraninda tum direkt linkler ve sekmeli onizlemeler tekrar one cikarildi
- `Gelişmis Embed` ust kart metinleri son kullaniciya daha teknik ve daha acik anlatacak sekilde sadeleştirildi
- `Player Sablonlari Studyosu` kalici kutuphane + taslak duzenleyici + gercek onizleme tezgahi modeline tasindi
- `ABR Profilleri` katman olusturucu secimli cozum / bitrate paketi mantigina yaklastirildi

## 1. Urun Vizyonu

FluxStream, tek binary ile kurulan, cok protokollu ingest alip HLS/DASH
merkezli dagitim yapan, oynatici/embed, operasyon, kayit, arsiv ve panel
katmani tek urunde toplanmis bir canli yayin omurgasidir.

Ana hedefler:

- yayini kararlı ve guvenli almak
- dusuk bantta akici izleme saglamak
- ABR ile kaliteyi baglanti kosullarina gore evirmek
- kayit, arsiv, yedek ve operasyon akisini urun seviyesine tasimak
- playback guvenligi ve Linux urunlestirmesini cekirdege gommek

Konferans, chat, sanal sinif ve mesajlasma katmanlari cekirdek streaming
omurgasi yeterince olgunlastiktan sonra eklenecek.

## 2. Bugun Itibariyla Cekirdekte Olanlar

### 2.1 Ingest ve Dagitim

- RTMP, RTMPS, SRT, RTP, RTSP, WebRTC/WHIP, MPEG-TS ve HTTP Push ingest
- HLS, LL-HLS, DASH, HTTP-FLV, MP4, WebM ve ses cikislari
- FFmpeg tabanli live transcode ve ABR merdiveni
- stream lifecycle, analytics, subscriber fanout ve recording akisi
- OBS Enhanced RTMP / multitrack ingest

### 2.2 Player, Embed ve Operasyon

- player, embed, iframe ve direct link akisi kararlilasti
- template preview gercek gomulu player akisi ile hizalandi
- `Embed Studyosu` ve `Gelişmis Embed` urunlesti
- `Player Sablonlari Studyosu` canli preview + acik modal modeli ile yenilendi
- QoE debug overlay, heartbeat telemetrisi ve kalici SQLite telemetry
- `Operasyon Merkezi` ve sekmeli canli tanilama merkezi
- ham HLS / MPD inceleme, kullanim rehberi ve debug akislari

### 2.3 Multitrack Video ve Audio

- OBS multitrack video katmanlari HLS varyantlarina baglanabiliyor
- DASH repack HLS varyantlarini representation olarak mapleyebiliyor
- alternate audio group ve player tarafinda ses izi secimi var
- kalite gecisi ve audio switch verisi telemetry / rapora yaziliyor
- `audio-only DASH` icin ayri `audio.mpd` ve init segment uretilebiliyor

### 2.4 Gozlemlenebilirlik ve Tanilama

- Prometheus `/metrics`
- OTel benzeri `/api/observability/otel`
- QoE riskli yayinlar, esik tabanli uyarilar ve housekeeping
- track bazli bitrate / runtime analytics
- `Teshis ve Tedavi Merkezi`
- `Hazir / Bekliyor / Kapali / Opsiyonel / Sorunlu` mantigi

### 2.5 Kayit, Arsiv ve Yedek

- `Depolama ve Arsiv Merkezi`
- kayit, arsiv ve sistem yedeklerini tek merkezden yonetme
- kayit icin varsayilan `mp4`
- guvenli `TS capture + finalize/remux` modeli
- `MP4 Hazirla` arka plan isi
- lokal, S3/MinIO, SFTP ve rclone tabanli bulut hedefleri
- ayri kayit hedefi ve ayri sistem yedegi hedefi
- zamanlama, hedef seviyesi ve soguk katman hazirligi
- ayni VPS uzerinde MinIO + SFTP saha testi

### 2.6 Urunlestirme ve Admin Katmani

- `Dashboard V2`
- `Streams V2`
- `Quick Settings V2`
- `Genel Ayarlar Merkezi`
- `Domain ve Embed Merkezi`
- `Giris Protokolleri Studyosu`
- `Cikis Formatlari Studyosu`
- `Security Studio`
- `Health & Alerts` studio katmani
- `Transkod / FFmpeg Merkezi`
- `Izleyici Merkezi`
- `Transcode Isleri Merkezi`
- `Token Merkezi`
- `Logo ve Marka Merkezi`

## 3. Bu Fazda Kapanan Teknik Paket

Bu son fazda admin panel geneline urun dili tasindi.

Kapananlar:

- ortak `studio` CSS ve kontrol stili
- global textarea / input / select / monospace audit'i
- legacy sayfalarin studio wrapper ile birlestirilmesi
- `Logo ve Marka Merkezi` ile upload tabanli varlik kutuphanesi
- `Player Sablonlari` icin logo upload ve acik modal kaydetme akisi
- `Gelişmis Embed` icin daha guclu laboratuvar ve hizli test aksiyonlari
- `Teshis ve Tedavi Merkezi` ile tani + aksiyon birlesimi
- `Bakim ve Yedek` ile `Depolama ve Arsiv Merkezi` rol ayriminin netlestirilmesi

## 4. Canli Saha Ogrenimleri

### 4.1 Multitrack Mikro Segment Sorunu

OBS multitrack yayininda gorulen mikro segment sorununun kok nedeni RTMP
chunk reader tarafindaki timestamp delta birikimiydi.

Kalici duzeltmeler:

- CSID bazli mutlak timestamp birikimi
- HLS segment duration fallback korumasi
- DASH `SegmentTimeline` fallback korumasi
- HLS master playlistin saglikli varyantlari yeniden ilan etmesi

Sonuc:

- mikro `EXTINF` segmentleri ortadan kalkti
- DASH `SegmentTimeline` tutarli hale geldi
- `360p + 1080p` ABR katmanlari saglikli sekilde ilan edildi

### 4.2 Recording ve Storage Sorunlari

Saha testinde storage ekranindaki tam sayfa donma / renderer crash zinciri
ve kayittan MP4 hazirlama sorunu goruldu.

Kalici duzeltmeler:

- storage aksiyonlarinda tam rerender kaldirildi
- preview teardown ve parcali yenileme akisi sertlestirildi
- `MP4 Hazirla` arka plan isi haline getirildi
- recording TS paketleme mantigi HLS ile hizalandi
- yeni kayitlar icin temiz remux kaynagi uretilmeye baslandi

## 5. Bugunku Uretim Degerlendirmesi

Bugun icin en dogru tanim:

- iyi bir tek-node medya sunucusu
- urunlesmis bir yayin cekirdegi
- operasyon, depolama ve player/embed paneli olan HLS merkezli dagitim urunu

Bu haliyle su alanlar icin ciddi bicimde kullanilabilir:

- kurum ici TV
- webcast
- webinar
- radyo ve audio streaming
- markali player / embed dagitimi

Enterprise seviyeye cikarmak icin hala acik kalan ana farklar:

- multi-node origin-edge
- daha derin playback policy ve DRM
- RBAC / SSO / audit
- gercek dis ortam storage ve failover testleri
- yuk testi ve soak testi

## 6. Siradaki Buyuk Fazlar

### 6.1 Admin Studio V2 Sonrasi Kisa Vade

Bu faz kapandi. Simdi en dogru siradaki isler:

- `audio-only DASH` akisini gercek audio-only kaynakla tarayici, dash.js ve VLC tarafinda dogrulamak
- `Embed Studyosu`, `Player Sablonlari Studyosu` ve `Analitik Merkezi` ekranlarini canli veriyle uzun sureli operator kullanim testine sokmak
- playback guvenligi V1 akisini domain/IP/token zorlamasi ile canli stream policy senaryolarinda sertlestirmek
- `Logo ve Marka Merkezi` ile player sablonlari arasindaki varlik akisini saha kullaniminda ince ayarlamak
- `Bakim ve Yedek` ile `Depolama ve Arsiv Merkezi` rol ayrimini gerekirse daha da sadeleştirmek

### 6.2 Harici Storage ve Saha Testleri

Planlanan alt basliklar:

- harici AWS S3 bucket ile gercek saha testi
- rclone tabanli `Google Drive`, `OneDrive` ve `Dropbox` akisini gercek hesaplarla dogrulama
- ayni VPS uzerindeki MinIO ve SFTP laboratuvar hedeflerini uzun sureli senaryolarla sertlestirme
- buyuk dosya, uzun sureli kayit, servis restart ve gec finalize senaryolari
- eski bozuk `TS` kayitlar icin kurtarma / uyari akisi

### 6.3 Playback Guvenligi V2 ve DRM Hazirligi

Planlanan alt basliklar:

- signed playback politikasini daha zengin presetlerle genisletmek
- oturum bagli watermark ve izleme izi tarafini saha verisiyle sertlestirmek
- AES-128 HLS key servis ve policy modeli icin zemin hazirlamak
- DRM abstraction katmani ve lisans baglantilari icin enterprise tasarim cikarmak

## 7. Sonraki Kisa Vade Sertlestirme Basliklari

- `Depolama ve Arsiv Merkezi` teknik terimlerini daha da azalt
- kullaniciya secim yardimi ve hazir preset sihirbazi ekle
- harici AWS S3 ve populer bulut hedeflerinin gercek saha testlerini yap
- `audio-only DASH` akisini farkli istemcilerde sertlestir
- buyuk dosya, uzun sureli kayit, servis restart ve gec finalize senaryolarini teste sok
- onceki bozuk `TS` kayitlar icin kurtarma / uyari akisi ekle

## 8. Tam DRM Hazirligi

- AES-128 HLS key servisi
- DRM abstraction
- Widevine / FairPlay / PlayReady hazirligi

## 9. Orta Vade Buyuk Fazlar

- multi-node origin-edge mimarisi
- RBAC, audit log ve SSO
- SSAI ve monetizasyon
- uzun sureli soak test / yuk testi
