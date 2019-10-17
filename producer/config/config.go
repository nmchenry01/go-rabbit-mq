package config

import (
	"github.com/spf13/viper"
)

type Configurations struct {
	InboundProducerConfigurations  ProducerConfigurations
	OutboundProducerConfigurations ProducerConfigurations
}

type ProducerConfigurations struct {
	URL          string
	Count        int
	ExchangeName string
}

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
