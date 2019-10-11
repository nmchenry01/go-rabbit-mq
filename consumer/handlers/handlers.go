package handlers

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/nmchenry/go-rabbit-mq/consumer/models"
	"github.com/nmchenry/go-rabbit-mq/consumer/utils"
	"github.com/streadway/amqp"
)

func InboundHandler(msg amqp.Delivery) {
	log.Println("Received a message from inbound")
	v := models.Pacs008Message{}
	err := xml.Unmarshal(msg.Body, &v)
	utils.FailOnError(err, "Failed to unmarshal XML")

	fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)
}

func OutboundHandler(msg amqp.Delivery) {
	log.Println("Received a message from outbound")
	v := models.Pacs008Message{}
	err := xml.Unmarshal(msg.Body, &v)
	utils.FailOnError(err, "Failed to unmarshal XML")

	fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)
}
