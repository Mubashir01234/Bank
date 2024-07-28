package utils

// Configuration loading and parsing

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string `mapstructure:"log_level"`
	// Other configuration fields...
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("logger")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("Failed to read configuration: ", err)
	}

	return config, nil
}

type KafkaConsumerConfig struct {
	KafkaBrokers string
	GroupId      string
}

func LoadKafkaConsumerConfig(path string) (*KafkaConsumerConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("kafka")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &KafkaConsumerConfig{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
