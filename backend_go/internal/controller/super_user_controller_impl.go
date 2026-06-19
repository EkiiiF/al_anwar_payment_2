package controller

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/exception"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/gofiber/fiber/v3"
)

// SuperUserControllerImpl adalah implementasi SuperUserController.
// Controller ini mendelegasikan semua business logic ke service layer
// dan hanya bertanggung jawab atas HTTP request/response handling.
type SuperUserControllerImpl struct {
	Service service.SuperUserService
}

// NewSuperUserController membuat instance baru SuperUserController
// dengan dependency injection untuk service layer.
func NewSuperUserController(s service.SuperUserService) SuperUserController {
	return &SuperUserControllerImpl{Service: s}
}

// ─── Dashboard ──────────────────────────────────────────────────────────────

// GetDashboardStats mengembalikan statistik ringkasan dashboard.
func (ctrl *SuperUserControllerImpl) GetDashboardStats(ctx fiber.Ctx) error {
	year := fiber.Query[int](ctx, "year", 0)
	dateRange := ctx.Query("range")
	stats, err := ctrl.Service.GetDashboardStats(year, dateRange)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(stats))
}

// ─── Manajemen Santri ───────────────────────────────────────────────────────

// GetStudents mengembalikan seluruh data santri, dengan filter pencarian opsional.
func (ctrl *SuperUserControllerImpl) GetStudents(ctx fiber.Ctx) error {
	search := ctx.Query("search")
	students, err := ctrl.Service.GetStudents(search)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(students))
}

// GetStudentsPaginated mengembalikan data santri dengan pagination.
func (ctrl *SuperUserControllerImpl) GetStudentsPaginated(ctx fiber.Ctx) error {
	search := ctx.Query("search")
	page, limit := utils.ParsePagination(ctx)

	students, total, err := ctrl.Service.GetStudentsPaginated(search, page, limit)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Paginated("students", students, page, limit, total))
}

// CreateStudent menambahkan data santri baru ke sistem.
func (ctrl *SuperUserControllerImpl) CreateStudent(ctx fiber.Ctx) error {
	var req request.CreateStudentRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.CreateStudent(ctx.Context(), req, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.Created("Santri berhasil ditambahkan"))
}

// UpdateStudent memperbarui data santri berdasarkan ID dari URL parameter.
func (ctrl *SuperUserControllerImpl) UpdateStudent(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var req request.UpdateStudentRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.UpdateStudent(ctx.Context(), id, req, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Data santri berhasil diperbarui"))
}

// DeleteStudent menghapus data santri berdasarkan ID dari URL parameter.
func (ctrl *SuperUserControllerImpl) DeleteStudent(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.DeleteStudent(ctx.Context(), id, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Santri berhasil dihapus"))
}

// ToggleStudentStatus mengubah status aktif/non-aktif santri berdasarkan ID.
func (ctrl *SuperUserControllerImpl) ToggleStudentStatus(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	status, err := ctrl.Service.ToggleStudentStatus(ctx.Context(), id, operatorID, ip, ua)
	if err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Status santri diubah menjadi " + status))
}

// ─── Manajemen Invoice ──────────────────────────────────────────────────────

// GetInvoices mengembalikan seluruh data invoice dengan filter opsional.
func (ctrl *SuperUserControllerImpl) GetInvoices(ctx fiber.Ctx) error {
	status := ctx.Query("status")
	month := fiber.Query[int](ctx, "month")
	year := fiber.Query[int](ctx, "year")

	invoices, err := ctrl.Service.GetInvoices(status, month, year)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(invoices))
}

// GetInvoicesPaginated mengembalikan data invoice dengan pagination dan filter.
func (ctrl *SuperUserControllerImpl) GetInvoicesPaginated(ctx fiber.Ctx) error {
	status := ctx.Query("status")
	month := fiber.Query[int](ctx, "month")
	year := fiber.Query[int](ctx, "year")
	page, limit := utils.ParsePagination(ctx)

	invoices, total, err := ctrl.Service.GetInvoicesPaginated(status, month, year, page, limit)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Paginated("invoices", invoices, page, limit, total))
}

// GetStudentsWithInvoicesPaginated mengembalikan data santri beserta invoice terkait.
func (ctrl *SuperUserControllerImpl) GetStudentsWithInvoicesPaginated(ctx fiber.Ctx) error {
	status := ctx.Query("status")
	month := fiber.Query[int](ctx, "month")
	year := fiber.Query[int](ctx, "year")
	search := ctx.Query("search")
	page, limit := utils.ParsePagination(ctx)

	students, total, err := ctrl.Service.GetStudentsWithInvoicesPaginated(status, month, year, search, page, limit)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Paginated("students", students, page, limit, total))
}

// GenerateInvoices men-generate invoice bulanan berdasarkan bulan dan tahun.
func (ctrl *SuperUserControllerImpl) GenerateInvoices(ctx fiber.Ctx) error {
	var req request.GenerateInvoiceRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	count, err := ctrl.Service.GenerateInvoices(ctx.Context(), req, operatorID, ip, ua)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(fiber.Map{"count": count}))
}

// GenerateSemesterInvoices men-generate invoice ujian semester secara manual.
func (ctrl *SuperUserControllerImpl) GenerateSemesterInvoices(ctx fiber.Ctx) error {
	var req request.GenerateSemesterInvoiceRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	count, err := ctrl.Service.GenerateSemesterInvoices(ctx.Context(), req.Semester, req.HijriYear, operatorID, ip, ua)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(fiber.Map{"count": count}))
}

