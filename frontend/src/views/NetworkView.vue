<template>
  <div class="space-y-6 animate-fade-in">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-white">Network</h2>
      <select v-model="selectedNode" class="input max-w-[180px]">
        <option v-for="n in nodesStore.nodes" :key="n.node" :value="n.node">{{ n.node }}</option>
      </select>
    </div>

    <div class="table-container">
      <table class="table">
        <thead>
          <tr>
            <th>Interface</th>
            <th>Type</th>
            <th>Status</th>
            <th>Address</th>
            <th>Gateway</th>
            <th>MTU</th>
            <th>Comment</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="loading">
            <tr v-for="i in 5" :key="i">
              <td colspan="7"><SkeletonLoader height="1.25rem" /></td>
            </tr>
          </template>
          <template v-else>
            <tr v-for="iface in interfaces" :key="iface.iface">
              <td class="font-mono text-sm text-white">{{ iface.iface }}</td>
              <td class="text-xs text-gray-400 uppercase">{{ iface.type }}</td>
              <td>
                <span :class="iface.active ? 'badge-running' : 'badge-stopped'">
                  {{ iface.active ? 'Up' : 'Down' }}
                </span>
              </td>
              <td class="font-mono text-xs text-gray-300">{{ iface.cidr || iface.address || '—' }}</td>
              <td class="font-mono text-xs text-gray-400">{{ iface.gateway || '—' }}</td>
              <td class="text-xs text-gray-400">{{ iface.mtu || '—' }}</td>
              <td class="text-xs text-gray-500">{{ iface.comments || '—' }}</td>
            </tr>
            <tr v-if="!interfaces.length && !loading">
              <td colspan="7" class="text-center text-gray-500 py-8">No interfaces found</td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { nodesApi } from '@/api'
import { useNodesStore } from '@/stores/nodes'
import type { NetworkInterface } from '@/types'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'

const nodesStore = useNodesStore()
const selectedNode = ref('')
const interfaces = ref<NetworkInterface[]>([])
const loading = ref(false)

async function fetchNetwork(node: string) {
  if (!node) return
  loading.value = true
  try {
    interfaces.value = await nodesApi.network(node)
  } catch {
    interfaces.value = []
  } finally {
    loading.value = false
  }
}

watch(selectedNode, fetchNetwork)

onMounted(async () => {
  await nodesStore.fetchNodes()
  if (nodesStore.nodes.length) {
    selectedNode.value = nodesStore.nodes[0].node
  }
})
</script>
