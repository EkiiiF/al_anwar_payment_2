package repository

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type SuperUserRepository interface {
	// Dashboard
	CountActiveStudents(db *gorm.DB) (int64, error)
	CountUnpaidInvoices(db *gorm.DB) (int64, error)
	CountPaidInvoices(db *gorm.DB) (int64, error)
	SumIncomeThisMonth(db *gorm.DB, month int, year int) (float64, error)
	GetMonthlyIncomeForYear(db *gorm.DB, year int) ([]struct {
		Month int
		Total float64
	}, error)

	// Student CRUD
	FindAllStudents(db *gorm.DB, search string) ([]domain.Student, error)
	FindAllStudentsPaginated(db *gorm.DB, search string, page int, limit int) ([]domain.Student, int64, error)
	FindStudentByID(db *gorm.DB, id string) (domain.Student, error)
	CreateStudent(db *gorm.DB, student *domain.Student) error
	SaveStudent(db *gorm.DB, student *domain.Student) error
	DeleteStudent(db *gorm.DB, id string) error

	// Guardian management
	FindRoleByName(db *gorm.DB, name string) (domain.Role, error)
	CreateUser(db *gorm.DB, user *domain.User) error
	CreateGuardian(db *gorm.DB, guardian *domain.Guardian) error
	FindGuardianByStudentID(db *gorm.DB, studentID string) (domain.Guardian, error)
	FindGuardianByID(db *gorm.DB, id string) (domain.Guardian, error)
	SaveGuardian(db *gorm.DB, guardian *domain.Guardian) error
	SaveUser(db *gorm.DB, user *domain.User) error

	// Status Type CRUD
	FindAllStatusTypes(db *gorm.DB) ([]domain.StudentStatusType, error)
	FindStatusTypeByID(db *gorm.DB, id string) (domain.StudentStatusType, error)
	CreateStatusType(db *gorm.DB, statusType *domain.StudentStatusType) error
	SaveStatusType(db *gorm.DB, statusType *domain.StudentStatusType) error
	DeleteStatusType(db *gorm.DB, id string) error

	// Category CRUD
	FindAllCategories(db *gorm.DB) ([]domain.Category, error)
	FindCategoryByID(db *gorm.DB, id string) (domain.Category, error)
	CreateCategory(db *gorm.DB, category *domain.Category) error
	SaveCategory(db *gorm.DB, category *domain.Category) error
	DeleteCategory(db *gorm.DB, id string) error

	// Invoice generation
	FindStudentsForBilling(db *gorm.DB) ([]domain.Student, error)
	FindActiveCategories(db *gorm.DB) ([]domain.Category, error)
	FindCategoryByName(db *gorm.DB, name string) (domain.Category, error)
	CountInvoice(db *gorm.DB, studentID string, categoryID string, month int, year int) (int64, error)
	CountInvoiceByHijri(db *gorm.DB, studentID string, categoryID string, hijriMonth int, hijriYear int) (int64, error)
	CreateInvoice(db *gorm.DB, invoice *domain.Invoice) error
	FindAllInvoices(db *gorm.DB, status string, month int, year int) ([]domain.Invoice, error)
	FindStudentsWithInvoicesPaginated(db *gorm.DB, status string, month int, year int, page int, limit int) ([]domain.Student, int64, error)
	GetSemesterPaymentStats(db *gorm.DB, hijriYear int) ([]struct {
		Semester     int
		AcademicYear string
		Total        float64
		InvoiceCount int
		PaidCount    int
		UnpaidCount  int
	}, error)

	// User management
	FindAllUsers(db *gorm.DB) ([]domain.User, error)
	FindAllUsersPaginated(db *gorm.DB, page int, limit int) ([]domain.User, int64, error)
	FindAllRoles(db *gorm.DB) ([]domain.Role, error)
	UpdateUserRole(db *gorm.DB, userID string, roleID string) error
	ToggleUserActive(db *gorm.DB, userID string) (bool, error)

	// Invoice and log pagination
	FindAllInvoicesPaginated(db *gorm.DB, status string, month int, year int, page int, limit int) ([]domain.Invoice, int64, error)
	FindAllActivityLogsPaginated(db *gorm.DB, page int, limit int) ([]domain.ActivityLog, int64, error)
}
