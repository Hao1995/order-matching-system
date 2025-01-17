package matchingengine

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
	"github.com/Hao1995/order-matching-system/pkg/cf"
	"github.com/Hao1995/order-matching-system/pkg/logger"
)

const (
	TOP_TICK_NUMBER = 5
)

var (
	getUUID = func() string {
		return uuid.NewString()
	}

	now = func() time.Time {
		return time.Now()
	}

	ErrNoBuyOrder = errors.New("no buy order exist")
)

type MatchingEngine struct {
	orderBook *OrderBook
}

func NewMatchingEngine(orderBook *OrderBook) *MatchingEngine {
	return &MatchingEngine{
		orderBook: orderBook,
	}
}

func (me *MatchingEngine) CancelOrder(orderID string) error {
	return me.orderBook.RemoveOrder(orderID)
}

func (me *MatchingEngine) PlaceOrder(order *Order) []events.MatchingTransaction {
	logger.Debug("handling order ...", zap.Any("order", *order))

	var transactions []events.MatchingTransaction

	priceLevel := me.orderBook.GetPriceLevels(GetOppositeSide(order.Side))

	for !priceLevel.IsDummyNode {
		logger.Debug("valid priceLevel\n", zap.Any("priceLevel", *priceLevel))

		if !me.isPriceMatch(priceLevel.Price, order) {
			logger.Debug("price is not matched, break")
		}

		targetOrder := priceLevel.headOrder.Next
		logger.Debug("get targetOrder", zap.Any("targetOrder", *targetOrder))

		for !targetOrder.IsDummyNode && order.RemainingQuantity > 0 {
			dealQuantity := min(order.RemainingQuantity, targetOrder.Quantity)

			order.RemainingQuantity -= dealQuantity
			priceLevel.Quantity -= dealQuantity
			targetOrder.RemainingQuantity -= dealQuantity

			transaction := events.MatchingTransaction{
				ID:          getUUID(),
				Symbol:      order.Symbol,
				BuyOrderID:  me.getOrderBySide(SideBUY, order, targetOrder).ID,
				SellOrderID: me.getOrderBySide(SideSELL, order, targetOrder).ID,
				Price:       cf.Min(order.Price, targetOrder.Price),
				Quantity:    dealQuantity,
				CreatedAt:   now(),
			}
			transactions = append(transactions, transaction)

			targetOrderID := targetOrder.ID
			logger.Debug("match!", zap.Any("transaction", transaction))

			if targetOrder.RemainingQuantity == 0 {
				logger.Debug("remove targetOrder", zap.String("targetOrderID", targetOrderID))
				priceLevel.Remove(targetOrderID)
				targetOrder = targetOrder.Next
			} else {
				logger.Debug("current order is fully matched", zap.Int64("remainingQuantity", order.RemainingQuantity))
			}
		}

		if priceLevel.IsEmpty() {
			logger.Debug("current priceLevel is fully matched, remove it", zap.Float64("price", priceLevel.Price))
			me.orderBook.RemovePriceLevel(GetOppositeSide(order.Side), priceLevel.Price)
			priceLevel = priceLevel.Next
		} else {
			logger.Debug("the remaining orders in the current priceLevel can't be matched, break", zap.Float64("price", priceLevel.Price))
			break
		}
	}

	if order.RemainingQuantity > 0 {
		logger.Debug("add current order into the OrderBook", zap.Int64("remainingQuantity", order.RemainingQuantity))
		me.orderBook.AddOrder(order)
	}

	return transactions
}

func (me *MatchingEngine) isPriceMatch(targetPrice float64, order *Order) bool {
	if order.Side == SideBUY {
		return targetPrice <= order.Price
	} else {
		return order.Price >= targetPrice
	}
}

func (me *MatchingEngine) getOrderBySide(side Side, order1, order2 *Order) *Order {
	if order1.Side == side {
		return order1
	} else {
		return order2
	}
}

func (me *MatchingEngine) GetTopTicks(side Side, k int8) []events.MatchingTick {
	var topTicks []events.MatchingTick

	var priceLevel *PriceLevel
	if side == SideBUY {
		priceLevel = me.orderBook.buyHead.Next
	} else {
		priceLevel = me.orderBook.sellHead.Next
	}

	var count int8
	for priceLevel != nil && count < k {
		topTicks = append(topTicks, events.MatchingTick{
			Price:    priceLevel.Price,
			Quantity: priceLevel.Quantity,
		})
		priceLevel = priceLevel.Next
		count++
	}

	return topTicks
}
