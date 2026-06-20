<script lang="ts">
  import { Badge, Card, EmptyState, Modal, Button } from '$lib/components';
  import { formatRupiah, getHijriMonthName } from '$lib/utils';
  import type { Invoice } from '$lib/types';
  import { Calendar } from 'lucide-svelte';

  let {
    pondokInvoices = [],
    selectedIds = new Set(),
    toggleSelect,
    selectAllPondok,
    deselectAllPondok,
    onCheckout,
    paying = false,
    totalSelected = 0,
    showPondokModal = $bindable(false),
    onSwitchToSemester
  }: {
    pondokInvoices: Invoice[];
    selectedIds: Set<string>;
    toggleSelect: (id: string) => void;
    selectAllPondok: () => void;
    deselectAllPondok: () => void;
    onCheckout: () => void;
    paying: boolean;
    totalSelected: number;
    showPondokModal: boolean;
    onSwitchToSemester: () => void;
  } = $props();

  const allPondokSelected = $derived(
    pondokInvoices.length > 0 && pondokInvoices.every(i => selectedIds.has(i.id))
  );

  const selectedCount = $derived(
    pondokInvoices.filter(i => selectedIds.has(i.id)).length
  );

  const totalAmount = $derived(
    pondokInvoices.reduce((sum, i) => sum + i.amount_due, 0)
  );

  // Sort monthly invoices from earliest due date to latest
  const sortedPondokInvoices = $derived(
    [...pondokInvoices].sort((a, b) => new Date(a.due_date).getTime() - new Date(b.due_date).getTime())
  );
</script>

<div class="space-y-3">
  {#if pondokInvoices.length === 0}
    <Card class="border-slate-200 bg-slate-50/20 !p-6">
      <EmptyState
        title="Tagihan Bulanan Lunas!"
        description="Alhamdulillah, tidak ada tagihan Syahriyyah Pondok yang aktif saat ini."
        class="bg-transparent"
      />
    </Card>
  {:else}
    <button
      type="button"
      onclick={() => showPondokModal = true}
      class="w-full text-left p-4 rounded-2xl border border-slate-200 transition-all duration-200 bg-white hover:bg-slate-50/50 flex items-center justify-between gap-4 cursor-pointer shadow-sm"
    >
      <div class="flex items-center gap-3.5 min-w-0">
        <div class="w-11 h-11 rounded-xl bg-emerald-50 text-emerald-800 flex items-center justify-center flex-shrink-0 border border-emerald-100">
          <Calendar size={20} />
        </div>
        <div class="min-w-0">
          <h3 class="font-extrabold text-slate-800 text-sm sm:text-base leading-tight">Syahriyyah Pondok (Bulanan)</h3>
          <p class="text-[11px] text-slate-400 font-medium mt-1 leading-snug truncate">Iuran operasional bulanan wajib pondok</p>
          <div class="flex items-center gap-1.5 mt-2">
            <span class="w-1.5 h-1.5 rounded-full {selectedCount > 0 ? 'bg-emerald-500' : 'bg-slate-400'}"></span>
            <span class="text-[10px] text-slate-550 font-bold">
              {selectedCount > 0 ? `${selectedCount} Bulan Terpilih` : `${pondokInvoices.length} Tagihan Aktif`}
            </span>
          </div>
        </div>
      </div>

      <div class="text-right flex-shrink-0 pl-2">
        <p class="text-lg sm:text-xl font-black text-slate-900 leading-none">{formatRupiah(totalAmount)}</p>
        <span class="text-[10px] text-emerald-700 bg-emerald-50 border border-emerald-150 px-2 py-0.5 rounded font-extrabold inline-block mt-2">
          Detail →
        </span>
      </div>
    </button>
  {/if}
</div>

<Modal bind:open={showPondokModal} title="Detail Tagihan: Syahriyyah Pondok (Bulanan)" size="md">
  {#snippet footer()}
    <div class="flex gap-3 justify-between items-center w-full">
      <button
        type="button"
        onclick={onSwitchToSemester}
        class="text-xs font-bold text-blue-600 hover:text-blue-700 hover:underline cursor-pointer"
      >
        Lihat Tagihan Semester →
      </button>

      {#if selectedCount > 0}
        <Button
          type="button"
          onclick={onCheckout}
          disabled={paying}
          variant="primary"
        >
          {#snippet children()}
            {paying ? 'Memproses...' : `Bayar (${formatRupiah(totalSelected)})`}
          {/snippet}
        </Button>
      {/if}
    </div>
  {/snippet}

  <div class="space-y-2.5">
    <div class="flex justify-between items-center pb-2.5 border-b border-slate-100 flex-shrink-0">
      <span class="text-xs font-bold text-slate-500">{sortedPondokInvoices.length} Tagihan</span>
      <button
        type="button"
        onclick={allPondokSelected ? deselectAllPondok : selectAllPondok}
        class="text-xs font-bold text-emerald-700 hover:text-emerald-800 hover:underline cursor-pointer"
      >
        {allPondokSelected ? 'Batalkan Semua' : 'Pilih Semua'}
      </button>
    </div>

    {#each sortedPondokInvoices as invoice (invoice.id)}
      {@const isSelected = selectedIds.has(invoice.id)}
      <div
        onclick={() => toggleSelect(invoice.id)}
        onkeydown={(e) => {
          if (e.key === ' ' || e.key === 'Enter') {
            e.preventDefault();
            toggleSelect(invoice.id);
          }
        }}
        class="group relative p-3 rounded-xl border transition-all duration-200 cursor-pointer select-none
          {isSelected
            ? 'bg-emerald-50/40 border-emerald-600 shadow-sm shadow-emerald-50/30'
            : 'bg-white border-slate-200 hover:bg-slate-50 hover:border-slate-300'}"
        role="checkbox"
        aria-checked={isSelected}
        tabindex="0"
      >
        <div class="flex items-center justify-between gap-3 w-full">
          <div class="flex items-center gap-3 min-w-0">
            <div class="w-5 h-5 rounded-md border-2 flex-shrink-0 flex items-center justify-center transition-all duration-200
              {isSelected ? 'border-emerald-800 bg-emerald-800' : 'border-slate-300 group-hover:border-emerald-700'}"
            >
              {#if isSelected}
                <svg class="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3.5" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
              {/if}
            </div>

            <div class="min-w-0">
              <p class="font-bold text-slate-800 text-sm sm:text-base leading-tight">
                {getHijriMonthName(invoice.hijri_month)} {invoice.hijri_year} H
              </p>
              <p class="text-[10px] text-slate-400 font-medium mt-0.5">
                Jatuh tempo: {new Date(invoice.due_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}
              </p>
            </div>
          </div>

          <div class="text-right flex-shrink-0">
            <p class="text-sm sm:text-base font-extrabold text-emerald-800">{formatRupiah(invoice.amount_due)}</p>
          </div>
        </div>
      </div>
    {/each}
  </div>
</Modal>
