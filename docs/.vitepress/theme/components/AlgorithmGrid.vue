<script setup lang="ts">
import { useData, withBase } from 'vitepress'
import { computed } from 'vue'

const { localeIndex } = useData()
const lang = computed(() => localeIndex.value === 'root' ? 'en' : localeIndex.value)
const algoLink = (slug: string) => withBase(`/${lang.value}/algorithms/${slug}`)

interface Algorithm {
  id: string
  name: string
  slug: string
  icon: string
  description: string
  compression: 'Medium' | 'High' | 'Variable' | 'Low'
  speed: 'Fast' | 'Medium' | 'Slow' | 'Very Fast'
  bestFor: string[]
  compressionLevel: 'medium' | 'high' | 'variable'
  speedLevel: 'fast' | 'medium' | 'slow' | 'very-fast'
}

const algorithms: Algorithm[] = [
  {
    id: 'huffman',
    name: 'Huffman Coding',
    slug: 'huffman',
    icon: '🌳',
    description: 'Optimal prefix codes based on symbol frequency. The classic approach to lossless compression.',
    compression: 'Medium',
    speed: 'Fast',
    bestFor: ['Text files', 'General data', 'Natural language'],
    compressionLevel: 'medium',
    speedLevel: 'fast'
  },
  {
    id: 'arithmetic',
    name: 'Arithmetic Coding',
    slug: 'arithmetic',
    icon: '🧮',
    description: 'Entire message encoded as a single number. Achieves entropy limit for maximum compression.',
    compression: 'High',
    speed: 'Medium',
    bestFor: ['Maximum compression', 'Statistical data', 'Archival storage'],
    compressionLevel: 'high',
    speedLevel: 'medium'
  },
  {
    id: 'range',
    name: 'Range Coder',
    slug: 'range',
    icon: '🎯',
    description: 'Integer-based arithmetic coding. Production-ready balance of speed and compression.',
    compression: 'High',
    speed: 'Fast',
    bestFor: ['Production systems', 'Real-time compression', 'Balanced workloads'],
    compressionLevel: 'high',
    speedLevel: 'fast'
  },
  {
    id: 'rle',
    name: 'Run-Length Encoding',
    slug: 'rle',
    icon: '📏',
    description: 'Simple and fast compression for repetitive data. Often used as preprocessing.',
    compression: 'Variable',
    speed: 'Very Fast',
    bestFor: ['Bitmap images', 'Log files', 'Preprocessing step'],
    compressionLevel: 'variable',
    speedLevel: 'very-fast'
  }
]

const getCompressionBadgeClass = (level: string) => {
  switch (level) {
    case 'high': return 'ck-badge-compression-high'
    case 'medium': return 'ck-badge-compression-medium'
    default: return 'ck-badge-compression-variable'
  }
}

const getSpeedBadgeClass = (level: string) => {
  switch (level) {
    case 'very-fast': return 'ck-badge-speed-very-fast'
    case 'fast': return 'ck-badge-speed-fast'
    case 'medium': return 'ck-badge-speed-medium'
    default: return 'ck-badge-speed-slow'
  }
}
</script>

<template>
  <div class="ck-algorithm-grid">
    <a 
      v-for="(algo, index) in algorithms" 
      :key="algo.id"
      :href="algoLink(algo.slug)"
      class="ck-algorithm-card ck-hover-lift"
      :style="{ animationDelay: `${index * 100}ms` }"
    >
      <div class="card-header">
        <div class="algorithm-icon">{{ algo.icon }}</div>
        <div class="algorithm-title">
          <h3>{{ algo.name }}</h3>
          <div class="badges">
            <span :class="['badge', getCompressionBadgeClass(algo.compressionLevel)]">
              {{ algo.compression }} Compression
            </span>
            <span :class="['badge', getSpeedBadgeClass(algo.speedLevel)]">
              {{ algo.speed }} Speed
            </span>
          </div>
        </div>
      </div>
      
      <p class="algorithm-description">{{ algo.description }}</p>
      
      <div class="best-for">
        <span class="label">Best for:</span>
        <div class="tags">
          <span v-for="use in algo.bestFor" :key="use" class="tag">{{ use }}</span>
        </div>
      </div>
      
      <div class="card-footer">
        <span class="learn-more">Learn more →</span>
      </div>
    </a>
  </div>
