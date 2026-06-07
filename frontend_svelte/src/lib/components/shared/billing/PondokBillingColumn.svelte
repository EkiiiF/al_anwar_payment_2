<script lang="ts">
  import { Badge, Card, EmptyState } from '$lib/components';
  import { formatRupiah, getHijriMonthName } from '$lib/utils';
  import type { Invoice } from '$lib/types';

  let {
    pondokInvoices = [],
    selectedIds = new Set(),
    toggleSelect,
    selectAllPondok,
    deselectAllPondok
  }: {
    pondokInvoices: Invoice[];
    selectedIds: Set<string>;
    toggleSelect: (id: string) => void;
    selectAllPondok: () => void;
    deselectAllPondok: () => void;
  } = $props();

  const allPondokSelected = $derived(
    pondokInvoices.length > 0 && pondokInvoices.every(i => selectedIds.has(i.id))
  );
</script>

<div class="space-y-4">
  <div class="bg-gradient-to-r from-emerald-800 to-teal-700 rounded-2xl p-4.5 text-white shadow-sm flex items-center justify-between border border-emerald-900/10">
    <div class="min-w-0">
      <h2 class="font-extrabold text-sm tracking-wide uppercase leading-tight">Syahriyyah Pondok (Bulanan)</h2>
      <p class="text-[10px] text-emerald-100 mt-1 leading-snug">Iuran operasional bulanan wajib pondok pesantren.</p>
    </div>
    {#if pondokInvoices.length > 0}
      <button
        type="button"
        onclick={allPondokSelected ? deselectAllPondok : selectAllPondok}
        class="flex-shrink-0 px-3 py-1.5 rounded-lg bg-white/10 hover:bg-white/20 active:bg-white/30 text-xs font-bold transition-all border border-white/10"
      >
        {allPondokSelected ? 'Batalkan Semua' : 'Pilih Semua'}
      </button>
    {/if}
  </div>

  {#if pondokInvoices.length === 0}
    <Card class="border-slate-150 bg-slate-50/20 !p-6">
      <EmptyState
        title="Tagihan Bulanan Lunas!"
        description="Alhamdulillah, tidak ada tagihan Syahriyyah Pondok yang aktif saat ini."
        class="bg-transparent"
      />
    </Card>
  {:else}
    <div class="space-y-3">
      {#each pondokInvoices as invoice (invoice.id)}
        {@const isSelected = selectedIds.has(invoice.id)}
        <div
          onclick={() => toggleSelect(invoice.id)}
          class="group relative p-4.5 rounded-2xl border transition-all duration-200 cursor-pointer select-none
            {isSelected
              ? 'bg-emerald-50/40 border-emerald-600 shadow-md shadow-emerald-50/50'
              : 'bg-white border-slate-200 hover:bg-slate-50 hover:border-slate-300'}"
          role="checkbox"
          aria-checked={isSelected}
          tabindex="0"
        >
          <div class="flex items-center justify-between gap-4">
            <div class="flex items-center gap-3.5 min-w-0">
              <div class="w-5 h-5 rounded-md border-2 flex-shrink-0 flex items-center justify-center transition-all duration-200
                {isSelected ? 'border-emerald-800 bg-emerald-800' : 'border-slate-350 group-hover:border-emerald-700'}"
              >
                {#if isSelected}
                  <svg class="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3.5" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                {/if}
              </div>

              <div class="min-w-0">
                <p class="text-[9px] text-slate-400 font-mono leading-none mb-1 mt-0.5">{invoice.invoice_number}</p>
                <p class="font-extrabold text-slate-900 text-sm tracking-tight leading-snug">
                  {getHijriMonthName(invoice.hijri_month)} {invoice.hijri_year} H
                </p>
                <p class="text-[10px] text-slate-500 font-semibold mt-1">
                  Jatuh tempo: {new Date(invoice.due_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}
                </p>
              </div>
            </div>

            <div class="text-right flex-shrink-0">
              <p class="text-base font-black text-emerald-850">{formatRupiah(invoice.amount_due)}</p>
              <Badge label="Bulanan" variant="info" class="mt-1.5 text-[9px] py-0.5" />
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
