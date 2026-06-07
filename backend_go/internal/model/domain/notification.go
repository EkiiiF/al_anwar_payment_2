package domain

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        string         `json:"id"         gorm:"type:char(36);primaryKey"`
	UserID    string         `json:"user_id"    gorm:"type:char(36);not null;index"`
	User      User           `json:"user"       gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Title     string         `json:"title"      gorm:"type:varchar(150);not null"`
	Message   string         `json:"message"    gorm:"type:text;not null"`
	Type      *string        `json:"type"       gorm:"type:varchar(50);index"`
	IsRead    bool           `json:"is_read"    gorm:"default:false;index"`
	ReadAt    *time.Time     `json:"read_at"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (n *Notification) TableName() string {
	return "notifications"
}
