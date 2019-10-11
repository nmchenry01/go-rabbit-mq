package handlers

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/nmchenry/go-rabbit-mq/consumer/models"
	"github.com/streadway/amqp"
)

func InboundHandler(msg amqp.Delivery) error {
	log.Println("Received a message from inbound")
	v := models.Pacs008Message{}
	err := xml.Unmarshal(msg.Body, &v)
	if err != nil {
		return err
	}

	fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)

	return nil
}

func OutboundHandler(msg amqp.Delivery) error {
	log.Println("Received a message from outbound")
	v := models.Pacs008Message{}
	err := xml.Unmarshal(msg.Body, &v)
	if err != nil {
		return err
	}

	fmt.Printf("SignatureValue: %#v\n", v.AppHdr.HeadSignature.Signature.SignatureValue)

	return nil
}
