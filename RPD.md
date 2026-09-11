# SawitOS — Requirements & Product Design

> **Document Type:** Requirements & Product Design
> **Project:** SawitOS
> **Current Version:** 0.1.0-dev
> **Status:** Active Development
> **Primary Goal:** Lightweight, headless, web-managed server operating system

---

# 1. Project Overview

## 1.1 Product Name

**SawitOS**

## 1.2 Product Type

SawitOS is a lightweight, headless Linux server operating system built initially on top of Debian Stable.

SawitOS is NOT intended to replace the Linux kernel, Debian package ecosystem, systemd, or other mature upstream components.

Instead, SawitOS provides a unified management layer consisting of:

* `sawitctl`
* `sawitd`
* `sawit-agent`
* Sawit Web
* SawitOS installer
* SawitOS configuration
* SawitOS security defaults
* SawitOS server profiles

---

# 2. Product Vision

SawitOS should provide a server administration experience where common server operations can be performed through either:

```text
Sawit Web
```

or:

```bash
sawitctl
```

The system must remain usable through SSH even when the Web UI is unavailable.

### Core principle

> Linux handles the operating system. SawitOS provides the unified management experience.

---

# 3. Primary Goals

SawitOS MUST prioritize:

1. Lightweight resource usage
2. Headless operation
3. Web-based administration
4. CLI administration
5. Security
6. Reliability
7. Modularity
8. Maintainability
9. Debian compatibility
10. Production-oriented architecture

---

# 4. Non-Goals

The following are NOT goals for the initial versions:

* Creating a new Linux kernel
* Rewriting systemd
* Rewriting APT/dpkg
* Creating a new filesystem
* Creating a complete desktop environment
* Replacing every Linux command
* Creating a proprietary package format
* Building every server application into the OS
* Creating a general-purpose desktop OS

Do NOT implement these unless explicitly approved in a future RPD revision.

---

# 5. Target Users

Primary users:

* Homelab administrators
* VPS users
* Small server administrators
* Developers
* Students learning Linux/server administration
* Self-hosting users
* Small organizations

---

# 6. Target Hardware

Initial target:

### Minimum

```text
CPU:     x86_64 / ARM64
RAM:     512 MB+
Storage: 4 GB+
Network: Ethernet or supported network interface
```

### Recommended

```text
CPU:     2+ cores
RAM:     2 GB+
Storage: 16 GB+
Network: Ethernet
```

SawitOS must not assume high-end hardware.

---

# 7. Operating Model

SawitOS is:

```text
HEADLESS FIRST
```

A desktop GUI is NOT installed by default.

After installation, the administrator primarily interacts using:

```text
SSH
+
Sawit Web
+
sawitctl
```

---

# 8. Architecture

```text
                         Browser
                            │
                           HTTPS
                            │
                            ▼
                     ┌─────────────┐
                     │  Sawit Web  │
                     └──────┬──────┘
                            │
                         REST/WS
                            │
                            ▼
                     ┌─────────────┐
                     │   sawitd    │
                     │             │
                     │ API/Auth    │
                     │ Events      │
                     │ Jobs        │
                     └──────┬──────┘
                            │
                       Unix Socket
                            │
                            ▼
                    ┌──────────────┐
                    │ sawit-agent  │
                    │              │
                    │ Privileged   │
                    │ Operations   │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              ▼            ▼            ▼
          systemd       nftables      storage
              │
              ▼
          Linux/Debian
```

---

# 9. Component Responsibilities

## 9.1 sawitctl

Purpose:

CLI interface for SawitOS.

Location:

```text
/usr/bin/sawitctl
```

Responsibilities:

* Query system state
* Manage services
* Manage packages
* Query network
* Query storage
* Manage firewall
* Manage containers
* Run health checks
* Manage updates

` sawitctl` should NOT directly duplicate privileged logic when that logic belongs in `sawit-agent`.

---

# 10. sawitd

Purpose:

Core SawitOS daemon and API server.

Responsibilities:

* REST API
* WebSocket
* Authentication
* Authorization
* Event handling
* Job management
* System state
* Communication with `sawit-agent`

Location:

```text
/usr/sbin/sawitd
```

Systemd unit:

```text
sawitd.service
```

---

# 11. sawit-agent

Purpose:

Privileged system operation layer.

Responsibilities:

* Service operations
* Network operations
* Firewall operations
* Storage operations
* Package operations
* Container operations

The agent MUST NOT accept arbitrary shell commands from HTTP requests.

Bad:

```json
{
  "command": "rm -rf /"
}
```

Good:

```json
{
  "operation": "service.restart",
  "service": "nginx"
}
```

All operations must be explicitly defined and validated.

---

# 12. Sawit Web

Sawit Web is the primary graphical management interface.

