# FluxStream Production Durum Raporu

Tarih: 10 Mart 2026

## Ozet

Bu turda dort ana alan tamamlandi:

- acilan `cmd` penceresinde startup banner ve servis adresleri geri getirildi
- kayip ikon/font zinciri duzeltildi
- ust bara `Yeniden Baslat` ve `Durdur` butonlari eklendi
- ham `DASH` zinciri canli repack mantigiyla calisir hale getirildi ve `HTTP-FLV` ham output yeniden benchmark edildi
- portable release icin `ffmpeg` runtime paketleme ve otomatik kesif eklendi
- Windows service paketi hazirlandi
- Linux systemd paketi hazirlandi
- Inno Setup tabanli Windows GUI installer uretildi
- Windows service install/start/stop/uninstall self-test basariyla gecti

Genel sonuc:

- `HLS`: production-ready
- `LL-HLS`: production-ready
- `DASH` ham output: production-ready
- `HTTP-FLV` ham output: production-ready
- `MP4` ham output: production-ready
- `WebM` ham output: production-ready
- `MP3/AAC/WAV/FLAC/Icecast` direkt linkler: production-ready
- admin/player asset zinciri: artik tamamen local vendor, CDN bagimliligi yok
- Windows portable release: `fluxstream.exe` + `ffmpeg/` runtime ile production-ready
- Windows service release paketi: hazir
- Windows GUI installer: hazir
- Linux systemd release paketi: hazir
- Windows GUI installer icinde opsiyonel kurulum modu / port / SSL sayfalari hazir
- Windows GUI installer icinde public domain alani da hazir

## Bu Turda Yapilanlar

### 1. Startup Konsol Ciktisi

`fluxstream.exe` acildiginda bos siyah pencere yerine su bilgiler yeniden yazdiriliyor:

- `FluxStream v2.0.0`
- `HTTP`, `RTMP`, `RTMPS`, `SRT`, `RTP`, `RTSP`, `WHIP`, `TS-UDP` endpoint listesi
- `Sunucu hazir. Yayin bekleniyor...`
- `Kapatmak icin Ctrl+C`

Bu bilgi dogrudan [main.go](C:/xampp/htdocs/stream/cmd/fluxstream/main.go) icinde console banner olarak yazdiriliyor.

### 2. Ikon ve Logo Duzeltmeleri

- `bootstrap-icons.css` font yolu duzeltildi
- local vendor font dosyalari `/static/fonts/...` altindan dogru servis ediliyor
- player ust markasindaki bozuk karakterli logo metni kaldirildi
- player artik local bootstrap icon ile `FluxStream` markasini gosteriyor

Ilgili dosyalar:

- [bootstrap-icons.css](C:/xampp/htdocs/stream/internal/web/static/vendor/bootstrap-icons.css)
- [player_html.go](C:/xampp/htdocs/stream/internal/web/player_html.go)
- [static_assets.go](C:/xampp/htdocs/stream/internal/web/static_assets.go)

### 3. Ust Bar Sunucu Kontrolleri

Admin ust barina iki buton eklendi:

- `Yeniden Baslat`
- `Durdur`

Bu butonlar:

- [admin_html.go](C:/xampp/htdocs/stream/internal/web/admin_html.go) icindeki `restartServer()` ve `stopServer()` fonksiyonlarini cagirir
- [server.go](C:/xampp/htdocs/stream/internal/web/server.go) uzerinden admin session kontroluyle korunur
- [main.go](C:/xampp/htdocs/stream/cmd/fluxstream/main.go) icindeki process control kanalina baglidir

Not:

- butonlar admin oturumu gerektirir
- restart sonrasi yeni surec `FLUXSTREAM_NO_BROWSER=1` ile acilir

### 4. Ham DASH Zinciri

Kok neden iki parcaydi:

1. canli `DASH` job'u publish aninda HLS manifest'i beklerken cok erken timeout oluyordu
2. `ffmpeg` segmentleri job output klasorune degil calisma dizinine yaziyordu
3. canli audio timeline live `DASH` manifest'te DTS uyarisina yol aciyordu

Yapilan duzeltmeler:

- `StartLiveDASH()` bloklayici degil, gecikmeli canli job modeline cevrildi
- bekleme suresi artirildi
- `ffmpeg` prosesinin `Dir` alani `jobOutputDir` olacak sekilde ayarlandi
- `DASH` argumanlari testte calisan profile gecirildi:
  - video `copy + tag:v avc1`
  - audio `aac` yeniden encode
  - `aresample=async=1:first_pts=0,asetpts=N/SR/TB`
  - `use_timeline=0`
  - `ldash=1`
  - `init-$RepresentationID$.m4s`
  - `chunk-$RepresentationID$-$Number%05d$.m4s`

Ilgili dosya:

