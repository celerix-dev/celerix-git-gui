type EventCallback = (data?: any) => void;
const listeners = new Map<string, EventCallback[]>();

export const eventBus = {
  on(event: string, callback: EventCallback) {
    if (!listeners.has(event)) {
      listeners.set(event, []);
    }
    listeners.get(event)?.push(callback);
  },
  off(event: string, callback: EventCallback) {
    const eventListeners = listeners.get(event);
    if (eventListeners) {
      const index = eventListeners.indexOf(callback);
      if (index !== -1) {
        eventListeners.splice(index, 1);
      }
    }
  },
  emit(event: string, data?: any) {
    listeners.get(event)?.forEach(callback => callback(data));
  }
};
