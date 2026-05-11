<template>
  <div class="min-h-screen bg-surface flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-10">
        <div class="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-accent/20 border border-accent/30 mb-4 shadow-glow">
          <ServerIcon class="w-7 h-7 text-accent" />
        </div>
        <h1 class="text-2xl font-bold text-white tracking-tight">Proxmox Modern Orchestrator</h1>
        <p class="text-gray-500 text-sm mt-1">Sign in to manage your infrastructure</p>
      </div>

      <!-- Auth card -->
      <div class="card p-8">
        <!-- Tab switcher -->
        <div class="flex rounded-lg bg-surface-200 p-1 mb-6 gap-1">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            :class="['flex-1 py-1.5 text-sm font-medium rounded-md transition-colors duration-150',
              activeTab === tab.id
                ? 'bg-surface-400 text-white shadow-sm'
                : 'text-gray-400 hover:text-gray-200']"
            @click="activeTab = tab.id"
          >
            {{ tab.label }}
          </button>
        </div>

        <!-- Password login form -->
        <form v-if="activeTab === 'password'" @submit.prevent="handlePasswordLogin" class="space-y-4">
          <div>
            <label class="label">Username</label>
            <input v-model="form.username" type="text" class="input" placeholder="root" autocomplete="username" required />
          </div>
          <div>
            <label class="label">Password</label>
            <input v-model="form.password" type="password" class="input" placeholder="••••••••" autocomplete="current-password" required />
          </div>
          <div>
            <label class="label">Realm</label>
            <select v-model="form.realm" class="input">
              <option value="pam">PAM (Linux)</option>
              <option value="pve">Proxmox VE</option>
            </select>
          </div>
          <div v-if="authStore.error" class="text-danger text-sm bg-danger/10 border border-danger/30 rounded-lg px-3 py-2">
            {{ authStore.error }}
          </div>
          <button type="submit" class="btn-primary w-full justify-center" :disabled="authStore.loading">
            <LoaderIcon v-if="authStore.loading" class="w-4 h-4 animate-spin" />
            <span>{{ authStore.loading ? 'Signing in...' : 'Sign in' }}</span>
          </button>
        </form>

        <!-- API Token form -->
        <form v-else @submit.prevent="handleTokenLogin" class="space-y-4">
          <div>
            <label class="label">Token ID</label>
            <input v-model="tokenForm.id" type="text" class="input" placeholder="user@pam!token-name" autocomplete="off" required />
          </div>
          <div>
            <label class="label">Token Secret</label>
            <input v-model="tokenForm.secret" type="password" class="input" placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" autocomplete="off" required />
          </div>
          <div v-if="authStore.error" class="text-danger text-sm bg-danger/10 border border-danger/30 rounded-lg px-3 py-2">
            {{ authStore.error }}
          </div>
          <button type="submit" class="btn-primary w-full justify-center" :disabled="authStore.loading">
            <LoaderIcon v-if="authStore.loading" class="w-4 h-4 animate-spin" />
            <span>{{ authStore.loading ? 'Authenticating...' : 'Authenticate' }}</span>
          </button>
        </form>
      </div>

      <p class="text-center text-xs text-gray-600 mt-6">
        PMO — Open source Proxmox alternative dashboard
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Server as ServerIcon, Loader as LoaderIcon } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

type TabId = 'password' | 'token'
const tabs = [
  { id: 'password' as TabId, label: 'Password' },
  { id: 'token' as TabId, label: 'API Token' },
]
const activeTab = ref<TabId>('password')

const form = ref({ username: 'root', password: '', realm: 'pam' })
const tokenForm = ref({ id: '', secret: '' })

async function handlePasswordLogin() {
  await authStore.login(form.value.username, form.value.password, form.value.realm)
  if (authStore.authenticated) {
    router.push('/')
  }
}

async function handleTokenLogin() {
  await authStore.tokenLogin(tokenForm.value.id, tokenForm.value.secret)
  if (authStore.authenticated) {
    router.push('/')
  }
}
</script>
