# VPS Deployment Report

Date: 2026-03-26
Host: `23.94.220.222`
Hostname: `host.kimediyoz.com.tr`
OS: `Ubuntu 24.04 LTS`

## Current Server State

- FluxStream Linux service installed at `/opt/fluxstream`
- systemd service name: `fluxstream`
- Current service health: `ok`
- Active listeners:
  - HTTP: `8844`
  - RTMP: `1935`
- UFW: inactive

## This Round

Bu turda sunucuya yeniden `temiz kurulum` yapildi.

Tamamlananlar:

1. Linux systemd paketi yeniden uretildi.
2. Paket sunucuya sifirdan kopyalandi.
3. Mevcut `fluxstream` servisi kaldirildi.
4. `/opt/fluxstream` temizlendi.
5. Paket icindeki `install.sh` ile yeniden kurulum yapildi.
6. `api/health` ve binary hash ile canli dogrulama alindi.

## Current Live Binary

- Linux SHA256:
  `2E18FF08103D166403832C5F57597567EF7F6910AB1BE6E534B4CA390D52570D`

## Validation

- `systemctl is-active fluxstream` => `active`
- `wget -qO- http://127.0.0.1:8844/api/health`
  => `{"status":"ok","uptime":...,"version":"2.0.0"}`
- `sha256sum /opt/fluxstream/fluxstream`
  => yerel package ile ayni hash

## Confirmed Product Areas On VPS

- OBS multitrack ingest zinciri
- HLS + DASH varyant uretimi
- QoE telemetry ve Operasyon Merkezi
- `Depolama ve Arsiv Merkezi`
- recording `mp4` varsayilani ve `TS capture + remux`
- MinIO + SFTP laboratuvar testi
- `Embed Studyosu`, `Analitik Merkezi`, `ABR Profilleri`
- `Admin Studio V2`
- stream bazinda sonradan `adaptive teslimat` akisi

## Next Useful VPS Checks

1. `adaptive teslimat` icin `live_now` akisinin canli stream uzerinde gozlenmesi
2. harici AWS S3 bucket ile gercek recording + backup upload / restore
3. `audio-only DASH` davranisinin farkli istemcilerde dogrulanmasi
4. playback guvenligi politikasinin canli stream policy senaryolariyla sertlestirilmesi
