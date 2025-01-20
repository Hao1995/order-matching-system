package pubsubkit

// Publisher defines the interface for publishing messages
type Publisher interface {
	Publish(value []byte) error
	Close() error
}

// Subscriber defines the interface for subscribing to messages
type Subscriber interface {
	Subscribe(handler func(value []byte) error) error
	Close() error
}
