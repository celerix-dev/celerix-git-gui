<script setup lang="ts">
import { useAlerts } from '@/composables/useAlerts';

const { alerts, removeAlert } = useAlerts();

const getIcon = (type: string) => {
  switch (type) {
    case 'success': return 'ti-circle-check';
    case 'danger': return 'ti-alert-circle';
    case 'warning': return 'ti-alert-triangle';
    case 'info': return 'ti-info-circle';
    default: return 'ti-info-circle';
  }
};
</script>

<template>
  <div class="alert-container position-fixed bottom-0 end-0 p-3" style="z-index: 2000;">
    <TransitionGroup name="alert-fade">
      <div 
        v-for="alert in alerts" 
        :key="alert.id" 
        :class="['alert alert-dismissible fade show shadow-sm border-0 d-flex align-items-start mb-2', `alert-${alert.type}`]"
        role="alert"
        style="min-width: 300px; max-width: 450px;"
      >
        <div class="alert-icon me-3">
          <i :class="['ti fs-4', getIcon(alert.type)]"></i>
        </div>
        <div class="flex-grow-1">
          <strong v-if="alert.title" class="d-block mb-1">{{ alert.title }}</strong>
          <div class="alert-message small">{{ alert.message }}</div>
        </div>
        <button type="button" class="btn-close ms-2" @click="removeAlert(alert.id)" aria-label="Close"></button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.alert-container {
  pointer-events: none;
}

.alert {
  pointer-events: auto;
}

.alert-fade-enter-active,
.alert-fade-leave-active {
  transition: all 0.3s ease;
}

.alert-fade-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.alert-fade-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
