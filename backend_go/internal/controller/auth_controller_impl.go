package controller

import (
	"strings"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/exception"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/request"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/service"
	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
	"github.com/gofiber/fiber/v3"
)

// AuthControllerImpl — implementasi AuthController.
type AuthControllerImpl struct {
	AuthService service.AuthService
}

// NewAuthController — konstruktor AuthController.
func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerImpl{AuthService: authService}
}

// Login — proses autentikasi dan menghasilkan JWT token.
func (c *AuthControllerImpl) Login(ctx fiber.Ctx) error {
	var req request.LoginRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format request tidak valid"}
	}

	res, err := c.AuthService.Login(req, ctx.IP(), ctx.Get("User-Agent"))
	if err != nil {
		return err
	}

	return ctx.JSON(response.Success(res))
}

// Logout — blacklist token dan logout user.
func (c *AuthControllerImpl) Logout(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return exception.ValidationError{Message: "Missing authorization token"}
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	userID, ip, ua := utils.GetOperatorInfo(ctx)

	if err := c.AuthService.Logout(token, userID, ip, ua); err != nil {
		return err
	}

	return ctx.JSON(response.SuccessMessage("Berhasil logout"))
}

// GetProfile — ambil profil user yang sedang login.
func (c *AuthControllerImpl) GetProfile(ctx fiber.Ctx) error {
	userID, err := utils.UserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	res, err := c.AuthService.GetProfile(userID)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(res))
}

// UpdateProfile — perbarui profil user yang sedang login.
func (c *AuthControllerImpl) UpdateProfile(ctx fiber.Ctx) error {
	userID, err := utils.UserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	var req request.UpdateProfileRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format request tidak valid"}
	}

	res, err := c.AuthService.UpdateProfile(userID, req)
	if err != nil {
		return err
	}
	return ctx.JSON(response.Success(res))
}

// ChangePassword — ganti password user yang sedang login.
func (c *AuthControllerImpl) ChangePassword(ctx fiber.Ctx) error {
	userID, err := utils.UserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	var req request.ChangePasswordRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return exception.ValidationError{Message: "Format request tidak valid"}
	}

	if err := c.AuthService.ChangePassword(userID, req); err != nil {
		return err
	}
	return ctx.JSON(response.SuccessMessage("Password berhasil diperbarui"))
}
