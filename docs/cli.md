# `sawitctl` CLI Reference Manual

`sawitctl` is the command-line interface for administering SawitOS servers over SSH or local terminal sessions.

## Global Flags

- `--json`: Output result in machine-readable JSON format.
- `--version`: Display binary version.
- `--quiet`: Suppress non-essential diagnostic output.
- `--help`: Print command help menu.

---

## Command Reference

### System Status & Info

#### `sawitctl status`
Displays live system summary (hostname, kernel, CPU %, memory %, disk %, load averages, and daemon state).

#### `sawitctl info`
Displays detailed OS, architecture, hardware CPU count, Go version, and uptime.

#### `sawitctl health`
Runs comprehensive system health audit across CPU, RAM, Disk, Network, Services, and Security settings.

---

### Service Management

#### `sawitctl service list`
Lists active systemd service units with load and sub-state.

#### `sawitctl service status <name>`
Shows detailed status of a specific systemd unit.

#### `sawitctl service start <name>` / `stop` / `restart` / `enable` / `disable`
Executes service state transitions using systemd.

#### `sawitctl service logs <name>`
Displays recent systemd log entries for the specified unit.

---

### Package Management

#### `sawitctl package search <query>`
Searches Debian APT package repositories.

#### `sawitctl package install <package>`
Installs a package via APT.

#### `sawitctl package remove <package>`
Removes a package via APT.

#### `sawitctl package update`
Updates the local APT package repository index.

---

### Firewall Management (`nftables`)

#### `sawitctl firewall status` / `firewall list`
Lists active `nftables` firewall rules and filter status.

#### `sawitctl firewall allow <port/protocol>`
Adds an input firewall rule to allow incoming traffic (e.g. `80/tcp`, `443/tcp`).

#### `sawitctl firewall deny <port/protocol>`
Adds an input firewall rule to block incoming traffic (e.g. `23/tcp`).

---

### Container Management (Docker / OCI)

#### `sawitctl container list`
Lists active Docker/OCI containers with ID, names, image, status, and ports.

#### `sawitctl container start <idOrName>` / `stop` / `restart` / `remove`
Controls container execution lifecycle.

#### `sawitctl container logs <idOrName>`
Displays stdout/stderr output lines for a container.

#### `sawitctl container images` / `volumes` / `networks`
Lists cached container images, storage volumes, and virtual bridge networks.

---

### Log Viewer (`journalctl`)

#### `sawitctl logs [--unit <unit>] [--lines <n>] [--priority <level>]`
Queries system and service logs via `journalctl`.

---

### Network & Storage

#### `sawitctl network status` / `interfaces` / `routes` / `dns`
Lists active network interfaces, IP addresses, routing table, and DNS resolvers.

#### `sawitctl storage list` / `usage` / `disks` / `mounts`
Lists filesystem mount points, storage capacity, and physical disk partitions.

---

### Security & Audit

#### `sawitctl security status`
Displays high-level security configuration (firewall state, root SSH login, password auth, pending security updates).

#### `sawitctl security audit`
Runs automated security posture audit check items (SSH hardening, firewall active ruleset, shadow file permissions, security patches) and outputs score percentage.

#### `sawitctl security updates`
Checks for pending Debian security repository updates.

---

### Server Profiles

#### `sawitctl profile list`
Lists all available server profile presets (`minimal`, `web`, `container`, `database`, `homelab`).

#### `sawitctl profile show <name>`
Inspects profile preset details (required APT packages, systemd services, firewall input rules).

#### `sawitctl profile apply <name>`
Applies a server profile preset onto the host system.


