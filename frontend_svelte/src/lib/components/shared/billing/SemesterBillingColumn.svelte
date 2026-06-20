<script lang="ts">
  import { Badge, Card, EmptyState, Modal, Button } from '$lib/components';
  import { formatRupiah, getHijriMonthName } from '$lib/utils';
  import type { Invoice } from '$lib/types';
  import { Layers } from 'lucide-svelte';

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
    selectSemester,
    onCheckout,
    paying = false,
    totalSelected = 0,
    activeGroup = $bindable(null),
    showModal = $bindable(false),
    onSwitchToPond
  }: {
    semesterInvoices: Invoice[];
    semesterGroups: SemGroup[];
    selectedIds: Set<string>;
    toggleSelect: (id: string) => void;
    selectAllSemester: () => void;
    deselectAllSemester: () => void;
    selectSemester: (group: SemGroup) => void;
    onCheckout: () => void;
    paying: boolean;
    totalSelected: number;
    activeGroup: SemGroup | null;
    showModal: boolean;
    onSwitchToPond: () => void;
  } = $props();

  const activeGroupAllSelected = $derived(
    activeGroup && activeGroup.invoices.every(i => selectedIds.has(i.id))
  );

  const selectedInActiveGroupCount = $derived(
    activeGroup 
      ? activeGroup.invoices.filter(i => selectedIds.has(i.id)).length 
      : 0
  );

  const sortedGroupInvoices = $derived(
    activeGroup 
      ? [...activeGroup.invoices].sort((a, b) => new Date(a.due_date).getTime() - new Date(b.due_date).getTime())
      : []
  );
</script>

<div class="space-y-3">
  {#if semesterGroups.length > 0}
    <div class="space-y-3">
      {#each semesterGroups as group}
        {@const selectedInGroupCount = group.invoices.filter(i => selectedIds.has(i.id)).length}
        {@const allSelected = selectedInGroupCount === group.invoices.length}
        <button
          type="button"
          onclick={() => { activeGroup = group; showModal = true; }}
          class="w-full text-left p-4 rounded-2xl border border-slate-200 transition-all duration-200 bg-white hover:bg-slate-50/50 flex items-center justify-between gap-4 cursor-pointer shadow-sm"
        >
          <div class="flex items-center gap-3.5 min-w-0">
            <div class="w-11 h-11 rounded-xl bg-purple-50 text-purple-800 flex items-center justify-center flex-shrink-0 border border-purple-100">
              <Layers size={20} />
            </div>
            <div class="min-w-0">
              <h3 class="font-extrabold text-slate-800 text-sm sm:text-base leading-tight">Paket Semester {group.semester}</h3>
              <p class="text-[11px] text-slate-400 font-medium mt-1 leading-snug truncate">Tahun Ajaran {group.academicYear}</p>
              <div class="flex items-center gap-1.5 mt-2">
                <span class="w-1.5 h-1.5 rounded-full {selectedInGroupCount > 0 ? 'bg-purple-500' : 'bg-slate-400'}"></span>
                <span class="text-[10px] text-slate-550 font-bold">
                  {selectedInGroupCount > 0 ? `${selectedInGroupCount} dari ${group.invoices.length} Terpilih` : `${group.invoices.length} Tagihan Paket`}
                </span>
              </div>
            </div>
          </div>

          <div class="text-right flex-shrink-0 pl-2">
            <p class="text-lg sm:text-xl font-black text-slate-900 leading-none">{formatRupiah(group.total)}</p>
            <span class="text-[10px] text-purple-700 bg-purple-50 border border-purple-150 px-2 py-0.5 rounded font-extrabold inline-block mt-2">
              Detail →
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
  {/if}
</div>

{#if activeGroup}
  <Modal 
    bind:open={showModal} 
    title={`Detail Tagihan: Paket Muhadhoroh Semester ${activeGroup.semester}`} 
    size="md"
    onclose={() => activeGroup = null}
  >
    {#snippet footer()}
      <div class="flex gap-3 justify-between items-center w-full">
        <button
          type="button"
          onclick={onSwitchToPond}
          class="text-xs font-bold text-emerald-600 hover:text-emerald-700 hover:underline cursor-pointer"
        >
          ← Lihat Tagihan Bulanan (Pondok)
        </button>

        {#if selectedInActiveGroupCount > 0}
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
        <span class="text-xs font-bold text-slate-500">{sortedGroupInvoices.length} Tagihan</span>
        <button
          type="button"
          onclick={() => selectSemester(activeGroup!)}
          class="text-xs font-bold text-blue-700 hover:text-blue-800 hover:underline cursor-pointer"
        >
          {activeGroupAllSelected ? 'Batalkan Semua' : 'Pilih Semua'}
        </button>
      </div>

      {#each sortedGroupInvoices as invoice (invoice.id)}
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
              ? 'bg-blue-50/40 border-blue-500 shadow-sm shadow-blue-50/30'
              : 'bg-white border-slate-200 hover:bg-slate-50 hover:border-slate-300'}"
          role="checkbox"
          aria-checked={isSelected}
          tabindex="0"
        >
          <div class="flex items-center justify-between gap-3 w-full">
            <div class="flex items-center gap-3 min-w-0">
              <div class="w-5 h-5 rounded-md border-2 flex-shrink-0 flex items-center justify-center transition-all duration-200
                {isSelected ? 'border-blue-800 bg-blue-800' : 'border-slate-300 group-hover:border-blue-700'}"
              >
                {#if isSelected}
                  <svg class="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3.5" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                {/if}
              </div>

              <div class="min-w-0">
                <p class="font-bold text-slate-800 text-sm sm:text-base leading-tight">
                  {#if invoice.category}
                    <span class="text-blue-800 font-extrabold">{invoice.category.name}</span>
                  {/if}
                  {#if invoice.hijri_month}
                    <span class="text-slate-300 font-semibold"> · </span>
                    <span class="text-slate-600 font-semibold text-xs sm:text-sm">{getHijriMonthName(invoice.hijri_month)} {invoice.hijri_year}H</span>
                  {/if}
                </p>
                <p class="text-[10px] text-slate-400 font-medium mt-0.5">
                  Jatuh tempo: {new Date(invoice.due_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}
                </p>
              </div>
            </div>

            <div class="text-right flex-shrink-0">
              <p class="text-sm sm:text-base font-extrabold text-blue-900">{formatRupiah(invoice.amount_due)}</p>
            </div>
          </div>
        </div>
      {/each}
    </div>
  </Modal>
{/if}
