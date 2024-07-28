package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/linkedin/goavro/v2"
)

type Topic interface {
	Name() string
	Schema() *goavro.Codec
}

func marshalAvro(m Topic) ([]byte, error) {
	textual, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the message: %v", err)
	}

	schema := m.Schema()

	native, _, err := schema.NativeFromTextual(textual)
	if err != nil {
		return nil, fmt.Errorf("native from textual: %w", err)
	}

	avroBytes, err := schema.BinaryFromNative(nil, native)
	if err != nil {
		return nil, fmt.Errorf("binary from native: %w", err)
	}

	return avroBytes, nil
}
