package service

import (
	"context"
	"fmt"

	"github.com/Mubashir01234/bank/bank_service/internal/errors"
	"github.com/Mubashir01234/bank/bank_service/internal/kafka"
	"github.com/Mubashir01234/bank/bank_service/internal/kafka/schemas"
	"github.com/Mubashir01234/bank/bank_service/internal/utils"
)

func (s *BankService) KafkaTopicsToSubscribe() []kafka.Topic {
	return []kafka.Topic{
		&schemas.BankStatement{},
	}
}

func (s *BankService) KafkaEventHandler(ctx context.Context, msg kafka.Topic) *errors.Error {
	switch v := msg.(type) {
	case nil:
		return errors.New("the message is nil")

	case *schemas.BankStatement:
		output, err := parseCSVFromBytes(v.StatementBuf)
		if err != nil {
			return errors.New(err.Error())
		}

		resp, err := utils.MapToJSON(output)
		if err != nil {
			return errors.New(err.Error())
		}

		fmt.Println("RESPONSE: ", resp)

	default:
		return errors.New(fmt.Sprintf("unknown message type: %T", v))
	}

	return nil
}
