// Package controller berisi HTTP handler untuk seluruh endpoint API.
// Controller bertanggung jawab atas parsing request, validasi input dasar,
// dan memanggil layer service untuk eksekusi business logic.
package controller

import "github.com/gofiber/fiber/v3"

// SuperUserController mendefinisikan kontrak HTTP handler untuk fitur admin/super user.
// Endpoint yang tersedia mencakup manajemen santri, invoice, kategori,
// status tagihan, pengaturan sistem, dan manajemen user.
type SuperUserController interface {
	// ─── Dashboard ──────────────────────────────────────────

	// GetDashboardStats mengembalikan statistik ringkasan dashboard admin.
	// Query: ?year=<int> (opsional, filter berdasarkan tahun).
	GetDashboardStats(ctx fiber.Ctx) error

	// ─── Manajemen Santri ───────────────────────────────────

	// GetStudents mengembalikan seluruh data santri.
	// Query: ?search=<string> (opsional, pencarian berdasarkan nama).
	GetStudents(ctx fiber.Ctx) error

	// GetStudentsPaginated mengembalikan data santri dengan pagination.
	// Query: ?search=<string>&page=<int>&limit=<int>.
	GetStudentsPaginated(ctx fiber.Ctx) error

	// CreateStudent menambahkan data santri baru.
	// Body: request.CreateStudentRequest.
	CreateStudent(ctx fiber.Ctx) error

	// UpdateStudent memperbarui data santri berdasarkan ID.
	// Params: :id. Body: request.UpdateStudentRequest.
	UpdateStudent(ctx fiber.Ctx) error

	// DeleteStudent menghapus data santri berdasarkan ID.
	// Params: :id.
	DeleteStudent(ctx fiber.Ctx) error

	// ToggleStudentStatus mengubah status aktif/non-aktif santri.
	// Params: :id.
	ToggleStudentStatus(ctx fiber.Ctx) error

	// ─── Manajemen Invoice ──────────────────────────────────

	// GetInvoices mengembalikan seluruh data invoice.
	// Query: ?status=<string>&month=<int>&year=<int>.
	GetInvoices(ctx fiber.Ctx) error

	// GetInvoicesPaginated mengembalikan data invoice dengan pagination.
	// Query: ?status=<string>&month=<int>&year=<int>&page=<int>&limit=<int>.
	GetInvoicesPaginated(ctx fiber.Ctx) error

	// GetStudentsWithInvoicesPaginated mengembalikan data santri beserta invoice terkait.
	// Query: ?status=<string>&month=<int>&year=<int>&page=<int>&limit=<int>.
	GetStudentsWithInvoicesPaginated(ctx fiber.Ctx) error

	// GenerateInvoices men-generate invoice bulanan secara manual.
	// Body: request.GenerateInvoiceRequest.
	GenerateInvoices(ctx fiber.Ctx) error

	// GenerateSemesterInvoices men-generate invoice ujian semester secara manual.
	// Body: request.GenerateSemesterInvoiceRequest.
	GenerateSemesterInvoices(ctx fiber.Ctx) error

	// ─── Manajemen Status Tagihan ───────────────────────────

	// GetStatusTypes mengembalikan seluruh tipe status tagihan.
	GetStatusTypes(ctx fiber.Ctx) error

	// CreateStatusType membuat tipe status tagihan baru.
	// Body: request.CreateStatusTypeRequest.
	CreateStatusType(ctx fiber.Ctx) error

	// UpdateStatusType memperbarui tipe status tagihan berdasarkan ID.
	// Params: :id. Body: request.CreateStatusTypeRequest.
	UpdateStatusType(ctx fiber.Ctx) error

	// DeleteStatusType menghapus tipe status tagihan berdasarkan ID.
	// Params: :id.
	DeleteStatusType(ctx fiber.Ctx) error

	// ─── Manajemen Kategori ─────────────────────────────────

	// GetCategories mengembalikan seluruh data kategori tagihan.
	GetCategories(ctx fiber.Ctx) error

	// CreateCategory membuat kategori tagihan baru.
	// Body: request.CreateCategoryRequest.
	CreateCategory(ctx fiber.Ctx) error

	// UpdateCategory memperbarui kategori tagihan berdasarkan ID.
	// Params: :id. Body: request.CreateCategoryRequest.
	UpdateCategory(ctx fiber.Ctx) error

	// DeleteCategory menghapus kategori tagihan berdasarkan ID.
	// Params: :id.
	DeleteCategory(ctx fiber.Ctx) error

	// ─── Pengaturan Sistem ──────────────────────────────────

	// GetSettings mengembalikan seluruh konfigurasi sistem.
	GetSettings(ctx fiber.Ctx) error

	// UpdateSetting memperbarui satu konfigurasi sistem.
	// Body: request.UpdateSettingRequest.
	UpdateSetting(ctx fiber.Ctx) error

	// ─── Manajemen User ─────────────────────────────────────

	// GetAllUsers mengembalikan seluruh data user.
	GetAllUsers(ctx fiber.Ctx) error

	// GetAllUsersPaginated mengembalikan data user dengan pagination.
	// Query: ?page=<int>&limit=<int>.
	GetAllUsersPaginated(ctx fiber.Ctx) error

	// GetAllRoles mengembalikan seluruh role yang tersedia.
	GetAllRoles(ctx fiber.Ctx) error

	// UpdateUserRole mengubah role user berdasarkan ID.
	// Params: :id. Body: request.UpdateUserRoleRequest.
	UpdateUserRole(ctx fiber.Ctx) error

	// ToggleUserActive mengubah status aktif/non-aktif user.
	// Params: :id.
	ToggleUserActive(ctx fiber.Ctx) error

	// ─── Activity Logs ──────────────────────────────────────

	// GetActivityLogsPaginated mengembalikan log aktivitas dengan pagination.
	// Query: ?page=<int>&limit=<int>.
	GetActivityLogsPaginated(ctx fiber.Ctx) error
}
