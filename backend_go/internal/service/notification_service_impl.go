package service

import (
	"context"
	"log"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"gorm.io/gorm"
)

type NotificationServiceImpl struct {
	DB *gorm.DB
}

func NewNotificationService(db *gorm.DB) NotificationService {
	return &NotificationServiceImpl{DB: db}
}

func (s *NotificationServiceImpl) Send(ctx context.Context, userID, title, message, category string) error {
	notification := domain.Notification{
		ID:       utils.GenerateID(),
		UserID:   userID,
		Title:    title,
		Message:  message,
		Type:     &category,
		IsRead:   false,
	}

	if err := s.DB.WithContext(ctx).Create(&notification).Error; err != nil {
		log.Printf("[NotificationService] Gagal simpan notifikasi: %v", err)
		return err
	}

	return nil
}

func (s *NotificationServiceImpl) GetNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	var notifications []domain.Notification
	err := s.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&notifications).Error
	return notifications, err
}
