package order

import (
	"context"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
)

type Producer interface {
	Publish(ctx context.Context, event *events.OrderEvent) error
	Close() error
}
