package proxmox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a thin Proxmox REST API client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	ticket     string // PVEAuthCookie value
	csrfToken  string
}

// NewClient creates a new Proxmox API client.
func NewClient(baseURL string, insecure bool) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}, //nolint:gosec
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

// LoginResponse is returned by /access/ticket.
type LoginResponse struct {
	Data struct {
		Ticket              string `json:"ticket"`
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
		Username            string `json:"username"`
	} `json:"data"`
}

// Login authenticates and stores the session ticket.
func (c *Client) Login(username, password string) error {
	return c.LoginFull(username, password, "pam")
}

// LoginFull authenticates with an explicit realm.
func (c *Client) LoginFull(username, password, realm string) error {
	form := url.Values{}
	form.Set("username", username+"@"+realm)
	form.Set("password", password)

	resp, err := c.httpClient.Post(
		c.baseURL+"/access/ticket",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed (%d): %s", resp.StatusCode, string(body))
	}

	var lr LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return fmt.Errorf("decoding login response: %w", err)
	}

	c.ticket = lr.Data.Ticket
	c.csrfToken = lr.Data.CSRFPreventionToken
	return nil
}

// SetToken allows using an API token instead of user/password login.
func (c *Client) SetToken(tokenID, secret string) {
	c.ticket = "PVEAPIToken=" + tokenID + "=" + secret
}

// IsAuthenticated returns whether the client has auth credentials.
func (c *Client) IsAuthenticated() bool {
	return c.ticket != ""
}

// get performs an authenticated GET request.
func (c *Client) get(path string, result interface{}) error {
	return c.do(http.MethodGet, path, nil, result)
}

// post performs an authenticated POST request with JSON body.
func (c *Client) post(path string, body interface{}, result interface{}) error {
	return c.do(http.MethodPost, path, body, result)
}

// put performs an authenticated PUT request with JSON body.
func (c *Client) put(path string, body interface{}, result interface{}) error {
	return c.do(http.MethodPut, path, body, result)
}

// delete performs an authenticated DELETE request.
func (c *Client) delete(path string, result interface{}) error {
	return c.do(http.MethodDelete, path, nil, result)
}

