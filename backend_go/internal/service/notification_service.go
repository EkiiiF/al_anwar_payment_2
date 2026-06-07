package service

import (
	"context"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
)

type NotificationService interface {
	Send(ctx context.Context, userID, title, message, category string) error
	GetNotifications(ctx context.Context, userID string) ([]domain.Notification, error)
}