Required sections:

```text
Dashboard
Services
Applications
Containers
Network
Storage
Security
Users
Updates
Logs
Terminal
Settings
```

The UI must be responsive.

Supported:

* Desktop
* Laptop
* Tablet
* Mobile

---

# 13. Dashboard Requirements

Dashboard MUST display:

```text
CPU
RAM
Disk
Network
System Load
Uptime
Temperature when available
Running Services
Containers
Security Status
Updates
```

Example:

```text
System
────────────────────────

CPU       34%
Memory    41%
Disk      52%
Load      0.87
Uptime    5d 12h

Health
────────────────────────

System       ✓
Network      ✓
Storage      ✓
Security     ✓
Services     ✓
```

---

# 14. Service Management

Required CLI:

```bash
sawitctl service list
sawitctl service status <name>
sawitctl service start <name>
sawitctl service stop <name>
sawitctl service restart <name>
sawitctl service enable <name>
sawitctl service disable <name>
sawitctl service logs <name>
```

Web UI MUST expose equivalent operations.

Service operations should use systemd.

SawitOS must NOT create a second service manager.

---

# 15. Package Management

SawitOS initially uses Debian APT/dpkg.

CLI:

```bash
sawitctl package search <query>
sawitctl package install <package>
sawitctl package remove <package>
sawitctl package update
sawitctl package upgrade
```

The package system must remain compatible with Debian.

Do not create a custom package manager during initial development.

---

# 16. Network Management

Required information:

```text
Interfaces
IPv4
IPv6
Routes
DNS
Hostname
Link status
Traffic statistics
```

CLI:

```bash
sawitctl network status
sawitctl network interfaces
sawitctl network routes
sawitctl network dns
```

Network configuration MUST prioritize remote-access safety.

A remote administrator must receive warnings before operations that may disconnect the current session.

---

# 17. Storage Management

Required:

```text
Disk information
Partitions
Filesystem
Mount points
Usage
```

CLI:

```bash
sawitctl storage list
sawitctl storage disks
sawitctl storage usage
sawitctl storage mounts
```

Destructive operations are NOT part of the initial implementation unless explicitly designed and confirmed.

Never automatically:

```text
format disk
delete partition
destroy filesystem
```

---

# 18. Firewall

Initial firewall backend:

```text
nftables
```

CLI:

```bash
sawitctl firewall status
sawitctl firewall list
sawitctl firewall allow <port>/<protocol>
sawitctl firewall deny <port>/<protocol>
```

Example:

```bash
sawitctl firewall allow 80/tcp
sawitctl firewall allow 443/tcp
```

Firewall changes must be validated before applying.

---

# 19. Security Requirements

SawitOS MUST follow:

```text
Least privilege
Secure defaults
Input validation
Authentication
Authorization
Auditability
Minimal attack surface
```

Default security requirements:

```text
Root SSH login disabled
Firewall enabled where practical
Secure password storage
Security updates available
Minimal packages
No plaintext secrets
```

---

# 20. Authentication

Sawit Web must have authentication.

Required:

* Login
* Logout
* Session management
* Password hashing
* Secure cookies
* Rate limiting
* Authorization

Passwords MUST NOT be stored as plaintext.

---

# 21. Authorization

SawitOS should eventually support roles:

```text
Administrator
Operator
Viewer
```

Initial version may implement only Administrator if required for simplicity.

Do not implement a fake permission system that is not enforced by the backend.

---

# 22. Audit Logging

Administrative operations should be recorded.

Example:

```text
Timestamp:
User:
Action:
Target:
Result:
Source IP:
```

Do not record:

```text
Passwords
API secrets
Tokens
Private keys
```

---

# 23. Web Terminal

Sawit Web may provide an interactive terminal.

Architecture:

```text
Browser
   │
WebSocket
   │
sawitd
   │
PTY
   │
Shell
```

Requirements:

* Interactive shell
* Terminal resizing
* ANSI colors
* Ctrl+C
* Scrollback
* Disconnect handling

Terminal access must be authenticated.

---

# 24. Container Support

Container support is modular.

Initial backend:

```text
Docker / OCI-compatible runtime
```

Required:

```text
Containers
Images
Volumes
Networks
Logs
Stats
```

Do NOT force Docker installation on the Minimal profile.

---

# 25. Health System

Binary:

```text
sawit-health
```

CLI:

```bash
sawitctl health
```

Checks:

```text
CPU
RAM
Disk
Filesystem
Network
Services
Failed systemd units
Container health
Security
Updates
```

Output:

```text
HEALTHY
WARNING
CRITICAL
```

---

# 26. Update System

Binary:

```text
sawit-update
```

Responsibilities:

* Check updates
* Security updates
* Package updates
* Reboot-required detection
* Update history

