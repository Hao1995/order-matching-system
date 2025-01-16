package matchingengine

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MatchingEngineTestSuite struct {
	suite.Suite

	symbol string
}

func TestMatchingEngineTestSuite(t *testing.T) {
	suite.Run(t, new(MatchingEngineTestSuite))
}

func (s *MatchingEngineTestSuite) SetupSuite() {
	s.symbol = "AAPL"
}

func (s *MatchingEngineTestSuite) SetupTest() {}

func (s *MatchingEngineTestSuite) TearDownTest() {}

func (s *MatchingEngineTestSuite) TearDownSuite() {}

func (s *MatchingEngineTestSuite) TestCancelOrder() {
	// arrange
	// insert order
	side := SideBUY
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")
	price := 25.89

	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              side,
		Price:             price,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	order2 := order1
	order2.ID = uuid.NewString()
	order2.CreatedAt = now.Add(1 * time.Minute)
	order2.UpdatedAt = now.Add(1 * time.Minute)

	order3 := order1
	order3.ID = uuid.NewString()
	order3.CreatedAt = now.Add(2 * time.Minute)
	order3.UpdatedAt = now.Add(2 * time.Minute)

	orderBook := NewOrderBook()
	orderBook.AddOrder(&order1)
	orderBook.AddOrder(&order2)
	orderBook.AddOrder(&order3)

	// act
	matchingEngine := NewMatchingEngine(orderBook)
	matchingEngine.CancelOrder(order2.ID)

	// assert
	priceLevel := orderBook.buyHead.Next
	s.Equal(price, priceLevel.Price)
	order := priceLevel.headOrder.Next
	s.Equal(&order1, order)
	s.Equal(&order3, order.Next)
	s.Equal(priceLevel.tailOrder, order.Next.Next)
}
