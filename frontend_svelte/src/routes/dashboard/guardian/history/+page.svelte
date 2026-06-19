<script lang="ts">
  import { onMount } from 'svelte';
  import { guardianApi } from '$lib/api';
  import { formatRupiah, formatDate, getHijriMonthName } from '$lib/utils';
  import { Spinner, Alert, Card, Badge, EmptyState, Button, Paginator } from '$lib/components';
  import type { Payment } from '$lib/types';
  import { Receipt, Printer, CheckCircle2, X, Clock, XCircle } from 'lucide-svelte';
  import logo from '$lib/assets/logo.png';

  let payments = $state<Payment[]>([]);
  let loading  = $state(true);
  let error    = $state('');
  let paying   = $state(false);

  const groupedTransactions = $derived(
    payments.reduce((acc, p) => {
      const key = p.external_id || p.id;
      if (!acc[key]) {
        acc[key] = {
          id: key,
          external_id: p.external_id,
          snap_token: p.snap_token,
          payment_date: p.payment_date,
          created_at: p.created_at,
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
      new Date(b.created_at || b.payment_date || 0).getTime() - new Date(a.created_at || a.payment_date || 0).getTime()
    )
  );

  let selectedTransaction = $state<any>(null);

  let page = $state(1);
  let limit = $state(5);
  const totalTransactions = $derived(transactionsList.length);
  const totalPages = $derived(Math.ceil(totalTransactions / limit) || 1);
  const paginatedTransactions = $derived(
    transactionsList.slice((page - 1) * limit, page * limit)
  );

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

  function resumePayment(snapToken: string) {
    type SnapWindow = Window & typeof globalThis & {
      snap?: { pay: (token: string, options: Record<string, unknown>) => void }
    };
    const w = window as SnapWindow;
    if (!w.snap) {
      alert('Midtrans Snap tidak tersedia. Silakan muat ulang halaman.');
      return;
    }
    paying = true;
    w.snap.pay(snapToken, {
      onSuccess: () => { window.location.reload(); },
      onPending: () => { window.location.reload(); },
      onError: (result: unknown) => {
        console.error('[Midtrans Error]', result);
        alert('Pembayaran gagal. Silakan coba lagi.');
        paying = false;
      },
      onClose: () => { paying = false; }
    });
  }

  function getStatusVariant(status: string): 'success' | 'warning' | 'danger' | 'default' {
    const map: Record<string, 'success' | 'warning' | 'danger' | 'default'> = {
      settlement: 'success', paid: 'success',
      pending: 'warning', pending_payment: 'warning',
      deny: 'danger', failure: 'danger', expire: 'danger'
    };
    return map[status] ?? 'default';
  }

  function translateTransactionStatus(status: string): string {
    const s = status?.toLowerCase() || '';
    if (s === 'settlement' || s === 'paid' || s === 'capture') {
      return 'Berhasil';
    } else if (s === 'pending' || s === 'pending_payment') {
      return 'Menunggu Pembayaran';
    } else if (s === 'expire' || s === 'expired') {
      return 'Kedaluwarsa';
    } else if (s === 'deny') {
      return 'Ditolak';
    } else if (s === 'cancel' || s === 'cancelled') {
      return 'Dibatalkan';
    } else if (s === 'failure' || s === 'failed') {
      return 'Gagal';
    }
    return status || '-';
  }

  function getReceiptStyle(status: string) {
    const s = status?.toLowerCase() || '';
    if (s === 'settlement' || s === 'paid' || s === 'capture') {
      return {
        bg: 'bg-emerald-50 border-b border-emerald-100',
        text: 'text-emerald-700',
        label: 'Pembayaran Berhasil',
        iconBg: 'bg-emerald-100 border-emerald-200',
        iconColor: 'text-emerald-600',
        badgeVariant: 'success' as const
      };
    } else if (s === 'pending' || s === 'pending_payment') {
      return {
        bg: 'bg-amber-50 border-b border-amber-100',
        text: 'text-amber-700',
        label: 'Menunggu Pembayaran',
        iconBg: 'bg-amber-100 border-amber-200',
        iconColor: 'text-amber-600',
        badgeVariant: 'warning' as const
      };
    } else {
      return {
        bg: 'bg-red-50 border-b border-red-100',
        text: 'text-red-700',
        label: 'Pembayaran Gagal',
        iconBg: 'bg-red-100 border-red-200',
        iconColor: 'text-red-600',
        badgeVariant: 'danger' as const
      };
    }
  }
</script>

<svelte:head>
  <title>Riwayat Pembayaran | Guardian - Al-Anwar Payment</title>
  <script src="https://app.sandbox.midtrans.com/snap/snap.js" data-client-key={import.meta.env.VITE_MIDTRANS_CLIENT_KEY}></script>
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
    <div class="space-y-4">
      <div class="grid grid-cols-1 gap-4">
        {#each paginatedTransactions as tx (tx.id)}
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
                        {getHijriMonthName(inv.hijri_month)} {inv.hijri_year} H
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
                  label={translateTransactionStatus(tx.transaction_status ?? '')}
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
      <Paginator
        page={page}
        limit={limit}
        total={totalTransactions}
        pages={totalPages}
        label="transaksi"
        onPageChange={(p) => page = p}
        onLimitChange={(l) => { limit = l; page = 1; }}
      />
    </div>
  {/if}
</div>

{#if selectedTransaction}
  {@const rStyle = getReceiptStyle(selectedTransaction.transaction_status ?? '')}
  <div
    class="fixed inset-0 bg-black/40 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    onclick={(e) => { if (e.target === e.currentTarget) selectedTransaction = null; }}
    role="presentation"
  >
    <div class="bg-white border border-slate-200 rounded-2xl w-full max-w-sm overflow-hidden shadow-2xl transition-all scale-100" id="receipt-print">
      <div class="p-6 text-center {rStyle.bg}">
        <img src={logo} alt="Logo Al-Anwar" class="w-12 h-12 object-contain mx-auto mb-3" />
        <div class="w-12 h-12 rounded-full {rStyle.iconBg} border flex items-center justify-center mx-auto mb-3">
          {#if rStyle.badgeVariant === 'success'}
            <CheckCircle2 class="{rStyle.iconColor} w-6 h-6" aria-hidden="true" />
          {:else if rStyle.badgeVariant === 'warning'}
            <Clock class="{rStyle.iconColor} w-6 h-6" aria-hidden="true" />
          {:else}
            <XCircle class="{rStyle.iconColor} w-6 h-6" aria-hidden="true" />
          {/if}
        </div>
        <p class="text-xs {rStyle.text} font-bold uppercase tracking-widest">{rStyle.label}</p>
        <p class="text-2xl font-black text-slate-900 mt-1">{formatRupiah(selectedTransaction.total_amount)}</p>
      </div>

      <div class="p-6 space-y-3">
        <div class="flex justify-between items-center gap-4">
          <span class="text-xs text-slate-500">No. Transaksi</span>
          <span class="text-xs font-bold text-slate-900 text-right font-mono">{selectedTransaction.external_id ?? selectedTransaction.id}</span>
        </div>
        <div class="flex justify-between items-start gap-4">
          <span class="text-xs text-slate-500">Detail Tagihan</span>
          <div class="text-right">
            {#each selectedTransaction.invoices as inv}
              <p class="text-xs font-bold text-slate-900">
                {#if inv.category}
                  <span class="text-emerald-800 text-[10px] font-bold block">{inv.category.name}</span>
                {/if}
                Tagihan {getHijriMonthName(inv.hijri_month)} {inv.hijri_year} H
              </p>
            {/each}
          </div>
        </div>
        <div class="flex justify-between items-center gap-4">
          <span class="text-xs text-slate-500">Metode Bayar</span>
          <span class="text-xs font-bold text-slate-900 text-right">{selectedTransaction.payment_method ?? '-'}</span>
        </div>
        {#if selectedTransaction.payment_date}
          <div class="flex justify-between items-center gap-4">
            <span class="text-xs text-slate-500">Tanggal Bayar</span>
            <span class="text-xs font-bold text-slate-900 text-right">{formatDate(selectedTransaction.payment_date, true)}</span>
          </div>
        {/if}
        <div class="flex justify-between items-center gap-4">
          <span class="text-xs text-slate-500">Status</span>
          <span class="text-xs font-bold {rStyle.text} text-right">{translateTransactionStatus(selectedTransaction.transaction_status ?? '')}</span>
        </div>
      </div>

      <div class="px-6">
        <div class="border-t border-dashed border-slate-200"></div>
      </div>

      <div class="p-4 text-center">
        <p class="text-xs text-slate-600 font-bold">Pondok Pesantren Al-Anwar</p>
        <p class="text-[10px] text-slate-400 mt-1 leading-normal">Dusun Kauman, Desa Selo, RT 05/RW 08<br/>Kecamatan Tawangharjo Kabupaten Grobogan</p>
      </div>

      <div class="flex gap-3 p-4 border-t border-slate-100">
        {#if selectedTransaction.transaction_status === 'settlement' || selectedTransaction.transaction_status === 'paid' || selectedTransaction.transaction_status === 'capture'}
          <Button
            onclick={printReceipt}
            variant="primary"
            class="flex-1"
          >
            {#snippet children()}
              <Printer size={15} aria-hidden="true" />
              <span>Cetak Struk</span>
            {/snippet}
          </Button>
        {:else if (selectedTransaction.transaction_status === 'pending' || selectedTransaction.transaction_status === 'pending_payment') && selectedTransaction.snap_token}
          <Button
            onclick={() => resumePayment(selectedTransaction.snap_token)}
            disabled={paying}
            variant="primary"
            class="flex-1"
          >
            {#snippet children()}Selesaikan Pembayaran{/snippet}
          </Button>
        {/if}
        <Button
          onclick={() => selectedTransaction = null}
          variant="outline"
          class="flex-1"
        >
          {#snippet children()}Tutup{/snippet}
        </Button>
      </div>
    </div>
  </div>
{/if}
