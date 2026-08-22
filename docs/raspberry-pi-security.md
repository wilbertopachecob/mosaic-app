# Raspberry Pi Security Checklist (Cloudflare Tunnel + Go)

Deployment hardening for running mosaic-app on a Raspberry Pi behind Cloudflare Tunnel.

Application-side controls are already implemented (see [Implemented in the app](#implemented-in-the-app)). This document focuses on infrastructure hardening for this specific context:

- Go backend serving HTTP locally on the Pi
- `cloudflared` exposing the app via Cloudflare Tunnel
- No direct public exposure of the Go service port

Cloudflare Tunnel works fine with a Go backend. The tunnel forwards HTTP traffic to your local Go service, regardless of backend language.

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

## Network layer (Tunnel-first profile)

### 1. Bind the app to localhost only

Do not expose the Go server directly to the internet. Keep it reachable only from localhost and route traffic through Cloudflare Tunnel.

**Recommended local path:**

```ini
Environment=SERVER_PORT=8080
# cloudflared routes to http://127.0.0.1:8080
```

In `/etc/cloudflared/config.yml`, route the hostname to the app:

```yaml
tunnel: your-tunnel-name
credentials-file: /etc/cloudflared/<tunnel-id>.json
ingress:
  - hostname: mosaic-generator.wilbertopachecob.dev
    service: http://127.0.0.1:8080
  - service: http_status:404
```

You can keep using one tunnel for multiple apps on the same Pi, for example:

- `paint.wilbertopachecob.dev` -> `http://127.0.0.1:3000`
- `tictactoe.wilbertopachecob.dev` -> `http://127.0.0.1:3001`
- `mosaic-generator.wilbertopachecob.dev` -> `http://127.0.0.1:8080`

### 2. Firewall (`ufw`) for tunnel setup

```bash
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow OpenSSH
sudo ufw enable
sudo ufw status
```

For Cloudflare Tunnel origin mode, you generally do **not** need inbound `80/443` open on the Pi. Keep port `8080` closed publicly.

### 3. SSH hardening

- [ ] Disable password login; use SSH keys only
- [ ] Consider a non-default SSH port (optional)
- [ ] Keep the Pi updated: `sudo apt update && sudo apt upgrade`

---

## Reverse proxy + TLS (optional fallback)

Use Caddy or Nginx only if you decide to expose the Pi directly (no tunnel). With Cloudflare Tunnel, `cloudflared` is typically enough.

### Caddy (recommended)

Install Caddy, then `/etc/caddy/Caddyfile`:

```caddyfile
mosaic-generator.wilbertopachecob.dev {
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

Caddy obtains and renews Let's Encrypt certificates automatically in direct-origin mode.

### Nginx (alternative)

```nginx
server {
    listen 443 ssl http2;
    server_name mosaic-generator.wilbertopachecob.dev;

    ssl_certificate     /etc/letsencrypt/live/mosaic-generator.wilbertopachecob.dev/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/mosaic-generator.wilbertopachecob.dev/privkey.pem;

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

## systemd service hardening + reboot persistence

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

### Ensure both services auto-start after Pi reboot

```bash
# App service
sudo systemctl enable mosaic
sudo systemctl restart mosaic

# Tunnel service
sudo systemctl enable cloudflared
sudo systemctl restart cloudflared
```

Validation:

```bash
sudo systemctl is-enabled mosaic cloudflared
sudo systemctl status mosaic cloudflared --no-pager
curl -sS http://127.0.0.1:8080/api/health
```

After a reboot:

```bash
sudo reboot
# after reconnecting via SSH:
sudo systemctl status mosaic cloudflared --no-pager
curl -I https://mosaic-generator.wilbertopachecob.dev/api/health
```

If you use PM2 for other apps on this same Pi, also persist PM2:

```bash
pm2 save
pm2 startup
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

## Cloudflare hardening for tunnel

If using a public domain:

1. Add the app hostname as a **Public Hostname** in your existing tunnel
2. Remove conflicting `A/AAAA/CNAME` records before creating tunnel DNS routes
3. Keep the hostname proxied through Cloudflare
4. Add a **WAF rate limiting rule** on `POST /api/file/upload`
5. Enable **Bot Fight Mode** or **Super Bot Fight Mode**

This matches your troubleshooting pattern: DNS conflicts are a common failure point, while PM2/system services may still be healthy.

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
| 1 | Cloudflare Tunnel hostname + clean DNS conflicts | Low |
| 2 | Firewall — allow SSH only, keep 8080 private | Low |
| 3 | `systemd` persistence for `mosaic` + `cloudflared` | Low |
| 4 | systemd hardening (`mosaic` non-root + memory limit) | Low |
| 5 | Tight Pi `.env` (see above) | Low |
| 6 | Cloudflare WAF/Bot protections | Medium |
| 7 | CAPTCHA on upload | Medium |
| 8 | Monitoring / alerts | Medium |

---

## Threat model summary

| Threat | Mitigation |
| --- | --- |
| CPU exhaustion (many uploads) | App rate limit + concurrency cap + Cloudflare WAF rate limits |
| Memory exhaustion (huge images) | Dimension + file size caps |
| Expensive tileSize abuse | Server-side 5–100 bounds |
| Stored user photos | Not stored (in-memory only) |
| MITM | HTTPS at Cloudflare edge; optional strict origin TLS in direct-origin mode |
| Direct port scanning | Firewall; app on localhost only |
| Bot floods | Rate limit → CAPTCHA → Cloudflare |
| Pi crash from one bad request | systemd `MemoryMax`, `MAX_CONCURRENT_GENERATIONS=1` |

---

## References

- App config: `env.example`
- Middleware: `middleware/security.go`
- Upload validation: `handlers.go`, `validation.go`
