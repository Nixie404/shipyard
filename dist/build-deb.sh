#!/usr/bin/env bash
set -euo pipefail

# Build a .deb package for yardctl
# Usage: ./dist/build-deb.sh [version]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
VERSION="${1:-0.1.0}"
ARCH="$(dpkg --print-architecture 2>/dev/null || echo amd64)"
PKG_NAME="yardctl"
PKG_DIR="${PROJECT_DIR}/dist/${PKG_NAME}_${VERSION}_${ARCH}"

echo "📦 Building ${PKG_NAME} ${VERSION} (${ARCH}) .deb package..."

# ── Build binary ──
echo "🔨 Compiling..."
cd "$PROJECT_DIR"

BUILD_FLAGS=""
LDFLAGS="-X shipyard/cmd.Version=${VERSION} \
    -X shipyard/cmd.Commit=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) \
    -X shipyard/cmd.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)"

if [ "${DEBUG:-0}" = "1" ]; then
    echo "🪲 Building with debug symbols (optimizations disabled)..."
    BUILD_FLAGS="-gcflags=all=-N -l"
else
    LDFLAGS="-s -w $LDFLAGS"
fi

CGO_ENABLED=0 go build \
  -trimpath \
  ${BUILD_FLAGS} \
  -ldflags "${LDFLAGS}" \
  -o "${PKG_NAME}"

# ── Assemble package tree ──
echo "📁 Assembling package structure..."
rm -rf "$PKG_DIR"

# Binary
install -Dm755 "${PKG_NAME}" "${PKG_DIR}/usr/bin/${PKG_NAME}"

# Systemd units
install -Dm644 "dist/yardctl-sync.service" "${PKG_DIR}/lib/systemd/system/yardctl-sync.service"
install -Dm644 "dist/yardctl-sync.timer"   "${PKG_DIR}/lib/systemd/system/yardctl-sync.timer"

# Default config
install -dm755 "${PKG_DIR}/etc/yardctl"
cat > "${PKG_DIR}/etc/yardctl/config.json" <<EOF
{
  "data_dir": "/var/lib/yardctl",
  "repos_dir": "/var/lib/yardctl/repos",
  "repos": [],
  "stacks": []
}
EOF

# Data directories
install -dm755 "${PKG_DIR}/var/lib/yardctl/repos"

# Shell completions
install -dm755 "${PKG_DIR}/usr/share/bash-completion/completions"
install -dm755 "${PKG_DIR}/usr/share/zsh/vendor-completions"
install -dm755 "${PKG_DIR}/usr/share/fish/vendor_completions.d"
./${PKG_NAME} completion bash > "${PKG_DIR}/usr/share/bash-completion/completions/${PKG_NAME}"
./${PKG_NAME} completion zsh  > "${PKG_DIR}/usr/share/zsh/vendor-completions/_${PKG_NAME}"
./${PKG_NAME} completion fish > "${PKG_DIR}/usr/share/fish/vendor_completions.d/${PKG_NAME}.fish"

# Man page
install -Dm644 "man/yardctl.1" "${PKG_DIR}/usr/share/man/man1/yardctl.1"
gzip -9 "${PKG_DIR}/usr/share/man/man1/yardctl.1"

# ── DEBIAN control files ──
install -dm755 "${PKG_DIR}/DEBIAN"

# Installed size (in KB)
INSTALLED_SIZE=$(du -sk "${PKG_DIR}" | cut -f1)

cat > "${PKG_DIR}/DEBIAN/control" <<EOF
Package: ${PKG_NAME}
Version: ${VERSION}
Architecture: ${ARCH}
Maintainer: nixie
Description: Docker-centered automated deployment from git repositories
 yardctl (Shipyard) pulls configuration from git repositories,
 discovers Docker Compose stacks, and manages their lifecycle.
 Includes a systemd timer for automatic periodic sync.
Depends: git, docker.io | docker-ce, docker-compose-plugin | docker-compose
Recommends: openssh-client
Section: admin
Priority: optional
Homepage: https://github.com/Nixie404/yardctl
Installed-Size: ${INSTALLED_SIZE}
EOF

# Conffiles — mark config as conffile so dpkg won't overwrite user edits
cat > "${PKG_DIR}/DEBIAN/conffiles" <<EOF
/etc/yardctl/config.json
EOF

# Post-install: reload systemd
cat > "${PKG_DIR}/DEBIAN/postinst" <<'EOF'
#!/bin/sh
set -e
if [ -d /run/systemd/system ]; then
    systemctl daemon-reload
fi
echo ""
echo "yardctl installed! Get started with:"
echo "  yardctl init"
echo "  yardctl check"
echo ""
echo "To enable automatic sync:"
echo "  systemctl enable --now yardctl-sync.timer"
EOF
chmod 755 "${PKG_DIR}/DEBIAN/postinst"

# Pre-remove: stop timer
cat > "${PKG_DIR}/DEBIAN/prerm" <<'EOF'
#!/bin/sh
set -e
if [ -d /run/systemd/system ]; then
    systemctl stop yardctl-sync.timer 2>/dev/null || true
    systemctl disable yardctl-sync.timer 2>/dev/null || true
    systemctl stop yardctl-sync.service 2>/dev/null || true
fi
EOF
chmod 755 "${PKG_DIR}/DEBIAN/prerm"

# Post-remove: reload systemd
cat > "${PKG_DIR}/DEBIAN/postrm" <<'EOF'
#!/bin/sh
set -e
if [ -d /run/systemd/system ]; then
    systemctl daemon-reload
fi
EOF
chmod 755 "${PKG_DIR}/DEBIAN/postrm"

# ── Build .deb ──
echo "📦 Building .deb package using dpkg-deb..."
DEB_FILE="${PROJECT_DIR}/dist/${PKG_NAME}_${VERSION}_${ARCH}.deb"

# Ensure the parent directory for the .deb file exists
mkdir -p "$(dirname "$DEB_FILE")"

# Build the package
dpkg-deb --build "${PKG_DIR}" "${DEB_FILE}"

# Cleanup staging dir
rm -rf "$PKG_DIR"

echo ""
echo "✅ Package built: ${DEB_FILE}"
echo ""
echo "Install with:"
echo "  sudo dpkg -i ${DEB_FILE}"
echo "  sudo apt-get install -f  # fix any missing deps"
echo ""
echo "Or with apt directly:"
echo "  sudo apt install ./${PKG_NAME}_${VERSION}_${ARCH}.deb"
