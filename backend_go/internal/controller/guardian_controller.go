package controller

import "github.com/gofiber/fiber/v3"

// GuardianController — kontrak handler untuk wali santri.
type GuardianController interface {
	// GetDashboard — statistik dashboard wali santri.
	GetDashboard(c fiber.Ctx) error
	// GetInvoices — daftar invoice santri yang diampu.
	GetInvoices(c fiber.Ctx) error
	// GetPayments — riwayat pembayaran santri.
	GetPayments(c fiber.Ctx) error
	// GetNotifications — daftar notifikasi user.
	GetNotifications(c fiber.Ctx) error
	// MarkNotificationRead — tandai notifikasi sudah dibaca.
	MarkNotificationRead(c fiber.Ctx) error
}
