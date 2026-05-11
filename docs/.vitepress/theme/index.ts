import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import { h } from 'vue'
import './styles/vars.css'
import './styles/animations.css'
import './styles/components.css'
import './styles/custom.css'

// Import custom components
import StatsBar from './components/StatsBar.vue'
import AlgorithmGrid from './components/AlgorithmGrid.vue'
import BenchmarkChart from './components/BenchmarkChart.vue'
import CodeComparison from './components/CodeComparison.vue'
import CustomFooter from './components/CustomFooter.vue'

export default {
  extends: DefaultTheme,
  Layout: () => {
    return h(DefaultTheme.Layout, null, {
      'layout-bottom': () => h(CustomFooter),
    })
  },
  enhanceApp({ app }) {
    // Register custom components globally
    app.component('StatsBar', StatsBar)
    app.component('AlgorithmGrid', AlgorithmGrid)
    app.component('BenchmarkChart', BenchmarkChart)
    app.component('CodeComparison', CodeComparison)
  },
} satisfies Theme
