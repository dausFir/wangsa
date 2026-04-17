import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
  timeout: 15000,
})

// Track if a refresh is already in-flight to avoid parallel refresh storms
let isRefreshing = false
let refreshQueue = []  // callbacks waiting for the new token

function processQueue(error) {
  refreshQueue.forEach(cb => error ? cb.reject(error) : cb.resolve())
  refreshQueue = []
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config

    // Only attempt refresh on 401, and only once per request
    if (
      error.response?.status === 401 &&
      !original._retried &&
      !original.url?.includes('/auth/refresh') &&
      !original.url?.includes('/auth/login')
    ) {
      original._retried = true

      if (isRefreshing) {
        // Another refresh is already happening — queue this request
        return new Promise((resolve, reject) => {
          refreshQueue.push({ resolve, reject })
        }).then(() => api(original))
          .catch(err => Promise.reject(err))
      }

      isRefreshing = true

      try {
        // Attempt to get new access token using the refresh token cookie
        // The refresh cookie is scoped to /api/auth so it's sent automatically
        await api.post('/auth/refresh')
        processQueue(null)
        return api(original)  // retry the original failed request
      } catch (refreshError) {
        processQueue(refreshError)
        // Refresh failed — session truly expired, redirect to login
        if (window.location.pathname !== '/login') {
          window.location.href = '/login'
        }
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export default api
