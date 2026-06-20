<script lang="ts">
  import { onDestroy, onMount, tick } from 'svelte';
  import {
    Users, CheckCircle, Moon, BookOpen, GraduationCap,
    ArrowUpRight, DollarSign, FileX, TrendingUp, Activity
  } from 'lucide-svelte';
  import { superUserApi, pengasuhApi } from '$lib/api';
  import { auth } from '$lib/stores/auth';
  import { formatRupiah } from '$lib/utils';
  import { Spinner, Alert, StatCard } from '$lib/components';
  import type { SuperUserDashboardStats } from '$lib/types';
  import Chart from 'chart.js/auto';

  const isReadOnly = $derived($auth.user?.role?.name === 'pengasuh');
  const api = $derived(isReadOnly ? pengasuhApi : superUserApi);
  const roleColor = $derived(isReadOnly ? 'purple' : 'green');
  const chartColor = $derived(isReadOnly ? '#7C3AED' : '#065F46');
  const chartBgColor = $derived(isReadOnly ? 'rgba(124, 58, 237, 0.06)' : 'rgba(6, 95, 70, 0.06)');

  let stats = $state<SuperUserDashboardStats | null>(null);
  let loading = $state(true);
  let error = $state('');

  let systemYear = $state<number | null>(null);
  let selectedYear = $state<number | null>(null);
  let selectedRange = $state<string>('1yr');

  let availableYears = $derived(stats?.available_years || []);

  let chartCanvas = $state<HTMLCanvasElement | null>(null);
  let chartInstance: Chart | null = null;

  onDestroy(() => {
    if (chartInstance) {
      chartInstance.destroy();
    }
  });

  async function loadDashboardData() {
    loading = true;
    error = '';
    try {
      const res = await api.getDashboard(selectedYear || undefined, selectedRange || undefined);
      stats = res.data;
      if (stats?.current_hijri) {
        if (systemYear === null) {
          systemYear = stats.current_hijri.hijri_year;
        }
      }
    } catch (err: any) {
      error = err.message || 'Gagal memuat data dashboard.';
    } finally {
      loading = false;
    }
  }

  function selectRange(range: string) {
    selectedRange = range;
    selectedYear = null;
    loadDashboardData();
  }

  function selectYear(year: number) {
    selectedYear = year;
    selectedRange = '';
    loadDashboardData();
  }

  // Reactive chart rendering — triggers when canvas is bound AND data is available
  $effect(() => {
    if (chartCanvas && stats?.monthly_payments) {
      // Use tick to ensure DOM is fully updated before rendering
      tick().then(() => renderChart());
    }
  });

  function renderChart() {
    if (!chartCanvas || !stats?.monthly_payments) return;
    
    if (chartInstance) {
      chartInstance.destroy();
    }
 
    const monthNames = ['Muharram', 'Safar', "Rabi'ul Awal", "Rabi'ul Akhir", 'Jumadil Awal', 'Jumadil Akhir', 'Rajab', "Sya'ban", 'Ramadhan', 'Syawal', "Dzulqa'dah", 'Dzulhijjah'];
    
    let sortedPayments = [...stats.monthly_payments];
    if (!selectedRange) {
      sortedPayments.sort((a, b) => a.month - b.month);
    }
    
    const labels = sortedPayments.map(mp => {
      const mName = monthNames[mp.month - 1] || '';
      return selectedRange ? `${mName} ${mp.year}` : mName;
    });
    
    const data = sortedPayments.map(mp => mp.total);

    chartInstance = new Chart(chartCanvas, {
      type: 'line',
      data: {
        labels: labels,
        datasets: [{
          label: 'Total Pembayaran (Rp)',
          data: data,
          borderColor: chartColor,
          backgroundColor: chartBgColor,
          borderWidth: 2,
          fill: true,
          tension: 0.3,
          pointBackgroundColor: '#fff',
          pointBorderColor: chartColor,
          pointBorderWidth: 2,
          pointRadius: 3,
          pointHoverRadius: 5,
          pointHoverBackgroundColor: chartColor,
          pointHoverBorderColor: '#fff',
          pointHoverBorderWidth: 2
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false },
          tooltip: {
            backgroundColor: '#0F172A',
            titleColor: '#F8FAFC',
            bodyColor: '#E2E8F0',
            borderColor: '#334155',
            borderWidth: 1,
            padding: 10,
            cornerRadius: 8,
            displayColors: false,
            callbacks: {
              label: (context) => formatRupiah(context.raw as number)
            }
          }
        },
        scales: {
          y: {
            beginAtZero: true,
            border: { display: false },
            grid: { display: true, color: '#F1F5F9' },
            ticks: {
              color: '#94A3B8',
              font: { size: 11 },
              callback: (value) => {
                const num = Number(value);
                if (num >= 1000000) return 'Rp ' + (num / 1000000) + ' Jt';
                if (num >= 1000) return 'Rp ' + (num / 1000) + ' Rb';
                return value;
              }
            }
          },
          x: {
            border: { display: false },
            grid: { display: false },
            ticks: {
              color: '#94A3B8',
              font: { size: 11 }
            }
          }
        },
        interaction: { intersect: false, mode: 'index' }
      }
    });
  }

  let hasInitialLoaded = false;
  $effect(() => {
    if (!$auth.loading && $auth.user && !hasInitialLoaded) {
      hasInitialLoaded = true;
      loadDashboardData();
    }
  });
