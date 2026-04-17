import { ref, reactive, onMounted, onUnmounted } from 'vue'

/**
 * useZoomPan — wheel zoom + mouse/touch drag pan untuk elemen container.
 *
 * Usage:
 *   const { transform, containerRef, contentRef, resetZoom } = useZoomPan()
 *   <div ref="containerRef" style="overflow:hidden; cursor:grab">
 *     <div ref="contentRef" :style="transform">...</div>
 *   </div>
 */
export function useZoomPan({ minScale = 0.2, maxScale = 3, zoomSpeed = 0.001 } = {}) {
  const containerRef = ref(null)
  const contentRef   = ref(null)

  const state = reactive({ scale: 1, x: 0, y: 0 })
  const dragging = ref(false)
  const lastPos  = reactive({ x: 0, y: 0 })

  const transform = ref('transform: translate(0px, 0px) scale(1); transform-origin: 0 0;')

  function applyTransform() {
    transform.value = `transform: translate(${state.x}px, ${state.y}px) scale(${state.scale}); transform-origin: 0 0; transition: none;`
  }

  function clampScale(s) {
    return Math.min(maxScale, Math.max(minScale, s))
  }

  // Wheel zoom — zoom toward cursor position
  function onWheel(e) {
    e.preventDefault()
    const rect = containerRef.value.getBoundingClientRect()
    const mouseX = e.clientX - rect.left
    const mouseY = e.clientY - rect.top

    const delta = -e.deltaY * zoomSpeed
    const newScale = clampScale(state.scale + delta * state.scale)
    const ratio = newScale / state.scale

    state.x = mouseX - ratio * (mouseX - state.x)
    state.y = mouseY - ratio * (mouseY - state.y)
    state.scale = newScale
    applyTransform()
  }

  // Mouse drag
  function onMouseDown(e) {
    if (e.button !== 0) return
    dragging.value = true
    lastPos.x = e.clientX
    lastPos.y = e.clientY
    containerRef.value.style.cursor = 'grabbing'
  }
  function onMouseMove(e) {
    if (!dragging.value) return
    state.x += e.clientX - lastPos.x
    state.y += e.clientY - lastPos.y
    lastPos.x = e.clientX
    lastPos.y = e.clientY
    applyTransform()
  }
  function onMouseUp() {
    dragging.value = false
    if (containerRef.value) containerRef.value.style.cursor = 'grab'
  }

  // Touch pinch + pan
  let lastTouches = []
  function onTouchStart(e) {
    lastTouches = Array.from(e.touches)
  }
  function onTouchMove(e) {
    e.preventDefault()
    const touches = Array.from(e.touches)
    if (touches.length === 1 && lastTouches.length === 1) {
      // Single finger pan
      state.x += touches[0].clientX - lastTouches[0].clientX
      state.y += touches[0].clientY - lastTouches[0].clientY
      applyTransform()
    } else if (touches.length === 2 && lastTouches.length === 2) {
      // Two finger pinch zoom
      const prevDist = Math.hypot(lastTouches[1].clientX - lastTouches[0].clientX, lastTouches[1].clientY - lastTouches[0].clientY)
      const currDist = Math.hypot(touches[1].clientX - touches[0].clientX, touches[1].clientY - touches[0].clientY)
      const ratio = clampScale(state.scale * (currDist / prevDist)) / state.scale

      const rect = containerRef.value.getBoundingClientRect()
      const midX = ((touches[0].clientX + touches[1].clientX) / 2) - rect.left
      const midY = ((touches[0].clientY + touches[1].clientY) / 2) - rect.top

      state.x = midX - ratio * (midX - state.x)
      state.y = midY - ratio * (midY - state.y)
      state.scale = clampScale(state.scale * (currDist / prevDist))
      applyTransform()
    }
    lastTouches = touches
  }

  function resetZoom() {
    state.scale = 1
    state.x = 0
    state.y = 0
    transform.value = `transform: translate(0px, 0px) scale(1); transform-origin: 0 0; transition: transform 0.3s ease;`
  }

  function fitToScreen() {
    if (!containerRef.value || !contentRef.value) return
    const cRect = containerRef.value.getBoundingClientRect()
    const iRect = contentRef.value.getBoundingClientRect()
    const scaleX = cRect.width  / (iRect.width  / state.scale)
    const scaleY = cRect.height / (iRect.height / state.scale)
    const newScale = clampScale(Math.min(scaleX, scaleY) * 0.9)
    state.scale = newScale
    state.x = (cRect.width  - (iRect.width  / state.scale) * newScale) / 2
    state.y = (cRect.height - (iRect.height / state.scale) * newScale) / 2
    transform.value = `transform: translate(${state.x}px, ${state.y}px) scale(${state.scale}); transform-origin: 0 0; transition: transform 0.3s ease;`
  }

  onMounted(() => {
    const el = containerRef.value
    if (!el) return
    el.addEventListener('wheel',      onWheel,     { passive: false })
    el.addEventListener('mousedown',  onMouseDown)
    el.addEventListener('touchstart', onTouchStart, { passive: true })
    el.addEventListener('touchmove',  onTouchMove,  { passive: false })
    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup',   onMouseUp)
  })

  onUnmounted(() => {
    const el = containerRef.value
    if (el) {
      el.removeEventListener('wheel',      onWheel)
      el.removeEventListener('mousedown',  onMouseDown)
      el.removeEventListener('touchstart', onTouchStart)
      el.removeEventListener('touchmove',  onTouchMove)
    }
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup',   onMouseUp)
  })

  return { transform, containerRef, contentRef, resetZoom, fitToScreen, state }
}
