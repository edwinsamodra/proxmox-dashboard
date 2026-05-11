<template>
  <div class="space-y-6 animate-fade-in">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-white">Containers (LXC)</h2>
    </div>

    <div class="flex items-center gap-3">
      <input v-model="search" type="search" class="input max-w-xs" placeholder="Search containers..." />
      <select v-model="filterStatus" class="input max-w-[140px]">
        <option value="">All status</option>
        <option value="running">Running</option>
        <option value="stopped">Stopped</option>
      </select>
      <div class="ml-auto text-xs text-gray-500">{{ filteredContainers.length }} containers</div>
    </div>

    <div class="table-container">
      <table class="table">
        <thead>
          <tr>
            <th>CTID</th>
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
          <template v-if="containersStore.loading">
            <tr v-for="i in 6" :key="i">
              <td colspan="9"><SkeletonLoader height="1.25rem" /></td>
            </tr>
          </template>
          <template v-else>
            <tr v-for="ctr in filteredContainers" :key="`${ctr.node}-${ctr.vmid}`">
              <td class="font-mono text-gray-400 text-xs">{{ ctr.vmid }}</td>
              <td class="font-medium text-white">{{ ctr.name }}</td>
              <td class="text-xs text-gray-400">{{ ctr.node }}</td>
              <td><StatusBadge :status="ctr.status" dot /></td>
              <td>
                <div class="flex items-center gap-2">
                  <span class="text-xs w-8 text-right text-gray-300">{{ Math.round(ctr.cpu * 100) }}%</span>
                  <UsageBar :value="ctr.cpu * 100" class="w-20" />
                </div>
              </td>
              <td>
                <div class="flex items-center gap-2">
                  <span class="text-xs w-8 text-right text-gray-300">{{ Math.round((ctr.mem / (ctr.maxmem || 1)) * 100) }}%</span>
                  <UsageBar :value="ctr.mem" :max="ctr.maxmem || 1" class="w-20" />
                </div>
              </td>
              <td class="text-xs text-gray-400">{{ formatBytes(ctr.maxdisk) }}</td>
              <td class="text-xs text-gray-400">{{ formatUptime(ctr.uptime) }}</td>
              <td>
                <div class="flex items-center gap-1">
                  <button
                    v-if="ctr.status !== 'running'"
                    class="btn-success py-1 px-2 text-xs"
                    @click="performAction(ctr, 'start')"
                    title="Start"
                  ><PlayIcon class="w-3 h-3" /></button>
                  <button
                    v-if="ctr.status === 'running'"
                    class="btn-danger py-1 px-2 text-xs"
                    @click="performAction(ctr, 'stop')"
                    title="Stop"
                  ><SquareIcon class="w-3 h-3" /></button>
                  <button
                    v-if="ctr.status === 'running'"
                    class="btn-secondary py-1 px-2 text-xs"
                    @click="performAction(ctr, 'reboot')"
                    title="Reboot"
                  ><RefreshCwIcon class="w-3 h-3" /></button>
                  <button
                    class="btn-ghost py-1 px-2 text-xs text-danger"
                    @click="handleDelete(ctr)"
                    title="Delete"
                  ><Trash2Icon class="w-3 h-3" /></button>
                </div>
              </td>
            </tr>
            <tr v-if="!filteredContainers.length">
              <td colspan="9" class="text-center text-gray-500 py-8">No containers found</td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Play as PlayIcon, Square as SquareIcon,
  RefreshCw as RefreshCwIcon, Trash2 as Trash2Icon,
} from 'lucide-vue-next'
import { useContainersStore } from '@/stores/containers'
import { useNodesStore } from '@/stores/nodes'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import UsageBar from '@/components/ui/UsageBar.vue'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'
import type { Container, VMAction } from '@/types'

const containersStore = useContainersStore()
const nodesStore = useNodesStore()
const search = ref('')
const filterStatus = ref('')

const filteredContainers = computed(() =>
  containersStore.containers.filter(c => {
    const matchSearch = !search.value ||
      c.name?.toLowerCase().includes(search.value.toLowerCase()) ||
      String(c.vmid).includes(search.value)
    const matchStatus = !filterStatus.value || c.status === filterStatus.value
    return matchSearch && matchStatus
  })
)

async function performAction(ctr: Container, action: VMAction) {
  await containersStore.performAction(ctr.node!, ctr.vmid, action)
}

async function handleDelete(ctr: Container) {
  if (!confirm(`Delete container ${ctr.name}?`)) return
  await containersStore.deleteContainer(ctr.node!, ctr.vmid)
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
  await Promise.all([nodesStore.fetchNodes(), containersStore.fetchAll()])
})
</script>
