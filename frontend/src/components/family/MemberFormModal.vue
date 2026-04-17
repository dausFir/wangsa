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

watch(() => props.member, (m) => {
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
  } else {
    form.value     = blankForm()
    photoPreview.value = null
  }
  photoFile.value   = null
  pendingDelete.value = false
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
        toast.success(isEdit.value ? 'Data & foto berhasil diperbarui' : 'Anggota berhasil ditambahkan')
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
      toast.success(isEdit.value ? 'Data anggota diperbarui' : 'Anggota berhasil ditambahkan')
    } else {
      toast.success(isEdit.value ? 'Data anggota berhasil diperbarui' : 'Anggota keluarga berhasil ditambahkan')
    }

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
