package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/proxmox"
	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/tofu"
	wsHub "github.com/edwinsamodra/proxmox-dashboard/backend/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins in development; tighten in production.
		return true
	},
}

// NewRouter builds and returns the main HTTP ServeMux.
func NewRouter(px *proxmox.Client, hub *wsHub.Hub, tf *tofu.Engine) http.Handler {
	h := NewHandler(px, hub, tf)
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /api/health", h.HandleHealth)

	// Auth
	mux.HandleFunc("POST /api/auth/login", h.HandleLogin)
	mux.HandleFunc("POST /api/auth/token", h.HandleTokenAuth)

	// Version
	mux.HandleFunc("GET /api/version", h.HandleVersion)

	// Cluster
	mux.HandleFunc("GET /api/cluster/status", h.HandleGetClusterStatus)
	mux.HandleFunc("GET /api/cluster/firewall", h.HandleGetClusterFirewall)

	// Nodes
	mux.HandleFunc("GET /api/nodes", h.HandleGetNodes)
	mux.HandleFunc("GET /api/nodes/{node}/status", h.HandleGetNodeStatus)
	mux.HandleFunc("GET /api/nodes/{node}/tasks", h.HandleGetTasks)
	mux.HandleFunc("GET /api/nodes/{node}/tasks/log", h.HandleGetTaskLog)
	mux.HandleFunc("GET /api/nodes/{node}/network", h.HandleGetNetwork)
	mux.HandleFunc("GET /api/nodes/{node}/storage", h.HandleGetStorage)
	mux.HandleFunc("GET /api/nodes/{node}/storage/{storage}/content", h.HandleGetStorageContent)

	// VMs
	mux.HandleFunc("GET /api/vms", h.HandleGetAllVMs)
	mux.HandleFunc("GET /api/nodes/{node}/vms", h.HandleGetVMs)
	mux.HandleFunc("GET /api/nodes/{node}/vms/{vmid}/status", h.HandleGetVMStatus)
	mux.HandleFunc("GET /api/nodes/{node}/vms/{vmid}/config", h.HandleGetVMConfig)
	mux.HandleFunc("GET /api/nodes/{node}/vms/{vmid}/snapshots", h.HandleGetVMSnapshots)
	mux.HandleFunc("POST /api/nodes/{node}/vms/{vmid}/snapshots", h.HandleCreateSnapshot)
	mux.HandleFunc("POST /api/nodes/{node}/vms", h.HandleCreateVM)
	mux.HandleFunc("POST /api/nodes/{node}/vms/{vmid}/action/{action}", h.HandleVMAction)
	mux.HandleFunc("DELETE /api/nodes/{node}/vms/{vmid}", h.HandleDeleteVM)
	mux.HandleFunc("GET /api/nodes/{node}/vms/{vmid}/firewall", h.HandleGetVMFirewall)

	// Containers
	mux.HandleFunc("GET /api/containers", h.HandleGetAllContainers)
	mux.HandleFunc("GET /api/nodes/{node}/containers", h.HandleGetContainers)
	mux.HandleFunc("POST /api/nodes/{node}/containers/{vmid}/action/{action}", h.HandleContainerAction)
	mux.HandleFunc("DELETE /api/nodes/{node}/containers/{vmid}", h.HandleDeleteContainer)

	// Backup
	mux.HandleFunc("GET /api/backups", h.HandleGetBackups)

	// OpenTofu
	mux.HandleFunc("GET /api/tofu/status", h.HandleTofuStatus)
	mux.HandleFunc("GET /api/tofu/workspaces", h.HandleListWorkspaces)
	mux.HandleFunc("POST /api/tofu/workspaces", h.HandleCreateWorkspace)
	mux.HandleFunc("GET /api/tofu/workspaces/{name}/tf", h.HandleGetWorkspaceTF)
	mux.HandleFunc("PUT /api/tofu/workspaces/{name}/tf", h.HandleSaveWorkspaceTF)
	mux.HandleFunc("DELETE /api/tofu/workspaces/{name}", h.HandleDeleteWorkspace)
	mux.HandleFunc("POST /api/tofu/workspaces/{name}/init", h.HandleTofuInit)
	mux.HandleFunc("POST /api/tofu/workspaces/{name}/plan", h.HandleTofuPlan)
	mux.HandleFunc("POST /api/tofu/workspaces/{name}/apply", h.HandleTofuApply)
	mux.HandleFunc("POST /api/tofu/workspaces/{name}/destroy", h.HandleTofuDestroy)
	mux.HandleFunc("GET /api/tofu/workspaces/{name}/output", h.HandleTofuOutput)
	mux.HandleFunc("POST /api/tofu/generate/vm", h.HandleGenerateVMTF)
	mux.HandleFunc("POST /api/tofu/generate/container", h.HandleGenerateContainerTF)
	mux.HandleFunc("POST /api/tofu/provision/vm", h.HandleProvisionVM)

	// WebSocket
	mux.HandleFunc("GET /api/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		wsHub.NewClient(hub, conn)
	})
	mux.HandleFunc("GET /api/ws/stats", h.HandleWSStats)

	return corsMiddleware(loggingMiddleware(mux))
}

// corsMiddleware adds CORS headers to allow the SPA frontend.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Proxmox-Token-ID, X-Proxmox-Token-Secret")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs each request.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		if !strings.HasPrefix(r.URL.Path, "/api/ws") {
			// Skip WebSocket upgrade logs.
			_ = start
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
