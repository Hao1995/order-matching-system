//go:generate go-enum --marshal
package events

import "time"

// ENUM(Matching, Cancel)
type MatchingEventType string

type MatchingEvent struct {
	EventType MatchingEventType
	Data      interface{}
}

type MatchingOrder struct {
	ID                string
	Symbol            string
	Side              Side
	Price             float64
	Quantity          int64
	RemainingQuantity int64
	CanceledQuantity  int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type MatchingTransaction struct {
	ID          string
	Symbol      string
	BuyOrderID  string
	SellOrderID string
	Price       float64
	Quantity    int64
	CreatedAt   time.Time
}

type MatchingTick struct {
	Price    float64
	Quantity int64
}

type MatchingData struct {
	Order        interface{}
	Transactions []MatchingTransaction
	TopBuy       []MatchingTick
	TopSell      []MatchingTick
}
