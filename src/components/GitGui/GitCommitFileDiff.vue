<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { invoke } from '@tauri-apps/api/core';

const props = defineProps<{
  path: string;
  hash: string;
  repoPath: string;
}>();

const diffLines = ref<{ left: string, right: string, type: 'added' | 'removed' | 'normal' | 'header' }[]>([]);
const loading = ref(true);

const fetchDiff = async () => {
  loading.value = true;
  try {
    const rawDiff = await invoke<string>('get_commit_file_diff', { 
      path: props.repoPath, 
      hash: props.hash, 
      filePath: props.path 
    });
    parseDiff(rawDiff);
  } catch (err) {
    console.error('Failed to fetch commit file diff:', err);
  } finally {
    loading.value = false;
  }
};

const parseDiff = (diff: string) => {
  const lines = diff.split('\n');
  const result: { left: string, right: string, type: 'added' | 'removed' | 'normal' | 'header' | 'mixed' }[] = [];
  
  let headerFinished = false;
  
  let i = 0;
  while (i < lines.length) {
    const line = lines[i];
    
    if (!headerFinished) {
      if (line.startsWith('@@')) {
        headerFinished = true;
        result.push({ left: line, right: line, type: 'header' });
        i++;
        continue;
      }
      if (line.trim() !== '') {
        result.push({ left: line, right: '', type: 'header' });
      }
      i++;
      continue;
    }

    if (line.startsWith('@@')) {
      result.push({ left: line, right: line, type: 'header' });
      i++;
    } else if (line.startsWith('-')) {
      // Collect consecutive removals
      let removals: string[] = [];
      while (i < lines.length && lines[i].startsWith('-')) {
        removals.push(lines[i].substring(1));
        i++;
      }
      // Collect consecutive additions
      let additions: string[] = [];
      while (i < lines.length && lines[i].startsWith('+')) {
        additions.push(lines[i].substring(1));
        i++;
      }
      
      const max = Math.max(removals.length, additions.length);
      for (let j = 0; j < max; j++) {
        result.push({
          left: removals[j] || '',
          right: additions[j] || '',
          type: (removals[j] !== undefined && additions[j] !== undefined) ? 'normal' : (removals[j] !== undefined ? 'removed' : 'added')
        });
      }
    } else if (line.startsWith('+')) {
      result.push({ left: '', right: line.substring(1), type: 'added' });
      i++;
    } else {
      const content = line.startsWith(' ') ? line.substring(1) : line;
      result.push({ left: content, right: content, type: 'normal' });
      i++;
    }
  }
  diffLines.value = result as any;
};

onMounted(fetchDiff);
</script>

<template>
  <div class="diff-side-by-side border rounded overflow-hidden mt-2 bg-body">
    <div v-if="loading" class="p-3 text-center">
      <div class="spinner-border spinner-border-sm text-primary" role="status"></div>
    </div>
    <div v-else class="table-responsive">
      <table class="table table-sm table-borderless mb-0 font-monospace" style="font-size: 0.75rem; table-layout: fixed;">
        <tbody>
          <tr v-for="(line, idx) in diffLines" :key="idx" class="diff-row">
            <td class="diff-left border-end" :class="{'bg-danger-subtle': line.type === 'removed', 'bg-light-subtle': line.type === 'header'}">
              <pre class="mb-0 text-truncate">{{ line.left }}</pre>
            </td>
            <td class="diff-right" :class="{'bg-success-subtle': line.type === 'added', 'bg-light-subtle': line.type === 'header'}">
              <pre class="mb-0 text-truncate">{{ line.right }}</pre>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.font-monospace {
  font-family: var(--bs-font-monospace);
}
.diff-row:hover {
  background-color: rgba(0,0,0,0.02);
}
.diff-left, .diff-right {
  padding: 0 8px;
  overflow: hidden;
  white-space: pre;
}
pre {
  white-space: pre;
}
.bg-danger-subtle {
  background-color: rgba(220, 53, 69, 0.1) !important;
}
.bg-success-subtle {
  background-color: rgba(40, 167, 69, 0.1) !important;
}
</style>
