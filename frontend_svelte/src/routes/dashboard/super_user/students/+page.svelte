<script lang="ts">
  import { onMount } from 'svelte';
  import { UserPlus, Edit, Trash2, Search, Power, ChevronDown, ChevronUp, MapPin, Phone, Mail, User } from 'lucide-svelte';
  import { superUserApi, pengasuhApi } from '$lib/api';
  import { auth } from '$lib/stores/auth';
  import { formatRupiah, formatDate } from '$lib/utils';
  import { Button, Alert, Spinner, Modal, Input, Select, Badge, Card, DataTable, ConfirmDialog } from '$lib/components';
  import type { Student, Category, StudentStatusType } from '$lib/types';
  import { getMuhadhorohLabel, MUHADHOROH_OPTIONS, GUARDIAN_RELATION_OPTIONS } from '$lib/types/student.types';
  import { toast } from '$lib/stores/toast';

  const isReadOnly = $derived($auth.user?.role?.name === 'pengasuh');
  const api = $derived(isReadOnly ? pengasuhApi : superUserApi);

  let students   = $state<Student[]>([]);
  let statuses   = $state<StudentStatusType[]>([]);
  let categories = $state<Category[]>([]);
  let loading    = $state(true);
  let error      = $state('');
  let pagination = $state<{ page: number; limit: number; total: number; pages: number } | null>(null);
  let page = $state(1);
  let limit = $state(10);

  let search        = $state('');
  let searchTimeout: any;

  function handleSearch() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(fetchData, 500);
  }

  let expandedId = $state<string | null>(null);
  function toggleExpand(id: string) {
    expandedId = expandedId === id ? null : id;
  }

  let showModal     = $state(false);
  let showCancelConfirm = $state(false);
  let deleteConfirmOpen = $state(false);
  let deleteConfirmTitle = $state('');
  let deleteConfirmMessage = $state('');
  let deleteConfirmAction = $state<(() => void | Promise<void>) | null>(null);
  let deleting = $state(false);
  let isEditing     = $state(false);
  let submitting    = $state(false);
  let modalError    = $state('');
  let currentId     = $state('');

  let nis             = $state('');
  let fullName        = $state('');
  let statusTypeId    = $state('');
  let muhadhorohLevel = $state('1');
  let currentSemester = $state('1');
  let nik             = $state('');
  let birthDay        = $state('');
  let birthMonth      = $state('');
  let birthYear       = $state('');
  let gender          = $state('L');
  let addressLine     = $state('');
  let village         = $state('');
  let district        = $state('');
  let city            = $state('');
  let province        = $state('');
  let country         = $state('Indonesia');
  let postalCode      = $state('');
  let guardianName    = $state('');
  let guardianPhone   = $state('');
  let guardianEmail   = $state('');
  let guardianPassword= $state('');
  let guardianRelation= $state('');
  let selectedCategoryIds = $state<string[]>([]);

  let nisError = $state('');
  let fullNameError = $state('');
  let nikError = $state('');
  let birthDayError = $state('');
  let addressLineError = $state('');
  let postalCodeError = $state('');
  let guardianNameError = $state('');
  let guardianPhoneError = $state('');
  let guardianEmailError = $state('');
  let guardianPasswordError = $state('');

  function clearErrors() {
    nisError = '';
    fullNameError = '';
    nikError = '';
    birthDayError = '';
    addressLineError = '';
    postalCodeError = '';
    guardianNameError = '';
    guardianPhoneError = '';
    guardianEmailError = '';
    guardianPasswordError = '';
  }

  const statusOptions = $derived(
    statuses.map(s => ({ value: s.id, label: `${s.name} (Diskon ${s.discount_percentage}%)` }))
  );
  const semesterOptions = [
    { value: '1', label: 'Semester 1' },
    { value: '2', label: 'Semester 2' },
  ];

  async function fetchData() {
    loading = true;
    error = '';
    try {
      if (isReadOnly) {
        const resStdPag = await api.getStudentsPaginated(search, page, limit);
        students   = resStdPag.data?.students ?? [];
        pagination = resStdPag.data?.pagination ?? null;
      } else {
        const [resStdPag, resStat, resCat] = await Promise.all([
          api.getStudentsPaginated(search, page, limit),
          superUserApi.getStatusTypes(),
          superUserApi.getCategories()
        ]);
        students   = resStdPag.data?.students ?? [];
        pagination = resStdPag.data?.pagination ?? null;
        statuses   = resStat.data  ?? [];
        categories = resCat.data   ?? [];
      }
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat data.';
    } finally {
      loading = false;
    }
  }

  function handlePageChange(newPage: number) {
    page = newPage;
    fetchData();
  }

  function handleLimitChange(newLimit: number) {
    limit = newLimit;
    page = 1;
    fetchData();
  }

  let hasInitialLoaded = false;
  $effect(() => {
    if (!$auth.loading && $auth.user && !hasInitialLoaded) {
      hasInitialLoaded = true;
      fetchData();
    }
  });

  function resetForm() {
    currentId = ''; isEditing = false; modalError = '';
    nis = ''; fullName = ''; statusTypeId = '';
    muhadhorohLevel = '1'; currentSemester = '1';
    nik = ''; birthDay = ''; birthMonth = ''; birthYear = ''; gender = 'L';
    addressLine = ''; village = ''; district = ''; city = ''; province = ''; country = 'Indonesia'; postalCode = '';
    guardianName = ''; guardianPhone = ''; guardianEmail = ''; guardianPassword = ''; guardianRelation = '';
    selectedCategoryIds = [];
    clearErrors();
  }

  function handleCancel() {
    showCancelConfirm = true;
  }

  function openAdd() {
    resetForm();
    selectedCategoryIds = categories.filter(c => c.is_active).map(c => c.id);
    showModal = true;
  }

  function openEdit(std: Student) {
    resetForm();
    isEditing = true;
    currentId = std.id;
    nis          = std.student_number;
    fullName     = [std.name.first_name, std.name.middle_name, std.name.last_name].filter(Boolean).join(' ');
    statusTypeId = std.status_id;
    muhadhorohLevel = String(std.muhadhoroh_level ?? 1);
    currentSemester = String(std.current_semester ?? 1);
    selectedCategoryIds = std.billing_categories?.map(c => c.id) ?? [];
    nik = std.nik ?? '';
    if (std.birth_date) {
      const parts = std.birth_date.split('T')[0].split('-');
      if (parts.length === 3) {
        birthYear = parts[0];
        birthMonth = String(parseInt(parts[1]));
        birthDay = String(parseInt(parts[2]));
      }
    } else {
      birthDay = ''; birthMonth = ''; birthYear = '';
    }
    gender = std.gender ?? 'L';
    const primaryAddr = std.addresses?.find(a => a.is_primary);
    addressLine = primaryAddr?.address_line ?? '';
    village     = primaryAddr?.village ?? '';
    district    = primaryAddr?.district ?? '';
    city        = primaryAddr?.city ?? '';
    province    = primaryAddr?.province ?? '';
    country     = primaryAddr?.country ?? 'Indonesia';
    postalCode  = primaryAddr?.postal_code ?? '';
    
    const primaryGuardian = std.guardians?.[0];
    if (primaryGuardian) {
      guardianName     = [primaryGuardian.name.first_name, primaryGuardian.name.middle_name, primaryGuardian.name.last_name].filter(Boolean).join(' ');
      guardianPhone    = primaryGuardian.phone;
      guardianEmail    = primaryGuardian.email;
      guardianRelation = primaryGuardian.relation || 'Orang Tua';
    }
    showModal = true;
  }

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    clearErrors();
    modalError = '';

    let hasError = false;

    if (!nis.trim()) {
      nisError = 'Nomor Induk (NIS) wajib diisi.';
      hasError = true;
    } else if (!/^\d+$/.test(nis)) {
      nisError = 'NIS harus berupa angka saja.';
      hasError = true;
    } else if (nis.length < 5 || nis.length > 15) {
      nisError = 'NIS harus berjumlah 5 hingga 15 digit.';
      hasError = true;
    }

    if (!fullName.trim()) {
      fullNameError = 'Nama lengkap wajib diisi.';
      hasError = true;
    } else if (fullName.trim().length < 3) {
      fullNameError = 'Nama lengkap minimal 3 karakter.';
      hasError = true;
    }

    if (nik.trim() && !/^\d{16}$/.test(nik.trim())) {
      nikError = 'NIK harus berupa 16 digit angka.';
      hasError = true;
    }

    if (birthDay && birthMonth && birthYear) {
      const bDate = new Date(Number(birthYear), Number(birthMonth) - 1, Number(birthDay));
      const today = new Date();
      if (bDate > today) {
        birthDayError = 'Tanggal lahir tidak boleh lebih besar dari tanggal hari ini.';
        hasError = true;
      }
    }

    if (postalCode.trim() && !/^\d{5}$/.test(postalCode.trim())) {
      postalCodeError = 'Kode pos harus berupa 5 digit angka.';
      hasError = true;
    }

    if (!guardianName.trim()) {
      guardianNameError = 'Nama wali santri wajib diisi.';
      hasError = true;
    } else if (guardianName.trim().length < 3) {
      guardianNameError = 'Nama wali santri minimal 3 karakter.';
      hasError = true;
    }

    if (!guardianPhone.trim()) {
      guardianPhoneError = 'Nomor HP/WhatsApp wali wajib diisi.';
      hasError = true;
    } else if (!/^\+?[0-9]{11,15}$/.test(guardianPhone.trim().replace(/[\s-]/g, ''))) {
      guardianPhoneError = 'Format nomor HP tidak valid.';
      hasError = true;
    }

    if (!guardianEmail.trim()) {
      guardianEmailError = 'Email wali wajib diisi.';
      hasError = true;
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(guardianEmail.trim())) {
      guardianEmailError = 'Format email tidak valid.';
      hasError = true;
    }

    if (guardianPassword && guardianPassword.length < 8) {
      guardianPasswordError = 'Password minimal 8 karakter.';
      hasError = true;
    }

    if (hasError) {
      modalError = 'Mohon periksa kembali input form Anda.';
      return;
    }

    submitting = true;
    try {
      const birthDate = (birthYear && birthMonth && birthDay)
        ? `${birthYear}-${String(birthMonth).padStart(2, '0')}-${String(birthDay).padStart(2, '0')}`
        : '';
      const payload = {
        nis, full_name: fullName,
        status_type_id: statusTypeId,
        muhadhoroh_level: Number(muhadhorohLevel),
        nik, birth_date: birthDate, gender,
        address_line: addressLine, village, district, city, province, country, postal_code: postalCode,
        guardian_name: guardianName, guardian_phone: guardianPhone, guardian_email: guardianEmail,
        guardian_relation: guardianRelation,
        guardian_username: nis,
        create_new_guardian: true,
        guardian_password: guardianPassword || 'password123',
        category_ids: selectedCategoryIds
      };

      if (isEditing) {
        await superUserApi.updateStudent(currentId, payload);
      } else {
        await superUserApi.createStudent(payload);
      }

      showModal = false;
      await fetchData();
    } catch (e: unknown) {
      modalError = e instanceof Error ? e.message : 'Gagal menyimpan data.';
    } finally {
      submitting = false;
    }
  }

  async function toggleStudentStatus(id: string) {
    try {
      await superUserApi.toggleStudentStatus(id);
      toast.success('Status keaktifan santri berhasil diubah');
      await fetchData();
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Gagal mengubah status.');
    }
  }

  function handleDelete(id: string, name: string) {
    deleteConfirmTitle = 'Hapus Data Santri';
    deleteConfirmMessage = `Apakah Anda yakin ingin menghapus data santri "${name}"? Tindakan ini tidak dapat dibatalkan, namun seluruh riwayat pembayaran dan tagihan santri ini tetap aman sebagai histori.`;
    deleteConfirmAction = async () => {
      deleting = true;
      try {
        await superUserApi.deleteStudent(id);
        toast.success('Data santri berhasil dihapus');
        await fetchData();
      } catch (e: unknown) {
        toast.error(e instanceof Error ? e.message : 'Gagal menghapus santri.');
      } finally {
        deleting = false;
        deleteConfirmOpen = false;
      }
    };
    deleteConfirmOpen = true;
  }
