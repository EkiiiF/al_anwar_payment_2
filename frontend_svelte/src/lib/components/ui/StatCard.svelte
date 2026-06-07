<script lang="ts">
  import type { IconComponent } from '$lib/types/icons';

  let {
    title,
    value,
    subtitle = '',
    icon: Icon,
    color = 'emerald',
    accent = false
  }: {
    title: string;
    value: string | number;
    subtitle?: string;
    icon: IconComponent;
    color?: 'emerald' | 'teal' | 'amber' | 'red' | 'blue' | 'green' | 'purple';
    accent?: boolean;
  } = $props();

  const colorMap: Record<string, { icon: string; border: string }> = {
    emerald: { icon: 'bg-emerald-50 text-emerald-800',  border: 'border-emerald-100' },
    teal:    { icon: 'bg-teal-50 text-teal-700',        border: 'border-teal-100' },
    amber:   { icon: 'bg-amber-50 text-amber-600',      border: 'border-amber-100' },
    red:     { icon: 'bg-red-50 text-red-600',           border: 'border-red-100' },
    blue:    { icon: 'bg-blue-100 text-blue-600',        border: 'border-blue-200' },
    green:   { icon: 'bg-emerald-50 text-emerald-800',  border: 'border-emerald-100' },
    purple:  { icon: 'bg-purple-50 text-purple-700',    border: 'border-purple-100' }
  };

  const c = $derived(colorMap[color] || colorMap.emerald);
</script>

<div
  class="relative overflow-hidden rounded-xl border p-5 transition-colors duration-200
    {accent
      ? 'bg-emerald-800 border-emerald-800 text-white'
      : 'bg-white border-slate-200/80 hover:border-slate-300'}"
>
  <div class="relative z-10 flex items-start gap-4">
    <div
      class="p-2.5 rounded-lg flex-shrink-0
        {accent ? 'bg-white/15 text-white' : c.icon}"
    >
      <Icon size={20} />
    </div>
    <div class="min-w-0">
      <p class="text-xs font-medium uppercase tracking-wider mb-1 {accent ? 'text-emerald-100' : 'text-slate-500'}">
        {title}
      </p>
      <p class="text-2xl font-bold leading-none {accent ? 'text-white' : 'text-slate-900'}">
        {value}
      </p>
      {#if subtitle}
        <p class="text-xs mt-1.5 {accent ? 'text-emerald-200' : 'text-slate-500'}">{subtitle}</p>
      {/if}
    </div>
  </div>
</div>
