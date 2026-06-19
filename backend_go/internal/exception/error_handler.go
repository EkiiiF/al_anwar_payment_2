package exception

import (
	"log"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/web/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

// ErrorHandler — handler terpusat untuk seluruh error HTTP.
// Menangani NotFoundError, ValidationError, ForbiddenError, validator errors, dan fallback 500.
func ErrorHandler(ctx fiber.Ctx, err error) error {
	if e, ok := err.(NotFoundError); ok {
		return ctx.Status(fiber.StatusNotFound).JSON(response.NotFound(e.Message))
	}
	if e, ok := err.(*NotFoundError); ok {
		return ctx.Status(fiber.StatusNotFound).JSON(response.NotFound(e.Message))
	}

	if e, ok := err.(*fiber.Error); ok {
		if e.Code == fiber.StatusNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(response.NotFound("Route tidak ditemukan"))
		}
		return ctx.Status(e.Code).JSON(response.WebResponse{
			Code:    e.Code,
			Status:  "ERROR",
			Message: e.Message,
		})
	}

	if e, ok := err.(ValidationError); ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest(e.Message))
	}

	if e, ok := err.(ForbiddenError); ok {
		return ctx.Status(fiber.StatusForbidden).JSON(response.Forbidden(e.Message))
	}

	if reflectErr, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, fieldErr := range reflectErr {
			switch fieldErr.Tag() {
			case "required":
				errors[fieldErr.Field()] = fieldErr.Field() + " tidak boleh kosong"
			case "min":
				errors[fieldErr.Field()] = fieldErr.Field() + " terlalu pendek"
			case "max":
				errors[fieldErr.Field()] = fieldErr.Field() + " terlalu panjang"
			default:
				errors[fieldErr.Field()] = fieldErr.Field() + " tidak valid"
			}
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(response.WebResponse{
			Code:    fiber.StatusBadRequest,
			Status:  "BAD_REQUEST",
			Data:    errors,
			Message: "Validasi gagal",
		})
	}

	log.Printf("[ERROR] Internal server error: %v", err)
	return ctx.Status(fiber.StatusInternalServerError).JSON(
		response.InternalServerError("Terjadi kesalahan internal pada server"),
	)
}
