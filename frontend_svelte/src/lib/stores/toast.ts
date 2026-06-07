import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'warning' | 'info';

export interface ToastItem {
  id: number;
  type: ToastType;
  message: string;
  duration?: number;
}

function createToastStore() {
  const { subscribe, update } = writable<ToastItem[]>([]);
  let nextId = 0;

  function add(type: ToastType, message: string, duration = 3500): void {
    const id = nextId++;
    update((items) => [...items, { id, type, message, duration }]);
    setTimeout(() => remove(id), duration);
  }

  function remove(id: number): void {
    update((items) => items.filter((t) => t.id !== id));
  }

  return {
    subscribe,
    success: (msg: string, duration?: number) => add('success', msg, duration),
    error:   (msg: string, duration?: number) => add('error',   msg, duration ?? 5000),
    warning: (msg: string, duration?: number) => add('warning', msg, duration),
    info:    (msg: string, duration?: number) => add('info',    msg, duration),
    remove
  };
}

export const toast = createToastStore();
