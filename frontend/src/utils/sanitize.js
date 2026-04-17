/**
 * sanitizePhotoUrl — strips dangerous URL schemes before binding to <img src>.
 *
 * Even though the backend validates photo_url, a defense-in-depth approach
 * means we also sanitize on the frontend before any URL hits the DOM.
 *
 * Allowed: https://, http://, data:image/...;base64,...
 * Blocked: javascript:, vbscript:, data:text/html, blob:, etc.
 *
 * Returns null if the URL is unsafe so the component falls back to initials.
 */
export function sanitizePhotoUrl(url) {
  if (!url || typeof url !== 'string') return null

  const trimmed = url.trim()
  const lower   = trimmed.toLowerCase()

  // Allow safe data URLs (base64 images only)
  const safeDataPrefixes = [
    'data:image/jpeg;base64,',
    'data:image/jpg;base64,',
    'data:image/png;base64,',
    'data:image/gif;base64,',
    'data:image/webp;base64,',
  ]
  if (lower.startsWith('data:')) {
    if (safeDataPrefixes.some(p => lower.startsWith(p))) {
      return trimmed
    }
    return null  // data:text/html, data:application/*, etc. — blocked
  }

  // Allow http(s) URLs
  if (lower.startsWith('https://') || lower.startsWith('http://')) {
    return trimmed
  }

  // Block everything else: javascript:, vbscript:, blob:, file:, //evil.com, etc.
  return null
}
