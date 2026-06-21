# PACKAGE MANAGER

pnpm

# PEDOMAN DESAIN & KONSISTENSI UI/UX FRONTEND

Dokumen ini berfungsi sebagai acuan utama (single source of truth) untuk menjaga konsistensi visual, kegunaan, responsivitas, dan standarisasi antarmuka pengguna pada sistem pembayaran Al-Anwar.

---

## 1. Panduan Gaya Visual Global (Global Style Guide)

Seluruh halaman dan komponen dalam `frontend_svelte` wajib menggunakan token gaya yang seragam di bawah ini untuk menjaga estetika premium, modern, dan profesional. PNPM sebagai package manager. Tailwind sebagai style framework. SvelteKit sebagai framework. (SvelteKit + Tailwind + TypeScript)

### 1.1 Tipografi & Font

- **Font Utama**: `Plus Jakarta Sans` (diimpor dari Google Fonts).
- **Font alternatif**: `sans-serif` (system default fallbacks).
- **Ketebalan Font (Font Weights)**:
  - `font-normal` (400) — Untuk teks isi (body copy).
  - `font-medium` (500) — Untuk teks sekunder, label tombol sedang.
  - `font-semibold` (600) — Untuk sub-judul kecil dan label input.
  - `font-bold` (700) — Untuk header kartu, judul modal, dan tombol utama.
  - `font-black` (900) — Untuk judul halaman utama (main page title).

### 1.2 Palet Warna & Kode Token

Aplikasi menggunakan skema warna terkurasi berdasarkan peran (role) pengguna dan status transaksi:

| Kategori Warna             | Kode Hex  | Class Tailwind                     | Penjelasan Penggunaan                                       |
| :------------------------- | :-------- | :--------------------------------- | :---------------------------------------------------------- |
| **Tema Utama (Emerald)**   | `#065F46` | `bg-emerald-800`                   | Identitas Pondok Al-Anwar, dipakai di tombol utama & header |
| **Aksen Utama**            | `#10B981` | `text-emerald-500`                 | Indikator sukses, link aktif, toggle aktif                  |
| **Super User (Bendahara)** | `#15803d` | `text-green-700` / `bg-green-50`   | Skema warna modul keuangan dan manajemen data santri        |
| **Wali Santri (Guardian)** | `#1d4ed8` | `text-blue-700` / `bg-blue-50`     | Skema warna tagihan wali santri & gateway pembayaran        |
| **Pengasuh**               | `#6b21a8` | `text-purple-800` / `bg-purple-50` | Skema warna laporan perkembangan & data akademik            |

#### Status Pembayaran & Transaksi (Semantic Colors)

- **Paid / Lunas**: `text-emerald-600` / `bg-emerald-50` / `border-emerald-200`
- **Unpaid / Belum Bayar**: `text-amber-600` / `bg-amber-50` / `border-amber-200`
- **Pending / Menunggu**: `text-blue-600` / `bg-blue-50` / `border-blue-200`
- **Cancelled / Batal**: `text-red-600` / `bg-red-50` / `border-red-200`
- **Expired / Kedaluwarsa**: `text-gray-500` / `bg-gray-50` / `border-gray-200`

---

## 2. Struktur Tata Letak & Responsivitas (Layout & Responsiveness)

Desain antarmuka dirancang menggunakan prinsip _mobile-first_ agar tetap nyaman diakses dari smartphone wali santri maupun komputer desktop bendahara.

### 2.1 Grid dan Spacing Konsisten

- **Margin Halaman**: Gunakan `px-4 sm:px-6 lg:px-8` dan `py-6`.
- **Jarak Antar Baris**: Gunakan `space-y-4` atau `space-y-6`.
- **Form Columns Grid**:
  ```html
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <!-- Kolom Kiri -->
    <!-- Kolom Kanan -->
  </div>
  ```

### 2.2 Desain Tabel Responsif (`DataTable`)

Tabel data besar tidak boleh meluber pada layar ponsel. Gunakan trik kolom tersembunyi dengan breakpoints:

- Kolom kritis (Nama, NIS): selalu tampil (`table-cell`).
- Kolom sekunder (NIK, Gender): tampilkan hanya di layar tablet/desktop (`hidden md:table-cell`).
- Kolom tersier (Tanggal Daftar, Alamat): tampilkan hanya di layar besar (`hidden lg:table-cell`).
- Detail lengkap dapat dilihat melalui perluasan baris (`toggleExpand`).

### 2.3 Sidebar Navigasi Adaptif

- **Desktop**: Sidebar menetap di sisi kiri dengan lebar tetap (`w-64`) dan transisi yang halus.
- **Mobile**: Tombol menu hamburger di pojok atas akan memicu munculnya drawer panel overlay dengan transisi geser samping.

---

## 3. Katalog Komponen Dasar Reusable (Base Components)

Pola interaksi UI harus seragam di setiap halaman. Jangan membuat styling ad-hoc, gunakan komponen dasar di `src/lib/components/ui/`:

### 3.1 Button Component (`Button.svelte`)

Menyediakan variasi button dengan parameter pemuatan (loading state):

- **Varian**:
  - `primary` (hijau emerald gelap untuk simpan/daftar).
  - `secondary` (abu-abu/slate untuk aksi penunjang).
  - `outline` (border transparan untuk batal/tutup).
  - `danger` (merah untuk hapus/nonaktifkan).
