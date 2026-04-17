<template>
  <Teleport to="body">
    <div
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
      style="background: rgba(20,29,39,0.52); backdrop-filter: blur(4px);"
      @mousedown.self="confirmClose"
    >
      <Transition name="modal" appear>
        <div class="bg-white rounded-2xl shadow-modal w-full max-w-xl max-h-[92vh] flex flex-col">
          <!-- Header -->
          <div class="flex items-center justify-between px-6 py-5 border-b border-warm-gray-100 flex-shrink-0">
            <div>
              <div class="flex items-center gap-2">
                <h3 class="text-lg font-semibold text-navy">
                  {{ isEdit ? 'Edit Anggota' : 'Tambah Anggota Keluarga' }}
                </h3>
                <span v-if="isDirty" class="text-xs text-amber-600 bg-amber-50 px-2 py-0.5 rounded-full border border-amber-200">
                  belum disimpan
                </span>
              </div>
              <p class="text-xs text-warm-gray-400 mt-0.5">
                {{ isEdit ? 'Perbarui data anggota di bawah ini.' : 'Isi data untuk menambah anggota baru.' }}
              </p>
            </div>
            <button @click="confirmClose" class="btn-ghost p-2 -mr-2 flex-shrink-0">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Scrollable form body -->
          <div class="flex-1 overflow-y-auto">
            <form id="member-form" @submit.prevent="handleSubmit" class="px-6 py-5 space-y-5">

              <!-- ── Photo upload ── -->
              <div class="flex items-center gap-4">
                <!-- Avatar preview -->
                <div :class="['w-20 h-20 rounded-2xl flex-shrink-0 flex items-center justify-center overflow-hidden font-bold text-xl', form.gender === 'male' ? 'bg-blue-100 text-blue-700' : 'bg-pink-100 text-pink-600']">
                  <img
                    v-if="photoPreview"
                    :src="sanitizePhotoUrl(photoPreview)"
                    alt="Preview foto"
                    class="w-full h-full object-cover"
                  />
                  <span v-else>{{ previewInitials }}</span>
                </div>

                <!-- Upload controls -->
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-navy mb-1">Foto Anggota</p>
                  <p class="text-xs text-warm-gray-400 mb-2">JPEG, PNG, atau GIF · max 5MB · otomatis di-resize ke 400×400</p>
                  <div class="flex gap-2 flex-wrap">
                    <label class="btn-secondary text-sm cursor-pointer">
                      <svg class="w-4 h-4 inline mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                      </svg>
                      {{ photoPreview ? 'Ganti Foto' : 'Pilih Foto' }}
                      <input
                        type="file"
                        accept="image/*"
                        class="hidden"
                        @change="onFileSelect"
                      />
                    </label>
                    <button
                      v-if="photoPreview"
                      type="button"
                      @click="removePhoto"
                      class="btn-ghost text-sm text-red-400 hover:text-red-600"
                    >
                      Hapus
                    </button>
                  </div>
                  <p v-if="photoError" class="text-xs text-red-500 mt-1">{{ photoError }}</p>
                  <p v-if="photoUploading" class="text-xs text-warm-gray-400 mt-1 flex items-center gap-1.5">
                    <span class="spinner w-3 h-3 border-warm-gray-300 border-t-warm-gray-600 inline-block"></span>
                    Mengupload foto...
                  </p>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-x-4 gap-y-5">
                <div class="col-span-2">
                  <label class="form-label">Nama Lengkap <span class="text-terra">*</span></label>
                  <input v-model="form.full_name" class="form-input" required minlength="2" placeholder="Ahmad Budi Santoso" autocomplete="off" />
                </div>

                <div>
                  <label class="form-label">Nama Panggilan</label>
                  <input v-model="form.nickname" class="form-input" placeholder="Ahmad" />
                </div>

                <div>
                  <label class="form-label">Jenis Kelamin <span class="text-terra">*</span></label>
                  <select v-model="form.gender" class="form-select" required>
                    <option value="male">♂  Laki-laki</option>
                    <option value="female">♀  Perempuan</option>
                  </select>
                </div>

                <div>
                  <label class="form-label">Tanggal Lahir</label>
                  <input v-model="form.birth_date" type="date" class="form-input" />
                </div>

                <div>
                  <label class="form-label">Tempat Lahir</label>
                  <input v-model="form.birth_place" class="form-input" placeholder="Jakarta" />
                </div>

                <div>
                  <label class="form-label">Tanggal Meninggal</label>
                  <input v-model="form.death_date" type="date" class="form-input" />
                </div>

                <div>
                  <label class="form-label">Orang Tua</label>
                  <select v-model="form.parent_id" class="form-select">
                    <option :value="null">— Tidak ada (akar silsilah) —</option>
                    <option v-for="m in eligibleParents" :key="m.id" :value="m.id">{{ m.full_name }}</option>
                  </select>
                </div>

                <div class="col-span-2">
                  <label class="form-label">Catatan</label>
                  <textarea v-model="form.notes" class="form-input resize-none" rows="3" placeholder="Informasi tambahan..."></textarea>
                </div>
              </div>

              <!-- ── Address section ── -->
              <div class="pt-4 border-t border-warm-gray-100">
                <div class="flex items-center justify-between mb-4">
                  <h4 class="text-sm font-semibold text-navy flex items-center gap-2">
                    🏠 Alamat Tempat Tinggal
                    <span v-if="memberAddresses.length > 0" class="badge-gray text-xs">{{ memberAddresses.length }} alamat</span>
                  </h4>
                  <button 
                    type="button"
                    @click="showAddressForm = !showAddressForm"
                    class="text-xs text-terra hover:text-terra/80 font-medium flex items-center gap-1"
                  >
                    <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" :d="showAddressForm ? 'M19 9l-7 7-7-7' : 'M12 4v16m8-8H4'"/>
                    </svg>
                    {{ showAddressForm ? 'Sembunyikan' : 'Tambah Alamat' }}
                  </button>
                </div>
                
                <!-- Existing addresses -->
                <div v-if="memberAddresses.length > 0" class="space-y-2 mb-4">
                  <div 
                    v-for="addr in memberAddresses" 
                    :key="addr.id"
                    class="flex items-center justify-between p-3 bg-warm-gray-50 rounded-lg text-sm border border-warm-gray-100"
                  >
                    <div class="flex-1">
                      <div class="flex items-center gap-2 mb-1">
                        <span class="font-medium text-navy">{{ addr.label }}</span>
                        <span v-if="addr.is_current" class="badge-green text-xs">Aktif</span>
                      </div>
                      <p class="text-warm-gray-600 text-xs">
                        {{ [addr.street, addr.city, addr.province].filter(Boolean).join(', ') }}
                      </p>
                    </div>
                    <div class="flex gap-1">
                      <button
                        type="button"
                        @click="editAddress(addr)"
                        class="text-warm-gray-400 hover:text-navy text-xs p-1 rounded"
                        title="Edit alamat"
                      >
                        <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                        </svg>
                      </button>
                      <button
                        type="button"
                        @click="deleteAddress(addr.id)"
                        class="text-red-400 hover:text-red-600 text-xs p-1 rounded"
                        title="Hapus alamat"
                      >
                        <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>

                <!-- No addresses message for existing member -->
                <div v-if="isEdit && memberAddresses.length === 0 && !showAddressForm" class="text-center py-4 text-warm-gray-400 text-sm">
                  <p>🏠 Anggota ini belum memiliki alamat.</p>
                  <p class="text-xs mt-1">Klik "Tambah Alamat" untuk menambah tempat tinggal.</p>
                </div>

                <!-- Add/Edit address form -->
                <div v-if="showAddressForm" class="space-y-3 p-4 bg-warm-gray-50 rounded-lg border border-warm-gray-200">
                  <div class="flex items-center gap-2 text-xs text-warm-gray-600 mb-2">
                    <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
                    </svg>
                    {{ isEdit ? 'Edit alamat anggota keluarga' : 'Alamat akan disimpan setelah anggota dibuat' }}
                  </div>
                  
                  <div class="grid grid-cols-2 gap-3">
                    <div>
                      <label class="form-label text-xs">Label <span class="text-terra">*</span></label>
                      <input v-model="addressForm.label" class="form-input text-sm" placeholder="Rumah" required />
                    </div>
                    <div class="flex items-center gap-2 pt-6">
                      <input v-model="addressForm.is_current" type="checkbox" id="addr_current" class="w-4 h-4 accent-terra rounded" />
                      <label for="addr_current" class="text-xs font-medium text-navy cursor-pointer">Alamat aktif</label>
                    </div>
                  </div>
                  
                  <div>
                    <label class="form-label text-xs">Alamat Jalan</label>
                    <input v-model="addressForm.street" class="form-input text-sm" placeholder="Jl. Merdeka No. 17" />
                  </div>
                  
                  <div class="grid grid-cols-2 gap-3">
                    <div>
                      <label class="form-label text-xs">Kota <span class="text-terra">*</span></label>
                      <input v-model="addressForm.city" class="form-input text-sm" placeholder="Jakarta, Bandung, Surabaya..." required />
                    </div>
                    <div>
                      <label class="form-label text-xs">Provinsi</label>
                      <input v-model="addressForm.province" class="form-input text-sm" placeholder="DKI Jakarta, Jawa Barat..." />
                    </div>
                  </div>
                  
                  <!-- Coordinates for map -->
                  <div class="border border-warm-gray-200 rounded p-3 bg-white">
                    <div class="flex items-center justify-between mb-2">
                      <label class="text-xs font-medium text-navy">📍 Koordinat (opsional)</label>
                      <button
                        type="button"
                        @click="findMemberAddressCoordinates"
                        :disabled="!addressForm.city || addressGeocoding"
                        class="text-xs text-terra hover:text-terra/80 disabled:opacity-50"
                      >
                        {{ addressGeocoding ? 'Mencari...' : '🔍 Cari Otomatis' }}
                      </button>
                    </div>
                    
                    <div class="grid grid-cols-2 gap-2">
                      <input v-model.number="addressForm.latitude" type="number" step="any" class="form-input text-xs" placeholder="Latitude" />
                      <input v-model.number="addressForm.longitude" type="number" step="any" class="form-input text-xs" placeholder="Longitude" />
                    </div>
                    
                    <p class="text-xs text-warm-gray-500 mt-1">
                      Dibutuhkan agar muncul di peta domisili
                    </p>
                  </div>
                  
                  <div class="flex gap-2 pt-2">
                    <button type="button" @click="saveAddress" class="btn-secondary text-sm flex-1">
                      {{ editingAddress ? 'Update Alamat' : 'Simpan Alamat' }}
                    </button>
                    <button type="button" @click="cancelAddressForm" class="btn-ghost text-sm px-3">
                      Batal
                    </button>
                  </div>
                </div>
              </div>

              <div v-if="error" class="form-error">
                <svg class="w-4 h-4 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                {{ error }}
              </div>
            </form>
          </div>

          <!-- Footer -->
          <div class="flex gap-3 px-6 py-4 border-t border-warm-gray-100 flex-shrink-0">
            <button type="button" @click="confirmClose" class="btn-secondary flex-1">Batal</button>
            <button type="submit" form="member-form" :disabled="loading || photoUploading" class="btn-primary flex-1">
              <span v-if="loading" class="spinner w-4 h-4 border-white/30 border-t-white"></span>
              {{ loading ? 'Menyimpan...' : (isEdit ? 'Simpan Perubahan' : 'Tambah Anggota') }}
            </button>
          </div>
        </div>
      </Transition>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useFamilyStore } from '@/stores/family.js'
