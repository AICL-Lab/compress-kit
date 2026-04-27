---
layout: home
---

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vitepress'

const router = useRouter()

onMounted(() => {
  // Redirect to English version
  router.go('/en/')
})
</script>

<div class="redirect-notice">
  <p>Redirecting to <a href="/en/">English documentation</a>...</p>
  <p><a href="/zh/">中文文档</a></p>
</div>

<style>
.redirect-notice {
  text-align: center;
  padding: 4rem 2rem;
}
.redirect-notice p {
  margin: 1rem 0;
  font-size: 1.1rem;
}
</style>
