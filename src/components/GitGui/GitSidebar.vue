<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import type { RepoTab } from './types';

defineProps<{
  tab: RepoTab;
}>();

const sections = ref({
  localBranches: true,
  remoteBranches: true,
  tags: false,
  stashes: false,
  remotes: false
});

const toggleSection = (section: keyof typeof sections.value) => {
  sections.value[section] = !sections.value[section];
};

const emit = defineEmits<{
  (e: 'setSelection', selection: 'local-changes' | 'all-commits'): void;
  (e: 'refresh'): void;
  (e: 'jumpToCommit', branchName: string): void;
  (e: 'switchBranch', branchName: string): void;
  (e: 'checkoutRemote', remoteBranch: string): void;
  (e: 'checkoutTag', tagName: string): void;
  (e: 'createBranch', fromBranch: string): void;
  (e: 'createBranchFromTag', tagName: string): void;
  (e: 'deleteBranch', branchName: string): void;
  (e: 'deleteStash', index: number): void;
  (e: 'popStash', index: number): void;
}>();

const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  type: 'branch' as 'branch' | 'stash' | 'tag',
  branchName: '',
  tagName: '',
  isCurrent: false,
  stashIndex: -1
});

const showContextMenu = (event: MouseEvent, branchName: string, isCurrent: boolean) => {
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    type: 'branch',
    branchName,
    tagName: '',
    isCurrent,
    stashIndex: -1
  };
};

const showTagContextMenu = (event: MouseEvent, tagName: string) => {
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    type: 'tag',
    branchName: '',
    tagName,
    isCurrent: false,
    stashIndex: -1
  };
};

const showStashContextMenu = (event: MouseEvent, index: number) => {
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    type: 'stash',
    branchName: '',
    tagName: '',
    isCurrent: false,
    stashIndex: index
  };
};

const closeContextMenu = () => {
  contextMenu.value.show = false;
};

const handleCreateBranch = () => {
  emit('createBranch', contextMenu.value.branchName);
  closeContextMenu();
};

const handleDeleteBranch = () => {
  emit('deleteBranch', contextMenu.value.branchName);
  closeContextMenu();
};

const handleDeleteStash = () => {
  emit('deleteStash', contextMenu.value.stashIndex);
  closeContextMenu();
};

const handlePopStash = () => {
  emit('popStash', contextMenu.value.stashIndex);
  closeContextMenu();
};

const handleCheckoutTag = () => {
  emit('checkoutTag', contextMenu.value.tagName);
  closeContextMenu();
};

const handleCreateBranchFromTag = () => {
  emit('createBranchFromTag', contextMenu.value.tagName);
  closeContextMenu();
};

onMounted(() => {
  window.addEventListener('click', closeContextMenu);
});

onUnmounted(() => {
  window.removeEventListener('click', closeContextMenu);
});
</script>

