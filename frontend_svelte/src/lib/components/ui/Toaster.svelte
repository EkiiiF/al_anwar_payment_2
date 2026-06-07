<script lang="ts">
  import { toast, type ToastItem } from '$lib/stores/toast';
  import { CheckCircle2, XCircle, AlertTriangle, Info, X } from 'lucide-svelte';

  const iconMap = {
    success: CheckCircle2,
    error:   XCircle,
    warning: AlertTriangle,
    info:    Info
  } as const;

  const colorMap = {
    success: {
      wrapper: 'bg-white border-emerald-200',
      icon:    'text-emerald-700',
      bar:     'bg-emerald-600',
      title:   'text-emerald-800'
    },
    error: {
      wrapper: 'bg-white border-red-200',
      icon:    'text-red-600',
      bar:     'bg-red-500',
      title:   'text-red-600'
    },
    warning: {
      wrapper: 'bg-white border-amber-200',
      icon:    'text-amber-600',
      bar:     'bg-amber-500',
      title:   'text-amber-600'
    },
    info: {
      wrapper: 'bg-white border-blue-200',
      icon:    'text-blue-600',
      bar:     'bg-blue-500',
      title:   'text-blue-600'
    }
  } as const;

  const labelMap = {
    success: 'Berhasil',
    error:   'Terjadi Kesalahan',
    warning: 'Perhatian',
    info:    'Informasi'
  } as const;

  let items: ToastItem[] = $state([]);
  toast.subscribe((v) => { items = v; });
</script>

<div
  class="fixed bottom-5 right-5 z-[9999] flex flex-col gap-3 pointer-events-none"
  aria-live="polite"
  aria-atomic="false"
  role="region"
  aria-label="Notifikasi"
>
  {#each items as item (item.id)}
    {@const c = colorMap[item.type]}
    {@const Icon = iconMap[item.type]}
    <div
      class="pointer-events-auto w-80 max-w-[calc(100vw-2.5rem)] rounded-xl border {c.wrapper} animate-slide-in-right overflow-hidden"
      role="alert"
    >
      <div
        class="h-0.5 {c.bar} animate-shrink-progress"
        style="animation-duration: {item.duration ?? 3500}ms"
      ></div>

      <div class="flex items-start gap-3 p-3.5">
        <Icon size={18} class="flex-shrink-0 mt-0.5 {c.icon}" aria-hidden="true" />
        <div class="min-w-0 flex-1">
          <p class="text-xs font-semibold uppercase tracking-wide {c.title}">{labelMap[item.type]}</p>
          <p class="text-sm text-slate-600 mt-0.5 leading-snug">{item.message}</p>
        </div>
        <button
          onclick={() => toast.remove(item.id)}
          class="flex-shrink-0 p-1 rounded-md text-slate-400 hover:text-slate-600 hover:bg-slate-50 transition-colors duration-200"
          aria-label="Tutup notifikasi"
        >
          <X size={14} aria-hidden="true" />
        </button>
      </div>
    </div>
  {/each}
</div>

<style>
  @keyframes slide-in-right {
    from { opacity: 0; transform: translateX(1rem); }
    to   { opacity: 1; transform: translateX(0); }
  }
  @keyframes shrink-progress {
    from { width: 100%; }
    to   { width: 0%; }
  }
  .animate-slide-in-right {
    animation: slide-in-right 0.25s ease-out forwards;
  }
  .animate-shrink-progress {
    animation: shrink-progress linear forwards;
  }
</style>
