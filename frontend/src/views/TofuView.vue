<template>
  <div class="space-y-6 animate-fade-in">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">OpenTofu Infrastructure as Code</h2>
        <p class="text-sm text-gray-500 mt-0.5">Zero-Touch Provisioning — define infrastructure, not clicks</p>
      </div>
      <div class="flex items-center gap-3">
        <div
          class="flex items-center gap-2 text-xs px-3 py-1.5 rounded-lg border"
          :class="tofuStore.available
            ? 'bg-success/10 border-success/30 text-success'
            : 'bg-warning/10 border-warning/30 text-warning'"
        >
          <span class="w-1.5 h-1.5 rounded-full" :class="tofuStore.available ? 'bg-success animate-pulse' : 'bg-warning'" />
          {{ tofuStore.available ? 'OpenTofu available' : 'OpenTofu not found in PATH' }}
        </div>
        <button class="btn-primary" @click="showProvisionModal = true">
          <PlusIcon class="w-4 h-4" />
          New Provision
        </button>
      </div>
    </div>

    <!-- Two-pane layout: workspace list + editor -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 min-h-[600px]">
      <!-- Workspace list -->
      <div class="card flex flex-col">
        <div class="flex items-center justify-between px-4 py-3 border-b border-surface-300">
          <span class="text-sm font-semibold text-white">Workspaces</span>
          <button class="btn-ghost py-1 px-2 text-xs" @click="showNewWorkspaceModal = true">
            <PlusIcon class="w-3.5 h-3.5" />
            New
          </button>
        </div>
        <div class="flex-1 overflow-y-auto p-2 space-y-1">
          <div v-if="!tofuStore.workspaces.length" class="text-center text-gray-500 text-sm py-8">
            No workspaces yet
          </div>
          <button
            v-for="ws in tofuStore.workspaces"
            :key="ws.id"
            :class="[
              'w-full flex items-center justify-between px-3 py-2.5 rounded-lg text-sm transition-colors duration-100',
              tofuStore.currentWorkspace === ws.name
                ? 'bg-accent/15 text-accent'
                : 'text-gray-300 hover:bg-surface-200',
            ]"
            @click="tofuStore.selectWorkspace(ws.name)"
          >
            <div class="flex items-center gap-2">
              <FolderIcon class="w-4 h-4 opacity-60" />
              <span class="truncate font-medium">{{ ws.name }}</span>
            </div>
            <button
              class="p-1 rounded text-gray-500 hover:text-danger hover:bg-danger/10"
              @click.stop="deleteWorkspace(ws.name)"
            >
              <Trash2Icon class="w-3 h-3" />
            </button>
          </button>
        </div>
      </div>

      <!-- Editor + actions -->
      <div class="lg:col-span-2 card flex flex-col">
        <template v-if="tofuStore.currentWorkspace">
          <!-- Toolbar -->
          <div class="flex items-center justify-between px-4 py-3 border-b border-surface-300">
            <div class="flex items-center gap-2 text-sm font-medium text-white">
              <CodeIcon class="w-4 h-4 text-accent" />
              {{ tofuStore.currentWorkspace }} / main.tf
            </div>
            <div class="flex items-center gap-2">
              <button class="btn-ghost text-xs py-1.5 px-3" @click="tofuStore.saveCurrentTF()">
                <SaveIcon class="w-3.5 h-3.5" />
                Save
              </button>
              <button
                class="btn-secondary text-xs py-1.5 px-3"
                :disabled="tofuStore.executing"
                @click="runInit"
              >
                <TerminalIcon class="w-3.5 h-3.5" />
                Init
              </button>
              <button
                class="btn-secondary text-xs py-1.5 px-3"
                :disabled="tofuStore.executing"
                @click="runPlan"
              >
                <PlayIcon class="w-3.5 h-3.5" />
                Plan
              </button>
              <button
                class="btn-primary text-xs py-1.5 px-3"
                :disabled="tofuStore.executing"
                @click="runApply"
              >
                <ZapIcon class="w-3.5 h-3.5" />
                Apply
              </button>
              <button
                class="btn-danger text-xs py-1.5 px-3"
                :disabled="tofuStore.executing"
                @click="runDestroy"
              >
                <BombIcon class="w-3.5 h-3.5" />
                Destroy
              </button>
            </div>
          </div>

          <!-- TF editor (textarea-based) -->
          <div class="flex-1 flex flex-col min-h-0">
            <textarea
              v-model="tofuStore.currentTF"
              class="flex-1 font-mono text-xs bg-surface-DEFAULT text-green-400 p-4 resize-none
                     focus:outline-none focus:ring-1 focus:ring-accent/40 leading-relaxed
                     border-b border-surface-300"
              spellcheck="false"
              placeholder="# Write your Terraform/OpenTofu configuration here..."
              style="min-height: 300px;"
            />

            <!-- Output pane -->
            <div class="flex-shrink-0" style="max-height: 250px;">
              <div class="flex items-center justify-between px-4 py-2 border-t border-surface-300 bg-surface-200">
                <span class="text-xs font-semibold text-gray-400 uppercase tracking-wider">Output</span>
                <div v-if="tofuStore.executing" class="flex items-center gap-1.5 text-xs text-warning">
                  <LoaderIcon class="w-3 h-3 animate-spin" />
                  Executing...
                </div>
                <span v-else-if="tofuStore.lastResult" class="text-xs" :class="tofuStore.lastResult.success ? 'text-success' : 'text-danger'">
                  {{ tofuStore.lastResult.success ? '✓ Success' : '✗ Failed' }}
                </span>
              </div>
              <div
                class="overflow-y-auto bg-surface-DEFAULT font-mono text-xs px-4 py-3 text-gray-300"
                style="max-height: 200px;"
              >
                <pre v-if="tofuStore.lastResult" class="whitespace-pre-wrap break-words">{{ tofuStore.lastResult.output }}<span v-if="tofuStore.lastResult.error" class="text-danger">{{ tofuStore.lastResult.error }}</span></pre>
                <div v-else class="text-gray-600">Run init / plan / apply to see output here.</div>
              </div>
            </div>
          </div>
        </template>
        <template v-else>
          <div class="flex-1 flex items-center justify-center text-center p-12">
            <div>
              <CodeIcon class="w-12 h-12 text-gray-600 mx-auto mb-4" />
              <div class="text-gray-400 font-medium">Select or create a workspace</div>
              <p class="text-gray-600 text-sm mt-1">Use workspaces to organize your Terraform configurations</p>
              <button class="btn-primary mt-4" @click="showNewWorkspaceModal = true">
                <PlusIcon class="w-4 h-4" />
                Create Workspace
              </button>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- New Workspace Modal -->
    <Modal v-model="showNewWorkspaceModal" title="New Workspace" size="sm">
      <div class="space-y-4">
        <div>
          <label class="label">Workspace Name</label>
          <input v-model="newWorkspaceName" type="text" class="input" placeholder="my-vm-stack" />
        </div>
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showNewWorkspaceModal = false">Cancel</button>
        <button class="btn-primary" @click="createWorkspace">
          <PlusIcon class="w-4 h-4" />
          Create
        </button>
      </template>
    </Modal>

    <!-- Zero-Touch Provision Modal -->
    <Modal v-model="showProvisionModal" title="Zero-Touch VM Provisioning" size="xl">
      <div class="space-y-6">
        <div class="bg-accent/10 border border-accent/30 rounded-lg p-4 text-sm text-gray-300">
          <div class="font-semibold text-accent mb-1">🚀 Zero-Touch Provisioning</div>
          Fill in the form below and PMO will generate a Terraform plan, run <code class="font-mono bg-surface-300 px-1 rounded">tofu apply</code>,
          and provision your VM automatically — no SSH required.
        </div>

        <!-- Tab: VM / Container -->
        <div class="flex rounded-lg bg-surface-200 p-1 gap-1 w-fit">
          <button
            v-for="t in (['VM', 'Container'] as const)"
            :key="t"
            :class="['px-4 py-1.5 text-sm font-medium rounded-md transition-colors',
              provisionType === t ? 'bg-surface-400 text-white' : 'text-gray-400 hover:text-white']"
            @click="provisionType = t"
          >{{ t }}</button>
        </div>

        <div class="grid grid-cols-2 gap-4 text-sm">
          <!-- Common fields -->
          <div>
            <label class="label">Workspace Name</label>
            <input v-model="provision.workspace_name" class="input" placeholder="prod-web-01" />
          </div>
          <div>
            <label class="label">Node</label>
            <select v-model="provision.node_name" class="input">
              <option v-for="n in nodesStore.nodes" :key="n.node" :value="n.node">{{ n.node }}</option>
            </select>
          </div>
          <div>
            <label class="label">VMID</label>
            <input v-model.number="provision.vmid" type="number" class="input" placeholder="200" />
          </div>
          <div>
            <label class="label">{{ provisionType === 'VM' ? 'VM Name' : 'Hostname' }}</label>
            <input v-model="provision.name" class="input" placeholder="web-server-01" />
          </div>
          <div>
            <label class="label">Template / OS Template</label>
            <input v-model="provision.template" class="input" placeholder="ubuntu-22.04-cloud" />
          </div>
          <div>
            <label class="label">CPU Cores</label>
            <input v-model.number="provision.cores" type="number" class="input" placeholder="2" />
          </div>
          <div>
            <label class="label">Memory (MB)</label>
            <input v-model.number="provision.memory_mb" type="number" class="input" placeholder="2048" />
          </div>
          <div>
            <label class="label">Disk Size (GB)</label>
            <input v-model.number="provision.disk_size_gb" type="number" class="input" placeholder="20" />
          </div>
          <div>
            <label class="label">Storage Pool</label>
            <input v-model="provision.storage_pool" class="input" placeholder="local-lvm" />
          </div>
          <div>
            <label class="label">Network Bridge</label>
            <input v-model="provision.network" class="input" placeholder="vmbr0" />
          </div>
          <div>
            <label class="label">IP Config (Cloud-Init)</label>
            <input v-model="provision.ip_config" class="input" placeholder="ip=192.168.1.10/24,gw=192.168.1.1" />
          </div>
          <div v-if="provisionType === 'VM'">
            <label class="label">Cloud-Init User</label>
            <input v-model="provision.cloud_init_user" class="input" placeholder="ubuntu" />
          </div>
          <div class="col-span-2">
            <label class="label">SSH Public Keys</label>
            <textarea v-model="provision.ssh_keys" class="input font-mono text-xs" rows="3" placeholder="ssh-ed25519 AAAA... user@host" />
          </div>
        </div>

        <!-- Generated TF preview -->
        <div v-if="generatedTF" class="border border-surface-300 rounded-lg overflow-hidden">
          <div class="px-4 py-2 bg-surface-200 text-xs font-semibold text-gray-400 flex items-center justify-between">
            <span>Generated Terraform Configuration</span>
            <button class="text-accent hover:underline text-xs" @click="openInEditor">Open in Editor</button>
          </div>
          <pre class="p-4 text-xs font-mono text-green-400 bg-surface-DEFAULT overflow-auto max-h-48">{{ generatedTF }}</pre>
        </div>
      </div>

      <template #footer>
        <button class="btn-secondary" @click="handleGenerateTF">
          <CodeIcon class="w-4 h-4" />
          Preview TF
        </button>
        <button class="btn-secondary" @click="showProvisionModal = false">Cancel</button>
        <button
          class="btn-primary"
          :disabled="tofuStore.executing || !tofuStore.available"
          @click="handleProvision"
        >
          <ZapIcon class="w-4 h-4" :class="{ 'animate-pulse': tofuStore.executing }" />
          {{ tofuStore.executing ? 'Provisioning...' : 'Provision Now' }}
        </button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  Plus as PlusIcon,
  Code as CodeIcon,
  Folder as FolderIcon,
  Save as SaveIcon,
  Terminal as TerminalIcon,
  Play as PlayIcon,
  Zap as ZapIcon,
  Bomb as BombIcon,
  Trash2 as Trash2Icon,
  Loader as LoaderIcon,
} from 'lucide-vue-next'
import { useTofuStore } from '@/stores/tofu'
import { useNodesStore } from '@/stores/nodes'
import { tofuApi } from '@/api'
import Modal from '@/components/ui/Modal.vue'
import type { VMProvisionRequest } from '@/types'

