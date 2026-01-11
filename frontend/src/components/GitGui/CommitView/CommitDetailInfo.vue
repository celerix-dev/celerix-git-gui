<script setup lang="ts">
import { backend } from "../../../../wailsjs/go/models";
import GitCommit = backend.GitCommit;
import dayjs from "dayjs";

defineProps<{
  commit: GitCommit;
}>();

const formatDate = (date: any) => {
  return dayjs(date).format('D MMM YYYY HH:mm');
};
</script>

<template>
  <div class="commit-detail-info p-3 h-100 overflow-auto">
    <div class="d-flex justify-content-between align-items-start mb-2">
      <div>
        <h5 class="mb-1">{{ commit.subject }}</h5>
        <div class="text-muted small">
          <strong>{{ commit.authorName }}</strong> &lt;{{ commit.authorEmail }}&gt;
          <span class="mx-2">•</span>
          {{ formatDate(commit.date) }}
        </div>
      </div>
      <code class="bg-body px-2 py-1 rounded border">{{ commit.hash }}</code>
    </div>
    <div v-if="commit.body" class="mt-3 p-3 bg-body rounded border white-space-pre">{{ commit.body }}</div>
    <div class="mt-2 small text-muted">
      Parents: <code v-for="p in commit.parentHashes" :key="p" class="me-2">{{ p.substring(0, 7) }}</code>
    </div>
  </div>
</template>

<style scoped>
.white-space-pre {
    white-space: pre-wrap;
}
</style>
