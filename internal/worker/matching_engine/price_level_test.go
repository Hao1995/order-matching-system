package matchingengine

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type PriceLevelTestSuite struct {
	suite.Suite

	symbol string
}

func TestPriceLevelTestSuite(t *testing.T) {
	suite.Run(t, new(PriceLevelTestSuite))
}

func (suite *PriceLevelTestSuite) SetupSuite() {
	suite.symbol = "AAPL"
}

func (suite *PriceLevelTestSuite) SetupTest() {}

func (suite *PriceLevelTestSuite) TearDownTest() {}

func (suite *PriceLevelTestSuite) TearDownSuite() {}

func (suite *PriceLevelTestSuite) TestGetFirstOrder() {
	// arrange
	side := SideBUY
	price := 35.9
	priceLevel := NewPriceLevel(price, side)
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00+08:00")

	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            suite.symbol,
		Side:              SideBUY,
		Price:             price,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	order2 := order1
	order2.ID = uuid.NewString()
	order2.CreatedAt = now.Add(1 * time.Hour)
	order2.UpdatedAt = now.Add(1 * time.Hour)

	// add orders with ascending order
	priceLevel.headOrder.Next = &order1
	order1.Prev = priceLevel.headOrder
	priceLevel.orderByID[order1.ID] = &order1

	order1.Next = &order2
	order2.Prev = &order1
	priceLevel.orderByID[order2.ID] = &order2

	order2.Next = priceLevel.tailOrder
	priceLevel.tailOrder.Prev = &order2

	// act
	firstOrder := priceLevel.GetFirstOrder()

	// assert
	suite.Equal(&order1, firstOrder)
}

func (suite *PriceLevelTestSuite) TestAdd() {
	// arrange
	side := SideBUY
	price := 35.9
	priceLevel := NewPriceLevel(price, side)
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00+08:00")

	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            suite.symbol,
		Side:              SideBUY,
		Price:             price,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	order2 := order1
	order2.ID = uuid.NewString()
	order2.CreatedAt = now.Add(1 * time.Hour)
	order2.UpdatedAt = now.Add(1 * time.Hour)

	// act
	priceLevel.Add(&order2)
	priceLevel.Add(&order1)

	// assert
	suite.Equal(&order1, priceLevel.headOrder.Next)
	suite.Equal(&order2, priceLevel.headOrder.Next.Next)
}

func (suite *PriceLevelTestSuite) TestAdd_ErrNotMatchingPrice() {
	// arrange
	side := SideBUY
	price1 := 35.9
	price2 := 30.1
	priceLevel := NewPriceLevel(price1, side)
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00+08:00")

	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            suite.symbol,
		Side:              side,
		Price:             price2,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// act
	err := priceLevel.Add(&order1)

	// assert
	suite.ErrorIs(err, ErrNotMatchingPrice)
}

func (suite *PriceLevelTestSuite) TestRemove() {
	// arrange
	side := SideBUY
	price := 35.9
	priceLevel := NewPriceLevel(price, side)
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00+08:00")

	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            suite.symbol,
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
	order2.CreatedAt = now.Add(1 * time.Hour)
	order2.UpdatedAt = now.Add(1 * time.Hour)

	// add orders
	priceLevel.headOrder.Next = &order1
	order1.Prev = priceLevel.headOrder
	priceLevel.orderByID[order1.ID] = &order1

	order1.Next = &order2
	order2.Prev = &order1
	priceLevel.orderByID[order2.ID] = &order2

	order2.Next = priceLevel.tailOrder
	priceLevel.tailOrder.Prev = &order2

	// act
	err := priceLevel.Remove(order1.ID)

	// assert
	suite.ErrorIs(err, nil)
	suite.Equal(&order2, priceLevel.headOrder.Next)
	suite.True(order2.Next.IsDummyNode)
}

func (suite *PriceLevelTestSuite) TestRemove_OrderNotExist() {
	// arrange
	side := SideBUY
	price := 35.9
	priceLevel := NewPriceLevel(price, side)
	now, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00+08:00")

	// add orders
	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            suite.symbol,
		Side:              side,
		Price:             price,
		Quantity:          10,
		RemainingQuantity: 10,
		CanceledQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	priceLevel.headOrder.Next = &order1
	order1.Prev = priceLevel.headOrder
	priceLevel.orderByID[order1.ID] = &order1

	order1.Next = priceLevel.tailOrder
	priceLevel.tailOrder.Prev = &order1

	// act
	err := priceLevel.Remove(uuid.NewString())

	// assert
	suite.ErrorIs(err, ErrOrderNotExist)
}
