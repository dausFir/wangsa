<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto">
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black bg-opacity-50 transition-opacity" @click="handleClose"></div>
    
    <!-- Modal Container -->
    <div class="flex min-h-full items-center justify-center p-4">
      <div class="relative bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-hidden">
        
        <!-- Header -->
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">{{ title }}</h3>
          <button
            @click="handleClose"
            class="text-gray-400 hover:text-gray-600 p-1 rounded-full hover:bg-gray-100"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <!-- Form Content -->
        <form @submit.prevent="handleSave" class="flex flex-col h-full">
          <div class="flex-1 p-6 space-y-4 overflow-y-auto">
            
            <!-- Title Input -->
            <div>
              <label for="note-title" class="block text-sm font-medium text-gray-700 mb-1">
                Judul Catatan *
              </label>
              <input
                id="note-title"
                v-model="form.title"
                type="text"
                required
                placeholder="Masukkan judul catatan..."
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>

            <!-- Category Input -->
            <div>
              <label for="note-category" class="block text-sm font-medium text-gray-700 mb-1">
                Kategori
              </label>
              <div class="relative">
                <input
                  id="note-category"
                  v-model="form.category"
                  type="text"
                  list="categories-list"
                  placeholder="Pilih atau buat kategori baru..."
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
                <datalist id="categories-list">
                  <option v-for="category in categories" :key="category" :value="category"></option>
                </datalist>
              </div>
            </div>

            <!-- Pin Checkbox -->
            <div class="flex items-center">
              <input
                id="note-pinned"
                v-model="form.isPinned"
                type="checkbox"
                class="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <label for="note-pinned" class="ml-2 text-sm text-gray-700 flex items-center">
                <svg class="w-4 h-4 mr-1 text-yellow-500" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M17 4v7l2 3v2h-6v5l-1 1-1-1v-5H5v-2l2-3V4c0-1.1.9-2 2-2h6c1.1 0 2 .9 2 2zM9 4v7.75L7.5 14h9L15 11.75V4H9z"/>
                </svg>
                Pin catatan ini
              </label>
            </div>

            <!-- Content Editor -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Isi Catatan *
              </label>
              <RichTextEditor
                v-model="form.content"
                placeholder="Tulis isi catatan di sini..."
                :height="300"
              />
            </div>

          </div>

          <!-- Footer Actions -->
          <div class="flex items-center justify-end space-x-3 p-6 border-t border-gray-200 bg-gray-50">
            <button
              type="button"
              @click="handleClose"
              class="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Batal
            </button>
            <button
              type="submit"
              :disabled="loading || !form.title.trim() || !form.content.trim()"
              class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center space-x-2"
            >
              <div v-if="loading" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              <span>{{ note ? 'Simpan Perubahan' : 'Buat Catatan' }}</span>
            </button>
          </div>
        </form>

      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, watch } from 'vue'
import { useNotesStore } from '@/stores/notes.js'
import RichTextEditor from '@/components/ui/RichTextEditor.vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Tambah Catatan'
  },
  note: {
    type: Object,
    default: null
  },
  categories: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'save'])

const notesStore = useNotesStore()
const { form, loading } = notesStore

// Watch for note prop changes to populate form
watch(() => props.note, (newNote) => {
  if (newNote) {
    form.title = newNote.title
    form.content = newNote.content
    form.category = newNote.category || ''
    form.isPinned = newNote.is_pinned
  }
}, { immediate: true })

function handleClose() {
  emit('close')
}

async function handleSave() {
  try {
    await emit('save')
  } catch (error) {
    console.error('Error saving note:', error)
  }
}

// Close modal on Escape key
function handleKeydown(event) {
  if (event.key === 'Escape' && props.isOpen) {
    handleClose()
  }
}

// Add event listener when modal opens
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})
</script>