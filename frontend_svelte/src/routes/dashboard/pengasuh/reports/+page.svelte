<script lang="ts">
  import { onMount } from 'svelte';
  import { BarChart3, Download, Filter, Search, FileText } from 'lucide-svelte';
  import { pengasuhApi } from '$lib/api';
  import { formatRupiah, getMonthName, MONTH_NAMES } from '$lib/utils';
  import { Button, Alert, Spinner, Card, Badge, Select } from '$lib/components';
  import { toast } from '$lib/stores/toast';

  let filterPeriod   = $state('monthly');
  let filterMonth    = $state(String(new Date().getMonth() + 1));
  let filterYear     = $state(String(new Date().getFullYear()));
  let filterCategory = $state('');

  let reportData = $state<any>(null);
  let loading    = $state(false);
  let error      = $state('');
  let exporting  = $state(false);

  async function handleExport() {
    exporting = true;
    try {
      await pengasuhApi.exportReports({
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

  const categoryOptions = [
    { value: '',                      label: 'Semua Kategori' },
    { value: 'Syahriyyah Pondok',     label: 'Syahriyyah Pondok' },
    { value: 'Syahriyyah Muhadhoroh',  label: 'Syahriyyah Muhadhoroh' },
    { value: 'Daftar Ulang',          label: 'Daftar Ulang' },
    { value: 'Ujian Semester',        label: 'Ujian Semester' }
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
      const res = await pengasuhApi.getReports(filters);
      reportData = res.data;
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat laporan.';
    } finally {
      loading = false;
    }
  }

  onMount(fetchReport);
</script>

<svelte:head>
  <title>Laporan Keuangan | Dashboard Pengasuh</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-purple-100 border border-purple-200 text-purple-700 text-xs font-semibold uppercase tracking-wider mb-2">
        <span class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse" aria-hidden="true"></span>
        Mode Hanya Lihat
      </div>
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Laporan Keuangan & Monitoring</h1>
      <p class="text-gray-500 text-sm mt-1">Filter, analisis, dan ekspor data pemasukan pesantren secara real-time.</p>
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

  <Card class="border-l-4 border-l-purple-500">
    <div class="flex items-center gap-2 text-sm font-bold text-gray-700 mb-4">
      <Filter size={16} class="text-purple-600" />
      <span>Filter Laporan Keuangan</span>
    </div>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
      <Select id="filter-period" label="Periode" bind:value={filterPeriod} options={periodOptions} />
      {#if filterPeriod === 'monthly'}
        <Select id="filter-month" label="Bulan" bind:value={filterMonth} options={monthOptions} />
      {/if}
      <div>
        <label for="filter-year" class="text-xs font-semibold text-gray-600 uppercase tracking-wider ml-0.5 block mb-1.5">Tahun</label>
        <input id="filter-year" type="number" bind:value={filterYear}
          class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-gray-900 text-sm outline-none focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500" />
      </div>
      <Select id="filter-category" label="Kategori Tagihan" bind:value={filterCategory} options={categoryOptions} />
    </div>
    <div class="flex justify-end mt-4">
      <Button onclick={fetchReport} variant="primary" size="md" loading={loading} class="!bg-purple-600 hover:!bg-purple-700">
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
    {@const data = reportData as Record<string, any>}
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {#each [
        { key: 'total_invoices',    label: 'Total Faktur',     suffix: 'faktur', color: 'text-gray-900' },
        { key: 'paid_count',        label: 'Lunas (Settled)',  suffix: 'faktur', color: 'text-green-600' },
        { key: 'unpaid_count',      label: 'Belum Dibayar',    suffix: 'faktur', color: 'text-amber-600' },
        { key: 'total_amount_paid', label: 'Total Uang Masuk', isCurrency: true, color: 'text-purple-700' }
      ] as item}
        <Card class="bg-white hover:shadow-md transition-all duration-200">
          <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">{item.label}</p>
          <p class="text-2xl font-black {item.color}">
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

    {@const invoices = (data.invoices as any[]) ?? []}
    {#if invoices.length > 0}
      <Card padding={false}>
        <div class="p-5 border-b border-gray-100 flex items-center justify-between">
          <h2 class="font-bold text-gray-900 flex items-center gap-2">
            <FileText size={18} class="text-purple-600" />
            <span>Detail Tagihan Santri</span>
          </h2>
          <Badge label={`${invoices.length} baris`} variant="purple" />
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-sm" aria-label="Tabel detail laporan keuangan">
            <thead>
              <tr class="bg-gray-50 border-b border-gray-100">
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">No. Invoice</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Santri</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Periode Hijriah</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-right">Nominal</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              {#each invoices as inv}
                {@const i = inv as Record<string, any>}
                <tr class="hover:bg-purple-50/10 transition-colors">
                  <td class="px-5 py-3 font-mono text-xs text-purple-600">{String(i.invoice_number)}</td>
                  <td class="px-5 py-3 text-gray-900 font-medium">
                    {String(i.student?.name?.first_name || '')} {String(i.student?.name?.last_name || '')}
                  </td>
                  <td class="px-5 py-3 text-gray-500">
                    {#if i.hijri_month}
                      {getMonthName(Number(i.month))} {i.year} / {i.hijri_month}H {i.hijri_year}
                    {:else}
                      {getMonthName(Number(i.month))} {i.year}
                    {/if}
                  </td>
                  <td class="px-5 py-3 text-right font-black text-gray-900">{formatRupiah(Number(i.amount_due))}</td>
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
      </Card>
    {/if}
  {:else}
    <Card>
      <p class="text-center text-gray-500 py-8">Klik "Tampilkan Laporan" untuk memuat data analisis.</p>
    </Card>
  {/if}
</div>
