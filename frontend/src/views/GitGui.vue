<script setup lang="ts">
import GitStatusBar from "@/components/GitGui/GitStatusBar.vue";
import {computed, onMounted, ref, watch} from "vue";
import GitSettingsModal from "@/components/GitGui/GitSettingsModal.vue";
import GitInitModal from "@/components/GitGui/GitInitModal.vue";
import type {RepoTab} from "@/components/GitGui/types";
import * as App from "../../wailsjs/go/backend/App";
import {GetRepoStats} from "../../wailsjs/go/backend/App";
import draggable from "vuedraggable";
import Sidebar from "@/components/GitGui/Sidebar.vue";
import {backend} from "../../wailsjs/go/models";
import dayjs from "dayjs";
import RepoStats = backend.RepoStats;

const showSettings = ref(false);
const tabs = ref<RepoTab[]>([]);
const activeTabId = ref<string | null>(null);
const recentRepos = ref<{ name: string, path: string }[]>([]);
const homeDir = ref<string>('');
const currentRepoStats = ref<RepoStats | null>(null);
const readmeHtml = ref<string>('');
const showInitModal = ref(false);
const pendingRepoPath = ref<string | null>(null);
const pendingRepoName = ref<string | null>(null);

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

const formattedPath = computed(() => {
  if (!activeTab.value) return '';
  const path = activeTab.value.path;
  if (homeDir.value && path.startsWith(homeDir.value)) {
    return path.replace(homeDir.value, '~');
  }
  return path;
});

const openInFileManager = () => {
  if (activeTab.value) {
    App.OpenInFileManager(activeTab.value.path);
  }
};

const openInBrowser = (url: string) => {
  App.OpenInBrowser(url);
};

const STORAGE_KEY = 'git-gui-tabs';

const setActiveTab = (id: string) => {
  console.log('Setting active tab:', id);
  activeTabId.value = id;
  saveState();
};

const setActiveVerticalTab = (tabName: 'info' | 'commit' | 'placeholder1' | 'placeholder2') => {
  if (activeTab.value) {
    activeTab.value.activeVerticalTab = tabName;
    saveState();
  }
};

const closeTab = (id: string) => {
  console.log('Closing tab:', id);
  const index = tabs.value.findIndex(t => t.id === id);
  if (index !== -1) {
    tabs.value.splice(index, 1);
    if (activeTabId.value === id) {
      activeTabId.value = tabs.value.length > 0 ? tabs.value[tabs.value.length - 1].id : null;
      console.log('Active tab changed after close to:', activeTabId.value);
    }
    saveState();
  }
};

const saveState = () => {
  const state = {
    tabs: tabs.value.map(t => ({
      id: t.id,
      name: t.name,
      path: t.path,
      activeVerticalTab: t.activeVerticalTab
    })),
    activeTabId: activeTabId.value,
    recentRepos: recentRepos.value
  };
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
};

const addToRecent = (name: string, path: string) => {
  const existingIndex = recentRepos.value.findIndex(r => r.path === path);
  if (existingIndex !== -1) {
    recentRepos.value.splice(existingIndex, 1);
  }
  recentRepos.value.unshift({name, path});
  if (recentRepos.value.length > 5) {
    recentRepos.value = recentRepos.value.slice(0, 5);
  }
  saveState();
};

