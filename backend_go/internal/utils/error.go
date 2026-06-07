package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

func HandleError(err error, context string) bool {
	if err != nil {
		log.Printf("[ERROR] %s: %v", context, err)
		return true
	}
	return false
}

func FormatValidationErrors(err error) string {
	if err == nil {
		return ""
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string

		for _, e := range ve {
			msg := fmt.Sprintf("Field '%s' validation failed: %s", e.Field(), e.Tag())
			errorMessages = append(errorMessages, msg)
		}

		return strings.Join(errorMessages, "; ")
	}

	return err.Error()
}
