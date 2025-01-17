package matchingengine

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/Hao1995/order-matching-system/internal/common/models/events"
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

func (s *MatchingEngineTestSuite) TestPlaceOrder_BuyOrder_PriceAndQuantityAreFullyMatched() {
	// arrange
	nowTime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")
	price := 25.89

	// add sell orders
	sellOrder1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              SideSELL,
		Price:             price,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         nowTime,
		UpdatedAt:         nowTime,
	}

	orderBook := NewOrderBook()
	orderBook.AddOrder(&sellOrder1)

	// init buy order
	buyOrder1 := sellOrder1
	buyOrder1.ID = uuid.NewString()
	buyOrder1.Side = SideBUY

	// mock uuid
	uuidStr := uuid.NewString()
	getUUID = func() string {
		return uuidStr
	}

	// mock now
	now = func() time.Time {
		return nowTime
	}

	// act
	matchingEngine := NewMatchingEngine(orderBook)
	matchingEvents := matchingEngine.PlaceOrder(&buyOrder1)

	// assert
	// assert matched orders not exist
	s.Equal(orderBook.sellHead.Next, orderBook.sellTail)
	// assert matching events
	s.Equal(events.MatchingTransaction{
		ID:          uuidStr,
		Symbol:      s.symbol,
		BuyOrderID:  buyOrder1.ID,
		SellOrderID: sellOrder1.ID,
		Price:       price,
		Quantity:    buyOrder1.Quantity,
		CreatedAt:   nowTime,
	}, matchingEvents[0])
}
