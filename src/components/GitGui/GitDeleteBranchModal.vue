<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  show: boolean;
  branchName: string;
  hasRemote: boolean;
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'delete', data: { deleteRemote: boolean }): void;
}>();

const deleteRemote = ref(false);

watch(() => props.show, (newVal) => {
  if (newVal) {
    deleteRemote.value = false;
  }
});

const handleDelete = () => {
  emit('delete', {
    deleteRemote: deleteRemote.value
  });
};
</script>

<template>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" style="background: rgba(0,0,0,0.5)">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content border-0 shadow">
        <div class="modal-header border-bottom-0 pb-0">
          <h5 class="modal-title fw-bold text-danger">Delete Branch</h5>
          <button type="button" class="btn-close" @click="emit('close')" :disabled="loading"></button>
        </div>
        <div class="modal-body py-4">
          <p>Are you sure you want to delete the branch <span class="fw-bold text-primary">{{ branchName }}</span>?</p>
          <p class="text-muted small">This action cannot be undone if the branch has not been pushed to a remote.</p>
          
          <div class="form-check form-switch mt-4" :class="{ 'opacity-50': !hasRemote }">
            <input 
              class="form-check-input" 
              type="checkbox" 
              role="switch" 
              id="deleteRemote" 
              v-model="deleteRemote"
              :disabled="loading || !hasRemote"
            >
            <label class="form-check-label" for="deleteRemote">
              Also delete corresponding remote branch
              <span v-if="!hasRemote" class="d-block smaller text-muted">(No remote branch detected)</span>
            </label>
          </div>
        </div>
        <div class="modal-footer border-top-0 pt-0">
          <button type="button" class="btn btn-ghost" @click="emit('close')" :disabled="loading">Cancel</button>
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
.btn-ghost {
  background: transparent;
  border: none;
  color: inherit;
}
.btn-ghost:hover {
  background: rgba(0, 0, 0, 0.05);
}
.smaller {
  font-size: 0.75rem;
}
</style>
