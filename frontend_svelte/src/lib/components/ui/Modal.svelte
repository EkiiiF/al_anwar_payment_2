<script lang="ts">
  import type { Snippet } from 'svelte';
  import { X } from 'lucide-svelte';

  let {
    open = $bindable(false),
    title,
    size = 'md',
    onclose,
    children,
    footer
  }: {
    open?: boolean;
    title?: string;
    size?: 'sm' | 'md' | 'lg' | 'xl';
    onclose?: () => void;
    children: Snippet;
    footer?: Snippet;
  } = $props();

  const sizes: Record<string, string> = {
    sm: 'max-w-sm',
    md: 'max-w-lg',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl'
  };

  function close() {
    open = false;
    onclose?.();
  }

  function onBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) close();
  }

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') close();
  }
</script>

<svelte:window onkeydown={onKeydown} />

{#if open}
  <div
    class="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-fade-in"
    onclick={onBackdropClick}
    role="presentation"
  >
    <div
      class="bg-white border border-slate-200 rounded-2xl w-full {sizes[size]}
             flex flex-col max-h-[90vh] transform scale-95 sm:scale-100 transition-all duration-350 ease-out animate-fade-in-up"
    >
      {#if title}
        <div class="flex items-center justify-between px-6 py-4 border-b border-slate-100 flex-shrink-0">
          <h2 class="text-base font-semibold text-slate-900">{title}</h2>
          <button
            onclick={close}
            class="p-1.5 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-50 transition-colors duration-200"
            aria-label="Tutup modal"
          >
            <X size={18} />
          </button>
        </div>
      {/if}
 
      <div class="p-6 overflow-y-auto flex-1">
        {@render children()}
      </div>
 
      {#if footer}
        <div class="px-6 py-4 border-t border-slate-100 bg-slate-50 rounded-b-2xl flex-shrink-0">
          {@render footer()}
        </div>
      {/if}
    </div>
  </div>
{/if}
