package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          string         `json:"id" gorm:"type:char(36);primaryKey"`
	Name        string         `json:"name" gorm:"column:name;unique;not null;index"`
	Description string         `json:"description" gorm:"type:text"`
	Users       []User         `json:"users,omitempty" gorm:"foreignKey:RoleID;constraint:OnDelete:RESTRICT"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (r *Role) TableName() string {
	return "roles"
}
