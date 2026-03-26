# FluxStream Gorev Listesi

Tarih: 26 Mart 2026

## 0. Yeni Fazda Kapanan Buyuk Paket

- [x] `Embed Kodlari` ekranini `Embed Studyosu` seviyesine tasi
- [x] `Basit Mod` ve `Gelismis Mod` ayrimini ekle
- [x] hazir kullanim tipleri ekle:
  `Web sitesi`, `Haber portali`, `Kurumsal sayfa`, `Mobil uyumlu`,
  `Sadece ses`, `Gizli yayin`, `Token korumali`, `Dusuk gecikme`,
  `DASH`, `HLS`, `MP4 fallback`
- [x] kartli cikis tipleri ekle:
  `Iframe`, `Script embed`, `Player URL`, `Audio player`,
  `Popup player`, `Direct manifest`, `VLC linki`
- [x] canli onizleme, kullanim ozeti, kopyala/test/debug aksiyonlari ekle
- [x] stream bazli kaydedilebilir `Embed Profili` mantigini ekle
- [x] signed playback URL, sureli token, domain/IP kisiti ve watermark tabanli
  `Playback Guvenligi V1` omurgasini Embed Studyosu ile birlestir
- [x] `Analitik` ekranini `Analitik Merkezi` seviyesine tasi
- [x] KPI kartlari, zaman serileri, kalite/audio dagilimlari ve sorunlu yayinlar bolumunu ekle
- [x] JSON/CSV disa aktarma ve Operasyon Merkezi'ne hizli gecis ekle
- [x] `Teslimat / ABR` ekranini `ABR Profilleri ve Teslimat Merkezi` seviyesine tasi
- [x] form tabanli profil olusturucu, preset kutuphanesi, kaydet/cogalt/sil akisini ekle
- [x] tahmini CPU, upload, dusuk bant uyumu, teslimat saglik ozeti ve yayin bazli oneri motoru ekle
- [x] `audio-only DASH` icin player/embed/ABR UI gorunurlugunu guclendir
- [x] `DASH Ses` ve `HLS Ses` linklerini teslimat merkezi icinde gorunur hale getir

## 0.1 Bu Fazdan Sonra Acik Kalan Saha Dogrulamalari

- [ ] `audio-only DASH` akisini gercek audio-only kaynakla tarayici, dash.js ve VLC tarafinda saha testinden gecir
- [ ] yeni `Embed Studyosu` ve `Analitik Merkezi` ekranlarini canli veriyle uzun sureli operatör kullanim testine sok
- [ ] playback guvenligi V1 akisini domain/IP/token zorlamasi ile canli stream policy senaryolarinda dogrula

## 1. Bu Turda Kapanan Ana Basliklar

- [x] tum kayit baslatma akislarinda varsayilan formati `mp4` yap
- [x] kayit tarafinda `ham capture + finalize/remux` modelini oturt
- [x] `MP4 Hazirla` akisinin arka planda calisan is olarak devam etmesini sagla
- [x] yeni kayitlarin TS paketlemesini Annex-B / ADTS ve ilk gecerli keyframe guvenligi ile duzelt
- [x] `Depolama ve Arsiv Merkezi` ekranindaki tam sayfa kilitlenme / renderer crash zincirini kapat
- [x] kayit / yedek / arsiv aksiyonlarini tam sayfa yeniden cizmeden guvenli hale getir
- [x] sistem yedegi silme endpoint ve arayuz baglantisini calisir hale getir
- [x] `audio-only DASH` icin ayri `audio.mpd` ve audio-only init segment uret
- [x] ayni VPS uzerinde MinIO ile gercek S3-uyumlu upload / restore saha testi yap
- [x] ayni VPS uzerinde SFTP ile gercek upload / restore saha testi yap
- [x] MinIO / S3 upload yolundaki `Content-Length` eksigi hatasini kapat

## 2. Depolama ve Bulut Fazinda Kapananlar

