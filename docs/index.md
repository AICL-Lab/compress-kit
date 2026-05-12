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

  const base = import.meta.env.BASE_URL
  const userLang = navigator.language || navigator.userLanguage || ''
  const savedLang = localStorage.getItem('docs-lang-preference')

  // Use saved preference if exists, otherwise use browser language
  const targetLang = savedLang || (userLang.startsWith('zh') ? `${base}zh/` : `${base}en/`)
  const currentPath = window.location.pathname

  // Only redirect if not already on correct language path
  if (targetLang.endsWith('/zh/') && !currentPath.includes('/zh/')) {
    window.location.replace(`${base}zh/`)
  } else if (targetLang.endsWith('/en/') && !currentPath.includes('/en/')) {
    window.location.replace(`${base}en/`)
  }
})
</script>
