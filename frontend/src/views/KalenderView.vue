<template>
  <div class="space-y-6">
    <!-- Calendar Header -->
    <div class="card">
      <div class="flex items-center justify-between mb-6 flex-wrap gap-3">
        <div class="flex items-center gap-4">
          <div>
            <h3 class="text-lg font-semibold text-navy">{{ formatMonthYear(currentDate) }}</h3>
            <p class="text-sm text-warm-gray-400 mt-0.5">{{ eventsStore.events.length }} acara bulan ini</p>
          </div>
          
          <!-- Month Navigation -->
          <div class="flex items-center gap-1">
            <button @click="previousMonth" class="btn-ghost p-2 text-warm-gray-400 hover:text-navy">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
            </button>
            <button @click="goToToday" class="btn-ghost px-3 py-1.5 text-xs font-medium">Hari Ini</button>
            <button @click="nextMonth" class="btn-ghost p-2 text-warm-gray-400 hover:text-navy">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </button>
          </div>
        </div>
        
        <!-- View Toggle & Create Button -->
        <div class="flex items-center gap-3">
          <div class="flex rounded-lg border border-warm-gray-200 bg-warm-gray-50">
            <button
              @click="viewMode = 'calendar'"
              :class="['px-3 py-1.5 text-sm font-medium rounded-l-lg transition-all',
                viewMode === 'calendar' 
                  ? 'bg-white text-navy shadow-sm' 
                  : 'text-warm-gray-400 hover:text-navy']"
            >
              📅 Kalender
            </button>
            <button
              @click="viewMode = 'list'"
              :class="['px-3 py-1.5 text-sm font-medium rounded-r-lg transition-all',
                viewMode === 'list' 
                  ? 'bg-white text-navy shadow-sm' 
                  : 'text-warm-gray-400 hover:text-navy']"
            >
              📋 Daftar
            </button>
          </div>
          
          <button @click="openCreate" class="btn-primary">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
            </svg>
            Buat Acara
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="eventsStore.loading" class="flex items-center justify-center py-12 gap-3 text-warm-gray-400">
        <div class="spinner w-6 h-6 border-warm-gray-200 border-t-warm-gray-400"></div>
        <span>Memuat acara...</span>
      </div>

      <!-- Calendar Grid View -->
      <div v-else-if="viewMode === 'calendar'" class="calendar-grid">
        <!-- Day Headers -->
        <div class="grid grid-cols-7 gap-0 mb-2">
          <div
            v-for="day in dayHeaders"
            :key="day"
            class="text-center text-xs font-medium text-warm-gray-400 uppercase tracking-wide py-2"
          >
            {{ day }}
          </div>
        </div>
        
        <!-- Calendar Days -->
        <div class="grid grid-cols-7 gap-0 border border-warm-gray-200 rounded-xl overflow-hidden">
          <div
            v-for="day in calendarDays"
            :key="day.dateKey"
            @click="selectDay(day)"
            :class="[
              'min-h-[120px] border-r border-b border-warm-gray-100 p-2 cursor-pointer transition-colors relative',
              'hover:bg-warm-gray-50',
              !day.isCurrentMonth ? 'bg-warm-gray-25 text-warm-gray-300' : '',
              day.isToday ? 'bg-blue-50 border-blue-200' : '',
              day.isSelected ? 'bg-terra-50 border-terra-200' : ''
            ]"
          >
            <!-- Date Number -->
            <div
              :class="[
                'flex items-center justify-center w-6 h-6 text-sm font-medium mb-1',
                day.isToday ? 'bg-blue-500 text-white rounded-full' : '',
                !day.isCurrentMonth ? 'text-warm-gray-300' : 'text-warm-gray-700'
              ]"
            >
              {{ day.date }}
            </div>
            
            <!-- Events for this day -->
            <div class="space-y-1">
              <div
                v-for="event in getEventsForDay(day)"
                :key="event.id"
                @click.stop="openEventDetail(event)"
                :style="{ backgroundColor: event.color || '#CC6649' }"
                class="text-white text-xs px-2 py-1 rounded-md cursor-pointer 
                       hover:opacity-80 transition-opacity truncate font-medium"
                :title="event.title + (event.location ? ' - ' + event.location : '')"
              >
                {{ formatEventTime(event) }} {{ event.title }}
              </div>
              
              <!-- Show more indicator -->
              <div
                v-if="getEventsForDay(day).length > 3"
                class="text-xs text-warm-gray-400 font-medium"
              >
                +{{ getEventsForDay(day).length - 3 }} lainnya
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- List View (existing) -->
      <div v-else-if="viewMode === 'list'">
        <div v-if="!eventsStore.events.length" class="text-center py-16">
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
    </div>

    <!-- Event Detail Modal -->
    <Teleport to="body">
      <div
        v-if="showEventDetail"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        style="background: rgba(20,29,39,0.52); backdrop-filter: blur(4px);"
        @mousedown.self="showEventDetail = false"
      >
        <div class="bg-white rounded-2xl shadow-modal w-full max-w-md">
          <div class="flex items-center justify-between px-6 py-5 border-b border-warm-gray-100">
            <h3 class="text-lg font-semibold text-navy">Detail Acara</h3>
            <button @click="showEventDetail = false" class="btn-ghost p-2 -mr-2">✕</button>
          </div>

          <div v-if="selectedEvent" class="px-6 py-5 space-y-4">
            <!-- Event Title -->
            <div>
              <div class="flex items-center gap-3 mb-2">
                <div
                  class="w-4 h-4 rounded-full flex-shrink-0"
                  :style="{ backgroundColor: selectedEvent.color || '#CC6649' }"
                ></div>
                <h4 class="font-semibold text-navy">{{ selectedEvent.title }}</h4>
                <span v-if="selectedEvent.is_recurring" class="badge-blue text-xs">🔁 Rutin</span>
              </div>
            </div>

            <!-- Event Time -->
            <div class="flex items-center gap-2 text-sm text-warm-gray-600">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span>{{ fmtDateTime(selectedEvent.start_at) }}</span>
              <span v-if="selectedEvent.end_at">- {{ fmtDateTime(selectedEvent.end_at) }}</span>
            </div>

            <!-- Location -->
            <div v-if="selectedEvent.location" class="flex items-center gap-2 text-sm text-warm-gray-600">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
              </svg>
              <span>{{ selectedEvent.location }}</span>
            </div>

            <!-- Description -->
            <div v-if="selectedEvent.description" class="text-sm text-warm-gray-600">
              <p class="font-medium text-navy mb-1">Deskripsi:</p>
              <p>{{ selectedEvent.description }}</p>
            </div>

            <!-- Attendees -->
            <div v-if="selectedEvent.attendee_count > 0" class="flex items-center gap-2 text-sm text-warm-gray-600">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
              </svg>
              <span>{{ selectedEvent.attendee_count }} orang hadir</span>
            </div>

            <!-- Notes -->
            <div v-if="selectedEvent.notes" class="border-t border-warm-gray-100 pt-4">
              <p class="font-medium text-navy mb-2 text-sm">Catatan / Notulensi:</p>
              <div class="prose prose-sm max-w-none text-warm-gray-600" v-html="selectedEvent.notes"></div>
            </div>

            <!-- RSVP Section -->
            <div class="border-t border-warm-gray-100 pt-4">
              <div class="flex items-center gap-2 mb-3">
                <span class="text-sm font-medium text-navy">RSVP saya:</span>
                <div class="flex gap-2">
                  <button
                    v-for="opt in rsvpOptions"
                    :key="opt.value"
                    @click="setRSVP(selectedEvent.id, opt.value)"
                    :class="['text-xs px-2.5 py-1 rounded-full border transition-all font-medium',
                      myRSVP[selectedEvent.id] === opt.value
                        ? opt.activeClass
                        : 'border-warm-gray-200 text-warm-gray-400 hover:border-warm-gray-400']"
                  >{{ opt.label }}</button>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div v-if="auth.isSuperAdmin" class="flex gap-3 pt-2">
              <button
                @click="openEdit(selectedEvent); showEventDetail = false"
                class="btn-secondary flex-1"
              >
                Edit Acara
              </button>
              <button
                @click="handleDelete(selectedEvent.id); showEventDetail = false"
                class="btn-ghost flex-1 text-red-600 hover:bg-red-50"
              >
                Hapus
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

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
            
            <!-- Notes Section -->
            <div>
              <RichTextEditor
                v-model="eForm.notes"
                label="Catatan / Notulensi"
                placeholder="Tulis catatan rapat, notulensi, atau informasi tambahan..."
              />
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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useEventsStore } from '@/stores/events.js'
import { useAuthStore }   from '@/stores/auth.js'
import { useToast }       from '@/composables/useToast.js'
import { fmtDateTime }    from '@/utils/format.js'
import RichTextEditor from '@/components/ui/RichTextEditor.vue'
import api from '@/api/axios.js'

