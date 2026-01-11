<script setup lang="ts">
import { backend } from "../../../../wailsjs/go/models";
import GitCommit = backend.GitCommit;
import dayjs from "dayjs";
import GitGraph from "./GitGraph.vue";

defineProps<{
  commits: GitCommit[];
  loading: boolean;
  selectedCommit: GitCommit | null;
}>();

const emit = defineEmits<{
  (e: 'select', commit: GitCommit): void;
}>();

const formatDate = (date: any) => {
  return dayjs(date).format('D MMM YYYY HH:mm');
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
</script>

<template>
  <div class="commit-list-container border-bottom d-flex flex-column h-100">
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
              @click="emit('select', commit)"
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

.table-active {
    background-color: var(--bs-table-active-bg);
}
</style>
