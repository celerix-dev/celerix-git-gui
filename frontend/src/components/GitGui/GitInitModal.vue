<script setup lang="ts">
defineProps<{
  show: boolean;
  repoPath: string | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'initialize'): void;
}>();
</script>

<template>
  <div v-if="show" class="modal-backdrop fade show"></div>
  <div v-if="show" class="modal fade show d-block" tabindex="-1" role="dialog" style="background: rgba(0,0,0,0.5)">
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content shadow-lg border-0">
        <div class="modal-header">
          <h5 class="modal-title">Not a Git Repository</h5>
          <button type="button" class="btn-close" @click="emit('close')" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <p>The directory <code>{{ repoPath }}</code> is not a git repository.</p>
          <p>Would you like to initialize a new git repository in this directory?</p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="emit('close')">Cancel</button>
          <button type="button" class="btn btn-primary" @click="emit('initialize')">Initialize and Open</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>
