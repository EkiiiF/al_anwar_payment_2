<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    padding = true,
    class: className = '',
    header,
    children
  }: {
    padding?: boolean;
    class?: string;
    header?: Snippet;
    children: Snippet;
  } = $props();
</script>

<div
  class="bg-white border border-slate-200/80 rounded-xl
    {padding ? 'p-6' : ''}
    {className}"
>
  {#if header}
    <div class="mb-4">
      {@render header()}
    </div>
  {/if}
  {@render children()}
</div>
