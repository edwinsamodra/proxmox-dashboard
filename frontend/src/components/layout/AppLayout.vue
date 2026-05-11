<template>
  <div class="flex h-screen overflow-hidden bg-surface">
    <!-- Sidebar -->
    <AppSidebar />

    <!-- Main content area -->
    <div class="flex flex-col flex-1 overflow-hidden">
      <!-- Top bar -->
      <AppTopBar />

      <!-- Page content -->
      <main class="flex-1 overflow-y-auto p-6">
        <RouterView v-slot="{ Component }">
          <Transition name="fade" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import AppSidebar from './AppSidebar.vue'
import AppTopBar from './AppTopBar.vue'
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
