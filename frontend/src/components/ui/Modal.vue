<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Backdrop -->
        <div
          class="absolute inset-0 bg-black/60 backdrop-blur-sm"
          @click="$emit('update:modelValue', false)"
        />
        <!-- Dialog -->
        <div
          class="relative z-10 card w-full max-h-[90vh] overflow-y-auto animate-fade-in"
          :class="sizeClass"
        >
          <!-- Header -->
          <div class="flex items-center justify-between p-5 border-b border-surface-300">
            <h2 class="text-base font-semibold text-white">{{ title }}</h2>
            <button
              class="btn-ghost p-1.5 rounded-lg"
              @click="$emit('update:modelValue', false)"
            >
              <XIcon class="w-4 h-4" />
            </button>
          </div>
          <!-- Body -->
          <div class="p-5">
            <slot />
          </div>
          <!-- Footer -->
          <div v-if="$slots.footer" class="flex items-center justify-end gap-3 p-5 border-t border-surface-300">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { X as XIcon } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  modelValue: boolean
  title: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
}>(), { size: 'md' })

defineEmits<{ 'update:modelValue': [value: boolean] }>()

const sizeClass = computed(() => {
  if (props.size === 'sm') return 'max-w-sm'
  if (props.size === 'lg') return 'max-w-2xl'
  if (props.size === 'xl') return 'max-w-4xl'
  return 'max-w-lg'
})
</script>

<style scoped>
.modal-enter-active, .modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from, .modal-leave-to {
  opacity: 0;
}
</style>
