# VPS Deployment Report

Date: 2026-03-12
Host: `23.94.220.222`
Hostname: `host.kimediyoz.com.tr`
OS: `Ubuntu 24.04 LTS`

## Current Server State

- FluxStream Linux service installed at `/opt/fluxstream`
- systemd service name: `fluxstream`
- Service user: `fluxstream`
- Admin OS user created: `fluxadmin`
- Current service health: `ok`
- Active listeners after reset:
  - HTTP: `8844`
  - RTMP: `1935`
- UFW: inactive

## Completed in This Round

1. Linux package deployed and validated on VPS.
2. Separate TLS profile support added:
   - Web HTTPS profile
   - Stream RTMPS profile
3. Manual CRT/KEY upload path split:
   - Web cert path default: `data/certs/web/server.crt|key`
   - Stream cert path default: `data/certs/stream/server.crt|key`
4. Let's Encrypt runtime support added:
   - Web domain config
   - Stream domain config
   - HTTP-01 challenge listener on port `80`
5. Linux systemd package updated:
   - no BOM shell scripts
   - `CAP_NET_BIND_SERVICE` for low ports

## VPS Validation Results

### Base Install

- `systemctl status fluxstream` => running
- `curl http://127.0.0.1:8844/api/health` => `{"status":"ok"...}`

### Manual Web + Stream TLS

Validated with temporary self-signed certificates:

- HTTPS listener came up on `443`
- RTMPS listener came up on `1936`
- `curl -k https://127.0.0.1/api/health` => success
- `openssl s_client -connect 127.0.0.1:1936` => success

### Let's Encrypt Runtime Mode

Validated infrastructure only:

- HTTP-01 challenge listener came up on `80`
- HTTPS listener came up on `443`
- RTMPS listener came up on `1936`
- `/api/ssl/status` returned separate web/stream LE config state

Note:
- Actual public certificate issuance was not completed in this round.
- For real issuance, DNS must point to this VPS.
- If stream uses a different domain than web, that stream domain must also have a valid DNS record.

## Important DNS Note

If you want separate certificates:

- Web UI / embed domain example:
  - `host.kimediyoz.com.tr`
- Stream RTMPS domain example:
  - `stream.kimediyoz.com.tr`

Both domains must resolve to `23.94.220.222`.

## Recommended Next Development Order

1. i18n for installer + setup wizard + admin panel
2. analytics chart layout cleanup
3. recording preview inside records page
4. license subsystem
5. Linux-specific admin actions and packaging polish

