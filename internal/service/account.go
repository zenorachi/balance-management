package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zenorachi/balance-management/internal/entity"
	"github.com/zenorachi/balance-management/internal/repository"
)

type AccountService struct {
	repo repository.Account
}

func NewAccounts(repo repository.Account) *AccountService {
	return &AccountService{repo: repo}
}

func (a *AccountService) Create(ctx context.Context, userId int) (int, error) {
	if a.isAccountExists(ctx, userId) {
		return 0, entity.ErrAccountAlreadyExists
	}

	account := entity.NewAccount(userId)
	return a.repo.Create(ctx, account)
}

func (a *AccountService) DepositByID(ctx context.Context, id int, amount float64) error {
	_, err := a.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrAccountDoesNotExist
		}
		return err
	}

	if amount <= 0 {
		return entity.ErrAmountIsNegative
	}

	return a.repo.DepositByID(ctx, id, amount)
}

func (a *AccountService) Transfer(ctx context.Context, srcAccountId, dstAccountId int, amount float64) error {
	account, _ := a.repo.GetByID(ctx, srcAccountId)
	if account.Balance < amount {
		return entity.ErrNotEnoughMoney
	}

	err := a.repo.Transfer(ctx, srcAccountId, dstAccountId, amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.ErrAccountDoesNotExist
		}
		return err
	}

	return nil
}

func (a *AccountService) GetByID(ctx context.Context, id int) (entity.Account, error) {
	if !a.isAccountExists(ctx, id) {
		return entity.Account{}, entity.ErrAccountDoesNotExist
	}

	return a.repo.GetByID(ctx, id)
}

func (a *AccountService) isAccountExists(ctx context.Context, userId int) bool {
	_, err := a.repo.GetByUserID(ctx, userId)
	return !errors.Is(err, sql.ErrNoRows)
}
