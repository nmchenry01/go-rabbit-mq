package main

import (
	"log"

	"github.com/nmchenry/go-rabbit-mq/consumer/config"
	"github.com/nmchenry/go-rabbit-mq/consumer/handlers"
	"github.com/nmchenry/go-rabbit-mq/consumer/messageclient"
	"github.com/nmchenry/go-rabbit-mq/consumer/utils"
)

// func initClients(clients map[string]*client.MessageClient, configuration config.Configurations) (map[string]<-chan amqp.Delivery, error) {
// 	amqpChannels := make(map[string]<-chan amqp.Delivery)

// 	for key, client := range clients {
// 		err := client.Setup(configuration.Client.URL)
// 		if err != nil {
// 			return nil, err
// 		}

// 		amqpChannel, err := client.Consume()
// 		if err != nil {
// 			return nil, err
// 		}

// 		amqpChannels[key] = amqpChannel
// 	}

// 	return amqpChannels, nil
// }

// func processMessages(channels map[string]chan []byte) {
// 	for {
// 		select {
// 		case msg, ok := <-channels["inboundClient"]:
// 			if !ok {
// 				log.Printf("Inbound channel closed unexpectedly")
// 				return
// 			}

// 			err := handlers.InboundHandler(msg)
// 			if err != nil {
// 				log.Printf("There was a problem processing a message: %s", err)
// 			}
// 		case msg, ok := <-channels["outboundClient"]:
// 			if !ok {
// 				log.Println("Outbound channel closed unexpectedly")
// 				return
// 			}

// 			err := handlers.OutboundHandler(msg)
// 			if err != nil {
// 				log.Printf("There was a problem processing a message: %s", err)
// 			}
// 		}
// 	}
// }

// func setup(clients map[string]*client.MessageClient, configuration config.Configurations) {
// 	amqpChannels, err := initClients(clients, configuration)
// 	if err != nil {
// 		utils.FailOnError(err, "Failed to setup clients")
// 	}

// 	// Make sure all connections are closed/reset on restart
// 	for _, client := range clients {
// 		defer client.Connection.Close()
// 	}

// 	log.Println("Clients setup and ready to receive messages")

// 	processMessages(amqpChannels)
// }

func main() {
	log.Println(" [*] Waiting for messages. To exit press CTRL+C")

	// Get configuration
	configurations, err := config.Init()
	if err != nil {
		utils.FailOnError(err, "Failed to initialize app configurations")
	}

	// Init clients
	inboundRabbitMQClient, err := messageclient.NewRabbitMQClient(configurations, "inbound")
	if err != nil {
		log.Fatalf("There was a problem initializing a client: %s\n", err)
	}

	if inboundRabbitMQClient != nil {
		defer inboundRabbitMQClient.Disconnect()
	}

	outboundRabbitMQClient, err := messageclient.NewRabbitMQClient(configurations, "outbound")
	if err != nil {
		log.Fatalf("There was a problem initializing a client: %s\n", err)
	}

	if outboundRabbitMQClient != nil {
		defer outboundRabbitMQClient.Disconnect()
	}

	// Consume from the channels
	inboundChannel, err := inboundRabbitMQClient.Consume()
	if err != nil {
		log.Fatalf("There was a problem creating a message channel: %s\n", err)
	}

	outboundChannel, err := outboundRabbitMQClient.Consume()
	if err != nil {
		log.Fatalf("There was a problem creating a message channel: %s\n", err)
	}

	// Process the messages coming off the channels
	for {
		select {
		case msg, ok := <-inboundChannel:
			if !ok {
				log.Printf("Inbound channel closed unexpectedly")
				return
			}

			err := handlers.InboundHandler(msg)
			if err != nil {
				log.Printf("There was a problem processing a message: %s", err)
			}
		case msg, ok := <-outboundChannel:
			if !ok {
				log.Println("Outbound channel closed unexpectedly")
				return
			}

			err := handlers.OutboundHandler(msg)
			if err != nil {
				log.Printf("There was a problem processing a message: %s", err)
			}
		}
	}
}
