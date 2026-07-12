import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import { onMounted, watch, nextTick } from 'vue'
import { useData } from 'vitepress'
import './styles/vars.css'

// Import custom components (保留核心组件)
import AlgorithmGrid from './components/AlgorithmGrid.vue'
import BenchmarkChart from './components/BenchmarkChart.vue'
import CodeComparison from './components/CodeComparison.vue'

export default {
  extends: DefaultTheme,
  enhanceApp({ app }) {
    // Register custom components globally
    app.component('AlgorithmGrid', AlgorithmGrid)
    app.component('BenchmarkChart', BenchmarkChart)
    app.component('CodeComparison', CodeComparison)
  },
  setup() {
    // Mermaid 深色主题动态切换
    const { isDark } = useData()

    const updateMermaidTheme = async () => {
      await nextTick()
      const mermaid = (window as any).mermaid
      if (mermaid) {
        mermaid.initialize({
          startOnLoad: false,
          theme: isDark.value ? 'dark' : 'default',
          themeVariables: {
            primaryColor: isDark.value ? '#60a5fa' : '#2563eb',
            primaryTextColor: isDark.value ? '#f1f5f9' : '#1e293b',
            primaryBorderColor: isDark.value ? '#334155' : '#e2e8f0',
            lineColor: isDark.value ? '#64748b' : '#94a3b8',
            secondaryColor: isDark.value ? '#1e293b' : '#f8fafc',
            tertiaryColor: isDark.value ? '#0f172a' : '#f1f5f9'
          }
        })
        mermaid.run()
      }
    }

    onMounted(updateMermaidTheme)
    watch(isDark, updateMermaidTheme)
  }
} satisfies Theme
