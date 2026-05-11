import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { nodesApi } from '@/api'
import type { Node } from '@/types'

export const useNodesStore = defineStore('nodes', () => {
  const nodes = ref<Node[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const onlineNodes = computed(() => nodes.value.filter(n => n.status === 'online'))
  const nodeNames = computed(() => nodes.value.map(n => n.node))

  const totalCPU = computed(() =>
    nodes.value.reduce((sum, n) => sum + n.maxcpu, 0)
  )
  const usedCPU = computed(() =>
    nodes.value.reduce((sum, n) => sum + n.cpu * n.maxcpu, 0)
  )
  const totalMem = computed(() =>
    nodes.value.reduce((sum, n) => sum + n.maxmem, 0)
  )
  const usedMem = computed(() =>
    nodes.value.reduce((sum, n) => sum + n.mem, 0)
  )

  async function fetchNodes() {
    loading.value = true
    error.value = null
    try {
      nodes.value = await nodesApi.list()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch nodes'
    } finally {
      loading.value = false
    }
  }

  function updateFromWS(updated: Node[]) {
    nodes.value = updated
  }

  return {
    nodes, loading, error,
    onlineNodes, nodeNames,
    totalCPU, usedCPU, totalMem, usedMem,
    fetchNodes, updateFromWS,
  }
})
