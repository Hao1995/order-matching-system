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
	tickNum   int8

	symbol  string
	now     time.Time
	uuidStr string
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
	suite.tickNum = 5
	suite.matcher = NewMatcher(suite.orderBook, suite.tickNum)
}

func (suite *MatcherTestSuite) TestCancelOrder() {
	order := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  10,
		CreatedAt: time.Now(),
	}

	suite.orderBook.InsertOrder(order)
	matching, err := suite.matcher.CancelOrder(order.ID)
	suite.NoError(err)
	suite.Empty(suite.orderBook.BuyLevels)
	suite.NotNil(matching.BuyTicks)
	suite.NotNil(matching.SellTicks)
	suite.Len(matching.BuyTicks, 0)
	suite.Len(matching.SellTicks, 0)
}

func (suite *MatcherTestSuite) TestCreateOrder_NoMatch() {
	order := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeSell,
		Price:     200.0,
		Quantity:  15,
		CreatedAt: time.Now(),
	}

	matching := suite.matcher.CreateOrder(order)
	suite.Empty(matching.Transactions)
	suite.Equal(1, len(matching.SellTicks))
	suite.Equal(Tick{
		Price:    200.0,
		Quantity: 15,
	}, matching.SellTicks[0])
}

func (suite *MatcherTestSuite) TestCreateOrder_PartialFill() {
	sellOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeSell,
		Price:     100.0,
		Quantity:  10,
		CreatedAt: time.Now(),
	}
	suite.orderBook.InsertOrder(sellOrder)

	buyOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  5,
		CreatedAt: time.Now(),
	}

	matching := suite.matcher.CreateOrder(buyOrder)
	suite.Equal(1, len(matching.Transactions))
	suite.Equal(Transaction{
		ID:          getUUID(),
		Symbol:      suite.symbol,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		Price:       100.0,
		Quantity:    5,
		CreatedAt:   now(),
	}, matching.Transactions[0])
}

func (suite *MatcherTestSuite) TestCreateOrder_MultipleMatches() {
	sellOrder1 := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeSell,
		Price:     100.0,
		Quantity:  5,
		CreatedAt: time.Now(),
	}
	sellOrder2 := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeSell,
		Price:     100.0,
		Quantity:  5,
		CreatedAt: time.Now(),
	}
	suite.orderBook.InsertOrder(sellOrder1)
	suite.orderBook.InsertOrder(sellOrder2)

	buyOrder := Order{
		ID:        uuid.NewString(),
		Symbol:    suite.symbol,
		Type:      OrderTypeBuy,
		Price:     100.0,
		Quantity:  10,
		CreatedAt: time.Now(),
	}

	matching := suite.matcher.CreateOrder(buyOrder)
	suite.Equal(2, len(matching.Transactions))
	suite.Equal(Transaction{
		ID:          getUUID(),
		Symbol:      suite.symbol,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder1.ID,
		Price:       100.0,
		Quantity:    5,
		CreatedAt:   now(),
	}, matching.Transactions[0])
	suite.Equal(Transaction{
		ID:          getUUID(),
		Symbol:      suite.symbol,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder2.ID,
		Price:       100.0,
		Quantity:    5,
		CreatedAt:   now(),
	}, matching.Transactions[1])
}

func (suite *MatcherTestSuite) TestCreateOrder_ExceedTickLimit() {
	for i := 0; i < 10; i++ {
		order := Order{
			ID:        uuid.NewString(),
			Symbol:    suite.symbol,
			Type:      OrderTypeSell,
			Price:     float64(100 + i),
			Quantity:  int64(10 * (i + 1)),
			CreatedAt: time.Now(),
		}
		suite.orderBook.InsertOrder(order)
	}

	buyOrder := Order{
		ID:        uuid.NewString(),
		Type:      OrderTypeBuy,
		Price:     110.0,
		Quantity:  50,
		CreatedAt: time.Now(),
	}
	matching := suite.matcher.CreateOrder(buyOrder)

	suite.Len(matching.SellTicks, 5)
	suite.Equal([]Tick{
		{
			Price:    102.00,
			Quantity: 30,
		},
		{
			Price:    103.00,
			Quantity: 40,
		},
		{
			Price:    104.00,
			Quantity: 50,
		},
		{
			Price:    105.00,
			Quantity: 60,
		},
		{
			Price:    106.00,
			Quantity: 70,
		},
	}, matching.SellTicks)
}
