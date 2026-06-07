<script lang="ts">
  import { onMount } from 'svelte';
  import { guardianApi } from '$lib/api';
  import { formatDate } from '$lib/utils';
  import { Spinner, Alert, Card, EmptyState, Badge } from '$lib/components';
  import type { Notification } from '$lib/types';
  import { Bell, Check } from 'lucide-svelte';

  let notifications = $state<Notification[]>([]);
  let loading       = $state(true);
  let error         = $state('');

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
    <div class="space-y-3">
      {#each notifications as notif (notif.id)}
        <div
          class="group relative p-5 rounded-xl border transition-all
            {notif.is_read ? 'bg-white border-gray-200' : 'bg-blue-50 border-blue-200'}"
        >
          {#if !notif.is_read}
            <div class="absolute top-4 right-4 w-2 h-2 rounded-full bg-blue-500"></div>
          {/if}

          <div class="flex items-start gap-4 pr-6">
            <div class="p-2.5 rounded-xl flex-shrink-0 border
              {notif.type === 'payment_success' ? 'bg-green-100 border-green-200 text-green-700' :
               notif.type === 'billing'          ? 'bg-amber-100 border-amber-200 text-amber-700' :
               'bg-blue-100 border-blue-200 text-blue-700'}">
              <Bell size={18} aria-hidden="true" />
            </div>

            <div class="min-w-0 flex-1">
              <div class="flex items-start justify-between gap-2 flex-wrap">
                <p class="font-bold text-gray-900 text-sm">{notif.title}</p>
                <Badge
                  label={notif.type === 'payment_success' ? 'Pembayaran' : notif.type === 'billing' ? 'Tagihan' : 'Info'}
                  variant={getTypeVariant(notif.type)}
                />
              </div>
              <p class="text-gray-600 text-sm mt-1 leading-relaxed">{notif.message}</p>
              <div class="flex items-center justify-between mt-3">
                <p class="text-xs text-gray-400">{formatDate(notif.created_at, true)}</p>
                {#if !notif.is_read}
                  <button
                    onclick={() => markRead(notif.id)}
                    class="flex items-center gap-1 text-xs text-green-600 hover:text-green-700 transition-colors font-semibold"
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
  {/if}
</div>
