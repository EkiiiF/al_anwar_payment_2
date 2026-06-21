<script lang="ts">
  import { Calendar, Settings } from 'lucide-svelte';

  let {
    autoBillingEnabled = $bindable(false),
    autoSemesterBillingEnabled = $bindable(false),
    billingHijriDay = $bindable('1'),
    midtransEnvironment = $bindable('sandbox'),
    settingLoading = false,
    onToggle,
    onDayChange,
    onMidtransToggle
  }: {
    autoBillingEnabled: boolean;
    autoSemesterBillingEnabled: boolean;
    billingHijriDay: string;
    midtransEnvironment: string;
    settingLoading: boolean;
    onToggle: (key: string) => void;
    onDayChange: () => void;
    onMidtransToggle: () => void;
  } = $props();
</script>

<div class="bg-white border border-slate-200 rounded-2xl shadow-sm p-6 hover:shadow-md transition-all duration-200">
  <div class="flex items-start gap-4 mb-5">
    <div class="p-3 rounded-xl bg-blue-50 border border-blue-100 flex-shrink-0">
      <Settings size={22} class="text-blue-700" />
    </div>
    <div class="flex-1 min-w-0">
      <h2 class="font-extrabold text-slate-900 text-lg">Pengaturan Sistem & Gateway</h2>
      <p class="text-xs text-slate-500 mt-0.5">Atur jadwal billing otomatis dan status lingkungan pembayaran Midtrans.</p>
    </div>
  </div>

  <div class="space-y-4">
    <!-- Auto Billing Bulanan -->
    <div class="rounded-2xl p-4 border transition-all duration-300 shadow-sm hover:shadow-md {autoBillingEnabled ? 'bg-emerald-50/30 border-emerald-200/80' : 'bg-slate-50/50 border-slate-200'}">
      <div class="flex items-center justify-between gap-4 mb-4">
        <div class="flex-1 min-w-0">
          <p class="text-sm font-bold {autoBillingEnabled ? 'text-emerald-900' : 'text-slate-900'}">Auto Billing Syahriyyah</p>
          <p class="text-xs {autoBillingEnabled ? 'text-emerald-700/80' : 'text-slate-500'} mt-0.5 font-medium">Tagihan bulanan dibuat otomatis tiap awal bulan Hijriah.</p>
        </div>
        <button 
          type="button"
          onclick={() => onToggle('auto_generate_billing')}
          disabled={settingLoading}
          aria-label="Toggle Auto Billing Syahriyyah"
          title="Toggle Auto Billing Syahriyyah"
          class="relative inline-flex h-6.5 w-12 flex-shrink-0 items-center rounded-full transition-all focus:outline-none focus:ring-2 focus:ring-emerald-600/20 focus:ring-offset-2
                 {autoBillingEnabled ? 'bg-emerald-600 shadow-inner' : 'bg-slate-300'}"
        >
          <span class="inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform duration-200
                 {autoBillingEnabled ? 'translate-x-6' : 'translate-x-1'}"></span>
        </button>
      </div>

      <div class="bg-white border {autoBillingEnabled ? 'border-emerald-100' : 'border-slate-200/80'} rounded-xl p-3.5 space-y-3 shadow-inner">
        <div class="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
          <div class="flex flex-col gap-1 w-full md:w-1/2">
            <label for="billing-hijri-day-select" class="text-[10px] font-bold {autoBillingEnabled ? 'text-emerald-700/80' : 'text-slate-500'} uppercase tracking-wider">Hari Penerbitan Tagihan</label>
            <select
              id="billing-hijri-day-select"
              bind:value={billingHijriDay}
              onchange={onDayChange}
              disabled={settingLoading || !autoBillingEnabled}
              class="w-full px-3 py-2 rounded-lg bg-white border {autoBillingEnabled ? 'border-emerald-300 focus:ring-emerald-600/20 focus:border-emerald-600' : 'border-slate-300 focus:ring-emerald-800/20 focus:border-emerald-800'} text-sm font-semibold outline-none disabled:opacity-50 disabled:bg-slate-50 transition-all cursor-pointer"
            >
              {#each Array.from({ length: 30 }, (_, i) => String(i + 1)) as day}
                <option value={day}>Tanggal {day} Hijriah</option>
              {/each}
            </select>
          </div>
        </div>
      </div>
    </div>

    <!-- Auto Billing Semester -->
    <div class="rounded-2xl p-4 border transition-all duration-300 shadow-sm hover:shadow-md flex items-center justify-between gap-4 {autoSemesterBillingEnabled ? 'bg-blue-50/30 border-blue-200/80' : 'bg-slate-50/50 border-slate-200'}">
      <div class="flex-1 min-w-0">
        <p class="text-sm font-bold {autoSemesterBillingEnabled ? 'text-blue-900' : 'text-slate-900'}">Auto Billing Semester</p>
        <p class="text-xs {autoSemesterBillingEnabled ? 'text-blue-700/80' : 'text-slate-500'} mt-0.5 font-medium">Tagihan semester (ngaji/sekolah) diterbitkan otomatis di bulan Syawal & Rabi'ul Akhir.</p>
      </div>
      <button 
        type="button"
        onclick={() => onToggle('auto_generate_semester_billing')}
        disabled={settingLoading}
        aria-label="Toggle Auto Billing Semester"
        title="Toggle Auto Billing Semester"
        class="relative inline-flex h-6.5 w-12 flex-shrink-0 items-center rounded-full transition-all focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:ring-offset-2
               {autoSemesterBillingEnabled ? 'bg-blue-600 shadow-inner' : 'bg-slate-300'}"
      >
        <span class="inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform duration-200
               {autoSemesterBillingEnabled ? 'translate-x-6' : 'translate-x-1'}"></span>
      </button>
    </div>

    <!-- Midtrans Environment Toggle -->
    <div class="rounded-2xl p-4 border transition-all duration-300 shadow-sm hover:shadow-md flex items-center justify-between gap-4 {midtransEnvironment === 'production' ? 'bg-amber-50/30 border-amber-200/80' : 'bg-slate-50/50 border-slate-200'}">
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2">
          <p class="text-sm font-bold {midtransEnvironment === 'production' ? 'text-amber-900' : 'text-slate-900'}">Mode Midtrans Gateway</p>
          <span class="text-[9px] font-extrabold px-1.5 py-0.5 rounded uppercase border {midtransEnvironment === 'production' ? 'bg-green-50 text-green-700 border-green-200' : 'bg-blue-50 text-blue-700 border-blue-200'}">
            {midtransEnvironment === 'production' ? 'Live' : 'Sandbox'}
          </span>
        </div>
        <p class="text-xs {midtransEnvironment === 'production' ? 'text-amber-700/80' : 'text-slate-500'} mt-0.5 font-medium">
          {midtransEnvironment === 'production' 
            ? 'Menggunakan transaksi uang asli. Pastikan Server & Client Key Production sudah disetel.' 
            : 'Menggunakan transaksi simulasi uang mainan untuk kebutuhan uji coba.'}
        </p>
      </div>
      <button 
        type="button"
        onclick={onMidtransToggle}
        disabled={settingLoading}
        aria-label="Toggle Midtrans Environment"
        title="Toggle Midtrans Environment"
        class="relative inline-flex h-6.5 w-12 flex-shrink-0 items-center rounded-full transition-all focus:outline-none focus:ring-2 focus:ring-amber-600/20 focus:ring-offset-2
               {midtransEnvironment === 'production' ? 'bg-amber-600 shadow-inner' : 'bg-slate-300'}"
      >
        <span class="inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform duration-200
               {midtransEnvironment === 'production' ? 'translate-x-6' : 'translate-x-1'}"></span>
      </button>
    </div>
  </div>
</div>
