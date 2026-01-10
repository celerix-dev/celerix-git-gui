<script setup lang="ts">

import type {RepoTab} from "@/types/git.types";

defineProps<{
  activeTab: RepoTab | null;
}>();

</script>

<template>
  <div class="status-bar d-flex align-items-center justify-content-between px-3 border-bottom position-relative" style="height: 60px;">
    <!-- Left: Action Buttons -->
    <div class="d-flex align-items-center gap-2 z-1">
      <div class="btn-group">
        <button class="btn btn-sm btn-ghost d-flex flex-column align-items-center gap-1 px-3" :disabled="!activeTab">
          <i class="ti ti-refresh fs-4"></i>
          <span class="x-small">Fetch</span>
        </button>
        <button class="btn btn-sm btn-ghost d-flex flex-column align-items-center gap-1 px-3" :disabled="!activeTab">
          <i class="ti ti-download fs-4"></i>
          <span class="x-small">Pull</span>
        </button>
        <button class="btn btn-sm btn-ghost d-flex flex-column align-items-center gap-1 px-3" :disabled="!activeTab">
          <i class="ti ti-upload fs-4"></i>
          <span class="x-small">Push</span>
        </button>
      </div>

      <div class="btn-group h-100">
        <button class="btn btn-sm btn-ghost d-flex flex-column align-items-center gap-1 px-3" :disabled="!activeTab">
          <i class="ti ti-archive fs-4"></i>
          <span class="x-small">Stash</span>
        </button>
        <button type="button" class="btn btn-sm btn-ghost dropdown-toggle dropdown-toggle-split px-2" 
                data-bs-toggle="dropdown" aria-expanded="false" :disabled="!activeTab">
          <span class="visually-hidden">Toggle Dropdown</span>
        </button>
        <ul class="dropdown-menu shadow">
          <li><a class="dropdown-item" href="#">Stash changes</a></li>
          <li><a class="dropdown-item" href="#">Stash pop</a></li>
          <li><hr class="dropdown-divider"></li>
          <li><a class="dropdown-item disabled" href="#">No stashes found</a></li>
        </ul>
      </div>
    </div>

    <!-- Middle: Status Window (Centered) -->
    <div class="status-window d-flex flex-column justify-content-center border rounded bg-body-tertiary px-3 position-absolute start-50 translate-middle-x overflow-hidden" 
         style="width: 500px; height: 50px;">
      <div class="d-flex justify-content-between align-items-center small">
        <span class="text-truncate">{{ activeTab ? 'Ready' : 'No repository selected' }}</span>
        <span class="text-muted" v-if="activeTab">{{ activeTab.name }}</span>
      </div>
      <div class="progress position-absolute bottom-0 start-0 w-100 rounded-0" style="height: 4px;" v-if="activeTab">
        <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: 100%"></div>
      </div>
    </div>

    <!-- Right: Dummy Options -->
    <div class="d-flex align-items-center gap-2 z-1">
      <div class="btn-group">
        <button class="btn btn-sm btn-ghost d-flex flex-column align-items-center gap-1 px-2" title="Search">
          <i class="ti ti-search fs-4"></i>
          <span class="x-small">Search</span>
        </button>
        <button class="btn btn-sm btn-ghost d-flex flex-column align-items-center gap-1 px-2" title="Filter">
          <i class="ti ti-filter fs-4"></i>
          <span class="x-small">Filter</span>
        </button>
        <button class="btn btn-sm btn-ghost d-flex flex-column align-items-center gap-1 px-2" title="More options">
          <i class="ti ti-dots-vertical fs-4"></i>
          <span class="x-small">More</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.status-bar {
  background-color: var(--bs-body-bg);
}

.status-window {
  border-color: var(--bs-border-color) !important;
}

.x-small {
  font-size: 0.7rem;
  font-weight: 500;
  text-transform: uppercase;
}

.btn-ghost {
  border: 1px solid transparent;
  color: var(--bs-secondary-color);
}

.btn-ghost:hover:not(:disabled) {
  background-color: var(--bs-tertiary-bg);
  color: var(--bs-primary);
}

.btn-ghost:disabled {
  opacity: 0.5;
}
</style>