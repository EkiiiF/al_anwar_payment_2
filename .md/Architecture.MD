# SYSTEM DESIGN DOCUMENT: FINANCIAL ADMINISTRATION SYSTEM

**Pondok Pesantren Al-Anwar**
_Standard Architecture, Clean Coding, and Component Specification Blueprint_

---

## 1. Clean Architecture & Design Principles

Sistem ini dirancang menggunakan prinsip **Clean Architecture** (Separation of Concerns) agar kode tetap terpisah antara HTTP, business logic, dan data persistence.

### Layering Architecture Rule

```text
┌─────────────────────────────────────────────────────────┐
│                    1. INFRASTRUCTURE                    │
│      (Fiber Framework, MySQL Drivers, Midtrans SDK)     │
└────────────────────────────┬────────────────────────────┘
                             ▼
┌─────────────────────────────────────────────────────────┐
│                       2. HANDLER                        │
│         (HTTP Request/Response, Input Validation)       │
└────────────────────────────┬────────────────────────────┘
                             ▼
┌─────────────────────────────────────────────────────────┐
│                       3. SERVICE                        │
│            (Business Logic, Transaction Management)     │
└────────────────────────────┬────────────────────────────┘
                             ▼
┌─────────────────────────────────────────────────────────┐
│                      4. REPOSITORY                      │
│                (Database Queries, Data Persistence)     │
└────────────────────────────┬────────────────────────────┘
                             ▼
┌──────────────────────────────────────────────────────────────────────┐
│                       5. DOMAIN                                      │
│               (Entities, Models, Core Interfaces, Request/Response)  │
└──────────────────────────────────────────────────────────────────────┘
```

- **Arah Dependency:** Handler bergantung pada Service, Service bergantung pada Repository, dan Repository bergantung pada Domain.
- **Interface-Driven Development:** Komunikasi antar-layer menggunakan interface eksplisit di `internal/service` dan `internal/repository`.
- **Dependency Injection:** Semua dependensi dibangun di `backend_go/internal/router/router.go` menggunakan konstruktor `New...`, tanpa variabel global.

---

## 2. Directory Structure

Struktur direktori mengikuti layout Go standar dan implementasi SvelteKit yang ada.

### Backend (Go Fiber)

```text
backend_go/
├── cmd/server/           # Entry point aplikasi (main.go)
├── db/                   # Migrasi dan seeder
│   └── seeder/           # Seeder master, role, setting, user, dll.
├── internal/
│   ├── config/           # Setup env dan database
│   ├── controller/       # Fiber HTTP handler
│   ├── exception/        # Error handler custom
│   ├── model/            # Domain entity dan web request/response models
│   ├── repository/       # Akses database
│   ├── router/           # Pendaftaran rute dan dependency injection
│   ├── service/          # Business logic dan use case
│   ├── utils/            # JWT, password, formatter, hijri
│   └── worker/           # Background jobs auto billing, cleanup
├── pkg/
│   └── middleware/       # JWT auth dan RBAC middleware
├── Dockerfile            # Container definition
```

### Frontend (SvelteKit)

```text
frontend_svelte/
├── src/
│   ├── lib/
│   │   ├── api/          # API client dan helper request
│   │   ├── assets/       # Static assets
│   │   ├── components/   # Reusable UI components
│   │   ├── context/      # Svelte context API
│   │   ├── stores/       # Global state management
│   │   ├── types/        # TypeScript type definitions
│   │   └── utils/        # Utility helpers
│   ├── routes/
│   │   ├── api/          # SvelteKit server API proxy routes
│   │   ├── dashboard/    # Dashboard pages
│   │   ├── login/        # Login page
│   │   ├── +layout.svelte
│   │   └── +page.svelte  # Landing/root page
├── Dockerfile            # Container definition
├── package.json
├── pnpm-lock.yaml
├── tsconfig.json
├── vite.config.ts
```

---

## 3. Core Database Schema & Relations

Backend menggunakan **MySQL** dengan ORM **GORM**. Relasi entitas dibangun untuk menjaga integritas transaksi, invoice, user, dan notifikasi.

---

## 4. Backend Technical Stack & Security Guideline

- **Framework:** Go Fiber v3.
- **ORM:** GORM.
- **Database:** MySQL.
- **Validation:** `go-playground/validator/v10`.
- **Auth:** JWT Bearer token dengan middleware validasi.
- **RBAC:** Role-based access control di middleware `RequireRole`.
- **CORS:** AllowOrigins dikonfigurasi dari env.
- **Logout:** Token yang logout dimasukkan blacklist, dan worker membersihkan blacklist tersebut.
- **Error Handling:** Custom error handler di `internal/exception/error_handler.go`.

---

## 5. API Endpoints

Semua API berada di bawah prefix `/api/v1`.

### 🔓 Public Routes

- `GET /` — Health Check.
- `POST /payments/notification` — Midtrans webhook callback.
- `POST /auth/login` — Login dan pembuatan JWT.

### 🔐 Authenticated Routes

- `POST /auth/logout` — Logout dan blacklist token.
- `GET /auth/profile` — Ambil profil user.
- `PUT /auth/profile` — Update profil user.
- `PUT /auth/change-password` — Ganti password.

### 👤 Role: Guardian

- `GET /guardian/dashboard` — Dashboard guardian.
- `GET /guardian/invoices` — Daftar invoice guardian.
- `GET /guardian/payments` — Riwayat pembayaran.
- `GET /guardian/notifications` — Daftar notifikasi.
- `PATCH /guardian/notifications/:id/read` — Tandai notifikasi sudah dibaca.
- `POST /payments/create` — Buat transaksi Midtrans Snap.
- `GET /payments/status/:order_id` — Cek status pembayaran.

