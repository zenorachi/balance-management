package entity

import "time"

type Reserve struct {
	ID        int
	OrderID   int
	Amount    float64
	CreatedAt time.Time
}
