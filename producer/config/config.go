package config

import (
	"github.com/spf13/viper"
)

// Configurations - Configurations for application
type Configurations struct {
	InboundProducerConfigurations  ProducerConfigurations
	OutboundProducerConfigurations ProducerConfigurations
}

// ProducerConfigurations - Config struct for message producers
type ProducerConfigurations struct {
	URL          string
	Count        int
	ExchangeName string
}

// Init - Initialize configurations and return struct
func Init() (Configurations, error) {
	viper.SetConfigName("./config")
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
