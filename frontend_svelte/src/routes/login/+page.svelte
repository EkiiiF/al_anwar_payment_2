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

<div class="min-h-screen bg-gray-50 flex items-center justify-center px-4 py-12">
  <div class="w-full max-w-md animate-fade-in-up">
    <!-- Back Link -->
    <div class="mb-6 flex justify-start">
      <a
        href="/"
        class="inline-flex items-center gap-2 px-3 py-1.5 text-sm font-semibold text-gray-500 hover:text-green-700 bg-white border border-gray-200 rounded-xl shadow-sm hover:shadow transition-all duration-200"
      >
        <ArrowLeft size={16} aria-hidden="true" />
        Kembali ke Beranda
      </a>
    </div>

    <!-- Brand Header -->
    <div class="text-center mb-8">
      <a href="/" aria-label="Kembali ke beranda">
        <div class="flex flex-col items-center gap-3">
          <img
            src={logo}
            alt="Logo Pondok Pesantren Al-Anwar"
            class="w-20 h-20 object-contain"
          />
          <div>
            <h1 class="text-xl font-black text-gray-900">Pondok Pesantren Al-Anwar</h1>
            <p class="text-sm text-gray-500">Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan</p>
          </div>
        </div>
      </a>
    </div>

    <!-- Login Card -->
    <div class="bg-white border border-gray-200 rounded-2xl shadow-sm overflow-hidden">

      <!-- Top bar green -->
      <div class="h-1.5 bg-gradient-to-r from-green-500 to-green-600"></div>

      <div class="p-8">
        <h2 class="text-lg font-bold text-gray-900 mb-1">Masuk ke Portal</h2>
        <p class="text-sm text-gray-500 mb-6">Gunakan username atau NIS.</p>

        {#if error}
          <Alert type="error" message={error} class="mb-5" />
        {/if}

        <form onsubmit={handleLogin} class="space-y-4" novalidate>
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
          />

          <Button type="submit" variant="primary" size="lg" {loading} class="w-full !mt-6">
            {#snippet children()}
              <LogIn size={16} aria-hidden="true" />
              {loading ? 'Memproses...' : 'Masuk Sekarang'}
            {/snippet}
          </Button>
        </form>
      </div>
    </div>

    <p class="text-center text-xs text-gray-400 mt-6">
      &copy; {new Date().getFullYear()} Al-Anwar Payment. Semua hak dilindungi.
    </p>
  </div>
</div>