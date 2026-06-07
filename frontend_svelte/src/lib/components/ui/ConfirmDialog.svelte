<script lang="ts">
  import { AlertTriangle, Trash2, LogOut, X } from 'lucide-svelte';
  import Button from './Button.svelte';

  let {
    open        = $bindable(false),
    title       = 'Konfirmasi Tindakan',
    message     = 'Apakah Anda yakin ingin melanjutkan?',
    confirmText = 'Ya, Lanjutkan',
    cancelText  = 'Batal',
    variant     = 'danger' as 'danger' | 'warning' | 'info',
    onConfirm,
    loading     = false
  }: {
    open?: boolean;
    title?: string;
    message?: string;
    confirmText?: string;
    cancelText?: string;
    variant?: 'danger' | 'warning' | 'info';
    onConfirm: () => void | Promise<void>;
    loading?: boolean;
  } = $props();

  function close() { open = false; }

  const variantMap = {
    danger:  {
      icon:    Trash2,
      iconBg:  'bg-red-50',
      iconCls: 'text-red-600',
      btn:     'danger' as const
    },
    warning: {
      icon:    LogOut,
      iconBg:  'bg-amber-50',
      iconCls: 'text-amber-600',
      btn:     'primary' as const
    },
    info: {
      icon:    AlertTriangle,
      iconBg:  'bg-blue-100',
      iconCls: 'text-blue-600',
      btn:     'primary' as const
    }
  };

  const v = $derived(variantMap[variant]);
</script>

{#if open}
  <div
    class="fixed inset-0 bg-black/20 backdrop-blur-[2px] z-[9998] flex items-center justify-center p-4 animate-fade-in"
    onclick={(e) => { if (e.target === e.currentTarget) close(); }}
    aria-hidden="true"
  ></div>

  <div
    class="fixed inset-0 z-[9999] flex items-center justify-center p-4 pointer-events-none"
    role="dialog"
    aria-modal="true"
    aria-labelledby="confirm-title"
  >
    <div
      class="pointer-events-auto w-full max-w-md bg-white rounded-xl border border-slate-200/80 animate-fade-in-up overflow-hidden"
    >
      <div class="flex items-start gap-4 p-6">
        <div class="p-3 rounded-xl flex-shrink-0 {v.iconBg}">
          <v.icon size={22} class={v.iconCls} aria-hidden="true" />
        </div>
        <div class="min-w-0 flex-1">
          <h2 id="confirm-title" class="font-semibold text-slate-900 text-base">{title}</h2>
          <p class="text-sm text-slate-600 mt-1 leading-relaxed">{message}</p>
        </div>
        <button
          onclick={close}
          class="p-1.5 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-50 transition-colors duration-200 flex-shrink-0"
          aria-label="Tutup dialog"
        >
          <X size={16} />
        </button>
      </div>

      <div class="flex justify-end gap-2 px-6 pb-6">
        <Button onclick={close} variant="outline" size="md" disabled={loading}>
          {#snippet children()}{cancelText}{/snippet}
        </Button>
        <Button
          onclick={async () => { await onConfirm(); }}
          variant={variant === 'danger' ? 'danger' : 'primary'}
          size="md"
          {loading}
        >
          {#snippet children()}{confirmText}{/snippet}
        </Button>
      </div>
    </div>
  </div>
{/if}
