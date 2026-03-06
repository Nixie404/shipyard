# Maintainer: nixie
pkgname=yardctl
pkgver=0.1.0
pkgrel=1
pkgdesc='Docker-centered automated deployment tool that pulls configuration from git repositories and manages Docker Compose stacks'
arch=('x86_64' 'aarch64')
url='https://github.com/Nixie404/yardctl'
license=('MIT')
depends=('docker' 'docker-compose' 'git')
makedepends=('go')
backup=('etc/yardctl/config.json')
source=()

build() {
  cd "${srcdir}/.."
  export CGO_ENABLED=0
  export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
  go build -ldflags "-s -w \
    -X shipyard/cmd.Version=${pkgver} \
    -X shipyard/cmd.Commit=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) \
    -X shipyard/cmd.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o ${pkgname}
}

package() {
  cd "${srcdir}/.."

  # Binary
  install -Dm755 "${pkgname}" "${pkgdir}/usr/bin/${pkgname}"

  # Systemd units
  install -Dm644 "dist/yardctl-sync.service" "${pkgdir}/usr/lib/systemd/system/yardctl-sync.service"
  install -Dm644 "dist/yardctl-sync.timer"   "${pkgdir}/usr/lib/systemd/system/yardctl-sync.timer"

  # Default config directory
  install -dm755 "${pkgdir}/etc/yardctl"
  install -dm755 "${pkgdir}/var/lib/yardctl/repos"

  # Default config file
  cat > "${pkgdir}/etc/yardctl/config.json" <<EOF
{
  "data_dir": "/var/lib/yardctl",
  "repos_dir": "/var/lib/yardctl/repos",
  "repos": [],
  "stacks": []
}
EOF
  chmod 644 "${pkgdir}/etc/yardctl/config.json"

  # Shell completions
  install -dm755 "${pkgdir}/usr/share/bash-completion/completions"
  install -dm755 "${pkgdir}/usr/share/zsh/site-functions"
  install -dm755 "${pkgdir}/usr/share/fish/vendor_completions.d"
  ./${pkgname} completion bash > "${pkgdir}/usr/share/bash-completion/completions/${pkgname}"
  ./${pkgname} completion zsh  > "${pkgdir}/usr/share/zsh/site-functions/_${pkgname}"
  ./${pkgname} completion fish > "${pkgdir}/usr/share/fish/vendor_completions.d/${pkgname}.fish"

  # Man page
  install -Dm644 "man/yardctl.1" "${pkgdir}/usr/share/man/man1/yardctl.1"
}
