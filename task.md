# FluxStream Gorev Listesi

## 1. Tamamlanan Ana Teknik Isler

- [x] player, embed, iframe ve direct link akisini yeniden kararlilastir
- [x] player template preview'i gercek gomulu player gorunumuyle hizala
- [x] sahte `offline`, `403` ve framing problemlerini kapat
- [x] runtime lisans modelini ABR / RTMPS / recording / branding tarafina bagla
- [x] backup, restore plani ve Linux servis omurgasini urunlestir
- [x] OBS `Cok kanalli Video` baglantisini kabul edecek ingest temelini kur
- [x] Enhanced RTMP multitrack paketlerini parse et
- [x] `trackId` bilgisini ingest tarafina al
- [x] panelde OBS icin `Config Override JSON`, RTMP URL ve stream key yardimi ver
- [x] ayni OBS kurulum rehberini stream detay ekranina da ekle
- [x] player heartbeat tabanli QoE telemetrisi ekle
- [x] admin stream detay ekranina canli QoE karti ekle
- [x] diagnostics ekranina HLS varyant ve DASH representation sayaclarini ekle
- [x] log ve API hata metinleri icin normalize etme katmani ekle

## 2. Kapanan Kritik Canli Hata

- [x] OBS multitrack yayininda `1080p` varyantta gorulen mikro segment kok nedenini bul
- [x] RTMP chunk timestamp delta birikim hatasini kapat
- [x] HLS `EXTINF` mikro segment uretimini kalici olarak duzelt
- [x] DASH `SegmentTimeline` mikro surelerini kalici olarak duzelt
- [x] `master.m3u8` icinde saglikli `360p + 1080p` ABR varyantlarini yeniden ilan et

## 3. Bugun Kesinlesen Durum

- [x] OBS normal RTMP baglantisi calisiyor
- [x] OBS multitrack baglantisi calisiyor
- [x] OBS multitrack video katmanlari HLS varyantlarina baglanabiliyor
- [x] DASH repack, HLS varyantlarindan representation uretebiliyor
- [x] HLS oynatim kalitesi onceki problemli donemden belirgin sekilde daha iyi
- [x] admin tarafinda QoE ve stall davranisi artik gorulebiliyor

## 4. Simdi Acik Kalan Kisa Vade Isler

- [x] aktif kaynak / fallback durumunu debug panelde daha durust ve net goster
- [x] admin/API tarafinda gelen track listesini ve track metadata bilgisini ac
- [x] varsayilan video izi secimini urun seviyesine tasi
- [x] varsayilan audio izi secimini urun seviyesine tasi
- [x] MPD/HLS manifestlerini panelden ham metin olarak gorulebilir hale getir
- [x] telemetri ve harici oynatici kullanimini arayuzde daha anlasilir yap
- [x] `Operasyon Merkezi` menusunu ekle
- [x] `Canli Izleme ve Tanilama Merkezi` sayfa iskeletini kur
- [x] tum streamleri secip izleyebilen sol liste, filtre ve secim akisini ekle
- [x] arama alanini kaldirip tum streamleri gosteran selectbox secimi ekle
- [x] secim altyapisini ileride on-demand playlistleri de destekleyecek sekilde hazirla
- [x] sekmeli merkez panelde `Genel Durum`, `Player ve Teslimat`, `QoE ve Telemetri`, `Track ve ABR`, `Manifest ve Ham Veri`, `OBS ve Ingest`, `Teshis` alanlarini ac
- [x] mevcut `Kullanim ve Tanilama Rehberi`, `QoE ve Stall Telemetrisi` ve `Canli Track` kartlarini bu yeni merkeze daha gelismis sekilde tasi
- [x] `MP4 preview fix` isini ayni fazda kapat
- [x] Teshis ekraninda opsiyonel cikislari gereksiz `Sorunlu` gostermeyen akilli durum etiketlerini ekle
- [x] DASH preview fallback davranisini uzun canli testte yeniden dogrula
- [x] multitrack audio secimini son kullanici player tarafina tasi
- [x] HLS alternate audio group yapisini ekle
- [x] son kullanici player icinde secili ses izini kalici tercih ve fallback uyumlu hale getir
- [x] dusuk bant genisligi icin ABR profil merdivenini olcumle optimize et

