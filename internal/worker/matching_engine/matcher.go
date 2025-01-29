package matchingengine

import (
	"errors"
	"time"

	"github.com/google/uuid"
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

type Matcher struct {
	orderBook *OrderBook
	tickNum   int8
}

func NewMatcher(orderBook *OrderBook, tickNum int8) *Matcher {
	return &Matcher{
		orderBook: orderBook,
		tickNum:   tickNum,
	}
}

// CancelOrder delete the order from the order book
func (me *Matcher) CancelOrder(orderID string) (Matching, error) {
	if err := me.orderBook.DeleteOrder(orderID); err != nil {
		return Matching{}, err
	}

	var matching Matching
	matching.BuyTicks, matching.SellTicks = me.orderBook.GetTopTicks(me.tickNum)
	return matching, nil
}

// CreateOrder inserts a new order
func (me *Matcher) CreateOrder(order Order) Matching {
	var matching Matching
	matching.Transactions = me.matchOrder(order)
	matching.BuyTicks, matching.SellTicks = me.orderBook.GetTopTicks(me.tickNum)
	return matching
}

// matchOrder attempts to match an incoming order with existing orders
func (me *Matcher) matchOrder(order Order) []Transaction {
	// Use two pointers to sync the updates back to the OrderBook
	var matchingLevels **PriceLevel
	var priceComparator func(float64, float64) bool
	transactions := []Transaction{}

	if order.Type == OrderTypeBuy {
		matchingLevels = &me.orderBook.SellLevels
		priceComparator = func(price1, price2 float64) bool { return price1 >= price2 }
	} else {
		matchingLevels = &me.orderBook.BuyLevels
		priceComparator = func(price1, price2 float64) bool { return price1 <= price2 }
	}

	for *matchingLevels != nil && priceComparator(order.Price, (*matchingLevels).Price) {
		currentLevel := *matchingLevels

		for currentLevel.HeadOrders != nil && order.Quantity > 0 {
			matchedQuantity := min(order.Quantity, currentLevel.HeadOrders.Order.Quantity)

			transactions = append(transactions, Transaction{
				ID:     getUUID(),
				Symbol: order.Symbol,
				BuyOrderID: func() string {
					if order.Type == OrderTypeBuy {
						return order.ID
					} else {
						return currentLevel.HeadOrders.Order.ID
					}
				}(),
				SellOrderID: func() string {
					if order.Type == OrderTypeSell {
						return order.ID
					} else {
						return currentLevel.HeadOrders.Order.ID
					}
				}(),
				Price:     currentLevel.Price,
				Quantity:  matchedQuantity,
				CreatedAt: now(),
			})

			order.Quantity -= matchedQuantity
			currentLevel.HeadOrders.Order.Quantity -= matchedQuantity

			if currentLevel.HeadOrders.Order.Quantity == 0 {
				nextOrder := currentLevel.HeadOrders.Next
				me.orderBook.DeleteOrder(currentLevel.HeadOrders.Order.ID)
				currentLevel.HeadOrders = nextOrder
			}
		}

		if order.Quantity == 0 {
			break
		}
	}

	if order.Quantity > 0 {
		me.orderBook.InsertOrder(order)
	}

	return transactions
}
