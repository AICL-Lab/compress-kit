<template>
  <div class="feature-card">
    <slot />
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.feature-card {
  border-radius: 8px;
  padding: 1rem;
}
</style>
