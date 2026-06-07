package domain

import (
	"time"

	"gorm.io/gorm"
)

type Setting struct {
	Key         string         `json:"key" gorm:"primaryKey;type:varchar(100)"`
	Value       string         `json:"value" gorm:"type:text"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (s *Setting) TableName() string {
	return "settings"
}
