<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import type { GitCommit } from './types';

const props = defineProps<{
  commits: GitCommit[];
  selectedCommitHash: string | null;
}>();

const emit = defineEmits<{
  (e: 'select', hash: string): void;
  (e: 'addTag', hash: string): void;
  (e: 'createBranch', hash: string): void;
  (e: 'copyHash', hash: string): void;
}>();

const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  hash: ''
});

const showContextMenu = (event: MouseEvent, hash: string) => {
  if (!hash) return;
  
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    hash
  };
  
  // Also select the commit if it's not selected
  if (props.selectedCommitHash !== hash) {
    emit('select', hash);
  }
};

const closeContextMenu = () => {
  contextMenu.value.show = false;
};

const handleAddTag = () => {
  emit('addTag', contextMenu.value.hash);
  closeContextMenu();
};

const handleCreateBranch = () => {
  emit('createBranch', contextMenu.value.hash);
  closeContextMenu();
};

const handleCopyHash = () => {
  emit('copyHash', contextMenu.value.hash);
  closeContextMenu();
};

onMounted(() => {
  window.addEventListener('click', closeContextMenu);
});

onUnmounted(() => {
  window.removeEventListener('click', closeContextMenu);
});

const formatDate = (dateStr: string) => {
  if (!dateStr) return '';
  const timestamp = parseInt(dateStr);
  if (isNaN(timestamp)) return dateStr;
  return new Date(timestamp * 1000).toLocaleString();
};
</script>

<template>
  <div class="flex-grow-1 overflow-auto border-bottom" style="height: 55%;">
    <table class="table table-hover table-sm mb-0 align-middle table-commits">
      <thead class="sticky-top bg-light-subtle">
        <tr>
          <th class="ps-3">Message</th>
          <th style="width: 150px">Author</th>
          <th style="width: 150px">Date</th>
          <th style="width: 100px">Hash</th>
        </tr>
      </thead>
      <tbody>
        <tr 
          v-for="(commit, idx) in commits" 
          :key="commit.hash + '-' + idx"
          :data-commit-hash="commit.hash"
          :class="{ 'table-primary': commit.hash && selectedCommitHash === commit.hash }"
          class="cursor-pointer"
          @click="commit.hash && emit('select', commit.hash)"
          @contextmenu.prevent="commit.hash && showContextMenu($event, commit.hash)"
        >
          <td class="ps-3">
            <div v-if="commit.hash" class="d-flex align-items-center flex-wrap gap-1">
              <template v-for="branch in commit.branches" :key="branch">
                <span v-if="branch.startsWith('origin/')" class="badge bg-secondary-subtle text-secondary-emphasis border border-secondary-subtle small py-0 px-1" :title="branch">
                  <i class="ti ti-brand-github" style="font-size: 0.7rem;"></i>
                  <span v-if="branch === 'origin/HEAD'" class="ms-1" style="font-size: 0.6rem;">HEAD</span>
                </span>
                <span v-else class="badge bg-info-subtle text-info-emphasis border border-info-subtle small py-0 px-1">
                  <i class="ti ti-git-branch me-1" style="font-size: 0.7rem;"></i>{{ branch }}
                </span>
              </template>
              <span v-for="tag in commit.tags" :key="tag" class="badge bg-warning-subtle text-warning-emphasis border border-warning-subtle small py-0 px-1">
                <i class="ti ti-tag me-1" style="font-size: 0.7rem;"></i>{{ tag }}
              </span>
              <span class="text-truncate fw-medium">{{ commit.message }}</span>
            </div>
          </td>
          <td class="text-nowrap small text-muted">{{ commit.author }}</td>
          <td class="text-nowrap small text-muted">{{ formatDate(commit.date) }}</td>
          <td><code v-if="commit.hash" class="smaller text-muted" :title="commit.hash">{{ commit.hash.substring(0, 7) }}</code></td>
        </tr>
      </tbody>
    </table>

    <!-- Context Menu -->
    <div 
      v-if="contextMenu.show" 
      class="dropdown-menu show shadow-sm border position-fixed" 
      :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
    >
      <button class="dropdown-item d-flex align-items-center" @click="handleAddTag">
        <i class="ti ti-tag me-2"></i> Add Tag...
      </button>
      <button class="dropdown-item d-flex align-items-center" @click="handleCreateBranch">
        <i class="ti ti-git-branch me-2"></i> Create branch...
      </button>
      <div class="dropdown-divider"></div>
      <button class="dropdown-item d-flex align-items-center" @click="handleCopyHash">
        <i class="ti ti-copy me-2"></i> Copy SHA
      </button>
    </div>
  </div>
</template>

<style scoped>
.table-commits {
  border-collapse: separate;
  border-spacing: 0;
}
.table-commits tr {
  height: 32px;
}
.cursor-pointer {
  cursor: pointer;
}
.smaller {
  font-size: 0.75rem;
}
</style>
