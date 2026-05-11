import { defineStore } from 'pinia'
import { ref } from 'vue'
import { tofuApi } from '@/api'
import type { TofuWorkspace, TofuExecResult, VMProvisionRequest } from '@/types'

export const useTofuStore = defineStore('tofu', () => {
  const available = ref(false)
  const workspaces = ref<TofuWorkspace[]>([])
  const loading = ref(false)
  const executing = ref(false)
  const error = ref<string | null>(null)
  const lastResult = ref<TofuExecResult | null>(null)
  const currentWorkspace = ref<string | null>(null)
  const currentTF = ref('')

  async function checkStatus() {
    try {
      const status = await tofuApi.status()
      available.value = status.available
      workspaces.value = status.workspaces ?? []
    } catch {
      available.value = false
    }
  }

  async function fetchWorkspaces() {
    loading.value = true
    try {
      workspaces.value = await tofuApi.listWorkspaces()
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Failed to list workspaces'
    } finally {
      loading.value = false
    }
  }

  async function selectWorkspace(name: string) {
    currentWorkspace.value = name
    try {
      const result = await tofuApi.getWorkspaceTF(name)
      currentTF.value = result.content
    } catch {
      currentTF.value = ''
    }
  }

  async function createWorkspace(name: string, content = '') {
    await tofuApi.createWorkspace(name, content)
    await fetchWorkspaces()
  }

  async function saveCurrentTF() {
    if (!currentWorkspace.value) return
    await tofuApi.saveWorkspaceTF(currentWorkspace.value, currentTF.value)
  }

  async function runInit() {
    if (!currentWorkspace.value) return
    executing.value = true
    try {
      lastResult.value = await tofuApi.init(currentWorkspace.value)
    } finally {
      executing.value = false
    }
  }

  async function runPlan(vars?: Record<string, string>) {
    if (!currentWorkspace.value) return
    executing.value = true
    try {
      lastResult.value = await tofuApi.plan(currentWorkspace.value, vars)
    } finally {
      executing.value = false
    }
  }

  async function runApply(vars?: Record<string, string>) {
    if (!currentWorkspace.value) return
    executing.value = true
    try {
      lastResult.value = await tofuApi.apply(currentWorkspace.value, vars)
    } finally {
      executing.value = false
    }
  }

  async function runDestroy(vars?: Record<string, string>) {
    if (!currentWorkspace.value) return
    executing.value = true
    try {
      lastResult.value = await tofuApi.destroy(currentWorkspace.value, vars)
    } finally {
      executing.value = false
    }
  }

  async function deleteWorkspace(name: string) {
    await tofuApi.deleteWorkspace(name)
    if (currentWorkspace.value === name) {
      currentWorkspace.value = null
      currentTF.value = ''
    }
    workspaces.value = workspaces.value.filter(w => w.name !== name)
  }

  async function generateAndProvisionVM(req: VMProvisionRequest) {
    executing.value = true
    lastResult.value = null
    try {
      lastResult.value = await tofuApi.provisionVM(req)
      return lastResult.value
    } finally {
      executing.value = false
    }
  }

  return {
    available, workspaces, loading, executing, error, lastResult,
    currentWorkspace, currentTF,
    checkStatus, fetchWorkspaces, selectWorkspace,
    createWorkspace, saveCurrentTF,
    runInit, runPlan, runApply, runDestroy, deleteWorkspace,
    generateAndProvisionVM,
  }
})
