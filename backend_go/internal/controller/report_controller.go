package controller

import "github.com/gofiber/fiber/v3"

type ReportController interface {
	GetReports(ctx fiber.Ctx) error
	ExportReports(ctx fiber.Ctx) error
}
