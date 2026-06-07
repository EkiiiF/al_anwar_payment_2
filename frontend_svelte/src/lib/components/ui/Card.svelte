<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { HTMLAttributes } from 'svelte/elements';

  interface Props extends HTMLAttributes<HTMLDivElement> {
    padding?: boolean;
    title?: string;
    children: Snippet;
  }

  let {
    class: className = '',
    padding = true,
    title = '',
    children,
    ...rest
  }: Props = $props();
</script>

<div
  class="bg-white border border-slate-200/80 rounded-xl
    {padding ? 'p-6' : ''}
    {className}"
  {...rest}
>
  {#if title}
    <h3 class="text-sm font-semibold text-slate-900 mb-4">{title}</h3>
  {/if}
  {@render children()}
</div>
