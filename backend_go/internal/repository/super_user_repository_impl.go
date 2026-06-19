package repository

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"gorm.io/gorm"
)

type SuperUserRepositoryImpl struct{}

func NewSuperUserRepository() SuperUserRepository {
	return &SuperUserRepositoryImpl{}
}

func (r *SuperUserRepositoryImpl) CountActiveStudents(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&domain.Student{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

func (r *SuperUserRepositoryImpl) CountUnpaidInvoices(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&domain.Invoice{}).Where("status = ?", "unpaid").Count(&count).Error
	return count, err
}

func (r *SuperUserRepositoryImpl) CountPaidInvoices(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&domain.Invoice{}).Where("status = ?", "paid").Count(&count).Error
	return count, err
}

func (r *SuperUserRepositoryImpl) SumIncomeThisMonth(db *gorm.DB, month int, year int) (float64, error) {
	var endYear, endMonth int
	if month == 12 {
		endYear = year + 1
		endMonth = 1
	} else {
		endYear = year
		endMonth = month + 1
	}
	startDate := utils.HijriToGregorian(year, month, 1)
	endDate := utils.HijriToGregorian(endYear, endMonth, 1)

	var total float64
	err := db.Model(&domain.Payment{}).
		Select("COALESCE(SUM(amount_paid), 0)").
		Where("transaction_status = ? AND payment_date >= ? AND payment_date < ?", "settlement", startDate, endDate).
		Scan(&total).Error
	return total, err
}

func (r *SuperUserRepositoryImpl) GetMonthlyIncomeForYear(db *gorm.DB, year int) ([]struct {
	Month int
	Total float64
}, error) {
	startDate := utils.HijriToGregorian(year, 1, 1)
	endDate := utils.HijriToGregorian(year+1, 1, 1)

	var payments []domain.Payment
	err := db.Model(&domain.Payment{}).
		Where("transaction_status = ? AND payment_date >= ? AND payment_date < ?", "settlement", startDate, endDate).
		Find(&payments).Error
	if err != nil {
		return nil, err
	}

	monthlyMap := make(map[int]float64)
	for _, p := range payments {
		if p.PaymentDate == nil {
			continue
		}
		hDate := utils.GregorianToHijri(*p.PaymentDate)
		monthlyMap[hDate.Month] += p.AmountPaid
	}

	var results []struct {
		Month int
		Total float64
	}
	for m := 1; m <= 12; m++ {
		if val, ok := monthlyMap[m]; ok {
			results = append(results, struct {
				Month int
				Total float64
			}{Month: m, Total: val})
		}
	}
	return results, nil
}

func (r *SuperUserRepositoryImpl) GetAvailableYears(db *gorm.DB) ([]int, error) {
	var years []int
	err := db.Model(&domain.Invoice{}).Distinct("hijri_year").Order("hijri_year desc").Pluck("hijri_year", &years).Error
	if err != nil {
		return nil, err
	}
	if len(years) == 0 {
		// Fallback defaults if database has no invoices
		hijriNow := utils.GetCurrentHijriDate()
		for y := hijriNow.Year; y >= 1443; y-- {
			years = append(years, y)
		}
	}
	return years, nil
}

func (r *SuperUserRepositoryImpl) FindAllStudents(db *gorm.DB, search string) ([]domain.Student, error) {
	query := db.Preload("Guardians").Preload("Guardians.User").Preload("Status", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).Preload("Addresses").Preload("BillingCategories")

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR student_number LIKE ?", s, s, s)
	}

	var students []domain.Student
	err := query.Order("first_name asc").Find(&students).Error
	return students, err
}

func (r *SuperUserRepositoryImpl) FindAllStudentsPaginated(db *gorm.DB, search string, page int, limit int) ([]domain.Student, int64, error) {
	query := db.Model(&domain.Student{})

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR student_number LIKE ?", s, s, s)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var students []domain.Student
	err := query.Preload("Guardians").Preload("Guardians.User").Preload("Status", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).Preload("Addresses").Preload("BillingCategories").
		Order("first_name asc").
		Limit(limit).
		Offset(offset).
		Find(&students).Error
	return students, total, err
}

func (r *SuperUserRepositoryImpl) FindStudentByID(db *gorm.DB, id string) (domain.Student, error) {
	var student domain.Student
	err := db.Preload("Guardians").Preload("Guardians.User").Preload("Status", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).Preload("Addresses").Preload("BillingCategories").
		First(&student, "id = ?", id).Error
	return student, err
}

