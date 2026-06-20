<script lang="ts">
  import { goto } from '$app/navigation';
  import { ArrowRight, ChevronRight, Wallet, BellRing, LineChart, Users, CheckCircle } from 'lucide-svelte';
  import { Button } from '$lib/components';
  import logo from '$lib/assets/logo.png';
  import { onMount } from 'svelte';

  let mounted = $state(false);

  onMount(() => {
    const timer = setTimeout(() => {
      mounted = true;
    }, 50);
    return () => clearTimeout(timer);
  });

  const features = [
    {
      icon: Wallet,
      title: 'Pembayaran Transparan',
      description: 'Semua tagihan terekam jelas dan dapat diakses wali santri kapan saja, di mana saja.'
    },
    { 
      icon: BellRing,
      title: 'Notifikasi',
      description: 'Notifikasi otomatis untuk tagihan baru dan konfirmasi setelah pembayaran berhasil.'
    },
    {
      icon: LineChart,
      title: 'Laporan Instan',
      description: 'Dapat melihat laporan keuangan bulanan hingga tahunan.'
    },
    {
      icon: Users,
      title: 'Multi-Role',
      description: 'Akses terpisah dengan keamanan penuh.'
    }
  ];

  const highlights = [
    'Integrasi Midtrans Payment Gateway',
    'Tagihan otomatis setiap awal bulan',
    'Struk digital yang dapat dicetak',
    'Laporan keuangan'
  ];
</script>

<svelte:head>
  <title>Al-Anwar Payment | Sistem Administrasi Keuangan Pesantren</title>
  <meta
    name="description"
    content="Platform keuangan digital PONDOK PESANTREN PUTRA - PUTRI 'AL-ANWAR' Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan. Kelola tagihan, bayar online, dan pantau laporan keuangan secara real-time."
  />
</svelte:head>

