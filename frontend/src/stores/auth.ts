import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const authenticated = ref(false)
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function login(username: string, password: string, realm = 'pam') {
    loading.value = true
    error.value = null
    try {
      await authApi.login(username, password, realm)
      authenticated.value = true
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Login failed'
      authenticated.value = false
    } finally {
      loading.value = false
    }
  }

  async function tokenLogin(tokenID: string, tokenSecret: string) {
    loading.value = true
    error.value = null
    try {
      await authApi.tokenAuth(tokenID, tokenSecret)
      authenticated.value = true
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : 'Token auth failed'
      authenticated.value = false
    } finally {
      loading.value = false
    }
  }

  function logout() {
    authenticated.value = false
  }

  return { authenticated, loading, error, login, tokenLogin, logout }
})
