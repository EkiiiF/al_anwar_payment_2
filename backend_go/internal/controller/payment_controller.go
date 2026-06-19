package controller

import "github.com/gofiber/fiber/v3"

// PaymentController — kontrak handler untuk transaksi pembayaran.
type PaymentController interface {
	// CreateTransaction — buat transaksi Midtrans Snap.
	CreateTransaction(ctx fiber.Ctx) error
	// HandleNotification — proses webhook callback Midtrans.
	HandleNotification(ctx fiber.Ctx) error
	// CheckStatus — cek status pembayaran via Midtrans.
	CheckStatus(ctx fiber.Ctx) error
}
