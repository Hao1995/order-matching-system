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

func (suite *PriceLevelTestSuite) TestGetFirstOrder_WithAscendingCreatedAtOrder() {
	price := 35.9
	priceLevel := NewPriceLevel(price)
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
	order2.CreatedAt = now.Add(1 * time.Hour)
	order2.UpdatedAt = now.Add(1 * time.Hour)

	// add orders with ascending order
	priceLevel.Add(&order1)
	priceLevel.Add(&order2)

	firstOrder := priceLevel.GetFirstOrder()

	suite.Equal(&order1, firstOrder)
}

func (suite *PriceLevelTestSuite) TestGetFirstOrder_WithDescendingCreatedAtOrder() {
	price := 35.9
	priceLevel := NewPriceLevel(price)
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
	order2.CreatedAt = now.Add(1 * time.Hour)
	order2.UpdatedAt = now.Add(1 * time.Hour)

	// add orders with ascending order
	priceLevel.Add(&order2)
	priceLevel.Add(&order1)

	firstOrder := priceLevel.GetFirstOrder()

	suite.Equal(&order1, firstOrder)
}
