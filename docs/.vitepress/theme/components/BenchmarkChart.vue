<script setup lang="ts">
import { ref, computed, watchEffect } from 'vue'
import { useData } from 'vitepress'
import benchmarkData from '../../data/benchmarks.json'
import { getBenchmarkContent } from '../../data/site-content.mjs'
import { getBarHeight as getMetricBarHeight, isBestCompressionRatio } from './benchmarkChartMetrics.js'

interface BenchmarkResult {
  algorithm: string
  language: string
  dataset: string
  encodeTime: number
  decodeTime: number
  encodeSpeed: number
  decodeSpeed: number
  compressionRatio: number
  throughput: string
}

const results = ref<BenchmarkResult[]>(benchmarkData.results)
const { localeIndex } = useData()
const locale = computed(() => (localeIndex.value === 'root' ? 'en' : localeIndex.value))
const content = computed(() => getBenchmarkContent(locale.value))

const algorithms = computed(() => content.value.algorithms.map(entry => entry.id))
const languages = computed(() => content.value.languages.map(entry => entry.id))
const datasetOrder = ['textlike_10MiB', 'repetitive_10MiB', 'small_dictionary_like']
const languageColors = computed(() =>
  Object.fromEntries(content.value.languages.map(entry => [entry.id, entry.color]))
)
const languageNames = computed(() =>
  Object.fromEntries(content.value.languages.map(entry => [entry.id, entry.label]))
)
const algorithmNames = computed(() =>
  Object.fromEntries(content.value.algorithms.map(entry => [entry.id, entry.label]))
)
const datasetNames = computed(() =>
  Object.fromEntries(content.value.datasets.map(entry => [entry.id, entry.label]))
)
const metricNames = computed(() =>
  Object.fromEntries(content.value.metrics.map(entry => [entry.id, entry.label]))
)

const selectedMetric = ref<'encodeSpeed' | 'decodeSpeed' | 'compressionRatio'>('encodeSpeed')
const selectedDataset = ref('')

const datasetOptions = computed(() => {
  return [...new Set(results.value.map(result => result.dataset))].sort((a, b) => {
    const left = datasetOrder.indexOf(a)
    const right = datasetOrder.indexOf(b)
    if (left === -1 && right === -1) return a.localeCompare(b)
    if (left === -1) return 1
    if (right === -1) return -1
    return left - right
  })
})

watchEffect(() => {
  if (!datasetOptions.value.includes(selectedDataset.value)) {
    selectedDataset.value = datasetOptions.value[0] ?? ''
  }
})

const filteredResults = computed(() => {
  return results.value.filter(
    r => r.dataset === selectedDataset.value
  )
})

const groupedByAlgorithm = computed(() => {
  const grouped: Record<string, BenchmarkResult[]> = {}
  algorithms.value.forEach(algo => {
    grouped[algo] = filteredResults.value.filter(r => r.algorithm === algo)
  })
  return grouped
})

const maxValue = computed(() => {
  if (filteredResults.value.length === 0) return 100
  return Math.max(...filteredResults.value.map(r => r[selectedMetric.value]))
})

const compressionRatios = computed(() => filteredResults.value.map(r => r.compressionRatio))

const formatMetricName = (metric: string): string => {
  return metricNames.value[metric] || metric
}

const getBarHeight = (value: number): number => {
  return getMetricBarHeight(selectedMetric.value, value, maxValue.value)
}

const isCompressionRatioLeader = (value: number): boolean => {
  return isBestCompressionRatio(selectedMetric.value, value, compressionRatios.value)
}

const formatMetricValue = (value: number): string => {
  if (selectedMetric.value === 'compressionRatio') {
    return value < 0.01 ? value.toFixed(3) : value.toFixed(2)
  }
  return value.toFixed(value < 10 ? 1 : 0)
}
</script>

