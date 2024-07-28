package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Mubashir01234/bank/bank_service/internal/errors"
	"github.com/linkedin/goavro/v2"
)

type EventHandler interface {
	KafkaTopicsToSubscribe() []Topic
	KafkaEventHandler(ctx context.Context, message Topic) *errors.Error
}

type Topic interface {
	Name() string
	Schema() *goavro.Codec
}

func unmarshalAvro(b []byte, target Topic) error {
	schema := target.Schema()

	native, _, err := schema.NativeFromBinary(b)
	if err != nil {
		return fmt.Errorf("native from binary: %w", err)
	}

	textual, err := schema.TextualFromNative(nil, native)
	if err != nil {
		return fmt.Errorf("textual from native: %w", err)
	}

	return json.Unmarshal(textual, target)
}