</script>

<svelte:head>
  <title>Dashboard {isReadOnly ? 'Monitoring' : 'Keuangan'} | Al-Anwar Payment</title>
  <meta name="description" content="Pantau pembayaran, tagihan, dan statistik keuangan Pondok Pesantren Al-Anwar." />
</svelte:head>

<div class="space-y-6 flex-1 overflow-y-auto">
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div>
      {#if isReadOnly}
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-purple-100 border border-purple-200 text-purple-700 text-xs font-semibold uppercase tracking-wider mb-2">
          <span class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse" aria-hidden="true"></span>
          Mode Hanya Lihat
        </div>
      {/if}
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Dashboard {isReadOnly ? 'Monitoring' : 'Keuangan'}</h1>
      <p class="text-gray-500 text-sm mt-1">{isReadOnly ? 'Pantau pembayaran, tagihan, dan statistik santri secara real-time.' : 'Pantau pembayaran, tagihan, dan statistik semester pondok.'}</p>
    </div>
  </div>

  {#if error}
    <Alert type="error" message={error} />
  {/if}

  {#if loading}
    <div class="flex justify-center items-center py-20">
      <Spinner size="lg" />
    </div>
  {:else if stats}

    {#if stats.current_hijri}
      <div class="flex flex-wrap items-center justify-between gap-3 px-4 py-3 border {isReadOnly ? 'border-purple-100 bg-purple-50/40' : 'border-emerald-100 bg-emerald-50/40'} rounded-xl">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-lg {isReadOnly ? 'bg-purple-100' : 'bg-emerald-100'}">
            <Moon size={16} class={isReadOnly ? 'text-purple-800' : 'text-emerald-800'} />
          </div>
          <div>
            <p class="text-sm font-semibold {isReadOnly ? 'text-purple-800' : 'text-emerald-800'}">
              {stats.current_hijri.hijri_month_name} {stats.current_hijri.hijri_year} H
            </p>
            <p class="text-xs text-slate-500">{stats.current_hijri.semester_name} — TA {stats.current_hijri.academic_year_label}</p>
          </div>
        </div>
        <div class="flex gap-2">
          {#if stats.current_hijri.is_registration}
            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full border border-blue-200 bg-blue-50 text-blue-700 text-xs font-medium">
              <BookOpen size={12} /> Daftar Ulang
            </span>
          {/if}
          {#if stats.current_hijri.is_exam_month}
            <span class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full border border-amber-200 bg-amber-50 text-amber-600 text-xs font-medium">
              <GraduationCap size={12} /> Ujian
            </span>
          {/if}
        </div>
      </div>
    {/if}

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <StatCard
        title="Total Pemasukan Bulan Ini"
        value={formatRupiah(stats.total_income_mo)}
        subtitle="Pemasukan pembayaran berhasil"
        icon={DollarSign}
        color={isReadOnly ? 'purple' : 'emerald'}
        accent
      />
      <StatCard
        title="Tagihan Sudah Lunas"
        value={stats.paid_invoices}
        subtitle="Tagihan sudah dibayar"
        icon={CheckCircle}
        color={isReadOnly ? 'purple' : 'teal'}
      />
      <StatCard
        title="Tagihan Belum Lunas"
        value={stats.unpaid_invoices}
        subtitle="Tagihan belum diselesaikan"
        icon={FileX}
        color="amber"
      />
      <StatCard
        title="Total Santri Aktif"
        value={stats.total_students}
        subtitle="Santri terdaftar aktif"
        icon={Users}
        color={isReadOnly ? 'purple' : 'blue'}
      />
    </div>
    
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2">
        <div class="border border-slate-200/80 rounded-xl bg-white overflow-hidden">
          <div class="px-6 py-4 border-b border-slate-100 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div class="flex items-center gap-2">
              <Activity size={16} class={isReadOnly ? 'text-purple-700' : 'text-emerald-700'} />
              <h3 class="text-sm font-semibold text-slate-900">Pembayaran Berhasil</h3>
            </div>
            {#if stats}
              <div class="flex flex-wrap items-center gap-3">
                <!-- Range Button Group -->
                <div class="inline-flex rounded-lg bg-slate-100 p-0.5 border border-slate-200">
                  <button
                    type="button"
                    onclick={() => selectRange('6mo')}
                    class="px-2.5 py-1 text-[11px] font-bold rounded-md transition-all {selectedRange === '6mo' ? (isReadOnly ? 'bg-white text-purple-700 shadow-sm' : 'bg-white text-emerald-850 shadow-sm') : 'text-slate-500 hover:text-slate-800'}"
                  >
                    6 Bulan
                  </button>
                  <button
                    type="button"
                    onclick={() => selectRange('1yr')}
                    class="px-2.5 py-1 text-[11px] font-bold rounded-md transition-all {selectedRange === '1yr' ? (isReadOnly ? 'bg-white text-purple-700 shadow-sm' : 'bg-white text-emerald-850 shadow-sm') : 'text-slate-500 hover:text-slate-800'}"
                  >
                    1 Tahun
                  </button>
                  <button
                    type="button"
                    onclick={() => selectRange('3yr')}
                    class="px-2.5 py-1 text-[11px] font-bold rounded-md transition-all {selectedRange === '3yr' ? (isReadOnly ? 'bg-white text-purple-700 shadow-sm' : 'bg-white text-emerald-850 shadow-sm') : 'text-slate-500 hover:text-slate-800'}"
                  >
                    3 Tahun
                  </button>
                  <button
                    type="button"
                    onclick={() => selectRange('all')}
                    class="px-2.5 py-1 text-[11px] font-bold rounded-md transition-all {selectedRange === 'all' ? (isReadOnly ? 'bg-white text-purple-700 shadow-sm' : 'bg-white text-emerald-850 shadow-sm') : 'text-slate-500 hover:text-slate-800'}"
                  >
                    Semua
                  </button>
                </div>

                <!-- Year Dropdown -->
                <div class="relative min-w-[110px]">
                  <select 
                    id="dashboard-year"
                    value={selectedYear || ''}
                    onchange={(e) => {
                      const val = (e.target as HTMLSelectElement).value;
                      if (val) selectYear(Number(val));
                    }}
                    class="w-full pl-3 pr-8 py-1.5 text-xs font-semibold text-slate-700 bg-white border border-slate-200 rounded-lg focus:outline-none focus:ring-2 {isReadOnly ? 'focus:ring-purple-500/20 focus:border-purple-500' : 'focus:ring-emerald-850/20 focus:border-emerald-850'} transition-all appearance-none cursor-pointer"
                  >
                    <option value="" disabled={!!selectedRange}>Pilih Tahun</option>
                    {#each availableYears as yr}
                      <option value={yr}>{yr} H</option>
                    {/each}
                  </select>
                  <div class="absolute inset-y-0 right-0 flex items-center pr-2.5 pointer-events-none text-slate-400">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                    </svg>
                  </div>
                </div>
              </div>
            {/if}
          </div>
          <div class="p-5 h-[340px] w-full">
            <canvas bind:this={chartCanvas}></canvas>
          </div>
        </div>
      </div>

      <div class="lg:col-span-1 space-y-5">
        <div class="border border-slate-200/80 rounded-xl bg-white overflow-hidden">
          <div class="px-5 py-3 border-b border-slate-100 flex items-center gap-2">
            <TrendingUp size={16} class={isReadOnly ? 'text-purple-700' : 'text-emerald-700'} />
            <h3 class="text-sm font-semibold text-slate-900">Statistik Semester</h3>
          </div>
          {#if stats.semester_stats && stats.semester_stats.length > 0}
            <div class="divide-y divide-slate-100">
              {#each stats.semester_stats as sStat}
                <div class="px-5 py-4 space-y-2">
                  <div class="flex justify-between items-center">
                    <div>
                      <span class="text-sm font-semibold text-slate-900">{sStat.semester_name}</span>
                      {#if sStat.academic_year_label}
                        <p class="text-xs {isReadOnly ? 'text-purple-700' : 'text-emerald-700'}">TA {sStat.academic_year_label}</p>
                      {/if}
                    </div>
                    <span class="text-sm font-bold {isReadOnly ? 'text-purple-800' : 'text-emerald-800'}">{formatRupiah(sStat.total)}</span>
                  </div>
                  <div class="flex gap-3 text-xs">
                    <span class="text-slate-500">{sStat.invoice_count} tagihan</span>
                    <span class="{isReadOnly ? 'text-purple-700' : 'text-emerald-700'} font-medium">{sStat.paid_count} lunas</span>
                    {#if sStat.unpaid_count > 0}
                      <span class="text-amber-600 font-medium">{sStat.unpaid_count} belum</span>
                    {/if}
                  </div>
                  {#if sStat.invoice_count > 0}
                    <div class="w-full bg-slate-100 rounded-full h-1">
                      <div class={isReadOnly ? "bg-purple-600 h-1 rounded-full transition-all duration-500" : "bg-emerald-700 h-1 rounded-full transition-all duration-500"} style="width: {Math.round((sStat.paid_count / sStat.invoice_count) * 100)}%"></div>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {:else}
            <div class="px-5 py-10 text-center text-sm text-slate-400">
              Belum ada data semester berjalan.
            </div>
          {/if}
        </div>

        {#if !isReadOnly}
          <div class="border border-slate-200/80 rounded-xl bg-white p-5">
            <h3 class="text-sm font-semibold text-slate-900 mb-3">Aksi Cepat</h3>
            <div class="grid grid-cols-2 gap-3">
              <a
                href="/dashboard/super_user/billing"
                class="group flex items-center justify-between px-3.5 py-2.5 rounded-lg border border-emerald-100 bg-emerald-50/50 hover:bg-emerald-50 text-sm font-medium text-emerald-800 transition-colors duration-200"
              >
                Buat Tagihan
                <ArrowUpRight size={14} class="text-emerald-600 opacity-0 group-hover:opacity-100 transition-opacity" />
              </a>
              <a
                href="/dashboard/super_user/students"
                class="group flex items-center justify-between px-3.5 py-2.5 rounded-lg border border-slate-200 bg-slate-50/50 hover:bg-slate-50 text-sm font-medium text-slate-700 transition-colors duration-200"
              >
                Data Santri
                <ArrowUpRight size={14} class="text-slate-400 opacity-0 group-hover:opacity-100 transition-opacity" />
              </a>
            </div>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>
