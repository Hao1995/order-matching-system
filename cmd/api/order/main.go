package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/internal/api/order"
	"github.com/Hao1995/order-matching-system/pkg/logger"
)

func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("failed to parse config", err)
	}
}

func main() {
	defer logger.Sync()

	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// Kafka
	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Brokers...),
		Topic:                  cfg.Kafka.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	logger.Info("success create a Kafka writer", zap.String("topic", cfg.Kafka.Topic))
	defer w.Close()

	kafkaProducer := order.NewKafkaProducer(w)

	// Init Gin Router
	hlr := order.NewHandler(kafkaProducer, cfg.Kafka.Topic)
	router := gin.Default()
	order.RegisterRoutes(router, hlr)

	RunGinServer(ctx, stop, router)
}

func RunGinServer(ctx context.Context, stop context.CancelFunc, router *gin.Engine) {
	srv := &http.Server{
		Addr:    cfg.App.Port,
		Handler: router,
	}

	// Init a goroutine to run the server so that it won't block the graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server closed", zap.Error(err))
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	stop()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5s to finish the requests that are currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", zap.Error(err))
	}

	logger.Info("Server exiting")
}
