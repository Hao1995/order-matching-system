package matchingengine

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	SYMBOL = "APPL"
)

var (
	now time.Time
)

func init() {
	now, _ = time.Parse(time.RFC3339, "2024-01-01T00:00:00+08:00")
}

// TestGetFirstOrder
func TestGetFirstOrder(t *testing.T) {
	price := 35.9
	priceLevel := NewPriceLevel(price)

	order1 := Order{
		ID:                uuid.NewString(),
		Symbol:            SYMBOL,
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

	priceLevel.Add(&order1)
	priceLevel.Add(&order2)

	firstOrder := priceLevel.GetFirstOrder()

	assert.Equal(t, &order1, firstOrder)
}
