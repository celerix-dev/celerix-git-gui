<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import GitStatusBar from "@/components/GitGui/Navigation/GitStatusBar.vue";
import GitSettingsModal from "@/components/GitGui/Modals/GitSettingsModal.vue";
import GitInitModal from "@/components/GitGui/Modals/GitInitModal.vue";
import GitBranchModal from "@/components/GitGui/Modals/GitBranchModal.vue";
import GitTagModal from "@/components/GitGui/Modals/GitTagModal.vue";
import Sidebar from "@/components/GitGui/Navigation/Sidebar.vue";
import GitTabHeader from "@/components/GitGui/Navigation/GitTabHeader.vue";
import GitVerticalNav from "@/components/GitGui/Navigation/GitVerticalNav.vue";
import RepoInfoView from "@/components/GitGui/RepoInfoView/RepoInfoView.vue";
import ChangesView from "@/components/GitGui/ChangesView/ChangesView.vue";
import CommitView from "@/components/GitGui/CommitView/CommitView.vue";
import BranchView from "@/components/GitGui/BranchView/BranchView.vue";

import * as App from "../../wailsjs/go/backend/App";
import { useRepoTabs } from "@/composables/useRepoTabs";
import { useGitActions } from "@/composables/useGitActions";
import { useRepoStats } from "@/composables/useRepoStats";

const showSettings = ref(false);
const showBranchModal = ref(false);
const showTagModal = ref(false);
const branchModalFrom = ref('');
const tagModalFrom = ref('');
const modalLoading = ref(false);

const recentRepos = ref<{ name: string, path: string }[]>([]);
const homeDir = ref<string>('');

const refreshCounter = ref(0);

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

const triggerRefresh = () => {
  refreshCounter.value++;
};

const {
  showInitModal,
  pendingRepoPath,
  openRecentRepo,
  openRepo,
  initializeRepo,
  cancelInit,
  addToRecent,
  checkoutBranch,
  createBranch,
  createTag,
  fetchRepo,
  pullRepo,
  pushRepo
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

const refreshAll = (path: string) => {
  loadRepoInfo(path);
  triggerRefresh();
};

const handleCreateBranch = (data: { name: string, checkout: boolean }) => {
  if (activeTab.value) {
    const path = activeTab.value.path;
    modalLoading.value = true;
    createBranch(path, data.name, data.checkout).then(() => {
      showBranchModal.value = false;
      refreshAll(path);
    }).finally(() => {
      modalLoading.value = false;
    });
  }
};

const handleCreateTag = (data: { name: string, message: string }) => {
  if (activeTab.value) {
    const path = activeTab.value.path;
    modalLoading.value = true;
    createTag(path, data.name, data.message).then(() => {
      showTagModal.value = false;
      refreshAll(path);
    }).finally(() => {
      modalLoading.value = false;
    });
  }
};

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
           @checkout-branch="(name, isRemote) => {
             if (activeTab) {
               const path = activeTab.path;
               checkoutBranch(path, name, isRemote).then(() => refreshAll(path));
             }
           }"
           @new-branch="(from) => { branchModalFrom = from; showBranchModal = true; }"
           @new-tag="(from) => { tagModalFrom = from; showTagModal = true; }"
  />

  <div class="git-gui-container h-100 d-flex flex-column">
    <GitStatusBar
        :active-tab="activeTab"
        @fetch="() => { if (activeTab) { const path = activeTab.path; fetchRepo(path).then((res) => { if(res === 'open-settings') showSettings = true; refreshAll(path); }) } }"
        @pull="() => { if (activeTab) { const path = activeTab.path; pullRepo(path).then((res) => { if(res === 'open-settings') showSettings = true; refreshAll(path); }) } }"
        @push="() => { if (activeTab) { const path = activeTab.path; pushRepo(path).then((res) => { if(res === 'open-settings') showSettings = true; refreshAll(path); }) } }"
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
            @refresh-stats="refreshAll(activeTab.path)"
        />

        <CommitView
            v-else-if="activeTab.activeVerticalTab === 'commit'"
            :repo-path="activeTab.path"
            :refresh-counter="refreshCounter"
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

  <GitBranchModal
      :show="showBranchModal"
      :from-branch="branchModalFrom"
      :loading="modalLoading"
      @close="showBranchModal = false"
      @create="handleCreateBranch"
  />

  <GitTagModal
      :show="showTagModal"
      :from-branch="tagModalFrom"
      :loading="modalLoading"
      @close="showTagModal = false"
      @create="handleCreateTag"
  />
</template>

<style scoped>
</style>