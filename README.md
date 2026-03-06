# Shipyard (`yardctl`)

![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

Shipyard is a lightweight, Docker-centered automated deployment tool. It fetches your configuration from Git repositories (including private repos via SSH) and automatically discovers and manages Docker Compose stacks.

Shipyard's `yardctl` makes managing a fleet of self-hosted, Compose-based applications completely frictionless.

📚 **[Read the full documentation on the GitHub Wiki!](https://github.com/Nixie404/yardctl/wiki)**

## ✨ Features

- **Automated Discovery**: Clones any Git repository and automatically discovers `docker-compose.yml` configs inside it, mapping them into manageable "stacks".
- **GitOps Pulling**: A `yardctl sync` updates your Git repositories to `HEAD` and seamlessly re-evaluates all available stacks.
- **Smart Wrappers**: Thin, context-aware wrappers over Docker Compose commands (`deploy`, `logs`, `ps`, `exec`, `restart`) so you don't have to `cd` everywhere.
- **Auto-Syncing**: Includes a builtin Systemd timer that keeps everything synced automatically every 5 minutes.
- **Rootless/Privilege Escalation**: Safely detects `doas` or `sudo` to only request privileges exactly when needed.

## 📦 Installation

For full installation instructions across different platforms (Arch, Debian/Ubuntu, Source), check out the **[Installation Guide](https://github.com/Nixie404/yardctl/wiki/Installation)** on our wiki.

## 🚀 Quick Start & Typical Workflow

**1. Initialize the system:**
```bash
yardctl init
```
This sets up `/etc/yardctl/` for configuration and `/var/lib/yardctl/` for cloned data.

**2. Verify everything is green:**
```bash
yardctl check
```

**3. Add your deployment repository:**
```bash
yardctl repo add git@github.com:your-org/deployops.git
```

**4. Pull latest and sync:**
```bash
yardctl repo sync
```

**5. List discovered stacks:**
```bash
yardctl stack list
```

**6. Deploy a stack:**
```bash
yardctl deploy my-api
# or yardctl stack deploy my-api
```

**7. Check status and logs:**
```bash
yardctl ps
yardctl logs my-api -f
yardctl status
```

## 📚 CLI Reference

For a complete breakdown of every command and subcommand available in Shipyard, read the **[CLI Reference](https://github.com/Nixie404/yardctl/wiki/CLI-Reference)** on the wiki.

## 📄 License
MIT License.