const eventsStore = useEventsStore()
const auth        = useAuthStore()
const toast       = useToast()
const showForm    = ref(false)
const saving      = ref(false)
const formError   = ref('')
const editingId   = ref(null)
const viewMode    = ref('calendar') // 'calendar' or 'list'

// Calendar state
const currentDate = ref(new Date())
const selectedDay = ref(null)
const showEventDetail = ref(false)
const selectedEvent = ref(null)

const colorPresets = ['#CC6649', '#1E2A38', '#3B82F6', '#10B981', '#8B5CF6', '#F59E0B']

const blankForm = () => ({
  title: '', description: '', location: '',
  start_at: '', end_at: '', is_recurring: false, color: '#CC6649',
  notes: ''
})
const eForm = ref(blankForm())

// ── RSVP state (keyed by event_id) ──────────────────────────────
const myRSVP = reactive({})  // { [eventId]: 'yes' | 'no' | 'pending' }

const rsvpOptions = [
  { value: 'yes',     label: '✓ Hadir',   activeClass: 'border-green-400 bg-green-50 text-green-700' },
  { value: 'no',      label: '✗ Tidak',   activeClass: 'border-red-400 bg-red-50 text-red-700' },
  { value: 'pending', label: '? Belum',   activeClass: 'border-amber-400 bg-amber-50 text-amber-700' },
]

