package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zenorachi/balance-management/internal/entity"
	"github.com/zenorachi/balance-management/internal/repository"
)

type OrderService struct {
	orderRepo   repository.Order
	accountRepo repository.Account
	productRepo repository.Product
}

func NewOrder(orderRepo repository.Order, accountRepo repository.Account, productRepo repository.Product) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		accountRepo: accountRepo,
		productRepo: productRepo,
	}
}

func (o *OrderService) Create(ctx context.Context, order entity.Order) (int, error) {
	if !o.isAccountExists(ctx, order.AccountID) {
		return 0, entity.ErrAccountDoesNotExist
	}

	if len(order.Products) == 0 {
		return 0, entity.ErrEmptyOrder
	}

	for _, product := range order.Products {
		if !o.isProductExists(ctx, int(product)) {
			return 0, entity.ErrProductDoesNotExist
		}
	}

	var (
		amount  float64
		account entity.Account
		err     error
	)

	account, err = o.accountRepo.GetByID(ctx, order.AccountID)
	if err != nil {
		return 0, err
	}

	amount, err = o.getOrderAmount(ctx, order)
	if err != nil {
		return 0, err
	}

	if account.Balance < amount {
		return 0, entity.ErrNotEnoughMoney
	}

	order.Amount = amount
	return o.orderRepo.Create(ctx, order)
}

func (o *OrderService) CancelByID(ctx context.Context, id int) error {
	if !o.isOrderExists(ctx, id) {
		return entity.ErrOrderDoesNotExist
	}

	if !o.isOrderCanBeCancelled(ctx, id) {
		return entity.ErrOrderCannotBeCancelled
	}

	_, err := o.orderRepo.SetStatusByID(ctx, id, entity.StatusCancelled)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderService) GetAllByAccountID(ctx context.Context, accountId int) ([]entity.Order, error) {
	if !o.isAccountExists(ctx, accountId) {
		return nil, entity.ErrAccountDoesNotExist
	}

	return o.orderRepo.GetAllByAccountID(ctx, accountId)
}

func (o *OrderService) isOrderExists(ctx context.Context, id int) bool {
	_, err := o.orderRepo.GetByID(ctx, id)
	return !errors.Is(err, sql.ErrNoRows)
}

func (o *OrderService) isOrderCanBeCancelled(ctx context.Context, id int) bool {
	order, _ := o.orderRepo.GetByID(ctx, id)
	return order.Status == entity.StatusAccepted
}

func (o *OrderService) isAccountExists(ctx context.Context, id int) bool {
	_, err := o.accountRepo.GetByID(ctx, id)
	return !errors.Is(err, sql.ErrNoRows)
}

func (o *OrderService) isProductExists(ctx context.Context, id int) bool {
	_, err := o.productRepo.GetByID(ctx, id)
	return !errors.Is(err, sql.ErrNoRows)
}

func (o *OrderService) getOrderAmount(ctx context.Context, order entity.Order) (float64, error) {
	var (
		amount  float64
		product entity.Product
		err     error
	)

	for _, productID := range order.Products {
		product, err = o.productRepo.GetByID(ctx, int(productID))
		if err != nil {
			return 0, err
		}

		amount += product.Price
	}

	return amount, nil
}
