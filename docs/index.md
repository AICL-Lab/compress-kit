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
import { getLandingRedirectTarget, readSavedLocale } from './.vitepress/theme/utils/language-preference.mjs'

// Client-only execution for SSR compatibility
onBeforeMount(() => {
  if (typeof window === 'undefined') return

  const target = getLandingRedirectTarget({
    pathname: window.location.pathname,
    base: import.meta.env.BASE_URL,
    savedLocale: readSavedLocale(window.localStorage),
    browserLanguage: navigator.language || navigator.userLanguage || '',
  })

  if (target && target !== window.location.pathname) {
    window.location.replace(target)
  }
})
</script>
