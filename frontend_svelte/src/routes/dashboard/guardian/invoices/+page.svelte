<script lang="ts">
  import { onMount } from 'svelte';
  import { guardianApi } from '$lib/api';
  import { formatRupiah, getInvoiceStatusStyle, getHijriMonthName } from '$lib/utils';
  import { Spinner, Alert, Button, Card, Badge, EmptyState, CheckoutSummary, PondokBillingColumn, SemesterBillingColumn, InvoiceHistoryList } from '$lib/components';
  import type { Invoice } from '$lib/types';
  import { ShoppingCart, CreditCard, CheckSquare, Layers, Calendar } from 'lucide-svelte';

  let invoices    = $state<Invoice[]>([]);
  let loading     = $state(true);
  let error       = $state('');
  let paying      = $state(false);
  let payError    = $state('');

  let selectedIds = $state<Set<string>>(new Set());

  const unpaidInvoices = $derived(invoices.filter(i => i.status === 'unpaid'));

  type SemGroup = { semester: number; academicYear: string; invoices: Invoice[]; total: number; paidCount: number };
  const semesterGroups = $derived.by(() => {
    const groups = new Map<string, SemGroup>();
    for (const inv of unpaidInvoices) {
      if (inv.category?.name === 'Syahriyyah Pondok') {
        continue;
      }
      const key = `${inv.semester}-${inv.academic_year_label || 'unknown'}`;
      if (!groups.has(key)) {
        groups.set(key, {
          semester: inv.semester,
          academicYear: inv.academic_year_label || '-',
          invoices: [],
          total: 0,
          paidCount: 0
        });
      }
      const g = groups.get(key)!;
      g.invoices.push(inv);
      g.total += inv.amount_due;
    }
    return Array.from(groups.values());
  });

  const pondokInvoices = $derived(
    unpaidInvoices.filter(i => i.category?.name === 'Syahriyyah Pondok')
  );

  const semesterInvoices = $derived(
    unpaidInvoices.filter(i => i.category?.name !== 'Syahriyyah Pondok')
  );

  const total = $derived(
    unpaidInvoices
      .filter(i => selectedIds.has(i.id))
      .reduce((sum, i) => sum + i.amount_due, 0)
  );
  const selectedCount = $derived(selectedIds.size);

  function toggleSelect(id: string) {
    const next = new Set(selectedIds);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    selectedIds = next;
  }

  function clearSelection(){ selectedIds = new Set(); }

  function selectAllPondok() {
    const next = new Set(selectedIds);
    pondokInvoices.forEach(i => next.add(i.id));
    selectedIds = next;
  }

  function deselectAllPondok() {
    const next = new Set(selectedIds);
    pondokInvoices.forEach(i => next.delete(i.id));
    selectedIds = next;
  }

  function selectAllSemester() {
    const next = new Set(selectedIds);
    semesterInvoices.forEach(i => next.add(i.id));
    selectedIds = next;
  }

  function deselectAllSemester() {
    const next = new Set(selectedIds);
    semesterInvoices.forEach(i => next.delete(i.id));
    selectedIds = next;
  }

  function selectSemester(group: SemGroup) {
    const next = new Set(selectedIds);
    const allSelected = group.invoices.every(i => next.has(i.id));
    if (allSelected) {
      group.invoices.forEach(i => next.delete(i.id));
    } else {
      group.invoices.forEach(i => next.add(i.id));
    }
    selectedIds = next;
  }

  onMount(async () => {
    try {
      const res = await guardianApi.getInvoices();
      invoices = res.data ?? [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memuat tagihan.';
    } finally {
      loading = false;
    }
  });

  async function handleCheckout() {
    if (selectedCount === 0) return;
    paying   = true;
    payError = '';
    try {
      const res = await guardianApi.createPayment([...selectedIds]);
      const snapToken = res.data.token;

      type SnapWindow = Window & typeof globalThis & {
        snap?: { pay: (token: string, options: Record<string, unknown>) => void }
      };
      const w = window as SnapWindow;
      if (!w.snap) {
        payError = 'Midtrans Snap tidak tersedia. Silakan muat ulang halaman.';
        paying = false;
        return;
      }
      w.snap.pay(snapToken, {
        onSuccess: () => { window.location.href = '/dashboard/guardian/history'; },
        onPending: () => { window.location.href = '/dashboard/guardian/history'; },
        onError:   (result: unknown) => {
          console.error('[Midtrans Error]', result);
          payError = 'Pembayaran gagal. Silakan coba lagi.';
          paying = false;
        },
        onClose: () => { paying = false; }
      });
    } catch (e: unknown) {
      payError = e instanceof Error ? e.message : 'Gagal memulai pembayaran.';
      paying = false;
    }
  }

  async function refreshStatus(orderId: string) {
    loading = true;
    try {
      await guardianApi.checkPaymentStatus(orderId);
      const res = await guardianApi.getInvoices();
      invoices = res.data ?? [];
    } catch (e: unknown) {
      error = e instanceof Error ? e.message : 'Gagal memperbarui status.';
    } finally {
      loading = false;
    }
  }

  function resumePayment(snapToken: string) {
    type SnapWindow = Window & typeof globalThis & {
      snap?: { pay: (token: string, options: Record<string, unknown>) => void }
    };
    const w = window as SnapWindow;
    if (!w.snap) {
      payError = 'Midtrans Snap tidak tersedia. Silakan muat ulang halaman.';
      return;
    }
    paying = true;
    w.snap.pay(snapToken, {
      onSuccess: () => { window.location.href = '/dashboard/guardian/history'; },
      onPending: () => { window.location.href = '/dashboard/guardian/history'; },
      onError: (result: unknown) => {
        console.error('[Midtrans Error]', result);
        payError = 'Pembayaran gagal. Silakan coba lagi.';
        paying = false;
      },
      onClose: () => { paying = false; }
    });
  }
</script>

<svelte:head>
  <title>Tagihan Saya | Guardian - Al-Anwar Payment</title>
  <script src="https://app.sandbox.midtrans.com/snap/snap.js" data-client-key={import.meta.env.VITE_MIDTRANS_CLIENT_KEY}></script>
</svelte:head>

<div class="space-y-6">
  <div>
    <h1 class="text-2xl font-black text-gray-900 tracking-tight">Tagihan Saya</h1>
    <p class="text-gray-555 text-sm mt-1">Pilih satu atau beberapa tagihan untuk dibayar sekaligus.</p>
  </div>

  {#if error}
    <Alert type="error" message={error} />
  {:else if loading}
    <Spinner size="lg" />
  {:else}
    {#if selectedCount > 0}
      <CheckoutSummary
        selectedCount={selectedCount}
        total={total}
        paying={paying}
        payError={payError}
        onClear={clearSelection}
        onCheckout={handleCheckout}
      />
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 items-start">
      <PondokBillingColumn
        pondokInvoices={pondokInvoices}
        selectedIds={selectedIds}
        toggleSelect={toggleSelect}
        selectAllPondok={selectAllPondok}
        deselectAllPondok={deselectAllPondok}
      />

      <SemesterBillingColumn
        semesterInvoices={semesterInvoices}
        semesterGroups={semesterGroups}
        selectedIds={selectedIds}
        toggleSelect={toggleSelect}
        selectAllSemester={selectAllSemester}
        deselectAllSemester={deselectAllSemester}
        selectSemester={selectSemester}
      />
    </div>

    {@const otherInvoices = invoices.filter(i => i.status !== 'unpaid')}
    {#if otherInvoices.length > 0}
      <InvoiceHistoryList
        otherInvoices={otherInvoices}
        refreshStatus={refreshStatus}
        resumePayment={resumePayment}
      />
    {/if}
  {/if}
</div>
