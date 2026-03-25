# FluxStream Uygulama Plani

## 1. Urun Hedefi

FluxStream, tek binary calisan, yerelde ve Linux sunucuda kolay kurulabilen,
cok protokollu bir canli yayin sunucusudur.

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
- lisans yukleme ve dogrulama altyapisi mevcut
- runtime lisans modeli ABR, RTMPS, recording ve branding tarafina bagli
- backup, restore plani, Linux servis yonetimi ve dagitim akisi var

### 2.3 Oynatici ve Canli Izleme

- preview, direct link ve embed ekranlari yeniden kararlilastirildi
- `play` ve `embed` ekranlari iframe icinde kullanilabilir durumda
- player template preview gercek gomulu player akisina cekildi
- onceki `403`, framing ve sahte `offline` problemleri kapatildi

## 3. Bu Turun Ana Isi: OBS Cok Kanalli Video Uyumu

Bu turdaki hedef, OBS tarafinda `Cok kanalli Video` acildiginda
baglantinin kopmadan yayin kabul edebilmesiydi.

### 3.1 Kapatilan Teknik Sorun

Mevcut RTMP ingest sadece klasik FLV video/audio paketlerini okuyordu.
OBS cok kanalli video acildiginda, birincil izin yanina ek video ve ses izleri
Enhanced RTMP multitrack paketleri ile geliyordu.

Bu ek iz paketleri sunucuda hata uretip baglantiyi dusuruyordu.

### 3.2 Uygulanan Cozum

FLV okuyucuya Enhanced RTMP uyum katmani eklendi:

- Enhanced video tag algilama
- Enhanced audio tag algilama
- multitrack wrapper icinden `trackId` okuma
- H.264 `avc1` paketlerini klasik ic video formatina donusturme
- AAC `mp4a` paketlerini klasik ic ses formatina donusturme
- desteklenmeyen ek iz paketlerini baglantiyi dusurmeden yoksayma

RTMP handler ve stream manager tarafinda da:

- `track 0` / varsayilan iz akisa alinacak
- ek izler baglantiyi bozmadan yoksayilacak

### 3.3 Panel Entegrasyonu

Multitrack destegi sadece backend tarafinda birakilmadi.
Stream olusturma ve stream detay ekranlarina su yardim katmani eklendi:

- kopyalanabilir `Config Override JSON`
- secilen yayin moduna gore hazir OBS multitrack JSON on ayari
- yayin olustuktan sonra otomatik dolan `OBS RTMP URL`
- yayin olustuktan sonra otomatik dolan `OBS Yayin Anahtari`
- cok basit, adim adim OBS kurulum rehberi
- stream detay sayfasinda da ayni rehbere tekrar ulasabilme

Boylece son kullanici teknik dokuman aramadan panelden dogrudan
kopyala-yapistir ile kurulumu tamamlayabilir.

## 4. Bu Fazin Bilincli Siniri

Bu ilk fazda sunucu artik OBS cok kanalli video yayinini kabul eder,
ancak mevcut dagitim zinciri hala tek bir birincil video izi uzerinden calisir.

Yani su an:

- OBS multitrack baglantisi kopmuyor
- varsayilan video/ses izi akiyor
- ek izler gelecekte kullanilmak uzere taninmis oluyor
- ama ek izler henuz HLS master varyantlarina dogrudan baglanmiyor

Bu bilincli bir ilk adimdir. Once baglanti ve temel akisi kararlilastirdik.
Sonraki fazda bu izleri gercek ABR cikislarina baglayacagiz.

## 5. Sonraki Multitrack Fazlari

### 5.1 Faz 2

- ek video izlerini stream icinde sakla
- track metadatasini admin/API tarafina ac
- OBS tarafindan gelen kalite katmanlarini ayristir
- birincil iz yerine secilebilir varsayilan iz mantigi ekle

### 5.2 Faz 3

- OBS multitrack katmanlarini dogrudan ABR varyantlarina bagla
- HLS master playlist icine OBS kaynakli varyantlar yaz
- gerekmiyorsa sunucu tarafi yeniden encode maliyetini dusur
- transcode ile OBS katmanlari arasinda karma mod ekle

### 5.2.1 Bu Turda Kapanan Parca

Bu turda OBS tarafindan gelen cok kanalli video izleri,
live HLS override klasorunde dogrudan varyant playlist olarak yazilmaya baslandi.

Yeni davranis:

- `TrackID != 0` paketler artik ingest tarafinda tamamen atilmiyor
- live transcode katmani bu paketleri bootstrap olarak bellekte tutuyor
- cok kanalli video algilanirsa `ffmpeg` tabanli tek giris ABR yerine
  dogrudan multitrack HLS session aktif oluyor
- birincil OBS video izi kok `index.m3u8` olarak kalirken
  diger izler alt varyant dizinlerine yaziliyor
- kok `master.m3u8` bu varyantlari gercek ABR playlist olarak sunuyor

Bu sayede OBS'ten gelen kalite katmanlari ilk kez gercekten
HLS master playlist tarafina baglanmis oldu.

### 5.2.2 Bu Turda Kapanan Parca

Bu turda multitrack akis sadece HLS tarafinda birakilmadi.
DASH repack ve player izleme tarafi da buna gore guclendirildi.

Yeni davranis:

- canli DASH repack artik HLS master icindeki tum video izlerini mapliyor
- boylece OBS multitrack HLS kaynagi varsa DASH manifest tarafinda da
  coklu representation uretilebiliyor
- diagnostics ekranina `HLS varyant sayisi` ve
  `DASH representation sayisi` alanlari eklendi
