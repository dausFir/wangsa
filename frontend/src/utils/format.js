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
