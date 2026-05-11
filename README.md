# PMO — Proxmox Modern Orchestrator

> A cloud-native, hyperscaler-inspired alternative dashboard for Proxmox VE.
> Dark mode · Real-time · Infrastructure as Code

![Architecture](https://img.shields.io/badge/Frontend-Vue%203%20%2B%20TypeScript-42b883?logo=vuedotjs)
![Backend](https://img.shields.io/badge/Backend-Go%201.21-00ADD8?logo=go)
![IaC](https://img.shields.io/badge/IaC-OpenTofu-purple)
![Docker](https://img.shields.io/badge/Deploy-Docker-2496ED?logo=docker)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🌑 **Dark Mode First** | GCP / AWS Amplify-inspired UI with Tailwind CSS |
| ⚡ **Real-time Updates** | WebSocket hub pushes node/VM/container stats every 10s |
| 🖥️ **Multi-Node Cluster** | Manage all nodes in a single integrated dashboard |
| 🚀 **VM / LXC Lifecycle** | Start, stop, reboot, suspend, delete — with optimistic UI |
| 💾 **Storage Management** | Per-node storage usage, content browser |
| 🌐 **Network View** | Network interfaces, bridges, and routing |
| 🔥 **Firewall** | Cluster and VM-level firewall rule inspection |
| 📦 **Backup Jobs** | View and manage vzdump backup schedules |
| 🌿 **OpenTofu IaC** | **Killer Feature** — Zero-Touch Provisioning via Terraform/OpenTofu plans |
| 🔄 **Skeleton Loaders** | Every data-fetching state shows proper skeleton animations |
| 🐳 **Container Deploy** | < 100 MB RAM backend footprint via static Go binary |

---

## 🏗 Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│  Browser (SPA)                                                  │
│  Vue 3 + TypeScript + Pinia + Tailwind CSS + Vue Router         │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────┐ │
│  │  Dashboard   │  │  VM / LXC    │  │  OpenTofu IaC         │ │
│  │  Nodes       │  │  Storage     │  │  Workspace Editor     │ │
│  │  Cluster     │  │  Network     │  │  Plan/Apply/Destroy   │ │
│  └──────────────┘  └──────────────┘  └───────────────────────┘ │
└───────────────────────────────┬─────────────────────────────────┘
                                │ HTTP/WebSocket (same origin)
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│  PMO Backend  (Go — static binary)                              │
│  net/http · gorilla/websocket                                   │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────────┐ │
│  │  REST API    │  │  WS Hub      │  │  OpenTofu Engine      │ │
│  │  /api/*      │  │  /api/ws     │  │  init/plan/apply/     │ │
│  │              │  │  broadcaster │  │  destroy/output       │ │
│  └──────────────┘  └──────────────┘  └───────────────────────┘ │
│                                                                 │
│  ┌──────────────┐  ┌───────────────────────────────────────┐   │
│  │ Proxmox API  │  │  SSH Client (golang.org/x/crypto/ssh) │   │
│  │ Client       │  │  For low-level node access            │   │
│  └──────────────┘  └───────────────────────────────────────┘   │
└───────────────────────────────┬─────────────────────────────────┘
                                │ HTTPS (Proxmox REST API)
                                ▼
                  ┌─────────────────────────┐
                  │  Proxmox VE Cluster      │
                  │  Node 1 · Node 2 · ...   │
                  └─────────────────────────┘
```

---

## 🚀 Quick Start

### With Docker Compose (recommended)

```bash
cp config.example.yml config.yml
# Edit config.yml with your Proxmox host

# Set environment variables
export PROXMOX_HOST=192.168.1.100
export PROXMOX_USER=root
export PROXMOX_PASSWORD=your-password

docker compose up -d
```

Open **http://localhost:8080** in your browser.

### Development mode

**Backend:**
```bash
cd backend
export PMO_PROXMOX_HOST=192.168.1.100
export PROXMOX_USER=root
export PROXMOX_PASSWORD=your-password
go run ./cmd/server
# Listening on http://0.0.0.0:8080
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
# Dev server on http://localhost:5173 (proxies /api → :8080)
```

---

## 🔐 Authentication

PMO supports two authentication methods:

### Username + Password
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"root","password":"secret","realm":"pam"}'
```

### API Token (recommended)
```bash
curl -X POST http://localhost:8080/api/auth/token \
  -H 'X-Proxmox-Token-ID: root@pam!mytoken' \
  -H 'X-Proxmox-Token-Secret: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
```

Or set via environment variables at startup:
```bash
PROXMOX_TOKEN_ID=root@pam!mytoken
PROXMOX_TOKEN_SECRET=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

---

## 🌿 OpenTofu IaC (Zero-Touch Provisioning)

PMO includes a built-in Terraform/OpenTofu workspace manager. No more Click-Ops!

### How it works

1. Navigate to **OpenTofu IaC** in the sidebar
2. Click **New Provision** and fill in the VM parameters
3. PMO auto-generates the `.tf` configuration
4. Click **Provision Now** — PMO runs `tofu init` + `tofu apply` automatically
5. Your VM is created with Cloud-init, network config, and SSH keys pre-configured

### Manual workspace management

```bash
# Check status
GET /api/tofu/status

# Create workspace with your own .tf
POST /api/tofu/workspaces
{"name": "my-stack", "content": "...tf content..."}

# Run init / plan / apply
POST /api/tofu/workspaces/my-stack/init
POST /api/tofu/workspaces/my-stack/plan
POST /api/tofu/workspaces/my-stack/apply
```

### Required: Proxmox Terraform Provider

The generated configurations use the [telmate/proxmox](https://registry.terraform.io/providers/telmate/proxmox/latest) provider.

Install OpenTofu:
```bash
# Linux
wget https://github.com/opentofu/opentofu/releases/download/v1.6.0/tofu_1.6.0_linux_amd64.zip
unzip tofu_*.zip
sudo mv tofu /usr/local/bin/
```

---

## 📡 WebSocket Real-Time Events

Connect to `ws://localhost:8080/api/ws` to receive live updates:

| Event | Payload |
|-------|---------|
| `nodes_update` | `Node[]` — full node list with stats |
| `vms_update` | `VM[]` — all VMs with current status |
| `containers_update` | `Container[]` — all LXC containers |
| `vm_action` | `{ node, vmid, action, task }` |
| `container_action` | `{ node, vmid, action, task }` |
| `tofu_apply` | `{ workspace, success }` |
| `vm_provisioned` | `{ workspace, vm_name, success }` |

---

## 🔌 REST API Reference

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/health` | Health check |
| `GET` | `/api/version` | Proxmox version info |
| `POST` | `/api/auth/login` | Username/password login |
| `POST` | `/api/auth/token` | API token auth |
| `GET` | `/api/nodes` | List all nodes |
| `GET` | `/api/nodes/{node}/status` | Node detailed status |
| `GET` | `/api/cluster/status` | Cluster status |
| `GET` | `/api/vms` | All VMs (all nodes) |
| `GET` | `/api/nodes/{node}/vms` | VMs on a node |
| `POST` | `/api/nodes/{node}/vms/{vmid}/action/{action}` | VM lifecycle action |
| `DELETE` | `/api/nodes/{node}/vms/{vmid}` | Delete VM |
| `GET` | `/api/containers` | All containers |
| `POST` | `/api/nodes/{node}/containers/{vmid}/action/{action}` | Container action |
| `GET` | `/api/nodes/{node}/storage` | Storage list |
| `GET` | `/api/nodes/{node}/network` | Network interfaces |
| `GET` | `/api/cluster/firewall` | Cluster firewall rules |
| `GET` | `/api/backups` | Backup jobs |
| `GET` | `/api/tofu/status` | OpenTofu availability |
| `POST` | `/api/tofu/provision/vm` | Zero-touch VM provision |

---

## 🛠 Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Vite 5 · Vue 3.4 · TypeScript · Pinia · Vue Router 4 |
| Styling | Tailwind CSS 3 · Custom dark theme |
| Icons | Lucide Vue Next |
| Backend | Go 1.21 · `net/http` (no heavy frameworks) |
| WebSocket | `gorilla/websocket` |
| SSH | `golang.org/x/crypto/ssh` |
| IaC | OpenTofu / Terraform (telmate/proxmox provider) |
| Deploy | Docker multi-stage · Alpine base |
| Config | YAML + environment variables |

---

## 🐳 Production Deployment

```bash
# Build image (~50 MB)
docker build -t pmo:latest .

# Run with API token (most secure)
docker run -d \
  -p 8080:8080 \
  -e PMO_PROXMOX_HOST=192.168.1.100 \
  -e PROXMOX_TOKEN_ID=root@pam!pmo \
  -e PROXMOX_TOKEN_SECRET=your-secret \
  -v pmo-workspaces:/app/workspaces \
  --name pmo \
  pmo:latest
```

---

## 📁 Project Structure

```
proxmox-dashboard/
├── backend/
│   ├── cmd/server/main.go          # Entry point + HTTP server
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers.go         # All HTTP handlers
│   │   │   └── router.go           # Route definitions + middleware
│   │   ├── config/config.go        # Configuration loading
│   │   ├── proxmox/
│   │   │   ├── client.go           # Full Proxmox REST API client
│   │   │   └── types.go            # Proxmox data types
│   │   ├── ws/hub.go               # WebSocket hub + client
│   │   ├── tofu/engine.go          # OpenTofu workspace manager
│   │   └── ssh/client.go           # SSH client (advanced features)
│   ├── go.mod
│   └── go.sum
│
├── frontend/
│   ├── src/
│   │   ├── api/index.ts            # API client (fetch-based)
│   │   ├── assets/main.css         # Tailwind + component styles
│   │   ├── components/
│   │   │   ├── layout/             # AppLayout, AppSidebar, AppTopBar
│   │   │   └── ui/                 # StatusBadge, UsageBar, SkeletonLoader, Modal
│   │   ├── router/index.ts         # Vue Router routes
│   │   ├── stores/                 # Pinia stores
│   │   ├── types/index.ts          # TypeScript types
│   │   └── views/                  # Page components
│   ├── package.json
│   ├── tailwind.config.js
│   └── vite.config.ts
│
├── Dockerfile                      # Multi-stage build (Go + Node → Alpine)
├── docker-compose.yml
├── config.example.yml
└── README.md
```

---

## License

MIT — contributions welcome!
