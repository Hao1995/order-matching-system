package matchingengine

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
	"github.com/Hao1995/order-matching-system/pkg/cf"
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

	var transactions []events.MatchingTransaction

	priceLevel := me.orderBook.GetPriceLevels(GetOppositeSide(order.Side))

	for !priceLevel.IsDummyNode {
		if !me.isPriceMatch(priceLevel.Price, order) {
			break
		}

		targetOrder := priceLevel.headOrder.Next
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

			if targetOrder.RemainingQuantity == 0 {
				priceLevel.Remove(targetOrder.ID)
				targetOrder = targetOrder.Next
			}
		}

		if priceLevel.IsEmpty() {
			me.orderBook.RemovePriceLevel(GetOppositeSide(order.Side), priceLevel.Price)
			priceLevel = priceLevel.Next
		} else {
			break
		}
	}

	if order.RemainingQuantity > 0 {
		me.orderBook.AddOrder(order)
	}

	return transactions
}

func (me *MatchingEngine) isPriceMatch(targetPrice float64, order *Order) bool {
	if order.Side == SideBUY {
		return order.Price >= targetPrice
	} else {
		return order.Price <= targetPrice
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
