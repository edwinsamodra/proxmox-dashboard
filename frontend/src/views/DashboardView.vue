<template>
  <div class="space-y-6 animate-fade-in">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-semibold text-white">Cluster Overview</h2>
        <p class="text-sm text-gray-500 mt-0.5">Real-time infrastructure status</p>
      </div>
      <div class="flex items-center gap-2 text-xs text-gray-500">
        <span class="w-2 h-2 bg-success rounded-full animate-pulse" />
        Live
      </div>
    </div>

    <!-- Stats grid -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Nodes -->
      <div class="stat-card">
        <div class="flex items-center justify-between">
          <div class="stat-label">Nodes</div>
          <ServerIcon class="w-4 h-4 text-accent opacity-60" />
        </div>
        <template v-if="nodesStore.loading">
          <SkeletonLoader height="2.5rem" width="60%" />
        </template>
        <template v-else>
          <div class="stat-value">{{ nodesStore.onlineNodes.length }}<span class="text-base text-gray-500">/{{ nodesStore.nodes.length }}</span></div>
          <div class="text-xs text-gray-500">{{ nodesStore.nodes.length - nodesStore.onlineNodes.length }} offline</div>
        </template>
      </div>

      <!-- VMs -->
      <div class="stat-card">
        <div class="flex items-center justify-between">
          <div class="stat-label">Virtual Machines</div>
          <CpuIcon class="w-4 h-4 text-accent opacity-60" />
        </div>
        <template v-if="vmsStore.loading">
          <SkeletonLoader height="2.5rem" width="60%" />
        </template>
        <template v-else>
          <div class="stat-value">{{ vmsStore.runningVMs.length }}<span class="text-base text-gray-500">/{{ vmsStore.vms.length }}</span></div>
          <div class="text-xs text-success">{{ vmsStore.runningVMs.length }} running</div>
        </template>
      </div>

      <!-- Containers -->
      <div class="stat-card">
        <div class="flex items-center justify-between">
          <div class="stat-label">Containers</div>
          <BoxIcon class="w-4 h-4 text-accent opacity-60" />
        </div>
        <template v-if="containersStore.loading">
          <SkeletonLoader height="2.5rem" width="60%" />
        </template>
        <template v-else>
          <div class="stat-value">{{ containersStore.runningContainers.length }}<span class="text-base text-gray-500">/{{ containersStore.containers.length }}</span></div>
          <div class="text-xs text-success">{{ containersStore.runningContainers.length }} running</div>
        </template>
      </div>

      <!-- CPU Usage -->
      <div class="stat-card">
        <div class="flex items-center justify-between">
          <div class="stat-label">Avg CPU</div>
          <ActivityIcon class="w-4 h-4 text-accent opacity-60" />
        </div>
        <template v-if="nodesStore.loading">
          <SkeletonLoader height="2.5rem" width="60%" />
        </template>
        <template v-else>
          <div class="stat-value">{{ avgCpuPercent }}<span class="text-base text-gray-500">%</span></div>
          <UsageBar :value="avgCpuPercent" class="mt-1" />
        </template>
      </div>
    </div>

    <!-- Memory + storage row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- Memory usage -->
      <div class="card p-5">
        <div class="flex items-center justify-between mb-4">
          <div class="text-sm font-semibold text-white">Memory Usage</div>
          <span class="text-xs text-gray-500">{{ formatBytes(nodesStore.usedMem) }} / {{ formatBytes(nodesStore.totalMem) }}</span>
        </div>
        <template v-if="nodesStore.loading">
          <SkeletonLoader height="1rem" />
        </template>
        <template v-else>
          <UsageBar :value="nodesStore.usedMem" :max="nodesStore.totalMem || 1" />
          <div class="mt-3 grid grid-cols-3 gap-3">
            <div v-for="node in nodesStore.nodes" :key="node.node" class="text-xs">
              <div class="text-gray-400 truncate mb-1">{{ node.node }}</div>
              <UsageBar :value="node.mem" :max="node.maxmem || 1" />
              <div class="text-gray-500 mt-0.5">{{ Math.round((node.mem / (node.maxmem || 1)) * 100) }}%</div>
            </div>
          </div>
        </template>
      </div>

      <!-- Node status cards -->
      <div class="card p-5">
        <div class="text-sm font-semibold text-white mb-4">Node Status</div>
        <template v-if="nodesStore.loading">
          <div class="space-y-2">
            <SkeletonLoader v-for="i in 3" :key="i" height="3rem" />
          </div>
        </template>
        <template v-else>
          <div class="space-y-2.5">
            <div
              v-for="node in nodesStore.nodes"
              :key="node.node"
              class="flex items-center justify-between rounded-lg bg-surface-200 px-4 py-3"
            >
              <div class="flex items-center gap-3">
                <StatusBadge :status="node.status" dot />
                <div>
                  <div class="text-sm font-medium text-white">{{ node.node }}</div>
                  <div class="text-xs text-gray-500">{{ node.maxcpu }} vCPU · {{ formatBytes(node.maxmem) }}</div>
                </div>
              </div>
              <div class="text-right text-xs">
                <div class="text-white font-medium">{{ Math.round(node.cpu * 100) }}% CPU</div>
                <div class="text-gray-500">{{ formatUptime(node.uptime) }}</div>
              </div>
            </div>
            <div v-if="!nodesStore.nodes.length" class="text-center text-gray-500 text-sm py-4">
              No nodes found — check Proxmox connection
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- Recent VMs table -->
    <div class="table-container">
      <div class="flex items-center justify-between p-4 border-b border-surface-300">
        <div class="text-sm font-semibold text-white">Recent Virtual Machines</div>
        <RouterLink to="/vms" class="text-xs text-accent hover:underline">View all</RouterLink>
      </div>
      <table class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Node</th>
            <th>Status</th>
            <th>CPU</th>
            <th>Memory</th>
            <th>Uptime</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="vmsStore.loading">
            <tr v-for="i in 5" :key="i">
              <td colspan="7"><SkeletonLoader height="1.25rem" /></td>
            </tr>
          </template>
          <template v-else>
            <tr v-for="vm in recentVMs" :key="`${vm.node}-${vm.vmid}`">
              <td class="font-mono text-gray-400 text-xs">{{ vm.vmid }}</td>
              <td class="font-medium text-white">{{ vm.name }}</td>
              <td class="text-gray-400 text-xs">{{ vm.node }}</td>
              <td><StatusBadge :status="vm.status" dot /></td>
              <td>
                <div class="flex items-center gap-2">
                  <span class="text-xs text-gray-300">{{ Math.round(vm.cpu * 100) }}%</span>
                  <UsageBar :value="vm.cpu * 100" class="w-16" />
                </div>
              </td>
              <td>
                <div class="flex items-center gap-2">
                  <span class="text-xs text-gray-300">{{ Math.round((vm.mem / (vm.maxmem || 1)) * 100) }}%</span>
                  <UsageBar :value="vm.mem" :max="vm.maxmem || 1" class="w-16" />
                </div>
              </td>
              <td class="text-xs text-gray-400">{{ formatUptime(vm.uptime) }}</td>
            </tr>
            <tr v-if="!recentVMs.length">
              <td colspan="7" class="text-center text-gray-500 py-6">No VMs found</td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import {
  Server as ServerIcon,
  Cpu as CpuIcon,
  Box as BoxIcon,
  Activity as ActivityIcon,
} from 'lucide-vue-next'
import { useNodesStore } from '@/stores/nodes'
import { useVMsStore } from '@/stores/vms'
import { useContainersStore } from '@/stores/containers'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import UsageBar from '@/components/ui/UsageBar.vue'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'

const nodesStore = useNodesStore()
const vmsStore = useVMsStore()
const containersStore = useContainersStore()

const avgCpuPercent = computed(() => {
  if (!nodesStore.nodes.length) return 0
  const total = nodesStore.nodes.reduce((s, n) => s + n.cpu, 0)
  return Math.round((total / nodesStore.nodes.length) * 100)
})

const recentVMs = computed(() =>
  vmsStore.vms
    .slice()
    .sort((a, b) => (b.status === 'running' ? 1 : 0) - (a.status === 'running' ? 1 : 0))
    .slice(0, 8)
)

function formatBytes(bytes: number): string {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  let val = bytes
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
  await Promise.all([
    nodesStore.fetchNodes(),
    vmsStore.fetchAll(),
    containersStore.fetchAll(),
  ])
})
</script>
