<template>
  <div class="space-y-6">

    <!-- ── Header bar ── -->
    <div class="flex flex-wrap items-center justify-between gap-4">
      <p class="text-warm-gray-500 text-sm">
        <span class="font-semibold text-navy">{{ family.members.length }}</span> anggota terdaftar
      </p>
      <div class="flex gap-3">
        <button @click="showMarriageForm = !showMarriageForm" class="btn-secondary text-sm">
          ♥ Pernikahan
        </button>
        <button @click="openCreate" class="btn-primary">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
          </svg>
          Tambah Anggota
        </button>
      </div>
    </div>

    <!-- ── Search & filter bar ── -->
    <div class="flex flex-wrap gap-3 items-center">
      <!-- Search input -->
      <div class="relative flex-1 min-w-[200px]">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-warm-gray-400 pointer-events-none" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
        <input
          v-model="searchQuery"
          class="form-input pl-10"
          placeholder="Cari nama anggota..."
          type="search"
        />
      </div>

      <!-- Gender filter -->
      <select v-model="filterGender" class="form-select w-auto">
        <option value="">Semua jenis kelamin</option>
        <option value="male">♂ Laki-laki</option>
        <option value="female">♀ Perempuan</option>
      </select>

      <!-- Status filter -->
      <select v-model="filterStatus" class="form-select w-auto">
        <option value="">Semua status</option>
        <option value="alive">Masih hidup</option>
        <option value="deceased">Sudah meninggal</option>
      </select>

      <!-- Result count / clear -->
      <div v-if="isFiltering" class="flex items-center gap-2 text-sm text-warm-gray-500">
        <span>{{ filteredMembers.length }} hasil</span>
        <button @click="clearFilters" class="text-terra hover:underline text-xs">Reset</button>
      </div>
    </div>

    <!-- ── Marriage mini-form ── -->
    <Transition name="slide-up">
      <div v-if="showMarriageForm" class="card border border-pink-100 bg-pink-50/30">
        <h4 class="font-semibold text-navy mb-4">♥ Tambah Data Pernikahan</h4>
        <form @submit.prevent="submitMarriage" class="space-y-4">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="form-label">Suami <span class="text-terra">*</span></label>
              <select v-model="marriageForm.husband_id" class="form-select" required>
                <option :value="null">— Pilih —</option>
                <option v-for="m in maleMembers" :key="m.id" :value="m.id">{{ m.full_name }}</option>
              </select>
            </div>
            <div>
              <label class="form-label">Istri <span class="text-terra">*</span></label>
              <select v-model="marriageForm.wife_id" class="form-select" required>
                <option :value="null">— Pilih —</option>
                <option v-for="m in femaleMembers" :key="m.id" :value="m.id">{{ m.full_name }}</option>
              </select>
            </div>
            <div>
              <label class="form-label">Tanggal Menikah</label>
              <input v-model="marriageForm.marriage_date" type="date" class="form-input" />
            </div>
            <div>
              <label class="form-label">Tanggal Cerai</label>
              <input v-model="marriageForm.divorce_date" type="date" class="form-input" />
            </div>
          </div>
          <div v-if="marriageError" class="form-error">{{ marriageError }}</div>
          <div class="flex gap-3">
            <button type="button" @click="showMarriageForm = false" class="btn-secondary">Tutup</button>
            <button type="submit" :disabled="marriageSaving" class="btn-primary">
              <span v-if="marriageSaving" class="spinner w-4 h-4 border-white/30 border-t-white"></span>
              {{ marriageSaving ? 'Menyimpan...' : 'Simpan Pernikahan' }}
            </button>
          </div>
        </form>
      </div>
    </Transition>

    <!-- ── Loading ── -->
    <div v-if="family.loading" class="card flex flex-col items-center justify-center py-24 gap-4">
      <div class="spinner w-10 h-10 border-terra/20 border-t-terra border-[3px]"></div>
      <p class="text-warm-gray-400 text-sm">Memuat pohon silsilah...</p>
    </div>

    <!-- ── Search results (flat list) ── -->
    <div v-else-if="isFiltering" class="card divide-y divide-warm-gray-50">
      <div v-if="!filteredMembers.length" class="py-12 text-center text-warm-gray-400 text-sm">
        Tidak ada anggota yang cocok dengan pencarian.
      </div>
      <div
        v-for="m in filteredMembers"
        :key="m.id"
        class="flex items-center gap-4 px-2 py-3 hover:bg-warm-gray-50 rounded-xl cursor-pointer transition-colors"
        @click="selectedMember = m"
      >
        <div :class="['w-10 h-10 rounded-xl flex items-center justify-center font-semibold text-sm flex-shrink-0', m.gender === 'male' ? 'bg-blue-100 text-blue-700' : 'bg-pink-100 text-pink-600']">
          {{ avatarInitials(m) }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="font-medium text-navy text-sm truncate">{{ m.full_name }}</p>
          <p class="text-xs text-warm-gray-400">
            {{ m.gender === 'male' ? 'Laki-laki' : 'Perempuan' }}
            <span v-if="m.birth_date"> · {{ fmtDate(m.birth_date) }}</span>
            <span v-if="m.death_date" class="text-red-400"> · ✝ {{ fmtDate(m.death_date) }}</span>
          </p>
        </div>
        <svg class="w-4 h-4 text-warm-gray-300 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7"/>
        </svg>
      </div>
    </div>

    <!-- ── Empty (no members at all) ── -->
    <div
      v-else-if="!family.tree.length"
      class="card flex flex-col items-center justify-center py-24 text-center"
    >
      <div class="text-7xl mb-6">🌳</div>
      <h3 class="text-2xl font-bold text-navy mb-2">Silsilah Belum Ada</h3>
      <p class="text-warm-gray-500 mb-7 max-w-sm">Mulailah dengan menambahkan anggota keluarga pertama sebagai akar pohon silsilah.</p>
      <button @click="openCreate" class="btn-primary">Tambah Anggota Pertama</button>
    </div>

    <!-- ── Tree canvas with D3.js interactive visualization ── -->
    <div v-else class="card p-0 overflow-hidden relative" style="height: 75vh;">
      <InteractiveFamilyTree
        :tree-data="family.tree"
        :selected-member-id="selectedMember?.id"
        @member-select="selectedMember = $event"
        @member-edit="openEdit"
        @add-child="openAddChild"
      />
    </div>

    <!-- ── Detail panel ── -->
    <Transition name="slide-up">
      <div v-if="selectedMember" class="card border-2 border-terra/20">
        <div class="flex items-start justify-between gap-4 flex-wrap">
          <div class="flex items-center gap-4">
            <div :class="['w-16 h-16 rounded-2xl flex items-center justify-center text-2xl font-bold flex-shrink-0', selectedMember.gender === 'male' ? 'bg-blue-100 text-blue-700' : 'bg-pink-100 text-pink-600']">
              {{ avatarInitials(selectedMember) }}
            </div>
            <div>
              <h3 class="text-xl font-bold text-navy">{{ selectedMember.full_name }}</h3>
              <p v-if="selectedMember.nickname" class="text-warm-gray-500 text-sm mb-2">"{{ selectedMember.nickname }}"</p>
              <div class="flex flex-wrap gap-2">
                <span class="badge-gray capitalize">{{ selectedMember.gender === 'male' ? '♂ Laki-laki' : '♀ Perempuan' }}</span>
                <span v-if="selectedMember.birth_date" class="badge-blue">📅 {{ fmtDate(selectedMember.birth_date) }}</span>
                <span v-if="selectedMember.birth_place" class="badge-gray">📍 {{ selectedMember.birth_place }}</span>
                <span v-if="selectedMember.death_date" class="badge-red">✝ {{ fmtDate(selectedMember.death_date) }}</span>
              </div>
            </div>
          </div>
          <div class="flex gap-2 flex-shrink-0">
            <button @click="openEdit(selectedMember)" class="btn-secondary text-sm">Edit</button>
            <button v-if="auth.isSuperAdmin" @click="confirmDelete(selectedMember)" class="btn-danger text-sm">Hapus</button>
            <button @click="selectedMember = null" class="btn-ghost">✕</button>
          </div>
        </div>

        <div v-if="selectedMember.spouses?.length" class="mt-5 pt-4 border-t border-warm-gray-100">
          <p class="text-xs font-semibold text-warm-gray-400 uppercase tracking-wider mb-2.5">Pasangan</p>
          <div class="flex flex-wrap gap-2">
            <div v-for="s in selectedMember.spouses" :key="s.id" class="flex items-center gap-2 px-3 py-2 bg-pink-50 rounded-xl border border-pink-100">
              <span class="text-pink-400 text-sm">♥</span>
              <span class="text-sm font-medium text-navy">{{ s.full_name }}</span>
            </div>
          </div>
        </div>

        <p v-if="selectedMember.notes" class="mt-4 text-sm text-warm-gray-600 italic bg-warm-gray-50 px-4 py-3 rounded-xl border border-warm-gray-100">
          {{ selectedMember.notes }}
        </p>
      </div>
    </Transition>

    <!-- ── Member Form Modal ── -->
    <MemberFormModal
      v-if="showMemberForm"
      :member="editingMember"
      :default-parent-id="addChildParentId"
      :all-members="family.members"
      @close="closeMemberModal"
      @saved="closeMemberModal"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useFamilyStore } from '@/stores/family.js'