import { useToast }       from '@/composables/useToast.js'
import { useDirtyGuard }  from '@/composables/useDirtyGuard.js'
import api from '@/api/axios.js'
import { sanitizePhotoUrl } from '@/utils/sanitize.js'

const props = defineProps({
  member:          { type: Object, default: null  },
  defaultParentId: { type: Number, default: null  },
  allMembers:      { type: Array,  default: () => [] },
})

const emit   = defineEmits(['close', 'saved'])
const family = useFamilyStore()
const toast  = useToast()

const loading        = ref(false)
const error          = ref('')
const photoPreview   = ref(null)   // data URL shown in preview
const photoFile      = ref(null)   // pending File to upload after save
const photoUploading = ref(false)
const photoError     = ref('')
const pendingDelete  = ref(false)  // user clicked "Hapus" before saving

// Address handling
const showAddressForm = ref(false)
const memberAddresses = ref([])
const editingAddress = ref(null)
const addressGeocoding = ref(false)
const addressForm = ref({
  label: 'Rumah',
  street: '',
  city: '',
  province: '',
  latitude: null,
  longitude: null,
  is_current: true
})

const isEdit = computed(() => !!props.member)

const eligibleParents = computed(() =>
  props.allMembers.filter(m => !isEdit.value || m.id !== props.member?.id)
)

