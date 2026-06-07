package controller

import "github.com/gofiber/fiber/v3"

type PaymentController interface {
	CreateTransaction(ctx fiber.Ctx) error
	HandleNotification(ctx fiber.Ctx) error
	CheckStatus(ctx fiber.Ctx) error
}
