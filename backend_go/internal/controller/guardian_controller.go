package controller

import "github.com/gofiber/fiber/v3"

type GuardianController interface {
	GetDashboard(c fiber.Ctx) error
	GetInvoices(c fiber.Ctx) error
	GetPayments(c fiber.Ctx) error
	GetNotifications(c fiber.Ctx) error
	MarkNotificationRead(c fiber.Ctx) error
}
