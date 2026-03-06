#!/usr/bin/env bash
set -euo pipefail

# Install script for yardctl systemd sync timer
# Run as root or with sudo

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "📦 Installing yardctl sync service and timer..."

# Copy binary
if [ -f "${SCRIPT_DIR}/../yardctl" ]; then
    cp "${SCRIPT_DIR}/../yardctl" /usr/local/bin/yardctl
    chmod +x /usr/local/bin/yardctl
    echo "✅ Binary installed to /usr/local/bin/yardctl"
else
    echo "⚠  Binary not found. Build first with 'go build -o yardctl'"
    echo "   Assuming yardctl is already in PATH."
fi

# Copy systemd units
cp "${SCRIPT_DIR}/yardctl-sync.service" /etc/systemd/system/
cp "${SCRIPT_DIR}/yardctl-sync.timer" /etc/systemd/system/
echo "✅ Systemd units installed"

# Reload and enable
systemctl daemon-reload
systemctl enable yardctl-sync.timer
systemctl start yardctl-sync.timer

echo "✅ Timer enabled and started"
echo ""
echo "Check status with:"
echo "  systemctl status yardctl-sync.timer"
echo "  systemctl list-timers yardctl-sync.timer"
echo "  journalctl -u yardctl-sync.service"
