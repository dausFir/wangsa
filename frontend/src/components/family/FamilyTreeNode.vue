<template>
  <div class="flex flex-col items-center">

    <!-- ── Member Card ── -->
    <div
      class="tree-node-card group relative"
      :class="{ selected }"
      @click="$emit('select', member)"
    >
      <div class="flex items-center gap-3">
        <!-- Avatar -->
        <div
          :class="[
            'w-11 h-11 rounded-full flex-shrink-0 flex items-center justify-center',
            'font-semibold text-sm select-none overflow-hidden',
            member.gender === 'male'   ? 'bg-blue-100 text-blue-700'  : 'bg-pink-100 text-pink-700',
            member.death_date          ? 'opacity-50 grayscale'        : '',
          ]"
        >
          <img
            v-if="sanitizePhotoUrl(member.photo_url)"
            :src="sanitizePhotoUrl(member.photo_url)"
            :alt="member.full_name"
            class="w-full h-full object-cover"
            @error="e => e.target.style.display = 'none'"
          />
          <span v-else>{{ initials }}</span>
        </div>

        <!-- Info -->
        <div class="min-w-0 flex-1">
          <p class="font-semibold text-navy text-sm leading-tight truncate max-w-[130px]">
            {{ member.full_name }}
          </p>
          <p class="text-xs text-warm-gray-400 mt-0.5 truncate">
            <template v-if="member.nickname">"{{ member.nickname }}" · </template>
            {{ birthYear ?? '—' }}
          </p>
          <!-- Spouse pills -->
          <div v-if="member.spouses?.length" class="flex gap-1 mt-1 flex-wrap">
            <span
              v-for="s in member.spouses.slice(0, 2)"
              :key="s.id"
              class="inline-flex items-center gap-1 text-[10px] bg-pink-50 text-pink-500
                     px-1.5 py-0.5 rounded-md leading-tight border border-pink-100"
              :title="`Pasangan: ${s.full_name}`"
            >♥ {{ s.full_name?.split(' ')[0] }}</span>
            <span v-if="member.spouses.length > 2" class="text-[10px] text-warm-gray-400 px-1">
              +{{ member.spouses.length - 2 }}
            </span>
          </div>
        </div>
      </div>

      <!-- Death indicator -->
      <span
        v-if="member.death_date"
        class="absolute top-2 right-2 text-warm-gray-300 text-xs leading-none"
        title="Telah meninggal"
      >✝</span>

      <!-- Hover action buttons -->
      <div class="absolute -top-2.5 -right-2.5 hidden group-hover:flex gap-1 z-10">
        <button
          @click.stop="$emit('add-child', member)"
          class="w-6 h-6 rounded-full bg-terra text-white text-xs font-bold
                 flex items-center justify-center shadow-soft
                 hover:bg-terra-dark transition-colors"
          title="Tambah anak"
        >+</button>
        <button
          @click.stop="$emit('edit', member)"
          class="w-6 h-6 rounded-full bg-navy text-white text-[11px]
                 flex items-center justify-center shadow-soft
                 hover:bg-navy-light transition-colors"
          title="Edit anggota"
        >✎</button>
      </div>
    </div>

    <!-- ── Children subtree ── -->
    <template v-if="member.children?.length">
      <!-- Vertical stem down from card -->
      <div class="w-0.5 h-8 bg-warm-gray-300 flex-shrink-0"></div>

      <!-- Children row -->
      <div class="relative flex items-start" :style="{ gap: childGap }">
        <!-- Horizontal crossbar connecting children -->
        <div
          v-if="member.children.length > 1"
          class="absolute top-0 h-0.5 bg-warm-gray-300 pointer-events-none"
          :style="{ left: '50%', transform: 'translateX(-50%)', width: `calc(100% - ${parseInt(childGap) * 2}px)` }"
        ></div>

        <!-- Recursive child nodes -->
        <div
          v-for="child in member.children"
          :key="child.id"
          class="flex flex-col items-center"
        >
          <div class="w-0.5 h-8 bg-warm-gray-300 flex-shrink-0"></div>
          <FamilyTreeNode
            :member="child"
            :selected-id="selectedId"
            @select="$emit('select', $event)"
            @add-child="$emit('add-child', $event)"
            @edit="$emit('edit', $event)"
          />
        </div>
      </div>
    </template>

  </div>
</template>

<script setup>
import { computed } from 'vue'
import { sanitizePhotoUrl } from '@/utils/sanitize.js'

const props = defineProps({
  member:     { type: Object,  required: true },
  selectedId: { type: Number,  default: null  },
})

defineEmits(['select', 'add-child', 'edit'])

const selected  = computed(() => props.member.id === props.selectedId)
const initials  = computed(() =>
  (props.member.full_name ?? '??')
    .split(' ')
    .filter(Boolean)
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
)
const birthYear = computed(() =>
  props.member.birth_date
    ? new Date(props.member.birth_date).getFullYear()
    : null
)

// Tighter gap for large families, wider for small ones
const childGap = computed(() => {
  const count = props.member.children?.length ?? 0
  if (count > 5) return '24px'
  if (count > 3) return '32px'
  return '40px'
})
</script>
