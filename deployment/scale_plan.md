# FluxStream Scale Plan (Origin + Edge Cache + CDN)

Bu plan, tek origin sunucudan binlerce eşzamanlı izleyiciye çıkış almak için minimum güvenli mimariyi verir.

## 1. Hedef Mimari

1. Encoder (OBS) -> FluxStream Origin (RTMP ingest)
2. FluxStream Origin -> HLS/DASH segment üretimi
3. Nginx Edge (reverse proxy + segment cache)
4. CDN (Cloudflare/CloudFront/Fastly) -> izleyiciler

## 2. Neden Bu Mimari

1. Origin CPU ve disk I/O baskısını azaltır.
2. Segment dosyalarının tekrar tekrar origin’den çekilmesini engeller.
3. Coğrafi olarak dağıtılmış izleme deneyimi sağlar.
4. Ani trafik sıçramalarında (raid/spike) sistemi ayakta tutar.

## 3. Nginx Edge Kurulumu

1. `deployment/nginx/fluxstream-edge.conf` dosyasını edge sunucuda etkinleştir.
2. `proxy_pass` origin adresini (FluxStream) doğru porta ayarla.
3. `/hls/*.m3u8` ve `/dash/*.mpd` için cache kapalı kalmalı.
4. Segment dosyaları (`.ts`, `.m4s`, `.mp4`) için cache açık olmalı.

## 4. CDN Kuralları

1. Manifest dosyaları:
   - Cache: `no-store` / bypass
2. Segment dosyaları:
   - Cache TTL: 10-30 dakika
   - Stale-while-revalidate: açık
3. CORS:
   - `Access-Control-Allow-Origin: *` (public yayın için)

## 5. Origin Sertleştirme

1. RTMP ingest’i sadece encoder IP’lerine aç.
2. İzleyici trafiğini origin’e direkt verme; sadece edge/CDN üzerinden yayınla.
3. Disk:
   - HLS dizinini hızlı SSD/NVMe üzerinde tut.
4. İzleme:
   - CPU, disk IOPS, net throughput, 5xx, player error rate.

## 6. Player Tarafı (No External Dependency Hedefi)

Şu an bazı player path’leri CDN JS kütüphanesi kullanıyor. Dış bağımlılığı sıfırlamak için:

1. `hls.js`, `dash.js`, `mpegts.js` dosyalarını proje içine vendor olarak ekle.
2. Bu dosyaları `internal/web` tarafından yerel endpoint ile servis et.
3. HTML tarafında script URL’lerini sadece local path’e çevir.

Not: Bu adım mimari olarak zorunlu değil, operasyonel bağımsızlık için önerilir.

## 7. Performans Kabul Kriteri

1. 10 dakika sürekli yayında player donması olmamalı.
2. Ortalama canlı gecikme:
   - HLS: 6-12 saniye
   - LL-HLS: 2-5 saniye
3. 5xx oranı < %0.1
4. Rebuffer ratio < %1

