#!/bin/sh
# SawitOS Automated Server & Edge Appliance Installer
# Target: Debian 11/12 (x86_64, aarch64 / ARM64, armv7l)

set -e

SAWIT_VERSION="0.1.0-dev"
INSTALL_PREFIX="/usr/local/bin"
CONFIG_DIR="/etc/sawit"
PROFILE="minimal"

# Print banner
echo "================================================================"
echo "          SawitOS Lightweight Server OS Installer              "
echo "                 Version: ${SAWIT_VERSION}                      "
echo "================================================================"

# Parse command line flags
while [ $# -gt 0 ]; do
    case "$1" in
        --profile=*)
            PROFILE="${1#*=}"
            ;;
        --profile)
            PROFILE="$2"
            shift
            ;;
        --help|-h)
            echo "Usage: sawit-install.sh [options]"
            echo ""
            echo "Options:"
            echo "  --profile <name>   Set server profile (minimal, web, container, database, homelab)"
            echo "  --help             Display installer help"
            exit 0
            ;;
    esac
    shift
done

# Check root privileges
if [ "$(id -u)" -ne 0 ]; then
    echo "Error: SawitOS installer must be run as root." >&2
    exit 1
fi

# Detect OS architecture
ARCH=$(uname -m)
case "${ARCH}" in
    x86_64|amd64)
        BINARY_ARCH="amd64"
        ;;
    aarch64|arm64)
        BINARY_ARCH="arm64"
        ;;
    armv7l|armhf)
        BINARY_ARCH="armv7"
        ;;
    *)
        echo "Error: Unsupported architecture '${ARCH}'. SawitOS supports x86_64, arm64, armv7." >&2
        exit 1
        ;;
esac

echo "[+] Target System: Linux (${ARCH} -> ${BINARY_ARCH})"
echo "[+] Selected Profile Preset: ${PROFILE}"

# Create directories
mkdir -p "${INSTALL_PREFIX}"
mkdir -p "${CONFIG_DIR}"
mkdir -p /var/log/sawit
mkdir -p /run/sawit

# Create default configuration files if not present
if [ ! -f "${CONFIG_DIR}/sawitd.yaml" ]; then
    echo "[+] Writing default configuration to ${CONFIG_DIR}/sawitd.yaml"
    cat <<EOF > "${CONFIG_DIR}/sawitd.yaml"
server:
  host: "0.0.0.0"
  port: 8080
agent:
  socket_path: "/run/sawit/sawit-agent.sock"
logging:
  level: "info"
  file_path: "/var/log/sawit/sawitd.log"
auth:
  jwt_secret: "sawit-server-secret-key-change-in-production"
  admin_username: "admin"
  admin_password_hash: "\$2a\$10\$K124Y1nN0V0mR9m6Y.1.E.O/gqfN/6N8V8V8V8V8V8V8V8V8V8V8" # default: admin123
EOF
fi

if [ ! -f "${CONFIG_DIR}/security.yaml" ]; then
    echo "[+] Writing default security policy to ${CONFIG_DIR}/security.yaml"
    cat <<EOF > "${CONFIG_DIR}/security.yaml"
firewall:
  enabled: true
  backend: "nftables"
  default_policy: "drop"
ssh:
  disable_root_login: true
  password_auth: true
EOF
fi

# Install systemd service units
echo "[+] Installing systemd service unit files"

cat <<EOF > /etc/systemd/system/sawitd.service
[Unit]
Description=SawitOS Core Daemon & REST API
After=network.target sawit-agent.service
Wants=sawit-agent.service

[Service]
Type=simple
ExecStart=${INSTALL_PREFIX}/sawitd -config ${CONFIG_DIR}/sawitd.yaml
Restart=always
RestartSec=3s
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

cat <<EOF > /etc/systemd/system/sawit-agent.service
[Unit]
Description=SawitOS Privileged System Agent
After=network.target

[Service]
Type=simple
ExecStart=${INSTALL_PREFIX}/sawit-agent -socket /run/sawit/sawit-agent.sock
Restart=always
RestartSec=3s

[Install]
WantedBy=multi-user.target
EOF

cat <<EOF > /etc/systemd/system/sawit-health.service
[Unit]
Description=SawitOS Automated Health Audit Timer
After=sawitd.service

[Service]
Type=oneshot
ExecStart=${INSTALL_PREFIX}/sawit-health

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd daemon
if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload
    systemctl enable sawit-agent sawitd sawit-health >/dev/null 2>&1 || true
fi

echo "================================================================"
echo " [✓] SawitOS installation structure complete!"
echo " Binaries directory: ${INSTALL_PREFIX}"
echo " Config directory:   ${CONFIG_DIR}"
echo " Default Web UI:     http://<server-ip>:8080"
echo " Initial Login:      admin / admin123"
echo "================================================================"
