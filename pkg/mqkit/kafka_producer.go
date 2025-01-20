package mqkit

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer is responsible sending order data to the matching engine.
type KafkaProducer struct {
	writer *kafka.Writer
	// key is used as symbol to guarantee in order
	key []byte
}

// NewKafkaProducer creates a new KafkaProducer
func NewKafkaProducer(brokers []string, topic string, key []byte) Producer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
		key: key,
	}
}

// Publish sends a message to the Kafka topic.
func (kp *KafkaProducer) Publish(ctx context.Context, val []byte) error {
	msg := kafka.Message{
		Key:   kp.key,
		Value: val,
	}
	if err := kp.writer.WriteMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}

// Close closes the Kafka writer.
func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
