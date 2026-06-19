package repository

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"gorm.io/gorm"
)

// ReportRepositoryImpl — implementasi ReportRepository.
type ReportRepositoryImpl struct{}

// NewReportRepository — konstruktor ReportRepository.
func NewReportRepository() ReportRepository {
	return &ReportRepositoryImpl{}
}

func (r *ReportRepositoryImpl) GetInvoiceReports(db *gorm.DB, filters map[string]interface{}) ([]domain.Invoice, error) {
	query := db.Preload("Student").Preload("Category", func(db *gorm.DB) *gorm.DB { return db.Unscoped() })

	query = query.Where("status = ?", "paid")

	if year, _ := utils.InterfaceToInt(filters["year"]); year > 0 {
		query = query.Where("hijri_year = ?", year)
	}
	if month, _ := utils.InterfaceToInt(filters["month"]); month > 0 {
		query = query.Where("hijri_month = ?", month)
	}
	if categoryID, ok := filters["category_id"].(string); ok && categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var invoices []domain.Invoice
	err := query.Order("updated_at desc").Find(&invoices).Error
	return invoices, err
}

func (r *ReportRepositoryImpl) GetSummary(db *gorm.DB, filters map[string]interface{}) (map[string]interface{}, error) {
	var totalInvoices int64
	var paidCount int64
	var unpaidCount int64
	var totalAmountPaid float64

	year, _ := utils.InterfaceToInt(filters["year"])
	month, _ := utils.InterfaceToInt(filters["month"])

	countQuery := db.Model(&domain.Invoice{})
	if year > 0 {
		countQuery = countQuery.Where("hijri_year = ?", year)
	}
	if month > 0 {
		countQuery = countQuery.Where("hijri_month = ?", month)
	}

	countQuery.Session(&gorm.Session{}).Count(&totalInvoices)
	countQuery.Session(&gorm.Session{}).Where("status = ?", "paid").Count(&paidCount)
	countQuery.Session(&gorm.Session{}).Where("status = ?", "unpaid").Count(&unpaidCount)

	sumQuery := db.Model(&domain.Payment{}).
		Joins("JOIN invoices ON invoices.id = payments.invoice_id").
		Where("payments.transaction_status = ?", "settlement")

	if year > 0 {
		sumQuery = sumQuery.Where("invoices.hijri_year = ?", year)
	}
	if month > 0 {
		sumQuery = sumQuery.Where("invoices.hijri_month = ?", month)
	}

	sumQuery.Select("COALESCE(SUM(payments.amount_paid), 0)").Scan(&totalAmountPaid)

	return map[string]interface{}{
		"total_invoices":    totalInvoices,
		"paid_count":        paidCount,
		"unpaid_count":      unpaidCount,
		"total_amount_paid": totalAmountPaid,
	}, nil
}
