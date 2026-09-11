# SawitOS Installation Guide

## Development Build & Setup

### Building Binaries
To build the SawitOS management binaries on your local system:

```bash
git clone https://github.com/sawitos/sawit.git
cd sawit
make build
```

The compiled binaries will be located in `bin/`:
- `bin/sawitctl`
- `bin/sawitd`
- `bin/sawit-agent`
- `bin/sawit-health`
- `bin/sawit-update`

---

## Linux System Installation

To install SawitOS binaries, configurations, and systemd units onto a Debian Linux host:

```bash
sudo make install
```

This installs binaries to `/usr/bin`, `/usr/sbin`, `/usr/libexec`, and copies system config templates to `/etc/sawit/`.

### Enable Services

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now sawit-agent.service
sudo systemctl enable --now sawitd.service
```

Verify service status:

```bash
sawitctl status
sawitctl health
```
