# Raspberry Pi Security Checklist

Deployment hardening for running mosaic-app on a Raspberry Pi exposed to the internet.

Application-side controls are already implemented (see [Implemented in the app](#implemented-in-the-app)). This document covers the remaining infrastructure and optional features to add later.

## Implemented in the app

These are configured via `.env` (see `env.example`):

| Control | Env vars | Default |
| --- | --- | --- |
| Max upload size | `MAX_FILE_SIZE` | 10 MB |
| Max image dimensions | `MAX_IMAGE_WIDTH`, `MAX_IMAGE_HEIGHT` | 4096×4096 |
| Tile size bounds | `MIN_TILE_SIZE`, `MAX_TILE_SIZE` | 5–100 |
| Per-IP rate limit | `RATE_LIMIT_REQUESTS`, `RATE_LIMIT_WINDOW` | 10 / hour |
| Concurrent generations | `MAX_CONCURRENT_GENERATIONS` | 2 |
| Request timeout | `REQUEST_TIMEOUT` | 120s |
| Generic error messages | `PRODUCTION=true` | off |

User uploads are processed in memory only — nothing is written to disk.

## Recommended Pi `.env`

```bash
SERVER_PORT=8080
MAX_FILE_SIZE=5242880
MAX_IMAGE_WIDTH=2048
MAX_IMAGE_HEIGHT=2048
MIN_TILE_SIZE=5
MAX_TILE_SIZE=100
MAX_CONCURRENT_GENERATIONS=1
RATE_LIMIT_REQUESTS=5
RATE_LIMIT_WINDOW=3600
REQUEST_TIMEOUT=90
PRODUCTION=true
TILES_DIR=tiles
LOG_LEVEL=info
```

---

## TODO: Network layer

### 1. Bind the app to localhost only

Do not expose the Go server directly on `0.0.0.0`. Run it on `127.0.0.1:8080` and put a reverse proxy in front.

**Option A — systemd `Environment`:**

```ini
Environment=SERVER_PORT=8080
# Listen on localhost by changing main.go or using a reverse proxy that targets 127.0.0.1:8080
```

**Option B — reverse proxy only** (see below): proxy to `http://127.0.0.1:8080`.

### 2. Firewall (`ufw`)

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow OpenSSH
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable
sudo ufw status
```

Do **not** open port 8080 to the WAN. Only the reverse proxy (80/443) should be reachable.

### 3. SSH hardening

- [ ] Disable password login; use SSH keys only
- [ ] Consider a non-default SSH port (optional)
- [ ] Keep the Pi updated: `sudo apt update && sudo apt upgrade`

---

## TODO: Reverse proxy + TLS

Use Caddy (simplest) or Nginx in front of the app.

### Caddy (recommended)

Install Caddy, then `/etc/caddy/Caddyfile`:

```caddyfile
mosaic.example.com {
    reverse_proxy 127.0.0.1:8080

    # Optional: stricter body limit at the edge (5 MB)
    request_body {
        max_size 5MB
    }

    # Optional: extra rate limit at the edge
    rate_limit {
        zone upload {
            key {remote_host}
            events 5
            window 1h
        }
    }
}
```

Caddy obtains and renews Let's Encrypt certificates automatically.

### Nginx (alternative)

```nginx
server {
    listen 443 ssl http2;
    server_name mosaic.example.com;

    ssl_certificate     /etc/letsencrypt/live/mosaic.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/mosaic.example.com/privkey.pem;

    client_max_body_size 5M;

    location /api/file/upload {
        limit_req zone=upload burst=2 nodelay;
        proxy_pass http://127.0.0.1:8080;
        proxy_read_timeout 120s;
    }

    location / {
        proxy_pass http://127.0.0.1:8080;
    }
}

limit_req_zone $binary_remote_addr zone=upload:10m rate=5r/h;
```

---

## TODO: systemd service hardening

Create `/etc/systemd/system/mosaic.service`:

```ini
[Unit]
Description=Mosaic Generator
After=network.target

[Service]
Type=simple
User=mosaic
Group=mosaic
WorkingDirectory=/opt/mosaic-app
EnvironmentFile=/opt/mosaic-app/.env
ExecStart=/opt/mosaic-app/mosaic
Restart=on-failure
RestartSec=5

# Limit blast radius of a runaway process
MemoryMax=512M
CPUQuota=80%

# Reduce privileges
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/mosaic-app/tiles

[Install]
WantedBy=multi-user.target
```

Setup:

```bash
sudo useradd --system --no-create-home mosaic
sudo chown -R mosaic:mosaic /opt/mosaic-app
sudo systemctl daemon-reload
sudo systemctl enable --now mosaic
```

---

## TODO: CAPTCHA (optional)

Useful if the site is public and bots become a problem. App-side rate limiting is the first line; CAPTCHA adds friction for automated abuse.

### Cloudflare Turnstile (free)

1. Create a site at [Cloudflare Turnstile](https://dash.cloudflare.com/?to=/:account/turnstile)
2. Add the widget to the frontend before the Generate button
3. Send the token with the upload request (e.g. `captchaToken` form field)
4. Verify server-side before calling `generateMosaic`:

```go
// POST https://challenges.cloudflare.com/turnstile/v0/siteverify
// Body: secret=<SECRET>&response=<token>
```

Reject uploads when verification fails (`403 Forbidden`).

### hCaptcha

Same pattern: widget in frontend, server-side verification before processing.

---

## TODO: Private access instead of public internet

Best option if only you and a few people need access:

### Tailscale (easiest)

1. Install Tailscale on the Pi and your devices
2. Do **not** port-forward 80/443 on your router
3. Access the app via the Pi's Tailscale IP (e.g. `http://100.x.x.x:8080`)

No public exposure = minimal attack surface.

### WireGuard

Same idea with a self-hosted VPN. More setup, full control.

---

## TODO: Cloudflare proxy (optional)

If using a public domain:

1. Point DNS to Cloudflare (orange cloud / proxied)
2. Enable **Bot Fight Mode** or **Super Bot Fight Mode**
3. Add a **WAF rate limiting rule** on `POST /api/file/upload`
4. Use **Full (strict)** SSL between Cloudflare and your origin (Caddy/Nginx with a valid cert)

Absorbs volumetric traffic before it hits the Pi.

---

## TODO: Monitoring

- [ ] Watch CPU/RAM with `htop` or `btop`
- [ ] Log rotation for the mosaic service (`journalctl -u mosaic -f`)
- [ ] Alert if CPU stays at 100% (simple cron + email, or Uptime Kuma / healthchecks.io on `/api/health`)
- [ ] Disk space on the Pi (tile library growth)

---

## Priority order

| Priority | Task | Effort |
| --- | --- | --- |
| 1 | Reverse proxy + TLS (Caddy) | Low |
| 2 | Firewall — block 8080 from WAN | Low |
| 3 | systemd service + non-root user + memory limit | Low |
| 4 | Tight Pi `.env` (see above) | Low |
| 5 | Tailscale instead of public exposure (if audience is small) | Low |
| 6 | Cloudflare proxy + WAF | Medium |
| 7 | CAPTCHA on upload | Medium |
| 8 | Monitoring / alerts | Medium |

---

## Threat model summary

| Threat | Mitigation |
| --- | --- |
| CPU exhaustion (many uploads) | App rate limit + concurrency cap + proxy limits |
| Memory exhaustion (huge images) | Dimension + file size caps |
| Expensive tileSize abuse | Server-side 5–100 bounds |
| Stored user photos | Not stored (in-memory only) |
| MITM | TLS via Caddy/Nginx |
| Direct port scanning | Firewall; app on localhost only |
| Bot floods | Rate limit → CAPTCHA → Cloudflare |
| Pi crash from one bad request | systemd `MemoryMax`, `MAX_CONCURRENT_GENERATIONS=1` |

---

## References

- App config: `env.example`
- Middleware: `middleware/security.go`
- Upload validation: `handlers.go`, `validation.go`
