<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount, onDestroy } from 'svelte';
  import { LogIn, ArrowLeft } from 'lucide-svelte';
  import { authApi, tokenStore } from '$lib/api';
  import { auth } from '$lib/stores/auth';
  import { Button, Input, Alert } from '$lib/components';
  import { toast } from '$lib/stores/toast';
  import logo from '$lib/assets/logo.png';

  let username = $state('');
  let password = $state('');
  let error    = $state('');
  let loading  = $state(false);
  let usernameError = $state('');
  let passwordError = $state('');

  onMount(() => {
    if ($auth.isAuthenticated) {
      const intended = sessionStorage.getItem('intended_url');
      goto(intended || '/');
    }

    // Bersihkan sisa token dari implementasi lama
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    }
  });

  // Bersihkan data form saat komponen di-unmount
  onDestroy(() => {
    username = '';
    password = '';
    error = '';
    usernameError = '';
    passwordError = '';
  });

  async function handleLogin(e: SubmitEvent) {
    e.preventDefault();
    error   = '';
    usernameError = '';
    passwordError = '';
    loading = false;

    let hasError = false;
    if (!username.trim()) {
      usernameError = 'Username atau NIS tidak boleh kosong.';
      hasError = true;
    }
    if (!password) {
      passwordError = 'Password tidak boleh kosong.';
      hasError = true;
    }

    if (hasError) return;

    loading = true;

    try {
      const res = await authApi.login(username, password);
      const session = res.data;

      // Token disimpan di MEMORI saja
      tokenStore.set(session.access_token);

      // Fetch profil user
      const profileRes = await authApi.getProfile();
      const user = profileRes.data;

      // Update memory-only auth store
      auth.login(user, session.access_token);

      // Set HttpOnly cookie via SvelteKit server endpoint
      // Cookie ini TIDAK bisa diakses oleh JavaScript (document.cookie)
      // Digunakan untuk re-auth saat page refresh
      await fetch('/api/auth/session', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token: session.access_token })
      });

      // Bersihkan password dari memori setelah login berhasil
      password = '';

      const roleRoutes: Record<string, string> = {
        super_user:  '/dashboard/super_user',
        guardian:    '/dashboard/guardian',
        pengasuh:    '/dashboard/pengasuh'
      };

      const role = user.role.name;
      const roleLabels: Record<string, string> = {
        super_user:  'Super User',
        guardian:    'Guardian',
        pengasuh:    'Pengasuh'
      };

      if (roleRoutes[role]) {
        let intended = sessionStorage.getItem('intended_url');
        sessionStorage.removeItem('intended_url');
        
        // Ensure intended url is valid, otherwise use default
        if (!intended || !intended.startsWith('/dashboard')) {
          intended = roleRoutes[role];
        }

        toast.success(`Selamat datang, ${roleLabels[role] ?? role}! Login berhasil.`);
        goto(intended);
      } else {
        throw new Error('Role tidak dikenali oleh sistem.');
      }
    } catch (e: any) {
      error = e.message || 'Username atau password salah.';
      // Bersihkan password jika login gagal
      password = '';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Masuk | Al-Anwar Payment</title>
  <meta name="description" content="Masuk ke portal administrasi keuangan Al-Anwar. Login menggunakan username atau NIS santri." />
</svelte:head>

<div class="relative min-h-screen bg-slate-50 flex items-center justify-center px-4 py-6 sm:py-12 overflow-hidden">
  <!-- Ambient Background Blobs -->
  <div class="absolute inset-0 pointer-events-none overflow-hidden bg-slate-50" aria-hidden="true">
    <div class="gradient-blob gradient-blob-1 opacity-25"></div>
    <div class="gradient-blob gradient-blob-2 opacity-20"></div>
  </div>

  <div class="relative w-full max-w-md animate-fade-in-up z-10">
    <!-- Back Link -->
    <div class="mb-4 sm:mb-6 flex justify-start">
      <a
        href="/"
        class="group inline-flex items-center gap-2 px-3 py-1.5 text-xs font-bold text-slate-500 hover:text-emerald-800 bg-white border border-slate-200/80 rounded-xl shadow-sm hover:shadow hover:border-emerald-600/30 transition-all duration-300"
      >
        <ArrowLeft size={14} class="transition-transform duration-300 group-hover:-translate-x-1 text-slate-400 group-hover:text-emerald-800" aria-hidden="true" />
        Kembali
      </a>
    </div>

    <!-- Brand Header -->
    <div class="text-center mb-6 sm:mb-8">
      <a href="/" aria-label="Kembali ke beranda">
        <div class="flex flex-col items-center gap-2 sm:gap-3">
          <img
            src={logo}
            alt="Logo Pondok Pesantren Al-Anwar"
            class="w-14 h-14 sm:w-20 sm:h-20 object-contain"
          />
          <div>
            <h1 class="text-lg sm:text-xl font-black text-slate-900 leading-tight">Pondok Pesantren Al-Anwar</h1>
            <p class="text-[10px] sm:text-xs text-slate-500 mt-1 max-w-xs mx-auto">Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan</p>
          </div>
        </div>
      </a>
    </div>

    <!-- Login Card -->
    <div class="bg-white/95 backdrop-blur-sm border border-slate-200/80 rounded-2xl shadow-sm hover:shadow-md transition-all duration-300 overflow-hidden">

      <!-- Top bar green -->
      <div class="h-1.5 bg-gradient-to-r from-emerald-600 to-emerald-700"></div>

      <div class="p-5 sm:p-8">
        <h2 class="text-base sm:text-lg font-bold text-slate-900 mb-1">Masuk ke Portal</h2>
        <p class="text-xs sm:text-sm text-slate-500 mb-4 sm:mb-6">Gunakan username atau NIS.</p>

        {#if error}
          <Alert type="error" message={error} class="mb-5" />
        {/if}

        <form onsubmit={handleLogin} class="space-y-3 sm:space-y-4" novalidate>
          <Input
            id="username"
            label="Username / NIS"
            type="text"
            placeholder="Masukkan username atau NIS"
            bind:value={username}
            required
            autocomplete="username"
            error={usernameError}
            oninput={() => usernameError = ''}
            class="text-xs sm:text-sm"
          />

          <Input
            id="password"
            label="Password"
            type="password"
            placeholder="Masukkan password"
            bind:value={password}
            required
            autocomplete="current-password"
            error={passwordError}
            oninput={() => passwordError = ''}
            class="text-xs sm:text-sm"
          />

          <Button type="submit" variant="primary" size="md" {loading} class="w-full sm:py-2.5 !mt-5 sm:!mt-6">
            {#snippet children()}
              <LogIn size={15} class="w-3.5 h-3.5 sm:w-4 sm:h-4" aria-hidden="true" />
              {loading ? 'Memproses...' : 'Masuk Sekarang'}
            {/snippet}
          </Button>
        </form>
      </div>
    </div>

    <p class="text-center text-[10px] sm:text-xs text-slate-400 mt-5 sm:mt-6">
      &copy; {new Date().getFullYear()} Al-Anwar Payment. Semua hak dilindungi.
    </p>
  </div>
</div>