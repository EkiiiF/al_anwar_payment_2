package domain

import "time"

type ActivityLog struct {
	ID         string    `json:"id"          gorm:"type:char(36);primaryKey"`
	UserID     string    `json:"user_id"     gorm:"type:char(36);index;not null"`
	User       User      `json:"user"        gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Action     string    `json:"action"      gorm:"type:varchar(255);index;not null"`
	EntityName *string   `json:"entity_name" gorm:"type:varchar(50);index"`
	EntityID   *string   `json:"entity_id"   gorm:"type:char(36);index"`
	Payload    string    `json:"payload"     gorm:"type:json"`
	IPAddress  *string   `json:"ip_address"  gorm:"type:varchar(45)"`
	UserAgent  *string   `json:"user_agent"  gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"  gorm:"index;autoCreateTime;<-:create"`
}

func (a *ActivityLog) TableName() string {
	return "activity_logs"
}
