<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import type { GitStatusFile } from '@/types/git.types';
import * as App from '../../../wailsjs/go/backend/App';
import DiffViewer from './DiffViewer.vue';

const props = defineProps<{
  repoPath: string;
}>();

const unstagedFiles = ref<GitStatusFile[]>([]);
const stagedFiles = ref<GitStatusFile[]>([]);
const selectedFile = ref<GitStatusFile | null>(null);
const diffContent = ref<string>('');
const loading = ref(false);

const commitSubject = ref('');
const commitDescription = ref('');
const amend = ref(false);

watch(amend, async (newValue) => {
  if (newValue) {
    try {
      const lastCommits = await App.GetCommitHistory(props.repoPath, 1);
      if (lastCommits && lastCommits.length > 0) {
        const lastCommit = lastCommits[0];
        commitSubject.value = lastCommit.subject;
        commitDescription.value = lastCommit.body;
      }
    } catch (err) {
      console.error('Failed to load last commit for amend:', err);
    }
  }
});

const emit = defineEmits<{
  (e: 'refresh-stats'): void;
}>();

const commitChanges = async () => {
  if (!commitSubject.value.trim()) return;
  
  loading.value = true;
  try {
    await App.Commit(
      props.repoPath, 
      commitSubject.value, 
      commitDescription.value, 
      amend.value
    );
    
    // Clear inputs on success
    commitSubject.value = '';
    commitDescription.value = '';
    amend.value = false;
    
    // Refresh status
    await loadStatus();
    selectedFile.value = null;
    diffContent.value = '';
  } catch (err) {
    console.error('Failed to commit:', err);
    alert('Failed to commit: ' + err);
  } finally {
    loading.value = false;
  }
};

const loadStatus = async () => {
  loading.value = true;
  try {
    const status = await App.GetGitStatus(props.repoPath);
    // Sort all status items alphabetically by path before filtering
    status.sort((a, b) => a.path.localeCompare(b.path));

    unstagedFiles.value = status.filter(f => !f.is_staged);
    stagedFiles.value = status.filter(f => f.is_staged);
    emit('refresh-stats');
  } catch (err) {
    // console.error('Failed to load status:', err);
    // Clear lists on error to avoid showing stale data from previous repo
    unstagedFiles.value = [];
    stagedFiles.value = [];
  } finally {
    loading.value = false;
  }
};

const stageFile = async (file: GitStatusFile) => {
  try {
    await App.StageFile(props.repoPath, file.path);
    await loadStatus();
    if (selectedFile.value?.path === file.path) {
        selectFile(stagedFiles.value.find(f => f.path === file.path) || null);
    }
  } catch (err) {
    console.error('Failed to stage file:', err);
  }
};

const stageAll = async () => {
  try {
    await App.StageAll(props.repoPath);
    await loadStatus();
  } catch (err) {
    console.error('Failed to stage all files:', err);
  }
};

const unstageFile = async (file: GitStatusFile) => {
  try {
    await App.UnstageFile(props.repoPath, file.path);
    await loadStatus();
    if (selectedFile.value?.path === file.path) {
        selectFile(unstagedFiles.value.find(f => f.path === file.path) || null);
    }
  } catch (err) {
    console.error('Failed to unstage file:', err);
  }
};

const unstageAll = async () => {
  try {
    await App.UnstageAll(props.repoPath);
    await loadStatus();
  } catch (err) {
    console.error('Failed to unstage all files:', err);
  }
};

const selectFile = async (file: GitStatusFile | null) => {
  selectedFile.value = file;
  if (file) {
    try {
      diffContent.value = await App.GetFileDiff(props.repoPath, file.path, file.is_staged);
    } catch (err) {
      diffContent.value = 'Failed to load diff: ' + err;
    }
  } else {
    diffContent.value = '';
  }
};

onMounted(() => {
  loadStatus();
});

// We keep the watch for when repoPath changes but the component IS NOT recreated.
// However, with the :key in GitGui.vue, this watch might be redundant for tab switches,
// but it's still useful if repoPath changes for the SAME tab (e.g. repo renamed or similar).
watch(() => props.repoPath, () => {
  selectedFile.value = null;
  diffContent.value = '';
  loadStatus();
});

</script>

