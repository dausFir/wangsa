import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router/index.js'
import App from './App.vue'
import './style.css'

const app  = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Global error boundary — catches unhandled Vue component errors.
// Logs to console and shows a toast, rather than crashing the whole page.
app.config.errorHandler = (err, _instance, info) => {
  console.error('[Vue error]', info, err)

  // Dynamically import to avoid circular dependency at module init time
  import('@/stores/toast.js').then(({ useToastStore }) => {
    useToastStore().error(
      'Terjadi kesalahan tak terduga. Silakan refresh halaman jika masalah berlanjut.'
    )
  }).catch(() => {
    // Toast not available — already logged above
  })
}

// Global unhandled promise rejection catcher
window.addEventListener('unhandledrejection', (event) => {
  // Don't surface cancelled axios requests (e.g. navigating away mid-request)
  if (event.reason?.code === 'ERR_CANCELED') return
  console.error('[Unhandled rejection]', event.reason)
})

app.mount('#app')
