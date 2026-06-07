package repository

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(db *gorm.DB, payment *domain.Payment) error
	UpdatePayment(db *gorm.DB, payment *domain.Payment) error
	UpdatePaymentStatusByExternalID(db *gorm.DB, externalID string, status string, paymentMethod string, rawResponse string) error
	FindPaymentByExternalID(db *gorm.DB, externalID string) ([]domain.Payment, error)
	UpdateInvoiceStatus(db *gorm.DB, id string, status string) error
	UpdateInvoicesStatusByExternalID(db *gorm.DB, externalID string, status string) error
	FindInvoiceByID(db *gorm.DB, id string) (domain.Invoice, error)
	FindPaymentByInvoiceID(db *gorm.DB, invoiceID string) (domain.Payment, error)
}
