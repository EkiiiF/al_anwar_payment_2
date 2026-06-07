<script lang="ts">
  import { Card, Badge } from '$lib/components';
  import { formatRupiah, getInvoiceStatusStyle, getHijriMonthName, getMonthName } from '$lib/utils';
  import type { Invoice } from '$lib/types';

  let {
    otherInvoices = [],
    refreshStatus
  }: {
    otherInvoices: Invoice[];
    refreshStatus: (orderId: string) => void;
  } = $props();
</script>

<div class="space-y-3.5">
  <h2 class="font-extrabold text-slate-900 text-lg select-none">Riwayat Transaksi & Pembayaran</h2>
  
  <Card padding={false} class="border-slate-200 shadow-sm overflow-hidden bg-white rounded-2xl">
    <div class="divide-y divide-slate-100">
      {#each otherInvoices as invoice (invoice.id)}
        {@const style = getInvoiceStatusStyle(invoice.status)}
        <div class="flex flex-col sm:flex-row sm:items-center justify-between p-4.5 hover:bg-slate-50/50 transition-colors gap-3">
          <div class="min-w-0">
            <p class="font-bold text-slate-900 text-sm tracking-tight">
              {#if invoice.category}
                <span class="text-emerald-805 font-bold">{invoice.category.name}</span>
              {/if}
              <span class="text-slate-400 font-normal"> · </span>
              {#if invoice.hijri_month}
                {getHijriMonthName(invoice.hijri_month)} {invoice.hijri_year} H
              {:else}
                Tagihan {getMonthName(invoice.month)} {invoice.year}
              {/if}
            </p>
            <p class="text-[10px] font-mono text-slate-450 mt-0.5">{invoice.invoice_number}</p>
            {#if invoice.semester}
              <div class="flex items-center gap-2 mt-1">
                <Badge label={`Sem ${invoice.semester}`} variant="warning" class="text-[9px] py-0.5" />
                <span class="text-[10px] text-emerald-600 font-bold">TA {invoice.academic_year_label}</span>
              </div>
            {/if}
          </div>

          <div class="text-left sm:text-right flex sm:flex-col items-center sm:items-end justify-between sm:justify-start gap-1">
            <p class="font-extrabold text-slate-950 text-sm leading-none">{formatRupiah(invoice.amount_due)}</p>
            <div class="flex flex-col items-end gap-1.5 mt-1">
              <span class="text-xs font-bold {style.text} bg-slate-50 px-2 py-0.5 rounded border border-slate-100">{style.label}</span>
              {#if invoice.status === 'pending' && invoice.payments && invoice.payments.length > 0}
                {@const lastPayment = invoice.payments[0]}
                {#if lastPayment.external_id}
                  <button 
                    type="button"
                    onclick={(e) => { e.stopPropagation(); refreshStatus(lastPayment.external_id!); }}
                    class="text-[10px] font-extrabold text-teal-650 hover:text-teal-700 underline transition-colors"
                  >
                    Perbarui Status
                  </button>
                {/if}
              {/if}
            </div>
          </div>
        </div>
      {/each}
    </div>
  </Card>
</div>
