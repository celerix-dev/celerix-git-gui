<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  show: boolean;
  title: string;
  label: string;
  defaultValue?: string;
  placeholder?: string;
  confirmText?: string;
  loading?: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'confirm', value: string): void;
}>();

const inputValue = ref('');

watch(() => props.show, (newVal) => {
  if (newVal) {
    inputValue.value = props.defaultValue || '';
  }
});

const handleConfirm = () => {
  if (inputValue.value.trim()) {
    emit('confirm', inputValue.value.trim());
  }
};
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
          <div class="mb-3">
            <label class="form-label small text-muted text-uppercase fw-bold mb-1">{{ label }}</label>
            <input 
              type="text" 
              class="form-control" 
              v-model="inputValue" 
              :placeholder="placeholder"
              @keyup.enter="handleConfirm"
              :disabled="loading"
            />
          </div>
        </div>
        <div class="modal-footer border-top-0 pt-0">
          <button type="button" class="btn btn-ghost" @click="emit('close')" :disabled="loading">Cancel</button>
          <button 
            type="button" 
            class="btn btn-primary px-4" 
            :disabled="!inputValue.trim() || loading"
            @click="handleConfirm"
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
