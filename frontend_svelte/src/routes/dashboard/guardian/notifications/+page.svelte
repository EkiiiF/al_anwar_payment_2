<script lang="ts">
  import { onMount } from 'svelte';
  import { guardianApi } from '$lib/api';
  import { formatDate } from '$lib/utils';
  import { Spinner, Alert, Card, EmptyState, Badge, Paginator } from '$lib/components';
  import type { Notification } from '$lib/types';
  import { Bell, Check } from 'lucide-svelte';

  let notifications = $state<Notification[]>([]);
  let loading       = $state(true);
  let error         = $state('');

  let page = $state(1);
  let limit = $state(5);

  const totalNotifications = $derived(notifications.length);
  const totalPages = $derived(Math.ceil(totalNotifications / limit) || 1);
  const paginatedNotifications = $derived(
    notifications.slice((page - 1) * limit, page * limit)
  );

  onMount(async () => {
    try {
      const res = await guardianApi.getNotifications();
      notifications = res.data ?? [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat notifikasi.';
    } finally {
      loading = false;
    }
  });

  async function markRead(id: string) {
    try {
      await guardianApi.markNotificationRead(id);
      notifications = notifications.map(n => n.id === id ? { ...n, is_read: true } : n);
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('notification-read'));
      }
    } catch { /* abaikan */ }
  }

  function getTypeVariant(type: string | undefined): 'success' | 'warning' | 'info' {
    if (type === 'payment_success') return 'success';
    if (type === 'billing')         return 'warning';
    return 'info';
  }
</script>

<svelte:head>
  <title>Notifikasi | Guardian - Al-Anwar Payment</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Notifikasi</h1>
      <p class="text-gray-500 text-sm mt-1">Pemberitahuan tagihan dan konfirmasi pembayaran.</p>
    </div>
    {#if notifications.some(n => !n.is_read)}
      <Badge label="{notifications.filter(n => !n.is_read).length} Baru" variant="warning" dot />
    {/if}
  </div>

  {#if error}
    <Alert type="error" message={error} />
  {:else if loading}
    <Spinner size="lg" />
  {:else if notifications.length === 0}
    <Card>
      <EmptyState
        title="Tidak ada notifikasi"
        description="Notifikasi tagihan dan pembayaran akan muncul di sini."
      />
    </Card>
  {:else}
    <div class="space-y-4">
      <div class="space-y-3">
        {#each paginatedNotifications as notif (notif.id)}
        <div
          class="group relative p-4 rounded-xl border transition-all duration-200
            {notif.is_read ? 'bg-white border-slate-200/60 shadow-sm' : 'bg-blue-50/50 border-blue-200/60 shadow-sm'}"
        >
          {#if !notif.is_read}
            <div class="absolute top-4 right-4 w-2 h-2 rounded-full bg-blue-600"></div>
          {/if}

          <div class="flex items-start gap-3.5 pr-6">
            <div class="p-2 rounded-lg flex-shrink-0 border
              {notif.type === 'payment_success' ? 'bg-emerald-50 border-emerald-200 text-emerald-700' :
               notif.type === 'billing'          ? 'bg-amber-50 border-amber-200 text-amber-700' :
               'bg-blue-50 border-blue-200 text-blue-700'}">
              <Bell size={16} aria-hidden="true" />
            </div>

            <div class="min-w-0 flex-1">
              <div class="flex items-start justify-between gap-2 flex-wrap mb-1">
                <p class="font-bold text-slate-900 text-base">{notif.title}</p>
                <Badge
                  label={notif.type === 'payment_success' ? 'Pembayaran' : notif.type === 'billing' ? 'Tagihan' : 'Info'}
                  variant={getTypeVariant(notif.type)}
                  class="text-xs py-0.5"
                />
              </div>
              <p class="text-slate-600 text-sm mt-1 leading-relaxed">{notif.message}</p>
              <div class="flex items-center justify-between mt-3.5 pt-2 border-t border-slate-100">
                <p class="text-xs text-slate-400">{formatDate(notif.created_at, true)}</p>
                {#if !notif.is_read}
                  <button
                    onclick={() => markRead(notif.id)}
                    class="flex items-center gap-1 text-xs text-blue-700 hover:text-blue-800 transition-colors font-bold cursor-pointer"
                  >
                    <Check size={12} aria-hidden="true" />
                    Tandai dibaca
                  </button>
                {/if}
              </div>
            </div>
          </div>
        </div>
      {/each}
      </div>
      <Paginator
        page={page}
        limit={limit}
        total={totalNotifications}
        pages={totalPages}
        label="notifikasi"
        onPageChange={(p) => page = p}
        onLimitChange={(l) => { limit = l; page = 1; }}
      />
    </div>
  {/if}
</div>
