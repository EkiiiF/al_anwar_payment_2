<script lang="ts">
  import { onMount } from 'svelte';
  import { User, Mail, Phone, Lock, KeyRound, GraduationCap, MapPin, Info, Shield, Edit, Check } from 'lucide-svelte';
  import { authApi, guardianApi } from '$lib/api';
  import { getInitials, formatDate } from '$lib/utils';
  import { Button, Input, Select, Alert, Spinner, Card, Badge } from '$lib/components';
  import type { User as UserType, Student } from '$lib/types';
  import { getMuhadhorohLabel } from '$lib/types/student.types';

  let {
    roleLabel = 'Pengguna',
    roleColor = 'green'
  }: {
    roleLabel?: string;
    roleColor?: 'green' | 'blue' | 'purple';
  } = $props();

  let user         = $state<UserType | null>(null);
  let loading      = $state(true);
  let error        = $state('');

  // Guardian student data (complete student entity)
  let student      = $state<Student | null>(null);



  let showPwForm   = $state(false);
  let pwSaving     = $state(false);
  let pwMsg        = $state('');
  let pwMsgType    = $state<'success' | 'error'>('success');
  let currentPw    = $state('');
  let newPw        = $state('');
  let confirmPw    = $state('');

  // Edit profile states
  let editFirstName = $state('');
  let editMiddleName = $state('');
  let editLastName = $state('');
  let editEmail = $state('');
  let editPhone = $state('');
  let editGender = $state('');
  let editBirthDay = $state<number | ''>('');
  let editBirthMonth = $state('');
  let editBirthYear = $state<number | ''>('');
  let editAddress = $state('');

  let isEditing     = $state(false);
  let profileSaving = $state(false);
  let profileMsg    = $state('');
  let profileMsgType = $state<'success' | 'error'>('success');

  // Error states
  let editFirstNameError = $state('');
  let editEmailError = $state('');
  let editPhoneError = $state('');
  let birthDayError = $state('');
  
  let currentPwError = $state('');
  let newPwError = $state('');
  let confirmPwError = $state('');

  function clearProfileErrors() {
    editFirstNameError = '';
    editEmailError = '';
    editPhoneError = '';
    birthDayError = '';
  }

  function clearPwErrors() {
    currentPwError = '';
    newPwError = '';
    confirmPwError = '';
  }

  function startEdit() {
    if (!user) return;
    editFirstName = user.first_name || '';
    editMiddleName = user.middle_name || '';
    editLastName = user.last_name || '';
    editEmail = user.email || '';
    editPhone = user.phone_number || '';
    editGender = user.gender || 'L';
    if (user.birth_date) {
      const parts = user.birth_date.split('T')[0].split('-');
      if (parts.length === 3) {
        editBirthYear = parseInt(parts[0], 10);
        editBirthMonth = String(parseInt(parts[1], 10));
        editBirthDay = parseInt(parts[2], 10);
      } else {
        editBirthDay = '';
        editBirthMonth = '';
        editBirthYear = '';
      }
    } else {
      editBirthDay = '';
      editBirthMonth = '';
      editBirthYear = '';
    }
    editAddress = user.address || '';
    isEditing = true;
    profileMsg = '';
    clearProfileErrors();
  }

  function cancelEdit() {
    isEditing = false;
    profileMsg = '';
    clearProfileErrors();
  }

  async function handleUpdateProfile(e: SubmitEvent) {
    e.preventDefault();
    clearProfileErrors();
    profileMsg = '';

    let hasError = false;
    if (!editFirstName.trim()) {
      editFirstNameError = 'Nama depan wajib diisi.';
      hasError = true;
    }
    if (editEmail.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(editEmail.trim())) {
      editEmailError = 'Format email tidak valid.';
      hasError = true;
    }
    if (editPhone.trim() && !/^\+?[0-9]{9,15}$/.test(editPhone.trim().replace(/[\s-]/g, ''))) {
      editPhoneError = 'Format nomor HP tidak valid (9-15 digit angka).';
      hasError = true;
    }

    let birthDateStr = '';
    if (editBirthDay || editBirthMonth || editBirthYear) {
      if (!editBirthDay || !editBirthMonth || !editBirthYear) {
        birthDayError = 'Lengkapi tanggal, bulan, dan tahun lahir.';
        hasError = true;
      } else {
        const day = Number(editBirthDay);
        const month = Number(editBirthMonth);
        const year = Number(editBirthYear);
        if (day < 1 || day > 31) {
          birthDayError = 'Tanggal harus di antara 1 dan 31.';
          hasError = true;
        } else if (year < 1900 || year > new Date().getFullYear()) {
          birthDayError = 'Tahun lahir tidak valid.';
          hasError = true;
        } else {
          birthDateStr = `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}T00:00:00Z`;
        }
      }
    }

    if (hasError) {
      profileMsg = 'Mohon periksa kembali input form Anda.';
      profileMsgType = 'error';
      return;
    }

    profileSaving = true;
    try {
      const res = await authApi.updateProfile({
        first_name: editFirstName,
        middle_name: editMiddleName,
        last_name: editLastName,
        email: editEmail,
        phone: editPhone,
        gender: editGender,
        birth_date: birthDateStr,
        address: editAddress
      });
      user = res.data;
      profileMsg = 'Profil berhasil diperbarui!';
      profileMsgType = 'success';
      isEditing = false;
    } catch (e: unknown) {
      profileMsg = e instanceof Error ? e.message : 'Gagal memperbarui profil.';
      profileMsgType = 'error';
    } finally {
      profileSaving = false;
    }
  }

  const avatarColors = {
    green:  'from-green-500 to-green-700',
    blue:   'from-blue-500 to-blue-700',
    purple: 'from-purple-500 to-purple-700'
  };

  onMount(async () => {
    try {
      const res = await authApi.getProfile();
      user = res.data;

      // If guardian, also fetch student data
      if (roleColor === 'blue') {
        try {
          const dashRes = await guardianApi.getDashboard();
          student = dashRes.data?.student ?? null;
        } catch { /* ignore */ }
      }
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat data profil.';
    } finally {
      loading = false;
    }
  });



  async function handleChangePassword(e: SubmitEvent) {
    e.preventDefault();
    clearPwErrors();
    pwMsg = '';

    let hasError = false;
    if (!currentPw) {
      currentPwError = 'Password saat ini wajib diisi.';
      hasError = true;
    }
    if (!newPw) {
      newPwError = 'Password baru wajib diisi.';
      hasError = true;
    } else if (newPw.length < 8) {
      newPwError = 'Password baru minimal 8 karakter.';
      hasError = true;
    }
    if (!confirmPw) {
      confirmPwError = 'Konfirmasi password baru wajib diisi.';
      hasError = true;
    } else if (newPw !== confirmPw) {
      confirmPwError = 'Password baru dan konfirmasi tidak cocok.';
      hasError = true;
    }

    if (hasError) {
      pwMsg = 'Mohon periksa kembali input form Anda.';
      pwMsgType = 'error';
      return;
    }

    pwSaving = true;
    try {
      await authApi.changePassword({
        current_password: currentPw,
        new_password:     newPw,
        confirm_password: confirmPw
      });
      pwMsg     = 'Password berhasil diubah!';
      pwMsgType = 'success';
      currentPw = ''; newPw = ''; confirmPw = '';
      showPwForm = false;
    } catch (e: unknown) {
      pwMsg     = e instanceof Error ? e.message : 'Gagal mengubah password.';
      pwMsgType = 'error';
    } finally {
      pwSaving = false;
    }
  }

  const displayName = $derived(
    user ? [user.first_name, user.last_name].filter(Boolean).join(' ') || user.username : ''
  );

  // For guardian: show student name as primary display
  const studentDisplayName = $derived(
    student ? [student.name.first_name, student.name.middle_name, student.name.last_name].filter(Boolean).join(' ') : ''
  );
