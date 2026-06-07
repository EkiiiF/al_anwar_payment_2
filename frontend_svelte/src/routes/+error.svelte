<script lang="ts">
  import { page } from '$app/stores';
  import { Button } from '$lib/components';
  import logo from '$lib/assets/logo.png';
  import {
    FileQuestion, ShieldAlert, AlertTriangle,
    Lock, Wrench, ArrowLeft, Home, RefreshCw
  } from 'lucide-svelte';

  // Icon & deskripsi berdasarkan HTTP status
  function getErrorInfo(status: number) {
    switch (status) {
      case 401:
        return {
          icon:        Lock,
          color:       'amber',
          badge:       '401',
          title:       'Autentikasi Diperlukan',
          description: 'Anda perlu masuk untuk mengakses halaman ini.',
          action:      'login' as const
        };
      case 403:
        return {
          icon:        ShieldAlert,
          color:       'red',
          badge:       '403',
          title:       'Akses Ditolak',
          description: 'Anda tidak memiliki izin untuk mengakses halaman ini.',
          action:      'back' as const
        };
      case 404:
        return {
          icon:        FileQuestion,
          color:       'blue',
          badge:       '404',
          title:       'Halaman Tidak Ditemukan',
          description: 'Halaman yang Anda cari tidak ada atau telah dipindahkan.',
          action:      'home' as const
        };
      case 500:
        return {
          icon:        AlertTriangle,
          color:       'red',
          badge:       '500',
          title:       'Kesalahan Server Internal',
          description: 'Terjadi kesalahan pada server. Silakan coba lagi beberapa saat.',
          action:      'home' as const
        };
      case 503:
        return {
          icon:        Wrench,
          color:       'amber',
          badge:       '503',
          title:       'Sedang Dalam Pemeliharaan',
          description: 'Sistem sedang dalam pemeliharaan. Silakan coba beberapa saat lagi.',
          action:      'refresh' as const
        };
      default:
        return {
          icon:        AlertTriangle,
          color:       'red',
          badge:       String(status),
          title:       'Terjadi Kesalahan',
          description: 'Terjadi kesalahan yang tidak terduga. Silakan coba lagi.',
          action:      'home' as const
        };
    }
  }

  const info = $derived(getErrorInfo($page.status));

  const colorStyles = {
    amber: { bg: 'bg-amber-50', border: 'border-amber-200', icon: 'text-amber-600', badge: 'bg-amber-100 text-amber-700' },
    red:   { bg: 'bg-red-50',   border: 'border-red-200',   icon: 'text-red-600',   badge: 'bg-red-100 text-red-700' },
    blue:  { bg: 'bg-blue-50',  border: 'border-blue-200',  icon: 'text-blue-600',  badge: 'bg-blue-100 text-blue-700' },
    green: { bg: 'bg-green-50', border: 'border-green-200', icon: 'text-green-600', badge: 'bg-green-100 text-green-700' }
  };

  const cs = $derived(colorStyles[info.color as keyof typeof colorStyles] ?? colorStyles.red);
</script>

<svelte:head>
  <title>{$page.status} – {info.title} | Al-Anwar Payment</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 flex flex-col items-center justify-center p-4">
  <div class="w-full max-w-md bg-white border border-gray-200 rounded-2xl shadow-sm overflow-hidden">

    <div class="h-1.5 {info.color === 'amber' ? 'bg-amber-500' : info.color === 'blue' ? 'bg-blue-500' : 'bg-red-500'}"></div>

    <div class="p-8 text-center">
      <div class="flex justify-center mb-6">
        <img src={logo} alt="Logo Al-Anwar" class="w-14 h-14 object-contain" />
      </div>

      <span class="inline-flex px-3 py-1 rounded-full text-xs font-bold tracking-wider uppercase {cs.badge} mb-4">
        Error {info.badge}
      </span>

      <div class="flex justify-center mb-4">
        <div class="p-4 rounded-full {cs.bg} border {cs.border}">
          <info.icon size={28} class={cs.icon} aria-hidden="true" />
        </div>
      </div>

      <h1 class="text-xl font-bold text-gray-900 mb-2">{info.title}</h1>
      <p class="text-sm text-gray-500 leading-relaxed mb-6">
        {$page.error?.message || info.description}
      </p>

      <div class="flex justify-center gap-3">
        {#if info.action === 'login'}
          <Button href="/login" variant="primary">
            {#snippet children()}
              <Lock size={15} aria-hidden="true" />
              Masuk ke Sistem
            {/snippet}
          </Button>
        {:else if info.action === 'refresh'}
          <Button onclick={() => location.reload()} variant="primary">
            {#snippet children()}
              <RefreshCw size={15} aria-hidden="true" />
              Muat Ulang
            {/snippet}
          </Button>
        {:else if info.action === 'back'}
          <Button onclick={() => history.back()} variant="outline">
            {#snippet children()}
              <ArrowLeft size={15} aria-hidden="true" />
              Kembali
            {/snippet}
          </Button>
          <Button href="/" variant="primary">
            {#snippet children()}
              <Home size={15} aria-hidden="true" />
              Beranda
            {/snippet}
          </Button>
        {:else}
          <Button onclick={() => history.back()} variant="outline">
            {#snippet children()}
              <ArrowLeft size={15} aria-hidden="true" />
              Kembali
            {/snippet}
          </Button>
          <Button href="/" variant="primary">
            {#snippet children()}
              <Home size={15} aria-hidden="true" />
              Beranda
            {/snippet}
          </Button>
        {/if}
      </div>
    </div>

    <!-- Footer -->
    <div class="border-t border-gray-100 bg-gray-50 px-6 py-3">
      <p class="text-center text-xs text-gray-400">
        Pondok Pesantren Putra - Putri "Al-Anwar" Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan
      </p>
    </div>
  </div>
</div>
