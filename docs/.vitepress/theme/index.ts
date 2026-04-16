import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import { h } from 'vue'
import './styles/vars.css'
import './styles/custom.css'
import './styles/animations.css'
import './styles/components.css'

// Import custom components
import CustomFooter from './components/CustomFooter.vue'
import AnimatedHero from './components/AnimatedHero.vue'
import FeatureCard from './components/FeatureCard.vue'
import AlgorithmShowcase from './components/AlgorithmShowcase.vue'
import CodePlayground from './components/CodePlayground.vue'
import LanguageSwitcher from './components/LanguageSwitcher.vue'
import ReadingProgress from './components/ReadingProgress.vue'
import BackToTop from './components/BackToTop.vue'

export default {
  extends: DefaultTheme,
  Layout: () => {
    return h(DefaultTheme.Layout, null, {
      'layout-bottom': () => h(CustomFooter),
      'doc-top': () => h(ReadingProgress),
      'doc-bottom': () => h(BackToTop),
    })
  },
  enhanceApp({ app }) {
    // Register custom components
    app.component('AnimatedHero', AnimatedHero)
    app.component('FeatureCard', FeatureCard)
    app.component('AlgorithmShowcase', AlgorithmShowcase)
    app.component('CodePlayground', CodePlayground)
    app.component('LanguageSwitcher', LanguageSwitcher)
  },
} satisfies Theme
