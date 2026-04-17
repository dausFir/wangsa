<template>
  <div class="min-h-screen flex">

    <!-- ── Left decorative panel (desktop only) ── -->
    <div class="hidden lg:flex lg:w-1/2 bg-navy flex-col justify-between p-12 relative overflow-hidden">
      <!-- Background decoration -->
      <div class="absolute inset-0 opacity-[0.04] pointer-events-none select-none">
        <svg viewBox="0 0 400 500" class="w-full h-full" aria-hidden="true">
          <circle cx="320" cy="80"  r="200" fill="white"/>
          <circle cx="80"  cy="380" r="160" fill="white"/>
        </svg>
      </div>

      <!-- Logo -->
      <div class="flex items-center gap-3 relative">
        <div class="w-10 h-10 rounded-xl bg-terra flex items-center justify-center shadow-soft">
          <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
          </svg>
        </div>
        <span class="text-white font-bold text-xl tracking-tight">Wangsa</span>
      </div>

      <!-- Tagline -->
      <div class="relative">
        <h2 class="text-4xl font-bold text-white leading-snug mb-5">
          Satu platform<br/>untuk seluruh<br/>keluarga besar.
        </h2>
        <p class="text-white/50 text-lg leading-relaxed">
          Silsilah, kas, kalender, dan peta domisili —<br/>semuanya tersentralisasi dalam satu tempat.
        </p>
        <div class="flex gap-4 mt-8">
          <div class="flex items-center gap-2 text-white/40 text-sm">
            <span class="text-xl">🌳</span> Silsilah
          </div>
          <div class="flex items-center gap-2 text-white/40 text-sm">
            <span class="text-xl">💰</span> Kas
          </div>
          <div class="flex items-center gap-2 text-white/40 text-sm">
            <span class="text-xl">📅</span> Kalender
          </div>
          <div class="flex items-center gap-2 text-white/40 text-sm">
            <span class="text-xl">🗺️</span> Peta
          </div>
        </div>
      </div>

      <p class="text-white/20 text-sm relative">© 2025 Wangsa · Dibangun dengan ❤️</p>
    </div>

    <!-- ── Right: login form ── -->
    <div class="flex-1 flex items-center justify-center p-8 bg-surface">
      <div class="w-full max-w-md">

        <!-- Mobile logo -->
        <div class="flex items-center gap-3 mb-8 lg:hidden">
          <div class="w-9 h-9 rounded-xl bg-navy flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round"
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
            </svg>
          </div>
          <span class="font-bold text-navy text-xl">Wangsa</span>
        </div>

        <div class="mb-9">
          <h1 class="text-3xl font-bold text-navy">Selamat datang 👋</h1>
          <p class="text-warm-gray-500 mt-2 text-lg">Masuk untuk melanjutkan.</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-5">
          <div>
            <label class="form-label">Alamat Email</label>
            <input
              v-model="form.email"
              type="email"
              class="form-input"
              placeholder="nama@keluarga.com"
              required
              autocomplete="email"
            />
          </div>
          <div>
            <label class="form-label">Password</label>
            <input
              v-model="form.password"
              type="password"
              class="form-input"
              placeholder="••••••••"
              required
              autocomplete="current-password"
            />
          </div>

          <div v-if="error" class="form-error">
            <svg class="w-4 h-4 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            {{ error }}
          </div>

          <button type="submit" :disabled="loading" class="btn-primary w-full py-3.5 text-base">
            <span v-if="loading" class="spinner w-4 h-4 border-white/30 border-t-white"></span>
            {{ loading ? 'Memproses...' : 'Masuk' }}
          </button>
        </form>

        <p class="mt-8 text-center text-warm-gray-500 text-base">
          Belum punya akun?
          <RouterLink to="/register" class="text-terra font-semibold hover:underline ml-1">
            Daftar sekarang
          </RouterLink>
        </p>

      </div>
    </div>

  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'

const auth    = useAuthStore()
const router  = useRouter()
const loading = ref(false)
const error   = ref('')
const form    = ref({ email: '', password: '' })

async function handleLogin() {
  loading.value = true
  error.value   = ''
  try {
    await auth.login(form.value)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Login gagal. Periksa email dan password Anda.'
  } finally {
    loading.value = false
  }
}
</script>