const tofuStore = useTofuStore()
const nodesStore = useNodesStore()

const showNewWorkspaceModal = ref(false)
const showProvisionModal = ref(false)
const newWorkspaceName = ref('')
const provisionType = ref<'VM' | 'Container'>('VM')
const generatedTF = ref('')

const provision = ref({
  workspace_name: '',
  node_name: '',
  vmid: 200,
  name: '',
  template: '',
  cores: 2,
  memory_mb: 2048,
  disk_size_gb: 20,
  storage_pool: 'local-lvm',
  network: 'vmbr0',
  ip_config: 'ip=dhcp',
  cloud_init_user: 'ubuntu',
  ssh_keys: '',
})

async function createWorkspace() {
  if (!newWorkspaceName.value) return
  await tofuStore.createWorkspace(newWorkspaceName.value)
  showNewWorkspaceModal.value = false
  newWorkspaceName.value = ''
}

async function deleteWorkspace(name: string) {
  if (!confirm(`Delete workspace "${name}"?`)) return
  await tofuStore.deleteWorkspace(name)
}

async function runInit() {
  await tofuStore.runInit()
}

async function runPlan() {
  await tofuStore.runPlan()
}

async function runApply() {
  await tofuStore.runApply()
}

async function runDestroy() {
  if (!confirm('Run tofu destroy? This will remove all managed resources.')) return
  await tofuStore.runDestroy()
}

