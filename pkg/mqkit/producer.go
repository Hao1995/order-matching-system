package mqkit

import (
	"context"
)

type Producer interface {
	Publish(ctx context.Context, val []byte) error
	Close() error
}
