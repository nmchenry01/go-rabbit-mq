package main

import (
	"log"

	"github.com/streadway/amqp"

	"github.com/nmchenry/go-rabbit-mq/consumer/client"
	"github.com/nmchenry/go-rabbit-mq/consumer/handlers"
	"github.com/nmchenry/go-rabbit-mq/consumer/utils"
)

var url string = "amqp://guest:guest@localhost:5672/"

func initClients(clients map[string]*client.MessageClient) (map[string]<-chan amqp.Delivery, error) {
	amqpChannels := make(map[string]<-chan amqp.Delivery)

	for key, client := range clients {
		err := client.Setup(url)
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
				log.Printf("Inbound channel closed")
				return
			}

			handlers.InboundHandler(msg)
		case msg, ok := <-amqpChannels["outboundClient"]:
			if !ok {
				log.Println("Outbound channel closed")
				return
			}

			handlers.OutboundHandler(msg)
		}
	}
}

func setup(clients map[string]*client.MessageClient) {
	amqpChannels, err := initClients(clients)
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

	inboundClient := &client.MessageClient{ExchangeName: "inbound", QueueName: "inboundQueue"}
	outboundClient := &client.MessageClient{ExchangeName: "outbound", QueueName: "outboundQueue"}

	clients := map[string]*client.MessageClient{
		"inboundClient":  inboundClient,
		"outboundClient": outboundClient,
	}

	for {
		setup(clients)
	}
}
