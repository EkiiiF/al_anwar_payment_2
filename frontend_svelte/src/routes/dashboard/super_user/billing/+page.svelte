<script lang="ts">
  import { onMount } from 'svelte';
  import { PlayCircle, Receipt, Filter, Search, Calendar, CheckCircle2, AlertCircle, Clock, BookOpen, GraduationCap, Moon, Layers, ChevronLeft, ChevronRight, User, X, FileText, MapPin, Phone, Mail, Eye, EyeOff, ChevronDown, ChevronUp, Zap, Settings2, PanelTopClose, PanelTop } from 'lucide-svelte';
  import { superUserApi } from '$lib/api';
  import { formatRupiah, formatDate, getHijriMonthName, HIJRI_MONTH_NAMES } from '$lib/utils';
  import { Button, Alert, Spinner, Badge, EmptyState, Card, Select, Modal, BillingGenerator, AutoBillingSettings, StudentBillingModal, DataTable } from '$lib/components';
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

  let search = $state('');
  let searchTimeout: any;
  function handleSearch() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      page = 1;
      fetchData();
    }, 300);
  }

  let dynamicYearOptions = $state<{ value: string; label: string }[]>([]);
  let dynamicMonthOptions = $state<{ value: string; label: string }[]>([]);

  async function loadDynamicFilters() {
    try {
      const res = await superUserApi.getInvoices();
      const invoiceList = res.data ?? [];
      
      // Extract unique years from invoices in database
      const years = [...new Set(invoiceList.map(i => i.hijri_year))].filter(Boolean).sort((a, b) => b - a);
      dynamicYearOptions = years.map(y => ({ value: String(y), label: `${y} H` }));
      
      // Extract unique months from invoices in database
      const months = [...new Set(invoiceList.map(i => i.hijri_month))].filter(Boolean).sort((a, b) => a - b);
      dynamicMonthOptions = months.map(m => ({ value: String(m), label: getHijriMonthName(m) }));
    } catch (e) {
      console.error('Gagal memuat filter dinamis', e);
      // Fallback
      dynamicYearOptions = Array.from({ length: 5 }, (_, i) => {
        const y = 1445 + i;
        return { value: String(y), label: `${y} H` };
      });
      dynamicMonthOptions = HIJRI_MONTH_NAMES.map((name, i) => ({ value: String(i + 1), label: name }));
    }
  }

  async function fetchData() {
    loading = true;
    error   = '';
    try {
      const filters: any = {};
      if (statusFilter) filters.status = statusFilter;
      if (monthFilter)  filters.month  = monthFilter;
      if (yearFilter)   filters.year   = yearFilter;
      if (search)       filters.search = search;
      
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
        genHijriYear = String(hijriInfo.hijri_year);
        genYear = String(hijriInfo.hijri_year);
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
    loadDynamicFilters();
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
    limit;
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

  function handlePageChange(newPage: number) {
    page = newPage;
  }

  function handleLimitChange(newLimit: number) {
    limit = newLimit;
    page = 1;
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
      if (hijriInfo) {
        const selectedMonthNum = Number(genMonth);
        const selectedYearNum = Number(genYear);
        if (selectedMonthNum !== hijriInfo.hijri_month || selectedYearNum !== hijriInfo.hijri_year) {
          feedbackType = 'error';
          feedbackTitle = 'Bulan & Tahun Tidak Sesuai';
          feedbackMessage = `Anda tidak dapat membuat tagihan untuk bulan ${getHijriMonthName(selectedMonthNum)} ${selectedYearNum} H. Pembuatan tagihan manual hanya dapat dilakukan pada bulan aktif sistem saat ini (${hijriInfo.hijri_month_name} ${hijriInfo.hijri_year} H).`;
          showFeedbackModal = true;
          return;
        }
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
      if (hijriInfo) {
        const selectedSemNum = Number(genSemester);
        const selectedYearNum = Number(genHijriYear);
        if (selectedSemNum !== hijriInfo.semester || selectedYearNum !== hijriInfo.hijri_year) {
          feedbackType = 'error';
          feedbackTitle = 'Semester & Tahun Tidak Sesuai';
          feedbackMessage = `Anda tidak dapat membuat tagihan semester untuk Semester ${selectedSemNum} TA ${selectedYearNum} H. Pembuatan tagihan semester hanya dapat dilakukan sesuai dengan Semester dan Tahun aktif sistem saat ini (Semester ${hijriInfo.semester} TA ${hijriInfo.hijri_year} H).`;
          showFeedbackModal = true;
          return;
        }
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
        feedbackMessage = `Tagihan Syahriyyah Pondok berhasil dibuat untuk bulan ${getHijriMonthName(Number(genMonth))} ${genYear} H!`;
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
            <PanelTopClose size={15} class="text-slate-600" />
            <span>Sembunyikan Panel</span>
            <ChevronUp size={14} class="text-slate-400" />
          {:else}
            <PanelTop size={15} class="text-slate-600" />
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

  <div class="bg-white border border-slate-200/80 rounded-2xl shadow-sm p-5 hover:shadow-md transition-all duration-200">
    <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
      <!-- Search Input -->
      <div class="relative flex-1 max-w-lg w-full">
        <Search class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
        <input
          type="text"
          placeholder="Cari NIS atau nama santri..."
          class="w-full pl-11 pr-4 py-2.5 rounded-xl border border-slate-200 focus:ring-2 focus:ring-emerald-600 focus:border-emerald-600 transition-all text-sm bg-slate-50/50 hover:bg-slate-50"
          bind:value={search}
          oninput={handleSearch}
        />
      </div>

      <!-- Filters Selectors -->
      <div class="flex flex-wrap items-center gap-3 w-full lg:w-auto">
        <div class="flex items-center gap-1.5 text-xs font-bold text-slate-500 uppercase tracking-wider select-none">
          <Filter size={14} class="text-emerald-800" />
          <span>Filter</span>
        </div>

        <div class="w-full sm:w-40 flex-1 sm:flex-none">
          <Select
            id="f-status"
            bind:value={statusFilter}
            options={[
              { value: '', label: 'Semua Status' },
              { value: 'unpaid', label: 'Belum Lunas' },
              { value: 'paid', label: 'Lunas' }
            ]}
            class="w-full"
          />
        </div>

        <div class="w-full sm:w-40 flex-1 sm:flex-none">
          <Select
            id="f-month"
            bind:value={monthFilter}
            options={[{ value: '', label: 'Semua Bulan' }, ...dynamicMonthOptions]}
            class="w-full"
          />
        </div>

        <div class="w-full sm:w-40 flex-1 sm:flex-none">
          <Select
            id="f-year"
            bind:value={yearFilter}
            options={[{ value: '', label: 'Semua Tahun' }, ...dynamicYearOptions]}
            class="w-full"
          />
        </div>
      </div>
    </div>
  </div>

  {#if loading}
    <Spinner size="lg" />
  {:else}
    <DataTable
      pagination={pagination}
      onPageChange={handlePageChange}
      onLimitChange={handleLimitChange}
      isEmpty={students.length === 0}
      emptyTitle="Tidak ada santri dengan tagihan"
      emptyDescription="Coba ubah filter atau buat tagihan baru."
      paginationLabel="santri"
    >
      {#snippet children()}
        <table class="w-full text-sm text-left animate-fade-in" aria-label="Tabel data tagihan santri">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-100 whitespace-nowrap">
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">NIS</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">Nama Santri</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">Kategori Santri</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">Jumlah Tagihan</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">Status Pembayaran</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 whitespace-nowrap bg-white">
            {#each students as student (student.id)}
              {@const fullName = [student.name?.first_name, student.name?.middle_name, student.name?.last_name].filter(Boolean).join(' ')}
              <tr class="hover:bg-slate-50/80 transition-colors cursor-pointer" onclick={() => openStudentModal(student)}>
                <td class="px-5 py-4 font-mono text-green-600 font-semibold text-xs">{student.student_number}</td>
                <td class="px-5 py-4 font-semibold text-gray-900">{fullName}</td>
                <td class="px-5 py-4">
                  {#if student.status}
                    <Badge label={student.status.name} variant="info" />
                  {:else}
                    <span class="text-gray-400">-</span>
                  {/if}
                </td>
                <td class="px-5 py-4 font-semibold text-gray-700">
                  {student.invoices?.length ?? 0} Tagihan
                </td>
                <td class="px-5 py-4">
                  {#if student.invoices?.some(i => i.status === 'unpaid')}
                    <Badge label="Belum Lunas" variant="warning" dot />
                  {:else if student.invoices?.every(i => i.status === 'paid') && (student.invoices?.length ?? 0) > 0}
                    <Badge label="Lunas Semua" variant="success" dot />
                  {:else}
                    <Badge label="Tidak Ada Tagihan" variant="secondary" dot />
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/snippet}
    </DataTable>
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
          Apakah Anda yakin ingin menerbitkan tagihan manual untuk periode ini?
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
              Peringatan: Seluruh tagihan bulanan pada semester ini akan dibuat sekaligus.
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
        <Button variant="primary" onclick={executeGenerate} loading={generating} size="md">
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
          <p class="text-sm text-slate-600 leading-relaxed whitespace-pre-line">
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
          class="px-8"
        >
          {#snippet children()}
          {feedbackType === 'success' ? 'Selesai' : 'Mengerti'}
          {/snippet}
        </Button>
      </div>
    {/snippet}
  </Modal>
</div>