</script>

<svelte:head>
  <title>Data Santri | {isReadOnly ? 'Dashboard Pengasuh' : 'Dashboard Super User'}</title>
  <meta name="description" content="Manajemen data santri lengkap dengan kelas Muhadhoroh, semester, dan guardian." />
</svelte:head>

<div class="space-y-6 flex flex-col flex-1 min-h-0">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      {#if isReadOnly}
        <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-purple-100 border border-purple-200 text-purple-700 text-xs font-semibold uppercase tracking-wider mb-2">
          <span class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse" aria-hidden="true"></span>
          Mode Hanya Lihat
        </div>
      {/if}
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Data Santri</h1>
      <p class="text-gray-500 text-sm mt-1">Kelola informasi data santri.</p>
    </div>
    <div class="flex items-center gap-3 w-full sm:w-auto">
      <div class="relative flex-1 sm:w-72 md:w-80">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
        <input
          type="text"
          placeholder="Cari santri berdasarkan nama atau NIS..."
          class="w-full pl-10 pr-4 py-2 rounded-xl border border-gray-200 focus:ring-2 {isReadOnly ? 'focus:ring-purple-500 focus:border-purple-500' : 'focus:ring-green-500 focus:border-green-500'} transition-all text-sm bg-white"
          bind:value={search}
          oninput={handleSearch}
        />
      </div>
      {#if !isReadOnly}
        <Button onclick={openAdd} variant="primary" size="md" class="shrink-0">
          {#snippet children()}
            <UserPlus size={16} aria-hidden="true" />
            <span class="whitespace-nowrap">Tambah</span>
          {/snippet}
        </Button>
      {/if}
    </div>
  </div>

  {#if error}
    <Alert type="error" message={error} />
  {/if}

  {#if loading}
    <Spinner size="lg" />
  {:else}
    <DataTable
      pagination={pagination}
      onPageChange={handlePageChange}
      onLimitChange={handleLimitChange}
      isEmpty={students.length === 0}
      emptyTitle="Belum ada data santri"
      emptyDescription="Klik tombol 'Tambah' untuk mendaftarkan santri baru."
      paginationLabel="santri"
    >
      {#snippet children()}
        <table class="w-full text-sm text-left" aria-label="Tabel data santri">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-100 whitespace-nowrap">
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">NIS</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">Santri</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">Kelas & Semester</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider hidden md:table-cell">Wali Santri</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider hidden lg:table-cell">Status</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 whitespace-nowrap">
            {#each students as student (student.id)}
              <tr class="hover:bg-gray-50 transition-colors cursor-pointer" onclick={() => toggleExpand(student.id)}>
                <td class="px-5 py-4 font-mono {isReadOnly ? 'text-purple-600' : 'text-green-600'} font-semibold text-xs">{student.student_number}</td>
                <td class="px-5 py-4">
                  <p class="font-semibold text-gray-900">
                    {[student.name.first_name, student.name.middle_name, student.name.last_name].filter(Boolean).join(' ')}
                  </p>
                  {#if student.status}
                    <Badge label={student.status.name} variant="purple" />
                  {/if}
                </td>
                <td class="px-5 py-4">
                  {#if student.muhadhoroh_level === 0}
                    <Badge label="Lulus" variant="success" />
                  {:else}
                    <div class="inline-block align-middle">
                      <p class="text-sm font-bold text-gray-900">Muhadhoroh {student.muhadhoroh_level}</p>
                      <p class="text-xs {isReadOnly ? 'text-purple-600' : 'text-emerald-600'} font-semibold">Semester {student.current_semester}</p>
                    </div>
                  {/if}
                </td>
                <td class="px-5 py-4 hidden md:table-cell">
                  {#if student.guardians && student.guardians.length > 0}
                    <p class="text-gray-900 font-medium">
                      {[student.guardians[0].name.first_name, student.guardians[0].name.middle_name, student.guardians[0].name.last_name].filter(Boolean).join(' ')}
                    </p>
                    <p class="text-xs text-gray-500 font-mono">{student.guardians[0].user?.username ?? '-'}</p>
                  {:else}
                    <p class="text-gray-400 italic">Tidak ada data</p>
                  {/if}
                </td>
                <td class="px-5 py-4 hidden lg:table-cell">
                  <Badge
                    label={student.is_active ? 'Aktif' : 'Nonaktif'}
                    variant={student.is_active ? 'success' : 'danger'}
                    dot
                  />
                </td>
                <td class="px-5 py-4 text-right">
                  <div class="inline-flex gap-2">
                    <button
                      onclick={(e) => { e.stopPropagation(); toggleExpand(student.id); }}
                      class="p-2 rounded-lg bg-gray-50 text-gray-500 hover:bg-gray-100 border border-gray-200 transition-colors"
                      title="Lihat Detail"
                    >
                      {#if expandedId === student.id}
                        <ChevronUp size={15} />
                      {:else}
                        <ChevronDown size={15} />
                      {/if}
                    </button>
                    {#if !isReadOnly}
                      <button
                        onclick={(e) => { e.stopPropagation(); toggleStudentStatus(student.id); }}
                        class="p-2 rounded-lg {student.is_active ? 'bg-red-50 text-red-600 hover:bg-red-100 border-red-200' : 'bg-green-50 text-green-700 hover:bg-green-100 border-green-200'} border transition-colors"
                        title={student.is_active ? 'Nonaktifkan Santri' : 'Aktifkan Santri'}
                      >
                        <Power size={15} aria-hidden="true" />
                      </button>
                      <button
                        onclick={(e) => { e.stopPropagation(); openEdit(student); }}
                        class="p-2 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 border border-blue-200 transition-colors"
                        aria-label="Edit santri"
                      >
                        <Edit size={15} aria-hidden="true" />
                      </button>
                      <button
                        onclick={(e) => { e.stopPropagation(); handleDelete(student.id, [student.name.first_name, student.name.last_name].join(' ')); }}
                        class="p-2 rounded-lg bg-red-50 text-red-600 hover:bg-red-100 border border-red-200 transition-colors"
                        aria-label="Hapus santri"
                      >
                        <Trash2 size={15} aria-hidden="true" />
                      </button>
                    {/if}
                  </div>
                </td>
              </tr>

              {#if expandedId === student.id}
                <tr>
                  <td colspan="6" class="bg-gray-50/80 px-5 py-5">
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                      <div class="space-y-2">
                        <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                          <User size={12} /> Data Pribadi
                        </h4>
                        <div class="bg-white rounded-lg border border-gray-200 p-3 space-y-1.5 text-xs">
                          <div class="flex justify-between"><span class="text-gray-500">NIK</span><span class="font-semibold text-gray-900">{student.nik || '-'}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Tgl Lahir</span><span class="font-semibold text-gray-900">{student.birth_date ? formatDate(student.birth_date) : '-'}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Jenis Kelamin</span><span class="font-semibold text-gray-900">{student.gender === 'L' ? 'Laki-laki' : 'Perempuan'}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Kelas</span><span class="font-bold {isReadOnly ? 'text-purple-700' : 'text-emerald-700'}">{getMuhadhorohLabel(student.muhadhoroh_level, student.current_semester)}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Status</span><span class="font-semibold">{student.status?.name ?? '-'}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Diskon</span><span class="font-semibold text-purple-700">{student.status?.discount_percentage ?? 0}%</span></div>
                          
                          <div class="flex flex-col gap-1 border-t pt-2 mt-2">
                            <span class="text-gray-500 font-semibold block mb-0.5">Tagihan Aktif:</span>
                            <div class="flex flex-wrap gap-1">
                              {#if student.billing_categories && student.billing_categories.length > 0}
                                {#each student.billing_categories as bc}
                                  <span class="{isReadOnly ? 'bg-purple-50 text-purple-700 border-purple-200' : 'bg-green-50 text-green-700 border-green-200'} border px-2 py-0.5 rounded-full font-bold text-[10px]">
                                    {bc.name}
                                  </span>
                                {/each}
                              {:else}
                                <span class="bg-slate-50 text-slate-500 border border-slate-200 px-2 py-0.5 rounded-full font-semibold text-[10px]">
                                  Semua Tagihan (Default)
                                </span>
                              {/if}
                            </div>
                          </div>
                        </div>
                      </div>

                      <div class="space-y-2">
                        <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                          <Phone size={12} /> Data Wali Santri
                        </h4>
                        {#if student.guardians && student.guardians.length > 0}
                          {#each student.guardians as guardian}
                            <div class="bg-white rounded-lg border border-gray-200 p-3 space-y-1.5 text-xs">
                              <p class="font-bold text-gray-900">{[guardian.name.first_name, guardian.name.middle_name, guardian.name.last_name].filter(Boolean).join(' ')}</p>
                              <div class="flex items-center gap-1.5 text-gray-500"><Phone size={10} /> {guardian.phone || '-'}</div>
                              <div class="flex items-center gap-1.5 text-gray-500"><Mail size={10} /> {guardian.email || '-'}</div>
                              <div class="flex justify-between"><span class="text-gray-500">Username</span><span class="font-mono font-semibold">{guardian.user?.username ?? '-'}</span></div>
                              <div class="flex justify-between"><span class="text-gray-500">Relasi</span><span class="font-semibold">{guardian.relation}</span></div>
                            </div>
                          {/each}
                        {:else}
                          <p class="text-xs text-gray-400 italic">Belum ada data wali santri</p>
                        {/if}
                      </div>

                      <div class="space-y-2">
                        <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                          <MapPin size={12} /> Alamat
                        </h4>
                        {#if student.addresses && student.addresses.length > 0}
                          {#each student.addresses as addr}
                            <div class="bg-white rounded-lg border border-gray-200 p-3 space-y-1.5 text-xs">
                              <p class="font-semibold text-gray-900">{addr.address_line || '-'}</p>
                              {#if addr.village || addr.district}
                                <div class="flex justify-between"><span class="text-gray-500">Desa / Kecamatan</span><span class="font-semibold">{[addr.village, addr.district].filter(Boolean).join(', ')}</span></div>
                              {/if}
                              {#if addr.city || addr.province}
                                <div class="flex justify-between"><span class="text-gray-500">Kota / Provinsi</span><span class="font-semibold">{[addr.city, addr.province].filter(Boolean).join(', ')}</span></div>
                              {/if}
                              {#if addr.postal_code}
                                <div class="flex justify-between"><span class="text-gray-500">Kode Pos</span><span class="font-semibold">{addr.postal_code}</span></div>
                              {/if}
                              {#if addr.country && addr.country !== 'Indonesia'}
                                <div class="flex justify-between"><span class="text-gray-500">Negara</span><span class="font-semibold">{addr.country}</span></div>
                              {/if}
                              {#if addr.is_primary}
                                <Badge label="Utama" variant="info" />
                              {/if}
                            </div>
                          {/each}
                        {:else}
                          <p class="text-xs text-gray-400 italic">Belum ada data alamat</p>
                        {/if}
                      </div>
                    </div>
                  </td>
                </tr>
              {/if}
            {/each}
          </tbody>
        </table>
      {/snippet}
    </DataTable>
  {/if}
</div>

<Modal bind:open={showModal} title={isEditing ? 'Edit Data Santri' : 'Tambah Santri Baru'} size="xl" confirmCloseMessage="Apakah Anda yakin ingin membatalkan tindakan ini? Perubahan yang belum disimpan akan hilang.">
  {#snippet children()}
    {#if modalError}
      <Alert type="error" message={modalError} class="mb-5" />
    {/if}

    <form id="student-form" onsubmit={handleSubmit} class="space-y-6" novalidate autocomplete="off">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <fieldset class="space-y-4">
          <legend class="text-xs font-semibold text-green-700 uppercase tracking-widest pb-2 border-b border-gray-200 w-full">
            Informasi Santri
          </legend>

          <Input id="nis"      label="Nomor Induk (NIS)" bind:value={nis}      required error={nisError} oninput={() => nisError = ''} />
          <Input id="fullName" label="Nama Lengkap"       bind:value={fullName} required error={fullNameError} oninput={() => fullNameError = ''} />
          
          <div class="grid grid-cols-2 gap-4">
            <Input id="nik" label="NIK" bind:value={nik} error={nikError} oninput={() => nikError = ''} />
            <div class="flex flex-col gap-1.5">
              <label for="gender" class="text-xs font-semibold text-gray-600 uppercase tracking-wider ml-0.5">Jenis Kelamin</label>
              <select id="gender" bind:value={gender} class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-gray-900 text-sm focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none">
                <option value="L">Laki-laki</option>
                <option value="P">Perempuan</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-3 gap-3">
            <div class="flex flex-col gap-1.5">
              <label for="birthDay" class="text-xs font-medium text-slate-600 uppercase tracking-wider">Tanggal</label>
              <input 
                id="birthDay" 
                type="number" 
                bind:value={birthDay} 
                min="1" 
                max="31" 
                class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-gray-900 text-sm focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none" 
                placeholder="1-31" 
                oninput={() => birthDayError = ''} 
              />
            </div>
            <div class="flex flex-col gap-1.5">
              <label for="birthMonth" class="text-xs font-medium text-slate-600 uppercase tracking-wider">Bulan</label>
              <select id="birthMonth" bind:value={birthMonth} class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-gray-900 text-sm focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none">
                <option value="">--</option>
                {#each ['Januari','Februari','Maret','April','Mei','Juni','Juli','Agustus','September','Oktober','November','Desember'] as name, i}
                  <option value={String(i + 1)}>{name}</option>
                {/each}
              </select>
            </div>
            <div class="flex flex-col gap-1.5">
              <label for="birthYear" class="text-xs font-medium text-slate-600 uppercase tracking-wider">Tahun Lahir</label>
              <input 
                id="birthYear" 
                type="number" 
                bind:value={birthYear} 
                min="1900" 
                max={new Date().getFullYear()} 
                class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-gray-900 text-sm focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none" 
                placeholder="Tahun" 
              />
            </div>
            {#if birthDayError}
              <p class="text-xs text-red-500 col-span-3">{birthDayError}</p>
            {/if}
          </div>

          <div class="bg-gray-50 rounded-lg border border-gray-200 p-3 space-y-3">
            <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Alamat Lengkap</p>
            <Input id="addressLine" label="Alamat (Jalan/Gang/RT/RW)" bind:value={addressLine} placeholder="Jl. Raya No.1 RT 01/RW 02" error={addressLineError} oninput={() => addressLineError = ''} />
            <div class="grid grid-cols-2 gap-3">
              <Input id="village"  label="Desa / Kelurahan" bind:value={village}  placeholder="Desa" />
              <Input id="district" label="Kecamatan"        bind:value={district} placeholder="Kecamatan" />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <Input id="city"     label="Kota / Kabupaten" bind:value={city}     placeholder="Kota" />
              <Input id="province" label="Provinsi"         bind:value={province} placeholder="Provinsi" />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <Input id="postalCode" label="Kode Pos" bind:value={postalCode} placeholder="00000" error={postalCodeError} oninput={() => postalCodeError = ''} />
              <Input id="country"    label="Negara"   bind:value={country}    placeholder="Indonesia" />
            </div>
          </div>

          <Select
            id="muhadhorohLevel"
            label="Kelas Muhadhoroh"
            bind:value={muhadhorohLevel}
            options={MUHADHOROH_OPTIONS}
            required
          />

          <Select
            id="statusTypeId"
            label="Status Santri"
            bind:value={statusTypeId}
            options={statusOptions}
            placeholder="Pilih status..."
            required
          />

          <div class="space-y-2">
            <span class="text-xs font-semibold text-slate-600 uppercase tracking-wider block ml-0.5">Kategori Tagihan yang Berlaku</span>
            <div class="grid grid-cols-2 gap-3 bg-slate-50 p-4 rounded-xl border border-slate-200">
              {#each categories as category}
                <label class="flex items-center gap-2.5 text-sm font-semibold text-slate-705 cursor-pointer hover:text-green-800 hover:bg-white p-2 rounded-lg border border-transparent hover:border-slate-200 hover:shadow-sm transition-all select-none">
                  <input
                    type="checkbox"
                    value={category.id}
                    checked={selectedCategoryIds.includes(category.id)}
                    onchange={(e) => {
                      const checked = e.currentTarget.checked;
                      if (checked) {
                        if (!selectedCategoryIds.includes(category.id)) {
                          selectedCategoryIds.push(category.id);
                        }
                      } else {
                        selectedCategoryIds = selectedCategoryIds.filter(id => id !== category.id);
                      }
                    }}
                    class="rounded border-slate-300 text-green-700 focus:ring-green-500/20 focus:ring-2 w-4 h-4 cursor-pointer"
                  />
                  <span class="leading-none">{category.name}</span>
                </label>
              {/each}
            </div>
          </div>
        </fieldset>

        <fieldset class="space-y-4">
          <legend class="text-xs font-bold text-blue-700 uppercase tracking-widest pb-2 border-b border-gray-200 w-full">
            Data & Akun Wali Santri
          </legend>

          <Input id="guardianName"  label="Nama Wali Santri"  bind:value={guardianName}  required error={guardianNameError} oninput={() => guardianNameError = ''} />
          <Input id="guardianPhone" label="Nomor HP / WhatsApp"     bind:value={guardianPhone} required type="tel" error={guardianPhoneError} oninput={() => guardianPhoneError = ''} />
          <Input id="guardianEmail" label="Email"      bind:value={guardianEmail} required type="email" error={guardianEmailError} oninput={() => guardianEmailError = ''}
                 autocomplete="email" />
          <Select
            id="guardianRelation"
            label="Hubungan dengan Santri"
            bind:value={guardianRelation}
            options={GUARDIAN_RELATION_OPTIONS}
            required
          />

          <div class="bg-blue-50 border border-blue-200 rounded-xl p-4 space-y-3">
            <p class="text-xs text-blue-700 font-semibold">Kredensial Login Wali Santri </p>
            <Input id="guardianUsername" label="Username Login" value={nis} disabled
                   helper="Otomatis mengikuti Nomor Induk (NIS)" />
            <Input id="guardianPassword" label="Password Login"  bind:value={guardianPassword} type="password" error={guardianPasswordError} oninput={() => guardianPasswordError = ''}
                   autocomplete="new-password" helper={isEditing ? 'Kosongkan jika tidak ingin mengubah' : 'Default: password123'} />
          </div>
        </fieldset>
      </div>
    </form>
  {/snippet}

  {#snippet footer()}
    <div class="flex justify-end gap-3">
      <Button onclick={handleCancel} variant="outline" size="md">
        {#snippet children()}Batal{/snippet}
      </Button>
      <Button type="submit" form="student-form" variant="primary" size="md" loading={submitting}>
        {#snippet children()}
          {isEditing ? 'Simpan Perubahan' : 'Daftarkan Santri'}
        {/snippet}
      </Button>
    </div>
  {/snippet}
</Modal>

<ConfirmDialog
  bind:open={showCancelConfirm}
  title="Konfirmasi Tindakan"
  message="Apakah Anda yakin ingin membatalkan tindakan ini? Perubahan yang belum disimpan akan hilang."
  confirmText="Ya, Batalkan"
  cancelText="Kembali"
  variant="warning"
  onConfirm={() => {
    showCancelConfirm = false;
    showModal = false;
  }}
/>

<ConfirmDialog
  bind:open={deleteConfirmOpen}
  title={deleteConfirmTitle}
  message={deleteConfirmMessage}
  confirmText="Ya, Hapus"
  cancelText="Batal"
  variant="danger"
  loading={deleting}
  onConfirm={async () => {
    if (deleteConfirmAction) {
      await deleteConfirmAction();
    }
  }}
/>
