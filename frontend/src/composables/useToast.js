import { useToastStore } from '@/stores/toast.js'

// Convenience composable — import this instead of the store directly
// Usage: const toast = useToast(); toast.success('Berhasil!')
export function useToast() {
  return useToastStore()
}
