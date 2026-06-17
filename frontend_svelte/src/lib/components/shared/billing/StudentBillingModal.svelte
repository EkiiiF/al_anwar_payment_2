<script lang="ts">
  import { User, Phone, Mail, MapPin, FileText } from 'lucide-svelte';
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

  const fullName = $derived(
    selectedStudent
      ? [selectedStudent.name?.first_name, selectedStudent.name?.middle_name, selectedStudent.name?.last_name]
          .filter(Boolean)
          .join(' ')
      : ''
  );
</script>

<Modal open={open} onclose={onclose} size="xl">
  {#snippet children()}
    <div class="space-y-6">
      <div class="border-b border-slate-100 pb-4">
        <h2 class="text-xl font-extrabold text-slate-900 leading-tight">
          {fullName}
        </h2>
        <div class="flex items-center gap-2 mt-1">
          <p class="text-xs text-slate-500 font-mono">NIS: {selectedStudent?.student_number}</p>
          {#if selectedStudent?.status}
            <Badge label={selectedStudent.status.name} variant="purple" class="text-[10px] py-0.5" />
          {/if}
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
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
            <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Diskon</span><span class="font-bold text-purple-750 bg-purple-50 px-1.5 py-0.5 rounded border border-purple-100">{selectedStudent?.status?.discount_percentage ?? 0}%</span></div>
          </div>
        </div>

        <div class="space-y-2.5">
          <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider flex items-center gap-1.5 select-none">
            <Phone size={14} class="text-slate-400" />
            <span>Data Wali / Guardian</span>
          </h4>
          {#if selectedStudent?.guardians && selectedStudent.guardians.length > 0}
            <div class="space-y-3">
              {#each selectedStudent.guardians as guardian}
                <div class="bg-slate-50/50 rounded-xl border border-slate-200 p-4 space-y-2 text-xs">
                  <p class="font-extrabold text-slate-900 border-b border-slate-200/50 pb-1">
                    {[guardian.name?.first_name, guardian.name?.middle_name, guardian.name?.last_name].filter(Boolean).join(' ')}
                  </p>
                  <div class="flex items-center gap-1.5 text-slate-655 font-medium"><Phone size={11} class="text-slate-450" /> {guardian.phone || '-'}</div>
                  <div class="flex items-center gap-1.5 text-slate-655 font-medium"><Mail size={11} class="text-slate-450" /> {guardian.email || '-'}</div>
                  <div class="flex justify-between items-center pt-1 border-t border-slate-100"><span class="text-slate-500 font-medium">Username</span><span class="font-mono font-bold text-slate-800">{guardian.user?.username ?? '-'}</span></div>
                  <div class="flex justify-between items-center"><span class="text-slate-500 font-medium">Relasi</span><span class="font-bold text-slate-800">{guardian.relation}</span></div>
                </div>
              {/each}
            </div>
          {:else}
            <div class="bg-slate-50/30 rounded-xl border border-slate-200 border-dashed p-4 text-center">
              <p class="text-xs text-slate-400 italic">Belum ada data wali santri</p>
            </div>
          {/if}
        </div>

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

      <div class="border-t border-slate-100 pt-5">
        <h3 class="text-base font-extrabold text-slate-900 mb-4 flex items-center gap-2 select-none">
          <FileText size={18} class="text-emerald-800" /> 
          <span>Rincian Tagihan Santri</span>
        </h3>
        <div class="space-y-3 max-h-80 overflow-y-auto pr-1">
          {#each selectedStudent?.invoices ?? [] as invoice (invoice.id)}
            <Card class="border-slate-200 hover:border-slate-300 transition-colors shadow-sm !p-4 bg-white">
              <div class="flex items-center justify-between gap-4">
                <div class="min-w-0">
                  <p class="font-bold text-slate-900 text-sm">
                    {#if invoice.category}
                      <span class="text-emerald-800">{invoice.category.name}</span>
                    {/if}
                    {#if invoice.hijri_month}
                      <span class="text-slate-400 font-normal"> · </span>{getHijriMonthName(invoice.hijri_month)} {invoice.hijri_year} H
                    {/if}
                  </p>
                  <p class="text-[10px] text-slate-400 font-mono mt-0.5">{invoice.invoice_number}</p>
                  
                  <div class="flex items-center gap-2 mt-1">
                    {#if invoice.semester}
                      <Badge label={`Sem ${invoice.semester}`} variant="warning" class="text-[9px] py-0.5" />
                    {/if}
                    {#if invoice.academic_year_label}
                      <span class="text-[10px] text-emerald-600 font-bold">TA {invoice.academic_year_label}</span>
                    {/if}
                  </div>
                </div>
                <div class="text-right flex-shrink-0">
                  <p class="font-extrabold text-slate-900 text-base leading-none mb-1.5">{formatRupiah(invoice.amount_due)}</p>
                  <div class="flex flex-col items-end gap-1">
                    <Badge
                      label={invoice.status === 'paid' ? 'Lunas' : 'Belum Bayar'}
                      variant={invoice.status === 'paid' ? 'success' : 'warning'}
                      dot
                      class="text-[9px] py-0.5"
                    />
                    <p class="text-[10px] text-slate-500 font-medium mt-0.5">Jatuh Tempo: {formatDate(invoice.due_date)}</p>
                  </div>
                </div>
              </div>
            </Card>
          {:else}
            <EmptyState title="Tidak ada tagihan untuk santri ini" description="Santri ini belum memiliki riwayat tagihan." class="border-slate-150 bg-slate-50/10" />
          {/each}
        </div>
      </div>
    </div>
  {/snippet}
</Modal>
