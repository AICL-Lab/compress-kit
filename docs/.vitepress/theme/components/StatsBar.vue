<script setup lang="ts">
import { ref, onMounted } from 'vue'
import statsData from '../../data/stats.json'

interface Stats {
  stars: number
  algorithms: number
  languages: number
  implementations: number
}

const stats = ref<Stats>(statsData)
const animatedValues = ref<Record<string, number>>({
  stars: 0,
  algorithms: 0,
  languages: 0,
  implementations: 0
})

const formatNumber = (num: number): string => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

const animateValue = (key: keyof Stats, target: number, duration: number = 1000) => {
  const start = performance.now()
  const startValue = 0
  
  const step = (timestamp: number) => {
    const elapsed = timestamp - start
    const progress = Math.min(elapsed / duration, 1)
    
    // Easing function: easeOutQuart
    const easeOutQuart = 1 - Math.pow(1 - progress, 4)
    const current = Math.floor(startValue + (target - startValue) * easeOutQuart)
    
    animatedValues.value[key] = current
    
    if (progress < 1) {
      requestAnimationFrame(step)
    }
  }
  
  requestAnimationFrame(step)
}

onMounted(() => {
  // Stagger animations
  Object.entries(stats.value).forEach(([key, value], index) => {
    setTimeout(() => {
      animateValue(key as keyof Stats, value)
    }, index * 150)
  })
})
</script>

<template>
  <div class="ck-stats-bar">
    <div class="ck-stat ck-animate-fade-in-up" style="animation-delay: 100ms">
      <div class="ck-stat-value">
        {{ formatNumber(animatedValues.stars) }}+
      </div>
      <div class="ck-stat-label">GitHub Stars</div>
    </div>
    <div class="ck-stat ck-animate-fade-in-up" style="animation-delay: 200ms">
      <div class="ck-stat-value">{{ animatedValues.algorithms }}</div>
      <div class="ck-stat-label">Algorithms</div>
    </div>
    <div class="ck-stat ck-animate-fade-in-up" style="animation-delay: 300ms">
      <div class="ck-stat-value">{{ animatedValues.languages }}</div>
      <div class="ck-stat-label">Languages</div>
    </div>
    <div class="ck-stat ck-animate-fade-in-up" style="animation-delay: 400ms">
      <div class="ck-stat-value">{{ animatedValues.implementations }}</div>
      <div class="ck-stat-label">Implementations</div>
    </div>
  </div>
</template>

<style scoped>
.ck-stats-bar {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: clamp(2rem, 5vw, 4rem);
  padding: 2rem;
  margin: 1rem 0 3rem;
}

.ck-stat {
  text-align: center;
  opacity: 0;
  animation: fadeInUp 0.6s ease-out forwards;
}

.ck-stat-value {
  font-size: clamp(2.5rem, 5vw, 3.5rem);
  font-weight: 700;
  background: linear-gradient(135deg, #2563eb 0%, #0ea5e9 50%, #10b981 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
  letter-spacing: -0.02em;
}

.ck-stat-label {
  font-size: 0.875rem;
  color: var(--vp-c-text-3);
  margin-top: 0.5rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
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
  .ck-stats-bar {
    gap: 1.5rem;
  }
  
  .ck-stat-value {
    font-size: 2rem;
  }
  
  .ck-stat-label {
    font-size: 0.75rem;
  }
}
</style>
