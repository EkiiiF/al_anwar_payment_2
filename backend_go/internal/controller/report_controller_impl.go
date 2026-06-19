package controller

import (
	"fmt"
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/gofiber/fiber/v3"
)

// ReportControllerImpl — implementasi ReportController.
type ReportControllerImpl struct {
	Service service.ReportService
}

// NewReportController — konstruktor ReportController.
func NewReportController(s service.ReportService) ReportController {
	return &ReportControllerImpl{Service: s}
}

// GetReports — ambil rekap laporan keuangan.
func (ctrl *ReportControllerImpl) GetReports(ctx fiber.Ctx) error {
	filters := map[string]interface{}{
		"year":        ctx.Query("year"),
		"month":       ctx.Query("month"),
		"category_id": ctx.Query("category"),
	}

	data, err := ctrl.Service.GetReport(ctx.Context(), filters)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(data))
}

// ExportReports — export laporan keuangan ke CSV.
func (ctrl *ReportControllerImpl) ExportReports(ctx fiber.Ctx) error {
	filters := map[string]interface{}{
		"year":        ctx.Query("year"),
		"month":       ctx.Query("month"),
		"category_id": ctx.Query("category"),
	}

	csvData, err := ctrl.Service.ExportCSV(ctx.Context(), filters)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("laporan_keuangan_%s.csv", time.Now().Format("20060102_150405"))

	ctx.Set("Content-Type", "text/csv")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	return ctx.Send(csvData)
}
