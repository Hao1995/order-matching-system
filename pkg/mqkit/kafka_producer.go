package mqkit

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer is responsible sending order data to the matching engine.
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new KafkaProducer
func NewKafkaProducer(writer *kafka.Writer) Producer {
	return &KafkaProducer{
		writer: writer,
	}
}

// Publish sends a message to the Kafka topic.
func (op *KafkaProducer) Publish(ctx context.Context, val []byte) error {
	msg := kafka.Message{
		Value: val,
	}
	if err := op.writer.WriteMessages(ctx, msg); err != nil {
		return err
	}
	return nil
}

// Close closes the Kafka writer.
func (op *KafkaProducer) Close() error {
	return op.writer.Close()
}
