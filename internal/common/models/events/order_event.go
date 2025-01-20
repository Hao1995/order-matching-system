//go:generate go-enum --marshal
package events

import "time"

// ENUM(Buy, Sell)
type OrderType int

// ENUM(CreateOrder, CancelOrder)
type EventType string

type Event struct {
	EventType EventType
	Data      interface{}
}

type OrderEvent struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol,omitempty"`
	Type      OrderType `json:"type,omitempty"`
	Price     float64   `json:"price,omitempty"`
	Quantity  int64     `json:"quantity,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
