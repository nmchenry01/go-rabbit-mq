package main

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/nmchenry/go-rabbit-mq/consumer/client"
	"github.com/nmchenry/go-rabbit-mq/consumer/models"
	"github.com/nmchenry/go-rabbit-mq/consumer/utils"
)

var url string = "amqp://guest:guest@localhost:5672/"

func main() {
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	for {
		client := &client.MessageClient{}

		// Run setup procedure
		err := client.Setup(url)
		utils.FailOnError(err, "Failed to setup client")
		defer client.Connection.Close()
		defer client.Channel.Close()

		// Get message feed
		deliveries, err := client.Consume()
		utils.FailOnError(err, "Failed to consume from RabbitMQ")

		for msg := range deliveries {
			log.Printf("Received a message")

			v := models.Pacs008Message{}
			err := xml.Unmarshal(msg.Body, &v)
			utils.FailOnError(err, "Failed to unmarshal XML")

			// fmt.Printf("Message Struct: %#v\n", v)
			// fmt.Printf("XMLName: %#v\n", v.XMLName)
			// fmt.Printf("AppHdr: %#v\n", v.AppHdr)
			// fmt.Printf("HeadBizMsgIdr: %#v\n", v.AppHdr.BizMsgIdr)
			fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)
		}
	}
}
