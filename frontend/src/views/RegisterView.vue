<template>
  <div class="min-h-screen flex items-center justify-center bg-surface p-8">
    <div class="w-full max-w-md">

      <div class="text-center mb-10">
        <div class="w-14 h-14 rounded-2xl bg-navy flex items-center justify-center mx-auto mb-5 shadow-card">
          <svg class="w-7 h-7 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
          </svg>
        </div>
        <h1 class="text-3xl font-bold text-navy">Buat Akun</h1>
        <p class="text-warm-gray-500 mt-2 text-base max-w-xs mx-auto">
          Pendaftar pertama otomatis menjadi <strong class="text-navy">Super Admin</strong>.
        </p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-5">
        <div>
          <label class="form-label">Nama Lengkap</label>
          <input
            v-model="form.name"
            type="text"
            class="form-input"
            placeholder="Budi Santoso"
            required
            minlength="2"
            autocomplete="name"
          />
        </div>
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
            placeholder="Minimal 8 karakter"
            required
            minlength="8"
            autocomplete="new-password"
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
          {{ loading ? 'Mendaftar...' : 'Daftar Sekarang' }}
        </button>
      </form>

      <p class="mt-8 text-center text-warm-gray-500 text-base">
        Sudah punya akun?
        <RouterLink to="/login" class="text-terra font-semibold hover:underline ml-1">Masuk</RouterLink>
      </p>

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
const form    = ref({ name: '', email: '', password: '' })

async function handleRegister() {
  loading.value = true
  error.value   = ''
  try {
    await auth.register(form.value)
    router.push('/')
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Pendaftaran gagal. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}
</script>
