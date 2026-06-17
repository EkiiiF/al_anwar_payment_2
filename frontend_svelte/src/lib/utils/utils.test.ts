import { describe, it, expect } from 'vitest';
import {
  formatRupiah,
  formatDate,
  getMonthName,
  getInvoiceStatusStyle,
  getInitials,
  MONTH_NAMES,
  getHijriMonthName,
  HIJRI_MONTH_NAMES
} from '$lib/utils';

// ── formatRupiah ──────────────────────────────────────────
describe('formatRupiah', () => {
  it('✅ memformat angka positif menjadi Rupiah', () => {
    const result = formatRupiah(100000);
    expect(result).toContain('100.000');
    expect(result).toContain('Rp');
  });

  it('✅ memformat nol menjadi Rp 0', () => {
    expect(formatRupiah(0)).toContain('0');
  });

  it('✅ memformat nilai besar (jutaan)', () => {
    const result = formatRupiah(1500000);
    expect(result).toContain('1.500.000');
  });

  it('✅ memformat nilai desimal tanpa digit desimal', () => {
    // maximumFractionDigits: 0 → tidak ada desimal
    expect(formatRupiah(99999.99)).not.toContain(',');
  });

  it('✅ memformat nilai negatif', () => {
    const result = formatRupiah(-50000);
    expect(result).toContain('50.000');
  });
});

// ── formatDate ────────────────────────────────────────────
describe('formatDate', () => {
  it('✅ memformat tanggal valid menjadi bahasa Indonesia', () => {
    const result = formatDate('2024-01-15');
    expect(result).toContain('2024');
    expect(result).toMatch(/Januari|January/); // fallback
  });

  it('❌ mengembalikan - jika null/undefined', () => {
    expect(formatDate(null)).toBe('-');
    expect(formatDate(undefined)).toBe('-');
    expect(formatDate('')).toBe('-');
  });

  it('❌ mengembalikan - jika string bukan tanggal', () => {
    expect(formatDate('bukan-tanggal')).toBe('-');
  });

  it('✅ menyertakan waktu jika withTime = true', () => {
    const result = formatDate('2024-06-15T10:30:00Z', true);
    expect(result).toBeTruthy();
    expect(result).not.toBe('-');
  });
});

// ── getMonthName ──────────────────────────────────────────
describe('getMonthName', () => {
  it('✅ mengembalikan nama bulan yang benar', () => {
    expect(getMonthName(1)).toBe('Januari');
    expect(getMonthName(6)).toBe('Juni');
    expect(getMonthName(12)).toBe('Desember');
  });

  it('❌ mengembalikan - untuk bulan invalid', () => {
    expect(getMonthName(0)).toBe('-');
    expect(getMonthName(13)).toBe('-');
  });

  it('✅ MONTH_NAMES memiliki 12 elemen', () => {
    expect(MONTH_NAMES.length).toBe(12);
  });
});

// ── getInvoiceStatusStyle ─────────────────────────────────
describe('getInvoiceStatusStyle', () => {
  it('✅ mengembalikan style untuk status unpaid', () => {
    const s = getInvoiceStatusStyle('unpaid');
    expect(s.label).toBe('Belum Bayar');
    expect(s.text).toContain('amber');
  });

  it('✅ mengembalikan style untuk status paid', () => {
    const s = getInvoiceStatusStyle('paid');
    expect(s.label).toBe('Lunas');
    expect(s.text).toContain('emerald');
  });

  it('✅ mengembalikan style default untuk status tidak dikenal', () => {
    const s = getInvoiceStatusStyle('unknown_status');
    expect(s.label).toBe('unknown_status');
    expect(s.text).toContain('gray');
  });

  const allStatuses = ['unpaid', 'pending', 'paid', 'expired', 'cancelled'];
  allStatuses.forEach(status => {
    it(`✅ memiliki style untuk status: ${status}`, () => {
      const s = getInvoiceStatusStyle(status);
      expect(s.bg).toBeTruthy();
      expect(s.text).toBeTruthy();
      expect(s.label).toBeTruthy();
    });
  });
});

// ── getInitials ───────────────────────────────────────────
describe('getInitials', () => {
  it('✅ mengambil inisial dari nama lengkap', () => {
    expect(getInitials('Ahmad Fauzi')).toBe('AF');
  });

  it('✅ mengambil inisial dari satu kata', () => {
    expect(getInitials('Ahmad')).toBe('A');
  });

  it('✅ mengambil maksimal 2 inisial', () => {
    expect(getInitials('Ahmad Fauzi Ramadhan')).toBe('AF');
  });

  it('❌ mengembalikan ?? jika nama kosong', () => {
    expect(getInitials('')).toBe('??');
  });
});

// ── getHijriMonthName ─────────────────────────────────────
describe('getHijriMonthName', () => {
  it('✅ mengembalikan nama bulan Hijriah yang benar', () => {
    expect(getHijriMonthName(1)).toBe('Muharram');
    expect(getHijriMonthName(10)).toBe('Syawal');
    expect(getHijriMonthName(12)).toBe('Dzulhijjah');
  });

  it('❌ mengembalikan - untuk bulan invalid', () => {
    expect(getHijriMonthName(0)).toBe('-');
    expect(getHijriMonthName(13)).toBe('-');
  });

  it('✅ HIJRI_MONTH_NAMES memiliki 12 elemen', () => {
    expect(HIJRI_MONTH_NAMES.length).toBe(12);
  });
});