## 5. QoE ve Gozlemlenebilirlik

- [x] player heartbeat telemetrisi topla
- [x] stall / toparlanma / reconnect bilgisini runtime bellekte sakla
- [x] admin stream detay ekraninda canli QoE karti goster
- [x] diagnostics ekraninda multitrack HLS / DASH sayaclarini goster
- [x] `bufferSeekOverHole` ve `bufferStalledError` davranisini canli testte ayristir
- [x] canli testte mikro segment kaynagini RTMP chunk timestamp zincirine kadar indir
- [x] telemetrileri kalici depolamaya al
- [x] telemetrileri grafik ve zaman serisi olarak goster
- [x] stream detay ekraninda canli track runtime karti goster
- [x] track bazli bitrate / cozunurluk analytics'i ekle
- [x] Prometheus / OpenTelemetry cikisi uret
- [x] retention, alarm ve esik tabanli QoE uyarilari ekle
- [x] Teshis ekraninda `Hazir / Bekliyor / Kapali / Opsiyonel / Sorunlu` ayrimini ekle
- [x] DASH tarafinda coklu audio adaptation setini canli testle dogrula
- [x] alarm esiklerini gercek saha verisine gore ince ayarla
- [x] track bazli kalite gecisi ve ses izi degisimi raporunu ekle
- [x] kalite / ses gecisi verisini saglik ve rapor ekranlarinda daha derin kullan
- [x] `Depolama ve Arsiv Merkezi` menusunu olustur ve kayit / depolama ayarlarini tek ekranda birlestir
- [x] kayitlari `ham capture + izlenebilir MP4/MKV finalize/remux` modeliyle gercekten oynatilabilir hale getir
- [x] sistem yedeklerini de ayni depolama ve arsiv merkezi altina bagla
- [x] harici hedeflerde `S3`, `MinIO`, `SFTP` seceneklerini ayni arayuzden yonetilebilir hale getir
- [x] `Depolama ve Arsiv Merkezi` ekranindaki tam sayfa kilitlenme / renderer crash zincirini kapat
- [x] `MP4 Hazirla` akisinin arka plan isi olarak sayfa degisse bile devam etmesini sagla
- [x] yeni kayitlarin TS paketlemesini Annex-B / ADTS ve ilk gecerli keyframe guvenligiyle duzelt
- [x] sistem yedegi silme ve kayit arsiv aksiyonlarini tam sayfa yeniden cizmeden guvenli hale getir
- [x] tum kayit baslatma akislarinda varsayilan formati `mp4` yap
- [x] S3 / MinIO upload'larinda eksik `Content-Length` hatasini kapat

## 6. Urunlestirme ve Lisans

- [x] runtime lisans modeli ABR / RTMPS / recording / branding tarafina baglandi
- [x] Linux servis yonetimi panel ve CLI tarafinda hazirlandi
- [x] backup / restore / upgrade plani olusturuldu
- [ ] `max_nodes` enforcement tamamla
- [ ] maintenance expiry ve grace policy ekle
- [ ] lisans modelini ilerideki cluster mimarisi ile uyumlu hale getir
- [ ] rollback guvenli Linux upgrade akisini sertlestir
- [ ] `.deb` / paketli Linux dagitimini tamamla

## 7. Arayuz ve Dokumantasyon

- [x] dokuman dosyalarini Turkce tut
- [x] panelde OBS icin adim adim kurulum rehberi ver
- [x] multitrack ingest icin panelde kopyalanabilir JSON ver
- [x] `Operasyon Merkezi` icin rehber metinleri ve kullanim aciklamalarini ekle
- [ ] admin panelde kalan `de`, `es`, `fr` ceviri eksiklerini kapat
- [x] `production_status.md` dosyasini yeni milestone'a gore guncelle

