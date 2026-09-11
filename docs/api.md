# SawitOS Daemon REST API Reference (`v1`)

The `sawitd` daemon exposes a versioned HTTP REST API (`/api/v1/`) on port 8080 (default).

## Base URL

```text
http://127.0.0.1:8080/api/v1
```

---

## Endpoints

### 1. System & Health
- `GET /api/v1/ping`: Health ping.
- `GET /api/v1/system`: System status & resources summary.
- `GET /api/v1/info`: Detailed OS, hardware, and uptime metrics.
- `GET /api/v1/health`: Health evaluation report.

### 2. Service Management
- `GET /api/v1/services`: List systemd services.
- `POST /api/v1/services/start`: Start service (`{"service": "nginx"}`).
- `POST /api/v1/services/stop`: Stop service (`{"service": "nginx"}`).
- `POST /api/v1/services/restart`: Restart service (`{"service": "nginx"}`).

### 3. Package Management
- `GET /api/v1/packages/search?q=<query>`: Search package repository.
- `POST /api/v1/packages/install`: Install package (`{"package": "nginx"}`).
- `POST /api/v1/packages/remove`: Remove package (`{"package": "nginx"}`).
- `POST /api/v1/packages/update`: Refresh package index.

### 4. Firewall Management
- `GET /api/v1/firewall`: List `nftables` firewall rules.
- `POST /api/v1/firewall/allow`: Allow port/protocol (`{"rule": "80/tcp"}`).
- `POST /api/v1/firewall/deny`: Deny port/protocol (`{"rule": "23/tcp"}`).

### 5. Container Management (Docker / OCI)
- `GET /api/v1/containers`: List containers.
- `POST /api/v1/containers/start`: Start container (`{"container": "redis"}`).
- `POST /api/v1/containers/stop`: Stop container (`{"container": "redis"}`).
- `POST /api/v1/containers/restart`: Restart container (`{"container": "redis"}`).
- `POST /api/v1/containers/remove`: Remove container (`{"container": "redis"}`).
- `GET /api/v1/containers/logs?id=<id>`: Fetch container stdout/stderr logs.
- `GET /api/v1/containers/images`: List cached container images.
- `GET /api/v1/containers/volumes`: List container volumes.
- `GET /api/v1/containers/networks`: List container networks.

### 6. Logs Viewer
- `GET /api/v1/logs?unit=<unit>&lines=<n>&priority=<level>`: Query `journalctl` logs.

### 7. Network & Storage
- `GET /api/v1/network`: Query network interfaces, routing table, and DNS config.
- `GET /api/v1/storage`: Query filesystem mount usage and physical disk layout.

### 8. Security & Audit
- `GET /api/v1/security/status`: High-level security configuration summary.
- `GET /api/v1/security/audit`: Security posture audit evaluation report and score.
- `GET /api/v1/security/updates`: Unapplied Debian security repository patches count and packages.

### 9. Server Profiles
- `GET /api/v1/profiles`: List registered server profile presets.
- `POST /api/v1/profiles/apply`: Apply profile preset (`{"profile": "container"}`).


