<template>
  <div>
    <label v-if="label" :for="editorId" class="form-label">{{ label }}</label>
    <div 
      :id="editorId" 
      class="quill-editor border border-warm-gray-200 rounded-xl"
      style="min-height: 120px;"
    ></div>
    <p v-if="error" class="text-xs text-red-500 mt-1">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import Quill from 'quill'
import 'quill/dist/quill.snow.css'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Tulis catatan...'
  },
  error: {
    type: String,
    default: ''
  },
  readonly: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue'])

const editorId = ref(`quill-${Date.now()}`)
let quillInstance = null

const quillOptions = {
  theme: 'snow',
  placeholder: props.placeholder,
  readOnly: props.readonly,
  modules: {
    toolbar: [
      [{ 'header': [1, 2, 3, false] }],
      ['bold', 'italic', 'underline', 'strike'],
      [{ 'list': 'ordered'}, { 'list': 'bullet' }],
      [{ 'color': [] }, { 'background': [] }],
      [{ 'align': [] }],
      ['link'],
      ['clean']
    ]
  }
}

onMounted(() => {
  const editorElement = document.getElementById(editorId.value)
  if (editorElement) {
    quillInstance = new Quill(editorElement, quillOptions)
    
    // Set initial content
    if (props.modelValue) {
      quillInstance.root.innerHTML = props.modelValue
    }
    
    // Listen for content changes
    quillInstance.on('text-change', () => {
      const html = quillInstance.root.innerHTML
      emit('update:modelValue', html === '<p><br></p>' ? '' : html)
    })
  }
})

onBeforeUnmount(() => {
  if (quillInstance) {
    quillInstance = null
  }
})

// Watch for external changes to modelValue
watch(() => props.modelValue, (newValue) => {
  if (quillInstance && newValue !== quillInstance.root.innerHTML) {
    quillInstance.root.innerHTML = newValue || ''
  }
})
</script>

<style scoped>
:deep(.ql-toolbar) {
  border-top: 1px solid #e5e7eb;
  border-left: 1px solid #e5e7eb;
  border-right: 1px solid #e5e7eb;
  border-bottom: none;
  border-top-left-radius: 0.75rem;
  border-top-right-radius: 0.75rem;
  background-color: #fafafa;
}

:deep(.ql-container) {
  border-bottom: 1px solid #e5e7eb;
  border-left: 1px solid #e5e7eb;
  border-right: 1px solid #e5e7eb;
  border-top: none;
  border-bottom-left-radius: 0.75rem;
  border-bottom-right-radius: 0.75rem;
  font-family: system-ui, -apple-system, sans-serif;
}

:deep(.ql-editor) {
  min-height: 100px;
  padding: 12px 16px;
  font-size: 14px;
  line-height: 1.5;
}

:deep(.ql-editor.ql-blank::before) {
  color: #9ca3af;
  font-style: normal;
}

/* Custom button styling */
:deep(.ql-toolbar .ql-formats) {
  margin-right: 12px;
}

:deep(.ql-toolbar button) {
  padding: 4px;
  margin: 1px;
  border-radius: 4px;
}

:deep(.ql-toolbar button:hover) {
  background-color: #e5e7eb;
}

:deep(.ql-toolbar button.ql-active) {
  background-color: #3b82f6;
  color: white;
}
</style>