<script lang="ts">
  import { onMount } from 'svelte';
  import { superUserApi } from '$lib/api';
  import { formatDate } from '$lib/utils';
  import { Spinner, Alert, DataTable } from '$lib/components';
  import { Activity } from 'lucide-svelte';
  import type { ActivityLog } from '$lib/types';

  let logs    = $state<ActivityLog[]>([]);
  let loading = $state(true);
  let error   = $state('');
  let pagination = $state<{ page: number; limit: number; total: number; pages: number } | null>(null);
  let page = $state(1);
  let limit = $state(10);

  async function fetchData() {
    loading = true;
    error = '';
    try {
      const res = await superUserApi.getActivityLogsPaginated(page, limit);
      logs = res.data?.logs ?? [];
      pagination = res.data?.pagination ?? null;
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat log aktivitas.';
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

  onMount(fetchData);
</script>

<svelte:head>
  <title>Log Aktivitas | Dashboard Super User</title>
</svelte:head>

<div class="space-y-6 flex flex-col flex-1 min-h-0">
  <div class="flex-shrink-0">
    <h1 class="text-2xl font-black text-gray-900 tracking-tight">Log Aktivitas & Audit Trail</h1>
    <p class="text-gray-500 text-sm mt-1">Rekam jejak seluruh aktivitas pembayaran dan perubahan data pada sistem.</p>
  </div>

  {#if error}
    <Alert type="error" message={error} />
  {:else if loading}
    <Spinner size="lg" />
  {:else}
    <DataTable
      pagination={pagination}
      onPageChange={handlePageChange}
      onLimitChange={handleLimitChange}
      isEmpty={logs.length === 0}
      emptyTitle="Belum ada log aktivitas"
      emptyDescription="Log akan muncul saat terjadi aktivitas di sistem."
      paginationLabel="log"
    >
      {#snippet children()}
        <table class="w-full text-sm" aria-label="Tabel log aktivitas">
          <thead>
            <tr class="bg-gray-50 border-b border-gray-100">
              <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Waktu</th>
              <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Pengguna</th>
              <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left">Aksi</th>
              <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left hidden md:table-cell">Entitas</th>
              <th scope="col" class="px-5 py-3 text-xs font-semibold text-gray-500 uppercase tracking-wider text-left hidden lg:table-cell">IP Address</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            {#each logs as log (log.id)}
              <tr class="hover:bg-gray-50 transition-colors">
                <td class="px-5 py-3.5 text-gray-500 text-xs whitespace-nowrap">{formatDate(log.created_at, true)}</td>
                <td class="px-5 py-3.5">
                  <p class="font-semibold text-gray-900">{log.user?.first_name ?? 'Sistem'}</p>
                  <p class="text-xs text-gray-500 font-mono">{log.user?.username ?? '-'}</p>
                </td>
                <td class="px-5 py-3.5">
                  <div class="flex items-center gap-2">
                    <Activity size={14} class="text-green-600 flex-shrink-0" aria-hidden="true" />
                    <span class="text-gray-800">{log.action}</span>
                  </div>
                </td>
                <td class="px-5 py-3.5 hidden md:table-cell">
                  {#if log.entity_name}
                    <span class="text-xs font-mono text-gray-500">{log.entity_name}</span>
                  {:else}
                    <span class="text-gray-400">—</span>
                  {/if}
                </td>
                <td class="px-5 py-3.5 hidden lg:table-cell">
                  <span class="text-xs font-mono text-gray-400">{log.ip_address ?? '—'}</span>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/snippet}
    </DataTable>
  {/if}
</div>
