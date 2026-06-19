package domain

import (
	"time"

	"gorm.io/gorm"
)

type StudentStatusType struct {
	ID                 string         `json:"id"                  gorm:"type:char(36);primaryKey"`
	Name               string         `json:"name"                gorm:"type:varchar(50);unique;not null;index"`
	DiscountPercentage float64        `json:"discount_percentage" gorm:"type:decimal(5,2);default:0.00"`
	IsActiveBilling    bool           `json:"is_active_billing"   gorm:"index"`
	IsDefault          bool           `json:"is_default"          gorm:"default:false"`
	Description        string         `json:"description"         gorm:"type:text"`
	Students           []Student      `json:"students,omitempty"  gorm:"foreignKey:StatusID;constraint:OnDelete:RESTRICT"`
	CreatedAt          time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (s *StudentStatusType) TableName() string {
	return "student_status_types"
}
