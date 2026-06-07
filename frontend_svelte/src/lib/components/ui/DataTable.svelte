<script lang="ts">
  import type { Snippet } from 'svelte';
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
</script>

<div class="datatable">
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

<style>
  .datatable {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
    background: white;
    border: 1px solid rgba(148, 163, 184, 0.2);
    border-radius: 12px;
    overflow: hidden;
  }

  .datatable__header {
    flex-shrink: 0;
    border-bottom: 1px solid #f1f5f9;
  }

  .datatable__body {
    flex: 1;
    overflow: auto;
    min-height: 0;
  }

  /* Sticky thead within the scroll container */
  .datatable__body :global(thead) {
    position: sticky;
    top: 0;
    z-index: 10;
  }

  .datatable__body :global(thead tr) {
    background: #f9fafb;
  }

  .datatable__empty {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 3rem 1rem;
    min-height: 200px;
  }

  .datatable__footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 12px 16px;
    border-top: 1px solid #e2e8f0;
    background: white;
    flex-shrink: 0;
  }

  .datatable__footer-info {
    font-size: 0.875rem;
    color: #4b5563;
    white-space: nowrap;
  }

  .datatable__footer-nav {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .datatable__footer-page {
    padding: 6px 12px;
    background: #f0fdf4;
    color: #15803d;
    border-radius: 8px;
    font-size: 0.875rem;
    font-weight: 700;
    border: 1px solid #bbf7d0;
    white-space: nowrap;
  }

  @media (max-width: 640px) {
    .datatable__footer {
      flex-direction: column;
      gap: 8px;
    }
    .datatable__footer-info {
      font-size: 0.75rem;
    }
  }
</style>