</template>

<style scoped>
.ck-algorithm-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin: 2rem 0;
}

.ck-algorithm-card {
  display: flex;
  flex-direction: column;
  background: var(--vp-c-bg-elv);
  border: 1px solid var(--vp-c-divider);
  border-radius: 1rem;
  padding: 1.5rem;
  text-decoration: none;
  color: inherit;
  transition: all 0.2s ease;
  opacity: 0;
  animation: fadeInUp 0.5s ease-out forwards;
  position: relative;
  overflow: hidden;
}

.ck-algorithm-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #2563eb, #0ea5e9, #10b981);
  opacity: 0;
  transition: opacity 0.2s ease;
}

.ck-algorithm-card:hover::before {
  opacity: 1;
}

.ck-algorithm-card:hover {
  border-color: var(--vp-c-brand-3);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 4px 10px -4px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1rem;
}

.algorithm-icon {
  font-size: 2.5rem;
  line-height: 1;
  flex-shrink: 0;
}

.algorithm-title h3 {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0 0 0.5rem;
  color: var(--vp-c-text-1);
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
}

.badge {
  font-size: 0.65rem;
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  text-transform: uppercase;
  letter-spacing: 0.025em;
  white-space: nowrap;
}

.ck-badge-compression-high {
  background: rgba(16, 185, 129, 0.1);
  color: #059669;
}

.dark .ck-badge-compression-high {
  background: rgba(52, 211, 153, 0.15);
  color: #34d399;
}

.ck-badge-compression-medium {
  background: rgba(37, 99, 235, 0.1);
  color: #2563eb;
}

.dark .ck-badge-compression-medium {
  background: rgba(96, 165, 250, 0.15);
  color: #60a5fa;
}

.ck-badge-compression-variable {
  background: rgba(139, 92, 246, 0.1);
  color: #8b5cf6;
}

.dark .ck-badge-compression-variable {
  background: rgba(167, 139, 250, 0.15);
  color: #a78bfa;
}

.ck-badge-speed-very-fast,
.ck-badge-speed-fast {
  background: rgba(14, 165, 233, 0.1);
  color: #0ea5e9;
}

.dark .ck-badge-speed-very-fast,
.dark .ck-badge-speed-fast {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

.ck-badge-speed-medium {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.dark .ck-badge-speed-medium {
  background: rgba(250, 204, 21, 0.15);
  color: #facc15;
}

.algorithm-description {
  font-size: 0.9375rem;
  line-height: 1.6;
  color: var(--vp-c-text-2);
  margin: 0 0 1.25rem;
  flex-grow: 1;
}

.best-for {
  margin-bottom: 1rem;
}

.best-for .label {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--vp-c-text-3);
  margin-bottom: 0.5rem;
  display: block;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.375rem;
}

.tag {
  font-size: 0.8125rem;
  padding: 0.25rem 0.625rem;
  background: var(--vp-c-bg-alt);
  border: 1px solid var(--vp-c-divider);
  border-radius: 0.375rem;
  color: var(--vp-c-text-2);
}

.card-footer {
  margin-top: auto;
  padding-top: 1rem;
  border-top: 1px solid var(--vp-c-divider);
}

.learn-more {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--vp-c-brand-1);
  transition: color 0.2s ease;
}

.ck-algorithm-card:hover .learn-more {
  color: var(--vp-c-brand-2);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 640px) {
  .ck-algorithm-grid {
    grid-template-columns: 1fr;
  }
}
</style>
