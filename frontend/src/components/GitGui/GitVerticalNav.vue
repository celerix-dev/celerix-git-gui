<script setup lang="ts">
import type { RepoTab } from "@/types/git.types";

defineProps<{
  activeVerticalTab: string;
}>();

const emit = defineEmits<{
  (e: 'setActiveVerticalTab', tab: 'info' | 'local-changes' | 'commit' | 'placeholder1' | 'placeholder2'): void;
}>();
</script>

<template>
  <div class="vertical-tabs d-flex flex-column border-end py-2">
    <div
        :class="['vertical-tab-item', { active: activeVerticalTab === 'info' }]"
        @click="emit('setActiveVerticalTab', 'info')"
        title="Repository Info"
    >
      <i class="ti ti-info-circle fs-4"></i>
    </div>
    <div
        :class="['vertical-tab-item', { active: activeVerticalTab === 'local-changes' }]"
        @click="emit('setActiveVerticalTab', 'local-changes')"
        title="Local Changes"
    >
      <i class="ti ti-edit fs-4"></i>
    </div>
    <div
        :class="['vertical-tab-item', { active: activeVerticalTab === 'commit' }]"
        @click="emit('setActiveVerticalTab', 'commit')"
        title="Commit Info"
    >
      <i class="ti ti-git-commit fs-4"></i>
    </div>
    <div
        :class="['vertical-tab-item', { active: activeVerticalTab === 'placeholder1' }]"
        @click="emit('setActiveVerticalTab', 'placeholder1')"
        title="Branch Management"
    >
      <i class="ti ti-git-branch fs-4"></i>
    </div>
    <div
        :class="['vertical-tab-item', { active: activeVerticalTab === 'placeholder2' }]"
        @click="emit('setActiveVerticalTab', 'placeholder2')"
        title="Settings"
    >
      <i class="ti ti-settings fs-4"></i>
    </div>
  </div>
</template>

<style scoped>
.vertical-tabs {
  width: 50px;
  background-color: var(--bs-tertiary-bg);
}

.vertical-tab-item {
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--bs-secondary-color);
  transition: all 0.2s ease;
  border-right: 3px solid transparent;
}

.vertical-tab-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--bs-primary);
}

[data-bs-theme='dark'] .vertical-tab-item:hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.vertical-tab-item.active {
  color: var(--bs-primary);
  border-right-color: var(--bs-primary);
  background-color: var(--bs-body-bg);
}
</style>
