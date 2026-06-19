<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { BarChart3, Download, FileText, FileDown } from 'lucide-svelte';
  import { superUserApi } from '$lib/api';
  import { auth } from '$lib/stores/auth';
  import { formatRupiah, getHijriMonthName, HIJRI_MONTH_NAMES } from '$lib/utils';
  import { Button, Alert, Card, Badge, Select, Paginator, ConfirmDialog } from '$lib/components';
  import { toast } from '$lib/stores/toast';
  import type { Category } from '$lib/types';

  let filterPeriod   = $state('');
  let filterMonth    = $state('');
  let filterYear     = $state('');
  let filterCategory = $state('');

  let reportData = $state<any>(null);
  let loading    = $state(false);
  let error      = $state('');
  let exporting  = $state(false);

  let availableYears = $state<number[]>([]);
  let categories = $state<Category[]>([]);
  let categoryOptions = $derived([
    { value: '', label: 'Pilih Kategori...' },
    ...categories.map(c => ({ value: c.id, label: c.name }))
  ]);
  
  let invoices = $state<any[]>([]);
  let invPagination = $state<{ page: number; limit: number; total: number; pages: number } | null>(null);
  let invPage = $state(1);
  let invLimit = $state(50); // Default to 50 records

  let showExportDropdown = $state(false);
  let showConfirmModal = $state(false);
  let exportType = $state<'pdf' | 'excel'>('pdf');

  const exporterName = $derived(
    $auth.user 
      ? `${$auth.user.first_name} ${$auth.user.last_name ?? ''}`.trim() || $auth.user.username 
      : 'Sistem'
  );

  const periodOptions = [
    { value: '',          label: 'Pilih Periode...' },
    { value: 'monthly',   label: 'Bulanan' },
    { value: 'quarterly', label: 'Triwulan' },
    { value: 'semester',  label: 'Semester' },
    { value: 'yearly',    label: 'Tahunan' }
  ];

  const monthOptions = [
    { value: '', label: 'Pilih Bulan...' },
    ...HIJRI_MONTH_NAMES.map((n, i) => ({ value: String(i + 1), label: n }))
  ];

  async function fetchReport() {
    loading = true;
    error = '';
    try {
      const reportFilters: Record<string, string> = {
        ...(filterYear && { year: filterYear }),
        ...(filterPeriod === 'monthly' && filterMonth && { month: filterMonth }),
        ...(filterCategory && { category: filterCategory })
      };

      const invoiceFilters: Record<string, string> = {
        status: 'paid', // Always only paid/lunas
        ...(filterYear && { year: filterYear }),
        ...(filterPeriod === 'monthly' && filterMonth && { month: filterMonth })
      };

      const [resReport, resInv] = await Promise.all([
        superUserApi.getReports(reportFilters),
        superUserApi.getInvoicesPaginated(invoiceFilters, invPage, invLimit)
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
      const invoiceFilters: Record<string, string> = {
        status: 'paid', // Always only paid/lunas
        ...(filterYear && { year: filterYear }),
        ...(filterPeriod === 'monthly' && filterMonth && { month: filterMonth })
      };
      const res = await superUserApi.getInvoicesPaginated(invoiceFilters, invPage, invLimit);
      invoices = res.data?.invoices ?? [];
      invPagination = res.data?.pagination ?? null;
    } catch (e: unknown) {
      console.error(e);
    }
  }
  
  function handleInvPageChange(newPage: number) {
    invPage = newPage;
    fetchInvoices();
  }
  
  function handleInvLimitChange(newLimit: number) {
    invLimit = newLimit;
    invPage = 1;
    fetchInvoices();
  }

  // Reactive effect to load data when filters change
  $effect(() => {
    // If Bulanan is selected but month is empty, wait for month selection
    if (filterPeriod === 'monthly' && !filterMonth) {
      return;
    }
    invPage = 1;
    fetchReport();
  });

  async function executeExport() {
    showConfirmModal = false;
    if (exportType === 'pdf') {
      await tick();
      // Use setTimeout to allow Svelte to fully hide the ConfirmDialog before printing
      setTimeout(() => {
        window.print();
      }, 150);
    } else {
      exporting = true;
      try {
        await superUserApi.exportReports({
          ...(filterYear && { year: filterYear }),
          ...(filterPeriod === 'monthly' && filterMonth && { month: filterMonth }),
          ...(filterCategory && { category: filterCategory }),
          exporter: exporterName
        });
        toast.success('Laporan berhasil diekspor.');
      } catch (e: any) {
        toast.error(e.message || 'Gagal mengekspor laporan.');
      } finally {
        exporting = false;
      }
    }
  }

  onMount(async () => {
    try {
      const [dashboardRes, catRes] = await Promise.all([
        superUserApi.getDashboard(),
        superUserApi.getCategories()
      ]);
      
      if (dashboardRes.data?.available_years) {
        availableYears = dashboardRes.data.available_years;
      } else {
        const currentY = 1448;
        availableYears = [];
        for (let y = currentY; y >= 1443; y--) {
          availableYears.push(y);
        }
      }

      categories = catRes.data || [];
    } catch (e) {
      console.error(e);
      const currentY = 1448;
      availableYears = [];
      for (let y = currentY; y >= 1443; y--) {
        availableYears.push(y);
      }
    }
  });
</script>

<svelte:head>
  <title>Laporan Keuangan | Dashboard Super User</title>
</svelte:head>

<div class="space-y-6 flex flex-col flex-1 min-h-0">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 no-print">
    <div>
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Laporan Keuangan</h1>
      <p class="text-gray-500 text-sm mt-1">Analisis pemasukan pesantren secara real-time.</p>
    </div>
  </div>

  <div id="print-area" class="space-y-6 flex flex-col flex-1 min-h-0">
    <!-- Printable Header (Only visible when printing/PDF) -->
    <div class="hidden print:block border-b border-gray-200 pb-4 mb-4">
      <h1 class="text-2xl font-bold text-gray-900">Laporan Keuangan Pondok Pesantren Al-Anwar</h1>
      <div class="grid grid-cols-2 gap-4 mt-4 text-xs text-gray-600">
        <div>
          <p><span class="font-semibold text-gray-800">Nama Laporan:</span> Laporan Pemasukan Santri (Lunas)</p>
          <p>
            <span class="font-semibold text-gray-800">Periode:</span> 
            {#if filterPeriod}
              {filterPeriod === 'monthly' ? 'Bulanan' : filterPeriod === 'quarterly' ? 'Triwulan' : filterPeriod === 'semester' ? 'Semester' : 'Tahunan'} 
              {#if filterPeriod === 'monthly' && filterMonth}
                - {getHijriMonthName(Number(filterMonth))}
              {/if}
              - {filterYear} H
            {:else}
              Semua Periode
            {/if}
          </p>
          <p><span class="font-semibold text-gray-800">Kategori:</span> {categories.find(c => c.id === filterCategory)?.name ?? 'Semua Kategori'}</p>
        </div>
        <div>
          <p><span class="font-semibold text-gray-800">Pengekspor:</span> {exporterName} ({$auth.user?.role?.description ?? ($auth.user?.role?.name === 'super_user' ? 'Super User/Bendahara' : 'Pengasuh')})</p>
          <p><span class="font-semibold text-gray-800">Waktu Ekspor:</span> {new Date().toLocaleString('id-ID', { dateStyle: 'long', timeStyle: 'short' })}</p>
          <p><span class="font-semibold text-gray-800">Total Pemasukan:</span> {formatRupiah(Number(reportData?.total_amount_paid ?? 0))}</p>
        </div>
      </div>
    </div>

    {#if reportData}
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {#each [
          { key: 'total_invoices',    label: 'Total Tagihan',     suffix: 'tagihan' },
          { key: 'paid_count',        label: 'Lunas',            suffix: 'tagihan' },
          { key: 'unpaid_count',      label: 'Belum Bayar',      suffix: 'tagihan' },
          { key: 'total_amount_paid', label: 'Total Pemasukan',  isCurrency: true }
        ] as item}
          <Card>
            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">{item.label}</p>
            <p class="text-2xl font-black text-gray-900">
              {#if item.isCurrency}
                {formatRupiah(Number(reportData[item.key] ?? 0))}
              {:else}
                {reportData[item.key] ?? 0}
                {#if item.suffix}
                  <span class="text-sm font-normal text-gray-500"> {item.suffix}</span>
                {/if}
              {/if}
            </p>
          </Card>
        {/each}
      </div>
    {:else}
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {#each ['Total Tagihan', 'Lunas', 'Belum Bayar', 'Total Pemasukan'] as label}
          <Card>
            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">{label}</p>
            <div class="h-8 bg-gray-100 animate-pulse rounded w-24"></div>
          </Card>
        {/each}
      </div>
    {/if}

    {#if error}
      <Alert type="error" message={error} class="no-print" />
    {/if}

    <Card padding={false} class="flex-1 flex flex-col min-h-0">
      <div class="p-5 border-b border-gray-100 space-y-4 no-print">
        <div class="flex flex-col lg:flex-row lg:items-end justify-between gap-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 flex-1">
            <Select id="filter-period" label="Periode" bind:value={filterPeriod} options={periodOptions} />
            
            {#if filterPeriod === 'monthly'}
              <Select id="filter-month" label="Bulan" bind:value={filterMonth} options={monthOptions} />
            {:else}
              <div class="hidden md:block"></div>
            {/if}

            <div>
              <label for="filter-year" class="text-xs font-semibold text-gray-600 uppercase tracking-wider ml-0.5 block mb-1.5">Tahun</label>
              <div class="relative">
                <select
                  id="filter-year"
                  bind:value={filterYear}
                  class="w-full pl-3 pr-8 py-2.5 text-sm font-semibold text-gray-700 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500 transition-all appearance-none cursor-pointer"
                >
                  <option value="">Pilih Tahun...</option>
                  {#each availableYears as yr}
                    <option value={String(yr)}>{yr} H</option>
                  {/each}
                </select>
                <div class="absolute inset-y-0 right-0 flex items-center pr-2.5 pointer-events-none text-gray-400">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
              </div>
            </div>

            <Select id="filter-category" label="Kategori" bind:value={filterCategory} options={categoryOptions} />
          </div>

          <div class="flex items-center gap-3 shrink-0 w-full lg:w-auto justify-end">
            <!-- Dropdown Export Button -->
            <div class="relative inline-block text-left w-full sm:w-auto">
              <button 
                onclick={() => showExportDropdown = !showExportDropdown}
                disabled={exporting}
                class="inline-flex items-center justify-center gap-2 w-full px-4 py-2.5 text-sm font-bold text-gray-700 bg-white border border-gray-300 rounded-lg shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500 transition-all disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                <Download size={16} />
                <span>Ekspor Laporan</span>
                <svg class="w-4 h-4 ml-1 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              
              {#if showExportDropdown}
                <!-- Backdrop to close dropdown -->
                <button class="fixed inset-0 z-10 cursor-default" onclick={() => showExportDropdown = false} aria-label="Close dropdown"></button>
                
                <div class="origin-top-right absolute right-0 mt-2 w-48 rounded-xl shadow-lg bg-white border border-gray-100 ring-1 ring-black ring-opacity-5 focus:outline-none z-20 divide-y divide-gray-50">
                  <div class="py-1">
                    <button 
                      onclick={() => { showExportDropdown = false; exportType = 'pdf'; showConfirmModal = true; }}
                      class="flex items-center gap-2.5 w-full px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 transition-colors text-left cursor-pointer"
                    >
                      <FileText size={16} class="text-red-500" />
                      <span>Ekspor PDF</span>
                    </button>
                    <button 
                      onclick={() => { showExportDropdown = false; exportType = 'excel'; showConfirmModal = true; }}
                      class="flex items-center gap-2.5 w-full px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 transition-colors text-left cursor-pointer"
                    >
                      <FileDown size={16} class="text-green-600" />
                      <span>Ekspor Excel / CSV</span>
                    </button>
                  </div>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>

      {#if loading && invoices.length === 0}
        <div class="py-16 flex justify-center items-center">
          <div class="w-8 h-8 border-4 border-green-500 border-t-transparent rounded-full animate-spin"></div>
        </div>
      {:else if invoices.length > 0}
        <div class="overflow-x-auto flex-1">
          <table class="w-full text-sm" aria-label="Tabel detail laporan keuangan">
            <thead>
              <tr class="bg-gray-50 border-b border-gray-100 whitespace-nowrap">
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">No. Invoice</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Santri</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Periode</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-right">Nominal</th>
                <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-center">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              {#each invoices as i}
                <tr class="hover:bg-gray-50 transition-colors">
                  <td class="px-5 py-3 font-mono text-xs text-green-600">{String(i.invoice_number)}</td>
                  <td class="px-5 py-3 text-gray-900">{String(i.student?.full_name ?? '-')}</td>
                  <td class="px-5 py-3 text-gray-500">{getHijriMonthName(Number(i.hijri_month))} {i.hijri_year} H</td>
                  <td class="px-5 py-3 text-right font-bold text-gray-900">{formatRupiah(Number(i.amount_due))}</td>
                  <td class="px-5 py-3 text-center">
                    <Badge label="Lunas" variant="success" />
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        
        {#if invPagination}
          <div class="no-print">
            <Paginator
              page={invPage}
              limit={invLimit}
              total={invPagination.total}
              pages={invPagination.pages}
              label="tagihan"
              onPageChange={handleInvPageChange}
              onLimitChange={handleInvLimitChange}
            />
          </div>
        {/if}
      {:else}
        <div class="py-16 flex flex-col items-center justify-center text-center px-4">
          <p class="text-gray-500 text-sm">Tidak ada data laporan tagihan lunas.</p>
        </div>
      {/if}
    </Card>
  </div>
</div>

<ConfirmDialog
  bind:open={showConfirmModal}
  title="Konfirmasi Ekspor Laporan"
  message={`Apakah Anda yakin ingin mengekspor laporan keuangan dalam format ${exportType === 'pdf' ? 'PDF' : 'Excel/CSV'} untuk periode yang dipilih?`}
  confirmText="Ekspor Sekarang"
  cancelText="Batal"
  variant="info"
  onConfirm={executeExport}
/>

<style>
  @media print {
    /* Hide everything except the print-area container */
    body * {
      visibility: hidden;
    }
    #print-area, #print-area * {
      visibility: visible;
    }
    #print-area {
      position: absolute;
      left: 0;
      top: 0;
      width: 100%;
      background: white;
    }
    .no-print {
      display: none !important;
      visibility: hidden !important;
    }
  }
</style>
