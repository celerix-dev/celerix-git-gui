import { ref, computed } from 'vue';
import type { RepoTab } from '@/types/git.types';

export function useRepoTabs() {
  const tabs = ref<RepoTab[]>([]);
  const activeTabId = ref<string | null>(null);
  const STORAGE_KEY = 'git-gui-tabs';

  const activeTab = computed(() => tabs.value.find(t => t.id === activeTabId.value) || null);

  const newTabObject = (tabId: string, repoName: string, repoPath: string): RepoTab => {
    return {
      id: tabId,
      name: repoName,
      path: repoPath,
      remotes: [],
      statusFiles: [],
      activeVerticalTab: 'info',
      loading: true,
      error: null
    }
  };

  const saveState = (recentRepos: { name: string, path: string }[]) => {
    const state = {
      tabs: tabs.value.map(t => ({
        id: t.id,
        name: t.name,
        path: t.path,
        activeVerticalTab: t.activeVerticalTab
      })),
      activeTabId: activeTabId.value,
      recentRepos: recentRepos
    };
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  };

  const setActiveTab = (id: string, recentRepos: { name: string, path: string }[]) => {
    activeTabId.value = id;
    saveState(recentRepos);
  };

  const closeTab = (id: string, recentRepos: { name: string, path: string }[]) => {
    const index = tabs.value.findIndex(t => t.id === id);
    if (index !== -1) {
      tabs.value.splice(index, 1);
      if (activeTabId.value === id) {
        activeTabId.value = tabs.value.length > 0 ? tabs.value[tabs.value.length - 1].id : null;
      }
      saveState(recentRepos);
    }
  };

  const setActiveVerticalTab = (tabName: 'info' | 'local-changes' | 'commit' | 'placeholder1' | 'placeholder2', recentRepos: { name: string, path: string }[]) => {
    if (activeTab.value) {
      activeTab.value.activeVerticalTab = tabName;
      saveState(recentRepos);
    }
  };

  return {
    tabs,
    activeTabId,
    activeTab,
    newTabObject,
    saveState,
    setActiveTab,
    closeTab,
    setActiveVerticalTab,
    STORAGE_KEY
  };
}
