<template>
  <div class="animated-hero">
    <slot />
  </div>
</template>

<script setup>
</script>

<style scoped>
.animated-hero {
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
