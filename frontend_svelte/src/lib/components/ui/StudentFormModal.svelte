<script lang="ts">
  import { Modal, Input, Button, Alert, Select } from '$lib/components';
  import type { Student, StudentStatusType, Category } from '$lib/types';
  import { GUARDIAN_RELATION_OPTIONS } from '$lib/types/student.types';
  import { superUserApi } from '$lib/api';

  let { 
    open = $bindable(false), 
    student = null, 
    statuses = [],
    categories = [],
    onSaved 
  }: { 
    open: boolean; 
    student: Student | null; 
    statuses: StudentStatusType[];
    categories: Category[];
    onSaved: () => void;
  } = $props();

  let submitting = $state(false);
  let error = $state('');

  let nis = $state('');
  let fullName = $state('');
  let nik = $state('');
  let birthDate = $state('');
  let statusId = $state('');
  let gender = $state('L');

  let addressLine = $state('');
  let village = $state('');
  let district = $state('');
  let city = $state('');
  let province = $state('');
  let country = $state('Indonesia');
  let postalCode = $state('');

  let createGuardian = $state(false);
  let guardianName = $state('');
  let guardianPhone = $state('');
  let guardianEmail = $state('');
  let guardianRelation = $state('Orang Tua');

  $effect(() => {
    if (open) {
      error = '';
      if (student) {
        nis = student.student_number;
        fullName = [student.name.first_name, student.name.middle_name, student.name.last_name].filter(Boolean).join(' ');
        nik = student.nik || '';
        birthDate = student.birth_date ? student.birth_date.split('T')[0] : '';
        statusId = student.status_id;
        gender = student.gender || 'L';

        const primaryAddr = student.addresses?.find(a => a.is_primary);
        addressLine = primaryAddr?.address_line ?? '';
        village     = primaryAddr?.village ?? '';
        district    = primaryAddr?.district ?? '';
        city        = primaryAddr?.city ?? '';
        province    = primaryAddr?.province ?? '';
        country     = primaryAddr?.country ?? 'Indonesia';
        postalCode  = primaryAddr?.postal_code ?? '';
        
        const g = student.guardians?.[0];
        if (g) {
          createGuardian = true;
          guardianName = [g.name.first_name, g.name.last_name].filter(Boolean).join(' ');
          guardianPhone = g.phone || '';
          guardianEmail = g.email || '';
          guardianRelation = g.relation || 'Orang Tua';
        } else {
          createGuardian = false;
          guardianName = ''; guardianPhone = ''; guardianEmail = ''; guardianRelation = 'Orang Tua';
        }
      } else {
        nis = ''; fullName = ''; nik = ''; birthDate = '';
        statusId = statuses.length > 0 ? statuses[0].id : '';
        gender = 'L';
        addressLine = ''; village = ''; district = ''; city = ''; province = ''; country = 'Indonesia'; postalCode = '';
        createGuardian = false; guardianName = ''; guardianPhone = ''; guardianEmail = ''; guardianRelation = 'Orang Tua';
      }
    }
  });

  async function handleSave(e: SubmitEvent) {
    e.preventDefault();
    submitting = true;
    error = '';
    
    const payload = {
      nis,
      full_name: fullName,
      nik,
      birth_date: birthDate,
      status_type_id: statusId,
      gender,
      address_line: addressLine,
      village, district, city, province, country, postal_code: postalCode,
      create_new_guardian: createGuardian,
      guardian_name: guardianName,
      guardian_phone: guardianPhone,
      guardian_email: guardianEmail,
      guardian_relation: guardianRelation,
    };

    try {
      if (student) {
        await superUserApi.updateStudent(student.id, payload);
      } else {
        await superUserApi.createStudent(payload);
      }
      open = false;
      onSaved();
    } catch (err: any) {
      error = err.message || 'Gagal menyimpan data santri.';
    } finally {
      submitting = false;
    }
  }
</script>

