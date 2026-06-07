<script lang="ts">
  import { Calendar } from 'lucide-svelte';

  let {
    autoBillingEnabled = $bindable(false),
    autoSemesterBillingEnabled = $bindable(false),
    billingHijriDay = $bindable('1'),
    settingLoading = false,
    onToggle,
    onDayChange
  }: {
    autoBillingEnabled: boolean;
    autoSemesterBillingEnabled: boolean;
    billingHijriDay: string;
    settingLoading: boolean;
    onToggle: (key: string) => void;
    onDayChange: () => void;
  } = $props();
</script>

<div class="bg-white border border-slate-200 rounded-2xl shadow-sm p-6 hover:shadow-md transition-all duration-200">
  <div class="flex items-start gap-4 mb-5">
    <div class="p-3 rounded-xl bg-blue-50 border border-blue-100 flex-shrink-0">
      <Calendar size={22} class="text-blue-700" />
    </div>
    <div class="flex-1 min-w-0">
      <h2 class="font-extrabold text-slate-900 text-lg">Jadwal Auto Billing</h2>
      <p class="text-xs text-slate-500 mt-0.5">Otomatisasi pembuatan invoice tagihan berbasis kalender Hijriah.</p>
    </div>
  </div>

  <div class="space-y-4">
    <div class="bg-slate-50/50 rounded-2xl p-4 border border-slate-150 shadow-sm transition-all hover:bg-slate-50">
      <div class="flex items-center justify-between mb-4">
        <div>
          <p class="text-sm font-bold text-slate-900">Auto Billing Syahriyyah</p>
          <p class="text-xs text-slate-500 mt-0.5 font-medium">Tagihan bulanan dibuat otomatis tiap awal bulan Hijriah.</p>
        </div>
        <button 
          type="button"
          onclick={() => onToggle('auto_generate_billing')}
          disabled={settingLoading}
          aria-label="Toggle Auto Billing Syahriyyah"
          title="Toggle Auto Billing Syahriyyah"
          class="relative inline-flex h-6.5 w-12 items-center rounded-full transition-all focus:outline-none focus:ring-2 focus:ring-emerald-800/20 focus:ring-offset-2
                 {autoBillingEnabled ? 'bg-emerald-850 shadow-inner' : 'bg-slate-300'}"
        >
          <span class="inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform duration-200
                 {autoBillingEnabled ? 'translate-x-6' : 'translate-x-1'}"></span>
        </button>
      </div>

      <div class="bg-white border border-slate-200/80 rounded-xl p-3.5 space-y-3 shadow-inner">
        <div class="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-4">
          <div class="flex flex-col gap-1 w-full md:w-1/2">
            <label for="billing-hijri-day-select" class="text-[10px] font-bold text-slate-500 uppercase tracking-wider">Hari Penerbitan Tagihan</label>
            <select
              id="billing-hijri-day-select"
              bind:value={billingHijriDay}
              onchange={onDayChange}
              disabled={settingLoading || !autoBillingEnabled}
              class="w-full px-3 py-2 rounded-lg bg-white border border-slate-350 text-sm font-semibold outline-none focus:ring-2 focus:ring-emerald-800/20 focus:border-emerald-800 disabled:opacity-50 disabled:bg-slate-50 transition-all cursor-pointer"
            >
              {#each Array.from({ length: 30 }, (_, i) => String(i + 1)) as day}
                <option value={day}>Tanggal {day} Hijriah</option>
              {/each}
            </select>
          </div>
          <div class="flex-1 text-[11px] text-slate-500 leading-relaxed md:border-l md:border-slate-200 md:pl-4">
            <span class="font-bold text-slate-800 block mb-0.5">Sistem Notifikasi H-5:</span>
            Wali santri otomatis menerima notifikasi pengingat tepat 5 hari sebelum tanggal tagihan di atas. Bulan & tahun dinamis mengikuti bulan berjalan.
          </div>
        </div>
      </div>
    </div>

    <div class="bg-slate-50/50 rounded-2xl p-4 border border-slate-150 shadow-sm flex items-center justify-between transition-all hover:bg-slate-50">
      <div>
        <p class="text-sm font-bold text-slate-900">Auto Billing Semester</p>
        <p class="text-xs text-slate-500 mt-0.5 font-medium">Tagihan semester (ngaji/sekolah) diterbitkan otomatis di bulan Syawal & Rabi'ul Akhir.</p>
      </div>
      <button 
        type="button"
        onclick={() => onToggle('auto_generate_semester_billing')}
        disabled={settingLoading}
        aria-label="Toggle Auto Billing Semester"
        title="Toggle Auto Billing Semester"
        class="relative inline-flex h-6.5 w-12 items-center rounded-full transition-all focus:outline-none focus:ring-2 focus:ring-blue-700/20 focus:ring-offset-2
               {autoSemesterBillingEnabled ? 'bg-blue-800 shadow-inner' : 'bg-slate-300'}"
      >
        <span class="inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform duration-200
               {autoSemesterBillingEnabled ? 'translate-x-6' : 'translate-x-1'}"></span>
      </button>
    </div>
  </div>
</div>
