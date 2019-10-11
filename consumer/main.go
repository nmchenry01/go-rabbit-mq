package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"github.com/nmchenry/go-rabbit-mq/consumer/client"
	"github.com/nmchenry/go-rabbit-mq/consumer/models"
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
				log.Println("Inbound channel closed")
				return
			}

			log.Println("Received a message from inbound")
			v := models.Pacs008Message{}
			err := xml.Unmarshal(msg.Body, &v)
			utils.FailOnError(err, "Failed to unmarshal XML")

			fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)

		case msg, ok := <-amqpChannels["outboundClient"]:
			if !ok {
				log.Println("Outbound channel closed")
				return
			}

			log.Println("Received a message from outbound")
			v := models.Pacs008Message{}
			err := xml.Unmarshal(msg.Body, &v)
			utils.FailOnError(err, "Failed to unmarshal XML")

			fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)
		}
	}
}

func setup(clients map[string]*client.MessageClient) {
	amqpChannels, err := initClients(clients)
	if err != nil {
		utils.FailOnError(err, "Failed to setup clients")
	}

	for _, client := range clients {
		defer client.Connection.Close()
	}

	processMessages(amqpChannels)
}

func main() {
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

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
