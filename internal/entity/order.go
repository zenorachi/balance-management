package entity

import "time"

type Order struct {
	ID        int
	AccountID int
	Products  []uint8
	Amount    float64
	CreatedAt time.Time
	Status    string
}
