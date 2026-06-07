package utils

import (
	"encoding/json"
	"fmt"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/gorm"
)

func LogActivity(
	db *gorm.DB,
	userID, action, entityName, entityID string,
	payload interface{},
	ipAddress, userAgent string,
) {
	payloadJSON := EncodePayload(payload)

	log := domain.ActivityLog{
		ID:         GenerateID(),
		UserID:     userID,
		Action:     action,
		EntityName: StrPtr(entityName),
		EntityID:   StrPtr(entityID),
		Payload:    payloadJSON,
		IPAddress:  StrPtr(ipAddress),
		UserAgent:  StrPtr(userAgent),
	}

	if err := db.Create(&log).Error; err != nil {
		fmt.Printf("[WARN] LogActivity gagal: %v\n", err)
	}
}

func EncodePayload(payload interface{}) string {
	if payload == nil {
		return "{}"
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