// ── Calendar computed properties ────────────────────────────────
const dayHeaders = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab']

const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  // First day of the month
  const firstDay = new Date(year, month, 1)
  const firstDayOfWeek = firstDay.getDay() // 0 = Sunday, 1 = Monday, etc
  
  // Last day of the month
  const lastDay = new Date(year, month + 1, 0)
  const daysInMonth = lastDay.getDate()
  
  // Previous month days needed to fill the grid
  const previousMonth = new Date(year, month, 0) // Last day of previous month
  const daysFromPrevMonth = firstDayOfWeek === 0 ? 6 : firstDayOfWeek - 1 // Monday = 0 offset
  
  // Next month days needed to fill the grid
  const totalCells = 42 // 6 weeks × 7 days
  const daysFromNextMonth = totalCells - daysInMonth - daysFromPrevMonth
  
  const days = []
  
  // Previous month days
  for (let i = daysFromPrevMonth; i > 0; i--) {
    const date = previousMonth.getDate() - i + 1
    const dayDate = new Date(previousMonth.getFullYear(), previousMonth.getMonth(), date)
    days.push({
      date,
      dateKey: formatDateKey(dayDate),
      fullDate: dayDate,
      isCurrentMonth: false,
      isToday: isSameDay(dayDate, new Date()),
      isSelected: selectedDay.value && isSameDay(dayDate, selectedDay.value)
    })
  }
  
  // Current month days
  for (let date = 1; date <= daysInMonth; date++) {
    const dayDate = new Date(year, month, date)
    days.push({
      date,
      dateKey: formatDateKey(dayDate),
      fullDate: dayDate,
      isCurrentMonth: true,
      isToday: isSameDay(dayDate, new Date()),
      isSelected: selectedDay.value && isSameDay(dayDate, selectedDay.value)
    })
  }
  
  // Next month days
  for (let date = 1; date <= daysFromNextMonth; date++) {
    const dayDate = new Date(year, month + 1, date)
    days.push({
      date,
      dateKey: formatDateKey(dayDate),
      fullDate: dayDate,
      isCurrentMonth: false,
      isToday: isSameDay(dayDate, new Date()),
      isSelected: selectedDay.value && isSameDay(dayDate, selectedDay.value)
    })
  }
  
  return days
})

