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

type clients struct {
	inboundClient  *client.MessageClient
	outboundClient *client.MessageClient
}

func initClients(inboundClient *client.MessageClient, outboundClient *client.MessageClient) (<-chan amqp.Delivery, <-chan amqp.Delivery, error) {
	// Run setup procedure for inboundClient
	err := inboundClient.Setup(url)
	if err != nil {
		return nil, nil, err
	}

	// Run setup procedure for outboundClient
	err = outboundClient.Setup(url)
	if err != nil {
		return nil, nil, err
	}

	// Get message feed for inboundClient
	inbound, err := inboundClient.Consume()
	if err != nil {
		return nil, nil, err
	}

	// Get message feed for outboundClient
	outbound, err := outboundClient.Consume()
	if err != nil {
		return nil, nil, err
	}

	return inbound, outbound, nil
}

func processMessages(client *client.MessageClient, inbound <-chan amqp.Delivery, outbound <-chan amqp.Delivery) {
	for {
		select {
		case msg, ok := <-inbound:
			if !ok {
				log.Println("Channel closed")
				return
			}

			log.Println("Received a message from inbound")
			v := models.Pacs008Message{}
			err := xml.Unmarshal(msg.Body, &v)
			utils.FailOnError(err, "Failed to unmarshal XML")

			fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)

		case msg, ok := <-outbound:
			if !ok {
				log.Println("Channel closed")
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

func setup(clients clients) {
	inbound, outbound, err := initClients(clients.inboundClient, clients.outboundClient)
	if err != nil {
		utils.FailOnError(err, "Failed to setup clients")
	}
	defer clients.inboundClient.Connection.Close()
	defer clients.outboundClient.Connection.Close()

	processMessages(clients.inboundClient, inbound, outbound)
}

func main() {
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	inboundClient := &client.MessageClient{ExchangeName: "xml", QueueName: "xmlQueue"}
	outboundClient := &client.MessageClient{ExchangeName: "otherXml", QueueName: "otherXmlQueue"}

	clients := clients{inboundClient, outboundClient}

	for {
		setup(clients)
	}
}
