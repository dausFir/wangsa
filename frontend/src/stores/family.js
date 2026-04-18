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
      // Fetch members first to ensure we have member data
      await fetchMembers()
      
      // Try to fetch tree structure
      const { data } = await api.get('/family/tree')
      tree.value = data.data ?? []
      
      // Fallback: if tree is empty but we have members, use members as tree
      if (tree.value.length === 0 && members.value.length > 0) {
        console.warn('Tree endpoint returned empty, using members as fallback')
        // Find root members (those without parent_id)
        const rootMembers = members.value.filter(m => !m.parent_id)
        if (rootMembers.length > 0) {
          tree.value = rootMembers
        } else {
          // If no root members, just use all members as separate trees
          tree.value = [...members.value]
        }
      }
      
      console.log(`Loaded ${tree.value.length} root nodes, ${members.value.length} total members`)
    } catch (e) {
      console.error('Tree fetch error:', e)
      error.value = e.response?.data?.error ?? 'Gagal memuat silsilah'
      
      // Try to load just members as fallback
      try {
        await fetchMembers()
        if (members.value.length > 0) {
          tree.value = members.value.filter(m => !m.parent_id) || [...members.value]
          console.log('Using members fallback after tree error')
        }
      } catch (memberError) {
        console.error('Members fallback also failed:', memberError)
        throw e
      }
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
