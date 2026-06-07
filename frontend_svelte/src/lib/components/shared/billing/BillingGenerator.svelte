<script lang="ts">
  import { PlayCircle, Layers } from 'lucide-svelte';
  import { Button, Alert, Select } from '$lib/components';
  import { getHijriMonthName } from '$lib/utils';

  let {
    genMode = $bindable('monthly'),
    genMonth = $bindable('1'),
    genYear = $bindable(''),
    genSemester = $bindable('1'),
    genHijriYear = $bindable(''),
    generating = false,
    genMessage = '',
    genSuccess = false,
    onGenerate
  }: {
    genMode?: 'monthly' | 'semester';
    genMonth?: string;
    genYear?: string;
    genSemester?: string;
    genHijriYear?: string;
    generating: boolean;
    genMessage: string;
    genSuccess: boolean;
    onGenerate: () => void;
  } = $props();
</script>

<div class="bg-white border border-slate-200 rounded-2xl shadow-sm p-6 hover:shadow-md transition-all duration-200">
  <div class="flex items-start gap-4 mb-5">
    <div class="p-3 rounded-xl bg-emerald-50 border border-emerald-100 flex-shrink-0">
      <PlayCircle size={22} class="text-emerald-800" />
    </div>
    <div class="flex-1 min-w-0">
      <h2 class="font-extrabold text-slate-900 text-lg">Buat Tagihan Manual</h2>
      <p class="text-xs text-slate-500 mt-0.5">Menerbitkan tagihan secara massal untuk santri aktif tipe reguler.</p>
    </div>
  </div>

  {#if genMessage}
    <Alert type={genSuccess ? 'success' : 'error'} message={genMessage} class="mb-5 rounded-xl border-slate-150" />
  {/if}

  <div class="flex gap-1 bg-slate-100/80 rounded-xl p-1 mb-5 w-fit border border-slate-200/50">
    <button
      type="button"
      onclick={() => genMode = 'monthly'}
      class="px-4 py-2 rounded-lg text-sm font-semibold tracking-wide transition-all duration-200
        {genMode === 'monthly' ? 'bg-emerald-800 text-white shadow-sm' : 'text-slate-650 hover:bg-slate-200/50 hover:text-slate-900'}"
    >
      Syahriyyah Bulanan
    </button>
    <button
      type="button"
      onclick={() => genMode = 'semester'}
      class="px-4 py-2 rounded-lg text-sm font-semibold tracking-wide transition-all duration-200 flex items-center gap-1.5
        {genMode === 'semester' ? 'bg-emerald-800 text-white shadow-sm' : 'text-slate-650 hover:bg-slate-200/50 hover:text-slate-900'}"
    >
      <Layers size={14} /> Tagihan Semester
    </button>
  </div>

  {#if genMode === 'monthly'}
    <div class="flex flex-col sm:flex-row gap-4 items-end">
      <div class="flex-1 w-full">
        <Select id="m" label="Bulan Hijriah" bind:value={genMonth} options={[
          { value: '1', label: '1. Muharram' },
          { value: '2', label: '2. Safar' },
          { value: '3', label: '3. Rabi\'ul Awal' },
          { value: '4', label: '4. Rabi\'ul Akhir' },
          { value: '5', label: '5. Jumadil Awal' },
          { value: '6', label: '6. Jumadil Akhir' },
          { value: '7', label: '7. Rajab' },
          { value: '8', label: '8. Sya\'ban' },
          { value: '9', label: '9. Ramadhan' },
          { value: '10', label: '10. Syawal' },
          { value: '11', label: '11. Dzulqa\'dah' },
          { value: '12', label: '12. Dzulhijjah' }
        ]} class="w-full" />
      </div>
      <div class="flex flex-col gap-1.5 flex-1 w-full">
        <label for="y" class="text-xs font-semibold text-slate-600 uppercase tracking-wider ml-0.5">Tahun Hijriah</label>
        <input 
          id="y" 
          type="number" 
          bind:value={genYear} 
          class="w-full px-3 py-2.5 rounded-lg bg-white border border-slate-350 text-sm focus:ring-2 focus:ring-emerald-800/20 focus:border-emerald-800 outline-none transition-all" 
          placeholder="Contoh: 1447" 
        />
      </div>
      <Button onclick={onGenerate} variant="primary" loading={generating} class="w-full sm:w-auto bg-teal-600 hover:bg-teal-700 text-white font-bold h-[42px] px-6 rounded-lg shadow-sm">
        {#snippet children()}Generate{/snippet}
      </Button>
    </div>
  {:else}
    <div class="flex flex-col sm:flex-row gap-4 items-end">
      <div class="flex-1 w-full">
        <Select id="s" label="Semester" bind:value={genSemester} options={[
          { value: '1', label: 'Semester 1 (Ganjil)' },
          { value: '2', label: 'Semester 2 (Genap)' }
        ]} class="w-full" />
      </div>
      <div class="flex flex-col gap-1.5 flex-1 w-full">
        <label for="hy" class="text-xs font-semibold text-slate-600 uppercase tracking-wider ml-0.5">Tahun Hijriah</label>
        <input 
          id="hy" 
          type="number" 
          bind:value={genHijriYear} 
          class="w-full px-3 py-2.5 rounded-lg bg-white border border-slate-350 text-sm focus:ring-2 focus:ring-emerald-800/20 focus:border-emerald-800 outline-none transition-all" 
          placeholder="Contoh: 1447" 
        />
      </div>
      <Button onclick={onGenerate} variant="primary" loading={generating} class="w-full sm:w-auto bg-teal-600 hover:bg-teal-700 text-white font-bold h-[42px] px-6 rounded-lg shadow-sm">
        {#snippet children()}Generate{/snippet}
      </Button>
    </div>
  {/if}
</div>
