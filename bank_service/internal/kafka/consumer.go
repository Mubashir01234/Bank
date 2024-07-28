package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/Mubashir01234/bank/bank_service/internal/errors"
	"github.com/Mubashir01234/bank/bank_service/internal/utils"
	"github.com/gofiber/fiber/v2/log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	ctx      context.Context
	Consumer *kafka.Consumer
	schemas  utils.SafeMap[string, reflect.Type]
	brokers  string
	handler  handlerFunc
	wg       sync.WaitGroup
}

type handlerFunc func(ctx context.Context, message Topic) *errors.Error

func NewConsumer(ctx context.Context, config *utils.KafkaConsumerConfig, eventHandler EventHandler) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBrokers,
		"group.id":          config.GroupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("create Kafka consumer: %w", err)
	} else {
		log.Info("Kafka consumer created")
	}

	schemas := utils.NewSafeMap[string, reflect.Type]()
	for _, topic := range eventHandler.KafkaTopicsToSubscribe() {
		schemas.Store(topic.Name(), reflect.TypeOf(topic).Elem())
	}

	var topicNames []string
	for _, topic := range eventHandler.KafkaTopicsToSubscribe() {
		topicNames = append(topicNames, topic.Name())
	}

	if len(topicNames) > 0 {
		err = createAndSubscribeTopics(ctx, consumer, config, topicNames)
		if err != nil {
			return nil, err
		}
	}

	return &Consumer{
		ctx:      ctx,
		Consumer: consumer,
		brokers:  config.KafkaBrokers,
		schemas:  schemas,
		handler:  eventHandler.KafkaEventHandler,
	}, nil
}

func (c *Consumer) Start() {
	go c.startConsuming()
}

func (c *Consumer) Wait() {
	c.wg.Wait()
}

func (c *Consumer) startConsuming() {
	for {
		select {
		case <-c.ctx.Done():
			log.Info("Received stopped signal. Stopping Kafka consumer...")
			return
		default:
			switch event := c.Consumer.Poll(100).(type) {
			case *kafka.Message:
				if event.TopicPartition.Topic == nil {
					log.Error("Missing topic", "event", event.String())
					break
				}

				slog.Info("Received message event", slog.String("event", event.String()))

				if err := c.handleMessageEvent(event); err != nil {
					log.Error("Failed to handle message event", "error", err.Error())
				}
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover.
				log.Warn("Received error event", "event", event.String())

				if event.Code() == kafka.ErrAllBrokersDown {
					return
				}
			}
		}
	}
}

// handleMessage defines how to process each Kafka message
func (c *Consumer) handleMessageEvent(event *kafka.Message) error {
	c.wg.Add(1)
	defer c.wg.Done()

	topic := *event.TopicPartition.Topic
	value := event.Value

	message := reflect.New(c.schemas.Load(topic)).Interface().(Topic)

	if err := unmarshalAvro(value, message); err != nil {
		return fmt.Errorf("failed to unmarshal Avro: %w", err)
	}

	if hErr := c.handler(c.ctx, message); hErr != nil {
		log.Error(hErr.Error())
		return hErr
	}

	return nil
}

func createAndSubscribeTopics(ctx context.Context, consumer *kafka.Consumer, config *utils.KafkaConsumerConfig, topicNames []string) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": config.KafkaBrokers})
	if err != nil {
		log.Errorf("Failed to create Admin client: %s", err)
		return err
	}
	defer adminClient.Close()

	topicSpecs := make([]kafka.TopicSpecification, len(topicNames))
	for i, topic := range topicNames {
		topicSpecs[i] = kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	results, err := adminClient.CreateTopics(ctx, topicSpecs)
	if err != nil {
		log.Errorf("Failed to create topics: %s", err)
		return err
	}

	for _, result := range results {
		log.Info("Creating topic", "topic", result.Topic, "error", result.Error)
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			return result.Error
		}
	}

	err = consumer.SubscribeTopics(topicNames, nil)
	if err != nil {
		log.Error("Failed to subscribe to Kafka topics", err)
		return err
	}

	return nil
}
