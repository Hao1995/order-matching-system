package matchingengine

import (
	"errors"
	"math"
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

	// priceLevelByOrderID stores PriceLevel by order id
	priceLevelByOrderID map[string]*PriceLevel

	buyNodeByPrice  map[float64]*PriceLevel
	sellNodeByPrice map[float64]*PriceLevel
}

func NewOrderBook() *OrderBook {
	headBuyPriceLevel, tailBuyPriceLevel := NewPriceLevel(math.MaxFloat64), NewPriceLevel(0.0)
	headBuyPriceLevel.Next = tailBuyPriceLevel
	tailBuyPriceLevel.Prev = headBuyPriceLevel

	headSellPriceLevel, tailSellPriceLevel := NewPriceLevel(0.0), NewPriceLevel(math.MaxFloat64)
	headSellPriceLevel.Next = tailSellPriceLevel
	tailSellPriceLevel.Prev = headSellPriceLevel

	return &OrderBook{
		buyHead:             headBuyPriceLevel,
		buyTail:             tailBuyPriceLevel,
		sellHead:            headSellPriceLevel,
		sellTail:            tailSellPriceLevel,
		priceLevelByOrderID: make(map[string]*PriceLevel),
		buyNodeByPrice:      make(map[float64]*PriceLevel),
		sellNodeByPrice:     make(map[float64]*PriceLevel),
	}
}

func (ob *OrderBook) GetPriceLevels(side Side) *PriceLevel {
	if side == SideBUY {
		return ob.buyHead.Next
	} else {
		return ob.sellHead.Next
	}
}

// Add
func (ob *OrderBook) AddOrder(order *Order) error {
	var priceLevel *PriceLevel
	if order.Side == SideBUY {
		priceLevel = ob.buyHead
	} else {
		priceLevel = ob.sellHead
	}

	for priceLevel != nil {
		if priceLevel.Price == order.Price {
			priceLevel.Add(order)
			ob.priceLevelByOrderID[order.ID] = priceLevel
			return nil
		}

		if order.Side == SideBUY {
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
	if order.Side == SideBUY {
		ob.buyNodeByPrice[order.Price] = newPriceLevel
	} else {
		ob.sellNodeByPrice[order.Price] = newPriceLevel
	}

	tmpNode := priceLevel.Prev
	tmpNode.Next = newPriceLevel
	newPriceLevel.Prev = tmpNode
	newPriceLevel.Next = priceLevel
	priceLevel.Prev = newPriceLevel
	return nil
}

// RemoveOrder
func (ob *OrderBook) RemoveOrder(orderID string) error {
	priceLevel, found := ob.priceLevelByOrderID[orderID]
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

	delete(ob.priceLevelByOrderID, orderID)
	return nil
}

// RemovePriceLevel
func (ob *OrderBook) RemovePriceLevel(side Side, price float64) error {
	var priceHash map[float64]*PriceLevel
	if side == SideBUY {
		priceHash = ob.buyNodeByPrice
	} else {
		priceHash = ob.sellNodeByPrice
	}

	priceLevel, found := priceHash[price]
	if !found {
		return ErrOrderNotExist
	}

	prevNode := priceLevel.Prev
	nextNode := priceLevel.Next
	prevNode.Next = nextNode
	nextNode.Prev = prevNode

	if side == SideBUY {
		delete(ob.buyNodeByPrice, price)
	} else {
		delete(ob.sellNodeByPrice, price)
	}

	return nil
}
