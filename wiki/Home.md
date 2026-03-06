# Welcome to the yardctl Wiki

Shipyard (`yardctl`) is a docker-centered automated deployment tool that pulls configuration from Git repositories and automatically discovers and manages Docker Compose stacks.

## 📖 Documentation

* [🏁 Getting Started](https://github.com/Nixie404/yardctl/wiki/Getting-Started)
* [📦 Installation Guide](https://github.com/Nixie404/yardctl/wiki/Installation)
* [📚 CLI Reference](https://github.com/Nixie404/yardctl/wiki/CLI-Reference)

---

## What makes yardctl different?

Unlike massive enterprise continuous delivery tools like ArgoCD or Portainer, **yardctl** is designed to be purely CLI-driven, extremely lightweight, and directly hooked into `docker-compose`. 

There's no heavy web interface or database required.
It operates purely on Git templates and the local Docker socket, making it the perfect middle-ground for self-hosted instances, personal servers, and lightweight staging environments.

## Systemd Integration

yardctl installs a natively integrated systemd timer (`yardctl-sync.timer`) that runs `yardctl sync` every 5 minutes.
This fetches all upstream git repositories, searches for `docker-compose.yml` changes, and intelligently redeploys any stacks that have modifications.
