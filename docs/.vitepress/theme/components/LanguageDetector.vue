<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  // Only run on root path (landing page)
  const pathname = window.location.pathname
  if (pathname !== '/compress-kit/' && pathname !== '/compress-kit') {
    return
  }

  // Check if user has already made a language choice
  const savedLang = localStorage.getItem('compresskit-lang')
  if (savedLang) {
    return
  }

  // Detect browser language
  const browserLang = navigator.language || (navigator as any).userLanguage
  const preferredLang = browserLang?.toLowerCase().startsWith('zh') ? 'zh' : 'en'

  // Save preference
  localStorage.setItem('compresskit-lang', preferredLang)

  // Redirect to appropriate language
  window.location.href = `/compress-kit/${preferredLang}/`
})
</script>

<template>
  <!-- This component doesn't render anything visible -->
</template>
