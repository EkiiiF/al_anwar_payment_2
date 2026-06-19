<script lang="ts">
  import { Badge, Card, EmptyState } from '$lib/components';
  import { formatRupiah, getHijriMonthName } from '$lib/utils';
  import type { Invoice } from '$lib/types';

  export type SemGroup = { 
    semester: number; 
    academicYear: string; 
    invoices: Invoice[]; 
    total: number; 
    paidCount: number 
  };

  let {
    semesterInvoices = [],
    semesterGroups = [],
    selectedIds = new Set(),
    toggleSelect,
    selectAllSemester,
    deselectAllSemester,
    selectSemester
  }: {
    semesterInvoices: Invoice[];
    semesterGroups: SemGroup[];
    selectedIds: Set<string>;
    toggleSelect: (id: string) => void;
    selectAllSemester: () => void;
    deselectAllSemester: () => void;
    selectSemester: (group: SemGroup) => void;
  } = $props();

  const allSemSelected = $derived(
    semesterInvoices.length > 0 && semesterInvoices.every(i => selectedIds.has(i.id))
  );
</script>

<div class="space-y-4">
  <div class="bg-gradient-to-r from-blue-800 to-indigo-700 rounded-2xl p-4.5 text-white shadow-sm flex items-center justify-between border border-blue-900/10">
    <div class="min-w-0">
      <h2 class="font-bold text-base tracking-wide uppercase leading-tight">Syahriyyah Muhadhoroh & Semester</h2>
      <p class="text-xs text-blue-100 mt-1 leading-snug">Biaya ngaji/sekolah & biaya akademik per semester.</p>
    </div>
    {#if semesterInvoices.length > 0}
      <button
        type="button"
        onclick={allSemSelected ? deselectAllSemester : selectAllSemester}
        class="flex-shrink-0 px-3 py-1.5 rounded-lg bg-white/10 hover:bg-white/20 active:bg-white/30 text-xs font-bold transition-all border border-white/10"
      >
        {allSemSelected ? 'Batalkan Semua' : 'Pilih Semua'}
      </button>
    {/if}
  </div>

  {#if semesterGroups.length > 0}
    <div class="grid grid-cols-1 gap-3">
      {#each semesterGroups as group}
        {@const allSelected = group.invoices.every(i => selectedIds.has(i.id))}
        <button
          type="button"
          onclick={() => selectSemester(group)}
          class="text-left p-4.5 rounded-2xl border transition-all duration-200 hover:shadow-sm cursor-pointer select-none
            {allSelected
              ? 'bg-blue-50/40 border-blue-600 shadow-sm shadow-blue-50/30'
              : 'bg-white border-slate-200 hover:border-emerald-500/40 hover:bg-slate-50/30'}"
        >
          <div class="flex justify-between items-start gap-4">
            <div>
              <p class="text-xs text-slate-500 font-bold uppercase tracking-wider">
                Paket Muhadhoroh Semester {group.semester}
              </p>
              <p class="text-sm text-blue-800 font-bold mt-0.5">
                Tahun Ajaran {group.academicYear}
              </p>
            </div>
            <div class="text-right">
              <p class="text-xl font-bold text-slate-900 leading-none">{formatRupiah(group.total)}</p>
              <p class="text-xs text-slate-500 font-medium mt-1">{group.invoices.length} Tagihan Paket</p>
            </div>
          </div>
          <div class="mt-3.5 border-t border-slate-100 pt-2.5 flex justify-between items-center">
            <span class="text-xs font-bold {allSelected ? 'text-blue-800' : 'text-slate-500'} flex items-center gap-1.5">
              <span>{allSelected ? '✓ Paket Terpilih' : 'Klik untuk memilih lunas paket ini'}</span>
            </span>
          </div>
        </button>
      {/each}
    </div>
  {/if}

  {#if semesterInvoices.length === 0}
    <Card class="border-slate-200 bg-slate-50/20 !p-6">
      <EmptyState
        title="Tagihan Semester Lunas!"
        description="Alhamdulillah, tidak ada tagihan Syahriyyah Muhadhoroh atau tagihan semester yang aktif."
        class="bg-transparent"
      />
    </Card>
  {:else}
    <div class="space-y-3">
      {#each semesterInvoices as invoice (invoice.id)}
        {@const isSelected = selectedIds.has(invoice.id)}
        <div
          onclick={() => toggleSelect(invoice.id)}
          class="group relative p-4.5 rounded-2xl border transition-all duration-200 cursor-pointer select-none
            {isSelected
              ? 'bg-blue-50/40 border-blue-500 shadow-sm shadow-blue-50/30'
              : 'bg-white border-slate-200 hover:bg-slate-50 hover:border-slate-300'}"
          role="checkbox"
          aria-checked={isSelected}
          tabindex="0"
        >
          <div class="flex items-center justify-between gap-4">
            <div class="flex items-start gap-3.5 min-w-0">
              <div class="w-5 h-5 rounded-md border-2 flex-shrink-0 flex items-center justify-center transition-all duration-200 mt-0.5
                {isSelected ? 'border-blue-800 bg-blue-800' : 'border-slate-300 group-hover:border-blue-700'}"
              >
                {#if isSelected}
                  <svg class="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3.5" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                {/if}
              </div>

              <div class="min-w-0">
                <p class="text-xs text-slate-400 font-mono leading-none mb-1 mt-0.5">{invoice.invoice_number}</p>
                <p class="font-bold text-slate-900 text-base tracking-tight leading-snug">
                  {#if invoice.category}
                    <span class="text-blue-700">{invoice.category.name}</span>
                  {/if}
                  {#if invoice.hijri_month}
                    · {getHijriMonthName(invoice.hijri_month)} {invoice.hijri_year} H
                  {/if}
                </p>
                
                <div class="flex items-center gap-2 mt-1.5">
                  {#if invoice.semester}
                    <Badge label={`Sem ${invoice.semester}`} variant="warning" class="text-xs py-0.5" />
                  {/if}
                  {#if invoice.academic_year_label}
                    <span class="text-xs text-emerald-600 font-bold">TA {invoice.academic_year_label}</span>
                  {/if}
                </div>
                <p class="text-xs text-slate-500 font-medium mt-1">
                  Jatuh tempo: {new Date(invoice.due_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}
                </p>
              </div>
            </div>

            <div class="text-right flex-shrink-0 self-start">
              <p class="text-lg font-bold text-slate-900">{formatRupiah(invoice.amount_due)}</p>
              <Badge label="Semester" variant="warning" class="mt-1.5 text-xs py-0.5" />
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
