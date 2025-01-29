package mqkit

import (
	"context"
)

type Consumer interface {
	// Consume receives a message from the message broker and manually
	// commit messages after completely executing the handler
	Consume(ctx context.Context, handler func(val []byte) error) error
	// Close closes the stream, preventing the program from reading any more
	// messages from it.
	Close() error
}
