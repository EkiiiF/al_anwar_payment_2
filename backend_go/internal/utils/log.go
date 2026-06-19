// Package utils — helper untuk log dan encoding payload.
package utils

import (
	"encoding/json"
)

// EncodePayload mengkonversi payload ke JSON string.
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

// StrPtr mengembalikan pointer dari string, nil jika kosong.
func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
