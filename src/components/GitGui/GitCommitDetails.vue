<script setup lang="ts">
import { ref, watch } from 'vue';
import { invoke } from '@tauri-apps/api/core';
import type { GitCommit, GitCommitFile } from './types';
import GitCommitFileDiff from './GitCommitFileDiff.vue';

const props = defineProps<{
  commit: GitCommit | null;
  repoPath?: string;
  avatars: Map<string, string>;
}>();

const emit = defineEmits<{
  (e: 'copyHash', hash: string): void;
  (e: 'selectCommit', hash: string): void;
}>();

const activeTab = ref<'commit' | 'changes'>('commit');
const commitFiles = ref<GitCommitFile[]>([]);
const filesLoading = ref(false);
const expandedFile = ref<string | null>(null);

watch(() => props.commit, async (newCommit) => {
  activeTab.value = 'commit';
  commitFiles.value = [];
  expandedFile.value = null;
  if (newCommit && props.repoPath) {
    fetchCommitFiles();
  }
}, { immediate: true });

const toggleFile = (path: string) => {
  if (expandedFile.value === path) {
    expandedFile.value = null;
  } else {
    expandedFile.value = path;
  }
};

const fetchCommitFiles = async () => {
  if (!props.commit || !props.repoPath) return;
  filesLoading.value = true;
  try {
    commitFiles.value = await invoke('get_commit_files', { path: props.repoPath, hash: props.commit.hash });
  } catch (err) {
    console.error('Failed to fetch commit files:', err);
  } finally {
    filesLoading.value = false;
  }
};

const getStatusBadgeClass = (status: string) => {
  const s = status.trim().toUpperCase();
  if (s === 'M') return 'bg-warning text-dark';
  if (s === 'A') return 'bg-success';
  if (s === 'D') return 'bg-danger';
  if (s === 'R') return 'bg-info text-dark';
  return 'bg-secondary';
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '';
  const timestamp = parseInt(dateStr);
  if (isNaN(timestamp)) return dateStr;
  return new Date(timestamp * 1000).toLocaleString();
};
</script>

