package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/proxmox"
	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/tofu"
	"github.com/edwinsamodra/proxmox-dashboard/backend/internal/ws"
)

// Handler holds dependencies for all HTTP handlers.
type Handler struct {
	proxmox *proxmox.Client
	hub     *ws.Hub
	tofu    *tofu.Engine
}

// NewHandler creates a new Handler.
func NewHandler(px *proxmox.Client, hub *ws.Hub, tf *tofu.Engine) *Handler {
	return &Handler{proxmox: px, hub: hub, tofu: tf}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("api: encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func pathParam(r *http.Request, name string) string {
	// Compatible with Go 1.22 ServeMux path params.
	return r.PathValue(name)
}

func intParam(r *http.Request, name string) (int, error) {
	s := pathParam(r, name)
	if s == "" {
		return 0, fmt.Errorf("missing path parameter %q", name)
	}
	return strconv.Atoi(s)
}

// ─────────────────────────────────────────────────────────────────────────────
// Auth
// ─────────────────────────────────────────────────────────────────────────────

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Realm    string `json:"realm"`
}

// HandleLogin authenticates to Proxmox and stores the session.
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Realm == "" {
		req.Realm = "pam"
	}
	if err := h.proxmox.LoginFull(req.Username, req.Password, req.Realm); err != nil {
		writeError(w, http.StatusUnauthorized, "authentication failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"authenticated": true})
}

// HandleTokenAuth sets a Proxmox API token for authentication.
func (h *Handler) HandleTokenAuth(w http.ResponseWriter, r *http.Request) {
	tokenID := r.Header.Get("X-Proxmox-Token-ID")
	tokenSecret := r.Header.Get("X-Proxmox-Token-Secret")
	if tokenID == "" || tokenSecret == "" {
		writeError(w, http.StatusBadRequest, "X-Proxmox-Token-ID and X-Proxmox-Token-Secret headers required")
		return
	}
	h.proxmox.SetToken(tokenID, tokenSecret)
	writeJSON(w, http.StatusOK, map[string]bool{"authenticated": true})
}

// ─────────────────────────────────────────────────────────────────────────────
// Version / Health
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) HandleVersion(w http.ResponseWriter, r *http.Request) {
	version, err := h.proxmox.GetVersion()
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, version)
}

// ─────────────────────────────────────────────────────────────────────────────
// Nodes
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleGetNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.proxmox.GetNodes()
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, nodes)
}

func (h *Handler) HandleGetNodeStatus(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	status, err := h.proxmox.GetNodeStatus(node)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (h *Handler) HandleGetClusterStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.proxmox.GetClusterStatus()
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, status)
}

// ─────────────────────────────────────────────────────────────────────────────
// VMs
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleGetVMs(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vms, err := h.proxmox.GetVMs(node)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, vms)
}

func (h *Handler) HandleGetAllVMs(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.proxmox.GetNodes()
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	names := make([]string, len(nodes))
	for i, n := range nodes {
		names[i] = n.Node
	}
	vms, err := h.proxmox.GetAllVMs(names)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, vms)
}

func (h *Handler) HandleGetVMStatus(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	status, err := h.proxmox.GetVMStatus(node, vmid)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func (h *Handler) HandleVMAction(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	action := pathParam(r, "action")

	validActions := map[string]bool{
		"start": true, "stop": true, "reboot": true,
		"shutdown": true, "suspend": true, "resume": true,
		"reset": true,
	}
	if !validActions[action] {
		writeError(w, http.StatusBadRequest, "invalid action: "+action)
		return
	}

	result, err := h.proxmox.VMAction(node, vmid, action)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	h.hub.Broadcast("vm_action", map[string]interface{}{
		"node":   node,
		"vmid":   vmid,
		"action": action,
		"task":   result,
	})

	writeJSON(w, http.StatusOK, map[string]string{"task": result})
}

func (h *Handler) HandleCreateVM(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")

	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	// Accept JSON body and convert to url.Values.
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
		params := url.Values{}
		for k, v := range body {
			params.Set(k, fmt.Sprintf("%v", v))
		}
		result, err := h.proxmox.CreateVM(node, params)
		if err != nil {
			writeError(w, http.StatusBadGateway, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"task": result})
		return
	}

	result, err := h.proxmox.CreateVM(node, r.Form)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"task": result})
}

func (h *Handler) HandleDeleteVM(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.proxmox.DeleteVM(node, vmid)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"task": result})
}

