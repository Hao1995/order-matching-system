package matchingengine

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/pkg/logger"
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
func (op *KafkaConsumer) Consume(ctx context.Context) ([]byte, error) {
	msg, err := op.reader.ReadMessage(ctx)
	if err != nil {
		logger.Error("failed to write message", zap.Error(err))
		return []byte{}, err
	}
	logger.Info("received message", zap.ByteString("val", msg.Value))
	return msg.Value, nil
}

// Close closes the Kafka reader.
func (op *KafkaConsumer) Close() error {
	logger.Info("closing Kafka consumer")
	return op.reader.Close()
}
