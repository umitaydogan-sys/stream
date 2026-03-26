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
    E --> F[Embed ve Analitik Studyo Fazı]
    F --> G[Playback Guvenligi]
    G --> H[DRM]
    H --> I[Origin-Edge ve Enterprise Faz]
    I --> J[Konferans Chat Sanal Sinif]
```

## 1.1 Son Fazda Nereye Geldik

```mermaid
flowchart TD
    S1[Embed Kodlari] --> S2[Embed Studyosu]
    S3[Basit Analitik] --> S4[Analitik Merkezi]
    S5[Ham ABR JSON] --> S6[ABR Profil Studyosu]
    S7[Temel Token Uretimi] --> S8[Playback Guvenligi V1]
    S9[Audio-only DASH Cekirdegi] --> S10[Audio teslimat gorunurlugu]
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
      Embed ve Analitik Studyo: 3
      ABR Profil Merkezi: 3
      Audio-only DASH sertlestirme: 3
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
| Embed + Analitik + ABR Studyo | Tamamlandi | embed, analitik ve ABR ekranlari urun seviyesine tasindi |
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
    Ingest ve HLS/DASH: [0.82, 0.90]
    OBS Multitrack: [0.78, 0.88]
    QoE ve Operasyon Merkezi: [0.76, 0.86]
    Recording ve Remux: [0.70, 0.84]
    Depolama ve Arsiv Merkezi: [0.68, 0.82]
    Embed ve Analitik Studyo: [0.28, 0.88]
    ABR Profil Merkezi: [0.24, 0.86]
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
- `Embed Stüdyosu`, `Analitik Merkezi` ve `ABR Profilleri` ekranlarinin urun seviyesine tasinmasi
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
    section Urunlestirme
    Embed Studyosu                   :active, emb1, 2026-03-26, 10d
    Analitik Merkezi                 :ana1, 2026-03-29, 10d
    ABR Profil Merkezi               :abr1, 2026-04-01, 10d
    Audio-only DASH sertlestirme     :dasha, 2026-04-05, 8d
    section Guvenlik
    Signed URL ve token policy       :sec1, 2026-04-10, 12d
    Hotlink / watermark / rate limit :sec2, 2026-04-14, 10d
    section Depolama Sertlestirme
    Harici AWS S3 saha testi         :s3, 2026-04-18, 7d
    Drive / OneDrive / Dropbox test  :drv, 2026-04-21, 10d
    Recording soak ve restart testi  :rec, 2026-04-24, 10d
    section DRM
    AES-128 ve key service           :drm1, 2026-05-02, 12d
    section Mimari
    Origin-edge lite tasarimi        :oe1, 2026-05-12, 12d
    RBAC / audit / SSO               :auth1, 2026-05-18, 14d
    section Sonrasi
    Konferans / chat / sanal sinif   :future, 2026-06-01, 20d
```

## 6. Once Neyi Bitirecegiz

1. `Embed Stüdyosu` ekranini urun seviyesine tasiyacagiz.
2. `Analitik Merkezi` ekranini daha guclu KPI ve grafiklerle yeniden kuracagiz.
3. `ABR Profilleri ve Teslimat Merkezi`ni form tabanli profil mantigina gecirecegiz.
4. Ayni faz icinde `audio-only DASH` istemci sertlestirmesini kapatacagiz.
5. Ayni faz icinde playback guvenligi v1 katmanini ekleyecegiz.
6. Sonra harici AWS S3 ve populer bulut hedeflerinin gercek saha testlerine donecegiz.
7. Daha sonra DRM ve origin-edge mimarisine gececegiz.

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
