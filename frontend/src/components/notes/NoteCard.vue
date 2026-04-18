<template>
  <div :class="['bg-white rounded-lg border shadow-sm p-4 hover:shadow-md transition-shadow',
    note.is_pinned ? 'border-yellow-300 bg-yellow-50' : 'border-gray-200']">
    
    <!-- Header -->
    <div class="flex items-start justify-between mb-3">
      <div class="flex-1">
        <h3 class="font-medium text-gray-900 mb-1">{{ note.title }}</h3>
        <div class="flex items-center space-x-2 text-sm text-gray-500">
          <span v-if="note.category" 
            class="inline-flex items-center px-2 py-1 rounded-full text-xs bg-terra/10 text-terra-dark">
            {{ note.category }}
          </span>
          <span>{{ formatDate(note.updated_at) }}</span>
        </div>
      </div>
      
      <!-- Actions Menu -->
      <div class="relative">
        <button
          @click="toggleMenu"
          class="p-1 text-gray-400 hover:text-gray-600 rounded-full hover:bg-gray-100"
        >
          <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z"/>
          </svg>
        </button>
        
        <!-- Dropdown Menu -->
        <div v-if="showMenu" 
          class="absolute right-0 mt-1 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-10">
          
          <button
            @click="handlePin"
            class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center"
          >
            <svg v-if="!note.is_pinned" class="w-4 h-4 mr-2 text-yellow-500" fill="currentColor" viewBox="0 0 24 24">
              <path d="M17 4v7l2 3v2h-6v5l-1 1-1-1v-5H5v-2l2-3V4c0-1.1.9-2 2-2h6c1.1 0 2 .9 2 2zM9 4v7.75L7.5 14h9L15 11.75V4H9z"/>
            </svg>
            <svg v-else class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"/>
            </svg>
            {{ note.is_pinned ? 'Lepas Pin' : 'Pin Catatan' }}
          </button>
          
          <button
            @click="handleEdit"
            class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
            Edit Catatan
          </button>
          
          <button
            @click="handleDelete"
            class="w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-red-50 flex items-center"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
            Hapus Catatan
          </button>
        </div>
      </div>
    </div>

    <!-- Content Preview -->
    <div class="prose prose-sm max-w-none text-gray-600 line-clamp-3" 
         v-html="truncateHtml(note.content)">
    </div>
    
    <!-- Expand/Collapse for longer content -->
    <button
      v-if="isContentLong"
      @click="expanded = !expanded"
      class="text-blue-600 text-sm mt-2 hover:underline"
    >
      {{ expanded ? 'Tampilkan lebih sedikit' : 'Tampilkan selengkapnya' }}
    </button>
    
    <!-- Full content when expanded -->
    <div v-if="expanded" class="prose prose-sm max-w-none mt-3 pt-3 border-t border-gray-200" 
         v-html="note.content">
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { format } from '@/utils/format.js'

const props = defineProps({
  note: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['edit', 'delete', 'togglePin'])

const showMenu = ref(false)
const expanded = ref(false)

const isContentLong = computed(() => {
  // Check if content is longer than ~150 characters (rough estimate)
  return props.note.content && props.note.content.length > 150
})

function toggleMenu() {
  showMenu.value = !showMenu.value
}

function closeMenu() {
  showMenu.value = false
}

function handlePin() {
  emit('togglePin', props.note.id)
  closeMenu()
}

function handleEdit() {
  emit('edit', props.note)
  closeMenu()
}

function handleDelete() {
  emit('delete', props.note.id)
  closeMenu()
}

function formatDate(dateString) {
  return format.relativeTime(dateString)
}

function truncateHtml(html) {
  if (!html) return ''
  
  // Strip HTML tags for length calculation
  const textOnly = html.replace(/<[^>]*>/g, '')
  
  if (textOnly.length <= 150) {
    return html
  }
  
  // Try to truncate at word boundary within HTML
  const words = textOnly.split(' ')
  let truncated = ''
  let charCount = 0
  
  for (const word of words) {
    if (charCount + word.length > 150) break
    truncated += (truncated ? ' ' : '') + word
    charCount += word.length + 1
  }
  
  // Find the truncated text in HTML and cut there
  const index = html.indexOf(truncated)
  if (index !== -1) {
    return html.substring(0, index + truncated.length) + '...'
  }
  
  return html.substring(0, 150) + '...'
}

// Close menu when clicking outside
function handleClickOutside(event) {
  if (!event.target.closest('.relative')) {
    closeMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>