func (r *SuperUserRepositoryImpl) CreateStudent(db *gorm.DB, student *domain.Student) error {
	return db.Create(student).Error
}

func (r *SuperUserRepositoryImpl) SaveStudent(db *gorm.DB, student *domain.Student) error {
	return db.Save(student).Error
}

func (r *SuperUserRepositoryImpl) DeleteStudent(db *gorm.DB, id string) error {
	return db.Delete(&domain.Student{}, "id = ?", id).Error
}

func (r *SuperUserRepositoryImpl) FindRoleByName(db *gorm.DB, name string) (domain.Role, error) {
	var role domain.Role
	err := db.Where("name = ?", name).First(&role).Error
	return role, err
}

func (r *SuperUserRepositoryImpl) CreateUser(db *gorm.DB, user *domain.User) error {
	return db.Create(user).Error
}

func (r *SuperUserRepositoryImpl) CreateGuardian(db *gorm.DB, guardian *domain.Guardian) error {
	return db.Create(guardian).Error
}

func (r *SuperUserRepositoryImpl) FindGuardianByStudentID(db *gorm.DB, studentID string) (domain.Guardian, error) {
	var guardian domain.Guardian
	err := db.Preload("User").Where("student_id = ?", studentID).First(&guardian).Error
	return guardian, err
}

func (r *SuperUserRepositoryImpl) SaveGuardian(db *gorm.DB, guardian *domain.Guardian) error {
	return db.Save(guardian).Error
}

func (r *SuperUserRepositoryImpl) FindGuardianByID(db *gorm.DB, id string) (domain.Guardian, error) {
	var guardian domain.Guardian
	err := db.Preload("User").Preload("Student").Preload("Student.Addresses").
		Where("id = ?", id).First(&guardian).Error
	return guardian, err
}

func (r *SuperUserRepositoryImpl) SaveUser(db *gorm.DB, user *domain.User) error {
	return db.Save(user).Error
}

func (r *SuperUserRepositoryImpl) FindAllStatusTypes(db *gorm.DB) ([]domain.StudentStatusType, error) {
	var statuses []domain.StudentStatusType
	err := db.Find(&statuses).Error
	return statuses, err
}

func (r *SuperUserRepositoryImpl) FindStatusTypeByID(db *gorm.DB, id string) (domain.StudentStatusType, error) {
	var status domain.StudentStatusType
	err := db.First(&status, "id = ?", id).Error
	return status, err
}

func (r *SuperUserRepositoryImpl) CreateStatusType(db *gorm.DB, statusType *domain.StudentStatusType) error {
	return db.Create(statusType).Error
}

func (r *SuperUserRepositoryImpl) SaveStatusType(db *gorm.DB, statusType *domain.StudentStatusType) error {
	return db.Save(statusType).Error
}

func (r *SuperUserRepositoryImpl) DeleteStatusType(db *gorm.DB, id string) error {
	return db.Delete(&domain.StudentStatusType{}, "id = ?", id).Error
}

func (r *SuperUserRepositoryImpl) FindAllCategories(db *gorm.DB) ([]domain.Category, error) {
	var categories []domain.Category
	err := db.Find(&categories).Error
	return categories, err
}

func (r *SuperUserRepositoryImpl) FindCategoryByID(db *gorm.DB, id string) (domain.Category, error) {
	var category domain.Category
	err := db.First(&category, "id = ?", id).Error
	return category, err
}

func (r *SuperUserRepositoryImpl) CreateCategory(db *gorm.DB, category *domain.Category) error {
	return db.Create(category).Error
}

func (r *SuperUserRepositoryImpl) SaveCategory(db *gorm.DB, category *domain.Category) error {
	return db.Save(category).Error
}

func (r *SuperUserRepositoryImpl) DeleteCategory(db *gorm.DB, id string) error {
	return db.Delete(&domain.Category{}, "id = ?", id).Error
}

func (r *SuperUserRepositoryImpl) FindStudentsForBilling(db *gorm.DB) ([]domain.Student, error) {
	var students []domain.Student
	err := db.Preload("Status", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).Preload("Guardians").Preload("BillingCategories").
		Joins("JOIN student_status_types ON student_status_types.id = students.status_id").
		Where("students.is_active = ? AND student_status_types.is_active_billing = ?", true, true).
		Find(&students).Error
	return students, err
}

