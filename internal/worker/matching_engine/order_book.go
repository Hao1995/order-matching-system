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

func (pl *OrderBook) GetPriceLevels(side Side) *PriceLevel {
	if side == SideBUY {
		return pl.buyHead.Next
	} else {
		return pl.sellHead.Next
	}
}

// Add
func (pl *OrderBook) AddOrder(order *Order) error {
	var priceLevel *PriceLevel
	if order.Side == SideBUY {
		priceLevel = pl.buyHead
	} else {
		priceLevel = pl.sellHead
	}

	for priceLevel != nil {
		if priceLevel.Price == order.Price {
			priceLevel.Add(order)
			pl.priceLevelByOrderID[order.ID] = priceLevel
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
		pl.buyNodeByPrice[order.Price] = newPriceLevel
	} else {
		pl.sellNodeByPrice[order.Price] = newPriceLevel
	}

	tmpNode := priceLevel.Prev
	tmpNode.Next = newPriceLevel
	newPriceLevel.Prev = tmpNode
	newPriceLevel.Next = priceLevel
	priceLevel.Prev = newPriceLevel
	return nil
}

// RemoveOrder
func (pl *OrderBook) RemoveOrder(orderID string) error {
	priceLevel, found := pl.priceLevelByOrderID[orderID]
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

	delete(pl.priceLevelByOrderID, orderID)
	return nil
}

// RemovePriceLevel
func (pl *OrderBook) RemovePriceLevel(side Side, price float64) error {
	var priceHash map[float64]*PriceLevel
	if side == SideBUY {
		priceHash = pl.buyNodeByPrice
	} else {
		priceHash = pl.sellNodeByPrice
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
		delete(pl.buyNodeByPrice, price)
	} else {
		delete(pl.sellNodeByPrice, price)
	}

	return nil
}
