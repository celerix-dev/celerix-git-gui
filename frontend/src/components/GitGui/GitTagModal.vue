<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  show: boolean;
  fromBranch: string;
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'create', data: { name: string; message: string }): void;
}>();

const tagName = ref('');
const tagMessage = ref('');

watch(() => props.show, (newVal) => {
  if (newVal) {
    tagName.value = '';
    tagMessage.value = '';
  }
});

const handleCreate = () => {
  if (tagName.value.trim()) {
    emit('create', {
      name: tagName.value.trim(),
      message: tagMessage.value.trim()
    });
  }
};
</script>

<template>
  <div v-if="show" class="modal-backdrop fade show"></div>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" role="dialog">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content border-0 shadow-lg">
        <div class="modal-header border-bottom-0 pb-0">
          <h5 class="modal-title fw-bold d-flex align-items-center">
            <i class="ti ti-tag me-2 text-primary"></i>
            Add Tag
          </h5>
          <button type="button" class="btn-close" @click="emit('close')" :disabled="loading"></button>
        </div>
        <div class="modal-body py-4">
          <div class="mb-3">
            <label class="form-label small text-muted text-uppercase fw-bold mb-1">Creating from</label>
            <div class="d-flex align-items-center p-2 bg-body-tertiary rounded border">
              <i class="ti ti-git-branch me-2 text-primary"></i>
              <span class="fw-semibold">{{ fromBranch }}</span>
            </div>
          </div>
          
          <div class="mb-3">
            <label for="tagName" class="form-label small text-muted text-uppercase fw-bold mb-1">Tag Name</label>
            <input 
              type="text" 
              class="form-control" 
              id="tagName" 
              v-model="tagName" 
              placeholder="e.g. v1.0.0"
              @keyup.enter="handleCreate"
              :disabled="loading"
              autofocus
            >
          </div>

          <div class="mb-3">
            <label for="tagMessage" class="form-label small text-muted text-uppercase fw-bold mb-1">Message (optional)</label>
            <textarea 
              class="form-control" 
              id="tagMessage" 
              v-model="tagMessage" 
              placeholder="Enter tag message..."
              rows="3"
              :disabled="loading"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer border-top-0 pt-0">
          <button type="button" class="btn btn-secondary px-4" @click="emit('close')" :disabled="loading">Cancel</button>
          <button 
            type="button" 
            class="btn btn-primary px-4" 
            :disabled="!tagName.trim() || loading"
            @click="handleCreate"
          >
            <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
            Add Tag
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
