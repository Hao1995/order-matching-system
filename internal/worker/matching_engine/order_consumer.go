package matchingengine

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
	"github.com/Hao1995/order-matching-system/pkg/logger"
)

// OrderConsumer is responsible for receiving data from the message queue.
type OrderConsumer struct {
	r *kafka.Reader
}

// Consume receives a message from a Kafka topic
func (oc *OrderConsumer) Consume(ctx context.Context) (events.OrderEvent, error) {
	var event events.OrderEvent

	msg, err := oc.r.ReadMessage(ctx)
	if err != nil {
		logger.Error("failed to consume the message", zap.Error(err))
		return event, err
	}

	if err := json.Unmarshal(msg.Value, &event); err != nil {
		logger.Error("failed to unmarshal the message", zap.Error(err))
		return event, err
	}

	return event, nil
}

// Close closes the kafka consumer
func (oc *OrderConsumer) Close() error {
	logger.Info("closing kafka consumer")
	return oc.r.Close()
}
