<script setup lang="ts">
import { computed } from "vue";
import draggable from "vuedraggable";
import type { RepoTab } from "@/types/git.types";

const props = defineProps<{
  tabs: RepoTab[];
  activeTabId: string | null;
}>();

const emit = defineEmits<{
  (e: 'update:tabs', value: RepoTab[]): void;
  (e: 'setActiveTab', id: string): void;
  (e: 'closeTab', id: string): void;
  (e: 'saveState'): void;
}>();

const localTabs = computed({
  get: () => props.tabs,
  set: (value) => emit('update:tabs', value)
});

</script>

<template>
  <div class="card border-radius-0 m-1 mb-0 ms-0 me-0 border-start-0 p-0 flex-shrink-0" v-if="tabs.length > 0"
       style="border-radius:0">
    <div class="card-body p-2">
      <draggable
          v-model="localTabs"
          item-key="id"
          tag="ul"
          class="nav nav-pills nav-justified"
          ghost-class="ghost-tab"
          @end="emit('saveState')"
      >
        <template #item="{ element: tab }">
          <li class="nav-item">
            <div
                :class="['nav-link d-flex align-items-center cursor-pointer', { active: activeTabId === tab.id }]"
                @click="emit('setActiveTab', tab.id)"
            >
              <i class="ti ti-x ms-1 close-icon" @click.stop="emit('closeTab', tab.id)"></i>
              <div class="ms-2 flex-grow-1 d-flex justify-content-center">{{ tab.name }}</div>
            </div>
          </li>
        </template>
      </draggable>
    </div>
  </div>
</template>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}

.nav-item {
  cursor: move;
}

.ghost-tab {
  opacity: 0.5;
  background: var(--bs-light);
}

.close-icon {
  visibility: hidden;
  opacity: 0;
  transition: opacity 0.2s, visibility 0.2s;
}

.nav-link:hover .close-icon {
  visibility: visible;
  opacity: 1;
}

.nav-link:not(.active):hover {
  background-color: rgba(0, 0, 0, 0.06);
}

[data-bs-theme='dark'] .nav-link:not(.active):hover {
  background-color: rgba(255, 255, 255, 0.05);
}

.close-icon:hover {
  color: var(--bs-danger);
}
</style>
