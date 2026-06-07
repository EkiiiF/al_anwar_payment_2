package repository

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type LogRepositoryImpl struct{}

func NewLogRepository() LogRepository {
	return &LogRepositoryImpl{}
}

func (r *LogRepositoryImpl) Create(db *gorm.DB, log *domain.ActivityLog) error {
	return db.Create(log).Error
}

func (r *LogRepositoryImpl) FindAll(db *gorm.DB) ([]domain.ActivityLog, error) {
	var logs []domain.ActivityLog
	err := db.Preload("User").Preload("User.Role").Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (r *LogRepositoryImpl) DeleteOlderThan(db *gorm.DB, days int) error {
	return db.Exec("DELETE FROM activity_logs WHERE created_at < DATE_SUB(NOW(), INTERVAL ? DAY)", days).Error
}
