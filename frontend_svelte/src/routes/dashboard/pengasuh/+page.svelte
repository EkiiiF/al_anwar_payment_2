<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import {
    Users, FileX, DollarSign, CheckCircle2,
    Moon, BookOpen, GraduationCap,
    TrendingUp, ArrowUpRight, Activity, Calendar
  } from 'lucide-svelte';
  import { pengasuhApi } from '$lib/api';
  import { formatRupiah, getMonthName } from '$lib/utils';
  import { Spinner, Alert, StatCard, Card, Badge } from '$lib/components';
  import type { SuperUserDashboardStats } from '$lib/types';
  import Chart from 'chart.js/auto';

  let stats = $state<SuperUserDashboardStats | null>(null);
  let loading = $state(true);
  let error = $state('');

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
      const res = await pengasuhApi.getDashboard();
      stats = res.data;
      
      if (stats?.monthly_payments) {
        setTimeout(renderChart, 0); 
      }
    } catch (err: any) {
      error = err.message || 'Gagal memuat data monitoring.';
    } finally {
      loading = false;
    }
  }

  function renderChart() {
    if (!chartCanvas || !stats?.monthly_payments) return;
    
    if (chartInstance) {
      chartInstance.destroy();
    }

    const months = ['Muharram', 'Safar', "Rabi'ul Awal", "Rabi'ul Akhir", 'Jumadil Awal', 'Jumadil Akhir', 'Rajab', "Sya'ban", 'Ramadhan', 'Syawal', "Dzulqa'dah", 'Dzulhijjah'];
    const data = stats.monthly_payments.sort((a, b) => a.month - b.month).map(mp => mp.total);

    chartInstance = new Chart(chartCanvas, {
      type: 'line',
      data: {
        labels: months,
        datasets: [{
          label: 'Total Pembayaran (Rp)',
          data: data,
          borderColor: '#7C3AED',
          backgroundColor: 'rgba(124, 58, 237, 0.06)',
          borderWidth: 2,
          fill: true,
          tension: 0.3,
          pointBackgroundColor: '#fff',
          pointBorderColor: '#7C3AED',
          pointBorderWidth: 2,
          pointRadius: 3,
          pointHoverRadius: 5,
          pointHoverBackgroundColor: '#7C3AED',
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
                if (value >= 1000000) return 'Rp ' + (value as number / 1000000) + ' Jt';
                if (value >= 1000) return 'Rp ' + (value as number / 1000) + ' Rb';
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

  onMount(loadDashboardData);
</script>

<svelte:head>
  <title>Dashboard Monitoring | Pengasuh - Al-Anwar Payment</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-purple-100 border border-purple-200 text-purple-700 text-xs font-semibold uppercase tracking-wider mb-2">
        <span class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse" aria-hidden="true"></span>
        Mode Hanya Lihat
      </div>
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Dashboard Monitoring Pengasuh</h1>
      <p class="text-gray-500 text-sm mt-1">Pantau perkembangan status akademik (Muhadhoroh), rekapitulasi dana masuk, dan data santri.</p>
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
      <Card class="border-purple-100 bg-gradient-to-r from-purple-50/40 to-indigo-50/40 shadow-sm">
        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
          <div class="flex items-start gap-4">
            <div class="p-2.5 rounded-xl bg-purple-100 border border-purple-200">
              <Moon size={20} class="text-purple-700" />
            </div>
            <div>
              <h2 class="font-bold text-gray-900">Periode Hijriah Aktif</h2>
              <p class="text-lg font-bold text-purple-800 mt-0.5">
                {stats.current_hijri.hijri_month_name} {stats.current_hijri.hijri_year} H
              </p>
              <p class="text-xs text-gray-500 mt-1">{stats.current_hijri.semester_name} — Tahun Ajaran {stats.current_hijri.academic_year_label}</p>
            </div>
          </div>
          <div class="flex flex-wrap gap-2">
            {#if stats.current_hijri.is_registration}
              <span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-blue-100 text-blue-700 text-xs font-bold">
                <BookOpen size={14} /> Daftar Ulang
              </span>
            {/if}
            {#if stats.current_hijri.is_exam_month}
              <span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-amber-100 text-amber-700 text-xs font-bold">
                <GraduationCap size={14} /> Ujian Semester
              </span>
            {/if}
            <span class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-purple-100 text-purple-700 text-xs font-bold">
              Semester {stats.current_hijri.semester}
            </span>
          </div>
        </div>
      </Card>
    {/if}

    <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
      <StatCard title="Total Santri"          value={stats.total_students}              subtitle="Santri terdaftar aktif"    icon={Users}        color="blue" />
      <StatCard title="Tagihan Belum Dibayar"  value={stats.unpaid_invoices}             subtitle="Faktur belum diselesaikan"  icon={FileX}        color="amber" />
      <StatCard title="Tagihan Lunas"         value={stats.paid_invoices}               subtitle="Faktur berhasil dibayar"   icon={CheckCircle2} color="green" />
      <StatCard
        title="Uang Masuk Bulan Ini"
        value={formatRupiah(stats.total_income_mo)}
        subtitle="Dana terkumpul bulan ini"
        icon={DollarSign}
        color="purple"
        accent
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2">
        <div class="border border-slate-200/80 rounded-xl bg-white overflow-hidden">
          <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
            <h3 class="text-sm font-semibold text-slate-900 flex items-center gap-2">
              <Activity size={16} class="text-purple-600" />
              Tren Pembayaran Bulanan
            </h3>
            <span class="text-xs text-gray-400">Tahun Ini</span>
          </div>
          <div class="p-5 h-[340px] w-full">
            <canvas bind:this={chartCanvas}></canvas>
          </div>
        </div>
      </div>

      <div class="lg:col-span-1 space-y-5">
        <div class="border border-slate-200/80 rounded-xl bg-white overflow-hidden">
          <div class="px-5 py-3 border-b border-slate-100 flex items-center gap-2">
            <TrendingUp size={16} class="text-purple-600" />
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
                        <p class="text-xs text-purple-700">TA {sStat.academic_year_label}</p>
                      {/if}
                    </div>
                    <span class="text-sm font-bold text-purple-800">{formatRupiah(sStat.total)}</span>
                  </div>
                  <div class="flex gap-3 text-xs">
                    <span class="text-slate-500">{sStat.invoice_count} tagihan</span>
                    <span class="text-emerald-700 font-medium">{sStat.paid_count} lunas</span>
                    {#if sStat.unpaid_count > 0}
                      <span class="text-amber-600 font-medium">{sStat.unpaid_count} belum</span>
                    {/if}
                  </div>
                  {#if sStat.invoice_count > 0}
                    <div class="w-full bg-slate-100 rounded-full h-1">
                      <div class="bg-purple-600 h-1 rounded-full transition-all duration-500" style="width: {Math.round((sStat.paid_count / sStat.invoice_count) * 100)}%"></div>
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

        <div class="border border-slate-200/80 rounded-xl bg-white p-5">
          <h3 class="text-sm font-semibold text-slate-900 mb-3">Menu Monitoring</h3>
          <div class="grid grid-cols-2 gap-3">
            <a
              href="/dashboard/pengasuh/students"
              class="group flex items-center justify-between px-3.5 py-2.5 rounded-lg border border-purple-100 bg-purple-50/50 hover:bg-purple-50 text-sm font-medium text-purple-800 transition-colors duration-200"
            >
              Data Santri
              <ArrowUpRight size={14} class="text-purple-600 opacity-0 group-hover:opacity-100 transition-opacity" />
            </a>
            <a
              href="/dashboard/pengasuh/reports"
              class="group flex items-center justify-between px-3.5 py-2.5 rounded-lg border border-slate-200 bg-slate-50/50 hover:bg-slate-50 text-sm font-medium text-slate-700 transition-colors duration-200"
            >
              Laporan Keuangan
              <ArrowUpRight size={14} class="text-slate-400 opacity-0 group-hover:opacity-100 transition-opacity" />
            </a>
          </div>
        </div>
      </div>
    </div>

    <Card class="border-purple-200 bg-purple-50/50">
      <div class="flex items-start gap-3">
        <div class="w-2 h-2 rounded-full bg-purple-500 mt-1.5 flex-shrink-0"></div>
        <p class="text-sm text-purple-800">
          <strong>Akses Monitoring Pengasuh:</strong> Anda berada dalam mode peninjauan laporan keuangan, status semester santri (Muhadhoroh), dan rekapitulasi dana masuk pesantren. Seluruh aksi edit, hapus, tambah, dan generate hanya dapat diubah oleh <strong>Super User (Bendahara)</strong>.
        </p>
      </div>
    </Card>
  {/if}
</div>