const blankForm = () => ({
  full_name: '', nickname: null, gender: 'male',
  birth_date: null, birth_place: null, death_date: null,
  parent_id: props.defaultParentId ?? null, notes: null,
})

const form = ref(blankForm())
const { isDirty, setInitial, confirmClose } = useDirtyGuard(form, emit)

// Preview initials from current name input
const previewInitials = computed(() =>
  (form.value.full_name || '?')
    .split(' ').filter(Boolean).map(w => w[0]).join('').toUpperCase().slice(0, 2)
)

watch(() => props.member, async (m) => {
  if (m) {
    form.value = {
      full_name:   m.full_name   ?? '',
      nickname:    m.nickname    ?? null,
      gender:      m.gender      ?? 'male',
      birth_date:  m.birth_date  ?? null,
      birth_place: m.birth_place ?? null,
      death_date:  m.death_date  ?? null,
      parent_id:   m.parent_id   ?? null,
      notes:       m.notes       ?? null,
    }
    photoPreview.value = sanitizePhotoUrl(m.photo_url) ?? null
    await fetchMemberAddresses(m.id)
  } else {
    form.value = blankForm()
    photoPreview.value = null
    memberAddresses.value = []
  }
  photoFile.value = null
  pendingDelete.value = false
  showAddressForm.value = false
  resetAddressForm()
  setTimeout(() => setInitial(form.value), 0)
}, { immediate: true })