async function handleGenerateTF() {
  const req: VMProvisionRequest = {
    workspace_name: provision.value.workspace_name || 'preview',
    node_name: provision.value.node_name,
    vmid: provision.value.vmid,
    vm_name: provision.value.name,
    template: provision.value.template,
    cores: provision.value.cores,
    memory_mb: provision.value.memory_mb,
    disk_size_gb: provision.value.disk_size_gb,
    storage_pool: provision.value.storage_pool,
    network: provision.value.network,
    ip_config: provision.value.ip_config,
    ssh_keys: provision.value.ssh_keys,
    cloud_init_user: provision.value.cloud_init_user,
    tags: [],
  }
  const result = await tofuApi.generateVMTF(req)
  generatedTF.value = result.content
}

async function openInEditor() {
  if (!provision.value.workspace_name) return
  await tofuStore.createWorkspace(provision.value.workspace_name, generatedTF.value)
  await tofuStore.selectWorkspace(provision.value.workspace_name)
  showProvisionModal.value = false
}

async function handleProvision() {
  if (!provision.value.workspace_name) {
    alert('Please enter a workspace name')
    return
  }
  const req: VMProvisionRequest = {
    workspace_name: provision.value.workspace_name,
    node_name: provision.value.node_name,
    vmid: provision.value.vmid,
    vm_name: provision.value.name,
    template: provision.value.template,
    cores: provision.value.cores,
    memory_mb: provision.value.memory_mb,
    disk_size_gb: provision.value.disk_size_gb,
    storage_pool: provision.value.storage_pool,
    network: provision.value.network,
    ip_config: provision.value.ip_config,
    ssh_keys: provision.value.ssh_keys,
    cloud_init_user: provision.value.cloud_init_user,
    tags: [],
  }
  const result = await tofuStore.generateAndProvisionVM(req)
  if (result?.success) {
    showProvisionModal.value = false
    alert('VM provisioned successfully!')
  }
}

onMounted(async () => {
  await Promise.all([
    tofuStore.checkStatus(),
    tofuStore.fetchWorkspaces(),
    nodesStore.fetchNodes(),
  ])
  if (nodesStore.nodes.length && !provision.value.node_name) {
    provision.value.node_name = nodesStore.nodes[0].node
  }
})
</script>