<Modal bind:open title={student ? 'Edit Santri Lengkap' : 'Tambah Santri Lengkap'} size="lg">
  {#snippet children()}
    {#if error}<Alert type="error" message={error} class="mb-4" />{/if}
    <form id="student-form" onsubmit={handleSave} class="space-y-6">
      
      <div>
        <h3 class="text-sm font-bold text-gray-900 border-b pb-2 mb-3">Data Pribadi Santri</h3>
        <div class="grid grid-cols-2 gap-4">
          <Input id="nis" label="NIS" bind:value={nis} required />
          <Input id="nik" label="NIK" bind:value={nik} />
        </div>
        <div class="grid grid-cols-2 gap-4 mt-3">
          <Input id="full_name" label="Nama Lengkap" bind:value={fullName} required />
          <Input id="birth_date" type="date" label="Tanggal Lahir" bind:value={birthDate} />
        </div>
        <div class="grid grid-cols-2 gap-4 mt-3">
          <div class="space-y-1.5">
            <label for="gender" class="block text-sm font-medium text-gray-700">Jenis Kelamin</label>
            <select id="gender" bind:value={gender} class="w-full border-gray-300 rounded-lg shadow-sm focus:border-green-500 focus:ring-green-500 text-sm">
              <option value="L">Laki-laki</option>
              <option value="P">Perempuan</option>
            </select>
          </div>
          <div class="space-y-1.5">
            <label for="status_id" class="block text-sm font-medium text-gray-700">Tipe Status / Diskon</label>
            <select id="status_id" bind:value={statusId} class="w-full border-gray-300 rounded-lg shadow-sm focus:border-green-500 focus:ring-green-500 text-sm" required>
              {#each statuses as stat}
                <option value={stat.id}>{stat.name}</option>
              {/each}
            </select>
          </div>
        </div>
      </div>

      <div>
        <h3 class="text-sm font-bold text-gray-900 border-b pb-2 mb-3">Alamat Santri</h3>
        <Input id="address" label="Alamat (Jalan/Gang/RT/RW)" bind:value={addressLine} placeholder="Jl. Raya No.1 RT 01/RW 02" />
        <div class="grid grid-cols-2 gap-4 mt-3">
          <Input id="village"  label="Desa / Kelurahan" bind:value={village}  placeholder="Desa" />
          <Input id="district" label="Kecamatan"        bind:value={district} placeholder="Kecamatan" />
        </div>
        <div class="grid grid-cols-2 gap-4 mt-3">
          <Input id="city"     label="Kota / Kabupaten" bind:value={city}     placeholder="Kota" />
          <Input id="province" label="Provinsi"         bind:value={province} placeholder="Provinsi" />
        </div>
        <div class="grid grid-cols-2 gap-4 mt-3">
          <Input id="postalCode" label="Kode Pos" bind:value={postalCode} placeholder="00000" />
          <Input id="country"    label="Negara"   bind:value={country}    placeholder="Indonesia" />
        </div>
      </div>

      <div class="bg-gray-50 p-4 rounded-xl border border-gray-200">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-sm font-bold text-gray-900">Data Wali Santri</h3>
          <div class="flex items-center gap-2">
            <label for="create_guardian" class="text-xs font-semibold text-gray-700 cursor-pointer">Sertakan Wali Santri</label>
            <button 
              type="button" 
              role="switch"
              aria-checked={createGuardian}
              id="create_guardian"
              onclick={() => createGuardian = !createGuardian}
              class="relative inline-flex h-5 w-10 items-center rounded-full transition-colors focus:outline-none {createGuardian ? 'bg-green-600' : 'bg-gray-300'}"
            >
              <span class="sr-only">Toggle Sertakan Wali Santri</span>
              <span class="inline-block h-3 w-3 transform rounded-full bg-white transition-transform {createGuardian ? 'translate-x-6' : 'translate-x-1'}"></span>
            </button>
          </div>
        </div>

        {#if createGuardian}
          <div class="grid grid-cols-2 gap-4">
            <Input id="g_name" label="Nama Wali" bind:value={guardianName} required={createGuardian} />
            <Input id="g_phone" label="No. Telepon / WhatsApp" bind:value={guardianPhone} />
            <Input id="g_email" type="email" label="Email (Opsional)" bind:value={guardianEmail} />
            <Select
              id="g_relation"
              label="Hubungan dengan Santri"
              bind:value={guardianRelation}
              options={GUARDIAN_RELATION_OPTIONS}
            />
          </div>
          {#if !student}
            <p class="text-[11px] text-gray-500 mt-2">
              * Username wali santri akan menggunakan NIS santri secara otomatis dan password default "password123".
            </p>
          {/if}
        {/if}
      </div>

    </form>
  {/snippet}
  {#snippet footer()}
    <div class="flex justify-end gap-2 w-full">
      <Button onclick={() => open = false} variant="outline" size="md">{#snippet children()}Batal{/snippet}</Button>
      <Button type="submit" form="student-form" variant="primary" size="md" loading={submitting}>
        {#snippet children()}{student ? 'Simpan Perubahan' : 'Tambah Santri'}{/snippet}
      </Button>
    </div>
  {/snippet}
</Modal>
