package messageproducer

// MessageProducer - Interface for message producers
type MessageProducer interface {
	Disconnect() error
	Send(message []byte) error
}
