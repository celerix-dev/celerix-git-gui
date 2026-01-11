<script setup lang="ts">
defineProps<{
  readmeHtml: string;
}>();
</script>

<template>
  <div class="flex-grow-1 overflow-auto p-3 pt-2 m-3 scroll-container">
    <!-- render README.md if it exists -->
    <div v-if="readmeHtml" class="readme-content" v-html="readmeHtml"></div>
    <p v-else class="text-muted italic">No README found for this repository.</p>
  </div>
</template>

<style scoped>
.readme-content :deep(h1), .readme-content :deep(h2), .readme-content :deep(h3) {
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid var(--bs-border-color);
  padding-bottom: 0.3rem;
}

.readme-content :deep(pre) {
  background-color: var(--bs-tertiary-bg);
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
}

.readme-content :deep(code) {
  font-family: var(--bs-font-monospace);
  background-color: var(--bs-tertiary-bg);
  padding: 0.2rem 0.4rem;
  border-radius: 0.25rem;
}

.readme-content :deep(img) {
  max-width: 100%;
}

.readme-content :deep(blockquote) {
  border-left: 4px solid var(--bs-border-color);
  padding-left: 1rem;
  color: var(--bs-secondary-color);
}
</style>
