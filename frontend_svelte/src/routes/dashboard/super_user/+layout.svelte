<script lang="ts">
  import { LayoutDashboard, Users, FileText, Settings, BarChart3, ClipboardList, Bell } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import SidebarLayout from '$lib/components/layout/SidebarLayout.svelte';
  import { auth } from '$lib/stores/auth';
  import { toast } from '$lib/stores/toast';

  let { children } = $props();

  onMount(() => {
    if (!$auth.loading && $auth.isAuthenticated) {
      const role = $auth.user?.role.name;
      if (role !== 'super_user') {
        toast.error('Akses ditolak: Anda bukan Super User.');
        if (role === 'guardian') goto('/dashboard/guardian');
        else if (role === 'pengasuh') goto('/dashboard/pengasuh');
        else goto('/login');
      }
    }
  });

  const userName = $derived([$auth.user?.first_name, $auth.user?.last_name].filter(Boolean).join(' ') || $auth.user?.username || '');
  const userEmail = $derived($auth.user?.email ?? '');

  const menuItems = [
    { name: 'Dashboard',     path: '/dashboard/super_user',               icon: LayoutDashboard },
    { name: 'Data Santri',   path: '/dashboard/super_user/students',      icon: Users },
    { name: 'Tagihan',       path: '/dashboard/super_user/billing',       icon: FileText },
    { name: 'Master Data',   path: '/dashboard/super_user/master',        icon: Settings },
    { name: 'Laporan',       path: '/dashboard/super_user/reports',       icon: BarChart3 },
    { name: 'Log Aktivitas', path: '/dashboard/super_user/logs',          icon: ClipboardList }
  ];
</script>

<SidebarLayout
  role="super_user"
  roleLabel="Super User (Bendahara)"
  roleColor="emerald"
  {menuItems}
  {userName}
  {userEmail}
  notifPath="/dashboard/super_user/notifications"
  notifEndpoint="/api/v1/admin/notifications"
>
  {#snippet children()}
    {@render children()}
  {/snippet}
</SidebarLayout>
