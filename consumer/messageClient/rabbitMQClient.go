package messageclient

import (
	"errors"

	"github.com/nmchenry/go-rabbit-mq/consumer/config"
	"github.com/streadway/amqp"
)

// RabbitMQClient - A struct for interacting with RabbitMQ
type RabbitMQClient struct {
	exchangeName string
	queueName    string
	url          string
	connection   *amqp.Connection
	channel      *amqp.Channel
	queue        amqp.Queue
}

// NewRabbitMQClient - Initializes a new RabbitMQClient (connection, channels, topography, etc)
func NewRabbitMQClient(configurations config.Configurations, messageDirection string) (MessageClient, error) {
	rabbitMQConfigurations, err := selectConfigurations(configurations, messageDirection)
	if err != nil {
		return nil, err
	}

	connection, channel, queue, err := setup(rabbitMQConfigurations)

	newClient := &RabbitMQClient{
		exchangeName: rabbitMQConfigurations.ExchangeName,
		queueName:    rabbitMQConfigurations.QueueName,
		url:          rabbitMQConfigurations.URL,
		connection:   connection,
		channel:      channel,
		queue:        queue,
	}

	return newClient, nil
}

// Consume - Creates a go channel from which to read messages
// TODO: Will have to modify in order to support message acknowledgment
func (rabbitMQClient *RabbitMQClient) Consume() (chan []byte, error) {
	if rabbitMQClient.connection == nil || rabbitMQClient.channel == nil {
		return nil, errors.New("AMQP connection/client not intialized")
	}

	messageChannel := make(chan []byte)

	// Get message channel from channel
	deliveries, err := rabbitMQClient.channel.Consume(
		rabbitMQClient.queueName, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		return nil, err
	}

	// Start up Goroutine and convert messages coming in from RabbitMQ
	go func() {
		for msg := range deliveries {
			messageChannel <- msg.Body
		}
		close(messageChannel)
	}()

	return messageChannel, nil
}

// Disconnect - Closes the client connection(s)
func (rabbitMQClient *RabbitMQClient) Disconnect() error {
	err := rabbitMQClient.channel.Close()
	if err != nil {
		return err
	}

	err = rabbitMQClient.connection.Close()
	if err != nil {
		return err
	}

	return nil
}

// Restart - Restarts the client connection in the event of failure or error
func (rabbitMQClient *RabbitMQClient) Restart() (MessageClient, error) {
	// Close the connections on the old client
	err := rabbitMQClient.Disconnect()
	if err != nil {
		return nil, err
	}

	// Use the existing configurations from the old client
	existingConfigurations := config.RabbitMQConfigurations{
		URL:          rabbitMQClient.url,
		ExchangeName: rabbitMQClient.exchangeName,
		QueueName:    rabbitMQClient.queueName,
	}

	connection, channel, queue, err := setup(existingConfigurations)
	if err != nil {
		return nil, err
	}

	newClient := &RabbitMQClient{
		exchangeName: rabbitMQClient.exchangeName,
		queueName:    rabbitMQClient.queueName,
		url:          rabbitMQClient.url,
		connection:   connection,
		channel:      channel,
		queue:        queue,
	}

	return newClient, nil
}

func setup(rabbitMQConfigurations config.RabbitMQConfigurations) (*amqp.Connection, *amqp.Channel, amqp.Queue, error) {
	connection, err := connect(rabbitMQConfigurations.URL)
	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	channel, err := createChannel(connection)
	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	queue, err := initializeTopography(channel, rabbitMQConfigurations.ExchangeName, rabbitMQConfigurations.QueueName)
	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	return connection, channel, queue, nil
}

func connect(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func createChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, err
}

func initializeTopography(channel *amqp.Channel, exchangeName string, queueName string) (amqp.Queue, error) {
	err := channel.ExchangeDeclare(
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

	queue, err := channel.QueueDeclare(
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
	err = channel.QueueBind(
		queueName,    // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	return queue, nil
}

func selectConfigurations(configurations config.Configurations, messageDirection string) (config.RabbitMQConfigurations, error) {
	if messageDirection == "inbound" {
		return configurations.InboundRabbitMQConfigurations, nil
	}

	if messageDirection == "outbound" {
		return configurations.OutboundRabbitMQConfigurations, nil
	}

	return config.RabbitMQConfigurations{}, errors.New("Message direction must be set to either 'inbound' or 'outbound'")
}
