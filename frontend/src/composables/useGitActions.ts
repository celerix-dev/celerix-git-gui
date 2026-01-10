import { ref, type Ref } from 'vue';
import * as App from '../../wailsjs/go/backend/App';
import type { RepoTab } from '@/types/git.types';

export function useGitActions(
  tabs: Ref<RepoTab[]>,
  activeTabId: Ref<string | null>,
  recentRepos: Ref<{ name: string, path: string }[]>,
  newTabObject: (tabId: string, repoName: string, repoPath: string) => RepoTab,
  saveState: (recent: { name: string, path: string }[]) => void,
  setActiveTab: (id: string, recent: { name: string, path: string }[]) => void
) {
  const showInitModal = ref(false);
  const pendingRepoPath = ref<string | null>(null);
  const pendingRepoName = ref<string | null>(null);

  const addToRecent = (name: string, path: string) => {
    const existingIndex = recentRepos.value.findIndex(r => r.path === path);
    if (existingIndex !== -1) {
      recentRepos.value.splice(existingIndex, 1);
    }
    recentRepos.value.unshift({ name, path });
    if (recentRepos.value.length > 5) {
      recentRepos.value = recentRepos.value.slice(0, 5);
    }
    saveState(recentRepos.value);
  };

  const openRecentRepo = async (repo: { name: string, path: string }) => {
    const existingTab = tabs.value.find(t => t.path === repo.path);
    if (existingTab) {
      setActiveTab(existingTab.id, recentRepos.value);
      return;
    }

    const isRepo = await App.IsGitRepo(repo.path);
    if (!isRepo) {
      pendingRepoPath.value = repo.path;
      pendingRepoName.value = repo.name;
      showInitModal.value = true;
      return;
    }

    const tabId = Math.random().toString(36).substring(7);
    const newTab = newTabObject(tabId, repo.name, repo.path);

    tabs.value.push(newTab);
    activeTabId.value = tabId;
    addToRecent(repo.name, repo.path);
    saveState(recentRepos.value);
  };

  const openRepo = async () => {
    try {
      const selected = await App.SelectDirectory('Open Git Repository');

      if (selected) {
        const path = selected;
        if (!path) return;

        const isRepo = await App.IsGitRepo(path);
        const repoName = path.split(/[\\/]/).pop() || path;

        if (!isRepo) {
          pendingRepoPath.value = path;
          pendingRepoName.value = repoName;
          showInitModal.value = true;
          return;
        }

        const existingTab = tabs.value.find(t => t.path === path);
        if (existingTab) {
          setActiveTab(existingTab.id, recentRepos.value);
          return;
        }

        const tabId = Math.random().toString(36).substring(7);
        const newTab = newTabObject(tabId, repoName, path);

        tabs.value.push(newTab);
        activeTabId.value = tabId;
        addToRecent(repoName, path);
        saveState(recentRepos.value);
      }
    } catch (err) {
      console.error('Failed to open repository:', err);
    }
  };

  const initializeRepo = async () => {
    if (!pendingRepoPath.value || !pendingRepoName.value) return;

    try {
      await App.GitInit(pendingRepoPath.value);
      const path = pendingRepoPath.value;
      const name = pendingRepoName.value;

      showInitModal.value = false;
      pendingRepoPath.value = null;
      pendingRepoName.value = null;

      const tabId = Math.random().toString(36).substring(7);
      const newTab = newTabObject(tabId, name, path);
      tabs.value.push(newTab);
      activeTabId.value = tabId;
      addToRecent(name, path);
      saveState(recentRepos.value);
    } catch (err) {
      alert("Failed to initialize repository: " + err);
    }
  };

  const cancelInit = () => {
    showInitModal.value = false;
    pendingRepoPath.value = null;
    pendingRepoName.value = null;
  };

  return {
    showInitModal,
    pendingRepoPath,
    pendingRepoName,
    openRecentRepo,
    openRepo,
    initializeRepo,
    cancelInit,
    addToRecent
  };
}
