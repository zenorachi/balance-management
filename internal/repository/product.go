package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zenorachi/balance-management/internal/entity"
)

type ProductsRepository struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) *ProductsRepository {
	return &ProductsRepository{db: db}
}

func (p *ProductsRepository) Create(ctx context.Context, product entity.Product) (int, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, err
	}

	var (
		id    int
		query = fmt.Sprintf("INSERT INTO %s (name, price) VALUES ($1, $2) RETURNING id",
			collectionProducts)
	)

	err = tx.QueryRowContext(ctx, query, product.Name, product.Price).Scan(&id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (p *ProductsRepository) GetByID(ctx context.Context, id int) (entity.Product, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Product{}, err
	}

	var (
		product entity.Product
		query   = fmt.Sprintf("SELECT * FROM %s WHERE id = $1",
			collectionProducts)
	)

	err = tx.QueryRowContext(ctx, query, id).
		Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		_ = tx.Rollback()
		return entity.Product{}, err
	}

	return product, tx.Commit()
}

func (p *ProductsRepository) GetByName(ctx context.Context, name string) (entity.Product, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return entity.Product{}, err
	}

	var (
		product entity.Product
		query   = fmt.Sprintf("SELECT * FROM %s WHERE name = $1",
			collectionProducts)
	)

	err = tx.QueryRowContext(ctx, query, name).
		Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		_ = tx.Rollback()
		return entity.Product{}, err
	}

	return product, tx.Commit()
}

func (p *ProductsRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}

	var (
		products []entity.Product
		query    = fmt.Sprintf("SELECT * FROM %s",
			collectionProducts)
	)

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var product entity.Product
		if err = rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, tx.Commit()
}
