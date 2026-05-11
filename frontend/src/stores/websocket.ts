import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { WSMessage } from '@/types'
import { useNodesStore } from './nodes'
import { useVMsStore } from './vms'
import { useContainersStore } from './containers'

type WSStatus = 'connecting' | 'connected' | 'disconnected' | 'error'

export const useWSStore = defineStore('ws', () => {
  const status = ref<WSStatus>('disconnected')
  const messages = ref<WSMessage[]>([])
  let socket: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let reconnectDelay = 1000

  function connect(url?: string) {
    const wsUrl = url ?? buildWsUrl()
    status.value = 'connecting'

    socket = new WebSocket(wsUrl)

    socket.onopen = () => {
      status.value = 'connected'
      reconnectDelay = 1000
    }

    socket.onmessage = (event) => {
      try {
        const msg: WSMessage = JSON.parse(event.data)
        messages.value.push(msg)
        if (messages.value.length > 200) messages.value.shift()
        handleMessage(msg)
      } catch {
        // ignore malformed messages
      }
    }

    socket.onclose = () => {
      status.value = 'disconnected'
      scheduleReconnect(wsUrl)
    }

    socket.onerror = () => {
      status.value = 'error'
      socket?.close()
    }
  }

  function handleMessage(msg: WSMessage) {
    const nodesStore = useNodesStore()
    const vmsStore = useVMsStore()
    const containersStore = useContainersStore()

    switch (msg.type) {
      case 'nodes_update':
        nodesStore.updateFromWS(msg.payload as never)
        break
      case 'vms_update':
        vmsStore.updateFromWS(msg.payload as never)
        break
      case 'containers_update':
        containersStore.updateFromWS(msg.payload as never)
        break
    }
  }

  function scheduleReconnect(url: string) {
    if (reconnectTimer) clearTimeout(reconnectTimer)
    reconnectTimer = setTimeout(() => {
      reconnectDelay = Math.min(reconnectDelay * 2, 30000)
      connect(url)
    }, reconnectDelay)
  }

  function disconnect() {
    if (reconnectTimer) clearTimeout(reconnectTimer)
    socket?.close()
    socket = null
    status.value = 'disconnected'
  }

  function buildWsUrl() {
    const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    return `${proto}//${host}/api/ws`
  }

  return { status, messages, connect, disconnect }
})
