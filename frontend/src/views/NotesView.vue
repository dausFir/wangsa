<template>
  <div class="notes-view">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center space-x-3">
        <svg class="w-6 h-6 text-terra" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
        <h1 class="text-2xl font-bold text-gray-800">Catatan Keluarga</h1>
      </div>
      <button
        @click="openCreateModal"
        class="bg-terra hover:bg-terra-dark text-white px-4 py-2 rounded-lg flex items-center space-x-2 transition-colors"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        <span>Tambah Catatan</span>
      </button>
    </div>

    <!-- Category Filter -->
    <div class="mb-6">
      <div class="flex flex-wrap gap-2">
        <button
          @click="fetchNotes('')"
          :class="['px-3 py-1 rounded-full text-sm transition-colors',
            selectedCategory === '' 
              ? 'bg-terra text-white' 
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300']"
        >
          Semua
        </button>
        <button
          v-for="category in categories"
          :key="category"
          @click="filterByCategory(category)"
          :class="['px-3 py-1 rounded-full text-sm transition-colors',
            selectedCategory === category
              ? 'bg-terra text-white'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300']"
        >
          {{ category }}
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-terra"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-red-600 text-center py-12">
      {{ error }}
    </div>

    <!-- Empty State -->
    <div v-else-if="notes.length === 0" class="text-center py-12 text-gray-500">
      <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
      </svg>
      <p class="text-lg font-medium">Belum ada catatan</p>
      <p class="text-sm">Buat catatan pertama untuk menyimpan informasi keluarga</p>
    </div>

    <!-- Notes List -->
    <div v-else class="space-y-4">
      <!-- Pinned Notes -->
      <div v-if="pinnedNotes.length > 0">
        <h2 class="text-sm font-medium text-gray-600 mb-3 flex items-center">
          <svg class="w-4 h-4 mr-2 text-yellow-500" fill="currentColor" viewBox="0 0 24 24">
            <path d="M17 4v7l2 3v2h-6v5l-1 1-1-1v-5H5v-2l2-3V4c0-1.1.9-2 2-2h6c1.1 0 2 .9 2 2zM9 4v7.75L7.5 14h9L15 11.75V4H9z"/>
          </svg>
          Catatan Terpin
        </h2>
        <div class="space-y-3 mb-6">
          <NoteCard
            v-for="note in pinnedNotes"
            :key="note.id"
            :note="note"
            @edit="openEditModal"
            @delete="deleteNote"
            @togglePin="togglePin"
          />
        </div>
      </div>

      <!-- Regular Notes -->
      <div v-if="unpinnedNotes.length > 0">
        <h2 v-if="pinnedNotes.length > 0" class="text-sm font-medium text-gray-600 mb-3">
          Catatan Lainnya
        </h2>
        <div class="space-y-3">
          <NoteCard
            v-for="note in unpinnedNotes"
            :key="note.id"
            :note="note"
            @edit="openEditModal"
            @delete="deleteNote"
            @togglePin="togglePin"
          />
        </div>
      </div>
    </div>

    <!-- Note Form Modal -->
    <NoteFormModal
      :isOpen="isModalOpen"
      :title="modalTitle"
      :note="editingNote"
      :categories="categories"
      @close="closeModal"
      @save="saveNote"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useNotesStore } from '@/stores/notes.js'
import NoteCard from '@/components/notes/NoteCard.vue'
import NoteFormModal from '@/components/notes/NoteFormModal.vue'

const notesStore = useNotesStore()
const selectedCategory = ref('')

// Destructure reactive state with storeToRefs
const {
  notes,
  categories,
  loading,
  error,
  pinnedNotes,
  unpinnedNotes,
  isModalOpen,
  modalTitle,
  editingNote
} = storeToRefs(notesStore)

// Destructure actions directly (they don't need reactivity)
const {
  fetchNotes,
  deleteNote,
  togglePin,
  openCreateModal,
  openEditModal,
  closeModal,
  saveNote,
  init
} = notesStore

function filterByCategory(category) {
  selectedCategory.value = category
  fetchNotes(category)
}

onMounted(() => {
  init()
})
</script>