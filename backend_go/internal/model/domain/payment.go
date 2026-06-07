package domain

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID                string         `json:"id"                 gorm:"type:char(36);primaryKey"`
	InvoiceID         string         `json:"invoice_id"         gorm:"type:char(36);not null;index"`
	Invoice           Invoice        `json:"invoice"            gorm:"foreignKey:InvoiceID;references:ID;constraint:OnDelete:CASCADE"`
	AttemptNumber     int            `json:"attempt_number"     gorm:"default:1"`
	ExternalID        *string        `json:"external_id"        gorm:"type:varchar(100);index"`
	SnapToken         *string        `json:"snap_token"         gorm:"type:varchar(255)"`
	AmountPaid        float64        `json:"amount_paid"        gorm:"type:decimal(15,2);not null"`
	PaymentMethod     *string        `json:"payment_method"     gorm:"type:varchar(50)"`
	TransactionStatus *string        `json:"transaction_status" gorm:"type:varchar(50);index"`
	PaymentDate       *time.Time     `json:"payment_date"`
	RawResponse       string         `json:"raw_response"       gorm:"type:json"`
	EvidenceURL       *string        `json:"evidence_url"       gorm:"type:varchar(255)"`
	Notes             string         `json:"notes"              gorm:"type:text"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (p *Payment) TableName() string {
	return "payments"
}
