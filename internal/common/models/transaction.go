package models

import "time"

type Transaction struct {
	ID          string
	Symbol      string
	BuyOrderID  string
	SellOrderID string
	Price       float64
	Quantity    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
