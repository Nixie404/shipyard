# CLI Reference

## Root Commands
* `yardctl init` — Setup system directories
* `yardctl check` — Comprehensive environment/permissions health check
* `yardctl doctor` — Quick dependencies check
* `yardctl status` — Dashboard of repositories, stacks, and containers
* `yardctl sync` — Pull all repos and rediscover stacks
* `yardctl deploy <stack>` — Deploy a named stack
* `yardctl logs <stack> [service]` — View logs
* `yardctl ps [stack]` — List running containers
* `yardctl exec <stack> <service> [command]` — Exec into a container
* `yardctl restart <stack> [service]` — Restart a service
* `yardctl events` — Stream Docker events
* `yardctl purge` — Uninstall/Delete all yardctl data safely

## Subcommands
* **`yardctl repo [add|list|remove|sync]`** — Manage Git configuration sources
* **`yardctl stack [list|info|deploy|destroy]`** — Manage Docker compose deployments
* **`yardctl config [view|edit]`** — Manage the main `yardctl` configuration json
* **`yardctl context [list|use]`** — Manage target endpoints (currently `local` only)

---
**Navigation**
* [🏠 Home](https://github.com/Nixie404/yardctl/wiki)
* [📦 Installation Guide](https://github.com/Nixie404/yardctl/wiki/Installation)
