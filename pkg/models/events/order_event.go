package events

import "time"

type Side int8

const (
	SideBuy Side = iota
	SideSell
)

type OrderEvent struct {
	EventType string
	Data      interface{}
}

type OrderCreateEvent struct {
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

type OrderCancelEvent struct {
	ID string
}
