<script setup lang="ts">
import { backend } from "../../../../wailsjs/go/models";
import RepoStats = backend.RepoStats;
import * as App from "../../../../wailsjs/go/backend/App";

import RepoStatsTable from './RepoStatsTable.vue';
import ReadmeViewer from './ReadmeViewer.vue';

defineProps<{
  repoName: string;
  formattedPath: string;
  currentRepoStats: RepoStats | null;
  formattedSizeMb: string;
  formattedFirstCommit: string;
  formattedLastCommit: string;
  readmeHtml: string;
}>();

const emit = defineEmits<{
  (e: 'openInFileManager'): void;
}>();

const openInBrowser = (url: string) => {
  App.OpenInBrowser(url);
};
</script>

<template>
  <div class="flex-grow-1 d-flex flex-column overflow-hidden">
    <!-- Fixed Header: Name, Path and Stats -->
    <div class="p-3 ps-4 pb-0">
      <h3>{{ repoName }}</h3>
      <div class="d-inline-flex">
        <p class="text-muted small d-flex align-items-center me-5">
          {{ formattedPath }}
          <i class="ti ti-external-link ms-2 cursor-pointer" @click="emit('openInFileManager')"
             title="Open in File Manager"></i>
        </p>
      </div>

      <RepoStatsTable 
        :current-repo-stats="currentRepoStats"
        :formatted-size-mb="formattedSizeMb"
        :formatted-first-commit="formattedFirstCommit"
        :formatted-last-commit="formattedLastCommit"
      />

      <hr class="mb-1"/>
      <div class="badge text-bg-secondary">README.md</div>
      <hr class="mt-1 mb-0"/>
    </div>

    <ReadmeViewer :readme-html="readmeHtml" />
  </div>
</template>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