### 👑 Role: Super User / Bendahara

- `GET /admin/dashboard` — Statistik dan ringkasan.
- `GET /admin/notifications` — Notifikasi.
- `PATCH /admin/notifications/:id/read` — Tandai notifikasi telah dibaca.
- `GET /admin/students` — Data santri.
- `GET /admin/students/paginated` — Data santri dengan pagination.
- `POST /admin/students` — Tambah santri.
- `PUT /admin/students/:id` — Update santri.
- `DELETE /admin/students/:id` — Hapus santri.
- `PATCH /admin/students/:id/toggle-status` — Toggle status aktif.
- `GET /admin/invoices` — Daftar invoice.
- `GET /admin/invoices/paginated` — Invoice dengan pagination.
- `GET /admin/invoices/students-paginated` — Invoice per santri.
- `POST /admin/invoices/generate` — Generate invoice manual.
- `POST /admin/invoices/generate-semester` — Generate invoice semester manual.
- `GET /admin/status-types` — List status tagihan.
- `POST /admin/status-types` — Buat status tagihan.
- `PUT /admin/status-types/:id` — Update status tagihan.
- `DELETE /admin/status-types/:id` — Hapus status tagihan.
- `GET /admin/categories` — List kategori.
- `POST /admin/categories` — Buat kategori.
- `PUT /admin/categories/:id` — Update kategori.
- `DELETE /admin/categories/:id` — Hapus kategori.
- `GET /admin/settings` — Ambil konfigurasi.
- `PUT /admin/settings` — Update konfigurasi.
- `GET /admin/users` — List user.
- `GET /admin/users/paginated` — User dengan pagination.
- `GET /admin/users/roles` — Daftar role user.
- `PATCH /admin/users/:id/role` — Ubah role user.
- `PATCH /admin/users/:id/toggle-active` — Toggle status aktif user.

### 👓 Role: Pengasuh (Read-only)

- `GET /pengasuh/dashboard` — Dashboard read-only.
- `GET /pengasuh/students` — Daftar santri.
- `GET /pengasuh/invoices` — Daftar invoice.

### 📊 Laporan & Auditing

- `GET /reports` — Rekap laporan.
- `GET /reports/export/csv` — Export laporan CSV.
- `GET /logs` — Activity logs.
- `GET /logs/paginated` — Activity logs pagination.

---

## 6. Frontend UI/UX Design (SvelteKit Clean Architecture)

Frontend dibangun dengan SvelteKit dan memisahkan logika dari presentasi.

- `src/lib/api/` untuk panggilan API dan helper request.
- `src/lib/components/` untuk reusable UI components.
- `src/lib/stores/` untuk state management.
- `src/lib/context/` untuk Svelte context.
- `src/lib/utils/` untuk helper format dan utilitas.
- `src/lib/types/` untuk definisi TypeScript.

Tidak ada styling inline berbasis blok `<style>` yang ditulis manual di dokumen ini; styling utama dapat dikelola melalui class utility atau CSS global.

---

## 7. System Flow & Feature Scopes

### A. Alur Billing & Kalender Hijriah

Sistem menggunakan kalender Hijriah untuk menghitung waktu billing dan mengirim notifikasi pre-billing.

- Worker otomatis memeriksa tanggal Hijriah setiap hari.
- Generate invoice Syahriyyah otomatis berjalan setelah validasi bulan Hijriah.
- Generate semester invoice berjalan pada awal semester atau bila dipicu manual.

### B. Komponen Tagihan Wajib

- **Daftar Ulang:** Dikenakan tepat di awal pembelajaran (Bulan **Syawal**) dengan nominal **Rp20.000**.
- **SPP (Syahriyah):** Dikenakan secara rutin di setiap bulan aktif pembelajaran dengan nominal **Rp40.000**.
- **Ujian Semester:** Dikenakan satu kali setiap per semester dengan nominal **Rp35.000**.

### C. Alur Transaksi Gateway

- Wali santri membuat transaksi Midtrans di endpoint `POST /payments/create`.
- Midtrans callback diproses di `POST /payments/notification`.
- Validasi signature Midtrans dilakukan di backend.
- Setelah validasi berhasil, status invoice diperbarui dan log aktivitas dicatat.

---

## 8. Deployment & Environment Container

Sistem sudah disiapkan untuk deployment container.

- **Backend:** `backend_go` menjalankan API pada port `8080`.
- **Frontend:** `frontend_svelte` menjalankan UI pada port `3000`.
- **Database:** MySQL berjalan pada port `3306`.

### Panduan Instalasi Lokal

1. Salin konfigurasi environment:
   ```bash
   cp backend_go/.env.example backend_go/.env
   cp frontend_svelte/.env.example frontend_svelte/.env
   ```
2. Jalankan Docker Compose:
   ```bash
   docker-compose up -d --build
   ```
3. Akses aplikasi:
   - Frontend: `http://localhost:3000`
   - Backend: `http://localhost:8080`

---

## 9. Background Jobs & Automation (Worker)

Backend menjalankan beberapa worker otomatis.

1. `StartTokenCleanupWorker` — membersihkan blacklist JWT setiap 1 jam.
2. `StartLogCleanupWorker` — membersihkan log aktivitas setiap 24 jam.
3. `StartAutoBillingWorker` — memeriksa kalender Hijriah setiap hari dan menjalankan:
   - auto billing Syahriyyah,
   - auto semester invoice,
   - pre-billing notification H-5.

Worker berjalan di latar belakang dan memperbarui sistem tanpa memblokir request HTTP.
