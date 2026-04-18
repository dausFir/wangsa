import { ref, computed, reactive } from 'vue'
import { defineStore } from 'pinia'
import api from '@/api/axios.js'

export const useNotesStore = defineStore('notes', () => {
  // State
  const notes = ref([])
  const categories = ref([])
  const loading = ref(false)
  const error = ref(null)
  
  // Form modal state
  const isModalOpen = ref(false)
  const editingNote = ref(null)
  const modalTitle = computed(() => 
    editingNote.value ? 'Edit Catatan' : 'Tambah Catatan'
  )

  // Form data
  const form = reactive({
    title: '',
    content: '',
    category: '',
    isPinned: false
  })

  // Getters
  const pinnedNotes = computed(() => 
    notes.value.filter(note => note.is_pinned)
  )
  
  const unpinnedNotes = computed(() => 
    notes.value.filter(note => !note.is_pinned)
  )
  
  const notesByCategory = computed(() => {
    return notes.value.reduce((groups, note) => {
      const category = note.category || 'Tanpa Kategori'
      if (!groups[category]) groups[category] = []
      groups[category].push(note)
      return groups
    }, {})
  })

  // Actions
  async function fetchNotes(categoryFilter = '') {
    loading.value = true
    error.value = null
    
    try {
      const params = categoryFilter ? { category: categoryFilter } : {}
      const response = await api.get('/notes', { params })
      notes.value = response.data.data || []
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal memuat catatan'
      console.error('Fetch notes error:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchCategories() {
    try {
      const response = await api.get('/notes/categories')
      categories.value = response.data.data || []
    } catch (err) {
      console.error('Fetch categories error:', err)
    }
  }

  async function createNote() {
    loading.value = true
    error.value = null
    
    try {
      const response = await api.post('/notes', {
        title: form.title,
        content: form.content,
        category: form.category || null,
        is_pinned: form.isPinned
      })
      
      // Add to local state
      notes.value.unshift(response.data.data)
      
      // Refresh categories if new category added
      if (form.category && !categories.value.includes(form.category)) {
        await fetchCategories()
      }
      
      // Reset form and close modal
      resetForm()
      closeModal()
      
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal membuat catatan'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateNote() {
    loading.value = true
    error.value = null
    
    try {
      const response = await api.put(`/notes/${editingNote.value.id}`, {
        title: form.title,
        content: form.content,
        category: form.category || null,
        is_pinned: form.isPinned
      })
      
      // Update local state
      const index = notes.value.findIndex(n => n.id === editingNote.value.id)
      if (index !== -1) {
        notes.value[index] = response.data.data
      }
      
      // Refresh categories if changed
      await fetchCategories()
      
      // Reset form and close modal
      resetForm()
      closeModal()
      
      return response.data.data
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal mengupdate catatan'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteNote(noteId) {
    if (!confirm('Yakin ingin menghapus catatan ini?')) return
    
    loading.value = true
    error.value = null
    
    try {
      await api.delete(`/notes/${noteId}`)
      
      // Remove from local state
      const index = notes.value.findIndex(n => n.id === noteId)
      if (index !== -1) {
        notes.value.splice(index, 1)
      }
      
      // Refresh categories
      await fetchCategories()
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal menghapus catatan'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function togglePin(noteId) {
    const note = notes.value.find(n => n.id === noteId)
    if (!note) return
    
    try {
      const response = await api.put(`/notes/${noteId}`, {
        ...note,
        is_pinned: !note.is_pinned
      })
      
      // Update local state
      const index = notes.value.findIndex(n => n.id === noteId)
      if (index !== -1) {
        notes.value[index] = response.data.data
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal mengubah status pin'
      throw err
    }
  }

  // Modal actions
  function openCreateModal() {
    editingNote.value = null
    resetForm()
    isModalOpen.value = true
  }

  function openEditModal(note) {
    editingNote.value = note
    form.title = note.title
    form.content = note.content
    form.category = note.category || ''
    form.isPinned = note.is_pinned
    isModalOpen.value = true
  }

  function closeModal() {
    isModalOpen.value = false
    editingNote.value = null
    resetForm()
  }

  function resetForm() {
    form.title = ''
    form.content = ''
    form.category = ''
    form.isPinned = false
  }

  async function saveNote() {
    if (editingNote.value) {
      await updateNote()
    } else {
      await createNote()
    }
  }

  // Initialize
  function init() {
    fetchNotes()
    fetchCategories()
  }

  return {
    // State
    notes,
    categories,
    loading,
    error,
    pinnedNotes,
    unpinnedNotes,
    notesByCategory,
    
    // Modal state
    isModalOpen,
    editingNote,
    modalTitle,
    form,
    
    // Actions
    fetchNotes,
    fetchCategories,
    createNote,
    updateNote,
    deleteNote,
    togglePin,
    saveNote,
    
    // Modal actions
    openCreateModal,
    openEditModal,
    closeModal,
    resetForm,
    
    // Init
    init
  }
})