// ─── Manajemen Status Tagihan ───────────────────────────────────────────────

// GetStatusTypes mengembalikan seluruh tipe status tagihan.
func (ctrl *SuperUserControllerImpl) GetStatusTypes(ctx fiber.Ctx) error {
	types, err := ctrl.Service.GetStatusTypes()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(types))
}

// CreateStatusType membuat tipe status tagihan baru.
func (ctrl *SuperUserControllerImpl) CreateStatusType(ctx fiber.Ctx) error {
	var req request.CreateStatusTypeRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	res, err := ctrl.Service.CreateStatusType(ctx.Context(), req, operatorID, ip, ua)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.Created(res))
}

// UpdateStatusType memperbarui tipe status tagihan berdasarkan ID.
func (ctrl *SuperUserControllerImpl) UpdateStatusType(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var req request.CreateStatusTypeRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	res, err := ctrl.Service.UpdateStatusType(ctx.Context(), id, req, operatorID, ip, ua)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(res))
}

// DeleteStatusType menghapus tipe status tagihan berdasarkan ID.
func (ctrl *SuperUserControllerImpl) DeleteStatusType(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.DeleteStatusType(ctx.Context(), id, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Tipe status berhasil dihapus"))
}

// ─── Manajemen Kategori ─────────────────────────────────────────────────────

// GetCategories mengembalikan seluruh data kategori tagihan.
func (ctrl *SuperUserControllerImpl) GetCategories(ctx fiber.Ctx) error {
	categories, err := ctrl.Service.GetCategories()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(categories))
}

// CreateCategory membuat kategori tagihan baru.
func (ctrl *SuperUserControllerImpl) CreateCategory(ctx fiber.Ctx) error {
	var req request.CreateCategoryRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	res, err := ctrl.Service.CreateCategory(ctx.Context(), req, operatorID, ip, ua)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(response.Created(res))
}

// UpdateCategory memperbarui kategori tagihan berdasarkan ID.
func (ctrl *SuperUserControllerImpl) UpdateCategory(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var req request.CreateCategoryRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	res, err := ctrl.Service.UpdateCategory(ctx.Context(), id, req, operatorID, ip, ua)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(res))
}

// DeleteCategory menghapus kategori tagihan berdasarkan ID.
func (ctrl *SuperUserControllerImpl) DeleteCategory(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.DeleteCategory(ctx.Context(), id, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Kategori berhasil dihapus"))
}

// ─── Pengaturan Sistem ──────────────────────────────────────────────────────

// GetSettings mengembalikan seluruh konfigurasi sistem.
func (ctrl *SuperUserControllerImpl) GetSettings(ctx fiber.Ctx) error {
	settings, err := ctrl.Service.GetSettings()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(settings))
}

// UpdateSetting memperbarui satu konfigurasi sistem berdasarkan key.
func (ctrl *SuperUserControllerImpl) UpdateSetting(ctx fiber.Ctx) error {
	var req request.UpdateSettingRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.UpdateSetting(ctx.Context(), req.Key, req.Value, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Setting berhasil diperbarui"))
}

// ─── Manajemen User ─────────────────────────────────────────────────────────

// GetAllUsers mengembalikan seluruh data user.
func (ctrl *SuperUserControllerImpl) GetAllUsers(ctx fiber.Ctx) error {
	users, err := ctrl.Service.GetAllUsers()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(users))
}

// GetAllUsersPaginated mengembalikan data user dengan pagination.
func (ctrl *SuperUserControllerImpl) GetAllUsersPaginated(ctx fiber.Ctx) error {
	page, limit := utils.ParsePagination(ctx)

	users, total, err := ctrl.Service.GetAllUsersPaginated(page, limit)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Paginated("users", users, page, limit, total))
}

// GetAllRoles mengembalikan seluruh role yang tersedia dalam sistem.
func (ctrl *SuperUserControllerImpl) GetAllRoles(ctx fiber.Ctx) error {
	roles, err := ctrl.Service.GetAllRoles()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(roles))
}

// UpdateUserRole mengubah role user berdasarkan ID user.
func (ctrl *SuperUserControllerImpl) UpdateUserRole(ctx fiber.Ctx) error {
	userID := ctx.Params("id")
	var req request.UpdateUserRoleRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	if err := ctrl.Service.UpdateUserRole(ctx.Context(), userID, req.RoleID); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Role user berhasil diperbarui"))
}

// ToggleUserActive mengubah status aktif/non-aktif user berdasarkan ID.
func (ctrl *SuperUserControllerImpl) ToggleUserActive(ctx fiber.Ctx) error {
	userID := ctx.Params("id")

	active, err := ctrl.Service.ToggleUserActive(ctx.Context(), userID)
	if err != nil {
		return err
	}

	status := "Nonaktif"
	if active {
		status = "Aktif"
	}
	return ctx.JSON(response.SuccessMessage("Status user diubah menjadi " + status))
}

// ─── Activity Logs ──────────────────────────────────────────────────────────

// GetActivityLogsPaginated mengembalikan log aktivitas sistem dengan pagination.
func (ctrl *SuperUserControllerImpl) GetActivityLogsPaginated(ctx fiber.Ctx) error {
	page, limit := utils.ParsePagination(ctx)

	logs, total, err := ctrl.Service.GetActivityLogsPaginated(page, limit)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Paginated("logs", logs, page, limit, total))
}
