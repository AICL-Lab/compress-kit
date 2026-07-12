---
layout: home
---

<script setup>
import { onBeforeMount } from 'vue'
import HomeLanding from '../.vitepress/theme/components/HomeLanding.vue'
import { persistLocale } from '../.vitepress/theme/utils/language-preference.mjs'

onBeforeMount(() => {
  if (typeof window !== 'undefined') {
    persistLocale(window.localStorage, 'zh')
  }
})
</script>

<HomeLanding locale="zh" />
