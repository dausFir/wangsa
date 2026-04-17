import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios.js'

export const useFamilyStore = defineStore('family', () => {
  const tree    = ref([])
  const members = ref([])
  const loading = ref(false)
  const error   = ref(null)

  async function fetchTree() {
    loading.value = true
    error.value   = null
    try {
      const { data } = await api.get('/family/tree')
      tree.value = data.data ?? []
    } catch (e) {
      error.value = e.response?.data?.error ?? 'Gagal memuat silsilah'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchMembers() {
    const { data } = await api.get('/family/members')
    members.value = data.data ?? []
  }

  async function createMember(payload) {
    const { data } = await api.post('/family/members', payload)
    await Promise.all([fetchTree(), fetchMembers()])
    return data.data
  }

  async function updateMember(id, payload) {
    const { data } = await api.put(`/family/members/${id}`, payload)
    await Promise.all([fetchTree(), fetchMembers()])
    return data.data
  }

  async function deleteMember(id) {
    await api.delete(`/family/members/${id}`)
    await Promise.all([fetchTree(), fetchMembers()])
  }

  async function createMarriage(payload) {
    const { data } = await api.post('/family/marriages', payload)
    await fetchTree()
    return data.data
  }

  async function deleteMarriage(id) {
    await api.delete(`/family/marriages/${id}`)
    await fetchTree()
  }

  return {
    tree, members, loading, error,
    fetchTree, fetchMembers,
    createMember, updateMember, deleteMember,
    createMarriage, deleteMarriage
  }
})
