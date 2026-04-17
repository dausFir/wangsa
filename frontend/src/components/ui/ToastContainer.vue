<template>
  <Teleport to="body">
    <div class="toast-stack" aria-live="polite" aria-atomic="false">
      <TransitionGroup name="toast" tag="div" class="toast-inner">
        <div
          v-for="toast in toastStore.toasts"
          :key="toast.id"
          :class="['toast-item', `toast-${toast.type}`]"
          role="alert"
        >
          <!-- Icon -->
          <svg v-if="toast.type === 'success'" class="toast-icon" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
          </svg>
          <svg v-else-if="toast.type === 'error'" class="toast-icon" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
          <svg v-else class="toast-icon" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
          </svg>

          <span class="toast-message">{{ toast.message }}</span>

          <button class="toast-close" @click="toastStore.remove(toast.id)" aria-label="Tutup">
            <svg viewBox="0 0 16 16" fill="currentColor" width="12" height="12">
              <path d="M3.72 3.72a.75.75 0 011.06 0L8 6.94l3.22-3.22a.75.75 0 111.06 1.06L9.06 8l3.22 3.22a.75.75 0 11-1.06 1.06L8 9.06l-3.22 3.22a.75.75 0 01-1.06-1.06L6.94 8 3.72 4.78a.75.75 0 010-1.06z"/>
            </svg>
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { useToastStore } from '@/stores/toast.js'
const toastStore = useToastStore()
</script>

<style scoped>
.toast-stack {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  z-index: 9999;
  pointer-events: none;
}
.toast-inner {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: flex-end;
}
.toast-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 500;
  max-width: 360px;
  pointer-events: all;
  border: 1px solid transparent;
}
.toast-success { background: #EAF3DE; color: #27500A; border-color: #C0DD97; }
.toast-error   { background: #FCEBEB; color: #791F1F; border-color: #F7C1C1; }
.toast-info    { background: #E6F1FB; color: #0C447C; border-color: #B5D4F4; }
.toast-icon    { width: 18px; height: 18px; flex-shrink: 0; }
.toast-message { flex: 1; line-height: 1.4; }
.toast-close {
  flex-shrink: 0;
  background: none;
  border: none;
  cursor: pointer;
  opacity: 0.5;
  padding: 2px;
  color: inherit;
  line-height: 0;
}
.toast-close:hover { opacity: 1; }

.toast-enter-active { transition: all 0.25s ease; }
.toast-leave-active { transition: all 0.2s ease; }
.toast-enter-from   { opacity: 0; transform: translateX(20px); }
.toast-leave-to     { opacity: 0; transform: translateX(20px); }
</style>
