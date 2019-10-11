package main

import (
	"log"

	"github.com/streadway/amqp"

	"github.com/nmchenry/go-rabbit-mq/consumer/client"
	"github.com/nmchenry/go-rabbit-mq/consumer/config"
	"github.com/nmchenry/go-rabbit-mq/consumer/handlers"
	"github.com/nmchenry/go-rabbit-mq/consumer/utils"
)

func initClients(clients map[string]*client.MessageClient, configuration config.Configurations) (map[string]<-chan amqp.Delivery, error) {
	amqpChannels := make(map[string]<-chan amqp.Delivery)

	for key, client := range clients {
		err := client.Setup(configuration.Client.URL)
		if err != nil {
			return nil, err
		}

		amqpChannel, err := client.Consume()
		if err != nil {
			return nil, err
		}

		amqpChannels[key] = amqpChannel
	}

	return amqpChannels, nil
}

func processMessages(amqpChannels map[string]<-chan amqp.Delivery) {
	for {
		select {
		case msg, ok := <-amqpChannels["inboundClient"]:
			if !ok {
				log.Printf("Inbound channel closed unexpectedly")
				return
			}

			handlers.InboundHandler(msg)
		case msg, ok := <-amqpChannels["outboundClient"]:
			if !ok {
				log.Println("Outbound channel closed unexpectedly")
				return
			}

			handlers.OutboundHandler(msg)
		}
	}
}

func setup(clients map[string]*client.MessageClient, configuration config.Configurations) {
	amqpChannels, err := initClients(clients, configuration)
	if err != nil {
		utils.FailOnError(err, "Failed to setup clients")
	}

	// Make sure all connections are closed/reset on restart
	for _, client := range clients {
		defer client.Connection.Close()
	}

	log.Println("Clients setup and ready to receive messages")

	processMessages(amqpChannels)
}

func main() {
	log.Println(" [*] Waiting for messages. To exit press CTRL+C")

	configuration, err := config.Init()
	if err != nil {
		utils.FailOnError(err, "Failed to initialize app configurations")
	}

	inboundClient := &client.MessageClient{ExchangeName: "inbound", QueueName: "inboundQueue"}
	outboundClient := &client.MessageClient{ExchangeName: "outbound", QueueName: "outboundQueue"}

	clients := map[string]*client.MessageClient{
		"inboundClient":  inboundClient,
		"outboundClient": outboundClient,
	}

	for {
		setup(clients, configuration)
	}
}
