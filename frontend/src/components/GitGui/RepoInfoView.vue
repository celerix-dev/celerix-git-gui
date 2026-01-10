<script setup lang="ts">
import { backend } from "../../../wailsjs/go/models";
import RepoStats = backend.RepoStats;
import * as App from "../../../wailsjs/go/backend/App";

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

      <table class="repo-stats-table">
        <tbody>
        <tr>
          <td class="text-end pe-2"><strong>Changed Files</strong></td>
          <td>
              <span :class="{ 'loading-placeholder': !currentRepoStats }">
                {{
                  currentRepoStats?.modifiedFiles ? currentRepoStats.modifiedFiles.length : (currentRepoStats ? 0 : '0')
                }}
              </span>
          </td>
        </tr>
        <tr>
          <td class="text-end pe-2"><strong>Repository Size</strong></td>
          <td>
              <span :class="{ 'loading-placeholder': !currentRepoStats }">
                {{ currentRepoStats ? formattedSizeMb + ' MB' : '0.0 MB' }}
              </span>
          </td>
        </tr>
        <tr>
          <td class="text-end pe-2"><strong>Commits</strong></td>
          <td>
              <span :class="{ 'loading-placeholder': !currentRepoStats }">
                {{ currentRepoStats?.commitCount ?? '0' }}
              </span>
          </td>
        </tr>
        <tr>
          <td class="text-end pe-2"><strong>Initial Commit</strong></td>
          <td>
              <span :class="{ 'loading-placeholder': !currentRepoStats }">
                {{ currentRepoStats ? formattedFirstCommit : '---' }}
              </span>
          </td>
        </tr>
        <tr>
          <td class="text-end pe-2"><strong>Last Commit</strong></td>
          <td>
              <span :class="{ 'loading-placeholder': !currentRepoStats }">
                {{ currentRepoStats ? formattedLastCommit : '---' }}
              </span>
          </td>
        </tr>
        </tbody>
      </table>

      <hr class="mb-1"/>
      <div class="badge text-bg-secondary">README.md</div>
      <hr class="mt-1 mb-0"/>
    </div>

    <!-- Scrollable README content -->
    <div class="flex-grow-1 overflow-auto p-3 pt-2 m-3 scroll-container">
      <!-- render README.md if it exists -->
      <div v-if="readmeHtml" class="readme-content" v-html="readmeHtml"></div>
      <p v-else class="text-muted italic">No README found for this repository.</p>
    </div>
  </div>
</template>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}

.loading-placeholder {
  opacity: 0.5;
  filter: blur(2px);
  transition: all 0.3s ease;
}

.repo-stats-table td {
  height: 24px;
}

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
