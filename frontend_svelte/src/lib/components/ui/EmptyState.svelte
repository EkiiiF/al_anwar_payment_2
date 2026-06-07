<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Inbox } from 'lucide-svelte';

  let {
    title = 'Tidak ada data',
    description = '',
    action
  }: {
    title?: string;
    description?: string;
    action?: Snippet;
  } = $props();
</script>

<div class="flex flex-col items-center justify-center py-14 px-4 text-center">
  <div class="p-4 rounded-full bg-slate-50 border border-slate-200/80 mb-4">
    <Inbox size={28} class="text-slate-400" />
  </div>
  <h3 class="text-sm font-semibold text-slate-700 mb-1">{title}</h3>
  {#if description}
    <p class="text-sm text-slate-500 max-w-xs">{description}</p>
  {/if}
  {#if action}
    <div class="mt-5">
      {@render action()}
    </div>
  {/if}
</div>
