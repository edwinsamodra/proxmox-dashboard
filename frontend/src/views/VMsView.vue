<template>
  <div class="space-y-6 animate-fade-in">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-white">Virtual Machines</h2>
      <button class="btn-primary" @click="showCreateModal = true">
        <PlusIcon class="w-4 h-4" />
        New VM
      </button>
    </div>

    <!-- Filters -->
    <div class="flex items-center gap-3">
      <input v-model="search" type="search" class="input max-w-xs" placeholder="Search VMs..." />
      <select v-model="filterStatus" class="input max-w-[140px]">
        <option value="">All status</option>
        <option value="running">Running</option>
        <option value="stopped">Stopped</option>
      </select>
      <div class="ml-auto text-xs text-gray-500">
        {{ filteredVMs.length }} VMs
      </div>
    </div>

    <!-- VMs table -->
    <div class="table-container">
      <table class="table">
        <thead>
          <tr>
            <th>VMID</th>
            <th>Name</th>
            <th>Node</th>
            <th>Status</th>
            <th>CPU</th>
            <th>Memory</th>
            <th>Disk</th>
            <th>Uptime</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="vmsStore.loading">
            <tr v-for="i in 8" :key="i">
              <td colspan="9"><SkeletonLoader height="1.25rem" /></td>
            </tr>
          </template>
          <template v-else>
            <tr v-for="vm in filteredVMs" :key="`${vm.node}-${vm.vmid}`">
              <td class="font-mono text-gray-400 text-xs">{{ vm.vmid }}</td>
              <td>
                <div class="font-medium text-white">{{ vm.name }}</div>
                <div v-if="vm.tags" class="text-xs text-gray-500 mt-0.5">{{ vm.tags }}</div>
              </td>
              <td class="text-xs text-gray-400">{{ vm.node }}</td>
              <td><StatusBadge :status="vm.status" dot /></td>
              <td>
                <div class="flex items-center gap-2">
                  <span class="text-xs w-8 text-right text-gray-300">{{ Math.round(vm.cpu * 100) }}%</span>
                  <UsageBar :value="vm.cpu * 100" class="w-20" />
                </div>
              </td>
              <td>
                <div class="flex items-center gap-2">
                  <span class="text-xs w-8 text-right text-gray-300">{{ Math.round((vm.mem / (vm.maxmem || 1)) * 100) }}%</span>
                  <UsageBar :value="vm.mem" :max="vm.maxmem || 1" class="w-20" />
                </div>
              </td>
              <td class="text-xs text-gray-400">{{ formatBytes(vm.maxdisk) }}</td>
              <td class="text-xs text-gray-400">{{ formatUptime(vm.uptime) }}</td>
              <td>
                <div class="flex items-center gap-1">
                  <button
                    v-if="vm.status !== 'running'"
                    class="btn-success py-1 px-2 text-xs"
                    :disabled="vmsStore.isActionLoading(vm.node!, vm.vmid)"
                    @click="performAction(vm, 'start')"
                    title="Start"
                  >
                    <PlayIcon class="w-3 h-3" />
                  </button>
                  <button
                    v-if="vm.status === 'running'"
                    class="btn-secondary py-1 px-2 text-xs"
                    :disabled="vmsStore.isActionLoading(vm.node!, vm.vmid)"
                    @click="performAction(vm, 'reboot')"
                    title="Reboot"
                  >
                    <RefreshCwIcon class="w-3 h-3" />
                  </button>
                  <button
                    v-if="vm.status === 'running'"
                    class="btn-danger py-1 px-2 text-xs"
                    :disabled="vmsStore.isActionLoading(vm.node!, vm.vmid)"
                    @click="performAction(vm, 'stop')"
                    title="Stop"
                  >
                    <SquareIcon class="w-3 h-3" />
                  </button>
                  <button
                    class="btn-ghost py-1 px-2 text-xs"
                    @click="selectVM(vm)"
                    title="Details"
                  >
                    <MoreHorizontalIcon class="w-3 h-3" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!filteredVMs.length">
              <td colspan="9" class="text-center text-gray-500 py-8">
                {{ vmsStore.loading ? 'Loading...' : 'No virtual machines found' }}
              </td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>

    <!-- VM Detail Modal -->
    <Modal v-model="showDetailModal" :title="`VM ${selectedVM?.vmid} — ${selectedVM?.name}`" size="lg">
      <template v-if="selectedVM">
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Status</div>
            <StatusBadge :status="selectedVM.status" dot />
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Node</div>
            <div class="text-white font-medium">{{ selectedVM.node }}</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">CPU</div>
            <div class="text-white font-medium">{{ selectedVM.cpus }} cores · {{ Math.round(selectedVM.cpu * 100) }}% used</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Memory</div>
            <div class="text-white font-medium">{{ formatBytes(selectedVM.mem) }} / {{ formatBytes(selectedVM.maxmem) }}</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Network In</div>
            <div class="text-white font-medium">{{ formatBytes(selectedVM.netin) }}</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Network Out</div>
            <div class="text-white font-medium">{{ formatBytes(selectedVM.netout) }}</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Disk Read</div>
            <div class="text-white font-medium">{{ formatBytes(selectedVM.diskread) }}</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Disk Write</div>
            <div class="text-white font-medium">{{ formatBytes(selectedVM.diskwrite) }}</div>
          </div>
        </div>
      </template>
      <template #footer>
        <button class="btn-danger" @click="handleDelete">
          <Trash2Icon class="w-4 h-4" />
          Delete VM
        </button>
        <button class="btn-secondary" @click="showDetailModal = false">Close</button>
      </template>
    </Modal>

    <!-- Create VM Modal (basic) -->
    <Modal v-model="showCreateModal" title="Create Virtual Machine" size="lg">
      <div class="space-y-4">
        <div class="bg-accent/10 border border-accent/30 rounded-lg p-4 text-sm text-accent">
          💡 For advanced provisioning with Cloud-init, network config and SSH keys, use the
          <RouterLink to="/tofu" class="underline font-medium" @click="showCreateModal = false">OpenTofu IaC</RouterLink> module.
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label">Node</label>
            <select v-model="createForm.node" class="input">
              <option v-for="n in nodesStore.nodes" :key="n.node" :value="n.node">{{ n.node }}</option>
            </select>
          </div>
          <div>
            <label class="label">VM ID</label>
            <input v-model.number="createForm.vmid" type="number" class="input" placeholder="100" />
          </div>
          <div>
            <label class="label">Name</label>
            <input v-model="createForm.name" type="text" class="input" placeholder="my-vm" />
          </div>
          <div>
            <label class="label">Memory (MB)</label>
            <input v-model.number="createForm.memory" type="number" class="input" placeholder="2048" />
          </div>
          <div>
            <label class="label">CPU Cores</label>
            <input v-model.number="createForm.cores" type="number" class="input" placeholder="2" />
          </div>
          <div>
            <label class="label">OS Type</label>
            <select v-model="createForm.ostype" class="input">
              <option value="l26">Linux 2.6+</option>
              <option value="win10">Windows 10</option>
              <option value="win11">Windows 11</option>
            </select>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showCreateModal = false">Cancel</button>
        <button class="btn-primary" @click="handleCreate">
          <PlusIcon class="w-4 h-4" />
          Create VM
        </button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Plus as PlusIcon, Play as PlayIcon, Square as SquareIcon,
  RefreshCw as RefreshCwIcon, MoreHorizontal as MoreHorizontalIcon,
  Trash2 as Trash2Icon,
} from 'lucide-vue-next'
import { useVMsStore } from '@/stores/vms'
import { useNodesStore } from '@/stores/nodes'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import UsageBar from '@/components/ui/UsageBar.vue'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'
import Modal from '@/components/ui/Modal.vue'
import type { VM, VMAction } from '@/types'

