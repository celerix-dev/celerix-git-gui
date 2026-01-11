<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  show: boolean;
  branchName: string;
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'delete', data: { branchName: string; deleteRemote: boolean }): void;
}>();

const deleteRemote = ref(false);

watch(() => props.show, (newVal) => {
  if (newVal) {
    deleteRemote.value = false;
  }
});

const handleDelete = () => {
  emit('delete', {
    branchName: props.branchName,
    deleteRemote: deleteRemote.value
  });
};
</script>

<template>
  <div v-if="show" class="modal-backdrop fade show"></div>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" role="dialog">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content border-0 shadow-lg">
        <div class="modal-header border-bottom-0 pb-0">
          <h5 class="modal-title fw-bold d-flex align-items-center">
            <i class="ti ti-trash me-2 text-danger"></i>
            Delete Branch
          </h5>
          <button type="button" class="btn-close" @click="emit('close')" :disabled="loading"></button>
        </div>
        <div class="modal-body py-4">
          <p class="text-muted mb-4">Delete local branch from your repository.</p>
          
          <div class="mb-3">
            <label class="form-label small text-muted text-uppercase fw-bold mb-1">Branch to delete</label>
            <div class="d-flex align-items-center p-2 bg-body-tertiary rounded border">
              <i class="ti ti-git-branch me-2 text-primary"></i>
              <span class="fw-semibold">{{ branchName }}</span>
            </div>
          </div>

          <div class="form-check form-switch mt-4">
            <input 
              class="form-check-input" 
              type="checkbox" 
              role="switch" 
              id="deleteRemote" 
              v-model="deleteRemote"
              :disabled="loading"
            >
            <label class="form-check-label" for="deleteRemote">Also delete remote branch</label>
          </div>
        </div>
        <div class="modal-footer border-top-0 pt-0">
          <button type="button" class="btn btn-secondary px-4" @click="emit('close')" :disabled="loading">Cancel</button>
          <button 
            type="button" 
            class="btn btn-danger px-4" 
            :disabled="loading"
            @click="handleDelete"
          >
            <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
            Delete Branch
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
