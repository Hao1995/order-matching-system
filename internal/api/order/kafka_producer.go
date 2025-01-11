package order

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
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
func (op *KafkaProducer) Publish(ctx context.Context, event *events.OrderEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		logger.Logger.Error("failed to convert event to byte array", zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Value: bytes,
	}
	if err := op.writer.WriteMessages(ctx, msg); err != nil {
		logger.Logger.Error("failed to write message", zap.Error(err))
		return err
	}
	logger.Logger.Info("message sent", zap.ByteString("event", bytes))
	return nil
}

// Close closes the Kafka writer.
func (op *KafkaProducer) Close() error {
	logger.Logger.Info("closing Kafka Kafkaproducer")
	return op.writer.Close()
}
