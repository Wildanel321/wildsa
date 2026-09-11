# Accessing SawitOS Remotely with Cloudflare Tunnel & Custom Domain

**Cloudflare Tunnel** (`cloudflared`) allows you to securely expose the SawitOS Web UI and REST API (`http://127.0.0.1:8080`) to the public internet using your own custom domain (e.g., `https://sawit.yourdomain.com`) without opening router ports or needing a static public IP (works through CGNAT).

---

## Architecture Overview

```text
SawitOS Server (Raspberry Pi / VPS)
┌──────────────────────────────────────────────┐
│  sawitd REST & WebSocket  (127.0.0.1:8080)   │
│                      ▲                       │
│                      │ Local Loopback        │
│                      ▼                       │
│  cloudflared daemon (Outbound Tunnel)        │
└──────────────────────┬───────────────────────┘
                       │ Encrypted TLS Tunnel
                       ▼
            Cloudflare Edge Network
                       │
                       ▼
       Client Browser / sawitctl remote
         (https://sawit.yourdomain.com)
```

---

## Step-by-Step Setup Guide

### Step 1: Install `cloudflared`

#### On Debian / Raspberry Pi OS (x86_64 / arm64 / armhf):
```bash
# Add Cloudflare GPG key and APT repository
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | sudo tee /etc/apt/keyrings/cloudflare-main.gpg >/dev/null
echo "deb [signed-by=/etc/apt/keyrings/cloudflare-main.gpg] https://pkg.cloudflare.com/cloudflared bookworm main" | sudo tee /etc/apt/sources.list.d/cloudflared.list

# Update package index and install cloudflared
sudo apt update && sudo apt install -y cloudflared
```

---

### Step 2: Authenticate Cloudflare Account

Run the login command to pair `cloudflared` with your Cloudflare domain:
```bash
cloudflared tunnel login
```
*Click the browser link printed on screen and select your domain.*

---

### Step 3: Create a Dedicated Tunnel

Create a new tunnel named `sawit-tunnel`:
```bash
cloudflared tunnel create sawit-tunnel
```
*Note down the Tunnel ID outputted (e.g. `a1b2c3d4-e5f6-7890-abcd-1234567890ab`).*

---

### Step 4: Route Custom Subdomain to Tunnel

Map your desired subdomain (e.g. `sawit.yourdomain.com`) to the tunnel:
```bash
cloudflared tunnel route dns sawit-tunnel sawit.yourdomain.com
```

---

### Step 5: Configure `cloudflared` Configuration File

Create `/etc/cloudflared/config.yml`:
```yaml
tunnel: a1b2c3d4-e5f6-7890-abcd-1234567890ab
credentials-file: /root/.cloudflared/a1b2c3d4-e5f6-7890-abcd-1234567890ab.json

ingress:
  - hostname: sawit.yourdomain.com
    service: http://127.0.0.1:8080
    originRequest:
      connectTimeout: 10s
      noTLSVerify: true
  - service: http_status:404
```

---

### Step 6: Install and Enable `cloudflared` Systemd Service

```bash
# Install systemd service
sudo cloudflared service install

# Start and enable service
sudo systemctl enable --now cloudflared
```

---

## Verification

1. Open `https://sawit.yourdomain.com` in your web browser.
2. Log in with your SawitOS credentials (`admin` / `admin123`).
3. Notice that:
   - **HTTPS SSL Certificate** is active and trusted automatically.
   - **Live WebSocket Telemetry** operates in real time.
   - No port forwarding or public IP configuration was required!
