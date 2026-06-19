package repository

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

// SettingRepository — kontrak akses data konfigurasi sistem.
type SettingRepository interface {
	FindAll(db *gorm.DB) ([]domain.Setting, error)
	FindByKey(db *gorm.DB, key string) (domain.Setting, error)
	Save(db *gorm.DB, setting *domain.Setting) error
}

type SettingRepositoryImpl struct{}

func NewSettingRepository() SettingRepository {
	return &SettingRepositoryImpl{}
}

func (r *SettingRepositoryImpl) FindAll(db *gorm.DB) ([]domain.Setting, error) {
	var settings []domain.Setting
	err := db.Find(&settings).Error
	return settings, err
}

func (r *SettingRepositoryImpl) FindByKey(db *gorm.DB, key string) (domain.Setting, error) {
	var setting domain.Setting
	err := db.Where("`key` = ?", key).First(&setting).Error
	return setting, err
}

func (r *SettingRepositoryImpl) Save(db *gorm.DB, setting *domain.Setting) error {
	return db.Save(setting).Error
}
