<script setup lang="ts">
import type { RepoTab } from './types';

defineProps<{
  activeTab: RepoTab | null;
}>();

const emit = defineEmits<{
  (e: 'fetch', tab: RepoTab): void;
  (e: 'pull', tab: RepoTab): void;
  (e: 'push', tab: RepoTab): void;
  (e: 'openSettings'): void;
  (e: 'openRepo'): void;
  (e: 'createBranch', branchName: string): void;
}>();
</script>

<template>
  <div class="d-flex align-items-center justify-content-between mb-3 p-2 rounded" style="z-index: 1024">
    <div class="d-flex align-items-center gap-4">
      <!-- Status Object -->
      <template v-if="activeTab">
        <div class="d-flex align-items-center gap-3 py-1 px-3 bg-body rounded border shadow-sm" style="min-height:38px;">
          <div class="d-flex align-items-center">
            <i class="ti ti-folder me-2 text-primary"></i>
            <span class="fw-bold small">{{ activeTab.name }}</span>
          </div>
          <div class="vr"></div>
          <div class="d-flex align-items-center">
            <i class="ti ti-git-branch me-2 text-success"></i>
            <span class="small">{{ activeTab.branches.find(b => b.is_current)?.name || 'Unknown' }}</span>
          </div>
          
          <template v-if="activeTab.operationLoading">
            <div class="vr"></div>
            <div class="d-flex flex-column justify-content-center" style="min-width: 120px">
              <div class="d-flex align-items-center text-primary mb-1">
                <div class="spinner-border spinner-border-sm me-2" role="status"></div>
                <span class="smaller fw-medium">{{ activeTab.operationName }}</span>
              </div>
              <div class="progress" style="height: 4px;">
                <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: 100%"></div>
              </div>
            </div>
          </template>
        </div>

          <div class="btn-group shadow-sm">
            <button
              class="btn btn-secondary d-flex align-items-center"
              @click="emit('fetch', activeTab)"
              :disabled="activeTab.operationLoading || activeTab.loading"
            >
              <i class="ti ti-cloud-download me-1"></i> Fetch
            </button>
            <button
              class="btn btn-secondary d-flex align-items-center"
              @click="emit('pull', activeTab)"
              :disabled="activeTab.operationLoading || activeTab.loading"
            >
              <i class="ti ti-arrow-down me-1"></i> Pull
            </button>
            <button
              class="btn btn-secondary d-flex align-items-center"
              @click="emit('push', activeTab)"
              :disabled="activeTab.operationLoading || activeTab.loading"
            >
              <i class="ti ti-arrow-up me-1"></i> Push
            </button>

            <div class="btn-group" role="group">
              <button type="button" class="btn btn-secondary dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">
                More
              </button>
              <ul class="dropdown-menu">
                <li><a class="dropdown-item" href="#" @click.prevent="emit('openSettings')"><i class="ti ti-settings me-1"></i> Settings</a></li>
                <li><a class="dropdown-item" href="#" @click.prevent="emit('openRepo')"><i class="ti ti-folder-open me-1"></i> Open Repository</a></li>
                <li><a class="dropdown-item" href="#" @click.prevent="emit('createBranch', activeTab.branches.find(b => b.is_current)?.name || 'HEAD')"><i class="ti ti-git-branch me-1"></i> Create branch...</a></li>
              </ul>
            </div>

          </div>

      </template>
      <template v-else>
<!--        <teleport to="#page-context">-->
          <div class="btn-group" role="group">
            <button class="btn btn-secondary" @click.prevent="emit('openSettings')"><i class="ti ti-settings me-1"></i> Settings</button>
            <button class="btn btn-secondary" @click.prevent="emit('openRepo')"><i class="ti ti-folder-open me-1"></i> Open Repository</button>
          </div>
<!--        </teleport>-->
      </template>
    </div>
    <div class="d-flex gap-2">

    </div>
  </div>
</template>

<style scoped>
.smaller {
  font-size: 0.75rem;
}
</style>
