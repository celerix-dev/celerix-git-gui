<script setup lang="ts">
import type { GitStatusFile } from '@/types/git.types';

defineProps<{
  title: string;
  files: GitStatusFile[];
  selectedFile: GitStatusFile | null;
  isStagedList: boolean;
}>();

const emit = defineEmits<{
  (e: 'select', file: GitStatusFile): void;
  (e: 'action', file: GitStatusFile): void;
  (e: 'actionAll'): void;
}>();
</script>

<template>
  <div class="section flex-grow-1 d-flex flex-column overflow-hidden" :class="{ 'border-top': isStagedList }">
    <div class="section-header px-3 py-2 bg-body-tertiary border-bottom d-flex justify-content-between align-items-center">
      <span class="fw-bold small">{{ title }} ({{ files.length }})</span>
      <button 
        v-if="files.length > 0" 
        class="btn btn-sm btn-ghost p-0" 
        @click="emit('actionAll')" 
        :title="isStagedList ? 'Unstage All' : 'Stage All'"
      >
        <i :class="['ti fs-5', isStagedList ? 'ti-minus text-danger' : 'ti-plus text-success']"></i>
      </button>
    </div>
    <div class="file-list flex-grow-1 overflow-auto">
      <div 
        v-for="file in files" 
        :key="file.path"
        :class="['file-item px-3 py-1 d-flex align-items-center cursor-pointer', { active: selectedFile?.path === file.path && selectedFile?.is_staged === isStagedList }]"
        @click="emit('select', file)"
      >
        <span :class="['status-square me-2', (file.status === '?' || file.status === 'A') ? 'A' : file.status]">
          {{ (file.status === '?' || file.status === 'A') ? '+' : (file.status === 'D' ? '-' : file.status) }}
        </span>
        <span class="file-path text-truncate flex-grow-1">{{ file.path }}</span>
        <button class="btn btn-sm btn-ghost p-0 ms-1 action-btn" @click.stop="emit('action', file)" :title="isStagedList ? 'Unstage' : 'Stage'">
          <i :class="['ti', isStagedList ? 'ti-minus text-danger' : 'ti-plus text-success']"></i>
        </button>
      </div>
      <div v-if="files.length === 0" class="p-3 text-center text-muted small">
        No {{ title.toLowerCase() }} changes
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

.action-btn {
    visibility: hidden;
}
.file-item:hover .action-btn {
    visibility: visible;
}

.btn-ghost {
  border: 1px solid transparent;
  background: transparent;
}
.btn-ghost:hover {
  background-color: var(--bs-secondary-bg);
}
</style>