- [x] `Depolama ve Arsiv Merkezi` menusunu olustur
- [x] kayit, arsiv ve sistem yedegini tek merkezde birlestir
- [x] basit ve gelismis mod ayrimini ekle
- [x] kayitlar ve sistem yedekleri icin ayri hedef tanimlayabil
- [x] isterse her iki akis icin ayni hedefi kullan
- [x] kayit ve yedek hedefleri icin ayri zamanlama mantigi ekle
- [x] hedef basina `standard / hot / cold` secenekleri ve soguk katmana gecis hazirligini ekle
- [x] `Yerel Disk`, `AWS S3`, `MinIO`, `Cloudflare R2`, `Backblaze B2`, `Wasabi`, `DigitalOcean Spaces`, `Linode Object Storage`, `Scaleway Object Storage`, `IDrive e2` kartlarini ekle
- [x] `SFTP` hedefini birinci sinif secenek olarak sun
- [x] `Google Drive`, `OneDrive`, `Dropbox`, `Google Cloud Storage`, `Azure Blob`, `Box`, `pCloud`, `MEGA`, `Nextcloud`, `WebDAV` gibi hedefler icin baglanti profili kartlarini ekle
- [x] S3 uyumlu saglayicilari tek backend motoru ile yonet
- [x] rclone tabanli genel bulut profili motoru ekle
- [x] kayit ve yedek hedefleri icin ayri `Baglantiyi Test Et` aksiyonu ekle
- [x] hedef kartlarinda kullanici dostu aciklama ve yonlendirme metinleri ekle
- [x] senkron ve donusum isleri icin ust ozet kartlari ekle
- [x] `Kayitlari Simdi Gonder` ve `Yedekleri Simdi Gonder` gibi daha anlasilir aksiyon metinleri ekle

## 3. Canli Yayin ve Oynatma Tarafinda Zaten Kapananlar

- [x] player, embed, iframe ve direct link akisini yeniden kararlilastir
- [x] player template preview'i gercek gomulu player gorunumuyle hizala
- [x] sahte `offline`, `403` ve framing problemlerini kapat
- [x] OBS `Cok kanalli Video` baglantisini kabul edecek ingest temelini kur
- [x] Enhanced RTMP multitrack paketlerini parse et
- [x] OBS multitrack video katmanlarini HLS varyantlarina bagla
- [x] DASH repack tarafini HLS varyantlariyla hizala
- [x] player heartbeat tabanli QoE telemetrisi ekle
- [x] admin stream detay ekranina canli QoE karti ekle
- [x] kalite gecisi ve ses izi degisimi raporunu ekle
- [x] Prometheus / OpenTelemetry cikisi uret
- [x] retention, alarm ve esik tabanli QoE uyarilari ekle
- [x] `Operasyon Merkezi` menusunu ve sekmeli canli tanilama merkezini ekle

## 4. Urunlestirme ve Linux Tarafi

- [x] runtime lisans modeli ABR / RTMPS / recording / branding tarafina baglandi
- [x] Linux servis yonetimi panel ve CLI tarafinda hazirlandi
- [x] backup / restore plani olusturuldu
- [x] temiz kurulum, kaldir-kur ve tekrar deploy akislari sahada denendi
- [ ] `max_nodes` enforcement tamamla
- [ ] maintenance expiry ve grace policy ekle
- [ ] rollback guvenli Linux upgrade akisini sertlestir
- [ ] `.deb` / paketli Linux dagitimini tamamla

## 5. Acik Kalan Kisa Vade Isler

- [ ] `Depolama ve Arsiv Merkezi` ekranini daha da sade, daha az teknik ve daha son kullanici odakli hale getir
- [ ] bulut baglanti profilleri icin adim adim sihirbaz ve hazir preset yardimlarini ekle
- [ ] rclone tabanli hedeflerde gercek `Google Drive`, `OneDrive` ve `Dropbox` saha testi al
- [ ] harici bir bucket ile gercek AWS S3 saha testi yap
- [ ] ayni VPS uzerindeki MinIO ve SFTP laboratuvar hedeflerini UI/UX akislariyla tekrar dogrula
- [ ] `audio-only DASH` davranisini farkli tarayicilar, VLC ve dash.js oyunculari ile canli testte sertlestir
- [ ] buyuk dosya, uzun sureli kayit ve servis restart senaryolarinda finalize/remux akisinin dayanikliligini arttir
- [ ] eski bozuk `TS` kayitlar icin kullaniciyi uyaran ve kurtarma yolunu gosteren akis ekle

