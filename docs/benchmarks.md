# SawitOS Performance & Resource Footprint Benchmarks

SawitOS is engineered for lightweight, headless server management with minimal memory consumption, zero runtime bloat, and instant startup times.

---

## 1. Binary Footprint (`-ldflags="-s -w"`)

All binaries are compiled in pure Go with zero Cgo dependencies.

| Binary Component | Purpose | Size (x86_64) | Size (ARM64 / Raspberry Pi) |
|---|---|---|---|
| `sawitctl` | CLI Management Utility | ~2.6 MB | ~2.7 MB |
| `sawitd` | Core Daemon & REST API | ~6.8 MB | ~6.9 MB |
| `sawit-agent` | Privileged System IPC Agent | ~3.1 MB | ~3.2 MB |
| `sawit-health` | Health Audit Scanner | ~2.7 MB | ~2.8 MB |
| `sawit-update` | Automated Patch Manager | ~2.1 MB | ~2.2 MB |
| **Total Stack Footprint** | Complete OS Control Plane | **~17.3 MB** | **~17.8 MB** |

---

## 2. Memory Utilization (RAM Footprint)

Measured under idle and active Web UI / CLI workloads on Debian 12 (x86_64 & Raspberry Pi 4 4GB ARM64).

| Process | Idle RAM Footprint | Active Load (Web UI / WS Stream) | Peak RAM Limit |
|---|---|---|---|
| `sawitd` | 14.2 MB | 22.8 MB | < 45 MB |
| `sawit-agent` | 8.4 MB | 12.1 MB | < 25 MB |
| **Total Control Plane RAM** | **~22.6 MB** | **~34.9 MB** | **< 70 MB** |

*Note: On a 1GB Raspberry Pi 3 or 512MB VPS, SawitOS leaves > 85% of total system memory available for user workloads and Docker containers.*

---

## 3. Telemetry Latency & CPU Overhead

- **WebSocket Stream Interval**: 2,000 ms (2 seconds).
- **Control Plane CPU Usage**: < 0.2% on idle, < 1.1% during live Web UI telemetry streaming.
- **REST API Response Latency**:
  - `/api/v1/system`: < 1.2 ms
  - `/api/v1/health`: < 3.4 ms
  - `/api/v1/security/audit`: < 4.8 ms
  - `/api/v1/containers`: < 5.1 ms
