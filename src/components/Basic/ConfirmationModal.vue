<script setup lang="ts">

defineProps<{
  show: boolean;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  variant?: 'primary' | 'danger' | 'warning' | 'info';
  loading?: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'confirm'): void;
}>();


// Optional: handle focus or escape key if needed, 
// but simple Bootstrap-like modal is usually enough.
</script>

<template>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" style="background: rgba(0,0,0,0.5); z-index: 1060;">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content border-0 shadow">
        <div class="modal-header border-bottom-0 pb-0">
          <h5 class="modal-title fw-bold">{{ title }}</h5>
          <button type="button" class="btn-close" @click="emit('close')" :disabled="loading"></button>
        </div>
        <div class="modal-body py-4">
          <p class="mb-0">{{ message }}</p>
        </div>
        <div class="modal-footer border-top-0 pt-0">
          <button type="button" class="btn btn-ghost" @click="emit('close')" :disabled="loading">
            {{ cancelText || 'Cancel' }}
          </button>
          <button 
            type="button" 
            :class="['btn px-4', variant ? `btn-${variant}` : 'btn-primary']"
            @click="emit('confirm')"
            :disabled="loading"
          >
            <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
            {{ confirmText || 'Confirm' }}
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
</style>
