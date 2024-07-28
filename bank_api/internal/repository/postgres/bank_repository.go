package postgres

import (
	"github.com/Mubashir01234/bank/bank_api/internal/domain"
)

// Assert interface implementation.
var _ domain.BankRepository = (*BankRepository)(nil)

// BankRepository is a struct that will implement the BankRepository interface.
type BankRepository struct {
}

// NewBankRepository creates a new instance of BankRepository.
func NewBankRepository() domain.BankRepository {
	return &BankRepository{}
}
