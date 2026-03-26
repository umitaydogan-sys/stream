# FluxStream

Tek binary ile calisan, admin paneli bulunan, HLS/DASH merkezli
canli yayin sunucusu.

## Bugun Neler Var

- RTMP, RTMPS, SRT, RTP, RTSP, WebRTC/WHIP, MPEG-TS ve HTTP Push ingest
- HLS, LL-HLS, DASH, HTTP-FLV, MP4, WebM ve ses cikislari
- OBS multitrack video ve audio dagitim zinciri
- admin paneli, setup wizard, stream CRUD, player/embed/template sistemi
- `Operasyon Merkezi`
- QoE, telemetry, Prometheus ve OTel-benzeri cikis
- recording, archive, backup ve `Depolama ve Arsiv Merkezi`
- Linux servis, deploy ve backup/restore omurgasi

## Bu Repo Nerede

Kaynak kod:

- `https://github.com/umitaydogan-sys/stream`

## Temel Dokumanlar

- [Uygulama Plani](implementation_plan.md)
- [Gorev Listesi](task.md)
- [Production Durumu](production_status.md)
- [Yol Haritasi](yol_haritasi.md)
- [Surec Kaydi](surec_kaydi_2026-03-24.md)

## Bugunku Konum

FluxStream su anda iyi bir tek-node medya sunucusu seviyesindedir.
Ozellikle:

- webcast
- kurum ici TV
- radyo
- markali player ve embed dagitimi

icin kullanilabilir durumdadir.

Tam enterprise seviye icin sonraki buyuk fazlar:

- playback guvenligi ve DRM
- harici storage ve failover sertlestirmesi
- origin-edge / cluster
- RBAC, audit log ve SSO