## 6. Yeni Buyuk Faz: Embed + Analitik + ABR Stüdyosu + Playback Guvenligi

### 6.1 Embed Stüdyosu

- [ ] `Embed Kodlari` ekranini `Embed Stüdyosu` seviyesine tası
- [ ] `Basit Mod` ve `Gelismis Mod` ayrimini ekle
- [ ] hazir kullanim tipleri ekle:
  `Web sitesi`, `Haber portalı`, `Kurumsal sayfa`, `Mobil uyumlu`,
  `Sadece ses`, `Gizli yayın`, `Token korumalı`, `Düşük gecikme`,
  `DASH`, `HLS`, `MP4 fallback`
- [ ] kartli cikis tipleri ekle:
  `Iframe`, `Script embed`, `Player URL`, `Audio player`,
  `Popup player`, `Direct manifest`, `VLC linki`
- [ ] her embed tipi icin canli onizleme ve `nerede kullanilir` kutusu ekle
- [ ] secilebilir opsiyonlari checkbox/select yapisina tasi:
  `responsive`, `autoplay`, `muted`, `poster`, `branding`,
  `watermark`, `audio-only`, `start quality`, `token`,
  `signed URL`, `referrer policy`
- [ ] stream bazli kaydedilebilir `Embed Profili` mantigi ekle
- [ ] `Kopyala`, `Paylaş`, `Test Et`, `Yeni sekmede aç`, `Debug ile aç` aksiyonlarini belirginlestir
- [ ] eksik veya gecersiz parametrelerde korumali uyari akislarini ekle
- [ ] `Embed Şablon Kütüphanesi` ekle
- [ ] `Paylaşım Paketleri` ekle
- [ ] `Gömme kodları için marka profili` ekle
- [ ] `Kısa link ve paylaşım linki üretimi` ekle

### 6.2 Analitik Merkezi

- [ ] `Analitik` ekranini `Analitik Merkezi` seviyesine tası
- [ ] ustte sabit `tarih araligi` ve `stream secimi` filtresini ekle
- [ ] KPI kartlari ekle:
  `aktif izleyici`, `tepe izleyici`, `ortalama buffer`, `stall`,
  `kalite geçişi`, `audio switch`, `hata oranı`, `en çok izlenen stream`
- [ ] tum streamler ve tek stream gorunumunu ayni sayfada destekle
- [ ] gelismis grafikler ekle:
  `izleyici zaman serisi`, `buffer trendi`, `stall trendi`,
  `kalite dağılımı`, `cihaz/oynatıcı kaynağı`,
  `audio track kullanımı`, `ABR katman dağılımı`
- [ ] `Sorunlu yayınlar` bolumu ekle
- [ ] `Kalite geçiş raporu` ve `Audio track değişim raporu` kartlari ekle
- [ ] `Canlı görünüm` ve `geçmiş rapor` modlarini ayir
- [ ] `CSV` ve `JSON` disa aktarma ekle
- [ ] ilgili streamin `Operasyon Merkezi` sayfasina hizli gecis ekle
- [ ] `Analitik alarm merkezi` ekle

### 6.3 ABR Profilleri ve Teslimat Merkezi

- [ ] `Teslimat / ABR` ekranini `ABR Profilleri ve Teslimat Merkezi` seviyesine tası
- [ ] JSON odakli kullanim yerine form tabanli profil olusturucu ekle
- [ ] `katman ekle`, `katman sil`, `surukle sirala` akisini ekle
- [ ] katman alanlarini form tabanli ver:
  `çözünürlük`, `bitrate`, `max bitrate`, `buffer`, `fps`,
  `preset`, `audio bitrate`
- [ ] hazir preset kartlari ekle:
  `Mobil`, `Dengeli`, `Dayanıklı`, `TV`, `Yüksek kalite`,
  `Audio-only`, `Radyo`, `Sadece düşük bant`
