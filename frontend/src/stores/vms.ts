import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { vmsApi } from '@/api'
import type { VM, VMAction } from '@/types'

export const useVMsStore = defineStore('vms', () => {
  const vms = ref<VM[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const actionLoading = ref<Record<string, boolean>>({})

  const runningVMs = computed(() => vms.value.filter(v => v.status === 'running'))
  const stoppedVMs = computed(() => vms.value.filter(v => v.status === 'stopped'))
  const templates = computed(() => vms.value.filter(v => v.template === 1))

  async function fetchAll() {
    loading.value = true
    error.value = null
    try {
      vms.value = await vmsApi.listAll()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch VMs'
    } finally {
      loading.value = false
    }
  }

  async function performAction(node: string, vmid: number, action: VMAction) {
    const key = `${node}-${vmid}`
    actionLoading.value[key] = true
    try {
      await vmsApi.action(node, vmid, action)
      // Optimistic UI update
      const vm = vms.value.find(v => v.vmid === vmid && v.node === node)
      if (vm) {
        if (action === 'start' || action === 'resume') vm.status = 'running'
        else if (action === 'stop' || action === 'shutdown') vm.status = 'stopped'
        else if (action === 'suspend') vm.status = 'suspended'
      }
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : `Action ${action} failed`
      throw e
    } finally {
      delete actionLoading.value[key]
    }
  }

  async function deleteVM(node: string, vmid: number) {
    await vmsApi.delete(node, vmid)
    vms.value = vms.value.filter(v => !(v.vmid === vmid && v.node === node))
  }

  function updateFromWS(updated: VM[]) {
    vms.value = updated
  }

  function isActionLoading(node: string, vmid: number) {
    return !!actionLoading.value[`${node}-${vmid}`]
  }

  return {
    vms, loading, error, actionLoading,
    runningVMs, stoppedVMs, templates,
    fetchAll, performAction, deleteVM, updateFromWS, isActionLoading,
  }
})