import { useAuthStore }   from '@/stores/auth.js'
import { useToast }       from '@/composables/useToast.js'
import { fmtDate }        from '@/utils/format.js'
import InteractiveFamilyTree from '@/components/family/InteractiveFamilyTree.vue'
import MemberFormModal  from '@/components/family/MemberFormModal.vue'

const family = useFamilyStore()
const auth   = useAuthStore()
const toast  = useToast()

const selectedMember   = ref(null)
const showMemberForm   = ref(false)
const editingMember    = ref(null)
const addChildParentId = ref(null)
const showMarriageForm = ref(false)
const marriageSaving   = ref(false)
const marriageError    = ref('')
const marriageForm     = ref({ husband_id: null, wife_id: null, marriage_date: null, divorce_date: null })

// ── Search & filter state ──
const searchQuery  = ref('')
const filterGender = ref('')
const filterStatus = ref('')

const isFiltering = computed(() =>
  searchQuery.value.trim() !== '' || filterGender.value !== '' || filterStatus.value !== ''
)

const filteredMembers = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  return family.members.filter(m => {
    // Name matching
    const matchName = !q || m.full_name.toLowerCase().includes(q) || (m.nickname ?? '').toLowerCase().includes(q)
    
    // Gender matching
    const matchGender = !filterGender.value || m.gender === filterGender.value
    
    // Status matching - fixed logic
    let matchStatus = true
    if (filterStatus.value === 'alive') {
      matchStatus = !m.death_date
    } else if (filterStatus.value === 'deceased') {
      matchStatus = !!m.death_date
    }
    // If filterStatus.value is empty string (Semua status), matchStatus stays true
    
    return matchName && matchGender && matchStatus
  })
})

