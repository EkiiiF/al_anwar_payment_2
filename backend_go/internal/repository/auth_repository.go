package repository

import (
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByUsername(db *gorm.DB, username string) (domain.User, error)
	FindUserByID(db *gorm.DB, id string) (domain.User, error)
	UpdateLastLogin(db *gorm.DB, userID string, t time.Time) error
	IsTokenBlacklisted(db *gorm.DB, token string) bool
	BlacklistToken(db *gorm.DB, token string, expiredAt time.Time) error
}
