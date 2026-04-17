<template>
  <div class="flex h-screen overflow-hidden bg-surface">

    <!-- Mobile backdrop overlay -->
    <Transition name="fade">
      <div
        v-if="sidebar.isOpen"
        class="fixed inset-0 z-30 bg-black/50 lg:hidden"
        @click="sidebar.close()"
        aria-hidden="true"
      />
    </Transition>

    <SidebarNav />

    <div class="flex-1 flex flex-col overflow-hidden min-w-0">
      <!-- Top Bar -->
      <header class="flex items-center justify-between px-4 lg:px-8 py-4 bg-white border-b border-warm-gray-100 shadow-soft flex-shrink-0">
        <div class="flex items-center gap-3 min-w-0">
          <!-- Hamburger — mobile only -->
          <button
            class="lg:hidden w-9 h-9 flex items-center justify-center rounded-xl text-warm-gray-400 hover:bg-warm-gray-100 transition-colors flex-shrink-0"
            @click="sidebar.toggle()"
            aria-label="Buka menu"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16"/>
            </svg>
          </button>
          <div class="min-w-0">
            <h1 class="text-base lg:text-xl font-semibold text-navy truncate">{{ pageTitle }}</h1>
            <p class="text-xs text-warm-gray-400 mt-0.5 hidden sm:block truncate">{{ pageSubtitle }}</p>
          </div>
        </div>

        <div class="flex items-center gap-2 lg:gap-3 flex-shrink-0">
          <div class="flex items-center gap-2.5 px-3 py-2 rounded-xl bg-warm-gray-50 border border-warm-gray-100">
            <div class="w-8 h-8 rounded-full bg-terra/15 flex items-center justify-center text-terra font-semibold text-sm flex-shrink-0 uppercase">
              {{ initials }}
            </div>
            <div class="hidden sm:block">
              <p class="text-sm font-medium text-navy leading-tight">{{ auth.user?.name }}</p>
              <p class="text-xs text-warm-gray-500 leading-tight">
                {{ auth.user?.role === 'super_admin' ? 'Super Admin' : 'Member' }}
              </p>
            </div>
          </div>
          <button
            @click="handleLogout"
            title="Keluar"
            class="w-9 h-9 flex items-center justify-center rounded-xl text-warm-gray-400 hover:bg-red-50 hover:text-red-500 transition-colors"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round"
                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
            </svg>
          </button>
        </div>
      </header>

      <!-- Scrollable content -->
      <main class="flex-1 overflow-y-auto p-4 lg:p-8">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter, RouterView } from 'vue-router'
import { useAuthStore }   from '@/stores/auth.js'
import { useSidebarStore } from '@/stores/sidebar.js'
import SidebarNav from './SidebarNav.vue'

const auth    = useAuthStore()
const sidebar = useSidebarStore()
const route   = useRoute()
const router  = useRouter()

const initials = computed(() =>
  (auth.user?.name ?? 'FH')
    .split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
)

const pageMeta = {
  Dashboard: { title: 'Dashboard',      subtitle: 'Selamat datang di Wangsa' },
  Silsilah:  { title: 'Pohon Silsilah', subtitle: 'Visualisasi dan kelola silsilah keluarga besar' },
  Kas:       { title: 'Kas Keluarga',   subtitle: 'Pencatatan keuangan bersama keluarga' },
  Peta:      { title: 'Peta Domisili',  subtitle: 'Sebaran tempat tinggal anggota keluarga' },
  Kalender:  { title: 'Kalender Acara', subtitle: 'Jadwal dan agenda keluarga' },
}

const pageTitle    = computed(() => pageMeta[route.name]?.title    ?? 'Wangsa')
const pageSubtitle = computed(() => pageMeta[route.name]?.subtitle ?? '')

async function handleLogout() {
  await auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.25s ease; }
.fade-enter-from, .fade-leave-to       { opacity: 0; }
</style>