function clearFilters() {
  searchQuery.value  = ''
  filterGender.value = ''
  filterStatus.value = ''
}

const maleMembers   = computed(() => family.members.filter(m => m.gender === 'male'))
const femaleMembers = computed(() => family.members.filter(m => m.gender === 'female'))

function avatarInitials(m) {
  return (m.full_name ?? '??').split(' ').filter(Boolean).map(w => w[0]).join('').toUpperCase().slice(0, 2)
}

function openCreate()         { editingMember.value = null; addChildParentId.value = null; showMemberForm.value = true }
function openAddChild(parent) { editingMember.value = null; addChildParentId.value = parent.id; showMemberForm.value = true }
function openEdit(member)     { editingMember.value = member; addChildParentId.value = null; showMemberForm.value = true }
function closeMemberModal()   { showMemberForm.value = false; editingMember.value = null; addChildParentId.value = null }

async function confirmDelete(member) {
  if (!confirm(`Hapus "${member.full_name}" dari silsilah?\n\nPeringatan: anak-anak dari anggota ini tidak akan ikut terhapus.`)) return
  await family.deleteMember(member.id)
  toast.success(`"${member.full_name}" berhasil dihapus`)
  selectedMember.value = null
}

async function submitMarriage() {
  marriageSaving.value = true
  marriageError.value  = ''
  try {
    await family.createMarriage({
      ...marriageForm.value,
      marriage_date: marriageForm.value.marriage_date || null,
      divorce_date:  marriageForm.value.divorce_date  || null,
    })
    toast.success('Data pernikahan berhasil ditambahkan')
    showMarriageForm.value = false
    marriageForm.value = { husband_id: null, wife_id: null, marriage_date: null, divorce_date: null }
  } catch (e) {
    marriageError.value = e.response?.data?.error ?? 'Gagal menyimpan data pernikahan.'
  } finally {
    marriageSaving.value = false
  }
}

onMounted(async () => {
  await Promise.all([family.fetchTree(), family.fetchMembers()])
})
</script>