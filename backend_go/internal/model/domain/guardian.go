package domain

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Guardian struct {
	ID        string         `json:"id" gorm:"type:char(36);primaryKey"`
	StudentID string         `json:"student_id" gorm:"type:char(36);not null;index"`
	Student   Student        `json:"student" gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE"`
	UserID    string         `json:"user_id" gorm:"type:char(36);not null;index"`
	User      User           `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Name      GuardianName   `json:"name" gorm:"embedded"`
	Phone     string         `json:"phone" gorm:"column:phone"`
	Email     string         `json:"email" gorm:"column:email"`
	Relation  string         `json:"relation" gorm:"type:varchar(50);column:relation"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (g *Guardian) TableName() string {
	return "guardians"
}

func (g *Guardian) FullName() string {
	parts := []string{}
	if g.Name.FirstName != "" {
		parts = append(parts, g.Name.FirstName)
	}
	if g.Name.MiddleName != "" {
		parts = append(parts, g.Name.MiddleName)
	}
	if g.Name.LastName != "" {
		parts = append(parts, g.Name.LastName)
	}
	return strings.Join(parts, " ")
}

type GuardianName struct {
	FirstName  string `json:"first_name" gorm:"column:first_name"`
	MiddleName string `json:"middle_name" gorm:"column:middle_name"`
	LastName   string `json:"last_name" gorm:"column:last_name"`
}
