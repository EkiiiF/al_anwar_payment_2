<script lang="ts">
  import type { Snippet } from 'svelte';
  import { AlertCircle, CheckCircle2, Info, AlertTriangle } from 'lucide-svelte';

  let {
    type = 'error',
    message = '',
    class: className = '',
    children
  }: {
    type?: 'error' | 'success' | 'warning' | 'info';
    message?: string;
    class?: string;
    children?: Snippet;
  } = $props();

  const styles: Record<string, { bg: string; border: string; text: string; icon: typeof AlertCircle }> = {
    error:   { bg: 'bg-red-50',       border: 'border-red-200',     text: 'text-red-600',     icon: AlertCircle },
    success: { bg: 'bg-emerald-50',   border: 'border-emerald-200', text: 'text-emerald-800', icon: CheckCircle2 },
    warning: { bg: 'bg-amber-50',     border: 'border-amber-200',   text: 'text-amber-600',   icon: AlertTriangle },
    info:    { bg: 'bg-blue-100',     border: 'border-blue-200',    text: 'text-blue-600',    icon: Info }
  };

  const s = $derived(styles[type]);
</script>

{#if message || children}
  <div class="flex items-start gap-3 {s.bg} {s.border} border {s.text} p-4 rounded-lg text-sm {className}">
    <s.icon size={16} class="flex-shrink-0 mt-0.5" />
    <div class="flex-1 min-w-0">
      {#if children}
        {@render children()}
      {:else}
        <p>{message}</p>
      {/if}
    </div>
  </div>
{/if}