- **Properti Wajib**: `loading` (spinner saat submit API), `disabled` (saat form tidak valid).

### 3.2 Input Component (`Input.svelte`)

Input field seragam dengan indikasi pesan kesalahan (error state):

- **Error State**: Jika properti `error` berisi pesan teks, border akan berganti warna merah (`border-red-400 focus:ring-red-200`) dan memunculkan teks deskripsi error di bawah input.
- **Helper Text**: Petunjuk tambahan berwarna abu-abu yang diletakkan di bawah input.

### 3.3 Card Component (`Card.svelte`)

Wadah/panel visual dengan shadow minimalis:

- **Style**: `bg-white border border-gray-200 rounded-2xl shadow-sm p-5 sm:p-6 transition-all duration-200`.

### 3.4 Modal Component (`Modal.svelte`)

Popup interaktif untuk formulir CRUD atau konfirmasi:

- **Style**: Latar belakang overlay buram (`backdrop-blur-sm bg-black/40`), transisi skala melompat lembut, dan memiliki tombol close di header kanan atas.
- **Ukuran**: `sm` (konfirmasi), `md` (form sedang), `lg` (form detail), `xl` (form santri baru).

---

## 4. Penamaan Deskripsi Fitur & Hak Akses (Role-Based Nomenclature)

Untuk menghindari kebingungan saat UAT, penamaan fitur dan deskripsi hak akses harus seragam:

1. **Super User (Bendahara)**
   - Deskripsi: Pengelola keuangan utama pondok.
   - Fitur: _Pendaftaran Santri_, _Pembayaran Tagihan Manual_, _Manajemen Kategori Biaya_, _Log Aktivitas_, _Laporan Keuangan_, _Pengaturan Auto-Billing_.
2. **Wali Santri (Guardian)**
   - Deskripsi: Orang tua / wali penanggung jawab biaya santri.
   - Fitur: _Portal Tagihan Anak_, _Gateway Pembayaran Instan (Midtrans)_, _Riwayat Transaksi_, _Notifikasi WA/Email_.
3. **Pengasuh**
   - Deskripsi: Dewan asatidz / pengasuh pondok pesantren.
   - Fitur: _Monitoring Status Keaktifan Santri_, _Statistik Distribusi Kelas_, _Laporan Perkembangan Santri_.

---

## 5. Sistem Kalender Islam (Bulan Hijriah)

Sebagai lembaga pendidikan Islam, siklus penagihan Al-Anwar berbasis kalender Hijriah. Penulisan nama bulan Islam diatur dengan urutan indeks 1 sampai 12 sebagai berikut:

| Indeks Bulan | Nama Bulan Hijriah | Keterangan Semester   |
| :----------: | :----------------- | :-------------------- |
|    **1**     | Muharram           | Semester 1 (Ganjil)   |
|    **2**     | Safar              | Semester 1 (Ganjil)   |
|    **3**     | Rabi'ul Awal       | Semester 1 (Ganjil)   |
|    **4**     | Rabi'ul Akhir      | Semester 2 (Genap)    |
|    **5**     | Jumadil Awal       | Semester 2 (Genap)    |
|    **6**     | Jumadil Akhir      | Semester 2 (Genap)    |
|    **7**     | Rajab              | Semester 2 (Genap)    |
|    **8**     | Sya'ban            | Semester 2 (Genap)    |
|    **9**     | Ramadhan           | Libur Akademik Pondok |
|    **10**    | Syawal             | Semester 1 (Ganjil)   |
|    **11**    | Dzulqa'dah         | Semester 1 (Ganjil)   |
|    **12**    | Dzulhijjah         | Semester 1 (Ganjil)   |

### Pengelompokan Siklus Semester

- **Semester 1 (Ganjil)**: Meliputi bulan **Syawal** sampai dengan **Rabi'ul Awal** tahun berikutnya.
- **Semester 2 (Genap)**: Meliputi bulan **Rabi'ul Akhir** sampai dengan **Sya'ban**.

---

## 6. Standar Prosedur Validasi Form & UAT Checklist

Setiap form input wajib mengimplementasikan validasi sisi klien (client-side validation) dengan kriteria:

1. **Pengecekan Wajib Isi (Required Fields)**: Memeriksa karakter kosong setelah _trim_. Tampilkan error spesifik per kolom, hindari pesan error umum.
2. **Pembersihan Pesan Error Terjadwal**: Gunakan _handler_ `oninput` pada `<Input>` untuk menghapus pesan kesalahan seketika setelah pengguna mulai mengetik ulang input tersebut.
3. **Validasi Tipe Data & Range**:
   - NIK harus bernilai tepat 16 digit angka.
   - NIS berkisar antara 5 hingga 15 digit angka.
   - Nomor HP wali harus berupa format telepon internasional/nasional yang valid (9 hingga 15 digit).
   - Angka nominal master data (biaya) harus bernilai positif (> 0).
   - Persentase diskon status santri harus berkisar antara 0% hingga 100%.
   - Tahun Hijriah pada Generator Tagihan harus berada pada rentang rasional yaitu 1300 H sampai dengan 1600 H.
