package matchingengine

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type MatcherTestSuite struct {
	suite.Suite
	matcher   *Matcher
	orderBook *OrderBook
}

func (suite *MatcherTestSuite) SetupTest() {
	suite.orderBook = NewOrderBook()
	suite.matcher = NewMatcher(suite.orderBook)
}

func (suite *MatcherTestSuite) TestCancelOrder() {
	order := Order{
		ID:                uuid.NewString(),
		Type:              OrderTypeBuy,
		Price:             100.0,
		Quantity:          10,
		RemainingQuantity: 10,
	}

	suite.orderBook.InsertOrder(order)
	err := suite.matcher.CancelOrder(order.ID)
	suite.NoError(err)
	suite.Empty(suite.orderBook.BuyLevels)
}

func (suite *MatcherTestSuite) TestMatchOrder_BuyOrder() {
	// Add a sell order
	sellOrder := Order{
		ID:                uuid.NewString(),
		Type:              OrderTypeSell,
		Price:             100.0,
		Quantity:          10,
		RemainingQuantity: 10,
	}
	suite.orderBook.InsertOrder(sellOrder)

	// Add a buy order
	buyOrder := Order{
		ID:                uuid.NewString(),
		Type:              OrderTypeBuy,
		Price:             100.0,
		Quantity:          5,
		RemainingQuantity: 5,
	}

	transactions := suite.matcher.MatchOrder(buyOrder)

	suite.Equal(1, len(transactions))
	transaction := transactions[0]
	suite.Equal(buyOrder.ID, transaction.BuyOrderID)
	suite.Equal(sellOrder.ID, transaction.SellOrderID)
	suite.Equal(float64(100.0), transaction.Price)
	suite.Equal(float64(5), transaction.Quantity)
	// Verify remaining sell order quantity
	suite.Equal(float64(5), suite.orderBook.SellLevels.HeadOrders.Order.RemainingQuantity)
}

func (suite *MatcherTestSuite) TestMatchOrder_SellOrder() {
	// Add a buy order
	buyOrder := Order{
		ID:       uuid.NewString(),
		Type:     OrderTypeBuy,
		Price:    100.0,
		Quantity: 10,
	}
	suite.orderBook.InsertOrder(buyOrder)

	// Add a sell order
	sellOrder := Order{
		ID:                uuid.NewString(),
		Type:              OrderTypeSell,
		Price:             100.0,
		Quantity:          5,
		RemainingQuantity: 5,
	}

	transactions := suite.matcher.MatchOrder(sellOrder)

	suite.Equal(1, len(transactions))
	transaction := transactions[0]
	suite.Equal(buyOrder.ID, transaction.BuyOrderID)
	suite.Equal(sellOrder.ID, transaction.SellOrderID)
	suite.Equal(float64(100.0), transaction.Price)
	suite.Equal(float64(5), transaction.Quantity)
	// Verify remaining buy order quantity
	suite.Equal(float64(5), suite.orderBook.BuyLevels.HeadOrders.Order.RemainingQuantity)
}

func (suite *MatcherTestSuite) TestMatchOrder_PartialFill() {
	// Add a sell order with larger quantity
	sellOrder := Order{
		ID:                uuid.NewString(),
		Type:              OrderTypeSell,
		Price:             100.0,
		Quantity:          10,
		RemainingQuantity: 10,
	}
	suite.orderBook.InsertOrder(sellOrder)

	// Add a buy order with smaller quantity
	buyOrder := Order{
		ID:                uuid.NewString(),
		Type:              OrderTypeBuy,
		Price:             100.0,
		Quantity:          5,
		RemainingQuantity: 5,
	}

	transactions := suite.matcher.MatchOrder(buyOrder)

	suite.Equal(1, len(transactions))
	transaction := transactions[0]
	suite.Equal(buyOrder.ID, transaction.BuyOrderID)
	suite.Equal(sellOrder.ID, transaction.SellOrderID)
	suite.Equal(float64(100.0), transaction.Price)
	suite.Equal(float64(5), transaction.Quantity)
	// Verify remaining sell order quantity
	suite.Equal(float64(5), suite.orderBook.SellLevels.HeadOrders.Order.RemainingQuantity)
}

func (suite *MatcherTestSuite) TestMatchOrder_NoMatch() {
	// Add a sell order
	sellOrder := Order{
		ID:                uuid.NewString(),
		Type:              OrderTypeSell,
		Price:             105.0,
		Quantity:          10,
		RemainingQuantity: 10,
	}
	suite.orderBook.InsertOrder(sellOrder)

	// Add a buy order with lower price
	buyOrder := Order{
		ID:       uuid.NewString(),
		Type:     OrderTypeBuy,
		Price:    100.0,
		Quantity: 5,
	}

	transactions := suite.matcher.MatchOrder(buyOrder)

	suite.Equal(0, len(transactions))
	// Verify orders are still in their respective levels
	suite.NotNil(suite.orderBook.SellLevels)
	suite.NotNil(suite.orderBook.BuyLevels)
}

func TestMatcherTestSuite(t *testing.T) {
	suite.Run(t, new(MatcherTestSuite))
}
