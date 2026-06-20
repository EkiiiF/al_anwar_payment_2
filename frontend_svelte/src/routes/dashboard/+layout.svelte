<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { tokenStore, authApi } from '$lib/api';
  import { auth } from '$lib/stores/auth';
  import { secureAuth } from '$lib/stores/secure-auth';
  import { initAuthContext } from '$lib/context/auth-context';
  import { toast } from '$lib/stores/toast';
  import Spinner from '$lib/components/ui/Spinner.svelte';

  let { data, children } = $props();
  const { actions, cleanup } = initAuthContext();

  onDestroy(() => {
    cleanup();
    tokenStore.clear();
    secureAuth.reset();
    if (import.meta.env.DEV) {
      secureAuth.audit();
    }
  });

  onMount(async () => {
    let token = tokenStore.get();
    if (!token && data.hasSession) {
      try {
        const res = await fetch('/api/auth/session');
        if (res.ok) {
          const session = await res.json();
            token = session.token;
            tokenStore.set(token);
          }
      } catch (e) {
        console.error('Gagal restorasi sesi:', e);
      }
    }
    
    if (!token) {
      const intended = $page.url.pathname;
      sessionStorage.setItem('intended_url', intended);
      toast.warning('Anda harus login terlebih dahulu.');
      goto('/login');
      return;
    }

    // Verifikasi token ke backend
    try {
      auth.setLoading(true);
      const res = await authApi.getProfile();
      auth.setUser(res.data);

      actions.setUser(res.data);
      actions.setToken(token);
    } catch (err) {
      // Token invalid — bersihkan semua
      cleanup();
      tokenStore.clear();
      secureAuth.reset();
      toast.error('Sesi tidak valid. Silakan login kembali.');
      goto('/login');
    }
  });
</script>

{#if !$auth.loading && $auth.isAuthenticated}
  {@render children()}
{:else}
  <div class="min-h-screen bg-white flex flex-col items-center justify-center p-4">
    <div class="text-center space-y-4 max-w-xs">
      <Spinner size="lg" />
      <div class="space-y-1">
        <p class="text-lg font-bold text-slate-900 tracking-tight">Memverifikasi Sesi</p>
        <p class="text-xs text-slate-500 leading-relaxed">
          Mohon tunggu sebentar, kami sedang memastikan keamanan akun Anda...
        </p>
      </div>
    </div>
  </div>
{/if}
