import { ref } from 'vue';

export interface Alert {
  id: string;
  message: string;
  type: 'success' | 'danger' | 'warning' | 'info';
  title?: string;
}

const alerts = ref<Alert[]>([]);

export function useAlerts() {
  const showAlert = (message: string, type: Alert['type'] = 'info', title?: string) => {
    const id = Math.random().toString(36).substring(7);
    const alert: Alert = { id, message, type, title };
    alerts.value.push(alert);

    // Auto-remove after 5 seconds
    setTimeout(() => {
      removeAlert(id);
    }, 5000);
  };

  const removeAlert = (id: string) => {
    const index = alerts.value.findIndex(a => a.id === id);
    if (index !== -1) {
      alerts.value.splice(index, 1);
    }
  };

  return {
    alerts,
    showAlert,
    removeAlert,
    showError: (message: string, title: string = 'Error') => showAlert(message, 'danger', title),
    showSuccess: (message: string, title: string = 'Success') => showAlert(message, 'success', title),
    showWarning: (message: string, title: string = 'Warning') => showAlert(message, 'warning', title),
    showInfo: (message: string, title: string = 'Info') => showAlert(message, 'info', title),
  };
}