<template>
  <div class="col-md-2 border-end bg-light-subtle h-100 overflow-auto">
    <div class="list-group list-group-flush border-bottom mb-3">
      <button 
        class="list-group-item list-group-item-action d-flex justify-content-between align-items-center py-3"
        :class="{ active: tab.sidebarSelection === 'local-changes' }"
        @click="emit('setSelection', 'local-changes')"
      >
        <div class="d-flex align-items-center">
          <i class="ti ti-file-diff me-2 fs-5"></i>
          <span>Local changes</span>
        </div>
        <span class="badge rounded-pill bg-primary" v-if="tab.statusFiles.length > 0">
          {{ tab.statusFiles.length }}
        </span>
      </button>
      <button 
        class="list-group-item list-group-item-action d-flex align-items-center py-3"
        :class="{ active: tab.sidebarSelection === 'all-commits' }"
        @click="emit('setSelection', 'all-commits')"
      >
        <i class="ti ti-history me-2 fs-5"></i>
        <span>All commits</span>
      </button>
    </div>

    <div class="px-3 pb-3">
      <div class="mb-4">
        <h6 class="mb-3 d-flex align-items-center justify-content-between cursor-pointer" @click="toggleSection('localBranches')">
          <span class="d-flex align-items-center">
            <i class="ti me-2" :class="sections.localBranches ? 'ti-chevron-down' : 'ti-chevron-right'"></i>
            <i class="ti ti-git-fork me-2"></i> Local Branches
          </span>
          <button 
            class="btn btn-sm btn-ghost p-1 lh-1" 
            @click.stop="emit('refresh')" 
            title="Refresh branches"
            :disabled="tab.loading"
          >
            <i class="ti ti-refresh" :class="{ 'spinner-border spinner-border-sm border-0': tab.loading }"></i>
          </button>
        </h6>
        <div class="list-group list-group-flush" v-if="sections.localBranches">
          <div 
            v-for="branch in tab.branches" 
            :key="branch.name"
            class="list-group-item list-group-item-action border-0 px-2 py-1 rounded mb-1 d-flex align-items-center cursor-pointer"
            :class="{ 'bg-primary-subtle fw-bold': branch.is_current }"
            @click="emit('jumpToCommit', branch.name)"
            @dblclick="emit('switchBranch', branch.name)"
            @contextmenu.prevent="showContextMenu($event, branch.name, branch.is_current)"
            :title="branch.is_current ? 'Current branch' : 'Click to jump, Double-click to switch'"
          >
            <i class="ti ti-git-branch me-2 small"></i>
            <span class="text-truncate">{{ branch.name }}</span>
          </div>
        </div>
      </div>

      <!-- Context Menu -->
      <div 
        v-if="contextMenu.show" 
        class="dropdown-menu show shadow-sm border position-fixed" 
        :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      >
        <template v-if="contextMenu.type === 'branch'">
          <button class="dropdown-item d-flex align-items-center" @click="handleCreateBranch">
            <i class="ti ti-git-branch me-2"></i> Create branch...
          </button>
          <button class="dropdown-item d-flex align-items-center text-danger" @click="handleDeleteBranch" :disabled="contextMenu.isCurrent">
            <i class="ti ti-trash me-2"></i> Delete branch...
          </button>
        </template>
        <template v-else-if="contextMenu.type === 'stash'">
          <button class="dropdown-item d-flex align-items-center" @click="handlePopStash">
            <i class="ti ti-arrow-up me-2"></i> Restore stash (pop)
          </button>
          <button class="dropdown-item d-flex align-items-center text-danger" @click="handleDeleteStash">
            <i class="ti ti-trash me-2"></i> Delete stash...
          </button>
        </template>
        <template v-else-if="contextMenu.type === 'tag'">
          <button class="dropdown-item d-flex align-items-center" @click="handleCheckoutTag">
            <i class="ti ti-git-check me-2"></i> Checkout tag
          </button>
          <button class="dropdown-item d-flex align-items-center" @click="handleCreateBranchFromTag">
            <i class="ti ti-git-branch me-2"></i> Create branch from tag...
          </button>
        </template>
      </div>

      <div class="mb-4">
        <h6 class="mb-3 d-flex align-items-center justify-content-between cursor-pointer" @click="toggleSection('remoteBranches')">
          <span class="d-flex align-items-center">
            <i class="ti me-2" :class="sections.remoteBranches ? 'ti-chevron-down' : 'ti-chevron-right'"></i>
            <i class="ti ti-cloud-download me-2"></i> Remote Branches
          </span>
          <button 
            class="btn btn-sm btn-ghost p-1 lh-1" 
            @click.stop="emit('refresh')" 
            title="Refresh branches"
            :disabled="tab.loading"
          >
            <i class="ti ti-refresh" :class="{ 'spinner-border spinner-border-sm border-0': tab.loading }"></i>
          </button>
        </h6>
        <div class="list-group list-group-flush" v-if="sections.remoteBranches">
          <div 
            v-for="rb in tab.remoteBranches" 
            :key="rb"
            class="list-group-item list-group-item-action border-0 px-2 py-1 rounded mb-1 d-flex align-items-center cursor-pointer"
            @dblclick="emit('checkoutRemote', rb)"
            title="Double-click to checkout"
          >
            <i class="ti ti-git-branch me-2 small text-muted"></i>
            <span class="text-truncate small">{{ rb }}</span>
          </div>
          <div v-if="tab.remoteBranches.length === 0" class="px-4 py-1 text-muted smaller">No remote branches</div>
        </div>
      </div>

      <div class="mb-4">
        <h6 class="mb-3 d-flex align-items-center justify-content-between cursor-pointer" @click="toggleSection('tags')">
          <span class="d-flex align-items-center">
            <i class="ti me-2" :class="sections.tags ? 'ti-chevron-down' : 'ti-chevron-right'"></i>
            <i class="ti ti-tag me-2"></i> Tags
          </span>
        </h6>
        <div class="list-group list-group-flush" v-if="sections.tags">
          <div 
            v-for="tag in tab.tags" 
            :key="tag"
            class="list-group-item list-group-item-action border-0 px-2 py-1 rounded mb-1 d-flex align-items-center cursor-pointer"
            @click="emit('jumpToCommit', tag)"
            @contextmenu.prevent="showTagContextMenu($event, tag)"
            title="Click to jump to tag, Right-click for options"
          >
            <i class="ti ti-tag me-2 small text-muted"></i>
            <span class="text-truncate small">{{ tag }}</span>
          </div>
          <div v-if="tab.tags.length === 0" class="px-4 py-1 text-muted smaller">No tags</div>
        </div>
      </div>

      <div class="mb-4">
        <h6 class="mb-3 d-flex align-items-center justify-content-between cursor-pointer" @click="toggleSection('stashes')">
          <span class="d-flex align-items-center">
            <i class="ti me-2" :class="sections.stashes ? 'ti-chevron-down' : 'ti-chevron-right'"></i>
            <i class="ti ti-archive me-2"></i> Stashes
          </span>
        </h6>
        <div class="list-group list-group-flush" v-if="sections.stashes">
          <div 
            v-for="stash in tab.stashes" 
            :key="stash.index"
            class="list-group-item list-group-item-action border-0 px-2 py-1 rounded mb-1 d-flex flex-column cursor-pointer"
            @contextmenu.prevent="showStashContextMenu($event, stash.index)"
          >
            <div class="d-flex align-items-center">
              <i class="ti ti-archive me-2 small text-muted"></i>
              <span class="text-truncate small fw-medium">stash@{{ stash.index }}</span>
            </div>
            <div class="text-muted smaller text-truncate ps-4" :title="stash.message">{{ stash.message }}</div>
          </div>
          <div v-if="tab.stashes.length === 0" class="px-4 py-1 text-muted smaller">No stashes</div>
        </div>
      </div>

      <div class="mb-4">
        <h6 class="mb-3 d-flex align-items-center justify-content-between cursor-pointer" @click="toggleSection('remotes')">
          <span class="d-flex align-items-center">
            <i class="ti me-2" :class="sections.remotes ? 'ti-chevron-down' : 'ti-chevron-right'"></i>
            <i class="ti ti-world me-2"></i> Remotes
          </span>
        </h6>
        <div class="list-group list-group-flush" v-if="sections.remotes">
          <div 
            v-for="remote in tab.remotes" 
            :key="remote.name"
            class="list-group-item list-group-item-action border-0 px-2 py-1 rounded mb-1 cursor-pointer"
          >
            <div class="fw-bold small">{{ remote.name }}</div>
            <div class="text-muted smaller text-truncate" :title="remote.url">{{ remote.url }}</div>
          </div>
          <div v-if="tab.remotes.length === 0" class="px-4 py-1 text-muted smaller">No remotes</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.btn-ghost {
  background: transparent;
  border: none;
  color: inherit;
}
.btn-ghost:hover {
  background: rgba(0, 0, 0, 0.05);
}
.list-group-item-action {
  transition: background-color 0.15s ease-in-out;
}
.dropdown-item:disabled {
  color: var(--bs-secondary-color, #6c757d) !important;
  background-color: transparent !important;
  cursor: default;
  opacity: 0.65;
}
.cursor-pointer {
  cursor: pointer;
}
.smaller {
  font-size: 0.75rem;
}
</style>