onMounted(() => setInitial(form.value))

// ── Photo handling ──────────────────────────────────────────────────────────

function onFileSelect(e) {
  photoError.value = ''
  const file = e.target.files?.[0]
  if (!file) return

  if (file.size > 5 * 1024 * 1024) {
    photoError.value = 'File terlalu besar. Maksimal 5MB.'
    return
  }
  if (!file.type.startsWith('image/')) {
    photoError.value = 'Hanya file gambar yang diizinkan.'
    return
  }

  photoFile.value    = file
  pendingDelete.value = false

  // Show local preview immediately — feels instant
  const reader = new FileReader()
  reader.onload = (ev) => { photoPreview.value = ev.target.result }
  reader.readAsDataURL(file)
}

function removePhoto() {
  photoPreview.value  = null
  photoFile.value     = null
  pendingDelete.value = !!props.member?.photo_url  // only need API call if was saved
}

// ── Address handling ───────────────────────────────────────────────────────

async function fetchMemberAddresses(memberId) {
  if (!memberId) return
  try {
    const { data } = await api.get('/addresses')
    memberAddresses.value = (data.data || []).filter(addr => addr.family_member_id === memberId)
  } catch (e) {
    console.error('Failed to fetch member addresses:', e)
  }
}

function resetAddressForm() {
  addressForm.value = {
    label: 'Rumah',
    street: '',
    city: '',
    province: '',
    latitude: null,
    longitude: null,
    is_current: true
  }
  editingAddress.value = null
  addressGeocoding.value = false
}

function editAddress(addr) {
  addressForm.value = {
    label: addr.label || 'Rumah',
    street: addr.street || '',
    city: addr.city || '',
    province: addr.province || '',
    latitude: addr.latitude || null,
    longitude: addr.longitude || null,
    is_current: addr.is_current || false
  }
  editingAddress.value = addr
  addressGeocoding.value = false
  showAddressForm.value = true
}

async function findMemberAddressCoordinates() {
  if (!addressForm.value.city) {
    toast.error('Isi kota terlebih dahulu')
    return
  }

  addressGeocoding.value = true

  try {
    // Strategy 1: Try detailed address first
    let result = await tryMemberGeocode([
      addressForm.value.street,
      addressForm.value.city,
      addressForm.value.province,
      'Indonesia'
    ].filter(Boolean))

    // Strategy 2: Fallback to city + province only
    if (!result && addressForm.value.province) {
      result = await tryMemberGeocode([addressForm.value.city, addressForm.value.province, 'Indonesia'])
    }

    // Strategy 3: Fallback to city only
    if (!result) {
      result = await tryMemberGeocode([addressForm.value.city, 'Indonesia'])
    }

    if (result) {
      addressForm.value.latitude = parseFloat(result.lat)
      addressForm.value.longitude = parseFloat(result.lon)
      
      toast.success('✅ Koordinat berhasil ditemukan!')
    } else {
      toast.error(`Lokasi "${addressForm.value.city}" tidak ditemukan. Coba gunakan nama kota yang lebih umum.`)
    }
  } catch (error) {
    console.error('Geocoding error:', error)
    toast.error('Gagal mencari koordinat. Silakan coba lagi.')
  } finally {
    addressGeocoding.value = false
  }
}

