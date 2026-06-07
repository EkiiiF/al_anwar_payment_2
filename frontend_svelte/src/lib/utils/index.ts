/**
 * Format angka menjadi format Rupiah Indonesia
 */
export function formatRupiah(amount: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0
  }).format(amount);
}

/**
 * Format tanggal ke format lokal Indonesia
 */
export function formatDate(dateStr: string | undefined | null, withTime = false): string {
  if (!dateStr) return '-';
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return '-';

  const options: Intl.DateTimeFormatOptions = {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    ...(withTime && { hour: '2-digit', minute: '2-digit' })
  };
  return new Intl.DateTimeFormat('id-ID', options).format(date);
}

/**
 * Nama bulan dalam Bahasa Indonesia
 */
export const MONTH_NAMES = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
] as const;

export function getMonthName(month: number): string {
  return MONTH_NAMES[month - 1] ?? '-';
}

/**
 * Warna badge berdasarkan status invoice
 */
export function getInvoiceStatusStyle(status: string): { bg: string; text: string; label: string } {
  const map: Record<string, { bg: string; text: string; label: string }> = {
    unpaid:    { bg: 'bg-amber-500/10',   text: 'text-amber-400',   label: 'Belum Bayar' },
    pending:   { bg: 'bg-blue-500/10',    text: 'text-blue-400',    label: 'Menunggu' },
    paid:      { bg: 'bg-emerald-500/10', text: 'text-emerald-400', label: 'Lunas' },
    expired:   { bg: 'bg-gray-500/10',    text: 'text-gray-400',    label: 'Kadaluwarsa' },
    cancelled: { bg: 'bg-red-500/10',     text: 'text-red-400',     label: 'Dibatalkan' }
  };
  return map[status] ?? { bg: 'bg-gray-500/10', text: 'text-gray-400', label: status };
}

/**
 * Mendapatkan inisial nama untuk avatar
 */
export function getInitials(name: string): string {
  if (!name) return '??';
  return name
    .split(' ')
    .slice(0, 2)
    .map((n) => n[0])
    .join('')
    .toUpperCase();
}

/**
 * Nama bulan Hijriah
 */
export const HIJRI_MONTH_NAMES = [
  'Muharram', 'Safar', "Rabi'ul Awal", "Rabi'ul Akhir",
  'Jumadil Awal', 'Jumadil Akhir', 'Rajab', "Sya'ban",
  'Ramadhan', 'Syawal', "Dzulqa'dah", 'Dzulhijjah'
] as const;

export function getHijriMonthName(month: number): string {
  return HIJRI_MONTH_NAMES[month - 1] ?? '-';
}

/**
 * Label semester berdasarkan nomor semester
 */
export function getSemesterLabel(semester: number): string {
  if (semester === 1) return 'Semester 1 (Syawal - Rabi\'ul Awal)';
  if (semester === 2) return 'Semester 2 (Rabi\'ul Akhir - Sya\'ban)';
  return '-';
}