- [ ] `Profili kaydet` ve `profil olarak çoğalt` akislarini ekle
- [ ] `JSON görünümü`nü sadece gelismis moda tasi
- [ ] profil bazli `tahmini CPU yükü`, `tahmini upload`, `düşük bant uyumu`, `önerilen kullanım` kartlarini ekle
- [ ] secilen profil icin beklenen HLS / DASH cikisini gosteren canli test kutusu ekle
- [ ] `varsayılan profil`, `stream bazlı özel profil`, `global profil kütüphanesi` mantigini ekle
- [ ] `Preset import/export` ekle
- [ ] `Stream’e profil bağla / profili miras al` mantigini ekle
- [ ] `Teslimat sağlık özeti` ekle
- [ ] `Yayın bazlı öneri motoru` ekle
- [ ] `A/B kalite testi` icin temel karsilastirma altyapisini ekle

### 6.4 Audio-only DASH Sertlestirme

- [ ] tarayici tarafinda `audio-only DASH` oynatimini dogrula
- [ ] dash.js istemcisi ile `audio-only DASH` oynatimini dogrula
- [ ] VLC ile `audio-only DASH` oynatimini dogrula
- [ ] `audio.mpd`, `manifest.mpd`, `init segment`, codec ve MIME basliklarini sertlestir
- [ ] `Sadece ses oynatici` UI'ini daha net hale getir
- [ ] `audio-only embed` ve `audio-only direct link` gorunurlugunu artir
- [ ] DASH ses cikisi icin `hazır / bekliyor / sorunlu` tanisini daha netlestir
- [ ] radyo ve podcast presetleri ekle

### 6.5 Playback Guvenligi V1

- [ ] signed playback URL destegi ekle
- [ ] signed manifest ve segment erisimi ekle
- [ ] sureli token mantigi ekle
- [ ] tek domain / referrer kisiti ekle
- [ ] iframe domain pinning ekle
- [ ] IP kisiti ekle
- [ ] tek kullanımlık token veya oturum bagli token mantigi ekle
- [ ] gorunur watermark ekle
- [ ] oturuma ozel izleme izi ekle
- [ ] embed guvenlik profilleri ekle

## 7. Playback Guvenligi ve DRM Fazlari

### 7.1 Kisa Vade Playback Guvenligi

- [ ] kisa omurlu signed playback URL destegi ekle
- [ ] manifest ve segment istekleri icin imzali token dogrulamasi ekle
- [ ] oturum bagli playback token mantigi kur
- [ ] domain / referrer / origin tabanli hotlink korumasi ekle
- [ ] IP / CIDR allowlist ve geo-kisit policy altyapisini ekle
- [ ] playback rate limit ve esik tabanli bloklama ekle
- [ ] gorunur watermark / dynamic overlay / oturum izi ekle
- [ ] playback auth olaylarini audit log tarafina bagla

### 7.2 Orta Vade Gelismis Playback Guvenligi

- [ ] AES-128 HLS sifreleme ve anahtar servis akisini ekle
- [ ] anahtar erisimini token / oturum / IP ile koru
- [ ] lisansli playback policy seti olustur
- [ ] embed domain pinning ve signed iframe mantigi ekle

### 7.3 Uzun Vade Tam DRM

- [ ] DRM abstraction layer tasarla
- [ ] Widevine / FairPlay / PlayReady entegrasyon noktalarini tasarla
- [ ] CENC / CMAF ve DRM lisans sunucusu baglantisi icin enterprise faz plani cikar
- [ ] tam DRM ozelliklerini lisans modeliyle eslestir

## 8. Buyuk Urun Eksikleri

- [ ] multi-node origin-edge mimarisi
- [ ] RBAC, audit log ve SSO
- [ ] SSAI ve monetizasyon omurgasi
- [ ] uzun sureli soak test ve yuk testi kapsamini artir
- [ ] playback guvenligi ile lisans katmanini ortak policy modeline bagla

## 9. Cekirdek Tamamlandiktan Sonra

- [ ] konferans odalari
- [ ] canli chat
- [ ] moderasyonlu soru-cevap
- [ ] sanal sinif rolleri ve yoklama
- [ ] breakout room mantigi
- [ ] takim ici mesajlasma
