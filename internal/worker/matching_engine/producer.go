package matchingengine

import (
	"context"
)

type Consumer interface {
	Consume(ctx context.Context) ([]byte, error)
	Close() error
}
