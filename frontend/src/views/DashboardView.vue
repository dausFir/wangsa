<template>
  <div class="space-y-7">

    <!-- ── Welcome Banner ── -->
    <div class="bg-navy rounded-2xl p-8 text-white relative overflow-hidden">
      <div class="absolute right-0 top-0 bottom-0 w-56 opacity-[0.04] pointer-events-none" aria-hidden="true">
        <svg viewBox="0 0 200 300" class="w-full h-full">
          <circle cx="100" cy="60"  r="120" fill="white"/>
          <circle cx="40"  cy="260" r="80"  fill="white"/>
        </svg>
      </div>
      <p class="text-white/50 text-sm mb-1.5 relative">Selamat datang kembali,</p>
      <h2 class="text-3xl font-bold mb-2 relative">{{ auth.user?.name }} 👋</h2>
      <p class="text-white/50 text-base relative max-w-md">
        Kelola keluarga besar Anda dengan mudah dan terorganisir dari satu tempat.
      </p>
    </div>

    <!-- ── Stats Grid ── -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 lg:gap-5">
      <div
        v-for="s in stats"
        :key="s.label"
        class="card hover:shadow-card hover:-translate-y-0.5 transition-all duration-200"
      >
        <div :class="['w-11 h-11 rounded-xl flex items-center justify-center mb-4 text-2xl', s.bg]">
          {{ s.icon }}
        </div>
        <p class="text-2xl lg:text-3xl font-bold text-navy leading-none mb-1">{{ s.value }}</p>
        <p class="text-sm text-warm-gray-500">{{ s.label }}</p>
      </div>
    </div>

    <!-- ── Quick Actions ── -->
    <div class="card">
      <h3 class="text-base font-semibold text-navy mb-5">Akses Cepat</h3>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
        <RouterLink
          v-for="a in quickActions"
          :key="a.name"
          :to="a.to"
          class="flex flex-col items-center gap-3 p-5 rounded-xl bg-warm-gray-50
                 border border-warm-gray-100 hover:bg-terra/5 hover:border-terra/20
                 transition-all duration-200 text-center group cursor-pointer"
        >
          <span class="text-3xl">{{ a.icon }}</span>
          <span class="text-sm font-medium text-warm-gray-700 group-hover:text-terra transition-colors leading-tight">
            {{ a.name }}
          </span>
        </RouterLink>
      </div>
    </div>

    <!-- ── Kas Highlight ── -->
    <div v-if="loaded" class="card border-l-4 border-l-terra flex flex-wrap items-center justify-between gap-4">
      <div>
        <p class="text-sm text-warm-gray-500 mb-1.5">Saldo Kas Keluarga</p>
        <p
          class="text-3xl font-bold"
          :class="kas.summary.balance >= 0 ? 'text-navy' : 'text-red-500'"
        >
          {{ fmtRp(kas.summary.balance) }}
        </p>
        <p class="text-xs text-warm-gray-400 mt-1">
          Pemasukan {{ fmtRp(kas.summary.income) }} · Pengeluaran {{ fmtRp(kas.summary.expense) }}
        </p>
      </div>
      <RouterLink to="/kas" class="btn-secondary text-sm flex-shrink-0">
        Lihat Rincian →
      </RouterLink>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore }   from '@/stores/auth.js'
import { useFamilyStore } from '@/stores/family.js'
import { useKasStore }    from '@/stores/kas.js'
import { useEventsStore } from '@/stores/events.js'
import { fmtRp }          from '@/utils/format.js'

const auth   = useAuthStore()
const family = useFamilyStore()
const kas    = useKasStore()
const events = useEventsStore()
const loaded = ref(false)

const stats = ref([
  { label: 'Anggota Keluarga', value: '—', icon: '👥', bg: 'bg-blue-50'   },
  { label: 'Saldo Kas',        value: '—', icon: '💰', bg: 'bg-green-50'  },
  { label: 'Acara Tercatat',   value: '—', icon: '📅', bg: 'bg-purple-50' },
  { label: 'Generasi',         value: '—', icon: '🌳', bg: 'bg-orange-50' },
])

const quickActions = [
  { name: 'Lihat Silsilah', to: '/silsilah', icon: '🌳' },
  { name: 'Tambah Kas',     to: '/kas',      icon: '💳' },
  { name: 'Buat Acara',     to: '/kalender', icon: '🎉' },
  { name: 'Peta Keluarga',  to: '/peta',     icon: '📍' },
]

function countGenerations(nodes, depth = 1) {
  if (!nodes?.length) return depth - 1
  return Math.max(...nodes.map(n => countGenerations(n.children, depth + 1)))
}

onMounted(async () => {
  await Promise.all([
    family.fetchTree(),
    family.fetchMembers(),
    kas.fetchAll(),
    events.fetchEvents(),
  ])
  stats.value[0].value = family.members.length || '0'
  stats.value[1].value = fmtRp(kas.summary.balance)
  stats.value[2].value = events.events.length || '0'   // was never populated before
  stats.value[3].value = countGenerations(family.tree) || '0'
  loaded.value = true
})
</script>
