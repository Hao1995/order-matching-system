package matchingengine

import (
	"errors"

	"github.com/Hao1995/order-matching-system/internal/common/models"
)

var (
	ErrNotMatchingPrice = errors.New("order price and current price are not matched")
	ErrOrderNotFound    = errors.New("order not found")
)

type OrderNode struct {
	*models.Order

	Next *OrderNode
	Prev *OrderNode
}

type PriceLevel struct {
	Price    float64
	Quantity int64

	// headOrder sorts by created_at in ascending order
	headOrder *OrderNode
	tailOrder *OrderNode

	// nodeHash stores the OrderNode by ID
	nodeHash map[string]*OrderNode

	Next *PriceLevel
	Prev *PriceLevel
}

func NewPriceLevel(price float64) *PriceLevel {
	return &PriceLevel{
		Price:     price,
		Quantity:  0,
		headOrder: &OrderNode{},
		tailOrder: &OrderNode{},
	}
}

// Add
func (pl *PriceLevel) Add(order *models.Order) error {
	if order.Price != pl.Price {
		return ErrNotMatchingPrice
	}
	pl.Quantity += order.RemainingQuantity

	tmpNode := pl.tailOrder.Prev
	newOrder := &OrderNode{
		Order: order,
		Prev:  tmpNode,
		Next:  pl.tailOrder,
	}
	tmpNode.Next = newOrder
	pl.tailOrder.Prev = newOrder
	return nil
}

// Remove
func (pl *PriceLevel) Remove(orderID string) error {
	orderNode, ok := pl.nodeHash[orderID]
	if !ok {
		return ErrOrderNotExist
	}

	prevNode := orderNode.Prev
	nextNode := orderNode.Next
	prevNode.Next = nextNode
	nextNode.Prev = prevNode

	pl.Quantity -= orderNode.RemainingQuantity
	orderNode = nil
	delete(pl.nodeHash, orderID)

	return nil
}

// IsEmpty checks if the price level has no orders
func (pl *PriceLevel) IsEmpty() bool {
	return pl.headOrder.Next == nil
}
