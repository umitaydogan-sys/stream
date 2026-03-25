# FluxStream Gorev Listesi

## 1. Bu Turda Tamamlananlar

- [x] OBS `Cok kanalli Video` acildiginda baglantiyi koparan ingest sorununu analiz et
- [x] RTMP/FLV tarafinda Enhanced RTMP multitrack paket yapisini netlestir
- [x] Enhanced video tag algilama ekle
- [x] Enhanced audio tag algilama ekle
- [x] multitrack wrapper icinden `trackId` oku
- [x] `avc1` H.264 paketlerini klasik ic video formatina donustur
- [x] `mp4a` AAC paketlerini klasik ic ses formatina donustur
- [x] ek iz paketlerini baglantiyi bozmadan yoksay
- [x] RTMP handler tarafinda sadece varsayilan izi akisa al
- [x] stream manager tarafinda ek izleri guvenli sekilde engelle
- [x] parser icin birim testleri ekle
- [x] `go test ./...` calistir
- [x] `go build ./cmd/fluxstream/` calistir
- [x] `go build ./cmd/fluxstream-license/` calistir
- [x] Windows portable paketi yeniden uret
- [x] stream olusturma ekranina kopyalanabilir OBS multitrack JSON alanini ekle
- [x] stream olusturma ekranina adim adim OBS kurulum rehberi ekle
- [x] yayin olustuktan sonra rehberi gercek RTMP URL ve stream key ile otomatik doldur
- [x] ayni OBS cok kanalli rehberini stream detay ekranina da ekle
- [x] Linux systemd paketini yeniden derle
- [x] VPS uzerinde FluxStream'i temiz kurulum olarak yeniden yukle
- [x] temiz kurulum sonrasi `api/health` ve `api/setup/status` ile dogrula
- [x] `production_status.md` dosyasini guncelle
- [x] `implementation_plan.md` dosyasini yeni duruma gore guncelle
- [x] `task.md` dosyasini yeni duruma gore guncelle
- [x] OBS multitrack video paketlerini ingest tarafinda tam gecir
- [x] live transcode katmaninda multitrack bootstrap bellegi ekle
- [x] OBS multitrack algilaninca dogrudan HLS varyant session ac
- [x] kok `master.m3u8` dosyasina OBS kaynakli gercek varyantlari yaz
- [x] cok kanalli ses paketlerini aktif varyantlara dagit
- [x] canli DASH repack icinde tum video izlerini maple
- [x] player tarafina heartbeat tabanli QoE telemetrisi ekle
- [x] admin stream detay ekranina canli QoE karti ekle
- [x] diagnostics ekranina HLS varyant ve DASH representation sayaci ekle
- [x] admin preview iframe'lerini `debug=1` ile ac
- [x] admin panel JS sentaks kontrolu calistir

## 2. Bu Tur Sonunda Kesinlesen Durum

- [x] OBS normal RTMP baglantisi calisiyor
- [x] OBS `Cok kanalli Video` icin gerekli Config Override JSON panelde veriliyor
- [x] stream olusturma ve stream detay ekranlarinda kullanici rehberi hazir
- [x] Linux VPS temiz kurulum durumuna getirildi
- [x] temiz kurulumda setup wizard yeniden aciliyor
- [x] `go test ./...` geciyor
- [x] `go build ./cmd/fluxstream/` geciyor
- [x] `go build ./cmd/fluxstream-license/` geciyor

## 3. Simdi Acik Kalan Kritik Teknik Isler

- [ ] ek video izlerini stream bellekte ayri ayri tut
- [ ] ek audio izlerini stream bellekte ayri ayri tut
- [ ] track metadata bilgisini API tarafina ac
- [ ] admin panelde gelen track listesi goster
- [ ] varsayilan video izi secimi ekle
- [ ] varsayilan audio izi secimi ekle
- [ ] dusuk bant genisligi icin ABR profil merdivenini olcumle optimize et

## 4. Gercek Multitrack ABR Faz

- [x] OBS'ten gelen kalite katmanlarini HLS varyantlarina dogrudan bagla
- [x] HLS master playlist icinde track kaynakli varyantlar yaz
- [x] DASH repack tarafinda tum video izlerini maple
- [x] RTMP chunk timestamp delta birikim hatasini kapat
- [x] `1080p` multitrack varyantinda mikro segment olusma kok nedenini kapat
- [x] bozuk `EXTINF` ve DASH `SegmentTimeline` uretimini kalici olarak duzelt
- [ ] transcode ile OBS varyantlarini karma kullan
- [ ] gereksiz yeniden encode maliyetini dusur
- [ ] track bazli bitrate / cozumunurluk analytics'i ekle
- [ ] cok kanalli audio track secimini player tarafina tasi

## 5. QoE ve Telemetri

- [x] player heartbeat telemetrisi topla
- [x] stall / toparlanma / reconnect bilgisini runtime bellekte sakla
- [x] admin stream detay ekraninda canli QoE karti goster
- [x] diagnostics ekraninda multitrack HLS / DASH sayaclarini goster
- [x] `bufferSeekOverHole` ve `bufferStalledError` davranisini canli testte ayristir
- [x] canli testte mikro segment kaynagini RTMP chunk timestamp zincirine kadar indir
- [ ] telemetrileri kalici depolamaya al
- [ ] telemetrileri grafik ve zaman serisi olarak goster
- [ ] Prometheus / OpenTelemetry cikisi uret

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
- [ ] admin panelde kalan `de`, `es`, `fr` ceviri eksiklerini kapat
- [x] multitrack ingest icin panelde kisa bilgilendirme ekle
- [x] OBS ayar orneklerini yardim ekranina ekle
- [ ] `production_status.md` dosyasini her ana faz sonunda guncel tut

## 8. Uretim Seviyesi Buyuk Eksikler

- [ ] multi-node origin-edge mimarisi
- [ ] S3 / MinIO archive ve restore akisi
- [ ] Prometheus / OpenTelemetry / alarm sistemi
- [ ] RBAC, audit log ve SSO
- [ ] DRM ve gelismis playback guvenligi
- [ ] SSAI ve monetizasyon omurgasi

## 9. Cekirdek Tamamlandiktan Sonra

- [ ] konferans odalari
- [ ] canli chat
- [ ] moderasyonlu soru-cevap
- [ ] sanal sinif rolleri ve yoklama
- [ ] breakout room mantigi
- [ ] takim ici mesajlasma
