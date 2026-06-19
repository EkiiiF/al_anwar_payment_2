<script lang="ts">
  import { onMount } from 'svelte';
  import { Plus, Edit, Trash2, Tag, Users2, UserCog, Power, Search } from 'lucide-svelte';
  import { superUserApi } from '$lib/api';
  import { formatRupiah } from '$lib/utils';
  import { Button, Alert, Spinner, Modal, Input, Badge, EmptyState, Card, DataTable, ConfirmDialog } from '$lib/components';
  import type { Category, StudentStatusType, User, Role, Student } from '$lib/types';
  import { toast } from '$lib/stores/toast';

  type Tab = 'categories' | 'statuses' | 'users';
  let activeTab = $state<Tab>('categories');

  let categories  = $state<Category[]>([]);
  let statuses    = $state<StudentStatusType[]>([]);
  let users       = $state<User[]>([]);
  let roles       = $state<Role[]>([]);
  let loading     = $state(true);
  let error       = $state('');
  let submitting  = $state(false);
  let userPagination = $state<{ page: number; limit: number; total: number; pages: number } | null>(null);
  let userPage = $state(1);
  let userLimit = $state(10);

  let showCatModal = $state(false);
  let showCatCancelConfirm = $state(false);
  let editCatId    = $state('');
  let catName      = $state('');
  let catAmount    = $state('');
  let catDesc      = $state('');
  let catIsActive  = $state(true);
  let catIsSemester = $state(false);
  let catError     = $state('');
  const isCatEditing = $derived(!!editCatId);

  let catNameError = $state('');
  let catAmountError = $state('');

  let applyToType = $state<'none' | 'all' | 'selected'>('none');
  let studentSearch = $state('');
  let studentList = $state<Student[]>([]);
  let selectedStudentIDs = $state<string[]>([]);
  
  async function loadStudentsForSelection() {
    try {
      const res = await superUserApi.getStudents();
      studentList = res.data ?? [];
    } catch (e: unknown) {
      console.error('Failed to load students:', e);
    }
  }

  $effect(() => {
    if (showCatModal && applyToType === 'selected' && studentList.length === 0) {
      loadStudentsForSelection();
    }
  });

  const filteredStudents = $derived(
    studentList.filter(s => {
      const first = s.name?.first_name ?? '';
      const middle = s.name?.middle_name ?? '';
      const last = s.name?.last_name ?? '';
      const fullName = `${first} ${middle} ${last}`.toLowerCase();
      return fullName.includes(studentSearch.toLowerCase()) || 
             (s.student_number && s.student_number.toLowerCase().includes(studentSearch.toLowerCase()));
    })
  );

  let showStatModal     = $state(false);
  let showStatCancelConfirm = $state(false);
  let editStatId        = $state('');
  let statName          = $state('');
  let statDiscount      = $state('0');
  let statActiveBilling = $state('true');
  let statDesc          = $state('');
  let statError         = $state('');
  const isStatEditing = $derived(!!editStatId);

  let statNameError = $state('');
  let statDiscountError = $state('');

  let deleteConfirmOpen = $state(false);
  let deleteConfirmTitle = $state('');
  let deleteConfirmMessage = $state('');
  let deleteConfirmAction = $state<(() => void | Promise<void>) | null>(null);
  let deleting = $state(false);

  $effect(() => {
    if (statName.trim().toLowerCase() === 'abdi dalem') {
      statActiveBilling = 'false';
      statDiscount = '100';
    }
  });

  async function fetchData() {
    loading = true; error = '';
    try {
      const [resCat, resStat, resUsersPag, resRoles] = await Promise.all([
        superUserApi.getCategories(),
        superUserApi.getStatusTypes(),
        superUserApi.getAllUsersPaginated(userPage, userLimit),
        superUserApi.getAllRoles()
      ]);
      categories = resCat.data ?? [];
      statuses   = resStat.data ?? [];
      users      = resUsersPag.data?.users ?? [];
      userPagination = resUsersPag.data?.pagination ?? null;
      roles      = resRoles.data ?? [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat data.';
    } finally {
      loading = false;
    }
  }
  onMount(fetchData);

  function handleUserPageChange(newPage: number) {
    userPage = newPage;
    fetchData();
  }

  function handleUserLimitChange(newLimit: number) {
    userLimit = newLimit;
    userPage = 1;
    fetchData();
  }

  const roleOptions = $derived(roles.map(r => ({ value: r.id, label: r.name })));

  async function handleUpdateRole(userId: string, roleId: string) {
    try {
      await superUserApi.updateUserRole(userId, roleId);
      toast.success('Role user berhasil diperbarui');
      await fetchData();
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Gagal update role.');
    }
  }

  async function handleToggleActive(userId: string) {
    try {
      await superUserApi.toggleUserActive(userId);
      toast.success('Status user diperbarui');
      await fetchData();
    } catch (e: unknown) {
      toast.error(e instanceof Error ? e.message : 'Gagal mengubah status.');
    }
  }

  function openAddCategory() {
    editCatId = ''; catName = ''; catAmount = ''; catDesc = ''; catIsActive = true; catError = '';
    catNameError = ''; catAmountError = '';
    applyToType = 'none'; studentSearch = ''; selectedStudentIDs = [];
    catIsSemester = false;
    showCatModal = true;
  }
  function openEditCategory(cat: Category) {
    editCatId = cat.id; catName = cat.name; catAmount = String(cat.base_amount);
    catDesc = cat.description ?? ''; catIsActive = cat.is_active; catError = '';
    catNameError = ''; catAmountError = '';
    catIsSemester = !!cat.is_semester;
    showCatModal = true;
  }
  async function handleSaveCategory(e?: Event) {
    if (e) e.preventDefault();
    catError = ''; catNameError = ''; catAmountError = '';
    
    const nameStr = String(catName).trim();
    const amountStr = String(catAmount).trim();

    let hasError = false;
    if (!nameStr) {
      catNameError = 'Nama kategori wajib diisi.';
      hasError = true;
    }
    if (!amountStr) {
      catAmountError = 'Nominal wajib diisi.';
      hasError = true;
    } else {
      const amt = parseFloat(amountStr);
      if (isNaN(amt) || amt <= 0) {
        catAmountError = 'Nominal harus berupa angka positif.';
        hasError = true;
      }
    }

    if (hasError) return;

    submitting = true;
    try {
      const payload = { 
        name: nameStr, 
        base_amount: parseFloat(amountStr), 
        description: catDesc, 
        is_fixed: true, 
        is_active: catIsActive,
        is_semester: catIsSemester,
        apply_to_all: applyToType === 'all',
        student_ids: applyToType === 'selected' ? selectedStudentIDs : []
      };
      if (isCatEditing) {
        await superUserApi.updateCategory(editCatId, payload);
        toast.success('Kategori berhasil diperbarui');
      } else {
        await superUserApi.createCategory(payload);
        toast.success('Kategori berhasil ditambahkan');
      }
      showCatModal = false; 
      await fetchData();
    } catch (err: unknown) {
      catError = err instanceof Error ? err.message : 'Gagal menyimpan.';
    } finally { submitting = false; }
  }
  function handleDeleteCategory(id: string, name: string) {
    deleteConfirmTitle = 'Hapus Kategori';
    deleteConfirmMessage = `Apakah Anda yakin ingin menghapus kategori "${name}"? Tindakan ini tidak dapat dibatalkan, namun seluruh riwayat tagihan santri yang menggunakan kategori ini akan tetap tersimpan sebagai histori.`;
    deleteConfirmAction = async () => {
      deleting = true;
      try {
        await superUserApi.deleteCategory(id);
        toast.success('Kategori berhasil dihapus');
        await fetchData();
      } catch (e: unknown) {
        toast.error(e instanceof Error ? e.message : 'Gagal menghapus kategori.');
      } finally {
        deleting = false;
        deleteConfirmOpen = false;
      }
    };
    deleteConfirmOpen = true;
  }

  function openAddStatus() {
    editStatId = ''; statName = ''; statDiscount = '0'; statActiveBilling = 'true'; statDesc = ''; statError = '';
    statNameError = ''; statDiscountError = '';
    showStatModal = true;
  }
  function openEditStatus(stat: StudentStatusType) {
    editStatId = stat.id; statName = stat.name; statDiscount = String(stat.discount_percentage);
    statActiveBilling = stat.is_active_billing ? 'true' : 'false'; statDesc = stat.description ?? ''; statError = '';
    statNameError = ''; statDiscountError = '';
    showStatModal = true;
  }
  async function handleSaveStatus(e?: Event) {
    if (e) e.preventDefault();
    statError = ''; statNameError = ''; statDiscountError = '';
    
    const nameStr = String(statName).trim();
    const discountStr = String(statDiscount).trim();

    let hasError = false;
    if (!nameStr) {
      statNameError = 'Nama status wajib diisi.';
      hasError = true;
    }
    if (!discountStr) {
      statDiscountError = 'Diskon wajib diisi.';
      hasError = true;
    } else {
      const disc = parseFloat(discountStr);
      if (isNaN(disc) || disc < 0 || disc > 100) {
        statDiscountError = 'Diskon harus berupa angka antara 0% hingga 100%.';
        hasError = true;
      }
    }

    if (hasError) return;

    submitting = true;
    try {
      const payload = { 
        name: nameStr, 
        discount_percentage: parseFloat(discountStr), 
        is_active_billing: statActiveBilling === 'true', 
        description: statDesc 
      };
      if (isStatEditing) {
        await superUserApi.updateStatusType(editStatId, payload);
        toast.success('Tipe status berhasil diperbarui');
      } else {
        await superUserApi.createStatusType(payload);
        toast.success('Tipe status berhasil ditambahkan');
      }
      showStatModal = false; 
      await fetchData();
    } catch (err: unknown) {
      statError = err instanceof Error ? err.message : 'Gagal menyimpan.';
    } finally { submitting = false; }
  }
  function handleDeleteStatus(id: string, name: string) {
    deleteConfirmTitle = 'Hapus Tipe Status';
    deleteConfirmMessage = `Apakah Anda yakin ingin menghapus tipe status "${name}"? Tindakan ini tidak dapat dibatalkan, namun data santri dengan status ini tetap aman.`;
    deleteConfirmAction = async () => {
      deleting = true;
      try {
        await superUserApi.deleteStatusType(id);
        toast.success('Tipe status berhasil dihapus');
        await fetchData();
      } catch (e: unknown) {
        toast.error(e instanceof Error ? e.message : 'Gagal menghapus tipe status.');
      } finally {
        deleting = false;
        deleteConfirmOpen = false;
      }
    };
    deleteConfirmOpen = true;
  }
  // Form submission programmatically handled directly via onclick handlers
</script>

<svelte:head>
  <title>Master Data | Super User</title>
</svelte:head>

<div class="space-y-6 flex flex-col flex-1 min-h-0">
  <div>
    <h1 class="text-2xl font-black text-gray-900 tracking-tight">Master Data</h1>
    <p class="text-gray-500 text-sm mt-1">Kelola kategori, status santri, data santri, dan manajemen pengguna.</p>
  </div>

  <div class="flex gap-1 bg-gray-100 rounded-xl p-1 w-fit">
    {#each [
      { key: 'categories', label: 'Kategori', icon: Tag },
      { key: 'statuses',   label: 'Status Santri', icon: Users2 },
      { key: 'users',      label: 'Manajemen User', icon: UserCog }
    ] as tab}
      <button
        onclick={() => activeTab = tab.key as Tab}
        class="flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-all
          {activeTab === tab.key ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-500 hover:text-gray-700'}"
      >
        <tab.icon size={14} />
        {tab.label}
      </button>
    {/each}
  </div>

  {#if error}<Alert type="error" message={error} />{/if}

  {#if loading}
    <div class="flex justify-center py-12"><Spinner size="lg" /></div>
  {:else}

    {#if activeTab === 'categories'}
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="font-bold text-gray-900">Kategori Pembayaran</h2>
          <Button onclick={openAddCategory} variant="primary" size="sm">
            {#snippet children()}<Plus size={14} /> Tambah{/snippet}
          </Button>
        </div>
        {#if categories.length === 0}
          <Card><EmptyState title="Belum ada kategori" description="Tambahkan kategori seperti Syahriyyah Pondok, Syahriyyah Muhadhoroh, dll." /></Card>
        {:else}
          <div class="space-y-2">
            {#each categories as cat (cat.id)}
              <div class="bg-white border border-gray-200 rounded-xl p-4 flex items-center justify-between gap-4 hover:shadow-sm transition-all group">
                <div class="min-w-0">
                  <div class="flex items-center gap-2 flex-wrap">
                    <p class="font-semibold text-gray-900 text-sm">{cat.name}</p>
                    <Badge label={cat.is_active ? 'Aktif' : 'Nonaktif'} variant={cat.is_active ? 'success' : 'danger'} dot />
                    {#if cat.is_semester}
                      <Badge label="Semester" variant="info" dot />
                    {:else}
                      <Badge label="Bulanan" variant="warning" dot />
                    {/if}
                  </div>
                  <p class="text-green-700 font-bold text-base mt-0.5">{formatRupiah(cat.base_amount)}</p>
                  {#if cat.description}<p class="text-xs text-gray-400 mt-0.5 truncate">{cat.description}</p>{/if}
                </div>
                <div class="flex gap-1.5 flex-shrink-0">
                  <button onclick={() => openEditCategory(cat)} class="p-1.5 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors" title="Edit"><Edit size={14} /></button>
                  <button onclick={() => handleDeleteCategory(cat.id, cat.name)} class="p-1.5 rounded-lg bg-red-50 text-red-600 hover:bg-red-100 transition-colors" title="Hapus"><Trash2 size={14} /></button>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>

    {:else if activeTab === 'statuses'}
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="font-bold text-gray-900">Tipe Status Santri</h2>
          <Button onclick={openAddStatus} variant="primary" size="sm">
            {#snippet children()}<Plus size={14} /> Tambah{/snippet}
          </Button>
        </div>
        {#if statuses.length === 0}
          <Card><EmptyState title="Belum ada tipe status" description="Tambahkan seperti Reguler, Abdi Dalem, dll." /></Card>
        {:else}
          <div class="space-y-2">
            {#each statuses as stat (stat.id)}
              <div class="bg-white border border-gray-200 rounded-xl p-4 flex items-center justify-between gap-4 hover:shadow-sm transition-all group">
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-2 flex-wrap">
                    <p class="font-semibold text-gray-900 text-sm">{stat.name}</p>
                    <Badge label={stat.is_active_billing ? 'Dikenakan Tagihan' : 'Gratis Tagihan'} variant={stat.is_active_billing ? 'success' : 'warning'} dot />
                  </div>
                  <p class="text-xs text-gray-500 mt-0.5">Diskon: <strong class="text-purple-700">{stat.discount_percentage}%</strong>{#if stat.description} · {stat.description}{/if}</p>
                </div>
                <div class="flex gap-1.5 flex-shrink-0">
                  <button onclick={() => openEditStatus(stat)} class="p-1.5 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors" title="Edit"><Edit size={14} /></button>
                  <button onclick={() => handleDeleteStatus(stat.id, stat.name)} class="p-1.5 rounded-lg bg-red-50 text-red-600 hover:bg-red-100 transition-colors" title="Hapus"><Trash2 size={14} /></button>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>

    {:else if activeTab === 'users'}
      <div class="space-y-4 flex-1 flex flex-col min-h-0">
        <div class="flex-shrink-0">
          <h2 class="font-bold text-gray-900">Manajemen Pengguna</h2>
          <p class="text-xs text-gray-500 mt-0.5">Kelola semua akun pengguna dan hak akses (role) mereka.</p>
        </div>
        <DataTable
          pagination={userPagination}
          onPageChange={handleUserPageChange}
          onLimitChange={handleUserLimitChange}
          isEmpty={users.length === 0}
          emptyTitle="Tidak ada data pengguna"
          emptyDescription="Belum ada pengguna terdaftar."
          paginationLabel="pengguna"
        >
          {#snippet children()}
            <table class="w-full text-sm">
              <thead>
                <tr class="bg-gray-50 border-b border-gray-100">
                  <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Username</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Role / Hak Akses</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Status</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Terdaftar</th>
                  <th class="px-4 py-3 text-right text-xs font-semibold text-gray-500 uppercase">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                {#each users as user (user.id)}
                  <tr class="hover:bg-gray-50 transition-colors">
                    <td class="px-4 py-3 font-mono text-sm font-semibold text-gray-900">{user.username}</td>
                    <td class="px-4 py-3">
                      <select
                        value={user.role_id}
                        onchange={(e) => handleUpdateRole(user.id, (e.target as HTMLSelectElement).value)}
                        class="text-xs px-2 py-1.5 border border-gray-200 rounded-lg bg-white focus:ring-2 focus:ring-green-500 focus:border-green-500 outline-none"
                      >
                        {#each roles as role}
                          <option value={role.id}>{role.name}</option>
                        {/each}
                      </select>
                    </td>
                    <td class="px-4 py-3">
                      <Badge label={user.is_active ? 'Aktif' : 'Nonaktif'} variant={user.is_active ? 'success' : 'danger'} dot />
                    </td>
                    <td class="px-4 py-3 text-xs text-gray-400">
                      {new Date(user.created_at).toLocaleDateString('id-ID')}
                    </td>
                    <td class="px-4 py-3 text-right">
                      <button
                        onclick={() => handleToggleActive(user.id)}
                        class="p-1.5 rounded-lg border transition-colors {user.is_active ? 'bg-red-50 text-red-600 hover:bg-red-100 border-red-200' : 'bg-green-50 text-green-700 hover:bg-green-100 border-green-200'}"
                        title={user.is_active ? 'Nonaktifkan akun' : 'Aktifkan akun'}
                      >
                        <Power size={14} />
                      </button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          {/snippet}
        </DataTable>
      </div>
    {/if}
  {/if}
</div>

<Modal bind:open={showCatModal} title={isCatEditing ? 'Edit Kategori' : 'Tambah Kategori'} size="sm" confirmCloseMessage="Apakah Anda yakin ingin membatalkan? Perubahan yang belum disimpan pada kategori ini akan hilang.">
  {#snippet children()}
    {#if catError}<Alert type="error" message={catError} class="mb-4" />{/if}
    <form id="cat-form" onsubmit={handleSaveCategory} class="space-y-4" novalidate>
      <Input id="cat-name" label="Nama Kategori" bind:value={catName} required error={catNameError} oninput={() => catNameError = ''} />
      <Input id="cat-amount" label="Nominal (Rp)" type="number" bind:value={catAmount} required error={catAmountError} oninput={() => catAmountError = ''} />
      <Input id="cat-desc" label="Deskripsi (opsional)" bind:value={catDesc} />
      
      {#if !isCatEditing}
        <div class="flex flex-col gap-1.5">
          <label for="apply-to" class="text-xs font-semibold text-gray-600 uppercase tracking-wider">Terapkan Kategori Baru Ke</label>
          <select id="apply-to" bind:value={applyToType}
            class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-sm outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500">
            <option value="none">Tidak satupun (atur nanti di data santri)</option>
            <option value="all">Semua santri aktif</option>
            <option value="selected">Santri pilihan</option>
          </select>
        </div>

        {#if applyToType === 'selected'}
          <div class="space-y-2 border border-gray-200 rounded-lg p-3 bg-gray-50">
            <div class="relative">
              <input
                type="text"
                placeholder="Cari santri..."
                bind:value={studentSearch}
                class="w-full pl-8 pr-3 py-1.5 rounded-md bg-white border border-gray-200 text-xs outline-none focus:ring-1 focus:ring-emerald-800"
              />
              <span class="absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400">
                <Search size={14} />
              </span>
            </div>

            <div class="max-h-40 overflow-y-auto divide-y divide-slate-100 border border-slate-200 rounded-md bg-white">
              {#if filteredStudents.length === 0}
                <p class="text-center text-xs text-gray-500 py-3">Tidak ada santri ditemukan</p>
              {:else}
                {#each filteredStudents as student}
                  {@const fullName = `${student.name?.first_name ?? ''} ${student.name?.middle_name ?? ''} ${student.name?.last_name ?? ''}`}
                  <label class="flex items-center gap-2 px-3 py-2 hover:bg-slate-50 cursor-pointer text-xs text-slate-700">
                    <input
                      type="checkbox"
                      value={student.id}
                      checked={selectedStudentIDs.includes(student.id)}
                      onchange={(e) => {
                        const target = e.target as HTMLInputElement;
                        if (target.checked) {
                          selectedStudentIDs.push(student.id);
                        } else {
                          selectedStudentIDs = selectedStudentIDs.filter(id => id !== student.id);
                        }
                      }}
                      class="rounded text-emerald-800 focus:ring-emerald-800 h-3.5 w-3.5"
                    />
                    <div class="flex flex-col">
                      <span class="font-medium text-slate-800">{fullName}</span>
                      <span class="text-[10px] text-slate-500">{student.student_number}</span>
                    </div>
                  </label>
                {/each}
              {/if}
            </div>
            
            <div class="flex justify-between items-center text-[10px] text-slate-500 pt-1">
              <span>Terpilih: {selectedStudentIDs.length} santri</span>
              <button 
                type="button" 
                onclick={() => selectedStudentIDs = studentList.map(s => s.id)}
                class="text-emerald-800 hover:underline font-semibold"
              >
                Pilih Semua
              </button>
            </div>
          </div>
        {/if}
      {/if}

      <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-200">
        <div class="flex flex-col">
          <p class="text-xs font-bold text-gray-700" id="cat-semester-label">Tipe Semester</p>
          <p class="text-[10px] text-gray-500">Aktifkan jika tagihan ini khusus semester (dihasilkan per semester)</p>
        </div>
        <button 
          type="button" 
          role="switch"
          aria-checked={catIsSemester}
          aria-labelledby="cat-semester-label"
          onclick={() => catIsSemester = !catIsSemester}
          class="relative inline-flex h-5 w-10 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 {catIsSemester ? 'bg-green-600' : 'bg-gray-200'}"
        >
          <span class="sr-only">Ubah tipe semester kategori</span>
          <span class="inline-block h-3 w-3 transform rounded-full bg-white transition-transform {catIsSemester ? 'translate-x-6' : 'translate-x-1'}"></span>
        </button>
      </div>

      <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-200">
        <p class="text-xs font-bold text-gray-700" id="cat-active-label">Status Aktif</p>
        <button 
          type="button" 
          role="switch"
          aria-checked={catIsActive}
          aria-labelledby="cat-active-label"
          onclick={() => catIsActive = !catIsActive}
          class="relative inline-flex h-5 w-10 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 {catIsActive ? 'bg-green-600' : 'bg-gray-200'}"
        >
          <span class="sr-only">Ubah status aktif kategori</span>
          <span class="inline-block h-3 w-3 transform rounded-full bg-white transition-transform {catIsActive ? 'translate-x-6' : 'translate-x-1'}"></span>
        </button>
      </div>
    </form>
  {/snippet}
  {#snippet footer()}
    <div class="flex justify-end gap-2">
      <Button onclick={() => showCatCancelConfirm = true} variant="outline" size="md">{#snippet children()}Batal{/snippet}</Button>
      <Button onclick={() => handleSaveCategory()} variant="primary" size="md" loading={submitting}>
        {#snippet children()}{isCatEditing ? 'Simpan' : 'Tambah'}{/snippet}
      </Button>
    </div>
  {/snippet}
</Modal>

<Modal bind:open={showStatModal} title={isStatEditing ? 'Edit Tipe Status' : 'Tambah Tipe Status'} size="sm" confirmCloseMessage="Apakah Anda yakin ingin membatalkan? Perubahan yang belum disimpan pada tipe status ini akan hilang.">
  {#snippet children()}
    {#if statError}<Alert type="error" message={statError} class="mb-4" />{/if}
    <form id="stat-form" onsubmit={handleSaveStatus} class="space-y-4" novalidate>
      <Input id="stat-name" label="Nama Status" bind:value={statName} required error={statNameError} oninput={() => statNameError = ''} />
      <Input id="stat-discount" label="Diskon (%)" type="number" bind:value={statDiscount} required error={statDiscountError} oninput={() => statDiscountError = ''} />
      <div class="flex flex-col gap-1.5">
        <label for="stat-billing" class="text-xs font-semibold text-gray-600 uppercase tracking-wider">Dikenakan Tagihan?</label>
        <select id="stat-billing" bind:value={statActiveBilling}
          class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-sm outline-none focus:ring-2 focus:ring-green-500/20 focus:border-green-500">
          <option value="true">Ya — Dikenakan tagihan Syahriyyah</option>
          <option value="false">Tidak — Bebas tagihan</option>
        </select>
      </div>
      <Input id="stat-desc" label="Deskripsi (opsional)" bind:value={statDesc} />
    </form>
  {/snippet}
  {#snippet footer()}
    <div class="flex justify-end gap-2">
      <Button onclick={() => showStatCancelConfirm = true} variant="outline" size="md">{#snippet children()}Batal{/snippet}</Button>
      <Button onclick={() => handleSaveStatus()} variant="primary" size="md" loading={submitting}>
        {#snippet children()}{isStatEditing ? 'Simpan' : 'Tambah'}{/snippet}
      </Button>
    </div>
  {/snippet}
</Modal>

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

<ConfirmDialog
  bind:open={showCatCancelConfirm}
  title="Konfirmasi Tindakan"
  message="Apakah Anda yakin ingin membatalkan? Perubahan yang belum disimpan pada kategori ini akan hilang."
  confirmText="Ya, Batalkan"
  cancelText="Kembali"
  variant="warning"
  onConfirm={() => {
    showCatCancelConfirm = false;
    showCatModal = false;
  }}
/>

<ConfirmDialog
  bind:open={showStatCancelConfirm}
  title="Konfirmasi Tindakan"
  message="Apakah Anda yakin ingin membatalkan? Perubahan yang belum disimpan pada tipe status ini akan hilang."
  confirmText="Ya, Batalkan"
  cancelText="Kembali"
  variant="warning"
  onConfirm={() => {
    showStatCancelConfirm = false;
    showStatModal = false;
  }}
/>
