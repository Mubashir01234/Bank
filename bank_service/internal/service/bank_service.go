package service

import (
	"context"

	"github.com/Mubashir01234/bank/bank_service/internal/errors"
	"github.com/Mubashir01234/bank/bank_service/internal/kafka"
)

// Assert interface implementation.
var _ IBankService = (*BankService)(nil)

type IBankService interface {
	KafkaTopicsToSubscribe() []kafka.Topic
	KafkaEventHandler(ctx context.Context, msg kafka.Topic) *errors.Error
}

// Service layer for bank
type BankService struct {
}

func NewBankService() IBankService {
	return &BankService{}
}
