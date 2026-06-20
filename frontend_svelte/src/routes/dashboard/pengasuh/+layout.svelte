<script lang="ts">
  import { LayoutDashboard, Users, BarChart3 } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import SidebarLayout from '$lib/components/layout/SidebarLayout.svelte';
  import { authApi } from '$lib/api';

  import { auth } from '$lib/stores/auth';
  import { toast } from '$lib/stores/toast';
  import { goto } from '$app/navigation';

  let { children } = $props();

  onMount(async () => {
    if (!$auth.loading && $auth.isAuthenticated) {
      if ($auth.user?.role.name !== 'pengasuh') {
        toast.error('Akses ditolak: Anda bukan Pengasuh.');
        const role = $auth.user?.role.name;
        if (role === 'super_user') goto('/dashboard/super_user');
        else if (role === 'guardian') goto('/dashboard/guardian');
        else goto('/login');
        return;
      }
    }
  });

  const userName = $derived([$auth.user?.first_name, $auth.user?.last_name].filter(Boolean).join(' ') || $auth.user?.username || '');
  const userEmail = $derived($auth.user?.email ?? '');

  const menuItems = [
    { name: 'Dashboard',   path: '/dashboard/pengasuh',          icon: LayoutDashboard },
    { name: 'Data Santri',  path: '/dashboard/pengasuh/students', icon: Users },
    { name: 'Monitoring',  path: '/dashboard/pengasuh/reports',  icon: BarChart3 }
  ];
</script>

<SidebarLayout
  role="pengasuh"
  roleLabel="Pengasuh"
  roleColor="teal"
  {menuItems}
  {userName}
  {userEmail}
>
  {#snippet children()}
    {@render children()}
  {/snippet}
</SidebarLayout>
