package repositories

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/internal/api/order/usecases"
	"github.com/Hao1995/order-matching-system/pkg/models/events"
)

// OrderProducer is responsible sending order data to the matching engine.
type OrderProducer struct {
	writer *kafka.Writer
}

// NewOrderProducer creates a new OrderProducer
func NewOrderProducer(writer *kafka.Writer) usecases.Producer {
	return &OrderProducer{
		writer: writer,
	}
}

// Publish sends a message to the Kafka topic.
func (op *OrderProducer) Publish(ctx context.Context, topic string, event *events.OrderEvent) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		zap.L().Error("failed to convert event to byte array", zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Value: bytes,
	}
	if err := op.writer.WriteMessages(ctx, msg); err != nil {
		zap.L().Error("failed to write message", zap.Error(err))
		return err
	}
	zap.L().Info("message sent", zap.String("topic", topic), zap.ByteString("event", bytes))
	return nil
}

// Close closes the Kafka writer.
func (op *OrderProducer) Close() error {
	zap.L().Info("closing Kafka producer")
	return op.writer.Close()
}
