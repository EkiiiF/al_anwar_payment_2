<script lang="ts">
  import { onMount } from 'svelte';
  import { Wallet, CheckCircle, Bell, ChevronRight } from 'lucide-svelte';
  import { guardianApi } from '$lib/api';
  import { formatRupiah, getInitials } from '$lib/utils';
  import { Spinner, Alert, StatCard, Card, Badge } from '$lib/components';
  import type { GuardianDashboardStats } from '$lib/types';

  let stats   = $state<GuardianDashboardStats | null>(null);
  let loading = $state(true);
  let error   = $state('');

  const lastPaymentTotal = $derived.by(() => {
    if (!stats?.recent_payments || stats.recent_payments.length === 0) return 0;
    const latest = stats.recent_payments[0];
    if (!latest.external_id) return latest.amount_paid;
    return stats.recent_payments
      .filter(p => p.external_id === latest.external_id)
      .reduce((sum, p) => sum + p.amount_paid, 0);
  });

  onMount(async () => {
    try {
      const res = await guardianApi.getDashboard();
      stats = res.data;
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat dashboard.';
    } finally {
      loading = false;
    }
  });
</script>

<svelte:head>
  <title>Dashboard | Guardian - Al-Anwar Payment</title>
</svelte:head>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-black text-gray-900 tracking-tight">Selamat Datang</h1>
    <p class="text-gray-500 text-sm mt-1">Pantau tagihan dan riwayat pembayaran santri Anda.</p>
  </div>

  {#if loading}
    <Spinner size="lg" />
  {:else if error}
    <Alert type="error" message={error} />
  {:else if stats}
    {@const student = stats.student}

    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <StatCard
        title="Tagihan Belum Bayar"
        value={stats.unpaid_count}
        subtitle="Perlu segera dilunasi"
        icon={Wallet}
        color={stats.unpaid_count > 0 ? 'amber' : 'green'}
        accent={stats.unpaid_count > 0}
      />
      <StatCard
        title="Total Tertunggak"
        value={formatRupiah(stats.total_unpaid)}
        subtitle="Jumlah yang harus dibayar"
        icon={Bell}
        color="amber"
      />
      {#if stats.recent_payments && stats.recent_payments.length > 0}
        <StatCard
          title="Pembayaran Terakhir"
          value={formatRupiah(lastPaymentTotal)}
          subtitle="Transaksi terakhir (total)"
          icon={CheckCircle}
          color="green"
        />
      {/if}
    </div>

    {#if stats.unpaid_count > 0}
      <div class="bg-amber-50 border border-amber-200 rounded-xl p-3.5 sm:p-4">
        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
          <div>
            <p class="font-bold text-amber-900 mb-0.5 text-sm">
              Anda memiliki <span class="text-amber-600">{stats.unpaid_count} tagihan</span> yang belum dibayar
            </p>
            <p class="text-sm text-amber-700">
              Total: <strong>{formatRupiah(stats.total_unpaid)}</strong>
            </p>
          </div>
          <a
            href="/dashboard/guardian/invoices"
            class="flex-shrink-0 inline-flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-700 text-white font-semibold rounded-lg text-xs sm:text-sm transition-all shadow-sm cursor-pointer"
          >
            Bayar Sekarang
            <ChevronRight size={16} aria-hidden="true" />
          </a>
        </div>
      </div>
    {:else}
      <div class="bg-green-50 border border-green-200 rounded-xl p-3.5 sm:p-4 flex items-center gap-3">
        <CheckCircle size={20} class="text-green-600 flex-shrink-0" />
        <div>
          <p class="font-semibold text-green-800 text-sm">Semua tagihan lunas!</p>
          <p class="text-xs sm:text-sm text-green-600">Tidak ada tagihan yang tertunggak saat ini.</p>
        </div>
      </div>
    {/if}
  {/if}
</div>
