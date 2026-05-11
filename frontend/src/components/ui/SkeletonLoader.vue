<template>
  <div class="skeleton" :style="{ height, width, borderRadius }" />
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  height?: string
  width?: string
  borderRadius?: string
}>(), {
  height: '1rem',
  width: '100%',
  borderRadius: '0.375rem',
})
</script>
