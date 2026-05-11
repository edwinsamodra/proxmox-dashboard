<template>
  <div class="space-y-6 animate-fade-in">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-white">Storage</h2>
      <select v-model="selectedNode" class="input max-w-[180px]">
        <option v-for="n in nodesStore.nodes" :key="n.node" :value="n.node">{{ n.node }}</option>
      </select>
    </div>

    <div v-if="loading" class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
      <SkeletonLoader v-for="i in 4" :key="i" height="10rem" />
    </div>
    <div v-else class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
      <div v-for="s in storages" :key="s.storage" class="card p-5">
        <div class="flex items-start justify-between mb-3">
          <div>
            <div class="font-semibold text-white">{{ s.storage }}</div>
            <div class="text-xs text-gray-500 mt-0.5 uppercase">{{ s.type }}</div>
          </div>
          <span :class="s.active ? 'badge-running' : 'badge-stopped'">
            {{ s.active ? 'Active' : 'Inactive' }}
          </span>
        </div>
        <div class="space-y-2">
          <div class="flex justify-between text-xs mb-1">
            <span class="text-gray-400">Usage</span>
            <span class="text-gray-300">{{ formatBytes(s.used) }} / {{ formatBytes(s.total) }}</span>
          </div>
          <UsageBar :value="s.used" :max="s.total || 1" />
          <div class="flex justify-between text-xs text-gray-500 mt-2">
            <span>Available: {{ formatBytes(s.avail) }}</span>
            <span>{{ s.content }}</span>
          </div>
        </div>
      </div>
      <div v-if="!storages.length && !loading" class="card p-8 text-center text-gray-500 col-span-full">
        No storage found for node {{ selectedNode }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { storageApi } from '@/api'
import { useNodesStore } from '@/stores/nodes'
import type { Storage } from '@/types'
import UsageBar from '@/components/ui/UsageBar.vue'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'

const nodesStore = useNodesStore()
const selectedNode = ref('')
const storages = ref<Storage[]>([])
const loading = ref(false)

async function fetchStorage(node: string) {
  if (!node) return
  loading.value = true
  try {
    storages.value = await storageApi.list(node)
  } catch {
    storages.value = []
  } finally {
    loading.value = false
  }
}

watch(selectedNode, fetchStorage)

function formatBytes(bytes: number): string {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0; let val = bytes
  while (val >= 1024 && i < units.length - 1) { val /= 1024; i++ }
  return `${val.toFixed(i > 1 ? 1 : 0)} ${units[i]}`
}

onMounted(async () => {
  await nodesStore.fetchNodes()
  if (nodesStore.nodes.length) {
    selectedNode.value = nodesStore.nodes[0].node
  }
})
</script>
