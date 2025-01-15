package matchingengine

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type OrderBookTestSuite struct {
	suite.Suite
}

func TestOrderBookTestSuite(t *testing.T) {
	suite.Run(t, new(OrderBookTestSuite))
}

func (s *OrderBookTestSuite) SetupSuite()    {}
func (s *OrderBookTestSuite) SetupTest()     {}
func (s *OrderBookTestSuite) TeardownSuite() {}
func (s *OrderBookTestSuite) TeardownTest()  {}

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
			s.NotEqual(t.expected, priceLevels)
		})
	}
}
