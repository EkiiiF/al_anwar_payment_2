<script lang="ts">
  import { onMount } from 'svelte';
  import { PlayCircle, Receipt, Filter, Search, Calendar, CheckCircle2, AlertCircle, Clock, BookOpen, GraduationCap, Moon, Layers, ChevronLeft, ChevronRight, User, X, FileText, MapPin, Phone, Mail, Eye, EyeOff, ChevronDown, ChevronUp, Zap, Settings2, PanelTopClose, PanelTop } from 'lucide-svelte';
  import { superUserApi } from '$lib/api';
  import { formatRupiah, formatDate, getHijriMonthName, HIJRI_MONTH_NAMES } from '$lib/utils';
  import { Button, Alert, Spinner, Badge, EmptyState, Card, Select, Modal, BillingGenerator, AutoBillingSettings, BillingFilters, StudentBillingCard, StudentBillingModal } from '$lib/components';
  import { toast } from '$lib/stores/toast';
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

  // Custom Modal Confirmation & Notification States
  let showConfirmModal           = $state(false);
  let confirmMode                = $state<'monthly' | 'semester'>('monthly');
  let showFeedbackModal          = $state(false);
  let feedbackType               = $state<'success' | 'error'>('success');
  let feedbackTitle              = $state('');
  let feedbackMessage            = $state('');

  const monthOptions = HIJRI_MONTH_NAMES.map((name, i) => ({ value: String(i + 1), label: name }));
  const yearOptions  = Array.from({ length: 5 }, (_, i) => {
    const y = 1445 + i;
    return { value: String(y), label: `${y} H` };
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
    const stored = localStorage.getItem('showBillingPanel');
    if (stored !== null) {
      showBillingPanel = stored === 'true';
    }
    fetchData();
    fetchSettings();
    fetchHijriInfo();
  });

  function toggleBillingPanel() {
    showBillingPanel = !showBillingPanel;
    localStorage.setItem('showBillingPanel', String(showBillingPanel));
  }

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

  function triggerGenerateConfirmation(mode: 'monthly' | 'semester') {
    if (mode === 'monthly') {
      const yr = Number(genYear);
      if (!genYear || isNaN(yr) || yr < 1300 || yr > 1600) {
        feedbackType = 'error';
        feedbackTitle = 'Input Tidak Valid';
        feedbackMessage = 'Tahun Hijriah wajib diisi dengan angka antara 1300 - 1600 H (Contoh: 1447).';
        showFeedbackModal = true;
        return;
      }
    } else {
      const yr = Number(genHijriYear);
      if (!genHijriYear || isNaN(yr) || yr < 1300 || yr > 1600) {
        feedbackType = 'error';
        feedbackTitle = 'Input Tidak Valid';
        feedbackMessage = 'Tahun Hijriah wajib diisi dengan angka antara 1300 - 1600 H (Contoh: 1447).';
        showFeedbackModal = true;
        return;
      }
    }
    confirmMode = mode;
    showConfirmModal = true;
  }

  function handleGenerate() {
    triggerGenerateConfirmation(genMode);
  }

  async function executeGenerate() {
    showConfirmModal = false;
    generating = true;
    genMessage  = '';
    genSuccess = false;
    
    if (confirmMode === 'monthly') {
      const yr = Number(genYear);
      try {
        await superUserApi.generateInvoices(Number(genMonth), yr);
        feedbackType = 'success';
        feedbackTitle = 'Pembuatan Tagihan Berhasil';
        feedbackMessage = `Tagihan Syahriyyah Pondok & Muhadhoroh berhasil dibuat untuk bulan ${getHijriMonthName(Number(genMonth))} ${genYear} H!`;
        genMessage = feedbackMessage;
        genSuccess = true;
        showFeedbackModal = true;
        await fetchData();
      } catch (e: unknown) {
        feedbackType = 'error';
        feedbackTitle = 'Pembuatan Tagihan Gagal';
        feedbackMessage = e instanceof Error ? e.message : 'Gagal membuat tagihan.';
        genMessage = feedbackMessage;
        genSuccess = false;
        showFeedbackModal = true;
      } finally {
        generating = false;
      }
    } else {
      const yr = Number(genHijriYear);
      try {
        await superUserApi.generateSemesterInvoices(Number(genSemester), yr);
        feedbackType = 'success';
        feedbackTitle = 'Pembuatan Tagihan Berhasil';
        feedbackMessage = `Tagihan Semester ${genSemester} berhasil dibuat untuk tahun ${yr} H! Semua tagihan bulan di semester ini telah diterbitkan.`;
        genMessage = feedbackMessage;
        genSuccess = true;
        showFeedbackModal = true;
        await fetchData();
      } catch (e: unknown) {
        feedbackType = 'error';
        feedbackTitle = 'Pembuatan Tagihan Gagal';
        feedbackMessage = e instanceof Error ? e.message : 'Gagal membuat tagihan semester.';
        genMessage = feedbackMessage;
        genSuccess = false;
        showFeedbackModal = true;
      } finally {
        generating = false;
      }
    }
  }

  async function toggleSetting(key: string) {
    settingLoading = true;
    try {
      if (key === 'auto_generate_billing') {
        const nextValue = !autoBillingEnabled;
        await superUserApi.updateSetting(key, String(nextValue));
        autoBillingEnabled = nextValue;
        toast.success(`Auto-billing bulanan berhasil ${nextValue ? 'diaktifkan' : 'dinonaktifkan'}.`);
      } else {
        const nextValue = !autoSemesterBillingEnabled;
        await superUserApi.updateSetting(key, String(nextValue));
        autoSemesterBillingEnabled = nextValue;
        toast.success(`Auto-billing semester berhasil ${nextValue ? 'diaktifkan' : 'dinonaktifkan'}.`);
      }
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : 'Gagal mengubah pengaturan.';
      toast.error(msg);
    } finally {
      settingLoading = false;
    }
  }

  async function updateHijriSetting(key: string, value: string) {
    settingLoading = true;
    try {
      await superUserApi.updateSetting(key, value);
      let label = key;
      if (key === 'billing_hijri_day') label = 'Hari jatuh tempo Hijriah';
      else if (key === 'billing_hijri_month') label = 'Bulan jatuh tempo Hijriah';
      else if (key === 'billing_hijri_year') label = 'Tahun jatuh tempo Hijriah';
      toast.success(`Pengaturan ${label} berhasil diperbarui.`);
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : 'Gagal mengubah pengaturan.';
      toast.error(msg);
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
      <p class="text-gray-500 text-sm mt-1">Buat tagihan berdasarkan kalender Hijriah dan pantau status pembayaran per semester.</p>
    </div>
    <div class="flex items-center gap-3">
      <Button
        variant="outline"
        onclick={toggleBillingPanel}
        class="flex items-center gap-2 text-xs font-semibold px-3 py-1.5 border border-slate-200 shadow-sm"
      >
        {#snippet children()}
          {#if showBillingPanel}
            <PanelTopClose size={15} class="text-slate-650" />
            <span>Sembunyikan Panel</span>
            <ChevronUp size={14} class="text-slate-400" />
          {:else}
            <PanelTop size={15} class="text-slate-650" />
            <span>Tampilkan Panel</span>
            <ChevronDown size={14} class="text-slate-400" />
          {/if}
        {/snippet}
      </Button>
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
          <p class="text-xs text-emerald-600 mt-0.5 truncate">Tagihan dibuat</p>
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
        <EmptyState title="Tidak ada santri dengan tagihan" description="Coba ubah filter atau buat tagihan baru." />
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

  <!-- Modal Konfirmasi Pembuatan Tagihan -->
  <Modal bind:open={showConfirmModal} title="Konfirmasi Pembuatan Tagihan" size="md">
    {#snippet children()}
      <div class="space-y-4 py-2">
        <p class="text-sm text-slate-600">
          Apakah Anda yakin ingin menerbitkan tagihan manual untuk periode berikut? Tindakan ini akan membuat tagihan massal bagi seluruh santri reguler yang aktif.
        </p>
        <div class="bg-slate-50 p-4 rounded-xl border border-slate-200/60 space-y-2">
          <div class="flex justify-between text-sm">
            <span class="text-slate-500">Jenis Tagihan:</span>
            <span class="font-bold text-slate-800">
              {confirmMode === 'monthly' ? 'Syahriyyah Bulanan' : 'Tagihan Semester'}
            </span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-slate-500">Periode Akademik:</span>
            <span class="font-bold text-slate-800">
              {confirmMode === 'monthly' 
                ? `${getHijriMonthName(Number(genMonth))} ${genYear} H` 
                : `Semester ${genSemester} (${genSemester === '1' ? 'Ganjil' : 'Genap'}) - Tahun ${genHijriYear} H`}
            </span>
          </div>
          {#if confirmMode === 'semester'}
            <p class="text-xs text-amber-600 font-semibold bg-amber-50 p-2.5 rounded-lg border border-amber-100 mt-2">
              ⚠️ Peringatan: Seluruh tagihan bulanan pada semester ini akan dibuat sekaligus. Pastikan periode semester ini sesuai dengan kalender Hijriah yang sedang berjalan.
            </p>
          {/if}
        </div>
      </div>
    {/snippet}
    {#snippet footer()}
      <div class="flex justify-end gap-3 w-full">
        <Button variant="outline" onclick={() => showConfirmModal = false} size="md">
          {#snippet children()}Batal{/snippet}
        </Button>
        <Button variant="primary" onclick={executeGenerate} loading={generating} size="md" class="!bg-emerald-800 hover:!bg-emerald-750 text-white font-bold">
          {#snippet children()}Ya, Buat Tagihan{/snippet}
        </Button>
      </div>
    {/snippet}
  </Modal>

  <!-- Modal Feedback Hasil Aksi -->
  <Modal bind:open={showFeedbackModal} title={feedbackTitle} size="md">
    {#snippet children()}
      <div class="flex flex-col items-center text-center space-y-4 py-4">
        {#if feedbackType === 'success'}
          <div class="w-16 h-16 bg-emerald-50 rounded-full flex items-center justify-center border border-emerald-100 animate-pulse animate-duration-1000">
            <CheckCircle2 size={36} class="text-emerald-600" />
          </div>
        {:else}
          <div class="w-16 h-16 bg-red-50 rounded-full flex items-center justify-center border border-red-100 animate-pulse animate-duration-1000">
            <AlertCircle size={36} class="text-red-600" />
          </div>
        {/if}
        <div class="space-y-2 max-w-sm">
          <h3 class="font-bold text-slate-900 text-lg">{feedbackTitle}</h3>
          <p class="text-sm text-slate-650 leading-relaxed whitespace-pre-line">
            {feedbackMessage}
          </p>
        </div>
      </div>
    {/snippet}
    {#snippet footer()}
      <div class="flex justify-center w-full">
        <Button 
          variant={feedbackType === 'success' ? 'primary' : 'danger'} 
          onclick={() => showFeedbackModal = false} 
          size="md"
          class={feedbackType === 'success' ? '!bg-emerald-800 hover:!bg-emerald-700 px-8 text-white font-bold' : 'px-8 font-bold'}
        >
          {#snippet children()}
            {feedbackType === 'success' ? 'Selesai' : 'Mengerti'}
          {/snippet}
        </Button>
      </div>
    {/snippet}
  </Modal>
</div>

<style>
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
</style>