func (r *SuperUserRepositoryImpl) FindActiveCategories(db *gorm.DB) ([]domain.Category, error) {
	var categories []domain.Category
	err := db.Where("is_active = ? AND is_fixed = ?", true, true).Find(&categories).Error
	return categories, err
}

func (r *SuperUserRepositoryImpl) FindCategoryByName(db *gorm.DB, name string) (domain.Category, error) {
	var category domain.Category
	err := db.Where("name = ? AND is_active = ?", name, true).First(&category).Error
	return category, err
}

func (r *SuperUserRepositoryImpl) CountInvoice(db *gorm.DB, studentID string, categoryID string, month int, year int) (int64, error) {
	var count int64
	err := db.Model(&domain.Invoice{}).
		Where("student_id = ? AND category_id = ? AND month = ? AND year = ?", studentID, categoryID, month, year).
		Count(&count).Error
	return count, err
}

func (r *SuperUserRepositoryImpl) CountInvoiceByHijri(db *gorm.DB, studentID string, categoryID string, hijriMonth int, hijriYear int) (int64, error) {
	var count int64
	err := db.Model(&domain.Invoice{}).
		Where("student_id = ? AND category_id = ? AND hijri_month = ? AND hijri_year = ?", studentID, categoryID, hijriMonth, hijriYear).
		Count(&count).Error
	return count, err
}

func (r *SuperUserRepositoryImpl) CreateInvoice(db *gorm.DB, invoice *domain.Invoice) error {
	return db.Create(invoice).Error
}

func (r *SuperUserRepositoryImpl) FindAllInvoices(db *gorm.DB, status string, month int, year int) ([]domain.Invoice, error) {
	query := db.Preload("Student").Preload("Category", func(db *gorm.DB) *gorm.DB { return db.Unscoped() })

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if month > 0 {
		query = query.Where("hijri_month = ?", month)
	}
	if year > 0 {
		query = query.Where("hijri_year = ?", year)
	}

	var invoices []domain.Invoice
	err := query.Order("created_at desc").Find(&invoices).Error
	return invoices, err
}

