package middleware

import (
	"strings"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"github.com/gofiber/fiber/v3"
)

// RequireRole memastikan user memiliki salah satu role yang diizinkan.
func RequireRole(roles ...string) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		user, ok := ctx.Locals("user").(*domain.User)
		if !ok || user == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"status":  "UNAUTHORIZED",
				"message": "Akses tidak terautentikasi",
			})
		}

		userRole := strings.ToLower(user.Role.Name)
		for _, role := range roles {
			if userRole == strings.ToLower(role) {
				return ctx.Next()
			}
		}

		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"status":  "FORBIDDEN",
			"message": "Anda tidak memiliki akses ke fitur ini",
		})
	}
}
