package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Mubashir01234/bank/bank_service/internal/kafka"
	"github.com/Mubashir01234/bank/bank_service/internal/service"
	"github.com/Mubashir01234/bank/bank_service/internal/utils"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	projectRoot := utils.FindProjectRoot("config")
	configPath := filepath.Join(projectRoot, "config", utils.GetEnviroment())

	config, err := utils.LoadConfig(configPath)
	if err != nil {
		utils.LogFatal("Failed to load configuration", err)
	}

	utils.InitializeLogger(config)
	ctx := context.Background()

	kafkaConfig, err := utils.LoadKafkaConsumerConfig(configPath)
	if err != nil {
		utils.LogFatal("Failed to load Kafka configuration", err)
	}

	bankService := service.NewBankService()

	kafkaConsumer, err := kafka.NewConsumer(ctx, kafkaConfig, bankService)
	if err != nil {
		utils.LogFatal("Failed to create Kafka consumer", err)
	}

	kafkaConsumer.Start()
	log.Info("Kafka consumer started")

	// wait for signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	kafkaConsumer.Wait()
}
