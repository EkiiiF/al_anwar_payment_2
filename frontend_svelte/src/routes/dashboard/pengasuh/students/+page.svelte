<script lang="ts">
  import { onMount } from 'svelte';
  import { Search, ChevronDown, ChevronUp, MapPin, Phone, Mail, User, Info, CheckCircle2, XCircle } from 'lucide-svelte';
  import { pengasuhApi } from '$lib/api';
  import { formatDate } from '$lib/utils';
  import { Alert, Spinner, Badge, EmptyState, Card } from '$lib/components';
  import type { Student } from '$lib/types';
  import { getMuhadhorohLabel } from '$lib/types/student.types';

  // ── Data State ─────────────────────────────────────────────
  let students = $state<Student[]>([]);
  let loading  = $state(true);
  let error    = $state('');

  // ── Search State ───────────────────────────────────────────
  let search = $state('');
  let searchTimeout: any;

  function handleSearch() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(fetchData, 500);
  }

  // ── Expanded Row (detail view) ─────────────────────────────
  let expandedId = $state<string | null>(null);
  function toggleExpand(id: string) {
    expandedId = expandedId === id ? null : id;
  }

  // ── Fetch ──────────────────────────────────────────────────
  async function fetchData() {
    loading = true;
    error = '';
    try {
      const res = await pengasuhApi.getStudents(search);
      students = res.data ?? [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat data santri.';
    } finally {
      loading = false;
    }
  }

  onMount(fetchData);
</script>

<svelte:head>
  <title>Data Santri | Dashboard Pengasuh</title>
  <meta name="description" content="Monitoring data santri Pondok Pesantren Al-Anwar secara real-time." />
</svelte:head>

<div class="space-y-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <div>
      <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-purple-100 border border-purple-200 text-purple-700 text-xs font-semibold uppercase tracking-wider mb-2">
        <span class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse" aria-hidden="true"></span>
        Mode Hanya Lihat
      </div>
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Data Santri</h1>
      <p class="text-gray-500 text-sm mt-1">Pantau perkembangan status akademik (Muhadhoroh) dan detail data wali santri.</p>
    </div>
  </div>

  <Card class="!p-4 border-l-4 border-l-purple-500">
    <div class="relative">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
      <input
        type="text"
        placeholder="Cari santri berdasarkan nama atau NIS..."
        class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-gray-200 focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500 transition-all text-sm"
        bind:value={search}
        oninput={handleSearch}
      />
    </div>
  </Card>

  {#if error}
    <Alert type="error" message={error} />
  {/if}

  {#if loading}
    <Spinner size="lg" />
  {:else}
    <Card padding={false}>
      <div class="overflow-x-auto">
        <table class="w-full text-sm text-left" aria-label="Tabel monitoring data santri">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-100">
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">NIS</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">Santri</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider">Kelas & Semester</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider hidden md:table-cell">Wali / Guardian</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider hidden lg:table-cell">Status Keaktifan</th>
              <th scope="col" class="px-5 py-3.5 text-xs font-semibold text-gray-500 uppercase tracking-wider text-right">Detail</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            {#each students as student (student.id)}
              <tr class="hover:bg-purple-50/20 transition-colors cursor-pointer" onclick={() => toggleExpand(student.id)}>
                <td class="px-5 py-4 font-mono text-purple-600 font-semibold text-xs">{student.student_number}</td>
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
                    <div>
                      <p class="text-sm font-bold text-gray-900">Muhadhoroh {student.muhadhoroh_level}</p>
                      <p class="text-xs text-purple-600 font-semibold">Semester {student.current_semester}</p>
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
                  <button
                    onclick={(e) => { e.stopPropagation(); toggleExpand(student.id); }}
                    class="p-2 rounded-lg bg-gray-50 text-gray-500 hover:bg-purple-50 hover:text-purple-600 border border-gray-200 transition-colors"
                    title="Lihat Detail"
                  >
                    {#if expandedId === student.id}
                      <ChevronUp size={15} />
                    {:else}
                      <ChevronDown size={15} />
                    {/if}
                  </button>
                </td>
              </tr>

              {#if expandedId === student.id}
                <tr>
                  <td colspan="6" class="bg-purple-50/10 px-5 py-5 border-l-4 border-l-purple-400">
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                      <div class="space-y-2">
                        <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                          <User size={12} class="text-purple-600" /> Data Pribadi
                        </h4>
                        <div class="bg-white rounded-lg border border-purple-100 p-3 space-y-1.5 text-xs">
                          <div class="flex justify-between"><span class="text-gray-500">NIK</span><span class="font-semibold text-gray-900">{student.nik || '-'}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Tgl Lahir</span><span class="font-semibold text-gray-900">{student.birth_date ? formatDate(student.birth_date) : '-'}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Gender</span><span class="font-semibold text-gray-900">{student.gender === 'M' ? 'Laki-laki' : 'Perempuan'}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Kelas</span><span class="font-bold text-purple-700">{getMuhadhorohLabel(student.muhadhoroh_level, student.current_semester)}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Kategori</span><span class="font-semibold text-gray-900">{student.status?.name ?? '-'}</span></div>
                          <div class="flex justify-between"><span class="text-gray-500">Diskon</span><span class="font-semibold text-purple-700">{student.status?.discount_percentage ?? 0}%</span></div>
                        </div>
                      </div>

                      <div class="space-y-2">
                        <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                          <Phone size={12} class="text-purple-600" /> Data Wali / Guardian
                        </h4>
                        {#if student.guardians && student.guardians.length > 0}
                          {#each student.guardians as guardian}
                            <div class="bg-white rounded-lg border border-purple-100 p-3 space-y-1.5 text-xs">
                              <p class="font-bold text-gray-900">{[guardian.name.first_name, guardian.name.middle_name, guardian.name.last_name].filter(Boolean).join(' ')}</p>
                              <div class="flex items-center gap-1.5 text-gray-500"><Phone size={10} /> {guardian.phone || '-'}</div>
                              <div class="flex items-center gap-1.5 text-gray-500"><Mail size={10} /> {guardian.email || '-'}</div>
                              <div class="flex justify-between"><span class="text-gray-500">Username Login</span><span class="font-mono font-semibold text-purple-600">{guardian.user?.username ?? '-'}</span></div>
                              <div class="flex justify-between"><span class="text-gray-500">Hubungan</span><span class="font-semibold text-gray-900">{guardian.relation}</span></div>
                            </div>
                          {/each}
                        {:else}
                          <p class="text-xs text-gray-400 italic">Belum ada data guardian</p>
                        {/if}
                      </div>

                      <div class="space-y-2">
                        <h4 class="text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-1.5">
                          <MapPin size={12} class="text-purple-600" /> Alamat Lengkap
                        </h4>
                        {#if student.addresses && student.addresses.length > 0}
                          {#each student.addresses as addr}
                            <div class="bg-white rounded-lg border border-purple-100 p-3 space-y-1 text-xs">
                              <p class="font-semibold text-gray-900">{addr.address_line || '-'}</p>
                              {#if addr.city || addr.province}
                                <p class="text-gray-500">{[addr.city, addr.province].filter(Boolean).join(', ')}</p>
                              {/if}
                              {#if addr.is_primary}
                                <span class="inline-flex mt-1.5 px-2 py-0.5 rounded-full bg-purple-50 text-purple-600 text-[10px] font-bold border border-purple-100">Utama</span>
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
            {:else}
              <tr>
                <td colspan="6">
                  <EmptyState
                    title="Tidak ada data santri ditemukan"
                    description="Belum ada santri terdaftar atau kata kunci pencarian Anda tidak cocok."
                  />
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </Card>
  {/if}
</div>
