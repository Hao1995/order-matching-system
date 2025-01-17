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

func (s *MatchingEngineTestSuite) TestPlaceOrder_BuyOrder_MatchFromLowerPriceToHigherPrice() {
	// arrange
	nowTime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")
	price := 25.89

	// add sell orders
	sellOrder1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              SideSELL,
		Price:             price - 1,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         nowTime,
		UpdatedAt:         nowTime,
	}

	sellOrder2 := sellOrder1
	sellOrder2.ID = uuid.NewString()
	sellOrder2.Price = price - 2
	sellOrder2.CreatedAt = nowTime.Add(1 * time.Minute)
	sellOrder2.UpdatedAt = nowTime.Add(1 * time.Minute)

	sellOrder3 := sellOrder1
	sellOrder3.ID = uuid.NewString()
	sellOrder3.Price = price - 1
	sellOrder3.CreatedAt = nowTime.Add(1 * time.Minute)
	sellOrder3.UpdatedAt = nowTime.Add(1 * time.Minute)

	sellOrder4 := sellOrder1
	sellOrder4.ID = uuid.NewString()
	sellOrder4.Price = price
	sellOrder4.CreatedAt = nowTime.Add(1 * time.Minute)
	sellOrder4.UpdatedAt = nowTime.Add(1 * time.Minute)

	orderBook := NewOrderBook()
	orderBook.AddOrder(&sellOrder1)
	orderBook.AddOrder(&sellOrder2)
	orderBook.AddOrder(&sellOrder3)
	orderBook.AddOrder(&sellOrder4)

	// init buy order
	buyOrder1 := sellOrder1
	buyOrder1.ID = uuid.NewString()
	buyOrder1.Price = price
	buyOrder1.Quantity = 30
	buyOrder1.RemainingQuantity = 30
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
	// assert only left non-matched order
	s.Equal(&sellOrder4, orderBook.sellHead.Next.headOrder.Next)
	// assert matching events
	s.Equal(3, len(matchingEvents))
	s.Equal(events.MatchingTransaction{
		ID:          uuidStr,
		Symbol:      s.symbol,
		BuyOrderID:  buyOrder1.ID,
		SellOrderID: sellOrder2.ID,
		Price:       sellOrder2.Price,
		Quantity:    sellOrder2.Quantity,
		CreatedAt:   nowTime,
	}, matchingEvents[0])
	s.Equal(events.MatchingTransaction{
		ID:          uuidStr,
		Symbol:      s.symbol,
		BuyOrderID:  buyOrder1.ID,
		SellOrderID: sellOrder1.ID,
		Price:       sellOrder1.Price,
		Quantity:    sellOrder1.Quantity,
		CreatedAt:   nowTime,
	}, matchingEvents[1])
	s.Equal(events.MatchingTransaction{
		ID:          uuidStr,
		Symbol:      s.symbol,
		BuyOrderID:  buyOrder1.ID,
		SellOrderID: sellOrder3.ID,
		Price:       sellOrder3.Price,
		Quantity:    sellOrder3.Quantity,
		CreatedAt:   nowTime,
	}, matchingEvents[2])
}

func (s *MatchingEngineTestSuite) TestPlaceOrder_BuyOrder_InsertRemainingOrderToOrderBook() {
	// arrange
	nowTime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")
	price := 25.89
	var sellQuantity, buyQuantity int64 = 10, 30

	// add sell orders
	sellOrder1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              SideSELL,
		Price:             price,
		Quantity:          sellQuantity,
		RemainingQuantity: sellQuantity,
		CanceledQuantity:  0,
		CreatedAt:         nowTime,
		UpdatedAt:         nowTime,
	}

	orderBook := NewOrderBook()
	orderBook.AddOrder(&sellOrder1)

	// init buy order
	buyOrder1 := sellOrder1
	buyOrder1.ID = uuid.NewString()
	buyOrder1.Price = price
	buyOrder1.Quantity = buyQuantity
	buyOrder1.RemainingQuantity = buyQuantity
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
	s.Equal(1, len(matchingEvents))
	s.Equal(events.MatchingTransaction{
		ID:          uuidStr,
		Symbol:      s.symbol,
		BuyOrderID:  buyOrder1.ID,
		SellOrderID: sellOrder1.ID,
		Price:       sellOrder1.Price,
		Quantity:    sellOrder1.Quantity,
		CreatedAt:   nowTime,
	}, matchingEvents[0])
	// assert the remaining buy order is in the order book
	s.Equal(&buyOrder1, orderBook.buyHead.Next.headOrder.Next)
	s.Equal(int64(20), buyOrder1.RemainingQuantity)
}