## 8. Uretim Seviyesi Buyuk Eksikler

- [ ] multi-node origin-edge mimarisi
- [x] S3 / MinIO archive ve restore akisi
- [x] SFTP tabanli arsiv ve yedek yukleme akisi
- [x] kayit / yedek / harici hedef yonetimini tek merkezde birlestir
- [x] Prometheus / OpenTelemetry / alarm sistemi
- [ ] RBAC, audit log ve SSO
- [ ] DRM ve gelismis playback guvenligi
- [ ] SSAI ve monetizasyon omurgasi
- [ ] uzun sureli soak test ve yuk testi kapsamini artir
- [x] ayni VPS uzerinde MinIO + SFTP ile gercek upload / geri yukleme saha testi yap
- [ ] harici bir bucket ile gercek AWS S3 saha testi yap
- [ ] MinIO ve SFTP akisini daha uzun sureli arsiv / geri yukleme testleriyle sertlestir
- [ ] kayit oynatim akisinda buyuk dosya, uzun sure ve kesintili finalize senaryolarini sertlestir
- [ ] eski bozuk `TS` kayitlar icin kullaniciyi dogru yonlendiren kurtarma / uyari akislarini ekle
- [ ] `Depolama ve Arsiv Merkezi` ekranini teknik terimleri azaltarak daha sade ve son kullanici dostu hale getir
- [ ] `Google Drive` ve `OneDrive` gibi populer cloud arsiv hedeflerini ekle
- [x] `audio-only DASH` icin ayri `audio.mpd` ve audio-only init segment uret

## 9. Siradaki Oncelikli Faz

- [ ] `Depolama ve Arsiv Merkezi` ekranini daha sade, daha az teknik ve son kullanici dostu hale getir
- [ ] harici `S3` bucket ile saha testi al
- [ ] ayni VPS uzerindeki `MinIO` ve `SFTP` laboratuvar hedeflerini UI/UX akislariyla yeniden dogrula
- [ ] `audio-only DASH` davranisini VLC, dash.js audio player ve farkli tarayicilarda canli testle sertlestir
- [ ] buyuk dosya, uzun sureli kayit ve servis restart senaryolarinda finalize/remux akisini sertlestir

## 10. Playback Guvenligi ve DRM Fazlari

### 10.1 Dusuk Butce / Kisa Vade Playback Guvenligi

- [ ] kisa omurlu signed playback URL destegi ekle
- [ ] manifest ve segment istekleri icin imzali token dogrulamasi ekle
- [ ] oturum bagli playback token mantigi kur
- [ ] domain / referrer / origin tabanli hotlink korumasi ekle
- [ ] IP / CIDR allowlist ve geo-kisit policy altyapisini ekle
- [ ] playback rate limit ve esik tabanli bloklama ekle
- [ ] gorunur watermark / dynamic overlay / oturum izi ekle
- [ ] playback auth olaylarini audit log tarafina bagla

### 10.2 Orta Vade Gelismis Playback Guvenligi

- [ ] AES-128 HLS sifreleme ve anahtar servis akisini ekle
- [ ] anahtar erisimini token / oturum / IP ile koru
- [ ] lisansli playback policy seti olustur
- [ ] embed domain pinning ve signed iframe mantigi ekle

### 10.3 Uzun Vade Tam DRM

- [ ] DRM abstraction layer tasarla
- [ ] Widevine / FairPlay / PlayReady entegrasyon noktalarini tasarla
- [ ] CENC / CMAF ve DRM lisans sunucusu baglantisi icin enterprise faz plani cikar
- [ ] tam DRM ozelliklerini lisans modeliyle eslestir

## 11. Cekirdek Tamamlandiktan Sonra

- [ ] konferans odalari
- [ ] canli chat
- [ ] moderasyonlu soru-cevap
- [ ] sanal sinif rolleri ve yoklama
- [ ] breakout room mantigi
- [ ] takim ici mesajlasma
