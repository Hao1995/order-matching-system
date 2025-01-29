package mqkit

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer is responsible sending order data to the matching engine.
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer creates a new KafkaConsumer
func NewKafkaConsumer(brokers []string, topic string) Consumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			MaxWait: 3 * time.Second,
		}),
	}
}

// Consume receives a message to the Kafka topic and manually
// commit messages after completely executing the handler
func (op *KafkaConsumer) Consume(ctx context.Context, handler func(val []byte) error) error {
	msg, err := op.reader.FetchMessage(ctx)
	if err != nil {
		return err
	}

	if err := handler(msg.Value); err != nil {
		return err
	}

	if err := op.reader.CommitMessages(ctx, msg); err != nil {
		return err
	}

	return nil
}

// Close closes the Kafka reader.
func (op *KafkaConsumer) Close() error {
	return op.reader.Close()
}
