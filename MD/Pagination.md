Paginator
Bernavigasi di antara beberapa halaman konten.

jenis huruf

Menyalin
import { Paginator } from '@skeletonlabs/skeleton';
Sumber
Sumber Halaman
Penggunaan
Alat peraga
Acara
Demo

jenis huruf

Menyalin
const source = [ /* any array of objects */ ];
jenis huruf

Menyalin
let paginationSettings = {
	page: 0,
	limit: 5,
	size: source.length,
	amounts: [1,2,5,10],
} satisfies PaginationSettings;
html

Menyalin
<Paginator
	bind:settings={paginationSettings}
	showFirstLastButtons="{true}"
	showPreviousNextButtons="{true}"
/>
Baris Numerik
Gunakan ini showNumeralsuntuk mengganti informasi teks dengan deretan tombol yang memungkinkan Anda menavigasi halaman dengan cepat. Kami merekomendasikan maxNumeralsjumlah yang kecil jika Anda berencana untuk mendukung perangkat seluler dengan ruang layar terbatas.


Nama	Simbol	Nomor atom
M N	Mangan	25
Rf	Rutherfordium	104
Di dalam	Flerovium	114

3 Barang


1

2

3

4

...

9

Paginasi Sisi Klien
Setelah komponen paginator Anda disiapkan, Anda perlu membagi konten sumber lokal Anda. Ini dapat dilakukan menggunakan properti reaktif Svelte yang dipasangkan dengan metode slice JavaScript .

jenis huruf

Menyalin
$: paginatedSource = source.slice(
	paginationSettings.page * paginationSettings.limit,
	paginationSettings.page * paginationSettings.limit + paginationSettings.limit
);
html

Menyalin
<ul>
	{#each paginatedSource as row}
		<li>{row}</li>
	{/each}
</ul>
Atau kombinasikan dengan komponen Tabel .

jenis huruf

Menyalin
let tableHeaders: string[] = ['Positions', 'Name', 'Weight', 'Symbol'];
html

Menyalin
<Table source={{ head: tableHeaders, body: paginatedSource }} />
Paginasi Sisi Server
Gunakan event pagedan amountuntuk memberi tahu server Anda tentang pembaruan paginasi.

jenis huruf

Menyalin
function onPageChange(e: CustomEvent): void {
	console.log('event:page', e.detail);
}

function onAmountChange(e: CustomEvent): void {
	console.log('event:amount', e.detail);
}
html

Menyalin
<Paginator ... on:page={onPageChange} on:amount={onAmountChange}></Paginator>
Menangani Reaktivitas
Gunakan teknik berikut jika Anda ingin memperbarui data paginasi secara reaktif. Pastikan untuk memperbarui paginationSettingssecara langsung, karena memperbarui referensi ke sumber tidak akan memicu reaktivitas .

jenis huruf

Menyalin
let paginationSettings = {
    page: 0,
    limit: 5,
    size: source.length,
    amounts: [1, 2, 5, 10],
} satisfies PaginationSettings;

$: paginationSettings.size = source.length;
Lihat Juga
Gunakan model berbasis data untuk membuat tabel presentasi sederhana.