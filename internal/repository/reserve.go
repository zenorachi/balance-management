package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zenorachi/balance-management/internal/entity"
)

type ReserveRepository struct {
	db *sql.DB
}

func NewReserve(db *sql.DB) *ReserveRepository {
	return &ReserveRepository{db: db}
}

func (r *ReserveRepository) Create(ctx context.Context, reserve entity.Reserve) (int, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}

	var (
		id                  int
		amount              float64
		queryGetOrderAmount = fmt.Sprintf("SELECT amount FROM %s WHERE id = $1",
			collectionOrders)
		queryCreateReserve = fmt.Sprintf("INSERT INTO %s (order_id, amount) VALUES ($1, $2) RETURNING id",
			collectionReserves)
		queryUpdateStatus = fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2 RETURNING id, account_id, product_id, created_at, status",
			collectionOrders)
		//queryGetProduct = fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		//	collectionProducts)
		queryWithdraw = fmt.Sprintf("UPDATE %s SET balance = balance - $1 WHERE id = $2",
			collectionAccounts)
		order entity.Order
		//product entity.Product
	)

	err = tx.QueryRowContext(ctx, queryGetOrderAmount, reserve.OrderID).Scan(&amount)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	err = tx.QueryRowContext(ctx, queryCreateReserve, reserve.OrderID, amount).Scan(&id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	err = tx.QueryRowContext(ctx, queryUpdateStatus, entity.StatusProcessing, reserve.OrderID).
		Scan(&order.ID, &order.AccountID, &order.Products, &order.CreatedAt, &order.Status)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	//err = tx.QueryRowContext(ctx, queryGetProduct, order.ProductID).Scan(&product.ID, &product.Name, &product.Price)
	//if err != nil {
	//	_ = tx.Rollback()
	//	return 0, err
	//}

	_, err = tx.ExecContext(ctx, queryWithdraw, amount, order.AccountID)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *ReserveRepository) GetByID(ctx context.Context, id int) (entity.Reserve, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Reserve{}, err
	}

	var (
		query   = fmt.Sprintf("SELECT * FROM %s WHERE id = $1", collectionReserves)
		reserve entity.Reserve
	)

	err = tx.QueryRowContext(ctx, query, id).Scan(&reserve.ID, &reserve.OrderID, &reserve.Amount, &reserve.CreatedAt)
	if err != nil {
		_ = tx.Rollback()
		return entity.Reserve{}, err
	}

	return reserve, tx.Commit()
}

func (r *ReserveRepository) ConfirmRevenueByID(ctx context.Context, id int) (int, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}

	var (
		operationId        int
		queryDeleteReserve = fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING id, order_id, created_at",
			collectionReserves)
		queryUpdateStatus = fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2 RETURNING id, account_id, products, amount, created_at, status",
			collectionOrders)
		queryCreateOperation = fmt.Sprintf("INSERT INTO %s (account_id, order_id, amount, type, order_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			collectionOperations)
		//queryGetProduct = fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		//	collectionProducts)
		reserve entity.Reserve
		order   entity.Order
		//product   entity.Product
		operation entity.Operation
	)

	err = tx.QueryRowContext(ctx, queryDeleteReserve, id).Scan(&reserve.ID, &reserve.OrderID, &reserve.CreatedAt)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	err = tx.QueryRowContext(ctx, queryUpdateStatus, entity.StatusConfirmed, reserve.OrderID).
		Scan(&order.ID, &order.AccountID, &order.Products, &order.Amount, &order.CreatedAt, &order.Status)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	//err = tx.QueryRowContext(ctx, queryGetProduct, order.ID).Scan(&product.ID, &product.Name, &product.Price)
	//if err != nil {
	//	_ = tx.Rollback()
	//	return 0, err
	//}

	operation = entity.NewOperation(order.AccountID, order.ID, order.Amount, entity.TypeRevenue, order.CreatedAt)
	err = tx.
		QueryRowContext(
			ctx, queryCreateOperation, operation.AccountID, operation.OrderID, operation.Amount, operation.OperationType, operation.OrderDate).
		Scan(&operationId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return operationId, tx.Commit()
}

func (r *ReserveRepository) ConfirmRefundByID(ctx context.Context, id int) (int, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}

	var (
		queryDeleteReserve = fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING id, order_id, created_at",
			collectionReserves)
		queryUpdateStatus = fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2 RETURNING id, account_id, products, amount, created_at, status",
			collectionOrders)
		queryCreateOperation = fmt.Sprintf("INSERT INTO %s (account_id, order_id, amount, type, order_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			collectionOperations)
		//queryGetProduct = fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
		//	collectionProducts)
		queryDeposit = fmt.Sprintf("UPDATE %s SET balance = balance + $1 WHERE id = $2",
			collectionAccounts)
		reserve entity.Reserve
		order   entity.Order
		//product     entity.Product
		operation   entity.Operation
		operationId int
	)

	err = tx.QueryRowContext(ctx, queryDeleteReserve, id).Scan(&reserve.ID, &reserve.OrderID, &reserve.CreatedAt)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	err = tx.QueryRowContext(ctx, queryUpdateStatus, entity.StatusCancelled, reserve.OrderID).
		Scan(&order.ID, &order.AccountID, &order.Products, &order.Amount, &order.CreatedAt, &order.Status)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	//err = tx.QueryRowContext(ctx, queryGetProduct, order.ProductID).Scan(&product.ID, &product.Name, &product.Price)
	//if err != nil {
	//	_ = tx.Rollback()
	//	return 0, err
	//}

	_, err = tx.ExecContext(ctx, queryDeposit, order.Amount, order.AccountID)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	operation = entity.NewOperation(order.AccountID, order.ID, order.Amount, entity.TypeRefund, order.CreatedAt)
	err = tx.
		QueryRowContext(
			ctx, queryCreateOperation, operation.AccountID, operation.OrderID, operation.Amount, operation.OperationType, operation.OrderDate).
		Scan(&operationId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return operationId, tx.Commit()
}