SawitOS MUST NOT silently perform dangerous system upgrades.

---

# 27. CLI Design

All SawitOS-specific CLI operations use:

```text
sawitctl
```

The command hierarchy must remain consistent.

Example:

```text
sawitctl
├── system
├── service
├── package
├── network
├── storage
├── firewall
├── container
├── security
├── update
├── health
└── user
```

Global options:

```text
--help
--version
--json
--quiet
```

---

# 28. JSON Output

Commands that expose structured information SHOULD support:

```bash
sawitctl status --json
```

Example:

```json
{
  "hostname": "sawit-server",
  "uptime": 123456,
  "cpu_usage": 23.4,
  "memory_usage": 41.2
}
```

This allows automation and future integrations.

---

# 29. API

API version:

```text
/api/v1/
```

Examples:

```text
GET  /api/v1/system
GET  /api/v1/health
GET  /api/v1/services
POST /api/v1/services/{name}/start
POST /api/v1/services/{name}/stop
POST /api/v1/services/{name}/restart

GET /api/v1/network
GET /api/v1/storage
GET /api/v1/containers
GET /api/v1/logs
GET /api/v1/updates
```

All privileged endpoints require authentication and authorization.

---

# 30. WebSocket

WebSocket should be used for:

```text
Live metrics
Terminal
Live logs
Events
Long-running tasks
```

Do not use WebSocket for everything unnecessarily.

---

# 31. Configuration

SawitOS configuration directory:

```text
/etc/sawit/
```

Persistent state:

```text
/var/lib/sawit/
```

Runtime state:

```text
/run/sawit/
```

Logs:

```text
/var/log/sawit/
```

Do not duplicate standard Linux configuration unnecessarily.

---

# 32. systemd Integration

SawitOS uses systemd.

Example units:

```text
sawitd.service
sawit-agent.service
sawit-health.service
```

Services must:

* Have correct dependencies
* Restart safely
* Use appropriate permissions
* Avoid unnecessary startup work
* Have resource controls where appropriate

---

# 33. Performance Requirements

Performance goals:

```text
Low idle CPU
Low idle RAM
Fast boot
Small disk footprint
Low Web API overhead
Low monitoring overhead
```

Target:

```text
Minimal installation:
<512 MB idle RAM where realistically achievable
```

This is a target, NOT a fake benchmark requirement.

Measurements must be performed on actual test systems.

Do not remove important security or reliability features merely to achieve an arbitrary RAM number.

---

# 34. Technology Stack

## Core

```text
Go
```

Primary components:

```text
sawitctl
sawitd
sawit-agent
sawit-health
sawit-update
```

## Web

```text
Next.js
React
TypeScript
Tailwind CSS
```

## Communication

```text
REST
WebSocket
JSON
Unix sockets
```

---

# 35. Repository Structure

Expected structure:

```text
sawitos/
├── cmd/
│   ├── sawitctl/
│   ├── sawitd/
│   ├── sawit-agent/
│   ├── sawit-health/
│   └── sawit-update/
│
├── internal/
│   ├── system/
│   ├── services/
│   ├── packages/
│   ├── network/
│   ├── storage/
│   ├── firewall/
│   ├── containers/
│   ├── security/
│   └── users/
│
├── web/
│   ├── app/
│   ├── components/
│   ├── lib/
│   └── styles/
│
├── installer/
├── packaging/
├── systemd/
├── configs/
├── scripts/
├── tests/
├── docs/
├── build/
│
├── Makefile
├── RPD.md
├── README.md
├── LICENSE
└── VERSION
```

The actual structure may evolve, but architectural changes must be documented.

---

# 36. Server Profiles

Installer profiles:

```text
Minimal
Web Server
Container Host
Database Server
Game Server
Homelab
Custom
```

Profiles must be configuration/package presets.

They must NOT create completely separate operating systems.

---

# 37. First Boot

After installation, SawitOS should display:

```text
SawitOS Server

Hostname:
sawit-server

Web Management:
https://SERVER-IP/

SSH:
ssh <user>@SERVER-IP
```

The Web UI should be reachable without requiring a desktop environment.

---

# 38. Installation Flow

Initial installation:

```text
Boot ISO
   ↓
Language
   ↓
Keyboard
   ↓
Network
   ↓
Hostname
   ↓
Disk
   ↓
User
   ↓
Server Profile
   ↓
Installation
   ↓
Reboot
   ↓
First Boot
```

---

# 39. Build System

SawitOS must eventually produce:

```text
SawitOS ISO
```

Future targets:

```text
x86_64
ARM64
Cloud image
VM image
```

Builds should be reproducible as far as practical.

---

# 40. Testing

Every important feature requires tests.

### Unit tests

