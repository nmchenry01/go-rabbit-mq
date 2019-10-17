package config

import (
	"github.com/spf13/viper"
)

type Configurations struct {
	InboundRabbitMQConfigurations  RabbitMQConfigurations
	OutboundRabbitMQConfigurations RabbitMQConfigurations
}

type RabbitMQConfigurations struct {
	URL          string
	ExchangeName string
	QueueName    string
}

func Init() (Configurations, error) {
	viper.SetConfigName("./consumer/config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	var configuration Configurations

	if err := viper.ReadInConfig(); err != nil {
		return configuration, err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}
