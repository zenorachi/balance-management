package entity

import "time"

type Account struct {
	ID        int       `json:"id,omitempty"`
	UserID    int       `json:"user_id"`
	Balance   float64   `json:"amount,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(userId int) Account {
	return Account{UserID: userId}
}
