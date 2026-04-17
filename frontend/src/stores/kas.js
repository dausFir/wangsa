import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api/axios.js'

const PAGE_SIZE = 20

export const useKasStore = defineStore('kas', () => {
  const transactions = ref([])
  const categories   = ref([])
  const summary      = ref({ income: 0, expense: 0, balance: 0 })
  const loading      = ref(false)
  const loadingMore  = ref(false)
  const error        = ref(null)
  const offset       = ref(0)
  const _lastPageFull = ref(false)  // true if last page had exactly PAGE_SIZE items
  const canLoadMore   = computed(() => _lastPageFull.value && !loadingMore.value)

  async function fetchAll() {
    loading.value = true
    error.value   = null
    offset.value  = 0
    try {
      const [txRes, catRes, sumRes] = await Promise.all([
        api.get('/kas/transactions', { params: { limit: PAGE_SIZE, offset: 0 } }),
        api.get('/kas/categories'),
        api.get('/kas/summary')
      ])
      const txData = txRes.data.data ?? []
      transactions.value  = txData
      categories.value    = catRes.data.data ?? []
      summary.value       = sumRes.data.data ?? { income: 0, expense: 0, balance: 0 }
      offset.value        = txData.length
      _lastPageFull.value = txData.length === PAGE_SIZE
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal memuat data kas'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function loadMore() {
    if (!canLoadMore.value) return
    loadingMore.value = true
    try {
      const { data } = await api.get('/kas/transactions', {
        params: { limit: PAGE_SIZE, offset: offset.value }
      })
      const newTxs = data.data ?? []
      transactions.value  = [...transactions.value, ...newTxs]
      offset.value       += newTxs.length
      _lastPageFull.value = newTxs.length === PAGE_SIZE
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal memuat transaksi'
    } finally {
      loadingMore.value = false
    }
  }

  async function createTransaction(payload) {
    try {
      const { data } = await api.post('/kas/transactions', payload)
      await fetchAll()
      return data.data
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal mencatat transaksi'
      throw e
    }
  }

  async function updateTransaction(id, payload) {
    try {
      const { data } = await api.put(`/kas/transactions/${id}`, payload)
      await fetchAll()
      return data.data
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal memperbarui transaksi'
      throw e
    }
  }

  async function deleteTransaction(id) {
    try {
      await api.delete(`/kas/transactions/${id}`)
      await fetchAll()
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal menghapus transaksi'
      throw e
    }
  }

  return {
    transactions, categories, summary, loading, loadingMore, error, canLoadMore,
    fetchAll, loadMore, createTransaction, updateTransaction, deleteTransaction
  }
})
