<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import * as App from '../../../wailsjs/go/backend/App';
import { backend } from "../../../wailsjs/go/models";
import GitCommit = backend.GitCommit;
import dayjs from "dayjs";
import GitGraph from "./GitGraph.vue";

const props = defineProps<{
  repoPath: string;
  refreshCounter?: number;
}>();

const commits = ref<GitCommit[]>([]);
const loading = ref(false);
const selectedCommit = ref<GitCommit | null>(null);
const detailsHeight = ref(parseInt(localStorage.getItem('commit-details-height') || '400'));
const isResizing = ref(false);
let startY = 0;
let startHeight = 0;

const startResizing = (e: MouseEvent) => {
  isResizing.value = true;
  startY = e.pageY;
  startHeight = detailsHeight.value;
  document.addEventListener('mousemove', handleMouseMove);
  document.addEventListener('mouseup', stopResizing);
  document.body.style.cursor = 'row-resize';
  document.body.classList.add('resizing');
};

const handleMouseMove = (e: MouseEvent) => {
  if (!isResizing.value) return;
  
  const delta = startY - e.pageY;
  const newHeight = startHeight + delta;
  
  const container = document.querySelector('.commit-view');
  const maxHeight = container ? container.clientHeight - 100 : 800;
  
  if (newHeight > 100 && newHeight < maxHeight) {
    detailsHeight.value = newHeight;
  }
};

const stopResizing = () => {
  isResizing.value = false;
  document.removeEventListener('mousemove', handleMouseMove);
  document.removeEventListener('mouseup', stopResizing);
  document.body.style.cursor = '';
  document.body.classList.remove('resizing');
  localStorage.setItem('commit-details-height', detailsHeight.value.toString());
};

const loadCommits = async () => {
  loading.value = true;
  try {
    commits.value = await App.GetCommitHistory(props.repoPath, 100);
    if (commits.value.length > 0 && !selectedCommit.value) {
      selectedCommit.value = commits.value[0];
    }
  } catch (err) {
    console.error('Failed to load commits:', err);
    commits.value = [];
  } finally {
    loading.value = false;
  }
};

const formatDate = (date: any) => {
  return dayjs(date).format('D MMM YYYY HH:mm');
};

const selectCommit = (commit: GitCommit) => {
  selectedCommit.value = commit;
};

const getRefClass = (ref: string) => {
  if (ref === 'HEAD') return 'bg-info-subtle text-info-emphasis border border-info';
  if (ref.includes('/') || ref === 'main' || ref === 'master' || ref === 'develop') {
      return 'bg-success-subtle text-success-emphasis border border-success';
  }
  return 'bg-warning-subtle text-warning-emphasis border border-warning';
};

