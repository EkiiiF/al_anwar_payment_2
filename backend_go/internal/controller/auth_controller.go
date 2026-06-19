// Package controller — HTTP handler untuk autentikasi.
package controller

import "github.com/gofiber/fiber/v3"

// AuthController — kontrak handler autentikasi user.
type AuthController interface {
	// Login — autentikasi dan buat JWT token.
	Login(ctx fiber.Ctx) error
	// Logout — invalidasi token dan blacklist.
	Logout(ctx fiber.Ctx) error
	// GetProfile — ambil data profil user yang sedang login.
	GetProfile(ctx fiber.Ctx) error
	// UpdateProfile — perbarui data profil user.
	UpdateProfile(ctx fiber.Ctx) error
	// ChangePassword — ganti password user.
	ChangePassword(ctx fiber.Ctx) error
}
