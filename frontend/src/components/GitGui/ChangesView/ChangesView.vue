<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import type { GitStatusFile } from '@/types/git.types';
import * as App from '../../../../wailsjs/go/backend/App';
import FileStatusList from './FileStatusList.vue';
import CommitSection from './CommitSection.vue';
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
const commitSectionRef = ref<any>(null);

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

const commitChanges = async (data: { subject: string; description: string; amend: boolean }) => {
  if (!data.subject.trim()) return;
  
  loading.value = true;
  try {
    await App.Commit(
      props.repoPath, 
      data.subject, 
      data.description, 
      data.amend
    );
    
    // Clear inputs on success
    if (commitSectionRef.value) {
        commitSectionRef.value.clearInputs();
    }
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
      <FileStatusList 
        title="UNSTAGED"
        :files="unstagedFiles"
        :selected-file="selectedFile"
        :is-staged-list="false"
        @select="selectFile"
        @action="stageFile"
        @action-all="stageAll"
      />

      <FileStatusList 
        title="STAGED"
        :files="stagedFiles"
        :selected-file="selectedFile"
        :is-staged-list="true"
        @select="selectFile"
        @action="unstageFile"
        @action-all="unstageAll"
      />
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

      <CommitSection 
        ref="commitSectionRef"
        :staged-count="stagedFiles.length"
        v-model:is-amend="amend"
        :initial-subject="commitSubject"
        :initial-description="commitDescription"
        @commit="commitChanges"
      />
    </div>
  </div>
</template>

<style scoped>
.section-header {
  height: 38px;
  font-size: 0.875rem;
}
</style>
