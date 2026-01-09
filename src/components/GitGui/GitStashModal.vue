<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
  show: boolean;
  files: string[];
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'confirm', files: string[], message: string): void;
}>();

const stashMessage = ref('');
const selectedFiles = ref<string[]>([...props.files]);

// Watch for changes in props.files and update selectedFiles
import { watch } from 'vue';
watch(() => props.files, (newFiles) => {
  selectedFiles.value = [...newFiles];
}, { immediate: true });

const toggleFile = (file: string) => {
  const index = selectedFiles.value.indexOf(file);
  if (index === -1) {
    selectedFiles.value.push(file);
  } else {
    selectedFiles.value.splice(index, 1);
  }
};

const handleConfirm = () => {
  emit('confirm', selectedFiles.value, stashMessage.value);
};
</script>

<template>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" style="background: rgba(0,0,0,0.5)">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content border-0 shadow">
        <div class="modal-header border-bottom-0 pb-0">
          <h5 class="modal-title fw-bold">Stash Changes</h5>
          <button type="button" class="btn-close" @click="emit('close')"></button>
        </div>
        <div class="modal-body py-4">
          <div class="mb-4">
            <label class="form-label small text-muted text-uppercase fw-bold mb-1">Stash Message (optional)</label>
            <input 
              type="text" 
              class="form-control" 
              placeholder="e.g. Working on feature X" 
              v-model="stashMessage"
              @keyup.enter="handleConfirm"
              autofocus
            />
          </div>
          <div>
            <label class="form-label small text-muted text-uppercase fw-bold mb-1">Files to Stash ({{ selectedFiles.length }})</label>
            <div class="list-group list-group-flush border rounded overflow-auto" style="max-height: 250px;">
              <div 
                v-for="file in files" 
                :key="file" 
                class="list-group-item d-flex align-items-center py-2 px-3"
              >
                <div class="form-check mb-0">
                  <input 
                    class="form-check-input" 
                    type="checkbox" 
                    :id="'file-' + file"
                    :checked="selectedFiles.includes(file)"
                    @change="toggleFile(file)"
                  >
                  <label class="form-check-label small text-truncate d-block cursor-pointer" :for="'file-' + file" style="max-width: 400px;">
                    {{ file }}
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer border-top-0 pt-0">
          <button type="button" class="btn btn-ghost" @click="emit('close')">Cancel</button>
          <button 
            type="button" 
            class="btn btn-primary px-4" 
            :disabled="selectedFiles.length === 0"
            @click="handleConfirm"
          >
            Stash Changes
          </button>
        </div>
      </div>
    </div>
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
</style>
