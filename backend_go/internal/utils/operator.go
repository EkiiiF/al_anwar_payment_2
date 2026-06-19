package utils

import "github.com/gofiber/fiber/v3"

// GetOperatorInfo — ambil user ID, IP, dan User-Agent dari context.
func GetOperatorInfo(ctx fiber.Ctx) (userID string, ip string, ua string) {
	idVal := ctx.Locals("user_id")
	if idVal != nil {
		if id, ok := idVal.(string); ok {
			userID = id
		}
	}

	ip = ctx.IP()
	ua = ctx.Get("User-Agent")
	return
}
