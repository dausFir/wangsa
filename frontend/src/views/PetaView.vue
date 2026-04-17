<template>
  <div class="space-y-6">
    <div class="card p-0 overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 pt-6 pb-4 flex-wrap gap-3">
        <div>
          <h3 class="text-lg font-semibold text-navy">Peta Domisili</h3>
          <p class="text-sm text-warm-gray-400 mt-0.5">{{ addresses.length }} alamat terdaftar</p>
        </div>
        <button @click="openAdd" class="btn-primary">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
          </svg>
          Tambah Alamat
        </button>
      </div>

      <!-- Leaflet Map -->
      <div id="wangsa-map" style="height: 320px; z-index: 0;"></div>

      <!-- Address list -->
      <div class="px-6 pb-6 pt-4">
        <div v-if="loading" class="text-center py-10 text-warm-gray-400">Memuat alamat...</div>
        <div v-else class="space-y-3">
          <div
            v-for="addr in addresses"
            :key="addr.id"
            class="flex items-start gap-4 p-4 bg-warm-gray-50 rounded-xl border border-warm-gray-100
                   hover:border-warm-gray-200 transition-colors group cursor-pointer"
            @click="flyTo(addr)"
          >
            <div class="w-10 h-10 rounded-xl bg-terra/10 text-terra flex items-center justify-center flex-shrink-0 text-lg">📍</div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <p class="font-semibold text-navy text-sm">{{ addr.member_name ?? 'Umum' }}</p>
                <span class="badge-gray text-xs">{{ addr.label }}</span>
                <span v-if="addr.is_current" class="badge-green text-xs">Aktif</span>
              </div>
              <p class="text-sm text-warm-gray-600 mt-0.5 leading-snug">
                {{ [addr.street, addr.city, addr.province, addr.country].filter(Boolean).join(', ') }}
              </p>
              <p v-if="addr.latitude != null" class="text-xs text-warm-gray-400 mt-1 font-mono">
                {{ addr.latitude?.toFixed(6) }}, {{ addr.longitude?.toFixed(6) }}
              </p>
            </div>
            <div class="flex gap-1 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
              <button
                v-if="auth.isSuperAdmin"
                @click.stop="openEdit(addr)"
                class="btn-ghost p-1.5 text-warm-gray-400 hover:text-navy hover:bg-warm-gray-100"
                title="Edit alamat"
              >
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                </svg>
              </button>
              <button
                v-if="auth.isSuperAdmin"
                @click.stop="deleteAddr(addr.id)"
                class="btn-ghost p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50"
                title="Hapus alamat"
              >
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
          </div>
          <div v-if="!addresses.length" class="text-center py-12">
            <div class="text-4xl mb-3">🏘️</div>
            <p class="text-warm-gray-500">Belum ada alamat terdaftar.</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Add / Edit Modal -->
    <Teleport to="body">
      <div
        v-if="showForm"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        style="background: rgba(20,29,39,0.52); backdrop-filter: blur(4px);"
        @mousedown.self="closeForm"
      >
        <div class="bg-white rounded-2xl shadow-modal w-full max-w-md">
          <div class="flex items-center justify-between px-6 py-5 border-b border-warm-gray-100">
            <h3 class="text-lg font-semibold text-navy">{{ editingId ? 'Edit Alamat' : 'Tambah Alamat' }}</h3>
            <button @click="closeForm" class="btn-ghost p-2 -mr-2">✕</button>
          </div>
          <form @submit.prevent="handleSave" class="px-6 py-5 space-y-4">
            <div>
              <label class="form-label">Label <span class="text-terra">*</span></label>
              <input v-model="aForm.label" class="form-input" placeholder="Rumah / Kantor / Kos" required />
            </div>
            <div>
              <label class="form-label">Alamat Jalan</label>
              <input v-model="aForm.street" class="form-input" placeholder="Jl. Merdeka No. 17" />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="form-label">Kota <span class="text-terra">*</span></label>
                <input v-model="aForm.city" class="form-input" placeholder="Jakarta" required />
              </div>
              <div>
                <label class="form-label">Provinsi</label>
                <input v-model="aForm.province" class="form-input" placeholder="DKI Jakarta" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="form-label">Latitude</label>
                <input v-model.number="aForm.latitude" type="number" step="any" class="form-input" placeholder="-6.2088" />
              </div>
              <div>
                <label class="form-label">Longitude</label>
                <input v-model.number="aForm.longitude" type="number" step="any" class="form-input" placeholder="106.8456" />
              </div>
            </div>
            <div class="flex items-center gap-2.5">
              <input v-model="aForm.is_current" type="checkbox" id="is_current_chk" class="w-4 h-4 accent-terra rounded" />
              <label for="is_current_chk" class="text-sm font-medium text-navy cursor-pointer">Alamat aktif saat ini</label>
            </div>
            <div v-if="formError" class="form-error">{{ formError }}</div>
            <div class="flex gap-3 pt-1">
              <button type="button" @click="closeForm" class="btn-secondary flex-1">Batal</button>
              <button type="submit" :disabled="saving" class="btn-primary flex-1">
                <span v-if="saving" class="spinner w-4 h-4 border-white/30 border-t-white"></span>
                {{ saving ? 'Menyimpan...' : (editingId ? 'Simpan Perubahan' : 'Simpan Alamat') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useAuthStore } from '@/stores/auth.js'
import { useToast }    from '@/composables/useToast.js'
import api from '@/api/axios.js'

const auth      = useAuthStore()
const toast     = useToast()
const addresses = ref([])
const loading   = ref(false)
const showForm  = ref(false)
const saving    = ref(false)
const formError = ref('')
const editingId = ref(null)

const blankForm = () => ({
  label: 'Rumah', street: '', city: '', province: '',
  latitude: null, longitude: null, country: 'Indonesia', is_current: true
})
const aForm = ref(blankForm())

// ── Leaflet map ────────────────────────────────────────────────
let mapInstance = null
let markersLayer = null

async function initMap() {
  // Load Leaflet CSS + JS from CDN
  if (!document.getElementById('leaflet-css')) {
    const link = document.createElement('link')
    link.id   = 'leaflet-css'
    link.rel  = 'stylesheet'
    link.href = 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/leaflet.min.css'
    document.head.appendChild(link)
  }

  await new Promise((resolve, reject) => {
    if (window.L) return resolve()
    const s = document.createElement('script')
    s.src = 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/leaflet.min.js'
    s.onload  = resolve
    s.onerror = reject
    document.head.appendChild(s)
  })

  const L = window.L
  mapInstance = L.map('wangsa-map', { zoomControl: true }).setView([-2.5, 118], 5)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
    maxZoom: 19,
  }).addTo(mapInstance)

  markersLayer = L.layerGroup().addTo(mapInstance)
}

