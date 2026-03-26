# FluxStream Yol Haritasi

Tarih: 26 Mart 2026

Bu dokuman, cekirdegin nereden nereye geldigini ve bundan sonra hangi
fazlara girecegini tek yerde gormek icin hazirlandi.

## 1. Buyuk Resim

```mermaid
flowchart LR
    A[Ham Streaming Cekirdegi] --> B[Player ve Embed Kararliligi]
    B --> C[OBS Multitrack ve ABR]
    C --> D[QoE ve Operasyon Merkezi]
    D --> E[Kayit Arsiv ve Yedek Merkezi]
    E --> F[Playback Guvenligi]
    F --> G[DRM]
    G --> H[Origin-Edge ve Enterprise Faz]
    H --> I[Konferans Chat Sanal Sinif]
```

## 2. Nereden Nereye Geldik

### 2.1 Yolculuk Ozeti

```mermaid
journey
    title FluxStream Yolculugu
    section Baslangic
      Tek binary cekirdek: 4
      Temel admin panel: 4
      Stream CRUD ve basit player: 3
    section Kararlilik
      Preview ve embed sorunlarini kapatma: 5
      QoE debug ve telemetry: 5
      Operasyon Merkezi: 5
    section Multitrack
      OBS multitrack parse: 5
      HLS/DASH ABR zinciri: 5
      Mikro segment kok nedenini kapatma: 5
    section Depolama
      Recording remux modeli: 5
      Depolama ve Arsiv Merkezi: 5
      MinIO ve SFTP saha testi: 4
    section Siradaki Faz
      Playback guvenligi: 2
      DRM: 1
      Origin-edge: 1
```

### 2.2 Milestone Tablosu

| Faz | Durum | Kisa Not |
|---|---|---|
| Temel ingest ve dagitim | Tamamlandi | HLS, DASH, recording ve admin panel omurgasi oturdu |
| Player ve preview kararliligi | Tamamlandi | embed, iframe, direct link ve offline hatalari kapandi |
| OBS multitrack ve ABR | Tamamlandi | HLS/DASH varyant zinciri ve chunk timestamp kok nedeni kapandi |
| QoE ve Operasyon Merkezi | Tamamlandi | telemetry, track analytics, Prometheus ve teshis ekrani var |
| Depolama ve Arsiv Merkezi | Buyuk oranda tamamlandi | kayit, arsiv, yedek ve bulut hedefleri tek merkezde |
| Harici storage saha testi | Kismen tamamlandi | ayni VPS uzerinde MinIO + SFTP test edildi, gercek S3 sirada |
| Playback guvenligi | Baslamadi | signed URL, token, hotlink, watermark fazi acik |
| DRM | Baslamadi | AES-128, DRM abstraction ve lisans servisleri acik |
| Origin-edge / cluster | Baslamadi | dusuk butceye uygun lite model tasarlanacak |
| Konferans / chat / sanal sinif | Baslamadi | cekirdek streaming tarafi tamamen oturduktan sonra |

## 3. Bugunku Mimari Olgunluk

```mermaid
quadrantChart
    title FluxStream Bilesen Olgunlugu
    x-axis Dusuk olgunluk --> Yuksek olgunluk
    y-axis Dusuk urun degeri --> Yuksek urun degeri
    quadrant-1 Buyut
    quadrant-2 Guclu Alan
    quadrant-3 Arka Plan
    quadrant-4 Sertlestir
    Ingest ve HLS/DASH: [0.82, 0.90]
    OBS Multitrack: [0.78, 0.88]
    QoE ve Operasyon Merkezi: [0.76, 0.86]
    Recording ve Remux: [0.70, 0.84]
    Depolama ve Arsiv Merkezi: [0.68, 0.82]
    Playback Guvenligi: [0.20, 0.86]
    DRM: [0.10, 0.80]
    Origin-edge: [0.12, 0.88]
    RBAC ve SSO: [0.15, 0.72]
```

## 4. Bugunku Gercek Durum

### Guclu Alanlar

- tek-node canli dagitim cekirdegi artik guven veriyor
- OBS multitrack ve ABR omurgasi calisiyor
- QoE, track ve manifest gorunurlugu var
- recording ve depolama tarafi urun hissi vermeye basladi
- ayni urunde admin panel, operasyon, yedek ve arsiv bir arada

### Hala Acik Olanlar

- harici AWS S3 saha testi
- rclone tabanli populer bulut hedeflerinin gercek hesaplarla dogrulanmasi
- `audio-only DASH` istemci saha testleri
- playback guvenligi
- DRM
- origin-edge
- RBAC, audit log ve SSO

## 5. Siradaki Yol

```mermaid
gantt
    title FluxStream Sonraki Fazlar
    dateFormat  YYYY-MM-DD
    axisFormat  %d %b
    section Kisa Vade
    Depolama UX sadeleştirme         :active, ux1, 2026-03-26, 10d
    Harici AWS S3 saha testi         :s3, 2026-03-28, 7d
    Drive / OneDrive / Dropbox test  :drv, 2026-03-31, 10d
    Audio-only DASH sertlestirme     :dasha, 2026-04-02, 8d
    Recording soak ve restart testi  :rec, 2026-04-05, 10d
    section Guvenlik
    Signed URL ve token policy       :sec1, 2026-04-10, 12d
    Hotlink / watermark / rate limit :sec2, 2026-04-15, 10d
    AES-128 ve key service           :drm1, 2026-04-24, 12d
    section Mimari
    Origin-edge lite tasarimi        :oe1, 2026-05-02, 12d
    RBAC / audit / SSO               :auth1, 2026-05-08, 14d
    section Sonrasi
    Konferans / chat / sanal sinif   :future, 2026-05-20, 20d
```

## 6. Once Neyi Bitirecegiz

1. `Depolama ve Arsiv Merkezi`ni daha da sade hale getirecegiz.
2. Harici AWS S3 bucket ile gercek dis ortam dogrulamasi yapacagiz.
3. Drive / OneDrive / Dropbox gibi hedefleri gercek hesaplarla test edecegiz.
4. `audio-only DASH` ve recording finalize davranisini sertlestirecegiz.
5. Ardindan signed playback security fazina girecegiz.
6. Sonra DRM ve origin-edge mimarisine gececegiz.

## 7. Bu Cekirdegin Uzerine Neler Insa Edilebilir

### Bugunden Yarin Cikabilecek Urunler

- kurumsal TV ve kurum ici yayin platformu
- radyo ve audio streaming platformu
- markali webcast ve webinar urunu
- arsiv / catch-up ve VOD portali
- egitim yayini ve sinif ici canli ders omurgasi

### Cekirdek Tamamlandiktan Sonra

- konferans odalari
- canli chat
- moderasyonlu soru-cevap
- sanal sinif rolleri
- yoklama
- breakout room
- takim ici mesajlasma

## 8. Tek Cumlelik Ozet

FluxStream, ham bir streaming denemesi olmaktan cikti; artik iyi bir
tek-node medya sunucusu ve urunlesmeye yaklasmis bir yayin cekirdegi.
