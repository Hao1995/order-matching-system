package order

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hao1995/order-matching-system/internal/api/order"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

// @todo: extract as config
const (
	MQ_TOPIC = "APPLE_ORDER"
	PORT     = ":8080"
)

var (
	MQ_CONNECTIONS = []string{
		"localhost:9092", "localhost:9093", "localhost:9094",
	}
)

func main() {
	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// init KafkaPublisher
	w := &kafka.Writer{
		Addr:     kafka.TCP(MQ_CONNECTIONS...),
		Topic:    MQ_TOPIC,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	kafkaProducer := order.NewKafkaProducer(w)

	// Init Gin Router
	hlr := order.NewHandler(kafkaProducer, MQ_TOPIC)
	router := gin.Default()
	order.RegisterRoutes(router, hlr)

	RunGinServer(ctx, stop, router)
}

func RunGinServer(ctx context.Context, stop context.CancelFunc, router *gin.Engine) {
	srv := &http.Server{
		Addr:    PORT,
		Handler: router,
	}

	// Init a goroutine to run the server so that it won't block the graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5s to finish the requests that are currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
