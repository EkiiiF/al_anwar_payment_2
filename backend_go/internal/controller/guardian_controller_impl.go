package controller

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/gofiber/fiber/v3"
)

// GuardianControllerImpl — implementasi GuardianController.
type GuardianControllerImpl struct {
	GuardianService service.GuardianService
}

// NewGuardianController — konstruktor GuardianController.
func NewGuardianController(guardianService service.GuardianService) GuardianController {
	return &GuardianControllerImpl{GuardianService: guardianService}
}

// GetDashboard — statistik dashboard wali santri.
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

// GetInvoices — daftar invoice santri yang diampu.
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

// GetPayments — riwayat pembayaran santri.
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

// GetNotifications — daftar notifikasi user.
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

// MarkNotificationRead — tandai notifikasi sudah dibaca.
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
