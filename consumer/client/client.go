package client

import (
	"errors"

	"github.com/streadway/amqp"
)

type MessageClient struct {
	ExchangeName string
	QueueName    string
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	Queue        amqp.Queue
}

func connect(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func createChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, err
}

func initTopography(ch *amqp.Channel, exchangeName string, queueName string) (amqp.Queue, error) {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments)
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	// Bind queue to the exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	return q, nil
}

// Setup - Initializes a connection/channel/queue for the client
func (client *MessageClient) Setup(url string) error {
	conn, err := connect(url)
	if err != nil {
		return err
	}

	ch, err := createChannel(conn)
	if err != nil {
		return err
	}

	q, err := initTopography(ch, client.ExchangeName, client.QueueName)
	if err != nil {
		return err
	}

	client.Connection = conn
	client.Channel = ch
	client.Queue = q

	return nil
}

// Consume - Creates a go channel from which to read messages
func (client *MessageClient) Consume() (<-chan amqp.Delivery, error) {
	if client.Connection == nil || client.Channel == nil {
		return nil, errors.New("AMQP connection/client not intialized, please run client setup")
	}

	// Get message channel from channel
	deliveries, err := client.Channel.Consume(
		client.Queue.Name, // queue
		"",                // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}
