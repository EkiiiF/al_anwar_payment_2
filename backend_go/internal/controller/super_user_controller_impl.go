package controller

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/exception"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type SuperUserControllerImpl struct {
	Service service.SuperUserService
}

func NewSuperUserController(s service.SuperUserService) SuperUserController {
	return &SuperUserControllerImpl{Service: s}
}

func (ctrl *SuperUserControllerImpl) GetDashboardStats(ctx fiber.Ctx) error {
	stats, err := ctrl.Service.GetDashboardStats()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(stats))
}

func (ctrl *SuperUserControllerImpl) GetStudents(ctx fiber.Ctx) error {
	search := ctx.Query("search")
	students, err := ctrl.Service.GetStudents(search)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(students))
}

func (ctrl *SuperUserControllerImpl) GetStudentsPaginated(ctx fiber.Ctx) error {
	search := ctx.Query("search")
	page := fiber.Query[int](ctx, "page", 1)
	limit := fiber.Query[int](ctx, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	students, total, err := ctrl.Service.GetStudentsPaginated(search, page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Success(fiber.Map{
		"students": students,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	}))
}

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

func (ctrl *SuperUserControllerImpl) GetInvoicesPaginated(ctx fiber.Ctx) error {
	status := ctx.Query("status")
	month := fiber.Query[int](ctx, "month")
	year := fiber.Query[int](ctx, "year")
	page := fiber.Query[int](ctx, "page", 1)
	limit := fiber.Query[int](ctx, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	invoices, total, err := ctrl.Service.GetInvoicesPaginated(status, month, year, page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Success(fiber.Map{
		"invoices": invoices,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	}))
}

func (ctrl *SuperUserControllerImpl) GetStudentsWithInvoicesPaginated(ctx fiber.Ctx) error {
	status := ctx.Query("status")
	month := fiber.Query[int](ctx, "month")
	year := fiber.Query[int](ctx, "year")
	page := fiber.Query[int](ctx, "page", 1)
	limit := fiber.Query[int](ctx, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	students, total, err := ctrl.Service.GetStudentsWithInvoicesPaginated(status, month, year, page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Success(fiber.Map{
		"students": students,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	}))
}

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

func (ctrl *SuperUserControllerImpl) DeleteStudent(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.DeleteStudent(ctx.Context(), id, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Santri berhasil dihapus"))
}

func (ctrl *SuperUserControllerImpl) ToggleStudentStatus(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	status, err := ctrl.Service.ToggleStudentStatus(ctx.Context(), id, operatorID, ip, ua)
	if err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Status santri diubah menjadi " + status))
}

func (ctrl *SuperUserControllerImpl) GetStatusTypes(ctx fiber.Ctx) error {
	types, err := ctrl.Service.GetStatusTypes()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(types))
}

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

func (ctrl *SuperUserControllerImpl) DeleteStatusType(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.DeleteStatusType(ctx.Context(), id, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Tipe status berhasil dihapus"))
}

func (ctrl *SuperUserControllerImpl) GetCategories(ctx fiber.Ctx) error {
	categories, err := ctrl.Service.GetCategories()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(categories))
}

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

func (ctrl *SuperUserControllerImpl) DeleteCategory(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	operatorID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := ctrl.Service.DeleteCategory(ctx.Context(), id, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Kategori berhasil dihapus"))
}

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

func (ctrl *SuperUserControllerImpl) GenerateSemesterInvoices(ctx fiber.Ctx) error {
	var req struct {
		Semester  int `json:"semester"`
		HijriYear int `json:"hijri_year"`
	}
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

func (ctrl *SuperUserControllerImpl) GetSettings(ctx fiber.Ctx) error {
	settings, err := ctrl.Service.GetSettings()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(settings))
}

func (ctrl *SuperUserControllerImpl) UpdateSetting(ctx fiber.Ctx) error {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}

	operatorID, ip, ua := utils.GetOperatorInfo(ctx)
	if err := ctrl.Service.UpdateSetting(ctx.Context(), req.Key, req.Value, operatorID, ip, ua); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Setting berhasil diperbarui"))
}

func (ctrl *SuperUserControllerImpl) GetAllUsers(ctx fiber.Ctx) error {
	users, err := ctrl.Service.GetAllUsers()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(users))
}

func (ctrl *SuperUserControllerImpl) GetAllUsersPaginated(ctx fiber.Ctx) error {
	page := fiber.Query[int](ctx, "page", 1)
	limit := fiber.Query[int](ctx, "limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := ctrl.Service.GetAllUsersPaginated(page, limit)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Success(fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	}))
}

func (ctrl *SuperUserControllerImpl) GetAllRoles(ctx fiber.Ctx) error {
	roles, err := ctrl.Service.GetAllRoles()
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(roles))
}

func (ctrl *SuperUserControllerImpl) UpdateUserRole(ctx fiber.Ctx) error {
	userID := ctx.Params("id")
	var req struct {
		RoleID string `json:"role_id"`
	}
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format data tidak valid"}
	}
	if err := ctrl.Service.UpdateUserRole(ctx.Context(), userID, req.RoleID); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Role user berhasil diperbarui"))
}

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

func (ctrl *SuperUserControllerImpl) GetActivityLogsPaginated(ctx fiber.Ctx) error {
	page := fiber.Query[int](ctx, "page", 1)
	limit := fiber.Query[int](ctx, "limit", 10)
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	logs, total, err := ctrl.Service.GetActivityLogsPaginated(page, limit)
	if err != nil {
		return err
	}
	
	return ctx.JSON(response.Success(fiber.Map{
		"logs": logs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	}))
}
