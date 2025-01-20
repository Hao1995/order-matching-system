package pubsubkit

// Publisher defines the interface for publishing messages
type Publisher interface {
	Publish(topic string, key, value []byte) error
}

// Subscriber defines the interface for subscribing to messages
type Subscriber interface {
	Subscribe(topic string, handler func(key, value []byte) error) error
	Close() error
}
