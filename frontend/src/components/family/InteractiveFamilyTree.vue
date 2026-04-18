<template>
  <div ref="containerRef" class="family-tree-container">
    <div class="family-tree-controls">
      <button @click="centerTree" class="control-btn" title="Center Tree">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"/>
        </svg>
      </button>
      <button @click="zoomIn" class="control-btn" title="Zoom In">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
        </svg>
      </button>
      <button @click="zoomOut" class="control-btn" title="Zoom Out">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 12H6"/>
        </svg>
      </button>
      <div class="zoom-level">{{ Math.round(zoomLevel * 100) }}%</div>
      <button @click="toggleLayout" class="control-btn" :title="layoutMode === 'vertical' ? 'Switch to Horizontal' : 'Switch to Vertical'">
        {{ layoutMode === 'vertical' ? '↔' : '↕' }}
      </button>
    </div>
    <svg ref="svgRef" class="family-tree-svg"></svg>
    
    <!-- Enhanced Node tooltip -->
    <div 
      v-if="tooltip.visible" 
      :style="{ left: tooltip.x + 'px', top: tooltip.y + 'px' }"
      class="family-tree-tooltip"
    >
      <div class="tooltip-content">
        <!-- Header with name and gender icon -->
        <div class="tooltip-header">
          <div class="flex items-center gap-2">
            <span class="gender-icon" :class="tooltip.member?.gender === 'male' ? 'text-blue-500' : 'text-pink-500'">
              {{ tooltip.member?.gender === 'male' ? '♂' : '♀' }}
            </span>
            <h4 class="font-bold text-gray-800">{{ tooltip.member?.full_name }}</h4>
          </div>
          <div v-if="tooltip.member?.death_date" class="death-indicator" title="Meninggal">✝</div>
        </div>

        <!-- Birth/Death dates with age calculation -->
        <div v-if="tooltip.member?.birth_date || tooltip.member?.death_date" class="date-info">
          <div v-if="tooltip.member?.birth_date && tooltip.member?.death_date" class="text-sm text-gray-600">
            📅 {{ new Date(tooltip.member.birth_date).getFullYear() }} - {{ new Date(tooltip.member.death_date).getFullYear() }}
            ({{ new Date(tooltip.member.death_date).getFullYear() - new Date(tooltip.member.birth_date).getFullYear() }} tahun)
          </div>
          <div v-else-if="tooltip.member?.birth_date" class="text-sm text-gray-600">
            📅 Lahir {{ new Date(tooltip.member.birth_date).getFullYear() }}
            ({{ new Date().getFullYear() - new Date(tooltip.member.birth_date).getFullYear() }} tahun)
          </div>
        </div>

        <!-- Contact information -->
        <div v-if="tooltip.member?.address || tooltip.member?.phone" class="contact-info">
          <p v-if="tooltip.member?.address" class="text-sm text-gray-600">
            📍 {{ tooltip.member.address }}
          </p>
          <p v-if="tooltip.member?.phone" class="text-sm text-gray-600">
            📞 {{ tooltip.member.phone }}
          </p>
        </div>

        <!-- Marriage information -->
        <div v-if="tooltip.member?.marriage_info?.length" class="marriage-info">
          <p class="text-xs font-semibold text-pink-600 mb-2 flex items-center gap-1">
            💖 {{ tooltip.member.marriage_info.length }} Pernikahan
          </p>
          <div v-for="(marriage, index) in tooltip.member.marriage_info" :key="marriage.spouse_id" 
               class="marriage-item">
            <div class="flex items-start gap-2">
              <span class="marriage-number">{{ index + 1 }}.</span>
              <div class="marriage-details">
                <div class="spouse-name">{{ marriage.spouse_name }}</div>
                <div v-if="marriage.marriage_date" class="marriage-date">
                  💍 {{ new Date(marriage.marriage_date).getFullYear() }}
                </div>
                <div v-if="marriage.divorce_date" class="divorce-date">
                  💔 Bercerai {{ new Date(marriage.divorce_date).getFullYear() }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Children count -->
        <div v-if="tooltip.member?.children?.length" class="children-info">
          <p class="text-xs text-green-600 flex items-center gap-1">
            👶 {{ tooltip.member.children.length }} Anak
          </p>
        </div>

        <!-- Instructions -->
        <div class="tooltip-instructions">
          💡 Klik: pilih • Double-klik: edit • Klik kanan: tambah anak
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as d3 from 'd3'
import { fmtDate } from '@/utils/format.js'

const props = defineProps({
  treeData: {
    type: Array,
    default: () => []
  },
  selectedMemberId: {
    type: Number,
    default: null
  }
})

const emit = defineEmits(['member-select', 'member-edit', 'add-child'])

// Refs
const containerRef = ref(null)
const svgRef = ref(null)
const zoomLevel = ref(1)
const layoutMode = ref('vertical') // vertical or horizontal

// D3 instances
let svg, g, simulation, zoom
let nodes = [], links = [], marriageLinks = []

// Tooltip
const tooltip = ref({
  visible: false,
  x: 0,
  y: 0,
  member: null
})

// Tree configuration
const config = {
  width: 1200,
  height: 800,
  nodeRadius: 30,
  linkDistance: 120,
  marriageLinkDistance: 80,
  colors: {
    male: '#3B82F6',
    female: '#EC4899',
    deceased: '#6B7280',
    marriage: '#F59E0B',
    link: '#D1D5DB'
  }
}

onMounted(() => {
  initializeD3()
  if (props.treeData.length > 0) {
    updateTree()
  }
})

watch(() => props.treeData, () => {
  if (props.treeData.length > 0) {
    updateTree()
  }
}, { deep: true })

watch(() => props.selectedMemberId, (newId) => {
  highlightSelectedNode(newId)
})

function initializeD3() {
  // Setup SVG
  svg = d3.select(svgRef.value)
    .attr('viewBox', `0 0 ${config.width} ${config.height}`)
    .style('width', '100%')
    .style('height', '100%')

  // Setup zoom
  zoom = d3.zoom()
    .scaleExtent([0.1, 3])
    .on('zoom', (event) => {
      g.attr('transform', event.transform)
      zoomLevel.value = event.transform.k
    })

  svg.call(zoom)

  // Create main group
  g = svg.append('g')
    .attr('class', 'tree-group')

  // Add gradient definitions for nodes
  const defs = svg.append('defs')
  
  // Male gradient
  const maleGradient = defs.append('linearGradient')
    .attr('id', 'male-gradient')
    .attr('x1', '0%').attr('y1', '0%').attr('x2', '0%').attr('y2', '100%')
  maleGradient.append('stop').attr('offset', '0%').attr('stop-color', '#60A5FA')
  maleGradient.append('stop').attr('offset', '100%').attr('stop-color', '#3B82F6')
  
  // Female gradient  
  const femaleGradient = defs.append('linearGradient')
    .attr('id', 'female-gradient')
    .attr('x1', '0%').attr('y1', '0%').attr('x2', '0%').attr('y2', '100%')
  femaleGradient.append('stop').attr('offset', '0%').attr('stop-color', '#F472B6')
  femaleGradient.append('stop').attr('offset', '100%').attr('stop-color', '#EC4899')

  // Marriage arrow marker
  defs.append('marker')
    .attr('id', 'marriage-arrow')
    .attr('viewBox', '0 0 10 10')
    .attr('refX', 5)
    .attr('refY', 5)
    .attr('markerWidth', 6)
    .attr('markerHeight', 6)
    .attr('orient', 'auto')
    .append('path')
    .attr('d', 'M 0 0 L 10 5 L 0 10 z')
    .attr('fill', config.colors.marriage)

  // Heart pattern for marriage connections
  defs.append('pattern')
    .attr('id', 'heart-pattern')
    .attr('patternUnits', 'userSpaceOnUse')
    .attr('width', 20)
    .attr('height', 20)
    .append('text')
    .attr('x', 10)
    .attr('y', 15)
    .attr('text-anchor', 'middle')
    .attr('font-size', '12px')
    .attr('fill', config.colors.marriage)
    .text('♥')

  // Setup force simulation with improved forces
  simulation = d3.forceSimulation()
    .force('link', d3.forceLink()
      .id(d => d.id)
      .distance(d => d.type === 'marriage' ? config.marriageLinkDistance : config.linkDistance)
      .strength(d => d.type === 'marriage' ? 0.3 : 0.8)
    )
    .force('charge', d3.forceManyBody().strength(-1200))
    .force('center', d3.forceCenter(config.width / 2, config.height / 2))
    .force('collision', d3.forceCollide().radius(config.nodeRadius + 15))
    .force('x', d3.forceX(config.width / 2).strength(0.1))
    .force('y', d3.forceY().y(d => (d.level || 0) * 150 + 100).strength(0.3))
    .on('tick', ticked)
}

function processTreeData() {
  nodes = []
  links = []
  marriageLinks = []

  // Recursive function to extract all nodes
  function extractNodes(nodeList, level = 0) {
    for (const node of nodeList) {
      nodes.push({
        ...node,
        level,
        x: config.width / 2 + Math.random() * 100 - 50,
        y: (level + 1) * 150 + Math.random() * 50 - 25
      })

      // Create parent-child links
      if (node.children && node.children.length > 0) {
        for (const child of node.children) {
          links.push({
            source: node.id,
            target: child.id,
            type: 'family'
          })
        }
        extractNodes(node.children, level + 1)
      }

      // Create marriage links
      if (node.marriage_info && node.marriage_info.length > 0) {
        for (const marriage of node.marriage_info) {
          const existingLink = marriageLinks.find(link => 
            (link.source === node.id && link.target === marriage.spouse_id) ||
            (link.source === marriage.spouse_id && link.target === node.id)
          )
          
          if (!existingLink) {
            marriageLinks.push({
              source: node.id,
              target: marriage.spouse_id,
              type: 'marriage',
              marriage_date: marriage.marriage_date
            })
          }
        }
      }
    }
  }

  extractNodes(props.treeData)
  return { nodes, links: [...links, ...marriageLinks] }
}

function updateTree() {
  const { nodes: newNodes, links: newLinks } = processTreeData()
  
  // Update simulation
  simulation.nodes(newNodes)
  simulation.force('link').links(newLinks)

  updateVisualization(newNodes, newLinks)
  simulation.alpha(1).restart()
}

function updateVisualization(nodes, links) {
  // Separate family and marriage links
  const familyLinks = links.filter(l => l.type === 'family')
  const marriageLinks = links.filter(l => l.type === 'marriage')

  // Update family links (parent-child)
  const familyLink = g.selectAll('.family-link')
    .data(familyLinks, d => `family-${d.source.id || d.source}-${d.target.id || d.target}`)
    
  familyLink.exit().remove()
  
  const familyLinkEnter = familyLink.enter().append('line')
    .attr('class', 'family-link link')
    .attr('stroke-width', 2)
    .attr('stroke', config.colors.link)

  familyLink.merge(familyLinkEnter)

  // Update marriage links with special styling
  const marriageLink = g.selectAll('.marriage-link')
    .data(marriageLinks, d => `marriage-${d.source.id || d.source}-${d.target.id || d.target}`)
    
  marriageLink.exit().remove()
  
  const marriageLinkEnter = marriageLink.enter().append('g')
    .attr('class', 'marriage-link-group')

  // Add marriage line
  marriageLinkEnter.append('line')
    .attr('class', 'marriage-link link')
    .attr('stroke-width', 4)
    .attr('stroke', config.colors.marriage)
    .attr('stroke-dasharray', '8,4')
    .attr('opacity', 0.7)

  // Add heart icon in the middle of marriage line
  marriageLinkEnter.append('text')
    .attr('class', 'marriage-heart')
    .attr('text-anchor', 'middle')
    .attr('dy', '0.3em')
    .attr('font-size', '16px')
    .attr('fill', config.colors.marriage)
    .attr('stroke', 'white')
    .attr('stroke-width', 2)
    .attr('paint-order', 'stroke')
    .text('♥')

  marriageLink.merge(marriageLinkEnter)

  // Update nodes
  const node = g.selectAll('.node')
    .data(nodes, d => d.id)

  node.exit().remove()

  const nodeEnter = node.enter().append('g')
    .attr('class', 'node')
    .call(d3.drag()
      .on('start', dragstarted)
      .on('drag', dragged)
      .on('end', dragended))

  // Add node background circle for better visibility
  nodeEnter.append('circle')
    .attr('class', 'node-bg')
    .attr('r', config.nodeRadius + 2)
    .attr('fill', 'white')
    .attr('stroke', '#E5E7EB')
    .attr('stroke-width', 1)
    .attr('opacity', 0.9)

  // Add node circles with enhanced styling
  nodeEnter.append('circle')
    .attr('class', 'node-main')
    .attr('r', config.nodeRadius)
    .attr('fill', d => {
      if (d.death_date) return config.colors.deceased
      return d.gender === 'male' ? 'url(#male-gradient)' : 'url(#female-gradient)'
    })
    .attr('stroke', '#fff')
    .attr('stroke-width', 3)
    .style('cursor', 'pointer')
    .style('filter', 'drop-shadow(0 4px 8px rgba(0, 0, 0, 0.15))')

  // Add decorative ring for special members
  nodeEnter.append('circle')
    .attr('class', 'node-ring')
    .attr('r', config.nodeRadius + 6)
    .attr('fill', 'none')
    .attr('stroke', d => d.marriage_info?.length > 0 ? config.colors.marriage : 'transparent')
    .attr('stroke-width', 2)
    .attr('stroke-dasharray', '4,4')
    .attr('opacity', 0.6)

  // Add node labels (initials)
  nodeEnter.append('text')
    .attr('class', 'node-initials')
    .attr('text-anchor', 'middle')
    .attr('dy', '.3em')
    .attr('fill', 'white')
    .attr('font-size', '12px')
    .attr('font-weight', 'bold')
    .attr('pointer-events', 'none')
    .text(d => getInitials(d.full_name))

  // Add name labels below nodes
  nodeEnter.append('text')
    .attr('class', 'node-name')
    .attr('text-anchor', 'middle')
    .attr('dy', config.nodeRadius + 18)
    .attr('fill', '#374151')
    .attr('font-size', '13px')
    .attr('font-weight', 'bold')
    .attr('pointer-events', 'none')
    .text(d => {
      const firstName = d.full_name.split(' ')[0]
      return firstName.length > 10 ? firstName.substring(0, 8) + '...' : firstName
    })

  // Add birth year with enhanced styling
  nodeEnter.append('text')
    .attr('class', 'node-year')
    .attr('text-anchor', 'middle')
    .attr('dy', config.nodeRadius + 33)
    .attr('fill', '#6B7280')
    .attr('font-size', '11px')
    .attr('pointer-events', 'none')
    .text(d => {
      if (d.birth_date && d.death_date) {
        const birth = new Date(d.birth_date).getFullYear()
        const death = new Date(d.death_date).getFullYear()
        return `${birth}-${death}`
      } else if (d.birth_date) {
        return new Date(d.birth_date).getFullYear()
      }
      return ''
    })

  // Add death indicator
  nodeEnter.append('text')
    .attr('class', 'death-indicator')
    .attr('text-anchor', 'middle')
    .attr('dy', -config.nodeRadius - 5)
    .attr('fill', '#6B7280')
    .attr('font-size', '16px')
    .attr('pointer-events', 'none')
    .style('opacity', d => d.death_date ? 1 : 0)
    .text('✝')

  // Add marriage count indicator
  nodeEnter.append('text')
    .attr('class', 'marriage-count')
    .attr('text-anchor', 'middle')
    .attr('x', config.nodeRadius - 5)
    .attr('y', -config.nodeRadius + 5)
    .attr('fill', config.colors.marriage)
    .attr('font-size', '12px')
    .attr('font-weight', 'bold')
    .attr('pointer-events', 'none')
    .style('opacity', d => d.marriage_info?.length > 0 ? 1 : 0)
    .text(d => d.marriage_info?.length > 0 ? '♥' : '')

  // Add event listeners
  const allNodes = node.merge(nodeEnter)
  
  allNodes
    .on('click', (event, d) => {
      event.stopPropagation()
      emit('member-select', d)
      highlightSelectedNode(d.id)
    })
    .on('dblclick', (event, d) => {
      event.stopPropagation()
      emit('member-edit', d)
    })
    .on('mouseover', (event, d) => {
      showTooltip(event, d)
      highlightConnections(d.id)
    })
    .on('mouseout', (event, d) => {
      hideTooltip()
      clearHighlights()
    })
    .on('contextmenu', (event, d) => {
      event.preventDefault()
      emit('add-child', d)
    })
}

function highlightConnections(nodeId) {
  // Highlight connected nodes
  g.selectAll('.node')
    .style('opacity', d => {
      // Check if this node is connected to the hovered node
      const isConnected = links.some(link => 
        (link.source.id === nodeId && link.target.id === d.id) ||
        (link.target.id === nodeId && link.source.id === d.id) ||
        d.id === nodeId
      )
      return isConnected ? 1 : 0.3
    })

  // Highlight connected links
  g.selectAll('.link')
    .style('opacity', d => 
      (d.source.id === nodeId || d.target.id === nodeId) ? 1 : 0.1
    )
    .attr('stroke-width', d => {
      const isConnected = d.source.id === nodeId || d.target.id === nodeId
      return isConnected ? (d.type === 'marriage' ? 6 : 4) : (d.type === 'marriage' ? 4 : 2)
    })
}

function clearHighlights() {
  g.selectAll('.node').style('opacity', 1)
  g.selectAll('.link')
    .style('opacity', d => d.type === 'marriage' ? 0.7 : 0.6)
    .attr('stroke-width', d => d.type === 'marriage' ? 4 : 2)
}

function ticked() {
  // Update family links (simple lines)
  g.selectAll('.family-link')
    .attr('x1', d => d.source.x)
    .attr('y1', d => d.source.y)
    .attr('x2', d => d.target.x)
    .attr('y2', d => d.target.y)

  // Update marriage link groups
  g.selectAll('.marriage-link-group').each(function(d) {
    const group = d3.select(this)
    
    // Update marriage line
    group.select('.marriage-link')
      .attr('x1', d.source.x)
      .attr('y1', d.source.y)
      .attr('x2', d.target.x)
      .attr('y2', d.target.y)

    // Position heart in the middle of marriage line
    group.select('.marriage-heart')
      .attr('x', (d.source.x + d.target.x) / 2)
      .attr('y', (d.source.y + d.target.y) / 2)
  })

  // Update node positions
  g.selectAll('.node')
    .attr('transform', d => `translate(${d.x},${d.y})`)
}

function dragstarted(event, d) {
  if (!event.active) simulation.alphaTarget(0.3).restart()
  d.fx = d.x
  d.fy = d.y
}

function dragged(event, d) {
  d.fx = event.x
  d.fy = event.y
}

function dragended(event, d) {
  if (!event.active) simulation.alphaTarget(0)
  d.fx = null
  d.fy = null
}

function highlightSelectedNode(nodeId) {
  g.selectAll('.node circle')
    .attr('stroke', d => d.id === nodeId ? '#F59E0B' : '#fff')
    .attr('stroke-width', d => d.id === nodeId ? 4 : 3)
}

function showTooltip(event, member) {
  const [x, y] = d3.pointer(event, containerRef.value)
  tooltip.value = {
    visible: true,
    x: x + 10,
    y: y - 10,
    member
  }
}

function hideTooltip() {
  tooltip.value.visible = false
}

function getInitials(name) {
  return name.split(' ')
    .map(word => word.charAt(0))
    .slice(0, 2)
    .join('')
    .toUpperCase()
}

function centerTree() {
  const transform = d3.zoomIdentity.translate(config.width / 2, config.height / 2).scale(1)
  svg.transition()
    .duration(750)
    .call(zoom.transform, transform)
}

function zoomIn() {
  svg.transition()
    .duration(300)
    .call(zoom.scaleBy, 1.5)
}

function zoomOut() {
  svg.transition()
    .duration(300)
    .call(zoom.scaleBy, 1 / 1.5)
}

function toggleLayout() {
  layoutMode.value = layoutMode.value === 'vertical' ? 'horizontal' : 'vertical'
  updateTreeLayout()
}

function updateTreeLayout() {
  if (layoutMode.value === 'horizontal') {
    simulation.force('center', d3.forceCenter(config.height / 2, config.width / 2))
  } else {
    simulation.force('center', d3.forceCenter(config.width / 2, config.height / 2))
  }
  simulation.alpha(1).restart()
}

onUnmounted(() => {
  if (simulation) {
    simulation.stop()
  }
})
</script>

<style scoped>
.family-tree-container {
  @apply relative w-full h-full bg-white rounded-xl overflow-hidden;
}

.family-tree-controls {
  @apply absolute top-4 right-4 z-10 flex gap-2;
}

.control-btn {
  @apply w-10 h-10 bg-white/90 backdrop-blur-sm border border-warm-gray-200 rounded-lg
         flex items-center justify-center text-warm-gray-700 hover:bg-white hover:text-terra
         transition-all duration-200 shadow-sm hover:shadow-md;
}

.zoom-level {
  @apply px-3 py-2 bg-white/90 backdrop-blur-sm border border-warm-gray-200 rounded-lg
         text-xs font-semibold text-warm-gray-700 min-w-[50px] text-center;
}

.family-tree-svg {
  @apply w-full h-full;
}

.family-tree-tooltip {
  @apply absolute pointer-events-none z-20;
  animation: fadeInUp 0.2s ease-out;
}

.tooltip-content {
  @apply bg-gray-900/95 backdrop-blur-xl text-white p-4 rounded-xl shadow-2xl max-w-xs border border-white/10;
}

.tooltip-header {
  @apply flex items-start justify-between mb-3 pb-2 border-b border-white/20;
}

.gender-icon {
  @apply text-lg font-bold;
}

.death-indicator {
  @apply text-gray-400 text-sm;
}

.date-info {
  @apply mb-2;
}

.contact-info {
  @apply mb-3 space-y-1;
}

.marriage-info {
  @apply mb-3 p-2 bg-black/20 rounded-lg border border-pink-500/30;
}

.marriage-item {
  @apply mb-2 last:mb-0;
}

.marriage-number {
  @apply text-pink-400 font-semibold text-sm;
}

.marriage-details {
  @apply flex-1;
}

.spouse-name {
  @apply font-semibold text-pink-300 text-sm;
}

.marriage-date {
  @apply text-xs text-yellow-300 mt-1;
}

.divorce-date {
  @apply text-xs text-red-300 mt-1;
}

.children-info {
  @apply mb-2 p-2 bg-black/20 rounded-lg border border-green-500/30;
}

.tooltip-instructions {
  @apply text-xs text-gray-400 pt-2 mt-2 border-t border-white/20;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

:deep(.node circle) {
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.15));
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.node:hover circle) {
  transform: scale(1.15);
  filter: drop-shadow(0 8px 25px rgba(0, 0, 0, 0.25));
}

:deep(.node text) {
  transition: all 0.3s ease;
}

:deep(.node:hover .node-name) {
  font-size: 14px;
  font-weight: 800;
}

:deep(.node:hover .node-year) {
  font-size: 12px;
  font-weight: 600;
}

:deep(.link) {
  opacity: 0.6;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.family-link:hover) {
  opacity: 1;
  stroke-width: 4px !important;
  stroke: #3B82F6;
}

:deep(.marriage-link:hover) {
  opacity: 1;
  stroke-width: 6px !important;
  stroke-dasharray: 12,6;
}

:deep(.marriage-heart) {
  transition: all 0.3s ease;
}

:deep(.marriage-link-group:hover .marriage-heart) {
  font-size: 20px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
}

:deep(.node-ring) {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 0.6;
    transform: scale(1);
  }
  50% {
    opacity: 0.9;
    transform: scale(1.05);
  }
}
</style>