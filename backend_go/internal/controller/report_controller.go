package controller

import "github.com/gofiber/fiber/v3"

// ReportController — kontrak handler untuk laporan keuangan.
type ReportController interface {
	// GetReports — ambil rekap laporan keuangan.
	GetReports(ctx fiber.Ctx) error
	// ExportReports — export laporan ke CSV.
	ExportReports(ctx fiber.Ctx) error
}
