package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          string         `json:"id"          gorm:"type:char(36);primaryKey"`
	Name        string         `json:"name"        gorm:"type:varchar(50);unique;not null;index"`
	BaseAmount  float64        `json:"base_amount" gorm:"type:decimal(15,2);not null;default:0"`
	IsFixed     bool           `json:"is_fixed"    gorm:"default:true;index"`
	IsActive    bool           `json:"is_active"   gorm:"default:true;index"`
	Description string         `json:"description" gorm:"type:text"`
	Invoices    []Invoice      `json:"invoices,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (c *Category) TableName() string {
	return "categories"
}
