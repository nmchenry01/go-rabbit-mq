package main

import (
	"log"

	"github.com/nmchenry/go-rabbit-mq/consumer/config"
	"github.com/nmchenry/go-rabbit-mq/consumer/handlers"
	"github.com/nmchenry/go-rabbit-mq/consumer/messageclient"
	"github.com/nmchenry/go-rabbit-mq/consumer/utils"
)

func main() {
	log.Println(" [*] Waiting for messages. To exit press CTRL+C")

	// Get configuration
	configurations, err := config.Init()
	utils.FailOnError(err, "Failed to initialize app configurations")

	// Init clients
	inboundRabbitMQClient, err := messageclient.NewRabbitMQClient(configurations, "inbound")
	utils.FailOnError(err, "There was a problem initializing a client")

	if inboundRabbitMQClient != nil {
		defer inboundRabbitMQClient.Disconnect()
	}

	outboundRabbitMQClient, err := messageclient.NewRabbitMQClient(configurations, "outbound")
	utils.FailOnError(err, "There was a problem initializing a client")

	if outboundRabbitMQClient != nil {
		defer outboundRabbitMQClient.Disconnect()
	}

	// Consume from the channels
	inboundChannel, err := inboundRabbitMQClient.Consume()
	utils.FailOnError(err, "There was a problem creating a message channel")

	outboundChannel, err := outboundRabbitMQClient.Consume()
	utils.FailOnError(err, "There was a problem creating a message channel")

	// Process the messages coming off the channels
	// TODO: Potentially refactor handler logic to leverage goroutines
	for {
		select {
		case msg, ok := <-inboundChannel:
			if !ok {
				log.Printf("Inbound channel closed unexpectedly")
				return
			}

			err := handlers.InboundHandler(msg)
			utils.LogOnError(err, "There was a problem processing a message")
		case msg, ok := <-outboundChannel:
			if !ok {
				log.Println("Outbound channel closed unexpectedly")
				return
			}

			err := handlers.OutboundHandler(msg)
			utils.LogOnError(err, "There was a problem processing a message")
		}
	}
}
