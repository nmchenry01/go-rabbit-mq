package handlers

import (
	"encoding/xml"
	"fmt"
	"log"

	"github.com/nmchenry/go-rabbit-mq/consumer/models"
)

// InboundHandler - A function for handling inbound messages
func InboundHandler(msg []byte) error {
	log.Println("Received a message from inbound")
	v := models.Pacs008Message{}
	err := xml.Unmarshal(msg, &v)
	if err != nil {
		return err
	}

	fmt.Printf("Message type: %#v\n", v.AppHdr.MsgDefIdr)

	return nil
}

// OutboundHandler - A function for handling outbound messages from RabbitMQ
func OutboundHandler(msg []byte) error {
	log.Println("Received a message from outbound")
	v := models.Pacs008Message{}
	err := xml.Unmarshal(msg, &v)
	if err != nil {
		return err
	}

	fmt.Printf("Message type: %#v\n", v.AppHdr.MsgDefIdr)

	return nil
}
