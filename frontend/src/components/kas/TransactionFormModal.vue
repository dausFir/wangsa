<template>
  <Teleport to="body">
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      style="background: rgba(20,29,39,0.52); backdrop-filter: blur(4px);"
      @mousedown.self="confirmClose"
    >
      <Transition name="modal" appear>
        <div class="bg-white rounded-2xl shadow-modal w-full max-w-md">
          <div class="flex items-center justify-between px-6 py-5 border-b border-warm-gray-100">
            <div class="flex items-center gap-2">
              <h3 class="text-lg font-semibold text-navy">
                {{ isEdit ? 'Edit Transaksi' : 'Catat Transaksi Kas' }}
              </h3>
              <span v-if="isDirty" class="text-xs text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full border border-amber-200">belum disimpan</span>
            </div>
            <button @click="confirmClose" class="btn-ghost p-2 -mr-2">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <form @submit.prevent="handleSubmit" class="px-6 py-5 space-y-5">
            <!-- Type toggle -->
            <div>
              <p class="form-label">Jenis Transaksi</p>
              <div class="flex rounded-xl overflow-hidden border border-warm-gray-200 p-1 gap-1 bg-warm-gray-50">
                <button type="button" @click="form.type = 'income'" :class="['flex-1 py-2.5 text-sm font-medium rounded-lg transition-all', form.type === 'income' ? 'bg-green-500 text-white shadow-soft' : 'text-warm-gray-500 hover:text-navy']">↑  Pemasukan</button>
                <button type="button" @click="form.type = 'expense'" :class="['flex-1 py-2.5 text-sm font-medium rounded-lg transition-all', form.type === 'expense' ? 'bg-red-500 text-white shadow-soft' : 'text-warm-gray-500 hover:text-navy']">↓  Pengeluaran</button>
              </div>
            </div>

            <div>
              <label class="form-label">Jumlah (Rp) <span class="text-terra">*</span></label>
              <div class="relative">
                <span class="absolute left-4 top-1/2 -translate-y-1/2 text-warm-gray-400 text-sm font-medium pointer-events-none">Rp</span>
                <input v-model.number="form.amount" type="number" min="1" step="1000" class="form-input pl-10" placeholder="100.000" required />
              </div>
            </div>

            <div>
              <label class="form-label">Kategori</label>
              <select v-model="form.category_id" class="form-select">
                <option :value="null">— Tanpa Kategori —</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>

            <div>
              <label class="form-label">Keterangan</label>
              <input v-model="form.description" class="form-input" placeholder="Iuran bulan Januari 2025..." />
            </div>

            <div>
              <label class="form-label">Tanggal <span class="text-terra">*</span></label>
              <input v-model="form.date" type="date" class="form-input" required />
            </div>

            <div v-if="error" class="form-error">
              <svg class="w-4 h-4 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              {{ error }}
            </div>

            <div class="flex gap-3 pt-1">
              <button type="button" @click="confirmClose" class="btn-secondary flex-1">Batal</button>
              <button type="submit" :disabled="loading" class="btn-primary flex-1">
                <span v-if="loading" class="spinner w-4 h-4 border-white/30 border-t-white"></span>
                {{ loading ? 'Menyimpan...' : (isEdit ? 'Simpan Perubahan' : 'Simpan Transaksi') }}
              </button>
            </div>
          </form>
        </div>
      </Transition>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useKasStore }   from '@/stores/kas.js'
import { useToast }      from '@/composables/useToast.js'
import { useDirtyGuard } from '@/composables/useDirtyGuard.js'

const props = defineProps({
  categories: { type: Array,  default: () => [] },
  editTx:     { type: Object, default: null },     // null = create mode
})
const emit = defineEmits(['close', 'saved'])

const kas   = useKasStore()
const toast = useToast()

const loading  = ref(false)
const error    = ref('')
const today    = new Date().toISOString().split('T')[0]
const isEdit   = computed(() => !!props.editTx)

const blankForm = () => ({
  type: 'income', amount: null, category_id: null, description: '', date: today,
})
const form = ref(blankForm())

// Pre-fill when editing
watch(() => props.editTx, (tx) => {
  if (tx) {
    form.value = {
      type:        tx.type,
      amount:      tx.amount,
      category_id: tx.category_id ?? null,
      description: tx.description ?? '',
      date:        tx.date,
    }
  } else {
    form.value = blankForm()
  }
  setTimeout(() => setInitial(form.value), 0)
}, { immediate: true })

const { isDirty, setInitial, confirmClose } = useDirtyGuard(form, emit)
onMounted(() => setInitial(form.value))

async function handleSubmit() {
  loading.value = true
  error.value   = ''
  try {
    const payload = { ...form.value, description: form.value.description || null }
    if (isEdit.value) {
      await kas.updateTransaction(props.editTx.id, payload)
      toast.success('Transaksi berhasil diperbarui')
    } else {
      await kas.createTransaction(payload)
      toast.success('Transaksi berhasil dicatat')
    }
    emit('saved')
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Gagal menyimpan transaksi.'
  } finally {
    loading.value = false
  }
}
</script>
