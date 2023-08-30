package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zenorachi/balance-management/internal/entity"
)

type OperationRepository struct {
	db *sql.DB
}

func NewOperation(db *sql.DB) *OperationRepository {
	return &OperationRepository{db: db}
}

func (o *OperationRepository) Create(ctx context.Context, operation entity.Operation) (int, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}

	var (
		id    int
		query = fmt.Sprintf("INSERT INTO %s (account_id, order_id, amount, type, order_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			collectionOperations)
	)

	err = tx.
		QueryRowContext(
			ctx, query, operation.AccountID, operation.OrderID, operation.Amount, operation.OperationType, operation.OrderDate).
		Scan(&id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (o *OperationRepository) GetByID(ctx context.Context, id int) (entity.Operation, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Operation{}, err
	}

	var (
		operation entity.Operation
		query     = fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
			collectionOperations)
	)

	err = tx.QueryRowContext(ctx, query, id).
		Scan(&operation.ID, &operation.AccountID, &operation.OrderID, &operation.Amount,
			&operation.OperationType, &operation.OrderDate, &operation.Description)
	if err != nil {
		_ = tx.Rollback()
		return entity.Operation{}, err
	}

	return operation, tx.Commit()
}

func (o *OperationRepository) GetAll(ctx context.Context) ([]entity.Operation, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}

	var (
		operations []entity.Operation
		query      = fmt.Sprintf("SELECT * FROM %s",
			collectionOperations)
	)

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var operation entity.Operation

		err = rows.Scan(&operation.ID, &operation.AccountID, &operation.OrderID, &operation.Amount,
			&operation.OperationType, &operation.OrderDate, &operation.Description)
		if err != nil {
			return nil, err
		}

		operations = append(operations, operation)
	}

	return operations, tx.Commit()
}
