package service

import (
	"context"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
)

// LogService — kontrak business logic activity logging.
type LogService interface {
	Log(ctx context.Context, userID, action, entityName, entityID string, payload interface{}, ip, ua string)
	GetLogs(ctx context.Context) ([]domain.ActivityLog, error)
	CleanupLogs() error
}
