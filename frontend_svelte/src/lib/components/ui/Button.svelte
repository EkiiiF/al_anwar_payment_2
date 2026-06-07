<script lang="ts">
  import type { Snippet } from 'svelte';
  import { Loader2 } from 'lucide-svelte';

  let {
    variant = 'primary',
    size = 'md',
    href = '',
    type = 'button',
    form = '',
    onclick,
    disabled = false,
    loading = false,
    class: className = '',
    children
  }: {
    variant?: 'primary' | 'secondary' | 'outline' | 'danger' | 'ghost';
    size?: 'sm' | 'md' | 'lg';
    href?: string;
    type?: 'button' | 'submit' | 'reset';
    form?: string;
    onclick?: (e: MouseEvent) => void;
    disabled?: boolean;
    loading?: boolean;
    class?: string;
    children: Snippet;
  } = $props();

  const base =
    'inline-flex items-center justify-center gap-2 font-semibold rounded-lg transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-800/40 focus-visible:ring-offset-2 disabled:opacity-40 disabled:cursor-not-allowed select-none cursor-pointer';

  const variants: Record<string, string> = {
    primary:   'bg-emerald-800 hover:bg-emerald-700 active:bg-emerald-900 text-white',
    secondary: 'bg-teal-600 hover:bg-teal-700 active:bg-teal-800 text-white',
    outline:   'bg-transparent border border-slate-200 text-slate-700 hover:bg-slate-50 hover:text-slate-900 active:bg-slate-100',
    danger:    'bg-red-600 hover:bg-red-700 active:bg-red-800 text-white',
    ghost:     'bg-transparent text-slate-600 hover:bg-slate-50 hover:text-slate-900 active:bg-slate-100'
  };

  const sizes: Record<string, string> = {
    sm: 'px-3 py-1.5 text-xs',
    md: 'px-4 py-2 text-sm',
    lg: 'px-5 py-2.5 text-sm'
  };

  const classes = $derived(`${base} ${variants[variant]} ${sizes[size]} ${className}`);
  const isDisabled = $derived(disabled || loading);
</script>

{#if href && !isDisabled}
  <a {href} class={classes}>
    {@render children()}
  </a>
{:else}
  <button {type} form={form || undefined} {onclick} disabled={isDisabled} class={classes}>
    {#if loading}
      <Loader2 size={14} class="animate-spin" aria-hidden="true" />
    {/if}
    {@render children()}
  </button>
{/if}
