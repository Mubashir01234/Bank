package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/Mubashir01234/bank/bank_api/internal/domain"
	"github.com/Mubashir01234/bank/bank_api/internal/errors"
	"github.com/Mubashir01234/bank/bank_api/internal/kafka"
	"github.com/Mubashir01234/bank/bank_api/internal/kafka/schemas"
)

// Assert interface implementation.
var _ IBankService = (*BankService)(nil)

type IBankService interface {
	ManageBankStatement(file *multipart.FileHeader) *errors.Error
}

// Service layer for bank
type BankService struct {
	bankRepository domain.BankRepository
	kafkaProducer  *kafka.Producer
}

func NewBankService(bankRepository domain.BankRepository, kafkaProducer *kafka.Producer) IBankService {
	return &BankService{
		bankRepository: bankRepository,
		kafkaProducer:  kafkaProducer,
	}
}

func (b *BankService) ManageBankStatement(file *multipart.FileHeader) *errors.Error {
	uploadedFile, err := file.Open()
	if err != nil {
		return errors.Wrap(fmt.Errorf("failed to open file: %v", err))
	}

	defer uploadedFile.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, uploadedFile); err != nil {
		return errors.Wrap(fmt.Errorf("failed to read file: %v", err))
	}

	// sending messages to bank service using apache kafka
	if err := b.kafkaProducer.SendMessage(&schemas.BankStatement{
		StatementBuf: buf.Bytes(),
	}); err != nil {
		return errors.Wrap(fmt.Errorf("unable to send message to consumer: %v", err))
	}

	return nil
}
