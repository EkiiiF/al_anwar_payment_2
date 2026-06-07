<script lang="ts">
  import { ChevronDown } from 'lucide-svelte';

  let {
    id,
    label,
    value = $bindable(''),
    required = false,
    disabled = false,
    error = '',
    placeholder = 'Pilih...',
    options = [],
    class: className = '',
    onchange
  }: {
    id: string;
    label?: string;
    value?: string;
    required?: boolean;
    disabled?: boolean;
    error?: string;
    placeholder?: string;
    options: { value: string | number; label: string }[];
    class?: string;
    onchange?: () => void;
  } = $props();
</script>

<div class="flex flex-col gap-1.5 {className}">
  {#if label}
    <label for={id} class="text-xs font-medium text-slate-600 uppercase tracking-wider">
      {label}
      {#if required}<span class="text-red-500 ml-0.5">*</span>{/if}
    </label>
  {/if}

  <div class="relative">
    <select
      {id}
      {required}
      {disabled}
      {onchange}
      bind:value
      class="w-full px-3 py-2.5 pr-9 rounded-lg bg-white border text-slate-900 text-sm outline-none
        appearance-none cursor-pointer transition-colors duration-200
        focus:ring-1 focus:ring-emerald-800 focus:border-emerald-800
        disabled:opacity-40 disabled:cursor-not-allowed disabled:bg-slate-50
        {error ? 'border-red-500' : 'border-slate-200'}"
    >
      <option value="" disabled class="text-slate-400">{placeholder}</option>
      {#each options as opt}
        <option value={String(opt.value)}>{opt.label}</option>
      {/each}
    </select>
    <ChevronDown
      size={15}
      class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none"
      aria-hidden="true"
    />
  </div>

  {#if error}
    <p class="text-xs text-red-500">{error}</p>
  {/if}
</div>
