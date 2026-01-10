import { ref, computed, type Ref, type ComputedRef } from 'vue';
import type { RepoTab } from '@/types/git.types';
import * as App from '../../wailsjs/go/backend/App';
import { GetRepoStats } from '../../wailsjs/go/backend/App';
import { backend } from "../../wailsjs/go/models";
import dayjs from "dayjs";
import RepoStats = backend.RepoStats;

export function useRepoStats(activeTab: ComputedRef<RepoTab | null>, homeDir: Ref<string>) {
  const currentRepoStats = ref<RepoStats | null>(null);
  const readmeHtml = ref<string>('');

  const formattedPath = computed(() => {
    if (!activeTab.value) return '';
    const path = activeTab.value.path;
    if (homeDir.value && path.startsWith(homeDir.value)) {
      return path.replace(homeDir.value, '~');
    }
    return path;
  });

  const formattedSizeMb = computed(() => {
    if (!currentRepoStats.value) return ''
    return currentRepoStats.value.sizeMb.toFixed(1);
  })

  const formattedFirstCommit = computed(() => {
    if (!currentRepoStats.value) return ''
    return dayjs(currentRepoStats.value.firstCommit).format('D MMM YYYY [at] HH:mm');
  })

  const formattedLastCommit = computed(() => {
    if (!currentRepoStats.value) return ''
    return dayjs(currentRepoStats.value.lastCommit).format('D MMM YYYY [at] HH:mm');
  })

  async function loadRepoInfo(path: string) {
    try {
      currentRepoStats.value = await GetRepoStats(path);
      if (!currentRepoStats.value.isClean) {
        console.log("Uncommitted changes:", currentRepoStats.value.modifiedFiles);
      }
      readmeHtml.value = await App.GetRepoReadme(path);
    } catch (err) {
      console.error("Could not read git repo: " + err);
    }
  }

  return {
    currentRepoStats,
    readmeHtml,
    formattedPath,
    formattedSizeMb,
    formattedFirstCommit,
    formattedLastCommit,
    loadRepoInfo
  };
}
