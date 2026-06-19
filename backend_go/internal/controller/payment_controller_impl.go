package controller

import (
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/exception"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/gofiber/fiber/v3"
)

// PaymentControllerImpl — implementasi PaymentController.
type PaymentControllerImpl struct {
	Service service.PaymentService
}

// NewPaymentController — konstruktor PaymentController.
func NewPaymentController(s service.PaymentService) PaymentController {
	return &PaymentControllerImpl{Service: s}
}

// CreateTransaction — buat transaksi pembayaran via Midtrans Snap.
func (ctrl *PaymentControllerImpl) CreateTransaction(ctx fiber.Ctx) error {
	var req request.CreateTransactionRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format request tidak valid"}
	}

	if len(req.InvoiceIDs) == 0 {
		return exception.ValidationError{Message: "Pilih minimal satu tagihan"}
	}

	res, err := ctrl.Service.CreateTransaction(ctx.Context(), req.InvoiceIDs)
	if err != nil {
		return err
	}

	return ctx.JSON(response.Success(res))
}

// HandleNotification — proses webhook callback dari Midtrans.
func (ctrl *PaymentControllerImpl) HandleNotification(ctx fiber.Ctx) error {
	var notification map[string]interface{}
	if err := ctx.Bind().JSON(&notification); err != nil {
		return exception.ValidationError{Message: "Format notifikasi tidak valid"}
	}

	if err := ctrl.Service.HandleNotification(ctx.Context(), notification); err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{"status": "ok"})
}

// CheckStatus — cek dan perbarui status pembayaran.
func (ctrl *PaymentControllerImpl) CheckStatus(ctx fiber.Ctx) error {
	orderID := ctx.Params("order_id")
	if orderID == "" {
		return exception.ValidationError{Message: "Order ID diperlukan"}
	}

	if err := ctrl.Service.CheckStatus(ctx.Context(), orderID); err != nil {
		return err
	}

	return ctx.JSON(response.SuccessMessage("Status berhasil diperbarui"))
}
