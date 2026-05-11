import type {
  Node, VM, Container, Storage, NetworkInterface,
  FirewallRule, Snapshot, BackupJob, Task,
  TofuWorkspace, TofuExecResult,
  VMProvisionRequest, ContainerProvisionRequest,
  VMAction
} from '@/types'

const BASE = '/api'

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

async function request<T>(
  method: string,
  path: string,
  body?: unknown
): Promise<T> {
  const opts: RequestInit = {
    method,
    headers: { 'Content-Type': 'application/json' },
  }
  if (body !== undefined) {
    opts.body = JSON.stringify(body)
  }
  const res = await fetch(BASE + path, opts)
  if (!res.ok) {
    const data = await res.json().catch(() => ({ error: res.statusText }))
    throw new ApiError(res.status, data.error ?? res.statusText)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

const get = <T>(path: string) => request<T>('GET', path)
const post = <T>(path: string, body?: unknown) => request<T>('POST', path, body)
const put = <T>(path: string, body?: unknown) => request<T>('PUT', path, body)
const del = <T>(path: string) => request<T>('DELETE', path)

// ─── Auth ────────────────────────────────────────────────────────────────────

export const authApi = {
  login: (username: string, password: string, realm = 'pam') =>
    post<{ authenticated: boolean }>('/auth/login', { username, password, realm }),
  tokenAuth: (tokenID: string, tokenSecret: string) => {
    return fetch(BASE + '/auth/token', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Proxmox-Token-ID': tokenID,
        'X-Proxmox-Token-Secret': tokenSecret,
      },
    }).then(r => r.json())
  },
}

// ─── Nodes ───────────────────────────────────────────────────────────────────

export const nodesApi = {
  list: () => get<Node[]>('/nodes'),
  status: (node: string) => get<Record<string, unknown>>(`/nodes/${node}/status`),
  clusterStatus: () => get<Record<string, unknown>[]>('/cluster/status'),
  network: (node: string) => get<NetworkInterface[]>(`/nodes/${node}/network`),
  tasks: (node: string) => get<Task[]>(`/nodes/${node}/tasks`),
  taskLog: (node: string, upid: string) =>
    get<Record<string, unknown>[]>(`/nodes/${node}/tasks/log?upid=${encodeURIComponent(upid)}`),
}

// ─── VMs ─────────────────────────────────────────────────────────────────────

export const vmsApi = {
  listAll: () => get<VM[]>('/vms'),
  listForNode: (node: string) => get<VM[]>(`/nodes/${node}/vms`),
  status: (node: string, vmid: number) =>
    get<Record<string, unknown>>(`/nodes/${node}/vms/${vmid}/status`),
  config: (node: string, vmid: number) =>
    get<Record<string, unknown>>(`/nodes/${node}/vms/${vmid}/config`),
  action: (node: string, vmid: number, action: VMAction) =>
    post<{ task: string }>(`/nodes/${node}/vms/${vmid}/action/${action}`),
  create: (node: string, params: Record<string, unknown>) =>
    post<{ task: string }>(`/nodes/${node}/vms`, params),
  delete: (node: string, vmid: number) =>
    del<{ task: string }>(`/nodes/${node}/vms/${vmid}`),
  snapshots: (node: string, vmid: number) =>
    get<Snapshot[]>(`/nodes/${node}/vms/${vmid}/snapshots`),
  createSnapshot: (node: string, vmid: number, name: string, description?: string) =>
    post<{ task: string }>(`/nodes/${node}/vms/${vmid}/snapshots`, { name, description }),
  firewall: (node: string, vmid: number) =>
    get<FirewallRule[]>(`/nodes/${node}/vms/${vmid}/firewall`),
}

// ─── Containers ──────────────────────────────────────────────────────────────

export const containersApi = {
  listAll: () => get<Container[]>('/containers'),
  listForNode: (node: string) => get<Container[]>(`/nodes/${node}/containers`),
  action: (node: string, vmid: number, action: VMAction) =>
    post<{ task: string }>(`/nodes/${node}/containers/${vmid}/action/${action}`),
  delete: (node: string, vmid: number) =>
    del<{ task: string }>(`/nodes/${node}/containers/${vmid}`),
}

// ─── Storage ─────────────────────────────────────────────────────────────────

export const storageApi = {
  list: (node: string) => get<Storage[]>(`/nodes/${node}/storage`),
  content: (node: string, storage: string) =>
    get<Record<string, unknown>[]>(`/nodes/${node}/storage/${storage}/content`),
}

// ─── Firewall ────────────────────────────────────────────────────────────────

export const firewallApi = {
  clusterRules: () => get<FirewallRule[]>('/cluster/firewall'),
  vmRules: (node: string, vmid: number) =>
    get<FirewallRule[]>(`/nodes/${node}/vms/${vmid}/firewall`),
}

// ─── Backup ──────────────────────────────────────────────────────────────────

export const backupApi = {
  list: () => get<BackupJob[]>('/backups'),
}

// ─── OpenTofu ────────────────────────────────────────────────────────────────

export const tofuApi = {
  status: () => get<{ available: boolean; workspaces: TofuWorkspace[] }>('/tofu/status'),
  listWorkspaces: () => get<TofuWorkspace[]>('/tofu/workspaces'),
  createWorkspace: (name: string, content = '') =>
    post<{ workspace: string }>('/tofu/workspaces', { name, content }),
  getWorkspaceTF: (name: string) =>
    get<{ content: string }>(`/tofu/workspaces/${name}/tf`),
  saveWorkspaceTF: (name: string, content: string) =>
    put<{ saved: boolean }>(`/tofu/workspaces/${name}/tf`, { content }),
  deleteWorkspace: (name: string) =>
    del<{ deleted: boolean }>(`/tofu/workspaces/${name}`),
  init: (name: string) =>
    post<TofuExecResult>(`/tofu/workspaces/${name}/init`),
  plan: (name: string, vars?: Record<string, string>) =>
    post<TofuExecResult>(`/tofu/workspaces/${name}/plan`, vars),
  apply: (name: string, vars?: Record<string, string>) =>
    post<TofuExecResult>(`/tofu/workspaces/${name}/apply`, vars),
  destroy: (name: string, vars?: Record<string, string>) =>
    post<TofuExecResult>(`/tofu/workspaces/${name}/destroy`, vars),
  output: (name: string) =>
    get<Record<string, unknown>>(`/tofu/workspaces/${name}/output`),
  generateVMTF: (req: VMProvisionRequest) =>
    post<{ content: string }>('/tofu/generate/vm', req),
  generateContainerTF: (req: ContainerProvisionRequest) =>
    post<{ content: string }>('/tofu/generate/container', req),
  provisionVM: (req: VMProvisionRequest) =>
    post<TofuExecResult>('/tofu/provision/vm', req),
}

export { ApiError }
