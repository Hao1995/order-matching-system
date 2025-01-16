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

func (s *OrderBookTestSuite) TestAdd_BuySide() {
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
	err := orderBook.Add(&order1)
	s.ErrorIs(err, nil)

	err = orderBook.Add(&order2)
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

func (s *OrderBookTestSuite) TestAdd_SellSide() {
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
	err := orderBook.Add(&order1)
	s.ErrorIs(err, nil)

	err = orderBook.Add(&order2)
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
