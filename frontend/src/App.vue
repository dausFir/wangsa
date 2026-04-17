<template>
  <!-- Global loading overlay saat router guard sedang restore session -->
  <div v-if="auth.isInitializing" class="app-init-overlay">
    <div class="app-init-spinner"></div>
  </div>

  <template v-else>
    <RouterView v-slot="{ Component }">
      <Transition name="page" mode="out-in">
        <component :is="Component" />
      </Transition>
    </RouterView>
  </template>

  <!-- Global toast container — selalu ada, di atas segalanya -->
  <ToastContainer />
</template>

<script setup>
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'
import ToastContainer from '@/components/ui/ToastContainer.vue'

const auth = useAuthStore()
</script>

<style scoped>
.app-init-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f7f4;
  z-index: 9998;
}
.app-init-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid #e8e4dd;
  border-top-color: #CC6649;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