<template>
  <div class="p-0 overflow-hidden d-flex flex-column" style="height: 45%;">
    <div v-if="commit" class="h-100 d-flex flex-column">
      <!-- Tabs for Commit Details -->
      <div class="d-flex border-bottom bg-light-subtle px-3 py-2 gap-3">
        <div 
          class="small pb-1 cursor-pointer" 
          :class="activeTab === 'commit' ? 'fw-bold border-bottom border-primary border-2' : 'text-muted'"
          @click="activeTab = 'commit'"
          style="cursor: pointer;"
        >
          Commit
        </div>
        <div 
          class="small pb-1 cursor-pointer" 
          :class="activeTab === 'changes' ? 'fw-bold border-bottom border-primary border-2' : 'text-muted'"
          @click="activeTab = 'changes'"
          style="cursor: pointer;"
        >
          Changes <span v-if="commitFiles.length" class="badge rounded-pill bg-secondary smaller ms-1">{{ commitFiles.length }}</span>
        </div>
      </div>

      <div v-if="activeTab === 'commit'" class="flex-grow-1 overflow-auto p-4">
        <div class="d-flex gap-4 mb-4">
          <div class="author-avatar rounded bg-secondary-subtle d-flex align-items-center justify-content-center overflow-hidden" style="width: 64px; height: 64px;">
            <img v-if="avatars.get(commit.author_email || '')" 
                 :src="avatars.get(commit.author_email || '')" 
                 :alt="commit.author" 
                 class="w-100 h-100 object-fit-cover" />
            <i v-else class="ti ti-user fs-1 text-secondary"></i>
          </div>
          <div class="flex-grow-1">
            <div class="text-muted smaller mb-1">AUTHOR</div>
            <div class="d-flex align-items-baseline gap-2">
              <h5 class="mb-0">{{ commit.author }}</h5>
              <span class="text-muted small">{{ commit.author_email }}</span>
            </div>
            <div class="text-muted small">{{ formatDate(commit.date) }}</div>
          </div>
        </div>

        <div class="row g-3 mb-4">
          <div class="col-md-12" v-if="commit.branches.length">
            <div class="text-muted smaller mb-1">REFS</div>
            <div class="d-flex flex-wrap gap-1">
              <span v-for="branch in commit.branches" :key="branch" class="badge bg-secondary-subtle text-secondary-emphasis border border-secondary-subtle">
                <i class="ti ti-git-branch me-1"></i>{{ branch }}
              </span>
            </div>
          </div>
          <div class="col-md-12">
            <div class="text-muted smaller mb-1">SHA</div>
            <div class="d-flex align-items-center gap-2">
               <code class="small bg-body-tertiary p-1 rounded border">{{ commit.hash }}</code>
               <button class="btn btn-sm btn-ghost p-1" @click="emit('copyHash', commit.hash)" title="Copy Hash" style="cursor: pointer;">
                 <i class="ti ti-copy text-muted"></i>
               </button>
               <i class="ti ti-brand-github text-muted cursor-pointer ms-2" title="View on Remote" style="cursor: pointer;"></i>
            </div>
          </div>
          <div class="col-md-12" v-if="commit.parents.length">
            <div class="text-muted smaller mb-1">PARENTS</div>
            <div class="d-flex flex-wrap gap-2">
              <button v-for="parent in commit.parents" :key="parent" class="btn btn-xs btn-outline-secondary font-monospace" @click="emit('selectCommit', parent)">
                {{ parent.substring(0, 7) }}
              </button>
            </div>
          </div>
        </div>

        <div class="commit-body-full p-3 rounded bg-light-subtle border">
          <h6 class="mb-2">{{ commit.message }}</h6>
          <div class="whitespace-pre-wrap small text-body-secondary">
            {{ commit.body || 'No additional description.' }}
          </div>
        </div>
      </div>

      <div v-else-if="activeTab === 'changes'" class="flex-grow-1 overflow-auto">
        <div v-if="filesLoading" class="h-100 d-flex align-items-center justify-content-center">
          <div class="spinner-border spinner-border-sm text-primary" role="status"></div>
        </div>
        <div v-else class="list-group list-group-flush">
          <div v-for="file in commitFiles" :key="file.path" class="list-group-item list-group-item-action py-2 px-3">
            <div class="d-flex align-items-center cursor-pointer" @click="toggleFile(file.path)" style="cursor: pointer;">
              <span :class="['badge me-3 font-monospace smaller', getStatusBadgeClass(file.status)]" style="width: 24px">
                {{ file.status.trim() }}
              </span>
              <div class="text-truncate small flex-grow-1" :title="file.path">
                {{ file.path }}
              </div>
              <i class="ti text-muted smaller ms-2" :class="expandedFile === file.path ? 'ti-chevron-down' : 'ti-chevron-right'"></i>
            </div>
            
            <GitCommitFileDiff 
              v-if="expandedFile === file.path"
              :path="file.path"
              :hash="commit.hash"
              :repo-path="repoPath!"
            />
          </div>
          <div v-if="commitFiles.length === 0" class="p-4 text-center text-muted small">
            No files found in this commit.
          </div>
        </div>
      </div>
    </div>
    <div v-else class="h-100 d-flex align-items-center justify-content-center text-muted">
      Select a commit to view details
    </div>
  </div>
</template>

<style scoped>
.author-avatar {
  border: 1px solid var(--bs-border-color);
}
.cursor-not-allowed {
  cursor: not-allowed;
  opacity: 0.5;
}
.smaller {
  font-size: 0.75rem;
}
.btn-ghost {
  background: transparent;
  border: none;
  color: inherit;
}
.btn-ghost:hover {
  background: rgba(0, 0, 0, 0.05);
}
.whitespace-pre-wrap {
  white-space: pre-wrap;
}
.font-monospace {
  font-family: var(--bs-font-monospace);
}
</style>