const vmsStore = useVMsStore()
const nodesStore = useNodesStore()

const search = ref('')
const filterStatus = ref('')
const showDetailModal = ref(false)
const showCreateModal = ref(false)
const selectedVM = ref<VM | null>(null)

const createForm = ref({
  node: '',
  vmid: 100,
  name: '',
  memory: 2048,
  cores: 2,
  ostype: 'l26',
})

const filteredVMs = computed(() =>
  vmsStore.vms.filter(vm => {
    const matchSearch = !search.value ||
      vm.name?.toLowerCase().includes(search.value.toLowerCase()) ||
      String(vm.vmid).includes(search.value)
    const matchStatus = !filterStatus.value || vm.status === filterStatus.value
    return matchSearch && matchStatus
  })
)

async function performAction(vm: VM, action: VMAction) {
  try {
    await vmsStore.performAction(vm.node!, vm.vmid, action)
  } catch {
    // error handled in store
  }
}

function selectVM(vm: VM) {
  selectedVM.value = vm
  showDetailModal.value = true
}

async function handleDelete() {
  if (!selectedVM.value) return
  if (!confirm(`Delete VM ${selectedVM.value.name}?`)) return
  await vmsStore.deleteVM(selectedVM.value.node!, selectedVM.value.vmid)
  showDetailModal.value = false
}

async function handleCreate() {
  // Basic creation — for full IaC use OpenTofu
  showCreateModal.value = false
}

function formatBytes(bytes: number): string {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0; let val = bytes
  while (val >= 1024 && i < units.length - 1) { val /= 1024; i++ }
  return `${val.toFixed(i > 1 ? 1 : 0)} ${units[i]}`
}

function formatUptime(seconds: number): string {
  if (!seconds) return '—'
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  if (d > 0) return `${d}d ${h}h`
  const m = Math.floor((seconds % 3600) / 60)
  if (h > 0) return `${h}h ${m}m`
  return `${m}m`
}

onMounted(async () => {
  await Promise.all([nodesStore.fetchNodes(), vmsStore.fetchAll()])
  if (!createForm.value.node && nodesStore.nodes.length) {
    createForm.value.node = nodesStore.nodes[0].node
  }
})
</script>
