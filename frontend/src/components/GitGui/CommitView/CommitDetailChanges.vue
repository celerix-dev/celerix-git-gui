<script setup lang="ts">
import { ref } from 'vue';
import * as App from '../../../../wailsjs/go/backend/App';
import { backend } from "../../../../wailsjs/go/models";
import CommitFileChange = backend.CommitFileChange;
import SideBySideDiff from "./SideBySideDiff.vue";

const props = defineProps<{
  repoPath: string;
  commitHash: string;
  changes: CommitFileChange[];
  loading: boolean;
}>();

const expandedFiles = ref<Set<string>>(new Set());
const fileDiffs = ref<Record<string, string>>({});

const toggleFile = async (path: string) => {
  if (expandedFiles.value.has(path)) {
    expandedFiles.value.delete(path);
  } else {
    expandedFiles.value.add(path);
    if (!fileDiffs.value[path] && props.commitHash) {
      try {
        fileDiffs.value[path] = await App.GetCommitFileDiff(props.repoPath, props.commitHash, path);
      } catch (err) {
        console.error('Failed to load file diff:', err);
      }
    }
  }
};

defineExpose({
    clearExpanded: () => {
        expandedFiles.value.clear();
        fileDiffs.value = {};
    }
});
</script>

<template>
  <div class="commit-detail-changes h-100 overflow-auto">
    <div v-if="loading" class="text-center py-4 text-muted">
        <div class="spinner-border spinner-border-sm me-2" role="status"></div>
        Loading changes...
    </div>
    <div v-else class="commit-changes-container">
        <table class="table table-hover table-sm mb-0">
            <thead>
                <tr class="small text-muted">
                    <th style="width: 30px;"></th>
                    <th style="width: 40px;">Status</th>
                    <th>Path</th>
                </tr>
            </thead>
            <tbody>
                <template v-for="file in changes" :key="file.path">
                    <tr class="cursor-pointer align-middle" @click="toggleFile(file.path)">
                        <td class="text-center">
                            <i :class="['ti', expandedFiles.has(file.path) ? 'ti-chevron-down' : 'ti-chevron-right']"></i>
                        </td>
                        <td>
                            <span :class="['status-square', file.status]">{{ file.status === 'A' ? '+' : (file.status === 'D' ? '-' : 'M') }}</span>
                        </td>
                        <td class="text-truncate file-path">{{ file.path }}</td>
                    </tr>
                    <tr v-if="expandedFiles.has(file.path)">
                        <td colspan="3" class="p-0 border-0">
                            <div class="file-change-diff p-3 bg-body-tertiary border-bottom">
                                <SideBySideDiff v-if="fileDiffs[file.path]" :diff="fileDiffs[file.path]" />
                                <div v-else class="text-center py-2 text-muted small">
                                    <div class="spinner-border spinner-border-sm me-1" role="status"></div>
                                    Loading diff...
                                </div>
                            </div>
                        </td>
                    </tr>
                </template>
            </tbody>
        </table>
        <div v-if="changes.length === 0" class="text-center py-4 text-muted">
            No file changes in this commit.
        </div>
    </div>
  </div>
</template>

<style scoped>
.status-square {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border-radius: 3px;
    font-size: 11px;
    font-weight: bold;
    color: #fff;
}

.status-square.M { background-color: #ff9f43; } /* Modified - Orange */
.status-square.A { background-color: #28c76f; } /* Added - Green */
.status-square.D { background-color: #ea5455; } /* Deleted - Red */

.commit-changes-container .table thead th {
    background-color: var(--bs-body-bg);
    border-bottom: 1px solid var(--bs-border-color);
    position: sticky;
    top: 0;
    z-index: 2;
}

.commit-changes-container .table tbody tr:not(.border-0) td {
    border-bottom: 1px solid var(--bs-border-color);
}

.commit-changes-container .table tr {
    height: auto;
}

.commit-changes-container .table tr td {
    height: auto;
    padding-top: 4px !important;
    padding-bottom: 4px !important;
}
</style>
