<script setup lang="ts">
import { ref, computed } from 'vue'

interface CodeExample {
  language: string
  label: string
  code: string
  filename: string
}

const props = defineProps<{
  examples: CodeExample[]
}>()

const activeIndex = ref(0)
const copied = ref(false)

const activeExample = computed(() => props.examples[activeIndex.value])

const languageClass = (lang: string): string => {
  const classes: Record<string, string> = {
    cpp: 'ck-lang-cpp',
    go: 'ck-lang-go',
    rust: 'ck-lang-rust'
  }
  return classes[lang] || ''
}

const languageIcon = (lang: string): string => {
  const icons: Record<string, string> = {
    cpp: '⚡',
    go: '🔵',
    rust: '⚙️'
  }
  return icons[lang] || '💻'
}

const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(activeExample.value.code)
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}
</script>

<template>
  <div class="code-comparison">
    <div class="comparison-header">
      <div class="language-tabs">
        <button
          v-for="(example, index) in examples"
          :key="example.language"
          class="tab-button"
          :class="{ 
            active: activeIndex === index,
            [languageClass(example.language)]: true 
          }"
          @click="activeIndex = index"
        >
          <span class="tab-icon">{{ languageIcon(example.language) }}</span>
          <span class="tab-label">{{ example.label }}</span>
        </button>
      </div>
      <button 
        class="copy-button"
        :class="{ copied }"
        @click="copyCode"
        :title="copied ? 'Copied!' : 'Copy code'"
      >
        <span class="copy-icon">📋</span>
        <span class="copy-text">{{ copied ? 'Copied!' : 'Copy' }}</span>
      </button>
    </div>
    
    <div class="code-content">
      <div class="code-header">
        <span class="filename">{{ activeExample.filename }}</span>
      </div>
      <div class="code-wrapper">
        <pre :class="`language-${activeExample.language}`"><code>{{ activeExample.code }}</code></pre>
      </div>
    </div>
  </div>
</template>

<style scoped>
.code-comparison {
  border: 1px solid var(--vp-c-divider);
  border-radius: 1rem;
  overflow: hidden;
  margin: 1.5rem 0;
  background: var(--vp-c-bg-elv);
}

.comparison-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: var(--vp-c-bg-alt);
  border-bottom: 1px solid var(--vp-c-divider);
  flex-wrap: wrap;
  gap: 0.5rem;
}

.language-tabs {
  display: flex;
  gap: 0.25rem;
  flex-wrap: wrap;
}

.tab-button {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  border: 1px solid transparent;
  border-radius: 0.5rem;
  background: transparent;
  color: var(--vp-c-text-2);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-button:hover {
  background: var(--vp-c-bg-elv);
  color: var(--vp-c-text-1);
}

.tab-button.active {
  background: var(--vp-c-bg-elv);
  border-color: var(--vp-c-divider);
  color: var(--vp-c-text-1);
}

.tab-button.active.ck-lang-cpp {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.1);
}

.tab-button.active.ck-lang-go {
  border-color: #00add8;
  background: rgba(0, 173, 216, 0.1);
}

.tab-button.active.ck-lang-rust {
  border-color: #de6e4b;
  background: rgba(222, 110, 75, 0.1);
}

.tab-icon {
  font-size: 0.875rem;
}

.copy-button {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 0.5rem;
  background: var(--vp-c-bg-elv);
  color: var(--vp-c-text-2);
  font-size: 0.8125rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-button:hover {
  background: var(--vp-c-bg-alt);
  color: var(--vp-c-text-1);
}

.copy-button.copied {
  background: rgba(16, 185, 129, 0.1);
  border-color: #10b981;
  color: #059669;
}

.dark .copy-button.copied {
  background: rgba(52, 211, 153, 0.15);
  color: #34d399;
}

.code-content {
  background: #1e1e1e;
}

.code-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  background: #252526;
  border-bottom: 1px solid #333;
}

.filename {
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.8125rem;
  color: #858585;
}

.code-wrapper {
  overflow-x: auto;
}

.code-wrapper pre {
  margin: 0;
  padding: 1rem;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.8125rem;
  line-height: 1.7;
  color: #d4d4d4;
  background: #1e1e1e;
}

@media (max-width: 640px) {
  .comparison-header {
    flex-direction: column;
    align-items: stretch;
  }

  .language-tabs {
    justify-content: center;
  }

  .tab-button {
    padding: 0.5rem 0.625rem;
    font-size: 0.8125rem;
  }

  .tab-label {
    display: none;
  }
}
</style>
