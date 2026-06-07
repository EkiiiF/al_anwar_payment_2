package repository

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetInvoiceReports(db *gorm.DB, filters map[string]interface{}) ([]domain.Invoice, error)
	GetSummary(db *gorm.DB, filters map[string]interface{}) (map[string]interface{}, error)
}
