package service

import (
	"context"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
)

// GuardianService — kontrak business logic wali santri.
type GuardianService interface {
	GetDashboardStats(ctx context.Context, userID string) (response.GuardianDashboardStats, error)
	GetInvoices(ctx context.Context, userID string) ([]domain.Invoice, error)
	GetPayments(ctx context.Context, userID string) ([]domain.Payment, error)
	GetNotifications(ctx context.Context, userID string) ([]domain.Notification, error)
	MarkNotificationRead(ctx context.Context, notificationID string, userID string) error
}