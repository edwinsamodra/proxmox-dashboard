<template>
  <div class="space-y-6 animate-fade-in">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-white">Backup & Restore</h2>
      <button class="btn-secondary" @click="fetchBackups">
        <RefreshCwIcon class="w-4 h-4" :class="{ 'animate-spin': loading }" />
        Refresh
      </button>
    </div>

    <div class="table-container">
      <div class="p-4 border-b border-surface-300 text-sm font-semibold text-white">Backup Jobs</div>
      <table class="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Node</th>
            <th>Storage</th>
            <th>Schedule</th>
            <th>Mode</th>
            <th>Compression</th>
            <th>VMs</th>
            <th>Enabled</th>
          </tr>
        </thead>
        <tbody>
          <template v-if="loading">
            <tr v-for="i in 3" :key="i">
              <td colspan="8"><SkeletonLoader height="1.25rem" /></td>
            </tr>
          </template>
          <template v-else>
            <tr v-for="job in backups" :key="job.id">
              <td class="font-mono text-xs text-gray-400">{{ job.id || '—' }}</td>
              <td class="text-sm text-gray-300">{{ job.node || 'all' }}</td>
              <td class="text-sm text-gray-300">{{ job.storage || '—' }}</td>
              <td class="font-mono text-xs text-gray-300">{{ job.starttime || '—' }}</td>
              <td class="text-xs text-gray-400 uppercase">{{ job.mode || '—' }}</td>
              <td class="text-xs text-gray-400 uppercase">{{ job.compress || 'none' }}</td>
              <td class="font-mono text-xs text-gray-400">{{ job.vmid || 'all' }}</td>
              <td>
                <span :class="job.enabled !== 0 ? 'badge-running' : 'badge-stopped'">
                  {{ job.enabled !== 0 ? 'Yes' : 'No' }}
                </span>
              </td>
            </tr>
            <tr v-if="!backups.length && !loading">
              <td colspan="8" class="text-center text-gray-500 py-8">No backup jobs configured</td>
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
import { backupApi } from '@/api'
import type { BackupJob } from '@/types'
import SkeletonLoader from '@/components/ui/SkeletonLoader.vue'

const backups = ref<BackupJob[]>([])
const loading = ref(false)

async function fetchBackups() {
  loading.value = true
  try {
    backups.value = await backupApi.list()
  } catch {
    backups.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchBackups)
</script>
