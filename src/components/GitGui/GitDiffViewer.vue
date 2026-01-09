<script setup lang="ts">
import { onBeforeUnmount, ref, watch, nextTick, onMounted } from 'vue';
import { EditorState, RangeSetBuilder } from '@codemirror/state';
import { EditorView, keymap, Decoration, lineNumbers } from '@codemirror/view';
import { defaultKeymap } from '@codemirror/commands';
import { oneDark } from '@codemirror/theme-one-dark';

const props = defineProps<{
  content: string | null;
  loading: boolean;
  filePath: string | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const editorContainer = ref<HTMLElement | null>(null);
const editorView = ref<EditorView | null>(null);
const currentTheme = ref<'light' | 'dark'>(
  (document.documentElement.getAttribute('data-bs-theme') as 'light' | 'dark') || 'light'
);

// Watch for theme changes on document element
let themeObserver: MutationObserver | null = null;

onMounted(() => {
  themeObserver = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      if (mutation.attributeName === 'data-bs-theme') {
        const newTheme = document.documentElement.getAttribute('data-bs-theme') as 'light' | 'dark';
        if (newTheme && newTheme !== currentTheme.value) {
          currentTheme.value = newTheme;
        }
      }
    });
  });

  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-bs-theme']
  });
});

const addedLine = Decoration.line({ attributes: { class: "cm-diff-added" } });
const removedLine = Decoration.line({ attributes: { class: "cm-diff-removed" } });
const headerLine = Decoration.line({ attributes: { class: "cm-diff-header" } });

const diffHighlightExtension = EditorView.decorations.compute(["doc"], state => {
  const builder = new RangeSetBuilder<Decoration>();
  for (let i = 1; i <= state.doc.lines; i++) {
    const line = state.doc.line(i);
    const text = line.text;
    if (text.startsWith('+') && !text.startsWith('+++')) {
      builder.add(line.from, line.from, addedLine);
    } else if (text.startsWith('-') && !text.startsWith('---')) {
      builder.add(line.from, line.from, removedLine);
    } else if (text.startsWith('@@') || text.startsWith('diff') || text.startsWith('---') || text.startsWith('+++')) {
      builder.add(line.from, line.from, headerLine);
    }
  }
  return builder.finish();
});

const renderEditor = (content: string) => {
  if (!editorContainer.value) return;
  
  if (editorView.value) {
    editorView.value.destroy();
  }

  const extensions = [
    EditorView.editable.of(false),
    EditorState.readOnly.of(true),
    lineNumbers(),
    diffHighlightExtension,
    keymap.of(defaultKeymap),
    EditorView.theme({
      "&": { height: "100%" },
      ".cm-scroller": { overflow: "auto" },
      ".cm-diff-added": { backgroundColor: "rgba(40, 167, 69, 0.2)" },
      ".cm-diff-removed": { backgroundColor: "rgba(220, 53, 69, 0.2)" },
      ".cm-diff-header": { 
        backgroundColor: "rgba(0, 123, 255, 0.1)", 
        color: currentTheme.value === 'dark' ? "#6cb6ff" : "#0056b3", 
        fontWeight: "bold" 
      },
      ".cm-gutterElement": { color: currentTheme.value === 'dark' ? "#5c6370" : "#999" }
    })
  ];

  if (currentTheme.value === 'dark') {
    extensions.push(oneDark);
  }

  const state = EditorState.create({
    doc: content,
    extensions
  });

  editorView.value = new EditorView({
    state,
    parent: editorContainer.value
  });
};

watch(() => props.content, (newContent) => {
  if (newContent !== null) {
    nextTick(() => {
      renderEditor(newContent);
    });
  }
}, { immediate: true });

watch(currentTheme, () => {
  if (props.content !== null) {
    renderEditor(props.content);
  }
});

onBeforeUnmount(() => {
  if (editorView.value) {
    editorView.value.destroy();
  }
  if (themeObserver) {
    themeObserver.disconnect();
  }
});
</script>

<template>
  <div class="h-100 d-flex flex-column bg-body border-start">
    <div class="p-2 border-bottom d-flex justify-content-between align-items-center bg-light-subtle">
      <span class="small fw-bold text-truncate">{{ filePath }}</span>
      <button class="btn btn-sm btn-ghost p-0" @click="emit('close')">
        <i class="ti ti-x"></i>
      </button>
    </div>
    
    <div class="flex-grow-1 overflow-hidden position-relative">
      <div v-if="loading" class="position-absolute top-50 start-50 translate-middle" style="z-index: 10">
        <div class="spinner-border text-primary" role="status"></div>
      </div>
      <div ref="editorContainer" class="h-100 selectable codemirror-container"></div>
    </div>
    
    <!-- Slot for Commit Pane -->
    <slot name="footer"></slot>
  </div>
</template>

<style scoped>
.codemirror-container :deep(.cm-editor) {
  height: 100%;
}
:deep(.cm-diff-added) {
  background-color: rgba(40, 167, 69, 0.2) !important;
}
:deep(.cm-diff-removed) {
  background-color: rgba(220, 53, 69, 0.2) !important;
}
:deep(.cm-diff-header) {
  background-color: rgba(0, 123, 255, 0.1) !important;
  color: #6cb6ff !important;
}
.btn-ghost {
  background: transparent;
  border: none;
  color: inherit;
}
.btn-ghost:hover {
  background: rgba(0, 0, 0, 0.05);
}
</style>
