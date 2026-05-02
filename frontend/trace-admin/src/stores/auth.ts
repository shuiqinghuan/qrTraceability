import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('access_token') || '')
  const refreshToken = ref(localStorage.getItem('refresh_token') || '')
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  const login = async (username: string, password: string) => {
    loading.value = true
    try {
      const { data } = await authApi.login({ username, password })
      token.value = data.access
      refreshToken.value = data.refresh
      localStorage.setItem('access_token', data.access)
      localStorage.setItem('refresh_token', data.refresh)
      router.push('/admin')
      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    try {
      await authApi.logout(refreshToken.value)
    } catch {
      // ignore
    }
    token.value = ''
    refreshToken.value = ''
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    router.push('/admin/login')
  }

  return {
    token,
    refreshToken,
    isAuthenticated,
    loading,
    login,
    logout
  }
})
