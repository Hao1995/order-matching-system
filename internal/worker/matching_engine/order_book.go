package matchingengine

import (
	"errors"

	"github.com/Hao1995/order-matching-system/internal/common/models"
)

var (
	ErrOrderNotExist      = errors.New("order not exists in the order book")
	ErrPriceLevelNotExist = errors.New("price_level not exists in the order book")
)

type OrderBook struct {
	buyHead *PriceLevel
	buyTail *PriceLevel

	sellHead *PriceLevel
	sellTail *PriceLevel

	// priceLevelHash stores PriceLevel by order id
	priceLevelHash map[string]*PriceLevel

	buyPriceHash  map[float64]*PriceLevel
	sellPriceHash map[float64]*PriceLevel
}

func (pl *OrderBook) GetPriceLevels(side models.Side) *PriceLevel {
	if side == models.SideBUY {
		return pl.buyHead.Next
	} else {
		return pl.sellHead.Next
	}
}

// Add
func (pl *OrderBook) Add(order *models.Order) error {
	var priceLevel *PriceLevel
	if order.Side == models.SideBUY {
		priceLevel = pl.buyHead
	} else {
		priceLevel = pl.sellHead
	}

	for priceLevel != nil {
		if priceLevel.Price == order.Price {
			priceLevel.Add(order)
			pl.priceLevelHash[order.ID] = priceLevel
			return nil
		}

		if order.Side == models.SideBUY {
			if priceLevel.Price < order.Price {
				break
			}
		} else {
			if priceLevel.Price > order.Price {
				break
			}
		}

		priceLevel = priceLevel.Next
	}

	newPriceLevel := NewPriceLevel(order.Price)
	newPriceLevel.Add(order)

	tmpNode := priceLevel.Prev
	tmpNode.Next = newPriceLevel
	newPriceLevel.Prev = tmpNode
	newPriceLevel.Next = priceLevel
	priceLevel.Prev = newPriceLevel
	return nil
}

// RemoveOrder
func (pl *OrderBook) RemoveOrder(orderID string) error {
	priceLevel, found := pl.priceLevelHash[orderID]
	if !found {
		return ErrOrderNotExist
	}

	priceLevel.Remove(orderID)

	if priceLevel.IsEmpty() {
		tmpNode := priceLevel.Prev
		nextNode := priceLevel.Next
		tmpNode.Next = nextNode
		nextNode.Prev = tmpNode
	}

	delete(pl.priceLevelHash, orderID)
	return nil
}

// RemovePriceLevel
func (pl *OrderBook) RemovePriceLevel(side models.Side, price float64) error {
	var priceHash map[float64]*PriceLevel
	if side == models.SideBUY {
		priceHash = pl.buyPriceHash
	} else {
		priceHash = pl.sellPriceHash
	}

	priceLevel, found := priceHash[price]
	if !found {
		return ErrOrderNotExist
	}

	prevNode := priceLevel.Prev
	nextNode := priceLevel.Next
	prevNode.Next = nextNode
	nextNode.Prev = prevNode

	if side == models.SideBUY {
		delete(pl.buyPriceHash, price)
	} else {
		delete(pl.sellPriceHash, price)
	}

	return nil
}
