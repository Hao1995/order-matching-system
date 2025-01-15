//go:generate go-enum --marshal
package matchingengine

import "time"

// ENUM(BUY, SELL)
type Side int

type Order struct {
	ID                string
	Symbol            string
	Side              Side
	Price             float64
	Quantity          int64
	RemainingQuantity int64
	CanceledQuantity  int64
	CreatedAt         time.Time
	UpdatedAt         time.Time

	IsDummyNode bool

	Next *Order
	Prev *Order
}

func GetOppositeSide(side Side) Side {
	if side == SideBUY {
		return SideSELL
	} else {
		return SideBUY
	}
}
