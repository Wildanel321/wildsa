# SawitOS — Lightweight Headless Server Operating System

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8.svg)](https://golang.org)
[![Debian Compatible](https://img.shields.io/badge/Debian-11%2F12-A80030.svg)](https://debian.org)
[![Architecture](https://img.shields.io/badge/Architecture-x86__64%20%7C%20ARM64%20%7C%20ARMv7-success.svg)](docs/benchmarks.md)

> **SawitOS** is a lightweight, secure, modular, headless Linux server OS with a unified Web UI and CLI management system built on top of Debian Stable. Designed for VPS, dedicated servers, homelabs, and Raspberry Pi edge devices.

---

## 🌟 Core Features

- **Headless-First & Ultra-Lightweight**: Minimal memory footprint (<60 MB total control plane RAM), leaving >85% system memory for user workloads.
- **Cross-Architecture Support**: Runs natively on `x86_64` (AMD64) and `ARM64` / `ARMv7` (Raspberry Pi 3/4/5, Orange Pi, ARM Cloud VPS).
- **Unified Control System**:
  - `sawitctl` CLI for terminal & SSH administration.
  - `sawitd` Daemon exposing versioned REST API (`/api/v1`) & live WebSockets (`/api/v1/ws`).
  - `sawit-web` Next.js 14 Web UI with live telemetry, container controls, and audit reports.
- **Modular Server Profiles**: Instantly switch deployment profiles (`minimal`, `web`, `container`, `database`, `homelab`).
- **Privileged Agent Architecture**: Secure IPC separation over `/run/sawit/sawit-agent.sock` without arbitrary shell injection risks.
- **Automated Security & Audit**: `nftables` stateful firewall enforcement, SSH hardening audits, and Debian security patch detection.
- **Docker / OCI Runtime Native**: Full container, image, volume, and network management.
- **Cloudflare Tunnel Ready**: Easy remote domain access (`https://sawit.yourdomain.com`) without open ports or CGNAT router issues.

---

## 🏗 System Architecture

```text
                        Sawit Web UI / sawitctl CLI
                                     │
                     HTTP REST / WebSockets / Terminal
                                     │
                                     ▼
                           sawitd (Core Daemon)
                                     │
                        Unix Domain Socket (IPC)
                                     │
                                     ▼
                       sawit-agent (Privileged Agent)
                                     │
        ┌───────────────┬────────────┴───┬───────────────┬───────────────┐
        ▼               ▼                ▼               ▼               ▼
     systemd         nftables           APT            Docker        Journalctl
  (Services)       (Firewall)       (Packages)      (Containers)       (Logs)
```

---

## ⚡ Automated One-Line Installation

To install SawitOS management stack on any Debian 11/12 server or Raspberry Pi:

```bash
curl -fsSL https://raw.githubusercontent.com/sawitos/sawit/main/installer/sawit-install.sh | sudo sh -s -- --profile=container
```

*Default Web UI:* `http://<server-ip>:8080` (Credentials: `admin` / `admin123`)

---

## 🎯 Server Profile Presets

| Profile | Title | Target Use Case | Included Components |
|---|---|---|---|
| `minimal` | Minimal Core Server | VPS / Low-RAM Edge Nodes | `sawitd`, `sawit-agent`, SSH, `nftables` |
| `web` | Web Application Host | Web hosting & reverse proxies | Minimal + Nginx, SSL Certbot, HTTP/HTTPS ports |
| `container` | Docker / OCI Host | Containerized workloads & microservices | Minimal + Docker engine, containerd, Compose |
| `database` | Database Server | Isolated database instances | Minimal + PostgreSQL / MariaDB, port isolation |
| `homelab` | Homelab & Appliance | Raspberry Pi, NAS & Home Servers | Minimal + Docker, NFS/CIFS, Avahi mDNS |

*Apply a profile via CLI:*
```bash
sawitctl profile apply container
```

---

## 💻 CLI Reference (`sawitctl`)

```bash
# System status & hardware info
sawitctl status [--json]
sawitctl info
sawitctl health

# Server Profile management
sawitctl profile list
sawitctl profile show container
sawitctl profile apply web

# Service management
sawitctl service list
sawitctl service restart nginx
sawitctl service logs sawitd

# Package management
sawitctl package search redis
sawitctl package install redis-server

# Firewall management (nftables)
sawitctl firewall status
sawitctl firewall allow 8080/tcp

# Container management (Docker / OCI)
sawitctl container list
sawitctl container start redis-server
sawitctl container logs redis-server

# Security posture audit
sawitctl security status
sawitctl security audit
sawitctl security updates
```

---

## 📊 Resource Footprint & Benchmarks

| Component | RAM Footprint (Idle) | RAM Footprint (Active) | Binary Size (ARM64) |
|---|---|---|---|
| `sawitd` | 14.2 MB | 22.8 MB | ~6.9 MB |
| `sawit-agent` | 8.4 MB | 12.1 MB | ~3.2 MB |
| **Total Control Plane** | **~22.6 MB** | **~34.9 MB** | **~10.1 MB** |

See [docs/benchmarks.md](docs/benchmarks.md) for full benchmark reports and performance metrics.

---

## 🔒 Remote Access via Cloudflare Tunnel

Expose SawitOS safely over HTTPS with your custom domain (`https://sawit.yourdomain.com`) without opening ports on your home router:

Follow the step-by-step guide: [docs/cloudflare-tunnel.md](docs/cloudflare-tunnel.md)

---

## 🛠 Building from Source

### Prerequisites
- Go 1.22+
- Node.js 18+ & npm (for Web UI)
- Linux or Windows/macOS (development)

### Build Go Binaries & Web UI
```bash
# Build Go backends into bin/
make build

# Cross-compile release binaries for Raspberry Pi / Linux ARM64
make build-arm64

# Build Next.js Web UI
make build-web

# Run full unit test suite
make test
```

---

## 📚 Documentation Index

- [`docs/cli.md`](docs/cli.md): `sawitctl` CLI reference manual
- [`docs/api.md`](docs/api.md): `sawitd` REST API documentation
- [`docs/benchmarks.md`](docs/benchmarks.md): Memory and CPU benchmark metrics
- [`docs/cloudflare-tunnel.md`](docs/cloudflare-tunnel.md): Remote access & domain setup guide

---

## 📄 License

SawitOS is licensed under the [MIT License](LICENSE).
