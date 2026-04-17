<template>
  <div class="space-y-6">

    <div v-if="kas.error" class="rounded-xl bg-red-50 border border-red-200 px-4 py-3 text-sm text-red-700 flex items-center gap-2">
      <svg class="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
      </svg>
      {{ kas.error }}
    </div>

    <!-- Summary cards -->
    <div class="grid grid-cols-3 gap-4 lg:gap-5">
      <div class="card text-center">
        <p class="text-xs lg:text-sm text-warm-gray-500 mb-2">Total Pemasukan</p>
        <p class="text-xl lg:text-2xl font-bold text-green-600">{{ fmtRp(kas.summary.income) }}</p>
      </div>
      <div class="card text-center border-2" :class="kas.summary.balance >= 0 ? 'border-terra/25' : 'border-red-200'">
        <p class="text-xs lg:text-sm text-warm-gray-500 mb-2">Saldo</p>
        <p class="text-2xl lg:text-3xl font-bold" :class="kas.summary.balance >= 0 ? 'text-navy' : 'text-red-500'">
          {{ fmtRp(kas.summary.balance) }}
        </p>
      </div>
      <div class="card text-center">
        <p class="text-xs lg:text-sm text-warm-gray-500 mb-2">Total Pengeluaran</p>
        <p class="text-xl lg:text-2xl font-bold text-red-500">{{ fmtRp(kas.summary.expense) }}</p>
      </div>
    </div>

    <!-- Transaction list -->
    <div class="card">
      <div class="flex items-center justify-between mb-6 flex-wrap gap-3">
        <div>
          <h3 class="text-lg font-semibold text-navy">Riwayat Transaksi</h3>
          <p class="text-sm text-warm-gray-400 mt-0.5">
            Menampilkan {{ kas.transactions.length }} transaksi
            <span v-if="kas.canLoadMore"> · ada lebih banyak</span>
          </p>
        </div>
        <button @click="openCreate" class="btn-primary">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
          </svg>
          Catat Transaksi
        </button>
      </div>

      <div v-if="kas.loading" class="flex items-center justify-center py-14 gap-3 text-warm-gray-400">
        <div class="spinner w-6 h-6 border-warm-gray-200 border-t-warm-gray-400"></div>
        <span>Memuat transaksi...</span>
      </div>

      <div v-else-if="!kas.transactions.length" class="text-center py-16">
        <div class="text-5xl mb-4">💸</div>
        <p class="text-warm-gray-500">Belum ada transaksi. Mulai catat kas keluarga!</p>
      </div>

      <div v-else class="divide-y divide-warm-gray-50">
        <div
          v-for="tx in kas.transactions"
          :key="tx.id"
          class="flex items-center gap-4 py-3.5 px-2 rounded-xl hover:bg-warm-gray-50 transition-colors group"
        >
          <div :class="['w-10 h-10 rounded-xl flex items-center justify-center text-base flex-shrink-0 font-medium', tx.type === 'income' ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-500']">
            {{ tx.type === 'income' ? '↑' : '↓' }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="font-medium text-navy text-sm leading-snug truncate">
              {{ tx.description || tx.category_name || 'Transaksi' }}
            </p>
            <div class="flex items-center gap-2 mt-0.5 flex-wrap">
              <span class="text-xs text-warm-gray-400">{{ fmtDate(tx.date) }}</span>
              <span class="text-warm-gray-200">·</span>
              <span class="text-xs text-warm-gray-400">{{ tx.creator_name ?? 'Sistem' }}</span>
              <span v-if="tx.category_name" class="badge-gray text-[11px] px-1.5 py-px">{{ tx.category_name }}</span>
            </div>
          </div>
          <p :class="['text-base font-bold flex-shrink-0 tabular-nums', tx.type === 'income' ? 'text-green-600' : 'text-red-500']">
            {{ tx.type === 'income' ? '+' : '−' }}{{ fmtRp(tx.amount) }}
          </p>

          <!-- Actions: edit + delete (visible on hover) -->
          <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
            <button
              @click="openEdit(tx)"
              class="btn-ghost p-1.5 text-warm-gray-400 hover:text-navy hover:bg-warm-gray-100"
              title="Edit transaksi"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
            </button>
            <button
              v-if="auth.isSuperAdmin"
              @click="handleDelete(tx.id)"
              class="btn-ghost p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50"
              title="Hapus transaksi"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
        </div>
      </div>

      <div v-if="kas.canLoadMore || kas.loadingMore" class="mt-4 pt-4 border-t border-warm-gray-50 flex justify-center">
        <button @click="kas.loadMore()" :disabled="kas.loadingMore" class="btn-secondary flex items-center gap-2">
          <span v-if="kas.loadingMore" class="spinner w-4 h-4 border-warm-gray-300 border-t-warm-gray-600"></span>
          {{ kas.loadingMore ? 'Memuat...' : 'Muat lebih banyak' }}
        </button>
      </div>
    </div>

    <TransactionFormModal
      v-if="showForm"
      :categories="kas.categories"
      :edit-tx="editingTx"
      @close="closeForm"
      @saved="closeForm"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useKasStore }  from '@/stores/kas.js'
import { useAuthStore } from '@/stores/auth.js'
import { fmtRp, fmtDate } from '@/utils/format.js'
import { useToast }      from '@/composables/useToast.js'
import TransactionFormModal from '@/components/kas/TransactionFormModal.vue'

const kas       = useKasStore()
const auth      = useAuthStore()
const toast     = useToast()
const showForm  = ref(false)
const editingTx = ref(null)

function openCreate() { editingTx.value = null; showForm.value = true }
function openEdit(tx) { editingTx.value = tx;   showForm.value = true }
function closeForm()  { showForm.value = false; editingTx.value = null }

async function handleDelete(id) {
  if (confirm('Hapus transaksi ini? Tindakan tidak dapat dibatalkan.')) {
    await kas.deleteTransaction(id)
    toast.success('Transaksi berhasil dihapus')
  }
}

onMounted(() => kas.fetchAll())
</script>
