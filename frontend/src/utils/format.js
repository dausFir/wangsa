/**
 * Shared formatting utilities — import these instead of re-defining per view.
 * Usage: import { fmtRp, fmtDate, fmtDateTime } from '@/utils/format.js'
 */

const rpFormatter = new Intl.NumberFormat('id-ID', {
  style: 'currency',
  currency: 'IDR',
  maximumFractionDigits: 0
})

const dateFormatter = new Intl.DateTimeFormat('id-ID', {
  day: 'numeric',
  month: 'short',
  year: 'numeric'
})

const dateTimeFormatter = new Intl.DateTimeFormat('id-ID', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit'
})

export function fmtRp(n) {
  return rpFormatter.format(n ?? 0)
}

export function fmtDate(d) {
  if (!d) return '—'
  return dateFormatter.format(new Date(d))
}

export function fmtDateTime(d) {
  if (!d) return '—'
  return dateTimeFormatter.format(new Date(d))
}

export function relativeTime(d) {
  if (!d) return '—'
  
  const date = new Date(d)
  const now = new Date()
  const diffInSeconds = Math.floor((now - date) / 1000)
  
  if (diffInSeconds < 60) return 'Baru saja'
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} menit lalu`
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} jam lalu`
  if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)} hari lalu`
  
  // For older dates, just show the formatted date
  return fmtDate(d)
}

// Export as object for convenience
export const format = {
  rp: fmtRp,
  date: fmtDate,
  dateTime: fmtDateTime,
  relativeTime
}
