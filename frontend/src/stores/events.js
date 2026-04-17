import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios.js'

export const useEventsStore = defineStore('events', () => {
  const events  = ref([])
  const loading = ref(false)
  const error   = ref(null)

  async function fetchEvents(from = '', to = '') {
    loading.value = true
    error.value   = null
    try {
      const { data } = await api.get('/events', { params: { from, to } })
      events.value = data.data ?? []
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal memuat acara'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function createEvent(payload) {
    try {
      const { data } = await api.post('/events', payload)
      await fetchEvents()
      return data.data
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal membuat acara'
      throw e
    }
  }

  async function updateEvent(id, payload) {
    try {
      const { data } = await api.put(`/events/${id}`, payload)
      await fetchEvents()
      return data.data
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal memperbarui acara'
      throw e
    }
  }

  async function deleteEvent(id) {
    try {
      await api.delete(`/events/${id}`)
      await fetchEvents()
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal menghapus acara'
      throw e
    }
  }

  return { events, loading, error, fetchEvents, createEvent, updateEvent, deleteEvent }
})
