package matchingengine

import (
	"errors"
	"time"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrPriceLevelNotFound = errors.New("price level not found")
)

// Transaction represents the details of a matched order
type Transaction struct {
	ID          string
	Symbol      string
	BuyOrderID  string
	SellOrderID string
	Price       float64
	Quantity    int64
	CreatedAt   time.Time
}

// Tick represents the total quantity of a price
type Tick struct {
	Price    float64
	Quantity int64
}

// ENUM(Buy, Sell)
type OrderType int

// Order represents a buy or sell order
type Order struct {
	ID                string
	Symbol            string
	Type              OrderType
	Price             float64
	Quantity          int64
	RemainingQuantity int64
	CanceledQuantity  int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type OrderNode struct {
	Order      Order
	Next       *OrderNode
	Prev       *OrderNode
	PriceLevel *PriceLevel
}

// PriceLevel represents a price point in the order book
type PriceLevel struct {
	Type          OrderType
	Price         float64
	TotalQuantity int64
	// HeadOrders is the head of the doubly linked list of orders
	HeadOrders *OrderNode
	// TailOrders is the tail of the doubly linked list of orders
	TailOrders *OrderNode
	Prev       *PriceLevel
	Next       *PriceLevel
}

// OrderBook maintains buy and sell price levels along with auxiliary maps
type OrderBook struct {
	// BuyLevels is the head of the buy PriceLevels (sorted descending price)
	BuyLevels *PriceLevel
	// SellLevels is the Head of the sell PriceLevels (sorted ascending price)
	SellLevels *PriceLevel
	// orderMap maps Order.ID to OrderNode
	orderMap map[string]*OrderNode
	// buyPriceMap maps Buy Price to PriceLevel
	buyPriceMap map[float64]*PriceLevel
	// sellPriceMap maps Sell Price to PriceLevel
	sellPriceMap map[float64]*PriceLevel
}

// NewOrderBook initializes and returns a new OrderBook
func NewOrderBook() *OrderBook {
	return &OrderBook{
		BuyLevels:    nil,
		SellLevels:   nil,
		orderMap:     make(map[string]*OrderNode),
		buyPriceMap:  make(map[float64]*PriceLevel),
		sellPriceMap: make(map[float64]*PriceLevel),
	}
}

// InsertOrder inserts an order into the order book in the correct position
func (ob *OrderBook) InsertOrder(order Order) {
	if order.Type == OrderTypeBuy {
		ob.BuyLevels = ob.insertOrderToPriceLevel(ob.BuyLevels, order, ob.buyPriceMap, true)
	} else {
		ob.SellLevels = ob.insertOrderToPriceLevel(ob.SellLevels, order, ob.sellPriceMap, false)
	}
}

// insertOrderToPriceLevel inserts an order into the buy or sell price levels
func (ob *OrderBook) insertOrderToPriceLevel(headPriceLevel *PriceLevel, order Order, priceMap map[float64]*PriceLevel, isBuy bool) *PriceLevel {
	newOrderNode := &OrderNode{Order: order}

	// No PriceLevel: Create a new price level at the head
	if headPriceLevel == nil {
		newLevel := &PriceLevel{
			Type:          order.Type,
			Price:         order.Price,
			TotalQuantity: order.RemainingQuantity,
			HeadOrders:    newOrderNode,
			TailOrders:    newOrderNode,
			Next:          headPriceLevel,
		}

		if isBuy {
			ob.BuyLevels = newLevel
		} else {
			ob.SellLevels = newLevel
		}

		priceMap[order.Price] = newLevel
		ob.orderMap[order.ID] = newOrderNode
		newOrderNode.PriceLevel = newLevel
		return newLevel
	}

	currPriceLevel := headPriceLevel
	prevPriceLevel := currPriceLevel.Prev
	for currPriceLevel != nil {
		if currPriceLevel.Price == order.Price {
			// Find the same PriceLevel: Insert order to the tail
			currPriceLevel.TailOrders.Next = newOrderNode
			newOrderNode.Prev = currPriceLevel.TailOrders
			currPriceLevel.TailOrders = newOrderNode
			currPriceLevel.TotalQuantity += newOrderNode.Order.Quantity

			ob.orderMap[order.ID] = newOrderNode
			return headPriceLevel
		}

		prevPriceLevel = currPriceLevel
		currPriceLevel = currPriceLevel.Next
	}

	// PriceLevel not found: Create a new PriceLevel
	newLevel := &PriceLevel{
		Type:          order.Type,
		Price:         order.Price,
		TotalQuantity: order.RemainingQuantity,
		HeadOrders:    newOrderNode,
		TailOrders:    newOrderNode,
		Prev:          prevPriceLevel,
		Next:          nil,
	}

	prevPriceLevel.Next = newLevel

	priceMap[order.Price] = newLevel
	ob.orderMap[order.ID] = newOrderNode
	newOrderNode.PriceLevel = newLevel

	return headPriceLevel
}

// DeleteOrder deletes an order by ID in O(1) time
func (ob *OrderBook) DeleteOrder(orderID string) error {
	orderNode, exists := ob.orderMap[orderID]
	if !exists {
		return errors.New("order not found")
	}

	pl := orderNode.PriceLevel
	if pl == nil {
		return errors.New("price level not found")
	}

	// Adjust total quantity
	pl.TotalQuantity -= orderNode.Order.RemainingQuantity
	orderNode.Order.CanceledQuantity += orderNode.Order.RemainingQuantity
	orderNode.Order.RemainingQuantity = 0

	// Remove OrderNode from the orders linked list
	if orderNode.Prev != nil {
		orderNode.Prev.Next = orderNode.Next
	} else {
		pl.HeadOrders = orderNode.Next
	}

	if orderNode.Next != nil {
		orderNode.Next.Prev = orderNode.Prev
	} else {
		pl.TailOrders = orderNode.Prev
	}

	// Remove from orderMap
	delete(ob.orderMap, orderID)

	// Check if PriceLevel is empty
	if pl.HeadOrders == nil {
		ob.deletePriceLevel(pl)
	}

	return nil
}

// deletePriceLevel deletes a PriceLevel from the order book
func (ob *OrderBook) deletePriceLevel(pl *PriceLevel) {
	if pl.Type == OrderTypeBuy {
		if pl.Prev != nil {
			pl.Prev.Next = pl.Next
		} else {
			ob.BuyLevels = pl.Next
		}

		if pl.Next != nil {
			pl.Next.Prev = pl.Prev
		}
		delete(ob.buyPriceMap, pl.Price)
	} else {
		if pl.Prev != nil {
			pl.Prev.Next = pl.Next
		} else {
			ob.SellLevels = pl.Next
		}

		if pl.Next != nil {
			pl.Next.Prev = pl.Prev
		}
		delete(ob.sellPriceMap, pl.Price)
	}
}

// GetTopTicks returns the top N buy and sell ticks
func (ob *OrderBook) GetTopTicks(n int8) ([]Tick, []Tick) {
	buyTicks := []Tick{}
	sellTicks := []Tick{}

	current := ob.BuyLevels
	for i := 0; i < int(n) && current != nil; i++ {
		buyTicks = append(buyTicks, Tick{Price: current.Price, Quantity: current.TotalQuantity})
		current = current.Next
	}

	current = ob.SellLevels
	for i := 0; i < int(n) && current != nil; i++ {
		sellTicks = append(sellTicks, Tick{Price: current.Price, Quantity: current.TotalQuantity})
		current = current.Next
	}

	return buyTicks, sellTicks
}
