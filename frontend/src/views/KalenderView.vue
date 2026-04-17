<template>
  <div class="space-y-6">
    <div class="card">
      <div class="flex items-center justify-between mb-6 flex-wrap gap-3">
        <div>
          <h3 class="text-lg font-semibold text-navy">Acara Mendatang</h3>
          <p class="text-sm text-warm-gray-400 mt-0.5">{{ eventsStore.events.length }} acara terjadwal</p>
        </div>
        <button @click="openCreate" class="btn-primary">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
          </svg>
          Buat Acara
        </button>
      </div>

      <div v-if="eventsStore.loading" class="flex items-center justify-center py-12 gap-3 text-warm-gray-400">
        <div class="spinner w-6 h-6 border-warm-gray-200 border-t-warm-gray-400"></div>
        <span>Memuat acara...</span>
      </div>

      <div v-else-if="!eventsStore.events.length" class="text-center py-16">
        <div class="text-5xl mb-4">📅</div>
        <p class="text-warm-gray-500">Belum ada acara terjadwal.</p>
        <p class="text-warm-gray-400 text-sm mt-1">Buat acara keluarga pertama Anda!</p>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="event in eventsStore.events"
          :key="event.id"
          class="flex items-start gap-4 p-4 rounded-2xl border border-warm-gray-100
                 hover:border-warm-gray-200 hover:shadow-soft transition-all group"
        >
          <!-- Date badge -->
          <div
            class="w-14 h-14 rounded-xl flex flex-col items-center justify-center text-white flex-shrink-0"
            :style="{ background: event.color ?? '#CC6649' }"
          >
            <span class="text-2xl font-bold leading-none">{{ fmtDay(event.start_at) }}</span>
            <span class="text-[10px] uppercase opacity-80 font-medium tracking-wide">{{ fmtMonth(event.start_at) }}</span>
          </div>

          <!-- Info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 flex-wrap mb-1">
              <p class="font-semibold text-navy">{{ event.title }}</p>
              <span v-if="event.is_recurring" class="badge-blue text-xs">🔁 Rutin</span>
              <span v-if="event.attendee_count > 0" class="badge-gray text-xs">{{ event.attendee_count }} hadir</span>
            </div>
            <p class="text-sm text-warm-gray-500">{{ fmtDateTime(event.start_at) }}</p>
            <p v-if="event.location" class="text-xs text-warm-gray-400 mt-1 flex items-center gap-1">
              <span>📍</span> {{ event.location }}
            </p>
            <p v-if="event.description" class="text-sm text-warm-gray-500 mt-1.5 line-clamp-2">{{ event.description }}</p>

            <!-- My RSVP -->
            <div class="flex items-center gap-2 mt-3 flex-wrap">
              <span class="text-xs text-warm-gray-400 font-medium">RSVP saya:</span>
              <button
                v-for="opt in rsvpOptions"
                :key="opt.value"
                @click="setRSVP(event.id, opt.value)"
                :class="['text-xs px-2.5 py-1 rounded-full border transition-all font-medium',
                  myRSVP[event.id] === opt.value
                    ? opt.activeClass
                    : 'border-warm-gray-200 text-warm-gray-400 hover:border-warm-gray-400']"
              >{{ opt.label }}</button>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex gap-1 flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity mt-0.5">
            <button
              v-if="auth.isSuperAdmin"
              @click="openEdit(event)"
              class="btn-ghost p-1.5 text-warm-gray-400 hover:text-navy hover:bg-warm-gray-100"
              title="Edit acara"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
            </button>
            <button
              v-if="auth.isSuperAdmin"
              @click="handleDelete(event.id)"
              class="btn-ghost p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50"
              title="Hapus acara"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create / Edit Modal -->
    <Teleport to="body">
      <div
        v-if="showForm"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        style="background: rgba(20,29,39,0.52); backdrop-filter: blur(4px);"
        @mousedown.self="showForm = false"
      >
        <div class="bg-white rounded-2xl shadow-modal w-full max-w-md max-h-[90vh] overflow-y-auto">
          <div class="flex items-center justify-between px-6 py-5 border-b border-warm-gray-100 sticky top-0 bg-white z-10">
            <h3 class="text-lg font-semibold text-navy">{{ editingId ? 'Edit Acara' : 'Buat Acara Baru' }}</h3>
            <button @click="showForm = false" class="btn-ghost p-2 -mr-2">✕</button>
          </div>

          <form @submit.prevent="handleSave" class="px-6 py-5 space-y-4">
            <div>
              <label class="form-label">Judul Acara <span class="text-terra">*</span></label>
              <input v-model="eForm.title" class="form-input" required placeholder="Reuni Keluarga Besar 2025" />
            </div>
            <div>
              <label class="form-label">Deskripsi</label>
              <textarea v-model="eForm.description" class="form-input resize-none" rows="2" placeholder="Detail singkat acara..."></textarea>
            </div>
            <div>
              <label class="form-label">Lokasi</label>
              <input v-model="eForm.location" class="form-input" placeholder="Aula Desa Merdeka, Jawa Tengah" />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="form-label">Mulai <span class="text-terra">*</span></label>
                <input v-model="eForm.start_at" type="datetime-local" class="form-input" required />
              </div>
              <div>
                <label class="form-label">Selesai</label>
                <input v-model="eForm.end_at" type="datetime-local" class="form-input" />
              </div>
            </div>
            <div>
              <label class="form-label">Warna Label</label>
              <div class="flex items-center gap-3">
                <input v-model="eForm.color" type="color" class="w-10 h-10 rounded-xl cursor-pointer border border-warm-gray-200 p-0.5" />
                <div class="flex gap-2">
                  <button
                    v-for="c in colorPresets" :key="c" type="button"
                    @click="eForm.color = c"
                    class="w-7 h-7 rounded-full border-2 transition-all"
                    :style="{ background: c, borderColor: eForm.color === c ? '#1E2A38' : 'transparent' }"
                  ></button>
                </div>
              </div>
            </div>
            <div class="flex items-center gap-2.5">
              <input v-model="eForm.is_recurring" type="checkbox" id="ev_recurring" class="w-4 h-4 accent-terra" />
              <label for="ev_recurring" class="text-sm font-medium text-navy cursor-pointer">Acara berulang tahunan</label>
            </div>

            <div v-if="formError" class="form-error">{{ formError }}</div>

            <div class="flex gap-3 pt-1">
              <button type="button" @click="showForm = false" class="btn-secondary flex-1">Batal</button>
              <button type="submit" :disabled="saving" class="btn-primary flex-1">
                <span v-if="saving" class="spinner w-4 h-4 border-white/30 border-t-white"></span>
                {{ saving ? 'Menyimpan...' : (editingId ? 'Simpan Perubahan' : 'Buat Acara') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useEventsStore } from '@/stores/events.js'
import { useAuthStore }   from '@/stores/auth.js'
import { useToast }       from '@/composables/useToast.js'
import { fmtDateTime }    from '@/utils/format.js'
import api from '@/api/axios.js'

const eventsStore = useEventsStore()
const auth        = useAuthStore()
const toast       = useToast()
const showForm    = ref(false)
const saving      = ref(false)
const formError   = ref('')
const editingId   = ref(null)

const colorPresets = ['#CC6649', '#1E2A38', '#3B82F6', '#10B981', '#8B5CF6', '#F59E0B']

const blankForm = () => ({
  title: '', description: '', location: '',
  start_at: '', end_at: '', is_recurring: false, color: '#CC6649'
})
const eForm = ref(blankForm())

// ── RSVP state (keyed by event_id) ──────────────────────────────
const myRSVP = reactive({})  // { [eventId]: 'yes' | 'no' | 'pending' }

const rsvpOptions = [
  { value: 'yes',     label: '✓ Hadir',   activeClass: 'border-green-400 bg-green-50 text-green-700' },
  { value: 'no',      label: '✗ Tidak',   activeClass: 'border-red-400 bg-red-50 text-red-700' },
  { value: 'pending', label: '? Belum',   activeClass: 'border-amber-400 bg-amber-50 text-amber-700' },
]

async function setRSVP(eventId, rsvp) {
  const userMemberId = auth.user?.id  // user id — we use user.id as the family_member_id proxy
  if (!userMemberId) return
  try {
    await api.put(`/events/${eventId}/attendees/${userMemberId}`, { rsvp })
    myRSVP[eventId] = rsvp
    toast.success(rsvp === 'yes' ? 'Konfirmasi hadir berhasil' : 'RSVP berhasil diperbarui')
    await eventsStore.fetchEvents()  // refresh attendee_count
  } catch (e) {
    toast.error(e.response?.data?.error ?? 'Gagal memperbarui RSVP')
  }
}

// ── Form helpers ────────────────────────────────────────────────
function openCreate() {
  editingId.value = null
  eForm.value     = blankForm()
  formError.value = ''
  showForm.value  = true
}

function openEdit(event) {
  editingId.value = event.id
  // Convert ISO to datetime-local format (YYYY-MM-DDTHH:MM)
  const toLocal = iso => iso ? iso.slice(0, 16) : ''
  eForm.value = {
    title:        event.title,
    description:  event.description ?? '',
    location:     event.location    ?? '',
    start_at:     toLocal(event.start_at),
    end_at:       toLocal(event.end_at),
    is_recurring: event.is_recurring,
    color:        event.color ?? '#CC6649',
  }
  formError.value = ''
  showForm.value  = true
}

const fmtDay   = d => new Date(d).getDate()
const fmtMonth = d => new Intl.DateTimeFormat('id-ID', { month: 'short' }).format(new Date(d))

async function handleSave() {
  saving.value    = true
  formError.value = ''
  const payload = {
    ...eForm.value,
    description: eForm.value.description || null,
    location:    eForm.value.location    || null,
    end_at:      eForm.value.end_at      || null,
  }
  try {
    if (editingId.value) {
      await eventsStore.updateEvent(editingId.value, payload)
      toast.success('Acara berhasil diperbarui')
    } else {
      await eventsStore.createEvent(payload)
      toast.success('Acara berhasil dibuat')
    }
    showForm.value = false
  } catch (e) {
    formError.value = e.response?.data?.error ?? 'Gagal menyimpan acara.'
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  if (confirm('Hapus acara ini?')) {
    await eventsStore.deleteEvent(id)
    toast.success('Acara berhasil dihapus')
  }
}

onMounted(() => eventsStore.fetchEvents())
</script>
