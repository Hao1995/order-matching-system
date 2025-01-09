package usecases

import (
	"context"

	"github.com/Hao1995/order-matching-system/pkg/models/events"
)

type Producer interface {
	Publish(ctx context.Context, topic string, event *events.OrderEvent) error
	Close() error
}
