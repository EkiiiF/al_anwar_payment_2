package controller

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/gofiber/fiber/v3"
)

// LogControllerImpl — implementasi LogController.
type LogControllerImpl struct {
	LogService service.LogService
}

// NewLogController — konstruktor LogController.
func NewLogController(logService service.LogService) LogController {
	return &LogControllerImpl{LogService: logService}
}

// GetLogs — ambil seluruh log aktivitas.
func (ctrl *LogControllerImpl) GetLogs(c fiber.Ctx) error {
	res, err := ctrl.LogService.GetLogs(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(response.Success(res))
}
