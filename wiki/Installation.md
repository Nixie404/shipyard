# Installation

Shipyard requires `git`, `docker`, and `docker-compose`.

## Arch Linux (AUR / `makepkg`)
```bash
git clone https://github.com/Nixie404/yardctl.git
cd yardctl
makepkg -sic
sudo systemctl enable --now yardctl-sync.timer
```

## Debian / Ubuntu
```bash
./dist/build-deb.sh
sudo apt install ./dist/yardctl_0.1.0_amd64.deb
sudo systemctl enable --now yardctl-sync.timer
```

## From Source
```bash
go build -o yardctl
sudo ./dist/install.sh
```

---
**Navigation**
* [🏠 Home](https://github.com/Nixie404/yardctl/wiki)
* [📚 CLI Reference](https://github.com/Nixie404/yardctl/wiki/CLI-Reference)