- [manager.go](C:/xampp/htdocs/stream/internal/transcode/manager.go)

### 5. Ham HTTP-FLV Benchmark

`HTTP-FLV` raw endpoint tekrar test edildi.

Sonuc:

- endpoint canli byte akisi veriyor
- indirilen ornek `ffprobe` tarafinda `flv` olarak tanindi
- video `h264`
- audio `aac`

### 6. MP4 / WebM Durumu

Bu turdaki benchmark'ta `MP4` ve `WebM` ham endpoint'leri yeniden dogrulandi.

Sonuc:

- `MP4` -> gecerli `mov,mp4,m4a,3gp,3g2,mj2`
- `WebM` -> gecerli `matroska,webm`

### 7. FFmpeg Runtime Paketleme

Transcode katmaninin sadece sistem `PATH` icindeki `ffmpeg`e bagimli kalmasi kaldirildi.

Yeni mantik:

- uygulama once config icindeki `ffmpeg_path` degerine bakar
- sonra `FLUXSTREAM_FFMPEG_PATH` environment degiskenini kontrol eder
- sonra `fluxstream.exe` yaninda ve alt klasorlerinde su adaylari arar:
  - `ffmpeg.exe`
  - `ffmpeg/ffmpeg.exe`
  - `bin/ffmpeg.exe`
  - `tools/ffmpeg.exe`
  - `tools/ffmpeg/ffmpeg.exe`
  - `data/tools/ffmpeg.exe`
- ancak bunlar yoksa sistem `PATH` fallback olur

Portable paketleme icin su script eklendi:

- [package_windows_portable.ps1](C:/xampp/htdocs/stream/deployment/package_windows_portable.ps1)

Bu script:

- `fluxstream.exe` derler
- sistemdeki `ffmpeg` runtime klasorunu bulur
- `ffmpeg.exe` ile gerekli `*.dll` dosyalarini `dist/fluxstream-windows-amd64-portable/ffmpeg/` altina kopyalar
- `README.txt` olusturur

Gercek dogrulama:

- `PATH` icinden `ffmpeg` kaldirilmis ortamda portable paket calistirildi
- `/api/transcode/status` icinde su path goruldu:
  - `C:\xampp\htdocs\stream\dist\fluxstream-windows-amd64-portable\ffmpeg\ffmpeg.exe`
- `ffmpeg_version` dogru dondu

### 8. Dagitim Paketleri

Uretilen artefaktlar:

- `dist/fluxstream-windows-amd64-portable`
- `dist/fluxstream-windows-amd64-service`
- `dist/FluxStream-Setup.exe`
- `dist/fluxstream-linux-amd64-systemd`

Windows service paketi:

- `install_service.ps1`
- `uninstall_service.ps1`
- `ffmpeg/` runtime klasoru

Linux systemd paketi:

- `fluxstream` linux binary
- `systemd/fluxstream.service`
- `install.sh`
- `uninstall.sh`

Not:

- Windows service komutlari gercek self-test ile dogrulandi
- self-test log dosyasi:
  - `dist/windows-service-selftest.log`
- Inno Setup 6 ile GUI installer derlendi
- installer icine opsiyonel su sayfalar eklendi:
  - servis mi manuel mod mu
  - custom HTTP / RTMP / HTTPS / RTMPS portlari
  - public domain / IP alani
  - CRT / KEY on yukleme ve SSL'yi aninda acma secenegi
- bu sayfalarda aciklama metni eklendi; ayarlarin kurulumdan sonra admin panelinden degistirilebildigi belirtiliyor
- installer secenekleri uygulama icinde `fluxstream.exe config set ...` ile yaziliyor
- izole paket testinde `config set` uygulanip ozel `http_port=9911` ile sunucu kaldirildi ve `/api/health` -> `200` dogrulandi
- admin paneline `Alan Adi / Embed` ayar sayfasi eklendi
- kopyala butonlari guvensiz origin fallback ile duzeltildi
- `localhost` kullanan embed/play URL'leri artik domain ayarina veya aktif request host'una gore uretiliyor
- `mp3/aac/ogg/wav/flac` raw audio cikislari gercek ffmpeg transcode ile duzeltildi
- Windows installer desktop icon varsayilani acik hale getirildi

## Gercek Zamanli Dogrulama

Canli benchmark 10 Mart 2026 tarihinde `live_f08222b9e8cbeceacf1a3296` key'i uzerinde yapildi.

### Sunucu ve UI

- `/api/health` -> `ok`
- `/static/vendor/bootstrap-icons.css` -> `200`
- `/static/fonts/bootstrap-icons.woff2` -> `200`
- `/play/{key}` HTML'i icinde:
  - `/static/vendor/bootstrap-icons.css`
  - `bi-lightning-charge-fill`
  - `Yayin cevrimdisi`
