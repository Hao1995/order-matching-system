//go:generate go-enum --marshal
package events

import "time"

// ENUM(Buy, Sell)
type OrderType int

// ENUM(CREATE, CANCEL)
type OrderEventType string

type OrderEvent struct {
	EventType OrderEventType
	Data      interface{}
}

type OrderCreateEvent struct {
	ID                string
	Symbol            string
	OrderType         OrderType
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
