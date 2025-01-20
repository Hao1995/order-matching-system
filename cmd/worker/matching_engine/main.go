package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/caarlos0/env/v11"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
	matchingengine "github.com/Hao1995/order-matching-system/internal/worker/matching_engine"
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
	defer stop()

	// Kafka
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		MaxWait: 3 * time.Second,
	})
	defer r.Close()
	consumer := matchingengine.NewKafkaConsumer(r)

	// OrderBook
	orderBook := matchingengine.NewOrderBook()

	// Matcher
	matcher := matchingengine.NewMatcher(orderBook)
	logger.Info("success create a Kafka reader", zap.String("topic", cfg.Kafka.Topic))

	go func() {
		for {
			// Consume event
			var event events.Event
			if err := retry.Do(
				func() error {
					val, err := consumer.Consume(context.Background())
					if err != nil {
						logger.Warn("failed to consume event from Kafka", zap.Error(err))
						return err
					}
					if err := json.Unmarshal(val, &event); err != nil {
						logger.Warn("failed to unmarshal event", zap.Error(err), zap.ByteString("val", val))
						return err
					}
					return nil
				},
				retry.Attempts(3),
			); err != nil {
				// Temporary leave end the service
				logger.Error("retry error achieve the max limit")
				return
			}

			orderEvent, ok := event.Data.(events.OrderEvent)
			if !ok {
				logger.Error("unexpected event order type, skip the event", zap.Any("data", event.Data))
				continue
			}

			data := orderEvent
			switch event.EventType {
			case events.EventTypeCreateOrder:
				order := matchingengine.Order{
					ID:        data.ID,
					Symbol:    data.Symbol,
					Type:      matchingengine.OrderType(data.Type),
					Price:     data.Price,
					Quantity:  data.Quantity,
					CreatedAt: data.CreatedAt,
				}

				transactions := matcher.MatchOrder(order)

				buyTicks, sellTicks := matcher.GetTopTicks(int8(cfg.TickNum))

				matchingEvent := events.MatchingEvent{
					Type:         events.MatchingEventTypeMatching,
					Order:        data,
					Transactions: convertToTransactionEvents(transactions),
					BuyTicks:     convertToTickEvents(buyTicks),
					SellTicks:    convertToTickEvents(sellTicks),
				}

				// @todo: ack and publish to matching Topic
				fmt.Println(matchingEvent)
			}
		}
	}()

	<-ctx.Done()
	logger.Info("received interrupt signals from the OS, end the process")
}

func convertToTransactionEvents(transactions []matchingengine.Transaction) []events.TransactionEvent {
	result := make([]events.TransactionEvent, 0, len(transactions))
	for _, transaction := range transactions {
		result = append(result, events.TransactionEvent{
			ID:          transaction.ID,
			Symbol:      transaction.Symbol,
			BuyOrderID:  transaction.BuyOrderID,
			SellOrderID: transaction.SellOrderID,
			Price:       transaction.Price,
			Quantity:    transaction.Quantity,
			CreatedAt:   transaction.CreatedAt,
		})
	}
	return result
}

func convertToTickEvents(ticks []matchingengine.Tick) []events.TickEvent {
	result := make([]events.TickEvent, 0, len(ticks))
	for _, tick := range ticks {
		result = append(result, events.TickEvent{
			Price:    tick.Price,
			Quantity: tick.Quantity,
		})
	}
	return result
}
