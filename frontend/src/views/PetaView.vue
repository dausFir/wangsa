<template>
  <div class="space-y-6">
    <div class="card p-0 overflow-hidden">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 pt-6 pb-4 flex-wrap gap-3">
        <div>
          <h3 class="text-lg font-semibold text-navy">Peta Domisili Keluarga</h3>
          <p class="text-sm text-warm-gray-400 mt-0.5">{{ filteredAddresses.length }} alamat terdaftar</p>
        </div>
        <div class="flex gap-2">
          <select v-model="filterMember" class="text-sm border border-warm-gray-200 rounded-lg px-3 py-1.5 bg-white">
            <option value="">Semua Anggota</option>
            <option value="null">Alamat Umum</option>
            <option v-for="member in familyMembers" :key="member.id" :value="member.id">
              {{ member.full_name }}
            </option>
          </select>
          <button @click="openAdd" class="btn-primary">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
            </svg>
            Tambah Alamat
          </button>
        </div>
      </div>

      <!-- Leaflet Map -->
      <div id="wangsa-map" style="height: 320px; z-index: 0;"></div>

      <!-- Address list -->
      <div class="px-6 pb-6 pt-4">
        <div v-if="loading" class="text-center py-10 text-warm-gray-400">Memuat alamat...</div>
        <div v-else class="space-y-3">
          <div
            v-for="addr in filteredAddresses"
            :key="addr.id"
            class="flex items-start gap-4 p-4 bg-warm-gray-50 rounded-xl border border-warm-gray-100
                   hover:border-warm-gray-200 transition-colors group cursor-pointer"
            @click="flyTo(addr)"
          >
            <div class="w-10 h-10 rounded-xl bg-terra/10 text-terra flex items-center justify-center flex-shrink-0 text-lg">🏠</div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <p class="font-semibold text-navy text-sm">
                  {{ addr.member_name ? addr.member_name : 'Alamat Umum' }}
                </p>
                <span class="badge-gray text-xs">{{ addr.label }}</span>
                <span v-if="addr.is_current" class="badge-green text-xs">Aktif</span>
                <span v-if="!addr.latitude" class="badge-orange text-xs">Belum ada koordinat</span>
              </div>
              <p class="text-sm text-warm-gray-600 mt-0.5 leading-snug">
                {{ [addr.street, addr.city, addr.province].filter(Boolean).join(', ') }}
              </p>
              <p v-if="addr.latitude != null" class="text-xs text-warm-gray-400 mt-1 font-mono">
                📍 {{ addr.latitude?.toFixed(6) }}, {{ addr.longitude?.toFixed(6) }}
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
          <div v-if="!filteredAddresses.length" class="text-center py-12">
            <div class="text-4xl mb-3">🏘️</div>
            <p v-if="filterMember" class="text-warm-gray-500">
              {{ filterMember === 'null' ? 'Belum ada alamat umum.' : 'Anggota ini belum memiliki alamat.' }}
            </p>
            <p v-else class="text-warm-gray-500">Belum ada alamat anggota keluarga terdaftar.</p>
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
              <label class="form-label">Anggota Keluarga</label>
              <select v-model="aForm.family_member_id" class="form-input">
                <option :value="null">Alamat Umum / Tidak terkait anggota</option>
                <option v-for="member in familyMembers" :key="member.id" :value="member.id">
                  {{ member.full_name }}
                </option>
              </select>
            </div>
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
                <input v-model="aForm.city" class="form-input" placeholder="Jakarta, Bandung, Surabaya..." required />
                <p class="text-xs text-warm-gray-500 mt-1">Gunakan nama kota yang umum dikenal</p>
              </div>
              <div>
                <label class="form-label">Provinsi</label>
                <input v-model="aForm.province" class="form-input" placeholder="DKI Jakarta, Jawa Barat..." />
              </div>
            </div>
            
            <!-- Coordinates Section -->
            <div class="border border-warm-gray-200 rounded-lg p-4 bg-warm-gray-50/50">
              <div class="flex items-center justify-between mb-3">
                <div>
                  <h5 class="text-sm font-medium text-navy">📍 Koordinat Lokasi</h5>
                  <p class="text-xs text-warm-gray-500">Diperlukan agar muncul di peta</p>
                </div>
                <button
                  type="button"
                  @click="findCoordinates"
                  :disabled="!aForm.city || geocoding"
                  class="btn-secondary text-xs py-1.5 px-3 disabled:opacity-50"
                >
                  <span v-if="geocoding" class="spinner w-3 h-3 border-warm-gray-400 border-t-warm-gray-700 mr-1.5"></span>
                  {{ geocoding ? 'Mencari...' : '🔍 Cari Otomatis' }}
                </button>
              </div>
              
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="form-label text-xs">Latitude</label>
                  <input v-model.number="aForm.latitude" type="number" step="any" class="form-input text-sm" placeholder="-6.2088" />
                </div>
                <div>
                  <label class="form-label text-xs">Longitude</label>
                  <input v-model.number="aForm.longitude" type="number" step="any" class="form-input text-sm" placeholder="106.8456" />
                </div>
              </div>
              
              <div v-if="geocodingError" class="mt-2 p-2 bg-red-50 border border-red-200 rounded text-xs text-red-600">
                {{ geocodingError }}
                <button
                  type="button"
                  @click="showCityList = !showCityList"
                  class="ml-2 underline hover:no-underline"
                >
                  {{ showCityList ? 'Sembunyikan' : 'Lihat daftar kota' }}
                </button>
                
                <div v-if="showCityList" class="mt-2 p-2 bg-white border border-warm-gray-200 rounded">
                  <p class="font-medium mb-2">Koordinat kota-kota besar:</p>
                  <div class="grid grid-cols-2 gap-1 text-xs">
                    <button
                      v-for="city in commonCities"
                      :key="city.name"
                      @click="useCommonCity(city)"
                      class="text-left p-1 hover:bg-warm-gray-100 rounded"
                    >
                      <strong>{{ city.name }}</strong><br>
                      <span class="text-warm-gray-500">{{ city.coords }}</span>
                    </button>
                  </div>
                </div>
              </div>
              
              <div v-if="!aForm.latitude" class="mt-2 p-2 bg-amber-50 border border-amber-200 rounded">
                <p class="text-xs text-amber-700 flex items-center gap-1.5">
                  <svg class="w-3 h-3 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
                  </svg>
                  Tanpa koordinat, alamat tidak akan muncul di peta
                </p>
              </div>
              
              <details class="mt-2">
                <summary class="text-xs text-warm-gray-500 cursor-pointer hover:text-warm-gray-700">Manual? Lihat cara mendapatkan koordinat</summary>
                <div class="mt-2 p-2 bg-white border border-warm-gray-200 rounded text-xs text-warm-gray-600 leading-relaxed">
                  <p><strong>Google Maps:</strong></p>
                  <p>1. Buka maps.google.com</p>
                  <p>2. Cari alamat atau klik lokasi di peta</p>
                  <p>3. Klik kanan → "What's here?"</p>
                  <p>4. Copy angka koordinat yang muncul</p>
                </div>
              </details>
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
import { ref, onMounted, onUnmounted, nextTick, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth.js'
import { useToast }    from '@/composables/useToast.js'
import api from '@/api/axios.js'

const auth      = useAuthStore()
const toast     = useToast()
const addresses = ref([])
const familyMembers = ref([])
const loading   = ref(false)
const showForm  = ref(false)
const saving    = ref(false)
const formError = ref('')
const editingId = ref(null)
const filterMember = ref('')
const geocoding = ref(false)
const geocodingError = ref('')
const showCityList = ref(false)

// Common Indonesian cities with coordinates
const commonCities = [
  { name: 'Jakarta', coords: '-6.2088, 106.8456', lat: -6.2088, lng: 106.8456 },
  { name: 'Bandung', coords: '-6.9175, 107.6191', lat: -6.9175, lng: 107.6191 },
  { name: 'Surabaya', coords: '-7.2575, 112.7521', lat: -7.2575, lng: 112.7521 },
  { name: 'Medan', coords: '3.5952, 98.6722', lat: 3.5952, lng: 98.6722 },
  { name: 'Yogyakarta', coords: '-7.7956, 110.3695', lat: -7.7956, lng: 110.3695 },
  { name: 'Semarang', coords: '-6.9667, 110.4167', lat: -6.9667, lng: 110.4167 },
  { name: 'Makassar', coords: '-5.1477, 119.4327', lat: -5.1477, lng: 119.4327 },
  { name: 'Palembang', coords: '-2.9761, 104.7754', lat: -2.9761, lng: 104.7754 }
]

const blankForm = () => ({
  family_member_id: null,
  label: 'Rumah', street: '', city: '', province: '',
  latitude: null, longitude: null, country: 'Indonesia', is_current: true
})
const aForm = ref(blankForm())

// Computed untuk filter alamat
const filteredAddresses = computed(() => {
  if (!filterMember.value) return addresses.value
  if (filterMember.value === 'null') {
    return addresses.value.filter(addr => addr.family_member_id == null)
  }
  return addresses.value.filter(addr => addr.family_member_id == filterMember.value)
})

// Watch filter changes untuk refresh map
watch(filteredAddresses, () => {
  nextTick(() => refreshMarkers())
}, { deep: true })

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

  const withCoords = filteredAddresses.value.filter(a => a.latitude != null && a.longitude != null)
  withCoords.forEach(addr => {
    const marker = L.marker([addr.latitude, addr.longitude])
    
    // Tentukan icon berdasarkan status
    const memberName = addr.member_name || 'Alamat Umum'
    const statusBadge = addr.is_current ? '<span style="color: #059669; font-size: 11px;">● Aktif</span>' : ''
    
    marker.bindPopup(`
      <div style="min-width: 150px;">
        <div style="font-weight: 600; color: #1e293b; margin-bottom: 4px;">${memberName}</div>
        <div style="font-size: 12px; color: #64748b; margin-bottom: 4px;">${addr.label} ${statusBadge}</div>
        <div style="font-size: 12px; color: #475569;">${[addr.street, addr.city].filter(Boolean).join(', ')}</div>
      </div>
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
async function fetchFamilyMembers() {
  try {
    const { data } = await api.get('/family/members')
    familyMembers.value = data.data ?? []
  } catch (e) {
    console.error('Failed to fetch family members:', e)
  }
}

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
  geocoding.value = false
  geocodingError.value = ''
  showCityList.value = false
  showForm.value  = true
}

function openEdit(addr) {
  editingId.value = addr.id
  aForm.value = {
    family_member_id: addr.family_member_id ?? null,
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
  geocoding.value = false
  geocodingError.value = ''
  showCityList.value = false
  showForm.value  = true
}

function closeForm() {
  showForm.value  = false
  editingId.value = null
  geocoding.value = false
  geocodingError.value = ''
  showCityList.value = false
}

// ── Geocoding ───────────────────────────────────────────────────────────────

async function findCoordinates() {
  if (!aForm.value.city) {
    geocodingError.value = 'Isi kota terlebih dahulu'
    return
  }

  geocoding.value = true
  geocodingError.value = ''

  try {
    // Strategy 1: Try detailed address first
    let result = await tryGeocode([
      aForm.value.street,
      aForm.value.city,
      aForm.value.province,
      'Indonesia'
    ].filter(Boolean))

    // Strategy 2: Fallback to city + province only
    if (!result && aForm.value.province) {
      result = await tryGeocode([aForm.value.city, aForm.value.province, 'Indonesia'])
    }

    // Strategy 3: Fallback to city only
    if (!result) {
      result = await tryGeocode([aForm.value.city, 'Indonesia'])
    }

    if (result) {
      aForm.value.latitude = parseFloat(result.lat)
      aForm.value.longitude = parseFloat(result.lon)
      
      // Show what was found
      const locationName = result.display_name.length > 80 
        ? result.display_name.substring(0, 80) + '...'
        : result.display_name
      
      toast.success(`✅ Koordinat ditemukan: ${locationName}`)
    } else {
      geocodingError.value = `Lokasi "${aForm.value.city}" tidak ditemukan. Coba gunakan nama kota yang lebih umum atau isi koordinat manual.`
    }
  } catch (error) {
    console.error('Geocoding error:', error)
    geocodingError.value = 'Gagal mencari koordinat. Silakan coba lagi atau isi manual.'
  } finally {
    geocoding.value = false
  }
}

async function tryGeocode(addressParts) {
  const searchQuery = addressParts.join(', ')
  console.log('🔍 Searching for:', searchQuery)
  
  try {
    // Try with Indonesia country code first
    let response = await fetch(
      `https://nominatim.openstreetmap.org/search?` + 
      `format=json&q=${encodeURIComponent(searchQuery)}&limit=3&countrycodes=id&addressdetails=1`,
      {
        headers: {
          'User-Agent': 'WangsaApp/1.0 (Family Tree App)'
        }
      }
    )
    
    if (!response.ok) throw new Error('Network error')
    
    let data = await response.json()
    console.log('🌍 Nominatim response:', data)
    
    // If no results with country code, try without it
    if (!data || data.length === 0) {
      console.log('🔄 Retrying without country code...')
      response = await fetch(
        `https://nominatim.openstreetmap.org/search?` + 
        `format=json&q=${encodeURIComponent(searchQuery)}&limit=3&addressdetails=1`,
        {
          headers: {
            'User-Agent': 'WangsaApp/1.0 (Family Tree App)'
          }
        }
      )
      
      if (response.ok) {
        data = await response.json()
        console.log('🌍 Nominatim response (no country filter):', data)
        
        // Filter for Indonesia manually
        data = data.filter(item => 
          item.display_name.toLowerCase().includes('indonesia') ||
          item.address?.country_code === 'id' ||
          item.address?.country === 'Indonesia'
        )
      }
    }
    
    return data && data.length > 0 ? data[0] : null
  } catch (error) {
    console.error('Geocoding request failed:', error)
    return null
  }
}

function useCommonCity(city) {
  aForm.value.latitude = city.lat
  aForm.value.longitude = city.lng
  geocodingError.value = ''
  showCityList.value = false
  toast.success(`✅ Menggunakan koordinat ${city.name}`)
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
  await Promise.all([fetchFamilyMembers(), fetchAddresses()])
})

onUnmounted(() => {
  if (mapInstance) { mapInstance.remove(); mapInstance = null }
})
</script>
