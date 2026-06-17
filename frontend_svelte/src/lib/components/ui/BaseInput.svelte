<script lang="ts">
  import { Eye, EyeOff } from 'lucide-svelte';

  let {
    id,
    label,
    type = 'text',
    placeholder = '',
    value = $bindable(''),
    required = false,
    disabled = false,
    error = $bindable(''),
    helper = '',
    autocomplete = undefined,
    class: className = '',
    oninput = undefined
  }: {
    id: string;
    label?: string;
    type?: string;
    placeholder?: string;
    value?: string;
    required?: boolean;
    disabled?: boolean;
    error?: string;
    helper?: string;
    autocomplete?: AutoFill;
    class?: string;
    oninput?: (e: Event) => void;
  } = $props();

  let showPassword = $state(false);
  const isPassword = $derived(type === 'password');
  const resolvedType = $derived(isPassword ? (showPassword ? 'text' : 'password') : type);

  function togglePassword() {
    showPassword = !showPassword;
  }

  function handleInput(e: Event) {
    error = '';
    oninput?.(e);
  }
</script>

<div class="flex flex-col gap-1.5 {className}">
  {#if label}
    <label for={id} class="text-xs font-medium text-slate-600 uppercase tracking-wider">
      {label}
      {#if required}<span class="text-red-500 ml-0.5">*</span>{/if}
    </label>
  {/if}

  <div class="relative">
    <input
      {id}
      type={resolvedType}
      {placeholder}
      {required}
      {disabled}
      {autocomplete}
      oninput={handleInput}
      bind:value
      class="w-full px-3 py-2.5 rounded-lg bg-white border text-slate-900 text-sm outline-none transition-colors duration-200
        placeholder:text-slate-400
        focus:ring-1 focus:ring-emerald-800 focus:border-emerald-800
        disabled:opacity-40 disabled:cursor-not-allowed disabled:bg-slate-50
        {isPassword ? 'pr-10' : ''}
        {error
          ? 'border-red-400 focus:border-red-400 focus:ring-2 focus:ring-red-200'
          : 'border-slate-200'}"
    />

    {#if isPassword}
      <button
        type="button"
        onclick={togglePassword}
        class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition-colors duration-200"
        aria-label={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
        tabindex="-1"
      >
        {#if showPassword}
          <EyeOff size={16} aria-hidden="true" />
        {:else}
          <Eye size={16} aria-hidden="true" />
        {/if}
      </button>
    {/if}
  </div>

  {#if error}
    <p class="text-xs text-red-500">{error}</p>
  {:else if helper}
    <p class="text-xs text-slate-500">{helper}</p>
  {/if}
</div>
