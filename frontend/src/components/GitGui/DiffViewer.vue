<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  diff: string;
}>();

const diffLines = computed(() => {
  if (!props.diff) return [];
  
  const lines = props.diff.split('\n');
  const result: any[] = [];
  
  let oldLineNum = 0;
  let newLineNum = 0;
  
  lines.forEach((line, index) => {
    // Skip empty line at the end if it's there
    if (index === lines.length - 1 && line === '') return;

    let type = 'normal';
    let oldNum: number | string = '';
    let newNum: number | string = '';
    
    if (line.startsWith('+')) {
      type = 'addition';
      newNum = newLineNum++;
    } else if (line.startsWith('-')) {
      type = 'deletion';
      oldNum = oldLineNum++;
    } else if (line.startsWith('@@')) {
      type = 'hunk';
      // Parse hunk header: @@ -oldStart,oldCount +newStart,newCount @@
      const match = line.match(/@@ -(\d+),(\d+) \+(\d+),(\d+) @@/);
      if (match) {
        oldLineNum = parseInt(match[1]);
        newLineNum = parseInt(match[3]);
      }
    } else {
      // Unchanged line
      oldNum = oldLineNum++;
      newNum = newLineNum++;
    }
    
    result.push({
      content: line,
      type,
      oldNum,
      newNum,
      id: index
    });
  });
  
  return result;
});
</script>

<template>
  <div class="diff-viewer font-monospace">
    <div 
      v-for="line in diffLines" 
      :key="line.id" 
      :class="['diff-line', line.type]"
    >
      <span class="line-number old">{{ line.oldNum }}</span>
      <span class="line-number new">{{ line.newNum }}</span>
      <span class="line-content">{{ line.content || ' ' }}</span>
    </div>
  </div>
</template>

<style scoped>
.diff-viewer {
  background-color: var(--bs-body-bg);
  color: var(--bs-body-color);
  line-height: 1.4;
  font-size: 0.75rem;
  white-space: pre-wrap;
  word-break: break-all;
  border: 1px solid var(--bs-border-color);
  border-radius: 4px;
}

.diff-line {
  display: flex;
  padding: 0;
}

.line-number {
  width: 40px;
  min-width: 40px;
  text-align: right;
  padding-right: 10px;
  color: var(--bs-secondary-color);
  border-right: 1px solid var(--bs-border-color);
  user-select: none;
  background-color: var(--bs-tertiary-bg);
  opacity: 0.6;
}

.line-number.new {
  border-right: 2px solid var(--bs-border-color);
}

.line-content {
  padding-left: 10px;
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

.hunk {
  background-color: rgba(0, 123, 255, 0.05);
  color: #007bff;
  opacity: 0.8;
}

.header {
  font-weight: bold;
  background-color: var(--bs-tertiary-bg);
}

.line-content {
  display: block;
}

[data-bs-theme='dark'] .addition {
  background-color: rgba(40, 167, 69, 0.2);
  color: #4cd137;
}

[data-bs-theme='dark'] .deletion {
  background-color: rgba(220, 53, 69, 0.2);
  color: #ff4d4d;
}

[data-bs-theme='dark'] .hunk {
  color: #00a8ff;
}
</style>
