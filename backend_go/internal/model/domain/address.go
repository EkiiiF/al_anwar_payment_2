package domain

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID          string         `json:"id" gorm:"type:char(36);primaryKey"`
	StudentID   string         `json:"student_id" gorm:"type:char(36);not null;index"`
	Student     Student        `json:"student" gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE"`
	AddressLine string         `json:"address_line" gorm:"type:text;not null;column:address_line"`
	Village     string         `json:"village" gorm:"type:varchar(100);column:village"`
	District    string         `json:"district" gorm:"type:varchar(100);column:district"`
	City        string         `json:"city" gorm:"type:varchar(100);column:city"`
	Province    string         `json:"province" gorm:"type:varchar(100);column:province"`
	Country     string         `json:"country" gorm:"type:varchar(100);default:'Indonesia';column:country"`
	PostalCode  string         `json:"postal_code" gorm:"type:varchar(10);column:postal_code"`
	IsPrimary   bool           `json:"is_primary" gorm:"column:is_primary;default:false;index:idx_student_primary"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (a *Address) TableName() string {
	return "addresses"
}
