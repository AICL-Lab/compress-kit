<script setup lang="ts">
import { computed } from 'vue'
import { withBase } from 'vitepress'
import { getHomepageContent } from '../../data/site-content.mjs'

const props = defineProps<{
  locale: 'en' | 'zh'
}>()

const content = computed(() => getHomepageContent(props.locale))
</script>

<template>
  <div class="home-header">
    <div class="home-header-left">
      <div class="home-logo">CK</div>
      <div>
        <span class="home-title">CompressKit</span>
        <span class="home-subtitle">{{ content.subtitle }}</span>
      </div>
    </div>
    <div class="home-nav">
      <a v-for="link in content.navLinks" :key="link.text" :href="withBase(link.link)">{{ link.text }}</a>
    </div>
  </div>

  <div class="home-intro-row">
    <div class="home-intro">{{ content.intro }}</div>
    <div class="home-stats">
      <span v-for="stat in content.stats" :key="stat"><strong>{{ stat }}</strong></span>
    </div>
  </div>

  <h2>{{ content.sections.algorithms }}</h2>

  <div class="feature-map">
    <div v-for="feature in content.featureCards" :key="feature.id" class="feature-card">
      <div class="feature-card-title">{{ feature.title }}</div>
      <div class="feature-card-desc">{{ feature.description }}</div>
      <div class="feature-tags">
        <template v-for="tag in feature.tags" :key="`${feature.id}-${tag.label}`">
          <a v-if="tag.link" :href="withBase(tag.link)" class="feature-tag">{{ tag.label }}</a>
          <span v-else class="feature-tag">{{ tag.label }}</span>
        </template>
      </div>
    </div>
  </div>

  <div class="quick-start">
    <div class="quick-start-title">{{ content.sections.quickStart }}</div>
    <div class="quick-start-content">
      <div class="command-block">
        <code>{{ content.quickStartCommand }}</code>
      </div>
    </div>
  </div>
</template>
