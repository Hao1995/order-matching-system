package logger

import (
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func init() {
	var err error
	once.Do(func() {
		env := os.Getenv("APP_ENV")
		var cfg zap.Config

		switch env {
		case "production":
			cfg = zap.NewProductionConfig()
		case "development":
			cfg = zap.NewDevelopmentConfig()
		default:
			cfg = zap.NewDevelopmentConfig()
			cfg.OutputPaths = []string{
				"stdout",
			}
		}
		logger, err = cfg.Build()
	})

	if err != nil {
		log.Fatal("failed to init logger", err)
	}
}

// Sync flushes any buffered log entries
func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

// Debug logs a message at DebugLevel
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Panic logs a message at PanicLevel.
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel.
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
