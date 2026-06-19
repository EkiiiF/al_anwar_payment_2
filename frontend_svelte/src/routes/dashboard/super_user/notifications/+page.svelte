<script lang="ts">
  import { onMount } from 'svelte';
  import { superUserApi } from '$lib/api';
  import { formatDate } from '$lib/utils';
  import { Spinner, Alert, Card, EmptyState, Badge } from '$lib/components';
  import type { Notification } from '$lib/types';
  import { Bell, Check } from 'lucide-svelte';

  let notifications = $state<Notification[]>([]);
  let loading       = $state(true);
  let error         = $state('');

  onMount(async () => {
    try {
      const res = await superUserApi.getNotifications();
      notifications = res.data ?? [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat notifikasi.';
    } finally {
      loading = false;
    }
  });

  async function markRead(id: string) {
    try {
      await superUserApi.markNotificationRead(id);
      notifications = notifications.map(n => n.id === id ? { ...n, is_read: true } : n);
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('notification-read'));
      }
    } catch { /* abaikan */ }
  }

  function getTypeVariant(title: string): 'success' | 'warning' | 'info' {
    if (title.includes('Berhasil') || title.includes('Masuk')) return 'success';
    if (title.includes('Gagal')) return 'warning';
    return 'info';
  }
</script>

<svelte:head>
  <title>Notifikasi Admin | Al-Anwar Payment</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-black text-gray-900 tracking-tight">Notifikasi Admin</h1>
      <p class="text-gray-500 text-sm mt-1">Laporan pembayaran santri dan sistem.</p>
    </div>
    {#if notifications.some(n => !n.is_read)}
      <Badge label="{notifications.filter(n => !n.is_read).length} Baru" variant="warning" dot />
    {/if}
  </div>

  {#if error}
    <Alert type="error" message={error} />
  {:else if loading}
    <div class="flex justify-center py-12">
      <Spinner size="lg" />
    </div>
  {:else if notifications.length === 0}
    <Card>
      <EmptyState
        title="Belum ada notifikasi"
        description="Semua aktivitas pembayaran santri akan muncul di sini."
      />
    </Card>
  {:else}
    <div class="space-y-3">
      {#each notifications as notif (notif.id)}
        <div
          class="group relative p-4 rounded-xl border transition-all duration-200
            {notif.is_read ? 'bg-white border-slate-200/60 shadow-sm' : 'bg-green-50/50 border-green-200/60 shadow-sm'}"
        >
          {#if !notif.is_read}
            <div class="absolute top-4 right-4 w-2 h-2 rounded-full bg-green-600"></div>
          {/if}

          <div class="flex items-start gap-3.5 pr-6">
            <div class="p-2 rounded-lg flex-shrink-0 border
              {getTypeVariant(notif.title) === 'success' ? 'bg-emerald-50 border-emerald-200 text-emerald-700' :
               getTypeVariant(notif.title) === 'warning' ? 'bg-amber-50 border-amber-200 text-amber-700' :
               'bg-blue-50 border-blue-200 text-blue-700'}">
              <Bell size={16} aria-hidden="true" />
            </div>

            <div class="min-w-0 flex-1">
              <div class="flex items-start justify-between gap-2 flex-wrap mb-1">
                <p class="font-bold text-slate-900 text-base">{notif.title}</p>
                <Badge
                  label={getTypeVariant(notif.title) === 'success' ? 'Pembayaran' : 'Sistem'}
                  variant={getTypeVariant(notif.title)}
                  class="text-xs py-0.5"
                />
              </div>
              <p class="text-slate-600 text-sm mt-1 leading-relaxed">{notif.message}</p>
              
              <div class="flex items-center justify-between mt-3.5 pt-2 border-t border-slate-100">
                <p class="text-xs text-slate-400">
                  {formatDate(notif.created_at, true)}
                </p>
                
                {#if !notif.is_read}
                  <button
                    onclick={() => markRead(notif.id)}
                    class="flex items-center gap-1 text-xs text-green-700 hover:text-green-800 transition-colors font-bold cursor-pointer"
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
