package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	once   sync.Once
)

// InitLogger initials logger instance
func InitLogger() error {
	var err error
	once.Do(func() {
		logger, err = zap.NewProduction()
		sugar = logger.Sugar()
	})
	return err
}

// Sync flushes any buffered log entries
func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}
