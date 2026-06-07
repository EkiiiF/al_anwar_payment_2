<script lang="ts">
  import { onMount } from 'svelte';
  import { BarChart3, Download, ChevronLeft, ChevronRight } from 'lucide-svelte';
  import { superUserApi } from '$lib/api';
  import { formatRupiah, getMonthName, MONTH_NAMES } from '$lib/utils';
  import { Button, Alert, Spinner, Card, Badge, Select } from '$lib/components';
  import { toast } from '$lib/stores/toast';

  let filterPeriod   = $state('monthly');
  let filterMonth    = $state(String(new Date().getMonth() + 1));
  let filterYear     = $state(String(new Date().getFullYear()));
  let filterCategory = $state('');

  let reportData = $state<unknown>(null);
  let loading    = $state(false);
  let error      = $state('');
  let exporting  = $state(false);
  
  let invoices = $state<unknown[]>([]);
  let invPagination = $state<{ page: number; limit: number; total: number; pages: number } | null>(null);
  let invPage = $state(1);
  let invLimit = $state(10);

  async function handleExport() {
    exporting = true;
    try {
      await superUserApi.exportReports({
        period: filterPeriod,
        year: filterYear,
        ...(filterPeriod === 'monthly' && { month: filterMonth }),
        ...(filterCategory && { category: filterCategory })
      });
      toast.success('Laporan berhasil diekspor.');
    } catch (e: any) {
      toast.error(e.message || 'Gagal mengekspor laporan.');
    } finally {
      exporting = false;
    }
  }

  const periodOptions = [
    { value: 'monthly',   label: 'Bulanan' },
    { value: 'quarterly', label: 'Triwulan' },
    { value: 'semester',  label: 'Semester' },
    { value: 'yearly',    label: 'Tahunan' }
  ];


  const monthOptions = MONTH_NAMES.map((n, i) => ({ value: String(i + 1), label: n }));

  async function fetchReport() {
    loading = true;
    error = '';
    reportData = null;
    try {
      const filters: Record<string, string> = {
        period: filterPeriod,
        year: filterYear,
        ...(filterPeriod === 'monthly' && { month: filterMonth }),
        ...(filterCategory && { category: filterCategory })
      };
      const [resReport, resInv] = await Promise.all([
        superUserApi.getReports(filters),
        superUserApi.getInvoicesPaginated({}, invPage, invLimit)
      ]);
      reportData = resReport.data;
      invoices = resInv.data?.invoices ?? [];
      invPagination = resInv.data?.pagination ?? null;
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat laporan.';
    } finally {
      loading = false;
    }
  }
  
  async function fetchInvoices() {
    try {
      const res = await superUserApi.getInvoicesPaginated({}, invPage, invLimit);
      invoices = res.data?.invoices ?? [];
      invPagination = res.data?.pagination ?? null;
    } catch (e: unknown) {
      console.error(e);
    }
  }
  
  function prevInvPage() {
    if (invPage > 1) {
      invPage--;
      fetchInvoices();
    }
  }
  
  function nextInvPage() {
    if (invPagination && invPage < invPagination.pages) {
      invPage++;
      fetchInvoices();
    }
  }

  onMount(fetchReport);
</script>

<svelte:head>
  <title>Laporan Keuangan | Dashboard Super User</title>
</svelte:head>

