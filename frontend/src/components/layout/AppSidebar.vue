<template>
  <aside class="w-60 flex-shrink-0 bg-surface-50 border-r border-surface-300 flex flex-col h-full">
    <!-- Logo / Brand -->
    <div class="flex items-center gap-3 px-4 py-5 border-b border-surface-300">
      <div class="w-8 h-8 rounded-lg bg-accent flex items-center justify-center">
        <ServerIcon class="w-4 h-4 text-white" />
      </div>
      <div>
        <div class="text-sm font-bold text-white tracking-tight">PMO</div>
        <div class="text-xs text-gray-500 leading-none">Proxmox Modern Orchestrator</div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-3 py-4 space-y-0.5 overflow-y-auto">
      <div class="px-3 mb-2">
        <span class="text-xs font-semibold text-gray-600 uppercase tracking-wider">Overview</span>
      </div>
      <RouterLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        custom
        v-slot="{ isActive, navigate }"
      >
        <button
          :class="['sidebar-item w-full', isActive ? 'active' : '']"
          @click="navigate"
        >
          <component :is="item.icon" class="w-4 h-4 flex-shrink-0" />
          <span>{{ item.label }}</span>
          <span
            v-if="item.badge !== undefined"
            class="ml-auto text-xs font-semibold px-1.5 py-0.5 rounded-md"
            :class="item.badge > 0 ? 'bg-accent/20 text-accent' : 'bg-surface-400 text-gray-500'"
          >
            {{ item.badge }}
          </span>
        </button>
      </RouterLink>

      <div class="px-3 mt-4 mb-2">
        <span class="text-xs font-semibold text-gray-600 uppercase tracking-wider">Infrastructure as Code</span>
      </div>
      <RouterLink to="/tofu" custom v-slot="{ isActive, navigate }">
        <button :class="['sidebar-item w-full', isActive ? 'active' : '']" @click="navigate">
          <CodeIcon class="w-4 h-4 flex-shrink-0" />
          <span>OpenTofu</span>
          <span class="ml-auto text-xs bg-accent/20 text-accent px-1.5 py-0.5 rounded-md font-semibold">IaC</span>
        </button>
      </RouterLink>
    </nav>

    <!-- Bottom: WS status -->
    <div class="px-4 py-3 border-t border-surface-300">
      <div class="flex items-center gap-2 text-xs">
        <span
          class="w-2 h-2 rounded-full"
          :class="wsStatusClass"
        />
        <span class="text-gray-500">{{ wsStatusLabel }}</span>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Server as ServerIcon,
  LayoutDashboard,
  Cpu,
  Box,
  HardDrive,
  Network,
  Shield,
  Archive,
  Code as CodeIcon,
} from 'lucide-vue-next'
import { useWSStore } from '@/stores/websocket'
import { useNodesStore } from '@/stores/nodes'
import { useVMsStore } from '@/stores/vms'
import { useContainersStore } from '@/stores/containers'

const wsStore = useWSStore()
const nodesStore = useNodesStore()
const vmsStore = useVMsStore()
const containersStore = useContainersStore()

const navItems = computed(() => [
  { to: '/', label: 'Dashboard', icon: LayoutDashboard },
  { to: '/nodes', label: 'Nodes', icon: ServerIcon, badge: nodesStore.nodes.length },
  { to: '/vms', label: 'Virtual Machines', icon: Cpu, badge: vmsStore.vms.length },
  { to: '/containers', label: 'Containers', icon: Box, badge: containersStore.containers.length },
  { to: '/storage', label: 'Storage', icon: HardDrive },
  { to: '/network', label: 'Network', icon: Network },
  { to: '/firewall', label: 'Firewall', icon: Shield },
  { to: '/backup', label: 'Backup', icon: Archive },
])

const wsStatusClass = computed(() => {
  if (wsStore.status === 'connected') return 'bg-success animate-pulse'
  if (wsStore.status === 'connecting') return 'bg-warning animate-pulse'
  return 'bg-danger'
})

const wsStatusLabel = computed(() => {
  if (wsStore.status === 'connected') return 'Live updates active'
  if (wsStore.status === 'connecting') return 'Connecting...'
  return 'Disconnected'
})
</script>
