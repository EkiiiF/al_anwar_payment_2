<script lang="ts">
  import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-svelte';

  interface Props {
    page: number;
    limit: number;
    total: number;
    pages: number;
    amounts?: number[];
    label?: string;
    onPageChange: (page: number) => void;
    onLimitChange?: (limit: number) => void;
  }

  let {
    page = 1,
    limit = 10,
    total = 0,
    pages = 1,
    amounts = [5, 10, 20, 50],
    label = 'data',
    onPageChange,
    onLimitChange
  }: Props = $props();

  const showStart = $derived(total === 0 ? 0 : ((page - 1) * limit) + 1);
  const showEnd = $derived(Math.min(page * limit, total));

  function getPageNumbers(current: number, totalPages: number) {
    const list: (number | string)[] = [];
    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) list.push(i);
    } else {
      list.push(1);
      if (current > 4) {
        list.push('...');
      }
      const start = Math.max(2, current - 2);
      const end = Math.min(totalPages - 1, current + 2);
      
      let adjustedStart = start;
      let adjustedEnd = end;
      if (current <= 4) {
        adjustedEnd = 5;
      } else if (current >= totalPages - 3) {
        adjustedStart = totalPages - 4;
      }
      
      for (let i = adjustedStart; i <= adjustedEnd; i++) {
        list.push(i);
      }
      
      if (current < totalPages - 3) {
        list.push('...');
      }
      list.push(totalPages);
    }
    return list;
  }

  const pageNumbers = $derived(getPageNumbers(page, pages));
</script>

<div class="flex flex-col sm:flex-row items-center justify-between gap-4 py-4 px-5 border-t border-slate-200 bg-white flex-shrink-0">
  <!-- Left side: Show amounts select + info text -->
  <div class="flex items-center gap-4 flex-wrap text-sm text-slate-600">
    {#if onLimitChange}
      <div class="flex items-center gap-2">
        <span>Tampilkan</span>
        <select
          value={limit}
          onchange={(e) => {
            const newLimit = Number((e.target as HTMLSelectElement).value);
            onLimitChange?.(newLimit);
          }}
          class="px-2 py-1.5 border border-slate-200 rounded-lg bg-white focus:ring-2 focus:ring-green-500 focus:border-green-500 outline-none text-xs font-semibold cursor-pointer"
        >
          {#each amounts as amt}
            <option value={amt}>{amt}</option>
          {/each}
        </select>
        <span>data</span>
      </div>
      <span class="text-slate-400">|</span>
    {/if}
    <p>
      Menampilkan <span class="font-bold text-slate-800">{showStart} - {showEnd}</span> dari <span class="font-bold text-slate-800">{total}</span> {label}
    </p>
  </div>

  <!-- Right side: Navigation buttons (First, Prev, Numerals, Next, Last) -->
  <div class="flex items-center gap-1.5 flex-wrap">
    <!-- First Button -->
    <button
      type="button"
      onclick={() => onPageChange(1)}
      disabled={page <= 1}
      class="w-9 h-9 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
      title="Halaman Pertama"
    >
      <ChevronsLeft size={16} />
    </button>

    <!-- Previous Button -->
    <button
      type="button"
      onclick={() => onPageChange(page - 1)}
      disabled={page <= 1}
      class="w-9 h-9 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
      title="Sebelumnya"
    >
      <ChevronLeft size={16} />
    </button>

    <!-- Numerals -->
    {#each pageNumbers as p}
      {#if p === '...'}
        <span class="w-9 h-9 flex items-center justify-center text-slate-400 text-sm">...</span>
      {:else}
        <button
          type="button"
          onclick={() => onPageChange(Number(p))}
          class="w-9 h-9 flex items-center justify-center rounded-lg border text-sm font-semibold transition-colors cursor-pointer
            {page === p 
              ? 'bg-emerald-800 border-emerald-800 text-white shadow-sm hover:bg-emerald-700' 
              : 'border-slate-200 text-slate-700 hover:bg-slate-50'}"
        >
          {p}
        </button>
      {/if}
    {/each}

    <!-- Next Button -->
    <button
      type="button"
      onclick={() => onPageChange(page + 1)}
      disabled={page >= pages}
      class="w-9 h-9 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
      title="Selanjutnya"
    >
      <ChevronRight size={16} />
    </button>

    <!-- Last Button -->
    <button
      type="button"
      onclick={() => onPageChange(pages)}
      disabled={page >= pages}
      class="w-9 h-9 flex items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
      title="Halaman Terakhir"
    >
      <ChevronsRight size={16} />
    </button>
  </div>
</div>
