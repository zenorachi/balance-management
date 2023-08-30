package entity

import (
	"errors"
)

var (
	ErrInvalidInput           = errors.New("invalid input")
	ErrUserAlreadyExists      = errors.New("user with such login already exist")
	ErrUserDoesNotExist       = errors.New("user does not exist")
	ErrAccountAlreadyExists   = errors.New("account already exist")
	ErrAccountDoesNotExist    = errors.New("account does not exist")
	ErrAmountIsNegative       = errors.New("amount is negative or equals to zero")
	ErrProductAlreadyExists   = errors.New("product already exists")
	ErrProductDoesNotExist    = errors.New("product/products does not exists")
	ErrPriceIsNegative        = errors.New("price is negative or equals to zero")
	ErrOrderDoesNotExist      = errors.New("order does not exists")
	ErrOrderCannotBeCancelled = errors.New("order cannot be cancelled")
	ErrEmptyOrder             = errors.New("order is empty")
	ErrNotEnoughMoney         = errors.New("not enough money")
	ErrReserveDoesNotExist    = errors.New("reserve does not exist")
	ErrOrderCannotBeProcessed = errors.New("order cannot be processed")
)
