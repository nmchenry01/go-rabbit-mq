package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/nmchenry/go-rabbit-mq/producer/config"
	"github.com/nmchenry/go-rabbit-mq/producer/messageproducer"
	"github.com/nmchenry/go-rabbit-mq/producer/utils"
)

func readData() ([][]byte, error) {
	files, err := ioutil.ReadDir("./data")
	if err != nil {
		return nil, err
	}

	var data [][]byte

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filePath := filepath.Join(fmt.Sprintf("%s/data", path), file.Name())

		log.Println(filePath)

		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		data = append(data, fileData)
	}

	return data, nil
}

func waitForInterval() {
	randomInt := rand.Intn(1000)
	time.Sleep(time.Duration(randomInt) * time.Millisecond)
	log.Printf("Waiting for %d milliseconds\n", randomInt)
}

func selectRandomData(data [][]byte) []byte {
	return data[rand.Intn(len(data))]
}

func sendMessageAndWait(messageProducer messageproducer.MessageProducer, data [][]byte) {
	for {
		waitForInterval()
		randomData := selectRandomData(data)
		err := messageProducer.Send(randomData)
		utils.LogOnError(err, "There was a problem sending a message")
	}
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

	forever := make(chan int)

	go sendMessageAndWait(inboundRabbitMQProducer, data)
	go sendMessageAndWait(outboundRabbitMQProducer, data)

	<-forever
}
