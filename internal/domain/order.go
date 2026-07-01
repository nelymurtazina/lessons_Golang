package domain

import "time"

type Order struct {
	ID        string
	UserID    string
	MarketID  string
	OrderType string
	Price     float64
	Quantity  int64
	Status    string
	CreatedAt time.Time
}