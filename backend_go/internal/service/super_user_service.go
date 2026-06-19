package service

import (
	"context"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
)

// SuperUserService — kontrak business logic admin/super user.
type SuperUserService interface {
	GetDashboardStats(year int, dateRange string) (response.SuperUserDashboardStats, error)
	GetStudents(search string) ([]domain.Student, error)
	GetStudentsPaginated(search string, page int, limit int) ([]domain.Student, int64, error)
	GetInvoices(status string, month int, year int) ([]domain.Invoice, error)
	GetStudentsWithInvoicesPaginated(status string, month int, year int, search string, page int, limit int) ([]domain.Student, int64, error)
	CreateStudent(ctx context.Context, req request.CreateStudentRequest, operatorID, ip, ua string) error
	UpdateStudent(ctx context.Context, id string, req request.UpdateStudentRequest, operatorID, ip, ua string) error
	DeleteStudent(ctx context.Context, id string, operatorID, ip, ua string) error
	ToggleStudentStatus(ctx context.Context, id string, operatorID, ip, ua string) (string, error)
	GetStatusTypes() ([]domain.StudentStatusType, error)
	CreateStatusType(ctx context.Context, req request.CreateStatusTypeRequest, operatorID, ip, ua string) (domain.StudentStatusType, error)
	UpdateStatusType(ctx context.Context, id string, req request.CreateStatusTypeRequest, operatorID, ip, ua string) (domain.StudentStatusType, error)
	DeleteStatusType(ctx context.Context, id string, operatorID, ip, ua string) error
	GetCategories() ([]domain.Category, error)
	CreateCategory(ctx context.Context, req request.CreateCategoryRequest, operatorID, ip, ua string) (domain.Category, error)
	UpdateCategory(ctx context.Context, id string, req request.CreateCategoryRequest, operatorID, ip, ua string) (domain.Category, error)
	DeleteCategory(ctx context.Context, id string, operatorID, ip, ua string) error
	GenerateInvoices(ctx context.Context, req request.GenerateInvoiceRequest, operatorID, ip, ua string) (int, error)
	GenerateSemesterInvoices(ctx context.Context, semester int, hijriYear int, operatorID, ip, ua string) (int, error)
	PromoteStudents(ctx context.Context, targetSemester int, targetHijriYear int) error
	SendPreBillingNotifications(ctx context.Context, nextMonth int, nextYear int, targetDay int) error
	GetSettings() ([]domain.Setting, error)
	UpdateSetting(ctx context.Context, key, value, operatorID, ip, ua string) error
	GetAllUsers() ([]domain.User, error)
	GetAllUsersPaginated(page int, limit int) ([]domain.User, int64, error)
	GetAllRoles() ([]domain.Role, error)
	UpdateUserRole(ctx context.Context, userID, roleID string) error
	ToggleUserActive(ctx context.Context, userID string) (bool, error)
	GetInvoicesPaginated(status string, month int, year int, page int, limit int) ([]domain.Invoice, int64, error)
	GetActivityLogsPaginated(page int, limit int) ([]domain.ActivityLog, int64, error)
}
