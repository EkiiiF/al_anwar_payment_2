package domain

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID              string      `json:"id" gorm:"type:char(36);primaryKey"`
	Name            StudentName `json:"name" gorm:"embedded"`
	Nik             string      `json:"nik" gorm:"column:nik;index"`
	BirthDate       string      `json:"birth_date" gorm:"column:birth_date"`
	StudentNumber   string      `json:"student_number" gorm:"column:student_number;unique;not null;index"`
	Gender          string      `json:"gender" gorm:"column:gender"`
	IsActive        bool        `json:"is_active" gorm:"column:is_active;default:true;index"`
	MuhadhorohLevel int         `json:"muhadhoroh_level" gorm:"column:muhadhoroh_level;type:tinyint;default:1;index"`
	CurrentSemester int         `json:"current_semester" gorm:"column:current_semester;type:tinyint;default:1"`
	StatusID string            `json:"status_id" gorm:"type:char(36);not null;index"`
	Status   StudentStatusType `json:"status" gorm:"foreignKey:StatusID;references:ID;constraint:OnDelete:RESTRICT"`
	Guardians []Guardian     `json:"guardians,omitempty" gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE"`
	Addresses []Address      `json:"addresses,omitempty" gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE"`
	Invoices  []Invoice      `json:"invoices,omitempty" gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	FullNameStr string `json:"full_name" gorm:"-"`
}

func (s *Student) AfterFind(tx *gorm.DB) (err error) {
	s.FullNameStr = s.FullName()
	return nil
}

func (s *Student) TableName() string {
	return "students"
}

func (s *Student) FullName() string {
	parts := []string{}
	if s.Name.FirstName != "" {
		parts = append(parts, s.Name.FirstName)
	}
	if s.Name.MiddleName != "" {
		parts = append(parts, s.Name.MiddleName)
	}
	if s.Name.LastName != "" {
		parts = append(parts, s.Name.LastName)
	}
	return strings.Join(parts, " ")
}

type StudentName struct {
	FirstName  string `json:"first_name" gorm:"column:first_name"`
	MiddleName string `json:"middle_name" gorm:"column:middle_name"`
	LastName   string `json:"last_name" gorm:"column:last_name"`
}
