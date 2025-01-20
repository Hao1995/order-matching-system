package matchingengine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type OrderBookTestSuite struct {
	suite.Suite
	orderBook *OrderBook
	symbol    string
	now       time.Time
}

func TestOrderBookTestSuite(t *testing.T) {
	suite.Run(t, new(OrderBookTestSuite))
}

func (suite *OrderBookTestSuite) SetupTest() {
	suite.orderBook = NewOrderBook()
	suite.symbol = "AAPL"
	suite.now = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
}

func (suite *OrderBookTestSuite) TestInsertOrder() {
	// Test inserting orders with different price levels
	orders := []Order{
		{
			ID:        "order1",
			Symbol:    suite.symbol,
			Type:      OrderTypeBuy,
			Price:     101.0,
			Quantity:  10,
			CreatedAt: suite.now,
		}, {
			ID:        "order2",
			Symbol:    suite.symbol,
			Type:      OrderTypeBuy,
			Price:     100.0,
			Quantity:  20,
			CreatedAt: suite.now,
		}, {
			ID:        "order3",
			Symbol:    suite.symbol,
			Type:      OrderTypeSell,
			Price:     102.0,
			Quantity:  25,
			CreatedAt: suite.now,
		}, {
			ID:        "order4",
			Symbol:    suite.symbol,
			Type:      OrderTypeSell,
			Price:     103.0,
			Quantity:  30,
			CreatedAt: suite.now,
		},
	}

	for _, order := range orders {
		suite.orderBook.InsertOrder(order)
	}

	// Verify that orders are inserted at the correct price levels
	suite.Equal(101.0, suite.orderBook.BuyLevels.Price)
	suite.Equal(102.0, suite.orderBook.SellLevels.Price)

	// Test inserting orders with the same price level
	duplicateOrder := Order{
		ID:        "order5",
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     101.0,
		Quantity:  7,
		CreatedAt: suite.now,
	}
	suite.orderBook.InsertOrder(duplicateOrder)
	suite.Equal(int64(17), suite.orderBook.BuyLevels.TotalQuantity)
}

func (suite *OrderBookTestSuite) TestDeleteOrder() {
	// Test deleting orders that exist
	order := Order{
		ID:        "order1",
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  15,
		CreatedAt: suite.now,
	}

	suite.orderBook.InsertOrder(order)
	err := suite.orderBook.DeleteOrder(order.ID)
	suite.NoError(err)

	// Verify the order is deleted correctly
	_, exists := suite.orderBook.orderMap[order.ID]
	suite.False(exists)

	// Test deleting an order that does not exist
	err = suite.orderBook.DeleteOrder("nonexistent_order")
	suite.Error(err)

	// Test deleting the only order in a price level
	order2 := Order{
		ID:        "order2",
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     101.0,
		Quantity:  10,
		CreatedAt: suite.now,
	}
	suite.orderBook.InsertOrder(order2)
	suite.Equal(101.0, suite.orderBook.BuyLevels.Price)
	err = suite.orderBook.DeleteOrder(order2.ID)
	suite.NoError(err)
	suite.Nil(suite.orderBook.BuyLevels)
}

func (suite *OrderBookTestSuite) TestGetTopTicks() {
	orders := []Order{
		{
			ID:        "order1",
			Symbol:    suite.symbol,
			Type:      OrderTypeBuy,
			Price:     101.0,
			Quantity:  10,
			CreatedAt: suite.now,
		},
		{
			ID:        "order2",
			Symbol:    suite.symbol,
			Type:      OrderTypeBuy,
			Price:     100.0,
			Quantity:  20,
			CreatedAt: suite.now,
		},
		{
			ID:        "order3",
			Symbol:    suite.symbol,
			Type:      OrderTypeSell,
			Price:     102.0,
			Quantity:  25,
			CreatedAt: suite.now,
		},
		{
			ID:        "order4",
			Symbol:    suite.symbol,
			Type:      OrderTypeSell,
			Price:     103.0,
			Quantity:  30,
			CreatedAt: suite.now,
		},
	}

	for _, order := range orders {
		suite.orderBook.InsertOrder(order)
	}

	// Test retrieving top N ticks with N less than total levels
	buyTicks, sellTicks := suite.orderBook.GetTopTicks(2)
	suite.Len(buyTicks, 2)
	suite.Len(sellTicks, 2)

	// Test retrieving top N ticks with N greater than total levels
	buyTicks, sellTicks = suite.orderBook.GetTopTicks(5)
	suite.Len(buyTicks, 2)
	suite.Len(sellTicks, 2)

	// Test retrieving top N ticks with exactly one level
	suite.orderBook = NewOrderBook()
	order := Order{
		ID:        "order5",
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     105.0,
		Quantity:  5,
		CreatedAt: suite.now,
	}
	suite.orderBook.InsertOrder(order)
	buyTicks, sellTicks = suite.orderBook.GetTopTicks(1)
	suite.Len(buyTicks, 1)
	suite.Len(sellTicks, 0)

	// Test retrieving top N ticks with empty order book
	suite.orderBook = NewOrderBook()
	buyTicks, sellTicks = suite.orderBook.GetTopTicks(2)
	suite.Len(buyTicks, 0)
	suite.Len(sellTicks, 0)
}
