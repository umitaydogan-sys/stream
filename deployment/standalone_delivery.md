# FluxStream Standalone Delivery Plan

Bu plan tek sunuculu kurulum icin yazildi. Ekstra origin sunucu gerekmez.

## 1. Temel Prensip

FluxStream tek basina calisir:

1. OBS -> `rtmp://sunucu:1935/live`
2. FluxStream -> HLS / LL-HLS / DASH / HTTP-FLV / MP4 / WebM / audio output
3. Izleyici -> dogrudan FluxStream

Bu modda harici cache, CDN veya ikinci sunucu zorunlu degildir.

## 2. Ne Zaman Ek Katman Gerekiyor

Tek sunucu senaryosunda ek katman sadece iki durumda faydalidir:

1. segment bazli izleyici sayisi artiyorsa
2. UI ve media isteklerini ayirmak istiyorsan

## 3. Ayni Makinede Opsiyonel Nginx Cache

Istersen ayni sunucuda Nginx reverse proxy calistirabilirsin.

Akis:

1. FluxStream `127.0.0.1:8844` uzerinde calisir
2. Nginx public portta dinler (`80` veya `443`)
3. Nginx:
   - `m3u8/mpd` isteklerini cache etmez
   - `ts/m4s/mp4` segmentlerini disk cache'e alir
   - `flv/mp4/webm/audio/icecast` gibi ham canli cikislari cache etmez

Hazir ornek konfig:

- [fluxstream-standalone.conf](C:/xampp/htdocs/stream/deployment/nginx/fluxstream-standalone.conf)

Not:

- Windows + XAMPP kullaniyorsan `80` portu baska servis tarafindan kullaniliyor olabilir
- bu durumda Nginx'i `8080` veya `8088` gibi bir porta al
- public erisim varsa daha sonra CDN origin olarak bu Nginx adresini tanimlayabilirsin

## 4. CDN Icin Ikinci Sunucu Gerekir mi

Hayir.

CDN kullanmak icin ikinci kendi sunucuna ihtiyacin yok. CDN mantigi su:

1. senin tek origin'in FluxStream veya onun onundeki Nginx olur
2. CDN bu origin'den segmentleri ceker
3. izleyiciye CDN servis eder

Yani:

- kendi tarafinda tek sunucu
- disarida opsiyonel CDN provider

Bu yuzden standalone kurulum bozulmaz. Sadece ihtiyac olursa ustune CDN eklenir.

## 5. Onerilen Kademeli Gecis

### Asama 1: Tam Standalone

- sadece `fluxstream.exe`
- public linkler dogrudan `:8844`

### Asama 2: Ayni Sunucuda Nginx

- FluxStream `127.0.0.1:8844`
- public trafik Nginx uzerinden
- segment cache aktif

### Asama 3: Opsiyonel CDN

- CDN origin = Nginx public domain
- manifest cache kapali
- segment cache acik

## 6. Hangi Cikislarda Cache Mantikli

Cache uygun:

- HLS segmentleri (`.ts`)
- DASH segmentleri (`.m4s`)
- static assetler (`/static/`)

Cache uygun degil:

- HLS manifest (`.m3u8`)
- DASH manifest (`.mpd`)
- HTTP-FLV
- MP4 live
- WebM live
- audio live stream endpoint'leri

## 7. Pratik Karar

Bugun icin en dogru kurulum:

1. FluxStream standalone calissin
2. tum ozellikler FluxStream uzerinden calismaya devam etsin
3. ihtiyac olursa ayni makinede Nginx cache katmani eklenebilsin
4. daha sonra trafik artarsa ayni origin'in onune CDN tanimlanabilsin

Bu model kullanicinin "tek sunucu" sinirini bozmaz.
