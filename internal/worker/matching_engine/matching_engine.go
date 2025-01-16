package matchingengine

import (
	"time"

	"github.com/google/uuid"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
)

const (
	TOP_TICK_NUMBER = 5
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
	var transactions []events.MatchingTransaction

	priceLevel := me.orderBook.GetPriceLevels(GetOppositeSide(order.Side))

	for priceLevel != nil {
		if me.isPriceMatch(priceLevel.Price, order) {
			targetOrder := priceLevel.headOrder.Next
			for targetOrder != nil && order.RemainingQuantity > 0 {
				dealQuantity := min(order.RemainingQuantity, targetOrder.Quantity)

				order.RemainingQuantity -= dealQuantity
				priceLevel.Quantity -= dealQuantity
				targetOrder.RemainingQuantity -= dealQuantity

				transactions = append(transactions, events.MatchingTransaction{
					ID:          uuid.New().String(),
					Symbol:      order.Symbol,
					BuyOrderID:  order.ID,
					SellOrderID: targetOrder.ID,
					Price:       priceLevel.Price,
					Quantity:    dealQuantity,
					CreatedAt:   time.Now(),
				})

				oldOrder := targetOrder
				targetOrder = targetOrder.Next

				if targetOrder.RemainingQuantity == 0 {
					priceLevel.Remove(oldOrder.ID)
				}
			}
		}

		if priceLevel.IsEmpty() {
			me.orderBook.RemovePriceLevel(GetOppositeSide(order.Side), priceLevel.Price)
		}
	}

	if order.RemainingQuantity > 0 {
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

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