func (s *MatchingEngineTestSuite) TestPlaceOrder_SellOrder_MatchFromHigherPriceToLowerPrice() {
	// arrange
	nowTime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")

	// add sell orders
	buyOrder1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              SideBUY,
		Price:             20,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         nowTime,
		UpdatedAt:         nowTime,
	}

	buyOrder2 := buyOrder1
	buyOrder2.ID = uuid.NewString()
	buyOrder2.Price = 30
	buyOrder2.CreatedAt = nowTime.Add(1 * time.Minute)
	buyOrder2.UpdatedAt = nowTime.Add(1 * time.Minute)

	buyOrder3 := buyOrder1
	buyOrder3.ID = uuid.NewString()
	buyOrder3.Price = 20
	buyOrder3.CreatedAt = nowTime.Add(1 * time.Minute)
	buyOrder3.UpdatedAt = nowTime.Add(1 * time.Minute)

	orderBook := NewOrderBook()
	orderBook.AddOrder(&buyOrder1)
	orderBook.AddOrder(&buyOrder2)
	orderBook.AddOrder(&buyOrder3)

	// init sell order
	sellOrder1 := buyOrder1
	sellOrder1.ID = uuid.NewString()
	sellOrder1.Side = SideSELL
	sellOrder1.Price = 20
	sellOrder1.Quantity = 30
	sellOrder1.RemainingQuantity = 30

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
	matchingEvents := matchingEngine.PlaceOrder(&sellOrder1)

	// assert
	// assert matched orders not exist
	s.Equal(orderBook.buyTail, orderBook.buyHead.Next)
	// assert matching events
	s.Equal(3, len(matchingEvents))
	s.Equal(events.MatchingTransaction{
		ID:          uuidStr,
		Symbol:      s.symbol,
		BuyOrderID:  buyOrder2.ID,
		SellOrderID: sellOrder1.ID,
		Price:       sellOrder1.Price,
		Quantity:    buyOrder2.Quantity,
		CreatedAt:   nowTime,
	}, matchingEvents[0])
	s.Equal(events.MatchingTransaction{
		ID:          uuidStr,
		Symbol:      s.symbol,
		BuyOrderID:  buyOrder1.ID,
		SellOrderID: sellOrder1.ID,
		Price:       sellOrder1.Price,
		Quantity:    buyOrder1.Quantity,
		CreatedAt:   nowTime,
	}, matchingEvents[1])
	s.Equal(events.MatchingTransaction{
		ID:          uuidStr,
		Symbol:      s.symbol,
		BuyOrderID:  buyOrder3.ID,
		SellOrderID: sellOrder1.ID,
		Price:       sellOrder1.Price,
		Quantity:    buyOrder3.Quantity,
		CreatedAt:   nowTime,
	}, matchingEvents[2])
	// assert sell order is fully matched
	s.Equal(int64(0), sellOrder1.RemainingQuantity)
}

func (s *MatchingEngineTestSuite) TestPlaceOrder_SellOrder_InsertOrderWhenNoMatchedOrders() {
	// arrange
	nowTime, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z08:00")

	// add sell orders
	buyOrder1 := Order{
		ID:                uuid.NewString(),
		Symbol:            s.symbol,
		Side:              SideBUY,
		Price:             30,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         nowTime,
		UpdatedAt:         nowTime,
	}

	orderBook := NewOrderBook()
	orderBook.AddOrder(&buyOrder1)

	// init sell order
	sellOrder1 := buyOrder1
	sellOrder1.ID = uuid.NewString()
	sellOrder1.Side = SideSELL
	sellOrder1.Price = 50

	// act
	matchingEngine := NewMatchingEngine(orderBook)
	matchingEvents := matchingEngine.PlaceOrder(&sellOrder1)

	// assert
	// assert buy orders still exist
	s.Equal(&buyOrder1, orderBook.buyHead.Next.headOrder.Next)
	// assert incoming order exist in order book
	s.Equal(&sellOrder1, orderBook.sellHead.Next.headOrder.Next)
	// assert no matched events
	s.Equal(0, len(matchingEvents))
}
