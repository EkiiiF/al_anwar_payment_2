package repository

import (
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type AuthRepositoryImpl struct{}

func NewAuthRepository() AuthRepository {
	return &AuthRepositoryImpl{}
}

func (r *AuthRepositoryImpl) FindUserByUsername(db *gorm.DB, username string) (domain.User, error) {
	var user domain.User
	err := db.Preload("Role").
		Where("username = ? AND is_active = ?", username, true).
		First(&user).Error
	return user, err
}

func (r *AuthRepositoryImpl) FindUserByID(db *gorm.DB, id string) (domain.User, error) {
	var user domain.User
	err := db.Preload("Role").
		Where("id = ? AND is_active = ?", id, true).
		First(&user).Error
	return user, err
}

func (r *AuthRepositoryImpl) UpdateLastLogin(db *gorm.DB, userID string, t time.Time) error {
	return db.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("last_login_at", t).Error
}

func (r *AuthRepositoryImpl) IsTokenBlacklisted(db *gorm.DB, token string) bool {
	var count int64
	db.Model(&domain.TokenBlacklist{}).Where("token = ?", token).Count(&count)
	return count > 0
}

func (r *AuthRepositoryImpl) BlacklistToken(db *gorm.DB, token string, expiredAt time.Time) error {
	bl := domain.TokenBlacklist{Token: token, ExpiredAt: expiredAt}
	return db.Create(&bl).Error
}
