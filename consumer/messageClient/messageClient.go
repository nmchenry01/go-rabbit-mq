package messageclient

// MessageClient - The interface the represents the functionality needed for a client to consume from a message bus
type MessageClient interface {
	Disconnect() error
	Restart() (MessageClient, error)
	Consume() (chan []byte, error)
}
