# FluxStream

Tek binary ile calisan, admin paneli bulunan, HLS/DASH merkezli canli yayin
sunucusu.

## Bugun Neler Var

- RTMP, RTMPS, SRT, RTP, RTSP, WebRTC/WHIP, MPEG-TS ve HTTP Push ingest
- HLS, LL-HLS, DASH, HTTP-FLV, MP4, WebM ve audio cikislari
- OBS multitrack video ve audio dagitimi
- `Embed Studyosu`, `Gelişmis Embed`, `Player Sablonlari Studyosu`
- `Analitik Merkezi`, `Operasyon Merkezi`, `Teshis ve Tedavi Merkezi`
- `ABR Profilleri ve Teslimat Merkezi`
- sonradan `adaptive teslimat` acma akisi
- recording, archive, backup ve `Depolama ve Arsiv Merkezi`
- playback guvenligi V1
- `Admin Studio V2`

## Kaynak Kod

- `https://github.com/umitaydogan-sys/stream`

## Temel Dokumanlar

- [Uygulama Plani](implementation_plan.md)
- [Gorev Listesi](task.md)
- [Production Durumu](production_status.md)
- [Yol Haritasi](yol_haritasi.md)
- [Surec Kaydi](surec_kaydi_2026-03-24.md)

## Bugunku Konum

FluxStream su anda iyi bir tek-node medya sunucusu ve urunlesmis bir yayin
cekirdegi seviyesindedir.

En mantikli sonraki fazlar:

- `audio-only DASH` saha sertlestirmesi
- harici AWS S3 / Drive / OneDrive / Dropbox testleri
- playback guvenligi V2
- DRM hazirligi
- origin-edge lite
