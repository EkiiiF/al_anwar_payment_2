package controller

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/gofiber/fiber/v3"
)

type GuardianControllerImpl struct {
	GuardianService service.GuardianService
}

func NewGuardianController(guardianService service.GuardianService) GuardianController {
	return &GuardianControllerImpl{GuardianService: guardianService}
}

func (ctrl *GuardianControllerImpl) GetDashboard(c fiber.Ctx) error {
	userID, err := utils.UserIDFromCtx(c)
	if err != nil {
		return err
	}
	res, err := ctrl.GuardianService.GetDashboardStats(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success(res))
}

func (ctrl *GuardianControllerImpl) GetInvoices(c fiber.Ctx) error {
	userID, err := utils.UserIDFromCtx(c)
	if err != nil {
		return err
	}
	res, err := ctrl.GuardianService.GetInvoices(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success(res))
}

func (ctrl *GuardianControllerImpl) GetPayments(c fiber.Ctx) error {
	userID, err := utils.UserIDFromCtx(c)
	if err != nil {
		return err
	}
	res, err := ctrl.GuardianService.GetPayments(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success(res))
}

func (ctrl *GuardianControllerImpl) GetNotifications(c fiber.Ctx) error {
	userID, err := utils.UserIDFromCtx(c)
	if err != nil {
		return err
	}
	res, err := ctrl.GuardianService.GetNotifications(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(response.Success(res))
}

func (ctrl *GuardianControllerImpl) MarkNotificationRead(c fiber.Ctx) error {
	userID, err := utils.UserIDFromCtx(c)
	if err != nil {
		return err
	}
	id := c.Params("id")
	if err := ctrl.GuardianService.MarkNotificationRead(c.Context(), id, userID); err != nil {
		return err
	}
	return c.JSON(response.SuccessMessage("Notifikasi ditandai dibaca"))
}
