<script lang="ts">
  import { User, MapPin, ChevronRight } from 'lucide-svelte';
  import { Card, Badge } from '$lib/components';
  import type { Student } from '$lib/types';

  let {
    student,
    onclick
  }: {
    student: Student;
    onclick: () => void;
  } = $props();

  const fullName = $derived(
    [student.name?.first_name, student.name?.middle_name, student.name?.last_name]
      .filter(Boolean)
      .join(' ')
  );
</script>

<div 
  onclick={onclick}
  class="group bg-white border border-slate-200 hover:border-emerald-600/60 rounded-2xl shadow-sm hover:shadow-md transition-all duration-200 border-l-4 border-l-emerald-800 cursor-pointer p-5"
>
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div class="flex items-start gap-4">
      <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-slate-50 to-slate-100/80 border border-slate-200/50 flex items-center justify-center flex-shrink-0 shadow-sm group-hover:from-emerald-55 group-hover:to-emerald-100 group-hover:border-emerald-200 transition-colors">
        <User size={22} class="text-slate-600 group-hover:text-emerald-800 transition-colors" />
      </div>
      <div class="space-y-1 min-w-0">
        <p class="font-extrabold text-slate-900 text-base leading-tight group-hover:text-emerald-950 transition-colors truncate">
          {fullName}
        </p>
        <div class="flex flex-wrap items-center gap-2">
          <span class="text-[10px] text-slate-500 font-mono bg-slate-100 px-2 py-0.5 rounded border border-slate-200/40">
            NIS: {student.student_number}
          </span>
          {#if student.status}
            <Badge label={student.status.name} variant="info" class="text-[10px] px-2 py-0.5 rounded" />
          {/if}
        </div>
        {#if student.addresses && student.addresses.length > 0}
          <p class="text-xs text-slate-500 truncate max-w-[280px] sm:max-w-md flex items-center gap-1 mt-0.5">
            <MapPin size={12} class="text-slate-400 flex-shrink-0" />
            <span class="truncate">{student.addresses[0].address_line}</span>
          </p>
        {/if}
      </div>
    </div>

    <div class="flex sm:flex-col items-end justify-between sm:justify-start gap-2 border-t border-slate-100 pt-3 sm:pt-0 sm:border-0">
      <div class="text-left sm:text-right">
        <p class="text-xs font-semibold text-slate-500 leading-none">
          {student.invoices?.length ?? 0} Tagihan
        </p>
        <div class="flex items-center gap-1.5 justify-end mt-1">
          {#if student.invoices?.some(i => i.status === 'unpaid')}
            <Badge label="Belum Lunas" variant="warning" dot />
          {/if}
          {#if student.invoices?.every(i => i.status === 'paid') && (student.invoices?.length ?? 0) > 0}
            <Badge label="Lunas Semua" variant="success" dot />
          {/if}
        </div>
      </div>
      <p class="text-xs text-emerald-800 font-bold flex items-center gap-0.5 group-hover:translate-x-0.5 transition-transform self-end">
        <span>Detail</span>
        <ChevronRight size={14} />
      </p>
    </div>
  </div>
</div>
