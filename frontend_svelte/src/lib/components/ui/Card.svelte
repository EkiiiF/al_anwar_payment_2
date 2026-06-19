<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { HTMLAttributes } from 'svelte/elements';

  interface Props extends HTMLAttributes<HTMLDivElement> {
    padding?: boolean;
    title?: string;
    class?: string;
    children: Snippet;
  }

  let {
    class: className = '',
    padding = true,
    title = '',
    children,
    ...restProps
  }: Props = $props();
</script>

<div
  class="bg-white border border-slate-200 rounded-xl shadow-sm transition-all duration-200
    {padding ? 'p-3.5 sm:p-4' : ''}
    {className}"
  {...restProps}
>
  {#if title}
    <h3 class="text-sm font-bold text-slate-900 mb-4">{title}</h3>
  {/if}
  {@render children()}
</div>
