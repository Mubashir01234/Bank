package api

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/Mubashir01234/bank/bank_api/internal/domain"
)

func NewApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			slog.Error("Failed to handle request", slog.Any("error", err))

			message := err.Error()

			var code int
			var ferr *fiber.Error

			if errors.As(err, &ferr) {
				code = ferr.Code
			} else {
				switch {
				case errors.Is(err, domain.ErrResourceNotFound):
					code = fiber.StatusNotFound
				case errors.Is(err, domain.ErrInGhostMode):
					code = fiber.StatusNotFound
				case errors.Is(err, domain.ErrInvalidParameter):
					code = fiber.StatusBadRequest
				case errors.Is(err, domain.ErrSessionNotPresent):
					code = fiber.StatusForbidden
				case errors.Is(err, domain.ErrSessionNotActive):
					code = fiber.StatusUnauthorized
				case errors.Is(err, domain.ErrInvalidAssetType):
					code = fiber.StatusBadRequest
				case errors.Is(err, domain.ErrResourceGone):
					code = fiber.StatusGone
				default:
					code = fiber.StatusInternalServerError
					// Don't use the error value here since we don't want to leak internal
					// information for these kind of errors.
					message = "internal server error"
				}
			}

			return ctx.Status(code).JSON(fiber.Error{
				Code:    code,
				Message: fmt.Sprintf("Failed to handle request: %v", message),
			})
		},
	})
}

type Dependencies struct {
	BankHandler *BankHandler
}

func SetUpRoutes(app *fiber.App, dependencies *Dependencies) {
	bankHandler := dependencies.BankHandler

	app.Get("/ping", handlePing)

	api := app.Group("/api/v1")

	// bank routes
	bank := api.Group("/bank")
	{
		bank.Post("/statement", bankHandler.ManageBankStatement)
	}
}

func handlePing(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Bank API",
	})
}
