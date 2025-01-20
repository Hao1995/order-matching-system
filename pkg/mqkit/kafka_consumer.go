package mqkit

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumer is responsible sending order data to the matching engine.
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer creates a new KafkaConsumer
func NewKafkaConsumer(reader *kafka.Reader) Consumer {
	return &KafkaConsumer{
		reader: reader,
	}
}

// Consume sends a message to the Kafka topic.
func (op *KafkaConsumer) Consume(ctx context.Context, handler func(key []byte, val []byte) error) error {

	msg, err := op.reader.FetchMessage(ctx)
	if err != nil {
		return err
	}

	if err := handler(msg.Key, msg.Value); err != nil {
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
