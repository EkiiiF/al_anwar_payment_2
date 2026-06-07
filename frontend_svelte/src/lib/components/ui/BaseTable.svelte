<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    columns,
    emptyMessage = 'Tidak ada data untuk ditampilkan.',
    isEmpty = false,
    class: className = '',
    headerRow,
    bodyRows
  }: {
    columns: string[];
    emptyMessage?: string;
    isEmpty?: boolean;
    class?: string;
    headerRow?: Snippet;
    bodyRows: Snippet;
  } = $props();
</script>

<div class="border border-slate-200/80 rounded-xl overflow-hidden {className}">
  <div class="overflow-x-auto">
    <table class="w-full text-sm text-left">
      <thead>
        {#if headerRow}
          {@render headerRow()}
        {:else}
          <tr class="bg-slate-50 border-b border-slate-200/60">
            {#each columns as col}
              <th class="px-4 py-3 text-xs font-medium text-slate-500 uppercase tracking-wider whitespace-nowrap">
                {col}
              </th>
            {/each}
          </tr>
        {/if}
      </thead>
      <tbody class="divide-y divide-slate-100">
        {#if isEmpty}
          <tr>
            <td colspan={columns.length} class="px-4 py-12 text-center text-sm text-slate-400">
              {emptyMessage}
            </td>
          </tr>
        {:else}
          {@render bodyRows()}
        {/if}
      </tbody>
    </table>
  </div>
</div>
