package mqkit

import (
	"context"
)

type Consumer interface {
	Consume(ctx context.Context, handler func(val []byte) error) error
	Close() error
}
