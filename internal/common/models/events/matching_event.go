//go:generate go-enum --marshal
package events

import "time"

// ENUM(Create, Cancel)
type MatchingEventType string

type MatchingEvent struct {
	Type         MatchingEventType  `json:"type"`
	Order        OrderEvent         `json:"order"`
	Transactions []TransactionEvent `json:"transactions,omitempty"`
	BuyTicks     []TickEvent        `json:"buy_ticks"`
	SellTicks    []TickEvent        `json:"sell_ticks"`
}

type TransactionEvent struct {
	ID          string    `jsno:"id"`
	Symbol      string    `jsno:"symbol"`
	BuyOrderID  string    `jsno:"buy_order_id"`
	SellOrderID string    `jsno:"sell_order_id"`
	Price       float64   `jsno:"price"`
	Quantity    int64     `jsno:"quantity"`
	CreatedAt   time.Time `jsno:"created_at"`
}

type TickEvent struct {
	Price    float64 `json:"price"`
	Quantity int64   `json:"quantity"`
}
