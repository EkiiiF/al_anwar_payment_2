<script lang="ts">
  import { CreditCard, ShoppingCart } from 'lucide-svelte';
  import { Button, Alert } from '$lib/components';
  import { formatRupiah } from '$lib/utils';

  let {
    selectedCount = 0,
    total = 0,
    paying = false,
    payError = '',
    onClear,
    onCheckout
  }: {
    selectedCount: number;
    total: number;
    paying: boolean;
    payError: string;
    onClear: () => void;
    onCheckout: () => void;
  } = $props();
</script>

<div class="sticky top-4 z-20 transition-all duration-300">
  <div class="bg-white border border-slate-200 shadow-sm rounded-xl p-3.5 sm:p-4 transition-all duration-300">
    {#if payError}
      <Alert type="error" message={payError} class="mb-4 rounded-xl" />
    {/if}
    
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <div class="w-12 h-12 rounded-xl bg-emerald-50 border border-emerald-100 flex items-center justify-center flex-shrink-0">
          <ShoppingCart size={22} class="text-emerald-800 animate-pulse" />
        </div>
        <div>
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">
            Total Terpilih ({selectedCount} Item)
          </p>
          <p class="text-2xl font-black text-emerald-800 mt-0.5 leading-none">
            {formatRupiah(total)}
          </p>
        </div>
      </div>

      <div class="flex flex-wrap items-center gap-2.5">
        <Button 
          onclick={onClear} 
          variant="outline" 
          size="md" 
          class="flex-1 sm:flex-none"
        >
          {#snippet children()}Batalkan Pilihan{/snippet}
        </Button>
        <Button 
          onclick={onCheckout} 
          variant="primary" 
          size="md" 
          loading={paying}
          class="flex-1 sm:flex-none flex items-center justify-center gap-2"
        >
          {#snippet children()}
            <CreditCard size={15} />
            <span>Bayar Sekarang</span>
          {/snippet}
        </Button>
      </div>
    </div>
  </div>
</div>
