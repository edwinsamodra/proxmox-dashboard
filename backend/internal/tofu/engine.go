package tofu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Engine manages OpenTofu workspaces and executions.
type Engine struct {
	workspaceDir string
	binaryPath   string
}

// NewEngine creates a new OpenTofu engine.
func NewEngine(workspaceDir, binaryPath string) *Engine {
	if binaryPath == "" {
		binaryPath = "tofu"
	}
	return &Engine{
		workspaceDir: workspaceDir,
		binaryPath:   binaryPath,
	}
}

// Workspace represents an OpenTofu workspace.
type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Plan      string    `json:"plan,omitempty"`
}

// VMProvisionRequest holds parameters for VM provisioning through OpenTofu.
type VMProvisionRequest struct {
	WorkspaceName string            `json:"workspace_name"`
	NodeName      string            `json:"node_name"`
	VMID          int               `json:"vmid"`
	VMName        string            `json:"vm_name"`
	Template      string            `json:"template"`
	Cores         int               `json:"cores"`
	Memory        int               `json:"memory_mb"`
	DiskSize      int               `json:"disk_size_gb"`
	StoragePool   string            `json:"storage_pool"`
	Network       string            `json:"network"`
	IPConfig      string            `json:"ip_config"`
	SSHKeys       string            `json:"ssh_keys"`
	CloudInitUser string            `json:"cloud_init_user"`
	Tags          []string          `json:"tags"`
	ExtraVars     map[string]string `json:"extra_vars,omitempty"`
}

// ContainerProvisionRequest holds parameters for LXC container provisioning.
type ContainerProvisionRequest struct {
	WorkspaceName string            `json:"workspace_name"`
	NodeName      string            `json:"node_name"`
	VMID          int               `json:"vmid"`
	Hostname      string            `json:"hostname"`
	Template      string            `json:"template"`
	Cores         int               `json:"cores"`
	Memory        int               `json:"memory_mb"`
	DiskSize      int               `json:"disk_size_gb"`
	StoragePool   string            `json:"storage_pool"`
	Network       string            `json:"network"`
	IPConfig      string            `json:"ip_config"`
	Password      string            `json:"password,omitempty"`
	SSHKeys       string            `json:"ssh_keys"`
	Tags          []string          `json:"tags"`
	ExtraVars     map[string]string `json:"extra_vars,omitempty"`
}

// ExecResult holds the result of a tofu command.
type ExecResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// EnsureWorkspaceDir ensures the workspace directory exists.
func (e *Engine) EnsureWorkspaceDir() error {
	return os.MkdirAll(e.workspaceDir, 0750)
}

// WorkspaceDir returns the base workspace directory.
func (e *Engine) WorkspaceDir() string {
	return e.workspaceDir
}

