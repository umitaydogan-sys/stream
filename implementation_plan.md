# FluxStream Uygulama Plani

Tarih: 26 Mart 2026

## 0. Bugun Kapanan Paket

Bu turda iki kritik alan birlikte kapatildi:

- `Adaptive Teslimat Sonradan Acma` urun akisi eklendi
- Windows ve Linux dagitimlari yeniden uretildi
- VPS'e temiz kurulum yapildi
- tum temel urun dokumanlari yeni duruma gore hizalandi

## 0.1 Yeni Adaptive Teslimat Akisi

Artik kaynak yayin tek kalite baslasa bile stream sonradan `adaptive teslimat`
moduna alinabiliyor.

Urun davranisi:

- `Streams` ekraninda hizli aksiyon: `Adaptiveye Al`
- `Stream Detayi` ekraninda ayri `Adaptive Teslimat` karti
- kullanici bir `profil seti` secebiliyor:
  - `balanced`
  - `mobile`
  - `resilient`
  - `radio`
- uygulama modu secebiliyor:
  - `Sonraki yayinda etkinlestir`
  - `Canli yayina hemen uygula`

Bu akista sunucu:

- stream policy icinde `enable_abr=true` yazar
- secilen `profile_set` degerini kaydeder
- istenirse canli HLS/DASH transcode zincirini yeni profil ile yeniden kurar

## 0.2 Bugunku Canli Durum

Yerelde dogrulananlar:

- `go test ./...`
- `go build ./cmd/fluxstream/`
- `go build ./cmd/fluxstream-license/`
- Windows portable package
- Windows service package
- Windows installer
- Linux systemd package

Canli host:

- servis: `active`
- health: `http://127.0.0.1:8844/api/health`
- canli Linux binary SHA256:
  `15D7CE047BC886ACB39C2B594C669C0B150F9C706E7179E36056216A369923F4`

Windows paket hashleri:

- portable / service `fluxstream.exe`:
  `7A99F1A7E5FC2A75247C83A2A7DEB459DD18A65B10C435A6AD9A3C5C2E339C55`
- installer `FluxStream-Setup.exe`:
  `F763F48D60E33FEEFC00C1D6DA7A0E99436613BD3145F888DD858FD623C9FB17`

## 1. Bugune Kadar Kapanan Buyuk Fazlar

### 1.1 Ingest ve Dagitim Cekirdegi

- RTMP, RTMPS, SRT, RTP, RTSP, WebRTC/WHIP, MPEG-TS ve HTTP Push ingest
- HLS, LL-HLS, DASH, HTTP-FLV, MP4, WebM ve audio cikislari
- FFmpeg tabanli live transcode
- OBS multitrack video ve audio akisi

### 1.2 Player, Embed ve ABR Uretkenligi

- `Embed Studyosu`
- `Gelismis Embed`
- `Player Sablonlari Studyosu`
- `ABR Profilleri ve Teslimat Merkezi`
- form tabanli ABR profil olusturma
- kalite paketleri ve preset tabanli katman secimi

### 1.3 Operasyon ve Analitik

- `Operasyon Merkezi`
- `Analitik Merkezi`
- QoE telemetry
- track analytics
- Prometheus ve OTel benzeri cikislar
- `Teshis ve Tedavi Merkezi`

### 1.4 Kayit, Arsiv ve Yedek

- `Depolama ve Arsiv Merkezi`
- varsayilan `mp4` recording
- `TS capture + finalize/remux`
- `MP4 Hazirla` arka plan isi
- MinIO ve SFTP saha testi

### 1.5 Admin Studio V2

- `Dashboard`
- `Streams`
- `Quick Settings`
- `Genel Ayarlar`
- `Security`
- `Health & Alerts`
- `Transkod / FFmpeg`
- `Izleyiciler`
- `Transcode Isleri`
- `Tokens`
- `Logo ve Marka Merkezi`

## 2. Bugunku Teknik Kazanim

Bu turun yeni cekirdek kazanimlari:

- tek kalite baslamis bir stream sonradan adaptive olarak isaretlenebiliyor
- stream bazinda profil secimi ve teslimat davranisi daha gorunur hale geldi
- canli stream icin kontrollu `ABR restart` akisi eklendi
- stream listesi ve detay ekraninda adaptive durum rozetleri eklendi
- yeni buildler hem Windows hem Linux icin yeniden uretildi
- VPS temiz kurulum ile paket dogrulamasi alindi

## 3. Bugunku Uretim Degerlendirmesi

FluxStream bugun icin:

- urunlesmis bir tek-node medya sunucusu
- admin paneli guclu bir yayin kontrol duzlemi
- player / embed / template / analytics / storage / security katmanlari olan
  bir yayin urunu

Bu haliyle su alanlarda rahat kullanilabilir:

- webcast
- kurum ici TV
- radyo ve audio streaming
- markali embed ve player dagitimi
- kayit ve arsiv tabanli yayin is akislari

## 4. Acik Kalan Kisa Vade Isler

- `adaptive teslimat` icin `live_now` akisinin saha etkisini canli testte gozlemle
- `audio-only DASH` akisini gercek audio-only kaynakla tarayici, dash.js ve VLC'de dogrula
- playback guvenligi V1'i domain / IP / token policy senaryolariyla sertlestir
- harici AWS S3 bucket ile gercek saha testi yap
- rclone tabanli `Google Drive`, `OneDrive` ve `Dropbox` hedeflerini gercek hesaplarla dogrula
- buyuk dosya, uzun sureli kayit ve servis restart senaryolarinda remux dayanikliligini arttir

## 5. Sonraki Buyuk Fazlar

### 5.1 Playback Guvenligi V2

- daha zengin signed playback presetleri
- oturum bagli watermark ve izleme izi
- daha guclu embed policy setleri
- playback auth olaylarini audit mantigina baglama

### 5.2 Storage ve Harici Bulut Sertlestirmesi

- gercek AWS S3 saha testi
- Drive / OneDrive / Dropbox saha testi
- buyuk dosya ve gec yukleme senaryolari
- kullanici dostu hata mesajlari

### 5.3 DRM Hazirligi

- AES-128 HLS key servisi
- tokenli key dagitimi
- DRM abstraction katmani
- Widevine / FairPlay / PlayReady hazirlik noktasi

### 5.4 Origin-Edge Lite

- dusuk butceye uygun iki node modeli
- tek VPS icinde origin-edge laboratuvari
- local + VPS topolojisi
- sonra harici ikinci node ile saha testi
