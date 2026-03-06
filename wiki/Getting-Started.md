# Getting Started with Shipyard

This guide will walk you through your first automated deployment using `yardctl`.

## 1. Prerequisites

Ensure you have the following installed on your system:
* **Docker** & **Docker Compose**
* **Git**
* **SSH Keys** (if using private repositories)

## 2. Initialization

Run the initialization command to set up the necessary directories and default configuration:

```bash
yardctl init
```

By default, this creates:
- `/etc/yardctl/config.json`: Master configuration.
- `/var/lib/yardctl/`: Data directory for repository clones.

## 3. Pre-Flight Check

Verify your environment is ready:

```bash
yardctl check
```

Address any issues marked with ❌. Warnings (⚠) are often non-critical but good to review.

## 4. Adding a Deployment Repository

Add a Git repository that contains your Docker Compose files. Shipyard will automatically scan all directories for `docker-compose.yml` or `compose.yml` files.

```bash
yardctl repo add git@github.com:your-user/your-deploy-repo.git
```

## 5. Syncing and Discovery

Pull the latest version of your repositories and discover new stacks:

```bash
yardctl sync
```

You can see all discovered stacks using:

```bash
yardctl stack list
```

## 6. Your First Deployment

Deploy a stack by name:

```bash
yardctl deploy <stack-name>
```

Shipyard uses `docker compose up -d` under the hood, but manages the paths and project names for you automatically.

**Note:** If you have a stack named `deploy-api`, you can just run `yardctl deploy api` as long as the name is unique! 🏷️

## 7. Monitoring

Check the status of your deployments:

```bash
# General status overview
yardctl status

# List containers for a stack
yardctl ps <stack-name>

# Tail logs
yardctl logs <stack-name> -f
```

## 8. Automating Updates

To keep your stacks in sync with your Git repository automatically, enable the Shipyard systemd timer:

```bash
sudo systemctl enable --now yardctl-sync.timer
```

Now, every 5 minutes, `yardctl` will pull your repositories and redeploy any stacks that have changed.

---
**Navigation**
* [🏠 Home](https://github.com/Nixie404/yardctl/wiki)
* [📚 CLI Reference](https://github.com/Nixie404/yardctl/wiki/CLI-Reference)
* [📦 Installation Guide](https://github.com/Nixie404/yardctl/wiki/Installation)
