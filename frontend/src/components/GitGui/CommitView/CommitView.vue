<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import * as App from '../../../../wailsjs/go/backend/App';
import { backend } from "../../../../wailsjs/go/models";
import GitCommit = backend.GitCommit;
import CommitList from "./CommitList.vue";
import CommitDetailInfo from "./CommitDetailInfo.vue";
import CommitDetailChanges from "./CommitDetailChanges.vue";

const props = defineProps<{
  repoPath: string;
  refreshCounter?: number;
}>();

const commits = ref<GitCommit[]>([]);
const loading = ref(false);
const selectedCommit = ref<GitCommit | null>(null);
const activeDetailTab = ref<'info' | 'changes'>('info');
const commitChanges = ref<any[]>([]);
const loadingChanges = ref<boolean>(false);
const changesRef = ref<any>(null);
const isMaximized = ref(false);

const toggleMaximize = () => {
  isMaximized.value = !isMaximized.value;
};

const loadCommitChanges = async (hash: string) => {
  loadingChanges.value = true;
  if (changesRef.value) {
    changesRef.value.clearExpanded();
  }
  try {
    commitChanges.value = await App.GetCommitChanges(props.repoPath, hash);
  } catch (err) {
    console.error('Failed to load commit changes:', err);
    commitChanges.value = [];
  } finally {
    loadingChanges.value = false;
  }
};

const loadCommits = async () => {
  loading.value = true;
  try {
    commits.value = await App.GetCommitHistory(props.repoPath, 100);
    if (commits.value.length > 0 && !selectedCommit.value) {
      selectedCommit.value = commits.value[0];
      loadCommitChanges(selectedCommit.value.hash);
    }
  } catch (err) {
    console.error('Failed to load commits:', err);
    commits.value = [];
  } finally {
    loading.value = false;
  }
};

const selectCommit = (commit: GitCommit) => {
  selectedCommit.value = commit;
  loadCommitChanges(commit.hash);
};

onMounted(() => {
  loadCommits();
});

watch(() => props.repoPath, () => {
  selectedCommit.value = null;
  loadCommits();
});

watch(() => props.refreshCounter, () => {
  loadCommits();
});
</script>

<template>
  <div class="commit-view h-100 d-flex flex-column overflow-hidden">
    <div class="flex-grow-1 d-flex flex-column overflow-hidden">
      <!-- Top: Commit List -->
      <div :style="{ height: isMaximized ? '35%' : '65%' }">
        <CommitList 
          :commits="commits" 
          :loading="loading" 
          :selected-commit="selectedCommit" 
          @select="selectCommit"
        />
      </div>

      <!-- Bottom: Commit Details -->
      <div v-if="selectedCommit" class="commit-details bg-body-tertiary d-flex flex-column overflow-hidden" :style="{ height: isMaximized ? '65%' : '35%' }">
          <div class="commit-details-tabs d-flex border-bottom bg-body px-3 justify-content-between align-items-center">
              <div class="d-flex">
                  <button 
                      :class="['btn btn-sm px-3 py-2 border-0 rounded-0', { 'active-tab border-bottom border-primary border-2 text-primary': activeDetailTab === 'info' }]"
                      @click="activeDetailTab = 'info'"
                  >
                      INFO
                  </button>
                  <button 
                      :class="['btn btn-sm px-3 py-2 border-0 rounded-0', { 'active-tab border-bottom border-primary border-2 text-primary': activeDetailTab === 'changes' }]"
                      @click="activeDetailTab = 'changes'"
                  >
                      CHANGES ({{ commitChanges.length }})
                  </button>
              </div>
              <button class="btn btn-sm btn-link text-muted p-0 me-2" @click="toggleMaximize" :title="isMaximized ? 'Restore' : 'Maximize'">
                  <i :class="['ti', isMaximized ? 'ti-arrows-minimize' : 'ti-arrows-maximize']"></i>
              </button>
          </div>

          <div class="flex-grow-1 d-flex flex-column overflow-hidden">
              <CommitDetailInfo 
                v-if="activeDetailTab === 'info'" 
                :commit="selectedCommit" 
              />

              <CommitDetailChanges 
                v-else-if="activeDetailTab === 'changes'" 
                ref="changesRef"
                :repo-path="repoPath"
                :commit-hash="selectedCommit.hash"
                :changes="commitChanges"
                :loading="loadingChanges"
              />
          </div>
      </div>
      <div v-else class="commit-details bg-body-tertiary d-flex align-items-center justify-content-center text-muted" :style="{ height: isMaximized ? '65%' : '35%' }">
          Select a commit to view details
      </div>
    </div>
  </div>
</template>

<style scoped>
.active-tab {
    font-weight: bold;
}
</style>
