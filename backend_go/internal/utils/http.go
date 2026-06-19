package utils

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/exception"
	"github.com/gofiber/fiber/v3"
)

// UserIDFromCtx — ambil user ID dari context JWT middleware.
func UserIDFromCtx(ctx fiber.Ctx) (string, error) {
	raw := ctx.Locals("user_id")
	if raw == nil {
		return "", exception.ValidationError{Message: "Missing authenticated user context"}
	}

	userID, ok := raw.(string)
	if !ok || userID == "" {
		return "", exception.ValidationError{Message: "Invalid authenticated user context"}
	}
	return userID, nil
}
