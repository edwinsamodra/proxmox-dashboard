// Proxmox API types matching backend Go structs

export interface Node {
  node: string
  status: 'online' | 'offline' | string
  cpu: number
  maxcpu: number
  mem: number
  maxmem: number
  disk: number
  maxdisk: number
  uptime: number
  level: string
  type: string
}

export interface VM {
  vmid: number
  name: string
  status: 'running' | 'stopped' | 'paused' | 'suspended' | string
  node?: string
  cpu: number
  cpus: number
  mem: number
  maxmem: number
  disk: number
  maxdisk: number
  netin: number
  netout: number
  diskread: number
  diskwrite: number
  uptime: number
  template?: number
  tags?: string
}

export interface Container {
  vmid: number
  name: string
  status: 'running' | 'stopped' | string
  node?: string
  cpu: number
  cpus: number
  mem: number
  maxmem: number
  swap: number
  maxswap: number
  disk: number
  maxdisk: number
  netin: number
  netout: number
  uptime: number
  tags?: string
}

export interface Storage {
  storage: string
  type: string
  status: string
  active: number
  used: number
  total: number
  avail: number
  content: string
  shared: number
  enabled: number
  path?: string
}

export interface NetworkInterface {
  iface: string
  type: string
  active: number
  method?: string
  address?: string
  netmask?: string
  gateway?: string
  bridge_ports?: string
  cidr?: string
  mtu?: number
  comments?: string
}

export interface FirewallRule {
  pos: number
  type: string
  action: string
  proto?: string
  source?: string
  dest?: string
  dport?: string
  sport?: string
  enable: number
  comment?: string
}

export interface Snapshot {
  name: string
  snaptime?: number
  description?: string
  parent?: string
  vmstate?: number
}

export interface BackupJob {
  id?: string
  comment?: string
  compress?: string
  dow?: string
  enabled?: number
  mailnotification?: string
  mode?: string
  node?: string
  starttime?: string
  storage?: string
  vmid?: string
}

export interface Task {
  upid: string
  node: string
  type: string
  status: string
  user: string
  starttime: number
  endtime?: number
  id?: string
  pstart?: number
}

// OpenTofu types

export interface TofuWorkspace {
  id: string
  name: string
  created_at: string
  plan?: string
}

export interface VMProvisionRequest {
  workspace_name: string
  node_name: string
  vmid: number
  vm_name: string
  template: string
  cores: number
  memory_mb: number
  disk_size_gb: number
  storage_pool: string
  network: string
  ip_config: string
  ssh_keys: string
  cloud_init_user: string
  tags: string[]
  extra_vars?: Record<string, string>
}

export interface ContainerProvisionRequest {
  workspace_name: string
  node_name: string
  vmid: number
  hostname: string
  template: string
  cores: number
  memory_mb: number
  disk_size_gb: number
  storage_pool: string
  network: string
  ip_config: string
  password?: string
  ssh_keys: string
  tags: string[]
  extra_vars?: Record<string, string>
}

export interface TofuExecResult {
  success: boolean
  output: string
  error?: string
}

// WebSocket message types

export type WSMessageType =
  | 'nodes_update'
  | 'vms_update'
  | 'containers_update'
  | 'vm_action'
  | 'container_action'
  | 'tofu_apply'
  | 'vm_provisioned'

export interface WSMessage<T = unknown> {
  type: WSMessageType
  payload: T
}

// UI state types

export type VMAction = 'start' | 'stop' | 'reboot' | 'shutdown' | 'suspend' | 'resume' | 'reset'

export interface ClusterStatus {
  name?: string
  nodes: number
  online: number
  quorate?: number
  version?: number
}

export interface DashboardStats {
  totalNodes: number
  onlineNodes: number
  totalVMs: number
  runningVMs: number
  totalContainers: number
  runningContainers: number
  totalStorageBytes: number
  usedStorageBytes: number
  cpuUsage: number
  memUsage: number
}
