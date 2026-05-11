import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { containersApi } from '@/api'
import type { Container, VMAction } from '@/types'

export const useContainersStore = defineStore('containers', () => {
  const containers = ref<Container[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const runningContainers = computed(() =>
    containers.value.filter(c => c.status === 'running')
  )
  const stoppedContainers = computed(() =>
    containers.value.filter(c => c.status === 'stopped')
  )

  async function fetchAll() {
    loading.value = true
    error.value = null
    try {
      containers.value = await containersApi.listAll()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch containers'
    } finally {
      loading.value = false
    }
  }

  async function performAction(node: string, vmid: number, action: VMAction) {
    try {
      await containersApi.action(node, vmid, action)
      const ctr = containers.value.find(c => c.vmid === vmid && c.node === node)
      if (ctr) {
        if (action === 'start' || action === 'resume') ctr.status = 'running'
        else if (action === 'stop' || action === 'shutdown') ctr.status = 'stopped'
      }
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : `Action ${action} failed`
      throw e
    }
  }

  async function deleteContainer(node: string, vmid: number) {
    await containersApi.delete(node, vmid)
    containers.value = containers.value.filter(
      c => !(c.vmid === vmid && c.node === node)
    )
  }

  function updateFromWS(updated: Container[]) {
    containers.value = updated
  }

  return {
    containers, loading, error,
    runningContainers, stoppedContainers,
    fetchAll, performAction, deleteContainer, updateFromWS,
  }
})
