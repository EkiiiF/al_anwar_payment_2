<script lang="ts">
  import { User, Phone, Mail, MapPin, FileText, ChevronDown, ChevronUp, CheckCircle, Clock } from 'lucide-svelte';
  import { Modal, Badge, Card, EmptyState } from '$lib/components';
  import { formatDate, formatRupiah, getHijriMonthName } from '$lib/utils';
  import type { Student } from '$lib/types';

  let {
    open = false,
    onclose,
    selectedStudent = null
  }: {
    open: boolean;
    onclose: () => void;
    selectedStudent: Student | null;
  } = $props();

  let showPaid = $state(false);
  let showDetails = $state(false);

  const fullName = $derived(
    selectedStudent
      ? [selectedStudent.name?.first_name, selectedStudent.name?.middle_name, selectedStudent.name?.last_name]
          .filter(Boolean)
          .join(' ')
      : ''
  );

  const initials = $derived(
    selectedStudent
      ? [selectedStudent.name?.first_name, selectedStudent.name?.last_name]
          .filter(Boolean)
          .map(n => n[0].toUpperCase())
          .join('')
      : ''
  );

  const unpaidInvoices = $derived(selectedStudent?.invoices?.filter(inv => inv.status !== 'paid') ?? []);
  const paidInvoices = $derived(selectedStudent?.invoices?.filter(inv => inv.status === 'paid') ?? []);

  const unpaidByCategory = $derived((() => {
    const groups: Record<string, { categoryName: string; invoices: any[] }> = {};
    for (const inv of unpaidInvoices) {
      const catName = inv.category?.name || 'Tagihan Lainnya';
      if (!groups[catName]) {
        groups[catName] = { categoryName: catName, invoices: [] };
      }
      groups[catName].invoices.push(inv);
    }
    return Object.values(groups);
  })());
</script>

