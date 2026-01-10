<script setup lang="ts">

import { ref } from 'vue';
import type { RepoTab } from "@/types/git.types";
import { backend } from "../../../wailsjs/go/models";
import RepoStats = backend.RepoStats;

const props = defineProps<{
  recentRepos: {name: string, path: string}[];
  activeTab: RepoTab | null;
  currentRepoStats: RepoStats | null;
}>();

const emit = defineEmits<{
  (e: 'openSettings'): void;
  (e: 'openRepo'): void;
  (e: 'openRecentRepo', repo: {name: string, path: string}): void;
  (e: 'selectVerticalTab', tab: 'info' | 'local-changes' | 'commit' | 'placeholder1' | 'placeholder2'): void;
}>();

const collapsed = ref({
  branches: false, // Branches default opened (false means NOT collapsed)
  remotes: true,
  tags: true,
  stashes: true
});

const collapsedRemotes = ref<Record<string, boolean>>({
  'origin': false // origin default opened
});

const toggleSection = (section: keyof typeof collapsed.value) => {
  collapsed.value[section] = !collapsed.value[section];
};

const toggleRemote = (remoteName: string) => {
  if (collapsedRemotes.value[remoteName] === undefined) {
    collapsedRemotes.value[remoteName] = false;
  } else {
    collapsedRemotes.value[remoteName] = !collapsedRemotes.value[remoteName];
  }
};

</script>

