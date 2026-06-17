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
export function getInvoiceStatusStyle(status: string): { bg: string; text: string; border: string; label: string } {
  const map: Record<string, { bg: string; text: string; border: string; label: string }> = {
    paid:      { bg: 'bg-emerald-50',  text: 'text-emerald-600', border: 'border-emerald-200', label: 'Lunas' },
    unpaid:    { bg: 'bg-amber-50',    text: 'text-amber-600',   border: 'border-amber-200',   label: 'Belum Bayar' },
    pending:   { bg: 'bg-blue-50',     text: 'text-blue-600',    border: 'border-blue-200',    label: 'Menunggu' },
    cancelled: { bg: 'bg-red-50',      text: 'text-red-600',     border: 'border-red-200',     label: 'Batal' },
    expired:   { bg: 'bg-gray-50',     text: 'text-gray-500',    border: 'border-gray-200',    label: 'Kedaluwarsa' }
  };
  return map[status] ?? { bg: 'bg-gray-50', text: 'text-gray-500', border: 'border-gray-200', label: status };
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

