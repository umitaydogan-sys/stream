# FluxStream Gorev Listesi

Tarih: 26 Mart 2026

## 0. Admin Studio V2 Fazinda Kapananlar

- [x] `Dashboard` sayfasini studio KPI ve operasyon giris katmanina tası
- [x] `Streams` sayfasini filtre, rozet, hizli aksiyon ve kart/tablo hibrit duzene tasi
- [x] `Quick Settings` ekranini preset ve hizli kullanim odakli hale getir
- [x] `Genel Ayarlar` ekranini daha buyuk kategori yapisi ve studio gorunumu ile genislet
- [x] `Gelişmis Embed` ekranini urunlestir
- [x] `Player Sablonlari` ekranini studio seviyesine tasi
- [x] player sablonlari modalinda `Kaydet ve Acik Kal` akisini ekle
- [x] player sablonlari icin logo upload ve varlik kutuphanesi ekle
- [x] `Domain ve Embed` ekranini studio katmanina tasi
- [x] `Giris Protokolleri` ekranini studio katmanina tasi
- [x] `Cikis Formatlari` ekranini studio katmanina tasi
- [x] `Security` ekranini risk odakli studio katmanina tasi
- [x] `Health & Alerts` ekranini ortak studio gorunumune tasi
- [x] `Transkod / FFmpeg` ekranini studio katmanina tasi
- [x] `Izleyiciler` ekranini studio katmanina tasi
- [x] `Transcode Isleri` ekranini studio katmanina tasi
- [x] `Diagnostics` ekranini `Teshis ve Tedavi Merkezi` seviyesine tasi
- [x] `Bakim ve Yedek` ekranini `Depolama ve Arsiv Merkezi` ile rol ayrimi net olacak sekilde urunlestir
- [x] `Tokens` ekranini birinci sinif urun bileseni haline getir
- [x] `Logo ve Marka` ekranini ekle ve medya varlik kutuphanesine bagla
- [x] tum textarea, input, select ve teknik metin bloklari icin ortak studio stil denetimi yap

## 0.1 Bu Fazdan Sonra Acik Kalan Kisa Saha Dogrulamalari

- [ ] `Embed Studyosu` ve `Gelişmis Embed` ekranlarini canli veriyle uzun sureli operator kullanim testine sok
- [ ] `Player Sablonlari Studyosu` icin upload edilen marka varliklarini farkli sablonlarla saha testinde dogrula
- [ ] `Analitik Merkezi` ve `Teshis ve Tedavi Merkezi` ekranlarini canli veriyle uzun sureli operator kullanim testine sok
- [ ] `Bakim ve Yedek` ile `Depolama ve Arsiv Merkezi` arasindaki rol ayrimini son kullanici bakisiyla tekrar dogrula

## 0.1.1 Son UI Polish Turunda Kapananlar

- [x] tum panelde input, select ve textarea gorunumunu daha kompakt ve daha kosegen hale getir
- [x] `GelisÌ§mis Embed` ust kart metinlerini son kullaniciya daha teknik ve daha acik olacak sekilde sadeleştir
- [x] `GelisÌ§mis Embed` ekraninda tum direkt linkleri ve sekmeli onizlemeleri yeniden one cikar
- [x] `Player Sablonlari Studyosu`nu modaldan cikarip kalici kutuphane + taslak duzenleyici modeline tasi
- [x] `ABR Profilleri` katman olusturucuyu secim odakli preset / paket akisi ile sadeleştir

## 0.2 Bu Turda Kapanan Cekirdek Sertlestirme

- [x] `Analitik Merkezi` acilisindeki eksik JS yardimci fonksiyon hatasini kapat
- [x] `require_signed_url` aktif streamlerde sorgu parametreli `v2` signed URL zorlamasi getir
- [x] domain / referrer eslesmesini gercek host ve subdomain sinirlari ile guvenli hale getir
- [x] tokenli HLS / DASH teslimatta `private, no-store` cache davranisini uygula
- [x] `audio.mpd`, `audio_init.mp4` ve `audio_*.m4s` icin audio odakli MIME / baslik davranisini sertlestir
- [x] teshis ekranina `Audio-only DASH manifest` ve `DASH ses representation` gorunurlugu ekle
- [x] admin asset yukleme / listeleme / silme API'lerini ekle
- [x] `/media-assets/` uzerinden logo ve marka varliklarini servis edilir hale getir

## 1. Kayit, Arsiv ve Storage Fazinda Kapananlar

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

## 2. Embed + Analitik + ABR + Playback Guvenligi Fazinda Kapananlar

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

- [ ] `audio-only DASH` akisini gercek audio-only kaynakla tarayici, dash.js ve VLC tarafinda saha testinden gecir
- [ ] playback guvenligi V1 akisini domain/IP/token zorlamasi ile canli stream policy senaryolarinda dogrula
- [ ] harici bir bucket ile gercek AWS S3 saha testi yap
- [ ] rclone tabanli hedeflerde gercek `Google Drive`, `OneDrive` ve `Dropbox` saha testi al
- [ ] ayni VPS uzerindeki MinIO ve SFTP laboratuvar hedeflerini uzun sureli senaryolarla tekrar dogrula
- [ ] buyuk dosya, uzun sureli kayit ve servis restart senaryolarinda finalize/remux akisinin dayanikliligini arttir
- [ ] eski bozuk `TS` kayitlar icin kullaniciyi uyaran ve kurtarma yolunu gosteren akis ekle
- [ ] `Depolama ve Arsiv Merkezi` ekranini daha da sade, daha az teknik ve daha son kullanici odakli hale getir
- [ ] `Logo ve Marka Merkezi` ile `Player Sablonlari Studyosu` arasindaki varlik akisini son kullanici testleriyle ince ayarla

## 6. Sonraki Buyuk Fazlar

### 6.1 Playback Guvenligi V2

- [ ] signed playback politikalarina daha zengin presetler ekle
- [ ] oturum bagli watermark ve izleme izi davranisini sertlestir
- [ ] daha guclu iframe domain pinning ve embed policy setleri ekle
- [ ] playback auth olaylarini audit log tarafina bagla

### 6.2 Harici Storage Sertlestirme

- [ ] harici AWS S3 ile tam saha testi yap
- [ ] gercek Drive / OneDrive / Dropbox baglantilari ile yukleme ve geri alma dogrulasi al
- [ ] storage hata mesajlarini son kullanici icin daha anlasilir hale getir
- [ ] soguk katman ve zamanlama akislarini saha verisiyle ince ayarla

### 6.3 DRM Hazirligi

- [ ] AES-128 HLS sifreleme ve anahtar servis akisini ekle
- [ ] anahtar erisimini token / oturum / IP ile koru
- [ ] DRM abstraction layer tasarla
- [ ] Widevine / FairPlay / PlayReady entegrasyon noktalarini tasarla

## 7. Buyuk Urun Eksikleri

- [ ] multi-node origin-edge mimarisi
- [ ] RBAC, audit log ve SSO
- [ ] SSAI ve monetizasyon omurgasi
- [ ] uzun sureli soak test ve yuk testi kapsamini artir
- [ ] playback guvenligi ile lisans katmanini ortak policy modeline bagla

## 8. Cekirdek Tamamlandiktan Sonra

- [ ] konferans odalari
- [ ] canli chat
- [ ] moderasyonlu soru-cevap
- [ ] sanal sinif rolleri ve yoklama
- [ ] breakout room mantigi
- [ ] takim ici mesajlasma
