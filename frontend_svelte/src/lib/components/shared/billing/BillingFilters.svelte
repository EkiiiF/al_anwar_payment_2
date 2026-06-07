<script lang="ts">
  import { Filter } from 'lucide-svelte';
  import { Select } from '$lib/components';

  let {
    statusFilter = $bindable(''),
    monthFilter = $bindable(''),
    yearFilter = $bindable(''),
    monthOptions = [],
    yearOptions = []
  }: {
    statusFilter: string;
    monthFilter: string;
    yearFilter: string;
    monthOptions: { value: string; label: string }[];
    yearOptions: { value: string; label: string }[];
  } = $props();
</script>

<div class="bg-white border border-slate-200 rounded-2xl shadow-sm p-4 hover:shadow-md transition-all duration-200">
  <div class="flex flex-wrap items-center gap-4">
    <div class="flex items-center gap-2 text-sm font-bold text-slate-700 select-none">
      <Filter size={16} class="text-emerald-800" />
      <span>Filter Data:</span>
    </div>
    
    <div class="w-full sm:w-44">
      <Select
        id="f-status"
        bind:value={statusFilter}
        options={[
          { value: '', label: 'Semua Status' },
          { value: 'unpaid', label: 'Belum Bayar' },
          { value: 'paid', label: 'Lunas' }
        ]}
        class="w-full"
      />
    </div>

    <div class="w-full sm:w-44">
      <Select
        id="f-month"
        bind:value={monthFilter}
        options={[{ value: '', label: 'Semua Bulan' }, ...monthOptions]}
        class="w-full"
      />
    </div>

    <div class="w-full sm:w-44">
      <Select
        id="f-year"
        bind:value={yearFilter}
        options={[{ value: '', label: 'Semua Tahun' }, ...yearOptions]}
        class="w-full"
      />
    </div>
  </div>
</div>
