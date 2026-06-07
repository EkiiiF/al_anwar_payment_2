package controller

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/gofiber/fiber/v3"
)

type LogControllerImpl struct {
	LogService service.LogService
}

func NewLogController(logService service.LogService) LogController {
	return &LogControllerImpl{LogService: logService}
}

func (ctrl *LogControllerImpl) GetLogs(c fiber.Ctx) error {
	res, err := ctrl.LogService.GetLogs(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(response.Success(res))
}