async function tryMemberGeocode(addressParts) {
  const searchQuery = addressParts.join(', ')
  console.log('🔍 Searching member address:', searchQuery)
  
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
    
    // If no results with country code, try without it
    if (!data || data.length === 0) {
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
    console.error('Member geocoding request failed:', error)
    return null
  }
}

async function saveAddress() {
  if (!addressForm.value.city) {
    toast.error('Kota wajib diisi')
    return
  }

  try {
    const payload = {
      family_member_id: props.member?.id || null,
      label: addressForm.value.label,
      street: addressForm.value.street || null,
      city: addressForm.value.city,
      province: addressForm.value.province || null,
      country: 'Indonesia',
      is_current: addressForm.value.is_current,
      latitude: addressForm.value.latitude ?? null,
      longitude: addressForm.value.longitude ?? null
    }

    if (editingAddress.value) {
      await api.put(`/addresses/${editingAddress.value.id}`, payload)
      toast.success('Alamat berhasil diperbarui')
    } else {
      await api.post('/addresses', payload)
      toast.success('Alamat berhasil ditambahkan')
    }

    if (props.member?.id) {
      await fetchMemberAddresses(props.member.id)
    }
    cancelAddressForm()
  } catch (e) {
    toast.error(e.response?.data?.error ?? 'Gagal menyimpan alamat')
  }
}

function cancelAddressForm() {
  showAddressForm.value = false
  addressGeocoding.value = false
  resetAddressForm()
}

async function deleteAddress(addressId) {
  if (!confirm('Hapus alamat ini?')) return
  
  try {
    await api.delete(`/addresses/${addressId}`)
    toast.success('Alamat berhasil dihapus')
    
    if (props.member?.id) {
      await fetchMemberAddresses(props.member.id)
    }
  } catch (e) {
    toast.error(e.response?.data?.error ?? 'Gagal menghapus alamat')
  }
}

// ── Submit ──────────────────────────────────────────────────────────────────

async function handleSubmit() {
  loading.value = true
  error.value   = ''

  try {
    // 1. Save / update member data (without photo_url — handled separately)
    const payload = { ...form.value }
    const optional = ['nickname', 'birth_date', 'birth_place', 'death_date', 'notes']
    for (const k of optional) { if (payload[k] === '') payload[k] = null }

    let savedMember
    if (isEdit.value) {
      savedMember = await family.updateMember(props.member.id, payload)
    } else {
      savedMember = await family.createMember(payload)
    }

    const memberId = savedMember?.id ?? props.member?.id

    // 2. Upload new photo if selected
    if (photoFile.value && memberId) {
      photoUploading.value = true
      try {
        const fd = new FormData()
        fd.append('photo', photoFile.value)
        await api.post(`/family/members/${memberId}/photo`, fd, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
      } catch (photoErr) {
        // Member saved OK but photo failed — warn, don't block
        toast.error('Data disimpan, tapi gagal upload foto: ' + (photoErr.response?.data?.error ?? 'Error tidak diketahui'))
      } finally {
        photoUploading.value = false
      }
    } else if (pendingDelete.value && memberId) {
      // 3. Delete photo if user clicked "Hapus"
      try {
        await api.delete(`/family/members/${memberId}/photo`)
      } catch { /* non-fatal */ }
    }

    // 4. Save address if form is filled and member was just created
    if (!isEdit.value && memberId && showAddressForm.value && addressForm.value.city) {
      try {
        const addressPayload = {
          family_member_id: memberId,
          label: addressForm.value.label,
          street: addressForm.value.street || null,
          city: addressForm.value.city,
          province: addressForm.value.province || null,
          country: 'Indonesia',
          is_current: addressForm.value.is_current,
          latitude: addressForm.value.latitude ?? null,
          longitude: addressForm.value.longitude ?? null
        }
        await api.post('/addresses', addressPayload)
      } catch (addrErr) {
        console.error('Failed to save address:', addrErr)
      }
    }

    toast.success(isEdit.value ? 'Data anggota berhasil diperbarui' : 'Anggota keluarga berhasil ditambahkan')

    // Refresh tree to show updated photo
    await Promise.all([family.fetchTree(), family.fetchMembers()])

    emit('saved')
    emit('close')
  } catch (e) {
    error.value = e.response?.data?.error ?? 'Terjadi kesalahan. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}
</script>