const openRecentRepo = async (repo: { name: string, path: string }) => {
  const existingTab = tabs.value.find(t => t.path === repo.path);
  if (existingTab) {
    setActiveTab(existingTab.id);
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
  saveState();
};

const openRepo = async () => {
  try {
    const selected = await App.SelectDirectory('Open Git Repository');

    if (selected) {
      const path = selected;
      if (!path) return;

      console.log('Opening repo at path:', path);

      const isRepo = await App.IsGitRepo(path);
      const repoName = path.split(/[\\/]/).pop() || path;

      if (!isRepo) {
        pendingRepoPath.value = path;
        pendingRepoName.value = repoName;
        showInitModal.value = true;
        return;
      }

      // Check if already open
      const existingTab = tabs.value.find(t => t.path === path);
      if (existingTab) {
        console.log('Repo already open, switching to tab:', existingTab.id);
        setActiveTab(existingTab.id);
        return;
      }

      const tabId = Math.random().toString(36).substring(7);
      const newTab = newTabObject(tabId, repoName, path);

      tabs.value.push(newTab);
      activeTabId.value = tabId;
      addToRecent(repoName, path);
      saveState();

      // await loadRepoData(tabId);
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

    // Open it
    const tabId = Math.random().toString(36).substring(7);
    const newTab = newTabObject(tabId, name, path);
    tabs.value.push(newTab);
    activeTabId.value = tabId;
    addToRecent(name, path);
    saveState();
  } catch (err) {
    alert("Failed to initialize repository: " + err);
  }
};

const cancelInit = () => {
  showInitModal.value = false;
  pendingRepoPath.value = null;
  pendingRepoName.value = null;
};

const loadState = async () => {
  const saved = localStorage.getItem(STORAGE_KEY);
  console.log('loadState: Found saved state:', !!saved);
  if (saved) {
    try {
      const state = JSON.parse(saved);
      console.log('loadState: Parsed state:', state);
      if (state.recentRepos && Array.isArray(state.recentRepos)) {
        recentRepos.value = state.recentRepos;
      }
      if (state.tabs && Array.isArray(state.tabs)) {
        console.log(`loadState: Restoring ${state.tabs.length} tabs`);
        tabs.value = state.tabs.map((t: any) => ({
          ...t,
          remotes: t.remotes || [],
          statusFiles: t.statusFiles || [],
          activeVerticalTab: t.activeVerticalTab || 'info',
          loading: true,
          error: null
        }));
        activeTabId.value = state.activeTabId;
        console.log('loadState: Restored activeTabId:', activeTabId.value);

        if (!activeTabId.value && tabs.value.length > 0) {
          activeTabId.value = tabs.value[0].id;
          console.log('loadState: Auto-setting activeTabId to first tab:', activeTabId.value);
        }

        // Reload data for each tab
        for (const tab of tabs.value) {
          // loadRepoData(tab.id);
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
    alert("Could not read git repo: " + err);
  }
}

</script>

<template>
  <Sidebar :recent-repos="recentRepos" @open-settings="showSettings = true"
           @open-repo="openRepo"
           @open-recent-repo="openRecentRepo"
  />

  <div class="git-gui-container h-100 d-flex flex-column">
    <GitStatusBar
        :active-tab="activeTab"
        class="mb-0 flex-shrink-0"
    />

    <!-- Tabs -->
    <div class="card border-radius-0 m-1 mb-0 ms-0 me-0 border-start-0 p-0 flex-shrink-0" v-if="tabs.length > 0"
         style="border-radius:0">
      <div class="card-body p-2">
        <draggable
            v-model="tabs"
            item-key="id"
            tag="ul"
            class="nav nav-pills nav-justified"
            ghost-class="ghost-tab"
            @end="saveState"
        >
          <template #item="{ element: tab }">
            <li class="nav-item">
              <div
                  :class="['nav-link d-flex align-items-center cursor-pointer', { active: activeTabId === tab.id }]"
                  @click="setActiveTab(tab.id)"
              >
                <i class="ti ti-x ms-1 close-icon" @click.stop="closeTab(tab.id)"></i>
                <div class="ms-2 flex-grow-1 d-flex justify-content-center">{{ tab.name }}</div>
              </div>
            </li>
          </template>
        </draggable>
      </div>
    </div>

    <!-- Tab Content -->
    <div v-if="activeTab" class="flex-grow-1 d-flex overflow-hidden">
      <!-- Vertical Tab Bar -->
      <div class="vertical-tabs d-flex flex-column border-end py-2">
        <div
            :class="['vertical-tab-item', { active: activeTab.activeVerticalTab === 'info' }]"
            @click="setActiveVerticalTab('info')"
            title="Repository Info"
        >
          <i class="ti ti-info-circle fs-4"></i>
        </div>
        <div
            :class="['vertical-tab-item', { active: activeTab.activeVerticalTab === 'commit' }]"
            @click="setActiveVerticalTab('commit')"
            title="Commit Info"
        >
          <i class="ti ti-git-commit fs-4"></i>
        </div>
        <div
            :class="['vertical-tab-item', { active: activeTab.activeVerticalTab === 'placeholder1' }]"
            @click="setActiveVerticalTab('placeholder1')"
            title="Branch Management"
        >
          <i class="ti ti-git-branch fs-4"></i>
        </div>
        <div
            :class="['vertical-tab-item', { active: activeTab.activeVerticalTab === 'placeholder2' }]"
            @click="setActiveVerticalTab('placeholder2')"
            title="Settings"
        >
          <i class="ti ti-settings fs-4"></i>
        </div>
      </div>

      <!-- Main Content Area -->
      <div class="flex-grow-1 d-flex flex-column overflow-hidden">
        <!-- Tab: Info -->
        <template v-if="activeTab.activeVerticalTab === 'info'">
          <!-- Fixed Header: Name, Path and Stats -->
          <div class="p-3 ps-4 pb-0">
            <h3>{{ activeTab.name }}</h3>
            <div class="d-inline-flex">
              <p class="text-muted small d-flex align-items-center me-5">
                {{ formattedPath }}
                <i class="ti ti-external-link ms-2 cursor-pointer" @click="openInFileManager"
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
        </template>

        <!-- Tab: Commit Info Placeholder -->
        <template v-else-if="activeTab.activeVerticalTab === 'commit'">
          <div class="h-100 d-flex flex-column align-items-center justify-content-center text-muted">
            <i class="ti ti-git-commit fs-1 mb-3"></i>
            <h3>Commit History</h3>
            <p>Commit history and details will be displayed here.</p>
          </div>
        </template>

        <!-- Placeholder Tabs -->
        <template v-else-if="activeTab.activeVerticalTab === 'placeholder1'">
          <div class="h-100 d-flex flex-column align-items-center justify-content-center text-muted">
            <i class="ti ti-git-branch fs-1 mb-3"></i>
            <h3>Branch Management</h3>
            <p>Branch operations and management will be available here.</p>
          </div>
        </template>

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
.cursor-pointer {
  cursor: pointer;
}

.nav-item {
  cursor: move;
}

.ghost-tab {
  opacity: 0.5;
  background: var(--bs-light);
}

.close-icon {
  visibility: hidden;
  opacity: 0;
  transition: opacity 0.2s, visibility 0.2s;
}

.nav-link:hover .close-icon {
  visibility: visible;
  opacity: 1;
}

.nav-link:not(.active):hover {
  background-color: rgba(0, 0, 0, 0.06);
}

[data-bs-theme='dark'] .nav-link:not(.active):hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.close-icon:hover {
  color: var(--bs-danger);
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

.loading-placeholder {
  opacity: 0.5;
  filter: blur(2px);
  transition: all 0.3s ease;
}

.repo-stats-table td {
  height: 24px;
}

.vertical-tabs {
  width: 50px;
  background-color: var(--bs-tertiary-bg);
}

.vertical-tab-item {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--bs-secondary-color);
  transition: all 0.2s ease;
  border-right: 3px solid transparent;
}

.vertical-tab-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--bs-primary);
}

[data-bs-theme='dark'] .vertical-tab-item:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.vertical-tab-item.active {
  color: var(--bs-primary);
  border-right-color: var(--bs-primary);
  background-color: var(--bs-body-bg);
}
</style>