// ── Calendar helper functions ───────────────────────────────────
function formatDateKey(date) {
  return date.toISOString().split('T')[0]
}

function isSameDay(date1, date2) {
  return date1.toDateString() === date2.toDateString()
}

function formatMonthYear(date) {
  return new Intl.DateTimeFormat('id-ID', { month: 'long', year: 'numeric' }).format(date)
}

function formatEventTime(event) {
  const startTime = new Date(event.start_at)
  return startTime.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', hour12: false })
}

function getEventsForDay(day) {
  return eventsStore.events.filter(event => {
    const eventDate = new Date(event.start_at)
    return isSameDay(eventDate, day.fullDate)
  }).slice(0, 3) // Limit to 3 events shown
}

// ── Calendar navigation ─────────────────────────────────────────
function previousMonth() {
  currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() - 1, 1)
}

function nextMonth() {
  currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1, 1)
}

function goToToday() {
  currentDate.value = new Date()
  selectedDay.value = new Date()
}

function selectDay(day) {
  selectedDay.value = day.fullDate
  // If day has no events and user is admin, open create form with pre-filled date
  if (getEventsForDay(day).length === 0 && auth.isSuperAdmin) {
    openCreateForDate(day.fullDate)
  }
}

function openCreateForDate(date) {
  editingId.value = null
  const isoDate = date.toISOString().slice(0, 16) // YYYY-MM-DDTHH:MM format
  eForm.value = {
    ...blankForm(),
    start_at: isoDate
  }
  formError.value = ''
  showForm.value = true
}

function openEventDetail(event) {
  selectedEvent.value = event
  showEventDetail.value = true
}

// ── Watch for month changes to fetch events ────────────────────
watch(currentDate, () => {
  fetchEventsForMonth()
}, { immediate: true })

function fetchEventsForMonth() {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  // Get first day of month and last day of month
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0, 23, 59, 59, 999) // End of last day
  
  // Format dates for API
  const from = firstDay.toISOString().split('T')[0] + 'T00:00:00'
  const to = lastDay.toISOString()
  
  eventsStore.fetchEvents(from, to)
}

// ── RSVP functionality ──────────────────────────────────────────
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
    notes:        event.notes ?? '',
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
    // Refresh events for current month
    fetchEventsForMonth()
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
    fetchEventsForMonth()
  }
}

onMounted(() => {
  fetchEventsForMonth()
})
</script>

<style scoped>
.calendar-grid {
  /* Calendar specific styles */
}

.calendar-day {
  transition: all 0.2s ease;
}

.calendar-day:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Custom scrollbar for calendar */
.calendar-grid::-webkit-scrollbar {
  width: 6px;
}

.calendar-grid::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 3px;
}

.calendar-grid::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

.calendar-grid::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Prose styling for notes */
:deep(.prose) {
  font-size: 14px;
  line-height: 1.6;
}

:deep(.prose h1) {
  font-size: 16px;
  font-weight: 600;
  color: #1E2A38;
  margin-top: 0;
  margin-bottom: 8px;
}

:deep(.prose h2) {
  font-size: 15px;
  font-weight: 600;
  color: #1E2A38;
  margin-top: 12px;
  margin-bottom: 6px;
}

:deep(.prose h3) {
  font-size: 14px;
  font-weight: 600;
  color: #1E2A38;
  margin-top: 10px;
  margin-bottom: 4px;
}

:deep(.prose p) {
  margin-top: 0;
  margin-bottom: 8px;
}

:deep(.prose ul, .prose ol) {
  margin-top: 8px;
  margin-bottom: 8px;
  padding-left: 16px;
}

:deep(.prose li) {
  margin-top: 2px;
  margin-bottom: 2px;
}

:deep(.prose strong) {
  color: #1E2A38;
  font-weight: 600;
}

:deep(.prose a) {
  color: #3B82F6;
  text-decoration: underline;
}

:deep(.prose a:hover) {
  color: #2563EB;
}
</style>
