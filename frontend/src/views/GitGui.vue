<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import GitStatusBar from "@/components/GitGui/GitStatusBar.vue";
import GitSettingsModal from "@/components/GitGui/GitSettingsModal.vue";
import GitInitModal from "@/components/GitGui/GitInitModal.vue";
import Sidebar from "@/components/GitGui/Sidebar.vue";
import GitTabHeader from "@/components/GitGui/GitTabHeader.vue";
import GitVerticalNav from "@/components/GitGui/GitVerticalNav.vue";
import RepoInfoView from "@/components/GitGui/RepoInfoView.vue";
import ChangesView from "@/components/GitGui/ChangesView.vue";
import CommitView from "@/components/GitGui/CommitView.vue";
import BranchView from "@/components/GitGui/BranchView.vue";

import * as App from "../../wailsjs/go/backend/App";
import { useRepoTabs } from "@/composables/useRepoTabs";
import { useGitActions } from "@/composables/useGitActions";
import { useRepoStats } from "@/composables/useRepoStats";

const showSettings = ref(false);
const recentRepos = ref<{ name: string, path: string }[]>([]);
const homeDir = ref<string>('');

const {
  tabs,
  activeTabId,
  activeTab,
  newTabObject,
  saveState,
  setActiveTab,
  closeTab,
  setActiveVerticalTab,
  STORAGE_KEY
} = useRepoTabs();

const {
  showInitModal,
  pendingRepoPath,
  openRecentRepo,
  openRepo,
  initializeRepo,
  cancelInit,
  addToRecent
} = useGitActions(tabs, activeTabId, recentRepos, newTabObject, (recent) => saveState(recent), (id, recent) => setActiveTab(id, recent));

const {
  currentRepoStats,
  readmeHtml,
  formattedPath,
  formattedSizeMb,
  formattedFirstCommit,
  formattedLastCommit,
  loadRepoInfo
} = useRepoStats(activeTab, homeDir);

const openInFileManager = () => {
  if (activeTab.value) {
    App.OpenInFileManager(activeTab.value.path);
  }
};

const loadState = async () => {
  const saved = localStorage.getItem(STORAGE_KEY);
  if (saved) {
    try {
      const state = JSON.parse(saved);
      if (state.recentRepos && Array.isArray(state.recentRepos)) {
        recentRepos.value = state.recentRepos;
      }
      if (state.tabs && Array.isArray(state.tabs)) {
        tabs.value = state.tabs.map((t: any) => ({
          ...t,
          remotes: t.remotes || [],
          statusFiles: t.statusFiles || [],
          activeVerticalTab: t.activeVerticalTab || 'info',
          loading: true,
          error: null
        }));
        activeTabId.value = state.activeTabId;

        if (!activeTabId.value && tabs.value.length > 0) {
          activeTabId.value = tabs.value[0].id;
        }
      }
    } catch (err) {
      console.error('Failed to load Git GUI state:', err);
    }
  }
};

watch(activeTabId, (newId) => {
  if (newId) {
    const tab = tabs.value.find(t => t.id === newId);
    if (tab) {
      currentRepoStats.value = null;
      loadRepoInfo(tab.path);
    }
  }
});

onMounted(async () => {
  await loadState();
  try {
    homeDir.value = await App.GetHomeDir();
  } catch (err) {
    console.error('Failed to get home directory:', err);
  }
});

</script>

<template>
  <Sidebar :recent-repos="recentRepos" :active-tab="activeTab" :current-repo-stats="currentRepoStats" @open-settings="showSettings = true"
           @open-repo="openRepo"
           @open-recent-repo="openRecentRepo"
           @select-vertical-tab="(tab) => setActiveVerticalTab(tab, recentRepos)"
  />

  <div class="git-gui-container h-100 d-flex flex-column">
    <GitStatusBar
        :active-tab="activeTab"
        class="mb-0 flex-shrink-0"
    />

    <GitTabHeader
        v-model:tabs="tabs"
        :active-tab-id="activeTabId"
        @set-active-tab="(id) => setActiveTab(id, recentRepos)"
        @close-tab="(id) => closeTab(id, recentRepos)"
        @save-state="saveState(recentRepos)"
    />

    <div v-if="activeTab" class="flex-grow-1 d-flex overflow-hidden">
      <GitVerticalNav
          :active-vertical-tab="activeTab.activeVerticalTab"
          @set-active-vertical-tab="(tab) => setActiveVerticalTab(tab, recentRepos)"
      />

      <div class="flex-grow-1 d-flex flex-column overflow-hidden" :key="activeTab.id">
        <RepoInfoView
            v-if="activeTab.activeVerticalTab === 'info'"
            :repo-name="activeTab.name"
            :formatted-path="formattedPath"
            :current-repo-stats="currentRepoStats"
            :formatted-size-mb="formattedSizeMb"
            :formatted-first-commit="formattedFirstCommit"
            :formatted-last-commit="formattedLastCommit"
            :readme-html="readmeHtml"
            @open-in-file-manager="openInFileManager"
        />

        <ChangesView
            v-else-if="activeTab.activeVerticalTab === 'local-changes'"
            :repo-path="activeTab.path"
            @refresh-stats="loadRepoInfo(activeTab.path)"
        />

        <CommitView
            v-else-if="activeTab.activeVerticalTab === 'commit'"
            :repo-path="activeTab.path"
        />

        <BranchView v-else-if="activeTab.activeVerticalTab === 'placeholder1'"/>

        <template v-else-if="activeTab.activeVerticalTab === 'placeholder2'">
          <div class="h-100 d-flex flex-column align-items-center justify-content-center text-muted">
            <i class="ti ti-settings fs-1 mb-3"></i>
            <h3>Repository Settings</h3>
            <p>Local repository settings and configurations.</p>
          </div>
        </template>
      </div>
    </div>

    <div v-else class="h-100 d-flex flex-column align-items-center justify-content-center text-muted">
      <i class="ti ti-git-branch fs-1 mb-3"></i>
      <p>No repository open. Click "Open Repository" to get started.</p>
    </div>
  </div>

  <GitSettingsModal
      :show="showSettings"
      @close="showSettings = false"
  />

  <GitInitModal
      :show="showInitModal"
      :repo-path="pendingRepoPath"
      @close="cancelInit"
      @initialize="initializeRepo"
  />
</template>

<style scoped>
</style>