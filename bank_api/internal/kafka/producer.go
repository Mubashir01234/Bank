package kafka

import (
	"fmt"

	"github.com/Mubashir01234/bank/bank_api/internal/utils"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(config *utils.KafkaConsumerConfig) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBrokers,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka producer failed to create: %v", err)
	}

	return &Producer{
		producer: producer,
	}, nil
}

func (p *Producer) SendMessage(message Topic) error {
	topic := message.Name()
	msgBody, err := marshalAvro(message)
	if err != nil {
		return fmt.Errorf("unable to marshal avro: %w", err)
	}

	ch := make(chan kafka.Event)
	defer close(ch)

	err = p.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: msgBody,
		},
		ch,
	)
	if err != nil {
		return fmt.Errorf("unable to send message: %w", err)
	}

	e := <-ch

	err = e.(*kafka.Message).TopicPartition.Error
	if err != nil {
		return fmt.Errorf("message delivery failed: %v", err)
	}

	return nil
}
