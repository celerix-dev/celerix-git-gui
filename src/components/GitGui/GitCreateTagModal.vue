<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  show: boolean;
  fromHash: string;
  loading: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'create', data: { name: string; message: string; pushAll: boolean }): void;
}>();

const tagName = ref('');
const tagMessage = ref('');
const pushAll = ref(false);

watch(() => props.show, (newVal) => {
  if (newVal) {
    tagName.value = '';
    tagMessage.value = '';
    pushAll.value = false;
  }
});

const handleCreate = () => {
  if (tagName.value.trim()) {
    emit('create', {
      name: tagName.value.trim(),
      message: tagMessage.value.trim(),
      pushAll: pushAll.value
    });
  }
};
</script>

<template>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" style="background: rgba(0,0,0,0.5)">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content border-0 shadow">
        <div class="modal-header border-bottom-0 pb-0">
          <h5 class="modal-title fw-bold">Add Tag</h5>
          <button type="button" class="btn-close" @click="emit('close')" :disabled="loading"></button>
        </div>
        <div class="modal-body py-4">
          <div class="mb-3">
            <label class="form-label small text-muted text-uppercase fw-bold mb-1">Creating from commit</label>
            <div class="d-flex align-items-center p-2 bg-light rounded border">
              <i class="ti ti-hash me-2 text-primary"></i>
              <span class="fw-semibold font-monospace">{{ fromHash }}</span>
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

          <div class="form-check form-switch mt-4">
            <input 
              class="form-check-input" 
              type="checkbox" 
              role="switch" 
              id="pushAll" 
              v-model="pushAll"
              :disabled="loading"
            >
            <label class="form-check-label" for="pushAll">Push to all remotes</label>
          </div>
          <div class="small text-muted mt-1" v-if="!pushAll">
            If unchecked, it will only push to the default remote (origin).
          </div>
        </div>
        <div class="modal-footer border-top-0 pt-0">
          <button type="button" class="btn btn-ghost" @click="emit('close')" :disabled="loading">Cancel</button>
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
.btn-ghost {
  background: transparent;
  border: none;
  color: inherit;
}
.btn-ghost:hover {
  background: rgba(0, 0, 0, 0.05);
}
</style>
