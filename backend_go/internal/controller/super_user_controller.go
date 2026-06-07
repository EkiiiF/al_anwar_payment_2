package controller

import "github.com/gofiber/fiber/v3"

type SuperUserController interface {
	GetDashboardStats(ctx fiber.Ctx) error
	GetStudents(ctx fiber.Ctx) error
	GetStudentsPaginated(ctx fiber.Ctx) error
	GetInvoices(ctx fiber.Ctx) error
	GetStudentsWithInvoicesPaginated(ctx fiber.Ctx) error
	GetInvoicesPaginated(ctx fiber.Ctx) error
	CreateStudent(ctx fiber.Ctx) error
	UpdateStudent(ctx fiber.Ctx) error
	DeleteStudent(ctx fiber.Ctx) error
	ToggleStudentStatus(ctx fiber.Ctx) error
	GetStatusTypes(ctx fiber.Ctx) error
	CreateStatusType(ctx fiber.Ctx) error
	UpdateStatusType(ctx fiber.Ctx) error
	DeleteStatusType(ctx fiber.Ctx) error
	GetCategories(ctx fiber.Ctx) error
	CreateCategory(ctx fiber.Ctx) error
	UpdateCategory(ctx fiber.Ctx) error
	DeleteCategory(ctx fiber.Ctx) error
	GenerateInvoices(ctx fiber.Ctx) error
	GenerateSemesterInvoices(ctx fiber.Ctx) error
	GetSettings(ctx fiber.Ctx) error
	UpdateSetting(ctx fiber.Ctx) error
	GetAllUsers(ctx fiber.Ctx) error
	GetAllUsersPaginated(ctx fiber.Ctx) error
	GetAllRoles(ctx fiber.Ctx) error
	UpdateUserRole(ctx fiber.Ctx) error
	ToggleUserActive(ctx fiber.Ctx) error
	GetActivityLogsPaginated(ctx fiber.Ctx) error
}
