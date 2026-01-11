<script setup lang="ts">
import { ref, watch } from 'vue';

const props = defineProps<{
  stagedCount: number;
  isAmend: boolean;
  initialSubject?: string;
  initialDescription?: string;
}>();

const emit = defineEmits<{
  (e: 'commit', data: { subject: string; description: string; amend: boolean }): void;
  (e: 'update:isAmend', value: boolean): void;
}>();

const commitSubject = ref(props.initialSubject || '');
const commitDescription = ref(props.initialDescription || '');
const amend = ref(props.isAmend);

watch(() => props.initialSubject, (val) => {
    if (val !== undefined) commitSubject.value = val;
});

watch(() => props.initialDescription, (val) => {
    if (val !== undefined) commitDescription.value = val;
});

watch(amend, (val) => {
    emit('update:isAmend', val);
});

watch(() => props.isAmend, (val) => {
    amend.value = val;
});

const handleCommit = () => {
    emit('commit', {
        subject: commitSubject.value,
        description: commitDescription.value,
        amend: amend.value
    });
};

const clearInputs = () => {
    commitSubject.value = '';
    commitDescription.value = '';
    amend.value = false;
};

defineExpose({ clearInputs });
</script>

<template>
  <div class="commit-section border-top p-3 bg-body-tertiary">
    <div class="mb-2">
      <input 
        v-model="commitSubject" 
        type="text" 
        class="form-control form-control-sm bg-body" 
        placeholder="Commit subject"
      />
    </div>
    <div class="mb-2">
      <textarea 
        v-model="commitDescription" 
        class="form-control form-control-sm bg-body" 
        rows="3" 
        placeholder="Description (optional)"
      ></textarea>
    </div>
    <div class="d-flex align-items-center justify-content-between">
      <div class="form-check">
        <input v-model="amend" class="form-check-input" type="checkbox" id="amendCheck">
        <label class="form-check-label small" for="amendCheck">
          Amend
        </label>
      </div>
      <button 
        class="btn btn-primary btn-sm px-4" 
        :disabled="!commitSubject.trim() || (stagedCount === 0 && !amend)"
        @click="handleCommit"
      >
        Commit
      </button>
    </div>
  </div>
</template>

<style scoped>
.commit-section .form-control {
  border-color: var(--bs-border-color);
}
</style>
