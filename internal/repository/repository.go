package repository

import (
	"context"
	"database/sql"
	"github.com/zenorachi/balance-management/internal/entity"
)

type (
	User interface {
		Create(ctx context.Context, user entity.User) (int, error)
		GetByLogin(ctx context.Context, login string) (entity.User, error)
		GetByCredentials(ctx context.Context, login, password string) (entity.User, error)
		GetByRefreshToken(ctx context.Context, refreshToken string) (entity.User, error)
		SetSession(ctx context.Context, userId int, session entity.Session) error
	}

	Account interface {
		Create(ctx context.Context, account entity.Account) (int, error)
		// DepositByUserID DepositByID(ctx context.Context, id int, amount float64) error
		DepositByID(ctx context.Context, id int, amount float64) error
		Transfer(ctx context.Context, srcAccountId, dstAccountId int, amount float64) error
		GetByID(ctx context.Context, id int) (entity.Account, error)
		GetByUserID(ctx context.Context, userId int) (entity.Account, error)
		//WithdrawFromBalanceByID(ctx context.Context, id int, amount float64) error
	}

	Product interface {
		Create(ctx context.Context, product entity.Product) (int, error)
		GetByID(ctx context.Context, id int) (entity.Product, error)
		GetByName(ctx context.Context, name string) (entity.Product, error)
		GetAll(ctx context.Context) ([]entity.Product, error)
	}

	Order interface {
		Create(ctx context.Context, order entity.Order) (int, error)
		GetByID(ctx context.Context, id int) (entity.Order, error)
		GetAllByAccountID(ctx context.Context, accountId int) ([]entity.Order, error)
		SetStatusByID(ctx context.Context, id int, status string) (entity.Order, error)
	}

	Reserve interface {
		Create(ctx context.Context, reserve entity.Reserve) (int, error)
		GetByID(ctx context.Context, id int) (entity.Reserve, error)
		ConfirmRevenueByID(ctx context.Context, id int) (int, error)
		ConfirmRefundByID(ctx context.Context, id int) (int, error)
	}

	Operation interface {
		Create(ctx context.Context, operation entity.Operation) (int, error)
		GetByID(ctx context.Context, id int) (entity.Operation, error)
		GetAll(ctx context.Context) ([]entity.Operation, error)
	}
)

type Repositories struct {
	User
	Account
	Product
	Order
	Reserve
	Operation
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		User:      NewUser(db),
		Account:   NewAccount(db),
		Product:   NewProduct(db),
		Order:     NewOrder(db),
		Reserve:   NewReserve(db),
		Operation: NewOperation(db),
	}
}
