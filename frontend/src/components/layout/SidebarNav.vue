<template>
  <aside :class="[
    'flex-shrink-0 flex flex-col bg-navy h-screen overflow-y-auto transition-transform duration-300 ease-in-out',
    'fixed inset-y-0 left-0 z-40 w-64',
    'lg:relative lg:translate-x-0 lg:w-60',
    isOpen ? 'translate-x-0' : '-translate-x-full'
  ]">
    <!-- Logo -->
    <div class="px-4 lg:px-5 py-6 border-b border-white/10 flex-shrink-0 flex items-center justify-between">
      <div class="flex items-center gap-3 min-w-0">
        <div class="w-9 h-9 rounded-xl bg-terra flex items-center justify-center flex-shrink-0 shadow-soft">
          <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round"
              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
          </svg>
        </div>
        <div class="min-w-0">
          <p class="font-bold text-white text-[15px] leading-tight truncate">Wangsa</p>
          <p class="text-xs text-white/40 leading-tight truncate">Pusat Keluarga Besar</p>
        </div>
      </div>
      <!-- Close button — mobile only -->
      <button
        class="lg:hidden w-8 h-8 flex items-center justify-center rounded-lg text-white/50 hover:text-white hover:bg-white/10 transition-colors flex-shrink-0"
        @click="sidebarStore.close()"
        aria-label="Tutup menu"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
    </div>

    <!-- Navigation links -->
    <nav class="flex-1 px-3 py-5 space-y-0.5">
      <RouterLink
        v-for="item in navItems"
        :key="item.name"
        :to="item.to"
        custom
        v-slot="{ isActive, navigate }"
      >
        <button
          @click="handleNav(navigate)"
          :class="[
            'w-full flex items-center gap-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-150 text-left',
            isActive
              ? 'bg-terra text-white shadow-soft'
              : 'text-white/55 hover:bg-white/[0.07] hover:text-white'
          ]"
        >
          <span class="text-[17px] leading-none w-5 flex-shrink-0 text-center">{{ item.emoji }}</span>
          <span class="truncate">{{ item.name }}</span>
          <span v-if="isActive" class="ml-auto w-1.5 h-1.5 rounded-full bg-white/50 flex-shrink-0"></span>
        </button>
      </RouterLink>
    </nav>

    <!-- Footer -->
    <div class="px-5 py-4 border-t border-white/10 flex-shrink-0">
      <p class="text-xs text-white/20 text-center tracking-wide">Wangsa v1.1.0</p>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useSidebarStore } from '@/stores/sidebar.js'

const sidebarStore = useSidebarStore()
const isOpen = computed(() => sidebarStore.isOpen)

const navItems = [
  { name: 'Dashboard', to: '/',         emoji: '🏠' },
  { name: 'Silsilah',  to: '/silsilah', emoji: '🌳' },
  { name: 'Kas',       to: '/kas',      emoji: '💰' },
  { name: 'Peta',      to: '/peta',     emoji: '🗺️' },
  { name: 'Kalender',  to: '/kalender', emoji: '📅' },
  { name: 'Catatan',   to: '/catatan',  emoji: '📝' },
]

function handleNav(navigate) {
  navigate()
  sidebarStore.close() // auto-close on mobile after nav
}
</script>