func (h *Handler) HandleGetVMConfig(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	config, err := h.proxmox.GetVMConfig(node, vmid)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, config)
}

func (h *Handler) HandleGetVMSnapshots(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	snapshots, err := h.proxmox.GetVMSnapshots(node, vmid)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, snapshots)
}

type createSnapshotRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *Handler) HandleCreateSnapshot(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var req createSnapshotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeError(w, http.StatusBadRequest, "snapshot name required")
		return
	}
	task, err := h.proxmox.CreateSnapshot(node, vmid, req.Name, req.Description)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"task": task})
}

// ─────────────────────────────────────────────────────────────────────────────
// Containers
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleGetContainers(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	ctrs, err := h.proxmox.GetContainers(node)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ctrs)
}

func (h *Handler) HandleGetAllContainers(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.proxmox.GetNodes()
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	names := make([]string, len(nodes))
	for i, n := range nodes {
		names[i] = n.Node
	}
	ctrs, err := h.proxmox.GetAllContainers(names)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ctrs)
}

func (h *Handler) HandleContainerAction(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	action := pathParam(r, "action")

	validActions := map[string]bool{
		"start": true, "stop": true, "reboot": true,
		"shutdown": true, "suspend": true, "resume": true,
	}
	if !validActions[action] {
		writeError(w, http.StatusBadRequest, "invalid action: "+action)
		return
	}

	result, err := h.proxmox.ContainerAction(node, vmid, action)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}

	h.hub.Broadcast("container_action", map[string]interface{}{
		"node":   node,
		"vmid":   vmid,
		"action": action,
		"task":   result,
	})

	writeJSON(w, http.StatusOK, map[string]string{"task": result})
}

func (h *Handler) HandleDeleteContainer(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.proxmox.DeleteContainer(node, vmid)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"task": result})
}

// ─────────────────────────────────────────────────────────────────────────────
// Storage
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleGetStorage(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	storages, err := h.proxmox.GetStorage(node)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, storages)
}

func (h *Handler) HandleGetStorageContent(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	storage := pathParam(r, "storage")
	content, err := h.proxmox.GetStorageContent(node, storage)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, content)
}

// ─────────────────────────────────────────────────────────────────────────────
// Network
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleGetNetwork(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	ifaces, err := h.proxmox.GetNetworkInterfaces(node)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, ifaces)
}

// ─────────────────────────────────────────────────────────────────────────────
// Firewall
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleGetClusterFirewall(w http.ResponseWriter, r *http.Request) {
	rules, err := h.proxmox.GetClusterFirewallRules()
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rules)
}

func (h *Handler) HandleGetVMFirewall(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	vmid, err := intParam(r, "vmid")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rules, err := h.proxmox.GetVMFirewallRules(node, vmid)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rules)
}

// ─────────────────────────────────────────────────────────────────────────────
// Backup
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleGetBackups(w http.ResponseWriter, r *http.Request) {
	backups, err := h.proxmox.GetBackups()
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, backups)
}

// ─────────────────────────────────────────────────────────────────────────────
// Tasks
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleGetTasks(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	tasks, err := h.proxmox.GetNodeTasks(node)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) HandleGetTaskLog(w http.ResponseWriter, r *http.Request) {
	node := pathParam(r, "node")
	upid := r.URL.Query().Get("upid")
	if upid == "" {
		writeError(w, http.StatusBadRequest, "upid query param required")
		return
	}
	log_, err := h.proxmox.GetTaskLog(node, upid)
	if err != nil {
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, log_)
}

// ─────────────────────────────────────────────────────────────────────────────
// OpenTofu
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleTofuStatus(w http.ResponseWriter, r *http.Request) {
	available := h.tofu.IsAvailable()
	workspaces, _ := h.tofu.ListWorkspaces()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"available":   available,
		"workspaces":  workspaces,
		"binary_path": os.Getenv("PMO_TOFU_BIN"),
	})
}

func (h *Handler) HandleListWorkspaces(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.tofu.ListWorkspaces()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, workspaces)
}

type createWorkspaceRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (h *Handler) HandleCreateWorkspace(w http.ResponseWriter, r *http.Request) {
	var req createWorkspaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		writeError(w, http.StatusBadRequest, "workspace name required")
		return
	}
	if err := h.tofu.CreateWorkspace(req.Name, req.Content); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"workspace": req.Name})
}

