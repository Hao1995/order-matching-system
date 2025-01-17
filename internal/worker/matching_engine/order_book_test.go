package matchingengine

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type OrderBookTestSuite struct {
	suite.Suite

	symbol string
}

func TestOrderBookTestSuite(t *testing.T) {
	suite.Run(t, new(OrderBookTestSuite))
}

func (s *OrderBookTestSuite) SetupSuite() {
	s.symbol = "AAPL"
}

func (s *OrderBookTestSuite) SetupTest() {}

func (s *OrderBookTestSuite) TeardownSuite() {}

func (s *OrderBookTestSuite) TeardownTest() {}

func (s *OrderBookTestSuite) TestGetPriceLevels() {
	// arrange
	orderBook := NewOrderBook()

	for _, t := range []struct {
		name     string
		side     Side
		expected *PriceLevel
	}{
		{
			name:     "test buy side",
			side:     SideBUY,
			expected: orderBook.buyHead.Next,
		},
		{
			name:     "test sell side",
			side:     SideSELL,
			expected: orderBook.sellHead.Next,
		},
	} {
		s.Run(t.name, func() {
			// act
			priceLevels := orderBook.GetPriceLevels(t.side)

			// assert
			s.Equal(t.expected, priceLevels)
		})
	}
}

func (s *OrderBookTestSuite) TestAddOrder_BuySide() {
	// arrange
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")
	side := SideBUY
	price1 := 25.89
	price2 := 30.33

	orderBook := NewOrderBook()
	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              side,
		Price:             price1,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	order2 := order1
	order2.Price = price2

	// act
	err := orderBook.AddOrder(&order1)
	s.ErrorIs(err, nil)

	err = orderBook.AddOrder(&order2)
	s.ErrorIs(err, nil)

	// assert
	buyPriceLevel := orderBook.buyHead.Next
	s.Equal(price2, buyPriceLevel.Price)
	s.Equal(&order2, buyPriceLevel.headOrder.Next)
	buyNode, ok := orderBook.buyNodeByPrice[price2]
	s.True(ok)
	s.Equal(buyPriceLevel, buyNode)

	buyPriceLevel = buyPriceLevel.Next
	s.Equal(price1, buyPriceLevel.Price)
	s.Equal(&order1, buyPriceLevel.headOrder.Next)
	buyNode, ok = orderBook.buyNodeByPrice[price1]
	s.True(ok)
	s.Equal(buyPriceLevel, buyNode)
}

func (s *OrderBookTestSuite) TestAddOrder_SellSide() {
	// arrange
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")
	side := SideSELL
	price1 := 30.33
	price2 := 25.89

	orderBook := NewOrderBook()
	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              side,
		Price:             price1,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	order2 := order1
	order2.Price = price2

	// act
	err := orderBook.AddOrder(&order1)
	s.ErrorIs(err, nil)

	err = orderBook.AddOrder(&order2)
	s.ErrorIs(err, nil)

	// assert
	sellPriceLevel := orderBook.sellHead.Next
	s.Equal(price2, sellPriceLevel.Price)
	s.Equal(&order2, sellPriceLevel.headOrder.Next)
	sellNode, ok := orderBook.sellNodeByPrice[price2]
	s.True(ok)
	s.Equal(sellPriceLevel, sellNode)

	sellPriceLevel = sellPriceLevel.Next
	s.Equal(price1, sellPriceLevel.Price)
	s.Equal(&order1, sellPriceLevel.headOrder.Next)
	sellNode, ok = orderBook.sellNodeByPrice[price1]
	s.True(ok)
	s.Equal(sellPriceLevel, sellNode)
}

func (s *OrderBookTestSuite) TestRemoveOrder() {
	// arrange
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")
	side := SideBUY
	price1 := 25.89
	price2 := 30.33

	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              side,
		Price:             price1,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	order2 := order1
	order2.Price = price2

	orderBook := NewOrderBook()
	priceLevel1 := NewPriceLevel(order1.Price, order1.Side)
	priceLevel1.Add(&order1)
	orderBook.priceLevelByOrderID[order1.ID] = priceLevel1

	priceLevel2 := NewPriceLevel(order2.Price, order2.Side)
	priceLevel2.Add(&order2)
	orderBook.priceLevelByOrderID[order2.ID] = priceLevel2

	orderBook.buyHead.Next = priceLevel2
	priceLevel2.Prev = orderBook.buyHead
	priceLevel2.Next = priceLevel1
	priceLevel1.Prev = priceLevel2
	priceLevel1.Next = orderBook.buyTail
	orderBook.buyTail.Prev = priceLevel1

	orderBook.buyNodeByPrice[priceLevel1.Price] = priceLevel1
	orderBook.buyNodeByPrice[priceLevel2.Price] = priceLevel2

	// act
	err := orderBook.RemoveOrder(order2.ID)
	s.ErrorIs(err, nil)

	// assert
	priceLevel := orderBook.buyHead.Next
	s.Equal(order1.Price, priceLevel.Price)
	s.Equal(orderBook.buyTail, priceLevel.Next)

	_, ok := orderBook.buyNodeByPrice[order2.Price]
	s.False(ok)
}

func (s *OrderBookTestSuite) TestRemovePriceLevel() {
	for _, t := range []struct {
		name string
		side Side
	}{
		{
			name: "buy side",
			side: SideBUY,
		},
		{
			name: "sell side",
			side: SideSELL,
		},
	} {
		s.Run(t.name, func() {
			// arrange
			now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")
			price := 25.89

			order1 := Order{
				ID:                uuid.NewString(),
				Symbol:            s.symbol,
				Side:              t.side,
				Price:             price,
				Quantity:          10,
				RemainingQuantity: 10,
				CanceledQuantity:  0,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			order2 := order1

			orderBook := NewOrderBook()
			priceLevel := NewPriceLevel(order1.Price, order1.Side)
			priceLevel.Add(&order1)
			orderBook.priceLevelByOrderID[order1.ID] = priceLevel
			priceLevel.Add(&order2)
			orderBook.priceLevelByOrderID[order2.ID] = priceLevel

			var head, tail *PriceLevel
			if t.side == SideBUY {
				head = orderBook.buyHead
				tail = orderBook.buyTail
				orderBook.buyNodeByPrice[priceLevel.Price] = priceLevel
			} else {
				head = orderBook.sellHead
				tail = orderBook.sellTail
				orderBook.sellNodeByPrice[priceLevel.Price] = priceLevel
			}

			head.Next = priceLevel
			priceLevel.Prev = head
			priceLevel.Next = tail
			tail.Prev = priceLevel

			// act
			err := orderBook.RemovePriceLevel(t.side, priceLevel.Price)
			s.ErrorIs(err, nil)

			// assert
			var ok bool
			if t.side == SideBUY {
				head = orderBook.buyHead
				tail = orderBook.buyTail
				_, ok = orderBook.buyNodeByPrice[priceLevel.Price]
			} else {
				head = orderBook.sellHead
				tail = orderBook.sellTail
				_, ok = orderBook.sellNodeByPrice[priceLevel.Price]
			}
			s.Equal(tail, head.Next)
			s.True(tail.IsDummyNode)
			s.False(ok)
		})
	}
}
