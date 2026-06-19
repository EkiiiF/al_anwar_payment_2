<script lang="ts">
  import type { Snippet } from 'svelte';
  import type { IconComponent } from '$lib/types/icons';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { LogOut, Menu, X, UserCircle } from 'lucide-svelte';
  import NotificationBell from '$lib/components/ui/NotificationBell.svelte';
  import { authApi, tokenStore } from '$lib/api';
  import { getInitials } from '$lib/utils';
  import { toast } from '$lib/stores/toast';
  import ConfirmDialog from '$lib/components/ui/ConfirmDialog.svelte';
  import logo from '$lib/assets/logo.png';

  let {
    role,
    roleLabel,
    roleColor = 'emerald',
    menuItems,
    userName = '',
    userEmail = '',
    notifPath = '',
    notifEndpoint = '',
    children
  }: {
    role: string;
    roleLabel: string;
    roleColor?: 'emerald' | 'teal' | 'blue';
    menuItems: { name: string; path: string; icon: IconComponent }[];
    userName?: string;
    userEmail?: string;
    notifPath?: string;
    notifEndpoint?: string;
    children: Snippet;
  } = $props();

  let isMobileOpen    = $state(false);
  let showLogoutConfirm = $state(false);
  let loggingOut      = $state(false);

  const roleThemeMap: Record<string, { badge: string; avatar: string; dot: string; activeBg: string }> = {
    super_user: { badge: 'bg-green-50 text-green-700',   avatar: 'bg-green-700',  dot: 'bg-green-600',  activeBg: 'bg-emerald-800 text-white' },
    guardian:   { badge: 'bg-blue-50 text-blue-700',     avatar: 'bg-blue-700',   dot: 'bg-blue-600',   activeBg: 'bg-blue-700 text-white' },
    pengasuh:   { badge: 'bg-purple-50 text-purple-800', avatar: 'bg-purple-700', dot: 'bg-purple-600', activeBg: 'bg-purple-800 text-white' }
  };
  const rc = $derived(roleThemeMap[role] || roleThemeMap.super_user);

  const profilePath = $derived(`/dashboard/${role}/profile`);

  async function handleLogout() {
    loggingOut = true;
    try { await authApi.logout(); } catch { /* ignore */ }
    try {
      await fetch('/api/auth/session', { method: 'DELETE' });
    } catch { /* ignore */ }
    tokenStore.clear();
    showLogoutConfirm = false;
    toast.info('Anda telah berhasil keluar. Sampai jumpa!');
    await new Promise(r => setTimeout(r, 800));
    goto('/login');
  }

  function openLogoutConfirm() {
    isMobileOpen = false;
    showLogoutConfirm = true;
  }

  function isActive(path: string): boolean {
    return $page.url.pathname === path;
  }

  function closeMobile() { isMobileOpen = false; }
</script>