// ListWorkspaces returns all workspaces.
func (e *Engine) ListWorkspaces() ([]Workspace, error) {
	entries, err := os.ReadDir(e.workspaceDir)
	if os.IsNotExist(err) {
		return []Workspace{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("listing workspaces: %w", err)
	}

	var workspaces []Workspace
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		info, _ := entry.Info()
		ws := Workspace{
			ID:        entry.Name(),
			Name:      entry.Name(),
			CreatedAt: info.ModTime(),
		}
		// Try to read plan file.
		planPath := filepath.Join(e.workspaceDir, entry.Name(), "main.tf")
		if data, err := os.ReadFile(planPath); err == nil {
			ws.Plan = string(data)
		}
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

// GenerateVMTerraform generates a Terraform configuration for a VM.
func (e *Engine) GenerateVMTerraform(req VMProvisionRequest) string {
	tags := `"` + strings.Join(req.Tags, `", "`) + `"`
	return fmt.Sprintf(`terraform {
  required_providers {
    proxmox = {
      source  = "telmate/proxmox"
      version = "~> 2.9"
    }
  }
}

provider "proxmox" {
  pm_api_url      = var.proxmox_api_url
  pm_user         = var.proxmox_user
  pm_password     = var.proxmox_password
  pm_tls_insecure = true
}

variable "proxmox_api_url" {}
variable "proxmox_user" {}
variable "proxmox_password" {}

resource "proxmox_vm_qemu" "%s" {
  name        = "%s"
  target_node = "%s"
  vmid        = %d
  clone       = "%s"
  cores       = %d
  memory      = %d
  scsihw      = "virtio-scsi-pci"
  boot        = "order=scsi0"
  agent       = 1
  os_type     = "cloud-init"
  ipconfig0   = "%s"
  sshkeys     = <<EOF
%s
EOF
  ciuser      = "%s"
  tags        = "%s"

  disk {
    slot    = "scsi0"
    size    = "%dG"
    type    = "disk"
    storage = "%s"
    iothread = 1
  }

  network {
    model  = "virtio"
    bridge = "%s"
  }
}

output "vm_id" {
  value = proxmox_vm_qemu.%s.vmid
}
`,
		sanitizeResourceName(req.VMName),
		req.VMName,
		req.NodeName,
		req.VMID,
		req.Template,
		req.Cores,
		req.Memory,
		req.IPConfig,
		req.SSHKeys,
		req.CloudInitUser,
		tags,
		req.DiskSize,
		req.StoragePool,
		req.Network,
		sanitizeResourceName(req.VMName),
	)
}

// GenerateContainerTerraform generates a Terraform configuration for an LXC container.
func (e *Engine) GenerateContainerTerraform(req ContainerProvisionRequest) string {
	return fmt.Sprintf(`terraform {
  required_providers {
    proxmox = {
      source  = "telmate/proxmox"
      version = "~> 2.9"
    }
  }
}

provider "proxmox" {
  pm_api_url      = var.proxmox_api_url
  pm_user         = var.proxmox_user
  pm_password     = var.proxmox_password
  pm_tls_insecure = true
}

variable "proxmox_api_url" {}
variable "proxmox_user" {}
variable "proxmox_password" {}

resource "proxmox_lxc" "%s" {
  hostname     = "%s"
  target_node  = "%s"
  vmid         = %d
  ostemplate   = "%s"
  cores        = %d
  memory       = %d
  swap         = 512
  unprivileged = true
  ssh_public_keys = <<EOF
%s
EOF

  rootfs {
    storage = "%s"
    size    = "%dG"
  }

  network {
    name   = "eth0"
    bridge = "%s"
    ip     = "%s"
  }
}

output "container_id" {
  value = proxmox_lxc.%s.vmid
}
`,
		sanitizeResourceName(req.Hostname),
		req.Hostname,
		req.NodeName,
		req.VMID,
		req.Template,
		req.Cores,
		req.Memory,
		req.SSHKeys,
		req.StoragePool,
		req.DiskSize,
		req.Network,
		req.IPConfig,
		sanitizeResourceName(req.Hostname),
	)
}

// CreateWorkspace creates a new workspace directory with a main.tf file.
func (e *Engine) CreateWorkspace(name, tfContent string) error {
	dir, err := e.workspacePath(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("creating workspace dir: %w", err)
	}
	tfPath := filepath.Join(dir, "main.tf")
	return os.WriteFile(tfPath, []byte(tfContent), 0600)
}

// GetWorkspaceTF returns the main.tf content of a workspace.
func (e *Engine) GetWorkspaceTF(name string) (string, error) {
	dir, err := e.workspacePath(name)
	if err != nil {
		return "", err
	}
	tfPath := filepath.Join(dir, "main.tf")
	data, err := os.ReadFile(tfPath)
	if err != nil {
		return "", fmt.Errorf("reading workspace tf: %w", err)
	}
	return string(data), nil
}

// SaveWorkspaceTF saves the main.tf content of a workspace.
func (e *Engine) SaveWorkspaceTF(name, content string) error {
	dir, err := e.workspacePath(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("creating workspace dir: %w", err)
	}
	tfPath := filepath.Join(dir, "main.tf")
	return os.WriteFile(tfPath, []byte(content), 0600)
}

// DeleteWorkspace removes a workspace directory.
func (e *Engine) DeleteWorkspace(name string) error {
	dir, err := e.workspacePath(name)
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

// Run executes a tofu command in a workspace.
func (e *Engine) Run(workspaceName string, args ...string) ExecResult {
	dir, err := e.workspacePath(workspaceName)
	if err != nil {
		return ExecResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	cmd := exec.Command(e.binaryPath, args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	combined := stdout.String() + stderr.String()

	if err != nil {
		return ExecResult{
			Success: false,
			Output:  combined,
			Error:   err.Error(),
		}
	}
	return ExecResult{
		Success: true,
		Output:  combined,
	}
}

func (e *Engine) workspacePath(name string) (string, error) {
	cleanName, err := sanitizeWorkspaceName(name)
	if err != nil {
		return "", err
	}

	baseDir := filepath.Clean(e.workspaceDir)
	targetDir := filepath.Join(baseDir, cleanName)
	rel, err := filepath.Rel(baseDir, targetDir)
	if err != nil {
		return "", fmt.Errorf("resolving workspace path: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid workspace path")
	}

	return targetDir, nil
}

var workspaceNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func sanitizeWorkspaceName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", fmt.Errorf("workspace name is required")
	}
	if filepath.IsAbs(trimmed) {
		return "", fmt.Errorf("invalid workspace name")
	}
	if strings.Contains(trimmed, "/") || strings.Contains(trimmed, "\\") {
		return "", fmt.Errorf("invalid workspace name")
	}

	clean := filepath.Clean(trimmed)
	if clean == "." || clean == ".." || clean != trimmed {
		return "", fmt.Errorf("invalid workspace name")
	}
	if !workspaceNamePattern.MatchString(clean) {
		return "", fmt.Errorf("invalid workspace name")
	}

	return clean, nil
}

// Init runs tofu init in the workspace.
func (e *Engine) Init(workspaceName string) ExecResult {
	return e.Run(workspaceName, "init", "-no-color")
}

// Plan runs tofu plan in the workspace.
func (e *Engine) Plan(workspaceName string, vars map[string]string) ExecResult {
	args := []string{"plan", "-no-color"}
	for k, v := range vars {
		args = append(args, fmt.Sprintf("-var=%s=%s", k, v))
	}
	return e.Run(workspaceName, args...)
}

// Apply runs tofu apply in the workspace.
func (e *Engine) Apply(workspaceName string, vars map[string]string) ExecResult {
	args := []string{"apply", "-auto-approve", "-no-color"}
	for k, v := range vars {
		args = append(args, fmt.Sprintf("-var=%s=%s", k, v))
	}
	return e.Run(workspaceName, args...)
}

// Destroy runs tofu destroy in the workspace.
func (e *Engine) Destroy(workspaceName string, vars map[string]string) ExecResult {
	args := []string{"destroy", "-auto-approve", "-no-color"}
	for k, v := range vars {
		args = append(args, fmt.Sprintf("-var=%s=%s", k, v))
	}
	return e.Run(workspaceName, args...)
}

// Output runs tofu output -json and returns the parsed result.
func (e *Engine) Output(workspaceName string) (map[string]interface{}, ExecResult) {
	result := e.Run(workspaceName, "output", "-json", "-no-color")
	if !result.Success {
		return nil, result
	}
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(result.Output), &out); err != nil {
		return nil, ExecResult{Success: false, Error: err.Error()}
	}
	return out, result
}

// IsAvailable checks whether the tofu binary is available in PATH.
func (e *Engine) IsAvailable() bool {
	_, err := exec.LookPath(e.binaryPath)
	return err == nil
}

// sanitizeResourceName replaces non-alphanumeric chars (except underscore) with underscore.
func sanitizeResourceName(name string) string {
	replacer := strings.NewReplacer(
		"-", "_",
		".", "_",
		" ", "_",
	)
	return replacer.Replace(name)
}
