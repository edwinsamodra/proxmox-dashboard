package proxmox

// Node represents a Proxmox cluster node.
type Node struct {
	Node           string  `json:"node"`
	Status         string  `json:"status"`
	CPU            float64 `json:"cpu"`
	MaxCPU         int     `json:"maxcpu"`
	Mem            int64   `json:"mem"`
	MaxMem         int64   `json:"maxmem"`
	Disk           int64   `json:"disk"`
	MaxDisk        int64   `json:"maxdisk"`
	Uptime         int64   `json:"uptime"`
	Level          string  `json:"level"`
	Type           string  `json:"type"`
	SSLFingerprint string  `json:"ssl_fingerprint,omitempty"`
}

// VM represents a QEMU virtual machine.
type VM struct {
	VMID      int     `json:"vmid"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	Node      string  `json:"node,omitempty"`
	CPU       float64 `json:"cpu"`
	CPUs      int     `json:"cpus"`
	Mem       int64   `json:"mem"`
	MaxMem    int64   `json:"maxmem"`
	Disk      int64   `json:"disk"`
	MaxDisk   int64   `json:"maxdisk"`
	NetIn     int64   `json:"netin"`
	NetOut    int64   `json:"netout"`
	DiskRead  int64   `json:"diskread"`
	DiskWrite int64   `json:"diskwrite"`
	Uptime    int64   `json:"uptime"`
	PID       int     `json:"pid,omitempty"`
	Template  int     `json:"template,omitempty"`
	Tags      string  `json:"tags,omitempty"`
}

// Container represents an LXC container.
type Container struct {
	VMID      int     `json:"vmid"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	Node      string  `json:"node,omitempty"`
	CPU       float64 `json:"cpu"`
	CPUs      int     `json:"cpus"`
	Mem       int64   `json:"mem"`
	MaxMem    int64   `json:"maxmem"`
	Swap      int64   `json:"swap"`
	MaxSwap   int64   `json:"maxswap"`
	Disk      int64   `json:"disk"`
	MaxDisk   int64   `json:"maxdisk"`
	NetIn     int64   `json:"netin"`
	NetOut    int64   `json:"netout"`
	Uptime    int64   `json:"uptime"`
	Tags      string  `json:"tags,omitempty"`
}

// Storage represents a Proxmox storage.
type Storage struct {
	Storage  string   `json:"storage"`
	Type     string   `json:"type"`
	Status   string   `json:"status"`
	Active   int      `json:"active"`
	Used     int64    `json:"used"`
	Total    int64    `json:"total"`
	Avail    int64    `json:"avail"`
	Content  string   `json:"content"`
	Shared   int      `json:"shared"`
	Enabled  int      `json:"enabled"`
	Path     string   `json:"path,omitempty"`
	NodeList []string `json:"nodes,omitempty"`
}

// NetworkInterface represents a network interface on a node.
type NetworkInterface struct {
	Iface   string `json:"iface"`
	Type    string `json:"type"`
	Active  int    `json:"active"`
	Method  string `json:"method,omitempty"`
	Address string `json:"address,omitempty"`
	Netmask string `json:"netmask,omitempty"`
	Gateway string `json:"gateway,omitempty"`
	Bridge  string `json:"bridge_ports,omitempty"`
	CIDR    string `json:"cidr,omitempty"`
	MTU     int    `json:"mtu,omitempty"`
	Comment string `json:"comments,omitempty"`
}

// FirewallRule represents a firewall rule.
type FirewallRule struct {
	Pos     int    `json:"pos"`
	Type    string `json:"type"`
	Action  string `json:"action"`
	Proto   string `json:"proto,omitempty"`
	Source  string `json:"source,omitempty"`
	Dest    string `json:"dest,omitempty"`
	DPort   string `json:"dport,omitempty"`
	Sport   string `json:"sport,omitempty"`
	Enable  int    `json:"enable"`
	Comment string `json:"comment,omitempty"`
}
