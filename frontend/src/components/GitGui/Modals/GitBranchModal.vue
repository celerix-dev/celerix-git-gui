<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  show: boolean;
  fromBranch: string;
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'create', data: { name: string; checkout: boolean }): void;
}>();

const branchName = ref('');
const checkoutAfter = ref(true);

watch(() => props.show, (newVal) => {
  if (newVal) {
    branchName.value = '';
    checkoutAfter.value = true;
  }
});

const handleCreate = () => {
  if (branchName.value.trim()) {
    emit('create', {
      name: branchName.value.trim(),
      checkout: checkoutAfter.value
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
            <i class="ti ti-git-branch me-2 text-primary"></i>
            Create Branch
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
            <label for="newBranchName" class="form-label small text-muted text-uppercase fw-bold mb-1">New Branch Name</label>
            <input 
              type="text" 
              class="form-control" 
              id="newBranchName" 
              v-model="branchName" 
              placeholder="e.g. feature/my-new-feature"
              @keyup.enter="handleCreate"
              :disabled="loading"
              autofocus
            >
          </div>

          <div class="form-check form-switch mt-4">
            <input 
              class="form-check-input" 
              type="checkbox" 
              role="switch" 
              id="checkoutAfter" 
              v-model="checkoutAfter"
              :disabled="loading"
            >
            <label class="form-check-label" for="checkoutAfter">Checkout branch after creation</label>
          </div>
        </div>
        <div class="modal-footer border-top-0 pt-0">
          <button type="button" class="btn btn-secondary px-4" @click="emit('close')" :disabled="loading">Cancel</button>
          <button 
            type="button" 
            class="btn btn-primary px-4" 
            :disabled="!branchName.trim() || loading"
            @click="handleCreate"
          >
            <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
            Create Branch
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
