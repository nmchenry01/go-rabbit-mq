package messageproducer

import (
	"errors"

	"github.com/nmchenry/go-rabbit-mq/producer/config"
	"github.com/streadway/amqp"
)

// RabbitMQProducer - Implementation of a message producer for RabbitMQ
type RabbitMQProducer struct {
	url          string
	count        int
	exchangeName string
	connection   *amqp.Connection
	channel      *amqp.Channel
}

// NewRabbitMQProducer - Creates a new message producer
func NewRabbitMQProducer(configurations config.Configurations, messageDirection string) (MessageProducer, error) {
	messageProducerConfigurations, err := selectConfigurations(configurations, messageDirection)
	if err != nil {
		return nil, err
	}

	connection, channel, err := setup(messageProducerConfigurations)

	newProducer := &RabbitMQProducer{
		url:          messageProducerConfigurations.URL,
		count:        messageProducerConfigurations.Count,
		exchangeName: messageProducerConfigurations.ExchangeName,
		connection:   connection,
		channel:      channel,
	}

	return newProducer, nil
}

// Send - Sends a message to the queue
func (rabbitMQProducer *RabbitMQProducer) Send(message []byte) error {
	// Send message(s)
	for i := 0; i < rabbitMQProducer.count; i++ {
		err := rabbitMQProducer.channel.Publish(
			rabbitMQProducer.exchangeName, // exchange
			"",                            // routing key
			false,                         // mandatory
			false,                         // immediate
			amqp.Publishing{
				ContentType: "application/xml",
				Body:        message,
			})

		if err != nil {
			return err
		}
	}

	return nil
}

// Disconnect - Disconnects from Queue
func (rabbitMQProducer *RabbitMQProducer) Disconnect() error {
	err := rabbitMQProducer.channel.Close()
	if err != nil {
		return err
	}

	err = rabbitMQProducer.connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func selectConfigurations(configurations config.Configurations, messageDirection string) (config.ProducerConfigurations, error) {
	if messageDirection == "inbound" {
		return configurations.InboundProducerConfigurations, nil
	}

	if messageDirection == "outbound" {
		return configurations.OutboundProducerConfigurations, nil
	}

	return config.ProducerConfigurations{}, errors.New("Message direction must be set to either 'inbound' or 'outbound'")
}

func setup(messageProducerConfigurations config.ProducerConfigurations) (*amqp.Connection, *amqp.Channel, error) {
	connection, err := connect(messageProducerConfigurations.URL)
	if err != nil {
		return nil, nil, err
	}

	channel, err := createChannel(connection)
	if err != nil {
		return nil, nil, err
	}

	err = createExchange(channel, messageProducerConfigurations.ExchangeName)
	if err != nil {
		return nil, nil, err
	}

	return connection, channel, nil
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

func createExchange(channel *amqp.Channel, exchangeName string) error {
	err := channel.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	return nil
}