<Modal open={open} onclose={onclose} title="Detail & Rincian Tagihan Santri" size="xl">
  {#snippet children()}
    <div class="space-y-6">

      <!-- Collapsible Detailed Student Info -->
      <div class="border border-slate-200 rounded-2xl overflow-hidden bg-slate-50/20">
        <button
          onclick={() => showDetails = !showDetails}
          class="w-full flex items-center justify-between px-5 py-4 bg-slate-50 hover:bg-slate-100/80 transition-all font-bold text-sm text-slate-700 cursor-pointer outline-none border-none"
        >
          <div class="flex items-center gap-2.5">
            <User size={16} class="text-emerald-700" />
            <span>{fullName} {showDetails ? '' : ''}</span>
          </div>
          <div class="flex items-center gap-1 text-slate-400">
            {#if showDetails}
              <ChevronUp size={16} />
            {:else}
              <ChevronDown size={16} />
            {/if}
          </div>
        </button>

        {#if showDetails}
          <div class="p-5 border-t border-slate-200 bg-white grid grid-cols-1 md:grid-cols-3 gap-6 transition-all">
            <!-- Data Pribadi -->
            <div class="space-y-2.5">
              <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider flex items-center gap-1.5 select-none">
                <User size={14} class="text-slate-400" />
                <span>Data Pribadi</span>
              </h4>
              <div class="bg-slate-50/50 rounded-xl border border-slate-200 p-4 space-y-2 text-xs">
                <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">NIK</span><span class="font-bold text-slate-800">{selectedStudent?.nik || '-'}</span></div>
                <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Tgl Lahir</span><span class="font-bold text-slate-800">{selectedStudent?.birth_date ? formatDate(selectedStudent.birth_date) : '-'}</span></div>
                <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Gender</span><span class="font-bold text-slate-800">{selectedStudent?.gender === 'M' || selectedStudent?.gender === 'L' ? 'Laki-laki' : 'Perempuan'}</span></div>
                <div class="flex justify-between items-start gap-2"><span class="text-slate-500 font-medium whitespace-nowrap">Kelas</span><span class="font-bold text-emerald-800 text-right">{selectedStudent?.muhadhoroh_level === 0 ? 'Lulus' : 'Muhadhoroh ' + selectedStudent?.muhadhoroh_level + ' Sem ' + selectedStudent?.current_semester}</span></div>
                <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Diskon</span><span class="font-bold text-purple-800 bg-purple-50 px-1.5 py-0.5 rounded border border-purple-100">{selectedStudent?.status?.discount_percentage ?? 0}%</span></div>
              </div>
            </div>

            <!-- Data Wali -->
            <div class="space-y-2.5">
              <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider flex items-center gap-1.5 select-none">
                <Phone size={14} class="text-slate-400" />
                <span>Data Wali</span>
              </h4>
              {#if selectedStudent?.guardians && selectedStudent.guardians.length > 0}
                <div class="space-y-3">
                  {#each selectedStudent.guardians as guardian}
                    <div class="bg-slate-50/50 rounded-xl border border-slate-200 p-4 space-y-2 text-xs">
                      <p class="font-extrabold text-slate-900 border-b border-slate-200/50 pb-1">
                        {[guardian.name?.first_name, guardian.name?.middle_name, guardian.name?.last_name].filter(Boolean).join(' ')}
                      </p>
                      <div class="flex items-center gap-1.5 text-slate-600 font-medium"><Phone size={11} class="text-slate-400" /> {guardian.phone || '-'}</div>
                      <div class="flex items-center gap-1.5 text-slate-600 font-medium"><Mail size={11} class="text-slate-400" /> {guardian.email || '-'}</div>
                      <div class="flex justify-between items-center pt-1 border-t border-slate-100"><span class="text-slate-500 font-medium">Username</span><span class="font-mono font-bold text-slate-800">{guardian.user?.username ?? '-'}</span></div>
                      <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Hubungan dengan santri</span><span class="font-bold text-slate-800">{guardian.relation}</span></div>
                    </div>
                  {/each}
                </div>
              {:else}
                <div class="bg-slate-50/30 rounded-xl border border-slate-200 border-dashed p-4 text-center">
                  <p class="text-xs text-slate-400 italic">Belum ada data wali santri</p>
                </div>
              {/if}
            </div>

            <!-- Alamat Domisili -->
            <div class="space-y-2.5">
              <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider flex items-center gap-1.5 select-none">
                <MapPin size={14} class="text-slate-400" />
                <span>Alamat Domisili</span>
              </h4>
              {#if selectedStudent?.addresses && selectedStudent.addresses.length > 0}
                <div class="space-y-3">
                  {#each selectedStudent.addresses as addr}
                    <div class="bg-slate-50/50 rounded-xl border border-slate-200 p-4 space-y-2 text-xs">
                      <p class="font-extrabold text-slate-900 leading-snug">{addr.address_line || '-'}</p>
                      {#if addr.village || addr.district}
                        <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Desa / Kec</span><span class="font-bold text-slate-800">{[addr.village, addr.district].filter(Boolean).join(', ')}</span></div>
                      {/if}
                      {#if addr.city || addr.province}
                        <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Kota / Prov</span><span class="font-bold text-slate-800">{[addr.city, addr.province].filter(Boolean).join(', ')}</span></div>
                      {/if}
                      {#if addr.postal_code}
                        <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Kode Pos</span><span class="font-bold text-slate-800">{addr.postal_code}</span></div>
                      {/if}
                      {#if addr.country && addr.country !== 'Indonesia'}
                        <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Negara</span><span class="font-bold text-slate-800">{addr.country}</span></div>
                      {/if}
                      {#if addr.is_primary}
                        <div class="pt-1 border-t border-slate-200/50 flex justify-between items-center">
                          <span class="text-slate-500 font-medium">Status</span>
                          <Badge label="Utama" variant="info" class="text-[9px] py-0.5" />
                        </div>
                      {/if}
                    </div>
                  {/each}
                </div>
              {:else}
                <div class="bg-slate-50/30 rounded-xl border border-slate-200 border-dashed p-4 text-center">
                  <p class="text-xs text-slate-400 italic">Belum ada data alamat</p>
                </div>
              {/if}
            </div>
          </div>
        {/if}
      </div>

      <!-- Billing Details Section -->
      <div class="border-t border-slate-100 pt-5 space-y-4">
        <h3 class="text-base font-black text-slate-800 flex items-center gap-2 select-none">
          <FileText size={18} class="text-emerald-700" /> 
          <span>Rincian Tagihan Aktif</span>
        </h3>

        {#if unpaidByCategory.length > 0}
          <div class="grid grid-cols-1 gap-4">
            {#each unpaidByCategory as group}
              <div class="bg-slate-50/50 border border-slate-200/85 rounded-2xl p-5 space-y-3">
                <div class="flex items-center justify-between border-b border-slate-200/60 pb-2.5">
                  <h4 class="font-black text-slate-900 text-sm flex items-center gap-2">
                    <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span>
                    <span>{group.categoryName}</span>
                  </h4>
                  <Badge label={`${group.invoices.length} Tagihan`} variant="warning" class="text-[10px] font-bold" />
                </div>

                <div class="overflow-x-auto">
                  <table class="w-full text-left text-xs border-collapse">
                    <thead>
                      <tr class="text-slate-400 font-semibold border-b border-slate-200/60 whitespace-nowrap">
                        <th class="pb-2">Bulan/Periode</th>
                        <th class="pb-2">Nomor Invoice</th>
                        <th class="pb-2">Semester</th>
                        <!-- <th class="pb-2">Jatuh Tempo</th> -->
                        <th class="pb-2 text-right">Jumlah</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-slate-200/60 text-slate-700 whitespace-nowrap">
                      {#each group.invoices as invoice (invoice.id)}
                        <tr class="hover:bg-slate-100/30 transition-colors">
                          <td class="py-3 font-semibold text-slate-800">
                            {#if invoice.hijri_month}
                              {getHijriMonthName(invoice.hijri_month)} {invoice.hijri_year} H
                            {:else}
                              {invoice.category?.name}
                            {/if}
                          </td>
                          <td class="py-3 font-mono text-[10px] text-slate-400">{invoice.invoice_number}</td>
                          <td class="py-3">
                            {#if invoice.semester}
                              <span class="bg-slate-100 text-slate-600 px-1.5 py-0.5 rounded font-bold">Sem {invoice.semester}</span>
                            {/if}
                            {#if invoice.academic_year_label}
                              <span class="text-emerald-600 font-bold ml-1">TA {invoice.academic_year_label}</span>
                            {/if}
                          </td>
                          <!-- <td class="py-3 text-slate-500">
                            <span class="inline-flex items-center gap-1">
                              <Clock size={10} class="text-amber-500" />
                              {formatDate(invoice.due_date)}
                            </span>
                          </td> -->
                          <td class="py-3 text-right font-black text-slate-900">{formatRupiah(invoice.amount_due)}</td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <EmptyState title="Tidak ada tagihan aktif" description="Santri ini tidak memiliki tagihan yang perlu dibayar." class="border-slate-200 bg-slate-50/10" />
        {/if}
      </div>

      <!-- Paid Invoices Section -->
      {#if paidInvoices.length > 0}
        <div class="border border-slate-200 rounded-2xl overflow-hidden bg-slate-50/20">
          <button
            onclick={() => showPaid = !showPaid}
            class="w-full flex items-center justify-between px-5 py-4 bg-slate-50 hover:bg-slate-100/80 transition-all font-bold text-sm text-slate-700 cursor-pointer outline-none border-none"
          >
            <div class="flex items-center gap-2.5">
              <CheckCircle size={16} class="text-green-600" />
              <span>Riwayat Tagihan Lunas ({paidInvoices.length})</span>
            </div>
            <div class="flex items-center gap-1 text-slate-400">
              <span class="text-xs font-semibold">{showPaid ? 'Sembunyikan' : 'Lihat Detail'}</span>
              {#if showPaid}
                <ChevronUp size={16} />
              {:else}
                <ChevronDown size={16} />
              {/if}
            </div>
          </button>

          {#if showPaid}
            <div class="p-5 border-t border-slate-200 bg-white max-h-72 overflow-y-auto">
              <table class="w-full text-left text-xs border-collapse">
                <thead>
                  <tr class="text-slate-400 font-semibold border-b border-slate-200/60 whitespace-nowrap">
                    <th class="pb-2">Kategori & Periode</th>
                    <th class="pb-2">Nomor Invoice</th>
                    <th class="pb-2">Semester</th>
                    <th class="pb-2">Jumlah</th>
                    <th class="pb-2 text-right">Status</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-200/60 text-slate-700 whitespace-nowrap">
                  {#each paidInvoices as invoice (invoice.id)}
                    <tr class="hover:bg-slate-100/30 transition-colors">
                      <td class="py-3 font-semibold text-slate-800">
                        <span class="text-emerald-700 mr-1">{invoice.category?.name}</span>
                        {#if invoice.hijri_month}
                          · {getHijriMonthName(invoice.hijri_month)} {invoice.hijri_year} H
                        {/if}
                      </td>
                      <td class="py-3 font-mono text-[10px] text-slate-400">{invoice.invoice_number}</td>
                      <td class="py-3">
                        {#if invoice.semester}
                          <span class="bg-slate-100 text-slate-600 px-1.5 py-0.5 rounded font-bold">Sem {invoice.semester}</span>
                        {/if}
                        {#if invoice.academic_year_label}
                          <span class="text-emerald-600 font-bold ml-1">TA {invoice.academic_year_label}</span>
                        {/if}
                      </td>
                      <td class="py-3 font-black text-slate-900">{formatRupiah(invoice.amount_due)}</td>
                      <td class="py-3 text-right">
                        <span class="inline-flex items-center gap-1 bg-green-50 text-green-700 border border-green-200 px-2 py-0.5 rounded-full font-bold text-[9px]">
                          Lunas
                        </span>
                      </td>
                    </tr>
                  {/each}
                </tbody>
              </table>
            </div>
          {/if}
        </div>
      {/if}

    </div>
  {/snippet}
</Modal>