func (c *Client) do(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	contentType := ""

	if body != nil {
		switch v := body.(type) {
		case url.Values:
			bodyReader = strings.NewReader(v.Encode())
			contentType = "application/x-www-form-urlencoded"
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return fmt.Errorf("marshaling request body: %w", err)
			}
			bodyReader = bytes.NewReader(data)
			contentType = "application/json"
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	if strings.HasPrefix(c.ticket, "PVEAPIToken=") {
		req.Header.Set("Authorization", c.ticket)
	} else if c.ticket != "" {
		req.Header.Set("Cookie", "PVEAuthCookie="+c.ticket)
		if method != http.MethodGet {
			req.Header.Set("CSRFPreventionToken", c.csrfToken)
		}
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("proxmox API error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Cluster / Node
// ─────────────────────────────────────────────────────────────────────────────

type apiListResponse struct {
	Data json.RawMessage `json:"data"`
}

// GetVersion returns PVE version info.
func (c *Client) GetVersion() (map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get("/version", &resp); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetNodes returns a list of nodes in the cluster.
func (c *Client) GetNodes() ([]Node, error) {
	var resp apiListResponse
	if err := c.get("/nodes", &resp); err != nil {
		return nil, err
	}
	var nodes []Node
	if err := json.Unmarshal(resp.Data, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNodeStatus returns detailed status for a single node.
func (c *Client) GetNodeStatus(node string) (map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get("/nodes/"+node+"/status", &resp); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetClusterStatus returns cluster-wide status.
func (c *Client) GetClusterStatus() ([]map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get("/cluster/status", &resp); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// VMs (QEMU)
// ─────────────────────────────────────────────────────────────────────────────

// GetVMs returns all VMs on a node.
func (c *Client) GetVMs(node string) ([]VM, error) {
	var resp apiListResponse
	if err := c.get("/nodes/"+node+"/qemu", &resp); err != nil {
		return nil, err
	}
	var vms []VM
	if err := json.Unmarshal(resp.Data, &vms); err != nil {
		return nil, err
	}
	return vms, nil
}

// GetAllVMs returns VMs from all nodes.
func (c *Client) GetAllVMs(nodes []string) ([]VM, error) {
	var all []VM
	for _, node := range nodes {
		vms, err := c.GetVMs(node)
		if err != nil {
			continue // best-effort
		}
		for i := range vms {
			vms[i].Node = node
		}
		all = append(all, vms...)
	}
	return all, nil
}

// GetVMStatus returns detailed runtime status for a VM.
func (c *Client) GetVMStatus(node string, vmid int) (map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get(fmt.Sprintf("/nodes/%s/qemu/%d/status/current", node, vmid), &resp); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// VMAction sends a lifecycle command (start/stop/reboot/shutdown/suspend/resume).
func (c *Client) VMAction(node string, vmid int, action string) (string, error) {
	var resp apiListResponse
	path := fmt.Sprintf("/nodes/%s/qemu/%d/status/%s", node, vmid, action)
	if err := c.post(path, url.Values{}, &resp); err != nil {
		return "", err
	}
	return string(resp.Data), nil
}

// CreateVM creates a new QEMU VM.
func (c *Client) CreateVM(node string, params url.Values) (string, error) {
	var resp apiListResponse
	if err := c.post("/nodes/"+node+"/qemu", params, &resp); err != nil {
		return "", err
	}
	return string(resp.Data), nil
}

// UpdateVMConfig updates VM configuration.
func (c *Client) UpdateVMConfig(node string, vmid int, params url.Values) error {
	var resp apiListResponse
	return c.put(fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmid), params, &resp)
}

// DeleteVM deletes a VM.
func (c *Client) DeleteVM(node string, vmid int) (string, error) {
	var resp apiListResponse
	if err := c.delete(fmt.Sprintf("/nodes/%s/qemu/%d", node, vmid), &resp); err != nil {
		return "", err
	}
	return string(resp.Data), nil
}

// GetVMConfig retrieves the VM configuration.
func (c *Client) GetVMConfig(node string, vmid int) (map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get(fmt.Sprintf("/nodes/%s/qemu/%d/config", node, vmid), &resp); err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Containers (LXC)
// ─────────────────────────────────────────────────────────────────────────────

// GetContainers returns all LXC containers on a node.
func (c *Client) GetContainers(node string) ([]Container, error) {
	var resp apiListResponse
	if err := c.get("/nodes/"+node+"/lxc", &resp); err != nil {
		return nil, err
	}
	var containers []Container
	if err := json.Unmarshal(resp.Data, &containers); err != nil {
		return nil, err
	}
	return containers, nil
}

// GetAllContainers returns containers from all nodes.
func (c *Client) GetAllContainers(nodes []string) ([]Container, error) {
	var all []Container
	for _, node := range nodes {
		ctrs, err := c.GetContainers(node)
		if err != nil {
			continue
		}
		for i := range ctrs {
			ctrs[i].Node = node
		}
		all = append(all, ctrs...)
	}
	return all, nil
}

// ContainerAction sends a lifecycle command to an LXC container.
func (c *Client) ContainerAction(node string, vmid int, action string) (string, error) {
	var resp apiListResponse
	path := fmt.Sprintf("/nodes/%s/lxc/%d/status/%s", node, vmid, action)
	if err := c.post(path, url.Values{}, &resp); err != nil {
		return "", err
	}
	return string(resp.Data), nil
}

// CreateContainer creates a new LXC container.
func (c *Client) CreateContainer(node string, params url.Values) (string, error) {
	var resp apiListResponse
	if err := c.post("/nodes/"+node+"/lxc", params, &resp); err != nil {
		return "", err
	}
	return string(resp.Data), nil
}

// DeleteContainer deletes an LXC container.
func (c *Client) DeleteContainer(node string, vmid int) (string, error) {
	var resp apiListResponse
	if err := c.delete(fmt.Sprintf("/nodes/%s/lxc/%d", node, vmid), &resp); err != nil {
		return "", err
	}
	return string(resp.Data), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Storage
// ─────────────────────────────────────────────────────────────────────────────

// GetStorage returns storage list for a node.
func (c *Client) GetStorage(node string) ([]Storage, error) {
	var resp apiListResponse
	if err := c.get("/nodes/"+node+"/storage", &resp); err != nil {
		return nil, err
	}
	var storages []Storage
	if err := json.Unmarshal(resp.Data, &storages); err != nil {
		return nil, err
	}
	return storages, nil
}

// GetStorageContent lists content (images, ISOs, etc.) on a storage.
func (c *Client) GetStorageContent(node, storage string) ([]map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get(fmt.Sprintf("/nodes/%s/storage/%s/content", node, storage), &resp); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Network
// ─────────────────────────────────────────────────────────────────────────────

// GetNetworkInterfaces returns network interfaces for a node.
func (c *Client) GetNetworkInterfaces(node string) ([]NetworkInterface, error) {
	var resp apiListResponse
	if err := c.get("/nodes/"+node+"/network", &resp); err != nil {
		return nil, err
	}
	var ifaces []NetworkInterface
	if err := json.Unmarshal(resp.Data, &ifaces); err != nil {
		return nil, err
	}
	return ifaces, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Firewall
// ─────────────────────────────────────────────────────────────────────────────

// GetClusterFirewallRules returns cluster-level firewall rules.
func (c *Client) GetClusterFirewallRules() ([]FirewallRule, error) {
	var resp apiListResponse
	if err := c.get("/cluster/firewall/rules", &resp); err != nil {
		return nil, err
	}
	var rules []FirewallRule
	if err := json.Unmarshal(resp.Data, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

// GetVMFirewallRules returns VM-level firewall rules.
func (c *Client) GetVMFirewallRules(node string, vmid int) ([]FirewallRule, error) {
	var resp apiListResponse
	if err := c.get(fmt.Sprintf("/nodes/%s/qemu/%d/firewall/rules", node, vmid), &resp); err != nil {
		return nil, err
	}
	var rules []FirewallRule
	if err := json.Unmarshal(resp.Data, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Backup / Snapshots
// ─────────────────────────────────────────────────────────────────────────────

// GetBackups returns vzdump backup jobs.
func (c *Client) GetBackups() ([]map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get("/cluster/backup", &resp); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetVMSnapshots returns snapshots for a VM.
func (c *Client) GetVMSnapshots(node string, vmid int) ([]map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get(fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", node, vmid), &resp); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreateSnapshot creates a VM snapshot.
func (c *Client) CreateSnapshot(node string, vmid int, name, description string) (string, error) {
	params := url.Values{}
	params.Set("snapname", name)
	if description != "" {
		params.Set("description", description)
	}
	var resp apiListResponse
	if err := c.post(fmt.Sprintf("/nodes/%s/qemu/%d/snapshot", node, vmid), params, &resp); err != nil {
		return "", err
	}
	return string(resp.Data), nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Tasks
// ─────────────────────────────────────────────────────────────────────────────

// GetNodeTasks returns recent tasks for a node.
func (c *Client) GetNodeTasks(node string) ([]map[string]interface{}, error) {
	var resp apiListResponse
	if err := c.get("/nodes/"+node+"/tasks", &resp); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetTaskLog returns the log output of a task (UPID).
func (c *Client) GetTaskLog(node, upid string) ([]map[string]interface{}, error) {
	var resp apiListResponse
	escapedUPID := url.PathEscape(upid)
	if err := c.get(fmt.Sprintf("/nodes/%s/tasks/%s/log", node, escapedUPID), &resp); err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
