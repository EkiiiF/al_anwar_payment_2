<script lang="ts">
  import { onMount } from 'svelte';
  import { guardianApi } from '$lib/api';
  import { formatRupiah, formatDate, getMonthName } from '$lib/utils';
  import { Spinner, Alert, Card, Badge, EmptyState } from '$lib/components';
  import type { Payment } from '$lib/types';
  import { Receipt, Printer, CheckCircle2, X } from 'lucide-svelte';
  import logo from '$lib/assets/logo.png';

  let payments = $state<Payment[]>([]);
  let loading  = $state(true);
  let error    = $state('');

  const groupedTransactions = $derived(
    payments.reduce((acc, p) => {
      const key = p.external_id || p.id;
      if (!acc[key]) {
        acc[key] = {
          id: key,
          external_id: p.external_id,
          payment_date: p.payment_date,
          payment_method: p.payment_method,
          transaction_status: p.transaction_status,
          total_amount: 0,
          invoices: []
        };
      }
      acc[key].total_amount += p.amount_paid;
      if (p.invoice) {
        acc[key].invoices.push(p.invoice);
      }
      return acc;
    }, {} as Record<string, any>)
  );
  
  const transactionsList = $derived(
    Object.values(groupedTransactions).sort((a, b) => 
      new Date(b.payment_date).getTime() - new Date(a.payment_date).getTime()
    )
  );

  let selectedTransaction = $state<any>(null);

  onMount(async () => {
    try {
      const res = await guardianApi.getPaymentHistory();
      payments = res.data ?? [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat riwayat pembayaran.';
    } finally {
      loading = false;
    }
  });

  function viewReceipt(tx: any) { selectedTransaction = tx; }
  function printReceipt() { window.print(); }

  function getStatusVariant(status: string): 'success' | 'warning' | 'danger' | 'default' {
    const map: Record<string, 'success' | 'warning' | 'danger' | 'default'> = {
      settlement: 'success', paid: 'success',
      pending: 'warning',
      deny: 'danger', failure: 'danger', expire: 'danger'
    };
    return map[status] ?? 'default';
  }
</script>

<svelte:head>
  <title>Riwayat Pembayaran | Guardian - Al-Anwar Payment</title>
</svelte:head>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-black text-gray-900 tracking-tight">Riwayat Pembayaran</h1>
    <p class="text-gray-500 text-sm mt-1">Daftar seluruh transaksi pembayaran dan bukti kwitansi digital.</p>
  </div>

  {#if error}
    <Alert type="error" message={error} />
  {:else if loading}
    <Spinner size="lg" />
  {:else if transactionsList.length === 0}
    <Card>
      <EmptyState
        title="Belum ada riwayat pembayaran"
        description="Riwayat akan muncul setelah Anda melakukan pembayaran pertama."
      />
    </Card>
  {:else}
    <div class="grid grid-cols-1 gap-4">
      {#each transactionsList as tx (tx.id)}
        <Card>
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div class="flex items-start gap-4">
              <div class="p-3 rounded-xl bg-green-100 border border-green-200 flex-shrink-0">
                <Receipt size={20} class="text-green-700" aria-hidden="true" />
              </div>
              <div>
                <p class="font-bold text-gray-900">
                  Pembayaran {tx.invoices.length} Tagihan
                </p>
                <div class="flex flex-wrap gap-1 mt-1">
                  {#each tx.invoices as inv}
                    <span class="text-[10px] bg-gray-100 text-gray-600 px-1.5 py-0.5 rounded border border-gray-200">
                      {getMonthName(inv.month)} {inv.year}
                    </span>
                  {/each}
                </div>
                <p class="text-xs font-mono text-gray-400 mt-1.5">{tx.external_id ?? tx.id}</p>
                <p class="text-xs text-gray-500 mt-1">{formatDate(tx.payment_date, true)}</p>
                {#if tx.payment_method}
                  <p class="text-xs text-gray-400 mt-0.5">Via: {tx.payment_method}</p>
                {/if}
              </div>
            </div>

            <div class="flex flex-col sm:items-end gap-2">
              <p class="text-xl font-black text-gray-900">{formatRupiah(tx.total_amount)}</p>
              <Badge
                label={tx.transaction_status === 'settlement' ? 'Berhasil' : (tx.transaction_status ?? 'pending')}
                variant={getStatusVariant(tx.transaction_status ?? '')}
                dot
              />
              <button
                onclick={() => viewReceipt(tx)}
                class="text-xs text-green-600 hover:text-green-700 underline underline-offset-2 transition-colors font-semibold"
              >
                Lihat Struk
              </button>
            </div>
          </div>
        </Card>
      {/each}
    </div>
  {/if}
</div>

{#if selectedTransaction}
  <div
    class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    onclick={(e) => { if (e.target === e.currentTarget) selectedTransaction = null; }}
    role="presentation"
  >
    <div class="bg-white border border-gray-200 rounded-2xl w-full max-w-sm overflow-hidden shadow-2xl" id="receipt-print">
      <div class="p-6 text-center border-b border-gray-100 bg-green-50">
        <img src={logo} alt="Logo Al-Anwar" class="w-12 h-12 object-contain mx-auto mb-3" />
        <div class="w-12 h-12 rounded-full bg-green-100 border border-green-200 flex items-center justify-center mx-auto mb-3">
          <CheckCircle2 class="text-green-600 w-6 h-6" aria-hidden="true" />
        </div>
        <p class="text-xs text-green-700 font-bold uppercase tracking-widest">Pembayaran Berhasil</p>
        <p class="text-2xl font-black text-gray-900 mt-1">{formatRupiah(selectedTransaction.total_amount)}</p>
      </div>

      <div class="p-6 space-y-3">
        <div class="flex justify-between items-center gap-4">
          <span class="text-xs text-gray-500">No. Transaksi</span>
          <span class="text-xs font-bold text-gray-900 text-right font-mono">{selectedTransaction.external_id ?? selectedTransaction.id}</span>
        </div>
        <div class="flex justify-between items-start gap-4">
          <span class="text-xs text-gray-500">Detail Tagihan</span>
          <div class="text-right">
            {#each selectedTransaction.invoices as inv}
              <p class="text-xs font-bold text-gray-900">{getMonthName(inv.month)} {inv.year}</p>
            {/each}
          </div>
        </div>
        <div class="flex justify-between items-center gap-4">
          <span class="text-xs text-gray-500">Metode Bayar</span>
          <span class="text-xs font-bold text-gray-900 text-right">{selectedTransaction.payment_method ?? '-'}</span>
        </div>
        <div class="flex justify-between items-center gap-4">
          <span class="text-xs text-gray-500">Tanggal</span>
          <span class="text-xs font-bold text-gray-900 text-right">{formatDate(selectedTransaction.payment_date, true)}</span>
        </div>
        <div class="flex justify-between items-center gap-4">
          <span class="text-xs text-gray-500">Status</span>
          <span class="text-xs font-bold text-gray-900 text-right">{selectedTransaction.transaction_status ?? '-'}</span>
        </div>
      </div>

      <div class="px-6">
        <div class="border-t border-dashed border-gray-200"></div>
      </div>

      <div class="p-4 text-center">
        <p class="text-xs text-gray-600 font-medium">Pondok Pesantren Al-Anwar</p>
        <p class="text-xs text-gray-400">Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan</p>
      </div>

      <div class="flex gap-3 p-4 border-t border-gray-100">
        <button
          onclick={printReceipt}
          class="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl bg-green-600 hover:bg-green-700 text-white font-bold text-sm transition-all"
        >
          <Printer size={15} aria-hidden="true" />
          Cetak Struk
        </button>
        <button
          onclick={() => selectedTransaction = null}
          class="flex-1 py-2.5 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-700 font-bold text-sm transition-all"
        >
          Tutup
        </button>
      </div>
    </div>
  </div>
{/if}
