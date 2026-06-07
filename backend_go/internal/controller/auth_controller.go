package controller

import "github.com/gofiber/fiber/v3"

type AuthController interface {
	Login(ctx fiber.Ctx) error
	Logout(ctx fiber.Ctx) error
	GetProfile(ctx fiber.Ctx) error
	UpdateProfile(ctx fiber.Ctx) error
	ChangePassword(ctx fiber.Ctx) error
}
