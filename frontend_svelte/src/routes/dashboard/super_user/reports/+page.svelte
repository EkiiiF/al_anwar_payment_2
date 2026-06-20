<script lang="ts">
  import { onMount, tick } from 'svelte';
  import logo from '$lib/assets/logo.png';
  import { BarChart3, Download, FileText, FileDown } from 'lucide-svelte';
  import { superUserApi, pengasuhApi } from '$lib/api';
  import { auth } from '$lib/stores/auth';
  import { formatRupiah, getHijriMonthName, HIJRI_MONTH_NAMES } from '$lib/utils';
  import { Button, Alert, Card, Badge, Select, Paginator, ConfirmDialog } from '$lib/components';
  import { toast } from '$lib/stores/toast';
  import type { Category } from '$lib/types';

  const isReadOnly = $derived($auth.user?.role?.name === 'pengasuh');
  const api = $derived(isReadOnly ? pengasuhApi : superUserApi);

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
        api.getReports(reportFilters),
        api.getInvoicesPaginated(invoiceFilters, invPage, invLimit)
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
      const res = await api.getInvoicesPaginated(invoiceFilters, invPage, invLimit);
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

  let printInvoices = $state<any[]>([]);

  function calculateRunningSums(invs: any[]): number[] {
    let sum = 0;
    return invs.map(i => {
      sum += Number(i.amount_due);
      return sum;
    });
  }

  function formatDate(dateStr: string): string {
    if (!dateStr) return '-';
    try {
      const d = new Date(dateStr);
      return d.toLocaleDateString('id-ID', { day: '2-digit', month: '2-digit', year: 'numeric' });
    } catch {
      return '-';
    }
  }

  // Reactive effect to load data when filters change
  $effect(() => {
    if ($auth.loading || !$auth.user) {
      return;
    }
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
      exporting = true;
      try {
        const invoiceFilters: Record<string, string> = {
          status: 'paid',
          ...(filterYear && { year: filterYear }),
          ...(filterPeriod === 'monthly' && filterMonth && { month: filterMonth }),
          ...(filterCategory && { category: filterCategory })
        };
        // Fetch ALL invoices for printing
        const totalPaid = reportData?.paid_count ?? 1000;
        const res = await api.getInvoicesPaginated(invoiceFilters, 1, totalPaid);
        printInvoices = res.data?.invoices ?? [];

        await tick();
        // Allow Svelte to render the print DOM
        setTimeout(() => {
          window.print();
          exporting = false;
        }, 300);
      } catch (e: any) {
        toast.error('Gagal mengambil data untuk cetak: ' + (e.message || 'Error'));
        exporting = false;
      }
    } else {
      exporting = true;
      try {
        await api.exportReports({
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

  let hasInitialLoaded = false;
  $effect(() => {
    if (!$auth.loading && $auth.user && !hasInitialLoaded) {
      hasInitialLoaded = true;
      (async () => {
        try {
          const [dashboardRes, catRes] = await Promise.all([
            api.getDashboard(),
            api.getCategories()
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
      })();
    }
  });
</script>

<svelte:head>
  <title>Laporan Keuangan | {isReadOnly ? 'Dashboard Pengasuh' : 'Dashboard Super User'}</title>
</svelte:head>

<div class="space-y-6 flex flex-col flex-1 min-h-0">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 no-print">
    <div>
      {#if isReadOnly}
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-purple-100 border border-purple-200 text-purple-700 text-xs font-semibold uppercase tracking-wider mb-2">
          <span class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse" aria-hidden="true"></span>
          Mode Hanya Lihat
        </div>
      {/if}
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Laporan Keuangan</h1>
      <p class="text-gray-500 text-sm mt-1">Analisis pemasukan pesantren secara real-time.</p>
    </div>
  </div>

  <div id="print-area" class="space-y-6 flex flex-col flex-1 min-h-0">
    <!-- Printable Document (Only shown in print/PDF) -->
    <div class="hidden print:block w-full text-black font-sans leading-relaxed">
      <!-- Kop Surat / Header Laporan -->
      <div class="flex items-center justify-between border-b-4 border-double border-slate-900 pb-4 mb-6">
        <div class="flex items-center gap-4 text-left">
          <img src={logo} alt="Logo Al-Anwar" class="w-16 h-16 object-contain" />
          <div>
            <h2 class="text-sm font-black uppercase tracking-wider text-slate-900">PONDOK PESANTREN PUTRA - PUTRI "AL-ANWAR"</h2>
            <p class="text-[10px] text-slate-600 leading-tight mt-1">
              Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan<br />
            </p>
          </div>
        </div>
      </div>

      <!-- Judul & Periode Laporan -->
      <div class="text-center mb-6">
        <h3 class="text-xs font-bold uppercase text-slate-800 tracking-wider">Laporan Realisasi Pemasukan Administrasi Keuangan</h3>
        <p class="text-[10px] text-slate-500 mt-1">
          Periode: 
          {#if filterPeriod}
            {filterPeriod === 'monthly' ? 'Bulanan' : filterPeriod === 'quarterly' ? 'Triwulan' : filterPeriod === 'semester' ? 'Semester' : 'Tahunan'} 
            {#if filterPeriod === 'monthly' && filterMonth}
              - {getHijriMonthName(Number(filterMonth))}
            {/if}
            - {filterYear} H
          {:else}
            Semua Periode
          {/if}
          {#if filterCategory}
            | Kategori: {categories.find(c => c.id === filterCategory)?.name ?? 'Semua'}
          {/if}
        </p>
      </div>

      <!-- Ringkasan / Rekapitulasi Metrik -->
      <div class="grid grid-cols-4 gap-4 mb-6 border border-slate-300 rounded-lg p-3 bg-slate-50/50">
        <div class="text-center border-r border-slate-200">
          <p class="text-[9px] uppercase font-bold text-slate-500">Total Tagihan</p>
          <p class="text-xs font-extrabold text-slate-800 mt-0.5">{reportData?.total_invoices ?? 0} Tagihan</p>
        </div>
        <div class="text-center border-r border-slate-200">
          <p class="text-[9px] uppercase font-bold text-slate-500">Tagihan Lunas</p>
          <p class="text-xs font-extrabold text-emerald-700 mt-0.5">{reportData?.paid_count ?? 0} Tagihan</p>
        </div>
        <div class="text-center border-r border-slate-200">
          <p class="text-[9px] uppercase font-bold text-slate-500">Belum Lunas</p>
          <p class="text-xs font-extrabold text-amber-700 mt-0.5">{reportData?.unpaid_count ?? 0} Tagihan</p>
        </div>
        <div class="text-center">
          <p class="text-[9px] uppercase font-bold text-slate-500">Total Keseluruhan</p>
          <p class="text-xs font-extrabold text-emerald-800 mt-0.5">{formatRupiah(Number(reportData?.total_amount_paid ?? 0))}</p>
        </div>
      </div>

      <!-- Tabel Utama Laporan Jurnal/Transaksi -->
      <table class="w-full text-[9px] border-collapse border border-slate-300 mb-8">
        <thead>
          <tr class="bg-slate-800 text-white border border-slate-300">
            <th scope="col" class="border border-slate-300 px-3 py-2 text-center w-8 font-bold">No</th>
            <th scope="col" class="border border-slate-300 px-3 py-2 text-left w-20 font-bold">Tanggal</th>
            <th scope="col" class="border border-slate-300 px-3 py-2 text-left w-24 font-bold">No. Referensi</th>
            <th scope="col" class="border border-slate-300 px-3 py-2 text-left w-28 font-bold">Kategori</th>
            <th scope="col" class="border border-slate-300 px-3 py-2 text-left font-bold">Nama Santri</th>
            <th scope="col" class="border border-slate-300 px-3 py-2 text-right w-24 font-bold">Nominal</th>
          </tr>
        </thead>
        <tbody>
          {#if printInvoices.length > 0}
            {#each printInvoices as i, idx}
              <tr class="border border-slate-200 {idx % 2 === 1 ? 'bg-slate-50/70' : 'bg-white'}">
                <td class="border border-slate-300 px-3 py-2 text-center">{idx + 1}</td>
                <td class="border border-slate-300 px-3 py-2 whitespace-nowrap">{formatDate(i.created_at)}</td>
                <td class="border border-slate-300 px-3 py-2 font-mono">{i.invoice_number}</td>
                <td class="border border-slate-300 px-3 py-2">{i.category?.name ?? 'SPP Pesantren'}</td>
                <td class="border border-slate-300 px-3 py-2 font-semibold text-slate-800">{i.student?.full_name ?? '-'}</td>
                <td class="border border-slate-300 px-3 py-2 text-right text-slate-900 font-semibold">{formatRupiah(Number(i.amount_due))}</td>
              </tr>
            {/each}
            <!-- Baris Total Footer -->
            <tr class="bg-slate-100 font-bold border-t border-slate-900">
              <td colspan="5" class="border border-slate-300 px-3 py-2.5 text-right font-extrabold text-slate-800">TOTAL KESELURUHAN</td>
              <td class="border border-slate-300 px-3 py-2.5 text-right text-slate-900 font-extrabold underline decoration-double">{formatRupiah(Number(reportData?.total_amount_paid ?? 0))}</td>
            </tr>
          {:else}
            <tr>
              <td colspan="6" class="border border-slate-300 px-3 py-6 text-center text-slate-500">Tidak ada data untuk periode ini.</td>
            </tr>
          {/if}
        </tbody>
      </table>

      <!-- Kolom Tanda Tangan -->
      <div class="mt-8 grid grid-cols-2 gap-8 text-center text-[10px] avoid-break">
        <div>
          <p class="text-slate-500">Mengetahui,</p>
          <p class="font-bold text-slate-800 uppercase mt-0.5">Pengasuh Pondok Pesantren</p>
          <div class="h-16"></div> 
          <div class="mx-auto w-48 border-b border-slate-800 font-bold"></div>
        </div>
        <div>
          <p class="text-slate-500">Grobogan, {new Date().toLocaleDateString('id-ID', { dateStyle: 'long' })}</p>
          <p class="font-bold text-slate-800 uppercase mt-0.5">Bendahara</p>
          <div class="h-16"></div>
          <div class="mx-auto w-48 border-b border-slate-800 font-bold"></div>
        </div>
      </div>
    </div>

    {#if reportData}
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 no-print">
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

    <Card padding={false} class="flex-1 flex flex-col min-h-0 no-print">
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
                  class="w-full pl-3 pr-8 py-2.5 text-sm font-semibold text-gray-700 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 {isReadOnly ? 'focus:ring-purple-500/20 focus:border-purple-500' : 'focus:ring-green-500/20 focus:border-green-500'} transition-all appearance-none cursor-pointer"
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
                class="inline-flex items-center justify-center gap-2 w-full px-4 py-2.5 text-sm font-bold text-gray-700 bg-white border border-gray-300 rounded-lg shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 {isReadOnly ? 'focus:ring-purple-500/20 focus:border-purple-500' : 'focus:ring-green-500/20 focus:border-green-500'} transition-all disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
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
          <div class="w-8 h-8 border-4 {isReadOnly ? 'border-purple-500' : 'border-green-500'} border-t-transparent rounded-full animate-spin"></div>
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
                  <td class="px-5 py-3 font-mono text-xs {isReadOnly ? 'text-purple-600' : 'text-green-600'}">{String(i.invoice_number)}</td>
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
