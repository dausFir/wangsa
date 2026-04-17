import { computed, ref } from 'vue'

/**
 * useDirtyGuard — deteksi perubahan form, konfirmasi sebelum tutup.
 *
 * Usage:
 *   const { isDirty, setInitial, confirmClose } = useDirtyGuard(form)
 *   - Panggil setInitial(form.value) sekali setelah form diinisialisasi
 *   - Ganti @mousedown.self="$emit('close')" dan @click="$emit('close')"
 *     dengan @mousedown.self="confirmClose" dan @click="confirmClose"
 *   - confirmClose() akan meminta konfirmasi jika form sudah diubah
 */
export function useDirtyGuard(formRef, emit) {
  const initialSnapshot = ref(null)

  function setInitial(values) {
    initialSnapshot.value = JSON.stringify(values)
  }

  const isDirty = computed(() => {
    if (initialSnapshot.value === null) return false
    return JSON.stringify(formRef.value) !== initialSnapshot.value
  })

  function confirmClose() {
    if (isDirty.value) {
      if (!confirm('Yakin ingin keluar?\nData yang sudah diisi akan hilang.')) return
    }
    emit('close')
  }

  return { isDirty, setInitial, confirmClose }
}
