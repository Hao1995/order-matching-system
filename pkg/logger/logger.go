package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
	once   sync.Once
)

// InitLogger initials logger instance
func InitLogger() error {
	var err error
	once.Do(func() {
		Logger, err = zap.NewProduction()
	})
	return err
}

// Sync flushes any buffered log entries
func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}
