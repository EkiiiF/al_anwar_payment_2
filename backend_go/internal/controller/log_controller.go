package controller

import "github.com/gofiber/fiber/v3"

// LogController — kontrak handler untuk activity logs.
type LogController interface {
	// GetLogs — ambil seluruh log aktivitas.
	GetLogs(c fiber.Ctx) error
}
