import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth.js'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/LoginView.vue'),
    meta: { public: true, title: 'Masuk' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/RegisterView.vue'),
    meta: { public: true, title: 'Daftar' }
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '',         name: 'Dashboard', component: () => import('@/views/DashboardView.vue'),  meta: { title: 'Dashboard' } },
      { path: 'silsilah', name: 'Silsilah',  component: () => import('@/views/SilsilahView.vue'),   meta: { title: 'Pohon Silsilah' } },
      { path: 'kas',      name: 'Kas',       component: () => import('@/views/KasView.vue'),         meta: { title: 'Kas Keluarga' } },
      { path: 'peta',     name: 'Peta',      component: () => import('@/views/PetaView.vue'),        meta: { title: 'Peta Domisili' } },
      { path: 'kalender', name: 'Kalender',  component: () => import('@/views/KalenderView.vue'),    meta: { title: 'Kalender Acara' } },
      { path: 'catatan',  name: 'Notes',     component: () => import('@/views/NotesView.vue'),       meta: { title: 'Catatan Keluarga' } },
    ]
  },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
})

let sessionRestored = false

router.beforeEach(async (to) => {
  const auth = useAuthStore()

  document.title = to.meta.title ? `${to.meta.title} — Wangsa` : 'Wangsa'

  // Restore session from cookie exactly once on first load
  if (!sessionRestored) {
    sessionRestored = true
    if (!auth.isLoggedIn) {
      await auth.initSession() // sets isInitializing = true → shows spinner in App.vue
    }
  }

  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return { name: 'Login' }
  }

  if (to.meta.public && auth.isLoggedIn) {
    return { path: '/' }
  }
})

export default router
