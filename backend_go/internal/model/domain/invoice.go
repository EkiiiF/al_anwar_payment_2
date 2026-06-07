package domain

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID                string         `json:"id"             gorm:"type:char(36);primaryKey"`
	StudentID         string         `json:"student_id"     gorm:"type:char(36);not null;index"`
	Student           Student        `json:"student"        gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE"`
	CategoryID        string         `json:"category_id"    gorm:"type:char(36);not null;index"`
	Category          Category       `json:"category"       gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:RESTRICT"`
	Payments          []Payment      `json:"payments"       gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`
	InvoiceNumber     string         `json:"invoice_number" gorm:"type:varchar(80);unique;not null;index"`
	Month             int            `json:"month"          gorm:"type:tinyint;not null"`
	Year              int            `json:"year"           gorm:"type:smallint;not null"`
	HijriMonth        int            `json:"hijri_month"         gorm:"type:tinyint;not null;default:0;index"`
	HijriYear         int            `json:"hijri_year"          gorm:"type:smallint;not null;default:0;index"`
	Semester          int            `json:"semester"             gorm:"type:tinyint;not null;default:0;index"`
	AcademicYearLabel string         `json:"academic_year_label"  gorm:"type:varchar(30);default:''"`
	AmountDue         float64        `json:"amount_due"     gorm:"type:decimal(15,2);not null"`
	Status            string         `json:"status"         gorm:"type:enum('unpaid','pending','paid','expired','cancelled');default:'unpaid';index"`
	DueDate           time.Time      `json:"due_date"`
	NotifiedAt        *time.Time     `json:"notified_at"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (i *Invoice) TableName() string {
	return "invoices"
}