<div class="h-screen bg-slate-50 flex overflow-hidden">

  <!-- ─── Desktop Sidebar ──────────────────────────────────── -->
  <aside
    class="w-60 h-screen bg-white border-r border-slate-200 hidden md:flex flex-col flex-shrink-0 sticky top-0"
    aria-label="Sidebar navigasi"
  >
    <!-- Brand -->
    <div class="px-5 py-4 border-b border-slate-100">
      <div class="flex items-center gap-3">
        <img src={logo} alt="Logo Al-Anwar" class="w-9 h-9 rounded-lg object-contain flex-shrink-0" />
        <div class="min-w-0">
          <p class="text-sm font-bold text-slate-900 leading-tight">Al-Anwar</p>
          <p class="text-xs text-slate-500 truncate">Pondok Pesantren</p>
        </div>
      </div>
      <div class="mt-3">
        <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium {rc.badge}">
          <span class="w-1.5 h-1.5 rounded-full {rc.dot}"></span>
          {roleLabel}
        </span>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-3 py-3 space-y-0.5 overflow-y-auto" aria-label="Menu {roleLabel}">
      {#each menuItems as item}
        <a
          href={item.path}
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors duration-200
            {isActive(item.path)
              ? rc.activeBg
              : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'}"
          aria-current={isActive(item.path) ? 'page' : undefined}
        >
          <item.icon size={17} aria-hidden="true" />
          {item.name}
        </a>
      {/each}
    </nav>

    <!-- Profile + Logout -->
    <div class="border-t border-slate-100 px-3 py-3 space-y-0.5">
      <!-- Profile Link -->
      <a
        href={profilePath}
        class="flex items-center gap-3 px-3 py-2 rounded-lg transition-colors duration-200
          {isActive(profilePath)
            ? rc.activeBg
            : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'}"
      >
        {#if userName}
          <div class="w-7 h-7 rounded-full {rc.avatar} flex items-center justify-center text-white text-xs font-semibold flex-shrink-0">
            {getInitials(userName)}
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-sm font-semibold leading-tight truncate {isActive(profilePath) ? 'text-white' : 'text-slate-900'}">{userName}</p>
            {#if userEmail}
              <p class="text-xs truncate {isActive(profilePath) ? 'text-white/80' : 'text-slate-400'}">{userEmail}</p>
            {/if}
          </div>
        {:else}
          <UserCircle size={17} aria-hidden="true" />
          <span class="text-sm font-medium">Profil Saya</span>
        {/if}
      </a>

      <!-- Logout -->
      <button
        onclick={openLogoutConfirm}
        class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium text-red-600 hover:bg-red-50 transition-colors duration-200"
      >
        <LogOut size={17} aria-hidden="true" />
        Keluar Sistem
      </button>
    </div>
  </aside>

  <!-- ─── Mobile Header ─── -->
  <header class="md:hidden fixed top-0 inset-x-0 bg-white border-b border-slate-200 z-40">
    <div class="flex items-center justify-between px-4 h-14">
      <!-- Hamburger KIRI -->
      <button
        onclick={() => isMobileOpen = !isMobileOpen}
        class="p-2 -ml-1 rounded-lg text-slate-500 hover:bg-slate-50 transition-colors duration-200"
        aria-label={isMobileOpen ? 'Tutup menu' : 'Buka menu'}
        aria-expanded={isMobileOpen}
        aria-controls="mobile-drawer"
      >
        {#if isMobileOpen}
          <X size={22} aria-hidden="true" />
        {:else}
          <Menu size={22} aria-hidden="true" />
        {/if}
      </button>

      <!-- Judul di tengah -->
      <p class="text-sm font-semibold text-slate-900">{roleLabel}</p>

      <!-- Notif Bell KANAN -->
      <div class="w-10 flex justify-end">
        {#if notifPath}
          <NotificationBell notifPath={notifPath} fetchEndpoint={notifEndpoint} />
        {/if}
      </div>
    </div>
  </header>

  <!-- ─── Mobile Drawer (slide dari KIRI) ─── -->
  {#if isMobileOpen}
    <!-- Backdrop -->
    <div
      class="md:hidden fixed inset-0 bg-black/20 z-30 animate-fade-in"
      onclick={closeMobile}
      aria-hidden="true"
    ></div>

    <!-- Drawer panel -->
    <aside
      id="mobile-drawer"
      class="md:hidden fixed top-0 left-0 h-full w-72 bg-white border-r border-slate-200 z-40 flex flex-col animate-slide-left"
    >
      <!-- Drawer header -->
      <div class="flex items-center gap-3 px-4 py-4 border-b border-slate-100">
        <img src={logo} alt="Logo Al-Anwar" class="w-9 h-9 object-contain" />
        <div class="min-w-0 flex-1">
          <p class="text-sm font-bold text-slate-900 leading-tight">Al-Anwar</p>
          <p class="text-xs text-slate-500">Pondok Pesantren</p>
        </div>
        <button
          onclick={closeMobile}
          class="p-1.5 rounded-lg text-slate-400 hover:bg-slate-50 transition-colors duration-200 flex-shrink-0"
          aria-label="Tutup menu"
        >
          <X size={18} aria-hidden="true" />
        </button>
      </div>

      <!-- Role badge -->
      <div class="px-4 py-2 border-b border-slate-100">
        <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium {rc.badge}">
          <span class="w-1.5 h-1.5 rounded-full {rc.dot}"></span>
          {roleLabel}
        </span>
      </div>

      <!-- Nav items -->
      <nav class="flex-1 p-3 space-y-0.5 overflow-y-auto">
        {#each menuItems as item}
          <a
            href={item.path}
            onclick={closeMobile}
            class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors duration-200
              {isActive(item.path) ? rc.activeBg : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'}"
            aria-current={isActive(item.path) ? 'page' : undefined}
          >
            <item.icon size={18} aria-hidden="true" />
            {item.name}
          </a>
        {/each}
      </nav>

      <!-- Profile + Logout di Drawer -->
      <div class="border-t border-slate-100 p-3 space-y-0.5">
        <a
          href={profilePath}
          onclick={closeMobile}
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors duration-200
            {isActive(profilePath) ? rc.activeBg : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'}"
        >
          {#if userName}
            <div class="w-8 h-8 rounded-full {rc.avatar} flex items-center justify-center text-white text-xs font-semibold flex-shrink-0">
              {getInitials(userName)}
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-sm font-semibold leading-tight truncate {isActive(profilePath) ? 'text-white' : 'text-slate-900'}">{userName}</p>
              {#if userEmail}
                <p class="text-xs truncate {isActive(profilePath) ? 'text-white/80' : 'text-slate-400'}">{userEmail}</p>
              {/if}
            </div>
          {:else}
            <UserCircle size={18} aria-hidden="true" />
            <span class="text-sm font-medium">Profil Saya</span>
          {/if}
        </a>

        <button
          onclick={openLogoutConfirm}
          class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-red-600 hover:bg-red-50 transition-colors duration-200"
        >
          <LogOut size={17} aria-hidden="true" />
          Keluar Sistem
        </button>
      </div>
    </aside>
  {/if}

  <!-- ─── Main Content ──────────────────────────────────────── -->
  <main class="flex-1 h-full flex flex-col overflow-hidden pt-14 md:pt-0" id="main-content">
    <!-- Top bar dengan notification bell (khusus Desktop) -->
    {#if notifPath}
      <div class="hidden md:flex justify-end px-8 pt-6 pb-0 flex-shrink-0">
        <NotificationBell notifPath={notifPath} fetchEndpoint={notifEndpoint} />
      </div>
    {/if}
    <div class="p-5 md:px-8 md:py-6 flex-1 overflow-y-auto min-h-0">
      {@render children()}
    </div>
  </main>
</div>

<!-- ─── Logout Confirmation Dialog ─────────────────────── -->
<ConfirmDialog
  bind:open={showLogoutConfirm}
  title="Konfirmasi Keluar"
  message="Apakah Anda yakin ingin keluar dari sistem? Anda perlu login kembali untuk mengakses dashboard."
  confirmText="Ya, Keluar"
  cancelText="Batal"
  variant="logout"
  loading={loggingOut}
  onConfirm={handleLogout}
/>
