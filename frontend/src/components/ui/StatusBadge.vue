<template>
  <span :class="badgeClass">
    <span v-if="dot" class="w-1.5 h-1.5 rounded-full" :class="dotClass" />
    <slot>{{ label }}</slot>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
  label?: string
  dot?: boolean
}>()

const badgeClass = computed(() => {
  const s = props.status?.toLowerCase()
  if (s === 'running' || s === 'online') return 'badge-running'
  if (s === 'stopped' || s === 'offline') return 'badge-stopped'
  if (s === 'paused' || s === 'suspended') return 'badge-paused'
  if (s === 'error' || s === 'unknown') return 'badge-error'
  return 'badge bg-surface-400 text-gray-400 border border-surface-400'
})

const dotClass = computed(() => {
  const s = props.status?.toLowerCase()
  if (s === 'running' || s === 'online') return 'bg-success animate-pulse'
  if (s === 'stopped' || s === 'offline') return 'bg-gray-500'
  if (s === 'paused' || s === 'suspended') return 'bg-warning animate-pulse'
  return 'bg-gray-500'
})
</script>