<template>
  <div class="changes-view h-100 d-flex overflow-hidden">
    <!-- Left Column: File Lists -->
    <div class="file-lists border-end d-flex flex-column" style="width: 375px;">
      <div class="section flex-grow-1 d-flex flex-column overflow-hidden">
        <div class="section-header px-3 py-2 bg-body-tertiary border-bottom d-flex justify-content-between align-items-center">
          <span class="fw-bold small">UNSTAGED ({{ unstagedFiles.length }})</span>
          <button 
            v-if="unstagedFiles.length > 0" 
            class="btn btn-sm btn-ghost p-0" 
            @click="stageAll" 
            title="Stage All"
          >
            <i class="ti ti-plus text-success fs-5"></i>
          </button>
        </div>
        <div class="file-list flex-grow-1 overflow-auto">
          <div 
            v-for="file in unstagedFiles" 
            :key="file.path"
            :class="['file-item px-3 py-1 d-flex align-items-center cursor-pointer', { active: selectedFile?.path === file.path && !selectedFile?.is_staged }]"
            @click="selectFile(file)"
          >
            <span :class="['status-square me-2', (file.status === '?' || file.status === 'A') ? 'A' : file.status]">
              {{ (file.status === '?' || file.status === 'A') ? '+' : (file.status === 'D' ? '-' : file.status) }}
            </span>
            <span class="file-path text-truncate flex-grow-1">{{ file.path }}</span>
            <button class="btn btn-sm btn-ghost p-0 ms-1 stage-btn" @click.stop="stageFile(file)" title="Stage">
              <i class="ti ti-plus text-success"></i>
            </button>
          </div>
          <div v-if="unstagedFiles.length === 0" class="p-3 text-center text-muted small">
            No unstaged changes
          </div>
        </div>
      </div>

      <div class="section flex-grow-1 d-flex flex-column overflow-hidden border-top">
        <div class="section-header px-3 py-2 bg-body-tertiary border-bottom d-flex justify-content-between align-items-center">
          <span class="fw-bold small">STAGED ({{ stagedFiles.length }})</span>
          <button 
            v-if="stagedFiles.length > 0" 
            class="btn btn-sm btn-ghost p-0" 
            @click="unstageAll" 
            title="Unstage All"
          >
            <i class="ti ti-minus text-danger fs-5"></i>
          </button>
        </div>
        <div class="file-list flex-grow-1 overflow-auto">
          <div 
            v-for="file in stagedFiles" 
            :key="file.path"
            :class="['file-item px-3 py-1 d-flex align-items-center cursor-pointer', { active: selectedFile?.path === file.path && selectedFile?.is_staged }]"
            @click="selectFile(file)"
          >
            <span :class="['status-square me-2', (file.status === '?' || file.status === 'A') ? 'A' : file.status]">
              {{ (file.status === '?' || file.status === 'A') ? '+' : (file.status === 'D' ? '-' : file.status) }}
            </span>
            <span class="file-path text-truncate flex-grow-1">{{ file.path }}</span>
            <button class="btn btn-sm btn-ghost p-0 ms-1 unstage-btn" @click.stop="unstageFile(file)" title="Unstage">
              <i class="ti ti-minus text-danger"></i>
            </button>
          </div>
          <div v-if="stagedFiles.length === 0" class="p-3 text-center text-muted small">
            No staged changes
          </div>
        </div>
      </div>
    </div>

    <!-- Right Column: Diff and Commit -->
    <div class="diff-commit-area flex-grow-1 d-flex flex-column overflow-hidden">
      <div class="diff-view flex-grow-1 overflow-auto bg-body font-monospace p-0">
        <div v-if="selectedFile" class="h-100 d-flex flex-column">
          <h6 class="section-header border-bottom px-3 mb-0 bg-body-tertiary d-flex align-items-center">{{ selectedFile.path }}</h6>
          <div class="flex-grow-1 overflow-auto p-3">
            <DiffViewer :diff="diffContent" />
          </div>
        </div>
        <div v-else class="h-100 d-flex align-items-center justify-content-center text-muted">
          Select a file to view changes
        </div>
      </div>

      <div class="commit-section border-top p-3 bg-body-tertiary">
        <div class="mb-2">
          <input 
            v-model="commitSubject" 
            type="text" 
            class="form-control form-control-sm bg-body" 
            placeholder="Commit subject"
          />
        </div>
        <div class="mb-2">
          <textarea 
            v-model="commitDescription" 
            class="form-control form-control-sm bg-body" 
            rows="3" 
            placeholder="Description (optional)"
          ></textarea>
        </div>
        <div class="d-flex align-items-center justify-content-between">
          <div class="form-check">
            <input v-model="amend" class="form-check-input" type="checkbox" id="amendCheck">
            <label class="form-check-label small" for="amendCheck">
              Amend
            </label>
          </div>
          <button 
            class="btn btn-primary btn-sm px-4" 
            :disabled="!commitSubject.trim() || (stagedFiles.length === 0 && !amend)"
            @click="commitChanges"
          >
            Commit
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.section-header {
  height: 38px;
  font-size: 0.875rem;
}
.file-item {
  font-size: 0.875rem;
}
.file-item:hover {
  background-color: var(--bs-tertiary-bg);
}
.file-item.active {
  background-color: var(--bs-primary-bg-subtle);
  color: var(--bs-primary-text-emphasis);
}
.status-square {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1.1rem;
  height: 1.1rem;
  min-width: 1.1rem;
  border-radius: 3px;
  font-weight: bold;
  font-size: 0.7rem;
  color: #fff !important;
}

.M { background-color: #f0ad4e; }
.A { background-color: #5cb85c; }
.D { background-color: #d9534f; }
.R { background-color: #5bc0de; }
.C { background-color: #5bc0de; }
.U { background-color: #d9534f; }


.stage-btn, .unstage-btn {
    visibility: hidden;
}
.file-item:hover .stage-btn, .file-item:hover .unstage-btn {
    visibility: visible;
}

.btn-ghost {
  border: 1px solid transparent;
  background: transparent;
}
.btn-ghost:hover {
  background-color: var(--bs-secondary-bg);
}

.commit-section .form-control {
  border-color: var(--bs-border-color);
}
</style>
