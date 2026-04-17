import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api/axios.js'

export const useAuthStore = defineStore('auth', () => {
  const user            = ref(null)
  const isInitializing  = ref(false)  // true only during the initial fetchMe on first load
  const isLoggedIn      = computed(() => !!user.value)
  const isSuperAdmin    = computed(() => user.value?.role === 'super_admin')

  async function register(payload) {
    const { data } = await api.post('/auth/register', payload)
    if (data.success) user.value = data.data.user
    return data
  }

  async function login(payload) {
    const { data } = await api.post('/auth/login', payload)
    if (data.success) user.value = data.data.user
    return data
  }

  async function logout() {
    try { await api.post('/auth/logout') }
    finally { user.value = null }
  }

  async function fetchMe() {
    const { data } = await api.get('/auth/me')
    if (data.success && data.data) user.value = data.data
    return data
  }

  // Called by router guard on very first navigation — shows fullscreen spinner in App.vue
  async function initSession() {
    isInitializing.value = true
    try { await fetchMe() } catch { /* unauthenticated — that's fine */ }
    finally { isInitializing.value = false }
  }

  return { user, isLoggedIn, isSuperAdmin, isInitializing, register, login, logout, fetchMe, initSession }
})
