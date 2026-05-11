<template>
  <header class="h-14 flex items-center justify-between px-6 border-b border-surface-300 bg-surface-50/80 backdrop-blur-sm flex-shrink-0">
    <!-- Page title -->
    <div class="flex items-center gap-2">
      <h1 class="text-sm font-semibold text-white">{{ pageTitle }}</h1>
      <span v-if="subtitle" class="text-gray-500 text-sm">/ {{ subtitle }}</span>
    </div>

    <!-- Right side controls -->
    <div class="flex items-center gap-3">
      <!-- Refresh button -->
      <button
        class="btn-ghost p-2 rounded-lg"
        title="Refresh data"
        @click="handleRefresh"
      >
        <RefreshCwIcon class="w-4 h-4" :class="{ 'animate-spin': refreshing }" />
      </button>

      <!-- Proxmox host indicator -->
      <div class="hidden md:flex items-center gap-2 text-xs text-gray-500 bg-surface-200 px-3 py-1.5 rounded-lg border border-surface-300">
        <ServerIcon class="w-3.5 h-3.5 text-accent" />
        <span>{{ proxmoxHost || 'Not configured' }}</span>
      </div>

      <!-- WS status dot -->
      <div
        class="w-2 h-2 rounded-full"
        :class="wsStore.status === 'connected' ? 'bg-success animate-pulse' : 'bg-gray-500'"
        :title="`WebSocket: ${wsStore.status}`"
      />
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { RefreshCw as RefreshCwIcon, Server as ServerIcon } from 'lucide-vue-next'
import { useWSStore } from '@/stores/websocket'
import { useNodesStore } from '@/stores/nodes'
import { useVMsStore } from '@/stores/vms'
import { useContainersStore } from '@/stores/containers'

const route = useRoute()
const wsStore = useWSStore()
const nodesStore = useNodesStore()
const vmsStore = useVMsStore()
const containersStore = useContainersStore()
const refreshing = ref(false)

const pageTitles: Record<string, string> = {
  '/': 'Dashboard',
  '/nodes': 'Nodes',
  '/vms': 'Virtual Machines',
  '/containers': 'Containers',
  '/storage': 'Storage',
  '/network': 'Network',
  '/firewall': 'Firewall',
  '/backup': 'Backup & Restore',
  '/tofu': 'OpenTofu IaC',
}

const pageTitle = computed(() => pageTitles[route.path] ?? route.name?.toString() ?? 'PMO')
const subtitle = computed(() => '')
const proxmoxHost = ((import.meta as unknown as { env: Record<string, string> }).env)?.VITE_PROXMOX_HOST ?? ''

async function handleRefresh() {
  refreshing.value = true
  try {
    await Promise.all([
      nodesStore.fetchNodes(),
      vmsStore.fetchAll(),
      containersStore.fetchAll(),
    ])
  } finally {
    setTimeout(() => { refreshing.value = false }, 500)
  }
}
</script>
