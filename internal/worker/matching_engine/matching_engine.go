package matchingengine

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Hao1995/order-matching-system/internal/common/models"
	"github.com/Hao1995/order-matching-system/internal/common/models/events"
)

const (
	TOP_TICK_NUMBER = 5
)

type MatchingEngine struct {
	OrderBook *OrderBook
}

func (me *MatchingEngine) CancelOrder(ctx context.Context, orderID string) events.Matching {
	var matching events.Matching

	// handle incoming order
	matching.Order = events.OrderCancelEvent{
		ID: orderID,
	}

	// handle matching
	me.OrderBook.RemoveOrder(orderID)

	// Get top ticks
	matching.TopBuy = me.GetTopTicks(models.SideBUY, TOP_TICK_NUMBER)
	matching.TopSell = me.GetTopTicks(models.SideSELL, TOP_TICK_NUMBER)

	return matching
}

func (me *MatchingEngine) PlaceOrder(ctx context.Context, order *models.Order) events.Matching {
	var matching events.Matching

	// handle incoming order
	matching.Order = events.OrderCreateEvent{
		ID:                order.ID,
		Symbol:            order.Symbol,
		Side:              events.Side(order.Side),
		Price:             order.Price,
		Quantity:          order.Quantity,
		RemainingQuantity: order.Quantity,
		CanceledQuantity:  0,
		CreatedAt:         order.CreatedAt,
		UpdatedAt:         order.UpdatedAt,
	}

	// handle matching
	matching.Transactions = me.match(order)

	// Get top ticks
	matching.TopBuy = me.GetTopTicks(models.SideBUY, TOP_TICK_NUMBER)
	matching.TopSell = me.GetTopTicks(models.SideSELL, TOP_TICK_NUMBER)

	return matching
}

func (me *MatchingEngine) match(order *models.Order) []events.MatchingTransaction {
	var transactions []events.MatchingTransaction

	priceLevel := me.OrderBook.GetPriceLevels(models.GetOppositeSide(order.Side))

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
			me.OrderBook.RemovePriceLevel(models.GetOppositeSide(order.Side), priceLevel.Price)
		}
	}

	if order.RemainingQuantity > 0 {
		me.OrderBook.Add(order)
	}

	return transactions
}

func (me *MatchingEngine) isPriceMatch(targetPrice float64, order *models.Order) bool {
	if order.Side == models.SideBUY {
		return targetPrice <= order.Price
	} else {
		return order.Price >= targetPrice
	}
}

func (me *MatchingEngine) GetTopTicks(side models.Side, k int8) []events.MatchingTick {
	var topTicks []events.MatchingTick

	var priceLevel *PriceLevel
	if side == models.SideBUY {
		priceLevel = me.OrderBook.buyHead.Next
	} else {
		priceLevel = me.OrderBook.sellHead.Next
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
