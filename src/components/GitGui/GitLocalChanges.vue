<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';
import type { GitStatusFile } from './types';

const props = defineProps<{
  statusFiles: GitStatusFile[];
  selectedFilePath: string | null;
  selectedFileStaged: boolean | null;
  selectedFiles: string[];
}>();

const emit = defineEmits<{
  (e: 'select', path: string, isStaged: boolean, multi?: boolean): void;
  (e: 'stage', path: string): void;
  (e: 'unstage', path: string): void;
  (e: 'stageAll'): void;
  (e: 'unstageAll'): void;
  (e: 'stashSelected', files: string[]): void;
  (e: 'discardSelected', files: string[]): void;
}>();

const unstagedFiles = computed(() => props.statusFiles.filter(f => !f.is_staged));
const stagedFiles = computed(() => props.statusFiles.filter(f => f.is_staged));

const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  files: [] as string[]
});

const showContextMenu = (event: MouseEvent, file: GitStatusFile) => {
  if (file.is_staged) return; // Only context menu for unstaged for now as per requirement
  
  let filesToMenu = [...props.selectedFiles];
  if (!filesToMenu.includes(file.path)) {
    filesToMenu = [file.path];
    emit('select', file.path, false, false);
  }

  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    files: filesToMenu
  };
};

const closeContextMenu = () => {
  contextMenu.value.show = false;
};

const handleStash = () => {
  emit('stashSelected', contextMenu.value.files);
  closeContextMenu();
};

const handleDiscard = () => {
  emit('discardSelected', contextMenu.value.files);
  closeContextMenu();
};

onMounted(() => {
  window.addEventListener('click', closeContextMenu);
});

onUnmounted(() => {
  window.removeEventListener('click', closeContextMenu);
});

const getStatusBadgeClass = (status: string) => {
  const s = status.trim();
  if (s === 'M') return 'bg-warning text-dark';
  if (s === 'A') return 'bg-success';
  if (s === 'D') return 'bg-danger';
  if (s === 'R') return 'bg-info text-dark';
  if (s === '??') return 'bg-secondary';
  if (s === 'UU') return 'bg-danger';
  return 'bg-light text-dark';
};
</script>

<template>
  <div :class="[selectedFilePath ? 'col-md-3' : 'col-md-12', 'h-100 border-end d-flex flex-column']">
    <div v-if="statusFiles.length === 0" class="h-100 d-flex flex-column align-items-center justify-content-center text-muted">
      <i class="ti ti-circle-check fs-2 mb-2 text-success"></i>
      <p>No local changes. Your working directory is clean.</p>
    </div>
    <template v-else>
      <!-- Unstaged Files (Top 60%) -->
      <div class="border-bottom d-flex flex-column" style="height: 60%;">
        <div style="min-height:30px" class="px-3 py-1 bg-light-subtle border-bottom small fw-bold d-flex justify-content-between align-items-center">
          <span>UNSTAGED ({{ unstagedFiles.length }})</span>
          <button 
            v-if="unstagedFiles.length > 0"
            class="btn btn-ghost p-0"
            @click.stop="emit('stageAll')" 
            title="Stage All Files"
          >
            <i class="ti ti-chevrons-down text-success"></i>
          </button>
        </div>
        <div class="flex-grow-1 overflow-auto list-group list-group-flush">
          <div 
            v-for="file in unstagedFiles" 
            :key="'unstaged-' + file.path" 
            class="list-group-item d-flex align-items-center cursor-pointer py-1 px-3"
            :class="{ 'bg-primary-subtle': selectedFiles.includes(file.path) && !selectedFileStaged }"
            @click="emit('select', file.path, false, $event.ctrlKey || $event.metaKey || $event.shiftKey)"
            @contextmenu.prevent="showContextMenu($event, file)"
          >
            <span :class="['badge me-2 font-monospace smaller', getStatusBadgeClass(file.status)]" style="width: 24px">
              {{ file.status.trim() }}
            </span>
            <div class="flex-grow-1 text-truncate">
              <div class="fw-medium small text-truncate">{{ file.path.split('/').pop() }}</div>
            </div>
            <button class="btn btn-xs btn-ghost p-0" @click.stop="emit('stage', file.path)" title="Stage File">
              <i class="ti ti-plus text-success"></i>
            </button>
          </div>
        </div>
      </div>

      <!-- Context Menu -->
      <div 
        v-if="contextMenu.show" 
        class="dropdown-menu show shadow-sm border position-fixed" 
        :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      >
        <button class="dropdown-item d-flex align-items-center" @click="handleStash">
          <i class="ti ti-archive me-2"></i> Stash ({{ contextMenu.files.length }} files)
        </button>
        <button class="dropdown-item d-flex align-items-center text-danger" @click="handleDiscard">
          <i class="ti ti-trash me-2"></i> Discard Changes
        </button>
      </div>

      <!-- Staged Files (Bottom 40%) -->
      <div class="d-flex flex-column" style="height: 40%;">
        <div style="min-height:30px" class="px-3 py-1 bg-light-subtle border-bottom small fw-bold d-flex justify-content-between align-items-center">
          <span>STAGED ({{ stagedFiles.length }})</span>
          <button 
            v-if="stagedFiles.length > 0"
            class="btn btn-ghost p-0"
            @click.stop="emit('unstageAll')" 
            title="Unstage All Files"
          >
            <i class="ti ti-chevrons-up text-danger"></i>
          </button>
        </div>
        <div class="flex-grow-1 overflow-auto list-group list-group-flush">
          <div 
            v-for="file in stagedFiles" 
            :key="'staged-' + file.path" 
            class="list-group-item d-flex align-items-center cursor-pointer py-1 px-3"
            :class="{ 'bg-primary-subtle': selectedFilePath === file.path && selectedFileStaged }"
            @click="emit('select', file.path, true)"
          >
            <span :class="['badge me-2 font-monospace smaller', getStatusBadgeClass(file.status)]" style="width: 24px">
              {{ file.status.trim() }}
            </span>
            <div class="flex-grow-1 text-truncate">
              <div class="fw-medium small text-truncate">{{ file.path.split('/').pop() }}</div>
            </div>
            <button class="btn btn-xs btn-ghost p-0" @click.stop="emit('unstage', file.path)" title="Unstage File">
              <i class="ti ti-minus text-danger"></i>
            </button>
          </div>
        </div>
      </div>
    </template>
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
.cursor-pointer {
  cursor: pointer;
}
.smaller {
  font-size: 0.75rem;
}
.font-monospace {
  font-family: var(--bs-font-monospace);
}
.btn-xs {
  padding: 0.1rem 0.4rem;
  font-size: 0.7rem;
  border-radius: 0.2rem;
}
</style>
