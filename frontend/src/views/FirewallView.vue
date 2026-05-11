<template>
  <div class="space-y-6 animate-fade-in">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-white">Firewall Rules</h2>
    </div>

    <div class="table-container">
      <div class="flex items-center justify-between p-4 border-b border-surface-300">
        <div class="text-sm font-semibold text-white">Cluster Firewall Rules</div>
        <button class="btn-secondary text-xs py-1.5" @click="fetchRules">
          <RefreshCwIcon class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
          Refresh
        </button>
      </div>
      <table class="table">
        <thead>
          <tr>
            <th>Pos</th>
            <th>Type</th>
            <th>Action</th>
            <th>Protocol</th>
            <th>Source</th>
            <th>Destination</th>
            <th>Dest Port</th>
            <th>Enabled</th>
            <th>Comment</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="loading">
            <tr v-for="i in 5" :key="i">
              <td colspan="9"><SkeletonLoader height="1.25rem" /></td>
            </tr>
          </template>
          <template v-else>
            <tr v-for="rule in rules" :key="rule.pos">
              <td class="text-xs text-gray-400 font-mono">{{ rule.pos }}</td>
              <td class="text-xs uppercase text-gray-400">{{ rule.type }}</td>
              <td>
                <span :class="rule.action === 'ACCEPT' ? 'badge-running' : 'badge-stopped'">
                  {{ rule.action }}
                </span>
              </td>
              <td class="text-xs text-gray-400">{{ rule.proto || 'any' }}</td>
              <td class="font-mono text-xs text-gray-300">{{ rule.source || 'any' }}</td>
              <td class="font-mono text-xs text-gray-300">{{ rule.dest || 'any' }}</td>
              <td class="font-mono text-xs text-gray-300">{{ rule.dport || '—' }}</td>
              <td>
                <span :class="rule.enable ? 'badge-running' : 'badge-stopped'">
                  {{ rule.enable ? 'Yes' : 'No' }}
                </span>
              </td>
              <td class="text-xs text-gray-500">{{ rule.comment || '—' }}</td>
            </tr>
            <tr v-if="!rules.length && !loading">
              <td colspan="9" class="text-center text-gray-500 py-8">No firewall rules configured</td>
            </tr>
          </template>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RefreshCw as RefreshCwIcon } from 'lucide-vue-next'
import { firewallApi } from '@/api'
import type { FirewallRule } from '@/types'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'

const rules = ref<FirewallRule[]>([])
const loading = ref(false)

async function fetchRules() {
  loading.value = true
  try {
    rules.value = await firewallApi.clusterRules()
  } catch {
    rules.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchRules)
</script>
