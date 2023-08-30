package service

import (
	"context"
	"github.com/zenorachi/balance-management/internal/entity"
	"time"

	"github.com/zenorachi/balance-management/pkg/auth"

	"github.com/zenorachi/balance-management/internal/repository"
	"github.com/zenorachi/balance-management/pkg/hash"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type (
	User interface {
		SignUp(ctx context.Context, login, email, password string) (int, error)
		SignIn(ctx context.Context, login, password string) (Tokens, error)
		RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	}

	Account interface {
		Create(ctx context.Context, userId int) (int, error)
		DepositByID(ctx context.Context, account int, amount float64) error
		Transfer(ctx context.Context, srcAccountId, dstAccountId int, amount float64) error
		GetByID(ctx context.Context, id int) (entity.Account, error)
	}

	Product interface {
		Create(ctx context.Context, product entity.Product) (int, error)
		GetByID(ctx context.Context, id int) (entity.Product, error)
		GetAll(ctx context.Context) ([]entity.Product, error)
	}

	Order interface {
		Create(ctx context.Context, order entity.Order) (int, error)
		CancelByID(ctx context.Context, id int) error
		GetAllByAccountID(ctx context.Context, accountId int) ([]entity.Order, error)
	}

	Reserve interface {
		Create(ctx context.Context, reserve entity.Reserve) (int, error)
		ConfirmRevenueByID(ctx context.Context, id int) (int, error)
		ConfirmRefundByID(ctx context.Context, id int) (int, error)
	}

	Operation interface {
		GetReportForUser(ctx context.Context, accountId int) ([]entity.Operation, error)
		GetReportForAccounting(ctx context.Context) ([]entity.Operation, error)
	}
)

type Services struct {
	User
	Account
	Product
	Order
	Reserve
	Operation
}

type Deps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func New(deps Deps) *Services {
	userService := NewUsers(deps.Repos.User, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)
	accountService := NewAccounts(deps.Repos.Account)
	productService := NewProduct(deps.Repos.Product)
	orderService := NewOrder(deps.Repos.Order, deps.Repos.Account, deps.Repos.Product)
	reserveService := NewReserve(deps.Repos.Reserve, deps.Repos.Account, deps.Repos.Product, deps.Repos.Order, deps.Repos.Operation)
	operationService := NewOperation(deps.Repos.Operation)

	return &Services{
		User:      userService,
		Account:   accountService,
		Product:   productService,
		Order:     orderService,
		Reserve:   reserveService,
		Operation: operationService,
	}
}
