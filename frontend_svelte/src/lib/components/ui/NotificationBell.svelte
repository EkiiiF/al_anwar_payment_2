<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Bell } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { apiClient } from '$lib/api';
  import type { Notification } from '$lib/types';

  let {
    notifPath = '/dashboard',
    fetchEndpoint = ''
  }: {
    notifPath?: string;
    fetchEndpoint?: string;
  } = $props();

  let unreadCount = $state(0);
  let intervalId: ReturnType<typeof setInterval>;

  async function fetchUnread() {
    if (!fetchEndpoint) return;
    try {
      const res = await apiClient.get<Notification[]>(fetchEndpoint);
      const notifications = res.data ?? [];
      unreadCount = notifications.filter(n => !n.is_read).length;
    } catch {
    }
  }

  onMount(() => {
    fetchUnread();
    intervalId = setInterval(fetchUnread, 30_000);
    if (typeof window !== 'undefined') {
      window.addEventListener('notification-read', fetchUnread);
    }
  });

  onDestroy(() => {
    clearInterval(intervalId);
    if (typeof window !== 'undefined') {
      window.removeEventListener('notification-read', fetchUnread);
    }
  });

  const badgeLabel = $derived(
    unreadCount >= 100 ? '99+' : String(unreadCount)
  );
</script>

<button
  onclick={() => goto(notifPath)}
  aria-label="Notifikasi{unreadCount > 0 ? `, ${unreadCount} belum dibaca` : ''}"
  title="Lihat Notifikasi"
  class="relative p-2 rounded-xl text-gray-500 hover:text-gray-900 hover:bg-gray-100 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-green-500/40"
>
  <Bell size={20} aria-hidden="true" />

  {#if unreadCount > 0}
    <span
      class="absolute -top-0.5 -right-0.5 min-w-[18px] h-[18px] px-1 flex items-center justify-center
             rounded-full bg-red-500 text-white text-[10px] font-bold leading-none
             ring-2 ring-white animate-pulse"
      aria-hidden="true"
    >
      {badgeLabel}
    </span>
  {/if}
</button>
