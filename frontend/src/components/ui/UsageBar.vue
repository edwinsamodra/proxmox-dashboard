<template>
  <div class="progress-bar">
    <div
      class="progress-fill"
      :class="fillClass"
      :style="{ width: clampedPercent + '%' }"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  value: number
  max?: number
  variant?: 'default' | 'success' | 'warning' | 'danger'
}>()

const percent = computed(() =>
  props.max ? (props.value / props.max) * 100 : props.value
)

const clampedPercent = computed(() =>
  Math.min(100, Math.max(0, percent.value))
)

const fillClass = computed(() => {
  const v = clampedPercent.value
  if (props.variant) {
    if (props.variant === 'success') return 'bg-success'
    if (props.variant === 'warning') return 'bg-warning'
    if (props.variant === 'danger') return 'bg-danger'
    return 'bg-accent'
  }
  // Auto color based on usage
  if (v > 85) return 'bg-danger'
  if (v > 65) return 'bg-warning'
  return 'bg-accent'
})
</script>
