package repository

import (
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type PaymentRepositoryImpl struct{}

func NewPaymentRepository() PaymentRepository {
	return &PaymentRepositoryImpl{}
}

func (r *PaymentRepositoryImpl) CreatePayment(db *gorm.DB, payment *domain.Payment) error {
	return db.Create(payment).Error
}

func (r *PaymentRepositoryImpl) UpdatePayment(db *gorm.DB, payment *domain.Payment) error {
	return db.Save(payment).Error
}

func (r *PaymentRepositoryImpl) UpdatePaymentStatusByExternalID(db *gorm.DB, externalID string, status string, paymentMethod string, rawResponse string) error {
	return db.Model(&domain.Payment{}).Where("external_id = ?", externalID).Updates(map[string]interface{}{
		"transaction_status": status,
		"payment_method":     paymentMethod,
		"raw_response":       rawResponse,
		"updated_at":         time.Now(),
	}).Error
}

func (r *PaymentRepositoryImpl) FindPaymentByExternalID(db *gorm.DB, externalID string) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := db.Preload("Invoice").Preload("Invoice.Student").Preload("Invoice.Student.Guardians").
		Where("external_id = ?", externalID).Find(&payments).Error
	return payments, err
}

func (r *PaymentRepositoryImpl) FindInvoiceByID(db *gorm.DB, id string) (domain.Invoice, error) {
	var invoice domain.Invoice
	err := db.Preload("Student").Preload("Student.Guardians").Preload("Category").
		Where("id = ?", id).First(&invoice).Error
	return invoice, err
}

func (r *PaymentRepositoryImpl) UpdateInvoiceStatus(db *gorm.DB, id string, status string) error {
	return db.Model(&domain.Invoice{}).Where("id = ?", id).Update("status", status).Error
}

func (r *PaymentRepositoryImpl) UpdateInvoicesStatusByExternalID(db *gorm.DB, externalID string, status string) error {
	var invoiceIDs []string
	db.Model(&domain.Payment{}).Where("external_id = ?", externalID).Pluck("invoice_id", &invoiceIDs)
	if len(invoiceIDs) == 0 {
		return nil
	}
	return db.Model(&domain.Invoice{}).Where("id IN ?", invoiceIDs).Update("status", status).Error
}

func (r *PaymentRepositoryImpl) FindPaymentByInvoiceID(db *gorm.DB, invoiceID string) (domain.Payment, error) {
	var payment domain.Payment
	err := db.Where("invoice_id = ?", invoiceID).Order("created_at desc").First(&payment).Error
	return payment, err
}
