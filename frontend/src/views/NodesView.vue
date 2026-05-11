<template>
  <div class="space-y-6 animate-fade-in">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-white">Nodes</h2>
      <button class="btn-secondary" @click="nodesStore.fetchNodes()">
        <RefreshCwIcon class="w-4 h-4" :class="{ 'animate-spin': nodesStore.loading }" />
        Refresh
      </button>
    </div>

    <!-- Node cards -->
    <template v-if="nodesStore.loading">
      <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
        <SkeletonLoader v-for="i in 3" :key="i" height="12rem" />
      </div>
    </template>
    <template v-else>
      <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
        <div
          v-for="node in nodesStore.nodes"
          :key="node.node"
          class="card p-5 hover:border-accent/30 transition-colors cursor-pointer"
          @click="selectNode(node)"
        >
          <div class="flex items-start justify-between mb-4">
            <div>
              <div class="font-semibold text-white text-base">{{ node.node }}</div>
              <div class="text-xs text-gray-500 mt-0.5">{{ node.type || 'node' }}</div>
            </div>
            <StatusBadge :status="node.status" dot />
          </div>

          <div class="space-y-3">
            <!-- CPU -->
            <div>
              <div class="flex justify-between text-xs mb-1">
                <span class="text-gray-400">CPU</span>
                <span class="text-gray-300">{{ Math.round(node.cpu * 100) }}% / {{ node.maxcpu }} cores</span>
              </div>
              <UsageBar :value="node.cpu * 100" />
            </div>
            <!-- Memory -->
            <div>
              <div class="flex justify-between text-xs mb-1">
                <span class="text-gray-400">Memory</span>
                <span class="text-gray-300">{{ formatBytes(node.mem) }} / {{ formatBytes(node.maxmem) }}</span>
              </div>
              <UsageBar :value="node.mem" :max="node.maxmem || 1" />
            </div>
            <!-- Disk -->
            <div>
              <div class="flex justify-between text-xs mb-1">
                <span class="text-gray-400">Root Disk</span>
                <span class="text-gray-300">{{ formatBytes(node.disk) }} / {{ formatBytes(node.maxdisk) }}</span>
              </div>
              <UsageBar :value="node.disk" :max="node.maxdisk || 1" />
            </div>
          </div>

          <div class="mt-4 pt-3 border-t border-surface-300 flex items-center justify-between text-xs text-gray-500">
            <span>Up {{ formatUptime(node.uptime) }}</span>
            <span>{{ node.level || '—' }}</span>
          </div>
        </div>
      </div>

      <div v-if="!nodesStore.nodes.length" class="card p-12 text-center text-gray-500">
        No nodes found. Check your Proxmox connection settings.
      </div>
    </template>

    <!-- Node detail modal -->
    <Modal v-if="selectedNode" v-model="showDetail" :title="selectedNode.node" size="lg">
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-3 text-sm">
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Status</div>
            <StatusBadge :status="selectedNode.status" dot />
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Uptime</div>
            <div class="text-white font-medium">{{ formatUptime(selectedNode.uptime) }}</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">CPU Cores</div>
            <div class="text-white font-medium">{{ selectedNode.maxcpu }} cores</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Total Memory</div>
            <div class="text-white font-medium">{{ formatBytes(selectedNode.maxmem) }}</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">CPU Usage</div>
            <div class="text-white font-medium">{{ Math.round(selectedNode.cpu * 100) }}%</div>
          </div>
          <div class="bg-surface-200 rounded-lg p-3">
            <div class="text-xs text-gray-500 mb-1">Memory Used</div>
            <div class="text-white font-medium">{{ formatBytes(selectedNode.mem) }}</div>
          </div>
        </div>
      </div>
      <template #footer>
        <button class="btn-secondary" @click="showDetail = false">Close</button>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RefreshCw as RefreshCwIcon } from 'lucide-vue-next'
import { useNodesStore } from '@/stores/nodes'
import StatusBadge from '@/components/ui/StatusBadge.vue'
import UsageBar from '@/components/ui/UsageBar.vue'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'
import Modal from '@/components/ui/Modal.vue'
import type { Node } from '@/types'

const nodesStore = useNodesStore()
const selectedNode = ref<Node | null>(null)
const showDetail = ref(false)

function selectNode(node: Node) {
  selectedNode.value = node
  showDetail.value = true
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

onMounted(() => nodesStore.fetchNodes())
</script>
