package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/api"
	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/config"
	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/proxmox"
	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/tofu"
	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/ws"
)

func main() {
	// Load configuration.
	cfgPath := os.Getenv("PMO_CONFIG")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Ensure OpenTofu workspace directory exists.
	if err := os.MkdirAll(cfg.OpenTofu.WorkspaceDir, 0750); err != nil {
		log.Printf("Warning: could not create tofu workspace dir: %v", err)
	}

	// Initialize Proxmox client.
	pxClient := proxmox.NewClient(cfg.Proxmox.BaseURL(), cfg.Proxmox.Insecure)

	// Auto-login if credentials are provided via environment.
	pxUser := os.Getenv("PROXMOX_USER")
	pxPass := os.Getenv("PROXMOX_PASSWORD")
	pxTokenID := os.Getenv("PROXMOX_TOKEN_ID")
	pxTokenSecret := os.Getenv("PROXMOX_TOKEN_SECRET")

	switch {
	case pxTokenID != "" && pxTokenSecret != "":
		pxClient.SetToken(pxTokenID, pxTokenSecret)
		log.Println("Proxmox: using API token authentication")
	case pxUser != "" && pxPass != "":
		realm := os.Getenv("PROXMOX_REALM")
		if realm == "" {
			realm = "pam"
		}
		if err := pxClient.LoginFull(pxUser, pxPass, realm); err != nil {
			log.Printf("Warning: Proxmox auto-login failed: %v", err)
		} else {
			log.Printf("Proxmox: authenticated as %s@%s", pxUser, realm)
		}
	default:
		log.Println("Proxmox: no credentials configured — use /api/auth/login or /api/auth/token")
	}

	// Initialize WebSocket hub.
	hub := ws.NewHub()
	go hub.Run()

	// Initialize OpenTofu engine.
	tofuBin := cfg.OpenTofu.BinaryPath
	if v := os.Getenv("PMO_TOFU_BIN"); v != "" {
		tofuBin = v
	}
	tofuEngine := tofu.NewEngine(cfg.OpenTofu.WorkspaceDir, tofuBin)

	// Start background poller that broadcasts node stats every 10 seconds.
	go pollAndBroadcast(pxClient, hub)

	// Build HTTP router.
	router := api.NewRouter(pxClient, hub, tofuEngine)

	addr := cfg.Server.Addr()
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown.
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("PMO Backend listening on http://%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	log.Println("Server stopped.")
}

// pollAndBroadcast polls Proxmox every 10 seconds and broadcasts stats via WebSocket.
func pollAndBroadcast(px *proxmox.Client, hub *ws.Hub) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !px.IsAuthenticated() {
			continue
		}
		nodes, err := px.GetNodes()
		if err != nil {
			continue
		}
		hub.Broadcast("nodes_update", nodes)

		// Also broadcast all VMs.
		nodeNames := make([]string, len(nodes))
		for i, n := range nodes {
			nodeNames[i] = n.Node
		}
		vms, err := px.GetAllVMs(nodeNames)
		if err == nil {
			hub.Broadcast("vms_update", vms)
		}
		ctrs, err := px.GetAllContainers(nodeNames)
		if err == nil {
			hub.Broadcast("containers_update", ctrs)
		}
	}
}
