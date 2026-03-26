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
    E --> F[Embed Analitik ve ABR Studyosu]
    F --> G[Admin Studio V2]
    G --> H[Playback Guvenligi V2]
    H --> I[DRM]
    I --> J[Origin-Edge ve Enterprise Faz]
    J --> K[Konferans Chat Sanal Sinif]
```

## 1.1 Son Iki Fazda Nereye Geldik

```mermaid
flowchart TD
    S1[Embed Kodlari] --> S2[Embed Studyosu]
    S3[Basit Analitik] --> S4[Analitik Merkezi]
    S5[Ham ABR JSON] --> S6[ABR Profil Studyosu]
    S7[Temel Token Uretimi] --> S8[Playback Guvenligi V1]
    S9[Eski Admin Sayfalari] --> S10[Admin Studio V2]
    S11[Logo URL Alani] --> S12[Logo ve Marka Merkezi]
    S13[Teshis Ozeti] --> S14[Teshis ve Tedavi Merkezi]
```

## 1.2 Son Cekirdek Sertlestirme

```mermaid
flowchart LR
    A[Signed URL Omurgasi] --> B[V2 query signed URL zorlamasi]
    C[Audio-only DASH] --> D[Audio MPD ve audio MIME sertlestirmesi]
    E[Teshis] --> F[Audio-only DASH manifest gorunurlugu]
    G[Maintenance] --> H[Storage ile rol ayriminin netlesmesi]
    I[Marka Varliklari] --> J[Admin asset upload ve media-assets servisi]
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
    section Studio
      Embed ve Analitik Studyosu: 5
      Admin Studio V2: 5
      Logo ve Marka Merkezi: 4
    section Siradaki Faz
      Audio-only DASH saha testi: 3
      Playback Guvenligi V2: 2
      Harici storage dogrulamasi: 2
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
| Embed + Analitik + ABR Studyosu | Tamamlandi | embed, analitik ve ABR ekranlari urun seviyesine tasindi |
| Admin Studio V2 | Tamamlandi | dashboard, streams, ayarlar, diagnostics ve markalama sayfalari urunlesti |
| Harici storage saha testi | Kismen tamamlandi | ayni VPS uzerinde MinIO + SFTP test edildi, gercek S3 sirada |
| Playback guvenligi | Buyuk oranda tamamlandi | signed URL, token, domain/IP kisiti ve watermark omurgasi panel tarafina baglandi |
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
    Ingest ve HLS/DASH: [0.84, 0.91]
    OBS Multitrack: [0.80, 0.89]
    QoE ve Operasyon Merkezi: [0.79, 0.88]
    Recording ve Remux: [0.74, 0.85]
    Depolama ve Arsiv Merkezi: [0.73, 0.84]
    Embed ve Analitik Studyosu: [0.76, 0.89]
    Admin Studio V2: [0.71, 0.90]
    Playback Guvenligi V1: [0.58, 0.86]
    DRM: [0.12, 0.80]
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
- studio katmani sayesinde kritik admin ekranlari daha tutarli

### Hala Acik Olanlar

- harici AWS S3 saha testi
- rclone tabanli populer bulut hedeflerinin gercek hesaplarla dogrulanmasi
- `audio-only DASH` istemci saha testleri
- playback guvenligi V2
- DRM
- origin-edge
- RBAC, audit log ve SSO

## 5. Siradaki Yol

```mermaid
gantt
    title FluxStream Sonraki Fazlar
    dateFormat  YYYY-MM-DD
    axisFormat  %d %b
    section Sertlestirme
    Audio-only DASH saha testi          :active, dasha, 2026-03-27, 7d
    Playback guvenligi canli policy     :secv2, 2026-03-29, 8d
    section Storage
    Harici AWS S3 saha testi            :s3, 2026-04-03, 6d
    Drive / OneDrive / Dropbox test     :drv, 2026-04-06, 10d
    Recording soak ve restart testi     :rec, 2026-04-10, 8d
    section Guvenlik
    AES-128 ve key service              :drm1, 2026-04-18, 10d
    DRM abstraction                     :drm2, 2026-04-24, 12d
    section Mimari
    Origin-edge lite tasarimi           :oe1, 2026-05-08, 12d
    RBAC / audit / SSO                  :auth1, 2026-05-18, 14d
    section Sonrasi
    Konferans / chat / sanal sinif      :future, 2026-06-05, 20d
```

## 6. Once Neyi Bitirecegiz

1. `audio-only DASH` akisini gercek audio-only kaynak ve farkli istemcilerle sertlestirecegiz.
2. `Playback Guvenligi V1` politikasini canli stream policy senaryolariyla daha katı hale getirecegiz.
3. Harici AWS S3 ve populer bulut hedeflerinin gercek saha testlerini alacagiz.
4. Storage ve recording finalize akisini buyuk dosya ve restart senaryolariyla test edecegiz.
5. Sonra AES-128, DRM ve origin-edge lite tasarimina gececegiz.

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