- admin ana HTML icinde:
  - `restartServer()`
  - `stopServer()`
  - `Yeniden Baslat`
  - `Durdur`

### DASH Raw Output

Canli yayin sirasinda su dosyalar olustu:

- `data/transcode/dash/live_f08222b9e8cbeceacf1a3296/init-0.m4s`
- `data/transcode/dash/live_f08222b9e8cbeceacf1a3296/init-1.m4s`
- `data/transcode/dash/live_f08222b9e8cbeceacf1a3296/chunk-0-00001.m4s`
- `data/transcode/dash/live_f08222b9e8cbeceacf1a3296/chunk-1-00001.m4s`
- `data/transcode/dash/live_f08222b9e8cbeceacf1a3296/manifest.mpd`

Canli probe sonucu:

- `/dash/live_f08222b9e8cbeceacf1a3296/manifest.mpd` -> `200`
- `ffmpeg -t 8 -i http://localhost:8844/dash/live_f08222b9e8cbeceacf1a3296/manifest.mpd -f null -` -> exit code `0`

Ek dogrulama:

- lokal decode testinde `WARN_COUNT=0`
- `non monotonically increasing dts` uyarisi yeniden uretilemedi
- canli manifest low-latency dynamic DASH formuna gecti

Karar:

- `DASH` ham output production-ready

### HTTP-FLV Raw Output

10 saniyelik canli capture sonucu:

- dosya boyutu: `909983` byte
- `ffprobe`:
  - container: `flv`
  - video: `h264`
  - audio: `aac`

Karar:

- `HTTP-FLV` ham output production-ready

### MP4 Raw Output

10 saniyelik canli capture sonucu:

- dosya boyutu: `1298054` byte
- `ffprobe format_name`: `mov,mp4,m4a,3gp,3g2,mj2`

Karar:

- `MP4` ham output production-ready

### WebM Raw Output

10 saniyelik canli capture sonucu:

- dosya boyutu: `3826824` byte
- `ffprobe format_name`: `matroska,webm`

Karar:

- `WebM` ham output production-ready

## Production Karari

### Production'a Alinabilir

- HLS
- LL-HLS
- DASH raw output
- HTTP-FLV raw output
- MP4 raw output
- WebM raw output
- MP3
- AAC
- WAV
- FLAC
- Icecast
- local vendor asset modeli
- admin ust bar stop/restart kontrolleri

### Izlenmesi Gereken Riskler

1. `DASH` audio timeline:
   - mevcut duzeltmeden sonra lokal decode testinde DTS uyarisi yok
   - yine de uzun sureli production izleme tavsiye edilir

2. `RTMPS`:
   - sertifika yoksa acilis log'unda hata gorunur
   - production'da `cert.pem` / `key.pem` saglanmadan kullanilmamali

3. Uzun sureli soak test:
   - `MP4`, `WebM`, `DASH`, `HTTP-FLV` icin 30-60 dakikalik yayin testi halen tavsiye edilir

## Sonraki Teknik Adimlar

1. `HTTP-FLV`, `MP4`, `WebM`, `DASH` icin otomatik self-test endpoint'i eklemek
2. tek sunuculu opsiyonel cache katmanini devreye almak
3. server metrics eklemek:
   - aktif publisher
   - aktif raw consumer
   - ffmpeg child process sayisi
   - segment olusum gecikmesi
4. uzun sureli soak test scriptlerini repoya almak
5. RTMPS sertifika kurulumu ve health check eklemek

## Standalone Dagitim

FluxStream tek basina calismaya devam eder. Harici cache veya CDN zorunlu degildir.

Opsiyonel tek-sunuculu reverse proxy/cache konfigi eklendi:

- [fluxstream-standalone.conf](C:/xampp/htdocs/stream/deployment/nginx/fluxstream-standalone.conf)
- [standalone_delivery.md](C:/xampp/htdocs/stream/deployment/standalone_delivery.md)

Bu modelde:

- FluxStream dogrudan kullanilabilir
- isteyen ayni makinede Nginx cache katmani acabilir
- daha sonra ayni tek origin'in onune CDN tanimlanabilir

## Kisa Sonuc

Bu tur sonunda sistemin eksik kalan operasyonel kisimlari toparlandi:

- startup CMD bilgileri geri geldi
- ikon/font zinciri duzeldi
- admin ust bar process kontrolu eklendi
- ham `DASH` zinciri artik gercek canli repack ile calisiyor ve audio DTS uyarisi kapatildi
- `HTTP-FLV`, `MP4` ve `WebM` ham output'lar yeniden benchmark edildi

Su anki durumla sistem, tek origin uzerinde production denemesi yapabilecek seviyede. Bir sonraki teknik esik, uzun sureli soak test ve CDN/edge dagitim katmani.
