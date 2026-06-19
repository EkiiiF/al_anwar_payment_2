package service

import (
	"context"
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"gorm.io/gorm"
)

type LogServiceImpl struct {
	DB      *gorm.DB
	LogRepo repository.LogRepository
}

func NewLogService(db *gorm.DB, logRepo repository.LogRepository) LogService {
	return &LogServiceImpl{DB: db, LogRepo: logRepo}
}

func (s *LogServiceImpl) Log(ctx context.Context, userID, action, entityName, entityID string, payload interface{}, ip, ua string) {
	payloadStr := utils.EncodePayload(payload)

	log := domain.ActivityLog{
		ID:         utils.GenerateID(),
		UserID:     userID,
		Action:     action,
		EntityName: utils.StrPtr(entityName),
		EntityID:   utils.StrPtr(entityID),
		Payload:    payloadStr,
		IPAddress:  utils.StrPtr(ip),
		UserAgent:  utils.StrPtr(ua),
		CreatedAt:  time.Now(),
	}

	if err := s.LogRepo.Create(s.DB, &log); err != nil {
		_ = err
	}
}

func (s *LogServiceImpl) GetLogs(ctx context.Context) ([]domain.ActivityLog, error) {
	return s.LogRepo.FindAll(s.DB)
}

func (s *LogServiceImpl) CleanupLogs() error {
	return s.LogRepo.DeleteOlderThan(s.DB, 7)
}