<template>
  <div class="benchmark-chart-container">
    <div class="chart-header">
      <h3 class="chart-title">{{ content.title }}</h3>
      <div class="chart-controls">
        <div class="control-group">
          <label>{{ content.datasetLabel }}</label>
          <select v-model="selectedDataset">
            <option v-for="dataset in datasetOptions" :key="dataset" :value="dataset">
              {{ datasetNames[dataset] || dataset }}
            </option>
          </select>
        </div>
        <div class="control-group">
          <label>{{ content.metricLabel }}</label>
          <select v-model="selectedMetric">
            <option
              v-for="metric in content.metricOptions"
              :key="metric.id"
              :value="metric.id"
            >
              {{ metric.label }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <div class="chart-legend">
      <div v-for="lang in languages" :key="lang" class="legend-item">
        <span class="legend-color" :style="{ backgroundColor: languageColors[lang] }"></span>
        <span class="legend-name">{{ languageNames[lang] }}</span>
      </div>
    </div>

    <div class="chart-content">
      <div v-for="algo in algorithms" :key="algo" class="algorithm-group">
        <div class="algorithm-label">{{ algorithmNames[algo] }}</div>
        <div class="bars-container">
          <div 
            v-for="result in groupedByAlgorithm[algo]" 
            :key="result.language"
            class="bar-wrapper"
              :title="`${languageNames[result.language]}: ${formatMetricValue(result[selectedMetric])}`"
            >
              <div
                class="bar"
                :class="{ 'bar-best': isCompressionRatioLeader(result.compressionRatio) }"
                :style="{
                  height: `${getBarHeight(result[selectedMetric])}%`,
                  backgroundColor: languageColors[result.language]
                }"
            >
              <span class="bar-value">{{ formatMetricValue(result[selectedMetric]) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="selectedMetric === 'compressionRatio'" class="metric-note">
      {{ content.compressionNote }}
    </div>
    <div class="metric-label">{{ formatMetricName(selectedMetric) }}</div>
  </div>
</template>

<style scoped>
.benchmark-chart-container {
  background: var(--vp-c-bg-elv);
  border: 1px solid var(--vp-c-divider);
  border-radius: 1rem;
  padding: 1.5rem;
  margin: 2rem 0;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.chart-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
  color: var(--vp-c-text-1);
}

.chart-controls {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.control-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.control-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--vp-c-text-2);
}

.control-group select {
  padding: 0.375rem 0.75rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 0.5rem;
  background: var(--vp-c-bg-alt);
  color: var(--vp-c-text-1);
  font-size: 0.875rem;
  cursor: pointer;
  outline: none;
}

.control-group select:focus {
  border-color: var(--vp-c-brand-1);
}

.chart-legend {
  display: flex;
  justify-content: center;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--vp-c-divider);
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}

.legend-name {
  font-size: 0.875rem;
  color: var(--vp-c-text-2);
}

.chart-content {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  min-height: 250px;
}

.algorithm-group {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.algorithm-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--vp-c-text-1);
  margin-bottom: 0.75rem;
  text-align: center;
}

.bars-container {
  display: flex;
  gap: 0.5rem;
  align-items: flex-end;
  height: 200px;
  flex-grow: 1;
  width: 100%;
  justify-content: center;
}

.bar-wrapper {
  flex: 1;
  max-width: 30px;
  height: 100%;
  display: flex;
  align-items: flex-end;
}

.bar {
  width: 100%;
  border-radius: 4px 4px 0 0;
  transition: height 0.3s ease;
  position: relative;
  opacity: 0.85;
}

.bar:hover {
  opacity: 1;
}

.bar-best {
  box-shadow: 0 0 0 2px var(--vp-c-brand-1);
  opacity: 1;
}

.bar-value {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--vp-c-text-2);
  white-space: nowrap;
}

.metric-note {
  text-align: center;
  font-size: 0.8125rem;
  color: var(--vp-c-text-2);
  margin-top: 1.25rem;
}

.metric-label {
  text-align: center;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--vp-c-text-3);
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--vp-c-divider);
}

@media (max-width: 768px) {
  .chart-header {
    flex-direction: column;
    align-items: stretch;
  }

  .chart-content {
    grid-template-columns: repeat(2, 1fr);
    gap: 1.5rem 0.5rem;
  }

  .bars-container {
    height: 150px;
  }

  .chart-legend {
    gap: 1rem;
  }
}
</style>
