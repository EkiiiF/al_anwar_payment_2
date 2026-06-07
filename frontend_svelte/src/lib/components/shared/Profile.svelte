<script lang="ts">
  import { onMount } from 'svelte';
  import { User, Mail, Phone, Lock, Edit, Save, X, Shield, KeyRound, Check } from 'lucide-svelte';
  import { authApi } from '$lib/api';
  import { getInitials } from '$lib/utils';
  import { Button, Input, Alert, Spinner, Card } from '$lib/components';
  import type { User as UserType } from '$lib/types';

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

  let editMode     = $state(false);
  let saving       = $state(false);
  let saveMsg      = $state('');
  let saveMsgType  = $state<'success' | 'error'>('success');

  let editFirstName  = $state('');
  let editMiddleName = $state('');
  let editLastName   = $state('');
  let editEmail      = $state('');
  let editPhone      = $state('');
  let editGender     = $state('L');
  let editBirthDate  = $state('');
  let editAddress    = $state('');

  let showPwForm   = $state(false);
  let pwSaving     = $state(false);
  let pwMsg        = $state('');
  let pwMsgType    = $state<'success' | 'error'>('success');
  let currentPw    = $state('');
  let newPw        = $state('');
  let confirmPw    = $state('');

  const avatarColors = {
    green:  'from-green-500 to-green-700',
    blue:   'from-blue-500 to-blue-700',
    purple: 'from-purple-500 to-purple-700'
  };

  onMount(async () => {
    try {
      const res = await authApi.getProfile();
      user = res.data;
      fillEditForm(user);
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat data profil.';
    } finally {
      loading = false;
    }
  });

  function fillEditForm(u: UserType) {
    editFirstName  = u.first_name ?? '';
    editMiddleName = u.middle_name ?? '';
    editLastName   = u.last_name  ?? '';
    editEmail      = u.email       ?? '';
    editPhone      = u.phone_number ?? '';
    editGender     = u.gender ?? 'L';
    editBirthDate  = u.birth_date ? u.birth_date.substring(0, 10) : '';
    editAddress    = u.address ?? '';
  }

  function openEdit() {
    if (user) fillEditForm(user);
    editMode = true;
    saveMsg  = '';
  }

  function cancelEdit() {
    editMode = false;
    saveMsg  = '';
  }

  async function handleSaveProfile(e: SubmitEvent) {
    e.preventDefault();
    saving  = true;
    saveMsg = '';
    try {
      const res = await authApi.updateProfile({
        first_name:  editFirstName,
        middle_name: editMiddleName,
        last_name:   editLastName,
        email:       editEmail,
        phone:       editPhone,
        gender:      editGender,
        birth_date:  editBirthDate,
        address:     editAddress
      });
      user       = res.data as UserType;
      editMode   = false;
      saveMsg    = 'Profil berhasil diperbarui!';
      saveMsgType = 'success';
    } catch (e: unknown) {
      saveMsg    = e instanceof Error ? e.message : 'Gagal menyimpan profil.';
      saveMsgType = 'error';
    } finally {
      saving = false;
    }
  }

  async function handleChangePassword(e: SubmitEvent) {
    e.preventDefault();
    if (newPw !== confirmPw) {
      pwMsg     = 'Password baru dan konfirmasi tidak cocok.';
      pwMsgType = 'error';
      return;
    }
    if (newPw.length < 8) {
      pwMsg     = 'Password baru minimal 8 karakter.';
      pwMsgType = 'error';
      return;
    }
    pwSaving = true;
    pwMsg    = '';
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
      <div class="flex flex-col sm:flex-row sm:items-center gap-5">
        <div class="relative flex-shrink-0">
          <div class="w-20 h-20 rounded-2xl bg-gradient-to-br {avatarColors[roleColor]} flex items-center justify-center text-white text-2xl font-black shadow-lg">
            {getInitials(displayName || user.username)}
          </div>
          <div class="absolute -bottom-1 -right-1 w-5 h-5 rounded-full bg-green-500 border-2 border-white flex items-center justify-center">
            <div class="w-2 h-2 rounded-full bg-white"></div>
          </div>
        </div>

        <div class="min-w-0 flex-1">
          <h2 class="text-xl font-black text-gray-900 leading-tight">{displayName || user.username}</h2>
          <p class="text-sm text-gray-500 mt-0.5 font-mono">@{user.username}</p>
          <div class="flex flex-wrap gap-2 mt-2">
            <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold
              {roleColor === 'green' ? 'bg-green-100 text-green-700' :
               roleColor === 'blue' ? 'bg-blue-100 text-blue-700' :
               'bg-purple-100 text-purple-700'}">
              <Shield size={11} aria-hidden="true" />
              {roleLabel}
            </span>
            <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-gray-100 text-gray-600">
              <Check size={11} aria-hidden="true" />
              Akun Aktif
            </span>
          </div>
        </div>

        {#if !editMode}
          <Button onclick={openEdit} variant="outline" size="sm" class="flex-shrink-0">
            {#snippet children()}
              <Edit size={14} aria-hidden="true" />
              Edit Profil
            {/snippet}
          </Button>
        {/if}
      </div>
    </Card>

    {#if saveMsg && !editMode}
      <Alert type={saveMsgType} message={saveMsg} />
    {/if}

    {#if editMode}
      <Card>
        <div class="flex items-center justify-between mb-5">
          <h3 class="font-bold text-gray-900 flex items-center gap-2">
            <Edit size={16} class="text-green-600" aria-hidden="true" />
            Edit Informasi Profil
          </h3>
          <button
            onclick={cancelEdit}
            class="p-1.5 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-100 transition-colors"
            aria-label="Batal edit"
          >
            <X size={16} />
          </button>
        </div>

        {#if saveMsg}
          <Alert type={saveMsgType} message={saveMsg} class="mb-4" />
        {/if}

        <form onsubmit={handleSaveProfile} id="profile-form" class="space-y-4" novalidate>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <Input id="edit-first" label="Nama Depan" bind:value={editFirstName} placeholder="Nama depan" />
            <Input id="edit-middle" label="Nama Tengah" bind:value={editMiddleName} placeholder="Nama tengah (opsional)" />
            <Input id="edit-last"  label="Nama Belakang" bind:value={editLastName} placeholder="Nama belakang" />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <Input
              id="edit-email"
              label="Email"
              type="email"
              bind:value={editEmail}
              placeholder="email@example.com"
              helper="Email digunakan untuk notifikasi tagihan"
            />
            <Input
              id="edit-phone"
              label="Nomor HP / WhatsApp"
              type="tel"
              bind:value={editPhone}
              placeholder="08xxxxxxxxxx"
            />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="flex flex-col gap-1.5">
              <label for="edit-gender" class="text-xs font-semibold text-gray-600 uppercase tracking-wider ml-0.5">Jenis Kelamin</label>
              <select
                id="edit-gender"
                bind:value={editGender}
                class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-sm focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none"
              >
                <option value="L">Laki-laki</option>
                <option value="P">Perempuan</option>
              </select>
            </div>
            <Input
              id="edit-birthdate"
              label="Tanggal Lahir"
              type="date"
              bind:value={editBirthDate}
            />
          </div>

          <div class="flex flex-col gap-1.5">
            <label for="edit-address" class="text-xs font-semibold text-gray-600 uppercase tracking-wider ml-0.5">Alamat</label>
            <textarea
              id="edit-address"
              bind:value={editAddress}
              rows="3"
              placeholder="Alamat lengkap..."
              class="w-full px-3 py-2.5 rounded-lg bg-white border border-gray-300 text-sm focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none"
            ></textarea>
          </div>
        </form>

        <div class="flex justify-end gap-2 mt-5 pt-5 border-t border-gray-100">
          <Button onclick={cancelEdit} variant="outline" size="md">
            {#snippet children()}Batal{/snippet}
          </Button>
          <Button type="submit" form="profile-form" variant="primary" size="md" loading={saving}>
            {#snippet children()}
              <Save size={15} aria-hidden="true" />
              Simpan Perubahan
            {/snippet}
          </Button>
        </div>
      </Card>
    {:else}
      <Card>
        <h3 class="font-bold text-gray-900 mb-4 flex items-center gap-2">
          <User size={16} class="text-gray-500" aria-hidden="true" />
          Informasi Pribadi
        </h3>

        <div class="space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Nama Depan</p>
              <p class="text-sm font-semibold text-gray-900">{user.first_name || '—'}</p>
            </div>
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Nama Tengah</p>
              <p class="text-sm font-semibold text-gray-900">{user.middle_name || '—'}</p>
            </div>
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Nama Belakang</p>
              <p class="text-sm font-semibold text-gray-900">{user.last_name || '—'}</p>
            </div>
          </div>

          <div class="border-t border-gray-100"></div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Username / NIS</p>
              <div class="flex items-center gap-2">
                <User size={14} class="text-gray-400" aria-hidden="true" />
                <p class="text-sm font-mono font-semibold text-gray-900">{user.username}</p>
              </div>
            </div>
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Jenis Kelamin</p>
              <p class="text-sm font-semibold text-gray-900">
                {user.gender === 'P' ? 'Perempuan' : 'Laki-laki'}
              </p>
            </div>
          </div>

          <div class="border-t border-gray-100"></div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Email</p>
              <div class="flex items-center gap-2">
                <Mail size={14} class="text-gray-400" aria-hidden="true" />
                <p class="text-sm text-gray-900">{user.email || '—'}</p>
              </div>
            </div>
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Nomor HP / WhatsApp</p>
              <div class="flex items-center gap-2">
                <Phone size={14} class="text-gray-400" aria-hidden="true" />
                <p class="text-sm text-gray-900">{user.phone_number || '—'}</p>
              </div>
            </div>
          </div>

          <div class="border-t border-gray-100"></div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Tanggal Lahir</p>
              <p class="text-sm text-gray-900">
                {user.birth_date ? new Date(user.birth_date).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }) : '—'}
              </p>
            </div>
            <div>
              <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">Alamat</p>
              <p class="text-sm text-gray-900 whitespace-pre-line">{user.address || '—'}</p>
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
            />
            <Input
              id="new-pw"
              label="Password Baru"
              type="password"
              bind:value={newPw}
              required
              placeholder="Minimal 8 karakter"
              helper="Kombinasikan huruf, angka, dan simbol"
            />
            <Input
              id="confirm-pw"
              label="Konfirmasi Password Baru"
              type="password"
              bind:value={confirmPw}
              required
              placeholder="Ulangi password baru"
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
