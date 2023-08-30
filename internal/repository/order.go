package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zenorachi/balance-management/internal/entity"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrder(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) Create(ctx context.Context, order entity.Order) (int, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}

	var (
		id          int
		queryCreate = fmt.Sprintf("INSERT INTO %s (account_id, amount) VALUES ($1, $2) RETURNING id",
			collectionOrders)
		querySetProducts = fmt.Sprintf("UPDATE %s SET products = array_append(products, $1) WHERE id = $2;",
			collectionOrders)
	)

	err = tx.QueryRowContext(ctx, queryCreate, order.AccountID, order.Amount).Scan(&id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	for _, product := range order.Products {
		_, err = tx.ExecContext(ctx, querySetProducts, product, id)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	return id, tx.Commit()
}

func (o *OrderRepository) GetByID(ctx context.Context, id int) (entity.Order, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Order{}, err
	}

	var (
		order entity.Order
		query = fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
			collectionOrders)
	)

	err = tx.QueryRowContext(ctx, query, id).
		Scan(&order.ID, &order.AccountID, &order.Products, &order.Amount, &order.CreatedAt, &order.Status)

	if err != nil {
		_ = tx.Rollback()
		return entity.Order{}, err
	}

	return order, tx.Commit()
}

func (o *OrderRepository) GetAllByAccountID(ctx context.Context, accountId int) ([]entity.Order, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}

	var (
		orders []entity.Order
		query  = fmt.Sprintf("SELECT * FROM %s WHERE account_id = $1",
			collectionOrders)
	)

	rows, err := tx.QueryContext(ctx, query, accountId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var order entity.Order
		if err = rows.
			Scan(&order.ID, &order.AccountID, &order.Products, &order.Amount, &order.CreatedAt, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, tx.Commit()
}

func (o *OrderRepository) SetStatusByID(ctx context.Context, id int, status string) (entity.Order, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return entity.Order{}, err
	}

	var (
		order entity.Order
		query = fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2 RETURNING id, account_id, products, created_at, status",
			collectionOrders)
	)

	err = tx.QueryRowContext(ctx, query, status, id).
		Scan(&order.ID, &order.AccountID, &order.Products, &order.CreatedAt, &order.Status)
	if err != nil {
		_ = tx.Rollback()
		return entity.Order{}, err
	}

	return order, tx.Commit()
}
