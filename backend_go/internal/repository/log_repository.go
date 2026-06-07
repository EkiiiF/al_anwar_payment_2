package repository

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type LogRepository interface {
	Create(db *gorm.DB, log *domain.ActivityLog) error
	FindAll(db *gorm.DB) ([]domain.ActivityLog, error)
	DeleteOlderThan(db *gorm.DB, days int) error
}
