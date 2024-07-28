package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Mubashir01234/bank/bank_api/internal/service"
)

// BankHandler handles bet HTTP requests.
type BankHandler struct {
	bankService service.IBankService
}

// NewBankHandler creates new instance of BankHandler.
func NewBankHandler(bankService service.IBankService) *BankHandler {
	return &BankHandler{bankService: bankService}
}

func (h *BankHandler) ManageBankStatement(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file is required",
		})
	}

	bErr := h.bankService.ManageBankStatement(file)
	if bErr != nil {
		return bErr
	}

	return c.JSON(fiber.Map{
		"message": "Bank statement data was sent to the bank service successfully.",
	})
}
