<script setup lang="ts">
import { backend } from "../../../../wailsjs/go/models";
import RepoStats = backend.RepoStats;

defineProps<{
  currentRepoStats: RepoStats | null;
  formattedSizeMb: string;
  formattedFirstCommit: string;
  formattedLastCommit: string;
}>();
</script>

<template>
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
</template>

<style scoped>
.loading-placeholder {
  opacity: 0.5;
  filter: blur(2px);
  transition: all 0.3s ease;
}

.repo-stats-table td {
  height: 24px;
}
</style>
