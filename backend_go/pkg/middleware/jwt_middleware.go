// Package middleware — JWT auth dan RBAC middleware.
package middleware

import (
	"strings"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/repository"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// JWTMiddleware — validasi JWT token dari Authorization header.
func JWTMiddleware(db *gorm.DB, authRepo repository.AuthRepository) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"status":  "UNAUTHORIZED",
				"message": "Missing authorization header",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"status":  "UNAUTHORIZED",
				"message": "Format authorization tidak valid",
			})
		}

		tokenString := tokenParts[1]

		if authRepo.IsTokenBlacklisted(db, tokenString) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"status":  "UNAUTHORIZED",
				"message": "Token telah di-invalidasi. Silakan login ulang.",
			})
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"status":  "UNAUTHORIZED",
				"message": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		user, err := authRepo.FindUserByID(db, claims.UserID)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"status":  "UNAUTHORIZED",
				"message": "User tidak ditemukan atau nonaktif",
			})
		}

		ctx.Locals("user", &user)
		ctx.Locals("user_id", user.ID)
		ctx.Locals("user_role", user.Role.Name)
		return ctx.Next()
	}
}
