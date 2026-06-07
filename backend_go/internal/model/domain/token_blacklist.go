package domain

import "time"

type TokenBlacklist struct {
	Token     string    `json:"token"      gorm:"primaryKey;type:varchar(512);index"`
	ExpiredAt time.Time `json:"expired_at" gorm:"not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;<-:create"`
}

func (t *TokenBlacklist) TableName() string {
	return "token_blacklist"
}
