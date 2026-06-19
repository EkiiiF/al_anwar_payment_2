package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          string         `json:"id" gorm:"type:char(36);primaryKey"`
	Username    string         `json:"username" gorm:"type:varchar(50);unique;not null;index"`
	Password    string         `json:"-" gorm:"type:varchar(255);not null"`
	IsActive    bool           `json:"is_active" gorm:"column:is_active"`
	RoleID      string         `json:"role_id" gorm:"type:char(36);not null;index"`
	Role        Role           `json:"role" gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:RESTRICT"`
	FirstName   string         `json:"first_name" gorm:"type:varchar(100);not null"`
	MiddleName  string         `json:"middle_name,omitempty" gorm:"type:varchar(100)"`
	LastName    string         `json:"last_name,omitempty" gorm:"type:varchar(100)"`
	Email       string         `json:"email" gorm:"type:varchar(100)"`
	PhoneNumber string         `json:"phone_number" gorm:"type:varchar(20)"`
	Gender      string         `json:"gender" gorm:"type:char(1);default:'L'"`
	BirthDate   *time.Time     `json:"birth_date,omitempty"`
	Address     string         `json:"address,omitempty" gorm:"type:text"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	Information string         `json:"information,omitempty" gorm:"-"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (u *User) TableName() string {
	return "users"
}
