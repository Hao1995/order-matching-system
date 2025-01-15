package matchingengine

import (
	"errors"
)

var (
	ErrNotMatchingPrice = errors.New("order price and current price are not matched")
	ErrOrderNotFound    = errors.New("order not found")
)

type PriceLevel struct {
	Price    float64
	Quantity int64

	// headOrder sorts by created_at in ascending order
	headOrder *Order
	tailOrder *Order

	// orderByID stores the Order by ID
	orderByID map[string]*Order

	IsDummyNode bool
	Next        *PriceLevel
	Prev        *PriceLevel
}

func NewPriceLevel(price float64) *PriceLevel {
	headOrder, tailOrder := &Order{IsDummyNode: true}, &Order{IsDummyNode: true}
	headOrder.Next = tailOrder
	tailOrder.Prev = headOrder

	return &PriceLevel{
		Price:     price,
		Quantity:  0,
		orderByID: make(map[string]*Order),
		headOrder: headOrder,
		tailOrder: tailOrder,
	}
}

// GetFirstOrder retrieves the first order node
func (pl *PriceLevel) GetFirstOrder() *Order {
	return pl.headOrder.Next
}

// Add
func (pl *PriceLevel) Add(order *Order) error {
	if order.Price != pl.Price {
		return ErrNotMatchingPrice
	}
	pl.Quantity += order.RemainingQuantity

	currOrder := pl.tailOrder.Prev
	for !currOrder.IsDummyNode {
		if currOrder.CreatedAt.Before(order.CreatedAt) {
			break
		}
		currOrder = currOrder.Prev
	}

	tmpNode := currOrder.Next
	currOrder.Next = order
	order.Prev = currOrder
	tmpNode.Prev = order
	order.Next = tmpNode

	pl.orderByID[order.ID] = order

	return nil
}

// Remove
func (pl *PriceLevel) Remove(orderID string) error {
	order, ok := pl.orderByID[orderID]
	if !ok {
		return ErrOrderNotExist
	}

	prevNode := order.Prev
	nextNode := order.Next
	prevNode.Next = nextNode
	nextNode.Prev = prevNode

	pl.Quantity -= order.RemainingQuantity
	order = nil
	delete(pl.orderByID, orderID)

	return nil
}

// IsEmpty checks if the price level has no orders
func (pl *PriceLevel) IsEmpty() bool {
	return pl.headOrder.Next == nil
}
