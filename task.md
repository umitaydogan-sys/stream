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
- [ ] player tarafina QoE / stall telemetry katmani ekle
- [ ] dusuk bant genisligi icin ABR profil merdivenini olcumle optimize et

## 4. Gercek Multitrack ABR Faz

- [ ] OBS'ten gelen kalite katmanlarini HLS varyantlarina dogrudan bagla
- [ ] HLS master playlist icinde track kaynakli varyantlar yaz
- [ ] transcode ile OBS varyantlarini karma kullan
- [ ] gereksiz yeniden encode maliyetini dusur
- [ ] track bazli bitrate / cozumunurluk analytics'i ekle
- [ ] cok kanalli audio track secimini player tarafina tasi

## 5. Urunlestirme ve Lisans

- [x] runtime lisans modeli ABR / RTMPS / recording / branding tarafina baglandi
- [x] Linux servis yonetimi panel ve CLI tarafinda hazirlandi
- [x] backup / restore / upgrade plani olusturuldu
- [ ] `max_nodes` enforcement tamamla
- [ ] maintenance expiry ve grace policy ekle
- [ ] lisans modelini ilerideki cluster mimarisi ile uyumlu hale getir
- [ ] rollback guvenli Linux upgrade akisini sertlestir
- [ ] `.deb` / paketli Linux dagitimini tamamla

## 6. Arayuz ve Dokumantasyon

- [x] dokuman dosyalarini Turkce tut
- [ ] admin panelde kalan `de`, `es`, `fr` ceviri eksiklerini kapat
- [x] multitrack ingest icin panelde kisa bilgilendirme ekle
- [x] OBS ayar orneklerini yardim ekranina ekle
- [ ] `production_status.md` dosyasini her ana faz sonunda guncel tut

## 7. Uretim Seviyesi Buyuk Eksikler

- [ ] multi-node origin-edge mimarisi
- [ ] S3 / MinIO archive ve restore akisi
- [ ] Prometheus / OpenTelemetry / alarm sistemi
- [ ] RBAC, audit log ve SSO
- [ ] DRM ve gelismis playback guvenligi
- [ ] SSAI ve monetizasyon omurgasi

## 8. Cekirdek Tamamlandiktan Sonra

- [ ] konferans odalari
- [ ] canli chat
- [ ] moderasyonlu soru-cevap
- [ ] sanal sinif rolleri ve yoklama
- [ ] breakout room mantigi
- [ ] takim ici mesajlasma
