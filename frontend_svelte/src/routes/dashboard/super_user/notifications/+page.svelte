<script lang="ts">
  import { onMount } from 'svelte';
  import { superUserApi } from '$lib/api';
  import { formatDate } from '$lib/utils';
  import { Spinner, Alert, EmptyState, Badge } from '$lib/components';
  import type { Notification } from '$lib/types';
  import { Bell } from 'lucide-svelte';

  let notifications = $state<Notification[]>([]);
  let loading       = $state(true);
  let error         = $state('');
  let activeNotif   = $state<Notification | null>(null);

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

  async function handleSelectNotif(notif: Notification) {
    if (!notif.is_read) {
      await markRead(notif.id);
    }
    activeNotif = notif;
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
  {#if !activeNotif}
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
      <EmptyState
        title="Belum ada notifikasi"
        description="Semua aktivitas pembayaran santri akan muncul di sini."
      />
    {:else}
      <div class="divide-y divide-slate-100/80">
        {#each notifications as notif (notif.id)}
          <button
            type="button"
            onclick={() => handleSelectNotif(notif)}
            class="w-full text-left py-3 px-1.5 hover:bg-slate-50/60 transition-colors flex items-start gap-3 sm:gap-4 relative group cursor-pointer focus:outline-none"
          >
            {#if !notif.is_read}
              <div class="absolute top-1/2 -translate-y-1/2 left-0 w-1.5 h-1.5 rounded-full bg-emerald-600"></div>
            {/if}
            
            <div class="pl-3.5 flex items-start gap-3 sm:gap-4 w-full min-w-0">
              <div class="p-2 rounded-lg flex-shrink-0
                {getTypeVariant(notif.title) === 'success' ? 'bg-emerald-50 text-emerald-700' :
                 getTypeVariant(notif.title) === 'warning' ? 'bg-amber-50 text-amber-700' :
                 'bg-blue-50 text-blue-700'}">
                <Bell size={16} aria-hidden="true" />
              </div>

              <div class="min-w-0 flex-1">
                <div class="flex items-center justify-between gap-2 mb-0.5">
                  <p class="font-bold text-slate-800 text-sm sm:text-base truncate {notif.is_read ? 'font-semibold' : 'font-extrabold'}">
                    {notif.title}
                  </p>
                  <span class="text-[10px] text-slate-400 flex-shrink-0">{formatDate(notif.created_at, false)}</span>
                </div>
                <p class="text-slate-500 text-xs sm:text-sm leading-relaxed truncate sm:line-clamp-2 sm:whitespace-normal">
                  {notif.message}
                </p>
              </div>
            </div>
          </button>
        {/each}
      </div>
    {/if}
  {:else}
    <div class="space-y-6 max-w-2xl">
      <button
        type="button"
        onclick={() => activeNotif = null}
        class="inline-flex items-center gap-2 text-sm font-bold text-slate-600 hover:text-slate-900 transition-colors cursor-pointer"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        <span>Kembali ke Daftar Notifikasi</span>
      </button>

      <div class="py-4 space-y-4">
        <div class="flex items-center gap-3">
          <div class="p-3 rounded-xl border
            {getTypeVariant(activeNotif.title) === 'success' ? 'bg-emerald-50 border-emerald-100 text-emerald-700' :
             getTypeVariant(activeNotif.title) === 'warning' ? 'bg-amber-50 border-amber-100 text-amber-700' :
             'bg-blue-50 border-blue-100 text-blue-700'}">
            <Bell size={24} aria-hidden="true" />
          </div>
          <div>
            <h2 class="text-xl font-extrabold text-slate-900">{activeNotif.title}</h2>
            <p class="text-xs text-slate-400 mt-1">{formatDate(activeNotif.created_at, true)}</p>
          </div>
        </div>

        <div class="pt-4 border-t border-slate-100">
          <p class="text-slate-700 text-base leading-relaxed whitespace-pre-wrap">{activeNotif.message}</p>
        </div>
      </div>
    </div>
  {/if}
</div>
