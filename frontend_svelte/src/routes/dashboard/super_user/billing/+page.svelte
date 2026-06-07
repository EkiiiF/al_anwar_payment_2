<script lang="ts">
  import { onMount } from 'svelte';
  import { PlayCircle, Receipt, Filter, Search, Calendar, CheckCircle2, Clock, BookOpen, GraduationCap, Moon, Layers, ChevronLeft, ChevronRight, User, X, FileText, MapPin, Phone, Mail, Eye, EyeOff, ChevronDown, ChevronUp, Zap, Settings2, PanelTopClose, PanelTop } from 'lucide-svelte';
  import { superUserApi } from '$lib/api';
  import { formatRupiah, formatDate, getMonthName, MONTH_NAMES, getHijriMonthName } from '$lib/utils';
  import { Button, Alert, Spinner, Badge, EmptyState, Card, Select, Modal, BillingGenerator, AutoBillingSettings, BillingFilters, StudentBillingCard, StudentBillingModal } from '$lib/components';
  import type { Invoice, HijriMonthInfo, Student } from '$lib/types';

  let students   = $state<Student[]>([]);
  let loading    = $state(true);
  let error      = $state('');
  
  let page       = $state(1);
  let limit      = $state(10);
  let pagination = $state<{ page: number; limit: number; total: number; pages: number } | null>(null);

  let selectedStudent = $state<Student | null>(null);
  let isModalOpen = $state(false);

  let hijriInfo  = $state<HijriMonthInfo | null>(null);

  let statusFilter = $state('');
  let monthFilter  = $state('');
  let yearFilter   = $state('');

  let genMode    = $state<'monthly' | 'semester'>('monthly');
  
  let genMonth   = $state('1'); // Default Muharram (1)
  let genYear    = $state('');

  let genSemester  = $state('1');
  let genHijriYear = $state('');

  let generating = $state(false);
  let genMessage = $state('');
  let genSuccess = $state(false);
  
  let autoBillingEnabled         = $state(false);
  let autoSemesterBillingEnabled = $state(false);
  let billingHijriDay            = $state('1');
  let billingHijriMonth          = $state('0');
  let billingHijriYear           = $state('0');
  let settingLoading             = $state(false);
  let showBillingPanel           = $state(true);

  const monthOptions = MONTH_NAMES.map((name, i) => ({ value: String(i + 1), label: name }));
  const yearOptions  = Array.from({ length: 5 }, (_, i) => {
    const y = new Date().getFullYear() - 2 + i;
    return { value: String(y), label: String(y) };
  });

  async function fetchData() {
    loading = true;
    error   = '';
    try {
      const filters: any = {};
      if (statusFilter) filters.status = statusFilter;
      if (monthFilter)  filters.month  = monthFilter;
      if (yearFilter)   filters.year   = yearFilter;
      
      const res = await superUserApi.getStudentsWithInvoicesPaginated(filters, page, limit);
      students  = res.data?.students ?? [];
      pagination = res.data?.pagination ?? null;
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat data tagihan.';
    } finally {
      loading = false;
    }
  }

  async function fetchSettings() {
    try {
      const res = await superUserApi.getSettings();
      const settings = res.data as any[];
      const autoBill = settings.find(s => s.key === 'auto_generate_billing');
      if (autoBill) autoBillingEnabled = autoBill.value === 'true';
      const autoSem = settings.find(s => s.key === 'auto_generate_semester_billing');
      if (autoSem) autoSemesterBillingEnabled = autoSem.value === 'true';
      const hijriDay = settings.find(s => s.key === 'billing_hijri_day');
      if (hijriDay) billingHijriDay = hijriDay.value;
      const hijriMonth = settings.find(s => s.key === 'billing_hijri_month');
      if (hijriMonth) billingHijriMonth = hijriMonth.value;
      const hijriYear = settings.find(s => s.key === 'billing_hijri_year');
      if (hijriYear) billingHijriYear = hijriYear.value;
    } catch (e) { console.error('Gagal memuat setting', e); }
  }

  async function fetchHijriInfo() {
    try {
      const res = await superUserApi.getDashboard();
      if (res.data?.current_hijri) {
        hijriInfo = res.data.current_hijri;
        if (!genHijriYear && hijriInfo) {
          genHijriYear = String(hijriInfo.hijri_year);
        }
        if (!genYear && hijriInfo) {
          genYear = String(hijriInfo.hijri_year);
        }
        genMonth = String(hijriInfo.hijri_month);
      }
    } catch (e) { console.error('Gagal memuat info Hijriah', e); }
  }

  onMount(() => {
    fetchData();
    fetchSettings();
    fetchHijriInfo();
  });

  $effect(() => {
    statusFilter;
    monthFilter;
    yearFilter;
    page = 1;
    fetchData();
  });
  
  $effect(() => {
    page;
    fetchData();
  });

  function openStudentModal(student: Student) {
    console.log('Opening modal for student:', student);
    selectedStudent = student;
    isModalOpen = true;
  }

  function closeStudentModal() {
    isModalOpen = false;
    selectedStudent = null;
  }

  function prevPage() {
    if (page > 1) {
      page--;
    }
  }

  function nextPage() {
    if (pagination && page < pagination.pages) {
      page++;
    }
  }

  async function handleGenerate() {
    if (genMode === 'monthly') {
      await handleGenerateMonthly();
    } else {
      await handleGenerateSemester();
    }
  }

  async function handleGenerateMonthly() {
    if (!confirm(`Generate tagihan Syahriyyah Pondok & Muhadhoroh untuk bulan ${getHijriMonthName(Number(genMonth))} ${genYear} H?`)) return;

    generating = true;
    genMessage  = '';
    try {
      await superUserApi.generateInvoices(Number(genMonth), Number(genYear));
      genMessage = 'Tagihan Syahriyyah Pondok & Muhadhoroh berhasil dibuat!';
      genSuccess = true;
      await fetchData();
    } catch (e: unknown) {
      genMessage = e instanceof Error ? e.message : 'Gagal generate tagihan.';
      genSuccess = false;
    } finally {
      generating = false;
    }
  }

  async function handleGenerateSemester() {
    const yr = Number(genHijriYear) || (hijriInfo?.hijri_year ?? 0);
    if (!yr) {
      genMessage = 'Tahun Hijriah belum terdeteksi.';
      genSuccess = false;
      return;
    }
    if (!confirm(`Generate seluruh tagihan Semester ${genSemester} untuk tahun ${yr} H?\nSemua tagihan bulan di semester ini akan dibuat sekaligus. Wali santri dapat membayar paket atau dicicil per bulan.`)) return;

    generating = true;
    genMessage  = '';
    try {
      await superUserApi.generateSemesterInvoices(Number(genSemester), yr);
      genMessage = `Tagihan Semester ${genSemester} berhasil dibuat!`;
      genSuccess = true;
      await fetchData();
    } catch (e: unknown) {
      genMessage = e instanceof Error ? e.message : 'Gagal generate tagihan semester.';
      genSuccess = false;
    } finally {
      generating = false;
    }
  }

  async function toggleSetting(key: string) {
    settingLoading = true;
    try {
      if (key === 'auto_generate_billing') {
        const nextValue = !autoBillingEnabled;
        await superUserApi.updateSetting(key, String(nextValue));
        autoBillingEnabled = nextValue;
      } else {
        const nextValue = !autoSemesterBillingEnabled;
        await superUserApi.updateSetting(key, String(nextValue));
        autoSemesterBillingEnabled = nextValue;
      }
    } catch (e: unknown) {
      console.error(e instanceof Error ? e.message : 'Gagal mengubah setting.');
    } finally {
      settingLoading = false;
    }
  }

  async function updateHijriSetting(key: string, value: string) {
    settingLoading = true;
    try {
      await superUserApi.updateSetting(key, value);
    } catch (e: unknown) {
      console.error(e instanceof Error ? e.message : `Gagal mengubah setting ${key}.`);
    } finally {
      settingLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Manajemen Tagihan Semester | Dashboard Super User</title>
</svelte:head>

<div class="space-y-6 flex flex-col flex-1 min-h-0">
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div>
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Manajemen Tagihan Semester</h1>
      <p class="text-gray-500 text-sm mt-1">Generate tagihan berdasarkan kalender Hijriah dan pantau status pembayaran per semester.</p>
    </div>
    <div class="flex items-center gap-3">
      <button
        onclick={() => showBillingPanel = !showBillingPanel}
        class="billing-toggle-btn {showBillingPanel ? 'billing-toggle-btn--active' : ''}"
      >
        <span class="billing-toggle-btn__icon">
          {#if showBillingPanel}
            <PanelTopClose size={16} />
          {:else}
            <PanelTop size={16} />
          {/if}
        </span>
        <span class="billing-toggle-btn__label">
          {showBillingPanel ? 'Sembunyikan Panel' : 'Tampilkan Panel'}
        </span>
        <span class="billing-toggle-btn__indicator">
          {#if showBillingPanel}
            <ChevronUp size={14} />
          {:else}
            <ChevronDown size={14} />
          {/if}
        </span>
      </button>
    </div>
  </div>

  {#if error}
    <Alert type="error" message={error} />
  {/if}

  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
    <Card class="border-emerald-200 bg-emerald-50/20">
      <div class="flex items-start gap-3">
        <div class="p-2 rounded-lg bg-emerald-100 flex-shrink-0">
          <Moon size={20} class="text-emerald-700" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-xs font-bold text-emerald-700 uppercase tracking-wider">Kalender Hijriah</p>
          {#if hijriInfo}
            <p class="text-lg font-black text-emerald-900 mt-1 truncate" title="{hijriInfo.hijri_month_name} {hijriInfo.hijri_year} H">
              {hijriInfo.hijri_month_name} {hijriInfo.hijri_year} H
            </p>
            <p class="text-xs text-emerald-600 mt-0.5 truncate">
              Sem {hijriInfo.semester} ({hijriInfo.is_exam_month ? 'Ujian' : 'Reguler'})
            </p>
          {:else}
            <p class="text-xs text-gray-400 mt-1">Loading...</p>
          {/if}
        </div>
      </div>
    </Card>

    <Card class="border-emerald-200 bg-emerald-50/30">
      <div class="flex items-start gap-4">
        <div class="p-2 rounded-lg bg-emerald-100 flex-shrink-0">
          <Receipt size={20} class="text-emerald-700" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-xs font-bold text-emerald-700 uppercase tracking-wider">Total Tagihan</p>
          <p class="text-2xl font-black text-emerald-900 mt-1">
            {students.reduce((sum, s) => sum + (s.invoices?.length || 0), 0)}
          </p>
          <p class="text-xs text-emerald-600 mt-0.5 truncate">Tagihan di-generate</p>
        </div>
      </div>
    </Card>

    <Card class="border-amber-200 bg-amber-50/30">
      <div class="flex items-start gap-4">
        <div class="p-2 rounded-lg bg-amber-100 flex-shrink-0">
          <Clock size={20} class="text-amber-700" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-xs font-bold text-amber-700 uppercase tracking-wider">Belum Bayar</p>
          <p class="text-2xl font-black text-amber-900 mt-1">
            {students.filter(s => s.invoices?.some(i => i.status === 'unpaid')).length}
          </p>
          <p class="text-xs text-amber-600 mt-0.5 truncate">Santri belum lunas</p>
        </div>
      </div>
    </Card>

    <Card class="border-green-200 bg-green-50/30">
      <div class="flex items-start gap-4">
        <div class="p-2 rounded-lg bg-green-100 flex-shrink-0">
          <CheckCircle2 size={20} class="text-green-700" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-xs font-bold text-green-700 uppercase tracking-wider">Lunas</p>
          <p class="text-2xl font-black text-green-900 mt-1">
            {students.filter(s => s.invoices && s.invoices.length > 0 && s.invoices.every(i => i.status === 'paid')).length}
          </p>
          <p class="text-xs text-green-600 mt-0.5 truncate">Santri lunas semua</p>
        </div>
      </div>
    </Card>
  </div>

  <div class="billing-panel-wrapper {showBillingPanel ? 'billing-panel-wrapper--open' : 'billing-panel-wrapper--closed'}">
    <div class="billing-panel-inner">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <BillingGenerator
          bind:genMode={genMode}
          bind:genMonth={genMonth}
          bind:genYear={genYear}
          bind:genSemester={genSemester}
          bind:genHijriYear={genHijriYear}
          generating={generating}
          genMessage={genMessage}
          genSuccess={genSuccess}
          onGenerate={handleGenerate}
        />

        <AutoBillingSettings
          bind:autoBillingEnabled={autoBillingEnabled}
          bind:autoSemesterBillingEnabled={autoSemesterBillingEnabled}
          bind:billingHijriDay={billingHijriDay}
          settingLoading={settingLoading}
          onToggle={toggleSetting}
          onDayChange={() => updateHijriSetting('billing_hijri_day', billingHijriDay)}
        />
      </div>
    </div>
  </div>

  <BillingFilters
    bind:statusFilter={statusFilter}
    bind:monthFilter={monthFilter}
    bind:yearFilter={yearFilter}
    monthOptions={monthOptions}
    yearOptions={yearOptions}
  />

  {#if loading}
    <Spinner size="lg" />
  {:else}
    <div class="space-y-4 flex-1">
      {#each students as student (student.id)}
        <StudentBillingCard
          student={student}
          onclick={() => openStudentModal(student)}
        />
      {:else}
        <EmptyState title="Tidak ada santri dengan tagihan" description="Coba ubah filter atau generate tagihan baru." />
      {/each}
      
      {#if pagination && pagination.pages > 1}
        <div class="flex items-center justify-between bg-white p-4 rounded-xl border border-slate-200 shadow-sm">
          <p class="text-sm text-slate-650">
            Menampilkan <span class="font-bold text-slate-900">{((page - 1) * limit) + 1} - {Math.min(page * limit, pagination.total)}</span> dari <span class="font-bold text-slate-900">{pagination.total}</span> santri
          </p>
          <div class="flex items-center gap-2">
            <Button variant="secondary" onclick={prevPage} disabled={page <= 1} size="sm" class="text-xs font-semibold px-3 py-1.5 border border-slate-200">
              {#snippet children()}<ChevronLeft size={16} /> Sebelumnya{/snippet}
            </Button>
            <div class="px-3 py-1.5 bg-emerald-50 text-emerald-800 rounded-lg text-sm font-bold border border-emerald-100">
              Halaman {page} / {pagination.pages}
            </div>
            <Button variant="secondary" onclick={nextPage} disabled={!pagination || page >= pagination.pages} size="sm" class="text-xs font-semibold px-3 py-1.5 border border-slate-200">
              {#snippet children()}Selanjutnya <ChevronRight size={16} />{/snippet}
            </Button>
          </div>
        </div>
      {/if}
    </div>
  {/if}

  <StudentBillingModal
    open={isModalOpen}
    onclose={closeStudentModal}
    selectedStudent={selectedStudent}
  />
</div>

<style>
  /* ═══════════════════════════════════════════════════
     Modern Billing Panel Toggle Button
     ═══════════════════════════════════════════════════ */
  .billing-toggle-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    border-radius: 12px;
    border: 1.5px solid #e2e8f0;
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
    color: #475569;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
    white-space: nowrap;
    position: relative;
    overflow: hidden;
  }

  .billing-toggle-btn::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    background: linear-gradient(135deg, rgba(16, 185, 129, 0.08) 0%, rgba(5, 150, 105, 0.04) 100%);
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .billing-toggle-btn:hover {
    border-color: #a7f3d0;
    color: #065f46;
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.12), 0 1px 3px rgba(0, 0, 0, 0.06);
    transform: translateY(-1px);
  }

  .billing-toggle-btn:hover::before {
    opacity: 1;
  }

  .billing-toggle-btn:active {
    transform: translateY(0);
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  }

  .billing-toggle-btn--active {
    background: linear-gradient(135deg, #ecfdf5 0%, #d1fae5 100%);
    border-color: #6ee7b7;
    color: #065f46;
    box-shadow: 0 2px 8px rgba(16, 185, 129, 0.15), inset 0 1px 0 rgba(255, 255, 255, 0.6);
  }

  .billing-toggle-btn--active::before {
    opacity: 0;
  }

  .billing-toggle-btn--active:hover {
    background: linear-gradient(135deg, #d1fae5 0%, #a7f3d0 100%);
    border-color: #34d399;
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.2), inset 0 1px 0 rgba(255, 255, 255, 0.6);
  }

  .billing-toggle-btn__icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 8px;
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
    color: white;
    box-shadow: 0 2px 6px rgba(16, 185, 129, 0.3);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    flex-shrink: 0;
  }

  .billing-toggle-btn:hover .billing-toggle-btn__icon {
    box-shadow: 0 3px 10px rgba(16, 185, 129, 0.4);
    transform: scale(1.05);
  }

  .billing-toggle-btn__label {
    position: relative;
    z-index: 1;
    letter-spacing: 0.01em;
  }

  .billing-toggle-btn__indicator {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border-radius: 6px;
    background: rgba(0, 0, 0, 0.04);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    flex-shrink: 0;
  }

  .billing-toggle-btn--active .billing-toggle-btn__indicator {
    background: rgba(16, 185, 129, 0.12);
    color: #059669;
  }

  .billing-toggle-btn:hover .billing-toggle-btn__indicator {
    background: rgba(16, 185, 129, 0.15);
  }

  /* ═══════════════════════════════════════════════════
     Billing Panel Slide Animation
     ═══════════════════════════════════════════════════ */
  .billing-panel-wrapper {
    display: grid;
    transition: grid-template-rows 0.45s cubic-bezier(0.4, 0, 0.2, 1),
                opacity 0.35s cubic-bezier(0.4, 0, 0.2, 1),
                margin-bottom 0.45s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .billing-panel-wrapper--open {
    grid-template-rows: 1fr;
    opacity: 1;
    margin-bottom: 24px;
  }

  .billing-panel-wrapper--closed {
    grid-template-rows: 0fr;
    opacity: 0;
    margin-bottom: 0;
    pointer-events: none;
  }

  .billing-panel-inner {
    overflow: hidden;
    min-height: 0;
  }

  .billing-panel-wrapper--open .billing-panel-inner {
    overflow: visible;
  }

  /* Responsive: stack button text on very small screens */
  @media (max-width: 400px) {
    .billing-toggle-btn {
      padding: 6px 12px;
      font-size: 12px;
    }
    .billing-toggle-btn__icon {
      width: 24px;
      height: 24px;
    }
  }
</style>
