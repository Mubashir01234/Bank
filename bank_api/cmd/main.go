package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Mubashir01234/bank/bank_api/internal/api"
	"github.com/Mubashir01234/bank/bank_api/internal/kafka"
	"github.com/Mubashir01234/bank/bank_api/internal/repository/postgres"
	"github.com/Mubashir01234/bank/bank_api/internal/service"
	"github.com/Mubashir01234/bank/bank_api/internal/utils"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	projectRoot := utils.FindProjectRoot("config")
	configPath := filepath.Join(projectRoot, "config", utils.GetEnviroment())

	config, err := utils.LoadConfig(configPath)
	if err != nil {
		utils.LogFatal("Failed to load configuration", err)
	}

	utils.InitializeLogger(config)

	kafkaConfig, err := utils.LoadKafkaConsumerConfig(configPath)
	if err != nil {
		utils.LogFatal("Failed to load Kafka configuration", err)
	}

	kafkaProducer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		utils.LogFatal("Failed to create Kafka producer", err)
	}

	bankRepository := postgres.NewBankRepository()
	bankService := service.NewBankService(bankRepository, kafkaProducer)

	app := api.NewApp()
	app.Use(cors.New())

	api.SetUpRoutes(app, &api.Dependencies{
		BankHandler: api.NewBankHandler(bankService),
	})

	go func() {
		if err := app.Listen(":8080"); err != nil {
			utils.LogFatal("Failed to start HTTP server", err)
		}
	}()

	// wait for signal
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	// shut down fiber
	if err := app.Shutdown(); err != nil {
		utils.LogFatal("Failed to stop HTTP server", err)
	}
}
