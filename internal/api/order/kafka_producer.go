package order

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/pkg/logger"
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
		logger.Error("failed to write message", zap.Error(err))
		return err
	}
	logger.Info("message sent", zap.ByteString("val", val))
	return nil
}

// Close closes the Kafka writer.
func (op *KafkaProducer) Close() error {
	logger.Info("closing Kafka producer")
	return op.writer.Close()
}
