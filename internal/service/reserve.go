package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zenorachi/balance-management/internal/entity"
	"github.com/zenorachi/balance-management/internal/repository"
)

type ReserveService struct {
	reserveRepo   repository.Reserve
	accountRepo   repository.Account
	productRepo   repository.Product
	orderRepo     repository.Order
	operationRepo repository.Operation
}

func NewReserve(reserveRepo repository.Reserve, accountRepo repository.Account, productRepo repository.Product,
	orderRepo repository.Order, operationRepo repository.Operation) *ReserveService {
	return &ReserveService{
		reserveRepo:   reserveRepo,
		accountRepo:   accountRepo,
		productRepo:   productRepo,
		orderRepo:     orderRepo,
		operationRepo: operationRepo,
	}
}

func (r *ReserveService) Create(ctx context.Context, reserve entity.Reserve) (int, error) {
	if !r.isOrderExists(ctx, reserve.OrderID) {
		return 0, entity.ErrOrderDoesNotExist
	}

	if !r.isOrderCanBeProcessed(ctx, reserve.OrderID) {
		return 0, entity.ErrOrderCannotBeProcessed
	}

	return r.reserveRepo.Create(ctx, reserve)
}

func (r *ReserveService) ConfirmRevenueByID(ctx context.Context, id int) (int, error) {
	if !r.isReserveExists(ctx, id) {
		return 0, entity.ErrReserveDoesNotExist
	}

	operationId, err := r.reserveRepo.ConfirmRevenueByID(ctx, id)
	if err != nil {
		return 0, err
	}

	return operationId, nil
}

func (r *ReserveService) ConfirmRefundByID(ctx context.Context, id int) (int, error) {
	if !r.isReserveExists(ctx, id) {
		return 0, entity.ErrReserveDoesNotExist
	}

	operationId, err := r.reserveRepo.ConfirmRefundByID(ctx, id)
	if err != nil {
		return 0, err
	}

	return operationId, nil
}

func (r *ReserveService) isReserveExists(ctx context.Context, id int) bool {
	_, err := r.reserveRepo.GetByID(ctx, id)
	return !errors.Is(err, sql.ErrNoRows)
}

func (r *ReserveService) isOrderExists(ctx context.Context, id int) bool {
	_, err := r.orderRepo.GetByID(ctx, id)
	return !errors.Is(err, sql.ErrNoRows)
}

func (r *ReserveService) isOrderCanBeProcessed(ctx context.Context, id int) bool {
	order, _ := r.orderRepo.GetByID(ctx, id)
	return order.Status == entity.StatusAccepted
}

func (r *ReserveService) isOrderStatusProcessed(ctx context.Context, id int) bool {
	order, _ := r.orderRepo.GetByID(ctx, id)
	return order.Status == entity.StatusProcessing
}
