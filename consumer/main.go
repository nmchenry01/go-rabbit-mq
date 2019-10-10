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
	xmlClient      *client.MessageClient
	otherXMLClient *client.MessageClient
}

func setupAllClients(xmlClient *client.MessageClient, otherXMLClient *client.MessageClient) (<-chan amqp.Delivery, <-chan amqp.Delivery, error) {
	// Run setup procedure for xmlClient
	err := xmlClient.Setup(url)
	if err != nil {
		return nil, nil, err
	}

	// Run setup procedure for otherXMLClient
	err = otherXMLClient.Setup(url)
	if err != nil {
		return nil, nil, err
	}

	// Get message feed for xmlClient
	xmlMsgs, err := xmlClient.Consume()
	if err != nil {
		return nil, nil, err
	}

	// Get message feed for otherXMLClient
	otherXMLMsgs, err := otherXMLClient.Consume()
	if err != nil {
		return nil, nil, err
	}

	return xmlMsgs, otherXMLMsgs, nil
}

func handleMessages(client *client.MessageClient, xmlMsgs <-chan amqp.Delivery, otherXMLMsgs <-chan amqp.Delivery) {
	for {
		select {
		case msg, ok := <-xmlMsgs:
			if !ok {
				log.Println("Channel closed")
				return
			}

			log.Println("Received a message from xmlMsgs")
			v := models.Pacs008Message{}
			err := xml.Unmarshal(msg.Body, &v)
			utils.FailOnError(err, "Failed to unmarshal XML")

			fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)

			client.Channel.Close()

		case msg, ok := <-otherXMLMsgs:
			if !ok {
				log.Println("Channel closed")
				return
			}

			log.Println("Received a message from otherXMLMsgs")
			v := models.Pacs008Message{}
			err := xml.Unmarshal(msg.Body, &v)
			utils.FailOnError(err, "Failed to unmarshal XML")

			fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)
		}
	}
}

func processMessages(clients clients) {
	xmlMsgs, otherXMLMsgs, err := setupAllClients(clients.xmlClient, clients.otherXMLClient)
	if err != nil {
		utils.FailOnError(err, "Failed to setup clients")
	}
	defer clients.xmlClient.Connection.Close()
	defer clients.otherXMLClient.Connection.Close()

	handleMessages(clients.xmlClient, xmlMsgs, otherXMLMsgs)
}

func main() {
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	xmlClient := &client.MessageClient{ExchangeName: "xml", QueueName: "xmlQueue"}
	otherXMLClient := &client.MessageClient{ExchangeName: "otherXml", QueueName: "otherXmlQueue"}

	clients := clients{xmlClient, otherXMLClient}

	for {
		processMessages(clients)
	}
}