const getRefIcon = (ref: string) => {
  if (ref.includes('/') || ref === 'main' || ref === 'master' || ref === 'develop') return 'ti ti-git-branch';
  if (ref.startsWith('v')) return 'ti ti-tag';
  return 'ti ti-git-commit';
}

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
    <div class="flex-grow-1 d-flex overflow-hidden">
      <!-- Left: Commit List -->
      <div class="commit-list-container border-end d-flex flex-column" style="width: 100%;">
        <div class="commit-list-header px-3 py-2 bg-body-tertiary border-bottom d-flex align-items-center">
          <span class="fw-bold small">COMMITS</span>
        </div>
        
        <div class="commit-list flex-grow-1 overflow-auto bg-body position-relative scroll-container">
          <!-- Graph Layer -->
          <div class="graph-layer position-absolute start-0 top-0">
            <GitGraph :commits="commits" :row-height="38" />
          </div>

          <table class="table table-hover table-sm mb-0 position-relative" style="background: transparent;z-index:0">
            <thead>
              <tr class="small bg-body-tertiary">
                <th class="ps-3 border-bottom-0 fw-normal text-muted" style="width: 120px;">Graph</th>
                <th class="border-bottom-0 fw-normal text-muted" style="width: 80px;">Hash</th>
                <th class="border-bottom-0 fw-normal text-muted">Subject</th>
                <th class="border-bottom-0 fw-normal text-muted" style="width: 150px;">Author</th>
                <th class="pe-3 border-bottom-0 fw-normal text-muted text-end" style="width: 150px;">Date</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="commit in commits" 
                  :key="commit.hash" 
                  :class="['cursor-pointer align-middle', { 'table-active': selectedCommit?.hash === commit.hash }]"
                  @click="selectCommit(commit)"
              >
                <td class="ps-3" style="border-bottom:0">
                   <!-- Graph spacing -->
                </td>
                <td><code class="small">{{ commit.hash.substring(0, 7) }}</code></td>
                <td class="text-truncate" style="max-width: 0;">
                  <span class="fw-medium">{{ commit.subject }}</span>
                  <span v-if="commit.refs && commit.refs.length" class="ms-2">
                    <span v-for="ref in commit.refs.filter(r => r !== 'HEAD' && !r.endsWith('/HEAD'))" :key="ref" 
                          :class="['badge me-1 border-1', getRefClass(ref)]" style="border-radius:5px;">
                      <i :class="['ti', getRefIcon(ref),'pe-1']"></i>{{ ref }}
                    </span>
                  </span>
                </td>
                <td class="text-truncate" style="max-width: 0;">{{ commit.authorName }}</td>
                <td class="pe-3 text-end text-muted small">{{ formatDate(commit.date) }}</td>
              </tr>
              <tr v-if="!loading && commits.length === 0">
                <td colspan="5" class="text-center py-5 text-muted">
                  No commits found in this repository.
                </td>
              </tr>
              <tr v-if="loading">
                <td colspan="5" class="text-center py-5 text-muted">
                  <div class="spinner-border spinner-border-sm me-2" role="status"></div>
                  Loading history...
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    
    <!-- Resizer -->
    <div v-if="selectedCommit" class="resizer" @mousedown="startResizing"></div>

    <!-- Bottom: Commit Details -->
    <div v-if="selectedCommit" class="commit-details border-top p-3 bg-body-tertiary overflow-auto" :style="{ height: detailsHeight + 'px' }">
        <div class="d-flex justify-content-between align-items-start mb-2">
            <div>
                <h5 class="mb-1">{{ selectedCommit.subject }}</h5>
                <div class="text-muted small">
                    <strong>{{ selectedCommit.authorName }}</strong> &lt;{{ selectedCommit.authorEmail }}&gt;
                    <span class="mx-2">•</span>
                    {{ formatDate(selectedCommit.date) }}
                </div>
            </div>
            <code class="bg-body px-2 py-1 rounded border">{{ selectedCommit.hash }}</code>
        </div>
        <div v-if="selectedCommit.body" class="mt-3 p-3 bg-body rounded border white-space-pre">{{ selectedCommit.body }}</div>
        <div class="mt-2 small text-muted">
            Parents: <code v-for="p in selectedCommit.parentHashes" :key="p" class="me-2">{{ p.substring(0, 7) }}</code>
        </div>
    </div>
  </div>
</template>

<style scoped>
.commit-list-header {
  height: 38px;
}

.table > :not(caption) > * > * {
  padding: 0.5rem 0.5rem;
  border-bottom-width: 1px;
}

.table tr {
    height: 38px;
}

.table td {
    padding-top: 0 !important;
    padding-bottom: 0 !important;
    height: 38px;
    vertical-align: middle;
    background: transparent !important;
}

.graph-layer {
    pointer-events: none;
    margin-top: 38px; /* Height of the table header */
}

.table thead th {
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    position: sticky;
    top: 0;
    background-color: var(--bs-tertiary-bg);
    z-index: 3;
}

.badge {
    font-size: 0.7rem;
    font-weight: 600;
    padding: 0.2rem 0.5rem;
}

.white-space-pre {
    white-space: pre-wrap;
}

.graph-cell {
    width: 20px;
    height: 20px;
}

.resizer {
    height: 6px;
    background-color: var(--bs-border-color);
    cursor: row-resize;
    transition: background-color 0.2s;
    z-index: 10;
    flex-shrink: 0;
    position: relative;
}

.resizer::after {
    content: "";
    position: absolute;
    top: -5px;
    bottom: -5px;
    left: 0;
    right: 0;
}

.resizer:hover {
    background-color: var(--bs-primary);
}
</style>