<div class="min-h-screen bg-white">

  <!-- ─── Navbar ──────────────────────────────────────────── -->
  <header class="sticky top-0 z-50 bg-white/95 backdrop-blur-sm border-b border-slate-200/60">
    <nav class="max-w-7xl mx-auto px-4 sm:px-6 flex items-center justify-between h-16" aria-label="Navigasi utama">
      <a href="/" class="flex items-center gap-3 hover:opacity-80 transition-opacity">
        <img src={logo} alt="Logo Al-Anwar" class="w-9 h-9 object-contain" />
        <div class="hidden sm:block">
          <p class="text-sm font-bold text-slate-900 leading-tight">Al-Anwar Payment</p>
          <p class="text-xs text-slate-500">Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan</p>
        </div>
      </a>

      <!-- <div class="hidden md:flex items-center gap-6 text-sm font-medium text-slate-600">
        <a href="#fitur" class="hover:text-emerald-800 transition-colors duration-200">Fitur</a>
        <a href="#tentang" class="hover:text-emerald-800 transition-colors duration-200">Tentang</a>
      </div> -->

      <Button onclick={() => goto('/login')} variant="primary" size="md">
        {#snippet children()}
          Masuk
          <ArrowRight size={14} aria-hidden="true" />
        {/snippet}
      </Button>
    </nav>
  </header>

  <!-- ─── Hero ─────────────────────────────────────────────── -->
  <main>
    <section class="relative overflow-hidden py-20 md:py-28 px-4">
      <!-- Modern Animated Mesh Gradient Background -->
      <div class="absolute inset-0 pointer-events-none overflow-hidden bg-white" aria-hidden="true">
        <div class="gradient-blob gradient-blob-1"></div>
        <div class="gradient-blob gradient-blob-2"></div>
        <div class="gradient-blob gradient-blob-3"></div>
      </div>

      <div class="relative max-w-5xl mx-auto text-center">
        <!-- Logo besar -->
        <div class="flex justify-center mb-8 transition-all duration-[1200ms] cubic-bezier(0.16, 1, 0.3, 1) transform {mounted ? 'translate-y-0 opacity-100 scale-100' : 'translate-y-16 opacity-0 scale-90'}">
          <img src={logo} alt="Logo Pondok Pesantren Al-Anwar" class="w-40 h-40 object-contain" />
        </div>

        <h1 class="text-4xl md:text-6xl font-bold text-slate-900 mb-5 leading-[1.1] tracking-tight transition-all duration-[1200ms] delay-150 cubic-bezier(0.16, 1, 0.3, 1) transform {mounted ? 'translate-y-0 opacity-100' : 'translate-y-16 opacity-0'}">
          Administrasi Keuangan<br />
          <span class="gradient-text">Pesantren Digital</span>
        </h1>

        <p class="text-lg text-slate-500 mb-10 max-w-2xl mx-auto leading-relaxed transition-all duration-[1200ms] delay-300 cubic-bezier(0.16, 1, 0.3, 1) transform {mounted ? 'translate-y-0 opacity-100' : 'translate-y-16 opacity-0'}">
          Platform terintegrasi dalam satu ekosistem yang
          <strong class="text-slate-700">transparan</strong>,
          <strong class="text-slate-700">aman</strong>, dan
          <strong class="text-slate-700">efisien</strong>.
        </p>

        <!-- Highlights -->
        <div class="mt-12 w-full max-w-4xl mx-auto overflow-hidden relative py-4 mask-gradient transition-all duration-[1200ms] delay-[450ms] cubic-bezier(0.16, 1, 0.3, 1) transform {mounted ? 'translate-y-0 opacity-100' : 'translate-y-16 opacity-0'}">
          <div class="flex gap-8 animate-marquee whitespace-nowrap min-w-max">
            <!-- First set -->
            {#each highlights as item}
              <div class="flex items-center gap-2 text-sm font-medium text-slate-600 bg-slate-50 border border-slate-200/60 rounded-full px-5 py-2 shadow-sm hover:border-emerald-500 hover:text-emerald-800 transition-colors duration-300">
                <CheckCircle size={14} class="text-emerald-700 flex-shrink-0" aria-hidden="true" />
                {item}
              </div>
            {/each}
            <!-- Second set (for seamless loop) -->
            {#each highlights as item}
              <div class="flex items-center gap-2 text-sm font-medium text-slate-600 bg-slate-50 border border-slate-200/60 rounded-full px-5 py-2 shadow-sm hover:border-emerald-500 hover:text-emerald-800 transition-colors duration-300" aria-hidden="true">
                <CheckCircle size={14} class="text-emerald-700 flex-shrink-0" aria-hidden="true" />
                {item}
              </div>
            {/each}
            <!-- Third set (for seamless loop on wide screens) -->
            {#each highlights as item}
              <div class="flex items-center gap-2 text-sm font-medium text-slate-600 bg-slate-50 border border-slate-200/60 rounded-full px-5 py-2 shadow-sm hover:border-emerald-500 hover:text-emerald-800 transition-colors duration-300" aria-hidden="true">
                <CheckCircle size={14} class="text-emerald-700 flex-shrink-0" aria-hidden="true" />
                {item}
              </div>
            {/each}
          </div>
        </div>
      </div>
    </section>

    <!-- ─── Features ─────────────────────────────────────── -->
    <section id="#" class="py-20 px-4 bg-slate-50 border-t border-slate-100 scroll-mt-16">
      <div class="max-w-6xl mx-auto">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-slate-900 mb-3">Kenapa Al-Anwar Payment?</h2>
          <p class="text-slate-500 max-w-xl mx-auto">
            Dirancang khusus untuk kebutuhan administrasi lembaga pendidikan berbasis pesantren.
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          {#each features as feature, i}
            <div
              class="group bg-white border border-slate-200/80 rounded-xl p-6 hover:border-emerald-200 transition-colors duration-300 animate-fade-in-up"
              style="animation-delay: {i * 80}ms"
            >
              <div class="w-11 h-11 rounded-lg bg-emerald-50 border border-emerald-100 flex items-center justify-center mb-4 group-hover:bg-emerald-800 transition-colors duration-300">
                <feature.icon class="text-emerald-800 group-hover:text-white w-5 h-5 transition-colors duration-300" aria-hidden="true" />
              </div>
              <h3 class="text-slate-900 font-semibold text-sm mb-2">{feature.title}</h3>
              <p class="text-slate-500 text-sm leading-relaxed">{feature.description}</p>
            </div>
          {/each}
        </div>
      </div>
    </section>

    <!-- ─── Tentang ───────────────────────────────────────── -->
    <section id="#" class="py-20 px-4 bg-emerald-800 scroll-mt-16">
      <div class="max-w-3xl mx-auto text-center">
        <!-- <div class="flex justify-center mb-6">
          <img src={logo} alt="Logo Al-Anwar" class="w-14 h-14 object-contain" />
        </div> -->
        <h2 class="text-3xl font-bold text-white mb-3">Pondok Pesantren Al-Anwar</h2>
        <p class="text-emerald-100 leading-relaxed">
          Pondok Pesantren Putra - Putri Al-Anwar Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan — Platform ini dibangun untuk mempermudah pengelolaan administrasi keuangan santri secara digital, transparan, dan aman.
        </p>
        <p class="text-emerald-200/70 text-sm mt-3">Seluruh transaksi diproses melalui Midtrans sebagai payment gateway terpercaya.</p>

      </div>
    </section>
  </main>

  <!-- ─── Footer ──────────────────────────────────────────── -->
  <footer class="border-t border-slate-200/60 py-6 px-4 bg-white">
    <div class="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-3 text-sm text-slate-400">
      <div class="flex items-center gap-2">
        <img src={logo} alt="Logo Al-Anwar" class="w-5 h-5 object-contain" />
        <p>&copy; {new Date().getFullYear()} Al-Anwar Payment. Semua hak dilindungi.</p>
      </div>
      <p>Pondok Pesantren Putra - Putri Al-Anwar Dusun Kauman, Desa Selo, RT 05/RW 08 Kecamatan Tawangharjo Kabupaten Grobogan</p>
    </div>
  </footer>
</div>