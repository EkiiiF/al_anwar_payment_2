<script lang="ts">
  import { LayoutDashboard, Wallet, Bell, ClockIcon } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import SidebarLayout from '$lib/components/layout/SidebarLayout.svelte';
  import { auth } from '$lib/stores/auth';
  import { toast } from '$lib/stores/toast';

  let { children } = $props();

  onMount(() => {
    if (!$auth.loading && $auth.isAuthenticated) {
      if ($auth.user?.role.name !== 'guardian') {
        toast.error('Akses ditolak: Anda bukan Guardian.');
        const role = $auth.user?.role.name;
        if (role === 'super_user') goto('/dashboard/super_user');
        else if (role === 'pengasuh') goto('/dashboard/pengasuh');
        else goto('/login');
      }
    }
  });

  const userName = $derived([$auth.user?.user_profile?.first_name, $auth.user?.user_profile?.last_name].filter(Boolean).join(' ') || $auth.user?.username || '');
  const userEmail = $derived($auth.user?.user_profile?.email ?? '');

  const menuItems = [
    { name: 'Dashboard',     path: '/dashboard/guardian',               icon: LayoutDashboard },
    { name: 'Tagihan Saya',  path: '/dashboard/guardian/invoices',      icon: Wallet },
    { name: 'Riwayat Bayar', path: '/dashboard/guardian/history',       icon: ClockIcon }
  ];
</script>

<SidebarLayout
  role="guardian"
  roleLabel="Wali Santri (Guardian)"
  roleColor="blue"
  {menuItems}
  {userName}
  {userEmail}
  notifPath="/dashboard/guardian/notifications"
  notifEndpoint="/api/v1/guardian/notifications"
>
  {#snippet children()}
    {@render children()}
  {/snippet}
</SidebarLayout>
