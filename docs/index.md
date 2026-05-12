---
layout: home

hero:
  name: CompressKit
  text: ' '
  actions:
    - theme: brand
      text: English
      link: /en/
    - theme: alt
      text: 中文
      link: /zh/
---

<script setup>
import { onBeforeMount } from 'vue'

// Client-only execution for SSR compatibility
onBeforeMount(() => {
  if (typeof window === 'undefined') return

  const userLang = navigator.language || navigator.userLanguage || ''
  const savedLang = localStorage.getItem('docs-lang-preference')

  // Use saved preference if exists, otherwise use browser language
  const targetLang = savedLang || (userLang.startsWith('zh') ? '/zh/' : '/en/')
  const currentPath = window.location.pathname

  // Only redirect if not already on correct language path
  if (targetLang === '/zh/' && !currentPath.startsWith('/zh')) {
    window.location.replace('/zh/')
  } else if (targetLang === '/en/' && !currentPath.startsWith('/en')) {
    window.location.replace('/en/')
  }
})
</script>
