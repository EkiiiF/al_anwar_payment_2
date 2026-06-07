package repository

import (
	"context"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type GuardianRepositoryImpl struct{}

func NewGuardianRepository() GuardianRepository {
	return &GuardianRepositoryImpl{}
}

func (r *GuardianRepositoryImpl) FindGuardianByUserID(ctx context.Context, db *gorm.DB, userID string) (domain.Guardian, error) {
	var guardian domain.Guardian
	err := db.WithContext(ctx).
		Preload("Student").Preload("Student.Status").
		Where("user_id = ?", userID).
		First(&guardian).Error
	return guardian, err
}

func (r *GuardianRepositoryImpl) FindGuardiansByUserID(ctx context.Context, db *gorm.DB, userID string) ([]domain.Guardian, error) {
	var guardians []domain.Guardian
	err := db.WithContext(ctx).
		Preload("Student").Preload("Student.Status").
		Where("user_id = ?", userID).
		Find(&guardians).Error
	return guardians, err
}

func (r *GuardianRepositoryImpl) FindStudentByID(ctx context.Context, db *gorm.DB, studentID string) (domain.Student, error) {
	var student domain.Student
	err := db.WithContext(ctx).
		Preload("Status").
		Where("id = ?", studentID).
		First(&student).Error
	return student, err
}

func (r *GuardianRepositoryImpl) FindInvoicesByStudentID(ctx context.Context, db *gorm.DB, studentID string) ([]domain.Invoice, error) {
	var invoices []domain.Invoice
	err := db.WithContext(ctx).
		Preload("Category").Preload("Payments").
		Where("student_id = ?", studentID).
		Order("year desc, month desc").
		Find(&invoices).Error
	return invoices, err
}

func (r *GuardianRepositoryImpl) FindPaymentsByStudentID(ctx context.Context, db *gorm.DB, studentID string) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := db.WithContext(ctx).
		Preload("Invoice").
		Joins("JOIN invoices ON invoices.id = payments.invoice_id").
		Where("invoices.student_id = ?", studentID).
		Order("payments.created_at desc").
		Find(&payments).Error
	return payments, err
}

func (r *GuardianRepositoryImpl) FindNotificationsByUserID(ctx context.Context, db *gorm.DB, userID string) ([]domain.Notification, error) {
	var notifications []domain.Notification
	err := db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&notifications).Error
	return notifications, err
}

func (r *GuardianRepositoryImpl) FindUnpaidInvoicesByIDs(ctx context.Context, db *gorm.DB, studentID string, invoiceIDs []string) ([]domain.Invoice, error) {
	var invoices []domain.Invoice
	err := db.WithContext(ctx).
		Where("student_id = ? AND id IN ? AND status = ?", studentID, invoiceIDs, "unpaid").
		Find(&invoices).Error
	return invoices, err
}

func (r *GuardianRepositoryImpl) CreateNotification(ctx context.Context, db *gorm.DB, notification *domain.Notification) error {
	return db.WithContext(ctx).Create(notification).Error
}

func (r *GuardianRepositoryImpl) MarkNotificationAsRead(ctx context.Context, db *gorm.DB, notificationID string, userID string) error {
	return db.WithContext(ctx).Model(&domain.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}