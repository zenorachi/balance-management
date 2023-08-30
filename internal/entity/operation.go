package entity

import "time"

type Operation struct {
	ID            int
	AccountID     int
	OrderID       int
	Amount        float64
	OperationType string
	OrderDate     time.Time
	Description   string
}

func NewOperation(accountId, orderId int, amount float64, opType string, orderDate time.Time) Operation {
	return Operation{
		AccountID:     accountId,
		OrderID:       orderId,
		Amount:        amount,
		OperationType: opType,
		OrderDate:     orderDate,
	}
}
