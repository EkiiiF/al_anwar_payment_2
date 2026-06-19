<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { HTMLButtonAttributes } from 'svelte/elements';
  import { Loader2 } from 'lucide-svelte';

  interface Props extends Omit<HTMLButtonAttributes, 'onclick'> {
    variant?: 'primary' | 'secondary' | 'outline' | 'danger' | 'ghost';
    size?: 'sm' | 'md' | 'lg';
    href?: string;
    loading?: boolean;
    class?: string;
    children: Snippet;
    onclick?: (e: MouseEvent) => void | Promise<void>;
  }

  let {
    variant = 'primary',
    size = 'md',
    href = '',
    type = 'button',
    loading = false,
    disabled = false,
    class: className = '',
    children,
    onclick,
    ...restProps
  }: Props = $props();

  const base =
    'inline-flex items-center justify-center gap-2 font-bold rounded-lg transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-800/40 focus-visible:ring-offset-2 disabled:opacity-40 disabled:cursor-not-allowed select-none cursor-pointer';

  const variants = {
    primary:   'bg-emerald-800 hover:bg-emerald-700 active:bg-emerald-900 text-white border border-transparent shadow-sm',
    secondary: 'bg-slate-600 hover:bg-slate-700 active:bg-slate-800 text-white border border-transparent shadow-sm',
    outline:   'bg-transparent border border-slate-200 text-slate-700 hover:bg-slate-50 hover:text-slate-900 active:bg-slate-100',
    danger:    'bg-red-600 hover:bg-red-700 active:bg-red-800 text-white border border-transparent shadow-sm',
    ghost:     'bg-transparent text-slate-600 hover:bg-slate-50 hover:text-slate-900 active:bg-slate-100'
  };

  const sizes = {
    sm: 'px-3 py-1.5 text-xs',
    md: 'px-4 py-2 text-sm',
    lg: 'px-5 py-2.5 text-sm'
  };

  const classes = $derived(`${base} ${variants[variant]} ${sizes[size]} ${className}`);
  const isDisabled = $derived(disabled || loading);
</script>

{#if href && !isDisabled}
  <a {href} class={classes} {...restProps}>
    {@render children()}
  </a>
{:else}
  <button {type} {onclick} disabled={isDisabled} class={classes} {...restProps}>
    {#if loading}
      <Loader2 size={16} class="animate-spin flex-shrink-0" aria-hidden="true" />
    {/if}
    {@render children()}
  </button>
{/if}