</script>

<svelte:head>
  <title>Profil Saya | {roleLabel} - Al-Anwar Payment</title>
  <meta name="description" content="Kelola informasi profil dan ubah password akun Anda." />
</svelte:head>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-black text-gray-900 tracking-tight">Profil Saya</h1>
    <p class="text-gray-500 text-sm mt-1">Kelola informasi pribadi dan keamanan akun.</p>
  </div>

  {#if loading}
    <Spinner size="lg" />
  {:else if error}
    <Alert type="error" message={error} />
  {:else if user}

    <Card>
      <div class="flex items-center gap-5">
        <div class="relative flex-shrink-0">
          <div class="w-20 h-20 rounded-2xl bg-gradient-to-br {avatarColors[roleColor]} flex items-center justify-center text-white text-2xl font-black shadow-sm">
            {getInitials((roleColor === 'blue' && studentDisplayName) ? studentDisplayName : (displayName || user.username))}
          </div>
          <div class="absolute -bottom-1 -right-1 w-5 h-5 rounded-full bg-green-500 border-2 border-white flex items-center justify-center">
            <div class="w-2 h-2 rounded-full bg-white"></div>
          </div>
        </div>

        <div class="min-w-0 flex-1">
          <h2 class="text-xl font-black text-gray-900 leading-tight">
            {roleColor === 'blue' && studentDisplayName ? studentDisplayName : (displayName || user.username)}
          </h2>
          <p class="text-sm text-gray-500 mt-1.5">
            {#if roleColor === 'blue'}
              NIS: {student?.student_number ?? user.username}
            {:else}
              Username: @{user.username}
            {/if}
          </p>
        </div>
      </div>
    </Card>

    <!-- Detail Profil Pengguna Section (Visible only for Super User and Pengasuh) -->
    {#if roleColor === 'green' || roleColor === 'purple'}
      <Card>
        <div class="flex items-center justify-between mb-6 pb-4 border-b border-gray-100">
        <div class="flex items-center gap-3">
          <User size={20} class={roleColor === 'green' ? 'text-green-600' : roleColor === 'blue' ? 'text-blue-600' : 'text-purple-600'} aria-hidden="true" />
          <div>
            <h3 class="font-black text-gray-900 text-lg">Informasi Profil Lengkap</h3>
            <p class="text-gray-500 text-xs mt-0.5">Detail data diri akun Anda yang terdaftar di sistem.</p>
          </div>
        </div>
        {#if !isEditing && (roleColor === 'green' || roleColor === 'purple')}
          <Button onclick={startEdit} variant="outline" size="sm">
            {#snippet children()}
              <Edit size={14} aria-hidden="true" />
              Edit Profil
            {/snippet}
          </Button>
        {/if}
      </div>

      {#if profileMsg}
        <Alert type={profileMsgType} message={profileMsg} class="mb-4" />
      {/if}

      {#if isEditing}
        <form onsubmit={handleUpdateProfile} id="profile-form" class="space-y-4" novalidate>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Input
              id="first-name"
              label="Nama Depan"
              bind:value={editFirstName}
              required
              placeholder="Nama depan"
              error={editFirstNameError}
              oninput={() => editFirstNameError = ''}
            />
            <Input
              id="middle-name"
              label="Nama Tengah"
              bind:value={editMiddleName}
              placeholder="Nama tengah (opsional)"
            />
            <Input
              id="last-name"
              label="Nama Belakang"
              bind:value={editLastName}
              placeholder="Nama belakang (opsional)"
            />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              id="email"
              label="Email"
              type="email"
              bind:value={editEmail}
              placeholder="nama@domain.com"
              error={editEmailError}
              oninput={() => editEmailError = ''}
            />
            <Input
              id="phone"
              label="No. Telepon / WA"
              type="tel"
              bind:value={editPhone}
              placeholder="Contoh: 08123456789"
              error={editPhoneError}
              oninput={() => editPhoneError = ''}
            />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Select
              id="gender"
              label="Jenis Kelamin"
              bind:value={editGender}
              options={[
                { value: 'L', label: 'Laki-laki' },
                { value: 'P', label: 'Perempuan' }
              ]}
            />
            <div class="flex flex-col gap-1.5 w-full">
              <span class="text-xs font-semibold text-slate-600 uppercase tracking-wider ml-0.5">Tanggal Lahir</span>
              <div class="grid grid-cols-3 gap-3">
                <input 
                  id="birth-day" 
                  type="number" 
                  bind:value={editBirthDay} 
                  min="1" 
                  max="31" 
                  class="w-full px-3 py-2.5 rounded-lg bg-white border border-slate-200 text-slate-900 text-sm focus:ring-1 focus:ring-emerald-800 focus:border-emerald-800 outline-none transition-colors" 
                  placeholder="Tanggal" 
                  oninput={() => birthDayError = ''} 
                />
                <select 
                  id="birth-month" 
                  bind:value={editBirthMonth} 
                  class="w-full px-3 py-2.5 rounded-lg bg-white border border-slate-200 text-slate-900 text-sm focus:ring-1 focus:ring-emerald-800 focus:border-emerald-800 outline-none transition-colors"
                  onchange={() => birthDayError = ''}
                >
                  <option value="">Bulan</option>
                  {#each ['Januari','Februari','Maret','April','Mei','Juni','Juli','Agustus','September','Oktober','November','Desember'] as name, i}
                    <option value={String(i + 1)}>{name}</option>
                  {/each}
                </select>
                <input 
                  id="birth-year" 
                  type="number" 
                  bind:value={editBirthYear} 
                  min="1900" 
                  max={new Date().getFullYear()} 
                  class="w-full px-3 py-2.5 rounded-lg bg-white border border-slate-200 text-slate-900 text-sm focus:ring-1 focus:ring-emerald-800 focus:border-emerald-800 outline-none transition-colors" 
                  placeholder="Tahun" 
                  oninput={() => birthDayError = ''} 
                />
              </div>
              {#if birthDayError}
                <p class="text-xs text-red-500 mt-1 ml-0.5">{birthDayError}</p>
              {/if}
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="address" class="text-xs font-medium text-slate-600 uppercase tracking-wider">
              Alamat
            </label>
            <textarea
              id="address"
              bind:value={editAddress}
              placeholder="Alamat lengkap tempat tinggal"
              rows="3"
              class="w-full px-3 py-2.5 rounded-lg bg-white border text-slate-900 text-sm outline-none transition-colors duration-200
                placeholder:text-slate-400 border-slate-200
                focus:ring-1 focus:ring-emerald-800 focus:border-emerald-800
                disabled:opacity-40 disabled:cursor-not-allowed disabled:bg-slate-50"
            ></textarea>
          </div>
        </form>

        <div class="flex justify-end gap-2 mt-6 pt-4 border-t border-gray-100">
          <Button onclick={cancelEdit} variant="outline" size="sm">
            {#snippet children()}Batal{/snippet}
          </Button>
          <Button type="submit" form="profile-form" variant="primary" size="sm" loading={profileSaving}>
            {#snippet children()}
              <Check size={14} aria-hidden="true" />
              Simpan Profil
            {/snippet}
          </Button>
        </div>
      {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Col 1: Kontak & Akun -->
          <div class="space-y-4">
            <div>
              <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Nama Lengkap</span>
              <span class="font-medium text-slate-900 text-sm">{[user.first_name, user.middle_name, user.last_name].filter(Boolean).join(' ') || '-'}</span>
            </div>

            <div>
              <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Email</span>
              <div class="flex items-center gap-2 text-slate-900 text-sm font-medium">
                <Mail size={14} class="text-slate-400 flex-shrink-0" />
                <span>{user.email || '-'}</span>
              </div>
            </div>

            <div>
              <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">No. Telepon / WA</span>
              <div class="flex items-center gap-2 text-slate-900 text-sm font-medium">
                <Phone size={14} class="text-slate-400 flex-shrink-0" />
                <span>{user.phone_number || '-'}</span>
              </div>
            </div>
          </div>

          <!-- Col 2: Info Lainnya & Peran -->
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Username</span>
                <span class="font-medium text-slate-900 text-sm">@{user.username}</span>
              </div>
              <div>
                <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Hak Akses</span>
                <Badge 
                  label={roleLabel} 
                  variant={roleColor === 'green' ? 'success' : roleColor === 'blue' ? 'info' : 'purple'} 
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Jenis Kelamin</span>
                <span class="font-medium text-slate-900 text-sm">
                  {user.gender === 'L' ? 'Laki-laki' : user.gender === 'P' ? 'Perempuan' : '-'}
                </span>
              </div>
              <div>
                <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Tgl Lahir</span>
                <span class="font-medium text-slate-900 text-sm">
                  {user.birth_date ? formatDate(user.birth_date) : '-'}
                </span>
              </div>
            </div>

            <div>
              <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Alamat</span>
              <div class="flex items-start gap-2 text-slate-900 text-sm font-medium">
                <MapPin size={14} class="text-slate-400 mt-0.5 flex-shrink-0" />
                <p class="text-sm text-slate-900 whitespace-pre-line leading-relaxed">{user.address || '-'}</p>
              </div>
            </div>
          </div>
        </div>
      {/if}
      </Card>
    {/if}

    <!-- Guardian: Data Santri Section -->
    {#if roleColor === 'blue' && student}
      <Card>
        <div class="flex items-center gap-3 mb-6 pb-4 border-b border-gray-100">
          <GraduationCap size={20} class="text-blue-600" aria-hidden="true" />
          <div>
            <h3 class="font-black text-gray-900 text-lg">Data Santri Lengkap</h3>
            <p class="text-gray-500 text-xs mt-0.5">Informasi lengkap santri yang terhubung dengan akun wali Anda.</p>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Col 1: Data Pribadi -->
          <div class="space-y-3">
            <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5 pb-1 border-b border-gray-50">
              <User size={14} class="text-blue-600" /> Data Pribadi
            </h4>
            <div class="bg-blue-50/30 rounded-xl border border-blue-100 p-4 space-y-3.5 text-xs">
              <div>
                <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Nama Lengkap</span>
                <span class="font-medium text-slate-900 text-sm">{[student.name.first_name, student.name.middle_name, student.name.last_name].filter(Boolean).join(' ')}</span>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">NIS (No. Induk)</span>
                  <span class="font-medium text-blue-700 text-sm">{student.student_number}</span>
                </div>
                <div>
                  <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">NIK</span>
                  <span class="font-medium text-slate-900 text-sm">{student.nik || '-'}</span>
                </div>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Tgl Lahir</span>
                  <span class="font-medium text-slate-900 text-sm">{student.birth_date ? formatDate(student.birth_date) : '-'}</span>
                </div>
                <div>
                  <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Jenis Kelamin</span>
                  <span class="font-medium text-slate-900 text-sm">{(student.gender === 'L' || student.gender === 'M') ? 'Laki-laki' : 'Perempuan'}</span>
                </div>
              </div>
              <div>
                <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Kelas & Semester</span>
                <span class="font-medium text-slate-900 text-sm">{getMuhadhorohLabel(student.muhadhoroh_level, student.current_semester)}</span>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Kategori Status</span>
                  <div class="mt-0.5">
                    <Badge label={student.status?.name ?? '-'} variant="purple" />
                  </div>
                </div>
                <div>
                  <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Diskon SPP</span>
                  <span class="font-bold text-purple-700 text-sm">{student.status?.discount_percentage ?? 0}%</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Col 2: Data Wali / Guardian -->
          <div class="space-y-3">
            <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5 pb-1 border-b border-gray-50">
              <Phone size={14} class="text-blue-600" /> Data Wali Santri
            </h4>
            <div class="space-y-3">
              {#if student.guardians && student.guardians.length > 0}
                {#each student.guardians as g}
                  <div class="bg-blue-50/30 rounded-xl border border-blue-100 p-4 space-y-3 text-xs">
                    <div>
                      <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Nama Wali ({g.relation})</span>
                      <span class="font-medium text-slate-900 text-sm">{[g.name.first_name, g.name.middle_name, g.name.last_name].filter(Boolean).join(' ')}</span>
                    </div>
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">No. Telepon / WA</span>
                        <span class="font-medium text-slate-900 text-sm">{g.phone || '-'}</span>
                      </div>
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Email</span>
                        <span class="font-medium text-slate-900 text-sm truncate block">{g.email || '-'}</span>
                      </div>
                    </div>
                    {#if g.user?.username}
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Username Login</span>
                        <span class="font-medium text-blue-700 text-sm">@{g.user.username}</span>
                      </div>
                    {/if}
                  </div>
                {/each}
              {:else}
                <div class="bg-gray-50 rounded-xl border border-slate-200 p-4 text-center py-6 text-slate-400 italic text-xs">
                  Belum ada data wali santri
                </div>
              {/if}
            </div>
          </div>

          <!-- Col 3: Alamat Lengkap -->
          <div class="space-y-3">
            <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5 pb-1 border-b border-gray-50">
              <MapPin size={14} class="text-blue-600" /> Alamat Lengkap
            </h4>
            <div class="space-y-3">
              {#if student.addresses && student.addresses.length > 0}
                {#each student.addresses as addr}
                  <div class="bg-blue-50/30 rounded-xl border border-blue-100 p-4 space-y-3 text-xs relative overflow-hidden">
                    {#if addr.is_primary}
                      <span class="absolute top-2 right-2 px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 text-[10px] font-bold border border-blue-200">Utama</span>
                    {/if}
                    <div>
                      <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Alamat</span>
                      <p class="font-medium text-slate-900 text-sm whitespace-pre-line leading-relaxed">{addr.address_line || '-'}</p>
                    </div>
                    <div class="grid grid-cols-2 gap-3 mt-2 pt-2 border-t border-blue-100/50">
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Kelurahan / Desa</span>
                        <span class="font-medium text-slate-900 text-sm">{addr.village || '-'}</span>
                      </div>
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Kecamatan</span>
                        <span class="font-medium text-slate-900 text-sm">{addr.district || '-'}</span>
                      </div>
                    </div>
                    <div class="grid grid-cols-2 gap-3">
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Kota / Kabupaten</span>
                        <span class="font-medium text-slate-900 text-sm">{addr.city || '-'}</span>
                      </div>
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Provinsi</span>
                        <span class="font-medium text-slate-900 text-sm">{addr.province || '-'}</span>
                      </div>
                    </div>
                    <div class="grid grid-cols-2 gap-3">
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Kodepos</span>
                        <span class="font-medium text-slate-900 text-sm">{addr.postal_code || '-'}</span>
                      </div>
                      <div>
                        <span class="text-xs font-semibold text-slate-500 block uppercase tracking-wider mb-1">Negara</span>
                        <span class="font-medium text-slate-900 text-sm">{addr.country || '-'}</span>
                      </div>
                    </div>
                  </div>
                {/each}
              {:else}
                <div class="bg-gray-50 rounded-xl border border-slate-200 p-4 text-center py-6 text-slate-400 italic text-xs">
                  Belum ada data alamat lengkap
                </div>
              {/if}
            </div>
          </div>
        </div>
      </Card>
    {/if}


    <Card>
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="p-2 rounded-lg bg-amber-100 border border-amber-200">
            <Lock size={16} class="text-amber-700" aria-hidden="true" />
          </div>
          <div>
            <h3 class="font-bold text-gray-900 text-sm">Ubah Password</h3>
            <p class="text-xs text-gray-500">Pastikan password baru minimal 8 karakter.</p>
          </div>
        </div>
        {#if !showPwForm}
          <Button onclick={() => { showPwForm = true; pwMsg = ''; }} variant="outline" size="sm">
            {#snippet children()}
              <KeyRound size={14} aria-hidden="true" />
              Ubah
            {/snippet}
          </Button>
        {/if}
      </div>

      {#if showPwForm}
        <div class="mt-5 pt-5 border-t border-gray-100">
          {#if pwMsg}
            <Alert type={pwMsgType} message={pwMsg} class="mb-4" />
          {/if}

          <form onsubmit={handleChangePassword} id="pw-form" class="space-y-4" novalidate>
            <Input
              id="current-pw"
              label="Password Saat Ini"
              type="password"
              bind:value={currentPw}
              required
              placeholder="Password lama"
              error={currentPwError}
              oninput={() => currentPwError = ''}
            />
            <Input
              id="new-pw"
              label="Password Baru"
              type="password"
              bind:value={newPw}
              required
              placeholder="Minimal 8 karakter"
              helper="Kombinasikan huruf, angka, dan simbol"
              error={newPwError}
              oninput={() => newPwError = ''}
            />
            <Input
              id="confirm-pw"
              label="Konfirmasi Password Baru"
              type="password"
              bind:value={confirmPw}
              required
              placeholder="Ulangi password baru"
              error={confirmPwError}
              oninput={() => confirmPwError = ''}
            />
          </form>

          <div class="flex justify-end gap-2 mt-4">
            <Button onclick={() => { showPwForm = false; pwMsg = ''; currentPw = ''; newPw = ''; confirmPw = ''; }} variant="outline" size="sm">
              {#snippet children()}Batal{/snippet}
            </Button>
            <Button type="submit" form="pw-form" variant="primary" size="sm" loading={pwSaving}>
              {#snippet children()}
                <Lock size={14} aria-hidden="true" />
                Simpan Password
              {/snippet}
            </Button>
          </div>
        </div>
      {/if}
    </Card>

    <Card class="border-blue-200 bg-blue-50">
      <div class="flex items-start gap-3">
        <Shield size={18} class="text-blue-600 flex-shrink-0 mt-0.5" aria-hidden="true" />
        <div>
          <p class="text-sm font-semibold text-blue-900 mb-1">Keamanan Akun</p>
          <ul class="text-xs text-blue-700 space-y-1 list-disc list-inside">
            <li>Jangan bagikan password kepada siapapun termasuk petugas.</li>
            <li>Gunakan password yang unik dan tidak mudah ditebak.</li>
            <li>Selalu logout setelah menggunakan perangkat bersama.</li>
          </ul>
        </div>
      </div>
    </Card>
  {/if}
</div>
