package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os/signal"
	"syscall"

	"github.com/avast/retry-go/v4"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
	matchingengine "github.com/Hao1995/order-matching-system/internal/worker/matching_engine"
	"github.com/Hao1995/order-matching-system/pkg/logger"
	"github.com/Hao1995/order-matching-system/pkg/mqkit"
	"github.com/Hao1995/order-matching-system/pkg/pubsubkit"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
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

	// Message queue
	consumer := mqkit.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.App.OrderTopic)
	defer consumer.Close()

	// Pub/Sub
	publisher := pubsubkit.NewKafkaPublisher(cfg.Kafka.Brokers, cfg.App.MatchingTopic)
	defer publisher.Close()

	// OrderBook
	orderBook := matchingengine.NewOrderBook()

	// Matcher
	matcher := matchingengine.NewMatcher(orderBook, cfg.TickNum)
	logger.Info("success create a Kafka reader", zap.String("topic", cfg.App.OrderTopic))

	go func() {
		for {
			// Consume event
			handler := func(val []byte) error {
				var event events.Event
				if err := json.Unmarshal(val, &event); err != nil {
					logger.Error("failed to unmarshal event", zap.Error(err), zap.ByteString("val", val))
					return err
				}

				orderEvent, ok := event.Data.(events.OrderEvent)
				if !ok {
					logger.Error("unknown event type, skip the event", zap.Any("data", event.Data))
					return ErrUnknownEventType
				}

				var matching matchingengine.Matching
				var matchingEvent events.MatchingEvent
				switch event.EventType {
				case events.EventTypeCreateOrder:
					order := convertOrderEventToOrder(orderEvent)
					matching = matcher.CreateOrder(order)
					matchingEvent.Type = events.MatchingEventTypeCreate
				case events.EventTypeCancelOrder:
					var err error
					matching, err = matcher.CancelOrder(orderEvent.ID)
					if err != nil {
						logger.Error("failed to cancel order", zap.Error(err))
						return err
					}
					matchingEvent.Type = events.MatchingEventTypeCancel
				default:
					logger.Error("unknown event type", zap.String("eventType", event.EventType.String()))
					return ErrUnknownEventType
				}

				// Convert matching data to matching event
				matchingEvent.Order = orderEvent
				matchingEvent.Transactions = convertToTransactionEvents(matching.Transactions)
				matchingEvent.BuyTicks = convertToTickEvents(matching.BuyTicks)
				matchingEvent.SellTicks = convertToTickEvents(matching.SellTicks)

				// Publish matching event
				matchingMsg, err := json.Marshal(matchingEvent)
				if err != nil {
					logger.Error("failed to marshal matching event", zap.Error(err), zap.Any("matchingEvent", matchingEvent))
					return err
				}

				if err := publisher.Publish(matchingMsg); err != nil {
					logger.Error("failed to publish matching event", zap.Error(err))
					return err
				}
				return nil
			}

			// Retry consume messages by BackOffDelay
			if err := retry.Do(
				func() error {
					if err := consumer.Consume(context.Background(), handler); err != nil {
						logger.Warn("failed to consume event from Kafka", zap.Error(err))
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

func convertOrderEventToOrder(orderEvent events.OrderEvent) matchingengine.Order {
	return matchingengine.Order{
		ID:        orderEvent.ID,
		Symbol:    orderEvent.Symbol,
		Type:      matchingengine.OrderType(orderEvent.Type),
		Price:     orderEvent.Price,
		Quantity:  orderEvent.Quantity,
		CreatedAt: orderEvent.CreatedAt,
	}
}