<div class="space-y-6 flex flex-col flex-1 min-h-0">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Laporan Keuangan</h1>
      <p class="text-gray-500 text-sm mt-1">Filter dan analisis pemasukan berdasarkan periode, kategori, dan status.</p>
    </div>
    <Button 
      onclick={handleExport} 
      variant="outline" 
      size="md"
      loading={exporting}
    >
      {#snippet children()}
        <Download size={16} /> Ekspor Data (CSV)
      {/snippet}
    </Button>
  </div>

  <Card>
    <h2 class="text-sm font-bold text-gray-700 mb-4">Filter Laporan</h2>
    <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
      <Select id="filter-period" label="Periode" bind:value={filterPeriod} options={periodOptions} />
      {#if filterPeriod === 'monthly'}
        <Select id="filter-month" label="Bulan" bind:value={filterMonth} options={monthOptions} />
      {/if}
      <div>
        <label for="filter-year" class="text-xs font-semibold text-gray-600 uppercase tracking-wider ml-0.5 block mb-1.5">Tahun</label>
        <input id="filter-year" type="number" bind:value={filterYear}
          class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-gray-900 text-sm outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500" />
      </div>
    </div>
    <div class="flex justify-end mt-4">
      <Button onclick={fetchReport} variant="primary" size="md" loading={loading}>
        {#snippet children()}
          <BarChart3 size={16} />
          Tampilkan Laporan
        {/snippet}
      </Button>
    </div>
  </Card>

  {#if error}
    <Alert type="error" message={error} />
  {:else if loading}
    <Spinner size="lg" />
  {:else if reportData}
    {@const data = reportData as Record<string, unknown>}
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {#each [
        { key: 'total_invoices',    label: 'Total Faktur',     suffix: 'faktur' },
        { key: 'paid_count',        label: 'Lunas',            suffix: 'faktur' },
        { key: 'unpaid_count',      label: 'Belum Bayar',      suffix: 'faktur' },
        { key: 'total_amount_paid', label: 'Total Pemasukan',  isCurrency: true }
      ] as item}
        <Card>
          <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">{item.label}</p>
          <p class="text-2xl font-black text-gray-900">
            {#if item.isCurrency}
              {formatRupiah(Number(data[item.key] ?? 0))}
            {:else}
              {data[item.key] ?? 0}
              {#if item.suffix}
                <span class="text-sm font-normal text-gray-500"> {item.suffix}</span>
              {/if}
            {/if}
          </p>
        </Card>
      {/each}
    </div>

    {#if invoices.length > 0}
      <Card padding={false} class="flex-1 flex flex-col min-h-0">
        <div class="p-5 border-b border-gray-100">
          <h2 class="font-bold text-gray-900">Detail Tagihan</h2>
        </div>
        <div class="overflow-x-auto flex-1">
          <table class="w-full text-sm" aria-label="Tabel detail laporan keuangan">
            <thead>
              <tr class="bg-gray-50 border-b border-gray-100">
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">No. Invoice</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Santri</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Periode</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-right">Nominal</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              {#each invoices as inv}
                {@const i = inv as Record<string, unknown>}
                <tr class="hover:bg-gray-50 transition-colors">
                  <td class="px-5 py-3 font-mono text-xs text-green-600">{String(i.invoice_number)}</td>
                  <td class="px-5 py-3 text-gray-900">{String((i.student as Record<string, unknown>)?.full_name ?? '-')}</td>
                  <td class="px-5 py-3 text-gray-500">{getMonthName(Number(i.month))} {i.year}</td>
                  <td class="px-5 py-3 text-right font-bold text-gray-900">{formatRupiah(Number(i.amount_due))}</td>
                  <td class="px-5 py-3 text-center">
                    <Badge
                      label={({ paid: 'Lunas', unpaid: 'Belum Bayar', pending: 'Pending', expired: 'Kadaluwarsa', cancelled: 'Dibatalkan' } as Record<string,string>)[String(i.status)] ?? String(i.status)}
                      variant={(({ paid: 'success', unpaid: 'warning', pending: 'info', expired: 'default', cancelled: 'danger' } as Record<string,string>)[String(i.status)] ?? 'default') as 'success' | 'warning' | 'info' | 'default' | 'danger'}
                    />
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        
        {#if invPagination}
          <div class="flex items-center justify-between bg-white p-4 rounded-xl border-t border-gray-200">
            <p class="text-sm text-gray-600">
              Menampilkan <span class="font-bold text-gray-900">{((invPage - 1) * invLimit) + 1} - {Math.min(invPage * invLimit, invPagination.total)}</span> dari <span class="font-bold text-gray-900">{invPagination.total}</span> tagihan
            </p>
            <div class="flex items-center gap-2">
              <Button variant="secondary" onclick={prevInvPage} disabled={invPage <= 1} size="sm">
                {#snippet children()}<ChevronLeft size={16} /> Sebelumnya{/snippet}
              </Button>
              <div class="px-3 py-1.5 bg-green-50 text-green-700 rounded-lg text-sm font-bold border border-green-200">
                Halaman {invPage} / {invPagination.pages}
              </div>
              <Button variant="secondary" onclick={nextInvPage} disabled={!invPagination || invPage >= invPagination.pages} size="sm">
                {#snippet children()}Selanjutnya <ChevronRight size={16} />{/snippet}
              </Button>
            </div>
          </div>
        {/if}
      </Card>
    {/if}
  {:else}
    <Card>
      <p class="text-center text-gray-500 py-8">Klik "Tampilkan Laporan" untuk melihat data.</p>
    </Card>
  {/if}
</div>
