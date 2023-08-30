package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zenorachi/balance-management/internal/entity"
)

type AccountsRepository struct {
	db *sql.DB
}

func NewAccount(db *sql.DB) *AccountsRepository {
	return &AccountsRepository{db: db}
}

func (a *AccountsRepository) Create(ctx context.Context, account entity.Account) (int, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}

	var (
		id    int
		query = fmt.Sprintf("INSERT INTO %s (user_id) VALUES ($1) RETURNING id",
			collectionAccounts)
	)

	err = tx.QueryRowContext(ctx, query, account.UserID).Scan(&id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

//func (a *AccountsRepository) DepositByID(ctx context.Context, id int, amount float64) error {
//	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
//		Isolation: sql.LevelSerializable,
//		ReadOnly:  false,
//	})
//	if err != nil {
//		return err
//	}
//
//	query := "UPDATE accounts SET balance = balance + $1 WHERE id = $2"
//
//	_, err = tx.ExecContext(ctx, query, amount, id)
//	if err != nil {
//		_ = tx.Rollback()
//		return err
//	}
//
//	return tx.Commit()
//}

func (a *AccountsRepository) DepositByID(ctx context.Context, id int, amount float64) error {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET balance = balance + $1 WHERE id = $2",
		collectionAccounts)

	_, err = tx.ExecContext(ctx, query, amount, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (a *AccountsRepository) Transfer(ctx context.Context, srcAccountId, dstAccountId int, amount float64) error {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}

	var (
		queryWithdraw = fmt.Sprintf("UPDATE %s SET balance = balance - $1 WHERE id = $2",
			collectionAccounts)
		queryDeposit = fmt.Sprintf("UPDATE %s SET balance = balance + $1 WHERE id = $2",
			collectionAccounts)
	)

	_, err = tx.ExecContext(ctx, queryWithdraw, amount, srcAccountId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, queryDeposit, amount, dstAccountId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (a *AccountsRepository) GetByID(ctx context.Context, id int) (entity.Account, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Account{}, err
	}

	var (
		account entity.Account
		query   = fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
			collectionAccounts)
	)

	err = tx.QueryRowContext(ctx, query, id).
		Scan(&account.ID, &account.UserID, &account.Balance, &account.CreatedAt)
	if err != nil {
		_ = tx.Rollback()
		return entity.Account{}, err
	}

	return account, tx.Commit()
}

func (a *AccountsRepository) GetByUserID(ctx context.Context, userId int) (entity.Account, error) {
	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Account{}, err
	}

	var (
		account entity.Account
		query   = fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1",
			collectionAccounts)
	)

	err = tx.QueryRowContext(ctx, query, userId).
		Scan(&account.ID, &account.UserID, &account.Balance, &account.CreatedAt)
	if err != nil {
		_ = tx.Rollback()
		return entity.Account{}, err
	}

	return account, tx.Commit()
}

//func (a *AccountsRepository) WithdrawFromBalanceByID(ctx context.Context, id int, amount float64) error {
//	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{
//		Isolation: sql.LevelSerializable,
//		ReadOnly:  false,
//	})
//	if err != nil {
//		return err
//	}
//
//	account, err := a.GetByID(ctx, id)
//	if err != nil {
//		return err
//	}
//
//	if account.Balance < amount {
//		return entity.ErrNotEnoughMoney
//	}
//
//	query := "UPDATE accounts SET balance = balance - $1 WHERE id = $2"
//
//	_, err = tx.ExecContext(ctx, query, amount, id)
//	if err != nil {
//		_ = tx.Rollback()
//		return err
//	}
//
//	return tx.Commit()
//}
