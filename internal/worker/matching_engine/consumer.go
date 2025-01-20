package matchingengine

import (
	"context"
)

type Consumer interface {
	Consume(ctx context.Context, handler func(key []byte, val []byte) error) error
	Close() error
}