func (h *Handler) HandleGetWorkspaceTF(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	content, err := h.tofu.GetWorkspaceTF(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"content": content})
}

func (h *Handler) HandleSaveWorkspaceTF(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.tofu.SaveWorkspaceTF(name, req.Content); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"saved": true})
}

func (h *Handler) HandleDeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	if err := h.tofu.DeleteWorkspace(name); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (h *Handler) HandleTofuInit(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	result := h.tofu.Init(name)
	code := http.StatusOK
	if !result.Success {
		code = http.StatusInternalServerError
	}
	writeJSON(w, code, result)
}

func (h *Handler) HandleTofuPlan(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	vars := extractVars(r)
	result := h.tofu.Plan(name, vars)
	code := http.StatusOK
	if !result.Success {
		code = http.StatusInternalServerError
	}
	writeJSON(w, code, result)
}

func (h *Handler) HandleTofuApply(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	vars := extractVars(r)
	result := h.tofu.Apply(name, vars)
	code := http.StatusOK
	if !result.Success {
		code = http.StatusInternalServerError
	}
	h.hub.Broadcast("tofu_apply", map[string]interface{}{
		"workspace": name,
		"success":   result.Success,
	})
	writeJSON(w, code, result)
}

func (h *Handler) HandleTofuDestroy(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	vars := extractVars(r)
	result := h.tofu.Destroy(name, vars)
	code := http.StatusOK
	if !result.Success {
		code = http.StatusInternalServerError
	}
	writeJSON(w, code, result)
}

func (h *Handler) HandleTofuOutput(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	out, result := h.tofu.Output(name)
	if !result.Success {
		writeError(w, http.StatusInternalServerError, result.Error)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

// HandleGenerateVMTF generates a Terraform config from a VM provision request.
func (h *Handler) HandleGenerateVMTF(w http.ResponseWriter, r *http.Request) {
	var req tofu.VMProvisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	content := h.tofu.GenerateVMTerraform(req)
	writeJSON(w, http.StatusOK, map[string]string{"content": content})
}

// HandleGenerateContainerTF generates a Terraform config from a container provision request.
func (h *Handler) HandleGenerateContainerTF(w http.ResponseWriter, r *http.Request) {
	var req tofu.ContainerProvisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	content := h.tofu.GenerateContainerTerraform(req)
	writeJSON(w, http.StatusOK, map[string]string{"content": content})
}

// HandleProvisionVM is an all-in-one handler: generates TF, creates workspace, runs apply.
func (h *Handler) HandleProvisionVM(w http.ResponseWriter, r *http.Request) {
	var req tofu.VMProvisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.WorkspaceName == "" {
		req.WorkspaceName = "vm-" + req.VMName
	}

	tfContent := h.tofu.GenerateVMTerraform(req)
	if err := h.tofu.CreateWorkspace(req.WorkspaceName, tfContent); err != nil {
		writeError(w, http.StatusInternalServerError, "creating workspace: "+err.Error())
		return
	}

	initResult := h.tofu.Init(req.WorkspaceName)
	if !initResult.Success {
		writeError(w, http.StatusInternalServerError, "tofu init failed: "+initResult.Error)
		return
	}

	vars := map[string]string{
		"proxmox_api_url":  os.Getenv("PROXMOX_API_URL"),
		"proxmox_user":     os.Getenv("PROXMOX_USER"),
		"proxmox_password": os.Getenv("PROXMOX_PASSWORD"),
	}
	for k, v := range req.ExtraVars {
		vars[k] = v
	}

	applyResult := h.tofu.Apply(req.WorkspaceName, vars)
	h.hub.Broadcast("vm_provisioned", map[string]interface{}{
		"workspace": req.WorkspaceName,
		"vm_name":   req.VMName,
		"success":   applyResult.Success,
	})

	code := http.StatusCreated
	if !applyResult.Success {
		code = http.StatusInternalServerError
	}
	writeJSON(w, code, applyResult)
}

// ─────────────────────────────────────────────────────────────────────────────
// WebSocket stats
// ─────────────────────────────────────────────────────────────────────────────

func (h *Handler) HandleWSStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]int{
		"connected_clients": h.hub.ClientCount(),
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

func extractVars(r *http.Request) map[string]string {
	vars := make(map[string]string)
	for k, v := range r.URL.Query() {
		if strings.HasPrefix(k, "var_") {
			vars[strings.TrimPrefix(k, "var_")] = v[0]
		}
	}
	return vars
}