```text
CLI parsing
API
Authentication
Authorization
System information
Validation
```

### Integration tests

```text
Service management
Package management
Network queries
Storage queries
Firewall
Containers
```

### Security tests

```text
Command injection
Path traversal
Authentication bypass
Authorization bypass
CSRF
Session attacks
Privilege escalation
```

### System tests

```text
ISO boot
Installation
First boot
Web UI availability
sawitctl availability
systemd services
```

---

# 41. Development Rules

These rules are mandatory for AI coding agents.

## Rule 1 — Do not hallucinate

If functionality is not implemented, say:

```text
NOT IMPLEMENTED
```

Do not pretend that it works.

---

## Rule 2 — Do not create fake APIs

Every documented API endpoint must correspond to real backend functionality.

---

## Rule 3 — Do not create fake system statistics

CPU, RAM, disk, network, temperature, and service information must come from real system data.

---

## Rule 4 — Do not fake security

Never write:

```text
Security: SECURE
```

unless actual checks have been performed.

---

## Rule 5 — No arbitrary privileged execution

Never implement a generic:

```text
execute(command)
```

endpoint for the Web UI.

---

## Rule 6 — Verify before claiming completion

After implementing a feature:

```text
Build
↓
Run
↓
Test
↓
Verify
↓
Document
```

---

## Rule 7 — Preserve upstream compatibility

Do not modify Debian/Linux components without a clear requirement.

---

## Rule 8 — Minimal dependencies

Every dependency must have a reason.

Avoid adding libraries simply because they are convenient.

---

## Rule 9 — Security before convenience

Do not weaken:

```text
authentication
authorization
firewall
filesystem permissions
TLS
input validation
```

just to make development easier.

---

## Rule 10 — Keep the system recoverable

An administrator must always have a way to recover through:

```text
SSH
console
single-user/recovery environment
```

Do not make Web UI the only recovery mechanism.

---

# 42. Source of Truth

This document is the primary product specification.

When implementation decisions conflict with this document:

```text
RPD.md
    ↓
Architecture documentation
    ↓
Existing implementation
    ↓
Developer assumptions
```

Never silently invent requirements.

If a requested feature conflicts with RPD.md, explain the conflict before implementing it.

---

# 43. Feature Status

Use:

```text
[ ] Not started
[~] In progress
[x] Implemented
[!] Broken
[-] Deprecated
```

Current status:

```text
## Core

[x] sawitctl
[x] sawitd
[x] sawit-agent
[x] sawit-health
[x] sawit-update

## Web

[x] Authentication
[x] Dashboard
[x] Services
[x] Network
[x] Storage
[x] Security
[x] Logs
[ ] Terminal
[x] Settings

## Infrastructure

[x] systemd integration
[ ] Debian packaging
[ ] Installer
[ ] ISO build
[x] Automated tests

## Containers

[x] Container discovery
[x] Container lifecycle
[x] Images
[x] Volumes
[x] Networks
[x] Logs
```

Update this section as development progresses.

---

# 44. Versioning

SawitOS follows:

```text
MAJOR.MINOR.PATCH
```

Example:

```text
0.1.0
0.2.0
0.3.0
1.0.0
```

Before `1.0.0`:

```text
Development / Experimental
```

After `1.0.0`:

```text
Stable
```

---

# 45. Definition of Done

A feature is considered DONE only when:

```text
[ ] Requirements defined
[ ] Implementation exists
[ ] Code builds
[ ] Tests exist
[ ] Tests pass
[ ] Security reviewed
[ ] Error handling exists
[ ] CLI/API behavior documented
[ ] Web UI behavior documented where applicable
[ ] No fake functionality
```

---

# 46. Current Development Directive

The AI agent MUST NOT attempt to build all SawitOS features simultaneously.

Current priority:

```text
PHASE 1
```

Implement only:

```text
sawitctl
sawitd
sawit-agent
system information
health information
basic systemd integration
configuration structure
tests
documentation
```

Initial commands:

```bash
sawitctl status
sawitctl info
sawitctl health
```

Only after Phase 1 is working should development continue to:

```text
Phase 2 — Core Management
Phase 3 — Web Backend
Phase 4 — Sawit Web
Phase 5 — Containers
Phase 6 — Security
Phase 7 — Installer
Phase 8 — Optimization
```

---

# 47. Final Product Definition

SawitOS is successful when a clean machine can be installed with SawitOS and subsequently administered primarily through:

```text
https://SERVER-IP/
```

while maintaining full administrative capability through:

```bash
sawitctl
```

over SSH.

The final system must remain:

```text
Lightweight
Headless
Secure
Modular
Maintainable
Debian-compatible
Production-oriented
```

And most importantly:

> **Never claim a feature exists when it has not actually been implemented and tested.**