function refreshMarkers() {
  if (!mapInstance || !window.L) return
  const L = window.L
  markersLayer.clearLayers()

  const withCoords = addresses.value.filter(a => a.latitude != null && a.longitude != null)
  withCoords.forEach(addr => {
    const marker = L.marker([addr.latitude, addr.longitude])
    marker.bindPopup(`
      <b>${addr.member_name ?? 'Umum'}</b><br/>
      <span style="font-size:12px;color:#888">${addr.label}</span><br/>
      ${[addr.street, addr.city].filter(Boolean).join(', ')}
    `)
    markersLayer.addLayer(marker)
  })

  if (withCoords.length === 1) {
    mapInstance.setView([withCoords[0].latitude, withCoords[0].longitude], 13)
  } else if (withCoords.length > 1) {
    const group = L.featureGroup(markersLayer.getLayers())
    mapInstance.fitBounds(group.getBounds().pad(0.2))
  }
}

function flyTo(addr) {
  if (!mapInstance || addr.latitude == null) return
  mapInstance.flyTo([addr.latitude, addr.longitude], 14, { duration: 1 })
}

// ── Data ────────────────────────────────────────────────────────
async function fetchAddresses() {
  loading.value = true
  try {
    const { data } = await api.get('/addresses')
    addresses.value = data.data ?? []
    await nextTick()
    refreshMarkers()
  } finally {
    loading.value = false
  }
}

function openAdd() {
  editingId.value = null
  aForm.value     = blankForm()
  formError.value = ''
  showForm.value  = true
}

function openEdit(addr) {
  editingId.value = addr.id
  aForm.value = {
    label:     addr.label,
    street:    addr.street    ?? '',
    city:      addr.city,
    province:  addr.province  ?? '',
    latitude:  addr.latitude  ?? null,
    longitude: addr.longitude ?? null,
    country:   addr.country   ?? 'Indonesia',
    is_current: addr.is_current,
  }
  formError.value = ''
  showForm.value  = true
}

function closeForm() {
  showForm.value  = false
  editingId.value = null
}

async function handleSave() {
  saving.value    = true
  formError.value = ''
  const payload = {
    ...aForm.value,
    street:    aForm.value.street   || null,
    province:  aForm.value.province || null,
    latitude:  aForm.value.latitude  ?? null,
    longitude: aForm.value.longitude ?? null,
  }
  try {
    if (editingId.value) {
      await api.put(`/addresses/${editingId.value}`, payload)
      toast.success('Alamat berhasil diperbarui')
    } else {
      await api.post('/addresses', payload)
      toast.success('Alamat berhasil ditambahkan')
    }
    closeForm()
    await fetchAddresses()
  } catch (e) {
    formError.value = e.response?.data?.error ?? 'Gagal menyimpan alamat.'
  } finally {
    saving.value = false
  }
}

async function deleteAddr(id) {
  if (!confirm('Hapus alamat ini?')) return
  try {
    await api.delete(`/addresses/${id}`)
    toast.success('Alamat berhasil dihapus')
    await fetchAddresses()
  } catch (e) {
    toast.error(e.response?.data?.error ?? 'Gagal menghapus alamat.')
  }
}

onMounted(async () => {
  await initMap()
  await fetchAddresses()
})

onUnmounted(() => {
  if (mapInstance) { mapInstance.remove(); mapInstance = null }
})
</script>
