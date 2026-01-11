<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  diff: string;
}>();

const hunks = computed(() => {
  if (!props.diff) return [];
  
  const lines = props.diff.split('\n');
  const hunksList: any[] = [];
  let currentHunk: any = null;
  
  let oldLineNum = 0;
  let newLineNum = 0;
  
  lines.forEach((line) => {
    if (line.startsWith('@@')) {
      if (currentHunk) hunksList.push(currentHunk);
      
      const match = line.match(/@@ -(\d+),(\d+) \+(\d+),(\d+) @@/);
      if (match) {
        oldLineNum = parseInt(match[1]);
        newLineNum = parseInt(match[3]);
      }
      
      currentHunk = {
        header: line,
        lines: []
      };
    } else if (currentHunk) {
      if (line.startsWith('+')) {
        currentHunk.lines.push({ type: 'addition', content: line.substring(1), oldNum: '', newNum: newLineNum++ });
      } else if (line.startsWith('-')) {
        currentHunk.lines.push({ type: 'deletion', content: line.substring(1), oldNum: oldLineNum++, newNum: '' });
      } else if (line !== '' || lines.indexOf(line) !== lines.length - 1) {
        currentHunk.lines.push({ type: 'normal', content: line, oldNum: oldLineNum++, newNum: newLineNum++ });
      }
    }
  });
  
  if (currentHunk) hunksList.push(currentHunk);
  
  // Process lines into side-by-side rows
  return hunksList.map(hunk => {
    const rows: any[] = [];
    let i = 0;
    while (i < hunk.lines.length) {
      const line = hunk.lines[i];
      if (line.type === 'deletion') {
        // Look ahead for matching addition
        let nextAdditionIndex = -1;
        for (let j = i + 1; j < hunk.lines.length; j++) {
            if (hunk.lines[j].type !== 'deletion' && hunk.lines[j].type !== 'addition') break;
            if (hunk.lines[j].type === 'addition') {
                nextAdditionIndex = j;
                break;
            }
        }
        
        if (nextAdditionIndex !== -1) {
            // Found a matching addition, pair them
            const addition = hunk.lines.splice(nextAdditionIndex, 1)[0];
            rows.push({
                left: line,
                right: addition
            });
        } else {
            rows.push({
                left: line,
                right: null
            });
        }
      } else if (line.type === 'addition') {
        rows.push({
            left: null,
            right: line
        });
      } else {
        rows.push({
            left: line,
            right: line
        });
      }
      i++;
    }
    return { ...hunk, rows };
  });
});
</script>

<template>
  <div class="side-by-side-diff font-monospace">
    <div v-for="(hunk, hIdx) in hunks" :key="hIdx" class="hunk">
      <div class="hunk-header">{{ hunk.header }}</div>
      <div class="diff-table">
        <div v-for="(row, rIdx) in hunk.rows" :key="rIdx" class="diff-row">
          <div :class="['diff-cell left', row.left?.type]">
            <span class="line-number">{{ row.left?.oldNum }}</span>
            <span class="line-content">{{ row.left?.content || ' ' }}</span>
          </div>
          <div :class="['diff-cell right', row.right?.type]">
            <span class="line-number">{{ row.right?.newNum }}</span>
            <span class="line-content">{{ row.right?.content || ' ' }}</span>
          </div>
        </div>
      </div>
    </div>
    <div v-if="hunks.length === 0" class="p-3 text-center text-muted">
        No changes or binary file.
    </div>
  </div>
</template>

<style scoped>
.side-by-side-diff {
  background-color: var(--bs-body-bg);
  border: 1px solid var(--bs-border-color);
  border-radius: 4px;
  overflow-x: auto;
  font-size: 0.75rem;
}

.hunk-header {
  background-color: rgba(0, 123, 255, 0.05);
  color: #007bff;
  padding: 2px 10px;
  border-bottom: 1px solid var(--bs-border-color);
  border-top: 1px solid var(--bs-border-color);
}

.hunk:first-child .hunk-header {
  border-top: none;
}

.diff-table {
  display: flex;
  flex-direction: column;
  min-width: max-content;
  width: 100%;
}

.diff-row {
  display: flex;
  border-bottom: 1px solid rgba(var(--bs-border-color-rgb), 0.3);
}

.diff-row:last-child {
  border-bottom: none;
}

.diff-cell {
  flex: 0 0 50%;
  display: flex;
  min-width: 0;
}

.diff-cell.left {
  border-right: 1px solid var(--bs-border-color);
}

.line-number {
  width: 40px;
  min-width: 40px;
  text-align: right;
  padding-right: 10px;
  color: var(--bs-secondary-color);
  background-color: var(--bs-tertiary-bg);
  user-select: none;
  opacity: 0.6;
}

.line-content {
  padding-left: 8px;
  white-space: pre;
  flex-grow: 1;
}

.addition {
  background-color: rgba(40, 167, 69, 0.15);
  color: #28a745;
}

.deletion {
  background-color: rgba(220, 53, 69, 0.15);
  color: #dc3545;
}

[data-bs-theme='dark'] .addition {
  background-color: rgba(40, 167, 69, 0.2);
  color: #4cd137;
}

[data-bs-theme='dark'] .deletion {
  background-color: rgba(220, 53, 69, 0.2);
  color: #ff4d4d;
}
</style>
