<script lang="ts">
  let {
    label,
    variant = 'default',
    dot = false
  }: {
    label: string;
    variant?: 'default' | 'success' | 'warning' | 'danger' | 'info' | 'teal' | 'purple';
    dot?: boolean;
  } = $props();

  const variants: Record<string, string> = {
    default: 'bg-slate-100 text-slate-600',
    success: 'bg-emerald-50 text-emerald-800',
    warning: 'bg-amber-50 text-amber-600',
    danger:  'bg-red-50 text-red-600',
    info:    'bg-blue-100 text-blue-600',
    teal:    'bg-teal-50 text-teal-700',
    purple:  'bg-purple-50 text-purple-700'
  };

  const dotColors: Record<string, string> = {
    default: 'bg-slate-400',
    success: 'bg-emerald-600',
    warning: 'bg-amber-500',
    danger:  'bg-red-500',
    info:    'bg-blue-500',
    teal:    'bg-teal-500',
    purple:  'bg-purple-500'
  };
</script>

<span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium {variants[variant]}">
  {#if dot}
    <span class="w-1.5 h-1.5 rounded-full {dotColors[variant]}"></span>
  {/if}
  {label}
</span>
