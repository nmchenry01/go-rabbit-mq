package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/nmchenry/go-rabbit-mq/producer/config"
	"github.com/nmchenry/go-rabbit-mq/producer/messageproducer"
	"github.com/nmchenry/go-rabbit-mq/producer/utils"
)

// TODO: Extend to read in all possible message types
func readData() ([]byte, error) {
	data, err := ioutil.ReadFile("./data/pacs008.xml")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func waitForInterval() {
	randomInt := rand.Intn(1000)
	time.Sleep(time.Duration(randomInt) * time.Millisecond)
	log.Printf("Waiting for %d milliseconds\n", randomInt)
}

func main() {
	// Get configuration
	configurations, err := config.Init()
	utils.FailOnError(err, "Failed to initialize app configurations")

	// Read in File(s)
	data, err := readData()
	utils.FailOnError(err, "Failed to read data")

	// Initialize producer(s)
	inboundRabbitMQProducer, err := messageproducer.NewRabbitMQProducer(configurations, "inbound")
	utils.FailOnError(err, "Failed to initialize producer")
	defer inboundRabbitMQProducer.Disconnect()

	outboundRabbitMQProducer, err := messageproducer.NewRabbitMQProducer(configurations, "outbound")
	utils.FailOnError(err, "Failed to initialize producer")
	defer outboundRabbitMQProducer.Disconnect()

	for {
		waitForInterval()

		// Send a message on inbound queue
		inboundRabbitMQProducer.Send(data)
		utils.FailOnError(err, "Failed to publish a message")

		// Send a message on outbound queue
		outboundRabbitMQProducer.Send(data)
		utils.FailOnError(err, "Failed to publish a message")

	}
}
