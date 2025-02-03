//go:generate go-enum --marshal
package events

import "time"

// ENUM(CreateOrder, CancelOrder, Matching)
type EventType string

type Event struct {
	EventType EventType   `json:"event_type"`
	Data      interface{} `json:"data"`
}

type OrderEvent struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Type      string    `json:"type"`
	Price     float64   `json:"price"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