func (r *SuperUserRepositoryImpl) FindStudentsWithInvoicesPaginated(db *gorm.DB, status string, month int, year int, search string, page int, limit int) ([]domain.Student, int64, error) {
	studentIDsQuery := db.Model(&domain.Invoice{}).
		Select("DISTINCT invoices.student_id").
		Joins("JOIN students ON students.id = invoices.student_id")

	if status != "" {
		studentIDsQuery = studentIDsQuery.Where("invoices.status = ?", status)
	}
	if month > 0 {
		studentIDsQuery = studentIDsQuery.Where("invoices.hijri_month = ?", month)
	}
	if year > 0 {
		studentIDsQuery = studentIDsQuery.Where("invoices.hijri_year = ?", year)
	}
	if search != "" {
		s := "%" + search + "%"
		studentIDsQuery = studentIDsQuery.Where("students.first_name LIKE ? OR students.last_name LIKE ? OR students.student_number LIKE ?", s, s, s)
	}

	var total int64
	countQuery := db.Model(&domain.Invoice{}).
		Joins("JOIN students ON students.id = invoices.student_id")
	if status != "" {
		countQuery = countQuery.Where("invoices.status = ?", status)
	}
	if month > 0 {
		countQuery = countQuery.Where("invoices.hijri_month = ?", month)
	}
	if year > 0 {
		countQuery = countQuery.Where("invoices.hijri_year = ?", year)
	}
	if search != "" {
		s := "%" + search + "%"
		countQuery = countQuery.Where("students.first_name LIKE ? OR students.last_name LIKE ? OR students.student_number LIKE ?", s, s, s)
	}
	if err := countQuery.Distinct("invoices.student_id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var studentIDs []string
	if err := studentIDsQuery.Pluck("invoices.student_id", &studentIDs).Error; err != nil {
		return nil, 0, err
	}

	if len(studentIDs) == 0 {
		return []domain.Student{}, total, nil
	}

	offset := (page - 1) * limit

	studentsQuery := db.Preload("Status", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).Preload("Addresses").Preload("Guardians").Preload("Guardians.User").
		Where("id IN ?", studentIDs).
		Order("first_name asc").
		Limit(limit).
		Offset(offset)

	var students []domain.Student
	if err := studentsQuery.Find(&students).Error; err != nil {
		return nil, 0, err
	}

	for i := range students {
		invoiceQuery := db.Preload("Category", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).Where("student_id = ?", students[i].ID)
		var invoices []domain.Invoice
		if err := invoiceQuery.Order("hijri_year desc, hijri_month desc").Find(&invoices).Error; err == nil {
			students[i].Invoices = invoices
		}
	}

	return students, total, nil
}

func (r *SuperUserRepositoryImpl) GetSemesterPaymentStats(db *gorm.DB, hijriYear int) ([]struct {
	Semester     int
	AcademicYear string
	Total        float64
	InvoiceCount int
	PaidCount    int
	UnpaidCount  int
}, error) {
	type SemesterStat struct {
		Semester     int     `gorm:"column:semester"`
		AcademicYear string  `gorm:"column:academic_year_label"`
		Total        float64 `gorm:"column:total"`
		InvoiceCount int     `gorm:"column:invoice_count"`
		PaidCount    int     `gorm:"column:paid_count"`
		UnpaidCount  int     `gorm:"column:unpaid_count"`
	}

	var stats []SemesterStat
	err := db.Model(&domain.Invoice{}).
		Select(`semester, 
			academic_year_label, 
			COALESCE(SUM(CASE WHEN status = 'paid' THEN amount_due ELSE 0 END), 0) as total,
			COUNT(*) as invoice_count,
			SUM(CASE WHEN status = 'paid' THEN 1 ELSE 0 END) as paid_count,
			SUM(CASE WHEN status != 'paid' THEN 1 ELSE 0 END) as unpaid_count`).
		Where("hijri_year = ? AND semester > 0", hijriYear).
		Group("semester, academic_year_label").
		Order("semester asc").
		Scan(&stats).Error

	var result []struct {
		Semester     int
		AcademicYear string
		Total        float64
		InvoiceCount int
		PaidCount    int
		UnpaidCount  int
	}
	for _, s := range stats {
		result = append(result, struct {
			Semester     int
			AcademicYear string
			Total        float64
			InvoiceCount int
			PaidCount    int
			UnpaidCount  int
		}{
			Semester:     s.Semester,
			AcademicYear: s.AcademicYear,
			Total:        s.Total,
			InvoiceCount: s.InvoiceCount,
			PaidCount:    s.PaidCount,
			UnpaidCount:  s.UnpaidCount,
		})
	}

	return result, err
}

func (r *SuperUserRepositoryImpl) FindAllUsers(db *gorm.DB) ([]domain.User, error) {
	var users []domain.User
	err := db.Preload("Role").Order("created_at desc").Find(&users).Error
	return users, err
}

func (r *SuperUserRepositoryImpl) FindAllRoles(db *gorm.DB) ([]domain.Role, error) {
	var roles []domain.Role
	err := db.Find(&roles).Error
	return roles, err
}

func (r *SuperUserRepositoryImpl) UpdateUserRole(db *gorm.DB, userID string, roleID string) error {
	return db.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("role_id", roleID).Error
}

func (r *SuperUserRepositoryImpl) ToggleUserActive(db *gorm.DB, userID string) (bool, error) {
	var user domain.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return false, err
	}
	newStatus := !user.IsActive
	err := db.Model(&user).Update("is_active", newStatus).Error
	return newStatus, err
}

func (r *SuperUserRepositoryImpl) FindAllUsersPaginated(db *gorm.DB, page int, limit int) ([]domain.User, int64, error) {
	var total int64
	if err := db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var users []domain.User
	err := db.Preload("Role").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	return users, total, err
}

func (r *SuperUserRepositoryImpl) FindAllInvoicesPaginated(db *gorm.DB, status string, month int, year int, page int, limit int) ([]domain.Invoice, int64, error) {
	query := db.Model(&domain.Invoice{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if month > 0 {
		query = query.Where("hijri_month = ?", month)
	}
	if year > 0 {
		query = query.Where("hijri_year = ?", year)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var invoices []domain.Invoice
	err := query.Preload("Student").Preload("Category", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&invoices).Error
	return invoices, total, err
}

func (r *SuperUserRepositoryImpl) FindAllActivityLogsPaginated(db *gorm.DB, page int, limit int) ([]domain.ActivityLog, int64, error) {
	var total int64
	if err := db.Model(&domain.ActivityLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var logs []domain.ActivityLog
	err := db.Preload("User").Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	return logs, total, err
}
