package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/nmchenry/go-rabbit-mq/producer/utils"
	"github.com/streadway/amqp"
)

var url string = "amqp://guest:guest@localhost:5672/"
var count int = 10

func setup(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func main() {
	start := time.Now()

	// Read in XML file
	xml, err := ioutil.ReadFile("./producer/data/pacs008.xml")
	utils.FailOnError(err, "Failed to read XML")

	// Setup a connection with RabbitMQ
	conn, ch, err := setup(url)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"inbound", // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")

	for i := 0; i < count; i++ {

		// Send a message
		err = ch.Publish(
			"inbound", // exchange
			"",        // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "application/xml",
				Body:        xml,
			})
		utils.FailOnError(err, "Failed to publish a message")

	}

	elapsed := time.Since(start)
	log.Printf("Process took %s", elapsed)
}
