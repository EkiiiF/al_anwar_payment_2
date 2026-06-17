<script lang="ts">
  import type { Snippet } from 'svelte';
  import { onMount } from 'svelte';
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';
  import Button from './Button.svelte';
  import EmptyState from './EmptyState.svelte';

  let {
    pagination = null,
    onPrevPage,
    onNextPage,
    isEmpty = false,
    emptyTitle = 'Tidak ada data',
    emptyDescription = '',
    paginationLabel = 'data',
    headerSlot,
    children,
  }: {
    pagination?: { page: number; limit: number; total: number; pages: number } | null;
    onPrevPage?: () => void;
    onNextPage?: () => void;
    isEmpty?: boolean;
    emptyTitle?: string;
    emptyDescription?: string;
    paginationLabel?: string;
    headerSlot?: Snippet;
    children: Snippet;
  } = $props();

  const currentPage = $derived(pagination?.page ?? 1);
  const totalPages = $derived(pagination?.pages ?? 1);
  const limit = $derived(pagination?.limit ?? 10);
  const total = $derived(pagination?.total ?? 0);
  const showStart = $derived(((currentPage - 1) * limit) + 1);
  const showEnd = $derived(Math.min(currentPage * limit, total));

  // Dynamically calculate max-height so table fills to viewport bottom
  let containerEl: HTMLDivElement;
  let maxHeight = $state('none');

  function recalcHeight() {
    if (!containerEl) return;
    const rect = containerEl.getBoundingClientRect();
    const viewportH = window.innerHeight;
    const remaining = viewportH - rect.top - 24; // 24px bottom padding
    maxHeight = `${Math.max(remaining, 300)}px`;
  }

  onMount(() => {
    recalcHeight();
    const ro = new ResizeObserver(recalcHeight);
    ro.observe(document.body);
    window.addEventListener('resize', recalcHeight);
    return () => {
      ro.disconnect();
      window.removeEventListener('resize', recalcHeight);
    };
  });
</script>

<div
  class="datatable"
  bind:this={containerEl}
  style:max-height={maxHeight}
>
  {#if headerSlot}
    <div class="datatable__header">
      {@render headerSlot()}
    </div>
  {/if}

  <div class="datatable__body">
    {#if isEmpty}
      <div class="datatable__empty">
        <EmptyState title={emptyTitle} description={emptyDescription} />
      </div>
    {:else}
      {@render children()}
    {/if}
  </div>

  {#if pagination}
    <div class="datatable__footer">
      <p class="datatable__footer-info">
        Menampilkan <span class="font-bold text-gray-900">{showStart} - {showEnd}</span>
        dari <span class="font-bold text-gray-900">{total}</span> {paginationLabel}
      </p>
      <div class="datatable__footer-nav">
        <Button variant="secondary" onclick={onPrevPage} disabled={currentPage <= 1} size="sm">
          {#snippet children()}<ChevronLeft size={16} /> Sebelumnya{/snippet}
        </Button>
        <div class="datatable__footer-page">
          Halaman {currentPage} / {totalPages}
        </div>
        <Button variant="secondary" onclick={onNextPage} disabled={currentPage >= totalPages} size="sm">
          {#snippet children()}Selanjutnya <ChevronRight size={16} />{/snippet}
        </Button>
      </div>
    </div>
  {/if}
</div>

