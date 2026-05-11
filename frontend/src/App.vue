<template>
  <RouterView />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useWSStore } from '@/stores/websocket'

const wsStore = useWSStore()

onMounted(() => {
  // Connect WebSocket for real-time updates
  wsStore.connect()
})
</script>
