package messageproducer

type MessageProducer interface {
	Disconnect() error
	Send(message []byte) error
}
