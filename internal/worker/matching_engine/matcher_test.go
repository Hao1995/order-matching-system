package matchingengine

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MatcherTestSuite struct {
	suite.Suite
	matcher   *Matcher
	orderBook *OrderBook
	symbol    string
	now       time.Time
	uuidStr   string
}

func TestMatcherTestSuite(t *testing.T) {
	suite.Run(t, new(MatcherTestSuite))
}

func (suite *MatcherTestSuite) SetupSuite() {
	suite.symbol = "AAPL"
	suite.now = time.Date(2025, 1, 15, 18, 0, 0, 0, time.UTC)
	suite.uuidStr = uuid.NewString()

	now = func() time.Time {
		return suite.now
	}

	getUUID = func() string {
		return suite.uuidStr
	}
}
func (suite *MatcherTestSuite) SetupTest() {
	suite.orderBook = NewOrderBook()
	suite.matcher = NewMatcher(suite.orderBook)
}

func (suite *MatcherTestSuite) TestCancelOrder() {
	order := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  10,
		CreatedAt: suite.now,
	}

	suite.orderBook.InsertOrder(order)
	err := suite.matcher.CancelOrder(order.ID)
	suite.NoError(err)
	suite.Empty(suite.orderBook.BuyLevels)
}

func (suite *MatcherTestSuite) TestMatchOrder_BuyOrder() {
	// Add a sell order
	sellOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeSell,
		Price:     100.0,
		Quantity:  10,
		CreatedAt: suite.now,
	}
	suite.orderBook.InsertOrder(sellOrder)

	// Add a buy order
	buyOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  5,
		CreatedAt: suite.now,
	}

	transactions := suite.matcher.MatchOrder(buyOrder)

	suite.Equal(1, len(transactions))
	suite.Equal(Transaction{
		ID:          getUUID(),
		Symbol:      suite.symbol,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		Price:       100.0,
		Quantity:    5,
		CreatedAt:   now(),
	}, transactions[0])
	suite.Equal(int64(5), suite.orderBook.SellLevels.HeadOrders.Order.Quantity)
}

func (suite *MatcherTestSuite) TestMatchOrder_SellOrder() {
	// Add a buy order
	buyOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  10,
		CreatedAt: suite.now,
	}
	suite.orderBook.InsertOrder(buyOrder)

	// Add a sell order
	sellOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeSell,
		Price:     100.0,
		Quantity:  5,
		CreatedAt: suite.now,
	}

	transactions := suite.matcher.MatchOrder(sellOrder)

	suite.Equal(1, len(transactions))
	suite.Equal(Transaction{
		ID:          getUUID(),
		Symbol:      suite.symbol,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		Price:       100.0,
		Quantity:    5,
		CreatedAt:   now(),
	}, transactions[0])
	suite.Equal(int64(5), suite.orderBook.BuyLevels.HeadOrders.Order.Quantity)
}

func (suite *MatcherTestSuite) TestMatchOrder_PartialFill() {
	// Add a sell order with larger quantity
	sellOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeSell,
		Price:     100.0,
		Quantity:  5,
		CreatedAt: suite.now,
	}
	suite.orderBook.InsertOrder(sellOrder)

	// Add a buy order with smaller quantity
	buyOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  10,
		CreatedAt: suite.now,
	}

	transactions := suite.matcher.MatchOrder(buyOrder)

	suite.Equal(1, len(transactions))
	suite.Equal(Transaction{
		ID:          getUUID(),
		Symbol:      suite.symbol,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		Price:       100.0,
		Quantity:    5,
		CreatedAt:   now(),
	}, transactions[0])
	suite.Nil(suite.orderBook.SellLevels)
	suite.Equal(float64(100.0), suite.orderBook.BuyLevels.HeadOrders.Order.Price)
	suite.Equal(int64(5), suite.orderBook.BuyLevels.HeadOrders.Order.Quantity)
}

func (suite *MatcherTestSuite) TestMatchOrder_NoMatch() {
	// Add a sell order
	sellOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeSell,
		Price:     105.0,
		Quantity:  10,
		CreatedAt: suite.now,
	}
	suite.orderBook.InsertOrder(sellOrder)

	// Add a buy order with lower price
	buyOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  5,
		CreatedAt: suite.now,
	}

	transactions := suite.matcher.MatchOrder(buyOrder)

	suite.Equal(0, len(transactions))
	suite.Equal(sellOrder, suite.orderBook.SellLevels.HeadOrders.Order)
	suite.Equal(buyOrder, suite.orderBook.BuyLevels.HeadOrders.Order)
}