- player sayfasi 5 saniyelik heartbeat ile QoE telemetrisi gonderiyor
- stall, toparlanma, buffer, aktif kaynak, reconnect ve hata bilgisi
  sunucu tarafinda runtime olarak toplanıyor
- stream detay ekranina admin tarafinda canli `QoE ve Stall Telemetrisi`
  karti eklendi
- admin preview iframe'leri `debug=1` ile acilarak overlay
  dogrudan panel icinde de gorulebilir hale geldi

Bu sayede bir sonraki OBS testinde sadece "goruntu var mi yok mu"
degil, HLS ve DASH tarafinda kac katman olustugu ile
player tarafinda stall davranisinin nasil aktigi da ayni anda izlenebilir.

### 5.2.3 Canli Testte Ortaya Cikan Yeni Durum

Canli VPS testinde OBS multitrack yayininda yeni bir dar bogaz netlesti.

Bugun icin durum soyle:

- kok HLS `index.m3u8` akisi kararlilastirildi
- player tarafinda siyah ekran ve sahte `offline` dongusu buyuk oranda azaltildi
- fakat `1080p` multitrack varyanti halen mikro segmentler uretiyor
- bu varyantta `0.010`, `0.011`, `0.016` saniye gibi bozuk `EXTINF` degerleri goruluyor
- ayni sorun DASH `SegmentTimeline` tarafina da yansiyor
- player kalite yukseltince `bufferSeekOverHole` ve `bufferStalledError`
  hatalariyla tekrar donabiliyor

Bu nedenle gecici guvenli mod uygulandi:

- multitrack direct HLS master playlist su an yalnizca stabil ana kati ilan ediyor
- player bu sayede bozuk ust kalite varyanta cikamiyor
- yayin akiciligi korunuyor, fakat tam adaptif cok katman davranisi
  gecici olarak sinirlanmis oluyor

Bir sonraki teknik hedef artik nettir:

- `1080p` multitrack track icin bozuk segment sinirlarini kokten bul
- Enhanced RTMP `timestamp / CTS / keyframe` hattini dogrula
- HLS muxer segment bolme mantigini track bazinda yeniden gozden gecir
- DASH uretimini bozuk HLS varyantindan etkilenmeyecek hale getir

### 5.3 Faz 4

- cok izli audio secimi
- dil bazli audio izleri
- track bazli analytics
- track bazli recording / archive metadata

## 6. Dogrulama Durumu

Yerelde:

- `go test ./...` gecti
- yeni Enhanced RTMP parser testleri eklendi
- `go build ./cmd/fluxstream/` gecti
- `go build ./cmd/fluxstream-license/` gecti
- admin panel JS sentaks kontrolu gecti

Windows:

- portable paket yeniden uretildi
- Windows tarafinda yeniden test icin hazir

## 7. Acik Urunlestirme Basliklari

- `max_nodes` lisans enforcement
- maintenance expiry ve grace policy
- `de`, `es`, `fr` ceviri kapsaminin tamamlanmasi
- `.deb` paketleme ve rollback guvenli dagitim
- Linux upgrade ve restore akisinin daha da sertlestirilmesi

## 8. Cekirdek Sonrasi Buyuk Faz

Streaming cekirdegi yeterince olgunlastiginda siradaki urun katmanlari:

- konferans odalari
- canli chat
- moderasyonlu soru-cevap
- sanal sinif rolleri
- yoklama ve katilim
- breakout room mantigi
- takim ici mesajlasma

## 9. Bugunku Uretim Degerlendirmesi

24 Mart 2026 itibariyla FluxStream artik yalnizca prototip degildir.
Tek node kurulumda, admin paneli olan, kurulabilen, yayin alabilen,
oynatabilen, embed ve player linkleri uretebilen bir medya sunucusu haline geldi.

Bugun icin guclu oldugu alanlar:

- tek binary ile kolay kurulum
- zengin ingest protokol seti
- genis playback ve audio output matrisi
- admin paneli, setup wizard ve stream yonetimi
- recording, analytics ve player template sistemi
- runtime lisans, backup ve Linux servis omurgasi
- OBS cok kanalli video baglantisini kabul eden ilk faz destek

Bugun icin sinirli veya eksik oldugu alanlar:

- cok kanalli OBS izleri henuz dogrudan ABR varyantlarina bagli degil
- multi-node origin-edge cluster mimarisi yok
- S3/MinIO benzeri harici obje depolama akisi yok
- tam enterprise seviye RBAC, audit log ve SSO katmani yok
- DRM, SSAI ve gelismis monetizasyon katmani yok
- otomatik test kapsami cekirdege gore henuz dar

Kisa karar:

- tek sunuculu kurum ici TV, radyo, webcast ve beyaz etiketli yayin isleri icin kullanilabilir seviyeye yaklasti
- ama Wowza / Ant Media / Red5 / Nimble gibi ust seviye urunlerle ayni ligde diyebilmek icin cluster, storage, telemetry ve security taraflari daha da olgunlasmali

## 10. Uretim Icin Siradaki Kapatma Sirasi

Bir sonraki zorunlu kapatma sirasini su sekilde goruyorum:

1. OBS cok kanalli izleri gercek ABR varyantlarina bagla
2. player QoE ve stall telemetry katmani ekle
3. Prometheus / OpenTelemetry / alarm entegrasyonu ekle
4. S3 veya MinIO recording archive ve restore akisini ekle
5. multi-node origin-edge mimarisini tasarla ve uygula
6. RBAC, audit log ve SSO tarafini urunlestir
