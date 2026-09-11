# SawitOS Architecture & Technical Design

## Architectural Philosophy

SawitOS follows a strict **management layer abstraction** pattern on top of Debian Linux.

```text
Linux Kernel + systemd + standard Linux tools
                     │
                     ▼
             SawitOS Management Layer
     (sawitctl + sawitd + sawit-agent + Sawit Web)
```

SawitOS does not duplicate the Linux kernel, systemd init system, APT package manager, or network stack. Instead, it provides a cohesive administration experience across CLI, REST API, and Web UI.

---

## Security Model & Privilege Separation

### 1. `sawitd` (Unprivileged / Least Privilege Daemon)
- Listens for HTTP REST / WebSocket requests.
- Manages user sessions, authentication, authorization, and audit logs.
- Does NOT run arbitrary shell commands.
- Communicates with `sawit-agent` over a secure Unix socket (`/run/sawit/sawit-agent.sock`).

### 2. `sawit-agent` (Privileged Operations Layer)
- Runs with root privileges.
- Accepts ONLY strictly typed and validated IPC actions (e.g. `service.restart`, `health.check`).
- Rejects raw command execution payloads (`rm -rf`, `sh -c`).
- Mitigates command injection vulnerabilities by design.

---

## Component Layout

| Component | Binary Location | Service / Socket | Function |
| :--- | :--- | :--- | :--- |
| `sawitctl` | `/usr/bin/sawitctl` | CLI client | CLI administration tool |
| `sawitd` | `/usr/sbin/sawitd` | `sawitd.service` | Management API & web daemon |
| `sawit-agent` | `/usr/libexec/sawit-agent` | `sawit-agent.service` | Privileged system action handler |
| `sawit-health` | `/usr/bin/sawit-health` | `sawit-health.service` | Health check runner |
| `sawit-update` | `/usr/bin/sawit-update` | CLI tool | System update inspector |
