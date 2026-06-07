package repository

import (
	"context"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type GuardianRepository interface {
	FindGuardianByUserID(ctx context.Context, db *gorm.DB, userID string) (domain.Guardian, error)
	FindGuardiansByUserID(ctx context.Context, db *gorm.DB, userID string) ([]domain.Guardian, error)
	FindStudentByID(ctx context.Context, db *gorm.DB, studentID string) (domain.Student, error)
	FindInvoicesByStudentID(ctx context.Context, db *gorm.DB, studentID string) ([]domain.Invoice, error)
	FindPaymentsByStudentID(ctx context.Context, db *gorm.DB, studentID string) ([]domain.Payment, error)
	FindNotificationsByUserID(ctx context.Context, db *gorm.DB, userID string) ([]domain.Notification, error)
	FindUnpaidInvoicesByIDs(ctx context.Context, db *gorm.DB, studentID string, invoiceIDs []string) ([]domain.Invoice, error)
	CreateNotification(ctx context.Context, db *gorm.DB, notification *domain.Notification) error
	MarkNotificationAsRead(ctx context.Context, db *gorm.DB, notificationID string, userID string) error
}