<template>
  <nav class="sidebar offcanvas-start offcanvas-md" tabindex="-1" id="donkers-sidebar">
    <div class="offcanvas-body p-0 d-flex flex-column">
      <!-- Top Part: Management Actions & Active Repo -->
      <div class="sidebar-management-section p-2 border-bottom">
        <div class="sidebar-group">
          <h6 class="sidebar-header">Actions</h6>
          <div class="list-group list-group-flush">
            <a href="javascript:void(0);" class="list-group-item list-group-item-action d-flex align-items-center border-0" @click.prevent="emit('openRepo')">
              <i class="ti ti-folder-open me-2"></i>
              <span>Open Repository</span>
            </a>
            <a href="javascript:void(0);" class="list-group-item list-group-item-action d-flex align-items-center border-0" @click.prevent="emit('openSettings')">
              <i class="ti ti-settings me-2"></i>
              <span>Settings</span>
            </a>
          </div>

          <template v-if="props.activeTab">
            <h6 class="sidebar-header text-truncate mt-2" :title="props.activeTab.name">{{ props.activeTab.name }}</h6>
            <div class="list-group list-group-flush">
              <a href="javascript:void(0);" :class="['list-group-item list-group-item-action d-flex align-items-center border-0', { active: props.activeTab.activeVerticalTab === 'local-changes' }]" @click="emit('selectVerticalTab', 'local-changes')">
                <i class="ti ti-edit me-2"></i>
                <span>Local changes {{ props.currentRepoStats?.modifiedFiles?.length ? `(${props.currentRepoStats.modifiedFiles.length})` : '' }}</span>
              </a>
              <a href="javascript:void(0);" :class="['list-group-item list-group-item-action d-flex align-items-center border-0', { active: props.activeTab.activeVerticalTab === 'commit' }]" @click="emit('selectVerticalTab', 'commit')">
                <i class="ti ti-history me-2"></i>
                <span>All Commits</span>
              </a>
            </div>
          </template>
        </div>
      </div>

      <!-- Bottom Part: Content Navigation -->
      <ul class="sidebar-nav flex-grow-1 overflow-auto py-2" data-key="dashboard" aria-expanded="true">
        <li v-if="props.activeTab">
          <div class="px-3 mb-3 mt-2">
            <input type="text" class="form-control form-control-sm filter-input" placeholder="Filter"/>
          </div>
          <h6 class="sidebar-header collapsible" @click="toggleSection('branches')">
            <i :class="['ti ti-chevron-right me-1 transition-icon', { 'rotate-90': !collapsed.branches }]"></i>
            <span>Branches</span>
          </h6>
          <div v-if="!collapsed.branches" class="list-group list-group-flush">
            <a v-for="branch in props.currentRepoStats?.branches || []"
               :key="branch"
               href="javascript:void(0);"
               :class="['list-group-item list-group-item-action d-flex align-items-center border-0', { 'fw-bold': branch === props.currentRepoStats?.currentBranch }]"
            >
              <i :class="['ti me-2', branch === props.currentRepoStats?.currentBranch ? 'ti-check' : 'ti-git-branch']"></i>
              <span class="text-truncate">{{ branch }}</span>
            </a>
            <div v-if="!props.currentRepoStats?.branches?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
              <span class="ms-4">No branches found</span>
            </div>
          </div>
        </li>
        <li v-if="props.activeTab">
          <h6 class="sidebar-header collapsible" @click="toggleSection('remotes')">
            <i :class="['ti ti-chevron-right me-1 transition-icon', { 'rotate-90': !collapsed.remotes }]"></i>
            <span>Remotes</span>
          </h6>
          <div v-if="!collapsed.remotes" class="list-group list-group-flush">
            <template v-for="remote in props.currentRepoStats?.remotes || []" :key="remote.name">
              <div 
                class="list-group-item d-flex align-items-center border-0 cursor-pointer py-1"
                @click="toggleRemote(remote.name)"
              >
                <i :class="['ti ti-chevron-right me-1 transition-icon', { 'rotate-90': !collapsedRemotes[remote.name] }]"></i>
                <i class="ti ti-cloud me-2"></i>
                <span class="text-truncate">{{ remote.name }}</span>
              </div>
              <div v-if="!collapsedRemotes[remote.name]" class="remote-branches ms-3">
                <a v-for="branch in (remote.branches || []).filter(b => !b.endsWith('/HEAD'))"
                   :key="branch"
                   href="javascript:void(0);"
                   class="list-group-item list-group-item-action d-flex align-items-center border-0 py-1"
                >
                  <i class="ti ti-git-branch me-2"></i>
                  <span class="text-truncate">{{ branch.startsWith(remote.name + '/') ? branch.substring(remote.name.length + 1) : branch }}</span>
                </a>
                <div v-if="!remote.branches?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
                  <span class="ms-4">No remote branches</span>
                </div>
              </div>
            </template>
            <div v-if="!props.currentRepoStats?.remotes?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
              <span class="ms-4">No remotes found</span>
            </div>
          </div>
        </li>
        <li v-if="props.activeTab">
          <h6 class="sidebar-header collapsible" @click="toggleSection('tags')">
            <i :class="['ti ti-chevron-right me-1 transition-icon', { 'rotate-90': !collapsed.tags }]"></i>
            <span>Tags</span>
          </h6>
          <div v-if="!collapsed.tags" class="list-group list-group-flush">
            <a v-for="tag in props.currentRepoStats?.tags || []"
               :key="tag"
               href="javascript:void(0);"
               class="list-group-item list-group-item-action d-flex align-items-center border-0"
            >
              <i class="ti ti-tag me-2"></i>
              <span class="text-truncate">{{ tag }}</span>
            </a>
            <div v-if="!props.currentRepoStats?.tags?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
              <span class="ms-4">No tags found</span>
            </div>
          </div>
        </li>
        <li v-if="props.activeTab">
          <h6 class="sidebar-header collapsible" @click="toggleSection('stashes')">
            <i :class="['ti ti-chevron-right me-1 transition-icon', { 'rotate-90': !collapsed.stashes }]"></i>
            <span>Stashes</span>
          </h6>
          <div v-if="!collapsed.stashes" class="list-group list-group-flush">
            <a v-for="(stash, index) in props.currentRepoStats?.stashes || []"
               :key="index"
               href="javascript:void(0);"
               class="list-group-item list-group-item-action d-flex align-items-center border-0"
            >
              <i class="ti ti-archive me-2"></i>
              <span class="text-truncate">{{ stash }}</span>
            </a>
            <div v-if="!props.currentRepoStats?.stashes?.length" class="list-group-item d-flex align-items-center border-0 text-muted small py-0">
              <span class="ms-4">No stashes found</span>
            </div>
          </div>
        </li>
        <li>
          <hr class="sidebar-divider mx-3 mt-1 mb-2" v-if="props.activeTab">
          <h6 class="sidebar-header">Recent Repositories</h6>
          <div class="list-group list-group-flush">
            <a v-for="repo in props.recentRepos"
               :key="repo.path"
               href="javascript:void(0);"
               class="list-group-item list-group-item-action d-flex align-items-center border-0"
               @click="emit('openRecentRepo', repo)"
               :title="repo.path"
            >
              <i class="ti ti-folder me-2"></i>
              <span class="text-truncate">{{ repo.name }}</span>
            </a>
          </div>
        </li>
      </ul>
    </div>
  </nav>
</template>

<style scoped>
.sidebar-management-section {
  background-color: var(--bs-tertiary-bg);
}

.sidebar-header {
  height: 32px;
  display: flex;
  align-items: center;
  margin-bottom: 0;
  padding: 0 1rem;
}

.sidebar-divider {
  border-top: 1px solid var(--bs-border-color);
  opacity: 0.5;
  margin: 0.5rem 1rem;
}

.sidebar-header.collapsible {
  cursor: pointer;
  user-select: none;
}

.sidebar-header.collapsible:hover {
  background-color: var(--bs-tertiary-bg);
}

.transition-icon {
  transition: transform 0.2s ease;
  width: 16px;
  display: inline-block;
}
.rotate-90 {
  transform: rotate(90deg);
}
.remote-branches .list-group-item {
  padding-left: 2rem !important;
}

.list-group-item {
  border: none !important;
  height: 32px; /* Fixed height to prevent shifts */
  padding-left: 1.5rem !important; /* Increased indentation */
}

.list-group-item.active {
  border: none !important;
  margin-top: 0 !important;
}
.list-group-item i {
  width: 20px;
  text-align: center;
}
.filter-input {
  font-size: 0.8125rem;
  padding-top: 2px;
  padding-bottom: 2px;
  height: 28px;
}
</style>