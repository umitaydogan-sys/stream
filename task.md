# FluxStream Gorev Listesi

Tarih: 26 Mart 2026

## 0. Bu Turda Kapananlar

- [x] tek kalite baslayan bir streami sonradan `adaptive teslimat` moduna alma ozelligini ekle
- [x] `Streams` ekranina `Adaptiveye Al` hizli aksiyonunu ekle
- [x] `Stream Detayi` ekranina `Adaptive Teslimat` karti ekle
- [x] `Sonraki yayinda etkinlestir` ve `Canli yayina hemen uygula` akisini ekle
- [x] `balanced`, `mobile`, `resilient`, `radio` profil secimini urun akisina bagla
- [x] `ABR Profilleri ve Teslimat Merkezi` icine genel adaptive ac/kapa kontrolunu geri ekle
- [x] ayni ekrana `HTTPS` ve `RTMPS` secure teslimat hizli ayarlarini ekle
- [x] kayitli profillerin yanina hazir preset secbox + `Preseti Yukle` akisini ekle
- [x] Linux package'i yeniden uret
- [x] Windows portable package'i yeniden uret
- [x] Windows service package'i yeniden uret
- [x] Windows installer'i yeniden uret
- [x] VPS'e temiz kurulum yap
- [x] temel md dokumanlarini yeni surume gore guncelle
- [x] kokte kalan gereksiz build artefaktlarini temizle
- [x] `Gelismis Embed` ust kart metinlerini musteriye uygun teknik dile cek
- [x] tum panelde `input / select / textarea` yuzlerini daha kompakt ve daha net hale getir
- [x] HTTPS embed preview ve SSL durum kartlarini gercek listener durumu ile hizala

## 1. Daha Once Kapanan Buyuk Fazlar

- [x] `Embed Studyosu`
- [x] `Gelismis Embed`
- [x] `Player Sablonlari Studyosu`
- [x] `Analitik Merkezi`
- [x] `ABR Profilleri ve Teslimat Merkezi`
- [x] `Playback Guvenligi V1`
- [x] `Operasyon Merkezi`
- [x] `Depolama ve Arsiv Merkezi`
- [x] `Admin Studio V2`
- [x] `Logo ve Marka Merkezi`
- [x] `Teshis ve Tedavi Merkezi`

## 2. Kisa Vade Acik Isler

- [ ] `adaptive teslimat` icin `live_now` akisinin canli testini al
- [ ] `audio-only DASH` akisini tarayici, dash.js ve VLC tarafinda saha testinden gecir
- [ ] playback guvenligi V1'i domain / referrer / IP / token zorlamasi ile canli policy senaryolarinda test et
- [ ] harici AWS S3 bucket ile gercek recording + backup upload / restore dogrulamasi yap
- [ ] rclone tabanli `Google Drive`, `OneDrive` ve `Dropbox` hedeflerini gercek hesaplarla test et
- [ ] buyuk dosya ve uzun sureli recording remux dayanikliligini tekrar sertlestir
- [ ] onceki bozuk `TS` kayitlar icin kurtarma / uyari akisi ekle

## 3. Orta Vade Isler

- [ ] playback guvenligi V2 presetleri ekle
- [ ] audit log ve guvenlik olay kaydi ekle
- [ ] AES-128 HLS key service ekle
- [ ] DRM abstraction katmani tasarla
- [ ] RBAC, audit ve SSO backlog'unu ac

## 4. Buyuk Urun Eksikleri

- [ ] origin-edge lite
- [ ] multi-node cluster
- [ ] failover ve yuk testi
- [ ] uzun sureli soak test
- [ ] enterprise seviye policy ve lisans baglantisi

## 5. Cekirdek Tamamlandiktan Sonra

- [ ] konferans odalari
- [ ] canli chat
- [ ] moderasyonlu soru-cevap
- [ ] sanal sinif rolleri
- [ ] breakout room
- [ ] takim ici mesajlasma
